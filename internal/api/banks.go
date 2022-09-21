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
)

type ifBankService interface {
	GetList(ctx context.Context, pg, pgs int, gs1, gs2 string, ord int, dsc bool) (models.Bank_count, error)
	Add(ctx context.Context, ea models.Bank) (int, error)
	Upd(ctx context.Context, eu models.Bank) (int, error)
	Del(ctx context.Context, ed []int) ([]int, error)
	GetOne(ctx context.Context, i int) (models.Bank_count, error)
}

// HandleBanks godoc
// @Summary List banks
// @Description Get banks
// @Tags banks
// @Produce  json
// @Param page query int false "page number"
// @Param page_size query int false "page size"
// @Param bankname query string false "bankname search pattern"
// @Param bankdescr query string false "bankdescr search pattern"
// @Param ordering query string false "order by {bankname|bankdescr}"
// @Param desc query boolean false "descending order {true|false}"
// @Success 200 {object} models.Bank_count
// @Failure 500
// @Router /banks [get]
func HandleBanks(w http.ResponseWriter, r *http.Request) {
	var gs ifBankService
	gs = services.NewBankService(pgsql.BankStorage{})
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

	bn := ""
	bns, ok := query["bankname"]
	if ok && len(bns) > 0 {
		//case insensitive
		bn = strings.ToUpper(bns[0])
		//quotes
		re := regexp.MustCompile(`'`)
		bn = string(re.ReplaceAll([]byte(bn), []byte("''")))
	}

	bd := ""
	bds, ok := query["bankdescr"]
	if ok && len(bds) > 0 {
		bd = strings.ToUpper(bds[0])
		re := regexp.MustCompile(`'`)
		bd = string(re.ReplaceAll([]byte(bd), []byte("''")))
	}

	ord := 1
	ords, ok := query["ordering"]
	if !ok || len(ords) == 0 {
		ord = 1
	} else if ords[0] == "bankname" {
		ord = 2
	} else if ords[0] == "bankdescr" {
		ord = 3
	}

	dsc := false
	dscs, ok := query["desc"]
	if ok && len(dscs) > 0 {
		if !(dscs[0] == "0") {
			dsc = true
		}
	}

	out_arr, err := gs.GetList(ctx, pg, pgs, bn, bd, ord, dsc)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	out_count, err := json.Marshal(out_arr)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Write(out_count)

	return

}

// HandleAddBank godoc
// @Summary Add bank
// @Description Add bank
// @Tags banks
// @Accept json
// @Produce  json
// @Param ab body models.AddBank true "New bank"
// @Success 200 {object} models.Json_id
// @Failure 500
// @Router /banks_add [post]
func HandleAddBank(w http.ResponseWriter, r *http.Request) {
	var gs ifBankService
	gs = services.NewBankService(pgsql.BankStorage{})
	ctx := context.Background()

	ab := models.Bank{}
	body, err := ioutil.ReadAll(r.Body)

	defer r.Body.Close()

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	err = json.Unmarshal(body, &ab)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	ai, err := gs.Add(ctx, ab)

	if err != nil {
		log.Println("Failed execute ifBankService.Add: ", err)
	}

	output, err := json.Marshal(models.Json_id{Id: ai})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(output)

	return

}

// HandleUpdBank godoc
// @Summary Update bank
// @Description Update bank
// @Tags banks
// @Accept json
// @Produce  json
// @Param ub body models.Bank true "Update bank"
// @Success 200 {object} models.Json_id
// @Failure 500
// @Router /banks_upd [post]
func HandleUpdBank(w http.ResponseWriter, r *http.Request) {
	var gs ifBankService
	gs = services.NewBankService(pgsql.BankStorage{})
	ctx := context.Background()

	u := models.Bank{}
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
		log.Println("Failed execute ifBankService.Upd: ", err)
	}

	output, err := json.Marshal(models.Json_id{Id: ui})

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(output)

	return
}

// HandleDelBank godoc
// @Summary Delete banks
// @Description Delete banks
// @Tags banks
// @Accept json
// @Produce  json
// @Param db body models.Json_ids true "Delete banks"
// @Success 200 {object} models.Json_ids
// @Failure 500
// @Router /banks_del [post]
func HandleDelBank(w http.ResponseWriter, r *http.Request) {
	var gs ifBankService
	gs = services.NewBankService(pgsql.BankStorage{})
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
		log.Println("Failed execute ifBankService.Del: ", err)
	}

	output, err := json.Marshal(models.Json_ids{Ids: res})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(output)

	return
}

// HandleGetBank godoc
// @Summary Get bank
// @Description Get bank
// @Tags banks
// @Produce  json
// @Param id path int true "Bank by id"
// @Success 200 {object} models.Bank_count
// @Failure 500
// @Router /banks/{id} [get]
func HandleGetBank(w http.ResponseWriter, r *http.Request) {
	var gs ifBankService
	gs = services.NewBankService(pgsql.BankStorage{})
	ctx := context.Background()

	vars := mux.Vars(r)
	i, err := strconv.Atoi(vars["id"])
	if err != nil {
		i = 0
	}

	out_arr, err := gs.GetOne(ctx, i)
	if err != nil {
		log.Println("Failed execute ifBankService.GetOne: ", err)
	}

	out_count, err := json.Marshal(out_arr)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(out_count)

	return
}
