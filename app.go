package main

import (
    "log"
    "net/http"
    "time"

    "gowes/routes"
)

func loggingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        next.ServeHTTP(w, r)
        log.Printf("%s %s %s", r.Method, r.URL.Path, time.Since(start))
    })
}

// corsMiddleware menambahkan header CORS dan menangani preflight OPTIONS
func corsMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Tetapkan origin dinamis. Jika tidak ada, gunakan wildcard.
        origin := r.Header.Get("Origin")
        if origin == "" {
            origin = "*"
        }
        w.Header().Set("Access-Control-Allow-Origin", origin)
        w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
        // Echo header yang diminta pada preflight agar tidak gagal karena header kustom.
        reqHeaders := r.Header.Get("Access-Control-Request-Headers")
        if reqHeaders != "" {
            w.Header().Set("Access-Control-Allow-Headers", reqHeaders)
        } else {
            w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
        }
        // Kurangi frekuensi preflight di browser yang mendukung
        w.Header().Set("Access-Control-Max-Age", "86400")
        w.Header().Add("Vary", "Origin")
        w.Header().Add("Vary", "Access-Control-Request-Method")
        w.Header().Add("Vary", "Access-Control-Request-Headers")

        if r.Method == http.MethodOptions {
            // Preflight: cukup kembalikan 204/200 dengan header CORS
            log.Printf("CORS preflight origin=%s method=%s headers=%s",
                origin,
                r.Header.Get("Access-Control-Request-Method"),
                reqHeaders,
            )
            w.WriteHeader(http.StatusNoContent)
            return
        }
        next.ServeHTTP(w, r)
    })
}

func main() {
    mux := http.NewServeMux()
    routes.RegisterTodoRoutes(mux)

    server := &http.Server{
        Addr:         ":8080",
        Handler:      loggingMiddleware(corsMiddleware(mux)),
        ReadTimeout:  5 * time.Second,
        WriteTimeout: 10 * time.Second,
        IdleTimeout:  60 * time.Second,
    }

    log.Println("Server berjalan di http://localhost:8080")
    if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
        log.Fatalf("server error: %v", err)
    }
}
