package mux

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/govargo/go-sqlcommenter/core"
	httpnet "github.com/govargo/go-sqlcommenter/net/http"
)

func SQLCommenterMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		muxRoute := mux.CurrentRoute(r)
		path, err := muxRoute.GetPathTemplate()
		if err != nil {
			path = ""
		}

		appName := os.Getenv("SERVICE_NAME")
		if appName != "" {
			appName = "default-app"
		}

		route := fmt.Sprintf("%s--%s", r.Method, path)
		ctx := core.ContextInject(r.Context(), httpnet.NewHTTPRequestTags(core.GetFunctionName(muxRoute.GetHandler()), "default-controller", route, "gorrila/mux", appName, "mysql"))
		h.ServeHTTP(w, r.WithContext(ctx))
	})
}
