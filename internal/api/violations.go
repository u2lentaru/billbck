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

// HandleViolations godoc
// @Summary List violations
// @Description get violations list
// @Tags violations
// @Produce  json
// @Param page query int false "page number"
// @Param page_size query int false "page size"
// @Param violationname query string false "violationname search pattern"
// @Param ordering query string false "order by {id|violationname}"
// @Param desc query boolean false "descending order {true|false}"
// @Success 200 {array} models.Violation_count
// @Failure 500
// @Router /violations [get]
func (s *APG) HandleViolations(w http.ResponseWriter, r *http.Request) {
	gs := models.Violation{}
	ctx := context.Background()
	out_arr := []models.Violation{}

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
	gs1s, ok := query["violationname"]
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
	} else if ords[0] == "violationname" {
		ord = 2
	}

	dsc := false
	dscs, ok := query["desc"]
	if ok && len(dscs) > 0 {
		if !(dscs[0] == "0") {
			dsc = true
		}
	}

	rows, err := s.Dbpool.Query(ctx, "SELECT * from func_violations_get($1,$2,$3,$4,$5);", pg, pgs, gs1, ord, dsc)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&gs.Id, &gs.ViolationName)
		if err != nil {
			log.Println("failed to scan row:", err)
		}

		out_arr = append(out_arr, gs)
	}

	gsc := 0
	err = s.Dbpool.QueryRow(ctx, "SELECT * from func_violations_cnt($1);", gs1).Scan(&gsc)

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	auth := models.Auth{Create: true, Read: true, Update: true, Delete: true}
	out_count, err := json.Marshal(models.Violation_count{Values: out_arr, Count: gsc, Auth: auth})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Write(out_count)

	return
}

// HandleAddViolation godoc
// @Summary Add violation
// @Description add violation
// @Tags violations
// @Accept json
// @Produce  json
// @Param a body models.AddViolation true "New violation"
// @Success 200 {array} models.Json_id
// @Failure 500
// @Router /violations_add [post]
func (s *APG) HandleAddViolation(w http.ResponseWriter, r *http.Request) {
	a := models.AddViolation{}
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
	err = s.Dbpool.QueryRow(context.Background(), "SELECT func_violations_add($1);", a.ViolationName).Scan(&ai)

	if err != nil {
		log.Println("Failed execute func_violations_add: ", err)
	}

	output, err := json.Marshal(models.Json_id{Id: ai})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(output)

	return

}

// HandleUpdViolation godoc
// @Summary Update violation
// @Description update violation
// @Tags violations
// @Accept json
// @Produce  json
// @Param u body models.Violation true "Update violation"
// @Success 200 {array} models.Json_id
// @Failure 500
// @Router /violations_upd [post]
func (s *APG) HandleUpdViolation(w http.ResponseWriter, r *http.Request) {
	u := models.Violation{}
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
	err = s.Dbpool.QueryRow(context.Background(), "SELECT func_violations_upd($1,$2);", u.Id, u.ViolationName).Scan(&ui)

	if err != nil {
		log.Println("Failed execute func_violations_upd: ", err)
	}

	output, err := json.Marshal(models.Json_id{Id: ui})

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(output)

	return
}

// HandleDelViolation godoc
// @Summary Delete violations
// @Description delete violations
// @Tags violations
// @Accept json
// @Produce  json
// @Param d body models.Json_ids true "Delete violations"
// @Success 200 {array} models.Json_ids
// @Failure 500
// @Router /violations_del [post]
func (s *APG) HandleDelViolation(w http.ResponseWriter, r *http.Request) {
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
		err = s.Dbpool.QueryRow(context.Background(), "SELECT func_violations_del($1);", id).Scan(&i)
		res = append(res, i)

		if err != nil {
			log.Println("Failed execute func_violations_del: ", err)
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

// HandleGetViolation godoc
// @Summary Get violation
// @Description get violation
// @Tags violations
// @Produce  json
// @Param id path int true "Violation by id"
// @Success 200 {array} models.Violation_count
// @Failure 500
// @Router /violations/{id} [get]
func (s *APG) HandleGetViolation(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	i := vars["id"]
	g := models.Violation{}
	out_arr := []models.Violation{}

	err := s.Dbpool.QueryRow(context.Background(), "SELECT * from func_violation_get($1);", i).Scan(&g.Id, &g.ViolationName)

	if err != nil && err != pgx.ErrNoRows {
		log.Println("Failed execute from func_violation_get: ", err)
	}

	out_arr = append(out_arr, g)
	auth := models.Auth{Create: true, Read: true, Update: true, Delete: true}

	// output, err := json.Marshal(g)
	out_count, err := json.Marshal(models.Violation_count{Values: out_arr, Count: 1, Auth: auth})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(out_count)

	return

}
