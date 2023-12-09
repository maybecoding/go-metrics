package dbstorage

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/maybecoding/go-metrics.git/internal/server/metric"
	"github.com/maybecoding/go-metrics.git/pkg/logger"
)

type DBStorage struct {
	conn *pgx.Conn
}

func New(connStr string) *DBStorage {
	// Соединяемся с базой данных
	conn, err := pgx.Connect(context.Background(), connStr)
	if err != nil {
		logger.Log.Fatal().Err(err).Msg("can't connect to database")
	}
	// Запускаем миграции
	err = runMigrations(connStr)
	if err != nil {
		logger.Log.Fatal().Err(err).Msg("can't run migrations")
	}
	return &DBStorage{
		conn: conn,
	}
}

func (ds *DBStorage) Set(mt *metric.Metrics) error {
	var err error
	if mt.MType == metric.Gauge {
		err = ds.conn.QueryRow(context.Background(), sqlSetMetricGauge, mt.MType, mt.ID, *mt.Value).
			Scan(&mt.MType, &mt.ID, &mt.Value)
	} else { // Слой выше это проверяет
		err = ds.conn.QueryRow(context.Background(), sqlSetMetricCounter, mt.MType, mt.ID, *mt.Delta).
			Scan(&mt.MType, &mt.ID, &mt.Delta)
	}
	if err != nil {
		return fmt.Errorf("error due scan after metric update: %w", err)
	}
	return nil
}
func (ds *DBStorage) Get(mt *metric.Metrics) error {
	row := ds.conn.QueryRow(context.Background(), sqlGetMetric, mt.MType, mt.ID)

	err := row.Scan(&mt.MType, &mt.ID, &mt.Delta, &mt.Value)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return metric.ErrNoMetricValue
		}
		return fmt.Errorf("error due get metric: %w", err)
	}

	return nil
}
func (ds *DBStorage) GetAll() ([]*metric.Metrics, error) {
	rows, err := ds.conn.Query(context.Background(), sqlGetMetricAll)
	if err != nil {
		return nil, fmt.Errorf("error due select all metrics: %w", err)
	}
	mts := make([]*metric.Metrics, 0)
	for rows.Next() {
		mt := metric.Metrics{}
		err = rows.Scan(&mt.MType, &mt.ID, &mt.Delta, &mt.Value)
		if err != nil {
			return nil, fmt.Errorf("error due scan row from select all metrics: %w", err)
		}
		mts = append(mts, &mt)
	}

	return mts, nil
}

func (ds *DBStorage) Ping() error {
	return ds.conn.Ping(context.Background())
}
