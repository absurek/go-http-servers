package main

import "net/http"

const addr = ":8080"

func main() {
	mux := &http.ServeMux{}
	server := http.Server{
		Handler: mux,
		Addr:    addr,
	}

	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)

		w.Write([]byte("OK\n"))
	})

	mux.Handle("/app/", http.StripPrefix("/app/", http.FileServer(http.Dir("."))))

	server.ListenAndServe()
}
