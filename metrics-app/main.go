package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"math/rand/v2"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/redis/go-redis/v9"
)

// Injected at build time via ldflags.
var (
	version   = "dev"
	buildDate = "0"
)

const port = ":8080"

var (
	hostname string
	rdb      redis.UniversalClient

	appInfo = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "app_info",
		Help: "Information about application",
	}, []string{"version", "build_timestamp"})

	fakeLoad = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "fake_load",
		Help: "Fake load Gauge",
	})

	httpRequestsTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "http_requests_total",
		Help: "Total number of HTTP requests made.",
	}, []string{"code", "handler", "method"})

	httpDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "http_request_duration_seconds",
		Help:    "The HTTP request latencies in seconds.",
		Buckets: prometheus.DefBuckets,
	}, []string{"handler"})
)

// responseWriter wraps http.ResponseWriter to capture the status code.
type responseWriter struct {
	http.ResponseWriter
	status int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}

// instrument wraps a handler with Prometheus metrics tracking.
func instrument(name string, h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rw := &responseWriter{ResponseWriter: w, status: http.StatusOK}
		h(rw, r)
		httpDuration.WithLabelValues(name).Observe(time.Since(start).Seconds())
		httpRequestsTotal.WithLabelValues(strconv.Itoa(rw.status), name, strings.ToLower(r.Method)).Inc()
	}
}

func defaultHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Request from %s served by %s\n", r.RemoteAddr, hostname)
}

func healthzHandler(w http.ResponseWriter, _ *http.Request) {
	fmt.Fprint(w, "OK")
}

func sickzHandler(w http.ResponseWriter, _ *http.Request) {
	http.Error(w, "Internal Server Error", http.StatusInternalServerError)
}

func versionHandler(w http.ResponseWriter, _ *http.Request) {
	buildDateInt, _ := strconv.ParseInt(buildDate, 10, 64)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"commit_hash": version,
		"build_date":  buildDateInt,
	})
}

func getTrololo(w http.ResponseWriter, r *http.Request) {
	if rdb == nil {
		http.Error(w, "Redis not configured", http.StatusServiceUnavailable)
		return
	}
	val, err := rdb.Get(context.Background(), r.PathValue("key")).Result()
	if err == redis.Nil {
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}
	if err != nil {
		slog.Error("Redis GET failed", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	fmt.Fprint(w, val)
}

func putTrololo(w http.ResponseWriter, r *http.Request) {
	if rdb == nil {
		http.Error(w, "Redis not configured", http.StatusServiceUnavailable)
		return
	}
	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "HTTP header 'Content-Type: application/json' expected", http.StatusUnsupportedMediaType)
		return
	}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	if err := rdb.Set(context.Background(), r.PathValue("key"), body, 0).Err(); err != nil {
		slog.Error("Redis SET failed", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	fmt.Fprint(w, "Created")
}

func deleteTrololo(w http.ResponseWriter, r *http.Request) {
	if rdb == nil {
		http.Error(w, "Redis not configured", http.StatusServiceUnavailable)
		return
	}
	if err := rdb.Del(context.Background(), r.PathValue("key")).Err(); err != nil {
		slog.Error("Redis DEL failed", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func initRedis() {
	sentinelAddrs := os.Getenv("REDIS_SENTINEL_ADDRS")
	if sentinelAddrs == "" {
		slog.Warn("REDIS_SENTINEL_ADDRS not set, Redis features disabled")
		return
	}
	masterName := os.Getenv("REDIS_MASTER_NAME")
	if masterName == "" {
		masterName = "mymaster"
	}
	rdb = redis.NewFailoverClient(&redis.FailoverOptions{
		MasterName:    masterName,
		SentinelAddrs: strings.Split(sentinelAddrs, ","),
	})
	slog.Info("Redis sentinel initialized", "master", masterName, "sentinels", sentinelAddrs)
}

func updateFakeLoad() {
	for {
		fakeLoad.Set(rand.Float64() * 3)
		time.Sleep(5 * time.Second)
	}
}

func main() {
	var err error
	hostname, err = os.Hostname()
	if err != nil {
		hostname = "unknown"
	}

	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))

	appInfo.WithLabelValues(version, buildDate).Set(1)
	go updateFakeLoad()

	initRedis()

	mux := http.NewServeMux()
	mux.HandleFunc("GET /{$}", instrument("Default", defaultHandler))
	mux.HandleFunc("GET /healthz", instrument("Healthz", healthzHandler))
	mux.HandleFunc("GET /sickz", instrument("Sickz", sickzHandler))
	mux.HandleFunc("GET /version", instrument("Version", versionHandler))
	mux.Handle("GET /metrics", promhttp.Handler())
	mux.HandleFunc("GET /trollolo/{key}", instrument("GetTrololo", getTrololo))
	mux.HandleFunc("PUT /trollolo/{key}", instrument("PutTrololo", putTrololo))
	mux.HandleFunc("DELETE /trollolo/{key}", instrument("DeleteTrololo", deleteTrololo))

	slog.Info("Starting metrics-app", "port", port, "version", version)
	if err := http.ListenAndServe(port, mux); err != nil {
		slog.Error("Server failed", "error", err)
		os.Exit(1)
	}
}
