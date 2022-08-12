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
func (s *APG) HandlePayments(w http.ResponseWriter, r *http.Request) {
	gs := models.Payment{}
	ctx := context.Background()
	out_arr := []models.Payment{}

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

	rows, err := s.Dbpool.Query(ctx, "SELECT * from func_payments_get($1,$2,$3,$4,$5,$6);", pg, pgs, gs1, NullableString(gs2), ord, dsc)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&gs.Id, &gs.PaymentDate, &gs.Contract.Id, &gs.Object.Id, &gs.PaymentType.Id, &gs.ChargeType.Id, &gs.Cashdesk.Id,
			&gs.BundleNumber, &gs.Amount, &gs.Contract.ContractNumber, &gs.Object.ObjectName, &gs.PaymentType.PaymentTypeName,
			&gs.ChargeType.ChargeTypeName, &gs.Cashdesk.CashdeskName)
		if err != nil {
			log.Println("failed to scan row:", err)
		}

		out_arr = append(out_arr, gs)
	}

	gsc := 0
	err = s.Dbpool.QueryRow(ctx, "SELECT * from func_payments_cnt($1,$2);", gs1, NullableString(gs2)).Scan(&gsc)

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	auth := models.Auth{Create: true, Read: true, Update: true, Delete: true}
	out_count, err := json.Marshal(models.Payment_count{Values: out_arr, Count: gsc, Auth: auth})
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
func (s *APG) HandleAddPayment(w http.ResponseWriter, r *http.Request) {
	a := models.AddPayment{}
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
	err = s.Dbpool.QueryRow(context.Background(), "SELECT func_payments_add($1,$2,$3,$4,$5,$6,$7,$8);", a.PaymentDate, a.Contract.Id,
		a.Object.Id, a.PaymentType.Id, a.ChargeType.Id, a.Cashdesk.Id, a.BundleNumber, a.Amount).Scan(&ai)

	if err != nil {
		log.Println("Failed execute func_payments_add: ", err)
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
func (s *APG) HandleUpdPayment(w http.ResponseWriter, r *http.Request) {
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

	ui := 0
	err = s.Dbpool.QueryRow(context.Background(), "SELECT func_payments_upd($1,$2,$3,$4,$5,$6,$7,$8,$9);", u.Id, u.PaymentDate, u.Contract.Id,
		u.Object.Id, u.PaymentType.Id, u.ChargeType.Id, u.Cashdesk.Id, u.BundleNumber, u.Amount).Scan(&ui)

	if err != nil {
		log.Println("Failed execute func_payments_upd: ", err)
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
func (s *APG) HandleDelPayment(w http.ResponseWriter, r *http.Request) {
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
		err = s.Dbpool.QueryRow(context.Background(), "SELECT func_payments_del($1);", id).Scan(&i)
		res = append(res, i)

		if err != nil {
			log.Println("Failed execute func_payments_del: ", err)
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

// HandleGetPayment godoc
// @Summary Get payment
// @Description get payment
// @Tags payments
// @Produce  json
// @Param id path int true "Payment by id"
// @Success 200 {object} models.Payment_count
// @Failure 500
// @Router /payments/{id} [get]
func (s *APG) HandleGetPayment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	i := vars["id"]
	g := models.Payment{}
	out_arr := []models.Payment{}

	err := s.Dbpool.QueryRow(context.Background(), "SELECT * from func_payment_get($1);", i).Scan(&g.Id, &g.PaymentDate, &g.Contract.Id,
		&g.Object.Id, &g.PaymentType.Id, &g.ChargeType.Id, &g.Cashdesk.Id, &g.BundleNumber, &g.Amount, &g.Contract.ContractNumber,
		&g.Object.ObjectName, &g.PaymentType.PaymentTypeName, &g.ChargeType.ChargeTypeName, &g.Cashdesk.CashdeskName)

	if err != nil && err != pgx.ErrNoRows {
		log.Println("Failed execute from func_payment_get: ", err)
	}

	out_arr = append(out_arr, g)
	auth := models.Auth{Create: true, Read: true, Update: true, Delete: true}

	// output, err := json.Marshal(g)
	out_count, err := json.Marshal(models.Payment_count{Values: out_arr, Count: 1, Auth: auth})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(out_count)

	return

}
