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
	"github.com/jackc/pgx/v4"
	"github.com/u2lentaru/billbck/internal/models"
	"github.com/u2lentaru/billbck/internal/utils"
)

// HandleBalance godoc
// @Summary Get balance nodes
// @Description get balance nodes
// @Tags balance
// @Produce  json
// @Param page query int false "page number"
// @Param page_size query int false "page size"
// @Param pid query int true "Balance nodes by pid"
// @Param tid query int true "Balance node type id"
// @Success 200 {object} models.Balance
// @Failure 500
// @Router /balance [get]
func (s *APG) HandleBalance(w http.ResponseWriter, r *http.Request) {
	gs := models.Balance{}
	ctx := context.Background()
	out_arr := []models.Balance{}

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

	rows, err := s.Dbpool.Query(ctx, "SELECT * from func_balance_get($1,$2,$3,$4);", pg, pgs, gs1, gs2)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&gs.Id, &gs.PId, &gs.BName, &gs.BTypeId, &gs.BTypeName, &gs.ChildCount, &gs.ReqId)
		if err != nil {
			log.Println("failed to scan row:", err)
		}

		out_arr = append(out_arr, gs)
	}

	gsc := 0
	err = s.Dbpool.QueryRow(ctx, "SELECT * from func_balance_cnt($1,$2);", gs1, gs2).Scan(&gsc)

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	out_count, err := json.Marshal(models.Balance_count{Values: out_arr, Count: gsc})
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
func (s *APG) HandleGetBalance(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	i := vars["id"]
	tid := vars["tid"]
	g := models.Balance{}

	err := s.Dbpool.QueryRow(context.Background(), "SELECT * from func_balance_getbyid($1,$2);", i, tid).Scan(&g.Id, &g.PId, &g.BName, &g.BTypeId,
		&g.BTypeName, &g.ChildCount, &g.ReqId)

	if err != nil && err != pgx.ErrNoRows {
		log.Println("Failed execute from func_balance_get: ", err)
	}
	output, err := json.Marshal(g)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(output)

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
func (s *APG) HandleBalanceSum(w http.ResponseWriter, r *http.Request) {
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

	// now := time.Now()
	// currentYear, currentMonth, _ := now.Date()
	// currentLocation := now.Location()

	// firstOfPrevMonth := time.Date(currentYear, currentMonth-1, 1, 0, 0, 0, 0, currentLocation)
	// lastOfPrevMonth := firstOfPrevMonth.AddDate(0, 1, -1)

	// if currentMonth == 1 {
	// 	firstOfPrevMonth = time.Date(currentYear-1, 12, 1, 0, 0, 0, 0, currentLocation)
	// 	lastOfPrevMonth = firstOfPrevMonth.AddDate(0, 1, -1)
	// }

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

	sum := 0.0
	err := s.Dbpool.QueryRow(ctx, "SELECT * from func_balance_sum($1,$2,$3,$4);", gs1, gs2, gs3, gs4).Scan(&sum)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	out_count, err := json.Marshal(models.Json_sum{Sum: float32(sum)})
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
func (s *APG) HandleBalanceSum1(w http.ResponseWriter, r *http.Request) {
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

	sum := 0.0
	err := s.Dbpool.QueryRow(ctx, "SELECT * from func_balance_sum1($1,$2,$3,$4);", gs1, gs2, gs3, gs4).Scan(&sum)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	out_count, err := json.Marshal(models.Json_sum{Sum: float32(sum)})
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
func (s *APG) HandleBalanceSum0(w http.ResponseWriter, r *http.Request) {
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

	sum := 0.0
	err := s.Dbpool.QueryRow(ctx, "SELECT * from func_balance_sum0($1,$2,$3,$4);", gs1, gs2, gs3, gs4).Scan(&sum)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	out_count, err := json.Marshal(models.Json_sum{Sum: float32(sum)})
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
func (s *APG) HandleBalanceTab1(w http.ResponseWriter, r *http.Request) {
	gs := models.BalanceTab{}
	ctx := context.Background()
	out_arr := []models.BalanceTab{}

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

	rows, err := s.Dbpool.Query(ctx, "SELECT * from func_balance_tab1($1,$2,$3,$4,$5,$6);", pg, pgs, gs1, gs2, gs3, gs4)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&gs.Id, &gs.PId, &gs.BName, &gs.BTypeId, &gs.BTypeName, &gs.Sum, &gs.ReqId)
		if err != nil {
			log.Println("failed to scan row:", err)
		}

		out_arr = append(out_arr, gs)
	}

	ogsc := -1.0
	// ogsc := 0.0
	// err = s.Dbpool.QueryRow(ctx, "SELECT * from func_balance_sum1($1,$2,$3,$4);", gs1, gs2, gs3, gs4).Scan(&ogsc)

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	//input node sum
	igsc := 0.0
	err = s.Dbpool.QueryRow(ctx, "SELECT * from func_balance_sum($1,$2,$3,$4);", gs1, gs2, gs3, gs4).Scan(&igsc)

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	tab1cnt := 0
	err = s.Dbpool.QueryRow(ctx, "SELECT * from func_balance_tab1_cnt($1,$2);", gs1, gs2).Scan(&tab1cnt)

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	out_count, err := json.Marshal(models.BalanceTab_sum{Values: out_arr, InSum: igsc, OutSum: ogsc, Count: tab1cnt})
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
func (s *APG) HandleBalanceTab0(w http.ResponseWriter, r *http.Request) {
	gs := models.BalanceTab{}
	ctx := context.Background()
	out_arr := []models.BalanceTab{}

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

	rows, err := s.Dbpool.Query(ctx, "SELECT * from func_balance_tab0($1,$2,$3,$4,$5,$6);", pg, pgs, gs1, gs2, gs3, gs4)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&gs.Id, &gs.PId, &gs.BName, &gs.BTypeId, &gs.BTypeName, &gs.Sum, &gs.ReqId)
		if err != nil {
			log.Println("failed to scan row:", err)
		}

		out_arr = append(out_arr, gs)
	}

	ogsc := -1.0
	// ogsc := 0.0
	// err = s.Dbpool.QueryRow(ctx, "SELECT * from func_balance_sum0($1,$2,$3,$4);", gs1, gs2, gs3, gs4).Scan(&ogsc)

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	//input node sum
	igsc := 0.0
	err = s.Dbpool.QueryRow(ctx, "SELECT * from func_balance_sum($1,$2,$3,$4);", gs1, gs2, gs3, gs4).Scan(&igsc)

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	tab0cnt := 0
	err = s.Dbpool.QueryRow(ctx, "SELECT * from func_balance_tab0_cnt($1,$2);", gs1, gs2).Scan(&tab0cnt)

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	out_count, err := json.Marshal(models.BalanceTab_sum{Values: out_arr, InSum: igsc, OutSum: ogsc, Count: tab0cnt})
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
func (s *APG) HandleBalanceBranch(w http.ResponseWriter, r *http.Request) {
	gs := models.BalanceTab{}
	ctx := context.Background()
	out_arr := []models.BalanceTab{}

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

	rows, err := s.Dbpool.Query(ctx, "SELECT * from func_balance_getbranch($1,$2);", gs1, gs2)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&gs.Id, &gs.PId, &gs.BName, &gs.BTypeId, &gs.BTypeName, &gs.Sum, &gs.ReqId)
		if err != nil {
			log.Println("failed to scan row:", err)
		}

		out_arr = append(out_arr, gs)
	}

	out_count, err := json.Marshal(models.BalanceTab_sum{Values: out_arr, InSum: 0, OutSum: 0, Count: 0})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Write(out_count)
	return
}
