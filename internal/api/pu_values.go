package api

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/u2lentaru/billbck/internal/adapters/db/pgsql"
	"github.com/u2lentaru/billbck/internal/models"
	"github.com/u2lentaru/billbck/internal/services"
	"github.com/u2lentaru/billbck/internal/utils"
)

type ifPuValueService interface {
	GetList(ctx context.Context, pg, pgs, gs1, ord int, dsc bool) (models.PuValue_count, error)
	Add(ctx context.Context, ea models.PuValue) (int, error)
	Upd(ctx context.Context, eu models.PuValue) (int, error)
	Del(ctx context.Context, ed []int) ([]int, error)
	GetOne(ctx context.Context, i int) (models.PuValue_count, error)
	AskuePrev(ctx context.Context, af models.AskueFile) (models.PuValueAskue_count, error)
	AskueLoad(ctx context.Context, af models.AskueFile) (models.AskueLoadRes, error)
}

// HandlePuValues godoc
// @Summary List pu values
// @Description get pu values list
// @Tags puvalues
// @Produce  json
// @Param page query int false "page number"
// @Param page_size query int false "page size"
// @Param puid query int false "puid search pattern"
// @Param ordering query string false "order by {puid|valuedate}"
// @Param desc query boolean false "descending order {true|false}"
// @Success 200 {object} models.PuValue_count
// @Failure 500
// @Router /puvalues [get]
func HandlePuValues(w http.ResponseWriter, r *http.Request) {
	var gs ifPuValueService
	gs = services.NewPuValueService(pgsql.PuValueStorage{})
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
	gs1s, ok := query["puid"]
	if ok && len(gs1s) > 0 {
		t, err := strconv.Atoi(gs1s[0])
		if err == nil {
			gs1 = t
		}
	}

	ord := 1
	ords, ok := query["ordering"]
	if !ok || len(ords) == 0 {
		ord = 1
	} else if ords[0] == "puid" {
		ord = 2
	} else if ords[0] == "valuedate" {
		ord = 3
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

	out_arr.Auth = auth
	out_count, err := json.Marshal(out_arr)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Write(out_count)

	return

}

// HandleAddPuValue godoc
// @Summary Add pu value
// @Description add pu value
// @Tags puvalues
// @Accept json
// @Produce  json
// @Param a body models.PuValue true "New pu value. Significant params: PuId, ValueDate, PuValue"
// @Success 200 {object} models.Json_id
// @Failure 500
// @Router /puvalues_add [post]
func HandleAddPuValue(w http.ResponseWriter, r *http.Request) {
	var gs ifPuValueService
	gs = services.NewPuValueService(pgsql.PuValueStorage{})
	ctx := context.Background()

	a := models.PuValue{}
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
		log.Println("Failed execute ifPuValueService.Add: ", err)
	}

	output, err := json.Marshal(models.Json_id{Id: ai})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(output)

	return

}

// HandleUpdPuValue godoc
// @Summary Update pu value
// @Description update pu value
// @Tags puvalues
// @Accept json
// @Produce  json
// @Param u body models.PuValue true "Update pu value. Significant params: Id, ValueDate, PuValue"
// @Success 200 {object} models.Json_id
// @Failure 500
// @Router /puvalues_upd [post]
func HandleUpdPuValue(w http.ResponseWriter, r *http.Request) {
	var gs ifPuValueService
	gs = services.NewPuValueService(pgsql.PuValueStorage{})
	ctx := context.Background()

	u := models.PuValue{}
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
		log.Println("Failed execute ifPuValueService.Upd: ", err)
	}

	output, err := json.Marshal(models.Json_id{Id: ui})

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(output)

	return
}

// HandleDelPuValue godoc
// @Summary Delete pu values
// @Description delete pu values
// @Tags puvalues
// @Accept json
// @Produce  json
// @Param d body models.Json_ids true "Delete pu values"
// @Success 200 {object} models.Json_ids
// @Failure 500
// @Router /puvalues_del [post]
func HandleDelPuValue(w http.ResponseWriter, r *http.Request) {
	var gs ifPuValueService
	gs = services.NewPuValueService(pgsql.PuValueStorage{})
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
		log.Println("Failed execute ifPuValueService.Del: ", err)
	}

	output, err := json.Marshal(models.Json_ids{Ids: res})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(output)

	return
}

// HandleGetPuValue godoc
// @Summary Get pu value
// @Description get pu value
// @Tags puvalues
// @Produce  json
// @Param id path int true "Pu value Id"
// @Success 200 {object} models.PuValue_count
// @Failure 500
// @Router /puvalues/{id} [get]
func HandleGetPuValue(w http.ResponseWriter, r *http.Request) {
	var gs ifPuValueService
	gs = services.NewPuValueService(pgsql.PuValueStorage{})
	ctx := context.Background()
	auth := utils.GetAuth(r)

	vars := mux.Vars(r)
	i, err := strconv.Atoi(vars["id"])
	if err != nil {
		i = 0
	}

	out_arr, err := gs.GetOne(ctx, i)
	if err != nil {
		log.Println("Failed execute ifPuValueService.GetOne: ", err)
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

// HandlePuValuesAskuePreview godoc
// @Summary Preview askue pu values
// @Description preview askue pu values
// @Tags puvalues
// @Accept json
// @Produce  json
// @Param af body models.AskueFile true "Askue file to preview"
// @Success 200 {object} models.PuValueAskue_count
// @Failure 500
// @Router /puvalues_askue_prev [post]
func HandlePuValuesAskuePreview(w http.ResponseWriter, r *http.Request) {
	var gs ifPuValueService
	gs = services.NewPuValueService(pgsql.PuValueStorage{})
	ctx := context.Background()
	auth := utils.GetAuth(r)

	af := models.AskueFile{}
	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	err = json.Unmarshal(body, &af)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	out_arr, err := gs.AskuePrev(ctx, af)
	if err != nil {
		log.Println("Failed execute ifPuValueService.AskuePrev: ", err)
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

// HandlePuValuesAskue godoc
// @Summary Load askue pu values
// @Description load askue pu values
// @Tags puvalues
// @Accept json
// @Produce  json
// @Param af body models.AskueFile true "Askue file to load"
// @Success 200 {object} models.AskueLoadRes
// @Failure 500
// @Router /puvalues_askue [post]
func HandlePuValuesAskue(w http.ResponseWriter, r *http.Request) {
	var gs ifPuValueService
	gs = services.NewPuValueService(pgsql.PuValueStorage{})
	ctx := context.Background()

	af := models.AskueFile{}
	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	err = json.Unmarshal(body, &af)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	out_arr, err := gs.AskueLoad(ctx, af)
	if err != nil {
		log.Println("Failed execute ifPuValueService.AskueLoad: ", err)
	}

	out_count, err := json.Marshal(out_arr)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(out_count)

	return
}
