package main

import (
	"context"
	"errors"
	"io/fs"
	"log/slog"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/pgxstore"
	"github.com/alexedwards/scs/v2"

	"github.com/jared-wallace/website-go/internal/config"
	"github.com/jared-wallace/website-go/internal/database"
	adminhandler "github.com/jared-wallace/website-go/internal/handler/admin"
	apihandler "github.com/jared-wallace/website-go/internal/handler/api"
	bloghandler "github.com/jared-wallace/website-go/internal/handler/blog"
	"github.com/jared-wallace/website-go/internal/markdown"
	"github.com/jared-wallace/website-go/internal/middleware"
	postrepo "github.com/jared-wallace/website-go/internal/repository/post"
	postservice "github.com/jared-wallace/website-go/internal/service/post"
	"github.com/jared-wallace/website-go/internal/server"
	"github.com/jared-wallace/website-go/web"
)

// hostRouter dispatches HTTP requests to the admin or blog handler based on the
// Host header. Requests matching cfg.AdminHost go to admin; all others go to blog.
type hostRouter struct {
	blog      http.Handler
	admin     http.Handler
	adminHost string
}

func (hr *hostRouter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	host := r.Host
	if h, _, err := net.SplitHostPort(host); err == nil {
		host = h
	}
	if host == hr.adminHost {
		hr.admin.ServeHTTP(w, r)
		return
	}
	hr.blog.ServeHTTP(w, r)
}

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)

	cfg := config.Load()
	logger.Info("config loaded", "env", cfg.AppEnv, "port", cfg.Port)
	if cfg.AdminEmail == "" {
		logger.Info("admin credentials not configured; admin panel disabled")
	}

	// Ensure the image directory exists on first boot (Pitfall 3).
	if err := os.MkdirAll(cfg.ImageDir, 0755); err != nil {
		logger.Error("failed to create image directory", "path", cfg.ImageDir, "error", err)
		os.Exit(1)
	}
	logger.Info("image directory ready", "path", cfg.ImageDir)

	ctx, stop := server.GracefulShutdown()
	defer stop()

	pool, err := database.Connect(ctx, cfg.DatabaseURL)
	if err != nil {
		logger.Error("database connection failed", "error", err)
		os.Exit(1)
	}
	defer pool.Close()
	logger.Info("database connected")

	if err := database.RunMigrations(ctx, pool); err != nil {
		logger.Error("migrations failed", "error", err)
		os.Exit(1)
	}
	logger.Info("migrations applied")

	// Build dependency graph: renderer -> repo -> service -> handlers
	renderer := markdown.NewRenderer()
	repo := postrepo.New(pool)
	svc := postservice.New(repo, renderer)

	// --- Session Manager (SCS + Postgres-backed store) ---
	sessionManager := scs.New()
	sessionManager.Store = pgxstore.New(pool)
	sessionManager.IdleTimeout = 24 * time.Hour      // D-10: inactivity-based expiry
	sessionManager.Lifetime = 30 * 24 * time.Hour    // 30-day absolute backstop
	sessionManager.Cookie.HttpOnly = true
	sessionManager.Cookie.Secure = cfg.AppEnv == "production" // false in dev for http://localhost
	sessionManager.Cookie.SameSite = http.SameSiteLaxMode     // D-09, Pitfall 3
	sessionManager.Cookie.Name = "admin_session"

	// Rate limiter: 5 attempts per minute per IP (D-11)
	rl := middleware.NewRateLimiter(5, time.Minute)

	// --- Blog mux ---
	blog := bloghandler.New(svc)

	staticFS, err := fs.Sub(web.Static, "static")
	if err != nil {
		logger.Error("static fs sub failed", "error", err)
		os.Exit(1)
	}

	// Image file server — serves uploaded images from EBS volume on both muxes.
	imageServer := http.StripPrefix("/images/", http.FileServer(http.Dir(cfg.ImageDir)))

	blogMux := http.NewServeMux()
	blogMux.Handle("GET /static/", http.StripPrefix("/static/", http.FileServerFS(staticFS)))
	blogMux.Handle("GET /images/{path...}", imageServer)
	blogMux.HandleFunc("GET /{$}", blog.ListPosts)        // Home page (exact match)
	blogMux.HandleFunc("GET /posts", blog.ListPosts)       // /posts?page=N
	blogMux.HandleFunc("GET /posts/{slug}", blog.ShowPost)          // Single post
	blogMux.HandleFunc("POST /posts/{slug}/react", blog.React)      // Thumbs-up reaction
	blogMux.HandleFunc("GET /rss", blog.ServeRSS)
	blogMux.HandleFunc("GET /sitemap.xml", blog.ServeSitemap)
	blogMux.HandleFunc("GET /robots.txt", blog.ServeRobots)
	// --- API push endpoint (bearer token auth, on blog mux for public access) ---
	apiH := apihandler.New(svc, renderer)
	requireToken := middleware.RequireAPIToken(cfg.APIToken)
	blogMux.Handle("POST /api/push", requireToken(http.HandlerFunc(apiH.PushPost)))

	blogMux.HandleFunc("GET /{path...}", blog.NotFound) // Catch-all 404

	// --- Admin handler + mux ---
	adminH := adminhandler.New(svc, sessionManager, renderer, rl, cfg)

	adminMux := http.NewServeMux()
	adminMux.HandleFunc("GET /admin/login", adminH.LoginPage)
	adminMux.HandleFunc("POST /admin/login", adminH.LoginPost)
	adminMux.HandleFunc("POST /admin/logout", adminH.Logout)

	requireAuth := middleware.RequireSession(sessionManager)
	adminMux.Handle("GET /admin/posts", requireAuth(http.HandlerFunc(adminH.Dashboard)))
	adminMux.Handle("GET /admin/posts/new", requireAuth(http.HandlerFunc(adminH.NewPost)))
	adminMux.Handle("POST /admin/posts/new", requireAuth(http.HandlerFunc(adminH.SavePost)))
	adminMux.Handle("GET /admin/posts/{id}/edit", requireAuth(http.HandlerFunc(adminH.EditPost)))
	adminMux.Handle("POST /admin/posts/{id}/edit", requireAuth(http.HandlerFunc(adminH.SavePost)))
	adminMux.Handle("POST /admin/posts/{id}/delete", requireAuth(http.HandlerFunc(adminH.DeletePost)))
	adminMux.Handle("POST /admin/posts/{id}/restore", requireAuth(http.HandlerFunc(adminH.RestorePost)))
	adminMux.Handle("POST /admin/posts/{id}/publish", requireAuth(http.HandlerFunc(adminH.PublishPost)))
	adminMux.Handle("POST /admin/posts/{id}/unpublish", requireAuth(http.HandlerFunc(adminH.UnpublishPost)))
	adminMux.Handle("POST /admin/preview", requireAuth(http.HandlerFunc(adminH.Preview)))
	adminMux.Handle("POST /admin/images/upload", requireAuth(http.HandlerFunc(adminH.UploadImage)))
	adminMux.Handle("GET /images/{path...}", imageServer)
	adminMux.Handle("GET /static/", http.StripPrefix("/static/", http.FileServerFS(staticFS)))
	adminMux.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/admin/posts", http.StatusSeeOther)
	})

	// Wrap admin mux: CrossOriginProtection (CSRF) + session middleware
	cop := http.NewCrossOriginProtection()
	cop.AddTrustedOrigin("https://" + cfg.AdminHost)
	adminHandler := sessionManager.LoadAndSave(cop.Handler(adminMux))

	// --- Host router ---
	router := &hostRouter{
		blog:      blogMux,
		admin:     adminHandler,
		adminHost: cfg.AdminHost,
	}

	srv := server.New(cfg.Port, router)
	logger.Info("server starting", "port", cfg.Port)

	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("server error", "error", err)
			os.Exit(1)
		}
	}()

	<-ctx.Done()
	logger.Info("shutdown signal received")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Error("shutdown error", "error", err)
	}
	logger.Info("server stopped")
}
