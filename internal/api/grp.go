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

// HandleGRp godoc
// @Summary List grp
// @Description get grp list
// @Tags grp
// @Produce  json
// @Param page query int false "page number"
// @Param page_size query int false "page size"
// @Param grpname query string false "grpname search pattern"
// @Param ordering query string false "order by {id|grpname}"
// @Param desc query boolean false "descending order {true|false}"
// @Success 200 {object} models.GRp_count
// @Failure 500
// @Router /grp [get]
func (s *APG) HandleGRp(w http.ResponseWriter, r *http.Request) {
	gs := models.GRp{}
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
	gs1s, ok := query["grpname"]
	if ok && len(gs1s) > 0 {
		//case insensitive
		gs1 = strings.ToUpper(gs1s[0])
		//quotes
		re := regexp.MustCompile(`'`)
		gs1 = string(re.ReplaceAll([]byte(gs1), []byte("''")))
	}

	gsc := 0
	err := s.Dbpool.QueryRow(ctx, "SELECT * from func_grp_cnt($1);", gs1).Scan(&gsc)

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	out_arr := make([]models.GRp, 0,
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
	} else if ords[0] == "grpname" {
		ord = 2
	}

	dsc := false
	dscs, ok := query["desc"]
	if ok && len(dscs) > 0 {
		if !(dscs[0] == "0") {
			dsc = true
		}
	}

	rows, err := s.Dbpool.Query(ctx, "SELECT * from func_grp_get($1,$2,$3,$4,$5);", pg, pgs, gs1, ord, dsc)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&gs.Id, &gs.GRpName)
		if err != nil {
			log.Println("failed to scan row:", err)
		}

		out_arr = append(out_arr, gs)
	}

	auth := models.Auth{Create: true, Read: true, Update: true, Delete: true}
	out_count, err := json.Marshal(models.GRp_count{Values: out_arr, Count: gsc, Auth: auth})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Write(out_count)

	return

}

// HandleAddGRp godoc
// @Summary Add grp
// @Description add grp
// @Tags grp
// @Accept json
// @Produce  json
// @Param a body models.AddGRp true "New grp"
// @Success 200 {object} models.Json_id
// @Failure 500
// @Router /grp_add [post]
func (s *APG) HandleAddGRp(w http.ResponseWriter, r *http.Request) {
	a := models.AddGRp{}
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
	err = s.Dbpool.QueryRow(context.Background(), "SELECT func_grp_add($1);", a.GRpName).Scan(&ai)

	if err != nil {
		log.Println("Failed execute func_grp_add: ", err)
	}

	output, err := json.Marshal(models.Json_id{Id: ai})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(output)

	return

}

// HandleUpdGRp godoc
// @Summary Update grp
// @Description update grp
// @Tags grp
// @Accept json
// @Produce  json
// @Param u body models.GRp true "Update grp"
// @Success 200 {object} models.Json_id
// @Failure 500
// @Router /grp_upd [post]
func (s *APG) HandleUpdGRp(w http.ResponseWriter, r *http.Request) {
	u := models.GRp{}
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
	err = s.Dbpool.QueryRow(context.Background(), "SELECT func_grp_upd($1,$2);", u.Id, u.GRpName).Scan(&ui)

	if err != nil {
		log.Println("Failed execute func_grp_upd: ", err)
	}

	output, err := json.Marshal(models.Json_id{Id: ui})

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(output)

	return
}

// HandleDelGRp godoc
// @Summary Delete grp
// @Description delete grp
// @Tags grp
// @Accept json
// @Produce  json
// @Param d body models.Json_ids true "Delete grp"
// @Success 200 {object} models.Json_ids
// @Failure 500
// @Router /grp_del [post]
func (s *APG) HandleDelGRp(w http.ResponseWriter, r *http.Request) {
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
		err = s.Dbpool.QueryRow(context.Background(), "SELECT func_grp_del($1);", id).Scan(&i)
		res = append(res, i)

		if err != nil {
			log.Println("Failed execute func_grp_del: ", err)
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

// HandleGetGRp godoc
// @Summary Get grp
// @Description get grp
// @Tags grp
// @Produce  json
// @Param id path int true "GRp by id"
// @Success 200 {object} models.GRp_count
// @Failure 500
// @Router /grp/{id} [get]
func (s *APG) HandleGetGRp(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	i := vars["id"]
	g := models.GRp{}
	out_arr := []models.GRp{}

	err := s.Dbpool.QueryRow(context.Background(), "SELECT * from func_grp_getbyid($1);", i).Scan(&g.Id, &g.GRpName)

	if err != nil && err != pgx.ErrNoRows {
		log.Println("Failed execute from func_grp_getbyid: ", err)
	}

	out_arr = append(out_arr, g)
	auth := models.Auth{Create: true, Read: true, Update: true, Delete: true}

	// output, err := json.Marshal(g)
	out_count, err := json.Marshal(models.GRp_count{Values: out_arr, Count: 1, Auth: auth})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(out_count)
	return

}
