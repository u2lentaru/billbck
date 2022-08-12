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
	"github.com/jackc/pgx/v4"
	"github.com/u2lentaru/billbck/internal/models"
	"github.com/u2lentaru/billbck/internal/utils"
)

// HandleRequests godoc
// @Summary List of requests
// @Description get request list
// @Tags requests
// @Produce json
// @Param page query int false "page number"
// @Param page_size query int false "page size"
// @Param requestnumber query string false "requestnumber search pattern"
// @Param ordering query string false "order by {id|requestnumber}"
// @Param desc query boolean false "descending order {true|false}"
// @Success 200 {object} models.Request_count
// @Failure 500
// @Router /requests [get]
func (s *APG) HandleRequests(w http.ResponseWriter, r *http.Request) {
	gs := models.Request{}
	ctx := context.Background()
	out_arr := []models.Request{}

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

	gs1 := ""
	gs1s, ok := query["requestnumber"]
	if ok && len(gs1s) > 0 {
		//case insensitive
		gs1 = strings.ToUpper(gs1s[0])
		//quotes
		re := regexp.MustCompile(`'`)
		gs1 = string(re.ReplaceAll([]byte(gs1), []byte("''")))
	}

	ord := 1
	ords, ok := query["ordering"]
	if !ok || len(ords) == 0 {
		ord = 1
	} else if ords[0] == "requestnumber" {
		ord = 2
	}

	dsc := false
	dscs, ok := query["desc"]
	if ok && len(dscs) > 0 {
		if !(dscs[0] == "0") {
			dsc = true
		}
	}

	rows, err := s.Dbpool.Query(ctx, "SELECT * from func_requests_get($1,$2,$3,$4,$5);", pg, pgs, gs1, ord, dsc)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	defer rows.Close()

	var ci, cli, ai, oi sql.NullInt32
	var cn, cln, an, on sql.NullString

	for rows.Next() {
		err = rows.Scan(&gs.Id, &gs.RequestNumber, &gs.RequestDate, &ci, &gs.ServiceType.Id,
			&gs.RequestType.Id, &gs.RequestKind.Id, &cli, &gs.TermDate, &gs.Executive, &gs.Accept, &gs.Notes,
			&gs.Result.Id, &ai, &oi, &cn, &gs.ServiceType.ServiceTypeName, &gs.RequestType.RequestTypeName,
			&gs.RequestKind.RequestKindName, &cln, &gs.Result.ResultName, &an, &on)

		gs.Contract.Id = int(ci.Int32)
		gs.Contract.ContractNumber = cn.String
		gs.ClaimType.Id = int(cli.Int32)
		gs.ClaimType.ClaimTypeName = cln.String
		gs.Act.Id = int(ai.Int32)
		gs.Act.ActNumber = an.String
		gs.Object.Id = int(oi.Int32)
		gs.Object.ObjectName = on.String

		if err != nil {
			log.Println("failed to scan row:", err)
		}

		out_arr = append(out_arr, gs)
	}

	gsc := 0
	err = s.Dbpool.QueryRow(ctx, "SELECT * from func_requests_cnt($1);", gs1).Scan(&gsc)

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	auth := models.Auth{Create: true, Read: true, Update: true, Delete: true}
	out_count, err := json.Marshal(models.Request_count{Values: out_arr, Count: gsc, Auth: auth})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Write(out_count)

	return
}

// HandleAddRequest godoc
// @Summary Add request
// @Description add request
// @Tags requests
// @Accept json
// @Produce  json
// @Param a body models.AddRequest true "New request. Significant params: RequestNumber, RequestDate, Contract.Id(n), ServiceType.Id, RequestType.Id, RequestKind.Id, ClaimType.Id(n), TermDate, Executive, Accept, Notes(n), Result.Id, Act.Id(n), Object.Id(n)"
// @Success 200 {object} models.Json_id
// @Failure 500
// @Router /requests_add [post]
func (s *APG) HandleAddRequest(w http.ResponseWriter, r *http.Request) {
	a := models.AddRequest{}
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

	// log.Println("a: ", a)

	ai := 0
	err = s.Dbpool.QueryRow(context.Background(), "SELECT func_requests_add($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14);",
		a.RequestNumber, a.RequestDate, utils.NullableInt(int32(a.Contract.Id)), a.ServiceType.Id, a.RequestType.Id, a.RequestKind.Id,
		utils.NullableInt(int32(a.ClaimType.Id)), a.TermDate, a.Executive, a.Accept, a.Notes, a.Result.Id,
		utils.NullableInt(int32(a.Act.Id)), utils.NullableInt(int32(a.Object.Id))).Scan(&ai)

	if err != nil {
		log.Println("Failed execute func_requests_add: ", err)
	}

	output, err := json.Marshal(models.Json_id{Id: ai})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(output)

	return

}

// HandleUpdRequest godoc
// @Summary Update request
// @Description update request
// @Tags requests
// @Accept json
// @Produce  json
// @Param u body models.Request true "Update request. Significant params: Id, RequestNumber, RequestDate, Contract.Id(n), ServiceType.Id, RequestType.Id, RequestKind.Id, ClaimType.Id(n), TermDate, Executive, Accept, Notes(n), Result.Id, Act.Id(n), Object.Id(n)"
// @Success 200 {object} models.Json_id
// @Failure 500
// @Router /requests_upd [post]
func (s *APG) HandleUpdRequest(w http.ResponseWriter, r *http.Request) {
	u := models.Request{}
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
	err = s.Dbpool.QueryRow(context.Background(), "SELECT func_requests_upd($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15);",
		u.Id, u.RequestNumber, u.RequestDate, utils.NullableInt(int32(u.Contract.Id)), u.ServiceType.Id,
		u.RequestType.Id, u.RequestKind.Id, utils.NullableInt(int32(u.ClaimType.Id)), u.TermDate, u.Executive, u.Accept,
		u.Notes, u.Result.Id, utils.NullableInt(int32(u.Act.Id)), utils.NullableInt(int32(u.Object.Id))).Scan(&ui)

	if err != nil {
		log.Println("Failed execute func_requests_upd: ", err)
	}

	output, err := json.Marshal(models.Json_id{Id: ui})

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(output)

	return
}

// HandleDelRequest godoc
// @Summary Delete requests
// @Description delete requests
// @Tags requests
// @Accept json
// @Produce  json
// @Param d body models.Json_ids true "Delete requests"
// @Success 200 {object} models.Json_ids
// @Failure 500
// @Router /requests_del [post]
func (s *APG) HandleDelRequest(w http.ResponseWriter, r *http.Request) {
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
		err = s.Dbpool.QueryRow(context.Background(), "SELECT func_requests_del($1);", id).Scan(&i)
		res = append(res, i)

		if err != nil {
			log.Println("Failed execute func_requests_del: ", err)
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

// HandleGetRequest godoc
// @Summary Get request
// @Description get request
// @Tags requests
// @Produce  json
// @Param id path int true "Request by id"
// @Success 200 {object} models.Request_count
// @Failure 500
// @Router /requests/{id} [get]
func (s *APG) HandleGetRequest(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	i := vars["id"]
	g := models.Request{}
	out_arr := []models.Request{}

	var ci, cli, ai, oi sql.NullInt32
	var cn, cln, an, on sql.NullString

	err := s.Dbpool.QueryRow(context.Background(), "SELECT * from func_request_get($1);", i).Scan(&g.Id,
		&g.RequestNumber, &g.RequestDate, &ci, &g.ServiceType.Id, &g.RequestType.Id, &g.RequestKind.Id, &cli, &g.TermDate,
		&g.Executive, &g.Accept, &g.Notes, &g.Result.Id, &ai, &oi, &cn, &g.ServiceType.ServiceTypeName,
		&g.RequestType.RequestTypeName, &g.RequestKind.RequestKindName, &cln, &g.Result.ResultName, &an, &on)

	g.Contract.Id = int(ci.Int32)
	g.Contract.ContractNumber = cn.String
	g.ClaimType.Id = int(cli.Int32)
	g.ClaimType.ClaimTypeName = cln.String
	g.Act.Id = int(ai.Int32)
	g.Act.ActNumber = an.String
	g.Object.Id = int(oi.Int32)
	g.Object.ObjectName = on.String

	if err != nil && err != pgx.ErrNoRows {
		log.Println("Failed execute from func_request_get: ", err)
	}

	out_arr = append(out_arr, g)
	auth := models.Auth{Create: true, Read: true, Update: true, Delete: true}

	// output, err := json.Marshal(g)
	out_count, err := json.Marshal(models.Request_count{Values: out_arr, Count: 1, Auth: auth})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(out_count)

	return

}
