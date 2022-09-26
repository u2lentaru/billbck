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

type ifEquipmentService interface {
	GetList(ctx context.Context, pg, pgs, gs1 int, gs2 string, ord int, dsc bool) (models.Equipment_count, error)
	Add(ctx context.Context, ea models.Equipment) (int, error)
	Upd(ctx context.Context, eu models.Equipment) (int, error)
	Del(ctx context.Context, ed []int) ([]int, error)
	GetOne(ctx context.Context, i int) (models.Equipment_count, error)
	AddList(ctx context.Context, al models.Equipment_count) ([]int, error)
	DelObj(ctx context.Context, d []int) ([]int, error)
}

// HandleEquipment godoc
// @Summary List equipment
// @Description get equipment list
// @Tags equipment
// @Produce  json
// @Param page query int false "page number"
// @Param page_size query int false "page size"
// @Param objectid query int false "objectid"
// @Param objectname query string false "name search pattern"
// @Param ordering query string false "order by {id|objectname}"
// @Param desc query boolean false "descending order {true|false}"
// @Success 200 {object} models.Equipment_count
// @Failure 500
// @Router /equipment [get]
func HandleEquipment(w http.ResponseWriter, r *http.Request) {
	var gs ifEquipmentService
	gs = services.NewEquipmentService(pgsql.EquipmentStorage{})
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

	gs1 := 0
	gs1s, ok := query["objectid"]
	if ok && len(gs1s) > 0 {
		t, err := strconv.Atoi(gs1s[0])
		if err == nil {
			gs1 = t
		}
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

	ord := 1
	ords, ok := query["ordering"]
	if !ok || len(ords) == 0 {
		ord = 1
	} else if ords[0] == "objectname" {
		ord = 8
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

// HandleAddEquipment godoc
// @Summary Add equipment
// @Description add equipment
// @Tags equipment
// @Accept json
// @Produce  json
// @Param a body models.AddEquipment true "New equipment. Significant params: EquipmentType.Id, Object.Id, Qty, WorkingHours"
// @Success 200 {object} models.Json_id
// @Failure 500
// @Router /equipment_add [post]
func HandleAddEquipment(w http.ResponseWriter, r *http.Request) {
	var gs ifEquipmentService
	gs = services.NewEquipmentService(pgsql.EquipmentStorage{})
	ctx := context.Background()

	a := models.Equipment{}
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
		log.Println("Failed execute ifEquipmentService.Add: ", err)
	}

	output, err := json.Marshal(models.Json_id{Id: ai})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(output)

	return
}

// HandleUpdEquipment godoc
// @Summary Update equipment
// @Description update equipment
// @Tags equipment
// @Accept json
// @Produce  json
// @Param u body models.Equipment true "Update equipment. Significant params: Id, EquipmentType.Id, Object.Id, Qty, WorkingHours"
// @Success 200 {object} models.Json_id
// @Failure 500
// @Router /equipment_upd [post]
func HandleUpdEquipment(w http.ResponseWriter, r *http.Request) {
	var gs ifEquipmentService
	gs = services.NewEquipmentService(pgsql.EquipmentStorage{})
	ctx := context.Background()

	u := models.Equipment{}
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
		log.Println("Failed execute ifEquipmentService.Upd: ", err)
	}

	output, err := json.Marshal(models.Json_id{Id: ui})

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(output)

	return
}

// HandleDelEquipment godoc
// @Summary Delete equipment
// @Description delete equipment
// @Tags equipment
// @Accept json
// @Produce  json
// @Param d body models.Json_ids true "Delete equipment"
// @Success 200 {object} models.Json_ids
// @Failure 500
// @Router /equipment_del [post]
func HandleDelEquipment(w http.ResponseWriter, r *http.Request) {
	var gs ifEquipmentService
	gs = services.NewEquipmentService(pgsql.EquipmentStorage{})
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
		log.Println("Failed execute ifEquipmentService.Del: ", err)
	}

	output, err := json.Marshal(models.Json_ids{Ids: res})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(output)

	return

}

// HandleGetEquipment godoc
// @Summary Get equipment
// @Description get equipment
// @Tags equipment
// @Produce  json
// @Param id path int true "Equipment by id"
// @Success 200 {object} models.Equipment_count
// @Failure 500
// @Router /equipment/{id} [get]
func HandleGetEquipment(w http.ResponseWriter, r *http.Request) {
	var gs ifEquipmentService
	gs = services.NewEquipmentService(pgsql.EquipmentStorage{})
	ctx := context.Background()
	auth := utils.GetAuth(r)

	vars := mux.Vars(r)
	i, err := strconv.Atoi(vars["id"])
	if err != nil {
		i = 0
	}

	out_arr, err := gs.GetOne(ctx, i)
	if err != nil {
		log.Println("Failed execute ifEquipmentService.GetOne: ", err)
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

// HandleAddEquipmentList godoc
// @Summary Add equipment list
// @Description add equipment list
// @Tags equipment
// @Accept json
// @Produce  json
// @Param al body models.Equipment_count true "Add equipment list. Old equipment delete by first value Object.Id. Significant params: EquipmentType.Id, Object.Id, Qty, WorkingHours"
// @Success 200 {object} models.Json_ids
// @Failure 500
// @Router /equipment_addlist [post]
func HandleAddEquipmentList(w http.ResponseWriter, r *http.Request) {
	var gs ifEquipmentService
	gs = services.NewEquipmentService(pgsql.EquipmentStorage{})
	ctx := context.Background()

	al := models.Equipment_count{}
	body, err := ioutil.ReadAll(r.Body)

	defer r.Body.Close()

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	err = json.Unmarshal(body, &al)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	res, err := gs.AddList(ctx, al)

	if err != nil {
		log.Println("Failed execute ifEquipmentService.AddList: ", err)
	}

	output, err := json.Marshal(models.Json_ids{Ids: res})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(output)

	return
}

// HandleDelObjEquipment godoc
// @Summary Delete equipment by object id
// @Description delete equipment by object id
// @Tags equipment
// @Accept json
// @Produce  json
// @Param d body models.Json_ids true "Delete equipment by object id"
// @Success 200 {object} models.Json_ids
// @Failure 500
// @Router /equipment_delobj [post]
func HandleDelObjEquipment(w http.ResponseWriter, r *http.Request) {
	var gs ifEquipmentService
	gs = services.NewEquipmentService(pgsql.EquipmentStorage{})
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

	res, err := gs.DelObj(ctx, d.Ids)
	if err != nil {
		log.Println("Failed execute ifEquipmentService.DelObj: ", err)
	}

	output, err := json.Marshal(models.Json_ids{Ids: res})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(output)

	return

}
