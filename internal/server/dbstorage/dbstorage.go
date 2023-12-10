package dbstorage

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/maybecoding/go-metrics.git/internal/server/metric"
	"github.com/maybecoding/go-metrics.git/pkg/logger"
	"time"
)

type DBStorage struct {
	conn           *pgx.Conn
	ctx            context.Context
	retryIntervals []time.Duration
}

func New(connStr string, ctx context.Context, retryIntervals []time.Duration) *DBStorage {
	// Соединяемся с базой данных
	conn, err := pgx.Connect(ctx, connStr)
	if err != nil {
		logger.Log.Fatal().Err(err).Msg("can't connect to database")
	}
	// Запускаем миграции
	err = runMigrations(connStr)
	if err != nil {
		logger.Log.Fatal().Err(err).Msg("can't run migrations")
	}
	return &DBStorage{
		conn:           conn,
		ctx:            ctx,
		retryIntervals: retryIntervals,
	}
}
func (ds *DBStorage) get(mt *metric.Metrics) error {
	row := ds.conn.QueryRow(ds.ctx, sqlGetMetric, mt.MType, mt.ID)

	err := row.Scan(&mt.MType, &mt.ID, &mt.Delta, &mt.Value)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return metric.ErrNoMetricValue
		}
		return fmt.Errorf("error due get metric: %w", err)
	}

	return nil
}

func (ds *DBStorage) set(mt *metric.Metrics) error {
	var err error
	if mt.MType == metric.Gauge {
		err = ds.conn.QueryRow(ds.ctx, sqlSetMetricGauge, mt.MType, mt.ID, *mt.Value).
			Scan(&mt.MType, &mt.ID, &mt.Value)
	} else { // Слой выше это проверяет
		err = ds.conn.QueryRow(ds.ctx, sqlSetMetricCounter, mt.MType, mt.ID, *mt.Delta).
			Scan(&mt.MType, &mt.ID, &mt.Delta)
	}
	if err != nil {
		return fmt.Errorf("error due scan after metric update: %w", err)
	}
	return nil
}

func (ds *DBStorage) getAll() ([]*metric.Metrics, error) {
	rows, err := ds.conn.Query(ds.ctx, sqlGetMetricAll)
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
func (ds *DBStorage) setAll(mts []*metric.Metrics) error {
	cols := []string{"type", "name", "value", "delta"}
	cpRw := make([][]interface{}, 0, len(mts))
	for _, mt := range mts {
		cpRw = append(cpRw, []interface{}{mt.MType, mt.ID, mt.Value, mt.Delta})
	}
	tx, err := ds.conn.Begin(ds.ctx)
	if err != nil {
		return fmt.Errorf("error due begin transaction: %w", err)
	}
	defer func() {
		_ = tx.Rollback(ds.ctx)
	}()

	_, err = tx.Exec(ds.ctx, `
	drop table if exists metric_tmp;
	create local temporary table metric_tmp (
	    id int generated always as identity,
		type metric_type not null,
		name varchar(255) not null,
		value double precision null,
		delta int8 null
	)`)
	if err != nil {
		return fmt.Errorf("error due create tmp table for data loading: %w", err)
	}

	_, err = tx.CopyFrom(ds.ctx, pgx.Identifier{"metric_tmp"}, cols, pgx.CopyFromRows(cpRw))
	if err != nil {
		return fmt.Errorf("error due data loading into tmp table: %w", err)
	}

	_, err = tx.Exec(ds.ctx, `
	with gauge as (
		select type, name, value, null::int8 as delta
			 ,row_number() over (partition by type, name order by id desc) rn
		from metric_tmp
		where type = 'gauge'
	), counter as (
		select type, name, null, null::double precision as value, sum(delta) as delta
		from metric_tmp
		where type = 'counter'
		group by type, name
	)
	insert into metric (type, name, value, delta)
	select type, name, value, delta
	from gauge
	where rn = 1
	union all
	select type, name, value, delta
	from counter
	on conflict(type, name) do update set
		delta = metric.delta + EXCLUDED.delta,
		value = EXCLUDED.value;
	drop table if exists metric_tmp;
	`)

	if err != nil {
		return fmt.Errorf("error due data loading into metrics: %w", err)
	}

	err = tx.Commit(ds.ctx)
	if err != nil {
		return fmt.Errorf("error due commit transaction: %w", err)
	}

	return nil
}

func (ds *DBStorage) Ping() error {
	return ds.conn.Ping(ds.ctx)
}

func (ds *DBStorage) Get(mt *metric.Metrics) (err error) {
	for _, ri := range ds.retryIntervals {
		err = ds.get(mt)
		var pgErr *pgconn.PgError
		if err == nil || !errors.Is(err, pgErr) || !pgerrcode.IsConnectionException(pgErr.Code) {
			return err
		}
		select {
		case <-time.After(ri):
			logger.Log.Debug().Err(err).Dur("duration", ri).Msg("try to retry after")
		case <-ds.ctx.Done():
			return ds.ctx.Err()
		}
	}
	return
}

func (ds *DBStorage) Set(mt *metric.Metrics) (err error) {
	for _, ri := range ds.retryIntervals {
		err = ds.set(mt)
		var pgErr *pgconn.PgError
		if err == nil || !errors.Is(err, pgErr) || !pgerrcode.IsConnectionException(pgErr.Code) {
			return err
		}
		select {
		case <-time.After(ri):
			logger.Log.Debug().Err(err).Dur("duration", ri).Msg("try to retry after")
		case <-ds.ctx.Done():
			return ds.ctx.Err()
		}
	}
	return
}

func (ds *DBStorage) GetAll() (mts []*metric.Metrics, err error) {

	for _, ri := range ds.retryIntervals {
		mts, err = ds.getAll()
		var pgErr *pgconn.PgError
		if err == nil || !errors.Is(err, pgErr) || !pgerrcode.IsConnectionException(pgErr.Code) {
			return mts, err
		}
		select {
		case <-time.After(ri):
			logger.Log.Debug().Err(err).Dur("duration", ri).Msg("try to retry after")
		case <-ds.ctx.Done():
			return nil, ds.ctx.Err()
		}
	}
	return
}

func (ds *DBStorage) SetAll(mts []*metric.Metrics) (err error) {
	for _, ri := range ds.retryIntervals {
		err = ds.setAll(mts)
		var pgErr *pgconn.PgError
		if err == nil || !errors.Is(err, pgErr) || !pgerrcode.IsConnectionException(pgErr.Code) {
			return err
		}
		select {
		case <-time.After(ri):
			logger.Log.Debug().Err(err).Dur("duration", ri).Msg("try to retry after")
		case <-ds.ctx.Done():
			return ds.ctx.Err()
		}
	}
	return

}
