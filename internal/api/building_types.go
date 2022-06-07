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

// HandleBuildingTypes godoc
// @Summary List building types
// @Description get building types
// @Tags building types
// @Produce  json
// @Param page query int false "page number"
// @Param page_size query int false "page size"
// @Param buildingtypename query string false "buildingtypename search pattern"
// @Param ordering query string false "order by {id|buildingtypename}"
// @Param desc query boolean false "descending order {true|false}"
// @Success 200 {array} models.BuildingType_count
// @Failure 500
// @Router /building_types [get]
func (s *APG) HandleBuildingTypes(w http.ResponseWriter, r *http.Request) {
	gs := models.BuildingType{}
	ctx := context.Background()
	out_arr := []models.BuildingType{}

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

	gsn := ""
	gsns, ok := query["buildingtypename"]
	if ok && len(gsns) > 0 {
		//case insensitive
		gsn = strings.ToUpper(gsns[0])
		//quotes
		re := regexp.MustCompile(`'`)
		gsn = string(re.ReplaceAll([]byte(gsn), []byte("''")))
	}

	ord := 1
	ords, ok := query["ordering"]
	if !ok || len(ords) == 0 {
		ord = 1
	} else if ords[0] == "id" {
		ord = 1
	} else if ords[0] == "buildingtypename" {
		ord = 2
	}

	dsc := false
	dscs, ok := query["desc"]
	if ok && len(dscs) > 0 {
		if !(dscs[0] == "0") {
			dsc = true
		}
	}

	rows, err := s.Dbpool.Query(ctx, "SELECT * from func_building_types_get($1,$2,$3,$4,$5);", pg, pgs, gsn, ord, dsc)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&gs.Id, &gs.BuildingTypeName)
		if err != nil {
			log.Println("failed to scan row:", err)
		}

		out_arr = append(out_arr, gs)
	}

	gsc := 0
	err = s.Dbpool.QueryRow(ctx, "SELECT * from func_building_types_cnt($1);", gsn).Scan(&gsc)

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	auth := models.Auth{Create: true, Read: true, Update: true, Delete: true}
	out_count, err := json.Marshal(models.BuildingType_count{Values: out_arr, Count: gsc, Auth: auth})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Write(out_count)

	return
}

// HandleAddBuildingType godoc
// @Summary Add building type
// @Description add building type
// @Tags building types
// @Accept json
// @Produce  json
// @Param a body models.AddBuildingType true "New building type"
// @Success 200 {array} models.Json_id
// @Failure 500
// @Router /building_types_add [post]
func (s *APG) HandleAddBuildingType(w http.ResponseWriter, r *http.Request) {
	a := models.AddBuildingType{}
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
	err = s.Dbpool.QueryRow(context.Background(), "SELECT func_building_types_add($1);", a.BuildingTypeName).Scan(&ai)

	if err != nil {
		log.Println("Failed execute func_building_types_add: ", err)
	}

	output, err := json.Marshal(models.Json_id{Id: ai})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(output)

	return

}

// HandleUpdBuildingType godoc
// @Summary Update building type
// @Description update building type
// @Tags building types
// @Accept json
// @Produce  json
// @Param u body models.BuildingType true "Update building type"
// @Success 200 {array} models.Json_id
// @Failure 500
// @Router /building_types_upd [post]
func (s *APG) HandleUpdBuildingType(w http.ResponseWriter, r *http.Request) {
	u := models.BuildingType{}
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
	err = s.Dbpool.QueryRow(context.Background(), "SELECT func_building_types_upd($1,$2);", u.Id, u.BuildingTypeName).Scan(&ui)

	if err != nil {
		log.Println("Failed execute func_building_types_upd: ", err)
	}

	output, err := json.Marshal(models.Json_id{Id: ui})

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(output)

	return

}

// HandleDelBuildingType godoc
// @Summary Delete building types
// @Description delete building types
// @Tags building types
// @Accept json
// @Produce  json
// @Param d body models.Json_ids true "Delete building types"
// @Success 200 {array} models.Json_ids
// @Failure 500
// @Router /building_types_del [post]
func (s *APG) HandleDelBuildingType(w http.ResponseWriter, r *http.Request) {
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
		err = s.Dbpool.QueryRow(context.Background(), "SELECT func_building_types_del($1);", id).Scan(&i)
		res = append(res, i)

		if err != nil {
			log.Println("Failed execute func_building_types_del: ", err)
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

// HandleGetBuildingType godoc
// @Summary Get building type
// @Description get building type
// @Tags building types
// @Produce  json
// @Param id path int true "Building type by id"
// @Success 200 {array} models.BuildingType_count
// @Failure 500
// @Router /building_types/{id} [get]
func (s *APG) HandleGetBuildingType(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	i := vars["id"]
	g := models.BuildingType{}
	out_arr := []models.BuildingType{}

	err := s.Dbpool.QueryRow(context.Background(), "SELECT * from func_building_type_get($1);", i).Scan(&g.Id, &g.BuildingTypeName)

	if err != nil && err != pgx.ErrNoRows {
		log.Println("Failed execute from func_building_type_get: ", err)
	}

	out_arr = append(out_arr, g)
	auth := models.Auth{Create: true, Read: true, Update: true, Delete: true}

	// output, err := json.Marshal(g)
	out_count, err := json.Marshal(models.BuildingType_count{Values: out_arr, Count: 1, Auth: auth})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(out_count)

	return
}
