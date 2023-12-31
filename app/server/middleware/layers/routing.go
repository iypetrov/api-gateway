package layers

import (
	"context"
	"github.com/IliyaYavorovPetrov/api-gateway/app/server/middleware"
	"github.com/IliyaYavorovPetrov/api-gateway/app/server/routing"
	"log"
	"net/http"
)

func Routing(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), middleware.ContextKey(middleware.IsAuthNeededKey), false)

		reqKey := routing.CreateRoutingCfgHashKey(r.Method, "http://"+r.Host+r.RequestURI)
		rri, err := routing.GetRoutingCfgFromRequestKey(ctx, reqKey)
		if err != nil {
			log.Printf("no such request in the routing configuration")
			return
		}

		proxyReq, err := http.NewRequest(r.Method, rri.DestinationURL, r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadGateway)
			return
		}
		proxyReq.Header = r.Header
		ctx = context.WithValue(r.Context(), middleware.ContextKey(middleware.ProxyRequest), proxyReq)

		if rri.IsAuthNeeded {
			ctx = context.WithValue(r.Context(), middleware.ContextKey(middleware.IsAuthNeededKey), true)
		}

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
