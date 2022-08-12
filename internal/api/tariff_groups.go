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

// HandleTariffGroups godoc
// @Summary List tariffgroups
// @Description get tariffgroup list
// @Tags tariffgroups
// @Produce  json
// @Param page query int false "page number"
// @Param page_size query int false "page size"
// @Param tariffgroupname query string false "tariffgroupname search pattern"
// @Param ordering query string false "order by {id|tariffgroupname}"
// @Param desc query boolean false "descending order {true|false}"
// @Success 200 {object} models.TariffGroup_count
// @Failure 500
// @Router /tariffgroups [get]
func (s *APG) HandleTariffGroups(w http.ResponseWriter, r *http.Request) {
	gs := models.TariffGroup{}
	ctx := context.Background()
	out_arr := []models.TariffGroup{}

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
	gs1s, ok := query["tariffgroupname"]
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
	} else if ords[0] == "tariffgroupname" {
		ord = 2
	}

	dsc := false
	dscs, ok := query["desc"]
	if ok && len(dscs) > 0 {
		if !(dscs[0] == "0") {
			dsc = true
		}
	}

	rows, err := s.Dbpool.Query(ctx, "SELECT * from func_tariff_groups_get($1,$2,$3,$4,$5);", pg, pgs, gs1, ord, dsc)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&gs.Id, &gs.TariffGroupName)
		if err != nil {
			log.Println("failed to scan row:", err)
		}

		out_arr = append(out_arr, gs)
	}

	gsc := 0
	err = s.Dbpool.QueryRow(ctx, "SELECT * from func_tariff_groups_cnt($1);", gs1).Scan(&gsc)

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	auth := models.Auth{Create: true, Read: true, Update: true, Delete: true}
	out_count, err := json.Marshal(models.TariffGroup_count{Values: out_arr, Count: gsc, Auth: auth})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Write(out_count)

	return

}

// HandleAddTariffGroup godoc
// @Summary Add tariffgroup
// @Description add tariffgroup
// @Tags tariffgroups
// @Accept json
// @Produce  json
// @Param a body models.AddTariffGroup true "New tariffgroup"
// @Success 200 {object} models.Json_id
// @Failure 500
// @Router /tariffgroups_add [post]
func (s *APG) HandleAddTariffGroup(w http.ResponseWriter, r *http.Request) {
	a := models.AddTariffGroup{}
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
	err = s.Dbpool.QueryRow(context.Background(), "SELECT func_tariff_groups_add($1);", a.TariffGroupName).Scan(&ai)

	if err != nil {
		log.Println("Failed execute func_tariff_groups_add: ", err)
	}

	output, err := json.Marshal(models.Json_id{Id: ai})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(output)

	return

}

// HandleUpdTariffGroup godoc
// @Summary Update tariffgroup
// @Description update tariffgroup
// @Tags tariffgroups
// @Accept json
// @Produce  json
// @Param u body models.TariffGroup true "Update tariffgroup"
// @Success 200 {object} models.Json_id
// @Failure 500
// @Router /tariffgroups_upd [post]
func (s *APG) HandleUpdTariffGroup(w http.ResponseWriter, r *http.Request) {
	u := models.TariffGroup{}
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
	err = s.Dbpool.QueryRow(context.Background(), "SELECT func_tariff_groups_upd($1,$2);", u.Id, u.TariffGroupName).Scan(&ui)

	if err != nil {
		log.Println("Failed execute func_tariff_groups_upd: ", err)
	}

	output, err := json.Marshal(models.Json_id{Id: ui})

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(output)

	return
}

// HandleDelTariffGroup godoc
// @Summary Delete tariffgroups
// @Description delete tariffgroups
// @Tags tariffgroups
// @Accept json
// @Produce  json
// @Param d body models.Json_ids true "Delete tariffgroups"
// @Success 200 {object} models.Json_ids
// @Failure 500
// @Router /tariffgroups_del [post]
func (s *APG) HandleDelTariffGroup(w http.ResponseWriter, r *http.Request) {
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
		err = s.Dbpool.QueryRow(context.Background(), "SELECT func_tariff_groups_del($1);", id).Scan(&i)
		res = append(res, i)

		if err != nil {
			log.Println("Failed execute func_tariff_groups_del: ", err)
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

// HandleGetTariffGroup godoc
// @Summary Get tariffgroup
// @Description get tariffgroup
// @Tags tariffgroups
// @Produce  json
// @Param id path int true "TariffGroup by id"
// @Success 200 {object} models.TariffGroup_count
// @Failure 500
// @Router /tariffgroups/{id} [get]
func (s *APG) HandleGetTariffGroup(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	i := vars["id"]
	g := models.TariffGroup{}
	out_arr := []models.TariffGroup{}

	err := s.Dbpool.QueryRow(context.Background(), "SELECT * from func_tariff_group_get($1);", i).Scan(&g.Id, &g.TariffGroupName)

	if err != nil && err != pgx.ErrNoRows {
		log.Println("Failed execute from func_tariff_group_get: ", err)
	}

	out_arr = append(out_arr, g)
	auth := models.Auth{Create: true, Read: true, Update: true, Delete: true}

	// output, err := json.Marshal(g)
	out_count, err := json.Marshal(models.TariffGroup_count{Values: out_arr, Count: 1, Auth: auth})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(out_count)

	return

}
