package api

import (
	"context"
	"encoding/json"
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

type ifBalanceService interface {
	GetList(ctx context.Context, pg, pgs, gs1, gs2 int) (models.Balance_count, error)
	GetNode(ctx context.Context, gs1, gs2 string) (models.Balance, error)
	GetNodeSum(ctx context.Context, gs1, gs2 int, gs3, gs4 string) (models.Json_sum, error)
	GetNodeSumL1(ctx context.Context, gs1, gs2 int, gs3, gs4 string) (models.Json_sum, error)
	GetNodeSumL0(ctx context.Context, gs1, gs2 int, gs3, gs4 string) (models.Json_sum, error)
	GetTabL1(ctx context.Context, pg, pgs, gs1, gs2 int, gs3, gs4 string) (models.BalanceTab_sum, error)
	GetTabL0(ctx context.Context, pg, pgs, gs1, gs2 int, gs3, gs4 string) (models.BalanceTab_sum, error)
	GetBranch(ctx context.Context, gs1, gs2 int) (models.BalanceTab_sum, error)
}

// HandleBalance godoc
// @Summary Get balance nodes
// @Description get balance nodes
// @Tags balance
// @Produce  json
// @Param page query int false "page number"
// @Param page_size query int false "page size"
// @Param pid query int true "Balance nodes by pid"
// @Param tid query int true "Balance node type id"
// @Success 200 {object} models.Balance_count
// @Failure 500
// @Router /balance [get]
func HandleBalance(w http.ResponseWriter, r *http.Request) {
	var gs ifBalanceService
	gs = services.NewBalanceService(pgsql.BalanceStorage{})
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

	gs1 := 0
	gs1s, ok := query["pid"]
	if ok && len(gs1s) > 0 {
		t, err := strconv.Atoi(gs1s[0])
		if err == nil {
			gs1 = t
		}
	}

	gs2 := 0
	gs2s, ok := query["tid"]
	if ok && len(gs2s) > 0 {
		t, err := strconv.Atoi(gs2s[0])
		if err == nil {
			gs2 = t
		}
	}

	out_arr, err := gs.GetList(ctx, pg, pgs, gs1, gs2)
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

// HandleGetBalance godoc
// @Summary Get balance node
// @Description get balance node
// @Tags balance
// @Produce  json
// @Param id path int true "Balance node by id"
// @Param tid path int true "Balance node type id"
// @Success 200 {object} models.Balance
// @Failure 500
// @Router /balance/{id}/{tid} [get]
func HandleGetBalance(w http.ResponseWriter, r *http.Request) {
	var gs ifBalanceService
	gs = services.NewBalanceService(pgsql.BalanceStorage{})
	ctx := context.Background()

	vars := mux.Vars(r)
	i := vars["id"]
	tid := vars["tid"]

	out_arr, err := gs.GetNode(ctx, i, tid)
	if err != nil {
		log.Println("Failed execute ifBalanceService.GetNode: ", err)
	}

	out_count, err := json.Marshal(out_arr)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(out_count)

	return

}

// HandleBalanceSum godoc
// @Summary Get balance by node id, type and period
// @Description get balance by node id, type and period
// @Tags balance
// @Produce  json
// @Param id query int true "Balance node id"
// @Param tid query int true "Balance node type id"
// @Param sd query string false "Balance startdate. Default first day of previous month"
// @Param ed query string false "Balance enddate. Default last day of previous month"
// @Success 200 {object} models.Json_sum
// @Router /balance_sum [get]
func HandleBalanceSum(w http.ResponseWriter, r *http.Request) {
	var gs ifBalanceService
	gs = services.NewBalanceService(pgsql.BalanceStorage{})
	ctx := context.Background()

	query := r.URL.Query()

	gs1 := 0
	gs1s, ok := query["id"]
	if ok && len(gs1s) > 0 {
		t, err := strconv.Atoi(gs1s[0])
		if err == nil {
			gs1 = t
		}
	}

	gs2 := 0
	gs2s, ok := query["tid"]
	if ok && len(gs2s) > 0 {
		t, err := strconv.Atoi(gs2s[0])
		if err == nil {
			gs2 = t
		}
	}

	gs3 := utils.FirstOfPrevMonth().Format("2006-01-02")
	gs3s, ok := query["sd"]
	if ok && len(gs3s) > 0 {
		//case insensitive
		gs3 = strings.ToUpper(gs3s[0])
		//quotes
		re := regexp.MustCompile(`'`)
		gs3 = string(re.ReplaceAll([]byte(gs3), []byte("''")))
	}

	gs4 := utils.LastOfPrevMonth().Format("2006-01-02")
	gs4s, ok := query["ed"]
	if ok && len(gs4s) > 0 {
		//case insensitive
		gs4 = strings.ToUpper(gs4s[0])
		//quotes
		re := regexp.MustCompile(`'`)
		gs4 = string(re.ReplaceAll([]byte(gs4), []byte("''")))
	}

	sum, err := gs.GetNodeSum(ctx, gs1, gs2, gs3, gs4)

	if err != nil {
		log.Println("Failed execute ifBalanceService.GetNodeSum: ", err)
	}

	out_count, err := json.Marshal(sum)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Write(out_count)

	return
}

// HandleBalanceSum1 godoc
// @Summary Get balance 1 level down by node id, type and period
// @Description get balance 1 level down by node id, type and period
// @Tags balance
// @Produce  json
// @Param id query int true "Main node id"
// @Param tid query int true "Main node type id"
// @Param sd query string false "Balance startdate. Default first day of previous month"
// @Param ed query string false "Balance enddate. Default last day of previous month"
// @Success 200 {object} models.Json_sum
// @Router /balance_sum1 [get]
func HandleBalanceSum1(w http.ResponseWriter, r *http.Request) {
	var gs ifBalanceService
	gs = services.NewBalanceService(pgsql.BalanceStorage{})
	ctx := context.Background()

	query := r.URL.Query()

	gs1 := 0
	gs1s, ok := query["id"]
	if ok && len(gs1s) > 0 {
		t, err := strconv.Atoi(gs1s[0])
		if err == nil {
			gs1 = t
		}
	}

	gs2 := 0
	gs2s, ok := query["tid"]
	if ok && len(gs2s) > 0 {
		t, err := strconv.Atoi(gs2s[0])
		if err == nil {
			gs2 = t
		}
	}

	// gs3 := time.Now().Format("2006-01-02")
	gs3 := utils.FirstOfPrevMonth().Format("2006-01-02")
	gs3s, ok := query["sd"]
	if ok && len(gs3s) > 0 {
		//case insensitive
		gs3 = strings.ToUpper(gs3s[0])
		//quotes
		re := regexp.MustCompile(`'`)
		gs3 = string(re.ReplaceAll([]byte(gs3), []byte("''")))
	}

	// gs4 := time.Now().Format("2006-01-02")
	gs4 := utils.LastOfPrevMonth().Format("2006-01-02")
	gs4s, ok := query["ed"]
	if ok && len(gs4s) > 0 {
		//case insensitive
		gs4 = strings.ToUpper(gs4s[0])
		//quotes
		re := regexp.MustCompile(`'`)
		gs4 = string(re.ReplaceAll([]byte(gs4), []byte("''")))
	}

	sum, err := gs.GetNodeSumL1(ctx, gs1, gs2, gs3, gs4)

	if err != nil {
		log.Println("Failed execute ifBalanceService.GetNodeSumL1: ", err)
	}

	out_count, err := json.Marshal(sum)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Write(out_count)

	return
}

// HandleBalanceSum0 godoc
// @Summary Get endpoints balance by node id, type and period
// @Description get endpoints balance by node id, type and period
// @Tags balance
// @Produce  json
// @Param id query int true "Main node id"
// @Param tid query int true "Main node type id"
// @Param sd query string false "Balance startdate. Default first day of previous month"
// @Param ed query string false "Balance enddate. Default last day of previous month"
// @Success 200 {object} models.Json_sum
// @Router /balance_sum0 [get]
func HandleBalanceSum0(w http.ResponseWriter, r *http.Request) {
	var gs ifBalanceService
	gs = services.NewBalanceService(pgsql.BalanceStorage{})
	ctx := context.Background()

	query := r.URL.Query()

	gs1 := 0
	gs1s, ok := query["id"]
	if ok && len(gs1s) > 0 {
		t, err := strconv.Atoi(gs1s[0])
		if err == nil {
			gs1 = t
		}
	}

	gs2 := 0
	gs2s, ok := query["tid"]
	if ok && len(gs2s) > 0 {
		t, err := strconv.Atoi(gs2s[0])
		if err == nil {
			gs2 = t
		}
	}

	// gs3 := time.Now().Format("2006-01-02")
	gs3 := utils.FirstOfPrevMonth().Format("2006-01-02")
	gs3s, ok := query["sd"]
	if ok && len(gs3s) > 0 {
		//case insensitive
		gs3 = strings.ToUpper(gs3s[0])
		//quotes
		re := regexp.MustCompile(`'`)
		gs3 = string(re.ReplaceAll([]byte(gs3), []byte("''")))
	}

	// gs4 := time.Now().Format("2006-01-02")
	gs4 := utils.LastOfPrevMonth().Format("2006-01-02")
	gs4s, ok := query["ed"]
	if ok && len(gs4s) > 0 {
		//case insensitive
		gs4 = strings.ToUpper(gs4s[0])
		//quotes
		re := regexp.MustCompile(`'`)
		gs4 = string(re.ReplaceAll([]byte(gs4), []byte("''")))
	}

	sum, err := gs.GetNodeSumL0(ctx, gs1, gs2, gs3, gs4)

	if err != nil {
		log.Println("Failed execute ifBalanceService.GetNodeSumL0: ", err)
	}

	out_count, err := json.Marshal(sum)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Write(out_count)

	return
}

// HandleBalanceTab1 godoc
// @Summary Get balance table 1 level down by node id, type and period
// @Description get balance table 1 level down by node id, type and period
// @Tags balance
// @Produce  json
// @Param page query int false "page number"
// @Param page_size query int false "page size"
// @Param id query int true "Main node id"
// @Param tid query int true "Main node type id"
// @Param sd query string false "Balance startdate. Default first day of previous month"
// @Param ed query string false "Balance enddate. Default last day of previous month"
// @Success 200 {object} models.BalanceTab_sum
// @Router /balance_tab1 [get]
func HandleBalanceTab1(w http.ResponseWriter, r *http.Request) {
	var gs ifBalanceService
	gs = services.NewBalanceService(pgsql.BalanceStorage{})
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

	gs1 := 0
	gs1s, ok := query["id"]
	if ok && len(gs1s) > 0 {
		t, err := strconv.Atoi(gs1s[0])
		if err == nil {
			gs1 = t
		}
	}

	gs2 := 0
	gs2s, ok := query["tid"]
	if ok && len(gs2s) > 0 {
		t, err := strconv.Atoi(gs2s[0])
		if err == nil {
			gs2 = t
		}
	}

	// gs3 := time.Now().Format("2006-01-02")
	gs3 := utils.FirstOfPrevMonth().Format("2006-01-02")
	gs3s, ok := query["sd"]
	if ok && len(gs3s) > 0 {
		//case insensitive
		gs3 = strings.ToUpper(gs3s[0])
		//quotes
		re := regexp.MustCompile(`'`)
		gs3 = string(re.ReplaceAll([]byte(gs3), []byte("''")))
	}

	// gs4 := time.Now().Format("2006-01-02")
	gs4 := utils.LastOfPrevMonth().Format("2006-01-02")
	gs4s, ok := query["ed"]
	if ok && len(gs4s) > 0 {
		//case insensitive
		gs4 = strings.ToUpper(gs4s[0])
		//quotes
		re := regexp.MustCompile(`'`)
		gs4 = string(re.ReplaceAll([]byte(gs4), []byte("''")))
	}

	out_arr, err := gs.GetTabL1(ctx, pg, pgs, gs1, gs2, gs3, gs4)

	if err != nil {
		log.Println("Failed execute ifBalanceService.GetTabL1: ", err)
	}

	out_count, err := json.Marshal(out_arr)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Write(out_count)
	return
}

// HandleBalanceTab0 godoc
// @Summary Get endpoints balance table by node id, type and period
// @Description get endpoints balance table by node id, type and period
// @Tags balance
// @Produce  json
// @Param page query int false "page number"
// @Param page_size query int false "page size"
// @Param id query int true "Main node id"
// @Param tid query int true "Main node type id"
// @Param sd query string false "Balance startdate. Default first day of previous month"
// @Param ed query string false "Balance enddate. Default last day of previous month"
// @Success 200 {object} models.BalanceTab_sum
// @Router /balance_tab0 [get]
func HandleBalanceTab0(w http.ResponseWriter, r *http.Request) {
	var gs ifBalanceService
	gs = services.NewBalanceService(pgsql.BalanceStorage{})
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

	gs1 := 0
	gs1s, ok := query["id"]
	if ok && len(gs1s) > 0 {
		t, err := strconv.Atoi(gs1s[0])
		if err == nil {
			gs1 = t
		}
	}

	gs2 := 0
	gs2s, ok := query["tid"]
	if ok && len(gs2s) > 0 {
		t, err := strconv.Atoi(gs2s[0])
		if err == nil {
			gs2 = t
		}
	}

	gs3 := utils.FirstOfPrevMonth().Format("2006-01-02")
	gs3s, ok := query["sd"]
	if ok && len(gs3s) > 0 {
		//case insensitive
		gs3 = strings.ToUpper(gs3s[0])
		//quotes
		re := regexp.MustCompile(`'`)
		gs3 = string(re.ReplaceAll([]byte(gs3), []byte("''")))
	}

	gs4 := utils.LastOfPrevMonth().Format("2006-01-02")
	gs4s, ok := query["ed"]
	if ok && len(gs4s) > 0 {
		//case insensitive
		gs4 = strings.ToUpper(gs4s[0])
		//quotes
		re := regexp.MustCompile(`'`)
		gs4 = string(re.ReplaceAll([]byte(gs4), []byte("''")))
	}

	out_arr, err := gs.GetTabL0(ctx, pg, pgs, gs1, gs2, gs3, gs4)

	if err != nil {
		log.Println("Failed execute ifBalanceService.GetTabL0: ", err)
	}

	out_count, err := json.Marshal(out_arr)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Write(out_count)
	return
}

// HandleBalanceBranch godoc
// @Summary Get branch by node id and type
// @Description get Get branch by node id and type
// @Tags balance
// @Produce  json
// @Param id query int true "Node id"
// @Param tid query int true "Node type id"
// @Success 200 {object} models.BalanceTab_sum
// @Router /balance_branch [get]
func HandleBalanceBranch(w http.ResponseWriter, r *http.Request) {
	var gs ifBalanceService
	gs = services.NewBalanceService(pgsql.BalanceStorage{})
	ctx := context.Background()

	query := r.URL.Query()

	gs1 := 0
	gs1s, ok := query["id"]
	if ok && len(gs1s) > 0 {
		t, err := strconv.Atoi(gs1s[0])
		if err == nil {
			gs1 = t
		}
	}

	gs2 := 0
	gs2s, ok := query["tid"]
	if ok && len(gs2s) > 0 {
		t, err := strconv.Atoi(gs2s[0])
		if err == nil {
			gs2 = t
		}
	}

	out_arr, err := gs.GetBranch(ctx, gs1, gs2)

	if err != nil {
		log.Println("Failed execute ifBalanceService.GetBranch: ", err)
	}

	out_count, err := json.Marshal(out_arr)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Write(out_count)
	return
}
