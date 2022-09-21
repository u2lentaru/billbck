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

type ifCableResistanceService interface {
	GetList(ctx context.Context, pg, pgs int, gs1 string, ord int, dsc bool) (models.CableResistance_count, error)
	Add(ctx context.Context, ea models.CableResistance) (int, error)
	Upd(ctx context.Context, eu models.CableResistance) (int, error)
	Del(ctx context.Context, ed []int) ([]int, error)
	GetOne(ctx context.Context, i int) (models.CableResistance_count, error)
}

// HandleCableResistances godoc
// @Summary List cableresistances
// @Description get cableresistance list
// @Tags cableresistances
// @Produce  json
// @Param page query int false "page number"
// @Param page_size query int false "page size"
// @Param cableresistancename query string false "cableresistancename search pattern"
// @Param ordering query string false "order by {id|cableresistancename}"
// @Param desc query boolean false "descending order {true|false}"
// @Success 200 {object} models.CableResistance_count
// @Failure 500
// @Router /cableresistances [get]
func HandleCableResistances(w http.ResponseWriter, r *http.Request) {
	var gs ifCableResistanceService
	gs = services.NewCableResistanceService(pgsql.CableResistanceStorage{})
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
	gs1s, ok := query["cableresistancename"]
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
	} else if ords[0] == "cableresistancename" {
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

// HandleAddCableResistance godoc
// @Summary Add cableresistance
// @Description add cableresistance
// @Tags cableresistances
// @Accept json
// @Produce  json
// @Param a body models.AddCableResistance true "New cableresistance"
// @Success 200 {object} models.Json_id
// @Failure 500
// @Router /cableresistances_add [post]
func HandleAddCableResistance(w http.ResponseWriter, r *http.Request) {
	var gs ifCableResistanceService
	gs = services.NewCableResistanceService(pgsql.CableResistanceStorage{})
	ctx := context.Background()

	a := models.CableResistance{}
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
		log.Println("Failed execute ifCableResistanceService.Add: ", err)
	}

	output, err := json.Marshal(models.Json_id{Id: ai})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(output)

	return

}

// HandleUpdCableResistance godoc
// @Summary Update cableresistance
// @Description update cableresistance
// @Tags cableresistances
// @Accept json
// @Produce  json
// @Param u body models.CableResistance true "Update cableresistance"
// @Success 200 {object} models.Json_id
// @Failure 500
// @Router /cableresistances_upd [post]
func HandleUpdCableResistance(w http.ResponseWriter, r *http.Request) {
	var gs ifCableResistanceService
	gs = services.NewCableResistanceService(pgsql.CableResistanceStorage{})
	ctx := context.Background()

	u := models.CableResistance{}
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
		log.Println("Failed execute ifCableResistanceService.Upd: ", err)
	}

	output, err := json.Marshal(models.Json_id{Id: ui})

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(output)

	return

}

// HandleDelCableResistance godoc
// @Summary Delete cableresistances
// @Description delete cableresistances
// @Tags cableresistances
// @Accept json
// @Produce  json
// @Param d body models.Json_ids true "Delete cableresistances"
// @Success 200 {object} models.Json_ids
// @Failure 500
// @Router /cableresistances_del [post]
func HandleDelCableResistance(w http.ResponseWriter, r *http.Request) {
	var gs ifCableResistanceService
	gs = services.NewCableResistanceService(pgsql.CableResistanceStorage{})
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
		log.Println("Failed execute ifCableResistanceService.Del: ", err)
	}

	output, err := json.Marshal(models.Json_ids{Ids: res})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(output)

	return

}

// HandleGetCableResistance godoc
// @Summary Get cableresistance
// @Description get cableresistance
// @Tags cableresistances
// @Produce  json
// @Param id path int true "CableResistance by id"
// @Success 200 {object} models.CableResistance_count
// @Failure 500
// @Router /cableresistances/{id} [get]
func HandleGetCableResistance(w http.ResponseWriter, r *http.Request) {
	var gs ifCableResistanceService
	gs = services.NewCableResistanceService(pgsql.CableResistanceStorage{})
	ctx := context.Background()

	vars := mux.Vars(r)
	i, err := strconv.Atoi(vars["id"])
	if err != nil {
		i = 0
	}

	out_arr, err := gs.GetOne(ctx, i)
	if err != nil {
		log.Println("Failed execute ifCableResistanceService.GetOne: ", err)
	}

	out_count, err := json.Marshal(out_arr)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(out_count)

	return
}
