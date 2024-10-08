package handlers

import (
	"context"
	"database/sql"
	"net/http"
	"time"
)

type PingHandler struct {
	DB *sql.DB
}

func (ph PingHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	if err := ph.DB.PingContext(ctx); err != nil {
		w.WriteHeader(http.StatusInternalServerError)

	} else {
		w.WriteHeader(http.StatusOK)
	}
}
