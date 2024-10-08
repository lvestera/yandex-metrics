package handlers

import (
	"context"
	"database/sql"
	"net/http"
	"time"
)

type PingHandler struct {
	Db *sql.DB
}

func (ph PingHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	if err := ph.Db.PingContext(ctx); err != nil {
		w.WriteHeader(http.StatusInternalServerError)

	} else {
		w.WriteHeader(http.StatusOK)
	}
}
