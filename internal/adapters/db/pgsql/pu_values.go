package pgsql

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/u2lentaru/billbck/internal/models"
	"github.com/u2lentaru/billbck/pkg/pgclient"
	"github.com/xuri/excelize/v2"
)

//type PuValueStorage struct
type PuValueStorage struct {
	db *pgxpool.Pool
}

//func NewPuValueStorage(db *pgxpool.Pool) *PuValueStorage
func NewPuValueStorage(db *pgxpool.Pool) *PuValueStorage {
	return &PuValueStorage{db: db}
}

//func (est *PuValueStorage) GetList(ctx context.Context, pg, pgs, gs1, ord int, dsc bool) (models.PuValue_count, error)
func (est *PuValueStorage) GetList(ctx context.Context, pg, pgs, gs1, ord int, dsc bool) (models.PuValue_count, error) {
	dbpool := pgclient.WDB
	gs := models.PuValue{}

	gsc := 0
	err := dbpool.QueryRow(ctx, "SELECT * from func_pu_values_cnt($1);", gs1).Scan(&gsc)

	if err != nil {
		log.Println(err.Error(), "func_pu_values_cnt")
		return models.PuValue_count{Values: []models.PuValue{}, Count: gsc, Auth: models.Auth{}}, err
	}

	out_arr := make([]models.PuValue, 0,
		func() int {
			if gsc < pgs {
				return gsc
			} else {
				return pgs
			}
		}())

	rows, err := dbpool.Query(ctx, "SELECT * from func_pu_values_get($1,$2,$3,$4,$5);", pg, pgs, gs1, ord, dsc)
	if err != nil {
		log.Println(err.Error(), "func_pu_values_get")
		return models.PuValue_count{Values: []models.PuValue{}, Count: gsc, Auth: models.Auth{}}, err
	}

	defer rows.Close()

	for rows.Next() {
		pv := 0.0

		err = rows.Scan(&gs.Id, &gs.PuId, &gs.ValueDate, &pv)
		if err != nil {
			log.Println("failed to scan row:", err)
		}

		out_arr = append(out_arr, gs)
	}

	out_count := models.PuValue_count{Values: out_arr, Count: gsc, Auth: models.Auth{}}

	return out_count, nil
}

//func (est *PuValueStorage) Add(ctx context.Context, a models.PuValue) (int, error)
func (est *PuValueStorage) Add(ctx context.Context, a models.PuValue) (int, error) {
	dbpool := pgclient.WDB
	ai := 0

	pv, err := strconv.ParseFloat(a.PuValue, 32)
	if err != nil {
		pv = 0
	}

	err = dbpool.QueryRow(ctx, "SELECT func_pu_values_add($1,$2,$3);", a.PuId, a.ValueDate, pv).Scan(&ai)

	if err != nil {
		log.Println("Failed execute func_pu_values_add: ", err)
		return 0, err
	}

	return ai, nil
}

//func (est *PuValueStorage) Upd(ctx context.Context, u models.PuValue) (int, error)
func (est *PuValueStorage) Upd(ctx context.Context, u models.PuValue) (int, error) {
	dbpool := pgclient.WDB
	ui := 0

	pv, err := strconv.ParseFloat(u.PuValue, 32)
	if err != nil {
		pv = 0
	}

	err = dbpool.QueryRow(ctx, "SELECT func_pu_values_upd($1,$2,$3);", u.Id, u.ValueDate, pv).Scan(&ui)

	if err != nil {
		log.Println("Failed execute func_pu_values_upd: ", err)
		return 0, err
	}
	return ui, nil
}

//func (est *PuValueStorage) Del(ctx context.Context, d []int) ([]int, error)
func (est *PuValueStorage) Del(ctx context.Context, d []int) ([]int, error) {
	dbpool := pgclient.WDB
	res := []int{}
	i := 0
	for _, id := range d {
		err := dbpool.QueryRow(ctx, "SELECT func_pu_values_del($1);", id).Scan(&i)
		res = append(res, i)

		if err != nil {
			log.Println("Failed execute func_pu_values_del: ", err)
		}
	}
	return res, nil
}

//func (est *PuValueStorage) GetOne(ctx context.Context, i int) (models.PuValue_count, error)
func (est *PuValueStorage) GetOne(ctx context.Context, i int) (models.PuValue_count, error) {
	dbpool := pgclient.WDB
	out_arr := []models.PuValue{}
	g := models.PuValue{}

	pv := 0.0

	err := dbpool.QueryRow(ctx, "SELECT * from func_pu_value_get($1);", i).Scan(&g.Id, &g.PuId, &g.ValueDate, &pv)

	g.PuValue = strconv.FormatFloat(pv, 'f', 1, 32)

	if err != nil && err != pgx.ErrNoRows {
		log.Println("Failed execute func_pu_value_get: ", err)
		return models.PuValue_count{Values: []models.PuValue{}, Count: 0, Auth: models.Auth{}}, err
	}

	out_arr = append(out_arr, g)

	out_count := models.PuValue_count{Values: out_arr, Count: 1, Auth: models.Auth{}}
	return out_count, nil
}

//func (est *PuValueStorage) AskuePrev(ctx context.Context, af models.AskueFile) (models.PuValueAskue_count, error)
func (est *PuValueStorage) AskuePrev(ctx context.Context, af models.AskueFile) (models.PuValueAskue_count, error) {
	dbpool := pgclient.WDB
	g := models.AskueType{}
	gs := models.PuValueAskue{}
	out_arr := []models.PuValueAskue{}

	tmpfile, err := ioutil.TempFile("/", "askue.*.xslx")
	if err != nil {
		log.Println("Failed create temporary file! :" + err.Error())
	}

	// defer os.Remove(tmpfile.Name())
	defer func() {
		if err := os.Remove(tmpfile.Name()); err != nil {
			log.Println("Failed remove temporary file! :" + err.Error())
		}
	}()

	if _, err := tmpfile.Write(af.AskueFile); err != nil {
		tmpfile.Close()
		log.Println("Failed write to temporary file! :" + err.Error())
	}
	if err := tmpfile.Close(); err != nil {
		log.Println("Failed close temporary file! :" + err.Error())
	}
	f, err := excelize.OpenFile(tmpfile.Name())
	if err != nil {
		log.Println("Failed open temporary Excel file! :" + err.Error())
		return models.PuValueAskue_count{}, err
	}
	defer func() {
		if err := f.Close(); err != nil {
			log.Println(err.Error())
		}
	}()
	rows, err := f.GetRows(af.Sheet)
	if err != nil {
		log.Println("Failed GetRows Excel file! :" + err.Error())
		return models.PuValueAskue_count{}, err
	}
	err = dbpool.QueryRow(ctx, "SELECT * from func_askue_type_get($1);", af.AskueType).Scan(&g.Id, &g.AskueTypeName,
		&g.StartLine, &g.PuColumn, &g.ValueColumn, &g.DateColumn, &g.DateColumnArray)

	if err != nil && err != pgx.ErrNoRows {
		log.Println("Failed execute func_askue_type_get! :" + err.Error())
		return models.PuValueAskue_count{}, err
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
		err = dbpool.QueryRow(ctx, "SELECT func_pu_getbynumber($1);", gs.PuNumber).Scan(&pui)

		if err != nil {
			log.Println("Failed execute func_pu_getbynumber: " + err.Error())
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

	return models.PuValueAskue_count{Values: out_arr, Count: gsc, Auth: models.Auth{}}, nil
}

//func (est *PuValueStorage) AskueLoad(ctx context.Context, af models.AskueFile) (models.AskueLoadRes, error)
func (est *PuValueStorage) AskueLoad(ctx context.Context, af models.AskueFile) (models.AskueLoadRes, error) {
	dbpool := pgclient.WDB
	var procrec, dnrec int
	g := models.AskueType{}
	gs := models.PuValueAskue{}

	tmpfile, err := ioutil.TempFile("/", "askue.*.xslx")
	if err != nil {
		log.Println("Failed create temporary file! :" + err.Error())
	}

	// defer os.Remove(tmpfile.Name())
	defer func() {
		if err := os.Remove(tmpfile.Name()); err != nil {
			log.Println("Failed remove temporary file! :" + err.Error())
		}
	}()

	if _, err := tmpfile.Write(af.AskueFile); err != nil {
		tmpfile.Close()
		log.Println("Failed write to temporary file! :" + err.Error())
	}
	if err := tmpfile.Close(); err != nil {
		log.Println("Failed close temporary file! :" + err.Error())
	}
	f, err := excelize.OpenFile(tmpfile.Name())
	if err != nil {
		log.Println("Failed open temporary Excel file! :" + err.Error())
		return models.AskueLoadRes{Processed: 0, Denied: 0}, err
	}
	defer func() {
		if err := f.Close(); err != nil {
			log.Println(err.Error())
		}
	}()

	rows, err := f.GetRows(af.Sheet)
	if err != nil {
		log.Println("Failed GetRows Excel file! :" + err.Error())
		return models.AskueLoadRes{Processed: 0, Denied: 0}, err
	}
	err = dbpool.QueryRow(ctx, "SELECT * from func_askue_type_get($1);", af.AskueType).Scan(&g.Id, &g.AskueTypeName,
		&g.StartLine, &g.PuColumn, &g.ValueColumn, &g.DateColumn, &g.DateColumnArray)

	if err != nil && err != pgx.ErrNoRows {
		log.Println("Failed execute func_askue_type_get! :" + err.Error())
		return models.AskueLoadRes{Processed: 0, Denied: 0}, err
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
		err = dbpool.QueryRow(ctx, "SELECT func_pu_getbynumber($1);", gs.PuNumber).Scan(&pui)

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
			err = dbpool.QueryRow(ctx, "SELECT func_pu_values_add($1,$2,$3);", gs.PuId, gs.ValueDate, gs.PuValue).Scan(&pui)

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

	out_count := models.AskueLoadRes{Processed: procrec, Denied: dnrec}
	return out_count, nil
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
