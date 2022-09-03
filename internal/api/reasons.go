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

// HandleReasons godoc
// @Summary List reasons
// @Description get reason list
// @Tags reasons
// @Produce json
// @Param page query int false "page number"
// @Param page_size query int false "page size"
// @Param reasonname query string false "reasonname search pattern"
// @Param ordering query string false "order by {id|reasonname}"
// @Param desc query boolean false "descending order {true|false}"
// @Success 200 {object} models.Reason_count
// @Failure 500
// @Router /reasons [get]
func (s *APG) HandleReasons(w http.ResponseWriter, r *http.Request) {
	gs := models.Reason{}
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
	gs1s, ok := query["reasonname"]
	if ok && len(gs1s) > 0 {
		//case insensitive
		gs1 = strings.ToUpper(gs1s[0])
		//quotes
		re := regexp.MustCompile(`'`)
		gs1 = string(re.ReplaceAll([]byte(gs1), []byte("''")))
	}

	gsc := 0
	err := s.Dbpool.QueryRow(ctx, "SELECT * from func_reasons_cnt($1);", gs1).Scan(&gsc)

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	out_arr := make([]models.Reason, 0,
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
	} else if ords[0] == "reasonname" {
		ord = 2
	}

	dsc := false
	dscs, ok := query["desc"]
	if ok && len(dscs) > 0 {
		if !(dscs[0] == "0") {
			dsc = true
		}
	}

	rows, err := s.Dbpool.Query(ctx, "SELECT * from func_reasons_get($1,$2,$3,$4,$5);", pg, pgs, gs1, ord, dsc)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&gs.Id, &gs.ReasonName)
		if err != nil {
			log.Println("failed to scan row:", err)
		}

		out_arr = append(out_arr, gs)
	}

	auth := models.Auth{Create: true, Read: true, Update: true, Delete: true}
	out_count, err := json.Marshal(models.Reason_count{Values: out_arr, Count: gsc, Auth: auth})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Write(out_count)

	return
}

// HandleAddReason godoc
// @Summary Add reason
// @Description add reason
// @Tags reasons
// @Accept json
// @Produce  json
// @Param a body models.AddReason true "New reason"
// @Success 200 {object} models.Json_id
// @Failure 500
// @Router /reasons_add [post]
func (s *APG) HandleAddReason(w http.ResponseWriter, r *http.Request) {
	a := models.AddReason{}
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
	err = s.Dbpool.QueryRow(context.Background(), "SELECT func_reasons_add($1);", a.ReasonName).Scan(&ai)

	if err != nil {
		log.Println("Failed execute func_reasons_add: ", err)
	}

	output, err := json.Marshal(models.Json_id{Id: ai})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(output)

	return

}

// HandleUpdReason godoc
// @Summary Update reason
// @Description update reason
// @Tags reasons
// @Accept json
// @Produce  json
// @Param u body models.Reason true "Update reason"
// @Success 200 {object} models.Json_id
// @Failure 500
// @Router /reasons_upd [post]
func (s *APG) HandleUpdReason(w http.ResponseWriter, r *http.Request) {
	u := models.Reason{}
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
	err = s.Dbpool.QueryRow(context.Background(), "SELECT func_reasons_upd($1,$2);", u.Id, u.ReasonName).Scan(&ui)

	if err != nil {
		log.Println("Failed execute func_reasons_upd: ", err)
	}

	output, err := json.Marshal(models.Json_id{Id: ui})

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(output)

	return
}

// HandleDelReason godoc
// @Summary Delete reasons
// @Description delete reasons
// @Tags reasons
// @Accept json
// @Produce  json
// @Param d body models.Json_ids true "Delete reasons"
// @Success 200 {object} models.Json_ids
// @Failure 500
// @Router /reasons_del [post]
func (s *APG) HandleDelReason(w http.ResponseWriter, r *http.Request) {
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
		err = s.Dbpool.QueryRow(context.Background(), "SELECT func_reasons_del($1);", id).Scan(&i)
		res = append(res, i)

		if err != nil {
			log.Println("Failed execute func_reasons_del: ", err)
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

// HandleGetReason godoc
// @Summary Get reason
// @Description get reason
// @Tags reasons
// @Produce  json
// @Param id path int true "Reason by id"
// @Success 200 {object} models.Reason_count
// @Failure 500
// @Router /reasons/{id} [get]
func (s *APG) HandleGetReason(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	i := vars["id"]
	g := models.Reason{}
	out_arr := []models.Reason{}

	err := s.Dbpool.QueryRow(context.Background(), "SELECT * from func_reason_get($1);", i).Scan(&g.Id, &g.ReasonName)

	if err != nil && err != pgx.ErrNoRows {
		log.Println("Failed execute from func_reason_get: ", err)
	}

	out_arr = append(out_arr, g)
	auth := models.Auth{Create: true, Read: true, Update: true, Delete: true}

	// output, err := json.Marshal(g)
	out_count, err := json.Marshal(models.Reason_count{Values: out_arr, Count: 1, Auth: auth})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(out_count)

	return

}
