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

// HandleTransPwr godoc
// @Summary List power transformers
// @Description get power transformers list
// @Tags transpwr
// @Produce json
// @Param page query int false "page number"
// @Param page_size query int false "page size"
// @Param transpwrname query string false "transpwrname search pattern"
// @Param ordering query string false "order by {id|transpwrname}"
// @Param desc query boolean false "descending order {true|false}"
// @Success 200 {object} models.TransPwr_count
// @Failure 500
// @Router /transpwr [get]
func (s *APG) HandleTransPwr(w http.ResponseWriter, r *http.Request) {
	gs := models.TransPwr{}
	ctx := context.Background()
	out_arr := []models.TransPwr{}

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
	gs1s, ok := query["transpwrname"]
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
	} else if ords[0] == "transpwrname" {
		ord = 2
	}

	dsc := false
	dscs, ok := query["desc"]
	if ok && len(dscs) > 0 {
		if !(dscs[0] == "0") {
			dsc = true
		}
	}

	rows, err := s.Dbpool.Query(ctx, "SELECT * from func_trans_pwr_get($1,$2,$3,$4,$5);", pg, pgs, gs1, ord, dsc)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&gs.Id, &gs.TransPwrName, &gs.TransPwrType.Id, &gs.TransPwrType.TransPwrTypeName, &gs.TransPwrType.ShortCircuitPower,
			&gs.TransPwrType.IdlingLossPower, &gs.TransPwrType.NominalPower)
		if err != nil {
			log.Println("failed to scan row:", err)
		}

		out_arr = append(out_arr, gs)
	}

	gsc := 0
	err = s.Dbpool.QueryRow(ctx, "SELECT * from func_trans_pwr_cnt($1);", gs1).Scan(&gsc)

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	auth := models.Auth{Create: true, Read: true, Update: true, Delete: true}
	out_count, err := json.Marshal(models.TransPwr_count{Values: out_arr, Count: gsc, Auth: auth})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Write(out_count)

	return
}

// HandleAddTransPwr godoc
// @Summary Add power transformer
// @Description add power transformer
// @Tags transpwr
// @Accept json
// @Produce  json
// @Param a body models.AddTransPwr true "New power transformer. Significant params: TransPwrName, TransPwrType.Id"
// @Success 200 {object} models.Json_id
// @Failure 500
// @Router /transpwr_add [post]
func (s *APG) HandleAddTransPwr(w http.ResponseWriter, r *http.Request) {
	a := models.AddTransPwr{}
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
	err = s.Dbpool.QueryRow(context.Background(), "SELECT func_trans_pwr_add($1,$2);", a.TransPwrName,
		a.TransPwrType.Id).Scan(&ai)

	if err != nil {
		log.Println("Failed execute func_trans_pwr_add: ", err)
	}

	output, err := json.Marshal(models.Json_id{Id: ai})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(output)

	return

}

// HandleUpdTransPwr godoc
// @Summary Update power transformer
// @Description update power transformer
// @Tags transpwr
// @Accept json
// @Produce  json
// @Param u body models.TransPwr true "Update power transformer. Significant params: Id, TransPwrName, TransPwrType.Id"
// @Success 200 {object} models.Json_id
// @Failure 500
// @Router /transpwr_upd [post]
func (s *APG) HandleUpdTransPwr(w http.ResponseWriter, r *http.Request) {
	u := models.TransPwr{}
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
	err = s.Dbpool.QueryRow(context.Background(), "SELECT func_trans_pwr_upd($1,$2,$3);", u.Id, u.TransPwrName,
		u.TransPwrType.Id).Scan(&ui)

	if err != nil {
		log.Println("Failed execute func_trans_pwr_upd: ", err)
	}

	output, err := json.Marshal(models.Json_id{Id: ui})

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(output)

	return

}

// HandleDelTransPwr godoc
// @Summary Delete power transformers
// @Description delete power transformers
// @Tags transpwr
// @Accept json
// @Produce  json
// @Param d body models.Json_ids true "Delete power transformers"
// @Success 200 {object} models.Json_ids
// @Failure 500
// @Router /transpwr_del [post]
func (s *APG) HandleDelTransPwr(w http.ResponseWriter, r *http.Request) {
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
		err = s.Dbpool.QueryRow(context.Background(), "SELECT func_trans_pwr_del($1);", id).Scan(&i)
		res = append(res, i)

		if err != nil {
			log.Println("Failed execute func_trans_pwr_del: ", err)
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

// HandleGetTransPwr godoc
// @Summary Get power transformer
// @Description get power transformer
// @Tags transpwr
// @Produce  json
// @Param id path int true "Power transformer by id"
// @Success 200 {object} models.TransPwr_count
// @Failure 500
// @Router /transpwr/{id} [get]
func (s *APG) HandleGetTransPwr(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	i := vars["id"]
	g := models.TransPwr{}
	out_arr := []models.TransPwr{}

	err := s.Dbpool.QueryRow(context.Background(), "SELECT * from func_trans_pwr_getbyid($1);", i).Scan(&g.Id, &g.TransPwrName,
		&g.TransPwrType.Id, &g.TransPwrType.TransPwrTypeName, &g.TransPwrType.ShortCircuitPower, &g.TransPwrType.IdlingLossPower,
		&g.TransPwrType.NominalPower)

	if err != nil && err != pgx.ErrNoRows {
		log.Println("Failed execute from func_trans_pwr_getbyid: ", err)
	}

	out_arr = append(out_arr, g)
	auth := models.Auth{Create: true, Read: true, Update: true, Delete: true}

	out_count, err := json.Marshal(models.TransPwr_count{Values: out_arr, Count: 1, Auth: auth})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Write(out_count)

	return

}
