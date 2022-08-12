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
	"github.com/u2lentaru/billbck/internal/utils"
)

// HandleEquipment godoc
// @Summary List equipment
// @Description get equipment list
// @Tags equipment
// @Produce  json
// @Param page query int false "page number"
// @Param page_size query int false "page size"
// @Param objectid query int false "objectid"
// @Param objectname query string false "name search pattern"
// @Param ordering query string false "order by {id|objectname}"
// @Param desc query boolean false "descending order {true|false}"
// @Success 200 {object} models.Equipment_count
// @Failure 500
// @Router /equipment [get]
func (s *APG) HandleEquipment(w http.ResponseWriter, r *http.Request) {
	gs := models.Equipment{}
	ctx := context.Background()
	out_arr := []models.Equipment{}

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

	gs1 := 0
	gs1s, ok := query["objectid"]
	if ok && len(gs1s) > 0 {
		t, err := strconv.Atoi(gs1s[0])
		if err == nil {
			gs1 = t
		}
	}

	gs2 := ""
	gs2s, ok := query["objectname"]
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
	} else if ords[0] == "objectname" {
		ord = 8
	}

	dsc := false
	dscs, ok := query["desc"]
	if ok && len(dscs) > 0 {
		if !(dscs[0] == "0") {
			dsc = true
		}
	}

	rows, err := s.Dbpool.Query(ctx, "SELECT * from func_equipment_get($1,$2,$3,$4,$5,$6);", pg, pgs, gs2, utils.NullableInt(int32(gs1)), ord, dsc)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&gs.Id, &gs.EquipmentType.Id, &gs.Object.Id, &gs.Qty, &gs.WorkingHours, &gs.EquipmentType.EquipmentTypeName,
			&gs.EquipmentType.EquipmentTypePower, &gs.Object.ObjectName)
		if err != nil {
			log.Println("failed to scan row:", err)
			http.Error(w, err.Error(), 500)
		}

		out_arr = append(out_arr, gs)
	}

	gsc := 0
	err = s.Dbpool.QueryRow(ctx, "SELECT * from func_equipment_cnt($1,$2);", gs2, utils.NullableInt(int32(gs1))).Scan(&gsc)

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	auth := models.Auth{Create: true, Read: true, Update: true, Delete: true}
	out_count, err := json.Marshal(models.Equipment_count{Values: out_arr, Count: gsc, Auth: auth})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Write(out_count)

	return
}

// HandleAddEquipment godoc
// @Summary Add equipment
// @Description add equipment
// @Tags equipment
// @Accept json
// @Produce  json
// @Param a body models.AddEquipment true "New equipment. Significant params: EquipmentType.Id, Object.Id, Qty, WorkingHours"
// @Success 200 {object} models.Json_id
// @Failure 500
// @Router /equipment_add [post]
func (s *APG) HandleAddEquipment(w http.ResponseWriter, r *http.Request) {
	a := models.AddEquipment{}
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
	err = s.Dbpool.QueryRow(context.Background(), "SELECT func_equipment_add($1,$2,$3,$4);", a.EquipmentType.Id, a.Object.Id,
		a.Qty, a.WorkingHours).Scan(&ai)

	if err != nil {
		log.Println("Failed execute func_equipment_add: ", err)
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

// HandleUpdEquipment godoc
// @Summary Update equipment
// @Description update equipment
// @Tags equipment
// @Accept json
// @Produce  json
// @Param u body models.Equipment true "Update equipment. Significant params: Id, EquipmentType.Id, Object.Id, Qty, WorkingHours"
// @Success 200 {object} models.Json_id
// @Failure 500
// @Router /equipment_upd [post]
func (s *APG) HandleUpdEquipment(w http.ResponseWriter, r *http.Request) {
	u := models.Equipment{}
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
	err = s.Dbpool.QueryRow(context.Background(), "SELECT func_equipment_upd($1,$2,$3,$4,$5);", u.Id, u.EquipmentType.Id, u.Object.Id,
		u.Qty, u.WorkingHours).Scan(&ui)

	if err != nil {
		log.Println("Failed execute func_equipment_upd: ", err)
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

// HandleDelEquipment godoc
// @Summary Delete equipment
// @Description delete equipment
// @Tags equipment
// @Accept json
// @Produce  json
// @Param d body models.Json_ids true "Delete equipment"
// @Success 200 {object} models.Json_ids
// @Failure 500
// @Router /equipment_del [post]
func (s *APG) HandleDelEquipment(w http.ResponseWriter, r *http.Request) {
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
		err = s.Dbpool.QueryRow(context.Background(), "SELECT func_equipment_del($1);", id).Scan(&i)
		res = append(res, i)

		if err != nil {
			log.Println("Failed execute func_equipment_del: ", err)
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

// HandleGetEquipment godoc
// @Summary Get equipment
// @Description get equipment
// @Tags equipment
// @Produce  json
// @Param id path int true "Equipment by id"
// @Success 200 {object} models.Equipment_count
// @Failure 500
// @Router /equipment/{id} [get]
func (s *APG) HandleGetEquipment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	i := vars["id"]
	g := models.Equipment{}
	out_arr := []models.Equipment{}

	err := s.Dbpool.QueryRow(context.Background(), "SELECT * from func_equipment_getbyid($1);", i).Scan(&g.Id, &g.EquipmentType.Id, &g.Object.Id,
		&g.Qty, &g.WorkingHours, &g.EquipmentType.EquipmentTypeName, &g.EquipmentType.EquipmentTypePower, &g.Object.ObjectName)

	if err != nil && err != pgx.ErrNoRows {
		log.Println("Failed execute from func_equipment_getbyid: ", err)
		http.Error(w, err.Error(), 500)
	}

	out_arr = append(out_arr, g)
	auth := models.Auth{Create: true, Read: true, Update: true, Delete: true}

	// output, err := json.Marshal(g)
	out_count, err := json.Marshal(models.Equipment_count{Values: out_arr, Count: 1, Auth: auth})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(out_count)

	return
}

// HandleAddEquipmentList godoc
// @Summary Add equipment list
// @Description add equipment list
// @Tags equipment
// @Accept json
// @Produce  json
// @Param al body models.Equipment_count true "Add equipment list. Old equipment delete by first value Object.Id. Significant params: EquipmentType.Id, Object.Id, Qty, WorkingHours"
// @Success 200 {object} models.Json_ids
// @Failure 500
// @Router /equipment_addlist [post]
func (s *APG) HandleAddEquipmentList(w http.ResponseWriter, r *http.Request) {
	al := models.Equipment_count{}
	body, err := ioutil.ReadAll(r.Body)

	defer r.Body.Close()

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	err = json.Unmarshal(body, &al)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	res := []int{}
	i := 0
	first_value := true

	for _, a := range al.Values {

		if first_value {
			err = s.Dbpool.QueryRow(context.Background(), "SELECT func_equipment_delbyobj($1);", a.Object.Id).Scan(&i)

			if err != nil {
				log.Println("Failed execute func_equipment_delbyobj: ", err)
				http.Error(w, err.Error(), 500)
			}

			first_value = false
			i = 0
		}

		err = s.Dbpool.QueryRow(context.Background(), "SELECT func_equipment_add($1,$2,$3,$4);", a.EquipmentType.Id, a.Object.Id,
			a.Qty, a.WorkingHours).Scan(&i)
		res = append(res, i)

		if err != nil {
			log.Println("Failed execute func_equipment_add: ", err)
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

// HandleDelObjEquipment godoc
// @Summary Delete equipment by object id
// @Description delete equipment by object id
// @Tags equipment
// @Accept json
// @Produce  json
// @Param d body models.Json_ids true "Delete equipment by object id"
// @Success 200 {object} models.Json_ids
// @Failure 500
// @Router /equipment_delobj [post]
func (s *APG) HandleDelObjEquipment(w http.ResponseWriter, r *http.Request) {
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
		err = s.Dbpool.QueryRow(context.Background(), "SELECT func_equipment_delbyobj($1);", id).Scan(&i)
		res = append(res, i)

		if err != nil {
			log.Println("Failed execute func_equipment_delbyobj: ", err)
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
