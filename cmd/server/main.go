package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/spf13/viper"
	"movie-lib/internal/app"
	"movie-lib/internal/ports/httpserver"
	"movie-lib/internal/repo"
	"movie-lib/pkg/logger"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func SetConfigs(configPath string) error {
	viper.SetConfigFile(configPath)
	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("cannot read config file %w", err)
	}
	return nil
}

func ConnectToPostgres(ctx context.Context) (*pgx.Conn, error) {
	adsRepoUrl := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		viper.GetString("postgres-movie-lib.username"),
		viper.GetString("postgres-movie-lib.password"),
		viper.GetString("postgres-movie-lib.host"),
		viper.GetInt("postgres-movie-lib.port"),
		viper.GetString("postgres-movie-lib.dbname"),
		viper.GetString("postgres-movie-lib.sslmode"))

	// 30 attempts to connect to postgres when starting in docker container
	for i := 0; i < 30; i++ {
		conn, err := pgx.Connect(ctx, adsRepoUrl)
		if err != nil {
			time.Sleep(time.Second)
		} else {
			return conn, nil
		}
	}

	return nil, errors.New("unable to connect to postgres ads repo")
}

const (
	dockerConfigFile = "config/config-docker.yml"
	localConfigFile  = "config/config-local.yml"
)

//	@title						movie-lib
//	@version					1.0
//	@description				Swagger документация к API фильмотеки
//	@host						localhost:8080
//	@BasePath					/api/v1
//	@SecurityDefinitions.apikey	ApiKeyAuth
//	@in							header
//	@name						Authorization

func main() {
	ctx := context.Background()
	logs := logger.DefaultLogger(os.Stdout)

	isDocker := flag.Bool("docker", false, "flag if this project is running in docker container")
	flag.Parse()
	var configPath string
	if *isDocker {
		configPath = dockerConfigFile
	} else {
		configPath = localConfigFile
	}

	if err := SetConfigs(configPath); err != nil {
		logs.FatalLog(fmt.Sprintf("reading configs: %s", err.Error()))
	}

	conn, err := ConnectToPostgres(ctx)
	if err != nil {
		logs.FatalLog(fmt.Sprintf("connecting to postgres: %s", err.Error()))
	}
	logs.InfoLog("successfully connected to postgres")

	r := repo.New(conn)
	a := app.New(r, logs)

	host := viper.GetString("http-server.host")
	port := viper.GetInt("http-server.port")

	srv := httpserver.New(ctx, host, port, a, logs)

	go func() {
		_ = srv.ListenAndServe()
	}()
	logs.InfoLog("http server successfully started")

	// preparing graceful shutdown
	osSignals := make(chan os.Signal, 1)
	signal.Notify(osSignals, os.Interrupt, syscall.SIGTERM)
	signal.Notify(osSignals, os.Interrupt, syscall.SIGINT)

	// waiting for Ctrl+C
	<-osSignals

	shutdownCtx, cancel := context.WithTimeout(ctx, 30*time.Second) // 30s timeout to finish all active connections
	defer cancel()

	_ = srv.Shutdown(shutdownCtx)
	logs.InfoLog("successfully stopped http server")
}
