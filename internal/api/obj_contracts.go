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
	"github.com/jackc/pgx/v4"
	"github.com/u2lentaru/billbck/internal/models"
	"github.com/u2lentaru/billbck/internal/utils"
)

// HandleObjContracts godoc
// @Summary List objcontracts
// @Description get objcontract list
// @Tags objcontracts
// @Produce  json
// @Param page query int false "page number"
// @Param page_size query int false "page size"
// @Param objectid query string false "object id"
// @Param tid query string false "type id"
// @Param contractid query string false "contract id"
// @Param active query boolean false "enddate is null"
// @Param ordering query string false "order by {object|contract|startdate}"
// @Param desc query boolean false "descending order {true|false}"
// @Success 200 {object} models.ObjContract_count
// @Failure 500
// @Router /objcontracts [get]
func (s *APG) HandleObjContracts(w http.ResponseWriter, r *http.Request) {
	gs := models.ObjContract{}
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
	gs1s, ok := query["objectid"]
	if ok && len(gs1s) > 0 {
		if gs1t, err := strconv.Atoi(gs1s[0]); err != nil {
			gs1 = 0
		} else {
			gs1 = gs1t
		}
	}

	gs2 := 0
	gs2s, ok := query["tid"]
	if ok && len(gs2s) > 0 {
		if gs2t, err := strconv.Atoi(gs2s[0]); err != nil {
			gs2 = 0
		} else {
			gs1 = gs2t
		}
	}

	gs3 := 0
	gs3s, ok := query["contractid"]
	if ok && len(gs3s) > 0 {
		if gs3t, err := strconv.Atoi(gs3s[0]); err != nil {
			gs3 = 0
		} else {
			gs3 = gs3t
		}
	}

	gs4 := false
	gs4f := true
	gs4s, ok := query["active"]
	if ok && len(gs4s) > 0 {
		if gs4s[0] == "true" || gs4s[0] == "false" {
			gs4, _ = strconv.ParseBool(gs4s[0])
		} else {
			gs4f = false
		}
	} else {
		gs4f = false
	}

	gsc := 0
	err := s.Dbpool.QueryRow(ctx, "SELECT * from func_obj_contracts_cnt($1,$2,$3,$4);", utils.NullableInt(int32(gs1)),
		utils.NullableInt(int32(gs2)), utils.NullableInt(int32(gs3)), utils.NullableBool(gs4, gs4f)).Scan(&gsc)

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	out_arr := make([]models.ObjContract, 0,
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
	} else if ords[0] == "object" {
		ord = 6
	} else if ords[0] == "contract" {
		ord = 17
	} else if ords[0] == "startdate" {
		ord = 4
	}

	dsc := false
	dscs, ok := query["desc"]
	if ok && len(dscs) > 0 {
		if !(dscs[0] == "0") {
			dsc = true
		}
	}

	rows, err := s.Dbpool.Query(ctx, "SELECT * from func_obj_contracts_get($1,$2,$3,$4,$5,$6,$7,$8);", pg, pgs, utils.NullableInt(int32(gs1)),
		utils.NullableInt(int32(gs2)), utils.NullableInt(int32(gs3)), utils.NullableBool(gs4, gs4f), ord, dsc)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&gs.Id, &gs.Contract.Id, &gs.Object.Id, &gs.ObjTypeId, &gs.Startdate, &gs.Enddate, &gs.Object.ObjectName, &gs.Object.RegQty,
			&gs.Object.FlatNumber, &gs.Object.House.Id, &gs.Object.House.HouseNumber, &gs.Object.House.BuildingNumber,
			&gs.Object.House.Street.City.CityName, &gs.Object.House.Street.Id, &gs.Object.House.Street.StreetName, &gs.Object.TariffGroup.Id,
			&gs.Object.TariffGroup.TariffGroupName, &gs.Contract.ContractNumber, &gs.Contract.Startdate, &gs.Contract.Enddate,
			&gs.Contract.Customer.SubId, &gs.Contract.Customer.SubName, &gs.Contract.Consignee.SubId, &gs.Contract.Consignee.SubName)
		if err != nil {
			log.Println("failed to scan row:", err)
		}

		out_arr = append(out_arr, gs)
	}

	auth := models.Auth{Create: true, Read: true, Update: true, Delete: true}
	out_count, err := json.Marshal(models.ObjContract_count{Values: out_arr, Count: gsc, Auth: auth})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Write(out_count)

	return
}

// HandleAddObjContract godoc
// @Summary Add objcontract
// @Description add objcontract
// @Tags objcontracts
// @Accept json
// @Produce  json
// @Param a body models.AddObjContract true "New objcontract. Old objcontract of the object will be closed. Significant params: Contract.Id, Object.Id, ObjTypeId, Startdate"
// @Success 200 {object} models.Json_id
// @Failure 500
// @Router /objcontracts_add [post]
func (s *APG) HandleAddObjContract(w http.ResponseWriter, r *http.Request) {
	a := models.AddObjContract{}
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
	err = s.Dbpool.QueryRow(context.Background(), "SELECT func_obj_contracts_add($1,$2,$3,$4);",
		a.Contract.Id, a.Object.Id, a.ObjTypeId, a.Startdate).Scan(&ai)

	if err != nil {
		log.Println("Failed execute func_obj_contracts_add: ", err)
	}

	output, err := json.Marshal(models.Json_id{Id: ai})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(output)

	return

}

// HandleUpdObjContract godoc
// @Summary Update objcontract
// @Description update objcontract
// @Tags objcontracts
// @Accept json
// @Produce  json
// @Param u body models.ObjContract true "Update objcontract. Significant params: Contract.Id, Object.Id, ObjTypeId, Startdate, Enddate"
// @Success 200 {object} models.Json_id
// @Failure 500
// @Router /objcontracts_upd [post]
func (s *APG) HandleUpdObjContract(w http.ResponseWriter, r *http.Request) {
	u := models.ObjContract{}
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
	err = s.Dbpool.QueryRow(context.Background(), "SELECT func_obj_contracts_upd($1,$2,$3,$4,$5,$6);", u.Id,
		u.Contract.Id, u.Object.Id, u.ObjTypeId, u.Startdate, u.Enddate).Scan(&ui)

	if err != nil {
		log.Println("Failed execute func_obj_contracts_upd: ", err)
	}

	output, err := json.Marshal(models.Json_id{Id: ui})

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(output)

	return
}

// HandleDelObjContract godoc
// @Summary Close objcontract
// @Description close objcontract
// @Tags objcontracts
// @Accept json
// @Produce  json
// @Param d body models.IdClose true "Close objcontract. Significant params: Id, CloseDate"
// @Success 200 {object} models.Json_id
// @Failure 500
// @Router /objcontracts_del [post]
func (s *APG) HandleDelObjContract(w http.ResponseWriter, r *http.Request) {
	d := models.IdClose{}
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

	// res := []int{}
	i := 0
	_, err = time.Parse("2006-01-02", d.CloseDate)
	if err != nil {
		d.CloseDate = time.Now().Format("2006-01-02")
	}

	// for _, id := range d.Ids {
	err = s.Dbpool.QueryRow(context.Background(), "SELECT func_obj_contracts_del($1,$2);", d.Id, d.CloseDate).Scan(&i)
	// res = append(res, i)

	if err != nil {
		log.Println("Failed execute func_obj_contracts_del: ", err)
	}
	// }

	output, err := json.Marshal(models.Json_id{Id: i})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(output)

	return
}

// HandleGetObjContract godoc
// @Summary Get objcontract
// @Description get objcontract
// @Tags objcontracts
// @Produce  json
// @Param id path int true "ObjContract by id"
// @Param actualdate query string false "actual date"
// @Success 200 {object} models.ObjContract_count
// @Failure 500
// @Router /objcontracts/{id} [get]
func (s *APG) HandleGetObjContract(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	i := vars["id"]
	g := models.ObjContract{}
	out_arr := []models.ObjContract{}

	query := r.URL.Query()

	gs3 := time.Now().Format("2006-01-02")
	// log.Println(gs3)
	gs3s, ok := query["actualdate"]
	if ok && len(gs3s) > 0 {
		//case insensitive
		gs3 = strings.ToUpper(gs3s[0])
		//quotes
		re := regexp.MustCompile(`'`)
		gs3 = string(re.ReplaceAll([]byte(gs3), []byte("''")))
	}

	err := s.Dbpool.QueryRow(context.Background(), "SELECT * from func_obj_contract_get($1,$2);", i, gs3).Scan(&g.Id, &g.Contract.Id,
		&g.Object.Id, &g.ObjTypeId, &g.Startdate, &g.Enddate, &g.Object.ObjectName, &g.Object.RegQty, &g.Object.FlatNumber, &g.Object.House.Id,
		&g.Object.House.HouseNumber, &g.Object.House.BuildingNumber, &g.Object.House.Street.City.CityName, &g.Object.House.Street.Id,
		&g.Object.House.Street.StreetName, &g.Object.TariffGroup.Id, &g.Object.TariffGroup.TariffGroupName, &g.Contract.ContractNumber,
		&g.Contract.Startdate, &g.Contract.Enddate, &g.Contract.Customer.SubId, &g.Contract.Customer.SubName, &g.Contract.Consignee.SubId,
		&g.Contract.Consignee.SubName)

	if err != nil && err != pgx.ErrNoRows {
		log.Println("Failed execute from func_obj_contract_get: ", err)
	}

	out_arr = append(out_arr, g)
	auth := models.Auth{Create: true, Read: true, Update: true, Delete: true}

	// output, err := json.Marshal(g)
	out_count, err := json.Marshal(models.ObjContract_count{Values: out_arr, Count: 1, Auth: auth})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(out_count)

	return
}
