package storage

import (
	"context"
	"database/sql"
	"sync"

	"github.com/lvestera/yandex-metrics/internal/models"
	"github.com/lvestera/yandex-metrics/internal/server/logger"
)

type DBRepository struct {
	DB  *sql.DB
	rwm sync.RWMutex
}

func NewDBRepository(configStr string) (*DBRepository, error) {
	db, err := sql.Open("pgx", configStr)
	if err != nil {
		return nil, err
	}
	//defer db.Close()

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
	_, table_check := db.Query("select 1 from metrics;")
	if table_check != nil {
		//создать таблицу
		// ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		// defer cancel()
		_, err := db.ExecContext(context.Background(), "CREATE TABLE IF NOT EXISTS metrics ("+
			"id VARCHAR(100) PRIMARY KEY,"+
			"type varchar(10) NOT NULL,"+
			"delta int,"+
			"gauge double precision"+
			");")
		if err != nil {
			return nil, err
		}
	}

	return rep, nil
}

func (rep *DBRepository) GetMetrics() ([]models.Metric, error) {
	rep.rwm.Lock()
	defer rep.rwm.Unlock()

	metrics := make([]models.Metric, 0)

	rows, err := rep.DB.QueryContext(context.Background(),
		"SELECT * FROM metrics",
	)
	if err != nil {
		return metrics, err
	}

	defer rows.Close()

	var m models.Metric

	for rows.Next() {
		rows.Scan(&m.ID, &m.MType, &m.Delta, &m.Value)
		metrics = append(metrics, m)
	}

	return metrics, nil
}

func (rep *DBRepository) GetMetric(mtype string, name string) (m models.Metric, err error) {
	rep.rwm.Lock()
	defer rep.rwm.Unlock()

	//var m models.Metric
	err = rep.DB.QueryRowContext(
		context.Background(),
		"SELECT * FROM metrics WHERE ID=$1 AND TYPE=$2", name, mtype,
	).Scan(&m.ID, &m.MType, &m.Delta, &m.Value)
	if err != nil {
		return m, err
	}

	return m, nil
}

func (rep *DBRepository) AddMetric(m models.Metric) (bool, error) {
	rep.rwm.Lock()
	defer rep.rwm.Unlock()

	_, err := rep.DB.ExecContext(
		context.Background(),
		"INSERT INTO metrics (id, type, delta, gauge) VALUES ($1, $2, $3, $4) ON CONFLICT (id) DO UPDATE SET delta=CAST(metrics.delta AS INTEGER)+CAST($3 AS INTEGER), gauge=$4", m.ID, m.MType, m.Delta, m.Value,
	)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (dbr *DBRepository) SetGauges(gauges map[string]float64) {

}

func (dbr *DBRepository) Save(interval int) error {
	return nil
}

func (dbr *DBRepository) Ping() error {
	return dbr.DB.PingContext(context.Background())
}
