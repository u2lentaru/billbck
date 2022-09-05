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

// HandleStaff godoc
// @Summary List staff
// @Description get staff list
// @Tags staff
// @Produce  json
// @Param page query int false "page number"
// @Param page_size query int false "page size"
// @Param staffname query string false "staffname search pattern"
// @Param ordering query string false "order by {id|staffname}"
// @Param desc query boolean false "descending order {true|false}"
// @Success 200 {object} models.Staff_count
// @Failure 500
// @Router /staff [get]
func (s *APG) HandleStaff(w http.ResponseWriter, r *http.Request) {
	gs := models.Staff{}
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
	gs1s, ok := query["staffname"]
	if ok && len(gs1s) > 0 {
		//case insensitive
		gs1 = strings.ToUpper(gs1s[0])
		//quotes
		re := regexp.MustCompile(`'`)
		gs1 = string(re.ReplaceAll([]byte(gs1), []byte("''")))
	}

	gsc := 0
	err := s.Dbpool.QueryRow(ctx, "SELECT * from func_staff_cnt($1);", gs1).Scan(&gsc)

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	out_arr := make([]models.Staff, 0,
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
	} else if ords[0] == "staffname" {
		ord = 2
	}

	dsc := false
	dscs, ok := query["desc"]
	if ok && len(dscs) > 0 {
		if !(dscs[0] == "0") {
			dsc = true
		}
	}

	rows, err := s.Dbpool.Query(ctx, "SELECT * from func_staff_get($1,$2,$3,$4,$5);", pg, pgs, gs1, ord, dsc)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&gs.Id, &gs.StaffName, &gs.OrgInfo.Id, &gs.Phone, &gs.Notes, &gs.OrgInfo.OIName)
		if err != nil {
			log.Println("failed to scan row:", err)
		}

		out_arr = append(out_arr, gs)
	}

	auth := models.Auth{Create: true, Read: true, Update: true, Delete: true}
	out_count, err := json.Marshal(models.Staff_count{Values: out_arr, Count: gsc, Auth: auth})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Write(out_count)

	return
}

// HandleAddStaff godoc
// @Summary Add staff
// @Description add staff
// @Tags staff
// @Accept json
// @Produce  json
// @Param a body models.AddStaff true "New staff"
// @Success 200 {object} models.Json_id
// @Failure 500
// @Router /staff_add [post]
func (s *APG) HandleAddStaff(w http.ResponseWriter, r *http.Request) {
	a := models.AddStaff{}
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
	err = s.Dbpool.QueryRow(context.Background(), "SELECT func_staff_add($1,$2,$3,$4);", a.StaffName, a.OrgInfo.Id, a.Phone, a.Notes).Scan(&ai)

	if err != nil {
		log.Println("Failed execute func_staff_add: ", err)
	}

	output, err := json.Marshal(models.Json_id{Id: ai})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(output)

	return

}

// HandleUpdStaff godoc
// @Summary Update staff
// @Description update staff
// @Tags staff
// @Accept json
// @Produce  json
// @Param u body models.Staff true "Update staff"
// @Success 200 {object} models.Json_id
// @Failure 500
// @Router /staff_upd [post]
func (s *APG) HandleUpdStaff(w http.ResponseWriter, r *http.Request) {
	u := models.Staff{}
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
	err = s.Dbpool.QueryRow(context.Background(), "SELECT func_staff_upd($1,$2,$3,$4,$5);", u.Id, u.StaffName, u.OrgInfo.Id, u.Phone,
		u.Notes).Scan(&ui)

	if err != nil {
		log.Println("Failed execute func_staff_upd: ", err)
	}

	output, err := json.Marshal(models.Json_id{Id: ui})

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(output)

	return
}

// HandleDelStaff godoc
// @Summary Delete staff
// @Description delete staff
// @Tags staff
// @Accept json
// @Produce  json
// @Param d body models.Json_ids true "Delete staff"
// @Success 200 {object} models.Json_ids
// @Failure 500
// @Router /staff_del [post]
func (s *APG) HandleDelStaff(w http.ResponseWriter, r *http.Request) {
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
		err = s.Dbpool.QueryRow(context.Background(), "SELECT func_staff_del($1);", id).Scan(&i)
		res = append(res, i)

		if err != nil {
			log.Println("Failed execute func_staff_del: ", err)
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

// HandleGetStaff godoc
// @Summary Get staff
// @Description get staff
// @Tags staff
// @Produce  json
// @Param id path int true "Staff by id"
// @Success 200 {object} models.Staff_count
// @Failure 500
// @Router /staff/{id} [get]
func (s *APG) HandleGetStaff(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	i := vars["id"]
	g := models.Staff{}
	out_arr := []models.Staff{}

	err := s.Dbpool.QueryRow(context.Background(), "SELECT * from func_staff_getbyid($1);", i).Scan(&g.Id, &g.StaffName, &g.OrgInfo.Id,
		&g.Phone, &g.Notes, &g.OrgInfo.OIName)

	if err != nil && err != pgx.ErrNoRows {
		log.Println("Failed execute from func_staff_getbyid: ", err)
	}

	out_arr = append(out_arr, g)
	auth := models.Auth{Create: true, Read: true, Update: true, Delete: true}

	// output, err := json.Marshal(g)
	out_count, err := json.Marshal(models.Staff_count{Values: out_arr, Count: 1, Auth: auth})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(out_count)

	return
}
