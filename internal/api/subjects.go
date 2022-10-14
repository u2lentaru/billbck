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

type ifSubjectService interface {
	GetList(ctx context.Context, pg, pgs int, gs1, gs2, gs3 string, ord int, dsc bool) (models.Subject_count, error)
	Add(ctx context.Context, ea models.Subject) (int, error)
	Upd(ctx context.Context, eu models.Subject) (int, error)
	Del(ctx context.Context, ed models.SubjectClose) (int, error)
	GetOne(ctx context.Context, i int) (models.Subject_count, error)
	GetHist(ctx context.Context, i int) ([]string, error)
}

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
func HandleSubjects(w http.ResponseWriter, r *http.Request) {
	var gs ifSubjectService
	gs = services.NewSubjectService(pgsql.SubjectStorage{})
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

	out_arr, err := gs.GetList(ctx, pg, pgs, sjn, sjd, hc, ord, dsc)
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
func HandleAddSubject(w http.ResponseWriter, r *http.Request) {
	var gs ifSubjectService
	gs = services.NewSubjectService(pgsql.SubjectStorage{})
	ctx := context.Background()

	a := models.Subject{}
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
		log.Println("Failed execute ifSubjectService.Add: ", err)
	}

	output, err := json.Marshal(models.Json_id{Id: ai})
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
func HandleUpdSubject(w http.ResponseWriter, r *http.Request) {
	var gs ifSubjectService
	gs = services.NewSubjectService(pgsql.SubjectStorage{})
	ctx := context.Background()

	u := models.Subject{}
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
		log.Println("Failed execute ifSubjectService.Upd: ", err)
	}

	output, err := json.Marshal(models.Json_id{Id: ui})

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
func HandleDelSubject(w http.ResponseWriter, r *http.Request) {
	var gs ifSubjectService
	gs = services.NewSubjectService(pgsql.SubjectStorage{})
	ctx := context.Background()

	d := models.SubjectClose{}
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

	res, err := gs.Del(ctx, d)
	if err != nil {
		log.Println("Failed execute ifSubjectService.Del: ", err)
	}

	output, err := json.Marshal(models.Json_id{Id: res})
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
func HandleGetSubject(w http.ResponseWriter, r *http.Request) {
	var gs ifSubjectService
	gs = services.NewSubjectService(pgsql.SubjectStorage{})
	ctx := context.Background()
	auth := utils.GetAuth(r)

	vars := mux.Vars(r)
	i, err := strconv.Atoi(vars["id"])
	if err != nil {
		i = 0
	}

	out_arr, err := gs.GetOne(ctx, i)
	if err != nil {
		log.Println("Failed execute ifSubjectService.GetOne: ", err)
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

// HandleGetSubjectHist godoc
// @Summary Subject history
// @Description get subject history
// @Tags subjects
// @Produce  json
// @Param id path int true "Subject history by id"
// @Success 200 {object} string
// @Failure 500
// @Router /subjects_hist/{id} [get]
func HandleGetSubjectHist(w http.ResponseWriter, r *http.Request) {
	var gs ifSubjectService
	gs = services.NewSubjectService(pgsql.SubjectStorage{})
	ctx := context.Background()

	vars := mux.Vars(r)
	i, err := strconv.Atoi(vars["id"])
	if err != nil {
		i = 0
	}

	out_arr, err := gs.GetHist(ctx, i)
	if err != nil {
		log.Println("Failed execute ifSubjectService.GetHist: ", err)
	}

	out_count, err := json.Marshal(out_arr)

	w.Write(out_count)

	return
}
