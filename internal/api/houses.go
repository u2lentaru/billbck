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
)

// HandleHouses godoc
// @Summary List houses
// @Description get house list
// @Tags houses
// @Produce  json
// @Param page query int false "page number"
// @Param page_size query int false "page size"
// @Param streetname query string false "streetname search pattern"
// @Param housenumber query string false "housenumber search pattern"
// @Param streetid query int false "streetid search pattern"
// @Param ordering query string false "order by {housenumber|streetname}"
// @Param desc query boolean false "descending order {true|false}"
// @Success 200 {object} models.House_count
// @Failure 500
// @Router /houses [get]
func (s *APG) HandleHouses(w http.ResponseWriter, r *http.Request) {
	gs := models.House{}
	ctx := context.Background()
	out_arr := []models.House{}

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
	gs1s, ok := query["streetname"]
	if ok && len(gs1s) > 0 {
		//case insensitive
		gs1 = strings.ToUpper(gs1s[0])
		//quotes
		re := regexp.MustCompile(`'`)
		gs1 = string(re.ReplaceAll([]byte(gs1), []byte("''")))
	}

	gs2 := ""
	gs2s, ok := query["housenumber"]
	if ok && len(gs2s) > 0 {
		//case insensitive
		gs2 = strings.ToUpper(gs2s[0])
		//quotes
		re := regexp.MustCompile(`'`)
		gs2 = string(re.ReplaceAll([]byte(gs2), []byte("''")))
	}

	gs3 := 0
	gs3s, ok := query["streetid"]
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
	} else if ords[0] == "streetname" {
		ord = 15
	} else if ords[0] == "housenumber" {
		ord = 4
	}

	dsc := false
	dscs, ok := query["desc"]
	if ok && len(dscs) > 0 {
		if !(dscs[0] == "0") {
			dsc = true
		}
	}

	rows, err := s.Dbpool.Query(ctx, "SELECT * from func_houses_get($1,$2,$3,$4,$5,$6,$7);", pg, pgs, gs1, gs2, gs3, ord, dsc)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&gs.Id, &gs.BuildingType.Id, &gs.Street.Id, &gs.HouseNumber, &gs.BuildingNumber, &gs.RP.Id, &gs.Area.Id, &gs.Ksk.Id,
			&gs.Sector.Id, &gs.Connector.Id, &gs.InputType.Id, &gs.Reliability.Id, &gs.Voltage.Id, &gs.Notes, &gs.BuildingType.BuildingTypeName,
			&gs.Street.StreetName, &gs.Street.Created, &gs.Street.City.CityName, &gs.RP.RpName, &gs.Area.AreaName, &gs.Area.AreaNumber, &gs.Ksk.KskName,
			&gs.Sector.SectorName, &gs.Connector.ConnectorName, &gs.InputType.InputTypeName, &gs.Reliability.ReliabilityName,
			&gs.Voltage.VoltageName, &gs.Voltage.VoltageValue)
		if err != nil {
			log.Println("failed to scan row:", err)
		}

		out_arr = append(out_arr, gs)
	}

	gsc := 0
	err = s.Dbpool.QueryRow(ctx, "SELECT * from func_houses_cnt($1,$2,$3);", gs1, gs2, gs3).Scan(&gsc)

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	auth := models.Auth{Create: true, Read: true, Update: true, Delete: true}
	out_count, err := json.Marshal(models.House_count{Values: out_arr, Count: gsc, Auth: auth})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Write(out_count)

	return
}

// HandleAddHouse godoc
// @Summary Add house
// @Description add house
// @Tags houses
// @Accept json
// @Produce  json
// @Param a body models.AddHouse true "New house"
// @Success 200 {object} models.Json_id
// @Failure 500
// @Router /houses_add [post]
func (s *APG) HandleAddHouse(w http.ResponseWriter, r *http.Request) {
	a := models.AddHouse{}
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
	err = s.Dbpool.QueryRow(context.Background(), "SELECT func_houses_add($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13);",
		a.BuildingType.Id, a.Street.Id, a.HouseNumber, a.BuildingNumber, a.RP.Id, a.Area.Id, a.Ksk.Id, a.Sector.Id, a.Connector.Id,
		a.InputType.Id, a.Reliability.Id, a.Voltage.Id, a.Notes).Scan(&ai)

	if err != nil {
		log.Println("Failed execute func_houses_add: ", err)
	}

	output, err := json.Marshal(models.Json_id{Id: ai})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(output)

	return

}

// HandleUpdHouse godoc
// @Summary Update house
// @Description update house
// @Tags houses
// @Accept json
// @Produce  json
// @Param u body models.House true "Update house"
// @Success 200 {object} models.Json_id
// @Failure 500
// @Router /houses_upd [post]
func (s *APG) HandleUpdHouse(w http.ResponseWriter, r *http.Request) {
	u := models.House{}
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
	err = s.Dbpool.QueryRow(context.Background(), "SELECT func_houses_upd($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14);", u.Id,
		u.BuildingType.Id, u.Street.Id, u.HouseNumber, u.BuildingNumber, u.RP.Id, u.Area.Id, u.Ksk.Id, u.Sector.Id, u.Connector.Id,
		u.InputType.Id, u.Reliability.Id, u.Voltage.Id, u.Notes).Scan(&ui)

	if err != nil {
		log.Println("Failed execute func_houses_upd: ", err)
	}

	output, err := json.Marshal(models.Json_id{Id: ui})

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(output)

	return
}

// HandleDelHouse godoc
// @Summary Delete houses
// @Description delete houses
// @Tags houses
// @Accept json
// @Produce  json
// @Param d body models.Json_ids true "Delete houses"
// @Success 200 {object} models.Json_ids
// @Failure 500
// @Router /houses_del [post]
func (s *APG) HandleDelHouse(w http.ResponseWriter, r *http.Request) {
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
		err = s.Dbpool.QueryRow(context.Background(), "SELECT func_houses_del($1);", id).Scan(&i)
		res = append(res, i)

		if err != nil {
			log.Println("Failed execute func_houses_del: ", err)
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

// HandleGetHouse godoc
// @Summary Get house
// @Description get house
// @Tags houses
// @Produce  json
// @Param id path int true "House by id"
// @Success 200 {object} models.House_count
// @Failure 500
// @Router /houses/{id} [get]
func (s *APG) HandleGetHouse(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	i := vars["id"]
	g := models.House{}
	out_arr := []models.House{}

	err := s.Dbpool.QueryRow(context.Background(), "SELECT * from func_house_get($1);", i).Scan(&g.Id, &g.BuildingType.Id, &g.Street.Id,
		&g.HouseNumber, &g.BuildingNumber, &g.RP.Id, &g.Area.Id, &g.Ksk.Id, &g.Sector.Id, &g.Connector.Id, &g.InputType.Id, &g.Reliability.Id,
		&g.Voltage.Id, &g.Notes, &g.BuildingType.BuildingTypeName, &g.Street.StreetName, &g.Street.Created, &g.Street.City.CityName, &g.RP.RpName,
		&g.Area.AreaName, &g.Area.AreaNumber, &g.Ksk.KskName, &g.Sector.SectorName, &g.Connector.ConnectorName, &g.InputType.InputTypeName,
		&g.Reliability.ReliabilityName, &g.Voltage.VoltageName, &g.Voltage.VoltageValue)

	if err != nil && err != pgx.ErrNoRows {
		log.Println("Failed execute from func_house_get: ", err)
	}

	out_arr = append(out_arr, g)
	auth := models.Auth{Create: true, Read: true, Update: true, Delete: true}

	// output, err := json.Marshal(g)
	out_count, err := json.Marshal(models.House_count{Values: out_arr, Count: 1, Auth: auth})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(out_count)

	return
}
