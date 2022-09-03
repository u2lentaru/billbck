package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/u2lentaru/billbck/internal/api"
	"github.com/u2lentaru/billbck/internal/routes"
	"github.com/u2lentaru/billbck/internal/utils"

	// "github.com/mfuentesg/go-jwtmiddleware"

	httpSwagger "github.com/swaggo/http-swagger"
	_ "github.com/u2lentaru/billbck/docs"
	"github.com/urfave/negroni"
)

//PG - server struct
type PG struct {
	dbpool *pgxpool.Pool
}

// @title Billing Backend Server
// @version 1.0
// @description This is a backend server.
// @termsOfService http://swagger.io/terms/

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

//posterc.kz:44475 localhost:8080
// @host localhost:8080
// @BasePath /

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
// @description Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6InVzZXIxIiwicGFzc3dvcmQiOiJ1c2VyMSJ9.-qgJjYhayo7CT1YD1xLB36Xytf1HprRBeLbi5tZcOPE
func main() {
	ctx := context.Background()
	url := "postgres://postgres:postgres@" + os.Getenv("DB_HOST") + ":5432/postgres"
	// url := "postgres://postgres:postgres@" + os.Getenv("DB_HOST") + ":5432/billing"

	//ApiKeyAuth Bearer

	cfg, err := pgxpool.ParseConfig(url)
	if err != nil {
		log.Fatal(err)
	}

	cfg.MaxConns = 8
	cfg.MinConns = 1

	dbpool, err := pgxpool.ConnectConfig(ctx, cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer dbpool.Close()

	rows, err := dbpool.Query(ctx, "SELECT version();")
	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()

	for rows.Next() {
		v := ""
		err = rows.Scan(&v)

		if err != nil {
			log.Println("failed to scan row:", err)
		}

		log.Println("version:", v)
	}

	apg := api.APG{Dbpool: dbpool}
	route := mux.NewRouter()

	route.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler).Methods("GET", "OPTIONS")

	art := &apg
	routes.AddRoutes(route, art)

	// n := negroni.New(negroni.HandlerFunc(utils.AuthValidate))
	n := negroni.New(negroni.HandlerFunc(utils.MWSetupResponse))
	// n.Use(negroni.HandlerFunc(utils.AuthValidate))
	// n.UseHandler(route)

	log.Println("Server is listening at http://localhost:8080/")

	//go run . noauth - run without utils.AuthValidate middleware
	//docker run ... -e NOAUTH="TRUE" - run without utils.AuthValidate middleware
	noauth := false
	sna, ok := os.LookupEnv("NOAUTH")

	if !ok {
		noauth = false
	} else {
		noauth = (sna == "TRUE")
	}

	if !((len(os.Args) > 1 && os.Args[1] == "noauth") || noauth) {
		n.Use(negroni.HandlerFunc(utils.AuthValidate))
	}

	n.UseHandler(route)
	// log.Fatal(http.ListenAndServe(":8080", n))

	//GS start
	srv := &http.Server{
		Addr:    ":8080",
		Handler: n,
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()
	log.Print("Server Started")

	<-done
	log.Print("Server Stopped")

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer func() {
		// extra handling here
		log.Println("Sleep on")
		time.Sleep(time.Second * 1)
		log.Println("Sleep off")
		cancel()
	}()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server Shutdown Failed:%+v", err)
	}
	log.Print("Server Exited Properly")

}
