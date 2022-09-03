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

// HandleObjLines godoc
// @Summary List objlines
// @Description get objlines list
// @Tags objlines
// @Produce json
// @Param page query int false "page number"
// @Param page_size query int false "page size"
// @Param objectname query string false "objectname search pattern"
// @Param cableresistancename query string false "cableresistancename search pattern"
// @Param ordering query string false "order by {id|cableresistancename|objectname}"
// @Param desc query boolean false "descending order {true|false}"
// @Success 200 {object} models.ObjLine_count
// @Failure 500
// @Router /objlines [get]
func (s *APG) HandleObjLines(w http.ResponseWriter, r *http.Request) {
	// start := time.Now()
	gs := models.ObjLine{}
	ctx := context.Background()
	// out_arr := []models.ObjLine{}

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
	gs2s, ok := query["cableresistancename"]
	if ok && len(gs2s) > 0 {
		//case insensitive
		gs2 = strings.ToUpper(gs2s[0])
		//quotes
		re := regexp.MustCompile(`'`)
		gs2 = string(re.ReplaceAll([]byte(gs2), []byte("''")))
	}

	gsc := 0
	err := s.Dbpool.QueryRow(ctx, "SELECT * from func_obj_lines_cnt($1,$2);", gs1, gs2).Scan(&gsc)

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	out_arr := make([]models.ObjLine, 0,
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
	} else if ords[0] == "objectname" {
		ord = 2
	} else if ords[0] == "cableresistancename" {
		ord = 3
	}

	dsc := false
	dscs, ok := query["desc"]
	if ok && len(dscs) > 0 {
		if !(dscs[0] == "0") {
			dsc = true
		}
	}

	rows, err := s.Dbpool.Query(ctx, "SELECT * from func_obj_lines_get($1,$2,$3,$4,$5,$6);", pg, pgs, gs1, gs2, ord, dsc)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&gs.Id, &gs.ObjId, &gs.ObjTypeId, &gs.CableResistance.Id, &gs.LineLength, &gs.Startdate, &gs.Enddate,
			&gs.ObjName, &gs.CableResistance.CableResistanceName, &gs.CableResistance.Resistance, &gs.CableResistance.MaterialType)
		if err != nil {
			log.Println("failed to scan row:", err)
		}

		out_arr = append(out_arr, gs)
	}

	auth := models.Auth{Create: true, Read: true, Update: true, Delete: true}
	out_count, err := json.Marshal(models.ObjLine_count{Values: out_arr, Count: gsc, Auth: auth})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Write(out_count)

	// log.Printf("Work time %s", time.Since(start))

	return
}

// HandleAddObjLine godoc
// @Summary Add objline
// @Description add objline
// @Tags objlines
// @Accept json
// @Produce  json
// @Param a body models.AddObjLine true "New objtranscurr. Significant params: ObjId, ObjTypeId, CableResistance.Id, LineLength, Startdate"
// @Success 200 {object} models.Json_id
// @Failure 500
// @Router /objlines_add [post]
func (s *APG) HandleAddObjLine(w http.ResponseWriter, r *http.Request) {
	a := models.AddObjLine{}
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
	err = s.Dbpool.QueryRow(context.Background(), "SELECT func_obj_lines_add($1,$2,$3,$4,$5);", a.ObjId, a.ObjTypeId,
		a.CableResistance.Id, a.LineLength, a.Startdate).Scan(&ai)

	if err != nil {
		log.Println("Failed execute func_obj_lines_add: ", err)
	}

	output, err := json.Marshal(models.Json_id{Id: ai})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(output)

	return
}

// HandleUpdObjLine godoc
// @Summary Update objline
// @Description update objline
// @Tags objlines
// @Accept json
// @Produce  json
// @Param u body models.ObjLine true "Update objtranscurr. Significant params: Id, ObjId, ObjTypeId, CableResistance.Id, LineLength, Startdate, Enddate"
// @Success 200 {object} models.Json_id
// @Failure 500
// @Router /objlines_upd [post]
func (s *APG) HandleUpdObjLine(w http.ResponseWriter, r *http.Request) {
	u := models.ObjLine{}
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
	err = s.Dbpool.QueryRow(context.Background(), "SELECT func_obj_lines_upd($1,$2,$3,$4,$5,$6,$7);", u.Id, u.ObjId, u.ObjTypeId,
		u.CableResistance.Id, u.LineLength, u.Startdate, u.Enddate).Scan(&ui)

	if err != nil {
		log.Println("Failed execute func_obj_lines_upd: ", err)
	}

	output, err := json.Marshal(models.Json_id{Id: ui})

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(output)

	return
}

// HandleDelObjLine godoc
// @Summary Delete objlines
// @Description delete objlines
// @Tags objlines
// @Accept json
// @Produce  json
// @Param d body models.Json_ids true "Delete objlines"
// @Success 200 {object} models.Json_ids
// @Failure 500
// @Router /objlines_del [post]
func (s *APG) HandleDelObjLine(w http.ResponseWriter, r *http.Request) {
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
		err = s.Dbpool.QueryRow(context.Background(), "SELECT func_obj_lines_del($1);", id).Scan(&i)
		res = append(res, i)

		if err != nil {
			log.Println("Failed execute func_obj_lines_del: ", err)
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

// HandleGetObjLine godoc
// @Summary Get objline
// @Description get objline
// @Tags objlines
// @Produce  json
// @Param id path int true "Objline by id"
// @Success 200 {object} models.ObjLine_count
// @Failure 500
// @Router /objlines/{id} [get]
func (s *APG) HandleGetObjLine(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	i := vars["id"]
	g := models.ObjLine{}
	out_arr := []models.ObjLine{}

	err := s.Dbpool.QueryRow(context.Background(), "SELECT * from func_obj_line_get($1);", i).Scan(&g.Id, &g.ObjId,
		&g.ObjTypeId, &g.CableResistance.Id, &g.LineLength, &g.Startdate, &g.Enddate, &g.ObjName, &g.CableResistance.CableResistanceName,
		&g.CableResistance.Resistance, &g.CableResistance.MaterialType)

	if err != nil && err != pgx.ErrNoRows {
		log.Println("Failed execute from func_obj_line_get: ", err)
	}

	out_arr = append(out_arr, g)
	auth := models.Auth{Create: true, Read: true, Update: true, Delete: true}

	out_count, err := json.Marshal(models.ObjLine_count{Values: out_arr, Count: 1, Auth: auth})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Write(out_count)

	return

}

// HandleObjLinesByObj godoc
// @Summary Get objlines by object
// @Description get objlines by object
// @Tags objlines
// @Produce  json
// @Param objid query string false "obj&tgu id"
// @Param tid query string false "obj&tgu type id (obj - type = 0, tgu - type > 0)"
// @Success 200 {object} models.ObjLine_count
// @Failure 500
// @Router /objlines_obj [get]
func (s *APG) HandleObjLinesByObj(w http.ResponseWriter, r *http.Request) {
	// start := time.Now()
	g := models.ObjLine{}
	ctx := context.Background()
	// out_arr := []models.ObjLine{}

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

	gsc := 0
	err := s.Dbpool.QueryRow(ctx, "SELECT * from func_obj_lines_obj_cnt($1,$2);", gs1, gs2).Scan(&gsc)

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	out_arr := make([]models.ObjLine, 0, gsc)

	rows, err := s.Dbpool.Query(ctx, "SELECT * from func_obj_lines_obj($1,$2);", gs1, gs2)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&g.Id, &g.ObjId, &g.ObjTypeId, &g.CableResistance.Id, &g.LineLength, &g.Startdate, &g.Enddate,
			&g.ObjName, &g.CableResistance.CableResistanceName, &g.CableResistance.Resistance, &g.CableResistance.MaterialType)
		if err != nil {
			log.Println("failed to scan row:", err)
		}

		out_arr = append(out_arr, g)
	}

	// gsc := 0
	// err = s.Dbpool.QueryRow(ctx, "SELECT * from func_obj_lines_obj_cnt($1,$2);", gs1, gs2).Scan(&gsc)

	// if err != nil {
	// 	http.Error(w, err.Error(), 500)
	// 	return
	// }

	auth := models.Auth{Create: true, Read: true, Update: true, Delete: true}
	out_count, err := json.Marshal(models.ObjLine_count{Values: out_arr, Count: gsc, Auth: auth})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Write(out_count)

	// log.Printf("Work time %s", time.Since(start))

	return

}
