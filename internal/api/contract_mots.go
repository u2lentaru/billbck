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

type ifContractMotService interface {
	GetList(ctx context.Context, pg, pgs int, gs1, gs2 string, ord int, dsc bool) (models.ContractMot_count, error)
	Add(ctx context.Context, ea models.ContractMot) (int, error)
	Upd(ctx context.Context, eu models.ContractMot) (int, error)
	Del(ctx context.Context, ed []int) ([]int, error)
	GetOne(ctx context.Context, i int) (models.ContractMot_count, error)
}

// HandleContractMots godoc
// @Summary List contracts motives of termination
// @Description get contracts motives of termination list
// @Tags contract_mots
// @Produce  json
// @Param page query int false "page number"
// @Param page_size query int false "page size"
// @Param contractmotnameru query string false "contractmotnameru search pattern"
// @Param contractmotnamekz query string false "contractmotnamekz search pattern"
// @Param ordering query string false "order by {id|contractmotnameru|contractmotnamekz}"
// @Param desc query boolean false "descending order {true|false}"
// @Success 200 {object} models.ContractMot_count
// @Failure 500
// @Router /contractmots [get]
func HandleContractMots(w http.ResponseWriter, r *http.Request) {
	var gs ifContractMotService
	gs = services.NewContractMotService(pgsql.ContractMotStorage{})
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
	gs1s, ok := query["contractmotnameru"]
	if ok && len(gs1s) > 0 {
		//case insensitive
		gs1 = strings.ToUpper(gs1s[0])
		//quotes
		re := regexp.MustCompile(`'`)
		gs1 = string(re.ReplaceAll([]byte(gs1), []byte("''")))
	}

	gs2 := ""
	gs2s, ok := query["contractmotnamekz"]
	if ok && len(gs2s) > 0 {
		//case insensitive
		gs2 = strings.ToUpper(gs2s[0])
		//quotes
		re := regexp.MustCompile(`'`)
		gs2 = string(re.ReplaceAll([]byte(gs2), []byte("''")))
	}

	ord := 1
	ords, ok := query["ordering"]
	if !ok || len(ords) == 0 {
		ord = 1
	} else if ords[0] == "contractmotnameru" {
		ord = 2
	} else if ords[0] == "contractmotnamekz" {
		ord = 3
	}

	dsc := false
	dscs, ok := query["desc"]
	if ok && len(dscs) > 0 {
		if !(dscs[0] == "0") {
			dsc = true
		}
	}

	out_arr, err := gs.GetList(ctx, pg, pgs, gs1, gs2, ord, dsc)
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

// HandleAddContractMot godoc
// @Summary Add contracts motive of termination
// @Description add contracts motive of termination
// @Tags contract_mots
// @Accept json
// @Produce  json
// @Param a body models.AddContractMot true "New contracts motive of termination"
// @Success 200 {object} models.Json_id
// @Failure 500
// @Router /contractmots_add [post]
func HandleAddContractMot(w http.ResponseWriter, r *http.Request) {
	var gs ifContractMotService
	gs = services.NewContractMotService(pgsql.ContractMotStorage{})
	ctx := context.Background()

	a := models.ContractMot{}
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
		log.Println("Failed execute ifContractMotService.Add: ", err)
	}

	output, err := json.Marshal(models.Json_id{Id: ai})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(output)

	return

}

// HandleUpdContractMot godoc
// @Summary Update contracts motive of termination
// @Description update contracts motive of termination
// @Tags contract_mots
// @Accept json
// @Produce  json
// @Param u body models.ContractMot true "Update contracts motive of termination"
// @Success 200 {object} models.Json_id
// @Failure 500
// @Router /contractmots_upd [post]
func HandleUpdContractMot(w http.ResponseWriter, r *http.Request) {
	var gs ifContractMotService
	gs = services.NewContractMotService(pgsql.ContractMotStorage{})
	ctx := context.Background()

	u := models.ContractMot{}
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
		log.Println("Failed execute ifContractMotService.Upd: ", err)
	}

	output, err := json.Marshal(models.Json_id{Id: ui})

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(output)

	return

}

// HandleDelContractMot godoc
// @Summary Delete contracts motives of termination
// @Description delete contracts motives of termination
// @Tags contract_mots
// @Accept json
// @Produce  json
// @Param d body models.Json_ids true "Delete contracts motives of termination"
// @Success 200 {object} models.Json_ids
// @Failure 500
// @Router /contractmots_del [post]
func HandleDelContractMot(w http.ResponseWriter, r *http.Request) {
	var gs ifContractMotService
	gs = services.NewContractMotService(pgsql.ContractMotStorage{})
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
		log.Println("Failed execute ifContractMotService.Del: ", err)
	}

	output, err := json.Marshal(models.Json_ids{Ids: res})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(output)

	return

}

// HandleGetContractMot godoc
// @Summary Get contracts motive of termination
// @Description get contracts motive of termination
// @Tags contract_mots
// @Produce  json
// @Param id path int true "Contracts motive of termination by id"
// @Success 200 {object} models.ContractMot_count
// @Failure 500
// @Router /contractmots/{id} [get]
func HandleGetContractMot(w http.ResponseWriter, r *http.Request) {
	var gs ifContractMotService
	gs = services.NewContractMotService(pgsql.ContractMotStorage{})
	ctx := context.Background()
	auth := utils.GetAuth(r)

	vars := mux.Vars(r)
	i, err := strconv.Atoi(vars["id"])
	if err != nil {
		i = 0
	}

	out_arr, err := gs.GetOne(ctx, i)
	if err != nil {
		log.Println("Failed execute ifContractMotService.GetOne: ", err)
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
