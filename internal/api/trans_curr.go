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

type ifTransCurrService interface {
	GetList(ctx context.Context, pg, pgs int, gs1 string, ord int, dsc bool) (models.TransCurr_count, error)
	Add(ctx context.Context, ea models.TransCurr) (int, error)
	Upd(ctx context.Context, eu models.TransCurr) (int, error)
	Del(ctx context.Context, ed []int) ([]int, error)
	GetOne(ctx context.Context, i int) (models.TransCurr_count, error)
}

// HandleTransCurr godoc
// @Summary List current transformers
// @Description get current transformers list
// @Tags transcurr
// @Produce json
// @Param page query int false "page number"
// @Param page_size query int false "page size"
// @Param transcurrname query string false "transcurrname search pattern"
// @Param ordering query string false "order by {id|transcurrname}"
// @Param desc query boolean false "descending order {true|false}"
// @Success 200 {object} models.TransCurr_count
// @Failure 500
// @Router /transcurr [get]
func HandleTransCurr(w http.ResponseWriter, r *http.Request) {
	var gs ifTransCurrService
	gs = services.NewTransCurrService(pgsql.TransCurrStorage{})
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
	gs1s, ok := query["transcurrname"]
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
	} else if ords[0] == "transcurrname" {
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

// HandleAddTransCurr godoc
// @Summary Add current transformer
// @Description add current transformer
// @Tags transcurr
// @Accept json
// @Produce  json
// @Param a body models.AddTransCurr true "New current transformer. Significant params: TransCurrName, TransType.Id, CheckDate(n), NextCheckDate(n), ProdDate(n), Serial1(n), Serial2(n), Serial3(n)"
// @Success 200 {object} models.Json_id
// @Failure 500
// @Router /transcurr_add [post]
func HandleAddTransCurr(w http.ResponseWriter, r *http.Request) {
	var gs ifTransCurrService
	gs = services.NewTransCurrService(pgsql.TransCurrStorage{})
	ctx := context.Background()

	a := models.TransCurr{}
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
		log.Println("Failed execute ifTransCurrService.Add: ", err)
	}

	output, err := json.Marshal(models.Json_id{Id: ai})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(output)

	return

}

// HandleUpdTransCurr godoc
// @Summary Update current transformer
// @Description update current transformer
// @Tags transcurr
// @Accept json
// @Produce  json
// @Param u body models.TransCurr true "Update current transformer. Significant params: Id, TransCurrName, TransType.Id, CheckDate(n), NextCheckDate(n), ProdDate(n), Serial1(n), Serial2(n), Serial3(n)"
// @Success 200 {object} models.Json_id
// @Failure 500
// @Router /transcurr_upd [post]
func HandleUpdTransCurr(w http.ResponseWriter, r *http.Request) {
	var gs ifTransCurrService
	gs = services.NewTransCurrService(pgsql.TransCurrStorage{})
	ctx := context.Background()

	u := models.TransCurr{}
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
		log.Println("Failed execute ifTransCurrService.Upd: ", err)
	}

	output, err := json.Marshal(models.Json_id{Id: ui})

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(output)

	return

}

// HandleDelTransCurr godoc
// @Summary Delete current transformers
// @Description delete current transformers
// @Tags transcurr
// @Accept json
// @Produce  json
// @Param d body models.Json_ids true "Delete current transformers"
// @Success 200 {object} models.Json_ids
// @Failure 500
// @Router /transcurr_del [post]
func HandleDelTransCurr(w http.ResponseWriter, r *http.Request) {
	var gs ifTransCurrService
	gs = services.NewTransCurrService(pgsql.TransCurrStorage{})
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
		log.Println("Failed execute ifTransCurrService.Del: ", err)
	}

	output, err := json.Marshal(models.Json_ids{Ids: res})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(output)

	return

}

// HandleGetTransCurr godoc
// @Summary Get current transformer
// @Description get current transformer
// @Tags transcurr
// @Produce  json
// @Param id path int true "Current transformer by id"
// @Success 200 {object} models.TransCurr_count
// @Failure 500
// @Router /transcurr/{id} [get]
func HandleGetTransCurr(w http.ResponseWriter, r *http.Request) {
	var gs ifTransCurrService
	gs = services.NewTransCurrService(pgsql.TransCurrStorage{})
	ctx := context.Background()
	auth := utils.GetAuth(r)

	vars := mux.Vars(r)
	i, err := strconv.Atoi(vars["id"])
	if err != nil {
		i = 0
	}

	out_arr, err := gs.GetOne(ctx, i)
	if err != nil {
		log.Println("Failed execute ifTransCurrService.GetOne: ", err)
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
