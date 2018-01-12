package rest

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"github.com/pressly/lg"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"

	"github.com/ld100/goblet/pkg/domain/common"
	user "github.com/ld100/goblet/pkg/domain/user/rest"
)

var (
	serverStarted = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_handler_started_requests_total",
			Help: "Count of started requests.",
		},
		[]string{"name", "handler", "host", "path", "method"},
	)
	serverCompleted = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_handler_completed_requests_total",
			Help: "Count of completed requests.",
		},
		[]string{"name", "handler", "host", "path", "method", "status"},
	)
	serverLatency = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_handler_completed_latency_seconds",
			Help:    "Latency of completed requests.",
			Buckets: []float64{.01, .03, .1, .3, 1, 3, 10, 30, 100, 300},
		},
		[]string{"name", "handler", "host", "path", "method", "status"},
	)
	serverRequestSize = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_handler_request_size_bytes",
			Help:    "Size of received requests.",
			Buckets: prometheus.ExponentialBuckets(32, 32, 6),
		},
		[]string{"name", "handler", "host", "path", "method", "status"},
	)
	serverResponseSize = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_handler_response_size_bytes",
			Help:    "Size of sent responses.",
			Buckets: prometheus.ExponentialBuckets(32, 32, 6),
		},
		[]string{"name", "handler", "host", "path", "method", "status"},
	)
)

func Serve() {

	// Setup the logger
	logger := logrus.New()
	logger.Formatter = &logrus.JSONFormatter{}

	lg.RedirectStdlogOutput(logger)
	lg.DefaultLogger = logger

	lg.Infoln("Welcome")

	serverCtx := context.Background()
	serverCtx = lg.WithLoggerContext(serverCtx, logger)
	lg.Log(serverCtx).Infof("Booting up server, %s", "v1.0")

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	//r.Use(middleware.Logger)
	//r.Use(NewStructuredLogger(logger))
	r.Use(lg.RequestLogger(logger))
	r.Use(middleware.Recoverer)
	r.Use(middleware.URLFormat)
	r.Use(render.SetContentType(render.ContentTypeJSON))

	// Prometheus instrumentation handler
	r.Handle("/metrics", promhttp.Handler())

	// Instrumented routes
	r.Group(func(r chi.Router) {
		// PrometheusInstrumentation
		prometheusEnabled, _ := strconv.ParseBool(os.Getenv("PROMETHEUS_ENABLED"))
		if prometheusEnabled {
			prometheus.MustRegister(serverStarted)
			prometheus.MustRegister(serverCompleted)
			prometheus.MustRegister(serverLatency)
			prometheus.MustRegister(serverRequestSize)
			prometheus.MustRegister(serverResponseSize)
			r.Use(PrometheusMiddleware)
		}

		r.Get("/", common.RootController)
		r.Get("/ping", common.PingController)
		r.Get("/panic", common.PanicController)

		// RESTy routes for "user" resource
		r.Mount("/user", user.UserRouter())

		// RESTy routes for "sessions" resource
		r.Mount("/sessions", user.SessionRouter())

		// r.Handle(GET, "/6", Chain(mwIncreaseCounter).HandlerFunc(handlerPrintCounter))
		// r.With(mwIncreaseCounter).Get("/6", handlerPrintCounter)
	})

	bindIP := ""
	bindPort := os.Getenv("HTTP_PORT")
	bindAddr := fmt.Sprintf("%v:%v", bindIP, bindPort)

	http.ListenAndServe(bindAddr, r)
}

// This is entirely optional, but I wanted to demonstrate how you could easily
// add your own logic to the render.Respond method.
func init() {
	render.Respond = func(w http.ResponseWriter, r *http.Request, v interface{}) {
		if err, ok := v.(error); ok {

			// We set a default error status response code if one hasn't been set.
			if _, ok := r.Context().Value(render.StatusCtxKey).(int); !ok {
				w.WriteHeader(400)
			}

			// We log the error
			fmt.Printf("Logging err: %s\n", err.Error())

			// We change the response to not reveal the actual error message,
			// instead we can transform the message something more friendly or mapped
			// to some code / language, etc.
			render.DefaultResponder(w, r, render.M{"status": "error"})
			return
		}

		render.DefaultResponder(w, r, v)
	}
}

// TODO: Plug actual prometheus calls here
func PrometheusMiddleware(next http.Handler) http.Handler {
	start := time.Now().UTC().UnixNano()
	name := os.Getenv("APP_NAME")

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// "name", "handler", "host", "path", "method", "status"
		requestSize := float64(r.ContentLength)
		host := r.Host
		path := r.RequestURI
		method := r.Method

		serverStarted.With(prometheus.Labels{
			"name": name,
			"handler": "",
			"host": host,
			"path": path,
			"method": method,
			}).Inc()
		serverRequestSize.With(prometheus.Labels{
			"name": name,
			"handler": "",
			"host": host,
			"path": path,
			"method": method,
			"status": "",
		}).Observe(requestSize)

		ctx := context.WithValue(r.Context(), "prometheus", "1")
		next.ServeHTTP(w, r.WithContext(ctx))

		//	TODO: Log response body size and valid statuses
		// Response time in miliseconds
		var respTime float64
		respTime = float64((time.Now().UTC().UnixNano() - start) / int64(time.Millisecond))
		serverCompleted.With(prometheus.Labels{
			"name": name,
			"handler": "",
			"host": host,
			"path": path,
			"method": method,
			"status": "",
		}).Inc()
		serverLatency.With(prometheus.Labels{
			"name": name,
			"handler": "",
			"host": host,
			"path": path,
			"method": method,
			"status": "",
		}).Observe(respTime)
	})
}
