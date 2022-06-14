package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/dgrijalva/jwt-go/v4"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/u2lentaru/billbck/internal/api"
	"github.com/u2lentaru/billbck/internal/models"
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

	pgs := PG{dbpool}
	apg := api.APG{Dbpool: dbpool}
	route := mux.NewRouter()

	route.HandleFunc("/", handleRoot).Methods("GET", "OPTIONS")
	route.HandleFunc("/", pgs.handleLogin).Methods("GET", "OPTIONS")
	route.HandleFunc("/admin/", pgs.handleAdmin).Methods("GET", "OPTIONS")
	// route.HandleFunc("/login/", pgs.handleLogin).Methods("GET", "OPTIONS")
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
		// log.Fatal(http.ListenAndServe(":8080", route))
		n.Use(negroni.HandlerFunc(utils.AuthValidate))
	}
	// } else {
	// 	// log.Fatal(http.ListenAndServe(":8080", n))
	// }

	n.UseHandler(route)
	log.Fatal(http.ListenAndServe(":8080", n))
}

func (s *PG) handleAdmin(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(fmt.Sprintf("It's handleAdmin!\n")))

	ctx := context.Background()
	v := ""

	err := s.dbpool.QueryRow(ctx, "SELECT version();").Scan(&v)

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Write([]byte(fmt.Sprintf(v, "\n")))
	return
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(fmt.Sprintf("It's handleRoot!\n")))

	return
}

/*
// handleLogin godoc
// @Summary List user forms
// @Description Get user forms
// @Tags login
// @Produce  json
// @Success 200 {array} []string
// @Failure 500
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Router / [get]*/
func (s *PG) handleLogin(w http.ResponseWriter, r *http.Request) {
	type Claims struct {
		Username string `json:"username"`
		Password string `json:"password"`
		jwt.StandardClaims
	}

	utils.SetupResponse(&w)

	if (*r).Method == "OPTIONS" {
		return
	}

	ctx := context.Background()

	var jwtSecretKey = []byte("jwt_secret_key")

	claims := &Claims{}

	bearerToken := r.Header.Get("Authorization")

	if bearerToken == "" {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(fmt.Sprintf("Token is empty!\n")))

		return
	}

	tokenString := strings.Split(bearerToken, " ")[1]

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtSecretKey, nil
	})
	if err != nil {
		http.Error(w, "Error parse token (login): "+err.Error(), 500)
		return
	}

	if !token.Valid {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(fmt.Sprintf("Error: %s\n", err)))

		return
	}

	// out_arr := []string{}
	out_arr := []models.LoginForm{}

	rows, err := s.dbpool.Query(ctx, "SELECT * from func_get_user_forms($1);", claims.Username)
	// rows, err := s.dbpool.Query(ctx, "SELECT to_jsonb(func_get_user_forms($1));", claims.Username)

	if err != nil {
		http.Error(w, "Error SELECT * from func_get_user_forms: "+err.Error(), 500)
		return
	}

	defer rows.Close()

	// w.Write([]byte("["))

	for rows.Next() {
		f := models.LoginForm{}
		err = rows.Scan(&f.Form, &f.Rights, &f.UserId)
		// err = rows.Scan(&f)

		// w.Write([]byte(f + "\n"))

		if err != nil {
			log.Println("failed to scan row:", err)
		}

		out_arr = append(out_arr, f)
	}
	// w.Write([]byte("]"))

	output, err := json.Marshal(out_arr)
	if err != nil {
		http.Error(w, "Error marshal output: "+err.Error(), 500)
		return
	}

	w.Write(output)

	return

}
