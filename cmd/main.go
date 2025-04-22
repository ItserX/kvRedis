package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"

	"kvManager/internal/handlers"
	log "kvManager/internal/pkg/log"
	"kvManager/internal/storage"
)

func loadEnv() error {
	err := godotenv.Load()
	if err != nil {
		return err
	}
	return nil
}

func connectToRedis(addr string) (*redis.Client, error) {
	log.Logger.Infow("Connecting to Redis", "address", addr)

	opts := &redis.Options{
		Addr:     addr,
		Password: "",
		DB:       0,
	}

	client := redis.NewClient(opts)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := client.Ping(ctx).Result()
	if err != nil {
		log.Logger.Errorw("Failed to connect to Redis", "error", err, "address", addr)
		return nil, err
	}

	log.Logger.Info("Successfully connected to Redis")
	return client, nil
}

func setupRouter(client *redis.Client, logger *zap.SugaredLogger) *mux.Router {
	logger.Info("Setting up router and initializing storage")

	st := storage.NewRedisRepository(client)
	h := handlers.Handler{Repo: st}

	r := mux.NewRouter()
	r.HandleFunc("/kv", h.Add).Methods("POST")
	r.HandleFunc("/kv/{id}", h.Get).Methods("GET")
	r.HandleFunc("/kv/{id}", h.Update).Methods("PUT")
	r.HandleFunc("/kv/{id}", h.Delete).Methods("DELETE")

	logger.Info("Router setup completed")
	return r
}

func main() {
	err := log.SetupLogger()
	if err != nil {
		fmt.Printf("failed to initialize logger: %v", err)
		return
	}

	err = loadEnv()
	if err != nil {
		log.Logger.Errorw("Env load failing")
		return
	}

	appPort := os.Getenv("APP_PORT")
	redisAddr := os.Getenv("REDIS_ADDRESS")

	log.Logger.Info("Starting app")
	client, err := connectToRedis(redisAddr)
	if err != nil {
		return
	}

	defer func() {
		err := client.Close()
		if err != nil {
			log.Logger.Errorw("Connection to Redis is not closed", "error", err)
		}
	}()

	r := setupRouter(client, log.Logger)

	log.Logger.Infow("Starting HTTP server", "address", appPort)
	err = http.ListenAndServe(appPort, r)
	if err != nil {
		log.Logger.Errorw("HTTP server error", "error", err)
		return
	}
}
