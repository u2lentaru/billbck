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
func (s *APG) HandleBanks(w http.ResponseWriter, r *http.Request) {
	b := models.Bank{}
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

	bc := 0
	err := s.Dbpool.QueryRow(ctx, "SELECT * from func_banks_cnt($1,$2);", bn, bd).Scan(&bc)

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	out_arr := make([]models.Bank, 0,
		func() int {
			if bc < pgs {
				return bc
			} else {
				return pgs
			}
		}())

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

	rows, err := s.Dbpool.Query(ctx, "SELECT * from func_banks_get($1,$2,$3,$4,$5,$6);", pg, pgs, bn, bd, ord, dsc)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&b.Id, &b.BankName, &b.BankDescr, &b.Mfo)
		if err != nil {
			log.Println("failed to scan row:", err)
		}

		out_arr = append(out_arr, b)
	}

	auth := models.Auth{Create: true, Read: true, Update: true, Delete: true}
	out_count, err := json.Marshal(models.Bank_count{Values: out_arr, Count: bc, Auth: auth})
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
func (s *APG) HandleAddBank(w http.ResponseWriter, r *http.Request) {
	ab := models.AddBank{}
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

	abi := 0
	err = s.Dbpool.QueryRow(context.Background(), "SELECT func_banks_add($1,$2,$3);", ab.BankName, ab.BankDescr, ab.Mfo).Scan(&abi)

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	output, err := json.Marshal(models.Json_id{Id: abi})
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
func (s *APG) HandleUpdBank(w http.ResponseWriter, r *http.Request) {
	ub := models.Bank{}
	body, err := ioutil.ReadAll(r.Body)

	defer r.Body.Close()

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	err = json.Unmarshal(body, &ub)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	ubi := 0
	err = s.Dbpool.QueryRow(context.Background(), "SELECT func_banks_upd($1,$2,$3,$4);", ub.Id, ub.BankName, ub.BankDescr, ub.Mfo).Scan(&ubi)

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	output, err := json.Marshal(models.Json_id{Id: ubi})

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
func (s *APG) HandleDelBank(w http.ResponseWriter, r *http.Request) {
	db := models.Json_ids{}
	body, err := ioutil.ReadAll(r.Body)

	defer r.Body.Close()

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	err = json.Unmarshal(body, &db)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	res := []int{}
	bi := 0
	for _, id := range db.Ids {
		err = s.Dbpool.QueryRow(context.Background(), "SELECT func_banks_del($1);", id).Scan(&bi)
		res = append(res, bi)

		if err != nil {
			log.Println("Failed execute func_banks_del: ", err)
		}
	}

	// if err != nil {
	// 	http.Error(w, err.Error(), 500)
	// 	return
	// }

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
func (s *APG) HandleGetBank(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	bi := vars["id"]
	b := models.Bank{}
	out_arr := []models.Bank{}

	err := s.Dbpool.QueryRow(context.Background(), "SELECT * from func_bank_get($1);", bi).Scan(&b.Id, &b.BankName, &b.BankDescr, &b.Mfo)

	if err != nil && err != pgx.ErrNoRows {
		log.Println("Failed execute from func_bank_get: ", err)
		// http.Error(w, err.Error(), 500)
		// return
	}

	out_arr = append(out_arr, b)
	auth := models.Auth{Create: true, Read: true, Update: true, Delete: true}

	// output, err := json.Marshal(b)
	out_count, err := json.Marshal(models.Bank_count{Values: out_arr, Count: 1, Auth: auth})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(out_count)

	return
}
