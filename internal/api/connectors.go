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
	"github.com/u2lentaru/billbck/internal/adapters/db/pgsql"
	"github.com/u2lentaru/billbck/internal/models"
	"github.com/u2lentaru/billbck/internal/services"
	"github.com/u2lentaru/billbck/internal/utils"
)

type ifConnectorService interface {
	GetList(ctx context.Context, pg, pgs int, gs1 string, ord int, dsc bool) (models.Connector_count, error)
	Add(ctx context.Context, ea models.Connector) (int, error)
	Upd(ctx context.Context, eu models.Connector) (int, error)
	Del(ctx context.Context, ed []int) ([]int, error)
	GetOne(ctx context.Context, i int) (models.Connector_count, error)
}

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
func HandleConnectors(w http.ResponseWriter, r *http.Request) {
	var gs ifConnectorService
	gs = services.NewConnectorService(pgsql.ConnectorStorage{})
	ctx := context.Background()
	auth := utils.GetAuth(r)

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

	out_arr, err := gs.GetList(ctx, pg, pgs, gs1, ord, dsc)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	out_arr.Auth = auth
	out_count, err := json.Marshal(out_arr)
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
func HandleAddConnector(w http.ResponseWriter, r *http.Request) {
	var gs ifConnectorService
	gs = services.NewConnectorService(pgsql.ConnectorStorage{})
	ctx := context.Background()

	a := models.Connector{}
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

	ai, err := gs.Add(ctx, a)

	if err != nil {
		log.Println("Failed execute ifConnectorService.Add: ", err)
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
func HandleUpdConnector(w http.ResponseWriter, r *http.Request) {
	var gs ifConnectorService
	gs = services.NewConnectorService(pgsql.ConnectorStorage{})
	ctx := context.Background()

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

	ui, err := gs.Upd(ctx, u)

	if err != nil {
		log.Println("Failed execute ifConnectorService.Upd: ", err)
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
func HandleDelConnector(w http.ResponseWriter, r *http.Request) {
	var gs ifConnectorService
	gs = services.NewConnectorService(pgsql.ConnectorStorage{})
	ctx := context.Background()

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

	res, err := gs.Del(ctx, d.Ids)
	if err != nil {
		log.Println("Failed execute ifConnectorService.Del: ", err)
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
func HandleGetConnector(w http.ResponseWriter, r *http.Request) {
	var gs ifConnectorService
	gs = services.NewConnectorService(pgsql.ConnectorStorage{})
	ctx := context.Background()
	auth := utils.GetAuth(r)

	vars := mux.Vars(r)
	i, err := strconv.Atoi(vars["id"])
	if err != nil {
		i = 0
	}

	out_arr, err := gs.GetOne(ctx, i)
	if err != nil {
		log.Println("Failed execute ifConnectorService.GetOne: ", err)
	}

	out_arr.Auth = auth
	out_count, err := json.Marshal(out_arr)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Write(out_count)

	return

}
