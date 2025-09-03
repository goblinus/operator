package healthz

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/go-logr/logr"

	"github.com/goblinus/operator/src/services/healthz/models"
	"github.com/goblinus/operator/src/utils"
)

type HealthChecker struct {
	server       *http.Server
	logger       logr.Logger
	operatorName string
	port         string
	isReady      bool
}

func NewHealthChecker(logger logr.Logger, operatorName, port string) *HealthChecker {
	return &HealthChecker{
		logger:       logger,
		operatorName: operatorName,
		port:         port,
		isReady:      false,
	}
}

func (h *HealthChecker) Start(ctx context.Context) {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", h.healthHandler)
	mux.HandleFunc("/ready", h.readyHandler)
	mux.HandleFunc("/live", h.livenessHandler)

	h.server = &http.Server{
		Addr:    ":" + h.port,
		Handler: mux,
	}

	h.isReady = true

	go func() {
		h.logger.Info("Starting health check server", "port", h.port)
		if err := h.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			h.logger.Error(err, "Health check server failed")
		}
	}()

	// Graceful shutdown
	<-ctx.Done()
	h.Shutdown(context.Background())
}

func (h *HealthChecker) Shutdown(ctx context.Context) error {
	h.isReady = false
	return h.server.Shutdown(ctx)
}

func (h *HealthChecker) healthHandler(w http.ResponseWriter, r *http.Request) {
	ip, _ := utils.GetLocalIP()

	response := models.HealthResponse{
		Status:    "ok",
		Operator:  h.operatorName,
		Timestamp: time.Now(),
		IPAddress: ip,
		Version:   os.Getenv("VERSION"), // можно задать через env
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (h *HealthChecker) readyHandler(w http.ResponseWriter, r *http.Request) {
	if h.isReady {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	} else {
		w.WriteHeader(http.StatusServiceUnavailable)
	}
}

func (h *HealthChecker) livenessHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
