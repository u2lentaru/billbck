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

type ifSubTypeService interface {
	GetList(ctx context.Context, pg, pgs int, gs1, gs2 string, ord int, dsc bool) (models.SubType_count, error)
	Add(ctx context.Context, ea models.SubType) (int, error)
	Upd(ctx context.Context, eu models.SubType) (int, error)
	Del(ctx context.Context, ed []int) ([]int, error)
	GetOne(ctx context.Context, i int) (models.SubType_count, error)
}

// HandleSubTypes godoc
// @Summary List subjects types
// @Description get subjects types
// @Tags sub types
// @Produce  json
// @Param page query int false "page number"
// @Param page_size query int false "page size"
// @Param subtypename query string false "subtypename search pattern"
// @Param subtypedescr query string false "subtypedescr search pattern"
// @Param ordering query string false "order by {subtypename|subtypedescr}"
// @Param desc query boolean false "descending order {true|false}"
// @Success 200 {object} models.SubType_count
// @Failure 500
// @Router /sub_types [get]
func HandleSubTypes(w http.ResponseWriter, r *http.Request) {
	var gs ifSubTypeService
	gs = services.NewSubTypeService(pgsql.SubTypeStorage{})
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

	stn := ""
	stns, ok := query["subtypename"]
	if ok && len(stns) > 0 {
		//case insensitive
		stn = strings.ToUpper(stns[0])
		//quotes
		re := regexp.MustCompile(`'`)
		stn = string(re.ReplaceAll([]byte(stn), []byte("''")))
	}

	std := ""
	stds, ok := query["subtypedescr"]
	if ok && len(stds) > 0 {
		std = strings.ToUpper(stds[0])
		re := regexp.MustCompile(`'`)
		std = string(re.ReplaceAll([]byte(std), []byte("''")))
	}

	ord := 1
	ords, ok := query["ordering"]
	if !ok || len(ords) == 0 {
		ord = 1
	} else if ords[0] == "subtypename" {
		ord = 2
	} else if ords[0] == "subtypedescr" {
		ord = 3
	}

	dsc := false
	dscs, ok := query["desc"]
	if ok && len(dscs) > 0 {
		if !(dscs[0] == "0") {
			dsc = true
		}
	}

	out_arr, err := gs.GetList(ctx, pg, pgs, stn, std, ord, dsc)
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

// HandleAddSubType godoc
// @Summary Add subjects types
// @Description add subjects types
// @Tags sub types
// @Accept json
// @Produce  json
// @Param ast body models.AddSubType true "New subType"
// @Success 200 {object} models.Json_id
// @Failure 500
// @Router /sub_types_add [post]
func HandleAddSubType(w http.ResponseWriter, r *http.Request) {
	var gs ifSubTypeService
	gs = services.NewSubTypeService(pgsql.SubTypeStorage{})
	ctx := context.Background()

	a := models.SubType{}
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
		log.Println("Failed execute ifSubTypeService.Add: ", err)
	}

	output, err := json.Marshal(models.Json_id{Id: ai})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(output)

	return

}

// HandleUpdSubType godoc
// @Summary Update subjects types
// @Description update subjects types
// @Tags sub types
// @Accept json
// @Produce  json
// @Param ust body models.SubType true "Update subtype"
// @Success 200 {object} models.Json_id
// @Failure 500
// @Router /sub_types_upd [post]
func HandleUpdSubType(w http.ResponseWriter, r *http.Request) {
	var gs ifSubTypeService
	gs = services.NewSubTypeService(pgsql.SubTypeStorage{})
	ctx := context.Background()

	u := models.SubType{}
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
		log.Println("Failed execute ifSubTypeService.Upd: ", err)
	}

	output, err := json.Marshal(models.Json_id{Id: ui})

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(output)

	return

}

// HandleDelSubType godoc
// @Summary Delete subjects types
// @Description delete subjects types
// @Tags sub types
// @Accept json
// @Produce  json
// @Param dst body models.Json_ids true "Delete subtypes"
// @Success 200 {object} models.Json_ids
// @Failure 500
// @Router /sub_types_del [post]
func HandleDelSubType(w http.ResponseWriter, r *http.Request) {
	var gs ifSubTypeService
	gs = services.NewSubTypeService(pgsql.SubTypeStorage{})
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
		log.Println("Failed execute ifSubTypeService.Del: ", err)
	}

	output, err := json.Marshal(models.Json_ids{Ids: res})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(output)

	return

}

// HandleGetSubType godoc
// @Summary Subjects type
// @Description get subject type
// @Tags sub types
// @Produce  json
// @Param id path int true "Subjects type by id"
// @Success 200 {object} models.SubType_count
// @Failure 500
// @Router /sub_types/{id} [get]
func HandleGetSubType(w http.ResponseWriter, r *http.Request) {
	var gs ifSubTypeService
	gs = services.NewSubTypeService(pgsql.SubTypeStorage{})
	ctx := context.Background()
	auth := utils.GetAuth(r)

	vars := mux.Vars(r)
	i, err := strconv.Atoi(vars["id"])
	if err != nil {
		i = 0
	}

	out_arr, err := gs.GetOne(ctx, i)
	if err != nil {
		log.Println("Failed execute ifSubTypeService.GetOne: ", err)
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
