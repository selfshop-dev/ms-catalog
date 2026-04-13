package app

import (
	"fmt"

	"github.com/go-chi/chi/v5"

	config "github.com/selfshop-dev/lib-config"
	logger "github.com/selfshop-dev/lib-logger"

	"github.com/selfshop-dev/ms-catalog/docs"
	"github.com/selfshop-dev/ms-catalog/internal/db/dbstorage"
	"github.com/selfshop-dev/ms-catalog/internal/db/gen"
	"github.com/selfshop-dev/ms-catalog/internal/handler"
	v1handler "github.com/selfshop-dev/ms-catalog/internal/handler/v1"
	v1usecase "github.com/selfshop-dev/ms-catalog/internal/usecase/v1"
	"github.com/selfshop-dev/ms-catalog/pkg/health"
	"github.com/selfshop-dev/ms-catalog/pkg/server/middleware"
	"github.com/selfshop-dev/ms-catalog/pkg/swagger"
)

func NewRouter(
	q *gen.Queries, l *logger.Logger, c *config.Base, h *health.Collector,
) *chi.Mux {
	pr := dbstorage.NewProductAdapter(q)

	v1 := []handler.Route{
		v1handler.NewCreateProduct(v1usecase.NewCreateProduct(pr)),
		v1handler.NewDeleteProduct(v1usecase.NewDeleteProduct(pr)),
		v1handler.NewGetProductBySlug(v1usecase.NewGetProductBySlug(pr)),
		v1handler.NewGetProduct(v1usecase.NewGetProduct(pr)),
		v1handler.NewGetActiveProducts(v1usecase.NewGetActiveProducts(pr)),
		v1handler.NewUpdateProduct(v1usecase.NewUpdateProduct(pr)),
		v1handler.NewUpdateStatusProduct(v1usecase.NewUpdateStatusProduct(pr)),
	}

	r := chi.NewRouter()
	r.Use(
		middleware.RequestID(),
		middleware.Logging(l.Logger),
		middleware.Recover(l.Unsampled()),
		middleware.Timeout(c.Entry.HTTP.RequestTimeout),
	)

	r.Route("/health", func(r chi.Router) { r.Get("/alive", h.AliveHandler); r.Get("/ready", h.ReadyHandler) })
	r.Route("/api/v1", func(r chi.Router) {
		for _, e := range v1 {
			r.Method(e.Method, e.Path, e.Handler)
		}
	})

	r.Handle("/log/level", logger.LevelHandler(l.Level, l.Unsampled()))

	if c.IsDev() {
		docs.SwaggerInfo.Host = fmt.Sprintf("localhost:%d", c.Entry.HTTP.Port)
		r.Get("/swagger", swagger.ScalarUI)
		r.Get("/swagger/doc.json", swagger.Spec)
	}
	return r
}
