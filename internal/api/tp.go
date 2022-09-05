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

// HandleTp godoc
// @Summary List tp
// @Description get tp list
// @Tags tp
// @Produce  json
// @Param page query int false "page number"
// @Param page_size query int false "page size"
// @Param tpname query string false "tpname search pattern"
// @Param ordering query string false "order by {id|tpname}"
// @Param desc query boolean false "descending order {true|false}"
// @Success 200 {object} models.Tp_count
// @Failure 500
// @Router /tp [get]
func (s *APG) HandleTp(w http.ResponseWriter, r *http.Request) {
	gs := models.Tp{}
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
	gs1s, ok := query["tpname"]
	if ok && len(gs1s) > 0 {
		//case insensitive
		gs1 = strings.ToUpper(gs1s[0])
		//quotes
		re := regexp.MustCompile(`'`)
		gs1 = string(re.ReplaceAll([]byte(gs1), []byte("''")))
	}

	gsc := 0
	err := s.Dbpool.QueryRow(ctx, "SELECT * from func_tp_cnt($1);", gs1).Scan(&gsc)

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	out_arr := make([]models.Tp, 0,
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
	} else if ords[0] == "tpname" {
		ord = 2
	}

	dsc := false
	dscs, ok := query["desc"]
	if ok && len(dscs) > 0 {
		if !(dscs[0] == "0") {
			dsc = true
		}
	}

	rows, err := s.Dbpool.Query(ctx, "SELECT * from func_tp_get($1,$2,$3,$4,$5);", pg, pgs, gs1, ord, dsc)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&gs.Id, &gs.TpName, &gs.GRp.Id, &gs.GRp.GRpName)
		if err != nil {
			log.Println("failed to scan row:", err)
		}

		out_arr = append(out_arr, gs)
	}

	auth := models.Auth{Create: true, Read: true, Update: true, Delete: true}
	out_count, err := json.Marshal(models.Tp_count{Values: out_arr, Count: gsc, Auth: auth})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Write(out_count)

	return

}

// HandleAddTp godoc
// @Summary Add tp
// @Description add tp
// @Tags tp
// @Accept json
// @Produce  json
// @Param a body models.AddTp true "New tp"
// @Success 200 {object} models.Json_id
// @Failure 500
// @Router /tp_add [post]
func (s *APG) HandleAddTp(w http.ResponseWriter, r *http.Request) {
	a := models.AddTp{}
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
	err = s.Dbpool.QueryRow(context.Background(), "SELECT func_tp_add($1,$2);", a.TpName, a.GRp.Id).Scan(&ai)

	if err != nil {
		log.Println("Failed execute func_tp_add: ", err)
	}

	output, err := json.Marshal(models.Json_id{Id: ai})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(output)

	return

}

// HandleUpdTp godoc
// @Summary Update tp
// @Description update tp
// @Tags tp
// @Accept json
// @Produce  json
// @Param u body models.Tp true "Update tp"
// @Success 200 {object} models.Json_id
// @Failure 500
// @Router /tp_upd [post]
func (s *APG) HandleUpdTp(w http.ResponseWriter, r *http.Request) {
	u := models.Tp{}
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
	err = s.Dbpool.QueryRow(context.Background(), "SELECT func_tp_upd($1,$2,$3);", u.Id, u.TpName, u.GRp.Id).Scan(&ui)

	if err != nil {
		log.Println("Failed execute func_tp_upd: ", err)
	}

	output, err := json.Marshal(models.Json_id{Id: ui})

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(output)

	return
}

// HandleDelTp godoc
// @Summary Delete tp
// @Description delete tp
// @Tags tp
// @Accept json
// @Produce  json
// @Param d body models.Json_ids true "Delete tp"
// @Success 200 {object} models.Json_ids
// @Failure 500
// @Router /tp_del [post]
func (s *APG) HandleDelTp(w http.ResponseWriter, r *http.Request) {
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
		err = s.Dbpool.QueryRow(context.Background(), "SELECT func_tp_del($1);", id).Scan(&i)
		res = append(res, i)

		if err != nil {
			log.Println("Failed execute func_tp_del: ", err)
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

// HandleGetTp godoc
// @Summary Get tp
// @Description get tp
// @Tags tp
// @Produce  json
// @Param id path int true "Tp by id"
// @Success 200 {object} models.Tp_count
// @Failure 500
// @Router /tp/{id} [get]
func (s *APG) HandleGetTp(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	i := vars["id"]
	g := models.Tp{}
	out_arr := []models.Tp{}

	err := s.Dbpool.QueryRow(context.Background(), "SELECT * from func_tp_getbyid($1);", i).Scan(&g.Id, &g.TpName, &g.GRp.Id, &g.GRp.GRpName)

	if err != nil && err != pgx.ErrNoRows {
		log.Println("Failed execute from func_tp_getbyid: ", err)
	}

	out_arr = append(out_arr, g)
	auth := models.Auth{Create: true, Read: true, Update: true, Delete: true}

	// output, err := json.Marshal(g)
	out_count, err := json.Marshal(models.Tp_count{Values: out_arr, Count: 1, Auth: auth})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(out_count)

	return

}
