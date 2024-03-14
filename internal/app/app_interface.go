package app

import (
	"context"
	"movie-lib/internal/model"
)

type App interface {
	CreateMovie(ctx context.Context, userId uint64, movie model.Movie) (model.Movie, error)
	UpdateMovie(ctx context.Context, userId uint64, id uint64, upd model.UpdateMovie) (model.Movie, error)
	DeleteMovie(ctx context.Context, userId uint64, id uint64) error
	GetMovie(ctx context.Context, userId uint64, id uint64) (model.Movie, error)
	GetMovies(ctx context.Context, userId uint64, sortBy model.SortParam) ([]model.Movie, error)
	SearchMovies(ctx context.Context, userId uint64, pattern string) ([]model.Movie, error)

	CreateActor(ctx context.Context, userId uint64, actor model.Actor) (model.Actor, error)
	UpdateActor(ctx context.Context, userId uint64, id uint64, upd model.UpdateActor) (model.Actor, error)
	DeleteActor(ctx context.Context, userId uint64, id uint64) error
	GetActor(ctx context.Context, userId uint64, id uint64) (model.Actor, error)
	GetActors(ctx context.Context, userId uint64) ([]model.Actor, error)
}
