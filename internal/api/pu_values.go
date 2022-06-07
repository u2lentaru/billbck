package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v4"
	"github.com/u2lentaru/billbck/internal/models"
	"github.com/xuri/excelize/v2"
)

// HandlePuValues godoc
// @Summary List pu values
// @Description get pu values list
// @Tags puvalues
// @Produce  json
// @Param page query int false "page number"
// @Param page_size query int false "page size"
// @Param puid query int false "puid search pattern"
// @Param ordering query string false "order by {puid|valuedate}"
// @Param desc query boolean false "descending order {true|false}"
// @Success 200 {object} models.PuValue_count
// @Failure 500
// @Router /puvalues [get]
func (s *APG) HandlePuValues(w http.ResponseWriter, r *http.Request) {
	gs := models.PuValue{}
	ctx := context.Background()
	out_arr := []models.PuValue{}

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

	gs1 := 0
	gs1s, ok := query["puid"]
	if ok && len(gs1s) > 0 {
		t, err := strconv.Atoi(gs1s[0])
		if err == nil {
			gs1 = t
		}
	}

	ord := 1
	ords, ok := query["ordering"]
	if !ok || len(ords) == 0 {
		ord = 1
	} else if ords[0] == "puid" {
		ord = 2
	} else if ords[0] == "valuedate" {
		ord = 3
	}

	dsc := false
	dscs, ok := query["desc"]
	if ok && len(dscs) > 0 {
		if !(dscs[0] == "0") {
			dsc = true
		}
	}

	rows, err := s.Dbpool.Query(ctx, "SELECT * from func_pu_values_get($1,$2,$3,$4,$5);", pg, pgs, gs1, ord, dsc)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	defer rows.Close()

	for rows.Next() {
		pv := 0.0

		err = rows.Scan(&gs.Id, &gs.PuId, &gs.ValueDate, &pv)

		if err != nil {
			log.Println("failed to scan row:", err)
		}

		gs.PuValue = strconv.FormatFloat(pv, 'f', 1, 32)

		out_arr = append(out_arr, gs)
	}

	gsc := 0
	err = s.Dbpool.QueryRow(ctx, "SELECT * from func_pu_values_cnt($1);", gs1).Scan(&gsc)

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	auth := models.Auth{Create: true, Read: true, Update: true, Delete: true}

	out_count, err := json.Marshal(models.PuValue_count{Values: out_arr, Count: gsc, Auth: auth})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Write(out_count)

	return

}

// HandleAddPuValue godoc
// @Summary Add pu value
// @Description add pu value
// @Tags puvalues
// @Accept json
// @Produce  json
// @Param a body models.PuValue true "New pu value. Significant params: PuId, ValueDate, PuValue"
// @Success 200 {object} models.Json_id
// @Failure 500
// @Router /puvalues_add [post]
func (s *APG) HandleAddPuValue(w http.ResponseWriter, r *http.Request) {
	a := models.PuValue{}
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

	pv, err := strconv.ParseFloat(a.PuValue, 32)
	if err != nil {
		pv = 0
	}

	err = s.Dbpool.QueryRow(context.Background(), "SELECT func_pu_values_add($1,$2,$3);", a.PuId, a.ValueDate, pv).Scan(&ai)

	if err != nil {
		log.Println("Failed execute func_pu_values_add: ", err)
	}

	output, err := json.Marshal(models.Json_id{Id: ai})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(output)

	return

}

// HandleUpdPuValue godoc
// @Summary Update pu value
// @Description update pu value
// @Tags puvalues
// @Accept json
// @Produce  json
// @Param u body models.PuValue true "Update pu value. Significant params: Id, ValueDate, PuValue"
// @Success 200 {object} models.Json_id
// @Failure 500
// @Router /puvalues_upd [post]
func (s *APG) HandleUpdPuValue(w http.ResponseWriter, r *http.Request) {
	u := models.PuValue{}
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

	pv, err := strconv.ParseFloat(u.PuValue, 32)

	if err != nil {
		pv = 0
	}

	err = s.Dbpool.QueryRow(context.Background(), "SELECT func_pu_values_upd($1,$2,$3);", u.Id, u.ValueDate, pv).Scan(&ui)

	if err != nil {
		log.Println("Failed execute func_pu_values_upd: ", err)
	}

	output, err := json.Marshal(models.Json_id{Id: ui})

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(output)

	return
}

// HandleDelPuValue godoc
// @Summary Delete pu values
// @Description delete pu values
// @Tags puvalues
// @Accept json
// @Produce  json
// @Param d body models.Json_ids true "Delete pu values"
// @Success 200 {object} models.Json_ids
// @Failure 500
// @Router /puvalues_del [post]
func (s *APG) HandleDelPuValue(w http.ResponseWriter, r *http.Request) {
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
		err = s.Dbpool.QueryRow(context.Background(), "SELECT func_pu_values_del($1);", id).Scan(&i)
		res = append(res, i)

		if err != nil {
			log.Println("Failed execute func_pu_values_del: ", err)
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

// HandleGetPuValue godoc
// @Summary Get pu value
// @Description get pu value
// @Tags puvalues
// @Produce  json
// @Param id path int true "Pu value Id"
// @Success 200 {object} models.PuValue_count
// @Failure 500
// @Router /puvalues/{id} [get]
func (s *APG) HandleGetPuValue(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	i := vars["id"]
	g := models.PuValue{}
	out_arr := []models.PuValue{}

	pv := 0.0

	err := s.Dbpool.QueryRow(context.Background(), "SELECT * from func_pu_value_get($1);", i).Scan(&g.Id, &g.PuId, &g.ValueDate, &pv)

	g.PuValue = strconv.FormatFloat(pv, 'f', 1, 32)

	if err != nil && err != pgx.ErrNoRows {
		log.Println("Failed execute from func_pu_value_get: ", err)
	}

	out_arr = append(out_arr, g)
	auth := models.Auth{Create: true, Read: true, Update: true, Delete: true}

	// output, err := json.Marshal(g)
	out_count, err := json.Marshal(models.PuValue_count{Values: out_arr, Count: 1, Auth: auth})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(out_count)

	return

}

// HandlePuValuesAskuePreview godoc
// @Summary Preview askue pu values
// @Description preview askue pu values
// @Tags puvalues
// @Accept json
// @Produce  json
// @Param af body models.AskueFile true "Askue file to preview"
// @Success 200 {object} models.PuValueAskue_count
// @Failure 500
// @Router /puvalues_askue_prev [post]
func (s *APG) HandlePuValuesAskuePreview(w http.ResponseWriter, r *http.Request) {
	g := models.AskueType{}
	gs := models.PuValueAskue{}
	out_arr := []models.PuValueAskue{}

	af := models.AskueFile{}
	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	err = json.Unmarshal(body, &af)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	tmpfile, err := ioutil.TempFile("/", "askue.*.xslx")
	if err != nil {
		http.Error(w, "Failed create temporary file! :"+err.Error(), 500)
	}

	// defer os.Remove(tmpfile.Name())
	defer func() {
		if err := os.Remove(tmpfile.Name()); err != nil {
			http.Error(w, "Failed remove temporary file! :"+err.Error(), 500)
		}
	}()

	if _, err := tmpfile.Write(af.AskueFile); err != nil {
		tmpfile.Close()
		http.Error(w, "Failed write to temporary file! :"+err.Error(), 500)
	}
	if err := tmpfile.Close(); err != nil {
		http.Error(w, "Failed close temporary file! :"+err.Error(), 500)
	}
	f, err := excelize.OpenFile(tmpfile.Name())
	if err != nil {
		http.Error(w, "Failed open temporary Excel file! :"+err.Error(), 500)
		return
	}
	defer func() {
		if err := f.Close(); err != nil {
			http.Error(w, err.Error(), 500)
		}
	}()
	rows, err := f.GetRows(af.Sheet)
	if err != nil {
		http.Error(w, "Failed GetRows Excel file! :"+err.Error(), 500)
		return
	}
	err = s.Dbpool.QueryRow(context.Background(), "SELECT * from func_askue_type_get($1);", af.AskueType).Scan(&g.Id, &g.AskueTypeName,
		&g.StartLine, &g.PuColumn, &g.ValueColumn, &g.DateColumn, &g.DateColumnArray)

	if err != nil && err != pgx.ErrNoRows {
		http.Error(w, "Failed execute func_askue_type_get! :"+err.Error(), 500)
		return
	}
	vds := []string{}
	if g.DateColumnArray != nil {
		vds = strings.Split(*g.DateColumnArray, ",")
	}

	gsc := 0
	pui := 0
	pv := 0.0
	for i, row := range rows {
		if i < g.StartLine-1 {
			continue
		}
		gs.Valid = true
		gs.PuNumber = row[g.PuColumn-1]

		pv, err = strconv.ParseFloat(row[g.ValueColumn-1], 32)
		if err != nil {
			gs.Valid = false
			gs.Notes = "Incorrect PuValue: " + row[g.ValueColumn-1] + " "
		} else {
			gs.PuValue = float32(pv)
		}
		// gs.PuValue = row[g.ValueColumn-1]

		if g.DateColumnArray != nil {
			dcn := 0
			fdc := [3]bool{false, false, false}
			for i, vd := range vds {
				if i == 0 {
					dcn, err = strconv.Atoi(vd)
					if err != nil {
						// http.Error(w, "Incorrect DateColumnArray: "+err.Error(), 500)
						// return
						gs.ValueDate = "0000-00-00"
						fdc[i] = true
					} else {
						_, err = time.Parse("2006-01-02", row[dcn])

						if err != nil {
							fdc[i] = true
						}
						gs.ValueDate = row[dcn]
					}
				} else {
					dcn, err = strconv.Atoi(vd)
					if err != nil {
						fdc[i] = true
						continue
					}

					_, err = time.Parse("2006-01-02", row[dcn])

					if err != nil {
						fdc[i] = true
						continue
					}

					if row[dcn] > gs.ValueDate {
						gs.ValueDate = vd
					}
				}
			}

			if fdc[0] && fdc[1] && fdc[2] {
				gs.Valid = false
				gs.Notes = gs.Notes + "Incorrect ValueDate "
			}

		} else {
			_, err = time.Parse("2006-01-02", row[g.DateColumn-1])
			//  tdv, err := time.Parse("02.01.2006", row[g.DateColumn-1])

			if err != nil {
				gs.Valid = false
				gs.Notes = gs.Notes + "Incorrect ValueDate: " + row[g.DateColumn-1] + " "
			}

			gs.ValueDate = row[g.DateColumn-1]
			// gs.ValueDate = fmt.Sprintf(tdv.Format("2006-01-02"))
		}
		err = s.Dbpool.QueryRow(context.Background(), "SELECT func_pu_getbynumber($1);", gs.PuNumber).Scan(&pui)

		if err != nil {
			http.Error(w, "Failed execute func_pu_getbynumber: "+err.Error(), 500)
			// return
		}
		if pui == 0 {
			// http.Error(w, "Pu number does not exist!", 500)
			// return
			gs.Valid = false
			gs.Notes = gs.Notes + "Punumber does not exist"
		}

		gs.PuId = pui

		out_arr = append(out_arr, gs)
		gsc++
	}
	auth := models.Auth{Create: true, Read: true, Update: true, Delete: true}
	out_count, err := json.Marshal(models.PuValueAskue_count{Values: out_arr, Count: gsc, Auth: auth})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(out_count)

	return

}

// HandlePuValuesAskue godoc
// @Summary Load askue pu values
// @Description load askue pu values
// @Tags puvalues
// @Accept json
// @Produce  json
// @Param af body models.AskueFile true "Askue file to load"
// @Success 200 {object} models.AskueLoadRes
// @Failure 500
// @Router /puvalues_askue [post]
func (s *APG) HandlePuValuesAskue(w http.ResponseWriter, r *http.Request) {
	var procrec, dnrec int
	g := models.AskueType{}
	gs := models.PuValueAskue{}

	af := models.AskueFile{}
	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	err = json.Unmarshal(body, &af)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	tmpfile, err := ioutil.TempFile("/", "askue.*.xslx")
	if err != nil {
		http.Error(w, "Failed create temporary file! :"+err.Error(), 500)
	}

	// defer os.Remove(tmpfile.Name())
	defer func() {
		if err := os.Remove(tmpfile.Name()); err != nil {
			http.Error(w, "Failed remove temporary file! :"+err.Error(), 500)
		}
	}()

	if _, err := tmpfile.Write(af.AskueFile); err != nil {
		tmpfile.Close()
		http.Error(w, "Failed write to temporary file! :"+err.Error(), 500)
	}
	if err := tmpfile.Close(); err != nil {
		http.Error(w, "Failed close temporary file! :"+err.Error(), 500)
	}
	f, err := excelize.OpenFile(tmpfile.Name())
	if err != nil {
		http.Error(w, "Failed open temporary Excel file! :"+err.Error(), 500)
		return
	}
	defer func() {
		if err := f.Close(); err != nil {
			http.Error(w, err.Error(), 500)
		}
	}()

	rows, err := f.GetRows(af.Sheet)
	if err != nil {
		http.Error(w, "Failed GetRows Excel file! :"+err.Error(), 500)
		return
	}
	err = s.Dbpool.QueryRow(context.Background(), "SELECT * from func_askue_type_get($1);", af.AskueType).Scan(&g.Id, &g.AskueTypeName,
		&g.StartLine, &g.PuColumn, &g.ValueColumn, &g.DateColumn, &g.DateColumnArray)

	if err != nil && err != pgx.ErrNoRows {
		http.Error(w, "Failed execute func_askue_type_get! :"+err.Error(), 500)
		return
	}
	vds := []string{}
	if g.DateColumnArray != nil {
		vds = strings.Split(*g.DateColumnArray, ",")
	}

	pui := 0
	pv := 0.0
	for i, row := range rows {
		if i < g.StartLine-1 {
			continue
		}
		gs.Valid = true
		gs.PuNumber = row[g.PuColumn-1]

		pv, err = strconv.ParseFloat(row[g.ValueColumn-1], 32)
		if err != nil {
			gs.Valid = false
			gs.Notes = "Incorrect PuValue: " + row[g.ValueColumn-1] + " "
		} else {
			gs.PuValue = float32(pv)
		}
		// gs.PuValue = row[g.ValueColumn-1]

		if g.DateColumnArray != nil {
			dcn := 0
			fdc := [3]bool{false, false, false}
			for i, vd := range vds {
				if i == 0 {
					dcn, err = strconv.Atoi(vd)
					if err != nil {
						// http.Error(w, "Incorrect DateColumnArray: "+err.Error(), 500)
						// return
						gs.ValueDate = "0000-00-00"
						fdc[i] = true
					} else {
						_, err = time.Parse("2006-01-02", row[dcn])

						if err != nil {
							fdc[i] = true
						}
						gs.ValueDate = row[dcn]
					}
				} else {
					dcn, err = strconv.Atoi(vd)
					if err != nil {
						fdc[i] = true
						continue
					}

					_, err = time.Parse("2006-01-02", row[dcn])

					if err != nil {
						fdc[i] = true
						continue
					}

					if row[dcn] > gs.ValueDate {
						gs.ValueDate = vd
					}
				}
			}

			if fdc[0] && fdc[1] && fdc[2] {
				gs.Valid = false
				gs.Notes = gs.Notes + "Incorrect ValueDate "
			}

		} else {
			_, err = time.Parse("2006-01-02", row[g.DateColumn-1])
			//  tdv, err := time.Parse("02.01.2006", row[g.DateColumn-1])

			if err != nil {
				gs.Valid = false
				gs.Notes = gs.Notes + "Incorrect ValueDate: " + row[g.DateColumn-1] + " "
			}

			gs.ValueDate = row[g.DateColumn-1]
			// gs.ValueDate = fmt.Sprintf(tdv.Format("2006-01-02"))
		}
		err = s.Dbpool.QueryRow(context.Background(), "SELECT func_pu_getbynumber($1);", gs.PuNumber).Scan(&pui)

		if err != nil {
			// http.Error(w, "Failed execute func_pu_getbynumber: "+err.Error(), 500)
			log.Println("Failed execute func_pu_getbynumber: " + err.Error())
		}
		if pui == 0 {
			gs.Valid = false
			gs.Notes = gs.Notes + "Punumber does not exist"
		}

		gs.PuId = pui

		if gs.Valid {
			err = s.Dbpool.QueryRow(context.Background(), "SELECT func_pu_values_add($1,$2,$3);", gs.PuId, gs.ValueDate, gs.PuValue).Scan(&pui)

			if err != nil {
				log.Println("Failed execute func_pu_values_add: " + err.Error())
				dnrec++
			} else {
				if pui == 0 {
					log.Println("Can't add record into wt_pu_values!")
					dnrec++
				} else {
					procrec++
				}
			}

		} else {
			dnrec++
		}

	}

	out_count, err := json.Marshal(models.AskueLoadRes{Processed: procrec, Denied: dnrec})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(out_count)

	return
}

// func DVTrans(sd string) (string, bool) - "DD.MM "* -> "DD.MM.YYYY", true ; !"DD.MM "* - "", false
func DVTrans(sd string) (string, bool) {
	var fl bool
	if sd[5:6] == " " {
		log.Println("sd is ok.")
		fl = true
	} else {
		log.Println("Incorrect dv!")
		fl = false
	}

	dv := "16.01 07:41 (4)"
	dv = dv[:5]
	tdv, err := time.Parse("02.01", dv)

	if err != nil {
		log.Println("Incorrect dv!")
		return "", false
	}

	stdv := fmt.Sprintf(tdv.Format("02.01.2006"))
	log.Println(stdv[:5])
	y := fmt.Sprintf(time.Now().Format("02.01.2006"))
	log.Println("y: ", y)
	y = y[6:]
	stdv = stdv[:5]
	log.Println("y:=", y, " stdv:= ", stdv, " + ", stdv+"."+y)
	return stdv + "." + y, fl

}
