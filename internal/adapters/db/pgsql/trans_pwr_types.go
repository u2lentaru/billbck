package pgsql

import (
	"context"
	"log"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/u2lentaru/billbck/internal/models"
	"github.com/u2lentaru/billbck/pkg/pgclient"
)

//type TransPwrTypeStorage struct
type TransPwrTypeStorage struct {
	db *pgxpool.Pool
}

//func NewTransPwrTypeStorage(db *pgxpool.Pool) *TransPwrTypeStorage
func NewTransPwrTypeStorage(db *pgxpool.Pool) *TransPwrTypeStorage {
	return &TransPwrTypeStorage{db: db}
}

//func (est *TransPwrTypeStorage) GetList(ctx context.Context, pg, pgs int, gs1 string, ord int, dsc bool) (models.TransPwrType_count, error)
func (est *TransPwrTypeStorage) GetList(ctx context.Context, pg, pgs int, gs1 string, ord int, dsc bool) (models.TransPwrType_count, error) {
	dbpool := pgclient.WDB
	gs := models.TransPwrType{}

	gsc := 0
	err := dbpool.QueryRow(ctx, "SELECT * from func_trans_pwr_types_cnt($1);", gs1).Scan(&gsc)

	if err != nil {
		log.Println(err.Error(), "func_trans_pwr_types_cnt")
		return models.TransPwrType_count{Values: []models.TransPwrType{}, Count: gsc, Auth: models.Auth{}}, err
	}

	out_arr := make([]models.TransPwrType, 0,
		func() int {
			if gsc < pgs {
				return gsc
			} else {
				return pgs
			}
		}())

	rows, err := dbpool.Query(ctx, "SELECT * from func_trans_pwr_types_get($1,$2,$3,$4,$5);", pg, pgs, gs1, ord, dsc)
	if err != nil {
		log.Println(err.Error(), "func_trans_pwr_types_get")
		return models.TransPwrType_count{Values: []models.TransPwrType{}, Count: gsc, Auth: models.Auth{}}, err
	}

	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&gs.Id, &gs.TransPwrTypeName, &gs.ShortCircuitPower, &gs.IdlingLossPower, &gs.NominalPower)
		if err != nil {
			log.Println("failed to scan row:", err)
		}

		out_arr = append(out_arr, gs)
	}

	out_count := models.TransPwrType_count{Values: out_arr, Count: gsc, Auth: models.Auth{}}
	if err != nil {
		log.Println(err.Error())
		return models.TransPwrType_count{}, err
	}

	return out_count, nil
}

//func (est *TransPwrTypeStorage) Add(ctx context.Context, a models.TransPwrType) (int, error)
func (est *TransPwrTypeStorage) Add(ctx context.Context, a models.TransPwrType) (int, error) {
	dbpool := pgclient.WDB
	ai := 0

	err := dbpool.QueryRow(ctx, "SELECT func_trans_pwr_types_add($1,$2,$3,$4);", a.TransPwrTypeName,
		a.ShortCircuitPower, a.IdlingLossPower, a.NominalPower).Scan(&ai)

	if err != nil {
		log.Println("Failed execute func_trans_pwr_types_add: ", err)
		return 0, err
	}

	return ai, nil
}

//func (est *TransPwrTypeStorage) Upd(ctx context.Context, u models.TransPwrType) (int, error)
func (est *TransPwrTypeStorage) Upd(ctx context.Context, u models.TransPwrType) (int, error) {
	dbpool := pgclient.WDB
	ui := 0

	err := dbpool.QueryRow(ctx, "SELECT func_trans_pwr_types_upd($1,$2,$3,$4,$5);", u.Id, u.TransPwrTypeName,
		u.ShortCircuitPower, u.IdlingLossPower, u.NominalPower).Scan(&ui)

	if err != nil {
		log.Println("Failed execute func_trans_pwr_types_upd: ", err)
		return 0, err
	}
	return ui, nil
}

//func (est *TransPwrTypeStorage) Del(ctx context.Context, d []int) ([]int, error)
func (est *TransPwrTypeStorage) Del(ctx context.Context, d []int) ([]int, error) {
	dbpool := pgclient.WDB
	res := []int{}
	i := 0
	for _, id := range d {
		err := dbpool.QueryRow(ctx, "SELECT func_trans_pwr_types_del($1);", id).Scan(&i)
		res = append(res, i)

		if err != nil {
			log.Println("Failed execute func_trans_pwr_types_del: ", err)
		}
	}
	return res, nil
}

//func (est *TransPwrTypeStorage) GetOne(ctx context.Context, i int) (models.TransPwrType_count, error)
func (est *TransPwrTypeStorage) GetOne(ctx context.Context, i int) (models.TransPwrType_count, error) {
	dbpool := pgclient.WDB
	out_arr := []models.TransPwrType{}
	g := models.TransPwrType{}

	err := dbpool.QueryRow(ctx, "SELECT * from func_trans_pwr_type_get($1);", i).Scan(&g.Id, &g.TransPwrTypeName,
		&g.ShortCircuitPower, &g.IdlingLossPower, &g.NominalPower)

	if err != nil && err != pgx.ErrNoRows {
		log.Println("Failed execute func_trans_pwr_type_get: ", err)
		return models.TransPwrType_count{Values: []models.TransPwrType{}, Count: 0, Auth: models.Auth{}}, err
	}

	out_arr = append(out_arr, g)

	out_count := models.TransPwrType_count{Values: out_arr, Count: 1, Auth: models.Auth{}}
	return out_count, nil
}
