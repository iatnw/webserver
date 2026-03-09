package main

import "net/http"

func main() {
  const filePathRoot = "."
  const port = ":8080"

  mux := http.NewServeMux()
  mux.Handle("/", http.FileServer(http.Dir(filePathRoot)))
  mux.Handle("/logo.png", http.FileServer(http.Dir(filePathRoot+"/assets")))

  server := &http.Server{Addr: port, Handler: mux}
  err := server.ListenAndServe()
  if err != nil {
    return
  }
}
