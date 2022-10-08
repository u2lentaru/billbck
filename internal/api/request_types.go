package api

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v4"
	"github.com/u2lentaru/billbck/internal/models"
	"github.com/u2lentaru/billbck/internal/utils"
)

// HandleRequestTypes godoc
// @Summary List requesttypes
// @Description get requesttype list
// @Tags requesttypes
// @Produce json
// @Param page query int false "page number"
// @Param page_size query int false "page size"
// @Param requesttypename query string false "requesttypename search pattern"
// @Param rkid query string false "request kind id search pattern"
// @Param ordering query string false "order by {id|requesttypename}"
// @Param desc query boolean false "descending order {true|false}"
// @Success 200 {object} models.RequestType_count
// @Failure 500
// @Router /requesttypes [get]
func (s *APG) HandleRequestTypes(w http.ResponseWriter, r *http.Request) {
	gs := models.RequestType{}
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

	// gs1 := ""
	// gs1s, ok := query["requesttypename"]
	// if ok && len(gs1s) > 0 {
	// 	//case insensitive
	// 	gs1 = strings.ToUpper(gs1s[0])
	// 	//quotes
	// 	re := regexp.MustCompile(`'`)
	// 	gs1 = string(re.ReplaceAll([]byte(gs1), []byte("''")))
	// }

	gs1 := ""
	gs1s, ok := query["requesttypename"]
	if ok && len(gs1s) > 0 {
		gs1 = gs1s[0]
	}

	gs2 := ""
	gs2s, ok := query["rkid"]
	if ok && len(gs2s) > 0 {
		_, err := strconv.Atoi(gs2s[0])
		if err == nil {
			gs2 = gs2s[0]
		}
	}

	gsc := 0
	err := s.Dbpool.QueryRow(ctx, "SELECT * from func_request_types_cnt($1,$2);", gs1, utils.NullableString(gs2)).Scan(&gsc)

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	out_arr := make([]models.RequestType, 0,
		func() int {
			if gsc < pgs {
				return gsc
			} else {
				return pgs
			}
		}())

	ord := 1
	ords, ok := query["ordering"]
	if !ok || len(ords) == 0 {
		ord = 1
	} else if ords[0] == "requesttypename" {
		ord = 2
	}

	dsc := false
	dscs, ok := query["desc"]
	if ok && len(dscs) > 0 {
		if !(dscs[0] == "0") {
			dsc = true
		}
	}

	rows, err := s.Dbpool.Query(ctx, "SELECT * from func_request_types_get($1,$2,$3,$4,$5,$6);", pg, pgs, gs1, utils.NullableString(gs2), ord, dsc)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&gs.Id, &gs.RequestTypeName, &gs.RequestKind.Id, &gs.RequestKind.RequestKindName)
		if err != nil {
			log.Println("failed to scan row:", err)
		}

		out_arr = append(out_arr, gs)
	}

	auth := models.Auth{Create: true, Read: true, Update: true, Delete: true}
	out_count, err := json.Marshal(models.RequestType_count{Values: out_arr, Count: gsc, Auth: auth})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Write(out_count)

	return
}

// HandleAddRequestType godoc
// @Summary Add requesttype
// @Description add requesttype
// @Tags requesttypes
// @Accept json
// @Produce  json
// @Param a body models.AddRequestType true "New requesttype"
// @Success 200 {object} models.Json_id
// @Failure 500
// @Router /requesttypes_add [post]
func (s *APG) HandleAddRequestType(w http.ResponseWriter, r *http.Request) {
	a := models.AddRequestType{}
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

	ai := 0
	err = s.Dbpool.QueryRow(context.Background(), "SELECT func_request_types_add($1,$2);", a.RequestTypeName, a.RequestKind.Id).Scan(&ai)

	if err != nil {
		log.Println("Failed execute func_request_types_add: ", err)
	}

	output, err := json.Marshal(models.Json_id{Id: ai})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(output)

	return

}

// HandleUpdRequestType godoc
// @Summary Update requesttype
// @Description update requesttype
// @Tags requesttypes
// @Accept json
// @Produce  json
// @Param u body models.RequestType true "Update requesttype"
// @Success 200 {object} models.Json_id
// @Failure 500
// @Router /requesttypes_upd [post]
func (s *APG) HandleUpdRequestType(w http.ResponseWriter, r *http.Request) {
	u := models.RequestType{}
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

	ui := 0
	err = s.Dbpool.QueryRow(context.Background(), "SELECT func_request_types_upd($1,$2,$3);", u.Id, u.RequestTypeName, u.RequestKind.Id).Scan(&ui)

	if err != nil {
		log.Println("Failed execute func_request_types_upd: ", err)
	}

	output, err := json.Marshal(models.Json_id{Id: ui})

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(output)

	return
}

// HandleDelRequestType godoc
// @Summary Delete requesttypes
// @Description delete requesttypes
// @Tags requesttypes
// @Accept json
// @Produce  json
// @Param d body models.Json_ids true "Delete requesttypes"
// @Success 200 {object} models.Json_ids
// @Failure 500
// @Router /requesttypes_del [post]
func (s *APG) HandleDelRequestType(w http.ResponseWriter, r *http.Request) {
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
		err = s.Dbpool.QueryRow(context.Background(), "SELECT func_request_types_del($1);", id).Scan(&i)
		res = append(res, i)

		if err != nil {
			log.Println("Failed execute func_request_types_del: ", err)
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

// HandleGetRequestType godoc
// @Summary Get requesttype
// @Description get requesttype
// @Tags requesttypes
// @Produce  json
// @Param id path int true "RequestType by id"
// @Success 200 {object} models.RequestType_count
// @Failure 500
// @Router /requesttypes/{id} [get]
func (s *APG) HandleGetRequestType(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	i := vars["id"]
	g := models.RequestType{}
	out_arr := []models.RequestType{}

	err := s.Dbpool.QueryRow(context.Background(), "SELECT * from func_request_type_get($1);", i).Scan(&g.Id, &g.RequestTypeName,
		&g.RequestKind.Id, &g.RequestKind.RequestKindName)

	if err != nil && err != pgx.ErrNoRows {
		log.Println("Failed execute from func_request_type_get: ", err)
	}

	out_arr = append(out_arr, g)
	auth := models.Auth{Create: true, Read: true, Update: true, Delete: true}

	// output, err := json.Marshal(g)
	out_count, err := json.Marshal(models.RequestType_count{Values: out_arr, Count: 1, Auth: auth})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(out_count)

	return

}
