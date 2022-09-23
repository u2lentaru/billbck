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
	"github.com/u2lentaru/billbck/internal/adapters/db/pgsql"
	"github.com/u2lentaru/billbck/internal/models"
	"github.com/u2lentaru/billbck/internal/services"
	"github.com/u2lentaru/billbck/internal/utils"
)

type ifChargeService interface {
	GetList(ctx context.Context, pg, pgs int, gs1, gs2 string, ord int, dsc bool) (models.Charge_count, error)
	Add(ctx context.Context, ea models.Charge) (int, error)
	Upd(ctx context.Context, eu models.Charge) (int, error)
	Del(ctx context.Context, ed []int) ([]int, error)
	GetOne(ctx context.Context, i int) (models.Charge_count, error)
	ChargeRun(ctx context.Context, i int) (int, error)
}

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
func HandleCharges(w http.ResponseWriter, r *http.Request) {
	var gs ifChargeService
	gs = services.NewChargeService(pgsql.ChargeStorage{})
	ctx := context.Background()
	auth := utils.GetAuth(r)

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

	out_arr, err := gs.GetList(ctx, pg, pgs, gs1, gs2, ord, dsc)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	out_arr.Auth = auth
	out_count, err := json.Marshal(out_arr)
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
func HandleAddCharge(w http.ResponseWriter, r *http.Request) {
	var gs ifChargeService
	gs = services.NewChargeService(pgsql.ChargeStorage{})
	ctx := context.Background()

	a := models.Charge{}
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

	ai, err := gs.Add(ctx, a)

	if err != nil {
		log.Println("Failed execute ifChargeService.Add: ", err)
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
func HandleUpdCharge(w http.ResponseWriter, r *http.Request) {
	var gs ifChargeService
	gs = services.NewChargeService(pgsql.ChargeStorage{})
	ctx := context.Background()

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

	ui, err := gs.Upd(ctx, u)

	if err != nil {
		log.Println("Failed execute ifChargeService.Upd: ", err)
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
func HandleDelCharge(w http.ResponseWriter, r *http.Request) {
	var gs ifChargeService
	gs = services.NewChargeService(pgsql.ChargeStorage{})
	ctx := context.Background()

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

	res, err := gs.Del(ctx, d.Ids)
	if err != nil {
		log.Println("Failed execute ifChargeService.Del: ", err)
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
func HandleGetCharge(w http.ResponseWriter, r *http.Request) {
	var gs ifChargeService
	gs = services.NewChargeService(pgsql.ChargeStorage{})
	ctx := context.Background()
	auth := utils.GetAuth(r)

	vars := mux.Vars(r)
	i, err := strconv.Atoi(vars["id"])
	if err != nil {
		i = 0
	}

	out_arr, err := gs.GetOne(ctx, i)
	if err != nil {
		log.Println("Failed execute ifChargeService.GetOne: ", err)
	}

	out_arr.Auth = auth
	out_count, err := json.Marshal(out_arr)
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
	var gs ifChargeService
	gs = services.NewChargeService(pgsql.ChargeStorage{})
	ctx := context.Background()

	vars := mux.Vars(r)
	i, err := strconv.Atoi(vars["id"])
	if err != nil {
		i = 0
	}

	pr, err := gs.ChargeRun(ctx, i)
	if err != nil {
		log.Println("Failed execute ifChargeService.ChargeRun: ", err)
	}

	output, err := json.Marshal(models.Json_id{Id: pr})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(output)

	return
}
