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
	"github.com/u2lentaru/billbck/internal/utils"
)

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
// @Success 200 {array} models.SubBank_count
// @Failure 500
// @Router /sub_banks [get]
func (s *APG) HandleSubBanks(w http.ResponseWriter, r *http.Request) {
	sb := models.SubBank{}
	ctx := context.Background()
	out_arr := []models.SubBank{}

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

	rows, err := s.Dbpool.Query(ctx, "SELECT * from func_sub_banks_get($1,$2,$3,$4,$5,$6,$7);", pg, pgs, sbn, utils.NullableInt(int32(sbi)), an, ord, dsc)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&sb.Id, &sb.Sub.Id, &sb.Bank.Id, &sb.AccNumber, &sb.Active, &sb.Sub.SBName, &sb.Bank.BankName)
		if err != nil {
			log.Println("failed to scan row:", err)
		}

		out_arr = append(out_arr, sb)
	}

	sbc := 0
	err = s.Dbpool.QueryRow(ctx, "SELECT * from func_sub_banks_cnt($1,$2,$3);", sbn, utils.NullableInt(int32(sbi)), an).Scan(&sbc)

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	auth := models.Auth{Create: true, Read: true, Update: true, Delete: true}
	out_count, err := json.Marshal(models.SubBank_count{Values: out_arr, Count: sbc, Auth: auth})
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
// @Success 200 {array} models.Json_id
// @Failure 500
// @Router /sub_banks_add [post]
func (s *APG) HandleAddSubBank(w http.ResponseWriter, r *http.Request) {
	asb := models.AddSubBank{}
	body, err := ioutil.ReadAll(r.Body)

	defer r.Body.Close()

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	err = json.Unmarshal(body, &asb)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	asbi := 0
	err = s.Dbpool.QueryRow(context.Background(), "SELECT func_sub_banks_add($1,$2,$3);", asb.Sub.Id, asb.Bank.Id, asb.AccNumber).Scan(&asbi)

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	output, err := json.Marshal(models.Json_id{Id: asbi})
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
// @Success 200 {array} models.Json_id
// @Failure 500
// @Router /sub_banks_upd [post]
func (s *APG) HandleUpdSubBank(w http.ResponseWriter, r *http.Request) {
	usb := models.SubBank{}
	body, err := ioutil.ReadAll(r.Body)

	defer r.Body.Close()

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	err = json.Unmarshal(body, &usb)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	usbi := 0
	err = s.Dbpool.QueryRow(context.Background(), "SELECT func_sub_banks_upd($1,$2,$3,$4);", usb.Id, usb.Sub.Id, usb.Bank.Id,
		usb.AccNumber).Scan(&usbi)

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	output, err := json.Marshal(models.Json_id{Id: usbi})

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
// @Success 200 {array} models.Json_ids
// @Failure 500
// @Router /sub_banks_del [post]
func (s *APG) HandleDelSubBank(w http.ResponseWriter, r *http.Request) {
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
		err = s.Dbpool.QueryRow(context.Background(), "SELECT func_sub_banks_del($1);", id).Scan(&i)
		res = append(res, i)

		if err != nil {
			log.Println("Failed execute func_sub_banks_del: ", err)
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

// HandleGetSubBank godoc
// @Summary Get subject account
// @Description get subject account
// @Tags subject banks
// @Produce  json
// @Param id path int true "Subject account by id"
// @Success 200 {array} models.SubBank_count
// @Failure 500
// @Router /sub_banks/{id} [get]
func (s *APG) HandleGetSubBank(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	i := vars["id"]
	g := models.SubBank{}
	out_arr := []models.SubBank{}

	err := s.Dbpool.QueryRow(context.Background(), "SELECT * from func_sub_bank_get($1);", i).Scan(&g.Id, &g.Sub.Id, &g.Bank.Id,
		&g.AccNumber, &g.Active, &g.Sub.SBName, &g.Bank.BankName)

	if err != nil && err != pgx.ErrNoRows {
		log.Println("Failed execute from func_sub_bank_get: ", err)
	}

	out_arr = append(out_arr, g)
	auth := models.Auth{Create: true, Read: true, Update: true, Delete: true}

	// output, err := json.Marshal(g)
	out_count, err := json.Marshal(models.SubBank_count{Values: out_arr, Count: 1, Auth: auth})
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
// @Success 200 {array} models.Json_id
// @Failure 500
// @Router /sub_banks_setactive/{id} [post]
func (s *APG) HandleGetSubBankSetActive(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	i := vars["id"]

	sai := 0
	err := s.Dbpool.QueryRow(context.Background(), "SELECT func_sub_banks_set_active($1);", i).Scan(&sai)

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	output, err := json.Marshal(models.Json_id{Id: sai})

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(output)

	return

}
