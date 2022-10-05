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

type ifObjectService interface {
	GetList(ctx context.Context, pg, pgs int, gs1, gs2 string, gs3, gs3f bool, ord int, dsc bool) (models.Object_count, error)
	Add(ctx context.Context, ea models.Object) (int, error)
	Upd(ctx context.Context, eu models.Object) (int, error)
	Del(ctx context.Context, ed []int) ([]int, error)
	GetOne(ctx context.Context, i int) (models.Object_count, error)
	GetObjContract(ctx context.Context, i int, a string) (models.ObjContract_count, error)
	GetMff(ctx context.Context, i int) (models.Object_count, error)
}

// HandleObjects godoc
// @Summary List objects
// @Description get object list
// @Tags objects
// @Produce  json
// @Param page query int false "page number"
// @Param page_size query int false "page size"
// @Param objectname query string false "objectname search pattern"
// @Param houseid query string false "house id search pattern"
// @Param active query boolean false "active contract"
// @Param ordering query string false "order by {id|objectname}"
// @Param desc query boolean false "descending order {true|false}"
// @Success 200 {object} models.Object_count
// @Failure 500
// @Router /objects [get]
func HandleObjects(w http.ResponseWriter, r *http.Request) {
	var gs ifObjectService
	gs = services.NewObjectService(pgsql.ObjectStorage{})
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
	gs2s, ok := query["houseid"]
	if ok && len(gs2s) > 0 {
		_, err := strconv.Atoi(gs2s[0])
		if err == nil {
			gs2 = gs2s[0]
		}
	}

	gs3 := false
	gs3f := true
	gs3s, ok := query["active"]
	if ok && len(gs3s) > 0 {
		if gs3s[0] == "true" || gs3s[0] == "false" {
			gs3, _ = strconv.ParseBool(gs3s[0])
		} else {
			gs3f = false
		}
	} else {
		gs3f = false
	}

	ord := 1
	ords, ok := query["ordering"]
	if !ok || len(ords) == 0 {
		ord = 1
	} else if ords[0] == "objectname" {
		ord = 2
	}

	dsc := false
	dscs, ok := query["desc"]
	if ok && len(dscs) > 0 {
		if !(dscs[0] == "0") {
			dsc = true
		}
	}

	out_arr, err := gs.GetList(ctx, pg, pgs, gs1, gs2, gs3, gs3f, ord, dsc)
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

	// log.Printf("Work time %s", time.Since(start))

	return

}

// HandleAddObject godoc
// @Summary Add object
// @Description add object
// @Tags objects
// @Accept json
// @Produce  json
// @Param a body models.AddObject true "New object. Significant params: ObjectName, House.Id, FlatNumber(n), ObjType.Id, RegQty, Uzo.Id, TariffGroup.Id, CalculationType.Id, ObjStatus.Id, Notes(n), MffId(n)"
// @Success 200 {object} models.Json_id
// @Failure 500
// @Router /objects_add [post]
func HandleAddObject(w http.ResponseWriter, r *http.Request) {
	var gs ifObjectService
	gs = services.NewObjectService(pgsql.ObjectStorage{})
	ctx := context.Background()

	a := models.Object{}
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
		log.Println("Failed execute ifObjectService.Add: ", err)
	}

	output, err := json.Marshal(models.Json_id{Id: ai})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(output)

	return

}

// HandleUpdObject godoc
// @Summary Update object
// @Description update object
// @Tags objects
// @Accept json
// @Produce  json
// @Param u body models.Object true "Update object. Significant params: Id, ObjectName, House.Id, FlatNumber(n), ObjType.Id, RegQty, Uzo.Id, TariffGroup.Id, CalculationType.Id, ObjStatus.Id, Notes(n), MffId(n)"
// @Success 200 {object} models.Json_id
// @Failure 500
// @Router /objects_upd [post]
func HandleUpdObject(w http.ResponseWriter, r *http.Request) {
	var gs ifObjectService
	gs = services.NewObjectService(pgsql.ObjectStorage{})
	ctx := context.Background()

	u := models.Object{}
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
		log.Println("Failed execute ifObjectService.Upd: ", err)
	}

	output, err := json.Marshal(models.Json_id{Id: ui})

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(output)

	return
}

// HandleDelObject godoc
// @Summary Delete objects
// @Description delete objects
// @Tags objects
// @Accept json
// @Produce  json
// @Param d body models.Json_ids true "Delete objects"
// @Success 200 {object} models.Json_ids
// @Failure 500
// @Router /objects_del [post]
func HandleDelObject(w http.ResponseWriter, r *http.Request) {
	var gs ifObjectService
	gs = services.NewObjectService(pgsql.ObjectStorage{})
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
		log.Println("Failed execute ifObjectService.Del: ", err)
	}

	output, err := json.Marshal(models.Json_ids{Ids: res})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(output)

	return

}

// HandleGetObject godoc
// @Summary Get object
// @Description get object
// @Tags objects
// @Produce  json
// @Param id path int true "Object by id"
// @Success 200 {object} models.Object_count
// @Failure 500
// @Router /objects/{id} [get]
func HandleGetObject(w http.ResponseWriter, r *http.Request) {
	var gs ifObjectService
	gs = services.NewObjectService(pgsql.ObjectStorage{})
	ctx := context.Background()
	auth := utils.GetAuth(r)

	vars := mux.Vars(r)
	i, err := strconv.Atoi(vars["id"])
	if err != nil {
		i = 0
	}

	out_arr, err := gs.GetOne(ctx, i)
	if err != nil {
		log.Println("Failed execute ifObjectService.GetOne: ", err)
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

// HandleGetObjectContract godoc
// @Summary Get objects contract
// @Description get objects contract
// @Tags objects
// @Produce  json
// @Param id path int true "Object by id"
// @Param active query boolean false "active contracts"
// @Success 200 {object} models.ObjContract
// @Failure 500
// @Router /objects_getcontract/{id} [get]
func HandleGetObjectContract(w http.ResponseWriter, r *http.Request) {
	var gs ifObjectService
	gs = services.NewObjectService(pgsql.ObjectStorage{})
	ctx := context.Background()
	auth := utils.GetAuth(r)

	vars := mux.Vars(r)
	i, err := strconv.Atoi(vars["id"])
	if err != nil {
		i = 0
	}
	out_arr := models.ObjContract_count{}

	query := r.URL.Query()

	dsc := ""
	dscs, ok := query["active"]
	if ok && len(dscs) > 0 {
		if dscs[0] == "true" || dscs[0] == "false" {
			dsc = dscs[0]
		} else {
			dsc = ""
		}
	}

	out_arr, err = gs.GetObjContract(ctx, i, dsc)
	if err != nil {
		log.Println("Failed execute ifObjectService.GetObjContract: ", err)
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

// HandleGetObjectMff godoc
// @Summary Get main mff object by house id
// @Description get main mff object by house id
// @Tags objects
// @Produce  json
// @Param hid path int true "House id"
// @Success 200 {object} models.Object_count
// @Failure 500
// @Router /objects_mff/{hid} [get]
func HandleGetObjectMff(w http.ResponseWriter, r *http.Request) {
	var gs ifObjectService
	gs = services.NewObjectService(pgsql.ObjectStorage{})
	ctx := context.Background()
	auth := utils.GetAuth(r)

	vars := mux.Vars(r)
	i, err := strconv.Atoi(vars["hid"])
	if err != nil {
		i = 0
	}

	out_arr, err := gs.GetMff(ctx, i)
	if err != nil {
		log.Println("Failed execute ifObjectService.GetMff: ", err)
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
