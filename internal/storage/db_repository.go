package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/lvestera/yandex-metrics/internal/models"
	"github.com/lvestera/yandex-metrics/internal/server/logger"
)

const maxRetries = 3
const defaultDelay = 1

type DBRepository struct {
	DB  *sql.DB
	rwm sync.RWMutex
}

func NewDBRepository(configStr string) (*DBRepository, error) {
	db, err := sql.Open("pgx", configStr)
	if err != nil {
		return nil, err
	}

	//create repository
	rep := &DBRepository{DB: db}

	//проверить подключение к базе
	if err = rep.Ping(); err != nil {
		logger.Log.Info("Error db connection " + err.Error())
		return nil, err
	}

	logger.Log.Info("Db connection OK")
	logger.Log.Info("DB string " + configStr)

	//проверить есть ли таблица
	rows, tableCheck := db.Query("select 1 from metrics;")
	if tableCheck != nil {
		//создать таблицу
		// ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		// defer cancel()
		_, err := db.ExecContext(context.Background(), "CREATE TABLE IF NOT EXISTS metrics ("+
			"id VARCHAR(100) PRIMARY KEY,"+
			"type varchar(10) NOT NULL,"+
			"delta bigint,"+
			"gauge double precision"+
			");")
		if err != nil {
			return nil, err
		}
	} else {
		err = rows.Err()
		if err != nil {
			logger.Log.Error(err.Error())
			return nil, err
		}
	}

	return rep, nil
}

func (rep *DBRepository) GetMetrics() ([]models.Metric, error) {
	rep.rwm.Lock()
	defer rep.rwm.Unlock()

	metrics := make([]models.Metric, 0)

	delay := defaultDelay
	for i := 0; i < maxRetries; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(delay)*time.Second)
		defer cancel()
		rows, err := rep.DB.QueryContext(ctx,
			"SELECT * FROM metrics",
		)

		if err != nil {
			logger.Log.Error(fmt.Sprint("Error while reading from db (", i, " attempt): ", err.Error()))

		} else {
			var m models.Metric

			for rows.Next() {
				rows.Scan(&m.ID, &m.MType, &m.Delta, &m.Value)
				metrics = append(metrics, m)
			}

			err = rows.Err()
			if err != nil {
				return nil, err
			}

			return metrics, nil
		}

		defer rows.Close()

		delay += 2
	}

	return metrics, errors.New(fmt.Sprint("error while reading from database after ", maxRetries, " attempts"))
}

func (rep *DBRepository) GetMetric(mtype string, name string) (m models.Metric, err error) {
	rep.rwm.Lock()
	defer rep.rwm.Unlock()

	delay := defaultDelay
	for i := 0; i < maxRetries; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(delay)*time.Second)
		defer cancel()

		err = rep.DB.QueryRowContext(
			ctx,
			"SELECT * FROM metrics WHERE ID=$1 AND TYPE=$2", name, mtype,
		).Scan(&m.ID, &m.MType, &m.Delta, &m.Value)
		if err != nil {
			logger.Log.Error(fmt.Sprint("Error while reading from db (", i, " attempt): ", err.Error()))

		} else {
			return m, nil
		}

		delay += 2
	}

	return m, errors.New(fmt.Sprint("error while reading from database after ", maxRetries, " attempts"))
}

func (rep *DBRepository) AddMetrics(metrics []models.Metric) (int, error) {
	tx, err := rep.DB.Begin()
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	stmt, err := tx.PrepareContext(
		context.Background(),
		"INSERT INTO metrics (id, type, delta, gauge) VALUES ($1, $2, $3, $4) "+
			"ON CONFLICT (id) DO UPDATE SET delta=CAST(metrics.delta AS INTEGER)+CAST($3 AS INTEGER), gauge=$4")
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	count := 0
	for _, m := range metrics {
		_, err := stmt.ExecContext(context.Background(), m.ID, m.MType, m.Delta, m.Value)
		if err != nil {
			return count, err
		}
		count = count + 1

	}
	return count, tx.Commit()
}

func (rep *DBRepository) AddMetric(m models.Metric) (bool, error) {
	rep.rwm.Lock()
	defer rep.rwm.Unlock()

	_, err := rep.DB.ExecContext(
		context.Background(),
		"INSERT INTO metrics (id, type, delta, gauge) VALUES ($1, $2, $3, $4) ON CONFLICT (id) DO UPDATE SET delta=CAST(metrics.delta AS INTEGER)+CAST($3 AS INTEGER), gauge=$4",
		m.ID, m.MType, m.Delta, m.Value,
	)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (rep *DBRepository) SetGauges(gauges map[string]float64) {

}

func (rep *DBRepository) Save(interval int) error {
	return nil
}

func (rep *DBRepository) Ping() error {
	return rep.DB.PingContext(context.Background())
}
