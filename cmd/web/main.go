package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/cleomon/bookings/internal/config"
	"github.com/cleomon/bookings/internal/driver"
	"github.com/cleomon/bookings/internal/handlers"
	"github.com/cleomon/bookings/internal/helpers"
	"github.com/cleomon/bookings/internal/models"
	"github.com/cleomon/bookings/internal/render"
)

const portNumber = ":8080"

var app config.AppConfig
var session *scs.SessionManager
var infoLog *log.Logger
var errorLog *log.Logger

// main is the main function
func main() {
	db, err := run()
	if err != nil {
		log.Fatal(err)
	}
	defer db.SQL.Close()

	fmt.Printf("Staring application on port %s\n", portNumber)

	srv := &http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}

	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}

func run() (*driver.DB, error) {
	// What ma I going to put in the session
	gob.Register(models.Reservation{})
	gob.Register(models.User{})
	gob.Register(models.Room{})
	gob.Register(models.Restriction{})

	// change this to true when in production
	app.InProduction = false

	infoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	app.InfoLog = infoLog

	errorLog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	app.ErrorLog = errorLog

	// set up the session
	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	app.Session = session

	// Connect to Database
	log.Println("Connecting to the database...")
	// db, err := pgx.Connect(context.Background(), "postgres://postgres:4141@localhost:5432/test_connect")
	db, err := driver.ConnectSQL("postgres://postgres:4141@localhost:5432/bookings")
	if err != nil {
		log.Fatal("Cannot connect to database! Dying...")
	}

	err = db.SQL.Ping()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Cannot ping database!")
	}
	log.Println("Pingged database!")

	defer db.SQL.Close()
	log.Println("Connected to the database...")

	rows, err := db.SQL.Exec(`select * from rooms`)

	log.Println(rows.RowsAffected())
	if err != nil {
		fmt.Fprintf(os.Stderr, "horrível")
	}

	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("cannot create template cache")
		return nil, err
	}

	app.TemplateCache = tc
	app.UseCache = false

	repo := handlers.NewRepo(&app, db)
	handlers.NewHandlers(repo)
	render.NewRenderer(&app)
	helpers.NewHelpers(&app)

	return db, nil
}
