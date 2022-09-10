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

type ifDistributionZones interface {
	GetDistributionZones(context.Context, *pgxpool.Pool, int, int, string, int, bool) (models.DistributionZone_count, error)
}

type ifAddDistributionZone interface {
	AddDistributionZone(context.Context, *pgxpool.Pool) (int, error)
}

type ifUpdDistributionZone interface {
	UpdDistributionZone(context.Context, *pgxpool.Pool) (int, error)
}
type ifDelDistributionZone interface {
	DelDistributionZone(context.Context, *pgxpool.Pool, []int) ([]int, error)
}

type ifDistributionZone interface {
	GetDistributionZone(context.Context, *pgxpool.Pool, int) (models.DistributionZone_count, error)
}

// HandleDistributionZone godoc
// @Summary List distributionzone
// @Description get distributionzone list
// @Tags distributionzones
// @Produce json
// @Param page query int false "page number"
// @Param page_size query int false "page size"
// @Param distributionzonename query string false "distributionzonename search pattern"
// @Param ordering query string false "order by {id|distributionzonename}"
// @Param desc query boolean false "descending order {true|false}"
// @Success 200 {object} models.DistributionZone_count
// @Failure 500
// @Router /distributionzones [get]
func (s *APG) HandleDistributionZones(w http.ResponseWriter, r *http.Request) {
	var gs ifDistributionZones
	gs = models.NewDistributionZone()
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
	gs1s, ok := query["distributionzonename"]
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
	} else if ords[0] == "distributionzonename" {
		ord = 2
	}

	dsc := false
	dscs, ok := query["desc"]
	if ok && len(dscs) > 0 {
		if !(dscs[0] == "0") {
			dsc = true
		}
	}

	out_count, err := gs.GetDistributionZones(ctx, s.Dbpool, pg, pgs, gs1, ord, dsc)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	out_json, err := json.Marshal(out_count)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Write(out_json)

	return
}

// HandleAddDistributionZone godoc
// @Summary Add distributionzone
// @Description add distributionzone
// @Tags distributionzones
// @Accept json
// @Produce  json
// @Param a body models.AddDistributionZone true "New distributionzone"
// @Success 200 {object} models.Json_id
// @Failure 500
// @Router /distributionzones_add [post]
func (s *APG) HandleAddDistributionZone(w http.ResponseWriter, r *http.Request) {
	var a ifAddDistributionZone
	a = models.NewDistributionZone()
	ctx := context.Background()

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

	ai, err := a.AddDistributionZone(ctx, s.Dbpool)
	if err != nil {
		log.Println("Failed execute func_distribution_zones_add: ", err)
	}

	output, err := json.Marshal(models.Json_id{Id: ai})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(output)

	return
}

// HandleUpdDistributionZone godoc
// @Summary Update distributionzone
// @Description update distributionzone
// @Tags distributionzones
// @Accept json
// @Produce  json
// @Param u body models.DistributionZone true "Update distributionzone"
// @Success 200 {object} models.Json_id
// @Failure 500
// @Router /distributionzones_upd [post]
func (s *APG) HandleUpdDistributionZone(w http.ResponseWriter, r *http.Request) {
	u := models.DistributionZone{}
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
	err = s.Dbpool.QueryRow(context.Background(), "SELECT func_distribution_zones_upd($1,$2);", u.Id, u.DistributionZoneName).Scan(&ui)

	if err != nil {
		log.Println("Failed execute func_distribution_zones_upd: ", err)
	}

	output, err := json.Marshal(models.Json_id{Id: ui})

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(output)

	return
}

// HandleDelDistributionZone godoc
// @Summary Delete distributionzones
// @Description delete distributionzones
// @Tags distributionzones
// @Accept json
// @Produce  json
// @Param d body models.Json_ids true "Delete distributionzones"
// @Success 200 {object} models.Json_ids
// @Failure 500
// @Router /distributionzones_del [post]
func (s *APG) HandleDelDistributionZone(w http.ResponseWriter, r *http.Request) {
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
		err = s.Dbpool.QueryRow(context.Background(), "SELECT func_distribution_zones_del($1);", id).Scan(&i)
		res = append(res, i)

		if err != nil {
			log.Println("Failed execute func_distribution_zones_del: ", err)
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

// HandleGetDistributionZone godoc
// @Summary Get distributionzone
// @Description get distributionzone
// @Tags distributionzones
// @Produce  json
// @Param id path int true "DistributionZone by id"
// @Success 200 {object} models.DistributionZone_count
// @Failure 500
// @Router /distributionzones/{id} [get]
func (s *APG) HandleGetDistributionZone(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	i := vars["id"]
	g := models.DistributionZone{}
	out_arr := []models.DistributionZone{}

	err := s.Dbpool.QueryRow(context.Background(), "SELECT * from func_distribution_zone_get($1);", i).Scan(&g.Id, &g.DistributionZoneName)

	if err != nil && err != pgx.ErrNoRows {
		log.Println("Failed execute from func_distribution_zone_get: ", err)
	}

	out_arr = append(out_arr, g)
	auth := models.Auth{Create: true, Read: true, Update: true, Delete: true}

	out_count, err := json.Marshal(models.DistributionZone_count{Values: out_arr, Count: 1, Auth: auth})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(out_count)

	return

}

/*
// HandleDistributionZone godoc
// @Summary List distributionzone
// @Description get distributionzone list
// @Tags distributionzones
// @Produce json
// @Param page query int false "page number"
// @Param page_size query int false "page size"
// @Param distributionzonename query string false "distributionzonename search pattern"
// @Param ordering query string false "order by {id|distributionzonename}"
// @Param desc query boolean false "descending order {true|false}"
// @Success 200 {object} models.DistributionZone_count
// @Failure 500
// @Router /distributionzones [get]
func (s *APG) HandleDistributionZones(w http.ResponseWriter, r *http.Request) {
	gs := models.DistributionZone{}
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
	gs1s, ok := query["distributionzonename"]
	if ok && len(gs1s) > 0 {
		//case insensitive
		gs1 = strings.ToUpper(gs1s[0])
		//quotes
		re := regexp.MustCompile(`'`)
		gs1 = string(re.ReplaceAll([]byte(gs1), []byte("''")))
	}

	gsc := 0
	err := s.Dbpool.QueryRow(ctx, "SELECT * from func_distribution_zones_cnt($1);", gs1).Scan(&gsc)

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	out_arr := make([]models.DistributionZone, 0,
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
	} else if ords[0] == "distributionzonename" {
		ord = 2
	}

	dsc := false
	dscs, ok := query["desc"]
	if ok && len(dscs) > 0 {
		if !(dscs[0] == "0") {
			dsc = true
		}
	}

	rows, err := s.Dbpool.Query(ctx, "SELECT * from func_distribution_zones_get($1,$2,$3,$4,$5);", pg, pgs, gs1, ord, dsc)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&gs.Id, &gs.DistributionZoneName)
		if err != nil {
			log.Println("failed to scan row:", err)
		}

		out_arr = append(out_arr, gs)
	}

	auth := models.Auth{Create: true, Read: true, Update: true, Delete: true}
	out_count, err := json.Marshal(models.DistributionZone_count{Values: out_arr, Count: gsc, Auth: auth})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Write(out_count)

	return
}

// HandleAddDistributionZone godoc
// @Summary Add distributionzone
// @Description add distributionzone
// @Tags distributionzones
// @Accept json
// @Produce  json
// @Param a body models.AddDistributionZone true "New distributionzone"
// @Success 200 {object} models.Json_id
// @Failure 500
// @Router /distributionzones_add [post]
func (s *APG) HandleAddDistributionZone(w http.ResponseWriter, r *http.Request) {
	a := models.AddDistributionZone{}
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
	err = s.Dbpool.QueryRow(context.Background(), "SELECT func_distribution_zones_add($1);", a.DistributionZoneName).Scan(&ai)

	if err != nil {
		log.Println("Failed execute func_distribution_zones_add: ", err)
	}

	output, err := json.Marshal(models.Json_id{Id: ai})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(output)

	return
}

// HandleUpdDistributionZone godoc
// @Summary Update distributionzone
// @Description update distributionzone
// @Tags distributionzones
// @Accept json
// @Produce  json
// @Param u body models.DistributionZone true "Update distributionzone"
// @Success 200 {object} models.Json_id
// @Failure 500
// @Router /distributionzones_upd [post]
func (s *APG) HandleUpdDistributionZone(w http.ResponseWriter, r *http.Request) {
	u := models.DistributionZone{}
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
	err = s.Dbpool.QueryRow(context.Background(), "SELECT func_distribution_zones_upd($1,$2);", u.Id, u.DistributionZoneName).Scan(&ui)

	if err != nil {
		log.Println("Failed execute func_distribution_zones_upd: ", err)
	}

	output, err := json.Marshal(models.Json_id{Id: ui})

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(output)

	return
}

// HandleDelDistributionZone godoc
// @Summary Delete distributionzones
// @Description delete distributionzones
// @Tags distributionzones
// @Accept json
// @Produce  json
// @Param d body models.Json_ids true "Delete distributionzones"
// @Success 200 {object} models.Json_ids
// @Failure 500
// @Router /distributionzones_del [post]
func (s *APG) HandleDelDistributionZone(w http.ResponseWriter, r *http.Request) {
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
		err = s.Dbpool.QueryRow(context.Background(), "SELECT func_distribution_zones_del($1);", id).Scan(&i)
		res = append(res, i)

		if err != nil {
			log.Println("Failed execute func_distribution_zones_del: ", err)
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

// HandleGetDistributionZone godoc
// @Summary Get distributionzone
// @Description get distributionzone
// @Tags distributionzones
// @Produce  json
// @Param id path int true "DistributionZone by id"
// @Success 200 {object} models.DistributionZone_count
// @Failure 500
// @Router /distributionzones/{id} [get]
func (s *APG) HandleGetDistributionZone(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	i := vars["id"]
	g := models.DistributionZone{}
	out_arr := []models.DistributionZone{}

	err := s.Dbpool.QueryRow(context.Background(), "SELECT * from func_distribution_zone_get($1);", i).Scan(&g.Id, &g.DistributionZoneName)

	if err != nil && err != pgx.ErrNoRows {
		log.Println("Failed execute from func_distribution_zone_get: ", err)
	}

	out_arr = append(out_arr, g)
	auth := models.Auth{Create: true, Read: true, Update: true, Delete: true}

	// output, err := json.Marshal(g)
	out_count, err := json.Marshal(models.DistributionZone_count{Values: out_arr, Count: 1, Auth: auth})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(out_count)

	return

}
*/
