package httpserver

import "movie-lib/internal/model"

type createActorData struct {
	FirstName    string `json:"first_name"`
	SecondName   string `json:"second_name"`
	model.Gender `json:"gender"`
}

type updateActorData struct {
	FirstName    string `json:"first_name"`
	SecondName   string `json:"second_name"`
	model.Gender `json:"gender"`
}

type createMovieData struct {
	Title       string   `json:"title"`
	Description string   `json:"description"`
	ReleaseDate int64    `json:"release_date"`
	Rating      float64  `json:"rating"`
	ActorsId    []uint64 `json:"actors"`
}

type updateMovieData struct {
	Title       string   `json:"title"`
	Description string   `json:"description"`
	ReleaseDate int64    `json:"release_date"`
	Rating      float64  `json:"rating"`
	ActorsId    []uint64 `json:"actors"`
}
