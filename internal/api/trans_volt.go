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

// HandleTransVolt godoc
// @Summary List voltage transformers
// @Description get voltage transformers list
// @Tags transvolt
// @Produce json
// @Param page query int false "page number"
// @Param page_size query int false "page size"
// @Param transvoltname query string false "transvoltname search pattern"
// @Param ordering query string false "order by {id|transvoltname}"
// @Param desc query boolean false "descending order {true|false}"
// @Success 200 {object} models.TransVolt_count
// @Failure 500
// @Router /transvolt [get]
func (s *APG) HandleTransVolt(w http.ResponseWriter, r *http.Request) {
	gs := models.TransVolt{}
	ctx := context.Background()
	out_arr := []models.TransVolt{}

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
	gs1s, ok := query["transvoltname"]
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
	} else if ords[0] == "transvoltname" {
		ord = 2
	}

	dsc := false
	dscs, ok := query["desc"]
	if ok && len(dscs) > 0 {
		if !(dscs[0] == "0") {
			dsc = true
		}
	}

	rows, err := s.Dbpool.Query(ctx, "SELECT * from func_trans_volt_get($1,$2,$3,$4,$5);", pg, pgs, gs1, ord, dsc)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&gs.Id, &gs.TransVoltName, &gs.TransType.Id, &gs.CheckDate, &gs.NextCheckDate, &gs.ProdDate, &gs.Serial1,
			&gs.Serial2, &gs.Serial3, &gs.TransType.TransTypeName, &gs.TransType.Ratio, &gs.TransType.Class, &gs.TransType.MaxCurr,
			&gs.TransType.NomCurr)
		if err != nil {
			log.Println("failed to scan row:", err)
		}

		out_arr = append(out_arr, gs)
	}

	gsc := 0
	err = s.Dbpool.QueryRow(ctx, "SELECT * from func_trans_volt_cnt($1);", gs1).Scan(&gsc)

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	auth := models.Auth{Create: true, Read: true, Update: true, Delete: true}
	out_count, err := json.Marshal(models.TransVolt_count{Values: out_arr, Count: gsc, Auth: auth})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Write(out_count)

	return

}

// HandleAddTransVolt godoc
// @Summary Add voltage transformer
// @Description add voltage transformer
// @Tags transvolt
// @Accept json
// @Produce  json
// @Param a body models.AddTransVolt true "New voltage transformer. Significant params: TransVoltName, TransType.Id, CheckDate(n), NextCheckDate(n), ProdDate(n), Serial1(n), Serial2(n), Serial3(n)"
// @Success 200 {object} models.Json_id
// @Failure 500
// @Router /transvolt_add [post]
func (s *APG) HandleAddTransVolt(w http.ResponseWriter, r *http.Request) {
	a := models.AddTransVolt{}
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
	err = s.Dbpool.QueryRow(context.Background(), "SELECT func_trans_volt_add($1,$2,$3,$4,$5,$6,$7,$8);", a.TransVoltName,
		a.TransType.Id, a.CheckDate, a.NextCheckDate, a.ProdDate, a.Serial1, a.Serial2, a.Serial3).Scan(&ai)

	if err != nil {
		log.Println("Failed execute func_trans_volt_add: ", err)
	}

	output, err := json.Marshal(models.Json_id{Id: ai})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(output)

	return

}

// HandleUpdTransVolt godoc
// @Summary Update voltage transformer
// @Description update voltage transformer
// @Tags transvolt
// @Accept json
// @Produce  json
// @Param u body models.TransVolt true "Update voltage transformer. Significant params: Id, TransVoltName, TransType.Id, CheckDate(n), NextCheckDate(n), ProdDate(n), Serial1(n), Serial2(n), Serial3(n)"
// @Success 200 {object} models.Json_id
// @Failure 500
// @Router /transvolt_upd [post]
func (s *APG) HandleUpdTransVolt(w http.ResponseWriter, r *http.Request) {
	u := models.TransVolt{}
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
	err = s.Dbpool.QueryRow(context.Background(), "SELECT func_trans_volt_upd($1,$2,$3,$4,$5,$6,$7,$8,$9);", u.Id, u.TransVoltName,
		u.TransType.Id, u.CheckDate, u.NextCheckDate, u.ProdDate, u.Serial1, u.Serial2, u.Serial3).Scan(&ui)

	if err != nil {
		log.Println("Failed execute func_trans_volt_upd: ", err)
	}

	output, err := json.Marshal(models.Json_id{Id: ui})

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(output)

	return

}

// HandleDelTransVolt godoc
// @Summary Delete voltage transformers
// @Description delete voltage transformers
// @Tags transvolt
// @Accept json
// @Produce  json
// @Param d body models.Json_ids true "Delete voltage transformers"
// @Success 200 {object} models.Json_ids
// @Failure 500
// @Router /transvolt_del [post]
func (s *APG) HandleDelTransVolt(w http.ResponseWriter, r *http.Request) {
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
		err = s.Dbpool.QueryRow(context.Background(), "SELECT func_trans_volt_del($1);", id).Scan(&i)
		res = append(res, i)

		if err != nil {
			log.Println("Failed execute func_trans_volt_del: ", err)
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

// HandleGetTransVolt godoc
// @Summary Get voltage transformer
// @Description get voltage transformer
// @Tags transvolt
// @Produce  json
// @Param id path int true "Voltage transformer by id"
// @Success 200 {object} models.TransVolt_count
// @Failure 500
// @Router /transvolt/{id} [get]
func (s *APG) HandleGetTransVolt(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	i := vars["id"]
	g := models.TransVolt{}
	out_arr := []models.TransVolt{}

	err := s.Dbpool.QueryRow(context.Background(), "SELECT * from func_trans_volt_getbyid($1);", i).Scan(&g.Id, &g.TransVoltName,
		&g.TransType.Id, &g.CheckDate, &g.NextCheckDate, &g.ProdDate, &g.Serial1, &g.Serial2, &g.Serial3, &g.TransType.TransTypeName,
		&g.TransType.Ratio, &g.TransType.Class, &g.TransType.MaxCurr, &g.TransType.NomCurr)

	if err != nil && err != pgx.ErrNoRows {
		log.Println("Failed execute from func_trans_volt_getbyid: ", err)
	}

	out_arr = append(out_arr, g)
	auth := models.Auth{Create: true, Read: true, Update: true, Delete: true}

	out_count, err := json.Marshal(models.TransVolt_count{Values: out_arr, Count: 1, Auth: auth})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Write(out_count)

	return

}
