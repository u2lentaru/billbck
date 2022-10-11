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

type ifSectorService interface {
	GetList(ctx context.Context, pg, pgs int, gs1 string, ord int, dsc bool) (models.Sector_count, error)
	Add(ctx context.Context, ea models.Sector) (int, error)
	Upd(ctx context.Context, eu models.Sector) (int, error)
	Del(ctx context.Context, ed []int) ([]int, error)
	GetOne(ctx context.Context, i int) (models.Sector_count, error)
}

// HandleSectors godoc
// @Summary List sectors
// @Description get sector list
// @Tags sectors
// @Produce  json
// @Param page query int false "page number"
// @Param page_size query int false "page size"
// @Param sectorname query string false "sectorname search pattern"
// @Param ordering query string false "order by {id|sectorname}"
// @Param desc query boolean false "descending order {true|false}"
// @Success 200 {object} models.Sector_count
// @Failure 500
// @Router /sectors [get]
func HandleSectors(w http.ResponseWriter, r *http.Request) {
	var gs ifSectorService
	gs = services.NewSectorService(pgsql.SectorStorage{})
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
	gs1s, ok := query["sectorname"]
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
	} else if ords[0] == "sectorname" {
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

// HandleAddSector godoc
// @Summary Add sector
// @Description add sector
// @Tags sectors
// @Accept json
// @Produce  json
// @Param a body models.AddSector true "New sector"
// @Success 200 {object} models.Json_id
// @Failure 500
// @Router /sectors_add [post]
func HandleAddSector(w http.ResponseWriter, r *http.Request) {
	var gs ifSectorService
	gs = services.NewSectorService(pgsql.SectorStorage{})
	ctx := context.Background()

	a := models.Sector{}
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
		log.Println("Failed execute ifSectorService.Add: ", err)
	}

	output, err := json.Marshal(models.Json_id{Id: ai})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(output)

	return

}

// HandleUpdSector godoc
// @Summary Update sector
// @Description update sector
// @Tags sectors
// @Accept json
// @Produce  json
// @Param u body models.Sector true "Update sector"
// @Success 200 {object} models.Json_id
// @Failure 500
// @Router /sectors_upd [post]
func HandleUpdSector(w http.ResponseWriter, r *http.Request) {
	var gs ifSectorService
	gs = services.NewSectorService(pgsql.SectorStorage{})
	ctx := context.Background()

	u := models.Sector{}
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
		log.Println("Failed execute ifSectorService.Upd: ", err)
	}

	output, err := json.Marshal(models.Json_id{Id: ui})

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(output)

	return
}

// HandleDelSector godoc
// @Summary Delete sectors
// @Description delete sectors
// @Tags sectors
// @Accept json
// @Produce  json
// @Param d body models.Json_ids true "Delete sectors"
// @Success 200 {object} models.Json_ids
// @Failure 500
// @Router /sectors_del [post]
func HandleDelSector(w http.ResponseWriter, r *http.Request) {
	var gs ifSectorService
	gs = services.NewSectorService(pgsql.SectorStorage{})
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
		log.Println("Failed execute ifSectorService.Del: ", err)
	}

	output, err := json.Marshal(models.Json_ids{Ids: res})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(output)

	return
}

// HandleGetSector godoc
// @Summary Get sector
// @Description get sector
// @Tags sectors
// @Produce  json
// @Param id path int true "Sector by id"
// @Success 200 {object} models.Sector_count
// @Failure 500
// @Router /sectors/{id} [get]
func HandleGetSector(w http.ResponseWriter, r *http.Request) {
	var gs ifSectorService
	gs = services.NewSectorService(pgsql.SectorStorage{})
	ctx := context.Background()
	auth := utils.GetAuth(r)

	vars := mux.Vars(r)
	i, err := strconv.Atoi(vars["id"])
	if err != nil {
		i = 0
	}

	out_arr, err := gs.GetOne(ctx, i)
	if err != nil {
		log.Println("Failed execute ifSectorService.GetOne: ", err)
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
