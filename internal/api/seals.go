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

// HandleSeals godoc
// @Summary List seals
// @Description get seal list
// @Tags seals
// @Produce json
// @Param page query int false "page number"
// @Param page_size query int false "page size"
// @Param packetnumber query string false "packetnumber search pattern"
// @Param ordering query string false "order by {id|packetnumber}"
// @Param desc query boolean false "descending order {true|false}"
// @Success 200 {array} models.Seal_count
// @Failure 500
// @Router /seals [get]
func (s *APG) HandleSeals(w http.ResponseWriter, r *http.Request) {
	gs := models.Seal{}
	ctx := context.Background()
	out_arr := []models.Seal{}

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
	gs1s, ok := query["packetnumber"]
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
	} else if ords[0] == "packetnumber" {
		ord = 2
	}

	dsc := false
	dscs, ok := query["desc"]
	if ok && len(dscs) > 0 {
		if !(dscs[0] == "0") {
			dsc = true
		}
	}

	rows, err := s.Dbpool.Query(ctx, "SELECT * from func_seals_get($1,$2,$3,$4,$5);", pg, pgs, gs1, ord, dsc)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&gs.Id, &gs.PacketNumber, &gs.Area.Id, &gs.Staff.Id, &gs.SealType.Id, &gs.SealColour.Id, &gs.SealStatus.Id,
			&gs.IssueDate, &gs.ReportDate, &gs.Area.AreaName, &gs.Staff.StaffName, &gs.SealType.SealTypeName, &gs.SealColour.SealColourName,
			&gs.SealStatus.SealStatusName)
		if err != nil {
			log.Println("failed to scan row:", err)
		}

		out_arr = append(out_arr, gs)
	}

	gsc := 0
	err = s.Dbpool.QueryRow(ctx, "SELECT * from func_seals_cnt($1);", gs1).Scan(&gsc)

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	auth := models.Auth{Create: true, Read: true, Update: true, Delete: true}
	out_count, err := json.Marshal(models.Seal_count{Values: out_arr, Count: gsc, Auth: auth})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Write(out_count)

	return
}

// HandleAddSeal godoc
// @Summary Add seal
// @Description add seal
// @Tags seals
// @Accept json
// @Produce  json
// @Param a body models.AddSeal true "New seal. Significant params: PacketNumber, Area.Id, Staff.Id, SealType.Id, SealColour.Id, SealStatus.Id, IssueDate, ReportDate"
// @Success 200 {array} models.Json_id
// @Failure 500
// @Router /seals_add [post]
func (s *APG) HandleAddSeal(w http.ResponseWriter, r *http.Request) {
	a := models.AddSeal{}
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
	err = s.Dbpool.QueryRow(context.Background(), "SELECT func_seals_add($1,$2,$3,$4,$5,$6,$7,$8);", a.PacketNumber, a.Area.Id,
		a.Staff.Id, a.SealType.Id, a.SealColour.Id, a.SealStatus.Id, a.IssueDate, a.ReportDate).Scan(&ai)

	if err != nil {
		log.Println("Failed execute func_seals_add: ", err)
	}

	output, err := json.Marshal(models.Json_id{Id: ai})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(output)

	return

}

// HandleUpdSeal godoc
// @Summary Update seal
// @Description update seal
// @Tags seals
// @Accept json
// @Produce  json
// @Param u body models.Seal true "Update seal. Significant params: Id, PacketNumber, Area.Id, Staff.Id, SealType.Id, SealColour.Id, SealStatus.Id, IssueDate, ReportDate"
// @Success 200 {array} models.Json_id
// @Failure 500
// @Router /seals_upd [post]
func (s *APG) HandleUpdSeal(w http.ResponseWriter, r *http.Request) {
	u := models.Seal{}
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
	err = s.Dbpool.QueryRow(context.Background(), "SELECT func_seals_upd($1,$2,$3,$4,$5,$6,$7,$8,$9);", u.Id, u.PacketNumber, u.Area.Id,
		u.Staff.Id, u.SealType.Id, u.SealColour.Id, u.SealStatus.Id, u.IssueDate, u.ReportDate).Scan(&ui)

	if err != nil {
		log.Println("Failed execute func_seals_upd: ", err)
	}

	output, err := json.Marshal(models.Json_id{Id: ui})

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(output)

	return

}

// HandleDelSeal godoc
// @Summary Delete seals
// @Description delete seals
// @Tags seals
// @Accept json
// @Produce  json
// @Param d body models.Json_ids true "Delete seals"
// @Success 200 {array} models.Json_ids
// @Failure 500
// @Router /seals_del [post]
func (s *APG) HandleDelSeal(w http.ResponseWriter, r *http.Request) {
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
		err = s.Dbpool.QueryRow(context.Background(), "SELECT func_seals_del($1);", id).Scan(&i)
		res = append(res, i)

		if err != nil {
			log.Println("Failed execute func_seals_del: ", err)
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

// HandleGetSeal godoc
// @Summary Get seal
// @Description get seal
// @Tags seals
// @Produce  json
// @Param id path int true "Seal by id"
// @Success 200 {array} models.Seal_count
// @Failure 500
// @Router /seals/{id} [get]
func (s *APG) HandleGetSeal(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	i := vars["id"]
	g := models.Seal{}
	out_arr := []models.Seal{}

	err := s.Dbpool.QueryRow(context.Background(), "SELECT * from func_seal_get($1);", i).Scan(&g.Id, &g.PacketNumber, &g.Area.Id,
		&g.Staff.Id, &g.SealType.Id, &g.SealColour.Id, &g.SealStatus.Id, &g.IssueDate, &g.ReportDate, &g.Area.AreaName, &g.Staff.StaffName,
		&g.SealType.SealTypeName, &g.SealColour.SealColourName, &g.SealStatus.SealStatusName)

	if err != nil && err != pgx.ErrNoRows {
		log.Println("Failed execute from func_seal_get: ", err)
	}

	out_arr = append(out_arr, g)
	auth := models.Auth{Create: true, Read: true, Update: true, Delete: true}

	// output, err := json.Marshal(g)
	out_count, err := json.Marshal(models.Seal_count{Values: out_arr, Count: 1, Auth: auth})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(out_count)

	return

}
