package api

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v4"
	"github.com/u2lentaru/billbck/internal/models"
)

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
func (s *APG) HandleSubPu(w http.ResponseWriter, r *http.Request) {
	gs := models.Pu{}
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

	gs1 := 0
	gs1s, ok := query["parid"]
	if ok && len(gs1s) > 0 {
		t, err := strconv.Atoi(gs1s[0])
		if err == nil {
			gs1 = t
		}
	}

	gsc := 0
	err := s.Dbpool.QueryRow(ctx, "SELECT * from func_sub_pu_cnt($1);", gs1).Scan(&gsc)

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	out_arr := make([]models.Pu, 0,
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

	rows, err := s.Dbpool.Query(ctx, "SELECT * from func_sub_pu_get($1,$2,$3::int,$4,$5);", pg, pgs, gs1, ord, dsc)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&gs.Id, &gs.Startdate, &gs.Enddate, &gs.PuType.Id, &gs.PuType.PuTypeName, &gs.PuNumber, &gs.InstallDate,
			&gs.CheckInterval, &gs.InitialValue, &gs.DevStopped, &gs.Object.Id, &gs.PuObjectType, &gs.Object.ObjectName, &gs.Object.House.Id,
			&gs.Object.House.HouseNumber, &gs.Object.FlatNumber, &gs.Object.House.BuildingNumber, &gs.Object.RegQty,
			&gs.Object.House.Street.Id, &gs.Object.House.Street.StreetName, &gs.Object.House.Street.City.CityName,
			&gs.Object.House.BuildingType.BuildingTypeName, &gs.Object.House.Street.City.Id, &gs.Object.House.Street.Created,
			&gs.Object.House.Street.Closed, &gs.Pid)
		if err != nil {
			log.Println("failed to scan row:", err)
		}

		out_arr = append(out_arr, gs)
	}

	auth := models.Auth{Create: true, Read: true, Update: true, Delete: true}
	out_count, err := json.Marshal(models.Pu_count{Values: out_arr, Count: gsc, Auth: auth})
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
func (s *APG) HandleAddSubPu(w http.ResponseWriter, r *http.Request) {
	a := models.AddSubPu{}
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
	err = s.Dbpool.QueryRow(context.Background(), "SELECT func_sub_pu_add($1,$2);", a.ParId, a.SubId).Scan(&ai)

	if err != nil {
		log.Println("Failed execute func_sub_pu_add: ", err)
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
func (s *APG) HandleUpdSubPu(w http.ResponseWriter, r *http.Request) {
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

	ui := 0
	err = s.Dbpool.QueryRow(context.Background(), "SELECT func_sub_pu_upd($1,$2,$3);", u.Id, u.ParId, u.SubId).Scan(&ui)

	if err != nil {
		log.Println("Failed execute func_sub_pu_upd: ", err)
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
func (s *APG) HandleDelSubPu(w http.ResponseWriter, r *http.Request) {
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
		err = s.Dbpool.QueryRow(context.Background(), "SELECT func_sub_pu_del($1);", id).Scan(&i)
		res = append(res, i)

		if err != nil {
			log.Println("Failed execute func_sub_pu_del: ", err)
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

// HandleGetSubPu godoc
// @Summary Get subpu
// @Description get subpu
// @Tags subpu
// @Produce  json
// @Param id path int true "SubPu by id"
// @Success 200 {object} models.SubPu_count
// @Failure 500
// @Router /subpu/{id} [get]
func (s *APG) HandleGetSubPu(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	i := vars["id"]
	g := models.SubPu{}
	out_arr := []models.SubPu{}

	err := s.Dbpool.QueryRow(context.Background(), "SELECT * from func_sub_pu_getbyid($1);", i).Scan(&g.Id, &g.ParId, &g.SubId)

	if err != nil && err != pgx.ErrNoRows {
		log.Println("Failed execute from func_sub_pu_getbyid: ", err)
	}

	out_arr = append(out_arr, g)
	auth := models.Auth{Create: true, Read: true, Update: true, Delete: true}

	// output, err := json.Marshal(g)
	out_count, err := json.Marshal(models.SubPu_count{Values: out_arr, Count: 1, Auth: auth})
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
func (s *APG) HandlePrlSubPu(w http.ResponseWriter, r *http.Request) {
	gs := models.Pu{}
	ctx := context.Background()
	out_arr := []models.Pu{}

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

	rows, err := s.Dbpool.Query(ctx, "SELECT * from func_sub_pu_prl($1,$2,$3::int,$4::int,$5,$6);", pg, pgs, gs1, gs2, ord, dsc)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&gs.Id, &gs.Startdate, &gs.Enddate, &gs.PuType.Id, &gs.PuType.PuTypeName, &gs.PuNumber, &gs.InstallDate,
			&gs.CheckInterval, &gs.InitialValue, &gs.DevStopped, &gs.Object.Id, &gs.PuObjectType, &gs.Object.ObjectName, &gs.Object.House.Id,
			&gs.Object.House.HouseNumber, &gs.Object.FlatNumber, &gs.Object.House.BuildingNumber, &gs.Object.RegQty,
			&gs.Object.House.Street.Id, &gs.Object.House.Street.StreetName, &gs.Object.House.Street.City.CityName,
			&gs.Object.House.BuildingType.BuildingTypeName, &gs.Object.House.Street.City.Id, &gs.Object.House.Street.Created,
			&gs.Object.House.Street.Closed, &gs.Pid)
		if err != nil {
			log.Println("failed to scan row:", err)
		}

		out_arr = append(out_arr, gs)
	}

	gsc := 0
	err = s.Dbpool.QueryRow(ctx, "SELECT * from func_sub_pu_prl_cnt($1,$2);", gs1, gs2).Scan(&gsc)

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	out_count, err := json.Marshal(models.Pu_count{Values: out_arr, Count: gsc})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Write(out_count)

	return

}
