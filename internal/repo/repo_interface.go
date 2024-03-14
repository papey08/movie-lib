package repo

import (
	"context"
	"movie-lib/internal/model"
)

type Repo interface {
	CreateMovie(ctx context.Context, movie model.Movie) (model.Movie, error)
	UpdateMovie(ctx context.Context, id uint64, upd model.UpdateMovie) (model.Movie, error)
	DeleteMovie(ctx context.Context, id uint64) error
	GetMovie(ctx context.Context, id uint64) (model.Movie, error)
	GetMovies(ctx context.Context, sortBy model.SortParam) ([]model.Movie, error)
	SearchMovies(ctx context.Context, pattern string) ([]model.Movie, error)

	CreateActor(ctx context.Context, actor model.Actor) (model.Actor, error)
	UpdateActor(ctx context.Context, id uint64, upd model.UpdateActor) (model.Actor, error)
	DeleteActor(ctx context.Context, id uint64) error
	GetActor(ctx context.Context, id uint64) (model.Actor, error)
	GetActors(ctx context.Context) ([]model.Actor, error)

	GetUserRole(ctx context.Context, id uint64) (model.Role, error)
}
