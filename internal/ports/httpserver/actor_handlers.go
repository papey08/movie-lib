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

// @Summary		Добавление актёра
// @Description	Добавляет нового актёра
// @Tags			actors
// @Security		ApiKeyAuth
// @Accept			json
// @Produce		json
// @Param			input	body		createActorData	true	"Информация о новом актёре"
// @Success		200		{object}	actorResponse	"Информация об актёре"
// @Failure		400		{object}	actorResponse	"Неверный формат входных данных"
// @Failure		500		{object}	actorResponse	"Проблемы на стороне сервера"
// @Failure		401		{object}	actorResponse	"Ошибка авторизации"
// @Failure		403		{object}	actorResponse	"Ошибка авторизации"
// @Router			/actors/ [post]
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

// @Summary		Обновление полей актёра
// @Description	Обновляет поля актёра по id
// @Tags			actors
// @Security		ApiKeyAuth
// @Accept			json
// @Produce		json
// @Param			actor_id	query		string			true	"id актёра"
// @Param			input		body		updateActorData	true	"Новые поля"
// @Success		200			{object}	actorResponse	"Информация об актёре"
// @Failure		404			{object}	actorResponse	"Актёра не существует"
// @Failure		400			{object}	actorResponse	"Неверный формат входных данных"
// @Failure		500			{object}	actorResponse	"Проблемы на стороне сервера"
// @Failure		401			{object}	actorResponse	"Ошибка авторизации"
// @Failure		403			{object}	actorResponse	"Ошибка авторизации"
// @Router			/actors/ [put]
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

// @Summary		Удаление актёра
// @Description	Удаляет актёра по id
// @Tags			actors
// @Security		ApiKeyAuth
// @Produce		json
// @Param			actor_id	query		string			true	"id актёра"
// @Success		200			{object}	actorResponse	"Пустая структура"
// @Failure		404			{object}	actorResponse	"Актёра не существует"
// @Failure		400			{object}	actorResponse	"Неверный формат входных данных"
// @Failure		500			{object}	actorResponse	"Проблемы на стороне сервера"
// @Failure		401			{object}	actorResponse	"Ошибка авторизации"
// @Failure		403			{object}	actorResponse	"Ошибка авторизации"
// @Router			/actors/ [delete]
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

// @Summary		Получение актёра
// @Description	Возвращает актёра с указанным id
// @Tags			actors
// @Security		ApiKeyAuth
// @Produce		json
// @Param			actor_id	query		string			true	"id актёра"
// @Success		200			{object}	actorResponse	"Информация об актёре"
// @Failure		404			{object}	actorResponse	"Актёра не существует"
// @Failure		500			{object}	actorResponse	"Проблемы на стороне сервера"
// @Failure		400			{object}	actorResponse	"Неверный формат входных данных"
// @Failure		401			{object}	actorResponse	"Ошибка авторизации"
// @Failure		403			{object}	actorResponse	"Ошибка авторизации"
// @Router			/actors/ [get]
func getActorHandler(ctx context.Context, a app.App) http.HandlerFunc {
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

// @Summary		Получение списка актёров
// @Description	Возвращает список актёров
// @Tags			actors
// @Security		ApiKeyAuth
// @Produce		json
// @Success		200	{object}	actorListResponse	"Информация об актёрах"
// @Failure		500	{object}	actorListResponse	"Проблемы на стороне сервера"
// @Failure		401	{object}	actorListResponse	"Ошибка авторизации"
// @Failure		403	{object}	actorListResponse	"Ошибка авторизации"
// @Router			/actors/list/ [get]
func getActorsListHandler(ctx context.Context, a app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userId, err := strconv.ParseUint(r.Header.Get("Authorization"), 10, 64)
		if err != nil {
			http.Error(w, errorResponse(model.ErrUnauthorized), http.StatusUnauthorized)
			return
		}

		actors, err := a.GetActors(ctx, userId)

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
	}
}
