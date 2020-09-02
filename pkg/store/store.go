package store

import (
	"context"
	"database/sql"
	api "github.com/aibotsoft/gen/blackapi"
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

type Event struct {
	Id         string
	LeagueId   int64
	PeriodCode string
}

const getEventQ = `
select top 200 e.Id,
               e.LeagueId,
               PeriodCode
from dbo.Event e
         join dbo.League l on e.LeagueId = l.Id
         join dbo.EventPeriod ep on ep.EventId = e.Id
where l.SportId = 1
  and ep.IsActive = 1
order by e.Starts
`

func (s *Store) GetLiveEvents() (events []Event, err error) {
	err = s.db.Select(&events, getEventQ)
	return
}

func (s *Store) DeactivateHandicap(eventPeriodId int64, handicapCode int64) {
	_, _ = s.db.Exec("update dbo.Handicap set IsActive = 0 where EventPeriodId = @p1 and HandicapCode = @p2", eventPeriodId, handicapCode)
}
func (s *Store) SaveHandicap(eventPeriodId int64, handicapCode int64, away float64, home float64, margin float64, isActive bool) {
	_, err := s.db.Exec("dbo.uspSaveHandicap",
		sql.Named("EventPeriodId", eventPeriodId),
		sql.Named("HandicapCode", handicapCode),
		sql.Named("Away", away),
		sql.Named("Home", home),
		sql.Named("Margin", margin),
		sql.Named("IsActive", isActive),
	)
	if err != nil {
		s.log.Error(err)
	}
}
func (s *Store) SaveTotal(eventPeriodId int64, handicapCode int64, ove float64, under float64, margin float64, isActive bool) {
	_, err := s.db.Exec("dbo.uspSaveTotal",
		sql.Named("EventPeriodId", eventPeriodId),
		sql.Named("HandicapCode", handicapCode),
		sql.Named("Ove", ove),
		sql.Named("Under", under),
		sql.Named("Margin", margin),
		sql.Named("IsActive", isActive),
	)
	if err != nil {
		s.log.Error(err)
	}
}
func (s *Store) DeactivateTotal(eventPeriodId int64, handicapCode int64) {
	_, _ = s.db.Exec("update dbo.Total set IsActive = 0 where EventPeriodId = @p1 and HandicapCode = @p2", eventPeriodId, handicapCode)
}
func (s *Store) SaveDoubleChance(eventPeriodId int64, awayDraw float64, homeAway float64, homeDraw float64, margin float64, isActive bool) {
	_, err := s.db.Exec("dbo.uspSaveDoubleChance",
		sql.Named("EventPeriodId", eventPeriodId),
		sql.Named("AwayDraw", awayDraw),
		sql.Named("HomeAway", homeAway),
		sql.Named("HomeDraw", homeDraw),
		sql.Named("Margin", margin),
		sql.Named("IsActive", isActive),
	)
	if err != nil {
		s.log.Error(err)
	}
}
func (s *Store) SaveWinDrawWin(eventPeriodId int64, away float64, home float64, draw float64, margin float64, isActive bool) {
	_, err := s.db.Exec("dbo.uspSaveWinDrawWin",
		sql.Named("EventPeriodId", eventPeriodId),
		sql.Named("Away", away),
		sql.Named("Home", home),
		sql.Named("Draw", draw),
		sql.Named("Margin", margin),
		sql.Named("IsActive", isActive),
	)
	if err != nil {
		s.log.Error(err)
	}
}

type Trio struct {
	EventPeriodId int64
	BetTypeId int64
	APrice float64
	BPrice float64
	CPrice float64
	IsActive bool
}
func (s *Store) SaveTrio(t Trio) {
	_, err := s.db.Exec("dbo.uspSaveTrio",
		sql.Named("EventPeriodId", t.EventPeriodId),
		sql.Named("BetTypeId", t.BetTypeId),
		sql.Named("APrice", t.APrice),
		sql.Named("BPrice", t.BPrice),
		sql.Named("CPrice", t.CPrice),
		sql.Named("IsActive", t.IsActive),
	)
	if err != nil {
		s.log.Error(err)
	}
}
type Duo struct {
	EventPeriodId int64
	BetTypeId int64
	Code int64
	APrice float64
	BPrice float64
	IsActive bool
}
func (s *Store) SaveDuo(t Duo) {
	_, err := s.db.Exec("dbo.uspSaveDuo",
		sql.Named("EventPeriodId", t.EventPeriodId),
		sql.Named("BetTypeId", t.BetTypeId),
		sql.Named("Code", t.Code),
		sql.Named("APrice", t.APrice),
		sql.Named("BPrice", t.BPrice),
		sql.Named("IsActive", t.IsActive),
	)
	if err != nil {
		s.log.Error(err)
	}
}

func (s *Store) DeactivateWinDrawWin(eventPeriodId int64) {
	_, _ = s.db.Exec("update dbo.WinDrawWin set IsActive = 0 where EventPeriodId = @p1", eventPeriodId)
}

func (s *Store) DeactivateDoubleChance(eventPeriodId int64) {
	_, _ = s.db.Exec("update dbo.DoubleChance set IsActive = 0 where EventPeriodId = @p1", eventPeriodId)
}

const GetDemoSurebetQ = `

with t as (
    select EventPeriodId,
           HandicapCode,
           Margin,
           EventId,
           PeriodCode,
           concat('for,', 'ah,', 'h,', HandicapCode) HomeBetType,
           concat('for,', 'ah,', 'a,', HandicapCode) AwayBetType
    from Handicap h
             join EventPeriod ep on ep.Id = EventPeriodId
    where h.IsActive = 1
      and ep.IsActive = 1
      and Margin > @p1
    union all
    select EventPeriodId,
           HandicapCode,
           Margin,
           EventId,
           PeriodCode,
           concat('for,', 'ahover,', HandicapCode)  HomeBetType,
           concat('for,', 'ahunder,', HandicapCode) AwayBetType
    from Total t
             join EventPeriod ep on ep.Id = EventPeriodId
    where t.IsActive = 1
      and ep.IsActive = 1
      and Margin > @p1
)
select top 1 EventPeriodId, HandicapCode, Margin, EventId, PeriodCode, HomeBetType, AwayBetType
from t
order by Margin desc
`

type Surebet struct {
	EventId       string
	EventPeriodId int64
	HandicapCode  int64
	Away          float64
	Home          float64
	Margin        float64
	PeriodCode    string
	HomeBetType   string
	AwayBetType   string
}

func (s *Store) GetDemoSurebet(minMargin float64) (surebet Surebet, err error) {
	err = s.db.Get(&surebet, GetDemoSurebetQ, minMargin)
	return
}

func (s *Store) GetEventPeriodId(eventId string, code string) (eventPeriodId int64, err error) {
	got, b := s.Cache.Get(eventId + code)
	if b {
		return got.(int64), nil
	}
	err = s.db.Get(&eventPeriodId, "select Id from dbo.EventPeriod where EventId = @p1 and PeriodCode = @p2", eventId, code)
	if err != nil {
		return
	}
	s.Cache.SetWithTTL(eventId+code, eventPeriodId, 1, time.Hour)
	return
}

type Price struct {
	BetslipId string
	BetType   string
	Bookie    string
	EventId   string
	Sport     string
	Status    string
	Username  string
	Num       int
	Price     float64
	Min       float64
	Max       float64
}

func (s *Store) SavePrice(price Price) {
	_, err := s.db.Exec("dbo.uspSavePrice",
		sql.Named("BetslipId", price.BetslipId),
		sql.Named("Bookie", price.Bookie),
		sql.Named("BetType", price.BetType),
		sql.Named("Num", price.Num),
		sql.Named("Price", price.Price),
		sql.Named("Min", price.Min),
		sql.Named("Max", price.Max),
		sql.Named("Status", price.Status),
		sql.Named("IsActive", 1),
	)
	if err != nil {
		s.log.Error(err)
		s.log.Infow("price", "", price)
	}
}

func (s *Store) DeactivatePrice(price Price) {
	_, err := s.db.Exec("update dbo.Price set IsActive = 0 where BetslipId = @p1 and Bookie = @p2 and BetType = @p3",
		price.BetslipId, price.Bookie, price.BetType)
	if err != nil {
		s.log.Error(err)
	}
}

func (s *Store) SaveBetSlip(b api.BetSlipData) {
	//s.log.Infow("", "slip", b)
	_, err := s.db.Exec("dbo.uspSaveBetSlip",
		sql.Named("BetslipId", b.BetslipId),
		sql.Named("EventId", b.EventId),
		sql.Named("SportCode", b.Sport),
		sql.Named("BetType", b.BetType),
		sql.Named("BetTypeDes", b.BetTypeDescription),
		sql.Named("BetTypeTemp", b.BetTypeTemplate),
		sql.Named("EquivalentBets", b.EquivalentBets),
		sql.Named("MultipleAccounts", b.MultipleAccounts),
		sql.Named("IsOpen", b.IsOpen),
		sql.Named("ExpiryTs", b.ExpiryTs),
	)
	if err != nil {
		s.log.Error(err)
	}
}

func (s *Store) HasBetSlip(sportCode string, eventId string, betType string) (betSlipId string) {
	err := s.db.Get(&betSlipId, "select BetslipId from dbo.BetSlip where SportCode = @p1 and EventId = @p2 and BetType = @p3",
		sportCode, eventId, betType)
	if err == sql.ErrNoRows {
		return
	} else if err != nil {
		s.log.Error(err)
	}
	return
}

func (s *Store) DeleteBetSlips() {
	_, err := s.db.Exec("truncate table dbo.BetSlip")
	if err != nil {
		s.log.Error(err)
	}
}

const getPriceQ = `
with t as (
    select max(Price)                  Price,
           sum(Price * Max) / sum(Max) WeightedPrice,
           sum(Max)                    Volume,
           count(Price)                PriceCount,
           BetslipId
    from Price
    where Price > 0
      and IsActive = 1
      and BetslipId = @p1
    group by BetslipId
)
select top 1 p.BetslipId,
             bs.EventId,
             bs.SportCode PeriodCode,
             p.Price BestPrice,
             WeightedPrice,
             Min,
             Max,
             Volume,
             Bookie,
             PriceCount
from Price p
         join t on p.BetslipId = t.BetslipId and p.Price = t.Price
join BetSlip bs on bs.BetslipId = p.BetslipId
order by Max desc
`
type Side struct {
	SurebetId     int64
	BetslipId     string
	EventId       string
	PeriodCode    string
	Price         float64
	BestPrice     float64
	WeightedPrice float64
	Min           float64
	Max           float64
	Volume        float64
	Bookie        string
	PriceCount    int64
}

func (s *Store) GetPrice(betSlipId string) (side Side, err error) {
	err = s.db.Get(&side, getPriceQ, betSlipId)
	return
}

const saveSurebetQ = `
insert into dbo.Surebet (SurebetId, BetslipId, EventId, PeriodCode, Price, BestPrice, WeightedPrice, Min, Max, Volume, Bookie, PriceCount)
VALUES (:SurebetId, :BetslipId, :EventId, :PeriodCode, :Price, :BestPrice, :WeightedPrice, :Min, :Max, :Volume, :Bookie, :PriceCount) 
`

func (s *Store) SaveSurebet(price Side) (err error) {
	_, err = s.db.NamedExec(saveSurebetQ, price)
	return
}
func (s *Store) DeleteTotals() {
	_, err := s.db.Exec(`delete t from Total t join EventPeriod ep on Id = EventPeriodId where ep.IsActive = 0`)
	if err != nil {
		s.log.Error(err)
	}
}
func (s *Store) DeleteHandicaps() {
	_, err := s.db.Exec(`delete h from dbo.Handicap h join EventPeriod ep on Id = EventPeriodId where ep.IsActive = 0`)
	if err != nil {
		s.log.Error(err)
	}
}
func (s *Store) DeleteWinDrawWins() {
	_, err := s.db.Exec(`delete w from dbo.WinDrawWin w join EventPeriod ep on Id = EventPeriodId where ep.IsActive = 0`)
	if err != nil {
		s.log.Error(err)
	}
}
