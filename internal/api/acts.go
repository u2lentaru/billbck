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

// HandleActs godoc
// @Summary List acts
// @Description get act list
// @Tags acts
// @Produce  json
// @Param page query int false "page number"
// @Param page_size query int false "page size"
// @Param actnumber query string false "actnumber search pattern"
// @Param objectname query string false "objectname search pattern"
// @Param objectid query int false "objectid search pattern"
// @Param ordering query string false "order by {actnumber|objectname}"
// @Param desc query boolean false "descending order {true|false}"
// @Success 200 {array} models.Act_count
// @Failure 500
// @Router /acts [get]
func (s *APG) HandleActs(w http.ResponseWriter, r *http.Request) {
	gs := models.Act{}
	ctx := context.Background()
	out_arr := []models.Act{}

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
	gs1s, ok := query["actnumber"]
	if ok && len(gs1s) > 0 {
		//case insensitive
		gs1 = strings.ToUpper(gs1s[0])
		//quotes
		re := regexp.MustCompile(`'`)
		gs1 = string(re.ReplaceAll([]byte(gs1), []byte("''")))
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

	gs3 := 0
	gs3s, ok := query["objectid"]
	if ok && len(gs3s) > 0 {
		t, err := strconv.Atoi(gs3s[0])
		if err == nil {
			gs3 = t
		}
	}

	ord := 1
	ords, ok := query["ordering"]
	if !ok || len(ords) == 0 {
		ord = 1
	} else if ords[0] == "actnumber" {
		ord = 3
	} else if ords[0] == "objectname" {
		ord = 7
	}

	dsc := false
	dscs, ok := query["desc"]
	if ok && len(dscs) > 0 {
		if !(dscs[0] == "0") {
			dsc = true
		}
	}

	rows, err := s.Dbpool.Query(ctx, "SELECT * from func_acts_get($1,$2,$3,$4,$5,$6,$7);", pg, pgs, gs1, gs2, gs3, ord, dsc)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	defer rows.Close()

	for rows.Next() {
		// err = rows.Scan(&gs.Id, &gs.ActType.Id, &gs.ActNumber, &gs.ActDate, &gs.Object.Id, &gs.Staff.Id, &gs.Customer, &gs.Notes,
		// 	&gs.ActType.ActTypeName, &gs.Object.ObjectName, &gs.Object.FlatNumber, &gs.Object.RegQty, &gs.Object.House.Street.StreetName,
		// 	&gs.Object.House.Street.City.CityName, &gs.Object.House.HouseNumber, &gs.Object.House.BuildingNumber,
		// 	&gs.Object.TariffGroup.TariffGroupName, &gs.Staff.StaffName)

		err = rows.Scan(&gs.Id, &gs.ActType.Id, &gs.ActNumber, &gs.ActDate, &gs.Object.Id, &gs.Staff.Id, &gs.Notes, &gs.Activated,
			&gs.ActType.ActTypeName, &gs.Object.ObjectName, &gs.Object.FlatNumber, &gs.Object.RegQty, &gs.Object.House.Street.StreetName,
			&gs.Object.House.Street.City.CityName, &gs.Object.House.HouseNumber, &gs.Object.House.BuildingNumber,
			&gs.Object.TariffGroup.TariffGroupName, &gs.Staff.StaffName)

		if err != nil {
			log.Println("failed to scan row:", err)
		}

		out_arr = append(out_arr, gs)
	}

	gsc := 0
	err = s.Dbpool.QueryRow(ctx, "SELECT * from func_acts_cnt($1,$2,$3);", gs1, gs2, gs3).Scan(&gsc)

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	auth := models.Auth{Create: true, Read: true, Update: true, Delete: true}
	out_count, err := json.Marshal(models.Act_count{Values: out_arr, Count: gsc, Auth: auth})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Write(out_count)

	return

}

// HandleAddAct godoc
// @Summary Add act
// @Description add act
// @Tags acts
// @Accept json
// @Produce  json
// @Param a body models.AddAct true "New act. Significant params: ActType.Id, ActNumber, ActDate, Object.Id, Staff.Id, Notes"
// @Success 200 {array} models.Json_id
// @Failure 500
// @Router /acts_add [post]
func (s *APG) HandleAddAct(w http.ResponseWriter, r *http.Request) {
	a := models.AddAct{}
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
	err = s.Dbpool.QueryRow(context.Background(), "SELECT func_acts_add($1,$2,$3,$4,$5,$6);",
		a.ActType.Id, a.ActNumber, a.ActDate, a.Object.Id, a.Staff.Id, a.Notes).Scan(&ai)

	if err != nil {
		log.Println("Failed execute func_acts_add: ", err)
	}

	output, err := json.Marshal(models.Json_id{Id: ai})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(output)

	return

}

// HandleUpdAct godoc
// @Summary Update act
// @Description update act
// @Tags acts
// @Accept json
// @Produce  json
// @Param u body models.Act true "Update act. Significant params: Id, ActType.Id, ActNumber, ActDate, Object.Id, Staff.Id, Notes"
// @Success 200 {array} models.Json_id
// @Failure 500
// @Router /acts_upd [post]
func (s *APG) HandleUpdAct(w http.ResponseWriter, r *http.Request) {
	u := models.Act{}
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
	err = s.Dbpool.QueryRow(context.Background(), "SELECT func_acts_upd($1,$2,$3,$4,$5,$6,$7);", u.Id, u.ActType.Id, u.ActNumber,
		u.ActDate, u.Object.Id, u.Staff.Id, u.Notes).Scan(&ui)

	if err != nil {
		log.Println("Failed execute func_acts_upd: ", err)
	}

	output, err := json.Marshal(models.Json_id{Id: ui})

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(output)

	return
}

// HandleDelAct godoc
// @Summary Delete acts
// @Description delete acts
// @Tags acts
// @Accept json
// @Produce  json
// @Param d body models.Json_ids true "Delete acts"
// @Success 200 {array} models.Json_ids
// @Failure 500
// @Router /acts_del [post]
func (s *APG) HandleDelAct(w http.ResponseWriter, r *http.Request) {
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
		err = s.Dbpool.QueryRow(context.Background(), "SELECT func_acts_del($1);", id).Scan(&i)
		res = append(res, i)

		if err != nil {
			log.Println("Failed execute func_acts_del: ", err)
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

// HandleGetAct godoc
// @Summary Get act
// @Description get act
// @Tags acts
// @Produce  json
// @Param id path int true "Act by id"
// @Success 200 {array} models.Act_count
// @Failure 500
// @Router /acts/{id} [get]
func (s *APG) HandleGetAct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	i := vars["id"]
	g := models.Act{}

	out_arr := []models.Act{}

	err := s.Dbpool.QueryRow(context.Background(), "SELECT * from func_act_get($1);", i).Scan(&g.Id, &g.ActType.Id, &g.ActNumber,
		&g.ActDate, &g.Object.Id, &g.Staff.Id, &g.Notes, &g.Activated, &g.ActType.ActTypeName, &g.Object.ObjectName,
		&g.Object.FlatNumber, &g.Object.RegQty, &g.Object.House.Street.StreetName, &g.Object.House.Street.City.CityName,
		&g.Object.House.HouseNumber, &g.Object.House.BuildingNumber, &g.Object.TariffGroup.TariffGroupName, &g.Staff.StaffName)

	if err != nil && err != pgx.ErrNoRows {
		log.Println("Failed execute from func_act_get: ", err)
	}

	out_arr = append(out_arr, g)
	auth := models.Auth{Create: true, Read: true, Update: true, Delete: true}

	// output, err := json.Marshal(g)
	out_count, err := json.Marshal(models.Act_count{Values: out_arr, Count: 1, Auth: auth})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(out_count)

	return

}

// HandleActActivate godoc
// @Summary Act activation
// @Description act activation
// @Tags acts
// @Produce  json
// @Param actid query int true "actid"
// @Param activationdate query string true "activation date"
// @Success 200 {array} models.Json_id
// @Failure 500
// @Router /acts_activate [get]
func (s *APG) HandleActActivate(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	query := r.URL.Query()

	gs1 := 0
	gs1s, ok := query["actid"]
	if ok && len(gs1s) > 0 {
		t, err := strconv.Atoi(gs1s[0])
		if err == nil {
			gs1 = t
		}
	}

	gs2 := ""
	gs2s, ok := query["activationdate"]
	if ok && len(gs2s) > 0 {
		//case insensitive
		gs2 = strings.ToUpper(gs2s[0])
		//quotes
		re := regexp.MustCompile(`'`)
		//gs2, err := time.Parse("2006-01-02", gs2s[0])
		gs2 = string(re.ReplaceAll([]byte(gs2), []byte("''")))
	}

	gsc := 0
	err := s.Dbpool.QueryRow(ctx, "SELECT * from func_acts_activate($1,$2);", utils.NullableInt(int32(gs1)), utils.NullableString(gs2)).Scan(&gsc)

	if err != nil {
		// w.Write([]byte(err.Error()))
		http.Error(w, err.Error(), 500)
		return
	}

	output, err := json.Marshal(models.Json_id{Id: gsc})

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(output)

	return
}
