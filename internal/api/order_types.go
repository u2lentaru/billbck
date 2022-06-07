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

// HandleOrderTypes godoc
// @Summary List ordertypes
// @Description get ordertype list
// @Tags ordertypes
// @Produce json
// @Param page query int false "page number"
// @Param page_size query int false "page size"
// @Param ordertypename query string false "ordertypename search pattern"
// @Param ordering query string false "order by {id|ordertypename}"
// @Param desc query boolean false "descending order {true|false}"
// @Success 200 {array} models.OrderType_count
// @Failure 500
// @Router /ordertypes [get]
func (s *APG) HandleOrderTypes(w http.ResponseWriter, r *http.Request) {
	gs := models.OrderType{}
	ctx := context.Background()
	out_arr := []models.OrderType{}

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
	gs1s, ok := query["ordertypename"]
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
	} else if ords[0] == "ordertypename" {
		ord = 2
	}

	dsc := false
	dscs, ok := query["desc"]
	if ok && len(dscs) > 0 {
		if !(dscs[0] == "0") {
			dsc = true
		}
	}

	rows, err := s.Dbpool.Query(ctx, "SELECT * from func_order_types_get($1,$2,$3,$4,$5);", pg, pgs, gs1, ord, dsc)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&gs.Id, &gs.OrderTypeName)
		if err != nil {
			log.Println("failed to scan row:", err)
		}

		out_arr = append(out_arr, gs)
	}

	gsc := 0
	err = s.Dbpool.QueryRow(ctx, "SELECT * from func_order_types_cnt($1);", gs1).Scan(&gsc)

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	auth := models.Auth{Create: true, Read: true, Update: true, Delete: true}
	out_count, err := json.Marshal(models.OrderType_count{Values: out_arr, Count: gsc, Auth: auth})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Write(out_count)

	return
}

// HandleAddOrderType godoc
// @Summary Add ordertype
// @Description add ordertype
// @Tags ordertypes
// @Accept json
// @Produce  json
// @Param a body models.AddOrderType true "New ordertype"
// @Success 200 {array} models.Json_id
// @Failure 500
// @Router /ordertypes_add [post]
func (s *APG) HandleAddOrderType(w http.ResponseWriter, r *http.Request) {
	a := models.AddOrderType{}
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
	err = s.Dbpool.QueryRow(context.Background(), "SELECT func_order_types_add($1);", a.OrderTypeName).Scan(&ai)

	if err != nil {
		log.Println("Failed execute func_order_types_add: ", err)
	}

	output, err := json.Marshal(models.Json_id{Id: ai})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(output)

	return

}

// HandleUpdOrderType godoc
// @Summary Update ordertype
// @Description update ordertype
// @Tags ordertypes
// @Accept json
// @Produce  json
// @Param u body models.OrderType true "Update ordertype"
// @Success 200 {array} models.Json_id
// @Failure 500
// @Router /ordertypes_upd [post]
func (s *APG) HandleUpdOrderType(w http.ResponseWriter, r *http.Request) {
	u := models.OrderType{}
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
	err = s.Dbpool.QueryRow(context.Background(), "SELECT func_order_types_upd($1,$2);", u.Id, u.OrderTypeName).Scan(&ui)

	if err != nil {
		log.Println("Failed execute func_order_types_upd: ", err)
	}

	output, err := json.Marshal(models.Json_id{Id: ui})

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(output)

	return
}

// HandleDelOrderType godoc
// @Summary Delete ordertypes
// @Description delete ordertypes
// @Tags ordertypes
// @Accept json
// @Produce  json
// @Param d body models.Json_ids true "Delete ordertypes"
// @Success 200 {array} models.Json_ids
// @Failure 500
// @Router /ordertypes_del [post]
func (s *APG) HandleDelOrderType(w http.ResponseWriter, r *http.Request) {
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
		err = s.Dbpool.QueryRow(context.Background(), "SELECT func_order_types_del($1);", id).Scan(&i)
		res = append(res, i)

		if err != nil {
			log.Println("Failed execute func_order_types_del: ", err)
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

// HandleGetOrderType godoc
// @Summary Get ordertype
// @Description get ordertype
// @Tags ordertypes
// @Produce  json
// @Param id path int true "OrderType by id"
// @Success 200 {array} models.OrderType_count
// @Failure 500
// @Router /ordertypes/{id} [get]
func (s *APG) HandleGetOrderType(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	i := vars["id"]
	g := models.OrderType{}
	out_arr := []models.OrderType{}

	err := s.Dbpool.QueryRow(context.Background(), "SELECT * from func_order_type_get($1);", i).Scan(&g.Id, &g.OrderTypeName)

	if err != nil && err != pgx.ErrNoRows {
		log.Println("Failed execute from func_order_type_get: ", err)
	}

	out_arr = append(out_arr, g)
	auth := models.Auth{Create: true, Read: true, Update: true, Delete: true}

	// output, err := json.Marshal(g)
	out_count, err := json.Marshal(models.OrderType_count{Values: out_arr, Count: 1, Auth: auth})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(out_count)

	return

}
