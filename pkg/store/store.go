package store

import (
	"context"
	"database/sql"
	"github.com/aibotsoft/micro/cache"
	"github.com/aibotsoft/micro/config"
	"github.com/dgraph-io/ristretto"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
	"time"
)

type Store struct {
	cfg   *config.Config
	log   *zap.SugaredLogger
	db    *sqlx.DB
	Cache *ristretto.Cache
}

func New(cfg *config.Config, log *zap.SugaredLogger, db *sqlx.DB) *Store {
	return &Store{log: log, db: db, Cache: cache.NewCache(cfg)}
}
func (s *Store) Close() {
	err := s.db.Close()
	if err != nil {
		s.log.Error(err)
	}
	s.Cache.Close()
}

type Token struct {
	Session     string
	CreatedAt   time.Time
	LastCheckAt time.Time
}
type Account struct {
	Id       int
	Username string
	Password string
}

func (s *Store) LoadToken(ctx context.Context) (token Token, err error) {
	err = s.db.GetContext(ctx, &token, "select top 1 Session, CreatedAt, LastCheckAt from dbo.Auth order by CreatedAt desc")
	return
}
func (s *Store) LoadAccount(ctx context.Context) (account Account, err error) {
	err = s.db.GetContext(ctx, &account, "select top 1 Id, Username, Password from dbo.Account")
	return
}
func (s *Store) SaveToken(ctx context.Context, token Token) (err error) {
	_, err = s.db.ExecContext(ctx, "insert into dbo.Auth (Session, CreatedAt, LastCheckAt) values (@p1, @p2, @p3)",
		token.Session, token.CreatedAt, token.LastCheckAt)
	return
}

func (s *Store) UpdateToken(ctx context.Context, token Token) (err error) {
	_, err = s.db.ExecContext(ctx, "update dbo.Auth set LastCheckAt = @p1 where Session = @p2",
		token.LastCheckAt, token.Session)
	return
}

func (s *Store) SaveSport(id int64, name string) {
	_, b := s.Cache.Get(id)
	if b {
		return
	}
	_, err := s.db.Exec("insert into dbo.Sport (Id, Name) select @p1, @p2 where not exists(select 1 from dbo.Sport where Name = @p2)", id, name)
	if err != nil {
		s.log.Error(err)
	} else {
		s.Cache.SetWithTTL(id, true, 1, time.Hour*12)
	}
}

func (s *Store) SaveLeague(id int64, name string, country string, sportId int64) {
	_, b := s.Cache.Get(id)
	if b {
		return
	}
	_, err := s.db.Exec("dbo.uspSaveLeague", id, name, country, sportId)
	if err != nil {
		s.log.Error(err)
	} else {
		s.Cache.SetWithTTL(id, true, 1, time.Hour*12)
	}
}

func (s *Store) SaveTeam(id int64, name string) {
	_, b := s.Cache.Get(id)
	if b {
		return
	}
	_, err := s.db.Exec("dbo.uspSaveTeam",
		sql.Named("Id", id),
		sql.Named("Name", name),
	)
	if err != nil {
		s.log.Error(err)
	} else {
		s.Cache.SetWithTTL(id, true, 1, time.Hour*12)
	}
}

const saveEventQ = `
insert into dbo.Event (Id, HomeId, AwayId, LeagueId, Starts) 
select @p1, @p2, @p3, @p4 , @p5
where not exists(select 1 from dbo.Event where Id = @p1)
`

func (s *Store) SaveEvent(id string, homeId int64, awayId int64, leagueId int64, starts string) {
	_, b := s.Cache.Get(id)
	if b {
		return
	}
	_, err := s.db.Exec(saveEventQ, id, homeId, awayId, leagueId, starts)
	if err != nil {
		s.log.Error(err)
	} else {
		s.Cache.SetWithTTL(id, true, 1, time.Hour*12)
	}
}

const saveEventPeriodQ = `
insert into dbo.EventPeriod (EventId, PeriodCode) 
select @p1, @p2
where not exists(select 1 from dbo.EventPeriod where EventId = @p1 and PeriodCode = @p2)
`

func (s *Store) SaveEventPeriod(eventId string, periodCode string, isActive bool) {
	_, err := s.db.Exec("dbo.uspSaveEventPeriod",
		sql.Named("EventId", eventId),
		sql.Named("PeriodCode", periodCode),
		sql.Named("IsActive", isActive),
	)
	if err != nil {
		s.log.Error(err)
	}
}

const getEventQ = `
select e.Id,
       e.LeagueId,
       PeriodCode
from dbo.Event e
         join dbo.League l on e.LeagueId = l.Id
         join dbo.EventPeriod ep on ep.EventId = e.Id
where l.SportId = 1 
  and ep.IsActive = 1
  and e.Starts < sysdatetimeoffset()
`

type Event struct {
	Id         string
	LeagueId   int64
	PeriodCode string
}

func (s *Store) GetLiveEvents() (events []Event, err error) {
	err = s.db.Select(&events, getEventQ)
	return
}
