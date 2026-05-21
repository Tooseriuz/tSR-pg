package main

import (
	"log"

	"github.com/tooseriuz/tsr-pg/apps/api/internal/adapters/config"
	httpadapter "github.com/tooseriuz/tsr-pg/apps/api/internal/adapters/http"
)

func main() {
	cfg := config.Load()
	router := httpadapter.NewRouter()

	log.Printf("api listening on %s", cfg.Address())
	if err := router.Run(cfg.Address()); err != nil {
		log.Fatal(err)
	}
}
