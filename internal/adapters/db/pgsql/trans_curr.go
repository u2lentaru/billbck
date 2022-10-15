package pgsql

import (
	"context"
	"log"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/u2lentaru/billbck/internal/models"
	"github.com/u2lentaru/billbck/pkg/pgclient"
)

//type TransCurrStorage struct
type TransCurrStorage struct {
	db *pgxpool.Pool
}

//func NewTransCurrStorage(db *pgxpool.Pool) *TransCurrStorage
func NewTransCurrStorage(db *pgxpool.Pool) *TransCurrStorage {
	return &TransCurrStorage{db: db}
}

//func (est *TransCurrStorage) GetList(ctx context.Context, pg, pgs int, gs1 string, ord int, dsc bool) (models.TransCurr_count, error)
func (est *TransCurrStorage) GetList(ctx context.Context, pg, pgs int, gs1 string, ord int, dsc bool) (models.TransCurr_count, error) {
	dbpool := pgclient.WDB
	gs := models.TransCurr{}

	gsc := 0
	err := dbpool.QueryRow(ctx, "SELECT * from func_trans_curr_cnt($1);", gs1).Scan(&gsc)

	if err != nil {
		log.Println(err.Error(), "func_trans_curr_cnt")
		return models.TransCurr_count{Values: []models.TransCurr{}, Count: gsc, Auth: models.Auth{}}, err
	}

	out_arr := make([]models.TransCurr, 0,
		func() int {
			if gsc < pgs {
				return gsc
			} else {
				return pgs
			}
		}())

	rows, err := dbpool.Query(ctx, "SELECT * from func_trans_curr_get($1,$2,$3,$4,$5);", pg, pgs, gs1, ord, dsc)
	if err != nil {
		log.Println(err.Error(), "func_trans_curr_get")
		return models.TransCurr_count{Values: []models.TransCurr{}, Count: gsc, Auth: models.Auth{}}, err
	}

	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&gs.Id, &gs.TransCurrName, &gs.TransType.Id, &gs.CheckDate, &gs.NextCheckDate, &gs.ProdDate, &gs.Serial1,
			&gs.Serial2, &gs.Serial3, &gs.TransType.TransTypeName, &gs.TransType.Ratio, &gs.TransType.Class, &gs.TransType.MaxCurr,
			&gs.TransType.NomCurr)
		if err != nil {
			log.Println("failed to scan row:", err)
		}

		out_arr = append(out_arr, gs)
	}

	out_count := models.TransCurr_count{Values: out_arr, Count: gsc, Auth: models.Auth{}}
	if err != nil {
		log.Println(err.Error())
		return models.TransCurr_count{}, err
	}

	return out_count, nil
}

//func (est *TransCurrStorage) Add(ctx context.Context, a models.TransCurr) (int, error)
func (est *TransCurrStorage) Add(ctx context.Context, a models.TransCurr) (int, error) {
	dbpool := pgclient.WDB
	ai := 0

	err := dbpool.QueryRow(ctx, "SELECT func_trans_curr_add($1,$2,$3,$4,$5,$6,$7,$8);", a.TransCurrName,
		a.TransType.Id, a.CheckDate, a.NextCheckDate, a.ProdDate, a.Serial1, a.Serial2, a.Serial3).Scan(&ai)

	if err != nil {
		log.Println("Failed execute func_trans_curr_add: ", err)
		return 0, err
	}

	return ai, nil
}

//func (est *TransCurrStorage) Upd(ctx context.Context, u models.TransCurr) (int, error)
func (est *TransCurrStorage) Upd(ctx context.Context, u models.TransCurr) (int, error) {
	dbpool := pgclient.WDB
	ui := 0

	err := dbpool.QueryRow(ctx, "SELECT func_trans_curr_upd($1,$2,$3,$4,$5,$6,$7,$8,$9);", u.Id, u.TransCurrName,
		u.TransType.Id, u.CheckDate, u.NextCheckDate, u.ProdDate, u.Serial1, u.Serial2, u.Serial3).Scan(&ui)

	if err != nil {
		log.Println("Failed execute func_trans_curr_upd: ", err)
		return 0, err
	}
	return ui, nil
}

//func (est *TransCurrStorage) Del(ctx context.Context, d []int) ([]int, error)
func (est *TransCurrStorage) Del(ctx context.Context, d []int) ([]int, error) {
	dbpool := pgclient.WDB
	res := []int{}
	i := 0
	for _, id := range d {
		err := dbpool.QueryRow(ctx, "SELECT func_trans_curr_del($1);", id).Scan(&i)
		res = append(res, i)

		if err != nil {
			log.Println("Failed execute func_trans_curr_del: ", err)
		}
	}
	return res, nil
}

//func (est *TransCurrStorage) GetOne(ctx context.Context, i int) (models.TransCurr_count, error)
func (est *TransCurrStorage) GetOne(ctx context.Context, i int) (models.TransCurr_count, error) {
	dbpool := pgclient.WDB
	out_arr := []models.TransCurr{}
	g := models.TransCurr{}

	err := dbpool.QueryRow(ctx, "SELECT * from func_trans_curr_getbyid($1);", i).Scan(&g.Id, &g.TransCurrName,
		&g.TransType.Id, &g.CheckDate, &g.NextCheckDate, &g.ProdDate, &g.Serial1, &g.Serial2, &g.Serial3, &g.TransType.TransTypeName,
		&g.TransType.Ratio, &g.TransType.Class, &g.TransType.MaxCurr, &g.TransType.NomCurr)

	if err != nil && err != pgx.ErrNoRows {
		log.Println("Failed execute func_trans_curr_getbyid: ", err)
		return models.TransCurr_count{Values: []models.TransCurr{}, Count: 0, Auth: models.Auth{}}, err
	}

	out_arr = append(out_arr, g)

	out_count := models.TransCurr_count{Values: out_arr, Count: 1, Auth: models.Auth{}}
	return out_count, nil
}
