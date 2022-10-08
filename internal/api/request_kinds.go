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

type ifRequestKindService interface {
	GetList(ctx context.Context, pg, pgs int, gs1 string, ord int, dsc bool) (models.RequestKind_count, error)
	Add(ctx context.Context, ea models.RequestKind) (int, error)
	Upd(ctx context.Context, eu models.RequestKind) (int, error)
	Del(ctx context.Context, ed []int) ([]int, error)
	GetOne(ctx context.Context, i int) (models.RequestKind_count, error)
}

// HandleRequestKinds godoc
// @Summary List requestkinds
// @Description get requestkind list
// @Tags requestkinds
// @Produce json
// @Param page query int false "page number"
// @Param page_size query int false "page size"
// @Param requestkindname query string false "requestkindname search pattern"
// @Param ordering query string false "order by {id|requestkindname}"
// @Param desc query boolean false "descending order {true|false}"
// @Success 200 {object} models.RequestKind_count
// @Failure 500
// @Router /requestkinds [get]
func HandleRequestKinds(w http.ResponseWriter, r *http.Request) {
	var gs ifRequestKindService
	gs = services.NewRequestKindService(pgsql.RequestKindStorage{})
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
	gs1s, ok := query["requestkindname"]
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
	} else if ords[0] == "requestkindname" {
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

// HandleAddRequestKind godoc
// @Summary Add requestkind
// @Description add requestkind
// @Tags requestkinds
// @Accept json
// @Produce  json
// @Param a body models.AddRequestKind true "New requestkind"
// @Success 200 {object} models.Json_id
// @Failure 500
// @Router /requestkinds_add [post]
func HandleAddRequestKind(w http.ResponseWriter, r *http.Request) {
	var gs ifRequestKindService
	gs = services.NewRequestKindService(pgsql.RequestKindStorage{})
	ctx := context.Background()

	a := models.RequestKind{}
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
		log.Println("Failed execute ifRequestKindService.Add: ", err)
	}

	output, err := json.Marshal(models.Json_id{Id: ai})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(output)

	return

}

// HandleUpdRequestKind godoc
// @Summary Update requestkind
// @Description update requestkind
// @Tags requestkinds
// @Accept json
// @Produce  json
// @Param u body models.RequestKind true "Update requestkind"
// @Success 200 {object} models.Json_id
// @Failure 500
// @Router /requestkinds_upd [post]
func HandleUpdRequestKind(w http.ResponseWriter, r *http.Request) {
	var gs ifRequestKindService
	gs = services.NewRequestKindService(pgsql.RequestKindStorage{})
	ctx := context.Background()

	u := models.RequestKind{}
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
		log.Println("Failed execute ifRequestKindService.Upd: ", err)
	}

	output, err := json.Marshal(models.Json_id{Id: ui})

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(output)

	return
}

// HandleDelRequestKind godoc
// @Summary Delete requestkinds
// @Description delete requestkinds
// @Tags requestkinds
// @Accept json
// @Produce  json
// @Param d body models.Json_ids true "Delete requestkinds"
// @Success 200 {object} models.Json_ids
// @Failure 500
// @Router /requestkinds_del [post]
func HandleDelRequestKind(w http.ResponseWriter, r *http.Request) {
	var gs ifRequestKindService
	gs = services.NewRequestKindService(pgsql.RequestKindStorage{})
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
		log.Println("Failed execute ifRequestKindService.Del: ", err)
	}

	output, err := json.Marshal(models.Json_ids{Ids: res})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(output)

	return

}

// HandleGetRequestKind godoc
// @Summary Get requestkind
// @Description get requestkind
// @Tags requestkinds
// @Produce  json
// @Param id path int true "RequestKind by id"
// @Success 200 {object} models.RequestKind_count
// @Failure 500
// @Router /requestkinds/{id} [get]
func HandleGetRequestKind(w http.ResponseWriter, r *http.Request) {
	var gs ifRequestKindService
	gs = services.NewRequestKindService(pgsql.RequestKindStorage{})
	ctx := context.Background()
	auth := utils.GetAuth(r)

	vars := mux.Vars(r)
	i, err := strconv.Atoi(vars["id"])
	if err != nil {
		i = 0
	}

	out_arr, err := gs.GetOne(ctx, i)
	if err != nil {
		log.Println("Failed execute ifRequestKindService.GetOne: ", err)
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
