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
	"time"

	"github.com/gorilla/mux"
	"github.com/u2lentaru/billbck/internal/adapters/db/pgsql"
	"github.com/u2lentaru/billbck/internal/models"
	"github.com/u2lentaru/billbck/internal/services"
	"github.com/u2lentaru/billbck/internal/utils"
)

type ifPuService interface {
	GetList(ctx context.Context, pg, pgs int, gs1, gs2, gs3, gs4, gs5, gs6, gs7 string, ord int, dsc bool) (models.Pu_count, error)
	Add(ctx context.Context, ea models.Pu) (int, error)
	Upd(ctx context.Context, eu models.Pu) (int, error)
	Del(ctx context.Context, ed []int) ([]int, error)
	GetOne(ctx context.Context, i int) (models.Pu_count, error)
	GetObj(ctx context.Context, gs1, gs2 string) (models.Pu_count, error)
}

// HandlePu godoc
// @Summary List pu
// @Description get pu list
// @Tags pu
// @Produce  json
// @Param page query int false "page number"
// @Param page_size query int false "page size"
// @Param objectname query string false "objectname search pattern"
// @Param streetname query string false "streetname search pattern"
// @Param actualdate query string false "actual date"
// @Param pid query string false "parent id search pattern"
// @Param houseid query string false "house id search pattern"
// @Param exid query string false "except pu id"
// @Param active query boolean false "active pu"
// @Param ordering query string false "order by {objectname|streetname|startdate}"
// @Param desc query boolean false "descending order {true|false}"
// @Success 200 {object} models.Pu_count
// @Failure 500
// @Router /pu [get]
func HandlePu(w http.ResponseWriter, r *http.Request) {
	var gs ifPuService
	gs = services.NewPuService(pgsql.PuStorage{})
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
	gs1s, ok := query["objectname"]
	if ok && len(gs1s) > 0 {
		//case insensitive
		gs1 = strings.ToUpper(gs1s[0])
		//quotes
		re := regexp.MustCompile(`'`)
		gs1 = string(re.ReplaceAll([]byte(gs1), []byte("''")))
	}

	gs2 := ""
	gs2s, ok := query["streetname"]
	if ok && len(gs2s) > 0 {
		//case insensitive
		gs2 = strings.ToUpper(gs2s[0])
		//quotes
		re := regexp.MustCompile(`'`)
		gs2 = string(re.ReplaceAll([]byte(gs2), []byte("''")))
	}

	gs3 := time.Now().Format("2006-01-02")
	gs3s, ok := query["actualdate"]
	if ok && len(gs3s) > 0 {
		//case insensitive
		gs3 = strings.ToUpper(gs3s[0])
		//quotes
		re := regexp.MustCompile(`'`)
		gs3 = string(re.ReplaceAll([]byte(gs3), []byte("''")))
	}

	gs4 := ""
	gs4s, ok := query["pid"]
	if ok && len(gs4s) > 0 {
		_, err := strconv.Atoi(gs4s[0])
		if err == nil {
			gs4 = gs4s[0]
		}
	}

	gs5 := ""
	gs5s, ok := query["houseid"]
	if ok && len(gs5s) > 0 {
		_, err := strconv.Atoi(gs5s[0])
		if err == nil {
			gs5 = gs5s[0]
		}
	}

	gs6 := ""
	gs6s, ok := query["exid"]
	if ok && len(gs6s) > 0 {
		_, err := strconv.Atoi(gs6s[0])
		if err == nil {
			gs6 = gs6s[0]
		}
	}

	gs7 := ""
	gs7s, ok := query["active"]
	if ok && len(gs7s) > 0 {
		if gs7s[0] == "true" || gs7s[0] == "false" {
			gs7 = gs7s[0]
		} else {
			gs7 = ""
		}
	}

	ord := 1
	ords, ok := query["ordering"]
	if !ok || len(ords) == 0 {
		ord = 1
	} else if ords[0] == "objectname" {
		ord = 5
	} else if ords[0] == "streetname" {
		ord = 12
	} else if ords[0] == "startdate" {
		ord = 2
	}

	dsc := false
	dscs, ok := query["desc"]
	if ok && len(dscs) > 0 {
		if !(dscs[0] == "0") {
			dsc = true
		}
	}

	out_arr, err := gs.GetList(ctx, pg, pgs, gs1, gs2, gs3, gs4, gs5, gs6, gs7, ord, dsc)
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

// HandleAddPu godoc
// @Summary Add pu
// @Description add pu
// @Tags pu
// @Accept json
// @Produce  json
// @Param a body models.AddPu true "New pu. Significant params: Object.Id, PuObjectType, PuType.Id, PuNumber, InstallDate, CheckInterval, InitialValue, DevStopped, Startdate, Pid"
// @Success 200 {object} models.Json_id
// @Failure 500
// @Router /pu_add [post]
func HandleAddPu(w http.ResponseWriter, r *http.Request) {
	var gs ifPuService
	gs = services.NewPuService(pgsql.PuStorage{})
	ctx := context.Background()

	a := models.Pu{}
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
		log.Println("Failed execute ifPuService.Add: ", err)
	}

	output, err := json.Marshal(models.Json_id{Id: ai})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(output)

	return

}

// HandleUpdPu godoc
// @Summary Update pu
// @Description update pu
// @Tags pu
// @Accept json
// @Produce  json
// @Param u body models.Pu true "Update pu. Significant params: Id, Object.Id, PuObjectType, PuType.Id, PuNumber, InstallDate, CheckInterval, InitialValue, DevStopped, Startdate, Enddate, Pid"
// @Success 200 {object} models.Json_id
// @Failure 500
// @Router /pu_upd [post]
func HandleUpdPu(w http.ResponseWriter, r *http.Request) {
	var gs ifPuService
	gs = services.NewPuService(pgsql.PuStorage{})
	ctx := context.Background()

	u := models.Pu{}
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
		log.Println("Failed execute ifPuService.Upd: ", err)
	}

	output, err := json.Marshal(models.Json_id{Id: ui})

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(output)

	return
}

// HandleDelPu godoc
// @Summary Delete pu list
// @Description delete pu list
// @Tags pu
// @Accept json
// @Produce  json
// @Param d body models.Json_ids true "Delete pu list"
// @Success 200 {object} models.Json_ids
// @Failure 500
// @Router /pu_del [post]
func HandleDelPu(w http.ResponseWriter, r *http.Request) {
	var gs ifPuService
	gs = services.NewPuService(pgsql.PuStorage{})
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
		log.Println("Failed execute ifPuService.Del: ", err)
	}

	output, err := json.Marshal(models.Json_ids{Ids: res})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(output)

	return
}

// HandleGetPu godoc
// @Summary Get pu
// @Description get pu
// @Tags pu
// @Produce  json
// @Param id path int true "Pu by id"
// @Success 200 {object} models.Pu_count
// @Failure 500
// @Router /pu/{id} [get]
func HandleGetPu(w http.ResponseWriter, r *http.Request) {
	var gs ifPuService
	gs = services.NewPuService(pgsql.PuStorage{})
	ctx := context.Background()
	auth := utils.GetAuth(r)

	vars := mux.Vars(r)
	i, err := strconv.Atoi(vars["id"])
	if err != nil {
		i = 0
	}

	out_arr, err := gs.GetOne(ctx, i)
	if err != nil {
		log.Println("Failed execute ifPuService.GetOne: ", err)
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

// HandlePuObj godoc
// @Summary List pu of obj&tgu
// @Description get pu of obj&tgu list
// @Tags pu
// @Produce json
// @Param objid query string false "obj&tgu id"
// @Param tid query string false "obj&tgu type id (obj - type = 0, tgu - type > 0)"
// @Success 200 {object} models.Pu_count
// @Failure 500
// @Router /pu_obj [get]
func HandlePuObj(w http.ResponseWriter, r *http.Request) {
	var gs ifPuService
	gs = services.NewPuService(pgsql.PuStorage{})
	ctx := context.Background()
	auth := utils.GetAuth(r)

	query := r.URL.Query()

	gs1 := ""
	gs1s, ok := query["objid"]
	if ok && len(gs1s) > 0 {
		_, err := strconv.Atoi(gs1s[0])
		if err == nil {
			gs1 = gs1s[0]
		}
	}

	gs2 := ""
	gs2s, ok := query["tid"]
	if ok && len(gs2s) > 0 {
		_, err := strconv.Atoi(gs2s[0])
		if err == nil {
			gs2 = gs2s[0]
		}
	}

	out_arr, err := gs.GetObj(ctx, gs1, gs2)
	if err != nil {
		log.Println("Failed execute ifPuService.GetObj: ", err)
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
