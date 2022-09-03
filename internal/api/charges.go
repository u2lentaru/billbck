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

// HandleCharges godoc
// @Summary List charges
// @Description get charge list
// @Tags charges
// @Produce  json
// @Param page query int false "page number"
// @Param page_size query int false "page size"
// @Param contractnumber query string false "contractnumber search pattern"
// @Param oid query string false "object id"
// @Param ordering query string false "order by {id|chargedate}"
// @Param desc query boolean false "descending order {true|false}"
// @Success 200 {object} models.Charge_count
// @Failure 500
// @Router /charges [get]
func (s *APG) HandleCharges(w http.ResponseWriter, r *http.Request) {
	gs := models.Charge{}
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
	gs1s, ok := query["contractnumber"]
	if ok && len(gs1s) > 0 {
		//case insensitive
		gs1 = strings.ToUpper(gs1s[0])
		//quotes
		re := regexp.MustCompile(`'`)
		gs1 = string(re.ReplaceAll([]byte(gs1), []byte("''")))
	}

	gs2 := ""
	gs2s, ok := query["oid"]
	if ok && len(gs2s) > 0 {
		_, err := strconv.Atoi(gs2s[0])
		if err == nil {
			gs2 = gs2s[0]
		}
	}

	gsc := 0
	err := s.Dbpool.QueryRow(ctx, "SELECT * from func_charges_cnt($1,$2);", gs1, NullableString(gs2)).Scan(&gsc)

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	out_arr := make([]models.Charge, 0,
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
	} else if ords[0] == "chargedate" {
		ord = 2
	}

	dsc := false
	dscs, ok := query["desc"]
	if ok && len(dscs) > 0 {
		if !(dscs[0] == "0") {
			dsc = true
		}
	}

	rows, err := s.Dbpool.Query(ctx, "SELECT * from func_charges_get($1,$2,$3,$4,$5,$6);", pg, pgs, gs1, NullableString(gs2), ord, dsc)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&gs.Id, &gs.ChargeDate, &gs.Contract.Id, &gs.Object.Id, &gs.ObjTypeId, &gs.Pu.Id, &gs.ChargeType.Id, &gs.Qty,
			&gs.TransLoss, &gs.Lineloss, &gs.Startdate, &gs.Enddate, &gs.Contract.ContractNumber, &gs.Object.ObjectName,
			&gs.Pu.PuNumber, &gs.ChargeType.ChargeTypeName)
		if err != nil {
			log.Println("failed to scan row:", err)
		}

		out_arr = append(out_arr, gs)
	}

	auth := models.Auth{Create: true, Read: true, Update: true, Delete: true}
	out_count, err := json.Marshal(models.Charge_count{Values: out_arr, Count: gsc, Auth: auth})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Write(out_count)

	return
}

// HandleAddCharge godoc
// @Summary Add charge
// @Description add charge
// @Tags charges
// @Accept json
// @Produce  json
// @Param a body models.AddCharge true "New charge. Significant params: ChargeDate, Contract.Id, Object.Id, ObjTypeId, Pu.Id, ChargeType.Id, Qty, TransLoss, Lineloss, Startdate, Enddate"
// @Success 200 {object} models.Json_id
// @Failure 500
// @Router /charges_add [post]
func (s *APG) HandleAddCharge(w http.ResponseWriter, r *http.Request) {
	a := models.AddCharge{}
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
	err = s.Dbpool.QueryRow(context.Background(), "SELECT func_charges_add($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11);", a.ChargeDate, a.Contract.Id,
		a.Object.Id, a.ObjTypeId, a.Pu.Id, a.ChargeType.Id, a.Qty, a.TransLoss, a.Lineloss, a.Startdate, a.Enddate).Scan(&ai)

	if err != nil {
		log.Println("Failed execute func_charges_add: ", err)
	}

	output, err := json.Marshal(models.Json_id{Id: ai})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(output)

	return

}

// HandleUpdCharge godoc
// @Summary Update charge
// @Description update charge
// @Tags charges
// @Accept json
// @Produce  json
// @Param u body models.Charge true "Update charge. Significant params: Id, ChargeDate, Contract.Id, Object.Id, ObjTypeId, Pu.Id, ChargeType.Id, Qty, TransLoss, Lineloss, Startdate, Enddate"
// @Success 200 {object} models.Json_id
// @Failure 500
// @Router /charges_upd [post]
func (s *APG) HandleUpdCharge(w http.ResponseWriter, r *http.Request) {
	u := models.Charge{}
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
	err = s.Dbpool.QueryRow(context.Background(), "SELECT func_charges_upd($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12);", u.Id, u.ChargeDate,
		u.Contract.Id, u.Object.Id, u.ObjTypeId, u.Pu.Id, u.ChargeType.Id, u.Qty, u.TransLoss, u.Lineloss, u.Startdate, u.Enddate).Scan(&ui)

	if err != nil {
		log.Println("Failed execute func_charges_upd: ", err)
	}

	output, err := json.Marshal(models.Json_id{Id: ui})

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(output)

	return

}

// HandleDelCharge godoc
// @Summary Delete charges
// @Description delete charges
// @Tags charges
// @Accept json
// @Produce  json
// @Param d body models.Json_ids true "Delete charges"
// @Success 200 {object} models.Json_ids
// @Failure 500
// @Router /charges_del [post]
func (s *APG) HandleDelCharge(w http.ResponseWriter, r *http.Request) {
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
		err = s.Dbpool.QueryRow(context.Background(), "SELECT func_charges_del($1);", id).Scan(&i)
		res = append(res, i)

		if err != nil {
			log.Println("Failed execute func_charges_del: ", err)
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

// HandleGetCharge godoc
// @Summary Get charge
// @Description get charge
// @Tags charges
// @Produce  json
// @Param id path int true "Charge by id"
// @Success 200 {object} models.Charge_count
// @Failure 500
// @Router /charges/{id} [get]
func (s *APG) HandleGetCharge(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	i := vars["id"]
	g := models.Charge{}
	out_arr := []models.Charge{}

	err := s.Dbpool.QueryRow(context.Background(), "SELECT * from func_charge_get($1);", i).Scan(&g.Id, &g.ChargeDate, &g.Contract.Id,
		&g.Object.Id, &g.ObjTypeId, &g.Pu.Id, &g.ChargeType.Id, &g.Qty, &g.TransLoss, &g.Lineloss, &g.Startdate, &g.Enddate,
		&g.Contract.ContractNumber, &g.Object.ObjectName, &g.Pu.PuNumber, &g.ChargeType.ChargeTypeName)

	if err != nil && err != pgx.ErrNoRows {
		log.Println("Failed execute from func_charge_get: ", err)
	}

	out_arr = append(out_arr, g)
	auth := models.Auth{Create: true, Read: true, Update: true, Delete: true}

	// output, err := json.Marshal(g)
	out_count, err := json.Marshal(models.Charge_count{Values: out_arr, Count: 1, Auth: auth})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(out_count)

	return
}

// HandleChargeRun godoc
// @Summary Charge run by period id
// @Description charge run
// @Tags charges
// @Produce  json
// @Param id path int true "Charge run by period id"
// @Success 200 {object} models.Json_id
// @Failure 500
// @Router /charges_run/{id} [get]
func (s *APG) HandleChargeRun(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	i := vars["id"]

	rc := 0
	err := s.Dbpool.QueryRow(context.Background(), "SELECT * from func_charges_run($1);", i).Scan(&rc)

	if err != nil {
		log.Println("Failed execute func_charges_run: ", err)
	}

	output, err := json.Marshal(models.Json_id{Id: rc})

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(output)

	return
}
