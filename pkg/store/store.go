package store

import (
	"context"
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
	err = s.db.GetContext(ctx, &token, "select top 1 Session, CreatedAt, LastCheckAt from dbo.Auth")
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
