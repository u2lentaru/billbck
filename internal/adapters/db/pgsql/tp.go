package pgsql

import (
	"context"
	"log"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/u2lentaru/billbck/internal/models"
	"github.com/u2lentaru/billbck/pkg/pgclient"
)

//type TpStorage struct
type TpStorage struct {
	db *pgxpool.Pool
}

//func NewTpStorage(db *pgxpool.Pool) *TpStorage
func NewTpStorage(db *pgxpool.Pool) *TpStorage {
	return &TpStorage{db: db}
}

//func (est *TpStorage) GetList(ctx context.Context, pg, pgs int, gs1 string, ord int, dsc bool) (models.Tp_count, error)
func (est *TpStorage) GetList(ctx context.Context, pg, pgs int, gs1 string, ord int, dsc bool) (models.Tp_count, error) {
	dbpool := pgclient.WDB
	gs := models.Tp{}

	gsc := 0
	err := dbpool.QueryRow(ctx, "SELECT * from func_tp_cnt($1);", gs1).Scan(&gsc)

	if err != nil {
		log.Println(err.Error(), "func_tp_cnt")
		return models.Tp_count{Values: []models.Tp{}, Count: gsc, Auth: models.Auth{}}, err
	}

	out_arr := make([]models.Tp, 0,
		func() int {
			if gsc < pgs {
				return gsc
			} else {
				return pgs
			}
		}())

	rows, err := dbpool.Query(ctx, "SELECT * from func_tp_get($1,$2,$3,$4,$5);", pg, pgs, gs1, ord, dsc)
	if err != nil {
		log.Println(err.Error(), "func_tp_get")
		return models.Tp_count{Values: []models.Tp{}, Count: gsc, Auth: models.Auth{}}, err
	}

	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&gs.Id, &gs.TpName, &gs.GRp.Id, &gs.GRp.GRpName)
		if err != nil {
			log.Println("failed to scan row:", err)
		}

		out_arr = append(out_arr, gs)
	}

	out_count := models.Tp_count{Values: out_arr, Count: gsc, Auth: models.Auth{}}
	if err != nil {
		log.Println(err.Error())
		return models.Tp_count{}, err
	}

	return out_count, nil
}

//func (est *TpStorage) Add(ctx context.Context, a models.Tp) (int, error)
func (est *TpStorage) Add(ctx context.Context, a models.Tp) (int, error) {
	dbpool := pgclient.WDB
	ai := 0

	err := dbpool.QueryRow(ctx, "SELECT func_tp_add($1,$2);", a.TpName, a.GRp.Id).Scan(&ai)

	if err != nil {
		log.Println("Failed execute func_tp_add: ", err)
		return 0, err
	}

	return ai, nil
}

//func (est *TpStorage) Upd(ctx context.Context, u models.Tp) (int, error)
func (est *TpStorage) Upd(ctx context.Context, u models.Tp) (int, error) {
	dbpool := pgclient.WDB
	ui := 0

	err := dbpool.QueryRow(ctx, "SELECT func_tp_upd($1,$2,$3);", u.Id, u.TpName, u.GRp.Id).Scan(&ui)

	if err != nil {
		log.Println("Failed execute func_tp_upd: ", err)
		return 0, err
	}
	return ui, nil
}

//func (est *TpStorage) Del(ctx context.Context, d []int) ([]int, error)
func (est *TpStorage) Del(ctx context.Context, d []int) ([]int, error) {
	dbpool := pgclient.WDB
	res := []int{}
	i := 0
	for _, id := range d {
		err := dbpool.QueryRow(ctx, "SELECT func_tp_del($1);", id).Scan(&i)
		res = append(res, i)

		if err != nil {
			log.Println("Failed execute func_tp_del: ", err)
		}
	}
	return res, nil
}

//func (est *TpStorage) GetOne(ctx context.Context, i int) (models.Tp_count, error)
func (est *TpStorage) GetOne(ctx context.Context, i int) (models.Tp_count, error) {
	dbpool := pgclient.WDB
	out_arr := []models.Tp{}
	g := models.Tp{}

	err := dbpool.QueryRow(ctx, "SELECT * from func_tp_getbyid($1);", i).Scan(&g.Id, &g.TpName, &g.GRp.Id, &g.GRp.GRpName)

	if err != nil && err != pgx.ErrNoRows {
		log.Println("Failed execute func_tp_getbyid: ", err)
		return models.Tp_count{Values: []models.Tp{}, Count: 0, Auth: models.Auth{}}, err
	}

	out_arr = append(out_arr, g)

	out_count := models.Tp_count{Values: out_arr, Count: 1, Auth: models.Auth{}}
	return out_count, nil
}
