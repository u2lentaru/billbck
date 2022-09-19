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
)

type ifActDetailService interface {
	GetList(ctx context.Context, pg, pgs, nm, ord int, dsc bool) (models.ActDetail_count, error)
	Add(ctx context.Context, ea models.ActDetail) (int, error)
	Upd(ctx context.Context, eu models.ActDetail) (int, error)
	Del(ctx context.Context, ed []int) ([]int, error)
	GetOne(ctx context.Context, i int) (models.ActDetail_count, error)
}

// HandleActDetails godoc
// @Summary List act details
// @Description get act details list
// @Tags actdetails
// @Produce  json
// @Param page query int false "page number"
// @Param page_size query int false "page size"
// @Param actid query int false "actid search pattern"
// @Param ordering query string false "order by {punumber|installdate}"
// @Param desc query boolean false "descending order {true|false}"
// @Success 200 {object} models.ActDetail_count
// @Failure 500
// @Router /actdetails [get]
func HandleActDetails(w http.ResponseWriter, r *http.Request) {
	var gs ifActDetailService
	gs = services.NewActDetailService(pgsql.ActDetailStorage{})
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

	gs1 := 0
	gs1s, ok := query["actid"]
	if ok && len(gs1s) > 0 {
		t, err := strconv.Atoi(gs1s[0])
		if err == nil {
			gs1 = t
		}
	}

	ord := 10
	ords, ok := query["ordering"]
	if !ok || len(ords) == 0 {
		ord = 10
	} else if ords[0] == "punumber" {
		ord = 12
	} else if ords[0] == "installdate" {
		ord = 15
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

// HandleAddActDetail godoc
// @Summary Add act detail
// @Description add act detail
// @Tags actdetails
// @Accept json
// @Produce  json
// @Param a body models.AddActDetail true "New act detail. Significant params: Act.Id, PuId, SealNumber, AdPuValue, ActDetailDate, PuNumber, InstallDate, CheckInterval, InitialValue, DevStopped, Startdate, Enddate, Pid, Seal.Id, SealDate, Notes, PuType.Id, Conclusion.Id, ConclusionNumber, ShutdownType.Id, CustomerPhone, CustomerPos, Reason.Id, Violation.Id, Customer"
// @Success 200 {object} models.Json_id
// @Failure 500
// @Router /actdetails_add [post]
func HandleAddActDetail(w http.ResponseWriter, r *http.Request) {
	var gs ifActDetailService
	gs = services.NewActDetailService(pgsql.ActDetailStorage{})
	ctx := context.Background()

	a := models.ActDetail{}
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
		log.Println("Failed execute ifActDetailService.Add: ", err)
	}

	output, err := json.Marshal(models.Json_id{Id: ai})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(output)

	return
}

// HandleUpdActDetail godoc
// @Summary Update act detail
// @Description update act detail
// @Tags actdetails
// @Accept json
// @Produce  json
// @Param u body models.ActDetail true "Update act detail. Significant params: Id, Act.Id, PuId, SealNumber, AdPuValue, ActDetailDate, PuNumber, InstallDate, CheckInterval, InitialValue, DevStopped, Startdate, Enddate, Pid, Seal.Id, SealDate, Notes, PuType.Id, Conclusion.Id, ConclusionNumber, ShutdownType.Id, CustomerPhone, CustomerPos, Reason.Id, Violation.Id, Customer"
// @Success 200 {object} models.Json_id
// @Failure 500
// @Router /actdetails_upd [post]
func HandleUpdActDetail(w http.ResponseWriter, r *http.Request) {
	var gs ifActDetailService
	gs = services.NewActDetailService(pgsql.ActDetailStorage{})
	ctx := context.Background()

	u := models.ActDetail{}
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
		log.Println("Failed execute ifActDetailService.Upd: ", err)
	}

	output, err := json.Marshal(models.Json_id{Id: ui})

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(output)

	return

}

// HandleDelActDetail godoc
// @Summary Delete act details
// @Description delete act details
// @Tags actdetails
// @Accept json
// @Produce  json
// @Param d body models.Json_ids true "Delete act details"
// @Success 200 {object} models.Json_ids
// @Failure 500
// @Router /actdetails_del [post]
func HandleDelActDetail(w http.ResponseWriter, r *http.Request) {
	var gs ifActDetailService
	gs = services.NewActDetailService(pgsql.ActDetailStorage{})
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
		log.Println("Failed execute ifActDetailService.Del: ", err)
	}

	output, err := json.Marshal(models.Json_ids{Ids: res})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(output)

	return
}

// HandleGetActDetail godoc
// @Summary Get act detail
// @Description get act detail
// @Tags actdetails
// @Produce  json
// @Param id path int true "Act detail by id"
// @Success 200 {object} models.ActDetail_count
// @Failure 500
// @Router /actdetails/{id} [get]
func HandleGetActDetail(w http.ResponseWriter, r *http.Request) {
	var gs ifActDetailService
	gs = services.NewActDetailService(pgsql.ActDetailStorage{})
	ctx := context.Background()

	vars := mux.Vars(r)
	i, err := strconv.Atoi(vars["id"])
	if err != nil {
		i = 0
	}

	out_arr, err := gs.GetOne(ctx, i)
	if err != nil {
		log.Println("Failed execute ifActDetailService.GetOne: ", err)
	}

	out_count, err := json.Marshal(out_arr)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(out_count)

	return
}
