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
	"github.com/u2lentaru/billbck/internal/utils"
)

// HandlePeriods godoc
// @Summary List periods
// @Description get period list
// @Tags periods
// @Produce json
// @Param page query int false "page number"
// @Param page_size query int false "page size"
// @Param periodname query string false "periodname search pattern"
// @Param ordering query string false "order by {id|periodname|startdate}"
// @Param desc query boolean false "descending order {true|false}"
// @Success 200 {array} models.Period_count
// @Failure 500
// @Router /periods [get]
func (s *APG) HandlePeriods(w http.ResponseWriter, r *http.Request) {
	gs := models.Period{}
	ctx := context.Background()
	out_arr := []models.Period{}

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
	gs1s, ok := query["periodname"]
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
	} else if ords[0] == "periodname" {
		ord = 2
	} else if ords[0] == "stratdate" {
		ord = 4
	}

	dsc := false
	dscs, ok := query["desc"]
	if ok && len(dscs) > 0 {
		if !(dscs[0] == "0") {
			dsc = true
		}
	}

	rows, err := s.Dbpool.Query(ctx, "SELECT * from func_periods_get($1,$2,$3,$4,$5);", pg, pgs, gs1, ord, dsc)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&gs.Id, &gs.PeriodName, &gs.Startdate, &gs.Enddate)
		if err != nil {
			log.Println("failed to scan row:", err)
		}

		out_arr = append(out_arr, gs)
	}

	gsc := 0
	err = s.Dbpool.QueryRow(ctx, "SELECT * from func_periods_cnt($1);", gs1).Scan(&gsc)

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// auth := models.Auth{Create: true, Read: true, Update: true, Delete: true}
	out_count, err := json.Marshal(models.Period_count{Values: out_arr, Count: gsc, Auth: utils.GetAuth(r, "periods")})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Write(out_count)

	return
}

// HandleAddPeriod godoc
// @Summary Add period
// @Description add period
// @Tags periods
// @Accept json
// @Produce  json
// @Param a body models.AddPeriod true "New period. Significant params: PeriodName, Startdate, Enddate"
// @Success 200 {array} models.Json_id
// @Failure 500
// @Router /periods_add [post]
func (s *APG) HandleAddPeriod(w http.ResponseWriter, r *http.Request) {
	a := models.AddPeriod{}
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
	err = s.Dbpool.QueryRow(context.Background(), "SELECT func_periods_add($1,$2,$3);", a.PeriodName, a.Startdate, a.Enddate).Scan(&ai)

	if err != nil {
		log.Println("Failed execute func_periods_add: ", err)
	}

	output, err := json.Marshal(models.Json_id{Id: ai})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(output)

	return

}

// HandleUpdPeriod godoc
// @Summary Update period
// @Description update period
// @Tags periods
// @Accept json
// @Produce  json
// @Param u body models.Period true "Update period. Significant params: Id, PeriodName, Startdate, Enddate"
// @Success 200 {array} models.Json_id
// @Failure 500
// @Router /periods_upd [post]
func (s *APG) HandleUpdPeriod(w http.ResponseWriter, r *http.Request) {
	u := models.Period{}
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
	err = s.Dbpool.QueryRow(context.Background(), "SELECT func_periods_upd($1,$2,$3,$4);", u.Id, u.PeriodName, u.Startdate, u.Enddate).Scan(&ui)

	if err != nil {
		log.Println("Failed execute func_periods_upd: ", err)
	}

	output, err := json.Marshal(models.Json_id{Id: ui})

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(output)

	return
}

// HandleDelPeriod godoc
// @Summary Delete periods
// @Description delete periods
// @Tags periods
// @Accept json
// @Produce  json
// @Param d body models.Json_ids true "Delete periods"
// @Success 200 {array} models.Json_ids
// @Failure 500
// @Router /periods_del [post]
func (s *APG) HandleDelPeriod(w http.ResponseWriter, r *http.Request) {
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
		err = s.Dbpool.QueryRow(context.Background(), "SELECT func_periods_del($1);", id).Scan(&i)
		res = append(res, i)

		if err != nil {
			log.Println("Failed execute func_periods_del: ", err)
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

// HandleGetPeriod godoc
// @Summary Get period
// @Security ApiKeyAuth
// @Description get period
// @Tags periods
// @Produce  json
// @Param id path int true "Period by id"
// @Success 200 {array} models.Period_count
// @Failure 500
// @Router /periods/{id} [get]
func (s *APG) HandleGetPeriod(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	i := vars["id"]
	g := models.Period{}
	out_arr := []models.Period{}

	err := s.Dbpool.QueryRow(context.Background(), "SELECT * from func_period_get($1);", i).Scan(&g.Id, &g.PeriodName, &g.Startdate,
		&g.Enddate)

	if err != nil && err != pgx.ErrNoRows {
		log.Println("Failed execute from func_period_get: ", err)
	}

	out_arr = append(out_arr, g)
	// auth := models.Auth{Create: true, Read: true, Update: true, Delete: true}
	// auth := utils.GetAuth(r)

	// output, err := json.Marshal(g)
	out_count, err := json.Marshal(models.Period_count{Values: out_arr, Count: 1, Auth: utils.GetAuth(r, "periods")})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(out_count)

	return

}
