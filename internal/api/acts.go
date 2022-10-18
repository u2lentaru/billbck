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

type ifActService interface {
	GetList(ctx context.Context, pg, pgs int, gs1, gs2 string, gs3, ord int, dsc bool) (models.Act_count, error)
	Add(ctx context.Context, ea models.Act) (int, error)
	Upd(ctx context.Context, eu models.Act) (int, error)
	Del(ctx context.Context, ed []int) ([]int, error)
	GetOne(ctx context.Context, i int) (models.Act_count, error)
	Activate(ctx context.Context, i int, d string) (int, error)
}

// HandleActs godoc
// @Summary List acts
// @Description get act list
// @Tags acts
// @Produce  json
// @Param page query int false "page number"
// @Param page_size query int false "page size"
// @Param actnumber query string false "actnumber search pattern"
// @Param objectname query string false "objectname search pattern"
// @Param objectid query int false "objectid search pattern"
// @Param ordering query string false "order by {actnumber|objectname}"
// @Param desc query boolean false "descending order {true|false}"
// @Success 200 {object} models.Act_count
// @Failure 500
// @Router /acts [get]
func HandleActs(w http.ResponseWriter, r *http.Request) {
	var gs ifActService
	gs = services.NewActService(pgsql.ActStorage{})
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
	gs1s, ok := query["actnumber"]
	if ok && len(gs1s) > 0 {
		//case insensitive
		gs1 = strings.ToUpper(gs1s[0])
		//quotes
		re := regexp.MustCompile(`'`)
		gs1 = string(re.ReplaceAll([]byte(gs1), []byte("''")))
	}

	gs2 := ""
	gs2s, ok := query["objectname"]
	if ok && len(gs2s) > 0 {
		//case insensitive
		gs2 = strings.ToUpper(gs2s[0])
		//quotes
		re := regexp.MustCompile(`'`)
		gs2 = string(re.ReplaceAll([]byte(gs2), []byte("''")))
	}

	gs3 := 0
	gs3s, ok := query["objectid"]
	if ok && len(gs3s) > 0 {
		t, err := strconv.Atoi(gs3s[0])
		if err == nil {
			gs3 = t
		}
	}

	ord := 1
	ords, ok := query["ordering"]
	if !ok || len(ords) == 0 {
		ord = 1
	} else if ords[0] == "actnumber" {
		ord = 3
	} else if ords[0] == "objectname" {
		ord = 7
	}

	dsc := false
	dscs, ok := query["desc"]
	if ok && len(dscs) > 0 {
		if !(dscs[0] == "0") {
			dsc = true
		}
	}

	out_arr, err := gs.GetList(ctx, pg, pgs, gs1, gs2, gs3, ord, dsc)
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

// HandleAddAct godoc
// @Summary Add act
// @Description add act
// @Tags acts
// @Accept json
// @Produce  json
// @Param a body models.AddAct true "New act. Significant params: ActType.Id, ActNumber, ActDate, Object.Id, Staff.Id, Notes"
// @Success 200 {object} models.Json_id
// @Failure 500
// @Router /acts_add [post]
func HandleAddAct(w http.ResponseWriter, r *http.Request) {
	var gs ifActService
	gs = services.NewActService(pgsql.ActStorage{})
	ctx := context.Background()

	a := models.Act{}
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
		log.Println("Failed execute ifActService.Add: ", err)
	}

	output, err := json.Marshal(models.Json_id{Id: ai})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(output)

	return

}

// HandleUpdAct godoc
// @Summary Update act
// @Description update act
// @Tags acts
// @Accept json
// @Produce  json
// @Param u body models.Act true "Update act. Significant params: Id, ActType.Id, ActNumber, ActDate, Object.Id, Staff.Id, Notes"
// @Success 200 {object} models.Json_id
// @Failure 500
// @Router /acts_upd [post]
func HandleUpdAct(w http.ResponseWriter, r *http.Request) {
	var gs ifActService
	gs = services.NewActService(pgsql.ActStorage{})
	ctx := context.Background()

	u := models.Act{}
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
		log.Println("Failed execute ifActService.Upd: ", err)
	}

	output, err := json.Marshal(models.Json_id{Id: ui})

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(output)

	return
}

// HandleDelAct godoc
// @Summary Delete acts
// @Description delete acts
// @Tags acts
// @Accept json
// @Produce  json
// @Param d body models.Json_ids true "Delete acts"
// @Success 200 {object} models.Json_ids
// @Failure 500
// @Router /acts_del [post]
func HandleDelAct(w http.ResponseWriter, r *http.Request) {
	var gs ifActService
	gs = services.NewActService(pgsql.ActStorage{})
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
		log.Println("Failed execute ifActService.Del: ", err)
	}

	output, err := json.Marshal(models.Json_ids{Ids: res})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(output)

	return
}

// HandleGetAct godoc
// @Summary Get act
// @Description get act
// @Tags acts
// @Produce  json
// @Param id path int true "Act by id"
// @Success 200 {object} models.Act_count
// @Failure 500
// @Router /acts/{id} [get]
func HandleGetAct(w http.ResponseWriter, r *http.Request) {
	var gs ifActService
	gs = services.NewActService(pgsql.ActStorage{})
	ctx := context.Background()
	auth := utils.GetAuth(r)

	vars := mux.Vars(r)
	i, err := strconv.Atoi(vars["id"])
	if err != nil {
		i = 0
	}

	out_arr, err := gs.GetOne(ctx, i)
	if err != nil {
		log.Println("Failed execute ifActService.GetOne: ", err)
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

// HandleActActivate godoc
// @Summary Act activation
// @Description act activation
// @Tags acts
// @Produce  json
// @Param actid query int true "actid"
// @Param activationdate query string true "activation date"
// @Success 200 {object} models.Json_id
// @Failure 500
// @Router /acts_activate [get]
func HandleActActivate(w http.ResponseWriter, r *http.Request) {
	var gs ifActService
	gs = services.NewActService(pgsql.ActStorage{})
	ctx := context.Background()

	query := r.URL.Query()

	gs1 := 0
	gs1s, ok := query["actid"]
	if ok && len(gs1s) > 0 {
		t, err := strconv.Atoi(gs1s[0])
		if err == nil {
			gs1 = t
		}
	}

	gs2 := ""
	gs2s, ok := query["activationdate"]
	if ok && len(gs2s) > 0 {
		//case insensitive
		gs2 = strings.ToUpper(gs2s[0])
		//quotes
		re := regexp.MustCompile(`'`)
		//gs2, err := time.Parse("2006-01-02", gs2s[0])
		gs2 = string(re.ReplaceAll([]byte(gs2), []byte("''")))
	}

	ai, err := gs.Activate(ctx, gs1, gs2)

	if err != nil {
		log.Println("Failed execute ifActService.Activate: ", err)
	}

	output, err := json.Marshal(models.Json_id{Id: ai})

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(output)

	return
}
