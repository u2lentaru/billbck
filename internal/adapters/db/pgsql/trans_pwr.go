package pgsql

import (
	"context"
	"log"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/u2lentaru/billbck/internal/models"
	"github.com/u2lentaru/billbck/pkg/pgclient"
)

//type TransPwrStorage struct
type TransPwrStorage struct {
	db *pgxpool.Pool
}

//func NewTransPwrStorage(db *pgxpool.Pool) *TransPwrStorage
func NewTransPwrStorage(db *pgxpool.Pool) *TransPwrStorage {
	return &TransPwrStorage{db: db}
}

//func (est *TransPwrStorage) GetList(ctx context.Context, pg, pgs int, gs1 string, ord int, dsc bool) (models.TransPwr_count, error)
func (est *TransPwrStorage) GetList(ctx context.Context, pg, pgs int, gs1 string, ord int, dsc bool) (models.TransPwr_count, error) {
	dbpool := pgclient.WDB
	gs := models.TransPwr{}

	gsc := 0
	err := dbpool.QueryRow(ctx, "SELECT * from func_trans_pwr_cnt($1);", gs1).Scan(&gsc)

	if err != nil {
		log.Println(err.Error(), "func_trans_pwr_cnt")
		return models.TransPwr_count{Values: []models.TransPwr{}, Count: gsc, Auth: models.Auth{}}, err
	}

	out_arr := make([]models.TransPwr, 0,
		func() int {
			if gsc < pgs {
				return gsc
			} else {
				return pgs
			}
		}())

	rows, err := dbpool.Query(ctx, "SELECT * from func_trans_pwr_get($1,$2,$3,$4,$5);", pg, pgs, gs1, ord, dsc)
	if err != nil {
		log.Println(err.Error(), "func_trans_pwr_get")
		return models.TransPwr_count{Values: []models.TransPwr{}, Count: gsc, Auth: models.Auth{}}, err
	}

	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&gs.Id, &gs.TransPwrName, &gs.TransPwrType.Id, &gs.TransPwrType.TransPwrTypeName, &gs.TransPwrType.ShortCircuitPower,
			&gs.TransPwrType.IdlingLossPower, &gs.TransPwrType.NominalPower)
		if err != nil {
			log.Println("failed to scan row:", err)
		}

		out_arr = append(out_arr, gs)
	}

	out_count := models.TransPwr_count{Values: out_arr, Count: gsc, Auth: models.Auth{}}
	return out_count, nil
}

//func (est *TransPwrStorage) Add(ctx context.Context, a models.TransPwr) (int, error)
func (est *TransPwrStorage) Add(ctx context.Context, a models.TransPwr) (int, error) {
	dbpool := pgclient.WDB
	ai := 0

	err := dbpool.QueryRow(ctx, "SELECT func_trans_pwr_add($1,$2);", a.TransPwrName, a.TransPwrType.Id).Scan(&ai)

	if err != nil {
		log.Println("Failed execute func_trans_pwr_add: ", err)
		return 0, err
	}

	return ai, nil
}

//func (est *TransPwrStorage) Upd(ctx context.Context, u models.TransPwr) (int, error)
func (est *TransPwrStorage) Upd(ctx context.Context, u models.TransPwr) (int, error) {
	dbpool := pgclient.WDB
	ui := 0

	err := dbpool.QueryRow(ctx, "SELECT func_trans_pwr_upd($1,$2,$3);", u.Id, u.TransPwrName, u.TransPwrType.Id).Scan(&ui)

	if err != nil {
		log.Println("Failed execute func_trans_pwr_upd: ", err)
		return 0, err
	}
	return ui, nil
}

//func (est *TransPwrStorage) Del(ctx context.Context, d []int) ([]int, error)
func (est *TransPwrStorage) Del(ctx context.Context, d []int) ([]int, error) {
	dbpool := pgclient.WDB
	res := []int{}
	i := 0
	for _, id := range d {
		err := dbpool.QueryRow(ctx, "SELECT func_trans_pwr_del($1);", id).Scan(&i)
		res = append(res, i)

		if err != nil {
			log.Println("Failed execute func_trans_pwr_del: ", err)
		}
	}
	return res, nil
}

//func (est *TransPwrStorage) GetOne(ctx context.Context, i int) (models.TransPwr_count, error)
func (est *TransPwrStorage) GetOne(ctx context.Context, i int) (models.TransPwr_count, error) {
	dbpool := pgclient.WDB
	out_arr := []models.TransPwr{}
	g := models.TransPwr{}

	err := dbpool.QueryRow(ctx, "SELECT * from func_trans_pwr_getbyid($1);", i).Scan(&g.Id, &g.TransPwrName, &g.TransPwrType.Id,
		&g.TransPwrType.TransPwrTypeName, &g.TransPwrType.ShortCircuitPower, &g.TransPwrType.IdlingLossPower, &g.TransPwrType.NominalPower)

	if err != nil && err != pgx.ErrNoRows {
		log.Println("Failed execute func_trans_pwr_getbyid: ", err)
		return models.TransPwr_count{Values: []models.TransPwr{}, Count: 0, Auth: models.Auth{}}, err
	}

	out_arr = append(out_arr, g)

	out_count := models.TransPwr_count{Values: out_arr, Count: 1, Auth: models.Auth{}}
	return out_count, nil
}
