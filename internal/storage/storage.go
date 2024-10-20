package storage

import (
	"github.com/lvestera/yandex-metrics/internal/server/config"
)

type Storage struct {
	Cfg *config.Config
}

func (st Storage) InitRepository() (Repository, error) {
	var r Repository
	var err error
	if len(st.Cfg.DBConfig) > 0 {
		r, err = NewDBRepository(st.Cfg.DBConfig)
		if err != nil {
			return nil, err
		}

	} else {
		r, err = NewMemStorage(st.Cfg.Restore, st.Cfg.FileStoragePath)
		if err != nil {
			return nil, err
		}
	}

	return r, nil
}

type Ping interface {
	Ping() error
}
