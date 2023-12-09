package dbstorage

import (
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/maybecoding/go-metrics.git/internal/server/metric"
	"github.com/maybecoding/go-metrics.git/pkg/logger"
)

type DbStorage struct {
	conn *pgx.Conn
}

func New(connStr string) *DbStorage {
	conn, err := pgx.Connect(context.Background(), connStr)
	if err != nil {
		logger.Log.Fatal().Err(err).Msg("can't connect to database")
	}
	return &DbStorage{
		conn: conn,
	}
}

func (ds *DbStorage) Set(_ *metric.Metrics) error {
	logger.Log.Fatal().Msg("method Set is not implemented in DbStorage")
	return nil
}
func (ds *DbStorage) Get(_ *metric.Metrics) error {
	logger.Log.Fatal().Msg("method Get is not implemented in DbStorage")
	return nil
}
func (ds *DbStorage) GetAll() ([]*metric.Metrics, error) {
	logger.Log.Fatal().Msg("method GetAll is not implemented in DbStorage")
	return nil, nil
}

func (ds *DbStorage) Ping() error {
	return ds.conn.Ping(context.Background())
}
