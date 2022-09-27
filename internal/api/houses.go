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

type ifHouseService interface {
	GetList(ctx context.Context, pg, pgs int, gs1, gs2 string, gs3, ord int, dsc bool) (models.House_count, error)
	Add(ctx context.Context, ea models.House) (int, error)
	Upd(ctx context.Context, eu models.House) (int, error)
	Del(ctx context.Context, ed []int) ([]int, error)
	GetOne(ctx context.Context, i int) (models.House_count, error)
}

// HandleHouses godoc
// @Summary List houses
// @Description get house list
// @Tags houses
// @Produce  json
// @Param page query int false "page number"
// @Param page_size query int false "page size"
// @Param streetname query string false "streetname search pattern"
// @Param housenumber query string false "housenumber search pattern"
// @Param streetid query int false "streetid search pattern"
// @Param ordering query string false "order by {housenumber|streetname}"
// @Param desc query boolean false "descending order {true|false}"
// @Success 200 {object} models.House_count
// @Failure 500
// @Router /houses [get]
func HandleHouses(w http.ResponseWriter, r *http.Request) {
	var gs ifHouseService
	gs = services.NewHouseService(pgsql.HouseStorage{})
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
	gs1s, ok := query["streetname"]
	if ok && len(gs1s) > 0 {
		//case insensitive
		gs1 = strings.ToUpper(gs1s[0])
		//quotes
		re := regexp.MustCompile(`'`)
		gs1 = string(re.ReplaceAll([]byte(gs1), []byte("''")))
	}

	gs2 := ""
	gs2s, ok := query["housenumber"]
	if ok && len(gs2s) > 0 {
		//case insensitive
		gs2 = strings.ToUpper(gs2s[0])
		//quotes
		re := regexp.MustCompile(`'`)
		gs2 = string(re.ReplaceAll([]byte(gs2), []byte("''")))
	}

	gs3 := 0
	gs3s, ok := query["streetid"]
	if ok && len(gs3s) > 0 {
		t, err := strconv.Atoi(gs3s[0])
		if err == nil {
			gs3 = t
		}
	}

	ord := 1
	ords, ok := query["ordering"]
	if !ok || len(ords) == 0 {
		ord = 1
	} else if ords[0] == "streetname" {
		ord = 15
	} else if ords[0] == "housenumber" {
		ord = 4
	}

	dsc := false
	dscs, ok := query["desc"]
	if ok && len(dscs) > 0 {
		if !(dscs[0] == "0") {
			dsc = true
		}
	}

	out_arr, err := gs.GetList(ctx, pg, pgs, gs1, gs2, gs3, ord, dsc)
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

// HandleAddHouse godoc
// @Summary Add house
// @Description add house
// @Tags houses
// @Accept json
// @Produce  json
// @Param a body models.AddHouse true "New house"
// @Success 200 {object} models.Json_id
// @Failure 500
// @Router /houses_add [post]
func HandleAddHouse(w http.ResponseWriter, r *http.Request) {
	var gs ifHouseService
	gs = services.NewHouseService(pgsql.HouseStorage{})
	ctx := context.Background()

	a := models.House{}
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
		log.Println("Failed execute ifHouseService.Add: ", err)
	}

	output, err := json.Marshal(models.Json_id{Id: ai})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(output)

	return

}

// HandleUpdHouse godoc
// @Summary Update house
// @Description update house
// @Tags houses
// @Accept json
// @Produce  json
// @Param u body models.House true "Update house"
// @Success 200 {object} models.Json_id
// @Failure 500
// @Router /houses_upd [post]
func HandleUpdHouse(w http.ResponseWriter, r *http.Request) {
	var gs ifHouseService
	gs = services.NewHouseService(pgsql.HouseStorage{})
	ctx := context.Background()

	u := models.House{}
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
		log.Println("Failed execute ifHouseService.Upd: ", err)
	}

	output, err := json.Marshal(models.Json_id{Id: ui})

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(output)

	return
}

// HandleDelHouse godoc
// @Summary Delete houses
// @Description delete houses
// @Tags houses
// @Accept json
// @Produce  json
// @Param d body models.Json_ids true "Delete houses"
// @Success 200 {object} models.Json_ids
// @Failure 500
// @Router /houses_del [post]
func HandleDelHouse(w http.ResponseWriter, r *http.Request) {
	var gs ifHouseService
	gs = services.NewHouseService(pgsql.HouseStorage{})
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
		log.Println("Failed execute ifHouseService.Del: ", err)
	}

	output, err := json.Marshal(models.Json_ids{Ids: res})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(output)

	return
}

// HandleGetHouse godoc
// @Summary Get house
// @Description get house
// @Tags houses
// @Produce  json
// @Param id path int true "House by id"
// @Success 200 {object} models.House_count
// @Failure 500
// @Router /houses/{id} [get]
func HandleGetHouse(w http.ResponseWriter, r *http.Request) {
	var gs ifHouseService
	gs = services.NewHouseService(pgsql.HouseStorage{})
	ctx := context.Background()
	auth := utils.GetAuth(r)

	vars := mux.Vars(r)
	i, err := strconv.Atoi(vars["id"])
	if err != nil {
		i = 0
	}

	out_arr, err := gs.GetOne(ctx, i)
	if err != nil {
		log.Println("Failed execute ifHouseService.GetOne: ", err)
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
