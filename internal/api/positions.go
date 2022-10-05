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

type ifPositionService interface {
	GetList(ctx context.Context, pg, pgs int, gs1 string, ord int, dsc bool) (models.Position_count, error)
	Add(ctx context.Context, ea models.Position) (int, error)
	Upd(ctx context.Context, eu models.Position) (int, error)
	Del(ctx context.Context, ed []int) ([]int, error)
	GetOne(ctx context.Context, i int) (models.Position_count, error)
}

// HandlePositions godoc
// @Summary List positions
// @Description get positions
// @Tags positions
// @Produce  json
// @Param page query int false "page number"
// @Param page_size query int false "page size"
// @Param positionname query string false "positionname search pattern"
// @Param ordering query string false "order by {id|positionname}"
// @Param desc query boolean false "descending order {true|false}"
// @Success 200 {object} models.Position_count
// @Failure 500
// @Router /positions [get]
func HandlePositions(w http.ResponseWriter, r *http.Request) {
	var gs ifPositionService
	gs = services.NewPositionService(pgsql.PositionStorage{})
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

	pn := ""
	pns, ok := query["positionname"]
	if ok && len(pns) > 0 {
		//case insensitive
		pn = strings.ToUpper(pns[0])
		//quotes
		re := regexp.MustCompile(`'`)
		pn = string(re.ReplaceAll([]byte(pn), []byte("''")))
	}

	ord := 1
	ords, ok := query["ordering"]
	if !ok || len(ords) == 0 {
		ord = 1
	} else if ords[0] == "id" {
		ord = 1
	} else if ords[0] == "positionname" {
		ord = 2
	}

	dsc := false
	dscs, ok := query["desc"]
	if ok && len(dscs) > 0 {
		if !(dscs[0] == "0") {
			dsc = true
		}
	}

	out_arr, err := gs.GetList(ctx, pg, pgs, pn, ord, dsc)
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

// HandleAddPosition godoc
// @Summary Add position
// @Description add position
// @Tags positions
// @Accept json
// @Produce  json
// @Param ap body models.AddPosition true "New Position"
// @Success 200 {object} models.Json_id
// @Failure 500
// @Router /positions_add [post]
func HandleAddPosition(w http.ResponseWriter, r *http.Request) {
	var gs ifPositionService
	gs = services.NewPositionService(pgsql.PositionStorage{})
	ctx := context.Background()

	a := models.Position{}
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
		log.Println("Failed execute ifPositionService.Add: ", err)
	}

	output, err := json.Marshal(models.Json_id{Id: ai})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(output)

	return

}

// HandleUpdPosition godoc
// @Summary Update position
// @Description update position
// @Tags positions
// @Accept json
// @Produce  json
// @Param up body models.Position true "Update position"
// @Success 200 {object} models.Json_id
// @Failure 500
// @Router /positions_upd [post]
func HandleUpdPosition(w http.ResponseWriter, r *http.Request) {
	var gs ifPositionService
	gs = services.NewPositionService(pgsql.PositionStorage{})
	ctx := context.Background()

	u := models.Position{}
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
		log.Println("Failed execute ifPositionService.Upd: ", err)
	}

	output, err := json.Marshal(models.Json_id{Id: ui})

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(output)

	return

}

// HandleDelPosition godoc
// @Summary Delete positions
// @Description delete positions
// @Tags positions
// @Accept json
// @Produce  json
// @Param dp body models.Json_ids true "Delete positions"
// @Success 200 {object} models.Json_ids
// @Failure 500
// @Router /positions_del [post]
func HandleDelPosition(w http.ResponseWriter, r *http.Request) {
	var gs ifPositionService
	gs = services.NewPositionService(pgsql.PositionStorage{})
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
		log.Println("Failed execute ifPositionService.Del: ", err)
	}

	output, err := json.Marshal(models.Json_ids{Ids: res})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(output)

	return
}

// HandleGetPosition godoc
// @Summary Get position
// @Description get position
// @Tags positions
// @Produce  json
// @Param id path int true "Position by id"
// @Success 200 {object} models.Position_count
// @Failure 500
// @Router /positions/{id} [get]
func HandleGetPosition(w http.ResponseWriter, r *http.Request) {
	var gs ifPositionService
	gs = services.NewPositionService(pgsql.PositionStorage{})
	ctx := context.Background()
	auth := utils.GetAuth(r)

	vars := mux.Vars(r)
	i, err := strconv.Atoi(vars["id"])
	if err != nil {
		i = 0
	}

	out_arr, err := gs.GetOne(ctx, i)
	if err != nil {
		log.Println("Failed execute ifPositionService.GetOne: ", err)
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
