#!/usr/bin/env bash
set -e

DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$DIR"

cat <<'EOF' > main.go
package main

import (
	"fmt"
	"net/http"
	"os"
)

func main() {
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "OK: user=%d\n", os.Getuid())
	})
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Go Microservice v1.0\n")
	})

	fmt.Println("Server listening on :8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Fprintf(os.Stderr, "Server error: %v\n", err)
		os.Exit(1)
	}
}
EOF

cat <<'EOF' > go.mod
module example.com/microservice

go 1.22
EOF

cat <<'EOF' > Dockerfile.orig
FROM golang:1.23-alpine
WORKDIR /app
COPY . .
RUN go build -o /app/server .
EXPOSE 8080
CMD ["/app/server"]
EOF

echo "[PASS] Environment prepared for img-01-multistage-build in $DIR"
