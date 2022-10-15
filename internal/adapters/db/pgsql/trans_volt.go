package pgsql

import (
	"context"
	"log"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/u2lentaru/billbck/internal/models"
	"github.com/u2lentaru/billbck/pkg/pgclient"
)

//type TransVoltStorage struct
type TransVoltStorage struct {
	db *pgxpool.Pool
}

//func NewTransVoltStorage(db *pgxpool.Pool) *TransVoltStorage
func NewTransVoltStorage(db *pgxpool.Pool) *TransVoltStorage {
	return &TransVoltStorage{db: db}
}

//func (est *TransVoltStorage) GetList(ctx context.Context, pg, pgs int, gs1 string, ord int, dsc bool) (models.TransVolt_count, error)
func (est *TransVoltStorage) GetList(ctx context.Context, pg, pgs int, gs1 string, ord int, dsc bool) (models.TransVolt_count, error) {
	dbpool := pgclient.WDB
	gs := models.TransVolt{}

	gsc := 0
	err := dbpool.QueryRow(ctx, "SELECT * from func_trans_volt_cnt($1);", gs1).Scan(&gsc)

	if err != nil {
		log.Println(err.Error(), "func_trans_volt_cnt")
		return models.TransVolt_count{Values: []models.TransVolt{}, Count: gsc, Auth: models.Auth{}}, err
	}

	out_arr := make([]models.TransVolt, 0,
		func() int {
			if gsc < pgs {
				return gsc
			} else {
				return pgs
			}
		}())

	rows, err := dbpool.Query(ctx, "SELECT * from func_trans_volt_get($1,$2,$3,$4,$5);", pg, pgs, gs1, ord, dsc)
	if err != nil {
		log.Println(err.Error(), "func_trans_volt_get")
		return models.TransVolt_count{Values: []models.TransVolt{}, Count: gsc, Auth: models.Auth{}}, err
	}

	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&gs.Id, &gs.TransVoltName, &gs.TransType.Id, &gs.CheckDate, &gs.NextCheckDate, &gs.ProdDate, &gs.Serial1,
			&gs.Serial2, &gs.Serial3, &gs.TransType.TransTypeName, &gs.TransType.Ratio, &gs.TransType.Class, &gs.TransType.MaxCurr,
			&gs.TransType.NomCurr)
		if err != nil {
			log.Println("failed to scan row:", err)
		}

		out_arr = append(out_arr, gs)
	}

	out_count := models.TransVolt_count{Values: out_arr, Count: gsc, Auth: models.Auth{}}
	if err != nil {
		log.Println(err.Error())
		return models.TransVolt_count{}, err
	}

	return out_count, nil
}

//func (est *TransVoltStorage) Add(ctx context.Context, a models.TransVolt) (int, error)
func (est *TransVoltStorage) Add(ctx context.Context, a models.TransVolt) (int, error) {
	dbpool := pgclient.WDB
	ai := 0

	err := dbpool.QueryRow(ctx, "SELECT func_trans_volt_add($1,$2,$3,$4,$5,$6,$7,$8);", a.TransVoltName,
		a.TransType.Id, a.CheckDate, a.NextCheckDate, a.ProdDate, a.Serial1, a.Serial2, a.Serial3).Scan(&ai)

	if err != nil {
		log.Println("Failed execute func_trans_volt_add: ", err)
		return 0, err
	}

	return ai, nil
}

//func (est *TransVoltStorage) Upd(ctx context.Context, u models.TransVolt) (int, error)
func (est *TransVoltStorage) Upd(ctx context.Context, u models.TransVolt) (int, error) {
	dbpool := pgclient.WDB
	ui := 0

	err := dbpool.QueryRow(ctx, "SELECT func_trans_volt_upd($1,$2,$3,$4,$5,$6,$7,$8,$9);", u.Id, u.TransVoltName,
		u.TransType.Id, u.CheckDate, u.NextCheckDate, u.ProdDate, u.Serial1, u.Serial2, u.Serial3).Scan(&ui)

	if err != nil {
		log.Println("Failed execute func_trans_volt_upd: ", err)
		return 0, err
	}
	return ui, nil
}

//func (est *TransVoltStorage) Del(ctx context.Context, d []int) ([]int, error)
func (est *TransVoltStorage) Del(ctx context.Context, d []int) ([]int, error) {
	dbpool := pgclient.WDB
	res := []int{}
	i := 0
	for _, id := range d {
		err := dbpool.QueryRow(ctx, "SELECT func_trans_volt_del($1);", id).Scan(&i)
		res = append(res, i)

		if err != nil {
			log.Println("Failed execute func_trans_volt_del: ", err)
		}
	}
	return res, nil
}

//func (est *TransVoltStorage) GetOne(ctx context.Context, i int) (models.TransVolt_count, error)
func (est *TransVoltStorage) GetOne(ctx context.Context, i int) (models.TransVolt_count, error) {
	dbpool := pgclient.WDB
	out_arr := []models.TransVolt{}
	g := models.TransVolt{}

	err := dbpool.QueryRow(ctx, "SELECT * from func_trans_volt_getbyid($1);", i).Scan(&g.Id, &g.TransVoltName,
		&g.TransType.Id, &g.CheckDate, &g.NextCheckDate, &g.ProdDate, &g.Serial1, &g.Serial2, &g.Serial3, &g.TransType.TransTypeName,
		&g.TransType.Ratio, &g.TransType.Class, &g.TransType.MaxCurr, &g.TransType.NomCurr)

	if err != nil && err != pgx.ErrNoRows {
		log.Println("Failed execute func_trans_volt_getbyid: ", err)
		return models.TransVolt_count{Values: []models.TransVolt{}, Count: 0, Auth: models.Auth{}}, err
	}

	out_arr = append(out_arr, g)

	out_count := models.TransVolt_count{Values: out_arr, Count: 1, Auth: models.Auth{}}
	return out_count, nil
}
