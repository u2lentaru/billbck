package api

import (
	"context"
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v4"
	"github.com/u2lentaru/billbck/internal/models"
)

// HandleActDetails godoc
// @Summary List act details
// @Description get act details list
// @Tags actdetails
// @Produce  json
// @Param page query int false "page number"
// @Param page_size query int false "page size"
// @Param actid query int false "actid search pattern"
// @Param ordering query string false "order by {punumber|installdate}"
// @Param desc query boolean false "descending order {true|false}"
// @Success 200 {object} models.ActDetail_count
// @Failure 500
// @Router /actdetails [get]
func (s *APG) HandleActDetails(w http.ResponseWriter, r *http.Request) {
	gs := models.ActDetail{}
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
	gs1s, ok := query["actid"]
	if ok && len(gs1s) > 0 {
		t, err := strconv.Atoi(gs1s[0])
		if err == nil {
			gs1 = t
		}
	}

	gsc := 0
	err := s.Dbpool.QueryRow(ctx, "SELECT * from func_act_details_cnt($1);", gs1).Scan(&gsc)

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	out_arr := make([]models.ActDetail, 0,
		func() int {
			if gsc < pgs {
				return gsc
			} else {
				return pgs
			}
		}())

	ord := 10
	ords, ok := query["ordering"]
	if !ok || len(ords) == 0 {
		ord = 10
	} else if ords[0] == "punumber" {
		ord = 12
	} else if ords[0] == "installdate" {
		ord = 15
	}

	dsc := false
	dscs, ok := query["desc"]
	if ok && len(dscs) > 0 {
		if !(dscs[0] == "0") {
			dsc = true
		}
	}

	rows, err := s.Dbpool.Query(ctx, "SELECT * from func_act_details_get($1,$2,$3,$4,$5);", pg, pgs, gs1, ord, dsc)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	defer rows.Close()

	var pti, si, ci, sti, ri, vi sql.NullInt32
	var ccn, ptn, stn, rn, vn sql.NullString

	for rows.Next() {
		// err = rows.Scan(&gs.Act.Id, &gs.Act.ActType.Id, &gs.Act.ActNumber, &gs.Act.ActDate, &gs.Act.Object.Id, &gs.Act.Staff.Id,
		// 	&gs.Act.Customer, &gs.Act.Notes, &gs.Act.ActType.ActTypeName, &gs.Act.Object.ObjectName, &gs.Act.Object.FlatNumber,
		// 	&gs.Act.Object.RegQty, &gs.Act.Object.House.Street.StreetName, &gs.Act.Object.House.Street.City.CityName,
		// 	&gs.Act.Object.House.HouseNumber, &gs.Act.Object.House.BuildingNumber, &gs.Act.Object.TariffGroup.TariffGroupName,
		// 	&gs.Act.Staff.StaffName, &gs.Id, &gs.ActDetailDate, &gs.PuId, &gs.PuType.Id, &gs.PuNumber, &gs.InstallDate, &gs.CheckInterval,
		// 	&gs.InitialValue, &gs.DevStopped, &gs.Startdate, &gs.Enddate, &gs.Pid, &gs.AdPuValue, &gs.Seal.Id, &gs.Seal.SealNumber,
		// 	&gs.SealDate, &gs.Conclusion.Id, &gs.Conclusion.ConclusionName, &gs.ShutdownType.Id, &gs.CustomerPhone, &gs.CustomerPos,
		// 	&gs.Notes, &gs.PuType.PuTypeName, &gs.Conclusion.ConclusionName, &gs.ShutdownType.ShutdownTypeName)

		err = rows.Scan(&gs.Act.Id, &gs.Act.ActType.Id, &gs.Act.ActNumber, &gs.Act.ActDate, &gs.Act.Object.Id, &gs.Act.Staff.Id,
			&gs.Act.Notes, &gs.Act.Activated, &gs.Act.ActType.ActTypeName, &gs.Act.Object.ObjectName, &gs.Act.Object.FlatNumber,
			&gs.Act.Object.RegQty, &gs.Act.Object.House.Street.StreetName, &gs.Act.Object.House.Street.City.CityName,
			&gs.Act.Object.House.HouseNumber, &gs.Act.Object.House.BuildingNumber, &gs.Act.Object.TariffGroup.TariffGroupName,
			&gs.Act.Staff.StaffName, &gs.Id, &gs.ActDetailDate, &gs.PuId, &pti, &gs.PuNumber, &gs.InstallDate, &gs.CheckInterval,
			&gs.InitialValue, &gs.DevStopped, &gs.Startdate, &gs.Enddate, &gs.Pid, &gs.AdPuValue, &si, &gs.SealNumber,
			&gs.SealDate, &ci, &gs.ConclusionNumber, &sti, &ri, &vi, &gs.Customer, &gs.CustomerPhone, &gs.CustomerPos,
			&gs.Notes, &ptn, &ccn, &stn, &rn, &vn)

		if err != nil {
			log.Println("failed to scan row:", err)
		}

		gs.Reason.Id = int(ri.Int32)
		gs.Violation.Id = int(vi.Int32)
		gs.Reason.ReasonName = rn.String
		gs.Violation.ViolationName = vn.String
		gs.PuType.Id = int(pti.Int32)
		gs.Seal.Id = int(si.Int32)
		// gs.Seal.SealNumber = ssn.String
		gs.Conclusion.Id = int(ci.Int32)
		gs.Conclusion.ConclusionName = ccn.String
		gs.ShutdownType.Id = int(sti.Int32)
		gs.PuType.PuTypeName = ptn.String
		gs.ShutdownType.ShutdownTypeName = stn.String

		out_arr = append(out_arr, gs)
	}

	auth := models.Auth{Create: true, Read: true, Update: true, Delete: true}

	out_count, err := json.Marshal(models.ActDetail_count{Values: out_arr, Count: gsc, Auth: auth})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Write(out_count)

	return

}

// HandleAddActDetail godoc
// @Summary Add act detail
// @Description add act detail
// @Tags actdetails
// @Accept json
// @Produce  json
// @Param a body models.AddActDetail true "New act detail. Significant params: Act.Id, PuId, SealNumber, AdPuValue, ActDetailDate, PuNumber, InstallDate, CheckInterval, InitialValue, DevStopped, Startdate, Enddate, Pid, Seal.Id, SealDate, Notes, PuType.Id, Conclusion.Id, ConclusionNumber, ShutdownType.Id, CustomerPhone, CustomerPos, Reason.Id, Violation.Id, Customer"
// @Success 200 {object} models.Json_id
// @Failure 500
// @Router /actdetails_add [post]
func (s *APG) HandleAddActDetail(w http.ResponseWriter, r *http.Request) {
	a := models.AddActDetail{}
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
	err = s.Dbpool.QueryRow(context.Background(),
		"SELECT func_act_details_add($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19,$20,$21,$22,$23,$24,$25);",
		a.Act.Id, a.PuId, a.SealNumber, a.AdPuValue, a.ActDetailDate, a.PuNumber, a.InstallDate, a.CheckInterval, a.InitialValue,
		a.DevStopped, a.Startdate, a.Enddate, a.Pid, a.Seal.Id, a.SealDate, a.Notes, a.PuType.Id, a.Conclusion.Id,
		a.ConclusionNumber, a.ShutdownType.Id, a.CustomerPhone, a.CustomerPos, a.Reason.Id, a.Violation.Id, a.Customer).Scan(&ai)

	if err != nil {
		log.Println("Failed execute func_act_details_add: ", err)
	}

	output, err := json.Marshal(models.Json_id{Id: ai})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(output)

	return

}

// HandleUpdActDetail godoc
// @Summary Update act detail
// @Description update act detail
// @Tags actdetails
// @Accept json
// @Produce  json
// @Param u body models.ActDetail true "Update act detail. Significant params: Id, Act.Id, PuId, SealNumber, AdPuValue, ActDetailDate, PuNumber, InstallDate, CheckInterval, InitialValue, DevStopped, Startdate, Enddate, Pid, Seal.Id, SealDate, Notes, PuType.Id, Conclusion.Id, ConclusionNumber, ShutdownType.Id, CustomerPhone, CustomerPos, Reason.Id, Violation.Id, Customer"
// @Success 200 {object} models.Json_id
// @Failure 500
// @Router /actdetails_upd [post]
func (s *APG) HandleUpdActDetail(w http.ResponseWriter, r *http.Request) {
	u := models.ActDetail{}
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
	err = s.Dbpool.QueryRow(context.Background(),
		"SELECT func_act_details_upd($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19,$20,$21,$22,$23,$24,$25,$26);",
		u.Id, u.Act.Id, u.PuId, u.SealNumber, u.AdPuValue, u.ActDetailDate, u.PuNumber, u.InstallDate, u.CheckInterval, u.InitialValue,
		u.DevStopped, u.Startdate, u.Enddate, u.Pid, u.Seal.Id, u.SealDate, u.Notes, u.PuType.Id, u.Conclusion.Id,
		u.ConclusionNumber, u.ShutdownType.Id, u.CustomerPhone, u.CustomerPos, u.Reason.Id, u.Violation.Id, u.Customer).Scan(&ui)

	if err != nil {
		log.Println("Failed execute func_act_details_upd: ", err)
	}

	output, err := json.Marshal(models.Json_id{Id: ui})

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(output)

	return

}

// HandleDelActDetail godoc
// @Summary Delete act details
// @Description delete act details
// @Tags actdetails
// @Accept json
// @Produce  json
// @Param d body models.Json_ids true "Delete act details"
// @Success 200 {object} models.Json_ids
// @Failure 500
// @Router /actdetails_del [post]
func (s *APG) HandleDelActDetail(w http.ResponseWriter, r *http.Request) {
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
		err = s.Dbpool.QueryRow(context.Background(), "SELECT func_act_details_del($1);", id).Scan(&i)
		res = append(res, i)

		if err != nil {
			log.Println("Failed execute func_act_details_del: ", err)
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

// HandleGetActDetail godoc
// @Summary Get act detail
// @Description get act detail
// @Tags actdetails
// @Produce  json
// @Param id path int true "Act detail by id"
// @Success 200 {object} models.ActDetail_count
// @Failure 500
// @Router /actdetails/{id} [get]
func (s *APG) HandleGetActDetail(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	i := vars["id"]
	g := models.ActDetail{}

	out_arr := []models.ActDetail{}

	var pti, si, ci, sti, ri, vi sql.NullInt32
	var ccn, ptn, stn, rn, vn sql.NullString

	// err := s.Dbpool.QueryRow(context.Background(), "SELECT * from func_act_detail_get($1);", i).Scan(&g.Act.Id, &g.Act.ActType.Id,
	// 	&g.Act.ActNumber, &g.Act.ActDate, &g.Act.Object.Id, &g.Act.Staff.Id, &g.Act.Customer, &g.Act.Notes, &g.Act.ActType.ActTypeName,
	// 	&g.Act.Object.ObjectName, &g.Act.Object.FlatNumber, &g.Act.Object.RegQty, &g.Act.Object.House.Street.StreetName,
	// 	&g.Act.Object.House.Street.City.CityName, &g.Act.Object.House.HouseNumber, &g.Act.Object.House.BuildingNumber,
	// 	&g.Act.Object.TariffGroup.TariffGroupName, &g.Act.Staff.StaffName, &g.Id, &g.ActDetailDate, &g.PuId, &g.PuType.Id, &g.PuNumber,
	// 	&g.InstallDate, &g.CheckInterval, &g.InitialValue, &g.DevStopped, &g.Startdate, &g.Enddate, &g.Pid, &g.AdPuValue, &g.Seal.Id,
	// 	&g.Seal.SealNumber, &g.SealDate, &g.Conclusion.Id, &g.Conclusion.ConclusionName, &g.ShutdownType.Id, &g.CustomerPhone, &g.CustomerPos,
	// 	&g.Notes, &g.PuType.PuTypeName, &g.Conclusion.ConclusionName, &g.ShutdownType.ShutdownTypeName)

	err := s.Dbpool.QueryRow(context.Background(), "SELECT * from func_act_detail_get($1);", i).Scan(&g.Act.Id, &g.Act.ActType.Id,
		&g.Act.ActNumber, &g.Act.ActDate, &g.Act.Object.Id, &g.Act.Staff.Id, &g.Act.Notes, &g.Act.Activated, &g.Act.ActType.ActTypeName,
		&g.Act.Object.ObjectName, &g.Act.Object.FlatNumber, &g.Act.Object.RegQty, &g.Act.Object.House.Street.StreetName,
		&g.Act.Object.House.Street.City.CityName, &g.Act.Object.House.HouseNumber, &g.Act.Object.House.BuildingNumber,
		&g.Act.Object.TariffGroup.TariffGroupName, &g.Act.Staff.StaffName, &g.Id, &g.ActDetailDate, &g.PuId, &pti, &g.PuNumber,
		&g.InstallDate, &g.CheckInterval, &g.InitialValue, &g.DevStopped, &g.Startdate, &g.Enddate, &g.Pid, &g.AdPuValue, &si,
		&g.SealNumber, &g.SealDate, &ci, &g.ConclusionNumber, &sti, &ri, &vi, &g.Customer, &g.CustomerPhone, &g.CustomerPos,
		&g.Notes, &ptn, &ccn, &stn, &rn, &vn)

	if err != nil && err != pgx.ErrNoRows {
		log.Println("Failed execute from func_act_detail_get: ", err)
	}

	g.Reason.Id = int(ri.Int32)
	g.Violation.Id = int(vi.Int32)
	g.Reason.ReasonName = rn.String
	g.Violation.ViolationName = vn.String
	g.PuType.Id = int(pti.Int32)
	g.Seal.Id = int(si.Int32)
	// g.Seal.SealNumber = ssn.String
	g.Conclusion.Id = int(ci.Int32)
	g.Conclusion.ConclusionName = ccn.String
	g.ShutdownType.Id = int(sti.Int32)
	g.PuType.PuTypeName = ptn.String
	g.ShutdownType.ShutdownTypeName = stn.String

	out_arr = append(out_arr, g)
	auth := models.Auth{Create: true, Read: true, Update: true, Delete: true}

	// output, err := json.Marshal(g)
	out_count, err := json.Marshal(models.ActDetail_count{Values: out_arr, Count: 1, Auth: auth})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	// w.Write(output)
	w.Write(out_count)

	return
}
