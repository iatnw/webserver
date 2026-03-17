package main

import (
	"fmt"
	"net/http"
)

func (apiConfig *apiConfig) handlerMetrics(w http.ResponseWriter, _ *http.Request) {
	w.Header().Add("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf(`<html>
<body>
<h1>Welcome, Chirpy %s</h1>
<p> Chirpy has been visited %d times!</p>
</body>
</html>`, "Admin", apiConfig.fileserverHits.Load())))
}

func (apiConfig *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiConfig.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}
