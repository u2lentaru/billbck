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

type ifObjLineService interface {
	GetList(ctx context.Context, pg, pgs int, gs1, gs2 string, ord int, dsc bool) (models.ObjLine_count, error)
	Add(ctx context.Context, ea models.ObjLine) (int, error)
	Upd(ctx context.Context, eu models.ObjLine) (int, error)
	Del(ctx context.Context, ed []int) ([]int, error)
	GetOne(ctx context.Context, i int) (models.ObjLine_count, error)
	GetObj(ctx context.Context, gs1, gs2 string) (models.ObjLine_count, error)
}

// HandleObjLines godoc
// @Summary List objlines
// @Description get objlines list
// @Tags objlines
// @Produce json
// @Param page query int false "page number"
// @Param page_size query int false "page size"
// @Param objectname query string false "objectname search pattern"
// @Param cableresistancename query string false "cableresistancename search pattern"
// @Param ordering query string false "order by {id|cableresistancename|objectname}"
// @Param desc query boolean false "descending order {true|false}"
// @Success 200 {object} models.ObjLine_count
// @Failure 500
// @Router /objlines [get]
func HandleObjLines(w http.ResponseWriter, r *http.Request) {
	var gs ifObjLineService
	gs = services.NewObjLineService(pgsql.ObjLineStorage{})
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
	gs2s, ok := query["cableresistancename"]
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
	} else if ords[0] == "cableresistancename" {
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

// HandleAddObjLine godoc
// @Summary Add objline
// @Description add objline
// @Tags objlines
// @Accept json
// @Produce  json
// @Param a body models.AddObjLine true "New objtranscurr. Significant params: ObjId, ObjTypeId, CableResistance.Id, LineLength, Startdate"
// @Success 200 {object} models.Json_id
// @Failure 500
// @Router /objlines_add [post]
func HandleAddObjLine(w http.ResponseWriter, r *http.Request) {
	var gs ifObjLineService
	gs = services.NewObjLineService(pgsql.ObjLineStorage{})
	ctx := context.Background()

	a := models.ObjLine{}
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
		log.Println("Failed execute ifKskService.Add: ", err)
	}

	output, err := json.Marshal(models.Json_id{Id: ai})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(output)

	return
}

// HandleUpdObjLine godoc
// @Summary Update objline
// @Description update objline
// @Tags objlines
// @Accept json
// @Produce  json
// @Param u body models.ObjLine true "Update objtranscurr. Significant params: Id, ObjId, ObjTypeId, CableResistance.Id, LineLength, Startdate, Enddate"
// @Success 200 {object} models.Json_id
// @Failure 500
// @Router /objlines_upd [post]
func HandleUpdObjLine(w http.ResponseWriter, r *http.Request) {
	var gs ifObjLineService
	gs = services.NewObjLineService(pgsql.ObjLineStorage{})
	ctx := context.Background()

	u := models.ObjLine{}
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
		log.Println("Failed execute ifKskService.Upd: ", err)
	}

	output, err := json.Marshal(models.Json_id{Id: ui})

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(output)

	return
}

// HandleDelObjLine godoc
// @Summary Delete objlines
// @Description delete objlines
// @Tags objlines
// @Accept json
// @Produce  json
// @Param d body models.Json_ids true "Delete objlines"
// @Success 200 {object} models.Json_ids
// @Failure 500
// @Router /objlines_del [post]
func HandleDelObjLine(w http.ResponseWriter, r *http.Request) {
	var gs ifObjLineService
	gs = services.NewObjLineService(pgsql.ObjLineStorage{})
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
		log.Println("Failed execute ifKskService.Del: ", err)
	}

	output, err := json.Marshal(models.Json_ids{Ids: res})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(output)

	return
}

// HandleGetObjLine godoc
// @Summary Get objline
// @Description get objline
// @Tags objlines
// @Produce  json
// @Param id path int true "Objline by id"
// @Success 200 {object} models.ObjLine_count
// @Failure 500
// @Router /objlines/{id} [get]
func HandleGetObjLine(w http.ResponseWriter, r *http.Request) {
	var gs ifObjLineService
	gs = services.NewObjLineService(pgsql.ObjLineStorage{})
	ctx := context.Background()
	auth := utils.GetAuth(r)

	vars := mux.Vars(r)
	i, err := strconv.Atoi(vars["id"])
	if err != nil {
		i = 0
	}

	out_arr, err := gs.GetOne(ctx, i)
	if err != nil {
		log.Println("Failed execute ifObjLineService.GetOne: ", err)
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

// HandleObjLinesByObj godoc
// @Summary Get objlines by object
// @Description get objlines by object
// @Tags objlines
// @Produce  json
// @Param objid query string false "obj&tgu id"
// @Param tid query string false "obj&tgu type id (obj - type = 0, tgu - type > 0)"
// @Success 200 {object} models.ObjLine_count
// @Failure 500
// @Router /objlines_obj [get]
func (s *APG) HandleObjLinesByObj(w http.ResponseWriter, r *http.Request) {
	var gs ifObjLineService
	gs = services.NewObjLineService(pgsql.ObjLineStorage{})
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

	w.Write(out_count)

	return
}
