package main

import "net/http"

func (apiConfig *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {
	if apiConfig.platform != "dev" {
		respondWithError(w, http.StatusForbidden, "Forbidden", nil)
		return
	}
	err := apiConfig.db.DeleteAllUsers(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not delete users", err)
		return
	}
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Resetting metrics counter"))
	apiConfig.fileserverHits.Store(0)
}
