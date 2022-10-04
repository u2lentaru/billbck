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

type ifPaymentService interface {
	GetList(ctx context.Context, pg, pgs int, gs1, gs2 string, ord int, dsc bool) (models.Payment_count, error)
	Add(ctx context.Context, ea models.Payment) (int, error)
	Upd(ctx context.Context, eu models.Payment) (int, error)
	Del(ctx context.Context, ed []int) ([]int, error)
	GetOne(ctx context.Context, i int) (models.Payment_count, error)
}

// HandlePayments godoc
// @Summary List payments
// @Description get payments list
// @Tags payments
// @Produce  json
// @Param page query int false "page number"
// @Param page_size query int false "page size"
// @Param contractnumber query string false "contractnumber search pattern"
// @Param oid query string false "object id"
// @Param ordering query string false "order by {id|paymentdate}"
// @Param desc query boolean false "descending order {true|false}"
// @Success 200 {object} models.Payment_count
// @Failure 500
// @Router /payments [get]
func HandlePayments(w http.ResponseWriter, r *http.Request) {
	var gs ifPaymentService
	gs = services.NewPaymentService(pgsql.PaymentStorage{})
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
	gs1s, ok := query["contractnumber"]
	if ok && len(gs1s) > 0 {
		//case insensitive
		gs1 = strings.ToUpper(gs1s[0])
		//quotes
		re := regexp.MustCompile(`'`)
		gs1 = string(re.ReplaceAll([]byte(gs1), []byte("''")))
	}

	gs2 := ""
	gs2s, ok := query["oid"]
	if ok && len(gs2s) > 0 {
		_, err := strconv.Atoi(gs2s[0])
		if err == nil {
			gs2 = gs2s[0]
		}
	}

	ord := 1
	ords, ok := query["ordering"]
	if !ok || len(ords) == 0 {
		ord = 1
	} else if ords[0] == "paymentdate" {
		ord = 2
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

// HandleAddPayment godoc
// @Summary Add payment
// @Description add payment
// @Tags payments
// @Accept json
// @Produce  json
// @Param a body models.AddPayment true "New payment. Significant params: PaymentDate, Contract.Id, Object.Id, PaymentType.Id, ChargeType.Id, Cashdesk.Id, BundleNumber, Amount"
// @Success 200 {object} models.Json_id
// @Failure 500
// @Router /payments_add [post]
func HandleAddPayment(w http.ResponseWriter, r *http.Request) {
	var gs ifPaymentService
	gs = services.NewPaymentService(pgsql.PaymentStorage{})
	ctx := context.Background()

	a := models.Payment{}
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
		log.Println("Failed execute ifPaymentService.Add: ", err)
	}

	output, err := json.Marshal(models.Json_id{Id: ai})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(output)

	return
}

// HandleUpdPayment godoc
// @Summary Update payment
// @Description update payment
// @Tags payments
// @Accept json
// @Produce  json
// @Param u body models.Payment true "Update payment. Significant params: Id, PaymentDate, Contract.Id, Object.Id, PaymentType.Id, ChargeType.Id, Cashdesk.Id, BundleNumber, Amount"
// @Success 200 {object} models.Json_id
// @Failure 500
// @Router /payments_upd [post]
func HandleUpdPayment(w http.ResponseWriter, r *http.Request) {
	var gs ifPaymentService
	gs = services.NewPaymentService(pgsql.PaymentStorage{})
	ctx := context.Background()

	u := models.Payment{}
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
		log.Println("Failed execute ifAPaymentService.Upd: ", err)
	}

	output, err := json.Marshal(models.Json_id{Id: ui})

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(output)

	return
}

// HandleDelPayment godoc
// @Summary Delete payments
// @Description delete payments
// @Tags payments
// @Accept json
// @Produce  json
// @Param d body models.Json_ids true "Delete payments"
// @Success 200 {object} models.Json_ids
// @Failure 500
// @Router /payments_del [post]
func HandleDelPayment(w http.ResponseWriter, r *http.Request) {
	var gs ifPaymentService
	gs = services.NewPaymentService(pgsql.PaymentStorage{})
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
		log.Println("Failed execute ifPaymentService.Del: ", err)
	}

	output, err := json.Marshal(models.Json_ids{Ids: res})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(output)

	return

}

// HandleGetPayment godoc
// @Summary Get payment
// @Description get payment
// @Tags payments
// @Produce  json
// @Param id path int true "Payment by id"
// @Success 200 {object} models.Payment_count
// @Failure 500
// @Router /payments/{id} [get]
func HandleGetPayment(w http.ResponseWriter, r *http.Request) {
	var gs ifPaymentService
	gs = services.NewPaymentService(pgsql.PaymentStorage{})
	ctx := context.Background()
	auth := utils.GetAuth(r)

	vars := mux.Vars(r)
	i, err := strconv.Atoi(vars["id"])
	if err != nil {
		i = 0
	}

	out_arr, err := gs.GetOne(ctx, i)
	if err != nil {
		log.Println("Failed execute ifPaymentService.GetOne: ", err)
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
