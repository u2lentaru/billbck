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

// HandleStreets godoc
// @Summary List streets
// @Description get street list
// @Tags streets
// @Produce  json
// @Param page query int false "page number"
// @Param page_size query int false "page size"
// @Param streetname query string false "streetname search pattern"
// @Param cityid query int false "cityid search pattern"
// @Param ordering query string false "order by {id|streetname}"
// @Param desc query boolean false "descending order {true|false}"
// @Success 200 {object} models.Street_count
// @Failure 500
// @Router /streets [get]
func (s *APG) HandleStreets(w http.ResponseWriter, r *http.Request) {
	gs := models.Street{}
	ctx := context.Background()

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
	gs1s, ok := query["streetname"]
	if ok && len(gs1s) > 0 {
		//case insensitive
		gs1 = strings.ToUpper(gs1s[0])
		//quotes
		re := regexp.MustCompile(`'`)
		gs1 = string(re.ReplaceAll([]byte(gs1), []byte("''")))
	}

	gs2 := 0
	gs2s, ok := query["cityid"]
	if ok && len(gs2s) > 0 {
		t, err := strconv.Atoi(gs2s[0])
		if err == nil {
			gs2 = t
		}
	}

	gsc := 0
	err := s.Dbpool.QueryRow(ctx, "SELECT * from func_streets_cnt($1,$2);", gs1, gs2).Scan(&gsc)

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	out_arr := make([]models.Street, 0,
		func() int {
			if gsc < pgs {
				return gsc
			} else {
				return pgs
			}
		}())

	ord := 1
	ords, ok := query["ordering"]
	if !ok || len(ords) == 0 {
		ord = 1
	} else if ords[0] == "streetname" {
		ord = 2
	}

	dsc := false
	dscs, ok := query["desc"]
	if ok && len(dscs) > 0 {
		if !(dscs[0] == "0") {
			dsc = true
		}
	}

	rows, err := s.Dbpool.Query(ctx, "SELECT * from func_streets_get($1,$2,$3,$4,$5,$6);", pg, pgs, gs1, gs2, ord, dsc)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&gs.Id, &gs.StreetName, &gs.Created, &gs.Closed, &gs.City.CityName, &gs.City.Id)
		if err != nil {
			log.Println("failed to scan row:", err)
		}

		out_arr = append(out_arr, gs)
	}

	auth := models.Auth{Create: true, Read: true, Update: true, Delete: true}
	out_count, err := json.Marshal(models.Street_count{Values: out_arr, Count: gsc, Auth: auth})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Write(out_count)

	return

}

// HandleAddStreet godoc
// @Summary Add street
// @Description add street
// @Tags streets
// @Accept json
// @Produce  json
// @Param a body models.AddStreet true "New street"
// @Success 200 {object} models.Json_id
// @Failure 500
// @Router /streets_add [post]
func (s *APG) HandleAddStreet(w http.ResponseWriter, r *http.Request) {
	a := models.AddStreet{}
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
	err = s.Dbpool.QueryRow(context.Background(), "SELECT func_streets_add($1,$2,$3);", a.StreetName, a.City.Id, a.Created).Scan(&ai)

	if err != nil {
		log.Println("Failed execute func_streets_add: ", err)
	}

	output, err := json.Marshal(models.Json_id{Id: ai})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(output)

	return

}

// HandleUpdStreet godoc
// @Summary Update street
// @Description update street
// @Tags streets
// @Accept json
// @Produce  json
// @Param u body models.Street true "Update street"
// @Success 200 {object} models.Json_id
// @Failure 500
// @Router /streets_upd [post]
func (s *APG) HandleUpdStreet(w http.ResponseWriter, r *http.Request) {
	u := models.Street{}
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
	err = s.Dbpool.QueryRow(context.Background(), "SELECT func_streets_upd($1,$2,$3);", u.Id, u.StreetName, u.Created).Scan(&ui)

	if err != nil {
		log.Println("Failed execute func_streets_upd: ", err)
	}

	output, err := json.Marshal(models.Json_id{Id: ui})

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(output)

	return
}

// HandleDelStreet godoc
// @Summary Delete street
// @Description delete street
// @Tags streets
// @Accept json
// @Produce  json
// @Param d body models.StreetClose true "Delete street"
// @Success 200 {object} models.Json_ids
// @Failure 500
// @Router /streets_del [post]
func (s *APG) HandleDelStreet(w http.ResponseWriter, r *http.Request) {
	d := models.StreetClose{}
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

	i := 0
	err = s.Dbpool.QueryRow(context.Background(), "SELECT func_streets_del($1,$2);", d.Id, d.Close).Scan(&i)

	if err != nil {
		log.Println("Failed execute func_streets_del: ", err)
	}

	output, err := json.Marshal(models.Json_id{Id: i})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(output)

	return

}

// HandleGetStreet godoc
// @Summary Get street
// @Description get street
// @Tags streets
// @Produce  json
// @Param id path int true "Street by id"
// @Success 200 {object} models.Street_count
// @Failure 500
// @Router /streets/{id} [get]
func (s *APG) HandleGetStreet(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	i := vars["id"]
	g := models.Street{}
	out_arr := []models.Street{}

	err := s.Dbpool.QueryRow(context.Background(), "SELECT * from func_street_get($1);", i).Scan(&g.Id, &g.StreetName, &g.Created, &g.Closed,
		&g.City.CityName, &g.City.Id)

	if err != nil && err != pgx.ErrNoRows {
		log.Println("Failed execute from func_street_get: ", err)
	}

	out_arr = append(out_arr, g)
	auth := models.Auth{Create: true, Read: true, Update: true, Delete: true}

	// output, err := json.Marshal(g)
	out_count, err := json.Marshal(models.Street_count{Values: out_arr, Count: 1, Auth: auth})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(out_count)

	return

}
