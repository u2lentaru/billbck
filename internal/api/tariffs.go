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

// HandleTariffs godoc
// @Summary List tariffs
// @Description get tariff list
// @Tags tariffs
// @Produce  json
// @Param page query int false "page number"
// @Param page_size query int false "page size"
// @Param tariffname query string false "tariffname search pattern"
// @Param ordering query string false "order by {id|tariffname}"
// @Param desc query boolean false "descending order {true|false}"
// @Success 200 {object} models.Tariff_count
// @Failure 500
// @Router /tariffs [get]
func (s *APG) HandleTariffs(w http.ResponseWriter, r *http.Request) {
	gs := models.Tariff{}
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
	gs1s, ok := query["tariffname"]
	if ok && len(gs1s) > 0 {
		//case insensitive
		gs1 = strings.ToUpper(gs1s[0])
		//quotes
		re := regexp.MustCompile(`'`)
		gs1 = string(re.ReplaceAll([]byte(gs1), []byte("''")))
	}

	gsc := 0
	err := s.Dbpool.QueryRow(ctx, "SELECT * from func_tariffs_cnt($1);", gs1).Scan(&gsc)

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	out_arr := make([]models.Tariff, 0,
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
	} else if ords[0] == "tariffname" {
		ord = 2
	}

	dsc := false
	dscs, ok := query["desc"]
	if ok && len(dscs) > 0 {
		if !(dscs[0] == "0") {
			dsc = true
		}
	}

	rows, err := s.Dbpool.Query(ctx, "SELECT * from func_tariffs_get($1,$2,$3,$4,$5);", pg, pgs, gs1, ord, dsc)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&gs.Id, &gs.TariffName, &gs.TariffGroup.Id, &gs.Norma, &gs.Tariff, &gs.Startdate, &gs.Enddate, &gs.TariffGroup.TariffGroupName)
		if err != nil {
			log.Println("failed to scan row:", err)
		}

		//change null value to ""
		// if gs.Enddate == nil {
		// 	s := ""
		// 	gs.Enddate = &s
		// }

		out_arr = append(out_arr, gs)
	}

	auth := models.Auth{Create: true, Read: true, Update: true, Delete: true}
	out_count, err := json.Marshal(models.Tariff_count{Values: out_arr, Count: gsc, Auth: auth})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Write(out_count)

	return

}

// HandleAddTariff godoc
// @Summary Add tariff
// @Description add tariff
// @Tags tariffs
// @Accept json
// @Produce  json
// @Param a body models.AddTariff true "New tariff"
// @Success 200 {object} models.Json_id
// @Failure 500
// @Router /tariffs_add [post]
func (s *APG) HandleAddTariff(w http.ResponseWriter, r *http.Request) {
	a := models.AddTariff{}
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
	err = s.Dbpool.QueryRow(context.Background(), "SELECT func_tariffs_add($1,$2,$3,$4,$5);", a.TariffName, a.TariffGroup.Id, a.Norma, a.Tariff, a.Startdate).Scan(&ai)

	if err != nil {
		log.Println("Failed execute func_tariffs_add: ", err)
	}

	output, err := json.Marshal(models.Json_id{Id: ai})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(output)

	return

}

// HandleUpdTariff godoc
// @Summary Update tariff
// @Description update tariff
// @Tags tariffs
// @Accept json
// @Produce  json
// @Param u body models.Tariff true "Update tariff"
// @Success 200 {object} models.Json_id
// @Failure 500
// @Router /tariffs_upd [post]
func (s *APG) HandleUpdTariff(w http.ResponseWriter, r *http.Request) {
	u := models.Tariff{}
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
	if u.Enddate == nil {
		s := ""
		u.Enddate = &s
	}
	err = s.Dbpool.QueryRow(context.Background(), "SELECT func_tariffs_upd($1,$2,$3,$4,$5,$6,$7);", u.Id, u.TariffName, u.TariffGroup.Id, u.Norma, u.Tariff, u.Startdate, u.Enddate).Scan(&ui)

	if err != nil {
		log.Println("Failed execute func_tariffs_upd: ", err)
	}

	output, err := json.Marshal(models.Json_id{Id: ui})

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(output)

	return
}

// HandleDelTariff godoc
// @Summary Delete tariffs
// @Description delete tariffs
// @Tags tariffs
// @Accept json
// @Produce  json
// @Param d body models.Json_ids true "Delete tariffs"
// @Success 200 {object} models.Json_ids
// @Failure 500
// @Router /tariffs_del [post]
func (s *APG) HandleDelTariff(w http.ResponseWriter, r *http.Request) {
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
		err = s.Dbpool.QueryRow(context.Background(), "SELECT func_tariffs_del($1);", id).Scan(&i)
		res = append(res, i)

		if err != nil {
			log.Println("Failed execute func_tariffs_del: ", err)
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

// HandleGetTariff godoc
// @Summary Get tariff
// @Description get tariff
// @Tags tariffs
// @Produce  json
// @Param id path int true "Tariff by id"
// @Success 200 {object} models.Tariff_count
// @Failure 500
// @Router /tariffs/{id} [get]
func (s *APG) HandleGetTariff(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	i := vars["id"]
	g := models.Tariff{}
	out_arr := []models.Tariff{}

	err := s.Dbpool.QueryRow(context.Background(), "SELECT * from func_tariff_get($1);", i).Scan(&g.Id, &g.TariffName, &g.TariffGroup.Id, &g.Norma, &g.Tariff, &g.Startdate, &g.Enddate, &g.TariffGroup.TariffGroupName)

	if err != nil && err != pgx.ErrNoRows {
		log.Println("Failed execute from func_tariff_get: ", err)
	}

	out_arr = append(out_arr, g)
	auth := models.Auth{Create: true, Read: true, Update: true, Delete: true}

	// output, err := json.Marshal(g)
	out_count, err := json.Marshal(models.Tariff_count{Values: out_arr, Count: 1, Auth: auth})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(out_count)

	return

}
