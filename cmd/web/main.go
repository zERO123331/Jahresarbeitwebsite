package main

import (
	"Jahresarbeitwebsite/internal/models"
	"database/sql"
	"flag"
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"os"
	"time"

	_ "github.com/lib/pq"
)

type config struct {
	dsn  string
	port int
	env  string
}

type application struct {
	config        config
	logger        *slog.Logger
	templateCache map[string]*template.Template
	shop          *models.ShopModel
	users         *models.UserModel
	db            struct {
		dsn          string
		maxOpenConns int
		maxIdleConns int
		maxIdleTime  time.Duration
	}
}

func main() {
	var cfg config
	flag.IntVar(&cfg.port, "port", 4000, "HTTP port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")
	flag.StringVar(&cfg.dsn, "dsn", os.Getenv("SHOP_DB_DSN"), "PostgreSQL DSN")
	maxOpenConns := flag.Int("db-max-open-conns", 25, "maximum number of open connections to the database")
	maxIdleConns := flag.Int("db-max-idle-conns", 25, "maximum number of idle connections in the pool")
	maxIdleTime := flag.Duration("db-max-idle-time", 15*time.Minute, "maximum amount of time a connection may be idle before being closed")
	flag.Parse()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	templateCache, err := newTemplateCache()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	db, err := OpenDB(cfg.dsn)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	logger.Info("connection to database established")
	defer db.Close()

	app := &application{
		config:        cfg,
		logger:        logger,
		templateCache: templateCache,
		shop:          &models.ShopModel{DB: db},
		users:         &models.UserModel{DB: db},
		db: struct {
			dsn          string
			maxOpenConns int
			maxIdleConns int
			maxIdleTime  time.Duration
		}{
			dsn:          cfg.dsn,
			maxOpenConns: *maxOpenConns,
			maxIdleConns: *maxIdleConns,
			maxIdleTime:  *maxIdleTime,
		},
	}

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.port),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		ErrorLog:     slog.NewLogLogger(logger.Handler(), slog.LevelError),
	}

	logger.Info("starting server", "port", srv.Addr, "env", cfg.env)

	err = srv.ListenAndServe()

	logger.Error(err.Error())
	os.Exit(1)
}

func OpenDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		db.Close()
		return nil, err
	}
	return db, nil
}
