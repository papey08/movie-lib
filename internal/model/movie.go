package model

import "time"

type Movie struct {
	Id          uint64
	Title       string
	Description string
	ReleaseDate time.Time
	Rating      float64
	Actors      []*Actor
	ActorsId    []uint64
}

type UpdateMovie struct {
	Title       string
	Description string
	ReleaseDate time.Time
	Rating      float64
	Actors      []uint64
}

type SortParam string

const (
	None        SortParam = ""
	Title       SortParam = "title"
	Rating      SortParam = "rating"
	ReleaseDate SortParam = "release_date"
)
