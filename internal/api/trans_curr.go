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

// HandleTransCurr godoc
// @Summary List current transformers
// @Description get current transformers list
// @Tags transcurr
// @Produce json
// @Param page query int false "page number"
// @Param page_size query int false "page size"
// @Param transcurrname query string false "transcurrname search pattern"
// @Param ordering query string false "order by {id|transcurrname}"
// @Param desc query boolean false "descending order {true|false}"
// @Success 200 {object} models.TransCurr_count
// @Failure 500
// @Router /transcurr [get]
func (s *APG) HandleTransCurr(w http.ResponseWriter, r *http.Request) {
	gs := models.TransCurr{}
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
	gs1s, ok := query["transcurrname"]
	if ok && len(gs1s) > 0 {
		//case insensitive
		gs1 = strings.ToUpper(gs1s[0])
		//quotes
		re := regexp.MustCompile(`'`)
		gs1 = string(re.ReplaceAll([]byte(gs1), []byte("''")))
	}

	gsc := 0
	err := s.Dbpool.QueryRow(ctx, "SELECT * from func_trans_curr_cnt($1);", gs1).Scan(&gsc)

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	out_arr := make([]models.TransCurr, 0,
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
	} else if ords[0] == "transcurrname" {
		ord = 2
	}

	dsc := false
	dscs, ok := query["desc"]
	if ok && len(dscs) > 0 {
		if !(dscs[0] == "0") {
			dsc = true
		}
	}

	rows, err := s.Dbpool.Query(ctx, "SELECT * from func_trans_curr_get($1,$2,$3,$4,$5);", pg, pgs, gs1, ord, dsc)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&gs.Id, &gs.TransCurrName, &gs.TransType.Id, &gs.CheckDate, &gs.NextCheckDate, &gs.ProdDate, &gs.Serial1,
			&gs.Serial2, &gs.Serial3, &gs.TransType.TransTypeName, &gs.TransType.Ratio, &gs.TransType.Class, &gs.TransType.MaxCurr,
			&gs.TransType.NomCurr)
		if err != nil {
			log.Println("failed to scan row:", err)
		}

		out_arr = append(out_arr, gs)
	}

	auth := models.Auth{Create: true, Read: true, Update: true, Delete: true}
	out_count, err := json.Marshal(models.TransCurr_count{Values: out_arr, Count: gsc, Auth: auth})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Write(out_count)

	return

}

// HandleAddTransCurr godoc
// @Summary Add current transformer
// @Description add current transformer
// @Tags transcurr
// @Accept json
// @Produce  json
// @Param a body models.AddTransCurr true "New current transformer. Significant params: TransCurrName, TransType.Id, CheckDate(n), NextCheckDate(n), ProdDate(n), Serial1(n), Serial2(n), Serial3(n)"
// @Success 200 {object} models.Json_id
// @Failure 500
// @Router /transcurr_add [post]
func (s *APG) HandleAddTransCurr(w http.ResponseWriter, r *http.Request) {
	a := models.AddTransCurr{}
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
	err = s.Dbpool.QueryRow(context.Background(), "SELECT func_trans_curr_add($1,$2,$3,$4,$5,$6,$7,$8);", a.TransCurrName,
		a.TransType.Id, a.CheckDate, a.NextCheckDate, a.ProdDate, a.Serial1, a.Serial2, a.Serial3).Scan(&ai)

	if err != nil {
		log.Println("Failed execute func_trans_curr_add: ", err)
	}

	output, err := json.Marshal(models.Json_id{Id: ai})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(output)

	return

}

// HandleUpdTransCurr godoc
// @Summary Update current transformer
// @Description update current transformer
// @Tags transcurr
// @Accept json
// @Produce  json
// @Param u body models.TransCurr true "Update current transformer. Significant params: Id, TransCurrName, TransType.Id, CheckDate(n), NextCheckDate(n), ProdDate(n), Serial1(n), Serial2(n), Serial3(n)"
// @Success 200 {object} models.Json_id
// @Failure 500
// @Router /transcurr_upd [post]
func (s *APG) HandleUpdTransCurr(w http.ResponseWriter, r *http.Request) {
	u := models.TransCurr{}
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
	err = s.Dbpool.QueryRow(context.Background(), "SELECT func_trans_curr_upd($1,$2,$3,$4,$5,$6,$7,$8,$9);", u.Id, u.TransCurrName,
		u.TransType.Id, u.CheckDate, u.NextCheckDate, u.ProdDate, u.Serial1, u.Serial2, u.Serial3).Scan(&ui)

	if err != nil {
		log.Println("Failed execute func_trans_curr_upd: ", err)
	}

	output, err := json.Marshal(models.Json_id{Id: ui})

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(output)

	return

}

// HandleDelTransCurr godoc
// @Summary Delete current transformers
// @Description delete current transformers
// @Tags transcurr
// @Accept json
// @Produce  json
// @Param d body models.Json_ids true "Delete current transformers"
// @Success 200 {object} models.Json_ids
// @Failure 500
// @Router /transcurr_del [post]
func (s *APG) HandleDelTransCurr(w http.ResponseWriter, r *http.Request) {
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
		err = s.Dbpool.QueryRow(context.Background(), "SELECT func_trans_curr_del($1);", id).Scan(&i)
		res = append(res, i)

		if err != nil {
			log.Println("Failed execute func_trans_curr_del: ", err)
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

// HandleGetTransCurr godoc
// @Summary Get current transformer
// @Description get current transformer
// @Tags transcurr
// @Produce  json
// @Param id path int true "Current transformer by id"
// @Success 200 {object} models.TransCurr_count
// @Failure 500
// @Router /transcurr/{id} [get]
func (s *APG) HandleGetTransCurr(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	i := vars["id"]
	g := models.TransCurr{}
	out_arr := []models.TransCurr{}

	err := s.Dbpool.QueryRow(context.Background(), "SELECT * from func_trans_curr_getbyid($1);", i).Scan(&g.Id, &g.TransCurrName,
		&g.TransType.Id, &g.CheckDate, &g.NextCheckDate, &g.ProdDate, &g.Serial1, &g.Serial2, &g.Serial3, &g.TransType.TransTypeName,
		&g.TransType.Ratio, &g.TransType.Class, &g.TransType.MaxCurr, &g.TransType.NomCurr)

	if err != nil && err != pgx.ErrNoRows {
		log.Println("Failed execute from func_trans_curr_getbyid: ", err)
	}

	out_arr = append(out_arr, g)
	auth := models.Auth{Create: true, Read: true, Update: true, Delete: true}

	out_count, err := json.Marshal(models.TransCurr_count{Values: out_arr, Count: 1, Auth: auth})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Write(out_count)

	return
}
