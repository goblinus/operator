package operator

import (
	"context"
	"os"
	"strconv"
	"time"

	"github.com/goblinus/operator/src/services/healthz"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

const (
	defaultHealthCheckPort = "8080"
	defaultCheckInterval   = 30 * time.Second
)

type MyOperator struct {
	mgr           manager.Manager
	healthChecker *healthz.HealthChecker
	logger        ctrl.Logger
}

func NewMyOperator() (*MyOperator, error) {
	scheme := runtime.NewScheme()
	// добавьте ваши схемы здесь

	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Scheme:                 scheme,
		Port:                   9443,
		HealthProbeBindAddress: ":8081", // для стандартных проверок k8s
	})
	if err != nil {
		return nil, err
	}

	logger := ctrl.Log.WithName("operator")

	operator := &MyOperator{
		mgr:    mgr,
		logger: logger,
	}

	return operator, nil
}

func (o *MyOperator) Start(ctx context.Context) error {
	// Инициализация health checker
	port := getEnv("HEALTH_CHECK_PORT", defaultHealthCheckPort)
	operatorName := getEnv("OPERATOR_NAME", "my-operator")

	o.healthChecker = healthz.NewHealthChecker(o.logger, operatorName, port)

	// Запуск health checker
	go o.healthChecker.Start(ctx)

	// Запуск активных проверок
	go o.startActiveHealthChecks(ctx)

	// Добавление стандартных проверок для k8s
	if err := o.mgr.AddHealthzCheck("healthz", healthz.Ping); err != nil {
		return err
	}
	if err := o.mgr.AddReadyzCheck("readyz", healthz.Ping); err != nil {
		return err
	}

	o.logger.Info("Starting operator")
	return o.mgr.Start(ctx)
}

func (o *MyOperator) startActiveHealthChecks(ctx context.Context) {
	intervalStr := getEnv("HEALTH_CHECK_INTERVAL", "30")
	intervalSec, _ := strconv.Atoi(intervalStr)
	interval := time.Duration(intervalSec) * time.Second

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			o.performHealthCheck()
		case <-ctx.Done():
			return
		}
	}
}

func (o *MyOperator) performHealthCheck() {
	// Здесь можно добавить логику проверки состояния оператора
	// Например, проверка подключения к Kubernetes API, базам данных и т.д.

	status := "healthy"
	// Добавьте ваши проверки здесь

	o.logger.Info("Health check performed", "status", status)
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
