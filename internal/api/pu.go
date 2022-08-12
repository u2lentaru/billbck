package api

import (
	"context"
	"database/sql"
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
)

func NullableString(s string) sql.NullString {
	return sql.NullString{String: s, Valid: s != ""}
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
func (s *APG) HandlePu(w http.ResponseWriter, r *http.Request) {
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

	rows, err := s.Dbpool.Query(ctx, "SELECT * from func_pu_get($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11);", pg, pgs, gs1, gs2, gs3, NullableString(gs4),
		NullableString(gs5), NullableString(gs6), NullableString(gs7), ord, dsc)
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
	err = s.Dbpool.QueryRow(ctx, "SELECT * from func_pu_cnt($1,$2,$3,$4,$5,$6,$7);", gs1, gs2, gs3, NullableString(gs4), NullableString(gs5),
		NullableString(gs6), NullableString(gs7)).Scan(&gsc)

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
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
func (s *APG) HandleAddPu(w http.ResponseWriter, r *http.Request) {
	a := models.AddPu{}
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
	err = s.Dbpool.QueryRow(context.Background(), "SELECT func_pu_add($1,$2,$3,$4,$5,$6,$7,$8,$9,$10);",
		a.Object.Id, a.PuObjectType, a.PuType.Id, a.PuNumber, a.InstallDate, a.CheckInterval, a.InitialValue, a.DevStopped, a.Startdate,
		a.Pid).Scan(&ai)

	if err != nil {
		log.Println("Failed execute func_pu_add: ", err)
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
func (s *APG) HandleUpdPu(w http.ResponseWriter, r *http.Request) {
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

	ui := 0
	err = s.Dbpool.QueryRow(context.Background(), "SELECT func_pu_upd($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12);", u.Id, u.Object.Id, u.PuObjectType,
		u.PuType.Id, u.PuNumber, u.InstallDate, u.CheckInterval, u.InitialValue, u.DevStopped, u.Startdate, u.Enddate, u.Pid).Scan(&ui)

	if err != nil {
		log.Println("Failed execute func_pu_upd: ", err)
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
func (s *APG) HandleDelPu(w http.ResponseWriter, r *http.Request) {
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
		err = s.Dbpool.QueryRow(context.Background(), "SELECT func_pu_del($1);", id).Scan(&i)
		res = append(res, i)

		if err != nil {
			log.Println("Failed execute func_pu_del: ", err)
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

// HandleGetPu godoc
// @Summary Get pu
// @Description get pu
// @Tags pu
// @Produce  json
// @Param id path int true "Pu by id"
// @Success 200 {object} models.Pu_count
// @Failure 500
// @Router /pu/{id} [get]
func (s *APG) HandleGetPu(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	i := vars["id"]
	g := models.Pu{}
	out_arr := []models.Pu{}

	err := s.Dbpool.QueryRow(context.Background(), "SELECT * from func_pu_getbyid($1);", i).Scan(&g.Id, &g.Startdate, &g.Enddate,
		&g.PuType.Id, &g.PuType.PuTypeName, &g.PuNumber, &g.InstallDate, &g.CheckInterval, &g.InitialValue, &g.DevStopped, &g.Object.Id,
		&g.PuObjectType, &g.Object.ObjectName, &g.Object.House.Id, &g.Object.House.HouseNumber, &g.Object.FlatNumber,
		&g.Object.House.BuildingNumber, &g.Object.RegQty, &g.Object.House.Street.Id, &g.Object.House.Street.StreetName,
		&g.Object.House.Street.City.CityName, &g.Object.House.BuildingType.BuildingTypeName, &g.Object.House.Street.City.Id,
		&g.Object.House.Street.Created, &g.Object.House.Street.Closed, &g.Pid)

	if err != nil && err != pgx.ErrNoRows {
		log.Println("Failed execute from func_pu_getbyid: ", err)
	}

	out_arr = append(out_arr, g)
	auth := models.Auth{Create: true, Read: true, Update: true, Delete: true}

	// output, err := json.Marshal(g)
	out_count, err := json.Marshal(models.Pu_count{Values: out_arr, Count: 1, Auth: auth})
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
func (s *APG) HandlePuObj(w http.ResponseWriter, r *http.Request) {
	gs := models.Pu{}
	ctx := context.Background()

	out_arr := []models.Pu{}

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

	rows, err := s.Dbpool.Query(ctx, "SELECT * from func_pu_obj($1,$2);", gs1, gs2)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	defer rows.Close()

	gsc := 0

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
		gsc++
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
