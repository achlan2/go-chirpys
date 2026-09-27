package main

import "net/http"

func (cfg *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {

	if cfg.platform != "dev" {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("Reset is only allowed in dev environtment."))
		return
	}
	cfg.fileserverHits.Store(0)
	err := cfg.db.Reset(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("Failed to reset the databse: " + err.Error()))
		return
	}
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hits reset to 0"))
}
