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

// HandlePositions godoc
// @Summary List positions
// @Description get positions
// @Tags positions
// @Produce  json
// @Param page query int false "page number"
// @Param page_size query int false "page size"
// @Param positionname query string false "positionname search pattern"
// @Param ordering query string false "order by {id|positionname}"
// @Param desc query boolean false "descending order {true|false}"
// @Success 200 {object} models.Position_count
// @Failure 500
// @Router /positions [get]
func (s *APG) HandlePositions(w http.ResponseWriter, r *http.Request) {
	p := models.Position{}
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

	pn := ""
	pns, ok := query["positionname"]
	if ok && len(pns) > 0 {
		//case insensitive
		pn = strings.ToUpper(pns[0])
		//quotes
		re := regexp.MustCompile(`'`)
		pn = string(re.ReplaceAll([]byte(pn), []byte("''")))
	}

	pc := 0
	err := s.Dbpool.QueryRow(ctx, "SELECT * from func_positions_cnt($1);", pn).Scan(&pc)

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	out_arr := make([]models.Position, 0,
		func() int {
			if pc < pgs {
				return pc
			} else {
				return pgs
			}
		}())

	ord := 1
	ords, ok := query["ordering"]
	if !ok || len(ords) == 0 {
		ord = 1
	} else if ords[0] == "id" {
		ord = 1
	} else if ords[0] == "positionname" {
		ord = 2
	}

	dsc := false
	dscs, ok := query["desc"]
	if ok && len(dscs) > 0 {
		if !(dscs[0] == "0") {
			dsc = true
		}
	}

	rows, err := s.Dbpool.Query(ctx, "SELECT * from func_positions_get($1,$2,$3,$4,$5);", pg, pgs, pn, ord, dsc)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&p.Id, &p.PositionName)
		if err != nil {
			log.Println("failed to scan row:", err)
		}

		out_arr = append(out_arr, p)
	}

	auth := models.Auth{Create: true, Read: true, Update: true, Delete: true}
	out_count, err := json.Marshal(models.Position_count{Values: out_arr, Count: pc, Auth: auth})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Write(out_count)

	return
}

// HandleAddPosition godoc
// @Summary Add position
// @Description add position
// @Tags positions
// @Accept json
// @Produce  json
// @Param ap body models.AddPosition true "New Position"
// @Success 200 {object} models.Json_id
// @Failure 500
// @Router /positions_add [post]
func (s *APG) HandleAddPosition(w http.ResponseWriter, r *http.Request) {
	ap := models.AddPosition{}
	body, err := ioutil.ReadAll(r.Body)

	defer r.Body.Close()

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	err = json.Unmarshal(body, &ap)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	api := 0
	err = s.Dbpool.QueryRow(context.Background(), "SELECT func_positions_add($1);", ap.PositionName).Scan(&api)

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	output, err := json.Marshal(models.Json_id{Id: api})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(output)

	return

}

// HandleUpdPosition godoc
// @Summary Update position
// @Description update position
// @Tags positions
// @Accept json
// @Produce  json
// @Param up body models.Position true "Update position"
// @Success 200 {object} models.Json_id
// @Failure 500
// @Router /positions_upd [post]
func (s *APG) HandleUpdPosition(w http.ResponseWriter, r *http.Request) {
	up := models.Position{}
	body, err := ioutil.ReadAll(r.Body)

	defer r.Body.Close()

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	err = json.Unmarshal(body, &up)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	upi := 0
	err = s.Dbpool.QueryRow(context.Background(), "SELECT func_positions_upd($1,$2);", up.Id, up.PositionName).Scan(&upi)

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	output, err := json.Marshal(models.Json_id{Id: upi})

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(output)

	return

}

// HandleDelPosition godoc
// @Summary Delete positions
// @Description delete positions
// @Tags positions
// @Accept json
// @Produce  json
// @Param dp body models.Json_ids true "Delete positions"
// @Success 200 {object} models.Json_ids
// @Failure 500
// @Router /positions_del [post]
func (s *APG) HandleDelPosition(w http.ResponseWriter, r *http.Request) {
	dp := models.Json_ids{}
	body, err := ioutil.ReadAll(r.Body)

	defer r.Body.Close()

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	err = json.Unmarshal(body, &dp)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	res := []int{}
	pi := 0
	for _, id := range dp.Ids {
		err = s.Dbpool.QueryRow(context.Background(), "SELECT func_positions_del($1);", id).Scan(&pi)
		res = append(res, pi)

		if err != nil {
			log.Println("Failed execute func_positions_del: ", err)
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

// HandleGetPosition godoc
// @Summary Get position
// @Description get position
// @Tags positions
// @Produce  json
// @Param id path int true "Position by id"
// @Success 200 {object} models.Position_count
// @Failure 500
// @Router /positions/{id} [get]
func (s *APG) HandleGetPosition(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	pi := vars["id"]
	p := models.Position{}
	out_arr := []models.Position{}

	err := s.Dbpool.QueryRow(context.Background(), "SELECT * from func_position_get($1);", pi).Scan(&p.Id, &p.PositionName)

	if err != nil && err != pgx.ErrNoRows {
		log.Println("Failed execute from func_position_get: ", err)
		// http.Error(w, err.Error(), 500)
		// return
	}

	out_arr = append(out_arr, p)
	auth := models.Auth{Create: true, Read: true, Update: true, Delete: true}

	// output, err := json.Marshal(p)
	out_count, err := json.Marshal(models.Position_count{Values: out_arr, Count: 1, Auth: auth})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(out_count)

	return

}
