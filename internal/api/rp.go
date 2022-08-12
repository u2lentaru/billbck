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

// HandleRp godoc
// @Summary List rp
// @Description Get rp list
// @Tags rp
// @Produce  json
// @Param page query int false "page number"
// @Param page_size query int false "page size"
// @Param rpname query string false "rpname search pattern"
// @Param invnumber query string false "invnumber search pattern"
// @Param ordering query string false "order by {rpname|invnumber}"
// @Param desc query boolean false "descending order {true|false}"
// @Success 200 {object} models.Rp_count
// @Failure 500
// @Router /rp [get]
func (s *APG) HandleRp(w http.ResponseWriter, r *http.Request) {
	gs := models.Rp{}
	ctx := context.Background()
	out_arr := []models.Rp{}

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
	gs1s, ok := query["rpname"]
	if ok && len(gs1s) > 0 {
		//case insensitive
		gs1 = strings.ToUpper(gs1s[0])
		//quotes
		re := regexp.MustCompile(`'`)
		gs1 = string(re.ReplaceAll([]byte(gs1), []byte("''")))
	}

	gs2 := ""
	gs2s, ok := query["invnumber"]
	if ok && len(gs2s) > 0 {
		//case insensitive
		gs2 = strings.ToUpper(gs2s[0])
		//quotes
		re := regexp.MustCompile(`'`)
		gs2 = string(re.ReplaceAll([]byte(gs2), []byte("''")))
	}

	ord := 1
	ords, ok := query["ordering"]
	if !ok || len(ords) == 0 {
		ord = 1
	} else if ords[0] == "rpname" {
		ord = 2
	} else if ords[0] == "invnumber" {
		ord = 3
	}

	dsc := false
	dscs, ok := query["desc"]
	if ok && len(dscs) > 0 {
		if !(dscs[0] == "0") {
			dsc = true
		}
	}

	rows, err := s.Dbpool.Query(ctx, "SELECT * from func_rp_get($1,$2,$3,$4,$5,$6);", pg, pgs, gs1, gs2, ord, dsc)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&gs.Id, &gs.RpName, &gs.InvNumber, &gs.InputVoltage.Id, &gs.OutputVoltage1.Id, &gs.OutputVoltage2.Id,
			&gs.Tp.Id, &gs.InputVoltage.VoltageName, &gs.InputVoltage.VoltageValue, &gs.OutputVoltage1.VoltageName,
			&gs.OutputVoltage1.VoltageValue, &gs.OutputVoltage2.VoltageName, &gs.OutputVoltage2.VoltageValue,
			&gs.Tp.TpName, &gs.Tp.GRp.Id, &gs.Tp.GRp.GRpName)
		if err != nil {
			log.Println("failed to scan row:", err)
		}

		out_arr = append(out_arr, gs)
	}

	gsc := 0
	err = s.Dbpool.QueryRow(ctx, "SELECT * from func_rp_cnt($1,$2);", gs1, gs2).Scan(&gsc)

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	auth := models.Auth{Create: true, Read: true, Update: true, Delete: true}
	out_count, err := json.Marshal(models.Rp_count{Values: out_arr, Count: gsc, Auth: auth})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Write(out_count)

	return
}

// HandleAddRp godoc
// @Summary Add rp
// @Description add rp
// @Tags rp
// @Accept json
// @Produce  json
// @Param a body models.AddRp true "New rp"
// @Success 200 {object} models.Json_id
// @Failure 500
// @Router /rp_add [post]
func (s *APG) HandleAddRp(w http.ResponseWriter, r *http.Request) {
	a := models.AddRp{}
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
	err = s.Dbpool.QueryRow(context.Background(), "SELECT func_rp_add($1,$2,$3,$4,$5,$6);", a.RpName, a.InvNumber, a.InputVoltage.Id,
		a.OutputVoltage1.Id, a.OutputVoltage2.Id, a.Tp.Id).Scan(&ai)

	if err != nil {
		log.Println("Failed execute func_rp_add: ", err)
	}

	output, err := json.Marshal(models.Json_id{Id: ai})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(output)

	return

}

// HandleUpdRp godoc
// @Summary Update rp
// @Description update rp
// @Tags rp
// @Accept json
// @Produce  json
// @Param u body models.Rp true "Update rp"
// @Success 200 {object} models.Json_id
// @Failure 500
// @Router /rp_upd [post]
func (s *APG) HandleUpdRp(w http.ResponseWriter, r *http.Request) {
	u := models.Rp{}
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
	err = s.Dbpool.QueryRow(context.Background(), "SELECT func_rp_upd($1,$2,$3,$4,$5,$6,$7);", u.Id, u.RpName, u.InvNumber,
		u.InputVoltage.Id, u.OutputVoltage1.Id, u.OutputVoltage2.Id, u.Tp.Id).Scan(&ui)

	if err != nil {
		log.Println("Failed execute func_rp_upd: ", err)
	}

	output, err := json.Marshal(models.Json_id{Id: ui})

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(output)

	return
}

// HandleDelRp godoc
// @Summary Delete rp
// @Description delete rp
// @Tags rp
// @Accept json
// @Produce  json
// @Param d body models.Json_ids true "Delete rp"
// @Success 200 {object} models.Json_ids
// @Failure 500
// @Router /rp_del [post]
func (s *APG) HandleDelRp(w http.ResponseWriter, r *http.Request) {
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
		err = s.Dbpool.QueryRow(context.Background(), "SELECT func_rp_del($1);", id).Scan(&i)
		res = append(res, i)

		if err != nil {
			log.Println("Failed execute func_rp_del: ", err)
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

// HandleGetRp godoc
// @Summary Get rp
// @Description get rp
// @Tags rp
// @Produce  json
// @Param id path int true "Rp by id"
// @Success 200 {object} models.Rp_count
// @Failure 500
// @Router /rp/{id} [get]
func (s *APG) HandleGetRp(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	i := vars["id"]
	g := models.Rp{}
	out_arr := []models.Rp{}

	err := s.Dbpool.QueryRow(context.Background(), "SELECT * from func_rp_getbyid($1);", i).Scan(&g.Id, &g.RpName, &g.InvNumber,
		&g.InputVoltage.Id, &g.OutputVoltage1.Id, &g.OutputVoltage2.Id, &g.Tp.Id, &g.InputVoltage.VoltageName,
		&g.InputVoltage.VoltageValue, &g.OutputVoltage1.VoltageName, &g.OutputVoltage1.VoltageValue, &g.OutputVoltage2.VoltageName,
		&g.OutputVoltage2.VoltageValue, &g.Tp.TpName, &g.Tp.GRp.Id, &g.Tp.GRp.GRpName)

	if err != nil && err != pgx.ErrNoRows {
		log.Println("Failed execute from func_rp_getbyid: ", err)
	}

	out_arr = append(out_arr, g)
	auth := models.Auth{Create: true, Read: true, Update: true, Delete: true}

	// output, err := json.Marshal(g)
	out_count, err := json.Marshal(models.Rp_count{Values: out_arr, Count: 1, Auth: auth})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(out_count)

	return
}
