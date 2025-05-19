package main

import (
	"context"
	"database/sql"
	"flag"
	"html/template"
	"log/slog"
	"os"
	"time"

	"github.com/chahiim/ticket_tracker/internal/data"
	"github.com/golangcollege/sessions"
	_ "github.com/lib/pq"
)

type application struct {
	addr          *string
	tickets       *data.TicketModel
	users         *data.UsersModel
	session       *sessions.Session
	logger        *slog.Logger
	templateCache map[string]*template.Template
}

func main() {
	addr := flag.String("addr", "", "HTTP network address")
	dsn := flag.String("dsn", "", "PostgreSQL DSN")
	secret := flag.String("secret", "CidapE50eufaLsgdJ*20+jEhr0rw_uYh", "Secret key")

	flag.Parse()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	db, err := openDB(*dsn)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	session := sessions.New([]byte(*secret))
	session.Lifetime = 1 * time.Hour
	session.Secure = true

	logger.Info("database connection pool established")
	templateCache, err := newTemplateCache()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	defer db.Close()

	app := &application{
		addr:          addr,
		tickets:       &data.TicketModel{DB: db},
		users:         &data.UsersModel{DB: db},
		logger:        logger,
		session:       session,
		templateCache: templateCache,
	}

	err = app.serve()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}
