package main

import (
	"flag"
	"log/slog"
	"net"
	"net/http"
	"os"
	"strings"
)

const (
	defaultPort     = "8080"
	defaultHost     = "127.0.0.1"
	defaultLogLevel = "DEBUG"
)

var (
	port     string
	host     string
	logLevel string

	buildDate string
	commitSHA string

	log *slog.Logger
)

func init() {
	flag.StringVar(&port, "port", getEnv("PORT", defaultPort), "port number")
	flag.StringVar(&host, "host", getEnv("HOST", defaultHost), "host")
	flag.StringVar(&logLevel, "log_level", getEnv("LOG_LEVEL", defaultLogLevel), "log level")
}

func whoamiHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log.InfoContext(ctx, "received request",
		slog.String("method", r.Method),
		slog.String("request_uri", r.RequestURI),
	)
}

func getEnv(key string, defaultValue string) string {
	v, ok := os.LookupEnv(key)
	if !ok {
		return defaultValue
	}
	return v
}

func main() {
	var l slog.Level
	switch strings.ToLower(logLevel) {
	case "error":
		l = slog.LevelError
	case "info":
		l = slog.LevelInfo
	default:
		l = slog.LevelDebug
	}

	log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: l}))

	flag.Parse()

	mux := http.NewServeMux()
	mux.HandleFunc("/whoami", whoamiHandler)

	hostAndPort := net.JoinHostPort(host, port)

	log.Info("start http/1 server", slog.String("server_addr", hostAndPort))
	if err := http.ListenAndServe(hostAndPort, mux); err != nil && err != http.ErrServerClosed {
		log.Error("failed to start http/1 server", slog.String("error", err.Error()))
	}
}
