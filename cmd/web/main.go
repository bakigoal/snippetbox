package main

import (
	"context"
	"crypto/tls"
	"flag"
	"github.com/alexedwards/scs/pgxstore"
	"github.com/alexedwards/scs/v2"
	"github.com/bakigoal/snippetbox/internal/models"
	"github.com/go-playground/form/v4"
	"github.com/jackc/pgx/v5/pgxpool"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"
)

type application struct {
	infoLog        *log.Logger
	errorLog       *log.Logger
	snippets       *models.SnippetModel
	users          *models.UserModel
	templateCache  map[string]*template.Template
	formDecoder    *form.Decoder
	sessionManager *scs.SessionManager
}

func main() {
	addr := flag.String("addr", ":4000", "HTTP network address")
	defaultConn := "postgres://postgres:postgres@localhost:5432/go_dev?sslmode=disable&search_path=snippetbox"
	dsn := flag.String("dsn", defaultConn, "Postgres data source name")
	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	db, err := openDB(*dsn)
	if err != nil {
		errorLog.Fatal(err)
	}
	defer db.Close()

	templateCache, err := newTemplateCache()
	if err != nil {
		errorLog.Fatal(err)
	}

	app := &application{
		infoLog:        log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime),
		errorLog:       log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile),
		snippets:       &models.SnippetModel{DB: db},
		users:          &models.UserModel{DB: db},
		templateCache:  templateCache,
		formDecoder:    form.NewDecoder(),
		sessionManager: initSessionManager(db),
	}

	infoLog.Printf("Starting server on %s", *addr)
	srv := &http.Server{
		Addr:         *addr,
		ErrorLog:     errorLog,
		Handler:      app.routes(),
		TLSConfig:    getTlsConfig(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	err = srv.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem")
	errorLog.Fatal(err)
}

func getTlsConfig() *tls.Config {
	return &tls.Config{
		CurvePreferences: []tls.CurveID{tls.X25519, tls.CurveP256},
		MinVersion:       tls.VersionTLS12,
		MaxVersion:       tls.VersionTLS13,
	}
}

func openDB(dsn string) (*pgxpool.Pool, error) {
	conn, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		return nil, err
	}
	if err = conn.Ping(context.Background()); err != nil {
		return nil, err
	}
	return conn, nil
}

func initSessionManager(db *pgxpool.Pool) *scs.SessionManager {
	sessionManager := scs.New()
	sessionManager.Store = pgxstore.New(db)
	sessionManager.Lifetime = 12 * time.Hour
	return sessionManager
}
