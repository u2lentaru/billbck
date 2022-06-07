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

// HandleOrgInfos godoc
// @Summary List org_info
// @Description Get org_info list
// @Tags org_info
// @Produce  json
// @Param page query int false "page number"
// @Param page_size query int false "page size"
// @Param oiname query string false "oiname search pattern"
// @Param oifname query string false "oifname search pattern"
// @Param ordering query string false "order by {oiname|oifname}"
// @Param desc query boolean false "descending order {true|false}"
// @Success 200 {array} models.OrgInfo_count
// @Failure 500
// @Router /org_info [get]
func (s *APG) HandleOrgInfos(w http.ResponseWriter, r *http.Request) {
	oi := models.OrgInfo{}
	ctx := context.Background()
	out_arr := []models.OrgInfo{}

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

	oin := ""
	oins, ok := query["oiname"]
	if ok && len(oins) > 0 {
		//case insensitive
		oin = strings.ToUpper(oins[0])
		//quotes
		re := regexp.MustCompile(`'`)
		oin = string(re.ReplaceAll([]byte(oin), []byte("''")))
	}

	oifn := ""
	oifns, ok := query["oifname"]
	if ok && len(oifns) > 0 {
		oifn = strings.ToUpper(oifns[0])
		re := regexp.MustCompile(`'`)
		oifn = string(re.ReplaceAll([]byte(oifn), []byte("''")))
	}

	ord := 1
	ords, ok := query["ordering"]
	if !ok || len(ords) == 0 {
		ord = 1
	} else if ords[0] == "oiname" {
		ord = 2
	} else if ords[0] == "oifname" {
		ord = 7
	}

	dsc := false
	dscs, ok := query["desc"]
	if ok && len(dscs) > 0 {
		if !(dscs[0] == "0") {
			dsc = true
		}
	}

	rows, err := s.Dbpool.Query(ctx, "SELECT * from func_orgs_info_get($1,$2,$3,$4,$5,$6);", pg, pgs, oin, oifn, ord, dsc)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&oi.Id, &oi.OIName, &oi.OIBin, &oi.OIAddr, &oi.OIBank.Id, &oi.OIAccNumber, &oi.OIFName, &oi.OIBank.BankName)
		if err != nil {
			log.Println("failed to scan row:", err)
		}

		out_arr = append(out_arr, oi)
	}

	oic := 0
	err = s.Dbpool.QueryRow(ctx, "SELECT * from func_orgs_info_cnt($1,$2);", oin, oifn).Scan(&oic)

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	auth := models.Auth{Create: true, Read: true, Update: true, Delete: true}
	out_count, err := json.Marshal(models.OrgInfo_count{Values: out_arr, Count: oic, Auth: auth})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Write(out_count)

	return
}

// HandleAddOrgInfo godoc
// @Summary Add org_info
// @Description Add org_info
// @Tags org_info
// @Accept json
// @Produce  json
// @Param ab body models.AddOrgInfo true "New org_info"
// @Success 200 {array} models.Json_id
// @Failure 500
// @Router /org_info_add [post]
func (s *APG) HandleAddOrgInfo(w http.ResponseWriter, r *http.Request) {
	aoi := models.AddOrgInfo{}
	body, err := ioutil.ReadAll(r.Body)

	defer r.Body.Close()

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	err = json.Unmarshal(body, &aoi)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	aoii := 0
	err = s.Dbpool.QueryRow(context.Background(), "SELECT func_orgs_info_add($1,$2,$3,$4,$5,$6);", aoi.OIName, aoi.OIBin, aoi.OIAddr, aoi.OIBank.Id, aoi.OIAccNumber, aoi.OIFName).Scan(&aoii)

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	output, err := json.Marshal(models.Json_id{Id: aoii})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(output)

	return
}

// HandleUpdOrgInfo godoc
// @Summary Update org_info
// @Description Update org_info
// @Tags org_info
// @Accept json
// @Produce  json
// @Param ub body models.OrgInfo true "Update org_info"
// @Success 200 {array} models.Json_id
// @Failure 500
// @Router /org_info_upd [post]
func (s *APG) HandleUpdOrgInfo(w http.ResponseWriter, r *http.Request) {
	uoi := models.OrgInfo{}
	body, err := ioutil.ReadAll(r.Body)

	defer r.Body.Close()

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	err = json.Unmarshal(body, &uoi)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	uoii := 0
	err = s.Dbpool.QueryRow(context.Background(), "SELECT func_orgs_info_upd($1,$2,$3,$4,$5,$6,$7);", uoi.Id, uoi.OIName, uoi.OIBin, uoi.OIAddr, uoi.OIBank.Id, uoi.OIAccNumber, uoi.OIFName).Scan(&uoii)

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	output, err := json.Marshal(models.Json_id{Id: uoii})

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(output)

	return
}

// HandleDelOrgInfo godoc
// @Summary Delete org_info
// @Description Delete org_info
// @Tags org_info
// @Accept json
// @Produce  json
// @Param db body models.Json_ids true "Delete org_info"
// @Success 200 {array} models.Json_id
// @Failure 500
// @Router /org_info_del [post]
func (s *APG) HandleDelOrgInfo(w http.ResponseWriter, r *http.Request) {
	doi := models.Json_ids{}
	body, err := ioutil.ReadAll(r.Body)

	defer r.Body.Close()

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	err = json.Unmarshal(body, &doi)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	res := []int{}
	oii := 0
	for _, id := range doi.Ids {
		err = s.Dbpool.QueryRow(context.Background(), "SELECT func_orgs_info_del($1);", id).Scan(&oii)
		res = append(res, oii)

		if err != nil {
			log.Println("Failed execute func_orgs_info_del: ", err)
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

// HandleGetOrgInfo godoc
// @Summary Get org_info
// @Description Get org_info
// @Tags org_info
// @Produce  json
// @Param id path int true "OrgInfo by id"
// @Success 200 {array} models.OrgInfo_count
// @Failure 500
// @Router /org_info/{id} [get]
func (s *APG) HandleGetOrgInfo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	oii := vars["id"]
	oi := models.OrgInfo{}
	out_arr := []models.OrgInfo{}

	err := s.Dbpool.QueryRow(context.Background(), "SELECT * from func_org_info_get($1);", oii).Scan(&oi.Id, &oi.OIName, &oi.OIBin, &oi.OIAddr, &oi.OIBank.Id, &oi.OIAccNumber, &oi.OIFName, &oi.OIBank.BankName)

	if err != nil && err != pgx.ErrNoRows {
		log.Println("Failed execute from func_org_info_get: ", err)
		// http.Error(w, err.Error(), 500)
		// return
	}

	out_arr = append(out_arr, oi)
	auth := models.Auth{Create: true, Read: true, Update: true, Delete: true}

	// output, err := json.Marshal(oi)
	out_count, err := json.Marshal(models.OrgInfo_count{Values: out_arr, Count: 1, Auth: auth})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(out_count)

	return
}
