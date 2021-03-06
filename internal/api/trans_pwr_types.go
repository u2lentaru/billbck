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

// HandleTransPwrTypes godoc
// @Summary List transpwrtypes
// @Description get transpwrtype list
// @Tags transpwrtypes
// @Produce json
// @Param page query int false "page number"
// @Param page_size query int false "page size"
// @Param transpwrtypename query string false "transpwrtypename search pattern"
// @Param ordering query string false "order by {id|transpwrtypename}"
// @Param desc query boolean false "descending order {true|false}"
// @Success 200 {array} models.TransPwrType_count
// @Failure 500
// @Router /transpwrtypes [get]
func (s *APG) HandleTransPwrTypes(w http.ResponseWriter, r *http.Request) {
	gs := models.TransPwrType{}
	ctx := context.Background()
	out_arr := []models.TransPwrType{}

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
	gs1s, ok := query["transpwrtypename"]
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
	} else if ords[0] == "transpwrtypename" {
		ord = 2
	}

	dsc := false
	dscs, ok := query["desc"]
	if ok && len(dscs) > 0 {
		if !(dscs[0] == "0") {
			dsc = true
		}
	}

	rows, err := s.Dbpool.Query(ctx, "SELECT * from func_trans_pwr_types_get($1,$2,$3,$4,$5);", pg, pgs, gs1, ord, dsc)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&gs.Id, &gs.TransPwrTypeName, &gs.ShortCircuitPower, &gs.IdlingLossPower, &gs.NominalPower)
		if err != nil {
			log.Println("failed to scan row:", err)
		}

		out_arr = append(out_arr, gs)
	}

	gsc := 0
	err = s.Dbpool.QueryRow(ctx, "SELECT * from func_trans_pwr_types_cnt($1);", gs1).Scan(&gsc)

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	auth := models.Auth{Create: true, Read: true, Update: true, Delete: true}
	out_count, err := json.Marshal(models.TransPwrType_count{Values: out_arr, Count: gsc, Auth: auth})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Write(out_count)

	return
}

// HandleAddTransPwrType godoc
// @Summary Add transpwrtype
// @Description add transpwrtype
// @Tags transpwrtypes
// @Accept json
// @Produce  json
// @Param a body models.AddTransPwrType true "New transpwrtype. Significant params: TransPwrTypeName, ShortCircuitPower, IdlingLossPower, NominalPower"
// @Success 200 {array} models.Json_id
// @Failure 500
// @Router /transpwrtypes_add [post]
func (s *APG) HandleAddTransPwrType(w http.ResponseWriter, r *http.Request) {
	a := models.AddTransPwrType{}
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
	err = s.Dbpool.QueryRow(context.Background(), "SELECT func_trans_pwr_types_add($1,$2,$3,$4);", a.TransPwrTypeName,
		a.ShortCircuitPower, a.IdlingLossPower, a.NominalPower).Scan(&ai)

	if err != nil {
		log.Println("Failed execute func_trans_pwr_types_add: ", err)
	}

	output, err := json.Marshal(models.Json_id{Id: ai})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(output)

	return
}

// HandleUpdTransPwrType godoc
// @Summary Update transpwrtype
// @Description update transpwrtype
// @Tags transpwrtypes
// @Accept json
// @Produce  json
// @Param u body models.TransPwrType true "Update transpwrtype. Significant params: Id, TransPwrTypeName, ShortCircuitPower, IdlingLossPower, NominalPower"
// @Success 200 {array} models.Json_id
// @Failure 500
// @Router /transpwrtypes_upd [post]
func (s *APG) HandleUpdTransPwrType(w http.ResponseWriter, r *http.Request) {
	u := models.TransPwrType{}
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
	err = s.Dbpool.QueryRow(context.Background(), "SELECT func_trans_pwr_types_upd($1,$2,$3,$4,$5);", u.Id, u.TransPwrTypeName,
		u.ShortCircuitPower, u.IdlingLossPower, u.NominalPower).Scan(&ui)

	if err != nil {
		log.Println("Failed execute func_trans_pwr_types_upd: ", err)
	}

	output, err := json.Marshal(models.Json_id{Id: ui})

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(output)

	return
}

// HandleDelTransPwrType godoc
// @Summary Delete transpwrtypes
// @Description delete transpwrtypes
// @Tags transpwrtypes
// @Accept json
// @Produce  json
// @Param d body models.Json_ids true "Delete transpwrtypes"
// @Success 200 {array} models.Json_ids
// @Failure 500
// @Router /transpwrtypes_del [post]
func (s *APG) HandleDelTransPwrType(w http.ResponseWriter, r *http.Request) {
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
		err = s.Dbpool.QueryRow(context.Background(), "SELECT func_trans_pwr_types_del($1);", id).Scan(&i)
		res = append(res, i)

		if err != nil {
			log.Println("Failed execute func_trans_pwr_types_del: ", err)
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

// HandleGetTransPwrType godoc
// @Summary Get transpwrtype
// @Description get transpwrtype
// @Tags transpwrtypes
// @Produce  json
// @Param id path int true "TransPwrType by id"
// @Success 200 {array} models.TransPwrType_count
// @Failure 500
// @Router /transpwrtypes/{id} [get]
func (s *APG) HandleGetTransPwrType(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	i := vars["id"]
	g := models.TransPwrType{}
	out_arr := []models.TransPwrType{}

	err := s.Dbpool.QueryRow(context.Background(), "SELECT * from func_trans_pwr_type_get($1);", i).Scan(&g.Id, &g.TransPwrTypeName,
		&g.ShortCircuitPower, &g.IdlingLossPower, &g.NominalPower)

	if err != nil && err != pgx.ErrNoRows {
		log.Println("Failed execute from func_trans_pwr_type_get: ", err)
	}

	out_arr = append(out_arr, g)
	auth := models.Auth{Create: true, Read: true, Update: true, Delete: true}

	out_count, err := json.Marshal(models.TransPwrType_count{Values: out_arr, Count: 1, Auth: auth})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Write(out_count)

	return

}
