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
	"github.com/u2lentaru/billbck/internal/utils"
)

type ifSubPuService interface {
	GetList(ctx context.Context, pg, pgs, gs1, ord int, dsc bool) (models.Pu_count, error)
	Add(ctx context.Context, ea models.SubPu) (int, error)
	Upd(ctx context.Context, eu models.SubPu) (int, error)
	Del(ctx context.Context, ed []int) ([]int, error)
	GetOne(ctx context.Context, i int) (models.SubPu_count, error)
	GetPrl(ctx context.Context, pg, pgs, gs1, gs2, ord int, dsc bool) (models.Pu_count, error)
}

// HandleSubPu godoc
// @Summary List subpu
// @Description get subpu list
// @Tags subpu
// @Produce json
// @Param page query int false "page number"
// @Param page_size query int false "page size"
// @Param parid query int true "Subpu parid"
// @Param ordering query string false "order by {id|punamber}"
// @Param desc query boolean false "descending order {true|false}"
// @Success 200 {object} models.Pu_count
// @Failure 500
// @Router /subpu [get]
func HandleSubPu(w http.ResponseWriter, r *http.Request) {
	var gs ifSubPuService
	gs = services.NewSubPuService(pgsql.SubPuStorage{})
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
	gs1s, ok := query["parid"]
	if ok && len(gs1s) > 0 {
		t, err := strconv.Atoi(gs1s[0])
		if err == nil {
			gs1 = t
		}
	}

	ord := 1
	ords, ok := query["ordering"]
	if !ok || len(ords) == 0 {
		ord = 1
	} else if ords[0] == "id" {
		ord = 1
	} else if ords[0] == "punamber" {
		ord = 6
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

// HandleAddSubPu godoc
// @Summary Add subpu
// @Description add subpu
// @Tags subpu
// @Accept json
// @Produce  json
// @Param a body models.AddSubPu true "New subpu. Significant params: ParId, SubId"
// @Success 200 {object} models.Json_id
// @Failure 500
// @Router /subpu_add [post]
func HandleAddSubPu(w http.ResponseWriter, r *http.Request) {
	var gs ifSubPuService
	gs = services.NewSubPuService(pgsql.SubPuStorage{})
	ctx := context.Background()

	a := models.SubPu{}
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
		log.Println("Failed execute ifStaffService.Add: ", err)
	}

	output, err := json.Marshal(models.Json_id{Id: ai})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(output)

	return

}

// HandleUpdSubPu godoc
// @Summary Update subpu
// @Description update subpu
// @Tags subpu
// @Accept json
// @Produce  json
// @Param u body models.SubPu true "Update subpu. Significant params: Id, ParId, SubId"
// @Success 200 {object} models.Json_id
// @Failure 500
// @Router /subpu_upd [post]
func HandleUpdSubPu(w http.ResponseWriter, r *http.Request) {
	var gs ifSubPuService
	gs = services.NewSubPuService(pgsql.SubPuStorage{})
	ctx := context.Background()

	u := models.SubPu{}
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
		log.Println("Failed execute ifStaffService.Upd: ", err)
	}

	output, err := json.Marshal(models.Json_id{Id: ui})

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(output)

	return
}

// HandleDelSubPu godoc
// @Summary Delete subpu list
// @Description delete subpu list
// @Tags subpu
// @Accept json
// @Produce  json
// @Param d body models.Json_ids true "Delete subpu list"
// @Success 200 {object} models.Json_ids
// @Failure 500
// @Router /subpu_del [post]
func HandleDelSubPu(w http.ResponseWriter, r *http.Request) {
	var gs ifSubPuService
	gs = services.NewSubPuService(pgsql.SubPuStorage{})
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
		log.Println("Failed execute ifStaffService.Del: ", err)
	}

	output, err := json.Marshal(models.Json_ids{Ids: res})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(output)

	return
}

// HandleGetSubPu godoc
// @Summary Get subpu
// @Description get subpu
// @Tags subpu
// @Produce  json
// @Param id path int true "SubPu by id"
// @Success 200 {object} models.SubPu_count
// @Failure 500
// @Router /subpu/{id} [get]
func HandleGetSubPu(w http.ResponseWriter, r *http.Request) {
	var gs ifSubPuService
	gs = services.NewSubPuService(pgsql.SubPuStorage{})
	ctx := context.Background()
	auth := utils.GetAuth(r)

	vars := mux.Vars(r)
	i, err := strconv.Atoi(vars["id"])
	if err != nil {
		i = 0
	}

	out_arr, err := gs.GetOne(ctx, i)
	if err != nil {
		log.Println("Failed execute ifSubPuService.GetOne: ", err)
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

// HandlePrlSubPu godoc
// @Summary List subpu
// @Description get subpu list
// @Tags subpu
// @Produce json
// @Param page query int false "page number"
// @Param page_size query int false "page size"
// @Param houseid query int true "Subpu houseid"
// @Param subpuid query int true "Subpu id"
// @Param ordering query string false "order by {id|punamber}"
// @Param desc query boolean false "descending order {true|false}"
// @Success 200 {object} models.Pu_count
// @Failure 500
// @Router /subpu_prl [get]
func HandlePrlSubPu(w http.ResponseWriter, r *http.Request) {
	var gs ifSubPuService
	gs = services.NewSubPuService(pgsql.SubPuStorage{})
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
	gs1s, ok := query["houseid"]
	if ok && len(gs1s) > 0 {
		t, err := strconv.Atoi(gs1s[0])
		if err == nil {
			gs1 = t
		}
	}

	gs2 := 0
	gs2s, ok := query["subpuid"]
	if ok && len(gs2s) > 0 {
		t, err := strconv.Atoi(gs2s[0])
		if err == nil {
			gs2 = t
		}
	}

	ord := 1
	ords, ok := query["ordering"]
	if !ok || len(ords) == 0 {
		ord = 1
	} else if ords[0] == "id" {
		ord = 1
	} else if ords[0] == "punamber" {
		ord = 6
	}

	dsc := false
	dscs, ok := query["desc"]
	if ok && len(dscs) > 0 {
		if !(dscs[0] == "0") {
			dsc = true
		}
	}

	out_arr, err := gs.GetPrl(ctx, pg, pgs, gs1, gs2, ord, dsc)
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
