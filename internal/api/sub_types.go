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
	"github.com/u2lentaru/billbck/internal/models"
)

// HandleSubTypes godoc
// @Summary List subjects types
// @Description get subjects types
// @Tags sub types
// @Produce  json
// @Param page query int false "page number"
// @Param page_size query int false "page size"
// @Param subtypename query string false "subtypename search pattern"
// @Param subtypedescr query string false "subtypedescr search pattern"
// @Param ordering query string false "order by {subtypename|subtypedescr}"
// @Param desc query boolean false "descending order {true|false}"
// @Success 200 {object} models.SubType_count
// @Failure 500
// @Router /sub_types [get]
func (s *APG) HandleSubTypes(w http.ResponseWriter, r *http.Request) {
	st := models.SubType{}
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

	stn := ""
	stns, ok := query["subtypename"]
	if ok && len(stns) > 0 {
		//case insensitive
		stn = strings.ToUpper(stns[0])
		//quotes
		re := regexp.MustCompile(`'`)
		stn = string(re.ReplaceAll([]byte(stn), []byte("''")))
	}

	std := ""
	stds, ok := query["subtypedescr"]
	if ok && len(stds) > 0 {
		std = strings.ToUpper(stds[0])
		re := regexp.MustCompile(`'`)
		std = string(re.ReplaceAll([]byte(std), []byte("''")))
	}

	stc := 0
	err := s.Dbpool.QueryRow(ctx, "SELECT * from func_sub_types_cnt($1,$2);", stn, std).Scan(&stc)

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	out_arr := make([]models.SubType, 0,
		func() int {
			if stc < pgs {
				return stc
			} else {
				return pgs
			}
		}())

	ord := 1
	ords, ok := query["ordering"]
	if !ok || len(ords) == 0 {
		ord = 1
	} else if ords[0] == "subtypename" {
		ord = 2
	} else if ords[0] == "subtypedescr" {
		ord = 3
	}

	dsc := false
	dscs, ok := query["desc"]
	if ok && len(dscs) > 0 {
		if !(dscs[0] == "0") {
			dsc = true
		}
	}

	rows, err := s.Dbpool.Query(ctx, "SELECT * from func_sub_types_get($1,$2,$3,$4,$5,$6);", pg, pgs, stn, std, ord, dsc)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&st.Id, &st.SubTypeName, &st.SubTypeDescr)
		if err != nil {
			log.Println("failed to scan row:", err)
		}

		out_arr = append(out_arr, st)
	}

	auth := models.Auth{Create: true, Read: true, Update: true, Delete: true}
	out_count, err := json.Marshal(models.SubType_count{Values: out_arr, Count: stc, Auth: auth})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Write(out_count)

	return

}

// HandleAddSubType godoc
// @Summary Add subjects types
// @Description add subjects types
// @Tags sub types
// @Accept json
// @Produce  json
// @Param ast body models.AddSubType true "New subType"
// @Success 200 {object} models.Json_id
// @Failure 500
// @Router /sub_types_add [post]
func (s *APG) HandleAddSubType(w http.ResponseWriter, r *http.Request) {
	ast := models.AddSubType{}
	body, err := ioutil.ReadAll(r.Body)

	defer r.Body.Close()

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	err = json.Unmarshal(body, &ast)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	asti := 0
	err = s.Dbpool.QueryRow(context.Background(), "SELECT func_sub_types_add($1,$2);", ast.SubTypeDescr, ast.SubTypeName).Scan(&asti)

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	output, err := json.Marshal(models.Json_id{Id: asti})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(output)

	return

}

// HandleUpdSubType godoc
// @Summary Update subjects types
// @Description update subjects types
// @Tags sub types
// @Accept json
// @Produce  json
// @Param ust body models.SubType true "Update subtype"
// @Success 200 {object} models.Json_id
// @Failure 500
// @Router /sub_types_upd [post]
func (s *APG) HandleUpdSubType(w http.ResponseWriter, r *http.Request) {
	ust := models.SubType{}
	body, err := ioutil.ReadAll(r.Body)

	defer r.Body.Close()

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	err = json.Unmarshal(body, &ust)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	usti := 0
	err = s.Dbpool.QueryRow(context.Background(), "SELECT func_sub_types_upd($1,$2,$3);", ust.Id, ust.SubTypeName, ust.SubTypeDescr).Scan(&usti)

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	output, err := json.Marshal(models.Json_id{Id: usti})

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(output)

	return

}

// HandleDelSubType godoc
// @Summary Delete subjects types
// @Description delete subjects types
// @Tags sub types
// @Accept json
// @Produce  json
// @Param dst body models.Json_ids true "Delete subtypes"
// @Success 200 {object} models.Json_ids
// @Failure 500
// @Router /sub_types_del [post]
func (s *APG) HandleDelSubType(w http.ResponseWriter, r *http.Request) {
	dst := models.Json_ids{}
	body, err := ioutil.ReadAll(r.Body)

	defer r.Body.Close()

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	err = json.Unmarshal(body, &dst)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	res := []int{}
	sti := 0
	for _, id := range dst.Ids {
		err = s.Dbpool.QueryRow(context.Background(), "SELECT func_sub_types_del($1);", id).Scan(&sti)
		res = append(res, sti)

		if err != nil {
			log.Println("Failed execute func_sub_types_del: ", err)
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

// HandleGetSubType godoc
// @Summary Subjects type
// @Description get subject type
// @Tags sub types
// @Produce  json
// @Param id path int true "Subjects type by id"
// @Success 200 {object} models.SubType_count
// @Failure 500
// @Router /sub_types/{id} [get]
func (s *APG) HandleGetSubType(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	sti := vars["id"]
	st := models.SubType{}
	out_arr := []models.SubType{}

	err := s.Dbpool.QueryRow(context.Background(), "SELECT * from func_sub_type_get($1);", sti).Scan(&st.Id, &st.SubTypeName, &st.SubTypeDescr)

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	out_arr = append(out_arr, st)
	auth := models.Auth{Create: true, Read: true, Update: true, Delete: true}

	// output, err := json.Marshal(st)
	out_count, err := json.Marshal(models.SubType_count{Values: out_arr, Count: 1, Auth: auth})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(out_count)

	return

}
