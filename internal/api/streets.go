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

type ifStreetService interface {
	GetList(ctx context.Context, pg, pgs int, gs1 string, gs2, ord int, dsc bool) (models.Street_count, error)
	Add(ctx context.Context, ea models.Street) (int, error)
	Upd(ctx context.Context, eu models.Street) (int, error)
	Del(ctx context.Context, d models.StreetClose) (models.Json_id, error)
	GetOne(ctx context.Context, i int) (models.Street_count, error)
}

// HandleStreets godoc
// @Summary List streets
// @Description get street list
// @Tags streets
// @Produce  json
// @Param page query int false "page number"
// @Param page_size query int false "page size"
// @Param streetname query string false "streetname search pattern"
// @Param cityid query int false "cityid search pattern"
// @Param ordering query string false "order by {id|streetname}"
// @Param desc query boolean false "descending order {true|false}"
// @Success 200 {object} models.Street_count
// @Failure 500
// @Router /streets [get]
func HandleStreets(w http.ResponseWriter, r *http.Request) {
	var gs ifStreetService
	gs = services.NewStreetService(pgsql.StreetStorage{})
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
	gs1s, ok := query["streetname"]
	if ok && len(gs1s) > 0 {
		//case insensitive
		gs1 = strings.ToUpper(gs1s[0])
		//quotes
		re := regexp.MustCompile(`'`)
		gs1 = string(re.ReplaceAll([]byte(gs1), []byte("''")))
	}

	gs2 := 0
	gs2s, ok := query["cityid"]
	if ok && len(gs2s) > 0 {
		t, err := strconv.Atoi(gs2s[0])
		if err == nil {
			gs2 = t
		}
	}

	ord := 1
	ords, ok := query["ordering"]
	if !ok || len(ords) == 0 {
		ord = 1
	} else if ords[0] == "streetname" {
		ord = 2
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

// HandleAddStreet godoc
// @Summary Add street
// @Description add street
// @Tags streets
// @Accept json
// @Produce  json
// @Param a body models.AddStreet true "New street"
// @Success 200 {object} models.Json_id
// @Failure 500
// @Router /streets_add [post]
func HandleAddStreet(w http.ResponseWriter, r *http.Request) {
	var gs ifStreetService
	gs = services.NewStreetService(pgsql.StreetStorage{})
	ctx := context.Background()

	a := models.Street{}
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
		log.Println("Failed execute ifStreetService.Add: ", err)
	}

	output, err := json.Marshal(models.Json_id{Id: ai})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(output)

	return

}

// HandleUpdStreet godoc
// @Summary Update street
// @Description update street
// @Tags streets
// @Accept json
// @Produce  json
// @Param u body models.Street true "Update street"
// @Success 200 {object} models.Json_id
// @Failure 500
// @Router /streets_upd [post]
func HandleUpdStreet(w http.ResponseWriter, r *http.Request) {
	var gs ifStreetService
	gs = services.NewStreetService(pgsql.StreetStorage{})
	ctx := context.Background()

	u := models.Street{}
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
		log.Println("Failed execute ifStreetService.Upd: ", err)
	}

	output, err := json.Marshal(models.Json_id{Id: ui})

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(output)

	return
}

// HandleDelStreet godoc
// @Summary Delete street
// @Description delete street
// @Tags streets
// @Accept json
// @Produce  json
// @Param d body models.StreetClose true "Delete street"
// @Success 200 {object} models.Json_ids
// @Failure 500
// @Router /streets_del [post]
func HandleDelStreet(w http.ResponseWriter, r *http.Request) {
	var gs ifStreetService
	gs = services.NewStreetService(pgsql.StreetStorage{})
	ctx := context.Background()

	d := models.StreetClose{}
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

	i, err := gs.Del(ctx, d)
	if err != nil {
		log.Println("Failed execute ifStreetService.Del: ", err)
	}

	output, err := json.Marshal(i)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(output)

	return

}

// HandleGetStreet godoc
// @Summary Get street
// @Description get street
// @Tags streets
// @Produce  json
// @Param id path int true "Street by id"
// @Success 200 {object} models.Street_count
// @Failure 500
// @Router /streets/{id} [get]
func HandleGetStreet(w http.ResponseWriter, r *http.Request) {
	var gs ifStreetService
	gs = services.NewStreetService(pgsql.StreetStorage{})
	ctx := context.Background()
	auth := utils.GetAuth(r)

	vars := mux.Vars(r)
	i, err := strconv.Atoi(vars["id"])
	if err != nil {
		i = 0
	}

	out_arr, err := gs.GetOne(ctx, i)
	if err != nil {
		log.Println("Failed execute ifStreetService.GetOne: ", err)
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
