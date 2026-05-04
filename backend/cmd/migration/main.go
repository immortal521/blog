package main

import (
	"context"
	"fmt"
	"log"

	"blog-server/config"
	"blog-server/datastore"
	"blog-server/ent"
	"blog-server/ent/migrate"
	"blog-server/logger"
)

func main() {
	cfg := config.MustLoad()
	log := logger.NewLogger(cfg)
	ds, err := datastore.NewDataStore(cfg, log)
	if err != nil {
		log.Fatal(fmt.Sprintf("failed to connect to database: %v", err))
	}
	defer func() {
		if err := ds.Close(); err != nil {
			log.Info(fmt.Sprintf("failed to close database client: %v", err))
		}
	}()
	createSchema(ds.Client(context.Background()))
}

func createSchema(client *ent.Client) {
	if err := client.Schema.Create(
		context.Background(),
		migrate.WithDropIndex(true),
		migrate.WithDropColumn(true),
		migrate.WithForeignKeys(false),
	); err != nil {
		log.Fatalf("failed creating schema resources: %v", err)
	}
}
