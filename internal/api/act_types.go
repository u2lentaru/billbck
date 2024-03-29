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

type ifActTypeService interface {
	GetList(ctx context.Context, pg, pgs int, nm string, ord int, dsc bool) (models.ActType_count, error)
	Add(ctx context.Context, ea models.ActType) (int, error)
	Upd(ctx context.Context, eu models.ActType) (int, error)
	Del(ctx context.Context, ed []int) ([]int, error)
	GetOne(ctx context.Context, i int) (models.ActType_count, error)
}

// HandleActTypes godoc
// @Summary List acttypes
// @Description get acttype list
// @Tags acttypes
// @Produce  json
// @Param page query int false "page number"
// @Param page_size query int false "page size"
// @Param acttypename query string false "acttypename search pattern"
// @Param ordering query string false "order by {id|acttypename}"
// @Param desc query boolean false "descending order {true|false}"
// @Success 200 {object} models.ActType_count
// @Failure 500
// @Router /acttypes [get]
func HandleActTypes(w http.ResponseWriter, r *http.Request) {
	var gs ifActTypeService
	gs = services.NewActTypeService(pgsql.ActTypeStorage{})
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
	gs1s, ok := query["acttypename"]
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
	} else if ords[0] == "acttypename" {
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

// HandleAddActType godoc
// @Summary Add acttype
// @Description add acttype
// @Tags acttypes
// @Accept json
// @Produce  json
// @Param a body models.AddActType true "New acttype"
// @Success 200 {object} models.Json_id
// @Failure 500
// @Router /acttypes_add [post]
func HandleAddActType(w http.ResponseWriter, r *http.Request) {
	var gs ifActTypeService
	gs = services.NewActTypeService(pgsql.ActTypeStorage{})
	ctx := context.Background()

	a := models.ActType{}
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
		log.Println("Failed execute ifActTypeService.Add: ", err)
	}

	output, err := json.Marshal(models.Json_id{Id: ai})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(output)

	return

}

// HandleUpdActType godoc
// @Summary Update acttype
// @Description update acttype
// @Tags acttypes
// @Accept json
// @Produce  json
// @Param u body models.ActType true "Update acttype"
// @Success 200 {object} models.Json_id
// @Failure 500
// @Router /acttypes_upd [post]
func HandleUpdActType(w http.ResponseWriter, r *http.Request) {
	var gs ifActTypeService
	gs = services.NewActTypeService(pgsql.ActTypeStorage{})
	ctx := context.Background()

	u := models.ActType{}
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
		log.Println("Failed execute ifActTypeService.Upd: ", err)
	}

	output, err := json.Marshal(models.Json_id{Id: ui})

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(output)

	return
}

// HandleDelActType godoc
// @Summary Delete acttypes
// @Description delete acttypes
// @Tags acttypes
// @Accept json
// @Produce  json
// @Param d body models.Json_ids true "Delete acttypes"
// @Success 200 {object} models.Json_ids
// @Failure 500
// @Router /acttypes_del [post]
func HandleDelActType(w http.ResponseWriter, r *http.Request) {
	var gs ifActTypeService
	gs = services.NewActTypeService(pgsql.ActTypeStorage{})
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
		log.Println("Failed execute ifActTypeService.Del: ", err)
	}

	output, err := json.Marshal(models.Json_ids{Ids: res})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(output)

	return

}

// HandleGetActType godoc
// @Summary Get acttype
// @Description get acttype
// @Tags acttypes
// @Produce  json
// @Param id path int true "ActType by id"
// @Success 200 {object} models.ActType_count
// @Failure 500
// @Router /acttypes/{id} [get]
func HandleGetActType(w http.ResponseWriter, r *http.Request) {
	var gs ifActTypeService
	gs = services.NewActTypeService(pgsql.ActTypeStorage{})
	ctx := context.Background()
	auth := utils.GetAuth(r)

	vars := mux.Vars(r)
	i, err := strconv.Atoi(vars["id"])
	if err != nil {
		i = 0
	}

	out_arr, err := gs.GetOne(ctx, i)
	if err != nil {
		log.Println("Failed execute ifActTypeService.GetOne: ", err)
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

/*
// HandleActTypes godoc
// @Summary List acttypes
// @Description get acttype list
// @Tags acttypes
// @Produce  json
// @Param page query int false "page number"
// @Param page_size query int false "page size"
// @Param acttypename query string false "acttypename search pattern"
// @Param ordering query string false "order by {id|acttypename}"
// @Param desc query boolean false "descending order {true|false}"
// @Success 200 {object} models.ActType_count
// @Failure 500
// @Router /acttypes [get]
func (s *APG) HandleActTypes(w http.ResponseWriter, r *http.Request) {
	gs := models.ActType{}
	ctx := context.Background()

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
	gs1s, ok := query["acttypename"]
	if ok && len(gs1s) > 0 {
		//case insensitive
		gs1 = strings.ToUpper(gs1s[0])
		//quotes
		re := regexp.MustCompile(`'`)
		gs1 = string(re.ReplaceAll([]byte(gs1), []byte("''")))
	}

	gsc := 0
	err := s.Dbpool.QueryRow(ctx, "SELECT * from func_act_types_cnt($1);", gs1).Scan(&gsc)

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	out_arr := make([]models.ActType, 0,
		func() int {
			if gsc < pgs {
				return gsc
			} else {
				return pgs
			}
		}())

	ord := 1
	ords, ok := query["ordering"]
	if !ok || len(ords) == 0 {
		ord = 1
	} else if ords[0] == "acttypename" {
		ord = 2
	}

	dsc := false
	dscs, ok := query["desc"]
	if ok && len(dscs) > 0 {
		if !(dscs[0] == "0") {
			dsc = true
		}
	}

	rows, err := s.Dbpool.Query(ctx, "SELECT * from func_act_types_get($1,$2,$3,$4,$5);", pg, pgs, gs1, ord, dsc)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&gs.Id, &gs.ActTypeName)
		if err != nil {
			log.Println("failed to scan row:", err)
		}

		out_arr = append(out_arr, gs)
	}

	auth := models.Auth{Create: true, Read: true, Update: true, Delete: true}

	out_count, err := json.Marshal(models.ActType_count{Values: out_arr, Count: gsc, Auth: auth})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Write(out_count)

	return

}

// HandleAddActType godoc
// @Summary Add acttype
// @Description add acttype
// @Tags acttypes
// @Accept json
// @Produce  json
// @Param a body models.AddActType true "New acttype"
// @Success 200 {object} models.Json_id
// @Failure 500
// @Router /acttypes_add [post]
func (s *APG) HandleAddActType(w http.ResponseWriter, r *http.Request) {
	a := models.AddActType{}
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

	ai := 0
	err = s.Dbpool.QueryRow(context.Background(), "SELECT func_act_types_add($1);", a.ActTypeName).Scan(&ai)

	if err != nil {
		log.Println("Failed execute func_act_types_add: ", err)
	}

	output, err := json.Marshal(models.Json_id{Id: ai})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(output)

	return

}

// HandleUpdActType godoc
// @Summary Update acttype
// @Description update acttype
// @Tags acttypes
// @Accept json
// @Produce  json
// @Param u body models.ActType true "Update acttype"
// @Success 200 {object} models.Json_id
// @Failure 500
// @Router /acttypes_upd [post]
func (s *APG) HandleUpdActType(w http.ResponseWriter, r *http.Request) {
	u := models.ActType{}
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

	ui := 0
	err = s.Dbpool.QueryRow(context.Background(), "SELECT func_act_types_upd($1,$2);", u.Id, u.ActTypeName).Scan(&ui)

	if err != nil {
		log.Println("Failed execute func_act_types_upd: ", err)
	}

	output, err := json.Marshal(models.Json_id{Id: ui})

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(output)

	return
}

// HandleDelActType godoc
// @Summary Delete acttypes
// @Description delete acttypes
// @Tags acttypes
// @Accept json
// @Produce  json
// @Param d body models.Json_ids true "Delete acttypes"
// @Success 200 {object} models.Json_ids
// @Failure 500
// @Router /acttypes_del [post]
func (s *APG) HandleDelActType(w http.ResponseWriter, r *http.Request) {
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

	res := []int{}
	i := 0
	for _, id := range d.Ids {
		err = s.Dbpool.QueryRow(context.Background(), "SELECT func_act_types_del($1);", id).Scan(&i)
		res = append(res, i)

		if err != nil {
			log.Println("Failed execute func_act_types_del: ", err)
		}
	}

	output, err := json.Marshal(models.Json_ids{Ids: res})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(output)

	return

}

// HandleGetActType godoc
// @Summary Get acttype
// @Description get acttype
// @Tags acttypes
// @Produce  json
// @Param id path int true "ActType by id"
// @Success 200 {object} models.ActType_count
// @Failure 500
// @Router /acttypes/{id} [get]
func (s *APG) HandleGetActType(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	i := vars["id"]
	g := models.ActType{}

	out_arr := []models.ActType{}

	err := s.Dbpool.QueryRow(context.Background(), "SELECT * from func_act_type_get($1);", i).Scan(&g.Id, &g.ActTypeName)

	if err != nil && err != pgx.ErrNoRows {
		log.Println("Failed execute from func_act_type_get: ", err)
	}

	out_arr = append(out_arr, g)
	auth := models.Auth{Create: true, Read: true, Update: true, Delete: true}

	// output, err := json.Marshal(g)
	out_count, err := json.Marshal(models.ActType_count{Values: out_arr, Count: 1, Auth: auth})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	// w.Write(output)
	w.Write(out_count)

	return
}
*/
