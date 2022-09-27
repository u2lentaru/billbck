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
	"time"

	"github.com/gorilla/mux"
	"github.com/u2lentaru/billbck/internal/adapters/db/pgsql"
	"github.com/u2lentaru/billbck/internal/models"
	"github.com/u2lentaru/billbck/internal/services"
	"github.com/u2lentaru/billbck/internal/utils"
)

type ifObjContractService interface {
	GetList(ctx context.Context, pg, pgs, gs1, gs2, gs3 int, gs4, gs4f bool, ord int, dsc bool) (models.ObjContract_count, error)
	Add(ctx context.Context, ea models.ObjContract) (int, error)
	Upd(ctx context.Context, eu models.ObjContract) (int, error)
	Del(ctx context.Context, ed models.IdClose) (int, error)
	GetOne(ctx context.Context, i int, d string) (models.ObjContract_count, error)
}

// HandleObjContracts godoc
// @Summary List objcontracts
// @Description get objcontract list
// @Tags objcontracts
// @Produce  json
// @Param page query int false "page number"
// @Param page_size query int false "page size"
// @Param objectid query string false "object id"
// @Param tid query string false "type id"
// @Param contractid query string false "contract id"
// @Param active query boolean false "enddate is null"
// @Param ordering query string false "order by {object|contract|startdate}"
// @Param desc query boolean false "descending order {true|false}"
// @Success 200 {object} models.ObjContract_count
// @Failure 500
// @Router /objcontracts [get]
func HandleObjContracts(w http.ResponseWriter, r *http.Request) {
	var gs ifObjContractService
	gs = services.NewObjContractService(pgsql.ObjContractStorage{})
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

	gs1 := 0
	gs1s, ok := query["objectid"]
	if ok && len(gs1s) > 0 {
		if gs1t, err := strconv.Atoi(gs1s[0]); err != nil {
			gs1 = 0
		} else {
			gs1 = gs1t
		}
	}

	gs2 := 0
	gs2s, ok := query["tid"]
	if ok && len(gs2s) > 0 {
		if gs2t, err := strconv.Atoi(gs2s[0]); err != nil {
			gs2 = 0
		} else {
			gs1 = gs2t
		}
	}

	gs3 := 0
	gs3s, ok := query["contractid"]
	if ok && len(gs3s) > 0 {
		if gs3t, err := strconv.Atoi(gs3s[0]); err != nil {
			gs3 = 0
		} else {
			gs3 = gs3t
		}
	}

	gs4 := false
	gs4f := true
	gs4s, ok := query["active"]
	if ok && len(gs4s) > 0 {
		if gs4s[0] == "true" || gs4s[0] == "false" {
			gs4, _ = strconv.ParseBool(gs4s[0])
		} else {
			gs4f = false
		}
	} else {
		gs4f = false
	}

	ord := 1
	ords, ok := query["ordering"]
	if !ok || len(ords) == 0 {
		ord = 1
	} else if ords[0] == "object" {
		ord = 6
	} else if ords[0] == "contract" {
		ord = 17
	} else if ords[0] == "startdate" {
		ord = 4
	}

	dsc := false
	dscs, ok := query["desc"]
	if ok && len(dscs) > 0 {
		if !(dscs[0] == "0") {
			dsc = true
		}
	}

	out_arr, err := gs.GetList(ctx, pg, pgs, gs1, gs2, gs3, gs4, gs4f, ord, dsc)
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

// HandleAddObjContract godoc
// @Summary Add objcontract
// @Description add objcontract
// @Tags objcontracts
// @Accept json
// @Produce  json
// @Param a body models.AddObjContract true "New objcontract. Old objcontract of the object will be closed. Significant params: Contract.Id, Object.Id, ObjTypeId, Startdate"
// @Success 200 {object} models.Json_id
// @Failure 500
// @Router /objcontracts_add [post]
func HandleAddObjContract(w http.ResponseWriter, r *http.Request) {
	var gs ifObjContractService
	gs = services.NewObjContractService(pgsql.ObjContractStorage{})
	ctx := context.Background()

	a := models.ObjContract{}
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
		log.Println("Failed execute ifAreaService.Add: ", err)
	}

	output, err := json.Marshal(models.Json_id{Id: ai})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(output)

	return

}

// HandleUpdObjContract godoc
// @Summary Update objcontract
// @Description update objcontract
// @Tags objcontracts
// @Accept json
// @Produce  json
// @Param u body models.ObjContract true "Update objcontract. Significant params: Contract.Id, Object.Id, ObjTypeId, Startdate, Enddate"
// @Success 200 {object} models.Json_id
// @Failure 500
// @Router /objcontracts_upd [post]
func HandleUpdObjContract(w http.ResponseWriter, r *http.Request) {
	var gs ifObjContractService
	gs = services.NewObjContractService(pgsql.ObjContractStorage{})
	ctx := context.Background()

	u := models.ObjContract{}
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
		log.Println("Failed execute ifAreaService.Upd: ", err)
	}

	output, err := json.Marshal(models.Json_id{Id: ui})

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(output)

	return
}

// HandleDelObjContract godoc
// @Summary Close objcontract
// @Description close objcontract
// @Tags objcontracts
// @Accept json
// @Produce  json
// @Param d body models.IdClose true "Close objcontract. Significant params: Id, CloseDate"
// @Success 200 {object} models.Json_id
// @Failure 500
// @Router /objcontracts_del [post]
func HandleDelObjContract(w http.ResponseWriter, r *http.Request) {
	var gs ifObjContractService
	gs = services.NewObjContractService(pgsql.ObjContractStorage{})
	ctx := context.Background()

	d := models.IdClose{}
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

	i, err := gs.Del(ctx, d)
	if err != nil {
		log.Println("Failed execute ifAreaService.Del: ", err)
	}

	output, err := json.Marshal(models.Json_id{Id: i})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(output)

	return
}

// HandleGetObjContract godoc
// @Summary Get objcontract
// @Description get objcontract
// @Tags objcontracts
// @Produce  json
// @Param id path int true "ObjContract by id"
// @Param actualdate query string false "actual date"
// @Success 200 {object} models.ObjContract_count
// @Failure 500
// @Router /objcontracts/{id} [get]
func HandleGetObjContract(w http.ResponseWriter, r *http.Request) {
	var gs ifObjContractService
	gs = services.NewObjContractService(pgsql.ObjContractStorage{})
	ctx := context.Background()
	auth := utils.GetAuth(r)

	vars := mux.Vars(r)
	i, err := strconv.Atoi(vars["id"])
	if err != nil {
		i = 0
	}

	query := r.URL.Query()

	gs3 := time.Now().Format("2006-01-02")
	// log.Println(gs3)
	gs3s, ok := query["actualdate"]
	if ok && len(gs3s) > 0 {
		//case insensitive
		gs3 = strings.ToUpper(gs3s[0])
		//quotes
		re := regexp.MustCompile(`'`)
		gs3 = string(re.ReplaceAll([]byte(gs3), []byte("''")))
	}

	out_arr, err := gs.GetOne(ctx, i, gs3)
	if err != nil {
		log.Println("Failed execute ifObjContractService.GetOne: ", err)
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
