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

// HandleObjTransCurr godoc
// @Summary List objtranscurr
// @Description get objtranscurr list
// @Tags objtranscurr
// @Produce json
// @Param page query int false "page number"
// @Param page_size query int false "page size"
// @Param objectname query string false "objectname search pattern"
// @Param transcurrname query string false "transcurrname search pattern"
// @Param ordering query string false "order by {id|objectname|transcurrname}"
// @Param desc query boolean false "descending order {true|false}"
// @Success 200 {object} models.ObjTransCurr_count
// @Failure 500
// @Router /objtranscurr [get]
func (s *APG) HandleObjTransCurr(w http.ResponseWriter, r *http.Request) {
	gs := models.ObjTransCurr{}
	ctx := context.Background()
	out_arr := []models.ObjTransCurr{}

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
	gs2s, ok := query["transcurrname"]
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
	} else if ords[0] == "transcurrname" {
		ord = 3
	}

	dsc := false
	dscs, ok := query["desc"]
	if ok && len(dscs) > 0 {
		if !(dscs[0] == "0") {
			dsc = true
		}
	}

	rows, err := s.Dbpool.Query(ctx, "SELECT * from func_obj_trans_curr_get($1,$2,$3,$4,$5,$6);", pg, pgs, gs1, gs2, ord, dsc)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&gs.Id, &gs.ObjId, &gs.ObjTypeId, &gs.TransCurr.Id, &gs.Startdate, &gs.Enddate, &gs.ObjName,
			&gs.TransCurr.TransCurrName, &gs.TransCurr.TransType.Id, &gs.TransCurr.CheckDate, &gs.TransCurr.NextCheckDate,
			&gs.TransCurr.ProdDate, &gs.TransCurr.Serial1, &gs.TransCurr.Serial2, &gs.TransCurr.Serial3,
			&gs.TransCurr.TransType.TransTypeName, &gs.TransCurr.TransType.Ratio, &gs.TransCurr.TransType.Class,
			&gs.TransCurr.TransType.MaxCurr, &gs.TransCurr.TransType.NomCurr)
		if err != nil {
			log.Println("failed to scan row:", err)
		}

		out_arr = append(out_arr, gs)
	}

	gsc := 0
	err = s.Dbpool.QueryRow(ctx, "SELECT * from func_obj_trans_curr_cnt($1,$2);", gs1, gs2).Scan(&gsc)

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	auth := models.Auth{Create: true, Read: true, Update: true, Delete: true}
	out_count, err := json.Marshal(models.ObjTransCurr_count{Values: out_arr, Count: gsc, Auth: auth})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Write(out_count)

	return

}

// HandleAddObjTransCurr godoc
// @Summary Add objtranscurr
// @Description add objtranscurr
// @Tags objtranscurr
// @Accept json
// @Produce  json
// @Param a body models.AddObjTransCurr true "New objtranscurr. Significant params: ObjId, ObjTypeId, TransCurr.Id, Startdate"
// @Success 200 {object} models.Json_id
// @Failure 500
// @Router /objtranscurr_add [post]
func (s *APG) HandleAddObjTransCurr(w http.ResponseWriter, r *http.Request) {
	a := models.AddObjTransCurr{}
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
	err = s.Dbpool.QueryRow(context.Background(), "SELECT func_obj_trans_curr_add($1,$2,$3,$4);", a.ObjId, a.ObjTypeId,
		a.TransCurr.Id, a.Startdate).Scan(&ai)

	if err != nil {
		log.Println("Failed execute func_obj_trans_curr_add: ", err)
	}

	output, err := json.Marshal(models.Json_id{Id: ai})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(output)

	return
}

// HandleUpdObjTransCurr godoc
// @Summary Update objtranscurr
// @Description update objtranscurr
// @Tags objtranscurr
// @Accept json
// @Produce  json
// @Param u body models.ObjTransCurr true "Update objtranscurr. Significant params: Id, ObjId, ObjTypeId, TransCurr.Id, Startdate, Enddate"
// @Success 200 {object} models.Json_id
// @Failure 500
// @Router /objtranscurr_upd [post]
func (s *APG) HandleUpdObjTransCurr(w http.ResponseWriter, r *http.Request) {
	u := models.ObjTransCurr{}
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
	err = s.Dbpool.QueryRow(context.Background(), "SELECT func_obj_trans_curr_upd($1,$2,$3,$4,$5,$6);", u.Id, u.ObjId, u.ObjTypeId,
		u.TransCurr.Id, u.Startdate, u.Enddate).Scan(&ui)

	if err != nil {
		log.Println("Failed execute func_obj_trans_curr_upd: ", err)
	}

	output, err := json.Marshal(models.Json_id{Id: ui})

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(output)

	return

}

// HandleDelObjTransCurr godoc
// @Summary Delete objtranscurrs
// @Description delete objtranscurrs
// @Tags objtranscurr
// @Accept json
// @Produce  json
// @Param d body models.Json_ids true "Delete objtranscurrs"
// @Success 200 {object} models.Json_ids
// @Failure 500
// @Router /objtranscurr_del [post]
func (s *APG) HandleDelObjTransCurr(w http.ResponseWriter, r *http.Request) {
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
		err = s.Dbpool.QueryRow(context.Background(), "SELECT func_obj_trans_curr_del($1);", id).Scan(&i)
		res = append(res, i)

		if err != nil {
			log.Println("Failed execute func_obj_trans_curr_del: ", err)
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

// HandleGetObjTransCurr godoc
// @Summary Get objtranscurr
// @Description get objtranscurr
// @Tags objtranscurr
// @Produce  json
// @Param id path int true "Objtranscurr by id"
// @Success 200 {object} models.ObjTransCurr_count
// @Failure 500
// @Router /objtranscurr/{id} [get]
func (s *APG) HandleGetObjTransCurr(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	i := vars["id"]
	g := models.ObjTransCurr{}
	out_arr := []models.ObjTransCurr{}

	err := s.Dbpool.QueryRow(context.Background(), "SELECT * from func_obj_trans_curr_getbyid($1);", i).Scan(&g.Id, &g.ObjId,
		&g.ObjTypeId, &g.TransCurr.Id, &g.Startdate, &g.Enddate, &g.ObjName, &g.TransCurr.TransCurrName, &g.TransCurr.TransType.Id,
		&g.TransCurr.CheckDate, &g.TransCurr.NextCheckDate, &g.TransCurr.ProdDate, &g.TransCurr.Serial1, &g.TransCurr.Serial2,
		&g.TransCurr.Serial3, &g.TransCurr.TransType.TransTypeName, &g.TransCurr.TransType.Ratio, &g.TransCurr.TransType.Class,
		&g.TransCurr.TransType.MaxCurr, &g.TransCurr.TransType.NomCurr)

	if err != nil && err != pgx.ErrNoRows {
		log.Println("Failed execute from func_obj_trans_curr_getbyid: ", err)
	}

	out_arr = append(out_arr, g)
	auth := models.Auth{Create: true, Read: true, Update: true, Delete: true}

	out_count, err := json.Marshal(models.ObjTransCurr_count{Values: out_arr, Count: 1, Auth: auth})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Write(out_count)

	return
}

// HandleObjTransCurrByObj godoc
// @Summary Get objtranscurr by object
// @Description get objtranscurr by object
// @Tags objtranscurr
// @Produce  json
// @Param objid query string false "obj&tgu id"
// @Param tid query string false "obj&tgu type id (obj - type = 0, tgu - type > 0)"
// @Success 200 {object} models.ObjTransCurr_count
// @Failure 500
// @Router /objtranscurr_obj [get]
func (s *APG) HandleObjTransCurrByObj(w http.ResponseWriter, r *http.Request) {
	g := models.ObjTransCurr{}
	out_arr := []models.ObjTransCurr{}

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

	err := s.Dbpool.QueryRow(context.Background(), "SELECT * from func_obj_trans_curr_getbyobj($1,$2);", gs1, gs2).Scan(&g.Id, &g.ObjId,
		&g.ObjTypeId, &g.TransCurr.Id, &g.Startdate, &g.Enddate, &g.ObjName, &g.TransCurr.TransCurrName, &g.TransCurr.TransType.Id,
		&g.TransCurr.CheckDate, &g.TransCurr.NextCheckDate, &g.TransCurr.ProdDate, &g.TransCurr.Serial1, &g.TransCurr.Serial2,
		&g.TransCurr.Serial3, &g.TransCurr.TransType.TransTypeName, &g.TransCurr.TransType.Ratio, &g.TransCurr.TransType.Class,
		&g.TransCurr.TransType.MaxCurr, &g.TransCurr.TransType.NomCurr)

	if err != nil && err != pgx.ErrNoRows {
		log.Println("Failed execute from func_obj_trans_curr_getbyobj: ", err)
	}

	out_arr = append(out_arr, g)
	auth := models.Auth{Create: true, Read: true, Update: true, Delete: true}

	out_count, err := json.Marshal(models.ObjTransCurr_count{Values: out_arr, Count: 1, Auth: auth})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Write(out_count)

	return
}
