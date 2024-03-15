package repo

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"movie-lib/internal/model"
)

const (
	createActorQuery = `
		INSERT INTO "actors" ("first_name", "second_name", "gender") 
		VALUES ($1, $2, $3)
		RETURNING "id";`

	updateActorQuery = `
		UPDATE "actors"
		SET "first_name" = $2,
		    "second_name" = $3,
		    "gender" = $4
		WHERE "id" = $1;`

	getActorQuery = `
		SELECT * FROM "actors"
		WHERE "id" = $1;`

	getActorsQuery = `
		SELECT * FROM "actors";`

	getActorMoviesQuery = `
		SELECT "movies"."id", "movies"."title", "movies"."description", "movies"."release_date", "movies"."rating"
		FROM "movie-actor"
			INNER JOIN "movies" ON "movie-id" = "movies"."id"
		WHERE "movie-actor"."actor_id" = $1;`

	deleteActorQuery = `
		DELETE FROM "actors"
		WHERE "id" = $1;`

	deleteActorFromMoviesQuery = `
		DELETE FROM "movie-actor"
		WHERE "actor_id" = $1;`
)

func (r *repoImpl) CreateActor(ctx context.Context, actor model.Actor) (model.Actor, error) {
	if err := r.QueryRow(ctx, createActorQuery,
		actor.FirstName,
		actor.SecondName,
		actor.Gender,
	).Scan(&actor.Id); err != nil {
		return model.Actor{}, errors.Join(model.ErrDatabaseError, err)
	}
	return actor, nil
}

func (r *repoImpl) UpdateActor(ctx context.Context, id uint64, upd model.UpdateActor) (model.Actor, error) {
	if e, err := r.Exec(ctx, updateActorQuery,
		id,
		upd.FirstName,
		upd.SecondName,
		upd.Gender,
	); err != nil {
		return model.Actor{}, errors.Join(model.ErrDatabaseError, err)
	} else if e.RowsAffected() == 0 {
		return model.Actor{}, model.ErrActorNotExists
	}
	return r.GetActor(ctx, id)
}

func (r *repoImpl) DeleteActor(ctx context.Context, id uint64) error {
	if e, err := r.Exec(ctx, deleteActorQuery, id); err != nil {
		return errors.Join(model.ErrDatabaseError, err)
	} else if e.RowsAffected() == 0 {
		return model.ErrActorNotExists
	}

	if _, err := r.Exec(ctx, deleteActorFromMoviesQuery, id); err != nil {
		return errors.Join(model.ErrDatabaseError, err)
	}
	return nil
}

func (r *repoImpl) GetActor(ctx context.Context, id uint64) (model.Actor, error) {
	row := r.QueryRow(ctx, getActorQuery, id)
	var actor model.Actor
	if err := row.Scan(
		&actor.Id,
		&actor.FirstName,
		&actor.SecondName,
		&actor.Gender,
	); errors.Is(err, pgx.ErrNoRows) {
		return model.Actor{}, model.ErrActorNotExists
	} else if err != nil {
		return model.Actor{}, errors.Join(model.ErrDatabaseError, err)
	}

	var err error
	actor.Movies, err = r.getActorMovies(ctx, id)
	if err != nil {
		return model.Actor{}, err
	}
	return actor, nil
}

func (r *repoImpl) GetActors(ctx context.Context) ([]model.Actor, error) {
	rows, err := r.Query(ctx, getActorsQuery)
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
	for i := range actors {
		actors[i].Movies, _ = r.getActorMovies(ctx, actors[i].Id)
	}
	return actors, nil
}

func (r *repoImpl) getActorMovies(ctx context.Context, id uint64) ([]model.Movie, error) {
	rows, err := r.Query(ctx, getActorMoviesQuery, id)
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
	return movies, nil
}
