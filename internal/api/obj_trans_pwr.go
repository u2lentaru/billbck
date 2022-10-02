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

type ifObjTransPwrService interface {
	GetList(ctx context.Context, pg, pgs int, gs1, gs2 string, ord int, dsc bool) (models.ObjTransPwr_count, error)
	Add(ctx context.Context, ea models.ObjTransPwr) (int, error)
	Upd(ctx context.Context, eu models.ObjTransPwr) (int, error)
	Del(ctx context.Context, ed []int) ([]int, error)
	GetOne(ctx context.Context, i int) (models.ObjTransPwr_count, error)
	GetObj(ctx context.Context, gs1, gs2 string) (models.ObjTransPwr_count, error)
}

// HandleObjTransPwr godoc
// @Summary List objtranspwr
// @Description get objtranspwr list
// @Tags objtranspwr
// @Produce json
// @Param page query int false "page number"
// @Param page_size query int false "page size"
// @Param objectname query string false "objectname search pattern"
// @Param transpwrname query string false "transpwrname search pattern"
// @Param ordering query string false "order by {id|objectname|transpwrname}"
// @Param desc query boolean false "descending order {true|false}"
// @Success 200 {object} models.ObjTransPwr_count
// @Failure 500
// @Router /objtranspwr [get]
func HandleObjTransPwr(w http.ResponseWriter, r *http.Request) {
	var gs ifObjTransPwrService
	gs = services.NewObjTransPwrService(pgsql.ObjTransPwrStorage{})
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
	gs2s, ok := query["transpwrname"]
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
	} else if ords[0] == "transpwrname" {
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

// HandleAddObjTransPwr godoc
// @Summary Add objtranspwr
// @Description add objtranspwr
// @Tags objtranspwr
// @Accept json
// @Produce  json
// @Param a body models.AddObjTransPwr true "New objtranspwr. Significant params: ObjId, ObjTypeId, TransPwr.Id, Startdate"
// @Success 200 {object} models.Json_id
// @Failure 500
// @Router /objtranspwr_add [post]
func HandleAddObjTransPwr(w http.ResponseWriter, r *http.Request) {
	var gs ifObjTransPwrService
	gs = services.NewObjTransPwrService(pgsql.ObjTransPwrStorage{})
	ctx := context.Background()

	a := models.ObjTransPwr{}
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
		log.Println("Failed execute ifObjTransPwrService.Add: ", err)
	}

	output, err := json.Marshal(models.Json_id{Id: ai})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(output)

	return
}

// HandleUpdObjTransPwr godoc
// @Summary Update objtranspwr
// @Description update objtranspwr
// @Tags objtranspwr
// @Accept json
// @Produce  json
// @Param u body models.ObjTransPwr true "Update objtranspwr. Significant params: Id, ObjId, ObjTypeId, TransPwr.Id, Startdate, Enddate"
// @Success 200 {object} models.Json_id
// @Failure 500
// @Router /objtranspwr_upd [post]
func HandleUpdObjTransPwr(w http.ResponseWriter, r *http.Request) {
	var gs ifObjTransPwrService
	gs = services.NewObjTransPwrService(pgsql.ObjTransPwrStorage{})
	ctx := context.Background()

	u := models.ObjTransPwr{}
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
		log.Println("Failed execute ifObjTransPwrService.Upd: ", err)
	}

	output, err := json.Marshal(models.Json_id{Id: ui})

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(output)

	return

}

// HandleDelObjTransPwr godoc
// @Summary Delete objtranspwrs
// @Description delete objtranspwrs
// @Tags objtranspwr
// @Accept json
// @Produce  json
// @Param d body models.Json_ids true "Delete objtranspwrs"
// @Success 200 {object} models.Json_ids
// @Failure 500
// @Router /objtranspwr_del [post]
func HandleDelObjTransPwr(w http.ResponseWriter, r *http.Request) {
	var gs ifObjTransPwrService
	gs = services.NewObjTransPwrService(pgsql.ObjTransPwrStorage{})
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
		log.Println("Failed execute ifObjTransPwrService.Del: ", err)
	}

	output, err := json.Marshal(models.Json_ids{Ids: res})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(output)

	return
}

// HandleGetObjTransPwr godoc
// @Summary Get objtranspwr
// @Description get objtranspwr
// @Tags objtranspwr
// @Produce  json
// @Param id path int true "Objtranspwr by id"
// @Success 200 {object} models.ObjTransPwr_count
// @Failure 500
// @Router /objtranspwr/{id} [get]
func HandleGetObjTransPwr(w http.ResponseWriter, r *http.Request) {
	var gs ifObjTransPwrService
	gs = services.NewObjTransPwrService(pgsql.ObjTransPwrStorage{})
	ctx := context.Background()
	auth := utils.GetAuth(r)

	vars := mux.Vars(r)
	i, err := strconv.Atoi(vars["id"])
	if err != nil {
		i = 0
	}

	out_arr, err := gs.GetOne(ctx, i)
	if err != nil {
		log.Println("Failed execute ifObjTransPwrService.GetOne: ", err)
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

// HandleObjTransPwrByObj godoc
// @Summary Get objtranspwr by object
// @Description get objtranspwr by object
// @Tags objtranspwr
// @Produce  json
// @Param objid query string false "obj&tgu id"
// @Param tid query string false "obj&tgu type id (obj - type = 0, tgu - type > 0)"
// @Success 200 {object} models.ObjTransPwr_count
// @Failure 500
// @Router /objtranspwr_obj [get]
func HandleObjTransPwrByObj(w http.ResponseWriter, r *http.Request) {
	var gs ifObjTransPwrService
	gs = services.NewObjTransPwrService(pgsql.ObjTransPwrStorage{})
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
