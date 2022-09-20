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
)

type ifAskueTypeService interface {
	GetList(ctx context.Context, pg, pgs int, gs1 string, ord int, dsc bool) (models.AskueType_count, error)
	Add(ctx context.Context, ea models.AskueType) (int, error)
	Upd(ctx context.Context, eu models.AskueType) (int, error)
	Del(ctx context.Context, ed []int) ([]int, error)
	GetOne(ctx context.Context, i int) (models.AskueType_count, error)
}

// HandleAskueTypes godoc
// @Summary List askuetypes
// @Description get askuetypes list
// @Tags askuetypes
// @Produce  json
// @Param page query int false "page number"
// @Param page_size query int false "page size"
// @Param askuetypename query string false "askuetypename search pattern"
// @Param ordering query string false "order by {id|askuetypename}"
// @Param desc query boolean false "descending order {true|false}"
// @Success 200 {object} models.CableResistance_count
// @Failure 500
// @Router /askuetypes [get]
func HandleAskueTypes(w http.ResponseWriter, r *http.Request) {
	var gs ifAskueTypeService
	gs = services.NewAskueTypeService(pgsql.AskueTypeStorage{})
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
	gs1s, ok := query["askuetypename"]
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
	} else if ords[0] == "askuetypename" {
		ord = 2
	}

	dsc := false
	dscs, ok := query["desc"]
	if ok && len(dscs) > 0 {
		if !(dscs[0] == "0") {
			dsc = true
		}
	}

	out_arr, err := gs.GetList(ctx, pg, pgs, gs1, ord, dsc)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	out_count, err := json.Marshal(out_arr)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Write(out_count)

	return
}

// HandleAddAskueType godoc
// @Summary Add askuetype
// @Description add askuetype
// @Tags askuetypes
// @Accept json
// @Produce  json
// @Param a body models.AddAskueType true "New askuetype. Significant params: AskueTypeName, StartLine, PuColumn, ValueColumn, DateColumn, DateColumnArray(n)"
// @Success 200 {object} models.Json_id
// @Failure 500
// @Router /askuetypes_add [post]
func HandleAddAskueType(w http.ResponseWriter, r *http.Request) {
	var gs ifAskueTypeService
	gs = services.NewAskueTypeService(pgsql.AskueTypeStorage{})
	ctx := context.Background()

	a := models.AskueType{}
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
		log.Println("Failed execute ifAskueTypeService.Add: ", err)
	}

	output, err := json.Marshal(models.Json_id{Id: ai})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(output)

	return

}

// HandleUpdAskueType godoc
// @Summary Update askuetype
// @Description update askuetype
// @Tags askuetypes
// @Accept json
// @Produce  json
// @Param u body models.AskueType true "Update askuetype. Significant params: Id, AskueTypeName, StartLine, PuColumn, ValueColumn, DateColumn, DateColumnArray(n)"
// @Success 200 {object} models.Json_id
// @Failure 500
// @Router /askuetypes_upd [post]
func HandleUpdAskueType(w http.ResponseWriter, r *http.Request) {
	var gs ifAskueTypeService
	gs = services.NewAskueTypeService(pgsql.AskueTypeStorage{})
	ctx := context.Background()

	u := models.AskueType{}
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
		log.Println("Failed execute ifAskueTypeService.Upd: ", err)
	}

	output, err := json.Marshal(models.Json_id{Id: ui})

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(output)

	return

}

// HandleDelAskueType godoc
// @Summary Delete askuetypes
// @Description delete askuetypes
// @Tags askuetypes
// @Accept json
// @Produce  json
// @Param d body models.Json_ids true "Delete askuetypes"
// @Success 200 {object} models.Json_ids
// @Failure 500
// @Router /askuetypes_del [post]
func HandleDelAskueType(w http.ResponseWriter, r *http.Request) {
	var gs ifAskueTypeService
	gs = services.NewAskueTypeService(pgsql.AskueTypeStorage{})
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
		log.Println("Failed execute ifAskueTypeService.Del: ", err)
	}

	output, err := json.Marshal(models.Json_ids{Ids: res})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(output)

	return

}

// HandleGetAskueType godoc
// @Summary Get askuetype
// @Description get askuetype
// @Tags askuetypes
// @Produce  json
// @Param id path int true "AskueType by id"
// @Success 200 {object} models.AskueType_count
// @Failure 500
// @Router /askuetypes/{id} [get]
func HandleGetAskueType(w http.ResponseWriter, r *http.Request) {
	var gs ifAskueTypeService
	gs = services.NewAskueTypeService(pgsql.AskueTypeStorage{})
	ctx := context.Background()

	vars := mux.Vars(r)
	i, err := strconv.Atoi(vars["id"])
	if err != nil {
		i = 0
	}

	out_arr, err := gs.GetOne(ctx, i)
	if err != nil {
		log.Println("Failed execute ifAskueTypeService.GetOne: ", err)
	}

	out_count, err := json.Marshal(out_arr)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(out_count)

	return
}
