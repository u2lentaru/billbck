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

type ifSubBankService interface {
	GetList(ctx context.Context, pg, pgs int, gs1 string, gs2 int, gs3 string, ord int, dsc bool) (models.SubBank_count, error)
	Add(ctx context.Context, ea models.SubBank) (int, error)
	Upd(ctx context.Context, eu models.SubBank) (int, error)
	Del(ctx context.Context, ed []int) ([]int, error)
	GetOne(ctx context.Context, i int) (models.SubBank_count, error)
	SetActive(ctx context.Context, i int) (int, error)
}

// HandleSubBanks godoc
// @Summary List subject accounts
// @Description Get subject accounts
// @Tags subject banks
// @Produce  json
// @Param page query int false "page number"
// @Param page_size query int false "page size"
// @Param subname query string false "subname search pattern"
// @Param subid query string false "subid search pattern"
// @Param accnumber query string false "accnumber search pattern"
// @Param ordering query string false "order by {subname|accnumber}"
// @Param desc query boolean false "descending order {true|false}"
// @Success 200 {object} models.SubBank_count
// @Failure 500
// @Router /sub_banks [get]
func HandleSubBanks(w http.ResponseWriter, r *http.Request) {
	var gs ifSubBankService
	gs = services.NewSubBankService(pgsql.SubBankStorage{})
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

	sbn := ""
	sbns, ok := query["subname"]
	if ok && len(sbns) > 0 {
		//case insensitive
		sbn = strings.ToUpper(sbns[0])
		//quotes
		re := regexp.MustCompile(`'`)
		sbn = string(re.ReplaceAll([]byte(sbn), []byte("''")))
	}

	sbi := 0
	sbis, ok := query["custid"]
	if ok && len(sbis) > 0 {
		t, err := strconv.Atoi(sbis[0])
		if err == nil {
			sbi = t
		}
	}

	an := ""
	ans, ok := query["accnumber"]
	if ok && len(ans) > 0 {
		an = strings.ToUpper(ans[0])
		re := regexp.MustCompile(`'`)
		an = string(re.ReplaceAll([]byte(an), []byte("''")))
	}

	ord := 1
	ords, ok := query["ordering"]
	if !ok || len(ords) == 0 {
		ord = 1
	} else if ords[0] == "subname" {
		ord = 4
	} else if ords[0] == "accnumber" {
		ord = 3
	}

	dsc := false
	dscs, ok := query["desc"]
	if ok && len(dscs) > 0 {
		if !(dscs[0] == "0") {
			dsc = true
		}
	}

	out_arr, err := gs.GetList(ctx, pg, pgs, sbn, sbi, an, ord, dsc)
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

// HandleAddSubBank godoc
// @Summary Add subject account
// @Description Add subject account
// @Tags subject banks
// @Accept json
// @Produce  json
// @Param ab body models.AddSubBank true "New subject account. Sets the first account of the subject active. Significant params: Sub.Id, Bank.Id, AccNumber"
// @Success 200 {object} models.Json_id
// @Failure 500
// @Router /sub_banks_add [post]
func HandleAddSubBank(w http.ResponseWriter, r *http.Request) {
	var gs ifSubBankService
	gs = services.NewSubBankService(pgsql.SubBankStorage{})
	ctx := context.Background()

	a := models.SubBank{}
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
		log.Println("Failed execute ifSubBankService.Add: ", err)
	}

	output, err := json.Marshal(models.Json_id{Id: ai})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(output)

	return

}

// HandleUpdSubBank godoc
// @Summary Update subject account
// @Description Update subject account
// @Tags subject banks
// @Accept json
// @Produce  json
// @Param ub body models.SubBank true "Update subject account. Significant params: Id, Sub.Id, Bank.Id, AccNumber"
// @Success 200 {object} models.Json_id
// @Failure 500
// @Router /sub_banks_upd [post]
func HandleUpdSubBank(w http.ResponseWriter, r *http.Request) {
	var gs ifSubBankService
	gs = services.NewSubBankService(pgsql.SubBankStorage{})
	ctx := context.Background()

	u := models.SubBank{}
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
		log.Println("Failed execute ifSubBankService.Upd: ", err)
	}

	output, err := json.Marshal(models.Json_id{Id: ui})

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(output)

	return

}

// HandleDelSubBank godoc
// @Summary Delete subject accounts
// @Description delete subject accounts
// @Tags subject banks
// @Accept json
// @Produce  json
// @Param d body models.Json_ids true "Delete subject accounts"
// @Success 200 {object} models.Json_ids
// @Failure 500
// @Router /sub_banks_del [post]
func HandleDelSubBank(w http.ResponseWriter, r *http.Request) {
	var gs ifSubBankService
	gs = services.NewSubBankService(pgsql.SubBankStorage{})
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
		log.Println("Failed execute ifSubBankService.Del: ", err)
	}

	output, err := json.Marshal(models.Json_ids{Ids: res})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(output)

	return

}

// HandleGetSubBank godoc
// @Summary Get subject account
// @Description get subject account
// @Tags subject banks
// @Produce  json
// @Param id path int true "Subject account by id"
// @Success 200 {object} models.SubBank_count
// @Failure 500
// @Router /sub_banks/{id} [get]
func HandleGetSubBank(w http.ResponseWriter, r *http.Request) {
	var gs ifSubBankService
	gs = services.NewSubBankService(pgsql.SubBankStorage{})
	ctx := context.Background()
	auth := utils.GetAuth(r)

	vars := mux.Vars(r)
	i, err := strconv.Atoi(vars["id"])
	if err != nil {
		i = 0
	}

	out_arr, err := gs.GetOne(ctx, i)
	if err != nil {
		log.Println("Failed execute ifSubBankService.GetOne: ", err)
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

// HandleGetSubBankSetActive godoc
// @Summary Set active subject account
// @Description set active subject account
// @Tags subject banks
// @Produce  json
// @Param id path int true "Sets the active account of the subject by ID, sets inactive all other accounts of the subject."
// @Success 200 {object} models.Json_id
// @Failure 500
// @Router /sub_banks_setactive/{id} [post]
func HandleGetSubBankSetActive(w http.ResponseWriter, r *http.Request) {
	var gs ifSubBankService
	gs = services.NewSubBankService(pgsql.SubBankStorage{})
	ctx := context.Background()

	vars := mux.Vars(r)
	i, err := strconv.Atoi(vars["id"])
	if err != nil {
		i = 0
	}

	res, err := gs.SetActive(ctx, i)
	if err != nil {
		log.Println("Failed execute ifSubBankService.SetActive: ", err)
	}

	output, err := json.Marshal(models.Json_id{Id: res})

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(output)

	return

}
