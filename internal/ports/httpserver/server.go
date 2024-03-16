package httpserver

import (
	"context"
	"fmt"
	"github.com/swaggo/http-swagger"
	_ "movie-lib/docs"
	"movie-lib/internal/app"
	"movie-lib/pkg/logger"
	"net/http"
)

func handleMovies(ctx context.Context, a app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			createMovieHandler(ctx, a)(w, r)
		case http.MethodPut:
			updateMovieHandler(ctx, a)(w, r)
		case http.MethodDelete:
			deleteMovieHandler(ctx, a)(w, r)
		case http.MethodGet:
			getMovieHandler(ctx, a)(w, r)
		}
	}
}

func handleActors(ctx context.Context, a app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			createActorHandler(ctx, a)(w, r)
		case http.MethodPut:
			updateActorHandler(ctx, a)(w, r)
		case http.MethodDelete:
			deleteActorHandler(ctx, a)(w, r)
		case http.MethodGet:
			getActorHandler(ctx, a)(w, r)
		}
	}
}

func New(ctx context.Context, host string, port int, a app.App, logs logger.Logger) *http.Server {
	mux := http.NewServeMux()

	mux.Handle("/swagger/", httpSwagger.Handler(httpSwagger.URL(fmt.Sprintf("http://%s:%d/swagger/doc.json", "localhost", port))))

	mux.Handle("/api/v1/actors/", logMiddleware(handleActors(ctx, a), logs))
	mux.Handle("/api/v1/actors/list/", logMiddleware(getActorsListHandler(ctx, a), logs))
	mux.Handle("/api/v1/movies/", logMiddleware(handleMovies(ctx, a), logs))
	mux.Handle("/api/v1/movies/list/", logMiddleware(getMovieListHandler(ctx, a), logs))

	return &http.Server{
		Addr:    fmt.Sprintf("%s:%d", host, port),
		Handler: mux,
	}
}
