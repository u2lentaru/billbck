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

type ifObjTransCurrService interface {
	GetList(ctx context.Context, pg, pgs int, gs1, gs2 string, ord int, dsc bool) (models.ObjTransCurr_count, error)
	Add(ctx context.Context, ea models.ObjTransCurr) (int, error)
	Upd(ctx context.Context, eu models.ObjTransCurr) (int, error)
	Del(ctx context.Context, ed []int) ([]int, error)
	GetOne(ctx context.Context, i int) (models.ObjTransCurr_count, error)
	GetObj(ctx context.Context, gs1, gs2 string) (models.ObjTransCurr_count, error)
}

// HandleObjTransCurr godoc
// @Summary List objtranscurr
// @Description get objtranscurr list
// @Tags objtranscurr
// @Produce json
// @Param page query int false "page number"
// @Param page_size query int false "page size"
// @Param objectname query string false "objectname search pattern"
// @Param transcurrname query string false "transcurrname search pattern"
// @Param ordering query string false "order by {id|objectname|transcurrname}"
// @Param desc query boolean false "descending order {true|false}"
// @Success 200 {object} models.ObjTransCurr_count
// @Failure 500
// @Router /objtranscurr [get]
func HandleObjTransCurr(w http.ResponseWriter, r *http.Request) {
	var gs ifObjTransCurrService
	gs = services.NewObjTransCurrService(pgsql.ObjTransCurrStorage{})
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
	gs2s, ok := query["transcurrname"]
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
	} else if ords[0] == "transcurrname" {
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

// HandleAddObjTransCurr godoc
// @Summary Add objtranscurr
// @Description add objtranscurr
// @Tags objtranscurr
// @Accept json
// @Produce  json
// @Param a body models.AddObjTransCurr true "New objtranscurr. Significant params: ObjId, ObjTypeId, TransCurr.Id, Startdate"
// @Success 200 {object} models.Json_id
// @Failure 500
// @Router /objtranscurr_add [post]
func HandleAddObjTransCurr(w http.ResponseWriter, r *http.Request) {
	var gs ifObjTransCurrService
	gs = services.NewObjTransCurrService(pgsql.ObjTransCurrStorage{})
	ctx := context.Background()

	a := models.ObjTransCurr{}
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
		log.Println("Failed execute ifObjTransCurrService.Add: ", err)
	}

	output, err := json.Marshal(models.Json_id{Id: ai})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(output)

	return
}

// HandleUpdObjTransCurr godoc
// @Summary Update objtranscurr
// @Description update objtranscurr
// @Tags objtranscurr
// @Accept json
// @Produce  json
// @Param u body models.ObjTransCurr true "Update objtranscurr. Significant params: Id, ObjId, ObjTypeId, TransCurr.Id, Startdate, Enddate"
// @Success 200 {object} models.Json_id
// @Failure 500
// @Router /objtranscurr_upd [post]
func HandleUpdObjTransCurr(w http.ResponseWriter, r *http.Request) {
	var gs ifObjTransCurrService
	gs = services.NewObjTransCurrService(pgsql.ObjTransCurrStorage{})
	ctx := context.Background()

	u := models.ObjTransCurr{}
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
		log.Println("Failed execute ifObjTransCurrService.Upd: ", err)
	}

	output, err := json.Marshal(models.Json_id{Id: ui})

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(output)

	return

}

// HandleDelObjTransCurr godoc
// @Summary Delete objtranscurrs
// @Description delete objtranscurrs
// @Tags objtranscurr
// @Accept json
// @Produce  json
// @Param d body models.Json_ids true "Delete objtranscurrs"
// @Success 200 {object} models.Json_ids
// @Failure 500
// @Router /objtranscurr_del [post]
func HandleDelObjTransCurr(w http.ResponseWriter, r *http.Request) {
	var gs ifObjTransCurrService
	gs = services.NewObjTransCurrService(pgsql.ObjTransCurrStorage{})
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
		log.Println("Failed execute ifObjTransCurrService.Del: ", err)
	}

	output, err := json.Marshal(models.Json_ids{Ids: res})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(output)

	return
}

// HandleGetObjTransCurr godoc
// @Summary Get objtranscurr
// @Description get objtranscurr
// @Tags objtranscurr
// @Produce  json
// @Param id path int true "Objtranscurr by id"
// @Success 200 {object} models.ObjTransCurr_count
// @Failure 500
// @Router /objtranscurr/{id} [get]
func HandleGetObjTransCurr(w http.ResponseWriter, r *http.Request) {
	var gs ifObjTransCurrService
	gs = services.NewObjTransCurrService(pgsql.ObjTransCurrStorage{})
	ctx := context.Background()
	auth := utils.GetAuth(r)

	vars := mux.Vars(r)
	i, err := strconv.Atoi(vars["id"])
	if err != nil {
		i = 0
	}

	out_arr, err := gs.GetOne(ctx, i)
	if err != nil {
		log.Println("Failed execute ifObjTransCurrService.GetOne: ", err)
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

// HandleObjTransCurrByObj godoc
// @Summary Get objtranscurr by object
// @Description get objtranscurr by object
// @Tags objtranscurr
// @Produce  json
// @Param objid query string false "obj&tgu id"
// @Param tid query string false "obj&tgu type id (obj - type = 0, tgu - type > 0)"
// @Success 200 {object} models.ObjTransCurr_count
// @Failure 500
// @Router /objtranscurr_obj [get]
func HandleObjTransCurrByObj(w http.ResponseWriter, r *http.Request) {
	var gs ifObjTransCurrService
	gs = services.NewObjTransCurrService(pgsql.ObjTransCurrStorage{})
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
