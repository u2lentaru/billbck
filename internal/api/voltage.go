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
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/u2lentaru/billbck/internal/models"
)

type APG struct {
	Dbpool *pgxpool.Pool
}

// HandleVoltages godoc
// @Summary List voltages
// @Description get voltage list
// @Tags voltages
// @Produce  json
// @Param page query int false "page number"
// @Param page_size query int false "page size"
// @Param voltagename query string false "voltagename search pattern"
// @Param ordering query string false "order by {id|voltagename}"
// @Param desc query boolean false "descending order {true|false}"
// @Success 200 {array} models.Voltage_count
// @Failure 500
// @Router /voltages [get]
func (s *APG) HandleVoltages(w http.ResponseWriter, r *http.Request) {
	gs := models.Voltage{}
	ctx := context.Background()
	out_arr := []models.Voltage{}

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
	gs1s, ok := query["voltagename"]
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
	} else if ords[0] == "voltagename" {
		ord = 2
	}

	dsc := false
	dscs, ok := query["desc"]
	if ok && len(dscs) > 0 {
		if !(dscs[0] == "0") {
			dsc = true
		}
	}

	rows, err := s.Dbpool.Query(ctx, "SELECT * from func_voltages_get($1,$2,$3,$4,$5);", pg, pgs, gs1, ord, dsc)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&gs.Id, &gs.VoltageName, &gs.VoltageValue)
		if err != nil {
			log.Println("failed to scan row:", err)
		}

		out_arr = append(out_arr, gs)
	}

	gsc := 0
	err = s.Dbpool.QueryRow(ctx, "SELECT * from func_voltages_cnt($1);", gs1).Scan(&gsc)

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	auth := models.Auth{Create: true, Read: true, Update: true, Delete: true}
	out_count, err := json.Marshal(models.Voltage_count{Values: out_arr, Count: gsc, Auth: auth})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Write(out_count)

	return

}

// HandleAddVoltage godoc
// @Summary Add voltage
// @Description add voltage
// @Tags voltages
// @Accept json
// @Produce  json
// @Param a body models.AddVoltage true "New voltage"
// @Success 200 {array} models.Json_id
// @Failure 500
// @Router /voltages_add [post]
func (s *APG) HandleAddVoltage(w http.ResponseWriter, r *http.Request) {
	a := models.AddVoltage{}
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
	err = s.Dbpool.QueryRow(context.Background(), "SELECT func_voltages_add($1, $2);", a.VoltageName, a.VoltageValue).Scan(&ai)

	if err != nil {
		log.Println("Failed execute func_voltages_add: ", err)
	}

	output, err := json.Marshal(models.Json_id{Id: ai})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(output)

	return

}

// HandleUpdVoltage godoc
// @Summary Update voltage
// @Description update voltage
// @Tags voltages
// @Accept json
// @Produce  json
// @Param u body models.Voltage true "Update voltage"
// @Success 200 {array} models.Json_id
// @Failure 500
// @Router /voltages_upd [post]
func (s *APG) HandleUpdVoltage(w http.ResponseWriter, r *http.Request) {
	u := models.Voltage{}
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
	err = s.Dbpool.QueryRow(context.Background(), "SELECT func_voltages_upd($1,$2,$3);", u.Id, u.VoltageName, u.VoltageValue).Scan(&ui)

	if err != nil {
		log.Println("Failed execute func_voltages_upd: ", err)
	}

	output, err := json.Marshal(models.Json_id{Id: ui})

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(output)

	return

}

// HandleDelVoltage godoc
// @Summary Delete voltages
// @Description delete voltages
// @Tags voltages
// @Accept json
// @Produce  json
// @Param d body models.Json_ids true "Delete voltages"
// @Success 200 {array} models.Json_ids
// @Failure 500
// @Router /voltages_del [post]
func (s *APG) HandleDelVoltage(w http.ResponseWriter, r *http.Request) {
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
		err = s.Dbpool.QueryRow(context.Background(), "SELECT func_voltages_del($1);", id).Scan(&i)
		res = append(res, i)

		if err != nil {
			log.Println("Failed execute func_voltages_del: ", err)
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

// HandleGetVoltage godoc
// @Summary Get voltage
// @Description get voltage
// @Tags voltages
// @Produce  json
// @Param id path int true "Voltage by id"
// @Success 200 {array} models.Voltage_count
// @Failure 500
// @Router /voltages/{id} [get]
func (s *APG) HandleGetVoltage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	i := vars["id"]
	g := models.Voltage{}
	out_arr := []models.Voltage{}

	err := s.Dbpool.QueryRow(context.Background(), "SELECT * from func_voltage_get($1);", i).Scan(&g.Id, &g.VoltageName, &g.VoltageValue)

	if err != nil && err != pgx.ErrNoRows {
		log.Println("Failed execute from func_voltage_get: ", err)
	}

	out_arr = append(out_arr, g)
	auth := models.Auth{Create: true, Read: true, Update: true, Delete: true}

	// output, err := json.Marshal(g)
	out_count, err := json.Marshal(models.Voltage_count{Values: out_arr, Count: 1, Auth: auth})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(out_count)

	return

}
