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

type ifSealService interface {
	GetList(ctx context.Context, pg, pgs int, gs1 string, ord int, dsc bool) (models.Seal_count, error)
	Add(ctx context.Context, ea models.Seal) (int, error)
	Upd(ctx context.Context, eu models.Seal) (int, error)
	Del(ctx context.Context, ed []int) ([]int, error)
	GetOne(ctx context.Context, i int) (models.Seal_count, error)
}

// HandleSeals godoc
// @Summary List seals
// @Description get seal list
// @Tags seals
// @Produce json
// @Param page query int false "page number"
// @Param page_size query int false "page size"
// @Param packetnumber query string false "packetnumber search pattern"
// @Param ordering query string false "order by {id|packetnumber}"
// @Param desc query boolean false "descending order {true|false}"
// @Success 200 {object} models.Seal_count
// @Failure 500
// @Router /seals [get]
func HandleSeals(w http.ResponseWriter, r *http.Request) {
	var gs ifSealService
	gs = services.NewSealService(pgsql.SealStorage{})
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
	gs1s, ok := query["packetnumber"]
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
	} else if ords[0] == "packetnumber" {
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

// HandleAddSeal godoc
// @Summary Add seal
// @Description add seal
// @Tags seals
// @Accept json
// @Produce  json
// @Param a body models.AddSeal true "New seal. Significant params: PacketNumber, Area.Id, Staff.Id, SealType.Id, SealColour.Id, SealStatus.Id, IssueDate, ReportDate"
// @Success 200 {object} models.Json_id
// @Failure 500
// @Router /seals_add [post]
func HandleAddSeal(w http.ResponseWriter, r *http.Request) {
	var gs ifSealService
	gs = services.NewSealService(pgsql.SealStorage{})
	ctx := context.Background()

	a := models.Seal{}
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
		log.Println("Failed execute ifSealService.Add: ", err)
	}

	output, err := json.Marshal(models.Json_id{Id: ai})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(output)

	return

}

// HandleUpdSeal godoc
// @Summary Update seal
// @Description update seal
// @Tags seals
// @Accept json
// @Produce  json
// @Param u body models.Seal true "Update seal. Significant params: Id, PacketNumber, Area.Id, Staff.Id, SealType.Id, SealColour.Id, SealStatus.Id, IssueDate, ReportDate"
// @Success 200 {object} models.Json_id
// @Failure 500
// @Router /seals_upd [post]
func HandleUpdSeal(w http.ResponseWriter, r *http.Request) {
	var gs ifSealService
	gs = services.NewSealService(pgsql.SealStorage{})
	ctx := context.Background()

	u := models.Seal{}
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
		log.Println("Failed execute ifSealService.Upd: ", err)
	}

	output, err := json.Marshal(models.Json_id{Id: ui})

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(output)

	return

}

// HandleDelSeal godoc
// @Summary Delete seals
// @Description delete seals
// @Tags seals
// @Accept json
// @Produce  json
// @Param d body models.Json_ids true "Delete seals"
// @Success 200 {object} models.Json_ids
// @Failure 500
// @Router /seals_del [post]
func HandleDelSeal(w http.ResponseWriter, r *http.Request) {
	var gs ifSealService
	gs = services.NewSealService(pgsql.SealStorage{})
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
		log.Println("Failed execute ifSealService.Del: ", err)
	}

	output, err := json.Marshal(models.Json_ids{Ids: res})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(output)

	return

}

// HandleGetSeal godoc
// @Summary Get seal
// @Description get seal
// @Tags seals
// @Produce  json
// @Param id path int true "Seal by id"
// @Success 200 {object} models.Seal_count
// @Failure 500
// @Router /seals/{id} [get]
func HandleGetSeal(w http.ResponseWriter, r *http.Request) {
	var gs ifSealService
	gs = services.NewSealService(pgsql.SealStorage{})
	ctx := context.Background()
	auth := utils.GetAuth(r)

	vars := mux.Vars(r)
	i, err := strconv.Atoi(vars["id"])
	if err != nil {
		i = 0
	}

	out_arr, err := gs.GetOne(ctx, i)
	if err != nil {
		log.Println("Failed execute ifSealService.GetOne: ", err)
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
