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

// HandleCableResistances godoc
// @Summary List cableresistances
// @Description get cableresistance list
// @Tags cableresistances
// @Produce  json
// @Param page query int false "page number"
// @Param page_size query int false "page size"
// @Param cableresistancename query string false "cableresistancename search pattern"
// @Param ordering query string false "order by {id|cableresistancename}"
// @Param desc query boolean false "descending order {true|false}"
// @Success 200 {object} models.CableResistance_count
// @Failure 500
// @Router /cableresistances [get]
func (s *APG) HandleCableResistances(w http.ResponseWriter, r *http.Request) {
	// start := time.Now()
	// gs := models.CableResistance{}
	// ctx := context.Background()

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
	gs1s, ok := query["cableresistancename"]
	if ok && len(gs1s) > 0 {
		//case insensitive
		gs1 = strings.ToUpper(gs1s[0])
		//quotes
		re := regexp.MustCompile(`'`)
		gs1 = string(re.ReplaceAll([]byte(gs1), []byte("''")))
	}

	// gsc := 0
	// err := s.Dbpool.QueryRow(ctx, "SELECT * from func_cable_resistances_cnt($1);", gs1).Scan(&gsc)

	// if err != nil {
	// 	http.Error(w, err.Error(), 500)
	// 	return
	// }

	// out_arr := make([]models.CableResistance, 0,
	// 	func() int {
	// 		if gsc < pgs {
	// 			return gsc
	// 		} else {
	// 			return pgs
	// 		}
	// 	}())

	ord := 1
	ords, ok := query["ordering"]
	if !ok || len(ords) == 0 {
		ord = 1
	} else if ords[0] == "cableresistancename" {
		ord = 2
	}

	dsc := false
	dscs, ok := query["desc"]
	if ok && len(dscs) > 0 {
		if !(dscs[0] == "0") {
			dsc = true
		}
	}

	// rows, err := s.Dbpool.Query(ctx, "SELECT * from func_cable_resistances_get($1,$2,$3,$4,$5);", pg, pgs, gs1, ord, dsc)
	// if err != nil {
	// 	http.Error(w, err.Error(), 500)
	// 	return
	// }

	// defer rows.Close()

	// for rows.Next() {
	// 	err = rows.Scan(&gs.Id, &gs.CableResistanceName, &gs.Resistance, &gs.MaterialType)
	// 	if err != nil {
	// 		log.Println("failed to scan row:", err)
	// 	}

	// 	out_arr = append(out_arr, gs)
	// }

	// auth := models.Auth{Create: true, Read: true, Update: true, Delete: true}
	// out_count, err := json.Marshal(models.CableResistance_count{Values: out_arr, Count: gsc, Auth: auth})
	// if err != nil {
	// 	http.Error(w, err.Error(), 500)
	// 	return
	// }

	type ifCableResistance_count interface {
		GetCableResistances(*pgxpool.Pool, int, int, string, int, bool) error
	}

	var out_count ifCableResistance_count

	out_count = models.NewCableResistance_count()

	out_count.GetCableResistances(s.Dbpool, pg, pgs, gs1, ord, dsc)

	w_out_count, err := json.Marshal(out_count)

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Write(w_out_count)

	// log.Printf("Work time %s", time.Since(start))

	return

}

// HandleAddCableResistance godoc
// @Summary Add cableresistance
// @Description add cableresistance
// @Tags cableresistances
// @Accept json
// @Produce  json
// @Param a body models.AddCableResistance true "New cableresistance"
// @Success 200 {object} models.Json_id
// @Failure 500
// @Router /cableresistances_add [post]
func (s *APG) HandleAddCableResistance(w http.ResponseWriter, r *http.Request) {
	a := models.AddCableResistance{}
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
	err = s.Dbpool.QueryRow(context.Background(), "SELECT func_cable_resistances_add($1,$2,$3);", a.CableResistanceName,
		a.Resistance, a.MaterialType).Scan(&ai)

	if err != nil {
		log.Println("Failed execute func_cable_resistances_add: ", err)
	}

	output, err := json.Marshal(models.Json_id{Id: ai})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(output)

	return

}

// HandleUpdCableResistance godoc
// @Summary Update cableresistance
// @Description update cableresistance
// @Tags cableresistances
// @Accept json
// @Produce  json
// @Param u body models.CableResistance true "Update cableresistance"
// @Success 200 {object} models.Json_id
// @Failure 500
// @Router /cableresistances_upd [post]
func (s *APG) HandleUpdCableResistance(w http.ResponseWriter, r *http.Request) {
	u := models.CableResistance{}
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
	err = s.Dbpool.QueryRow(context.Background(), "SELECT func_cable_resistances_upd($1,$2,$3,$4);", u.Id, u.CableResistanceName,
		u.Resistance, u.MaterialType).Scan(&ui)

	if err != nil {
		log.Println("Failed execute func_cable_resistances_upd: ", err)
	}

	output, err := json.Marshal(models.Json_id{Id: ui})

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(output)

	return

}

// HandleDelCableResistance godoc
// @Summary Delete cableresistances
// @Description delete cableresistances
// @Tags cableresistances
// @Accept json
// @Produce  json
// @Param d body models.Json_ids true "Delete cableresistances"
// @Success 200 {object} models.Json_ids
// @Failure 500
// @Router /cableresistances_del [post]
func (s *APG) HandleDelCableResistance(w http.ResponseWriter, r *http.Request) {
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
		err = s.Dbpool.QueryRow(context.Background(), "SELECT func_cable_resistances_del($1);", id).Scan(&i)
		res = append(res, i)

		if err != nil {
			log.Println("Failed execute func_cable_resistances_del: ", err)
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

// HandleGetCableResistance godoc
// @Summary Get cableresistance
// @Description get cableresistance
// @Tags cableresistances
// @Produce  json
// @Param id path int true "CableResistance by id"
// @Success 200 {object} models.CableResistance_count
// @Failure 500
// @Router /cableresistances/{id} [get]
func (s *APG) HandleGetCableResistance(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	i := vars["id"]
	g := models.CableResistance{}
	out_arr := []models.CableResistance{}

	err := s.Dbpool.QueryRow(context.Background(), "SELECT * from func_cable_resistance_get($1);", i).Scan(&g.Id,
		&g.CableResistanceName, &g.Resistance, &g.MaterialType)

	if err != nil && err != pgx.ErrNoRows {
		log.Println("Failed execute from func_cable_resistance_get: ", err)
	}

	out_arr = append(out_arr, g)
	auth := models.Auth{Create: true, Read: true, Update: true, Delete: true}

	// output, err := json.Marshal(g)
	out_count, err := json.Marshal(models.CableResistance_count{Values: out_arr, Count: 1, Auth: auth})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(out_count)

	return
}
