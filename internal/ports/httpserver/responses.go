package httpserver

import (
	"encoding/json"
	"movie-lib/internal/model"
)

func errorResponse(err error) string {
	var resp actorResponse
	if err != nil {
		errStr := err.Error()
		resp.Err = &errStr
	}
	data, _ := json.Marshal(resp)
	return string(data)
}

func actorResponseOk(actor model.Actor) string {
	data := actorToActorData(actor)
	resp := actorResponse{
		Data: &data,
		Err:  nil,
	}
	body, _ := json.Marshal(resp)
	return string(body)
}

func movieResponseOk(movie model.Movie) string {
	data := movieToMovieData(movie)
	resp := movieResponse{
		Data: &data,
		Err:  nil,
	}
	body, _ := json.Marshal(resp)
	return string(body)
}

func actorListResponseOk(actors []model.Actor) string {
	data := actorsToActorListData(actors)
	resp := actorListResponse{
		Data: data,
		Err:  nil,
	}
	body, _ := json.Marshal(resp)
	return string(body)
}

func movieListResponseOk(movies []model.Movie) string {
	data := moviesToMovieListData(movies)
	resp := movieListResponse{
		Data: data,
		Err:  nil,
	}
	body, _ := json.Marshal(resp)
	return string(body)
}

func actorToActorData(actor model.Actor) actorData {
	data := actorData{
		Id:         actor.Id,
		FirstName:  actor.FirstName,
		SecondName: actor.SecondName,
		Gender:     actor.Gender,
	}

	data.Movies = make([]movieData, 0, len(actor.Movies))
	for _, movie := range actor.Movies {
		data.Movies = append(data.Movies, movieData{
			Id:          movie.Id,
			Title:       movie.Title,
			Description: movie.Description,
			ReleaseDate: movie.ReleaseDate.UTC().Unix(),
			Rating:      movie.Rating,
		})
	}
	return data
}

func movieToMovieData(movie model.Movie) movieData {
	data := movieData{
		Id:          movie.Id,
		Title:       movie.Title,
		Description: movie.Description,
		ReleaseDate: movie.ReleaseDate.UTC().Unix(),
		Rating:      movie.Rating,
	}

	data.Actors = make([]actorData, 0, len(movie.Actors))
	for _, actor := range movie.Actors {
		data.Actors = append(data.Actors, actorData{
			Id:         actor.Id,
			FirstName:  actor.FirstName,
			SecondName: actor.SecondName,
			Gender:     actor.Gender,
		})
	}
	return data
}

type actorData struct {
	Id           uint64 `json:"id"`
	FirstName    string `json:"first_name"`
	SecondName   string `json:"second_name"`
	model.Gender `json:"gender"`
	Movies       []movieData `json:"movies,omitempty"`
}

type actorResponse struct {
	Data *actorData `json:"data"`
	Err  *string    `json:"error"`
}

type movieData struct {
	Id          uint64      `json:"id"`
	Title       string      `json:"title"`
	Description string      `json:"description"`
	ReleaseDate int64       `json:"release_date"`
	Rating      float64     `json:"rating"`
	Actors      []actorData `json:"actors,omitempty"`
}

type movieResponse struct {
	Data *movieData `json:"data"`
	Err  *string    `json:"error"`
}

func actorsToActorListData(actors []model.Actor) []actorData {
	data := make([]actorData, 0, len(actors))
	for _, actor := range actors {
		data = append(data, actorToActorData(actor))
	}
	return data
}

type actorListResponse struct {
	Data []actorData `json:"data"`
	Err  *string     `json:"error"`
}

func moviesToMovieListData(movies []model.Movie) []movieData {
	data := make([]movieData, 0, len(movies))
	for _, movie := range movies {
		data = append(data, movieToMovieData(movie))
	}
	return data
}

type movieListResponse struct {
	Data []movieData `json:"data"`
	Err  *string     `json:"error"`
}
