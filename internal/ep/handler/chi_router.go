package handler

import (
	"net/http"
	"time"

	"github.com/ElfAstAhe/url-shortener2/internal/app/config"
	auditservice "github.com/ElfAstAhe/url-shortener2/internal/bll/service/audit"
	"github.com/ElfAstAhe/url-shortener2/internal/ep/facade"
	"github.com/ElfAstAhe/url-shortener2/internal/ep/middleware/audit"
	"github.com/ElfAstAhe/url-shortener2/internal/ep/middleware/compress"
	"github.com/ElfAstAhe/url-shortener2/internal/ep/middleware/iter14"
	mwarelog "github.com/ElfAstAhe/url-shortener2/internal/ep/middleware/logger"
	"github.com/ElfAstAhe/url-shortener2/pkg/logger"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type AppChiRouter struct {
	router        *chi.Mux
	log           logger.Logger
	conf          *config.Config
	toolFacade    facade.ToolFacade
	authFacade    facade.AuthFacade
	shortenFacade facade.ShortenFacade
	userFacade    facade.UserFacade
}

func NewAppChiRouter(toolFacade facade.ToolFacade, authFacade facade.AuthFacade, shortenFacade facade.ShortenFacade, userFacade facade.UserFacade, observers []auditservice.IncomeObserver, conf *config.Config, logger logger.Logger) *AppChiRouter {
	res := &AppChiRouter{
		router:        chi.NewRouter(),
		log:           logger.GetLogger("app router"),
		conf:          conf,
		toolFacade:    toolFacade,
		authFacade:    authFacade,
		shortenFacade: shortenFacade,
		userFacade:    userFacade,
	}

	res.setupMiddleware(observers, logger)
	res.setupRoutes()

	return res
}

func (cr *AppChiRouter) GetRouter() http.Handler {
	return cr.router
}

func (cr *AppChiRouter) setupMiddleware(observers []auditservice.IncomeObserver, logger logger.Logger) {
	// dev income request audit
	cr.router.Use(audit.NewDevIncomeMiddleware(logger, true).Handle)
	// jwt auth iter14
	cr.router.Use(iter14.NewJWTAuthIter14(nil, logger).Iter14Auth)
	// jwt auth
	// ..
	// requestID
	cr.router.Use(middleware.RequestID)
	// realIP
	cr.router.Use(middleware.RealIP)
	// compress
	cr.router.Use(compress.CustomCompress(compress.DefaultCompressionLevel, compress.ContentTypeApplicationJSON, compress.ContentTypeTextHTML))
	// decompress
	cr.router.Use(compress.CustomDecompress)
	// income/outcome logger
	cr.router.Use(mwarelog.NewHTTPReqRespLogger(logger).LogReqRes)
	// income audit
	cr.router.Use(audit.NewIncomeAuditMiddleware([]*audit.IncomeAuditPath{
		audit.NewIncomeAuditPath(http.MethodGet, "/"),
		audit.NewIncomeAuditPath(http.MethodPost, "/"),
		audit.NewIncomeAuditPath(http.MethodPost, "/api/shorten")},
		logger,
		observers...).Audit)
	// recoverer
	cr.router.Use(middleware.Recoverer)
	// timeout
	cr.router.Use(middleware.Timeout(20 * time.Second))
}

func (cr *AppChiRouter) setupRoutes() {
	// root
	cr.router.Route("/", func(r chi.Router) {
		r.Get("/{key}", cr.getRoot) // GET /{key}
		r.Post("/", cr.postRoot)    // POST /

		// ping sub-router
		r.Route("/ping", func(r chi.Router) {
			r.Get("/", cr.getPing) // GET /ping
		})

		// api sub-router
		r.Route("/api", func(r chi.Router) {
			r.Route("/shorten", func(r chi.Router) {
				r.Post("/", cr.postAPIShorten)           // POST /api/shorten
				r.Post("/batch", cr.postAPIShortenBatch) // POST /api/shorten/batch
			})

			// user sub-router
			r.Route("/user", func(r chi.Router) {
				// url sub-router
				r.Route("/urls", func(r chi.Router) {
					r.Get("/", cr.getAPIUserUrls)
					r.Delete("/", cr.deleteAPIUserUrls)
				})
			})
		})
	})
}
