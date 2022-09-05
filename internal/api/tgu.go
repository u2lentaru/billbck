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

// HandleTgu godoc
// @Summary List tgu
// @Description get tgu list
// @Tags tgu
// @Produce  json
// @Param page query int false "page number"
// @Param page_size query int false "page size"
// @Param tguname query string false "tguname search pattern"
// @Param ttid query string false "tgu type id"
// @Param ordering query string false "order by {id|tguname}"
// @Param desc query boolean false "descending order {true|false}"
// @Success 200 {object} models.Tgu_count
// @Failure 500
// @Router /tgu [get]
func (s *APG) HandleTgu(w http.ResponseWriter, r *http.Request) {
	gs := models.Tgu{}
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

	gs1 := ""
	gs1s, ok := query["tguname"]
	if ok && len(gs1s) > 0 {
		//case insensitive
		gs1 = strings.ToUpper(gs1s[0])
		//quotes
		re := regexp.MustCompile(`'`)
		gs1 = string(re.ReplaceAll([]byte(gs1), []byte("''")))
	}

	gs2 := ""
	gs2s, ok := query["ttid"]
	if ok && len(gs2s) > 0 {
		_, err := strconv.Atoi(gs2s[0])
		if err == nil {
			gs2 = gs2s[0]
		}
	}

	gsc := 0
	err := s.Dbpool.QueryRow(ctx, "SELECT * from func_tgu_cnt($1,$2);", gs1, utils.NullableString(gs2)).Scan(&gsc)

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	out_arr := make([]models.Tgu, 0,
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
	} else if ords[0] == "tguname" {
		ord = 3
	}

	dsc := false
	dscs, ok := query["desc"]
	if ok && len(dscs) > 0 {
		if !(dscs[0] == "0") {
			dsc = true
		}
	}

	rows, err := s.Dbpool.Query(ctx, "SELECT * from func_tgu_get($1,$2,$3,$4,$5,$6);", pg, pgs, gs1, utils.NullableString(gs2), ord, dsc)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	defer rows.Close()

	// var iid, o1id, o2id, iv, o1v, o2v *int

	var iid, o1id, o2id, iv, o1v, o2v sql.NullInt32
	var ivn, ov1n, ov2n sql.NullString

	for rows.Next() {
		// err = rows.Scan(&gs.Id, &gs.PId, &gs.TguName, &gs.TguType.Id, &gs.InvNumber, &gs.InputVoltage.Id, &gs.OutputVoltage1.Id,
		// 	&gs.OutputVoltage2.Id, &gs.TguType.TguTypeName, &gs.InputVoltage.VoltageValue, &gs.OutputVoltage1.VoltageValue,
		// 	&gs.OutputVoltage2.VoltageValue)

		// err = rows.Scan(&gs.Id, &gs.PId, &gs.TguName, &gs.TguType.Id, &gs.InvNumber, &gs.IVId, &gs.OV1Id,
		// 	&gs.OV2Id, &gs.TguType.TguTypeName, &gs.IVV, &gs.OV1V, &gs.OV2V)

		err = rows.Scan(&gs.Id, &gs.PId, &gs.TguName, &gs.TguType.Id, &gs.InvNumber, &iid, &o1id,
			&o2id, &gs.TguType.TguTypeName, &iv, &o1v, &o2v, &ivn, &ov1n, &ov2n)

		gs.InputVoltage.Id = int(iid.Int32)
		gs.InputVoltage.VoltageName = ivn.String
		gs.InputVoltage.VoltageValue = int(iv.Int32)
		gs.OutputVoltage1.Id = int(o1id.Int32)
		gs.OutputVoltage1.VoltageName = ov1n.String
		gs.OutputVoltage1.VoltageValue = int(o1v.Int32)
		gs.OutputVoltage2.Id = int(o2id.Int32)
		gs.OutputVoltage2.VoltageName = ov2n.String
		gs.OutputVoltage2.VoltageValue = int(o2v.Int32)
		// gs.InputVoltage.VoltageValue = *iv
		// gs.OutputVoltage1.Id = *o1id
		// gs.OutputVoltage1.VoltageValue = *o1v
		// gs.OutputVoltage2.Id = *o2id
		// gs.OutputVoltage2.VoltageValue = *o2v

		if err != nil {
			log.Println("failed to scan row:", err)
		}

		out_arr = append(out_arr, gs)
	}

	auth := models.Auth{Create: true, Read: true, Update: true, Delete: true}
	out_count, err := json.Marshal(models.Tgu_count{Values: out_arr, Count: gsc, Auth: auth})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Write(out_count)

	return
}

// HandleAddTgu godoc
// @Summary Add tgu
// @Description add tgu
// @Tags tgu
// @Accept json
// @Produce  json
// @Param a body models.AddTgu true "New tgu"
// @Success 200 {object} models.Json_id
// @Failure 500
// @Router /tgu_add [post]
func (s *APG) HandleAddTgu(w http.ResponseWriter, r *http.Request) {
	a := models.AddTgu{}
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
	err = s.Dbpool.QueryRow(context.Background(), "SELECT func_tgu_add($1,$2,$3,$4,$5,$6,$7);", a.PId, a.TguName, a.TguType.Id,
		a.InvNumber, a.InputVoltage.Id, a.OutputVoltage1.Id, a.OutputVoltage2.Id).Scan(&ai)

	if err != nil {
		log.Println("Failed execute func_tgu_add: ", err)
	}

	output, err := json.Marshal(models.Json_id{Id: ai})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(output)

	return
}

// HandleUpdTgu godoc
// @Summary Update tgu
// @Description update tgu
// @Tags tgu
// @Accept json
// @Produce  json
// @Param u body models.Tgu true "Update tgu"
// @Success 200 {object} models.Json_id
// @Failure 500
// @Router /tgu_upd [post]
func (s *APG) HandleUpdTgu(w http.ResponseWriter, r *http.Request) {
	u := models.Tgu{}
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
	err = s.Dbpool.QueryRow(context.Background(), "SELECT func_tgu_upd($1,$2,$3,$4,$5,$6,$7,$8);", u.Id, u.PId, u.TguName, u.TguType.Id,
		u.InvNumber, u.InputVoltage.Id, u.OutputVoltage1.Id, u.OutputVoltage2.Id).Scan(&ui)

	// err = s.Dbpool.QueryRow(context.Background(), "SELECT func_tgu_upd($1,$2,$3,$4,$5,$6,$7,$8);", u.Id, u.PId, u.TguName, u.TguType.Id,
	// 	u.InvNumber, u.IVId, u.OV1Id, u.OV2Id).Scan(&ui)

	if err != nil {
		log.Println("Failed execute func_tgu_upd: ", err)
	}

	output, err := json.Marshal(models.Json_id{Id: ui})

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(output)

	return

}

// HandleDelTgu godoc
// @Summary Delete tgu
// @Description delete tgu
// @Tags tgu
// @Accept json
// @Produce  json
// @Param d body models.Json_ids true "Delete tgu"
// @Success 200 {object} models.Json_ids
// @Failure 500
// @Router /tgu_del [post]
func (s *APG) HandleDelTgu(w http.ResponseWriter, r *http.Request) {
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
		err = s.Dbpool.QueryRow(context.Background(), "SELECT func_tgu_del($1);", id).Scan(&i)
		res = append(res, i)

		if err != nil {
			log.Println("Failed execute func_tgu_del: ", err)
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

// HandleGetTgu godoc
// @Summary Get tgu
// @Description get tgu
// @Tags tgu
// @Produce  json
// @Param id path int true "Tgu by id"
// @Success 200 {object} models.Tgu_count
// @Failure 500
// @Router /tgu/{id} [get]
func (s *APG) HandleGetTgu(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	i := vars["id"]
	g := models.Tgu{}
	out_arr := []models.Tgu{}

	var iid, o1id, o2id, iv, o1v, o2v sql.NullInt32
	var ivn, ov1n, ov2n sql.NullString

	// err := s.Dbpool.QueryRow(context.Background(), "SELECT * from func_tgu_getbyid($1);", i).Scan(&g.Id, &g.PId, &g.TguName, &g.TguType.Id,
	// 	&g.InvNumber, &g.InputVoltage.Id, &g.OutputVoltage1.Id, &g.OutputVoltage2.Id, &g.TguType.TguTypeName, &g.InputVoltage.VoltageValue,
	// 	&g.OutputVoltage1.VoltageValue, &g.OutputVoltage2.VoltageValue)

	err := s.Dbpool.QueryRow(context.Background(), "SELECT * from func_tgu_getbyid($1);", i).Scan(&g.Id, &g.PId, &g.TguName, &g.TguType.Id,
		&g.InvNumber, &iid, &o1id, &o2id, &g.TguType.TguTypeName, &iv, &o1v, &o2v, &ivn, &ov1n, &ov2n)

	g.InputVoltage.Id = int(iid.Int32)
	g.InputVoltage.VoltageName = ivn.String
	g.InputVoltage.VoltageValue = int(iv.Int32)
	g.OutputVoltage1.Id = int(o1id.Int32)
	g.OutputVoltage1.VoltageName = ov1n.String
	g.OutputVoltage1.VoltageValue = int(o1v.Int32)
	g.OutputVoltage2.Id = int(o2id.Int32)
	g.OutputVoltage2.VoltageName = ov2n.String
	g.OutputVoltage2.VoltageValue = int(o2v.Int32)

	if err != nil && err != pgx.ErrNoRows {
		log.Println("Failed execute from func_tgu_getbyid: ", err)
	}

	out_arr = append(out_arr, g)
	auth := models.Auth{Create: true, Read: true, Update: true, Delete: true}

	// output, err := json.Marshal(g)
	out_count, err := json.Marshal(models.Tgu_count{Values: out_arr, Count: 1, Auth: auth})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(out_count)

	return

}
