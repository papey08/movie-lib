package app

import (
	"context"
	"movie-lib/internal/model"
	"movie-lib/internal/repo"
)

type appImpl struct {
	r repo.Repo
}

func (a *appImpl) CreateMovie(ctx context.Context, userId uint64, movie model.Movie) (model.Movie, error) {
	if role, err := a.r.GetUserRole(ctx, userId); err != nil {
		return model.Movie{}, err
	} else if role != model.Admin {
		return model.Movie{}, model.ErrPermissionDenied
	}

	if !(len([]rune(movie.Title)) <= 150 && len([]rune(movie.Title)) >= 1) ||
		len([]rune(movie.Description)) > 1000 ||
		!(movie.Rating >= 0. && movie.Rating <= 10.) {
		return model.Movie{}, model.ErrValidationError
	}

	return a.r.CreateMovie(ctx, movie)
}

func (a *appImpl) UpdateMovie(ctx context.Context, userId uint64, id uint64, upd model.UpdateMovie) (model.Movie, error) {
	if role, err := a.r.GetUserRole(ctx, userId); err != nil {
		return model.Movie{}, err
	} else if role != model.Admin {
		return model.Movie{}, model.ErrPermissionDenied
	}

	if !(len([]rune(upd.Title)) <= 150 && len([]rune(upd.Title)) >= 1) ||
		len([]rune(upd.Description)) > 1000 ||
		!(upd.Rating >= 0. && upd.Rating <= 10.) {
		return model.Movie{}, model.ErrValidationError
	}

	return a.r.UpdateMovie(ctx, id, upd)
}

func (a *appImpl) DeleteMovie(ctx context.Context, userId uint64, id uint64) error {
	if role, err := a.r.GetUserRole(ctx, userId); err != nil {
		return err
	} else if role != model.Admin {
		return model.ErrPermissionDenied
	}

	return a.r.DeleteMovie(ctx, id)
}

func (a *appImpl) GetMovie(ctx context.Context, userId uint64, id uint64) (model.Movie, error) {
	if _, err := a.r.GetUserRole(ctx, userId); err != nil {
		return model.Movie{}, err
	}

	return a.r.GetMovie(ctx, id)
}

func (a *appImpl) GetMovies(ctx context.Context, userId uint64, sortBy model.SortParam) ([]model.Movie, error) {
	if _, err := a.r.GetUserRole(ctx, userId); err != nil {
		return []model.Movie{}, err
	}

	return a.r.GetMovies(ctx, sortBy)
}

func (a *appImpl) SearchMovies(ctx context.Context, userId uint64, pattern string) ([]model.Movie, error) {
	if _, err := a.r.GetUserRole(ctx, userId); err != nil {
		return []model.Movie{}, err
	}

	return a.r.SearchMovies(ctx, pattern)
}

func (a *appImpl) CreateActor(ctx context.Context, userId uint64, actor model.Actor) (model.Actor, error) {
	if role, err := a.r.GetUserRole(ctx, userId); err != nil {
		return model.Actor{}, err
	} else if role != model.Admin {
		return model.Actor{}, model.ErrPermissionDenied
	}

	return a.r.CreateActor(ctx, actor)
}

func (a *appImpl) UpdateActor(ctx context.Context, userId uint64, id uint64, upd model.UpdateActor) (model.Actor, error) {
	if role, err := a.r.GetUserRole(ctx, userId); err != nil {
		return model.Actor{}, err
	} else if role != model.Admin {
		return model.Actor{}, model.ErrPermissionDenied
	}

	return a.r.UpdateActor(ctx, id, upd)
}

func (a *appImpl) DeleteActor(ctx context.Context, userId uint64, id uint64) error {
	if role, err := a.r.GetUserRole(ctx, userId); err != nil {
		return err
	} else if role != model.Admin {
		return model.ErrPermissionDenied
	}

	return a.r.DeleteActor(ctx, id)
}

func (a *appImpl) GetActor(ctx context.Context, userId uint64, id uint64) (model.Actor, error) {
	if _, err := a.r.GetUserRole(ctx, userId); err != nil {
		return model.Actor{}, err
	}

	return a.r.GetActor(ctx, id)
}

func (a *appImpl) GetActors(ctx context.Context, userId uint64) ([]model.Actor, error) {
	if _, err := a.r.GetUserRole(ctx, userId); err != nil {
		return []model.Actor{}, err
	}

	return a.r.GetActors(ctx)
}
