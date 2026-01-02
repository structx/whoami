package main

import (
	"encoding/json"
	"flag"
	"net"
	"net/http"
	"os"
	"os/user"
	"strings"

	tea "github.com/structx/teapot"
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

	log *tea.Logger
)

func init() {
	flag.StringVar(&port, "port", getEnv("PORT", defaultPort), "port number")
	flag.StringVar(&host, "host", getEnv("HOST", defaultHost), "host")
	flag.StringVar(&logLevel, "log_level", getEnv("LOG_LEVEL", defaultLogLevel), "log level")
}

func getHostname() string {
	h, _ := os.Hostname()
	return h
}

func whoamiHandler(w http.ResponseWriter, r *http.Request) {
	ip, _, _ := net.SplitHostPort(r.RemoteAddr)
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		ip = strings.Split(xff, ",")[0]
	}

	u, err := user.Current()
	if err != nil {
		log.Error("current user", tea.Error(err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	type resp struct {
		Hostname   string      `json:"hostname"`
		UID        string      `json:"uid"`
		GID        string      `json:"gid"`
		IP         string      `json:"ip"`
		RemoteAddr string      `json:"remote_addr"`
		Protocol   string      `json:"protocol"`
		Method     string      `json:"method"`
		URL        string      `json:"url"`
		Header     http.Header `json:"header"`
	}

	rr := &resp{
		Hostname:   getHostname(),
		UID:        u.Uid,
		GID:        u.Gid,
		IP:         ip,
		RemoteAddr: r.RemoteAddr,
		Protocol:   r.Proto,
		Method:     r.Method,
		URL:        r.RequestURI,
		Header:     r.Header,
	}

	log.Infof("received request",
		tea.String("hostname", rr.Hostname),
		tea.String("uid", rr.UID),
		tea.String("gid", rr.GID),
		tea.String("ip", rr.IP),
		tea.String("remote_addr", rr.RemoteAddr),
		tea.String("protocol", rr.Protocol),
		tea.String("method", rr.Method),
		tea.String("url", rr.URL),
		tea.Any("header", rr.Header),
	)

	w.Header().Add("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(&rr); err != nil {
		log.Error("encode json response", tea.Error(err))
	}
}

func healthHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func getEnv(key string, defaultValue string) string {
	v, ok := os.LookupEnv(key)
	if !ok {
		return defaultValue
	}
	return v
}

func newLogger() *tea.Logger {
	var l tea.Level
	switch strings.ToLower(logLevel) {
	case "error":
		l = tea.ERROR
	case "info":
		l = tea.INFO
	default:
		l = tea.DEBUG
	}

	return tea.New(
		tea.WithLevel(l),
	)
}

func main() {
	log = newLogger()

	flag.Parse()

	mux := http.NewServeMux()
	mux.HandleFunc("/", whoamiHandler)
	mux.HandleFunc("/health", healthHandler)

	hostAndPort := net.JoinHostPort(host, port)

	log.Infof("start http/1 server", tea.String("server_addr", hostAndPort))
	if err := http.ListenAndServe(hostAndPort, mux); err != nil && err != http.ErrServerClosed {
		log.Error("failed to start http/1 server", tea.Error(err))
	}
}
