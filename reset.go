package main

import "net/http"

func (cfg *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {
	if cfg.platform != "dev" {
		http.Error(w, "Reset endpoint is only available in dev environment", http.StatusForbidden)
		respondWithError(w, http.StatusForbidden, "Reset endpoint is only available in dev environment", nil)
	}
	cfg.dbQueries.ResetUsers(r.Context())
	cfg.fileserverHits.Store(0)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hits reset to 0, Users table cleared"))
}