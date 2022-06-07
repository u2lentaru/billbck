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

// HandleAskueTypes godoc
// @Summary List askuetypes
// @Description get askuetypes list
// @Tags askuetypes
// @Produce  json
// @Param page query int false "page number"
// @Param page_size query int false "page size"
// @Param askuetypename query string false "askuetypename search pattern"
// @Param ordering query string false "order by {id|askuetypename}"
// @Param desc query boolean false "descending order {true|false}"
// @Success 200 {array} models.CableResistance_count
// @Failure 500
// @Router /askuetypes [get]
func (s *APG) HandleAskueTypes(w http.ResponseWriter, r *http.Request) {
	gs := models.AskueType{}
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

	out_arr := make([]models.AskueType, 0, pgs)

	gs1 := ""
	gs1s, ok := query["askuetypename"]
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
	} else if ords[0] == "askuetypename" {
		ord = 2
	}

	dsc := false
	dscs, ok := query["desc"]
	if ok && len(dscs) > 0 {
		if !(dscs[0] == "0") {
			dsc = true
		}
	}

	rows, err := s.Dbpool.Query(ctx, "SELECT * from func_askue_types_get($1,$2,$3,$4,$5);", pg, pgs, gs1, ord, dsc)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&gs.Id, &gs.AskueTypeName, &gs.StartLine, &gs.PuColumn, &gs.ValueColumn, &gs.DateColumn, &gs.DateColumnArray)
		if err != nil {
			log.Println("failed to scan row:", err)
		}

		out_arr = append(out_arr, gs)
	}

	gsc := 0
	err = s.Dbpool.QueryRow(ctx, "SELECT * from func_askue_types_cnt($1);", gs1).Scan(&gsc)

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	auth := models.Auth{Create: true, Read: true, Update: true, Delete: true}
	out_count, err := json.Marshal(models.AskueType_count{Values: out_arr, Count: gsc, Auth: auth})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Write(out_count)

	return
}

// HandleAddAskueType godoc
// @Summary Add askuetype
// @Description add askuetype
// @Tags askuetypes
// @Accept json
// @Produce  json
// @Param a body models.AddAskueType true "New askuetype. Significant params: AskueTypeName, StartLine, PuColumn, ValueColumn, DateColumn, DateColumnArray(n)"
// @Success 200 {array} models.Json_id
// @Failure 500
// @Router /askuetypes_add [post]
func (s *APG) HandleAddAskueType(w http.ResponseWriter, r *http.Request) {
	a := models.AddAskueType{}
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
	err = s.Dbpool.QueryRow(context.Background(), "SELECT func_askue_types_add($1,$2,$3,$4,$5,$6);", a.AskueTypeName, a.StartLine,
		a.PuColumn, a.ValueColumn, a.DateColumn, a.DateColumnArray).Scan(&ai)

	if err != nil {
		log.Println("Failed execute func_askue_types_add: ", err)
	}

	output, err := json.Marshal(models.Json_id{Id: ai})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(output)

	return

}

// HandleUpdAskueType godoc
// @Summary Update askuetype
// @Description update askuetype
// @Tags askuetypes
// @Accept json
// @Produce  json
// @Param u body models.AskueType true "Update askuetype. Significant params: Id, AskueTypeName, StartLine, PuColumn, ValueColumn, DateColumn, DateColumnArray(n)"
// @Success 200 {array} models.Json_id
// @Failure 500
// @Router /askuetypes_upd [post]
func (s *APG) HandleUpdAskueType(w http.ResponseWriter, r *http.Request) {
	u := models.AskueType{}
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
	err = s.Dbpool.QueryRow(context.Background(), "SELECT func_askue_types_upd($1,$2,$3,$4,$5,$6,$7);", u.Id, u.AskueTypeName,
		u.StartLine, u.PuColumn, u.ValueColumn, u.DateColumn, u.DateColumnArray).Scan(&ui)

	if err != nil {
		log.Println("Failed execute func_askue_types_upd: ", err)
	}

	output, err := json.Marshal(models.Json_id{Id: ui})

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(output)

	return

}

// HandleDelAskueType godoc
// @Summary Delete askuetypes
// @Description delete askuetypes
// @Tags askuetypes
// @Accept json
// @Produce  json
// @Param d body models.Json_ids true "Delete askuetypes"
// @Success 200 {array} models.Json_ids
// @Failure 500
// @Router /askuetypes_del [post]
func (s *APG) HandleDelAskueType(w http.ResponseWriter, r *http.Request) {
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
		err = s.Dbpool.QueryRow(context.Background(), "SELECT func_askue_types_del($1);", id).Scan(&i)
		res = append(res, i)

		if err != nil {
			log.Println("Failed execute func_askue_types_del: ", err)
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

// HandleGetAskueType godoc
// @Summary Get askuetype
// @Description get askuetype
// @Tags askuetypes
// @Produce  json
// @Param id path int true "AskueType by id"
// @Success 200 {array} models.AskueType_count
// @Failure 500
// @Router /askuetypes/{id} [get]
func (s *APG) HandleGetAskueType(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	i := vars["id"]
	g := models.AskueType{}
	out_arr := []models.AskueType{}

	err := s.Dbpool.QueryRow(context.Background(), "SELECT * from func_askue_type_get($1);", i).Scan(&g.Id, &g.AskueTypeName,
		&g.StartLine, &g.PuColumn, &g.ValueColumn, &g.DateColumn, &g.DateColumnArray)

	if err != nil && err != pgx.ErrNoRows {
		log.Println("Failed execute from func_askue_type_get: ", err)
	}

	out_arr = append(out_arr, g)
	auth := models.Auth{Create: true, Read: true, Update: true, Delete: true}

	// output, err := json.Marshal(g)
	out_count, err := json.Marshal(models.AskueType_count{Values: out_arr, Count: 1, Auth: auth})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(out_count)

	return
}
