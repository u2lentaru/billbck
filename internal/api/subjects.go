package api

import (
	"context"
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/u2lentaru/billbck/internal/models"
	"github.com/u2lentaru/billbck/internal/utils"
)

// HandleSubjects godoc
// @Summary List subjects
// @Description get subjects
// @Tags subjects
// @Produce  json
// @Param page query int false "page number"
// @Param page_size query int false "page size"
// @Param subname query string false "subname search pattern"
// @Param subdescr query string false "subdescr search pattern"
// @Param hideclosed query boolean false "hide closed, default true"
// @Param ordering query string false "order by {subname|subdescr}"
// @Param desc query boolean false "descending order {true|false}"
// @Success 200 {object} models.Subject_count
// @Failure 500
// @Router /subjects [get]
func (s *APG) HandleSubjects(w http.ResponseWriter, r *http.Request) {
	sj := models.Subject{}
	ctx := context.Background()
	out_arr := []models.Subject{}

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

	sjn := ""
	sjns, ok := query["subname"]
	if ok && len(sjns) > 0 {
		//case insensitive
		sjn = strings.ToUpper(sjns[0])
		//quotes
		re := regexp.MustCompile(`'`)
		sjn = string(re.ReplaceAll([]byte(sjn), []byte("''")))
	}

	sjd := ""
	sjds, ok := query["subdescr"]
	if ok && len(sjds) > 0 {
		sjd = strings.ToUpper(sjds[0])
		re := regexp.MustCompile(`'`)
		sjd = string(re.ReplaceAll([]byte(sjd), []byte("''")))
	}

	hc := "true"
	hcs, ok := query["hideclosed"]
	if ok && len(hcs) > 0 {
		hc = strings.ToLower(hcs[0])
		if hc != "false" {
			hc = "true"
		} else {
			hc = "false"
		}
	}

	ord := 1
	ords, ok := query["ordering"]
	if !ok || len(ords) == 0 {
		ord = 1
	} else if ords[0] == "subname" {
		ord = 5
	} else if ords[0] == "subdescr" {
		ord = 4
	}

	dsc := false
	dscs, ok := query["desc"]
	if ok && len(dscs) > 0 {
		if !(dscs[0] == "0") {
			dsc = true
		}
	}

	rows, err := s.Dbpool.Query(ctx, "SELECT * from func_sub_details_get($1,$2,$3,$4,$5,$6,$7);", pg, pgs, sjn, sjd, hc, ord, dsc)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	defer rows.Close()

	var shp, sap sql.NullInt32
	var shn, san sql.NullString

	for rows.Next() {
		err = rows.Scan(&sj.SubId, &sj.SubType.Id, &sj.SubPhys, &sj.SubDescr, &sj.SubName, &sj.SubBin, &shp, &sj.SubHeadName,
			&sap, &sj.SubAccName, &sj.SubAddr, &sj.SubPhone, &sj.SubStart, &sj.SubAccNumber, &sj.SubType.SubTypeName,
			&shn, &san, &sj.Job, &sj.Email, &sj.MobPhone, &sj.JobPhone, &sj.Notes)
		if err != nil {
			log.Println("failed to scan row:", err)
		}

		sj.SubHeadPos.Id = int(shp.Int32)
		sj.SubAccPos.Id = int(sap.Int32)
		sj.SubHeadPos.PositionName = shn.String
		sj.SubAccPos.PositionName = san.String

		out_arr = append(out_arr, sj)
	}

	sjc := 0
	err = s.Dbpool.QueryRow(ctx, "SELECT * from func_sub_details_cnt($1,$2,$3);", sjn, sjd, hc).Scan(&sjc)

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	auth := models.Auth{Create: true, Read: true, Update: true, Delete: true}
	out_count, err := json.Marshal(models.Subject_count{Values: out_arr, Count: sjc, Auth: auth})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Write(out_count)

	return

}

// HandleAddSubject godoc
// @Summary Add subject
// @Description add subject
// @Tags subjects
// @Accept json
// @Produce  json
// @Param asb body models.AddSubject true "New Subject. Significant params: SubType.Id, SubPhys, SubDescr, SubName, SubBin, SubHeadPos.Id(n), SubHeadName(n), SubAccPos.Id(n), SubAccName(n), SubAddr, SubPhone, SubStart, SubAccNumber, Job(n), Email(n), MobPhone(n), JobPhone(n), Notes(n)"
// @Success 200 {object} models.Json_id
// @Failure 500
// @Router /subjects_add [post]
func (s *APG) HandleAddSubject(w http.ResponseWriter, r *http.Request) {
	asb := models.AddSubject{}
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
	err = s.Dbpool.QueryRow(context.Background(), "SELECT func_sub_details_add($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18);",
		asb.SubType.Id, asb.SubPhys, asb.SubDescr, asb.SubName, asb.SubBin, utils.NullableInt(int32(asb.SubHeadPos.Id)), asb.SubHeadName,
		utils.NullableInt(int32(asb.SubAccPos.Id)), asb.SubAccName, asb.SubAddr, asb.SubPhone, asb.SubStart, asb.SubAccNumber, asb.Job,
		asb.Email, asb.MobPhone, asb.JobPhone, asb.Notes).Scan(&asbi)

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

// HandleUpdSubject godoc
// @Summary Update subject
// @Description update subject
// @Tags subjects
// @Accept  json
// @Produce  json
// @Param usb body models.Subject true "Update sybject. Significant params: SubId, SubType.Id, SubPhys, SubDescr, SubName, SubBin, SubHeadPos.Id(n), SubHeadName(n), SubAccPos.Id(n), SubAccName(n), SubAddr, SubPhone, SubStart, SubAccNumber, Job(n), Email(n), MobPhone(n), JobPhone(n), Notes(n)""
// @Success 200 {object} models.Json_id
// @Failure 500
// @Router /subjects_upd [post]
func (s *APG) HandleUpdSubject(w http.ResponseWriter, r *http.Request) {
	usb := models.Subject{}
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
	err = s.Dbpool.QueryRow(context.Background(), "SELECT func_sub_details_upd($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19);",
		usb.SubId, usb.SubType.Id, usb.SubPhys, usb.SubDescr, usb.SubName, usb.SubBin, utils.NullableInt(int32(usb.SubHeadPos.Id)),
		usb.SubHeadName, utils.NullableInt(int32(usb.SubAccPos.Id)), usb.SubAccName, usb.SubAddr, usb.SubPhone, usb.SubStart,
		usb.SubAccNumber, usb.Job, usb.Email, usb.MobPhone, usb.JobPhone, usb.Notes).Scan(&usbi)

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

// HandleDelSubject godoc
// @Summary Delete subjects
// @Description delete subjects
// @Tags subjects
// @Accept json
// @Produce  json
// @Param dsb body models.SubjectClose true "Close subject"
// @Success 200 {object} models.Json_id
// @Failure 500
// @Router /subjects_del [post]
func (s *APG) HandleDelSubject(w http.ResponseWriter, r *http.Request) {
	dsb := models.SubjectClose{}
	body, err := ioutil.ReadAll(r.Body)

	defer r.Body.Close()

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	err = json.Unmarshal(body, &dsb)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	sbi := 0
	err = s.Dbpool.QueryRow(context.Background(), "SELECT func_subjects_close($1,$2);", dsb.SubId, dsb.SubClose).Scan(&sbi)

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	output, err := json.Marshal(models.Json_id{Id: sbi})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(output)

	return
}

// HandleGetSubject godoc
// @Summary Subject
// @Description get subject
// @Tags subjects
// @Produce  json
// @Param id path int true "Subject by id"
// @Success 200 {object} models.Subject_count
// @Failure 500
// @Router /subjects/{id} [get]
func (s *APG) HandleGetSubject(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	sbi := vars["id"]
	sb := models.Subject{}
	out_arr := []models.Subject{}

	var shp, sap sql.NullInt32
	var shn, san sql.NullString

	err := s.Dbpool.QueryRow(context.Background(), "SELECT * from func_sub_detail_get($1);", sbi).Scan(&sb.SubId, &sb.SubType.Id, &sb.SubPhys,
		&sb.SubDescr, &sb.SubName, &sb.SubBin, &shp, &sb.SubHeadName, &sap, &sb.SubAccName, &sb.SubAddr, &sb.SubPhone,
		&sb.SubStart, &sb.SubAccNumber, &sb.SubType.SubTypeName, &shn, &san, &sb.Job, &sb.Email,
		&sb.MobPhone, &sb.JobPhone, &sb.Notes)

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	sb.SubHeadPos.Id = int(shp.Int32)
	sb.SubAccPos.Id = int(sap.Int32)
	sb.SubHeadPos.PositionName = shn.String
	sb.SubAccPos.PositionName = san.String

	out_arr = append(out_arr, sb)
	auth := models.Auth{Create: true, Read: true, Update: true, Delete: true}

	// output, err := json.Marshal(sb)
	out_count, err := json.Marshal(models.Subject_count{Values: out_arr, Count: 1, Auth: auth})

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Write(out_count)

	return
}

// HandleGetSubjectHist godoc
// @Summary Subject history
// @Description get subject history
// @Tags subjects
// @Produce  json
// @Param id path int true "Subject history by id"
// @Success 200 {object} string
// @Failure 500
// @Router /subjects_hist/{id} [get]
func (s *APG) HandleGetSubjectHist(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	sbi := vars["id"]
	// hist_arr := []models.Subject_hist{}
	hist_arr := []string{}

	// sb := models.Subject_hist{}
	h := ""
	rows, err := s.Dbpool.Query(context.Background(), "SELECT * from func_sub_detail_get_hist($1);", sbi)

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	defer rows.Close()

	qa := false
	w.Write([]byte("["))

	for rows.Next() {

		if qa {
			w.Write([]byte(","))
		}
		qa = true

		err = rows.Scan(&h)
		w.Write([]byte(h))

		if err != nil {
			log.Println("failed to scan row:", err)
		}
		hist_arr = append(hist_arr, h)
	}

	w.Write([]byte("]"))

	return
}
