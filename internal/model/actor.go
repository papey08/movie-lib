package model

type Gender string

const (
	Unknown Gender = ""
	Male    Gender = "male"
	Female  Gender = "female"
)

type Actor struct {
	Id         uint64
	FirstName  string
	SecondName string
	Gender
	Movies []*Movie
}

type UpdateActor struct {
	FirstName  string
	SecondName string
	Gender
}
