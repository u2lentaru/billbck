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

type ifRpService interface {
	GetList(ctx context.Context, pg, pgs int, gs1, gs2 string, ord int, dsc bool) (models.Rp_count, error)
	Add(ctx context.Context, ea models.Rp) (int, error)
	Upd(ctx context.Context, eu models.Rp) (int, error)
	Del(ctx context.Context, ed []int) ([]int, error)
	GetOne(ctx context.Context, i int) (models.Rp_count, error)
}

// HandleRp godoc
// @Summary List rp
// @Description Get rp list
// @Tags rp
// @Produce  json
// @Param page query int false "page number"
// @Param page_size query int false "page size"
// @Param rpname query string false "rpname search pattern"
// @Param invnumber query string false "invnumber search pattern"
// @Param ordering query string false "order by {rpname|invnumber}"
// @Param desc query boolean false "descending order {true|false}"
// @Success 200 {object} models.Rp_count
// @Failure 500
// @Router /rp [get]
func HandleRp(w http.ResponseWriter, r *http.Request) {
	var gs ifRpService
	gs = services.NewRpService(pgsql.RpStorage{})
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
	gs1s, ok := query["rpname"]
	if ok && len(gs1s) > 0 {
		//case insensitive
		gs1 = strings.ToUpper(gs1s[0])
		//quotes
		re := regexp.MustCompile(`'`)
		gs1 = string(re.ReplaceAll([]byte(gs1), []byte("''")))
	}

	gs2 := ""
	gs2s, ok := query["invnumber"]
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
	} else if ords[0] == "rpname" {
		ord = 2
	} else if ords[0] == "invnumber" {
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

// HandleAddRp godoc
// @Summary Add rp
// @Description add rp
// @Tags rp
// @Accept json
// @Produce  json
// @Param a body models.AddRp true "New rp"
// @Success 200 {object} models.Json_id
// @Failure 500
// @Router /rp_add [post]
func HandleAddRp(w http.ResponseWriter, r *http.Request) {
	var gs ifRpService
	gs = services.NewRpService(pgsql.RpStorage{})
	ctx := context.Background()

	a := models.Rp{}
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
		log.Println("Failed execute ifRpService.Add: ", err)
	}

	output, err := json.Marshal(models.Json_id{Id: ai})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(output)

	return

}

// HandleUpdRp godoc
// @Summary Update rp
// @Description update rp
// @Tags rp
// @Accept json
// @Produce  json
// @Param u body models.Rp true "Update rp"
// @Success 200 {object} models.Json_id
// @Failure 500
// @Router /rp_upd [post]
func HandleUpdRp(w http.ResponseWriter, r *http.Request) {
	var gs ifRpService
	gs = services.NewRpService(pgsql.RpStorage{})
	ctx := context.Background()

	u := models.Rp{}
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
		log.Println("Failed execute ifRpService.Upd: ", err)
	}

	output, err := json.Marshal(models.Json_id{Id: ui})

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(output)

	return
}

// HandleDelRp godoc
// @Summary Delete rp
// @Description delete rp
// @Tags rp
// @Accept json
// @Produce  json
// @Param d body models.Json_ids true "Delete rp"
// @Success 200 {object} models.Json_ids
// @Failure 500
// @Router /rp_del [post]
func HandleDelRp(w http.ResponseWriter, r *http.Request) {
	var gs ifRpService
	gs = services.NewRpService(pgsql.RpStorage{})
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
		log.Println("Failed execute ifRpService.Del: ", err)
	}

	output, err := json.Marshal(models.Json_ids{Ids: res})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(output)

	return
}

// HandleGetRp godoc
// @Summary Get rp
// @Description get rp
// @Tags rp
// @Produce  json
// @Param id path int true "Rp by id"
// @Success 200 {object} models.Rp_count
// @Failure 500
// @Router /rp/{id} [get]
func HandleGetRp(w http.ResponseWriter, r *http.Request) {
	var gs ifRpService
	gs = services.NewRpService(pgsql.RpStorage{})
	ctx := context.Background()
	auth := utils.GetAuth(r)

	vars := mux.Vars(r)
	i, err := strconv.Atoi(vars["id"])
	if err != nil {
		i = 0
	}

	out_arr, err := gs.GetOne(ctx, i)
	if err != nil {
		log.Println("Failed execute ifRpService.GetOne: ", err)
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
