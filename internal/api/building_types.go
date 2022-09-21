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

type ifBuildingTypeService interface {
	GetList(ctx context.Context, pg, pgs int, gs1 string, ord int, dsc bool) (models.BuildingType_count, error)
	Add(ctx context.Context, ea models.BuildingType) (int, error)
	Upd(ctx context.Context, eu models.BuildingType) (int, error)
	Del(ctx context.Context, ed []int) ([]int, error)
	GetOne(ctx context.Context, i int) (models.BuildingType_count, error)
}

// HandleBuildingTypes godoc
// @Summary List building types
// @Description get building types
// @Tags building types
// @Produce  json
// @Param page query int false "page number"
// @Param page_size query int false "page size"
// @Param buildingtypename query string false "buildingtypename search pattern"
// @Param ordering query string false "order by {id|buildingtypename}"
// @Param desc query boolean false "descending order {true|false}"
// @Success 200 {object} models.BuildingType_count
// @Failure 500
// @Router /building_types [get]
func HandleBuildingTypes(w http.ResponseWriter, r *http.Request) {
	var gs ifBuildingTypeService
	gs = services.NewBuildingTypeService(pgsql.BuildingTypeStorage{})
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

	gsn := ""
	gsns, ok := query["buildingtypename"]
	if ok && len(gsns) > 0 {
		//case insensitive
		gsn = strings.ToUpper(gsns[0])
		//quotes
		re := regexp.MustCompile(`'`)
		gsn = string(re.ReplaceAll([]byte(gsn), []byte("''")))
	}

	ord := 1
	ords, ok := query["ordering"]
	if !ok || len(ords) == 0 {
		ord = 1
	} else if ords[0] == "id" {
		ord = 1
	} else if ords[0] == "buildingtypename" {
		ord = 2
	}

	dsc := false
	dscs, ok := query["desc"]
	if ok && len(dscs) > 0 {
		if !(dscs[0] == "0") {
			dsc = true
		}
	}

	out_arr, err := gs.GetList(ctx, pg, pgs, gsn, ord, dsc)
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

// HandleAddBuildingType godoc
// @Summary Add building type
// @Description add building type
// @Tags building types
// @Accept json
// @Produce  json
// @Param a body models.AddBuildingType true "New building type"
// @Success 200 {object} models.Json_id
// @Failure 500
// @Router /building_types_add [post]
func HandleAddBuildingType(w http.ResponseWriter, r *http.Request) {
	var gs ifBuildingTypeService
	gs = services.NewBuildingTypeService(pgsql.BuildingTypeStorage{})
	ctx := context.Background()

	a := models.BuildingType{}
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
		log.Println("Failed execute ifBuildingTypeService.Add: ", err)
	}

	output, err := json.Marshal(models.Json_id{Id: ai})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(output)

	return

}

// HandleUpdBuildingType godoc
// @Summary Update building type
// @Description update building type
// @Tags building types
// @Accept json
// @Produce  json
// @Param u body models.BuildingType true "Update building type"
// @Success 200 {object} models.Json_id
// @Failure 500
// @Router /building_types_upd [post]
func HandleUpdBuildingType(w http.ResponseWriter, r *http.Request) {
	var gs ifBuildingTypeService
	gs = services.NewBuildingTypeService(pgsql.BuildingTypeStorage{})
	ctx := context.Background()

	u := models.BuildingType{}
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
		log.Println("Failed execute ifBuildingTypeService.Upd: ", err)
	}

	output, err := json.Marshal(models.Json_id{Id: ui})

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(output)

	return

}

// HandleDelBuildingType godoc
// @Summary Delete building types
// @Description delete building types
// @Tags building types
// @Accept json
// @Produce  json
// @Param d body models.Json_ids true "Delete building types"
// @Success 200 {object} models.Json_ids
// @Failure 500
// @Router /building_types_del [post]
func HandleDelBuildingType(w http.ResponseWriter, r *http.Request) {
	var gs ifBuildingTypeService
	gs = services.NewBuildingTypeService(pgsql.BuildingTypeStorage{})
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
		log.Println("Failed execute ifBuildingTypeService.Del: ", err)
	}

	output, err := json.Marshal(models.Json_ids{Ids: res})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(output)

	return
}

// HandleGetBuildingType godoc
// @Summary Get building type
// @Description get building type
// @Tags building types
// @Produce  json
// @Param id path int true "Building type by id"
// @Success 200 {object} models.BuildingType_count
// @Failure 500
// @Router /building_types/{id} [get]
func HandleGetBuildingType(w http.ResponseWriter, r *http.Request) {
	var gs ifBuildingTypeService
	gs = services.NewBuildingTypeService(pgsql.BuildingTypeStorage{})
	ctx := context.Background()

	vars := mux.Vars(r)
	i, err := strconv.Atoi(vars["id"])
	if err != nil {
		i = 0
	}

	out_arr, err := gs.GetOne(ctx, i)
	if err != nil {
		log.Println("Failed execute ifBuildingTypeService.GetOne: ", err)
	}

	out_count, err := json.Marshal(out_arr)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(out_count)

	return
}
