package utils

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/jackc/pgx/v4"
	"github.com/u2lentaru/billbck/internal/models"
)

func MWSetupResponse(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	SetupResponse(&w)
	if (*r).Method == "OPTIONS" {
		return
	}
	next(w, r)
}

//SetupResponse(w *http.ResponseWriter)
func SetupResponse(w *http.ResponseWriter) {
	// (*w).Header().Set("content-type", "application/json")
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Credentials", "true")
	(*w).Header().Set("Access-Control-Allow-Methods", "OPTIONS, GET, POST")
	(*w).Header().Set("Access-Control-Allow-Headers", "Content-Type, Depth, User-Agent, X-File-Size, X-Requested-With, If-Modified-Since, X-File-Name, Cache-ControlContent-Type, Accept, Authorization")
}

//Authentication, authorization and token validation middleware function.
func AuthValidate(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	type Claims struct {
		Username string `json:"username"`
		Password string `json:"password"`
		jwt.StandardClaims
	}
	// SetupResponse(&w)

	if (*r).Method == "OPTIONS" {
		next(w, r)
	}

	if r.URL.Path == "/" {
		next(w, r)
	}

	var jwtSecretKey = []byte("jwt_secret_key")

	claims := &Claims{}

	bearerToken := r.Header.Get("Authorization")

	if bearerToken == "" {
		w.Write([]byte(fmt.Sprintf("Token is empty!\n")))

		return
	}

	tokenString := strings.Split(bearerToken, " ")[1]

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtSecretKey, nil
	})

	if err != nil {
		http.Error(w, "Error parse token: "+err.Error(), 500)
		return
	}

	if !token.Valid {
		w.Write([]byte(fmt.Sprintf("Error: %s\n", err)))

		return
	}

	//Active Directory authentication github.com/korylprince/go-ad-auth/v3

	// config := &auth.Config{
	// 	Server:   "ldap.example.com",
	// 	Port:     389,
	// 	BaseDN:   "OU=Users,DC=example,DC=com",
	// 	Security: auth.SecurityStartTLS,
	// }

	// status, err := auth.Authenticate(config, claims.Username, claims.Password)

	// if err != nil {
	// 	//handle err
	// 	return
	// }

	// if !status {
	// 	//handle failed authentication
	// 	return
	// }

	url := "postgres://" + claims.Username + ":" + claims.Password + "@" + os.Getenv("DB_HOST") + ":5432/postgres"

	conn, err := pgx.Connect(context.Background(), url)
	defer conn.Close(context.Background())

	if err != nil {
		// w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(fmt.Sprintf("Unable to connect to database: %s\n", err)))
		return
	}

	url = "postgres://postgres:postgres@" + os.Getenv("DB_HOST") + ":5432/postgres"

	conn, err = pgx.Connect(context.Background(), url)
	if err != nil {
		w.Write([]byte(fmt.Sprintf("Unable to connect to database: %s\n", err)))
		return
	}
	defer conn.Close(context.Background())

	re := regexp.MustCompile(`/[0-9]+$`)
	funcname := string(re.ReplaceAll([]byte(r.URL.Path), []byte("")))

	if strings.Contains(r.URL.Path, "/swagger/") {
		funcname = "/swagger/"
	}

	user_auth := false
	// w.Write([]byte(fmt.Sprintf("claims.Username: %s funcname: %s\n", claims.Username, funcname)))

	err = conn.QueryRow(context.Background(), "SELECT func_auth_user($1,$2);", claims.Username, funcname).Scan(&user_auth)

	if err != nil {
		http.Error(w, "Error SELECT func_auth_user: "+err.Error(), 500)
		return
	}

	if !user_auth {
		w.Write([]byte(fmt.Sprintf("Access denied to %s\n", funcname)))

		return
	}

	next(w, r)

	return
}

//FirstOfPrevMonth() time.Time - returns first day of current previous month
func FirstOfPrevMonth() time.Time {
	now := time.Now()
	currentYear, currentMonth, _ := now.Date()
	currentLocation := now.Location()

	firstOfPrevMonth := time.Date(currentYear, currentMonth-1, 1, 0, 0, 0, 0, currentLocation)

	if currentMonth == 1 {
		firstOfPrevMonth = time.Date(currentYear-1, 12, 1, 0, 0, 0, 0, currentLocation)
	}
	return firstOfPrevMonth
}

//LastOfPrevMonth() time.Time - returns last day of current previous month
func LastOfPrevMonth() time.Time {
	now := time.Now()
	currentYear, currentMonth, _ := now.Date()
	currentLocation := now.Location()

	firstOfPrevMonth := time.Date(currentYear, currentMonth-1, 1, 0, 0, 0, 0, currentLocation)
	lastOfPrevMonth := firstOfPrevMonth.AddDate(0, 1, -1)

	if currentMonth == 1 {
		lastOfPrevMonth = firstOfPrevMonth.AddDate(0, 1, -1)
	}
	return lastOfPrevMonth
}

//NullableString(s string) sql.NullString - returns null for empty string
func NullableString(s string) sql.NullString {
	return sql.NullString{String: s, Valid: s != ""}
}

//NullableInt(s string) sql.NullInt - returns null for 0
func NullableInt(i int32) sql.NullInt32 {
	return sql.NullInt32{Int32: i, Valid: i != 0}
}

//NullableBool(s string) sql.NullBool - returns null for empty string
func NullableBool(b bool, n bool) sql.NullBool {
	return sql.NullBool{Bool: b, Valid: n}
}

//GetAuth(r *http.Request, formName string) models.Auth - returns Auth struct
func GetAuth(r *http.Request, formName string) models.Auth {
	////////////////////////////////////////////////////////////////////////////////////////////////////////
	if len(formName) > 0 {
		return models.Auth{Create: true, Read: true, Update: true, Delete: true}
	}
	////////////////////////////////////////////////////////////////////////////////////////////////////////

	type Claims struct {
		Username string `json:"username"`
		Password string `json:"password"`
		jwt.StandardClaims
	}

	g := models.Auth{}

	// noauth := false
	// sna, ok := os.LookupEnv("NOAUTH")

	// if !ok {
	// 	noauth = false
	// } else {
	// 	noauth = (sna == "TRUE")
	// }

	// if (len(os.Args) > 1 && os.Args[1] == "noauth") || noauth {
	// 	return models.Auth{Create: false, Read: true, Update: false, Delete: false}
	// }

	var jwtSecretKey = []byte("jwt_secret_key")

	claims := &Claims{}

	bearerToken := r.Header.Get("Authorization")

	if bearerToken == "" {
		log.Println("Token is empty!")
		return models.Auth{Create: false, Read: false, Update: false, Delete: false}
	}

	if len(strings.Split(bearerToken, " ")) < 2 {
		//Bearer lost
		log.Println("Incorrect token format!")
		return models.Auth{Create: false, Read: false, Update: false, Delete: false}
	}

	tokenString := strings.Split(bearerToken, " ")[1]

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtSecretKey, nil
	})

	if err != nil {
		log.Println("Error parse token: " + err.Error())
		return models.Auth{Create: false, Read: false, Update: false, Delete: false}
	}

	if !token.Valid {
		log.Printf("Error: %s\n", err)
		return models.Auth{Create: false, Read: false, Update: false, Delete: false}
	}

	url := "postgres://postgres:postgres@" + os.Getenv("DB_HOST") + ":5432/postgres"

	conn, err := pgx.Connect(context.Background(), url)
	if err != nil {
		log.Printf("Unable to connect to database: %s\n", err)
		return models.Auth{Create: false, Read: false, Update: false, Delete: false}
	}
	defer conn.Close(context.Background())

	// log.Printf("utils.GetAuth!!! claims.Username:%s, formName:%s\n", claims.Username, formName)

	err = conn.QueryRow(context.Background(), "SELECT * FROM func_user_getauth($1,$2);", claims.Username, formName).Scan(&g.Create, &g.Read,
		&g.Update, &g.Delete)

	if err != nil {
		log.Printf("Error SELECT * FROM func_user_getauth: %s\n", err)
		return models.Auth{Create: false, Read: false, Update: false, Delete: false}
	}

	// log.Printf("User: %s, passw:%s\n", claims.Username, claims.Password)
	return g
}
