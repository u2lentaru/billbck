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

// HandleObjTransVolt godoc
// @Summary List objtransvolt
// @Description get objtransvolt list
// @Tags objtransvolt
// @Produce json
// @Param page query int false "page number"
// @Param page_size query int false "page size"
// @Param objectname query string false "objectname search pattern"
// @Param transvoltname query string false "transvoltname search pattern"
// @Param ordering query string false "order by {id|objectname|transvoltname}"
// @Param desc query boolean false "descending order {true|false}"
// @Success 200 {object} models.ObjTransVolt_count
// @Failure 500
// @Router /objtransvolt [get]
func (s *APG) HandleObjTransVolt(w http.ResponseWriter, r *http.Request) {
	gs := models.ObjTransVolt{}
	ctx := context.Background()
	out_arr := []models.ObjTransVolt{}

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
	gs1s, ok := query["objectname"]
	if ok && len(gs1s) > 0 {
		//case insensitive
		gs1 = strings.ToUpper(gs1s[0])
		//quotes
		re := regexp.MustCompile(`'`)
		gs1 = string(re.ReplaceAll([]byte(gs1), []byte("''")))
	}

	gs2 := ""
	gs2s, ok := query["transvoltname"]
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
		ord = 2
	} else if ords[0] == "transvoltname" {
		ord = 3
	}

	dsc := false
	dscs, ok := query["desc"]
	if ok && len(dscs) > 0 {
		if !(dscs[0] == "0") {
			dsc = true
		}
	}

	rows, err := s.Dbpool.Query(ctx, "SELECT * from func_obj_trans_volt_get($1,$2,$3,$4,$5,$6);", pg, pgs, gs1, gs2, ord, dsc)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&gs.Id, &gs.ObjId, &gs.ObjTypeId, &gs.TransVolt.Id, &gs.Startdate, &gs.Enddate, &gs.ObjName,
			&gs.TransVolt.TransVoltName, &gs.TransVolt.TransType.Id, &gs.TransVolt.CheckDate, &gs.TransVolt.NextCheckDate,
			&gs.TransVolt.ProdDate, &gs.TransVolt.Serial1, &gs.TransVolt.Serial2, &gs.TransVolt.Serial3,
			&gs.TransVolt.TransType.TransTypeName, &gs.TransVolt.TransType.Ratio, &gs.TransVolt.TransType.Class,
			&gs.TransVolt.TransType.MaxCurr, &gs.TransVolt.TransType.NomCurr)
		if err != nil {
			log.Println("failed to scan row:", err)
		}

		out_arr = append(out_arr, gs)
	}

	gsc := 0
	err = s.Dbpool.QueryRow(ctx, "SELECT * from func_obj_trans_volt_cnt($1,$2);", gs1, gs2).Scan(&gsc)

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	auth := models.Auth{Create: true, Read: true, Update: true, Delete: true}
	out_count, err := json.Marshal(models.ObjTransVolt_count{Values: out_arr, Count: gsc, Auth: auth})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Write(out_count)

	return

}

// HandleAddObjTransVolt godoc
// @Summary Add objtransvolt
// @Description add objtransvolt
// @Tags objtransvolt
// @Accept json
// @Produce  json
// @Param a body models.AddObjTransVolt true "New objtransvolt. Significant params: ObjId, ObjTypeId, TransVolt.Id, Startdate"
// @Success 200 {object} models.Json_id
// @Failure 500
// @Router /objtransvolt_add [post]
func (s *APG) HandleAddObjTransVolt(w http.ResponseWriter, r *http.Request) {
	a := models.AddObjTransVolt{}
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
	err = s.Dbpool.QueryRow(context.Background(), "SELECT func_obj_trans_volt_add($1,$2,$3,$4);", a.ObjId, a.ObjTypeId,
		a.TransVolt.Id, a.Startdate).Scan(&ai)

	if err != nil {
		log.Println("Failed execute func_obj_trans_volt_add: ", err)
	}

	output, err := json.Marshal(models.Json_id{Id: ai})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(output)

	return
}

// HandleUpdObjTransVolt godoc
// @Summary Update objtransvolt
// @Description update objtransvolt
// @Tags objtransvolt
// @Accept json
// @Produce  json
// @Param u body models.ObjTransVolt true "Update objtransvolt. Significant params: Id, ObjId, ObjTypeId, TransVolt.Id, Startdate, Enddate"
// @Success 200 {object} models.Json_id
// @Failure 500
// @Router /objtransvolt_upd [post]
func (s *APG) HandleUpdObjTransVolt(w http.ResponseWriter, r *http.Request) {
	u := models.ObjTransVolt{}
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
	err = s.Dbpool.QueryRow(context.Background(), "SELECT func_obj_trans_volt_upd($1,$2,$3,$4,$5,$6);", u.Id, u.ObjId, u.ObjTypeId,
		u.TransVolt.Id, u.Startdate, u.Enddate).Scan(&ui)

	if err != nil {
		log.Println("Failed execute func_obj_trans_volt_upd: ", err)
	}

	output, err := json.Marshal(models.Json_id{Id: ui})

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(output)

	return

}

// HandleDelObjTransVolt godoc
// @Summary Delete objtransvolts
// @Description delete objtransvolts
// @Tags objtransvolt
// @Accept json
// @Produce  json
// @Param d body models.Json_ids true "Delete objtransvolts"
// @Success 200 {object} models.Json_ids
// @Failure 500
// @Router /objtransvolt_del [post]
func (s *APG) HandleDelObjTransVolt(w http.ResponseWriter, r *http.Request) {
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
		err = s.Dbpool.QueryRow(context.Background(), "SELECT func_obj_trans_volt_del($1);", id).Scan(&i)
		res = append(res, i)

		if err != nil {
			log.Println("Failed execute func_obj_trans_volt_del: ", err)
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

// HandleGetObjTransVolt godoc
// @Summary Get objtransvolt
// @Description get objtransvolt
// @Tags objtransvolt
// @Produce  json
// @Param id path int true "Objtransvolt by id"
// @Success 200 {object} models.ObjTransVolt_count
// @Failure 500
// @Router /objtransvolt/{id} [get]
func (s *APG) HandleGetObjTransVolt(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	i := vars["id"]
	g := models.ObjTransVolt{}
	out_arr := []models.ObjTransVolt{}

	err := s.Dbpool.QueryRow(context.Background(), "SELECT * from func_obj_trans_volt_getbyid($1);", i).Scan(&g.Id, &g.ObjId,
		&g.ObjTypeId, &g.TransVolt.Id, &g.Startdate, &g.Enddate, &g.ObjName, &g.TransVolt.TransVoltName, &g.TransVolt.TransType.Id,
		&g.TransVolt.CheckDate, &g.TransVolt.NextCheckDate, &g.TransVolt.ProdDate, &g.TransVolt.Serial1, &g.TransVolt.Serial2,
		&g.TransVolt.Serial3, &g.TransVolt.TransType.TransTypeName, &g.TransVolt.TransType.Ratio, &g.TransVolt.TransType.Class,
		&g.TransVolt.TransType.MaxCurr, &g.TransVolt.TransType.NomCurr)

	if err != nil && err != pgx.ErrNoRows {
		log.Println("Failed execute from func_obj_trans_volt_getbyid: ", err)
	}

	out_arr = append(out_arr, g)
	auth := models.Auth{Create: true, Read: true, Update: true, Delete: true}

	out_count, err := json.Marshal(models.ObjTransVolt_count{Values: out_arr, Count: 1, Auth: auth})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Write(out_count)

	return
}

// HandleObjTransVoltByObj godoc
// @Summary Get objtransvolt by object
// @Description get objtransvolt by object
// @Tags objtransvolt
// @Produce  json
// @Param objid query string false "obj&tgu id"
// @Param tid query string false "obj&tgu type id (obj - type = 0, tgu - type > 0)"
// @Success 200 {object} models.ObjTransVolt_count
// @Failure 500
// @Router /objtransvolt_obj [get]
func (s *APG) HandleObjTransVoltByObj(w http.ResponseWriter, r *http.Request) {
	g := models.ObjTransVolt{}
	out_arr := []models.ObjTransVolt{}

	query := r.URL.Query()

	gs1 := "0"
	gs1s, ok := query["objid"]
	if ok && len(gs1s) > 0 {
		_, err := strconv.Atoi(gs1s[0])
		if err == nil {
			gs1 = gs1s[0]
		}
	}

	gs2 := "0"
	gs2s, ok := query["tid"]
	if ok && len(gs2s) > 0 {
		_, err := strconv.Atoi(gs2s[0])
		if err == nil {
			gs2 = gs2s[0]
		}
	}

	err := s.Dbpool.QueryRow(context.Background(), "SELECT * from func_obj_trans_volt_getbyobj($1,$2);", gs1, gs2).Scan(&g.Id, &g.ObjId,
		&g.ObjTypeId, &g.TransVolt.Id, &g.Startdate, &g.Enddate, &g.ObjName, &g.TransVolt.TransVoltName, &g.TransVolt.TransType.Id,
		&g.TransVolt.CheckDate, &g.TransVolt.NextCheckDate, &g.TransVolt.ProdDate, &g.TransVolt.Serial1, &g.TransVolt.Serial2,
		&g.TransVolt.Serial3, &g.TransVolt.TransType.TransTypeName, &g.TransVolt.TransType.Ratio, &g.TransVolt.TransType.Class,
		&g.TransVolt.TransType.MaxCurr, &g.TransVolt.TransType.NomCurr)

	if err != nil && err != pgx.ErrNoRows {
		log.Println("Failed execute from func_obj_trans_volt_getbyobj: ", err)
	}

	out_arr = append(out_arr, g)
	auth := models.Auth{Create: true, Read: true, Update: true, Delete: true}

	out_count, err := json.Marshal(models.ObjTransVolt_count{Values: out_arr, Count: 1, Auth: auth})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Write(out_count)

	return

}
