package rest

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"github.com/pressly/lg"
	"github.com/sirupsen/logrus"
	// "github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/ld100/goblet/pkg/domain/common"
	user "github.com/ld100/goblet/pkg/domain/user/rest"
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

	r.Get("/", common.RootController)
	r.Get("/ping", common.PingController)
	r.Get("/panic", common.PanicController)

	// RESTy routes for "user" resource
	r.Mount("/user", user.UserRouter())

	// RESTy routes for "sessions" resource
	r.Mount("/sessions", user.SessionRouter())

	bindIP := ""
	bindPort := os.Getenv("HTTP_PORT")
	bindAddr := fmt.Sprintf("%v:%v", bindIP, bindPort)

	// Prometheus instrumentation handler
	http.Handle("/metrics", promhttp.Handler())

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
