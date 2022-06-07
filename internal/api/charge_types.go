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

// HandleChargeTypes godoc
// @Summary List chargetypes
// @Description get chargetype list
// @Tags chargetypes
// @Produce  json
// @Param page query int false "page number"
// @Param page_size query int false "page size"
// @Param chargetypename query string false "chargetypename search pattern"
// @Param ordering query string false "order by {id|chargetypename}"
// @Param desc query boolean false "descending order {true|false}"
// @Success 200 {array} models.ChargeType_count
// @Failure 500
// @Router /chargetypes [get]
func (s *APG) HandleChargeTypes(w http.ResponseWriter, r *http.Request) {
	gs := models.ChargeType{}
	ctx := context.Background()
	out_arr := []models.ChargeType{}

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
	gs1s, ok := query["chargetypename"]
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
	} else if ords[0] == "chargetypename" {
		ord = 2
	}

	dsc := false
	dscs, ok := query["desc"]
	if ok && len(dscs) > 0 {
		if !(dscs[0] == "0") {
			dsc = true
		}
	}

	rows, err := s.Dbpool.Query(ctx, "SELECT * from func_charge_types_get($1,$2,$3,$4,$5);", pg, pgs, gs1, ord, dsc)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&gs.Id, &gs.ChargeTypeName)
		if err != nil {
			log.Println("failed to scan row:", err)
		}

		out_arr = append(out_arr, gs)
	}

	gsc := 0
	err = s.Dbpool.QueryRow(ctx, "SELECT * from func_charge_types_cnt($1);", gs1).Scan(&gsc)

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	auth := models.Auth{Create: true, Read: true, Update: true, Delete: true}
	out_count, err := json.Marshal(models.ChargeType_count{Values: out_arr, Count: gsc, Auth: auth})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Write(out_count)

	return

}

// HandleAddChargeType godoc
// @Summary Add chargetype
// @Description add chargetype
// @Tags chargetypes
// @Accept json
// @Produce  json
// @Param a body models.AddChargeType true "New chargetype"
// @Success 200 {array} models.Json_id
// @Failure 500
// @Router /chargetypes_add [post]
func (s *APG) HandleAddChargeType(w http.ResponseWriter, r *http.Request) {
	a := models.AddChargeType{}
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
	err = s.Dbpool.QueryRow(context.Background(), "SELECT func_charge_types_add($1);", a.ChargeTypeName).Scan(&ai)

	if err != nil {
		log.Println("Failed execute func_charge_types_add: ", err)
	}

	output, err := json.Marshal(models.Json_id{Id: ai})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(output)

	return

}

// HandleUpdChargeType godoc
// @Summary Update chargetype
// @Description update chargetype
// @Tags chargetypes
// @Accept json
// @Produce  json
// @Param u body models.ChargeType true "Update chargetype"
// @Success 200 {array} models.Json_id
// @Failure 500
// @Router /chargetypes_upd [post]
func (s *APG) HandleUpdChargeType(w http.ResponseWriter, r *http.Request) {
	u := models.ChargeType{}
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
	err = s.Dbpool.QueryRow(context.Background(), "SELECT func_charge_types_upd($1,$2);", u.Id, u.ChargeTypeName).Scan(&ui)

	if err != nil {
		log.Println("Failed execute func_charge_types_upd: ", err)
	}

	output, err := json.Marshal(models.Json_id{Id: ui})

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(output)

	return
}

// HandleDelChargeType godoc
// @Summary Delete chargetypes
// @Description delete chargetypes
// @Tags chargetypes
// @Accept json
// @Produce  json
// @Param d body models.Json_ids true "Delete chargetypes"
// @Success 200 {array} models.Json_ids
// @Failure 500
// @Router /chargetypes_del [post]
func (s *APG) HandleDelChargeType(w http.ResponseWriter, r *http.Request) {
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
		err = s.Dbpool.QueryRow(context.Background(), "SELECT func_charge_types_del($1);", id).Scan(&i)
		res = append(res, i)

		if err != nil {
			log.Println("Failed execute func_charge_types_del: ", err)
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

// HandleGetChargeType godoc
// @Summary Get chargetype
// @Description get chargetype
// @Tags chargetypes
// @Produce  json
// @Param id path int true "ChargeType by id"
// @Success 200 {array} models.ChargeType_count
// @Failure 500
// @Router /chargetypes/{id} [get]
func (s *APG) HandleGetChargeType(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	i := vars["id"]
	g := models.ChargeType{}
	out_arr := []models.ChargeType{}

	err := s.Dbpool.QueryRow(context.Background(), "SELECT * from func_charge_type_get($1);", i).Scan(&g.Id, &g.ChargeTypeName)

	if err != nil && err != pgx.ErrNoRows {
		log.Println("Failed execute from func_charge_type_get: ", err)
	}

	out_arr = append(out_arr, g)
	auth := models.Auth{Create: true, Read: true, Update: true, Delete: true}
	// output, err := json.Marshal(g)
	out_count, err := json.Marshal(models.ChargeType_count{Values: out_arr, Count: 1, Auth: auth})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(out_count)

	return

}
