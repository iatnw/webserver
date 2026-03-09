package main

import "net/http"

func main() {
	const filePathRoot = "."
	const port = ":8080"

	mux := http.NewServeMux()
	mux.Handle("/app/", http.StripPrefix("/app", http.FileServer(http.Dir(filePathRoot))))
	mux.HandleFunc("/healthz", healthzHandler)

	server := &http.Server{Addr: port, Handler: mux}
	err := server.ListenAndServe()
	if err != nil {
		return
	}
}
