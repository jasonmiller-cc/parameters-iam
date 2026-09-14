// Command server is the entrypoint for the parameters-iam REST API service.
package main

import (
	"context"
	"flag"
	"net/http"
	"os"

	"github.com/jasonmiller-cc/parameters-core/pkg/auth"
	"github.com/jasonmiller-cc/parameters-core/pkg/health"
	corelog "github.com/jasonmiller-cc/parameters-core/pkg/log"
	"github.com/jasonmiller-cc/parameters-core/pkg/middleware"
	"github.com/jasonmiller-cc/parameters-core/pkg/server"
	"github.com/jasonmiller-cc/parameters-core/pkg/version"
	"github.com/jasonmiller-cc/parameters-iam/internal/api"
	"github.com/jasonmiller-cc/parameters-iam/internal/config"
	"github.com/jasonmiller-cc/parameters-iam/internal/service"
)

func main() {
	configPath := flag.String("config", "", "path to config.yaml (default: auto-detect)")
	flag.Parse()

	log := corelog.Service("parameters-iam")
	log.Info("starting", "version", version.String())

	cfg, err := config.Load(*configPath)
	if err != nil {
		log.Error("failed to load config", "error", err)
		os.Exit(1)
	}

	jwtCfg := &auth.JWTConfig{
		Secret: []byte(cfg.Auth.JWTSecret),
		Issuer: cfg.Auth.JWTIssuer,
	}

	iamSvc := service.NewIAMService(
		cfg.IAM.LDAPServiceURL,
		cfg.IAM.KerberosServiceURL,
		cfg.IAM.CAServiceURL,
		cfg.IAM.DefaultRoles,
	)

	mux := http.NewServeMux()

	// IAM API routes.
	h := api.New(iamSvc, jwtCfg)
	h.Register(mux)

	// Health endpoints (no auth required).
	hc := health.New("parameters-iam", version.Version)
	mux.HandleFunc("GET /healthz", health.LiveHandler())
	mux.HandleFunc("GET /readyz", hc.Handler())

	// Version endpoint.
	mux.HandleFunc("GET /version", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		info := version.Get()
		_, _ = w.Write([]byte(`{"version":"` + info.Version +
			`","commit":"` + info.Commit +
			`","build_time":"` + info.BuildTime + `"}`))
	})

	chain := middleware.Chain(
		middleware.Recover(log),
		middleware.RequestID,
		middleware.Logger(log),
		middleware.CORS("*"),
	)

	srv := server.New(cfg.Server, chain(mux), log)
	if err := srv.Run(context.Background()); err != nil {
		log.Error("server exited with error", "error", err)
		os.Exit(1)
	}
}
