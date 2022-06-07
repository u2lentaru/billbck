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

// HandleAreas godoc
// @Summary List areas
// @Description get areas
// @Tags areas
// @Produce  json
// @Param page query int false "page number"
// @Param page_size query int false "page size"
// @Param areanumber query string false "areanumber search pattern"
// @Param areaname query string false "areaname search pattern"
// @Param ordering query string false "order by {areanumber|areaname}"
// @Param desc query boolean false "descending order {true|false}"
// @Success 200 {array} models.Area_count
// @Failure 500
// @Router /areas [get]
func (s *APG) HandleAreas(w http.ResponseWriter, r *http.Request) {
	gs := models.Area{}
	ctx := context.Background()
	out_arr := []models.Area{}

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
	gs1s, ok := query["areanumber"]
	if ok && len(gs1s) > 0 {
		//case insensitive
		gs1 = strings.ToUpper(gs1s[0])
		//quotes
		re := regexp.MustCompile(`'`)
		gs1 = string(re.ReplaceAll([]byte(gs1), []byte("''")))
	}

	gs2 := ""
	gs2s, ok := query["areaname"]
	if ok && len(gs2s) > 0 {
		//case insensitive
		gs2 = strings.ToUpper(gs2s[0])
		//quotes
		re := regexp.MustCompile(`'`)
		gs2 = string(re.ReplaceAll([]byte(gs2), []byte("''")))
	}

	ord := 1
	ords, ok := query["ordering"]
	if !ok || len(ords) == 0 {
		ord = 1
	} else if ords[0] == "areanumber" {
		ord = 2
	} else if ords[0] == "areaname" {
		ord = 3
	}

	dsc := false
	dscs, ok := query["desc"]
	if ok && len(dscs) > 0 {
		if !(dscs[0] == "0") {
			dsc = true
		}
	}

	rows, err := s.Dbpool.Query(ctx, "SELECT * from func_areas_get($1,$2,$3,$4,$5,$6);", pg, pgs, gs1, gs2, ord, dsc)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&gs.Id, &gs.AreaNumber, &gs.AreaName)
		if err != nil {
			log.Println("failed to scan row:", err)
		}

		out_arr = append(out_arr, gs)
	}

	gsc := 0
	err = s.Dbpool.QueryRow(ctx, "SELECT * from func_areas_cnt($1,$2);", gs1, gs2).Scan(&gsc)

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	auth := models.Auth{Create: true, Read: true, Update: true, Delete: true}
	out_count, err := json.Marshal(models.Area_count{Values: out_arr, Count: gsc, Auth: auth})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Write(out_count)

	return
}

// HandleAddArea godoc
// @Summary Add area
// @Description add area
// @Tags areas
// @Accept json
// @Produce  json
// @Param a body models.AddArea true "New area"
// @Success 200 {array} models.Json_id
// @Failure 500
// @Router /areas_add [post]
func (s *APG) HandleAddArea(w http.ResponseWriter, r *http.Request) {
	a := models.AddArea{}
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
	err = s.Dbpool.QueryRow(context.Background(), "SELECT func_areas_add($1,$2);", a.AreaNumber, a.AreaName).Scan(&ai)

	if err != nil {
		log.Println("Failed execute func_areas_add: ", err)
	}

	output, err := json.Marshal(models.Json_id{Id: ai})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(output)

	return

}

// HandleUpdArea godoc
// @Summary Update area
// @Description update area
// @Tags areas
// @Accept json
// @Produce  json
// @Param u body models.Area true "Update area"
// @Success 200 {array} models.Json_id
// @Failure 500
// @Router /areas_upd [post]
func (s *APG) HandleUpdArea(w http.ResponseWriter, r *http.Request) {
	u := models.Area{}
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
	err = s.Dbpool.QueryRow(context.Background(), "SELECT func_areas_upd($1,$2,$3);", u.Id, u.AreaNumber, u.AreaName).Scan(&ui)

	if err != nil {
		log.Println("Failed execute func_areas_upd: ", err)
	}

	output, err := json.Marshal(models.Json_id{Id: ui})

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(output)

	return
}

// HandleDelArea godoc
// @Summary Delete areas
// @Description delete areas
// @Tags areas
// @Accept json
// @Produce  json
// @Param d body models.Json_ids true "Delete areas"
// @Success 200 {array} models.Json_ids
// @Failure 500
// @Router /areas_del [post]
func (s *APG) HandleDelArea(w http.ResponseWriter, r *http.Request) {
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
		err = s.Dbpool.QueryRow(context.Background(), "SELECT func_areas_del($1);", id).Scan(&i)
		res = append(res, i)

		if err != nil {
			log.Println("Failed execute func_areas_del: ", err)
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

// HandleGetArea godoc
// @Summary Get area
// @Description get area
// @Tags areas
// @Produce  json
// @Param id path int true "Area by id"
// @Success 200 {array} models.Area_count
// @Failure 500
// @Router /areas/{id} [get]
func (s *APG) HandleGetArea(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	i := vars["id"]
	g := models.Area{}
	out_arr := []models.Area{}

	err := s.Dbpool.QueryRow(context.Background(), "SELECT * from func_area_get($1);", i).Scan(&g.Id, &g.AreaNumber, &g.AreaName)

	if err != nil && err != pgx.ErrNoRows {
		log.Println("Failed execute from func_area_get: ", err)
	}

	out_arr = append(out_arr, g)
	auth := models.Auth{Create: true, Read: true, Update: true, Delete: true}

	// output, err := json.Marshal(g)
	out_count, err := json.Marshal(models.Area_count{Values: out_arr, Count: 1, Auth: auth})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(out_count)

	return
}
