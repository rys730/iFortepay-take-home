package postgres

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
	"github.com/rys730/iFortepay-take-home/internal/common/config"
)

type Pool interface {
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

type ShopDB struct {
	Conn Pool
}

func NewPostgres(dbCfg *config.DBConfig) *ShopDB {
	connString := dbCfg.GetConnectionString()
	pool, err := pgxpool.New(context.Background(), connString)
	if err != nil {
		log.Fatal().Err(err).Msg("failed creating db pool")
	}

	if err := pool.Ping(context.Background()); err != nil {
		log.Fatal().Err(err).Msg("failed pinging db")
	}
	return &ShopDB{
		Conn: pool,
	}
}
