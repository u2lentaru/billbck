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

// HandleConclusions godoc
// @Summary List conclusions
// @Description get conclusions list
// @Tags conclusions
// @Produce  json
// @Param page query int false "page number"
// @Param page_size query int false "page size"
// @Param conclusionname query string false "conclusionname search pattern"
// @Param ordering query string false "order by {id|conclusionname}"
// @Param desc query boolean false "descending order {true|false}"
// @Success 200 {array} models.Conclusion_count
// @Failure 500
// @Router /conclusions [get]
func (s *APG) HandleConclusions(w http.ResponseWriter, r *http.Request) {
	gs := models.Conclusion{}
	ctx := context.Background()
	out_arr := []models.Conclusion{}

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
	gs1s, ok := query["conclusionname"]
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
	} else if ords[0] == "conclusionname" {
		ord = 2
	}

	dsc := false
	dscs, ok := query["desc"]
	if ok && len(dscs) > 0 {
		if !(dscs[0] == "0") {
			dsc = true
		}
	}

	rows, err := s.Dbpool.Query(ctx, "SELECT * from func_conclusions_get($1,$2,$3,$4,$5);", pg, pgs, gs1, ord, dsc)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&gs.Id, &gs.ConclusionName)
		if err != nil {
			log.Println("failed to scan row:", err)
		}

		out_arr = append(out_arr, gs)
	}

	gsc := 0
	err = s.Dbpool.QueryRow(ctx, "SELECT * from func_conclusions_cnt($1);", gs1).Scan(&gsc)

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	auth := models.Auth{Create: true, Read: true, Update: true, Delete: true}
	out_count, err := json.Marshal(models.Conclusion_count{Values: out_arr, Count: gsc, Auth: auth})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Write(out_count)

	return
}

// HandleAddConclusion godoc
// @Summary Add conclusion
// @Description add conclusion
// @Tags conclusions
// @Accept json
// @Produce  json
// @Param a body models.AddConclusion true "New conclusion"
// @Success 200 {array} models.Json_id
// @Failure 500
// @Router /conclusions_add [post]
func (s *APG) HandleAddConclusion(w http.ResponseWriter, r *http.Request) {
	a := models.AddConclusion{}
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
	err = s.Dbpool.QueryRow(context.Background(), "SELECT func_conclusions_add($1);", a.ConclusionName).Scan(&ai)

	if err != nil {
		log.Println("Failed execute func_conclusions_add: ", err)
	}

	output, err := json.Marshal(models.Json_id{Id: ai})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(output)

	return

}

// HandleUpdConclusion godoc
// @Summary Update conclusion
// @Description update conclusion
// @Tags conclusions
// @Accept json
// @Produce  json
// @Param u body models.Conclusion true "Update conclusion"
// @Success 200 {array} models.Json_id
// @Failure 500
// @Router /conclusions_upd [post]
func (s *APG) HandleUpdConclusion(w http.ResponseWriter, r *http.Request) {
	u := models.Conclusion{}
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
	err = s.Dbpool.QueryRow(context.Background(), "SELECT func_conclusions_upd($1,$2);", u.Id, u.ConclusionName).Scan(&ui)

	if err != nil {
		log.Println("Failed execute func_conclusions_upd: ", err)
	}

	output, err := json.Marshal(models.Json_id{Id: ui})

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(output)

	return
}

// HandleDelConclusion godoc
// @Summary Delete conclusions
// @Description delete conclusions
// @Tags conclusions
// @Accept json
// @Produce  json
// @Param d body models.Json_ids true "Delete conclusions"
// @Success 200 {array} models.Json_ids
// @Failure 500
// @Router /conclusions_del [post]
func (s *APG) HandleDelConclusion(w http.ResponseWriter, r *http.Request) {
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
		err = s.Dbpool.QueryRow(context.Background(), "SELECT func_conclusions_del($1);", id).Scan(&i)
		res = append(res, i)

		if err != nil {
			log.Println("Failed execute func_conclusions_del: ", err)
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

// HandleGetConclusion godoc
// @Summary Get conclusion
// @Description get conclusion
// @Tags conclusions
// @Produce  json
// @Param id path int true "Conclusion by id"
// @Success 200 {array} models.Conclusion_count
// @Failure 500
// @Router /conclusions/{id} [get]
func (s *APG) HandleGetConclusion(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	i := vars["id"]
	g := models.Conclusion{}

	out_arr := []models.Conclusion{}

	err := s.Dbpool.QueryRow(context.Background(), "SELECT * from func_conclusion_get($1);", i).Scan(&g.Id, &g.ConclusionName)

	if err != nil && err != pgx.ErrNoRows {
		log.Println("Failed execute from func_conclusion_get: ", err)
	}

	out_arr = append(out_arr, g)
	auth := models.Auth{Create: true, Read: true, Update: true, Delete: true}

	// output, err := json.Marshal(g)
	out_count, err := json.Marshal(models.Conclusion_count{Values: out_arr, Count: 1, Auth: auth})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(out_count)

	return

}
