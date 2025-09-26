package postgres

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
	"github.com/rys730/iFortepay-take-home/internal/common/config"
)

type ShopDB struct {
	Conn *pgxpool.Pool
}

// TODO: env db creds
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
