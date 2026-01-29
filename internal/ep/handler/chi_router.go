package handler

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/ElfAstAhe/url-shortener2/internal/app/config"
	"github.com/ElfAstAhe/url-shortener2/internal/bll/service/audit"
	"github.com/ElfAstAhe/url-shortener2/internal/ep/facade"
	mwareapp "github.com/ElfAstAhe/url-shortener2/internal/ep/middleware"
	mwareaudit "github.com/ElfAstAhe/url-shortener2/internal/ep/middleware/audit"
	"github.com/ElfAstAhe/url-shortener2/internal/ep/middleware/auth"
	"github.com/ElfAstAhe/url-shortener2/internal/ep/middleware/compress"
	"github.com/ElfAstAhe/url-shortener2/internal/ep/middleware/iter14"
	mwarelog "github.com/ElfAstAhe/url-shortener2/internal/ep/middleware/logger"
	"github.com/ElfAstAhe/url-shortener2/pkg/logger"
)

type AppChiRouter struct {
	router        *chi.Mux
	log           logger.Logger
	conf          *config.AppConf
	toolFacade    facade.ToolFacade
	shortenFacade facade.ShortenFacade
	userFacade    facade.UserFacade
}

func NewAppChiRouter(toolFacade facade.ToolFacade, shortenFacade facade.ShortenFacade, userFacade facade.UserFacade, auditIncome audit.IncomePublisher, conf *config.AppConf, logger logger.Logger) *AppChiRouter {
	// new router
	res := &AppChiRouter{
		router:        chi.NewRouter(),
		log:           logger.GetLogger("app router"),
		conf:          conf,
		toolFacade:    toolFacade,
		shortenFacade: shortenFacade,
		userFacade:    userFacade,
	}

	// setup middleware
	res.setupMiddleware(auditIncome, logger)
	// mount
	res.router.Mount("/debug", middleware.Profiler())
	// routes
	res.setupRoutes()

	return res
}

func (cr *AppChiRouter) GetRouter() http.Handler {
	return cr.router
}

func (cr *AppChiRouter) setupMiddleware(auditIncome audit.IncomePublisher, logger logger.Logger) {
	// dev income request audit
	//	cr.router.Use(audit.NewDevIncomeMiddleware(logger, true).Handle)
	// jwt auth iter14
	cr.router.Use(iter14.NewJWTAuthIter14(iter14.NewAuthIter14PathMatchers([]*iter14.AuthIter14PathMatcher{
		// POST /api/shorten/batch
		iter14.NewAuthIter14PathMatcher(http.MethodPost, "/api/shorten/batch", mwareapp.PatternPostAPIShortenBatch, http.StatusUnauthorized, http.StatusUnauthorized, http.StatusInternalServerError),
		// GET /api/user/urls
		iter14.NewAuthIter14PathMatcher(http.MethodGet, "/api/user/urls", mwareapp.PatternGetAPIUserUrls, http.StatusNoContent, http.StatusUnauthorized, http.StatusInternalServerError),
	}, logger), logger).Iter14Auth)
	// jwt auth retriever
	cr.router.Use(auth.NewJWTAuthRetriever(logger).AuthRetriever)
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
	cr.router.Use(mwareaudit.NewIncomeAuditMiddleware(mwareaudit.NewIncomeAuditPathMatchers([]*mwareaudit.IncomeAuditPathMatcher{
		mwareaudit.NewIncomeAuditPathMatcher(http.MethodGet, "/", mwareapp.PatternGetRoot, mwareaudit.OriginalURLExtractGetRoot),
		mwareaudit.NewIncomeAuditPathMatcher(http.MethodPost, "/", mwareapp.PatternPostRoot, mwareaudit.OriginalURLExtractPostRoot),
		mwareaudit.NewIncomeAuditPathMatcher(http.MethodPost, "/api/shorten", mwareapp.PatternPostAPIShorten, mwareaudit.OriginalURLExtractPostAPIShorten),
	}), auditIncome, logger).Audit)
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
