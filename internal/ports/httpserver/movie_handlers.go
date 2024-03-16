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
	"time"
)

// @Summary		Добавление фильма
// @Description	Добавляет новый фильм
// @Tags			movies
// @Security		ApiKeyAuth
// @Accept			json
// @Produce		json
// @Param			input	body		createMovieData	true	"Информация о новом фильме"
// @Success		200		{object}	movieResponse	"Информация о фильме"
// @Failure		400		{object}	movieResponse	"Неверный формат входных данных"
// @Failure		404		{object}	movieResponse	"Актёра из списка не существует"
// @Failure		500		{object}	movieResponse	"Проблемы на стороне сервера"
// @Failure		401		{object}	movieResponse	"Ошибка авторизации"
// @Failure		403		{object}	movieResponse	"Ошибка авторизации"
// @Router			/movies [post]
func createMovieHandler(ctx context.Context, a app.App) http.HandlerFunc {
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
		var data createMovieData
		if err = json.Unmarshal(body, &data); err != nil {
			http.Error(w, errorResponse(model.ErrInvalidInput), http.StatusBadRequest)
			return
		}

		movie, err := a.CreateMovie(ctx, userId, model.Movie{
			Title:       data.Title,
			Description: data.Description,
			ReleaseDate: time.Unix(data.ReleaseDate, 0),
			Rating:      data.Rating,
			ActorsId:    data.ActorsId,
		})

		switch {
		case err == nil:
			w.WriteHeader(http.StatusOK)
			_, _ = fmt.Fprintf(w, movieResponseOk(movie))
		case errors.Is(err, model.ErrValidationError):
			http.Error(w, errorResponse(model.ErrValidationError), http.StatusBadRequest)
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

// @Summary		Обновление фильма
// @Description	Обновляет поля фильма по id
// @Tags			movies
// @Security		ApiKeyAuth
// @Accept			json
// @Produce		json
// @Param			movie_id	query		string			true	"id фильма"
// @Param			input		body		updateMovieData	true	"Новые поля"
// @Success		200			{object}	movieResponse	"Информация о фильме"
// @Failure		400			{object}	movieResponse	"Неверный формат входных данных"
// @Failure		404			{object}	movieResponse	"Фильма либо актёра из списка не существует"
// @Failure		500			{object}	movieResponse	"Проблемы на стороне сервера"
// @Failure		401			{object}	movieResponse	"Ошибка авторизации"
// @Failure		403			{object}	movieResponse	"Ошибка авторизации"
// @Router			/movies [put]
func updateMovieHandler(ctx context.Context, a app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userId, err := strconv.ParseUint(r.Header.Get("Authorization"), 10, 64)
		if err != nil {
			http.Error(w, errorResponse(model.ErrUnauthorized), http.StatusUnauthorized)
			return
		}
		movieId, err := strconv.ParseUint(r.URL.Query().Get("movie_id"), 10, 64)
		if err != nil {
			http.Error(w, errorResponse(model.ErrInvalidInput), http.StatusBadRequest)
			return
		}
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, errorResponse(model.ErrInvalidInput), http.StatusBadRequest)
			return
		}
		var data updateMovieData
		if err = json.Unmarshal(body, &data); err != nil {
			http.Error(w, errorResponse(model.ErrInvalidInput), http.StatusBadRequest)
			return
		}

		movie, err := a.UpdateMovie(ctx, userId, movieId, model.UpdateMovie{
			Title:       data.Title,
			Description: data.Description,
			ReleaseDate: time.Unix(data.ReleaseDate, 0),
			Rating:      data.Rating,
			Actors:      data.ActorsId,
		})

		switch {
		case err == nil:
			w.WriteHeader(http.StatusOK)
			_, _ = fmt.Fprintf(w, movieResponseOk(movie))
		case errors.Is(err, model.ErrMovieNotExists):
			http.Error(w, errorResponse(model.ErrMovieNotExists), http.StatusNotFound)
		case errors.Is(err, model.ErrActorNotExists):
			http.Error(w, errorResponse(model.ErrActorNotExists), http.StatusNotFound)
		case errors.Is(err, model.ErrValidationError):
			http.Error(w, errorResponse(model.ErrValidationError), http.StatusBadRequest)
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

// @Summary		Удаление фильма
// @Description	Удаление фильма по id
// @Tags			movies
// @Security		ApiKeyAuth
// @Accept			json
// @Produce		json
// @Param			movie_id	query		string			true	"id фильма"
// @Success		200			{object}	movieResponse	"Пустая структура"
// @Failure		400			{object}	movieResponse	"Неверный формат входных данных"
// @Failure		404			{object}	movieResponse	"Фильма не существует"
// @Failure		500			{object}	movieResponse	"Проблемы на стороне сервера"
// @Failure		401			{object}	movieResponse	"Ошибка авторизации"
// @Failure		403			{object}	movieResponse	"Ошибка авторизации"
// @Router			/movies [delete]
func deleteMovieHandler(ctx context.Context, a app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userId, err := strconv.ParseUint(r.Header.Get("Authorization"), 10, 64)
		if err != nil {
			http.Error(w, errorResponse(model.ErrUnauthorized), http.StatusUnauthorized)
			return
		}
		movieId, err := strconv.ParseUint(r.URL.Query().Get("movie_id"), 10, 64)
		if err != nil {
			http.Error(w, errorResponse(model.ErrInvalidInput), http.StatusBadRequest)
			return
		}

		err = a.DeleteMovie(ctx, userId, movieId)

		switch {
		case err == nil:
			w.WriteHeader(http.StatusOK)
			_, _ = fmt.Fprintf(w, errorResponse(nil))
		case errors.Is(err, model.ErrMovieNotExists):
			http.Error(w, errorResponse(model.ErrMovieNotExists), http.StatusNotFound)
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

// @Summary		Получение списка фильмов
// @Description	Возвращает список фильмов
// @Tags			movies
// @Security		ApiKeyAuth
// @Accept			json
// @Produce		json
// @Param			pattern	query		string				false	"Поиск по названию фильма/фамилии/имени актёра"
// @Param			sort_by	query		string				false	"Параметр для сортировки. Поддерживаемые параметры: title, rating, release_date"
// @Success		200		{object}	movieListResponse	"Информация о фильмах"
// @Failure		400		{object}	movieListResponse	"Неверный формат входных данных"
// @Failure		500		{object}	movieListResponse	"Проблемы на стороне сервера"
// @Failure		401		{object}	movieListResponse	"Ошибка авторизации"
// @Failure		403		{object}	movieListResponse	"Ошибка авторизации"
// @Router			/movies/list [get]
func getMovieListHandler(ctx context.Context, a app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userId, err := strconv.ParseUint(r.Header.Get("Authorization"), 10, 64)
		if err != nil {
			http.Error(w, errorResponse(model.ErrUnauthorized), http.StatusUnauthorized)
			return
		}

		if r.URL.Query().Has("pattern") {
			pattern := r.URL.Query().Get("pattern")

			var movies []model.Movie
			movies, err = a.SearchMovies(ctx, userId, pattern)

			switch {
			case err == nil:
				w.WriteHeader(http.StatusOK)
				_, _ = fmt.Fprintf(w, movieListResponseOk(movies))
			case errors.Is(err, model.ErrUserNotExists):
				http.Error(w, errorResponse(model.ErrUserNotExists), http.StatusForbidden)
			case errors.Is(err, model.ErrDatabaseError):
				http.Error(w, errorResponse(model.ErrDatabaseError), http.StatusInternalServerError)
			default:
				http.Error(w, errorResponse(model.ErrServiceError), http.StatusInternalServerError)
			}
		} else {
			sortParam := r.URL.Query().Get("sort_by")

			var movies []model.Movie
			movies, err = a.GetMovies(ctx, userId, model.SortParam(sortParam))

			switch {
			case err == nil:
				w.WriteHeader(http.StatusOK)
				_, _ = fmt.Fprintf(w, movieListResponseOk(movies))
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

// @Summary		Получение фильма по id
// @Description	Возвращает фильм с указанным id
// @Tags			movies
// @Security		ApiKeyAuth
// @Accept			json
// @Produce		json
// @Param			movie_id	query		string			true	"id фильма"
// @Success		200			{object}	movieResponse	"Пустая структура"
// @Failure		404			{object}	movieResponse	"Фильма не существует"
// @Failure		400			{object}	movieResponse	"Неверный формат входных данных"
// @Failure		500			{object}	movieResponse	"Проблемы на стороне сервера"
// @Failure		401			{object}	movieResponse	"Ошибка авторизации"
// @Failure		403			{object}	movieResponse	"Ошибка авторизации"
// @Router			/movies [get]
func getMovieHandler(ctx context.Context, a app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userId, err := strconv.ParseUint(r.Header.Get("Authorization"), 10, 64)
		if err != nil {
			http.Error(w, errorResponse(model.ErrUnauthorized), http.StatusUnauthorized)
			return
		}

		var movieId uint64
		movieId, err = strconv.ParseUint(r.URL.Query().Get("movie_id"), 10, 64)
		if err != nil {
			http.Error(w, errorResponse(model.ErrInvalidInput), http.StatusBadRequest)
			return
		}

		var movie model.Movie
		movie, err = a.GetMovie(ctx, userId, movieId)

		switch {
		case err == nil:
			w.WriteHeader(http.StatusOK)
			_, _ = fmt.Fprintf(w, movieResponseOk(movie))
		case errors.Is(err, model.ErrMovieNotExists):
			http.Error(w, errorResponse(model.ErrMovieNotExists), http.StatusNotFound)
		case errors.Is(err, model.ErrUserNotExists):
			http.Error(w, errorResponse(model.ErrUserNotExists), http.StatusForbidden)
		case errors.Is(err, model.ErrDatabaseError):
			http.Error(w, errorResponse(model.ErrDatabaseError), http.StatusInternalServerError)
		default:
			http.Error(w, errorResponse(model.ErrServiceError), http.StatusInternalServerError)
		}
	}
}
