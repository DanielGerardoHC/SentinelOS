package handlers

import (
	"net/http"

	"sentinelos/core/internal/system/config_engine"
)

func CommitHandler(w http.ResponseWriter, r *http.Request) {

	err := config_engine.Commit()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Write([]byte("commit successful"))
}