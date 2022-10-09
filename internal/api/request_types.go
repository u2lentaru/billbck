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

type ifRequestTypeService interface {
	GetList(ctx context.Context, pg, pgs int, gs1, gs2 string, ord int, dsc bool) (models.RequestType_count, error)
	Add(ctx context.Context, ea models.RequestType) (int, error)
	Upd(ctx context.Context, eu models.RequestType) (int, error)
	Del(ctx context.Context, ed []int) ([]int, error)
	GetOne(ctx context.Context, i int) (models.RequestType_count, error)
}

// HandleRequestTypes godoc
// @Summary List requesttypes
// @Description get requesttype list
// @Tags requesttypes
// @Produce json
// @Param page query int false "page number"
// @Param page_size query int false "page size"
// @Param requesttypename query string false "requesttypename search pattern"
// @Param rkid query string false "request kind id search pattern"
// @Param ordering query string false "order by {id|requesttypename}"
// @Param desc query boolean false "descending order {true|false}"
// @Success 200 {object} models.RequestType_count
// @Failure 500
// @Router /requesttypes [get]
func HandleRequestTypes(w http.ResponseWriter, r *http.Request) {
	var gs ifRequestTypeService
	gs = services.NewRequestTypeService(pgsql.RequestTypeStorage{})
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
	gs1s, ok := query["requesttypename"]
	if ok && len(gs1s) > 0 {
		gs1 = gs1s[0]
	}

	gs2 := ""
	gs2s, ok := query["rkid"]
	if ok && len(gs2s) > 0 {
		_, err := strconv.Atoi(gs2s[0])
		if err == nil {
			gs2 = gs2s[0]
		}
	}

	ord := 1
	ords, ok := query["ordering"]
	if !ok || len(ords) == 0 {
		ord = 1
	} else if ords[0] == "requesttypename" {
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

// HandleAddRequestType godoc
// @Summary Add requesttype
// @Description add requesttype
// @Tags requesttypes
// @Accept json
// @Produce  json
// @Param a body models.AddRequestType true "New requesttype"
// @Success 200 {object} models.Json_id
// @Failure 500
// @Router /requesttypes_add [post]
func HandleAddRequestType(w http.ResponseWriter, r *http.Request) {
	var gs ifRequestTypeService
	gs = services.NewRequestTypeService(pgsql.RequestTypeStorage{})
	ctx := context.Background()

	a := models.RequestType{}
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
		log.Println("Failed execute ifRequestTypeService.Add: ", err)
	}

	output, err := json.Marshal(models.Json_id{Id: ai})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(output)

	return

}

// HandleUpdRequestType godoc
// @Summary Update requesttype
// @Description update requesttype
// @Tags requesttypes
// @Accept json
// @Produce  json
// @Param u body models.RequestType true "Update requesttype"
// @Success 200 {object} models.Json_id
// @Failure 500
// @Router /requesttypes_upd [post]
func HandleUpdRequestType(w http.ResponseWriter, r *http.Request) {
	var gs ifRequestTypeService
	gs = services.NewRequestTypeService(pgsql.RequestTypeStorage{})
	ctx := context.Background()

	u := models.RequestType{}
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
		log.Println("Failed execute ifRequestTypeService.Upd: ", err)
	}

	output, err := json.Marshal(models.Json_id{Id: ui})

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(output)

	return
}

// HandleDelRequestType godoc
// @Summary Delete requesttypes
// @Description delete requesttypes
// @Tags requesttypes
// @Accept json
// @Produce  json
// @Param d body models.Json_ids true "Delete requesttypes"
// @Success 200 {object} models.Json_ids
// @Failure 500
// @Router /requesttypes_del [post]
func HandleDelRequestType(w http.ResponseWriter, r *http.Request) {
	var gs ifRequestTypeService
	gs = services.NewRequestTypeService(pgsql.RequestTypeStorage{})
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
		log.Println("Failed execute ifRequestTypeService.Del: ", err)
	}

	output, err := json.Marshal(models.Json_ids{Ids: res})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(output)

	return

}

// HandleGetRequestType godoc
// @Summary Get requesttype
// @Description get requesttype
// @Tags requesttypes
// @Produce  json
// @Param id path int true "RequestType by id"
// @Success 200 {object} models.RequestType_count
// @Failure 500
// @Router /requesttypes/{id} [get]
func HandleGetRequestType(w http.ResponseWriter, r *http.Request) {
	var gs ifRequestTypeService
	gs = services.NewRequestTypeService(pgsql.RequestTypeStorage{})
	ctx := context.Background()
	auth := utils.GetAuth(r)

	vars := mux.Vars(r)
	i, err := strconv.Atoi(vars["id"])
	if err != nil {
		i = 0
	}

	out_arr, err := gs.GetOne(ctx, i)
	if err != nil {
		log.Println("Failed execute ifRequestTypeService.GetOne: ", err)
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
