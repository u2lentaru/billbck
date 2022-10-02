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

type ifObjTransVoltService interface {
	GetList(ctx context.Context, pg, pgs int, gs1, gs2 string, ord int, dsc bool) (models.ObjTransVolt_count, error)
	Add(ctx context.Context, ea models.ObjTransVolt) (int, error)
	Upd(ctx context.Context, eu models.ObjTransVolt) (int, error)
	Del(ctx context.Context, ed []int) ([]int, error)
	GetOne(ctx context.Context, i int) (models.ObjTransVolt_count, error)
	GetObj(ctx context.Context, gs1, gs2 string) (models.ObjTransVolt_count, error)
}

// HandleObjTransVolt godoc
// @Summary List objtransvolt
// @Description get objtransvolt list
// @Tags objtransvolt
// @Produce json
// @Param page query int false "page number"
// @Param page_size query int false "page size"
// @Param objectname query string false "objectname search pattern"
// @Param transvoltname query string false "transvoltname search pattern"
// @Param ordering query string false "order by {id|objectname|transvoltname}"
// @Param desc query boolean false "descending order {true|false}"
// @Success 200 {object} models.ObjTransVolt_count
// @Failure 500
// @Router /objtransvolt [get]
func HandleObjTransVolt(w http.ResponseWriter, r *http.Request) {
	var gs ifObjTransVoltService
	gs = services.NewObjTransVoltService(pgsql.ObjTransVoltStorage{})
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
	gs1s, ok := query["objectname"]
	if ok && len(gs1s) > 0 {
		//case insensitive
		gs1 = strings.ToUpper(gs1s[0])
		//quotes
		re := regexp.MustCompile(`'`)
		gs1 = string(re.ReplaceAll([]byte(gs1), []byte("''")))
	}

	gs2 := ""
	gs2s, ok := query["transvoltname"]
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
	} else if ords[0] == "objectname" {
		ord = 2
	} else if ords[0] == "transvoltname" {
		ord = 3
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

// HandleAddObjTransVolt godoc
// @Summary Add objtransvolt
// @Description add objtransvolt
// @Tags objtransvolt
// @Accept json
// @Produce  json
// @Param a body models.AddObjTransVolt true "New objtransvolt. Significant params: ObjId, ObjTypeId, TransVolt.Id, Startdate"
// @Success 200 {object} models.Json_id
// @Failure 500
// @Router /objtransvolt_add [post]
func HandleAddObjTransVolt(w http.ResponseWriter, r *http.Request) {
	var gs ifObjTransVoltService
	gs = services.NewObjTransVoltService(pgsql.ObjTransVoltStorage{})
	ctx := context.Background()

	a := models.ObjTransVolt{}
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
		log.Println("Failed execute ifObjTransVoltService.Add: ", err)
	}

	output, err := json.Marshal(models.Json_id{Id: ai})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(output)

	return
}

// HandleUpdObjTransVolt godoc
// @Summary Update objtransvolt
// @Description update objtransvolt
// @Tags objtransvolt
// @Accept json
// @Produce  json
// @Param u body models.ObjTransVolt true "Update objtransvolt. Significant params: Id, ObjId, ObjTypeId, TransVolt.Id, Startdate, Enddate"
// @Success 200 {object} models.Json_id
// @Failure 500
// @Router /objtransvolt_upd [post]
func HandleUpdObjTransVolt(w http.ResponseWriter, r *http.Request) {
	var gs ifObjTransVoltService
	gs = services.NewObjTransVoltService(pgsql.ObjTransVoltStorage{})
	ctx := context.Background()

	u := models.ObjTransVolt{}
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
		log.Println("Failed execute ifObjTransVoltService.Upd: ", err)
	}

	output, err := json.Marshal(models.Json_id{Id: ui})

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(output)

	return

}

// HandleDelObjTransVolt godoc
// @Summary Delete objtransvolts
// @Description delete objtransvolts
// @Tags objtransvolt
// @Accept json
// @Produce  json
// @Param d body models.Json_ids true "Delete objtransvolts"
// @Success 200 {object} models.Json_ids
// @Failure 500
// @Router /objtransvolt_del [post]
func HandleDelObjTransVolt(w http.ResponseWriter, r *http.Request) {
	var gs ifObjTransVoltService
	gs = services.NewObjTransVoltService(pgsql.ObjTransVoltStorage{})
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
		log.Println("Failed execute ifObjTransVoltService.Del: ", err)
	}

	output, err := json.Marshal(models.Json_ids{Ids: res})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(output)

	return
}

// HandleGetObjTransVolt godoc
// @Summary Get objtransvolt
// @Description get objtransvolt
// @Tags objtransvolt
// @Produce  json
// @Param id path int true "Objtransvolt by id"
// @Success 200 {object} models.ObjTransVolt_count
// @Failure 500
// @Router /objtransvolt/{id} [get]
func HandleGetObjTransVolt(w http.ResponseWriter, r *http.Request) {
	var gs ifObjTransVoltService
	gs = services.NewObjTransVoltService(pgsql.ObjTransVoltStorage{})
	ctx := context.Background()
	auth := utils.GetAuth(r)

	vars := mux.Vars(r)
	i, err := strconv.Atoi(vars["id"])
	if err != nil {
		i = 0
	}

	out_arr, err := gs.GetOne(ctx, i)
	if err != nil {
		log.Println("Failed execute ifObjTransVoltService.GetOne: ", err)
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

// HandleObjTransVoltByObj godoc
// @Summary Get objtransvolt by object
// @Description get objtransvolt by object
// @Tags objtransvolt
// @Produce  json
// @Param objid query string false "obj&tgu id"
// @Param tid query string false "obj&tgu type id (obj - type = 0, tgu - type > 0)"
// @Success 200 {object} models.ObjTransVolt_count
// @Failure 500
// @Router /objtransvolt_obj [get]
func HandleObjTransVoltByObj(w http.ResponseWriter, r *http.Request) {
	var gs ifObjTransVoltService
	gs = services.NewObjTransVoltService(pgsql.ObjTransVoltStorage{})
	ctx := context.Background()
	auth := utils.GetAuth(r)

	query := r.URL.Query()

	gs1 := "0"
	gs1s, ok := query["objid"]
	if ok && len(gs1s) > 0 {
		_, err := strconv.Atoi(gs1s[0])
		if err == nil {
			gs1 = gs1s[0]
		}
	}

	gs2 := "0"
	gs2s, ok := query["tid"]
	if ok && len(gs2s) > 0 {
		_, err := strconv.Atoi(gs2s[0])
		if err == nil {
			gs2 = gs2s[0]
		}
	}

	out_arr, err := gs.GetObj(ctx, gs1, gs2)
	if err != nil {
		log.Println("Failed execute ifObjTransVoltService.GetObj: ", err)
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
