package main

import (
	"context"
	"embed"
	"errors"
	"io/fs"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ramadhantriyant/buku/internal/handlers"
	"github.com/ramadhantriyant/buku/internal/middlewares"
	"github.com/ramadhantriyant/buku/internal/models"
)

//go:embed all:ui/dist
var uiFiles embed.FS

func createServer(config *models.Config, port string) *http.Server {
	staticFS, err := fs.Sub(uiFiles, "ui/dist")
	if err != nil {
		log.Fatal("UI must be build first")
	}
	staticHandler := http.FileServer(http.FS(staticFS))

	mux := http.NewServeMux()
	h := handlers.New(config)

	// UI handler
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		staticHandler.ServeHTTP(w, r)
	})

	// Auth (public)
	mux.HandleFunc("POST /api/auth/register", h.Register)
	mux.HandleFunc("POST /api/auth/login", h.Login)
	mux.HandleFunc("POST /api/auth/refresh", h.RefreshToken)
	mux.HandleFunc("POST /api/auth/logout", h.Logout)

	// Password reset (public)
	mux.HandleFunc("POST /api/auth/forgot-password", h.RequestPasswordReset)
	mux.HandleFunc("POST /api/auth/reset-password", h.ResetPassword)

	// Protected routes
	authMiddleware := middlewares.AuthMiddleware(config.JWTSecret)

	// Profile
	mux.HandleFunc("GET /api/profile", authMiddleware(h.GetProfile))
	mux.HandleFunc("PUT /api/profile/password", authMiddleware(h.ChangePassword))

	// Categories (protected)
	mux.HandleFunc("GET /api/category", authMiddleware(h.ListCategory))
	mux.HandleFunc("GET /api/category/{id}", authMiddleware(h.GetCategory))
	mux.HandleFunc("POST /api/category", authMiddleware(h.CreateCategory))
	mux.HandleFunc("PUT /api/category/{id}", authMiddleware(h.UpdateCategory))
	mux.HandleFunc("DELETE /api/category/{id}", authMiddleware(h.DeleteCategory))

	// URLs (protected)
	mux.HandleFunc("GET /api/url", authMiddleware(h.ListURL))
	mux.HandleFunc("GET /api/url/{id}", authMiddleware(h.GetURL))
	mux.HandleFunc("POST /api/url", authMiddleware(h.CreateURL))
	mux.HandleFunc("PUT /api/url/{id}", authMiddleware(h.UpdateURL))
	mux.HandleFunc("DELETE /api/url/{id}", authMiddleware(h.DeleteURL))

	// Admin routes
	mux.HandleFunc("GET /api/admin/dashboard", authMiddleware(middlewares.AdminOnly(h.AdminDashboard)))

	handler := middlewares.Chain(mux, middlewares.CORS, middlewares.Logger, middlewares.ShouldJSON)
	return &http.Server{
		Addr:    port,
		Handler: handler,
	}
}

func runServer(ctx context.Context, server *http.Server, shutdownTimeout time.Duration) error {
	serverErr := make(chan error, 1)

	go func() {
		log.Println("Starting server...")
		if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			serverErr <- err
		}
		close(serverErr)
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-serverErr:
		return err
	case <-stop:
		log.Println("Shutting down...")
	case <-ctx.Done():
		log.Println("Context cancelled")
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()
	if err := server.Shutdown(shutdownCtx); err != nil {
		if closeErr := server.Close(); closeErr != nil {
			return errors.Join(err, closeErr)
		}
		return err
	}

	log.Println("Shutdown completed")

	return nil
}
