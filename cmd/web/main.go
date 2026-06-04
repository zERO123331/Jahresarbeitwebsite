package main

import (
	"Jahresarbeitwebsite/internal/cdn"
	"Jahresarbeitwebsite/internal/models"
	"context"
	"database/sql"
	"errors"
	"flag"
	"fmt"
	"html/template"
	"log/slog"
	"net/http/httputil"
	"net/url"
	"os"
	"time"

	"github.com/alexedwards/scs/postgresstore"
	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form/v4"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type config struct {
	port int
	env  string
	db   struct {
		dsn          string
		maxOpenConns int
		maxIdleConns int
		maxIdleTime  time.Duration
	}

	cdn struct {
		endpoint  string
		accessKey string
		secretKey string
		bucket    string
		secure    bool
	}

	limiter struct {
		rps     float64
		burst   int
		enabled bool
	}
}

type application struct {
	config         config
	logger         *slog.Logger
	templateCache  map[string]*template.Template
	formDecoder    *form.Decoder
	models         models.Models
	cdn            *cdn.CDN
	sessionManager *scs.SessionManager
	reverseProxy   *httputil.ReverseProxy
}

func main() {
	var cfg config
	flag.IntVar(&cfg.port, "port", 4000, "HTTP port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")

	flag.StringVar(&cfg.db.dsn, "dsn", os.Getenv("SHOP_DB_DSN"), "PostgreSQL DSN")

	flag.IntVar(&cfg.db.maxOpenConns, "db-max-open-conns", 25, "maximum number of open connections to the database")
	flag.IntVar(&cfg.db.maxIdleConns, "db-max-idle-conns", 25, "maximum number of idle connections in the pool")
	flag.DurationVar(&cfg.db.maxIdleTime, "db-max-idle-time", 15*time.Minute, "maximum amount of time a connection may be idle before being closed")
	flag.Float64Var(&cfg.limiter.rps, "limiter-rps", 10, "requests per second limit")
	flag.IntVar(&cfg.limiter.burst, "limiter-burst", 20, "maximum burst of requests")
	flag.BoolVar(&cfg.limiter.enabled, "limiter-enabled", true, "enable rate limiter")

	flag.StringVar(&cfg.cdn.endpoint, "cdn-endpoint", os.Getenv("SHOP_CDN_ENDPOINT"), "Minio endpoint")
	flag.StringVar(&cfg.cdn.accessKey, "cdn-acess-key", os.Getenv("SHOP_CDN_ACCESS_KEY"), "Minio access key")
	flag.StringVar(&cfg.cdn.secretKey, "cdn-secret-key", os.Getenv("SHOP_CDN_SECRET_KEY"), "Minio secret key")
	flag.StringVar(&cfg.cdn.bucket, "cdn-bucket", os.Getenv("SHOP_CDN_BUCKET"), "Minio bucket")
	flag.BoolVar(&cfg.cdn.secure, "cdn-secure", false, "Minio secure connection")
	flag.Parse()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	templateCache, err := newTemplateCache()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	formDecoder := form.NewDecoder()

	db, err := OpenDB(cfg)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	logger.Info("connection to database established")
	defer db.Close()

	migrationDriver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	migrator, err := migrate.NewWithDatabaseInstance("file://migrations", "postgres", migrationDriver)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	err = migrator.Up()
	if err != nil && err != migrate.ErrNoChange {
		logger.Error(err.Error())
		os.Exit(1)
	}
	if errors.Is(err, migrate.ErrNoChange) {
		logger.Info("no migrations to apply")
	} else {
		logger.Info("migrations applied")
	}

	cdnClient, err := OpenCDN(cfg)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	logger.Info("connection to CDN established")

	reverseProxy, err := imageProxy(cfg)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	reverseProxy.ErrorLog = slog.NewLogLogger(logger.Handler(), slog.LevelError)

	sessionManager := scs.New()
	sessionManager.Store = postgresstore.New(db)
	sessionManager.Lifetime = 12 * time.Hour
	sessionManager.Cookie.Secure = true

	app := &application{
		config:         cfg,
		logger:         logger,
		templateCache:  templateCache,
		formDecoder:    formDecoder,
		models:         models.NewModels(db),
		sessionManager: sessionManager,
		cdn:            cdn.New(cdnClient, cfg.cdn.bucket),
		reverseProxy:   reverseProxy,
	}

	err = app.serve()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
}

func OpenDB(cfg config) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.db.dsn)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(cfg.db.maxOpenConns)
	db.SetMaxIdleConns(cfg.db.maxIdleConns)
	db.SetConnMaxIdleTime(cfg.db.maxIdleTime)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}

func OpenCDN(cfg config) (*minio.Client, error) {
	client, err := minio.New(cfg.cdn.endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.cdn.accessKey, cfg.cdn.secretKey, ""),
		Secure: cfg.cdn.secure,
	})
	if err != nil {
		return nil, err
	}
	isOnline := client.IsOnline()
	if !isOnline {
		return nil, fmt.Errorf("minio is not online")
	}
	return client, nil
}

func imageProxy(cfg config) (*httputil.ReverseProxy, error) {
	cdnURL, err := url.Parse(fmt.Sprintf("http://%s", cfg.cdn.endpoint))
	if err != nil {
		return nil, err
	}
	proxy := httputil.NewSingleHostReverseProxy(cdnURL)
	return proxy, nil
}
