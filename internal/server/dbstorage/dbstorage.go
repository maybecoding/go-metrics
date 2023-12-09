package dbstorage

import (
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/maybecoding/go-metrics.git/internal/server/metric"
	"github.com/maybecoding/go-metrics.git/pkg/logger"
)

type DBStorage struct {
	conn *pgx.Conn
}

func New(connStr string) *DBStorage {
	conn, err := pgx.Connect(context.Background(), connStr)
	if err != nil {
		logger.Log.Fatal().Err(err).Msg("can't connect to database")
	}
	return &DBStorage{
		conn: conn,
	}
}

func (ds *DBStorage) Set(_ *metric.Metrics) error {
	logger.Log.Fatal().Msg("method Set is not implemented in DBStorage")
	return nil
}
func (ds *DBStorage) Get(_ *metric.Metrics) error {
	logger.Log.Fatal().Msg("method Get is not implemented in DBStorage")
	return nil
}
func (ds *DBStorage) GetAll() ([]*metric.Metrics, error) {
	logger.Log.Fatal().Msg("method GetAll is not implemented in DBStorage")
	return nil, nil
}

func (ds *DBStorage) Ping() error {
	return ds.conn.Ping(context.Background())
}
