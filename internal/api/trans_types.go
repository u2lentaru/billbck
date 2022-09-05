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

// HandleTransTypes godoc
// @Summary List transtypes
// @Description get transtype list
// @Tags transtypes
// @Produce json
// @Param page query int false "page number"
// @Param page_size query int false "page size"
// @Param transtypename query string false "transtypename search pattern"
// @Param ordering query string false "order by {id|transtypename}"
// @Param desc query boolean false "descending order {true|false}"
// @Success 200 {object} models.TransType_count
// @Failure 500
// @Router /transtypes [get]
func (s *APG) HandleTransTypes(w http.ResponseWriter, r *http.Request) {
	gs := models.TransType{}
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
	gs1s, ok := query["transtypename"]
	if ok && len(gs1s) > 0 {
		//case insensitive
		gs1 = strings.ToUpper(gs1s[0])
		//quotes
		re := regexp.MustCompile(`'`)
		gs1 = string(re.ReplaceAll([]byte(gs1), []byte("''")))
	}

	gsc := 0
	err := s.Dbpool.QueryRow(ctx, "SELECT * from func_trans_types_cnt($1);", gs1).Scan(&gsc)

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	out_arr := make([]models.TransType, 0,
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
	} else if ords[0] == "transtypename" {
		ord = 2
	}

	dsc := false
	dscs, ok := query["desc"]
	if ok && len(dscs) > 0 {
		if !(dscs[0] == "0") {
			dsc = true
		}
	}

	rows, err := s.Dbpool.Query(ctx, "SELECT * from func_trans_types_get($1,$2,$3,$4,$5);", pg, pgs, gs1, ord, dsc)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&gs.Id, &gs.TransTypeName, &gs.Ratio, &gs.Class, &gs.MaxCurr, &gs.NomCurr)
		if err != nil {
			log.Println("failed to scan row:", err)
		}

		out_arr = append(out_arr, gs)
	}

	auth := models.Auth{Create: true, Read: true, Update: true, Delete: true}
	out_count, err := json.Marshal(models.TransType_count{Values: out_arr, Count: gsc, Auth: auth})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Write(out_count)

	return
}

// HandleAddTransType godoc
// @Summary Add transtype
// @Description add transtype
// @Tags transtypes
// @Accept json
// @Produce  json
// @Param a body models.AddTransType true "New transtype. Significant params: TransTypeName, Ratio, Class, MaxCurr, NomCurr"
// @Success 200 {object} models.Json_id
// @Failure 500
// @Router /transtypes_add [post]
func (s *APG) HandleAddTransType(w http.ResponseWriter, r *http.Request) {
	a := models.AddTransType{}
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
	err = s.Dbpool.QueryRow(context.Background(), "SELECT func_trans_types_add($1,$2,$3,$4,$5);", a.TransTypeName, a.Ratio,
		a.Class, a.MaxCurr, a.NomCurr).Scan(&ai)

	if err != nil {
		log.Println("Failed execute func_trans_types_add: ", err)
	}

	output, err := json.Marshal(models.Json_id{Id: ai})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(output)

	return

}

// HandleUpdTransType godoc
// @Summary Update transtype
// @Description update transtype
// @Tags transtypes
// @Accept json
// @Produce  json
// @Param u body models.TransType true "Update transtype. Significant params: Id, TransTypeName, Ratio, Class, MaxCurr, NomCurr"
// @Success 200 {object} models.Json_id
// @Failure 500
// @Router /transtypes_upd [post]
func (s *APG) HandleUpdTransType(w http.ResponseWriter, r *http.Request) {
	u := models.TransType{}
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
	err = s.Dbpool.QueryRow(context.Background(), "SELECT func_trans_types_upd($1,$2,$3,$4,$5,$6);", u.Id, u.TransTypeName,
		u.Ratio, u.Class, u.MaxCurr, u.NomCurr).Scan(&ui)

	if err != nil {
		log.Println("Failed execute func_trans_types_upd: ", err)
	}

	output, err := json.Marshal(models.Json_id{Id: ui})

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(output)

	return
}

// HandleDelTransType godoc
// @Summary Delete transtypes
// @Description delete transtypes
// @Tags transtypes
// @Accept json
// @Produce  json
// @Param d body models.Json_ids true "Delete transtypes"
// @Success 200 {object} models.Json_ids
// @Failure 500
// @Router /transtypes_del [post]
func (s *APG) HandleDelTransType(w http.ResponseWriter, r *http.Request) {
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
		err = s.Dbpool.QueryRow(context.Background(), "SELECT func_trans_types_del($1);", id).Scan(&i)
		res = append(res, i)

		if err != nil {
			log.Println("Failed execute func_trans_types_del: ", err)
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

// HandleGetTransType godoc
// @Summary Get transtype
// @Description get transtype
// @Tags transtypes
// @Produce  json
// @Param id path int true "TransType by id"
// @Success 200 {object} models.TransType_count
// @Failure 500
// @Router /transtypes/{id} [get]
func (s *APG) HandleGetTransType(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	i := vars["id"]
	g := models.TransType{}
	out_arr := []models.TransType{}

	err := s.Dbpool.QueryRow(context.Background(), "SELECT * from func_trans_type_get($1);", i).Scan(&g.Id, &g.TransTypeName,
		&g.Ratio, &g.Class, &g.MaxCurr, &g.NomCurr)

	if err != nil && err != pgx.ErrNoRows {
		log.Println("Failed execute from func_trans_type_get: ", err)
	}

	out_arr = append(out_arr, g)
	auth := models.Auth{Create: true, Read: true, Update: true, Delete: true}

	out_count, err := json.Marshal(models.TransType_count{Values: out_arr, Count: 1, Auth: auth})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Write(out_count)

	return

}
