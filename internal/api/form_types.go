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
// @Success 200 {array} models.FormType_count
// @Failure 400,404
// @Failure 500
// @Router /form_types [get]
func (s *APG) HandleFormTypes(w http.ResponseWriter, r *http.Request) {
	// type FormType struct {
	// 	Id            int    `json:"id"`
	// 	FormTypeName  string `json:"formtypename"`
	// 	FormTypeDescr string `json:"formtypedescr"`
	// }
	// ft := FormType{}

	// utils.SetupResponse(&w)

	// if (*r).Method == "OPTIONS" {
	// 	return
	// }

	ft := models.FormType{}
	ctx := context.Background()
	out_arr := []models.FormType{}

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

	rows, err := s.Dbpool.Query(ctx, "SELECT * from func_get_form_types_flt($1,$2,$3,$4,$5,$6);", pg, pgs, ftn, ftd, ord, dsc)
	// rows, err := s.dbpool.Query(ctx, "SELECT * from func_get_form_types();")
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&ft.Id, &ft.FormTypeName, &ft.FormTypeDescr)
		if err != nil {
			log.Println("failed to scan row:", err)
		}

		out_arr = append(out_arr, ft)
	}

	// output, err := json.Marshal(out_arr)
	// if err != nil {
	// 	http.Error(w, err.Error(), 500)
	// 	return
	// }
	// w.Write(output)

	ftc := 0
	// err = s.dbpool.QueryRow(ctx, "SELECT count(*) from st_form_types;").Scan(&ftc)
	// err = s.dbpool.QueryRow(ctx, "SELECT * from func_cnt_form_types_flt($1,$2,$3,$4,$5,$6);", pg, pgs, ftn, ftd, ord, dsc).Scan(&ftc)
	err = s.Dbpool.QueryRow(ctx, "SELECT * from func_cnt_form_types_flt($1,$2);", ftn, ftd).Scan(&ftc)

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	auth := models.Auth{Create: true, Read: true, Update: true, Delete: true}
	out_count, err := json.Marshal(models.FormType_count{Values: out_arr, Count: ftc, Auth: auth})
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
// @Success 200 {array} models.Json_id
// @Failure 400,404
// @Failure 500
// @Router /form_types_add [post]
func (s *APG) HandleAddFormType(w http.ResponseWriter, r *http.Request) {
	//FormType struct AddFormType
	type FormType struct {
		FormTypeName  string `json:"formtypename"`
		FormTypeDescr string `json:"formtypedescr"`
	}

	// utils.SetupResponse(&w)

	// if (*r).Method == "OPTIONS" {
	// 	return
	// }

	ft := FormType{}
	body, err := ioutil.ReadAll(r.Body)

	defer r.Body.Close()

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	err = json.Unmarshal(body, &ft)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	fti := 0
	err = s.Dbpool.QueryRow(context.Background(), "SELECT func_add_form_type($1,$2);", ft.FormTypeName, ft.FormTypeDescr).Scan(&fti)

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	output, err := json.Marshal(models.Json_id{Id: fti})
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
// @Success 200 {array} models.Json_id
// @Failure 500
// @Router /form_types_upd [post]
func (s *APG) HandleUpdFormType(w http.ResponseWriter, r *http.Request) {
	//FormType struct UpdFormType
	type FormType struct {
		Id            int    `json:"id"`
		FormTypeName  string `json:"formtypename"`
		FormTypeDescr string `json:"formtypedescr"`
	}

	// utils.SetupResponse(&w)

	// if (*r).Method == "OPTIONS" {
	// 	return
	// }

	ft := FormType{}
	body, err := ioutil.ReadAll(r.Body)

	defer r.Body.Close()

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	err = json.Unmarshal(body, &ft)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	fti := 0
	err = s.Dbpool.QueryRow(context.Background(), "SELECT func_upd_form_type($1,$2,$3);", ft.Id, ft.FormTypeName, ft.FormTypeDescr).Scan(&fti)

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	output, err := json.Marshal(models.Json_id{Id: fti})

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
// @Success 200 {array} models.Json_ids
// @Failure 500
// @Router /form_types_del [post]
func (s *APG) HandleDelFormType(w http.ResponseWriter, r *http.Request) {
	// utils.SetupResponse(&w)

	// if (*r).Method == "OPTIONS" {
	// 	return
	// }

	ft := models.Json_ids{}
	body, err := ioutil.ReadAll(r.Body)

	defer r.Body.Close()

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	err = json.Unmarshal(body, &ft)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	res := []int{}
	fti := 0
	for _, id := range ft.Ids {
		err = s.Dbpool.QueryRow(context.Background(), "SELECT func_del_form_type($1);", id).Scan(&fti)
		// log.Printf("id: %v fti : %v", id, fti)
		res = append(res, fti)

		if err != nil {
			log.Println("Failed execute func_del_form_type: ", err)
		}
	}
	// err = s.dbpool.QueryRow(context.Background(), "SELECT func_del_form_type($1);", ft.Id).Scan(&fti)

	// if err != nil {
	// 	http.Error(w, err.Error(), 500)
	// 	return
	// }

	// output, err := json.Marshal(models.Json_id{Id: fti})
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
// @Success 200 {array} models.FormType_count
// @Failure 500
// @Router /form_types/{id} [get]
func (s *APG) HandleGetFormType(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	fti := vars["id"]
	fts := models.FormType{}
	out_arr := []models.FormType{}

	// err := s.dbpool.QueryRow(context.Background(), "SELECT * from func_get_form_type($1);", ft.Id).Scan(&fts.Id, &fts.FormTypeName, &fts.FormTypeDescr)
	err := s.Dbpool.QueryRow(context.Background(), "SELECT * from func_get_form_type($1);", fti).Scan(&fts.Id, &fts.FormTypeName, &fts.FormTypeDescr)

	if err != nil && err != pgx.ErrNoRows {
		log.Println("Failed execute from func_get_form_type: ", err)
		// http.Error(w, err.Error(), 500)
		// return
	}

	out_arr = append(out_arr, fts)
	auth := models.Auth{Create: true, Read: true, Update: true, Delete: true}

	// output, err := json.Marshal(fts)
	out_count, err := json.Marshal(models.FormType_count{Values: out_arr, Count: 1, Auth: auth})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Write(out_count)

	return

}
