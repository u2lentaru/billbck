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

// HandleConnectors godoc
// @Summary List connectors
// @Description get connector list
// @Tags connectors
// @Produce  json
// @Param page query int false "page number"
// @Param page_size query int false "page size"
// @Param connectorname query string false "connectorname search pattern"
// @Param ordering query string false "order by {id|connectorname}"
// @Param desc query boolean false "descending order {true|false}"
// @Success 200 {object} models.Connector_count
// @Failure 500
// @Router /connectors [get]
func (s *APG) HandleConnectors(w http.ResponseWriter, r *http.Request) {
	gs := models.Connector{}
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
	gs1s, ok := query["connectorname"]
	if ok && len(gs1s) > 0 {
		//case insensitive
		gs1 = strings.ToUpper(gs1s[0])
		//quotes
		re := regexp.MustCompile(`'`)
		gs1 = string(re.ReplaceAll([]byte(gs1), []byte("''")))
	}

	gsc := 0
	err := s.Dbpool.QueryRow(ctx, "SELECT * from func_connectors_cnt($1);", gs1).Scan(&gsc)

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	out_arr := make([]models.Connector, 0,
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
	} else if ords[0] == "connectorname" {
		ord = 2
	}

	dsc := false
	dscs, ok := query["desc"]
	if ok && len(dscs) > 0 {
		if !(dscs[0] == "0") {
			dsc = true
		}
	}

	rows, err := s.Dbpool.Query(ctx, "SELECT * from func_connectors_get($1,$2,$3,$4,$5);", pg, pgs, gs1, ord, dsc)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&gs.Id, &gs.ConnectorName)
		if err != nil {
			log.Println("failed to scan row:", err)
		}

		out_arr = append(out_arr, gs)
	}

	auth := models.Auth{Create: true, Read: true, Update: true, Delete: true}
	out_count, err := json.Marshal(models.Connector_count{Values: out_arr, Count: gsc, Auth: auth})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Write(out_count)

	return
}

// HandleAddConnector godoc
// @Summary Add connector
// @Description add connector
// @Tags connectors
// @Accept json
// @Produce  json
// @Param a body models.AddConnector true "New connector"
// @Success 200 {object} models.Json_id
// @Failure 500
// @Router /connectors_add [post]
func (s *APG) HandleAddConnector(w http.ResponseWriter, r *http.Request) {
	a := models.AddConnector{}
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
	err = s.Dbpool.QueryRow(context.Background(), "SELECT func_connectors_add($1);", a.ConnectorName).Scan(&ai)

	if err != nil {
		log.Println("Failed execute func_connectors_add: ", err)
	}

	output, err := json.Marshal(models.Json_id{Id: ai})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(output)

	return

}

// HandleUpdConnector godoc
// @Summary Update connector
// @Description update connector
// @Tags connectors
// @Accept json
// @Produce  json
// @Param u body models.Connector true "Update connector"
// @Success 200 {object} models.Json_id
// @Failure 500
// @Router /connectors_upd [post]
func (s *APG) HandleUpdConnector(w http.ResponseWriter, r *http.Request) {
	u := models.Connector{}
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
	err = s.Dbpool.QueryRow(context.Background(), "SELECT func_connectors_upd($1,$2);", u.Id, u.ConnectorName).Scan(&ui)

	if err != nil {
		log.Println("Failed execute func_connectors_upd: ", err)
	}

	output, err := json.Marshal(models.Json_id{Id: ui})

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(output)

	return
}

// HandleDelConnector godoc
// @Summary Delete connectors
// @Description delete connectors
// @Tags connectors
// @Accept json
// @Produce  json
// @Param d body models.Json_ids true "Delete connectors"
// @Success 200 {object} models.Json_ids
// @Failure 500
// @Router /connectors_del [post]
func (s *APG) HandleDelConnector(w http.ResponseWriter, r *http.Request) {
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
		err = s.Dbpool.QueryRow(context.Background(), "SELECT func_connectors_del($1);", id).Scan(&i)
		res = append(res, i)

		if err != nil {
			log.Println("Failed execute func_connectors_del: ", err)
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

// HandleGetConnector godoc
// @Summary Get connector
// @Description get connector
// @Tags connectors
// @Produce  json
// @Param id path int true "Connector by id"
// @Success 200 {object} models.Connector_count
// @Failure 500
// @Router /connectors/{id} [get]
func (s *APG) HandleGetConnector(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	i := vars["id"]
	g := models.Connector{}
	out_arr := []models.Connector{}

	err := s.Dbpool.QueryRow(context.Background(), "SELECT * from func_connector_get($1);", i).Scan(&g.Id, &g.ConnectorName)

	if err != nil && err != pgx.ErrNoRows {
		log.Println("Failed execute from func_connector_get: ", err)
	}

	out_arr = append(out_arr, g)
	auth := models.Auth{Create: true, Read: true, Update: true, Delete: true}

	// output, err := json.Marshal(g)
	out_count, err := json.Marshal(models.Connector_count{Values: out_arr, Count: 1, Auth: auth})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Write(out_count)

	return

}
