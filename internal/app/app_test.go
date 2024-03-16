package app

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"movie-lib/internal/model"
	"movie-lib/internal/repo"
	"movie-lib/pkg/logger"
	"os"
	"testing"
	"time"
)

const (
	adminUserId   = 1
	regularUserId = 2
)

var (
	ctx    = context.Background()
	actors = []model.Actor{
		{
			FirstName:  "TestActor01",
			SecondName: "",
			Gender:     model.Male,
			Movies:     make([]model.Movie, 0),
		},
		{
			FirstName:  "TestActor02",
			SecondName: "",
			Gender:     model.Male,
			Movies:     make([]model.Movie, 0),
		},
		{
			FirstName:  "TestActor03",
			SecondName: "",
			Gender:     model.Male,
			Movies:     make([]model.Movie, 0),
		},
		{
			FirstName:  "TestActor04",
			SecondName: "",
			Gender:     model.Male,
			Movies:     make([]model.Movie, 0),
		},
		{
			FirstName:  "TestActor05",
			SecondName: "",
			Gender:     model.Female,
			Movies:     make([]model.Movie, 0),
		},
	}
	movies = []model.Movie{
		{
			Title:       "TestMovie01",
			Description: "",
			ReleaseDate: time.Unix(1577826000, 0), // 2020-01-01
			Rating:      5,
		},
		{
			Title:       "TestMovie02",
			Description: "",
			ReleaseDate: time.Unix(1609448400, 0), // 2021-01-01
			Rating:      4,
		},
		{
			Title:       "TestMovie03",
			Description: "",
			ReleaseDate: time.Unix(1640984400, 0), // 2022-01-01
			Rating:      3,
		},
		{
			Title:       "TestMovie04",
			Description: "",
			ReleaseDate: time.Unix(1640984400, 0), // 2022-01-01
			Rating:      3,
		},
		{
			Title:       "",
			Description: "Invalid movie",
			ReleaseDate: time.Unix(1640984400, 0), // 2022-01-01
			Rating:      2,
		},
	}
)

type appTestSuite struct {
	suite.Suite

	moviesIdsToDelete []uint64
	actorsIdsToDelete []uint64

	conn    *pgx.Conn
	service App
}

func (s *appTestSuite) SetupSuite() {
	var err error
	// setting configs
	viper.SetConfigFile("../../config/config-local.yml")
	err = viper.ReadInConfig()
	if err != nil {
		s.Fail("unable to read configs")
	}
	for i := 0; i < 30; i++ { // 30 attempts to connect to postgres
		s.conn, err = pgx.Connect(context.Background(), fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
			viper.GetString("postgres-movie-lib.username"),
			viper.GetString("postgres-movie-lib.password"),
			viper.GetString("postgres-movie-lib.host"),
			viper.GetInt("postgres-movie-lib.port"),
			viper.GetString("postgres-movie-lib.dbname"),
			viper.GetString("postgres-movie-lib.sslmode"),
		))
		if err == nil {
			break
		}
		time.Sleep(time.Second)
	}
	if err != nil {
		s.Fail("unable to connect to postgres")
	}

	s.service = New(repo.New(s.conn), logger.DefaultLogger(os.Stdout))

	// add 3 actors to database
	for i := 0; i < 3; i++ {
		addedActor, _ := s.service.CreateActor(ctx, adminUserId, actors[i])
		actors[i].Id = addedActor.Id
		s.actorsIdsToDelete = append(s.actorsIdsToDelete, addedActor.Id)
	}
	movies[0].ActorsId = append(movies[0].ActorsId, actors[0].Id, actors[1].Id)
	movies[1].ActorsId = append(movies[1].ActorsId, actors[1].Id, actors[2].Id)
	movies[2].ActorsId = append(movies[2].ActorsId, actors[0].Id, actors[1].Id, 0)
}

func (s *appTestSuite) TearDownSuite() {
	for _, id := range s.moviesIdsToDelete {
		_ = s.service.DeleteMovie(ctx, adminUserId, id)
	}
	for _, id := range s.actorsIdsToDelete {
		_ = s.service.DeleteActor(ctx, adminUserId, id)
	}

	_ = s.conn.Close(context.Background())
}

type createMovieTest struct {
	description string
	user        uint64
	movie       *model.Movie
	err         error
}

func (s *appTestSuite) TestCreateMovie() {
	tests := []createMovieTest{
		{
			description: "successful creating of the movie",
			user:        adminUserId,
			movie:       &movies[0],
			err:         nil,
		},
		{
			description: "successful creating of the movie",
			user:        adminUserId,
			movie:       &movies[1],
			err:         nil,
		},
		{
			description: "successful creating of the movie",
			user:        adminUserId,
			movie:       &movies[3],
			err:         nil,
		},
		{
			description: "creating of the movie with non existing actor",
			user:        adminUserId,
			movie:       &movies[2],
			err:         model.ErrActorNotExists,
		},
		{
			description: "creating of the invalid movie",
			user:        adminUserId,
			movie:       &movies[4],
			err:         model.ErrValidationError,
		},
		{
			description: "creating of the movie with no admin rights",
			user:        regularUserId,
			movie:       &movies[1],
			err:         model.ErrPermissionDenied,
		},
		{
			description: "creating of the movie with non existing user",
			user:        0,
			movie:       &movies[0],
			err:         model.ErrUserNotExists,
		},
	}

	for _, test := range tests {
		s.T().Run(test.description, func(t *testing.T) {
			addedMovie, err := s.service.CreateMovie(ctx, test.user, *test.movie)
			if err == nil {
				test.movie.Id = addedMovie.Id
				s.moviesIdsToDelete = append(s.moviesIdsToDelete, addedMovie.Id)
			}
			assert.ErrorIs(s.T(), err, test.err)
		})
	}
}

type updateMovieTest struct {
	description string
	user        uint64
	id          uint64
	upd         model.UpdateMovie
	res         model.Movie
	err         error
}

func (s *appTestSuite) TestUpdateMovie() {
	tests := []updateMovieTest{
		{
			description: "successful updating of movie description",
			user:        adminUserId,
			id:          movies[0].Id,
			upd: model.UpdateMovie{
				Title:       movies[0].Title,
				Description: "aaa",
				ReleaseDate: movies[0].ReleaseDate,
				Rating:      movies[0].Rating,
				Actors:      movies[0].ActorsId,
			},
			res: model.Movie{
				Id:          movies[0].Id,
				Description: "aaa",
				Rating:      movies[0].Rating,
			},
			err: nil,
		},
		{
			description: "updating of non existing movie",
			user:        adminUserId,
			id:          0,
			upd: model.UpdateMovie{
				Title:       movies[0].Title,
				Description: "bbb",
				ReleaseDate: movies[0].ReleaseDate,
				Rating:      movies[0].Rating,
				Actors:      movies[0].ActorsId,
			},
			res: model.Movie{},
			err: model.ErrMovieNotExists,
		},
		{
			description: "adding to movie list of actors non existing actor",
			user:        adminUserId,
			id:          movies[0].Id,
			upd: model.UpdateMovie{
				Title:       movies[0].Title,
				Description: movies[0].Description,
				ReleaseDate: movies[0].ReleaseDate,
				Rating:      movies[0].Rating,
				Actors:      []uint64{0},
			},
			res: model.Movie{},
			err: model.ErrActorNotExists,
		},
		{
			description: "updating movie with no rights",
			user:        regularUserId,
			id:          movies[0].Id,
			upd: model.UpdateMovie{
				Title:       movies[0].Title,
				Description: "bbb",
				ReleaseDate: movies[0].ReleaseDate,
				Rating:      movies[0].Rating,
				Actors:      movies[0].ActorsId,
			},
			res: model.Movie{},
			err: model.ErrPermissionDenied,
		},
		{
			description: "updating movie with non existing user",
			user:        0,
			id:          movies[0].Id,
			upd: model.UpdateMovie{
				Title:       movies[0].Title,
				Description: "bbb",
				ReleaseDate: movies[0].ReleaseDate,
				Rating:      movies[0].Rating,
				Actors:      movies[0].ActorsId,
			},
			res: model.Movie{},
			err: model.ErrUserNotExists,
		},
		{
			description: "updating movie with invalid fields",
			user:        adminUserId,
			id:          movies[0].Id,
			upd: model.UpdateMovie{
				Title:       movies[0].Title,
				Description: movies[0].Description,
				ReleaseDate: movies[0].ReleaseDate,
				Rating:      500,
				Actors:      movies[0].ActorsId,
			},
			res: model.Movie{},
			err: model.ErrValidationError,
		},
	}

	for _, test := range tests {
		s.T().Run(test.description, func(t *testing.T) {
			updatedMovie, err := s.service.UpdateMovie(ctx, test.user, test.id, test.upd)
			assert.Equal(t, test.res.Id, updatedMovie.Id)
			assert.Equal(t, test.res.Description, updatedMovie.Description)
			assert.Equal(t, test.res.Rating, updatedMovie.Rating)
			assert.ErrorIs(t, err, test.err)
		})
	}
	movies[0].Description = "aaa"
}

type deleteMovieTest struct {
	description string
	user        uint64
	id          uint64
	err         error
}

func (s *appTestSuite) TestDeleteMovie() {
	tests := []deleteMovieTest{
		{
			description: "successful deleting of the movie",
			user:        adminUserId,
			id:          movies[1].Id,
			err:         nil,
		},
		{
			description: "deleting of non existing movie",
			user:        adminUserId,
			id:          0,
			err:         model.ErrMovieNotExists,
		},
		{
			description: "deleting of the movie with no rights",
			user:        regularUserId,
			id:          movies[1].Id,
			err:         model.ErrPermissionDenied,
		},
		{
			description: "deleting of the movie with non existing user",
			user:        0,
			id:          movies[1].Id,
			err:         model.ErrUserNotExists,
		},
	}

	for _, test := range tests {
		s.T().Run(test.description, func(t *testing.T) {
			err := s.service.DeleteMovie(ctx, test.user, test.id)
			assert.ErrorIs(t, err, test.err)
		})
	}
}

type getMovieTest struct {
	description string
	user        uint64
	id          uint64
	res         model.Movie
	err         error
}

func (s *appTestSuite) TestGetMovie() {
	tests := []getMovieTest{
		{
			description: "successful getting of the movie",
			user:        adminUserId,
			id:          movies[0].Id,
			res:         movies[0],
			err:         nil,
		},
		{
			description: "successful getting of the movie with no admin rights",
			user:        regularUserId,
			id:          movies[0].Id,
			res:         movies[0],
			err:         nil,
		},
		{
			description: "getting of the movie with non existing user",
			user:        0,
			id:          movies[0].Id,
			res:         model.Movie{},
			err:         model.ErrUserNotExists,
		},
		{
			description: "getting of the movie with non existing id",
			user:        adminUserId,
			id:          0,
			res:         model.Movie{},
			err:         model.ErrMovieNotExists,
		},
	}

	for _, test := range tests {
		s.T().Run(test.description, func(t *testing.T) {
			gotMovie, err := s.service.GetMovie(ctx, test.user, test.id)
			assert.Equal(t, test.res.Id, gotMovie.Id)
			assert.Equal(t, test.res.Title, gotMovie.Title)
			assert.ErrorIs(t, err, test.err)
		})
	}
}

type getMoviesTest struct {
	description string
	user        uint64
	sortBy      model.SortParam

	// moviesIdList содержит в себе правильный порядок следования тестовых
	// фильмов в списке всех фильмов
	moviesIdList []uint64

	// moviesIdSet содержит в себе id тестовых фильмов, которые должны быть в
	// списке всех фильмов
	moviesIdSet map[uint64]struct{}

	err error
}

func (s *appTestSuite) TestGetMovies() {
	tests := []getMoviesTest{
		{
			description:  "getting of movies list with default sort params",
			user:         adminUserId,
			sortBy:       "",
			moviesIdList: []uint64{movies[0].Id, movies[3].Id},
			moviesIdSet: map[uint64]struct{}{
				movies[0].Id: {},
				movies[3].Id: {},
			},
			err: nil,
		},
		{
			description:  "getting of movies list sorted by title",
			user:         adminUserId,
			sortBy:       model.Title,
			moviesIdList: []uint64{movies[0].Id, movies[3].Id},
			moviesIdSet: map[uint64]struct{}{
				movies[0].Id: {},
				movies[3].Id: {},
			},
			err: nil,
		},
		{
			description:  "getting of movies list sorted by release date",
			user:         adminUserId,
			sortBy:       model.ReleaseDate,
			moviesIdList: []uint64{movies[0].Id, movies[3].Id},
			moviesIdSet: map[uint64]struct{}{
				movies[0].Id: {},
				movies[3].Id: {},
			},
			err: nil,
		},
		{
			description:  "getting of movies list sorted by rating",
			user:         adminUserId,
			sortBy:       model.Rating,
			moviesIdList: []uint64{movies[3].Id, movies[0].Id},
			moviesIdSet: map[uint64]struct{}{
				movies[0].Id: {},
				movies[3].Id: {},
			},
			err: nil,
		},
		{
			description:  "getting of movies list with no admin rights",
			user:         regularUserId,
			sortBy:       "",
			moviesIdList: []uint64{movies[0].Id, movies[3].Id},
			moviesIdSet: map[uint64]struct{}{
				movies[0].Id: {},
				movies[3].Id: {},
			},
			err: nil,
		},
	}

	for _, test := range tests {
		s.T().Run(test.description, func(t *testing.T) {
			gotMoviesList, err := s.service.GetMovies(ctx, test.user, test.sortBy)
			assert.ErrorIs(s.T(), err, test.err)

			// Здесь происходит проверка на то, что тестовые фильмы в списке всех
			// фильмов располагаются в правильном порядке, т.е. правильно отсортированы
			moviesSequence := make([]uint64, 0, len(test.moviesIdList))
			for _, m := range gotMoviesList {
				if _, ok := test.moviesIdSet[m.Id]; ok {
					moviesSequence = append(moviesSequence, m.Id)
				}
			}
			assert.Len(s.T(), moviesSequence, len(test.moviesIdList))
			assert.Equal(s.T(), moviesSequence, test.moviesIdList)
		})
	}
}

type searchMoviesTest struct {
	description string
	user        uint64
	pattern     string

	// moviesIdSet содержит в себе id тестовых фильмов, которые должны быть в
	// списке найденных фильмов
	moviesIdSet map[uint64]struct{}

	err error
}

func (s *appTestSuite) TestSearchMovies() {
	tests := []searchMoviesTest{
		{
			description: "getting film by pattern",
			user:        adminUserId,
			pattern:     "TestMovie01",
			moviesIdSet: map[uint64]struct{}{
				movies[0].Id: {},
			},
			err: nil,
		},
		{
			description: "getting film by actor pattern",
			user:        adminUserId,
			pattern:     "TestActor01",
			moviesIdSet: map[uint64]struct{}{
				movies[0].Id: {},
			},
			err: nil,
		},
		{
			description: "getting film by pattern with no admin rights",
			user:        regularUserId,
			pattern:     "TestMovie01",
			moviesIdSet: map[uint64]struct{}{
				movies[0].Id: {},
			},
			err: nil,
		},
	}

	for _, test := range tests {
		s.T().Run(test.description, func(t *testing.T) {
			gotMoviesList, err := s.service.SearchMovies(ctx, test.user, test.pattern)
			assert.ErrorIs(s.T(), err, test.err)

			// Здесь происходит проверка на то, что все фильмы, которые нужно
			// найти, содержатся в полученном списке
			gotMoviesSet := make(map[uint64]struct{})
			for _, m := range gotMoviesList {
				if _, ok := test.moviesIdSet[m.Id]; ok {
					gotMoviesSet[m.Id] = struct{}{}
				}
			}
			assert.Len(s.T(), gotMoviesSet, len(test.moviesIdSet))
			assert.Equal(s.T(), test.moviesIdSet, gotMoviesSet)
		})
	}
}

type createActorTest struct {
	description string
	user        uint64
	actor       *model.Actor
	err         error
}

func (s *appTestSuite) TestCreateActor() {
	tests := []createActorTest{
		{
			description: "successful creation of the actor",
			user:        adminUserId,
			actor:       &actors[3],
			err:         nil,
		},
		{
			description: "creation of the actor with no admin rights",
			user:        regularUserId,
			actor:       &actors[4],
			err:         model.ErrPermissionDenied,
		},
		{
			description: "creation of the actor with non existing user",
			user:        0,
			actor:       &actors[4],
			err:         model.ErrUserNotExists,
		},
	}

	for _, test := range tests {
		s.T().Run(test.description, func(t *testing.T) {
			addedActor, err := s.service.CreateActor(ctx, test.user, *test.actor)
			if addedActor.Id != 0 {
				s.actorsIdsToDelete = append(s.actorsIdsToDelete, addedActor.Id)
				test.actor.Id = addedActor.Id
			}
			assert.ErrorIs(t, err, test.err)
		})
	}
}

type updateActorTest struct {
	description string
	user        uint64
	id          uint64
	upd         model.UpdateActor
	res         model.Actor
	err         error
}

func (s *appTestSuite) TestUpdateActor() {
	tests := []updateActorTest{
		{
			description: "successful update of actor name",
			user:        adminUserId,
			id:          actors[0].Id,
			upd: model.UpdateActor{
				FirstName:  actors[0].FirstName,
				SecondName: "NewSecondName",
				Gender:     actors[0].Gender,
			},
			res: model.Actor{
				Id:         actors[0].Id,
				FirstName:  actors[0].FirstName,
				SecondName: "NewSecondName",
				Gender:     actors[0].Gender,
			},
			err: nil,
		},
		{
			description: "update of non existing actor",
			user:        adminUserId,
			id:          0,
			upd: model.UpdateActor{
				FirstName:  "NewFirstName",
				SecondName: "NewSecondName",
				Gender:     model.Male,
			},
			res: model.Actor{},
			err: model.ErrActorNotExists,
		},
		{
			description: "update of actor with no admin rights",
			user:        regularUserId,
			id:          actors[0].Id,
			upd: model.UpdateActor{
				FirstName:  "NewFirstName",
				SecondName: "NewSecondName",
				Gender:     model.Male,
			},
			res: model.Actor{},
			err: model.ErrPermissionDenied,
		},
		{
			description: "update of actor with non existing user",
			user:        0,
			id:          actors[0].Id,
			upd: model.UpdateActor{
				FirstName:  "NewFirstName",
				SecondName: "NewSecondName",
				Gender:     model.Male,
			},
			res: model.Actor{},
			err: model.ErrUserNotExists,
		},
	}

	for _, test := range tests {
		s.T().Run(test.description, func(t *testing.T) {
			updatedActor, err := s.service.UpdateActor(ctx, test.user, test.id, test.upd)
			assert.Equal(t, test.res.Id, updatedActor.Id)
			assert.Equal(t, test.res.FirstName, updatedActor.FirstName)
			assert.Equal(t, test.res.SecondName, updatedActor.SecondName)
			assert.ErrorIs(t, err, test.err)
		})
	}
	actors[0].SecondName = "NewSecondName"
}

type deleteActorTest struct {
	description string
	user        uint64
	id          uint64
	err         error
}

func (s *appTestSuite) TestDeleteActor() {
	tests := []deleteActorTest{
		{
			description: "successful deleting of actor",
			user:        adminUserId,
			id:          actors[2].Id,
			err:         nil,
		},
		{
			description: "deleting of non existing actor",
			user:        adminUserId,
			id:          0,
			err:         model.ErrActorNotExists,
		},
		{
			description: "deleting of actor with no admin rights",
			user:        regularUserId,
			id:          actors[2].Id,
			err:         model.ErrPermissionDenied,
		},
		{
			description: "deleting of actor with non existing user",
			user:        0,
			id:          actors[2].Id,
			err:         model.ErrUserNotExists,
		},
	}

	for _, test := range tests {
		s.T().Run(test.description, func(t *testing.T) {
			err := s.service.DeleteActor(ctx, test.user, test.id)
			assert.ErrorIs(t, err, test.err)
		})
	}
}

type getActorTest struct {
	description string
	user        uint64
	id          uint64
	res         model.Actor
	err         error
}

func (s *appTestSuite) TestGetActor() {
	tests := []getActorTest{
		{
			description: "successful getting of actor",
			user:        adminUserId,
			id:          actors[0].Id,
			res:         actors[0],
			err:         nil,
		},
		{
			description: "getting of actor with no admin rights",
			user:        regularUserId,
			id:          actors[0].Id,
			res:         actors[0],
			err:         nil,
		},
		{
			description: "getting of non existing actor",
			user:        adminUserId,
			id:          0,
			res:         model.Actor{},
			err:         model.ErrActorNotExists,
		},
		{
			description: "getting of actor with non existing user",
			user:        0,
			id:          actors[0].Id,
			res:         model.Actor{},
			err:         model.ErrUserNotExists,
		},
	}

	for _, test := range tests {
		s.T().Run(test.description, func(t *testing.T) {
			gotActor, err := s.service.GetActor(ctx, test.user, test.id)
			assert.Equal(t, test.res.Id, gotActor.Id)
			assert.Equal(t, test.res.FirstName, gotActor.FirstName)
			assert.ErrorIs(t, err, test.err)
		})
	}
}

type getActorsTest struct {
	description string
	user        uint64
	actorsIds   map[uint64]struct{}
	err         error
}

func (s *appTestSuite) TestGetActors() {
	tests := []getActorsTest{
		{
			description: "successful getting of list of actors",
			user:        adminUserId,
			actorsIds: map[uint64]struct{}{
				actors[0].Id: {},
				actors[1].Id: {},
				actors[3].Id: {},
			},
			err: nil,
		},
		{
			description: "getting of list of actors with no admin rights",
			user:        regularUserId,
			actorsIds: map[uint64]struct{}{
				actors[0].Id: {},
				actors[1].Id: {},
				actors[3].Id: {},
			},
			err: nil,
		},
		{
			description: "getting of list of actors with non existing user",
			user:        0,
			actorsIds:   map[uint64]struct{}{},
			err:         model.ErrUserNotExists,
		},
	}

	for _, test := range tests {
		s.T().Run(test.description, func(t *testing.T) {
			actorsList, err := s.service.GetActors(ctx, test.user)
			actorsIdSet := make(map[uint64]struct{})
			for _, actor := range actorsList {
				if _, ok := test.actorsIds[actor.Id]; ok {
					actorsIdSet[actor.Id] = struct{}{}
				}
			}
			assert.Len(t, actorsIdSet, len(test.actorsIds))
			assert.ErrorIs(t, err, test.err)
		})
	}
}

func TestAppTestSuite(t *testing.T) {
	suite.Run(t, new(appTestSuite))
}
