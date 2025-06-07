package http

import (
	"net/http"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/gorilla/mux"
	"github.com/newrelic/go-agent/v3/newrelic"
)

// NewRelicMiddleware wraps each HTTP handler with New Relic transaction monitoring
func NewRelicMiddleware(app *newrelic.Application, logger log.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if app == nil {
				next.ServeHTTP(w, r)
				return
			}

			var routePattern string
			currentRoute := mux.CurrentRoute(r)
			if currentRoute != nil {
				if pattern, err := currentRoute.GetPathTemplate(); err == nil {
					routePattern = pattern
				}
			}

			if routePattern == "" {
				routePattern = r.URL.Path
			}

			txn := app.StartTransaction(routePattern)
			defer txn.End()

			r = newrelic.RequestWithTransactionContext(r, txn)

			w = txn.SetWebResponse(w)

			next.ServeHTTP(w, r)
		})
	}
}

func InitNewRelic(appName, licenseKey string, logger log.Logger) *newrelic.Application {
	if licenseKey == "" {
		level.Warn(logger).Log("msg", "New Relic license key is empty, skipping New Relic initialization")
		return nil
	}

	level.Info(logger).Log("msg", "Initializing New Relic", "app", appName, "license", licenseKey)

	app, err := newrelic.NewApplication(
		newrelic.ConfigAppName(appName),
		newrelic.ConfigLicense(licenseKey),
		newrelic.ConfigDistributedTracerEnabled(true),
	)

	if err != nil {
		level.Error(logger).Log("msg", "Failed to initialize New Relic", "err", err)
		return nil
	}

	level.Info(logger).Log("msg", "New Relic initialized successfully", "app", appName)
	return app
}
