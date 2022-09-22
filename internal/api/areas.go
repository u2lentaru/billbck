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

type ifAreaService interface {
	GetList(ctx context.Context, pg, pgs int, gs1, gs2 string, ord int, dsc bool) (models.Area_count, error)
	Add(ctx context.Context, ea models.Area) (int, error)
	Upd(ctx context.Context, eu models.Area) (int, error)
	Del(ctx context.Context, ed []int) ([]int, error)
	GetOne(ctx context.Context, i int) (models.Area_count, error)
}

// HandleAreas godoc
// @Summary List areas
// @Description get areas
// @Tags areas
// @Produce  json
// @Param page query int false "page number"
// @Param page_size query int false "page size"
// @Param areanumber query string false "areanumber search pattern"
// @Param areaname query string false "areaname search pattern"
// @Param ordering query string false "order by {areanumber|areaname}"
// @Param desc query boolean false "descending order {true|false}"
// @Success 200 {object} models.Area_count
// @Failure 500
// @Router /areas [get]
func HandleAreas(w http.ResponseWriter, r *http.Request) {
	var gs ifAreaService
	gs = services.NewAreaService(pgsql.AreaStorage{})
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
	gs1s, ok := query["areanumber"]
	if ok && len(gs1s) > 0 {
		//case insensitive
		gs1 = strings.ToUpper(gs1s[0])
		//quotes
		re := regexp.MustCompile(`'`)
		gs1 = string(re.ReplaceAll([]byte(gs1), []byte("''")))
	}

	gs2 := ""
	gs2s, ok := query["areaname"]
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
	} else if ords[0] == "areanumber" {
		ord = 2
	} else if ords[0] == "areaname" {
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

// HandleAddArea godoc
// @Summary Add area
// @Description add area
// @Tags areas
// @Accept json
// @Produce  json
// @Param a body models.AddArea true "New area"
// @Success 200 {object} models.Json_id
// @Failure 500
// @Router /areas_add [post]
func HandleAddArea(w http.ResponseWriter, r *http.Request) {
	var gs ifAreaService
	gs = services.NewAreaService(pgsql.AreaStorage{})
	ctx := context.Background()

	a := models.Area{}
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
		log.Println("Failed execute ifAreaService.Add: ", err)
	}

	output, err := json.Marshal(models.Json_id{Id: ai})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(output)

	return

}

// HandleUpdArea godoc
// @Summary Update area
// @Description update area
// @Tags areas
// @Accept json
// @Produce  json
// @Param u body models.Area true "Update area"
// @Success 200 {object} models.Json_id
// @Failure 500
// @Router /areas_upd [post]
func HandleUpdArea(w http.ResponseWriter, r *http.Request) {
	var gs ifAreaService
	gs = services.NewAreaService(pgsql.AreaStorage{})
	ctx := context.Background()

	u := models.Area{}
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
		log.Println("Failed execute ifAreaService.Upd: ", err)
	}

	output, err := json.Marshal(models.Json_id{Id: ui})

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(output)

	return
}

// HandleDelArea godoc
// @Summary Delete areas
// @Description delete areas
// @Tags areas
// @Accept json
// @Produce  json
// @Param d body models.Json_ids true "Delete areas"
// @Success 200 {object} models.Json_ids
// @Failure 500
// @Router /areas_del [post]
func HandleDelArea(w http.ResponseWriter, r *http.Request) {
	var gs ifAreaService
	gs = services.NewAreaService(pgsql.AreaStorage{})
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
		log.Println("Failed execute ifAreaService.Del: ", err)
	}

	output, err := json.Marshal(models.Json_ids{Ids: res})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(output)

	return
}

// HandleGetArea godoc
// @Summary Get area
// @Description get area
// @Tags areas
// @Produce  json
// @Param id path int true "Area by id"
// @Success 200 {object} models.Area_count
// @Failure 500
// @Router /areas/{id} [get]
func HandleGetArea(w http.ResponseWriter, r *http.Request) {
	var gs ifAreaService
	gs = services.NewAreaService(pgsql.AreaStorage{})
	ctx := context.Background()
	auth := utils.GetAuth(r)

	vars := mux.Vars(r)
	i, err := strconv.Atoi(vars["id"])
	if err != nil {
		i = 0
	}

	out_arr, err := gs.GetOne(ctx, i)
	if err != nil {
		log.Println("Failed execute ifAreaService.GetOne: ", err)
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
