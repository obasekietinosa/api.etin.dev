package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"api.etin.dev/internal/data"
	"api.etin.dev/pkg/openapi"
	_ "github.com/lib/pq"
)

const version = "1.0.0"

type config struct {
	port    int
	env     string
	dsn     string
	authKey string
}

type application struct {
	config  config
	logger  *log.Logger
	models  data.Models
	swagger []byte
}

func main() {
	var cfg config

	flag.IntVar(&cfg.port, "port", 4000, "API server port")
	flag.StringVar(&cfg.env, "env", "dev", "Environment (dev|stage|prod)")
	flag.StringVar(&cfg.dsn, "dsn", os.Getenv("WEBSITE_DB_DSN"), "PostgreSQL DSN")
	flag.StringVar(&cfg.authKey, "key", os.Getenv("WEBSITE_AUTH_KEY"), "Admin auth key")
	flag.Parse()

	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	if cfg.authKey == "" {
		logger.Fatal("No auth key provided")
	}

	db, err := openDB(cfg.dsn)
	if err != nil {
		logger.Fatal(err)
	}
	defer db.Close()

	logger.Printf("database connection pool established")

	swaggerDoc, err := openapi.Build(version)
	if err != nil {
		logger.Fatal(err)
	}

	app := &application{
		config:  cfg,
		logger:  logger,
		models:  data.NewModels(db),
		swagger: swaggerDoc,
	}

	addr := fmt.Sprintf(":%d", cfg.port)

	srv := &http.Server{
		Addr:         addr,
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Minute,
		WriteTimeout: 10 * time.Minute,
	}

	logger.Printf("starting %s server on %s", cfg.env, addr)
	err = srv.ListenAndServe()
	logger.Fatal(err)
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
