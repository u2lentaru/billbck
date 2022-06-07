package api

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v4"
	"github.com/u2lentaru/billbck/internal/models"
)

// HandleCities godoc
// @Summary List cities
// @Description get city list
// @Tags cities
// @Produce  json
// @Param page query int false "page number"
// @Param page_size query int false "page size"
// @Param cityname query string false "cityname search pattern"
// @Param ordering query string false "order by {id|cityname}"
// @Param desc query boolean false "descending order {true|false}"
// @Success 200 {array} models.City_count
// @Failure 500
// @Router /cities [get]
func (s *APG) HandleCities(w http.ResponseWriter, r *http.Request) {
	gs := models.City{}
	ctx := context.Background()
	out_arr := []models.City{}

	query := r.URL.Query()

	pg := 1
	spg, ok := query["page"]

	if ok && len(spg) > 0 {
		if pgt, err := strconv.Atoi(spg[0]); err != nil {
			pg = 1
		} else {
			pg = pgt
		}
	}

	pgs := 20
	spgs, ok := query["page_size"]
	if ok && len(spgs) > 0 {
		if pgst, err := strconv.Atoi(spgs[0]); err != nil {
			pgs = 20
		} else {
			pgs = pgst
		}
	}

	gs1 := ""
	gs1s, ok := query["cityname"]
	if ok && len(gs1s) > 0 {
		//case insensitive
		gs1 = strings.ToUpper(gs1s[0])
		//quotes
		re := regexp.MustCompile(`'`)
		gs1 = string(re.ReplaceAll([]byte(gs1), []byte("''")))
	}

	ord := 1
	ords, ok := query["ordering"]
	if !ok || len(ords) == 0 {
		ord = 1
	} else if ords[0] == "cityname" {
		ord = 2
	}

	dsc := false
	dscs, ok := query["desc"]
	if ok && len(dscs) > 0 {
		if !(dscs[0] == "0") {
			dsc = true
		}
	}

	rows, err := s.Dbpool.Query(ctx, "SELECT * from func_cities_get($1,$2,$3,$4,$5);", pg, pgs, gs1, ord, dsc)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&gs.Id, &gs.CityName)
		if err != nil {
			log.Println("failed to scan row:", err)
		}

		out_arr = append(out_arr, gs)
	}

	gsc := 0
	err = s.Dbpool.QueryRow(ctx, "SELECT * from func_cities_cnt($1);", gs1).Scan(&gsc)

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	auth := models.Auth{Create: true, Read: true, Update: true, Delete: true}
	out_count, err := json.Marshal(models.City_count{Values: out_arr, Count: gsc, Auth: auth})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Write(out_count)

	return
}

// HandleAddCity godoc
// @Summary Add city
// @Description add city
// @Tags cities
// @Accept json
// @Produce  json
// @Param a body models.AddCity true "New city"
// @Success 200 {array} models.Json_id
// @Failure 500
// @Router /cities_add [post]
func (s *APG) HandleAddCity(w http.ResponseWriter, r *http.Request) {
	a := models.AddCity{}
	body, err := ioutil.ReadAll(r.Body)

	defer r.Body.Close()

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	err = json.Unmarshal(body, &a)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	ai := 0
	err = s.Dbpool.QueryRow(context.Background(), "SELECT func_cities_add($1);", a.CityName).Scan(&ai)

	if err != nil {
		log.Println("Failed execute func_cities_add: ", err)
	}

	output, err := json.Marshal(models.Json_id{Id: ai})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(output)

	return

}

// HandleUpdCity godoc
// @Summary Update city
// @Description update city
// @Tags cities
// @Accept json
// @Produce  json
// @Param u body models.City true "Update city"
// @Success 200 {array} models.Json_id
// @Failure 500
// @Router /cities_upd [post]
func (s *APG) HandleUpdCity(w http.ResponseWriter, r *http.Request) {
	u := models.City{}
	body, err := ioutil.ReadAll(r.Body)

	defer r.Body.Close()

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	err = json.Unmarshal(body, &u)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	ui := 0
	err = s.Dbpool.QueryRow(context.Background(), "SELECT func_cities_upd($1,$2);", u.Id, u.CityName).Scan(&ui)

	if err != nil {
		log.Println("Failed execute func_cities_upd: ", err)
	}

	output, err := json.Marshal(models.Json_id{Id: ui})

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(output)

	return
}

// HandleDelCity godoc
// @Summary Delete cities
// @Description delete cities
// @Tags cities
// @Accept json
// @Produce  json
// @Param d body models.Json_ids true "Delete cities"
// @Success 200 {array} models.Json_ids
// @Failure 500
// @Router /cities_del [post]
func (s *APG) HandleDelCity(w http.ResponseWriter, r *http.Request) {
	d := models.Json_ids{}
	body, err := ioutil.ReadAll(r.Body)

	defer r.Body.Close()

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	err = json.Unmarshal(body, &d)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	res := []int{}
	i := 0
	for _, id := range d.Ids {
		err = s.Dbpool.QueryRow(context.Background(), "SELECT func_cities_del($1);", id).Scan(&i)
		res = append(res, i)

		if err != nil {
			log.Println("Failed execute func_cities_del: ", err)
		}
	}

	output, err := json.Marshal(models.Json_ids{Ids: res})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(output)

	return
}

// HandleGetCity godoc
// @Summary Get city
// @Description get city
// @Tags cities
// @Produce  json
// @Param id path int true "City by id"
// @Success 200 {array} models.City_count
// @Failure 500
// @Router /cities/{id} [get]
func (s *APG) HandleGetCity(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	i := vars["id"]
	g := models.City{}
	out_arr := []models.City{}

	err := s.Dbpool.QueryRow(context.Background(), "SELECT * from func_city_get($1);", i).Scan(&g.Id, &g.CityName)

	if err != nil && err != pgx.ErrNoRows {
		log.Println("Failed execute from func_city_get: ", err)
	}

	out_arr = append(out_arr, g)
	auth := models.Auth{Create: true, Read: true, Update: true, Delete: true}

	// output, err := json.Marshal(g)
	out_count, err := json.Marshal(models.City_count{Values: out_arr, Count: 1, Auth: auth})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(out_count)

	return
}
