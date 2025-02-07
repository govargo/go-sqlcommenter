package mux

import (
	"fmt"
	"net/http"

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

		route := fmt.Sprintf("%s--%s", r.Method, path)
		ctx := core.ContextInject(r.Context(), httpnet.NewHTTPRequestTags("gorrila/mux", route, core.GetFunctionName(muxRoute.GetHandler())))
		h.ServeHTTP(w, r.WithContext(ctx))
	})
}
