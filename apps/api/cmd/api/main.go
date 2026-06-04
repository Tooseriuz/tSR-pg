package main

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/tooseriuz/tsr-pg/apps/api/internal/adapters/config"
	httpadapter "github.com/tooseriuz/tsr-pg/apps/api/internal/adapters/http"
	"github.com/tooseriuz/tsr-pg/apps/api/internal/adapters/postgres"
	"github.com/tooseriuz/tsr-pg/apps/api/internal/adapters/storage/gcs"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	pool, err := pgxpool.New(context.Background(), cfg.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	imageStorage, err := gcs.NewImageStorage(context.Background(), gcs.Config{
		BucketName:          cfg.GCSBucketName,
		Endpoint:            cfg.GCSEndpoint,
		PublicBaseURL:       cfg.GCSPublicBaseURL,
		SignedURLAccessID:   cfg.GCSSignedURLAccessID,
		SignedURLPrivateKey: cfg.GCSSignedURLPrivateKey,
		SignedURLHostname:   cfg.GCSSignedURLHostname,
		SignedURLInsecure:   cfg.GCSSignedURLInsecure,
	})
	if err != nil {
		log.Fatal(err)
	}

	router := httpadapter.NewRouter(
		httpadapter.WithJourneyRepository(postgres.NewJourneyRepository(pool)),
		httpadapter.WithImageStorage(imageStorage),
	)

	log.Printf("api listening on %s", cfg.Address())
	if err := router.Run(cfg.Address()); err != nil {
		log.Fatal(err)
	}
}
