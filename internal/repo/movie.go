package repo

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"movie-lib/internal/model"
)

const (
	createMovieQuery = `
		INSERT INTO "movies" ("title", "description", "release_date", "rating")
		VALUES ($1, $2, $3, $4)
		RETURNING "id";`

	addActorToMovieQuery = `
		INSERT INTO "movie-actor" ("movie-id", actor_id) 
		VALUES ($1, $2);`

	updateMovieQuery = `
		UPDATE "movies"
		SET "title" = $2,
		    "description" = $3,
		    "release_date" = $4,
		    "rating" = $5
		WHERE "id" = $1;`

	deleteMovieFromActorsQuery = `
		DELETE FROM "movie-actor"
		WHERE "movie-id" = $1;`

	deleteMovieQuery = `
		DELETE FROM "movies"
		WHERE "id" = $1;`

	getMovieQuery = `
		SELECT * FROM "movies"
		WHERE "id" = $1;`

	getMoviesSortByDefaultQuery = `
		SELECT * FROM "movies"
		ORDER BY "rating" DESC;`

	getMoviesSortByTitleQuery = `
		SELECT * FROM "movies"
		ORDER BY "title";`

	getMoviesSortByRatingQuery = `
		SELECT * FROM "movies"
		ORDER BY "rating";`

	getMoviesSortByReleaseDateQuery = `
		SELECT * FROM "movies"
		ORDER BY "release_date";`

	getMoviesByPatternQuery = `
		SELECT "movies"."id", "movies"."title", "movies"."description", "movies"."release_date", "movies"."rating" FROM "movie-actor"
		INNER JOIN "movies" ON "movies"."id" = "movie-actor"."movie-id"
		INNER JOIN "actors" ON "actors"."id" = "movie-actor"."actor_id"
		WHERE "movies"."title" LIKE $1 OR
			  "actors"."first_name" LIKE $1 OR
			  "actors"."second_name" LIKE $1
		GROUP BY "movies"."id";`

	getMovieActorsQuery = `
		SELECT "actors"."id", "actors"."first_name", "actors"."second_name", "actors"."gender"
		FROM "movie-actor"
			INNER JOIN "actors" ON "movie-actor"."actor_id" = "actors"."id"
		WHERE "movie-actor"."movie-id" = $1;`
)

func (r *repoImpl) CreateMovie(ctx context.Context, movie model.Movie) (model.Movie, error) {
	for _, id := range movie.ActorsId {
		if _, err := r.GetActor(ctx, id); err != nil {
			return model.Movie{}, err
		}
	}

	if err := r.QueryRow(ctx, createMovieQuery,
		movie.Title,
		movie.Description,
		movie.ReleaseDate,
		movie.Rating,
	).Scan(&movie.Id); err != nil {
		return model.Movie{}, errors.Join(model.ErrDatabaseError, err)
	}

	for _, id := range movie.ActorsId {
		_, _ = r.Exec(ctx, addActorToMovieQuery, movie.Id, id)
	}

	return r.GetMovie(ctx, movie.Id)
}

func (r *repoImpl) UpdateMovie(ctx context.Context, id uint64, upd model.UpdateMovie) (model.Movie, error) {
	for _, actorId := range upd.Actors {
		if _, err := r.GetActor(ctx, actorId); err != nil {
			return model.Movie{}, err
		}
	}

	if e, err := r.Exec(ctx, updateMovieQuery,
		id,
		upd.Title,
		upd.Description,
		upd.ReleaseDate,
		upd.Rating,
	); err != nil {
		return model.Movie{}, errors.Join(model.ErrDatabaseError, err)
	} else if e.RowsAffected() == 0 {
		return model.Movie{}, model.ErrMovieNotExists
	}

	_, _ = r.Exec(ctx, deleteMovieFromActorsQuery, id)

	for _, actorId := range upd.Actors {
		_, _ = r.Exec(ctx, addActorToMovieQuery, id, actorId)
	}

	return r.GetMovie(ctx, id)
}

func (r *repoImpl) DeleteMovie(ctx context.Context, id uint64) error {
	if e, err := r.Exec(ctx, deleteMovieQuery, id); err != nil {
		return errors.Join(model.ErrDatabaseError, err)
	} else if e.RowsAffected() == 0 {
		return model.ErrMovieNotExists
	}

	if _, err := r.Exec(ctx, deleteMovieFromActorsQuery, id); err != nil {
		return errors.Join(model.ErrDatabaseError, err)
	}
	return nil
}

func (r *repoImpl) GetMovie(ctx context.Context, id uint64) (model.Movie, error) {
	row := r.QueryRow(ctx, getMovieQuery, id)
	var movie model.Movie
	if err := row.Scan(
		&movie.Id,
		&movie.Title,
		&movie.Description,
		&movie.ReleaseDate,
		&movie.Rating,
	); errors.Is(err, pgx.ErrNoRows) {
		return model.Movie{}, model.ErrMovieNotExists
	} else if err != nil {
		return model.Movie{}, errors.Join(model.ErrDatabaseError, err)
	}

	var err error
	movie.Actors, err = r.getMovieActors(ctx, id)
	if err != nil {
		return model.Movie{}, errors.Join(model.ErrDatabaseError, err)
	}
	return movie, nil
}

func (r *repoImpl) GetMovies(ctx context.Context, sortBy model.SortParam) ([]model.Movie, error) {
	var query string
	switch sortBy {
	case model.Title:
		query = getMoviesSortByTitleQuery
	case model.Rating:
		query = getMoviesSortByRatingQuery
	case model.ReleaseDate:
		query = getMoviesSortByReleaseDateQuery
	default:
		query = getMoviesSortByDefaultQuery
	}

	rows, err := r.Query(ctx, query)
	if err != nil {
		return []model.Movie{}, errors.Join(model.ErrDatabaseError, err)
	}
	defer rows.Close()

	movies := make([]model.Movie, 0)
	for rows.Next() {
		var movie model.Movie
		_ = rows.Scan(
			&movie.Id,
			&movie.Title,
			&movie.Description,
			&movie.ReleaseDate,
			&movie.Rating,
		)
		movies = append(movies, movie)
	}
	for i := range movies {
		movies[i].Actors, _ = r.getMovieActors(ctx, movies[i].Id)
	}
	return movies, nil
}

func (r *repoImpl) SearchMovies(ctx context.Context, pattern string) ([]model.Movie, error) {
	rows, err := r.Query(ctx, getMoviesByPatternQuery, "%"+pattern+"%")
	if err != nil {
		return []model.Movie{}, errors.Join(model.ErrDatabaseError, err)
	}
	defer rows.Close()

	movies := make([]model.Movie, 0)
	for rows.Next() {
		var movie model.Movie
		_ = rows.Scan(
			&movie.Id,
			&movie.Title,
			&movie.Description,
			&movie.ReleaseDate,
			&movie.Rating,
		)
		movies = append(movies, movie)
	}
	for i := range movies {
		movies[i].Actors, _ = r.getMovieActors(ctx, movies[i].Id)
	}
	return movies, nil
}

func (r *repoImpl) getMovieActors(ctx context.Context, id uint64) ([]model.Actor, error) {
	rows, err := r.Query(ctx, getMovieActorsQuery, id)
	if err != nil {
		return []model.Actor{}, errors.Join(model.ErrDatabaseError, err)
	}
	defer rows.Close()

	actors := make([]model.Actor, 0)
	for rows.Next() {
		var actor model.Actor
		_ = rows.Scan(
			&actor.Id,
			&actor.FirstName,
			&actor.SecondName,
			&actor.Gender,
		)
		actors = append(actors, actor)
	}
	return actors, nil
}
