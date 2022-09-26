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

type ifFormTypeService interface {
	GetList(ctx context.Context, pg, pgs int, gs1, gs2 string, ord int, dsc bool) (models.FormType_count, error)
	Add(ctx context.Context, ea models.FormType) (int, error)
	Upd(ctx context.Context, eu models.FormType) (int, error)
	Del(ctx context.Context, ed []int) ([]int, error)
	GetOne(ctx context.Context, i int) (models.FormType_count, error)
}

// HandleFormTypes godoc
// @Summary List form types
// @Description get form types
// @Tags form types
// @Produce  json
// @Param page query int false "page number"
// @Param page_size query int false "page size"
// @Param formtypename query string false "formtypename search pattern"
// @Param formtypedescr query string false "formtypedescr search pattern"
// @Param ordering query string false "order by {formtypename|formtypedescr}"
// @Param desc query boolean false "descending order {true|false}"
// @Success 200 {object} models.FormType_count
// @Failure 400,404
// @Failure 500
// @Router /form_types [get]
func HandleFormTypes(w http.ResponseWriter, r *http.Request) {
	var gs ifFormTypeService
	gs = services.NewFormTypeService(pgsql.FormTypeStorage{})
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

	ftn := ""
	ftns, ok := query["formtypename"]
	if ok && len(ftns) > 0 {
		//case insensitive
		ftn = strings.ToUpper(ftns[0])
		//quotes
		re := regexp.MustCompile(`'`)
		ftn = string(re.ReplaceAll([]byte(ftn), []byte("''")))

	}

	ftd := ""
	ftds, ok := query["formtypedescr"]
	if ok && len(ftds) > 0 {
		ftd = strings.ToUpper(ftds[0])
		re := regexp.MustCompile(`'`)
		ftd = string(re.ReplaceAll([]byte(ftd), []byte("''")))
	}

	ord := 1
	ords, ok := query["ordering"]
	if !ok || len(ords) == 0 {
		ord = 1
	} else if ords[0] == "formtypename" {
		ord = 2
	} else if ords[0] == "formtypedescr" {
		ord = 3
	}

	dsc := false
	dscs, ok := query["desc"]
	if ok && len(dscs) > 0 {
		if !(dscs[0] == "0") {
			dsc = true
		}
	}

	out_arr, err := gs.GetList(ctx, pg, pgs, ftn, ftd, ord, dsc)
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

// HandleAddFormType godoc
// @Summary Add form type
// @Description add form type
// @Tags form types
// @Accept  json
// @Produce  json
// @Param ft body models.FormType true "New FormType"
// @Success 200 {object} models.Json_id
// @Failure 400,404
// @Failure 500
// @Router /form_types_add [post]
func HandleAddFormType(w http.ResponseWriter, r *http.Request) {
	var gs ifFormTypeService
	gs = services.NewFormTypeService(pgsql.FormTypeStorage{})
	ctx := context.Background()

	a := models.FormType{}
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
		log.Println("Failed execute ifFormTypeService.Add: ", err)
	}

	output, err := json.Marshal(models.Json_id{Id: ai})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(output)

	return
}

// HandleUpdFormType godoc
// @Summary Update form types
// @Description update form types
// @Tags form types
// @Accept  json
// @Produce  json
// @Param ft body models.FormType true "Update formtype"
// @Success 200 {object} models.Json_id
// @Failure 500
// @Router /form_types_upd [post]
func HandleUpdFormType(w http.ResponseWriter, r *http.Request) {
	var gs ifFormTypeService
	gs = services.NewFormTypeService(pgsql.FormTypeStorage{})
	ctx := context.Background()

	u := models.FormType{}
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
		log.Println("Failed execute ifFormTypeService.Upd: ", err)
	}

	output, err := json.Marshal(models.Json_id{Id: ui})

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(output)

	return

}

// HandleDelFormType godoc
// @Summary Delete form types
// @Description delete form types
// @Tags form types
// @Accept json
// @Produce  json
// @Param ft body models.Json_ids true "Delete formtype"
// @Success 200 {object} models.Json_ids
// @Failure 500
// @Router /form_types_del [post]
func HandleDelFormType(w http.ResponseWriter, r *http.Request) {
	var gs ifFormTypeService
	gs = services.NewFormTypeService(pgsql.FormTypeStorage{})
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
		log.Println("Failed execute ifFormTypeService.Del: ", err)
	}

	output, err := json.Marshal(models.Json_ids{Ids: res})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(output)

	return
}

// HandleGetFormType godoc
// @Summary Form type
// @Description get form type
// @Tags form types
// @Produce  json
// @Param id path int true "Form type by id"
// @Success 200 {object} models.FormType_count
// @Failure 500
// @Router /form_types/{id} [get]
func HandleGetFormType(w http.ResponseWriter, r *http.Request) {
	var gs ifFormTypeService
	gs = services.NewFormTypeService(pgsql.FormTypeStorage{})
	ctx := context.Background()
	auth := utils.GetAuth(r)

	vars := mux.Vars(r)
	i, err := strconv.Atoi(vars["id"])
	if err != nil {
		i = 0
	}

	out_arr, err := gs.GetOne(ctx, i)
	if err != nil {
		log.Println("Failed execute ifFormTypeService.GetOne: ", err)
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
