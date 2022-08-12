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
	"github.com/jackc/pgx/v4"
	"github.com/u2lentaru/billbck/internal/models"
	"github.com/u2lentaru/billbck/internal/utils"
)

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
func (s *APG) HandleObjects(w http.ResponseWriter, r *http.Request) {
	// start := time.Now()
	gs := models.Object{}
	ctx := context.Background()
	// out_arr := []models.Object{}

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

	out_arr := make([]models.Object, 0, pgs)

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

	rows, err := s.Dbpool.Query(ctx, "SELECT * from func_objects_get($1,$2,$3,$4,$5,$6,$7);", pg, pgs, gs1, NullableString(gs2),
		utils.NullableBool(gs3, gs3f), ord, dsc)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&gs.Id, &gs.ObjectName, &gs.House.Id, &gs.FlatNumber, &gs.ObjType.Id, &gs.RegQty, &gs.Uzo.Id, &gs.TariffGroup.Id,
			&gs.Notes, &gs.CalculationType.Id, &gs.ObjStatus.Id, &gs.MffId, &gs.House.BuildingType.Id, &gs.House.Street.Id, &gs.House.HouseNumber,
			&gs.House.BuildingNumber, &gs.House.RP.Id, &gs.House.Area.Id, &gs.House.Ksk.Id, &gs.House.Sector.Id, &gs.House.Connector.Id,
			&gs.House.InputType.Id, &gs.House.Reliability.Id, &gs.House.Voltage.Id, &gs.House.BuildingType.BuildingTypeName,
			&gs.House.Street.StreetName, &gs.House.Street.Created, &gs.House.Street.City.CityName, &gs.House.RP.RpName, &gs.House.Area.AreaName,
			&gs.House.Area.AreaNumber, &gs.House.Ksk.KskName, &gs.House.Sector.SectorName, &gs.House.Connector.ConnectorName,
			&gs.House.InputType.InputTypeName, &gs.House.Reliability.ReliabilityName, &gs.House.Voltage.VoltageName, &gs.House.Voltage.VoltageValue,
			&gs.ObjType.ObjTypeName, &gs.Uzo.UzoName, &gs.Uzo.UzoValue, &gs.TariffGroup.TariffGroupName, &gs.CalculationType.CalculationTypeName,
			&gs.ObjStatus.ObjStatusName)
		if err != nil {
			log.Println("failed to scan row:", err)
		}

		out_arr = append(out_arr, gs)
	}

	gsc := 0
	err = s.Dbpool.QueryRow(ctx, "SELECT * from func_objects_cnt($1,$2,$3);", gs1, NullableString(gs2), utils.NullableBool(gs3, gs3f)).Scan(&gsc)

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	auth := models.Auth{Create: true, Read: true, Update: true, Delete: true}
	out_count, err := json.Marshal(models.Object_count{Values: out_arr, Count: gsc, Auth: auth})
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
func (s *APG) HandleAddObject(w http.ResponseWriter, r *http.Request) {
	a := models.AddObject{}
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
	err = s.Dbpool.QueryRow(context.Background(), "SELECT func_objects_add($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11);", a.ObjectName, a.House.Id,
		a.FlatNumber, a.ObjType.Id, a.RegQty, a.Uzo.Id, a.TariffGroup.Id, a.CalculationType.Id, a.ObjStatus.Id, a.Notes, a.MffId).Scan(&ai)

	if err != nil {
		log.Println("Failed execute func_objects_add: ", err)
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
func (s *APG) HandleUpdObject(w http.ResponseWriter, r *http.Request) {
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

	ui := 0
	err = s.Dbpool.QueryRow(context.Background(), "SELECT func_objects_upd($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12);", u.Id, u.ObjectName,
		u.House.Id, u.FlatNumber, u.ObjType.Id, u.RegQty, u.Uzo.Id, u.TariffGroup.Id, u.CalculationType.Id, u.ObjStatus.Id, u.Notes,
		u.MffId).Scan(&ui)

	if err != nil {
		log.Println("Failed execute func_objects_upd: ", err)
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
func (s *APG) HandleDelObject(w http.ResponseWriter, r *http.Request) {
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
		err = s.Dbpool.QueryRow(context.Background(), "SELECT func_objects_del($1);", id).Scan(&i)
		res = append(res, i)

		if err != nil {
			log.Println("Failed execute func_objects_del: ", err)
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

// HandleGetObject godoc
// @Summary Get object
// @Description get object
// @Tags objects
// @Produce  json
// @Param id path int true "Object by id"
// @Success 200 {object} models.Object_count
// @Failure 500
// @Router /objects/{id} [get]
func (s *APG) HandleGetObject(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	i := vars["id"]
	g := models.Object{}
	out_arr := []models.Object{}

	err := s.Dbpool.QueryRow(context.Background(), "SELECT * from func_object_get($1);", i).Scan(&g.Id, &g.ObjectName, &g.House.Id,
		&g.FlatNumber, &g.ObjType.Id, &g.RegQty, &g.Uzo.Id, &g.TariffGroup.Id, &g.Notes, &g.CalculationType.Id, &g.ObjStatus.Id, &g.MffId,
		&g.House.BuildingType.Id, &g.House.Street.Id, &g.House.HouseNumber, &g.House.BuildingNumber, &g.House.RP.Id, &g.House.Area.Id,
		&g.House.Ksk.Id, &g.House.Sector.Id, &g.House.Connector.Id, &g.House.InputType.Id, &g.House.Reliability.Id, &g.House.Voltage.Id,
		&g.House.BuildingType.BuildingTypeName, &g.House.Street.StreetName, &g.House.Street.Created, &g.House.Street.City.CityName,
		&g.House.RP.RpName, &g.House.Area.AreaName, &g.House.Area.AreaNumber, &g.House.Ksk.KskName, &g.House.Sector.SectorName,
		&g.House.Connector.ConnectorName, &g.House.InputType.InputTypeName, &g.House.Reliability.ReliabilityName, &g.House.Voltage.VoltageName,
		&g.House.Voltage.VoltageValue, &g.ObjType.ObjTypeName, &g.Uzo.UzoName, &g.Uzo.UzoValue, &g.TariffGroup.TariffGroupName,
		&g.CalculationType.CalculationTypeName, &g.ObjStatus.ObjStatusName)

	if err != nil && err != pgx.ErrNoRows {
		log.Println("Failed execute from func_object_get: ", err)
	}

	out_arr = append(out_arr, g)
	auth := models.Auth{Create: true, Read: true, Update: true, Delete: true}

	// output, err := json.Marshal(g)
	out_count, err := json.Marshal(models.Object_count{Values: out_arr, Count: 1, Auth: auth})
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
func (s *APG) HandleGetObjectContract(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	i := vars["id"]

	gs := models.ObjContract{}
	ctx := context.Background()
	out_arr := []models.ObjContract{}

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

	rows, err := s.Dbpool.Query(ctx, "SELECT * from func_object_getcontract($1,$2);", i, utils.NullableString(dsc))
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&gs.Id, &gs.Contract.Id, &gs.Object.Id, &gs.Startdate, &gs.Enddate, &gs.Object.ObjectName, &gs.Object.RegQty,
			&gs.Object.FlatNumber, &gs.Object.House.Id, &gs.Object.House.HouseNumber, &gs.Object.House.BuildingNumber,
			&gs.Object.House.Street.City.CityName, &gs.Object.House.Street.Id, &gs.Object.House.Street.StreetName, &gs.Object.TariffGroup.Id,
			&gs.Object.TariffGroup.TariffGroupName, &gs.Contract.ContractNumber, &gs.Contract.Startdate, &gs.Contract.Enddate,
			&gs.Contract.Customer.SubId, &gs.Contract.Customer.SubName, &gs.Contract.Consignee.SubId, &gs.Contract.Consignee.SubName)
		if err != nil {
			log.Println("failed to scan row:", err)
		}

		out_arr = append(out_arr, gs)
	}

	gsc := 0
	err = s.Dbpool.QueryRow(ctx, "SELECT * from func_object_getcontract_cnt($1,$2);", i, utils.NullableString(dsc)).Scan(&gsc)

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	out_count, err := json.Marshal(models.ObjContract_count{Values: out_arr, Count: gsc})
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
func (s *APG) HandleGetObjectMff(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	i := vars["hid"]
	g := models.Object{}
	out_arr := []models.Object{}

	err := s.Dbpool.QueryRow(context.Background(), "SELECT * from func_objects_mff($1);", i).Scan(&g.Id, &g.ObjectName, &g.House.Id,
		&g.FlatNumber, &g.ObjType.Id, &g.RegQty, &g.Uzo.Id, &g.TariffGroup.Id, &g.Notes, &g.CalculationType.Id, &g.ObjStatus.Id, &g.MffId)

	if err != nil && err != pgx.ErrNoRows {
		log.Println("Failed execute from func_objects_mff: ", err)
	}

	out_arr = append(out_arr, g)
	auth := models.Auth{Create: true, Read: true, Update: true, Delete: true}

	// output, err := json.Marshal(g)
	out_count, err := json.Marshal(models.Object_count{Values: out_arr, Count: 1, Auth: auth})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(out_count)

	return

}
