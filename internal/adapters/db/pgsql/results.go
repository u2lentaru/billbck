package pgsql

import (
	"context"
	"log"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/u2lentaru/billbck/internal/models"
	"github.com/u2lentaru/billbck/pkg/pgclient"
)

//type ResultStorage struct
type ResultStorage struct {
	db *pgxpool.Pool
}

//func NewResultStorage(db *pgxpool.Pool) *ResultStorage
func NewResultStorage(db *pgxpool.Pool) *ResultStorage {
	return &ResultStorage{db: db}
}

//func (est *ResultStorage) GetList(ctx context.Context, pg, pgs int, gs1 string, ord int, dsc bool) (models.Result_count, error)
func (est *ResultStorage) GetList(ctx context.Context, pg, pgs int, gs1 string, ord int, dsc bool) (models.Result_count, error) {
	dbpool := pgclient.WDB
	gs := models.Result{}

	gsc := 0
	err := dbpool.QueryRow(ctx, "SELECT * from func_results_cnt($1);", gs1).Scan(&gsc)

	if err != nil {
		log.Println(err.Error(), "func_results_cnt")
		return models.Result_count{Values: []models.Result{}, Count: gsc, Auth: models.Auth{}}, err
	}

	out_arr := make([]models.Result, 0,
		func() int {
			if gsc < pgs {
				return gsc
			} else {
				return pgs
			}
		}())

	rows, err := dbpool.Query(ctx, "SELECT * from func_results_get($1,$2,$3,$4,$5);", pg, pgs, gs1, ord, dsc)
	if err != nil {
		log.Println(err.Error(), "func_results_get")
		return models.Result_count{Values: []models.Result{}, Count: gsc, Auth: models.Auth{}}, err
	}

	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&gs.Id, &gs.ResultName)
		if err != nil {
			log.Println("failed to scan row:", err)
		}

		out_arr = append(out_arr, gs)
	}

	out_count := models.Result_count{Values: out_arr, Count: gsc, Auth: models.Auth{}}
	return out_count, nil
}

//func (est *ResultStorage) Add(ctx context.Context, a models.Result) (int, error)
func (est *ResultStorage) Add(ctx context.Context, a models.Result) (int, error) {
	dbpool := pgclient.WDB
	ai := 0

	err := dbpool.QueryRow(ctx, "SELECT func_results_add($1);", a.ResultName).Scan(&ai)

	if err != nil {
		log.Println("Failed execute func_results_add: ", err)
		return 0, err
	}

	return ai, nil
}

//func (est *ResultStorage) Upd(ctx context.Context, u models.Result) (int, error)
func (est *ResultStorage) Upd(ctx context.Context, u models.Result) (int, error) {
	dbpool := pgclient.WDB
	ui := 0

	err := dbpool.QueryRow(ctx, "SELECT func_results_upd($1,$2);", u.Id, u.ResultName).Scan(&ui)

	if err != nil {
		log.Println("Failed execute func_results_upd: ", err)
		return 0, err
	}
	return ui, nil
}

//func (est *ResultStorage) Del(ctx context.Context, d []int) ([]int, error)
func (est *ResultStorage) Del(ctx context.Context, d []int) ([]int, error) {
	dbpool := pgclient.WDB
	res := []int{}
	i := 0
	for _, id := range d {
		err := dbpool.QueryRow(ctx, "SELECT func_results_del($1);", id).Scan(&i)
		res = append(res, i)

		if err != nil {
			log.Println("Failed execute func_results_del: ", err)
		}
	}
	return res, nil
}

//func (est *ResultStorage) GetOne(ctx context.Context, i int) (models.Result_count, error)
func (est *ResultStorage) GetOne(ctx context.Context, i int) (models.Result_count, error) {
	dbpool := pgclient.WDB
	out_arr := []models.Result{}
	g := models.Result{}

	err := dbpool.QueryRow(ctx, "SELECT * from func_result_get($1);", i).Scan(&g.Id, &g.ResultName)

	if err != nil && err != pgx.ErrNoRows {
		log.Println("Failed execute func_result_get: ", err)
		return models.Result_count{Values: []models.Result{}, Count: 0, Auth: models.Auth{}}, err
	}

	out_arr = append(out_arr, g)

	out_count := models.Result_count{Values: out_arr, Count: 1, Auth: models.Auth{}}
	return out_count, nil
}
