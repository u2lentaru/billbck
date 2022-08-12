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

// HandleRequestKinds godoc
// @Summary List requestkinds
// @Description get requestkind list
// @Tags requestkinds
// @Produce json
// @Param page query int false "page number"
// @Param page_size query int false "page size"
// @Param requestkindname query string false "requestkindname search pattern"
// @Param ordering query string false "order by {id|requestkindname}"
// @Param desc query boolean false "descending order {true|false}"
// @Success 200 {object} models.RequestKind_count
// @Failure 500
// @Router /requestkinds [get]
func (s *APG) HandleRequestKinds(w http.ResponseWriter, r *http.Request) {
	gs := models.RequestKind{}
	ctx := context.Background()
	out_arr := []models.RequestKind{}

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
	gs1s, ok := query["requestkindname"]
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
	} else if ords[0] == "requestkindname" {
		ord = 2
	}

	dsc := false
	dscs, ok := query["desc"]
	if ok && len(dscs) > 0 {
		if !(dscs[0] == "0") {
			dsc = true
		}
	}

	rows, err := s.Dbpool.Query(ctx, "SELECT * from func_request_kinds_get($1,$2,$3,$4,$5);", pg, pgs, gs1, ord, dsc)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&gs.Id, &gs.RequestKindName)
		if err != nil {
			log.Println("failed to scan row:", err)
		}

		out_arr = append(out_arr, gs)
	}

	gsc := 0
	err = s.Dbpool.QueryRow(ctx, "SELECT * from func_request_kinds_cnt($1);", gs1).Scan(&gsc)

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	auth := models.Auth{Create: true, Read: true, Update: true, Delete: true}
	out_count, err := json.Marshal(models.RequestKind_count{Values: out_arr, Count: gsc, Auth: auth})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Write(out_count)

	return
}

// HandleAddRequestKind godoc
// @Summary Add requestkind
// @Description add requestkind
// @Tags requestkinds
// @Accept json
// @Produce  json
// @Param a body models.AddRequestKind true "New requestkind"
// @Success 200 {object} models.Json_id
// @Failure 500
// @Router /requestkinds_add [post]
func (s *APG) HandleAddRequestKind(w http.ResponseWriter, r *http.Request) {
	a := models.AddRequestKind{}
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
	err = s.Dbpool.QueryRow(context.Background(), "SELECT func_request_kinds_add($1);", a.RequestKindName).Scan(&ai)

	if err != nil {
		log.Println("Failed execute func_request_kinds_add: ", err)
	}

	output, err := json.Marshal(models.Json_id{Id: ai})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(output)

	return

}

// HandleUpdRequestKind godoc
// @Summary Update requestkind
// @Description update requestkind
// @Tags requestkinds
// @Accept json
// @Produce  json
// @Param u body models.RequestKind true "Update requestkind"
// @Success 200 {object} models.Json_id
// @Failure 500
// @Router /requestkinds_upd [post]
func (s *APG) HandleUpdRequestKind(w http.ResponseWriter, r *http.Request) {
	u := models.RequestKind{}
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
	err = s.Dbpool.QueryRow(context.Background(), "SELECT func_request_kinds_upd($1,$2);", u.Id, u.RequestKindName).Scan(&ui)

	if err != nil {
		log.Println("Failed execute func_request_kinds_upd: ", err)
	}

	output, err := json.Marshal(models.Json_id{Id: ui})

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(output)

	return
}

// HandleDelRequestKind godoc
// @Summary Delete requestkinds
// @Description delete requestkinds
// @Tags requestkinds
// @Accept json
// @Produce  json
// @Param d body models.Json_ids true "Delete requestkinds"
// @Success 200 {object} models.Json_ids
// @Failure 500
// @Router /requestkinds_del [post]
func (s *APG) HandleDelRequestKind(w http.ResponseWriter, r *http.Request) {
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
		err = s.Dbpool.QueryRow(context.Background(), "SELECT func_request_kinds_del($1);", id).Scan(&i)
		res = append(res, i)

		if err != nil {
			log.Println("Failed execute func_request_kinds_del: ", err)
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

// HandleGetRequestKind godoc
// @Summary Get requestkind
// @Description get requestkind
// @Tags requestkinds
// @Produce  json
// @Param id path int true "RequestKind by id"
// @Success 200 {object} models.RequestKind_count
// @Failure 500
// @Router /requestkinds/{id} [get]
func (s *APG) HandleGetRequestKind(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	i := vars["id"]
	g := models.RequestKind{}
	out_arr := []models.RequestKind{}

	err := s.Dbpool.QueryRow(context.Background(), "SELECT * from func_request_kind_get($1);", i).Scan(&g.Id, &g.RequestKindName)

	if err != nil && err != pgx.ErrNoRows {
		log.Println("Failed execute from func_request_kind_get: ", err)
	}

	out_arr = append(out_arr, g)
	auth := models.Auth{Create: true, Read: true, Update: true, Delete: true}

	// output, err := json.Marshal(g)
	out_count, err := json.Marshal(models.RequestKind_count{Values: out_arr, Count: 1, Auth: auth})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(out_count)

	return
}
