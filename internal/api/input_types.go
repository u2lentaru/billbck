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

type ifInputTypeService interface {
	GetList(ctx context.Context, pg, pgs int, gs1 string, ord int, dsc bool) (models.InputType_count, error)
	Add(ctx context.Context, ea models.InputType) (int, error)
	Upd(ctx context.Context, eu models.InputType) (int, error)
	Del(ctx context.Context, ed []int) ([]int, error)
	GetOne(ctx context.Context, i int) (models.InputType_count, error)
}

// HandleInputTypes godoc
// @Summary List input types
// @Description get input type list
// @Tags input types
// @Produce  json
// @Param page query int false "page number"
// @Param page_size query int false "page size"
// @Param inputtypename query string false "inputtypename search pattern"
// @Param ordering query string false "order by {id|inputtypename}"
// @Param desc query boolean false "descending order {true|false}"
// @Success 200 {object} models.InputType_count
// @Failure 500
// @Router /input_types [get]
func HandleInputTypes(w http.ResponseWriter, r *http.Request) {
	var gs ifInputTypeService
	gs = services.NewInputTypeService(pgsql.InputTypeStorage{})
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
	gs1s, ok := query["inputtypename"]
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
	} else if ords[0] == "inputtypename" {
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

	out_arr.Auth = auth
	out_count, err := json.Marshal(out_arr)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Write(out_count)

	return
}

// HandleAddInputType godoc
// @Summary Add input type
// @Description add input type
// @Tags input types
// @Accept json
// @Produce  json
// @Param a body models.AddInputType true "New input type"
// @Success 200 {object} models.Json_id
// @Failure 500
// @Router /input_types_add [post]
func HandleAddInputType(w http.ResponseWriter, r *http.Request) {
	var gs ifInputTypeService
	gs = services.NewInputTypeService(pgsql.InputTypeStorage{})
	ctx := context.Background()

	a := models.InputType{}
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
		log.Println("Failed execute ifInputTypeService.Add: ", err)
	}

	output, err := json.Marshal(models.Json_id{Id: ai})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(output)

	return

}

// HandleUpdInputType godoc
// @Summary Update input type
// @Description update input type
// @Tags input types
// @Accept json
// @Produce  json
// @Param u body models.InputType true "Update input type"
// @Success 200 {object} models.Json_id
// @Failure 500
// @Router /input_types_upd [post]
func HandleUpdInputType(w http.ResponseWriter, r *http.Request) {
	var gs ifInputTypeService
	gs = services.NewInputTypeService(pgsql.InputTypeStorage{})
	ctx := context.Background()

	u := models.InputType{}
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
		log.Println("Failed execute ifInputTypeService.Upd: ", err)
	}

	output, err := json.Marshal(models.Json_id{Id: ui})

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(output)

	return
}

// HandleDelInputType godoc
// @Summary Delete input types
// @Description delete input types
// @Tags input types
// @Accept json
// @Produce  json
// @Param d body models.Json_ids true "Delete input types"
// @Success 200 {object} models.Json_ids
// @Failure 500
// @Router /input_types_del [post]
func HandleDelInputType(w http.ResponseWriter, r *http.Request) {
	var gs ifInputTypeService
	gs = services.NewInputTypeService(pgsql.InputTypeStorage{})
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
		log.Println("Failed execute ifInputTypeService.Del: ", err)
	}

	output, err := json.Marshal(models.Json_ids{Ids: res})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(output)

	return
}

// HandleGetInputType godoc
// @Summary Get input type
// @Description get input type
// @Tags input types
// @Produce  json
// @Param id path int true "Input type by id"
// @Success 200 {object} models.InputType_count
// @Failure 500
// @Router /input_types/{id} [get]
func HandleGetInputType(w http.ResponseWriter, r *http.Request) {
	var gs ifInputTypeService
	gs = services.NewInputTypeService(pgsql.InputTypeStorage{})
	ctx := context.Background()
	auth := utils.GetAuth(r)

	vars := mux.Vars(r)
	i, err := strconv.Atoi(vars["id"])
	if err != nil {
		i = 0
	}

	out_arr, err := gs.GetOne(ctx, i)
	if err != nil {
		log.Println("Failed execute ifInputTypeService.GetOne: ", err)
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
