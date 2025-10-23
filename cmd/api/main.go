package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"api.etin.dev/internal/assets"
	"api.etin.dev/internal/data"
	_ "github.com/lib/pq"
)

type config struct {
	port          int
	env           string
	dsn           string
	adminEmail    string
	adminPassword string
	deployWebhook string
	cors          struct {
		trustedOrigins []string
	}
	cloudinary struct {
		cloudName string
		apiKey    string
		apiSecret string
		folder    string
	}
}

type application struct {
	config     config
	logger     *log.Logger
	models     data.Models
	assetModel assetSaver
	assets     assets.Uploader
	swagger    []byte
	sessions   *sessionManager
	httpClient *http.Client
}

func main() {
	var cfg config

	var corsTrustedOrigins string

	flag.IntVar(&cfg.port, "port", 4000, "API server port")
	flag.StringVar(&cfg.env, "env", "dev", "Environment (dev|stage|prod)")
	flag.StringVar(&cfg.dsn, "dsn", os.Getenv("WEBSITE_DB_DSN"), "PostgreSQL DSN")
	flag.StringVar(&cfg.adminEmail, "admin-email", os.Getenv("WEBSITE_ADMIN_EMAIL"), "Admin login email")
	flag.StringVar(&cfg.adminPassword, "admin-password", os.Getenv("WEBSITE_ADMIN_PASSWORD"), "Admin login password")
	flag.StringVar(&corsTrustedOrigins, "cors-trusted-origins", os.Getenv("WEBSITE_CORS_TRUSTED_ORIGINS"), "Space separated list of trusted CORS origins")
	flag.StringVar(&cfg.cloudinary.cloudName, "cloudinary-cloud-name", os.Getenv("WEBSITE_CLOUDINARY_CLOUD_NAME"), "Cloudinary cloud name")
	flag.StringVar(&cfg.cloudinary.apiKey, "cloudinary-api-key", os.Getenv("WEBSITE_CLOUDINARY_API_KEY"), "Cloudinary API key")
	flag.StringVar(&cfg.cloudinary.apiSecret, "cloudinary-api-secret", os.Getenv("WEBSITE_CLOUDINARY_API_SECRET"), "Cloudinary API secret")
	flag.StringVar(&cfg.cloudinary.folder, "cloudinary-folder", os.Getenv("WEBSITE_CLOUDINARY_FOLDER"), "Optional Cloudinary folder for uploads")
	flag.StringVar(&cfg.deployWebhook, "deploy-webhook-url", os.Getenv("WEBSITE_DEPLOY_WEBHOOK_URL"), "Optional URL to trigger frontend deployments")
	flag.Parse()

	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	cfg.cors.trustedOrigins = parseTrustedOrigins(corsTrustedOrigins)
	if len(cfg.cors.trustedOrigins) == 0 {
		cfg.cors.trustedOrigins = []string{"https://admin.etin.dev", "https://etin.dev"}
	}

	for i, origin := range cfg.cors.trustedOrigins {
		cfg.cors.trustedOrigins[i] = normalizeOrigin(origin)
	}

	if cfg.adminEmail == "" || cfg.adminPassword == "" {
		logger.Fatal("Admin credentials must be provided")
	}

	if cfg.cloudinary.cloudName == "" {
		logger.Fatal("Cloudinary cloud name must be provided")
	}

	if cfg.cloudinary.apiKey == "" {
		logger.Fatal("Cloudinary API key must be provided")
	}

	if cfg.cloudinary.apiSecret == "" {
		logger.Fatal("Cloudinary API secret must be provided")
	}

	db, err := openDB(cfg.dsn)
	if err != nil {
		logger.Fatal(err)
	}
	defer db.Close()

	logger.Printf("database connection pool established")

	uploader, err := assets.NewCloudinaryUploader(
		cfg.cloudinary.cloudName,
		cfg.cloudinary.apiKey,
		cfg.cloudinary.apiSecret,
		cfg.cloudinary.folder,
	)
	if err != nil {
		logger.Fatal(err)
	}

	models := data.NewModels(db)

	app := &application{
		config:     cfg,
		logger:     logger,
		models:     models,
		assetModel: models.Assets,
		assets:     uploader,
		swagger:    embeddedSwagger,
		sessions:   newSessionManager(24 * time.Hour),
		httpClient: &http.Client{Timeout: 10 * time.Second},
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

func parseTrustedOrigins(input string) []string {
	if input == "" {
		return nil
	}

	fields := strings.FieldsFunc(input, func(r rune) bool {
		switch r {
		case ',', ' ', '\n', '\t':
			return true
		default:
			return false
		}
	})

	var origins []string
	for _, field := range fields {
		trimmed := strings.TrimSpace(field)
		if trimmed != "" {
			origins = append(origins, trimmed)
		}
	}

	return origins
}
