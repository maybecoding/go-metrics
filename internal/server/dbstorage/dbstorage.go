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
func (ds *DBStorage) SetAll(mts []*metric.Metrics) error {
	cols := []string{"type", "name", "value", "delta"}
	cpRw := make([][]interface{}, 0, len(mts))
	for _, mt := range mts {
		cpRw = append(cpRw, []interface{}{mt.MType, mt.ID, mt.Value, mt.Delta})
	}
	tx, err := ds.conn.Begin(context.Background())
	if err != nil {
		return fmt.Errorf("error due begin transaction: %w", err)
	}
	defer func() {
		_ = tx.Rollback(context.Background())
	}()

	_, err = tx.Exec(context.Background(), `
	drop table if exists metric_tmp;
	create local temporary table metric_tmp (
		type metric_type not null,
		name varchar(255) not null,
		value double precision null,
		delta int8 null
	)`)
	if err != nil {
		return fmt.Errorf("error due create tmp table for data loading: %w", err)
	}

	_, err = tx.CopyFrom(context.Background(), pgx.Identifier{"metric_tmp"}, cols, pgx.CopyFromRows(cpRw))
	if err != nil {
		return fmt.Errorf("error due data loading into tmp table: %w", err)
	}

	_, err = tx.Exec(context.Background(), `insert into metric (type, name, value, delta)
		select type, name, value, delta
		from metric_tmp
		on conflict(type, name) do update set
			delta = metric.delta + EXCLUDED.delta,
			value = EXCLUDED.value;
		drop table if exists metric_tmp;
	`)

	if err != nil {
		return fmt.Errorf("error due data loading into metrics: %w", err)
	}

	err = tx.Commit(context.Background())
	if err != nil {
		return fmt.Errorf("error due commit transaction: %w", err)
	}

	return nil
}

func (ds *DBStorage) Ping() error {
	return ds.conn.Ping(context.Background())
}
