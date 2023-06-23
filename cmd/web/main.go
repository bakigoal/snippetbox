package main

import (
	"database/sql"
	"flag"
	"github.com/bakigoal/snippetbox/internal/models"
	"github.com/go-playground/form/v4"
	_ "github.com/lib/pq"
	"html/template"
	"log"
	"net/http"
	"os"
)

type application struct {
	infoLog       *log.Logger
	errorLog      *log.Logger
	snippets      *models.SnippetModel
	templateCache map[string]*template.Template
	formDecoder   *form.Decoder
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
		infoLog:       log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime),
		errorLog:      log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile),
		snippets:      &models.SnippetModel{DB: db},
		templateCache: templateCache,
		formDecoder:   form.NewDecoder(),
	}

	infoLog.Printf("Starting server on %s", *addr)
	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  app.routes(),
	}
	err = srv.ListenAndServe()
	errorLog.Fatal(err)
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
