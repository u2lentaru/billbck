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

// HandleSealColours godoc
// @Summary List sealcolours
// @Description get sealcolour list
// @Tags sealcolours
// @Produce json
// @Param page query int false "page number"
// @Param page_size query int false "page size"
// @Param sealcolourname query string false "sealcolourname search pattern"
// @Param ordering query string false "order by {id|sealcolourname}"
// @Param desc query boolean false "descending order {true|false}"
// @Success 200 {object} models.SealColour_count
// @Failure 500
// @Router /sealcolours [get]
func (s *APG) HandleSealColours(w http.ResponseWriter, r *http.Request) {
	gs := models.SealColour{}
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
	gs1s, ok := query["sealcolourname"]
	if ok && len(gs1s) > 0 {
		//case insensitive
		gs1 = strings.ToUpper(gs1s[0])
		//quotes
		re := regexp.MustCompile(`'`)
		gs1 = string(re.ReplaceAll([]byte(gs1), []byte("''")))
	}

	gsc := 0
	err := s.Dbpool.QueryRow(ctx, "SELECT * from func_seal_colours_cnt($1);", gs1).Scan(&gsc)

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	out_arr := make([]models.SealColour, 0,
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
	} else if ords[0] == "sealcolourname" {
		ord = 2
	}

	dsc := false
	dscs, ok := query["desc"]
	if ok && len(dscs) > 0 {
		if !(dscs[0] == "0") {
			dsc = true
		}
	}

	rows, err := s.Dbpool.Query(ctx, "SELECT * from func_seal_colours_get($1,$2,$3,$4,$5);", pg, pgs, gs1, ord, dsc)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&gs.Id, &gs.SealColourName)
		if err != nil {
			log.Println("failed to scan row:", err)
		}

		out_arr = append(out_arr, gs)
	}

	auth := models.Auth{Create: true, Read: true, Update: true, Delete: true}
	out_count, err := json.Marshal(models.SealColour_count{Values: out_arr, Count: gsc, Auth: auth})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Write(out_count)

	return
}

// HandleAddSealColour godoc
// @Summary Add sealcolour
// @Description add sealcolour
// @Tags sealcolours
// @Accept json
// @Produce  json
// @Param a body models.AddSealColour true "New sealcolour"
// @Success 200 {object} models.Json_id
// @Failure 500
// @Router /sealcolours_add [post]
func (s *APG) HandleAddSealColour(w http.ResponseWriter, r *http.Request) {
	a := models.AddSealColour{}
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
	err = s.Dbpool.QueryRow(context.Background(), "SELECT func_seal_colours_add($1);", a.SealColourName).Scan(&ai)

	if err != nil {
		log.Println("Failed execute func_seal_colours_add: ", err)
	}

	output, err := json.Marshal(models.Json_id{Id: ai})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(output)

	return

}

// HandleUpdSealColour godoc
// @Summary Update sealcolour
// @Description update sealcolour
// @Tags sealcolours
// @Accept json
// @Produce  json
// @Param u body models.SealColour true "Update sealcolour"
// @Success 200 {object} models.Json_id
// @Failure 500
// @Router /sealcolours_upd [post]
func (s *APG) HandleUpdSealColour(w http.ResponseWriter, r *http.Request) {
	u := models.SealColour{}
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
	err = s.Dbpool.QueryRow(context.Background(), "SELECT func_seal_colours_upd($1,$2);", u.Id, u.SealColourName).Scan(&ui)

	if err != nil {
		log.Println("Failed execute func_seal_colours_upd: ", err)
	}

	output, err := json.Marshal(models.Json_id{Id: ui})

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(output)

	return
}

// HandleDelSealColour godoc
// @Summary Delete sealcolours
// @Description delete sealcolours
// @Tags sealcolours
// @Accept json
// @Produce  json
// @Param d body models.Json_ids true "Delete sealcolours"
// @Success 200 {object} models.Json_ids
// @Failure 500
// @Router /sealcolours_del [post]
func (s *APG) HandleDelSealColour(w http.ResponseWriter, r *http.Request) {
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
		err = s.Dbpool.QueryRow(context.Background(), "SELECT func_seal_colours_del($1);", id).Scan(&i)
		res = append(res, i)

		if err != nil {
			log.Println("Failed execute func_seal_colours_del: ", err)
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

// HandleGetSealColour godoc
// @Summary Get sealcolour
// @Description get sealcolour
// @Tags sealcolours
// @Produce  json
// @Param id path int true "SealColour by id"
// @Success 200 {object} models.SealColour_count
// @Failure 500
// @Router /sealcolours/{id} [get]
func (s *APG) HandleGetSealColour(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	i := vars["id"]
	g := models.SealColour{}
	out_arr := []models.SealColour{}

	err := s.Dbpool.QueryRow(context.Background(), "SELECT * from func_seal_colour_get($1);", i).Scan(&g.Id, &g.SealColourName)

	if err != nil && err != pgx.ErrNoRows {
		log.Println("Failed execute from func_seal_colour_get: ", err)
	}

	out_arr = append(out_arr, g)
	auth := models.Auth{Create: true, Read: true, Update: true, Delete: true}

	// output, err := json.Marshal(g)
	out_count, err := json.Marshal(models.SealColour_count{Values: out_arr, Count: 1, Auth: auth})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(out_count)

	return

}
