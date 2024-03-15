package app

import (
	"context"
	"movie-lib/internal/model"
	"movie-lib/internal/repo"
	"movie-lib/pkg/logger"
)

type appImpl struct {
	r    repo.Repo
	logs logger.Logger
}

func (a *appImpl) CreateMovie(ctx context.Context, userId uint64, movie model.Movie) (model.Movie, error) {
	var err error
	defer func() {
		if err != nil {
			a.logs.ErrorLog(err.Error())
		}
	}()

	var role model.Role
	if role, err = a.r.GetUserRole(ctx, userId); err != nil {
		return model.Movie{}, err
	} else if role != model.Admin {
		return model.Movie{}, model.ErrPermissionDenied
	}

	if !(len([]rune(movie.Title)) <= 150 && len([]rune(movie.Title)) >= 1) ||
		len([]rune(movie.Description)) > 1000 ||
		!(movie.Rating >= 0. && movie.Rating <= 10.) {
		return model.Movie{}, model.ErrValidationError
	}

	movie, err = a.r.CreateMovie(ctx, movie)
	return movie, err
}

func (a *appImpl) UpdateMovie(ctx context.Context, userId uint64, id uint64, upd model.UpdateMovie) (model.Movie, error) {
	var err error
	defer func() {
		if err != nil {
			a.logs.ErrorLog(err.Error())
		}
	}()

	var role model.Role
	if role, err = a.r.GetUserRole(ctx, userId); err != nil {
		return model.Movie{}, err
	} else if role != model.Admin {
		return model.Movie{}, model.ErrPermissionDenied
	}

	if !(len([]rune(upd.Title)) <= 150 && len([]rune(upd.Title)) >= 1) ||
		len([]rune(upd.Description)) > 1000 ||
		!(upd.Rating >= 0. && upd.Rating <= 10.) {
		return model.Movie{}, model.ErrValidationError
	}

	var movie model.Movie
	movie, err = a.r.UpdateMovie(ctx, id, upd)
	return movie, err
}

func (a *appImpl) DeleteMovie(ctx context.Context, userId uint64, id uint64) error {
	var err error
	defer func() {
		if err != nil {
			a.logs.ErrorLog(err.Error())
		}
	}()

	var role model.Role
	if role, err = a.r.GetUserRole(ctx, userId); err != nil {
		return err
	} else if role != model.Admin {
		return model.ErrPermissionDenied
	}

	err = a.r.DeleteMovie(ctx, id)
	return err
}

func (a *appImpl) GetMovie(ctx context.Context, userId uint64, id uint64) (model.Movie, error) {
	var err error
	defer func() {
		if err != nil {
			a.logs.ErrorLog(err.Error())
		}
	}()

	if _, err = a.r.GetUserRole(ctx, userId); err != nil {
		return model.Movie{}, err
	}

	var movie model.Movie
	movie, err = a.r.GetMovie(ctx, id)
	return movie, err
}

func (a *appImpl) GetMovies(ctx context.Context, userId uint64, sortBy model.SortParam) ([]model.Movie, error) {
	var err error
	defer func() {
		if err != nil {
			a.logs.ErrorLog(err.Error())
		}
	}()

	if _, err = a.r.GetUserRole(ctx, userId); err != nil {
		return []model.Movie{}, err
	}

	var movies []model.Movie
	movies, err = a.r.GetMovies(ctx, sortBy)
	return movies, err
}

func (a *appImpl) SearchMovies(ctx context.Context, userId uint64, pattern string) ([]model.Movie, error) {
	var err error
	defer func() {
		if err != nil {
			a.logs.ErrorLog(err.Error())
		}
	}()

	if _, err = a.r.GetUserRole(ctx, userId); err != nil {
		return []model.Movie{}, err
	}

	var movies []model.Movie
	movies, err = a.r.SearchMovies(ctx, pattern)
	return movies, err
}

func (a *appImpl) CreateActor(ctx context.Context, userId uint64, actor model.Actor) (model.Actor, error) {
	var err error
	defer func() {
		if err != nil {
			a.logs.ErrorLog(err.Error())
		}
	}()

	var role model.Role
	if role, err = a.r.GetUserRole(ctx, userId); err != nil {
		return model.Actor{}, err
	} else if role != model.Admin {
		return model.Actor{}, model.ErrPermissionDenied
	}

	actor, err = a.r.CreateActor(ctx, actor)
	return actor, err
}

func (a *appImpl) UpdateActor(ctx context.Context, userId uint64, id uint64, upd model.UpdateActor) (model.Actor, error) {
	var err error
	defer func() {
		if err != nil {
			a.logs.ErrorLog(err.Error())
		}
	}()

	var role model.Role
	if role, err = a.r.GetUserRole(ctx, userId); err != nil {
		return model.Actor{}, err
	} else if role != model.Admin {
		return model.Actor{}, model.ErrPermissionDenied
	}

	var actor model.Actor
	actor, err = a.r.UpdateActor(ctx, id, upd)
	return actor, err
}

func (a *appImpl) DeleteActor(ctx context.Context, userId uint64, id uint64) error {
	var err error
	defer func() {
		if err != nil {
			a.logs.ErrorLog(err.Error())
		}
	}()

	var role model.Role
	if role, err = a.r.GetUserRole(ctx, userId); err != nil {
		return err
	} else if role != model.Admin {
		return model.ErrPermissionDenied
	}

	err = a.r.DeleteActor(ctx, id)
	return err
}

func (a *appImpl) GetActor(ctx context.Context, userId uint64, id uint64) (model.Actor, error) {
	var err error
	defer func() {
		if err != nil {
			a.logs.ErrorLog(err.Error())
		}
	}()

	if _, err = a.r.GetUserRole(ctx, userId); err != nil {
		return model.Actor{}, err
	}

	var actor model.Actor
	actor, err = a.r.GetActor(ctx, id)
	return actor, err
}

func (a *appImpl) GetActors(ctx context.Context, userId uint64) ([]model.Actor, error) {
	var err error
	defer func() {
		if err != nil {
			a.logs.ErrorLog(err.Error())
		}
	}()

	if _, err = a.r.GetUserRole(ctx, userId); err != nil {
		return []model.Actor{}, err
	}

	var actors []model.Actor
	actors, err = a.r.GetActors(ctx)
	return actors, err
}
