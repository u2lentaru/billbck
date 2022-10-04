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

type ifOrgInfoService interface {
	GetList(ctx context.Context, pg, pgs int, gs1, gs2 string, ord int, dsc bool) (models.OrgInfo_count, error)
	Add(ctx context.Context, ea models.OrgInfo) (int, error)
	Upd(ctx context.Context, eu models.OrgInfo) (int, error)
	Del(ctx context.Context, ed []int) ([]int, error)
	GetOne(ctx context.Context, i int) (models.OrgInfo_count, error)
}

// HandleOrgInfos godoc
// @Summary List org_info
// @Description Get org_info list
// @Tags org_info
// @Produce  json
// @Param page query int false "page number"
// @Param page_size query int false "page size"
// @Param oiname query string false "oiname search pattern"
// @Param oifname query string false "oifname search pattern"
// @Param ordering query string false "order by {oiname|oifname}"
// @Param desc query boolean false "descending order {true|false}"
// @Success 200 {object} models.OrgInfo_count
// @Failure 500
// @Router /org_info [get]
func HandleOrgInfos(w http.ResponseWriter, r *http.Request) {
	var gs ifOrgInfoService
	gs = services.NewOrgInfoService(pgsql.OrgInfoStorage{})
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

	oin := ""
	oins, ok := query["oiname"]
	if ok && len(oins) > 0 {
		//case insensitive
		oin = strings.ToUpper(oins[0])
		//quotes
		re := regexp.MustCompile(`'`)
		oin = string(re.ReplaceAll([]byte(oin), []byte("''")))
	}

	oifn := ""
	oifns, ok := query["oifname"]
	if ok && len(oifns) > 0 {
		oifn = strings.ToUpper(oifns[0])
		re := regexp.MustCompile(`'`)
		oifn = string(re.ReplaceAll([]byte(oifn), []byte("''")))
	}

	ord := 1
	ords, ok := query["ordering"]
	if !ok || len(ords) == 0 {
		ord = 1
	} else if ords[0] == "oiname" {
		ord = 2
	} else if ords[0] == "oifname" {
		ord = 7
	}

	dsc := false
	dscs, ok := query["desc"]
	if ok && len(dscs) > 0 {
		if !(dscs[0] == "0") {
			dsc = true
		}
	}

	out_arr, err := gs.GetList(ctx, pg, pgs, oin, oifn, ord, dsc)
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

// HandleAddOrgInfo godoc
// @Summary Add org_info
// @Description Add org_info
// @Tags org_info
// @Accept json
// @Produce  json
// @Param ab body models.AddOrgInfo true "New org_info"
// @Success 200 {object} models.Json_id
// @Failure 500
// @Router /org_info_add [post]
func HandleAddOrgInfo(w http.ResponseWriter, r *http.Request) {
	var gs ifOrgInfoService
	gs = services.NewOrgInfoService(pgsql.OrgInfoStorage{})
	ctx := context.Background()

	a := models.OrgInfo{}
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
		log.Println("Failed execute ifOrgInfoService.Add: ", err)
	}

	output, err := json.Marshal(models.Json_id{Id: ai})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(output)

	return
}

// HandleUpdOrgInfo godoc
// @Summary Update org_info
// @Description Update org_info
// @Tags org_info
// @Accept json
// @Produce  json
// @Param ub body models.OrgInfo true "Update org_info"
// @Success 200 {object} models.Json_id
// @Failure 500
// @Router /org_info_upd [post]
func HandleUpdOrgInfo(w http.ResponseWriter, r *http.Request) {
	var gs ifOrgInfoService
	gs = services.NewOrgInfoService(pgsql.OrgInfoStorage{})
	ctx := context.Background()

	u := models.OrgInfo{}
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
		log.Println("Failed execute ifOrgInfoService.Upd: ", err)
	}

	output, err := json.Marshal(models.Json_id{Id: ui})

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(output)

	return
}

// HandleDelOrgInfo godoc
// @Summary Delete org_info
// @Description Delete org_info
// @Tags org_info
// @Accept json
// @Produce  json
// @Param db body models.Json_ids true "Delete org_info"
// @Success 200 {object} models.Json_id
// @Failure 500
// @Router /org_info_del [post]
func HandleDelOrgInfo(w http.ResponseWriter, r *http.Request) {
	var gs ifOrgInfoService
	gs = services.NewOrgInfoService(pgsql.OrgInfoStorage{})
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
		log.Println("Failed execute ifOrgInfoService.Del: ", err)
	}

	output, err := json.Marshal(models.Json_ids{Ids: res})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(output)

	return
}

// HandleGetOrgInfo godoc
// @Summary Get org_info
// @Description Get org_info
// @Tags org_info
// @Produce  json
// @Param id path int true "OrgInfo by id"
// @Success 200 {object} models.OrgInfo_count
// @Failure 500
// @Router /org_info/{id} [get]
func HandleGetOrgInfo(w http.ResponseWriter, r *http.Request) {
	var gs ifOrgInfoService
	gs = services.NewOrgInfoService(pgsql.OrgInfoStorage{})
	ctx := context.Background()
	auth := utils.GetAuth(r)

	vars := mux.Vars(r)
	i, err := strconv.Atoi(vars["id"])
	if err != nil {
		i = 0
	}

	out_arr, err := gs.GetOne(ctx, i)
	if err != nil {
		log.Println("Failed execute ifOrgInfoService.GetOne: ", err)
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
