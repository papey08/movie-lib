package httpserver

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"movie-lib/internal/app"
	"movie-lib/internal/model"
	"net/http"
	"strconv"
)

func createActorHandler(ctx context.Context, a app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userId, err := strconv.ParseUint(r.Header.Get("Authorization"), 10, 64)
		if err != nil {
			http.Error(w, errorResponse(model.ErrUnauthorized), http.StatusUnauthorized)
			return
		}
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, errorResponse(model.ErrInvalidInput), http.StatusBadRequest)
			return
		}
		var data createActorData
		if err = json.Unmarshal(body, &data); err != nil {
			http.Error(w, errorResponse(model.ErrInvalidInput), http.StatusBadRequest)
			return
		}

		actor, err := a.CreateActor(ctx, userId, model.Actor{
			FirstName:  data.FirstName,
			SecondName: data.SecondName,
			Gender:     data.Gender,
		})

		switch {
		case err == nil:
			w.WriteHeader(http.StatusOK)
			_, _ = fmt.Fprint(w, actorResponseOk(actor))
		case errors.Is(err, model.ErrPermissionDenied):
			http.Error(w, errorResponse(model.ErrPermissionDenied), http.StatusForbidden)
		case errors.Is(err, model.ErrUserNotExists):
			http.Error(w, errorResponse(model.ErrUserNotExists), http.StatusForbidden)
		case errors.Is(err, model.ErrDatabaseError):
			http.Error(w, errorResponse(model.ErrDatabaseError), http.StatusInternalServerError)
		default:
			http.Error(w, errorResponse(model.ErrServiceError), http.StatusInternalServerError)
		}
	}
}

func updateActorHandler(ctx context.Context, a app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userId, err := strconv.ParseUint(r.Header.Get("Authorization"), 10, 64)
		if err != nil {
			http.Error(w, errorResponse(model.ErrUnauthorized), http.StatusUnauthorized)
			return
		}
		actorId, err := strconv.ParseUint(r.URL.Query().Get("actor_id"), 10, 64)
		if err != nil {
			http.Error(w, errorResponse(model.ErrInvalidInput), http.StatusBadRequest)
			return
		}
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, errorResponse(model.ErrInvalidInput), http.StatusBadRequest)
			return
		}
		var data updateActorData
		if err = json.Unmarshal(body, &data); err != nil {
			http.Error(w, errorResponse(model.ErrInvalidInput), http.StatusBadRequest)
			return
		}

		actor, err := a.UpdateActor(ctx, userId, actorId, model.UpdateActor{
			FirstName:  data.FirstName,
			SecondName: data.SecondName,
			Gender:     data.Gender,
		})

		switch {
		case err == nil:
			w.WriteHeader(http.StatusOK)
			_, _ = fmt.Fprint(w, actorResponseOk(actor))
		case errors.Is(err, model.ErrActorNotExists):
			http.Error(w, errorResponse(model.ErrActorNotExists), http.StatusNotFound)
		case errors.Is(err, model.ErrPermissionDenied):
			http.Error(w, errorResponse(model.ErrPermissionDenied), http.StatusForbidden)
		case errors.Is(err, model.ErrUserNotExists):
			http.Error(w, errorResponse(model.ErrUserNotExists), http.StatusForbidden)
		case errors.Is(err, model.ErrDatabaseError):
			http.Error(w, errorResponse(model.ErrDatabaseError), http.StatusInternalServerError)
		default:
			http.Error(w, errorResponse(model.ErrServiceError), http.StatusInternalServerError)
		}
	}
}

func deleteActorHandler(ctx context.Context, a app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userId, err := strconv.ParseUint(r.Header.Get("Authorization"), 10, 64)
		if err != nil {
			http.Error(w, errorResponse(model.ErrUnauthorized), http.StatusUnauthorized)
			return
		}
		actorId, err := strconv.ParseUint(r.URL.Query().Get("actor_id"), 10, 64)
		if err != nil {
			http.Error(w, errorResponse(model.ErrInvalidInput), http.StatusBadRequest)
			return
		}

		err = a.DeleteActor(ctx, userId, actorId)

		switch {
		case err == nil:
			w.WriteHeader(http.StatusOK)
			_, _ = fmt.Fprint(w, errorResponse(nil))
		case errors.Is(err, model.ErrActorNotExists):
			http.Error(w, errorResponse(model.ErrActorNotExists), http.StatusNotFound)
		case errors.Is(err, model.ErrPermissionDenied):
			http.Error(w, errorResponse(model.ErrPermissionDenied), http.StatusForbidden)
		case errors.Is(err, model.ErrUserNotExists):
			http.Error(w, errorResponse(model.ErrUserNotExists), http.StatusForbidden)
		case errors.Is(err, model.ErrDatabaseError):
			http.Error(w, errorResponse(model.ErrDatabaseError), http.StatusInternalServerError)
		default:
			http.Error(w, errorResponse(model.ErrServiceError), http.StatusInternalServerError)
		}
	}
}

func getActorHandler(ctx context.Context, a app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userId, err := strconv.ParseUint(r.Header.Get("Authorization"), 10, 64)
		if err != nil {
			http.Error(w, errorResponse(model.ErrUnauthorized), http.StatusUnauthorized)
			return
		}

		if r.URL.Query().Get("actor_id") == "" {
			var actors []model.Actor
			actors, err = a.GetActors(ctx, userId)

			switch {
			case err == nil:
				w.WriteHeader(http.StatusOK)
				_, _ = fmt.Fprint(w, actorListResponseOk(actors))
			case errors.Is(err, model.ErrUserNotExists):
				http.Error(w, errorResponse(model.ErrUserNotExists), http.StatusForbidden)
			case errors.Is(err, model.ErrDatabaseError):
				http.Error(w, errorResponse(model.ErrDatabaseError), http.StatusInternalServerError)
			default:
				http.Error(w, errorResponse(model.ErrServiceError), http.StatusInternalServerError)
			}
		} else {
			var actorId uint64
			actorId, err = strconv.ParseUint(r.URL.Path[len("/actor/"):], 10, 64)
			if err != nil {
				http.Error(w, errorResponse(model.ErrInvalidInput), http.StatusBadRequest)
				return
			}

			var actor model.Actor
			actor, err = a.GetActor(ctx, userId, actorId)

			switch {
			case err == nil:
				w.WriteHeader(http.StatusOK)
				_, _ = fmt.Fprint(w, actorResponseOk(actor))
			case errors.Is(err, model.ErrActorNotExists):
				http.Error(w, errorResponse(model.ErrActorNotExists), http.StatusNotFound)
			case errors.Is(err, model.ErrUserNotExists):
				http.Error(w, errorResponse(model.ErrUserNotExists), http.StatusForbidden)
			case errors.Is(err, model.ErrDatabaseError):
				http.Error(w, errorResponse(model.ErrDatabaseError), http.StatusInternalServerError)
			default:
				http.Error(w, errorResponse(model.ErrServiceError), http.StatusInternalServerError)
			}
		}
	}
}
