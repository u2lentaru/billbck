package api

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/u2lentaru/billbck/internal/models"
)

func TestHandleGetCableResistance(t *testing.T) {
	//url := "postgres://postgres:postgres@localhost:5432/postgres"

	// cfg, err := pgxpool.ParseConfig(url)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// dbpool, err := pgxpool.ConnectConfig(context.Background(), cfg)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// tst_srv := APG{Dbpool: dbpool}
	router := mux.NewRouter()
	router.HandleFunc("/cableresistances/{id:[0-9]+}", HandleGetCableResistance).Methods("GET", "OPTIONS")

	req, _ := http.NewRequest("GET", "/cableresistances/7", nil)

	w := httptest.NewRecorder()
	(*w).Header().Set("Content-type", "application/json")

	vars := map[string]string{
		"id": "7",
	}

	req = mux.SetURLVars(req, vars)

	// vars := mux.Vars(req)
	// i := vars["id"]

	// log.Println("vars: ", vars)
	// log.Println("vars[id]: ", i)

	// tst_srv.HandleGetCableResistance(w, req)
	router.ServeHTTP(w, req)

	res := w.Result()
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)

	log.Println("body: ", string(body[:]))

	rcr := models.CableResistance_count{}

	err = json.Unmarshal(body, &rcr)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// log.Println(rcr)

	// if rcr.Values[0].Id == 0 {
	// 	t.Errorf("expected {7,'95 мм',0.195, false} got %v", rcr.Values[0])
	// }

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, models.CableResistance{Id: 7, CableResistanceName: "95 мм", Resistance: 0.195, MaterialType: false}, rcr.Values[0])

}
