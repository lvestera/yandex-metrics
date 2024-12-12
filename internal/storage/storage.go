package storage

import (
	"github.com/lvestera/yandex-metrics/internal/server/config"
)

func NewStorageRepository(cfg *config.Config) (Repository, error) {
	var r Repository
	var err error
	if len(cfg.DBConfig) > 0 {
		r, err = NewDBRepository(cfg.DBConfig)
		if err != nil {
			return nil, err
		}

	} else {
		r, err = NewMemStorage(cfg.Restore, cfg.FileStoragePath)
		if err != nil {
			return nil, err
		}
	}

	return r, nil
}

type Ping interface {
	Ping() error
}
