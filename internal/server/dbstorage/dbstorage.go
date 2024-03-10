package dbstorage

import (
	"context"
	"errors"
	"fmt"
	"github.com/maybecoding/go-metrics.git/internal/server/entity"
	"time"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/maybecoding/go-metrics.git/internal/server/metricservice"
	"github.com/maybecoding/go-metrics.git/pkg/logger"
)

type DBStorage struct {
	conn           *pgxpool.Pool
	ctx            context.Context
	retryIntervals []time.Duration
}

func New(connStr string, ctx context.Context, retryIntervals []time.Duration) *DBStorage {
	// Соединяемся с базой данных
	pool, err := pgxpool.New(ctx, connStr)
	if err != nil {
		logger.Fatal().Err(err).Msg("can't connect to database")
	}

	dbs := &DBStorage{
		conn:           pool,
		ctx:            ctx,
		retryIntervals: retryIntervals,
	}
	return dbs
}
func (ds *DBStorage) get(mt *entity.Metrics) error {
	row := ds.conn.QueryRow(ds.ctx, sqlGetMetric, mt.MType, mt.ID)

	err := row.Scan(&mt.MType, &mt.ID, &mt.Delta, &mt.Value)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return metricservice.ErrNoMetricValue
		}
		return fmt.Errorf("error due get metric: %w", err)
	}

	return nil
}

func (ds *DBStorage) set(mt entity.Metrics) error {
	var err error
	if mt.MType == metricservice.Gauge {
		logger.Debug().Str("MType", mt.MType).Str("ID", mt.ID).Float64("Value", *mt.Value).Msg("set value input params")
		_, err = ds.conn.Exec(ds.ctx, sqlSetMetricGauge, mt.MType, mt.ID, *mt.Value)
	} else { // Слой выше это проверяет
		_, err = ds.conn.Exec(ds.ctx, sqlSetMetricCounter, mt.MType, mt.ID, *mt.Delta)
	}
	if err != nil {
		return fmt.Errorf("error due scan after metric update: %w", err)
	}
	return nil
}

func (ds *DBStorage) getAll() ([]*entity.Metrics, error) {
	rows, _ := ds.conn.Query(ds.ctx, `select name, type, delta, value from metric`)
	mts, err := pgx.CollectRows(rows, pgx.RowToAddrOfStructByPos[entity.Metrics])
	if err != nil {
		return nil, fmt.Errorf("error due scan row from select all metrics: %w", err)
	}
	return mts, nil
}

func (ds *DBStorage) Ping() error {
	pe := ds.conn.Ping(ds.ctx)
	logger.Debug().Err(pe).Msg("db storage ping result")
	return pe
}

func (ds *DBStorage) Get(mt *entity.Metrics) (err error) {
	for _, ri := range ds.retryIntervals {
		err = ds.get(mt)
		var pgErr *pgconn.PgError

		if err != nil && errors.As(err, &pgErr) && pgerrcode.IsConnectionException(pgErr.Code) {
			select {
			case <-time.After(ri):
				logger.Debug().Err(err).Dur("duration", ri).Msg("try to retry after")
			case <-ds.ctx.Done():
				return ds.ctx.Err()
			}
		} else {
			return err
		}
	}
	return
}

func (ds *DBStorage) Set(mt entity.Metrics) (err error) {
	for _, ri := range ds.retryIntervals {
		err = ds.set(mt)

		var pgErr *pgconn.PgError
		if err != nil && errors.As(err, &pgErr) && pgerrcode.IsConnectionException(pgErr.Code) {
			select {
			case <-time.After(ri):
				logger.Debug().Err(err).Dur("duration", ri).Msg("try to retry after")
			case <-ds.ctx.Done():
				return ds.ctx.Err()
			}
		} else {
			return err
		}

	}
	return
}

func (ds *DBStorage) GetAll() (mts []*entity.Metrics, err error) {

	for _, ri := range ds.retryIntervals {
		mts, err = ds.getAll()
		var pgErr *pgconn.PgError
		if err != nil && errors.As(err, &pgErr) && pgerrcode.IsConnectionException(pgErr.Code) {
			select {
			case <-time.After(ri):
				logger.Debug().Err(err).Dur("duration", ri).Msg("try to retry after")
			case <-ds.ctx.Done():
				return nil, ds.ctx.Err()
			}
		} else {
			return mts, err
		}
	}
	return
}

func (ds *DBStorage) SetAll(mts []entity.Metrics) (err error) {
	for _, m := range mts {
		err = ds.set(m)
		if err != nil {
			return
		}
	}
	return
}

func (ds *DBStorage) ConnectionClose() {
	ds.conn.Close()
}
