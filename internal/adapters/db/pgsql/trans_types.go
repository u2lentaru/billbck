package pgsql

import (
	"context"
	"log"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/u2lentaru/billbck/internal/models"
	"github.com/u2lentaru/billbck/pkg/pgclient"
)

//type TransTypeStorage struct
type TransTypeStorage struct {
	db *pgxpool.Pool
}

//func NewTransTypeStorage(db *pgxpool.Pool) *TransTypeStorage
func NewTransTypeStorage(db *pgxpool.Pool) *TransTypeStorage {
	return &TransTypeStorage{db: db}
}

//func (est *TransTypeStorage) GetList(ctx context.Context, pg, pgs int, gs1 string, ord int, dsc bool) (models.TransType_count, error)
func (est *TransTypeStorage) GetList(ctx context.Context, pg, pgs int, gs1 string, ord int, dsc bool) (models.TransType_count, error) {
	dbpool := pgclient.WDB
	gs := models.TransType{}

	gsc := 0
	err := dbpool.QueryRow(ctx, "SELECT * from func_trans_types_cnt($1);", gs1).Scan(&gsc)

	if err != nil {
		log.Println(err.Error(), "func_trans_types_cnt")
		return models.TransType_count{Values: []models.TransType{}, Count: gsc, Auth: models.Auth{}}, err
	}

	out_arr := make([]models.TransType, 0,
		func() int {
			if gsc < pgs {
				return gsc
			} else {
				return pgs
			}
		}())

	rows, err := dbpool.Query(ctx, "SELECT * from func_trans_types_get($1,$2,$3,$4,$5);", pg, pgs, gs1, ord, dsc)
	if err != nil {
		log.Println(err.Error(), "func_trans_types_get")
		return models.TransType_count{Values: []models.TransType{}, Count: gsc, Auth: models.Auth{}}, err
	}

	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&gs.Id, &gs.TransTypeName, &gs.Ratio, &gs.Class, &gs.MaxCurr, &gs.NomCurr)
		if err != nil {
			log.Println("failed to scan row:", err)
		}

		out_arr = append(out_arr, gs)
	}

	out_count := models.TransType_count{Values: out_arr, Count: gsc, Auth: models.Auth{}}
	return out_count, nil
}

//func (est *TransTypeStorage) Add(ctx context.Context, a models.TransType) (int, error)
func (est *TransTypeStorage) Add(ctx context.Context, a models.TransType) (int, error) {
	dbpool := pgclient.WDB
	ai := 0

	err := dbpool.QueryRow(ctx, "SELECT func_trans_types_add($1,$2,$3,$4,$5);", a.TransTypeName, a.Ratio, a.Class, a.MaxCurr, a.NomCurr).Scan(&ai)

	if err != nil {
		log.Println("Failed execute func_trans_types_add: ", err)
		return 0, err
	}

	return ai, nil
}

//func (est *TransTypeStorage) Upd(ctx context.Context, u models.TransType) (int, error)
func (est *TransTypeStorage) Upd(ctx context.Context, u models.TransType) (int, error) {
	dbpool := pgclient.WDB
	ui := 0

	err := dbpool.QueryRow(ctx, "SELECT func_trans_types_upd($1,$2,$3,$4,$5,$6);", u.Id, u.TransTypeName, u.Ratio, u.Class, u.MaxCurr,
		u.NomCurr).Scan(&ui)

	if err != nil {
		log.Println("Failed execute func_trans_types_upd: ", err)
		return 0, err
	}
	return ui, nil
}

//func (est *TransTypeStorage) Del(ctx context.Context, d []int) ([]int, error)
func (est *TransTypeStorage) Del(ctx context.Context, d []int) ([]int, error) {
	dbpool := pgclient.WDB
	res := []int{}
	i := 0
	for _, id := range d {
		err := dbpool.QueryRow(ctx, "SELECT func_trans_types_del($1);", id).Scan(&i)
		res = append(res, i)

		if err != nil {
			log.Println("Failed execute func_trans_types_del: ", err)
		}
	}
	return res, nil
}

//func (est *TransTypeStorage) GetOne(ctx context.Context, i int) (models.TransType_count, error)
func (est *TransTypeStorage) GetOne(ctx context.Context, i int) (models.TransType_count, error) {
	dbpool := pgclient.WDB
	out_arr := []models.TransType{}
	g := models.TransType{}

	err := dbpool.QueryRow(ctx, "SELECT * from func_trans_type_get($1);", i).Scan(&g.Id, &g.TransTypeName, &g.Ratio, &g.Class, &g.MaxCurr,
		&g.NomCurr)

	if err != nil && err != pgx.ErrNoRows {
		log.Println("Failed execute func_trans_type_get: ", err)
		return models.TransType_count{Values: []models.TransType{}, Count: 0, Auth: models.Auth{}}, err
	}

	out_arr = append(out_arr, g)

	out_count := models.TransType_count{Values: out_arr, Count: 1, Auth: models.Auth{}}
	return out_count, nil
}
