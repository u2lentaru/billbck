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

// HandleEquipmentTypes godoc
// @Summary List equipmenttypes
// @Description get equipmenttype list
// @Tags equipment types
// @Produce  json
// @Param page query int false "page number"
// @Param page_size query int false "page size"
// @Param equipmenttypename query string false "equipmenttypename search pattern"
// @Param ordering query string false "order by {id|equipmenttypename}"
// @Param desc query boolean false "descending order {true|false}"
// @Success 200 {object} models.EquipmentType_count
// @Failure 500
// @Router /equipmenttypes [get]
func (s *APG) HandleEquipmentTypes(w http.ResponseWriter, r *http.Request) {
	gs := models.EquipmentType{}
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
	gs1s, ok := query["equipmenttypename"]
	if ok && len(gs1s) > 0 {
		//case insensitive
		gs1 = strings.ToUpper(gs1s[0])
		//quotes
		re := regexp.MustCompile(`'`)
		gs1 = string(re.ReplaceAll([]byte(gs1), []byte("''")))
	}

	gsc := 0
	err := s.Dbpool.QueryRow(ctx, "SELECT * from func_equipment_types_cnt($1);", gs1).Scan(&gsc)

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	out_arr := make([]models.EquipmentType, 0,
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
	} else if ords[0] == "equipmenttypename" {
		ord = 2
	}

	dsc := false
	dscs, ok := query["desc"]
	if ok && len(dscs) > 0 {
		if !(dscs[0] == "0") {
			dsc = true
		}
	}

	rows, err := s.Dbpool.Query(ctx, "SELECT * from func_equipment_types_get($1,$2,$3,$4,$5);", pg, pgs, gs1, ord, dsc)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&gs.Id, &gs.EquipmentTypeName, &gs.EquipmentTypePower)
		if err != nil {
			log.Println("failed to scan row:", err)
			http.Error(w, err.Error(), 500)
		}

		out_arr = append(out_arr, gs)
	}

	auth := models.Auth{Create: true, Read: true, Update: true, Delete: true}
	out_count, err := json.Marshal(models.EquipmentType_count{Values: out_arr, Count: gsc, Auth: auth})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Write(out_count)

	return
}

// HandleAddEquipmentType godoc
// @Summary Add equipmenttype
// @Description add equipmenttype
// @Tags equipment types
// @Accept json
// @Produce  json
// @Param a body models.AddEquipmentType true "New equipmenttype. Significant params: EquipmentTypeName, EquipmentTypePower"
// @Success 200 {object} models.Json_id
// @Failure 500
// @Router /equipmenttypes_add [post]
func (s *APG) HandleAddEquipmentType(w http.ResponseWriter, r *http.Request) {
	a := models.AddEquipmentType{}
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
	err = s.Dbpool.QueryRow(context.Background(), "SELECT func_equipment_types_add($1,$2);", a.EquipmentTypeName,
		a.EquipmentTypePower).Scan(&ai)

	if err != nil {
		log.Println("Failed execute func_equipment_types_add: ", err)
		http.Error(w, err.Error(), 500)
	}

	output, err := json.Marshal(models.Json_id{Id: ai})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(output)

	return
}

// HandleUpdEquipmentType godoc
// @Summary Update equipmenttype
// @Description update equipmenttype
// @Tags equipment types
// @Accept json
// @Produce  json
// @Param u body models.EquipmentType true "Update equipmenttype. Significant params: Id, EquipmentTypeName, EquipmentTypePower"
// @Success 200 {object} models.Json_id
// @Failure 500
// @Router /equipmenttypes_upd [post]
func (s *APG) HandleUpdEquipmentType(w http.ResponseWriter, r *http.Request) {
	u := models.EquipmentType{}
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
	err = s.Dbpool.QueryRow(context.Background(), "SELECT func_equipment_types_upd($1,$2,$3);", u.Id, u.EquipmentTypeName,
		u.EquipmentTypePower).Scan(&ui)

	if err != nil {
		log.Println("Failed execute func_equipment_types_upd: ", err)
		http.Error(w, err.Error(), 500)
	}

	output, err := json.Marshal(models.Json_id{Id: ui})

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(output)

	return

}

// HandleDelEquipmentType godoc
// @Summary Delete equipmenttypes
// @Description delete equipmenttypes
// @Tags equipment types
// @Accept json
// @Produce  json
// @Param d body models.Json_ids true "Delete equipmenttypes"
// @Success 200 {object} models.Json_ids
// @Failure 500
// @Router /equipmenttypes_del [post]
func (s *APG) HandleDelEquipmentType(w http.ResponseWriter, r *http.Request) {
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
		err = s.Dbpool.QueryRow(context.Background(), "SELECT func_equipment_types_del($1);", id).Scan(&i)
		res = append(res, i)

		if err != nil {
			log.Println("Failed execute func_equipment_types_del: ", err)
			http.Error(w, err.Error(), 500)
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

// HandleGetEquipmentType godoc
// @Summary Get equipmenttype
// @Description get equipmenttype
// @Tags equipment types
// @Produce  json
// @Param id path int true "EquipmentType by id"
// @Success 200 {object} models.EquipmentType_count
// @Failure 500
// @Router /equipmenttypes/{id} [get]
func (s *APG) HandleGetEquipmentType(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	i := vars["id"]
	g := models.EquipmentType{}
	out_arr := []models.EquipmentType{}

	err := s.Dbpool.QueryRow(context.Background(), "SELECT * from func_equipment_type_get($1);", i).Scan(&g.Id, &g.EquipmentTypeName,
		&g.EquipmentTypePower)

	if err != nil && err != pgx.ErrNoRows {
		log.Println("Failed execute from func_equipment_type_get: ", err)
		http.Error(w, err.Error(), 500)
	}

	out_arr = append(out_arr, g)
	auth := models.Auth{Create: true, Read: true, Update: true, Delete: true}

	// output, err := json.Marshal(g)
	out_count, err := json.Marshal(models.EquipmentType_count{Values: out_arr, Count: 1, Auth: auth})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(out_count)

	return

}
