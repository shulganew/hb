package router

import (
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/go-chi/chi/v5"
	middleware "github.com/oapi-codegen/nethttp-middleware"
	"github.com/shulganew/hb.git/internal/config"
)

// Chi Router for application.
func Route(conf config.Config, swagger *openapi3.T) (r *chi.Mux) {
	r = chi.NewMux()
	// Use our validation middleware to check all requests against the
	// OpenAPI schema.
	r.Use(middleware.OapiRequestValidatorWithOptions(swagger, nil))
	return
}
