package pgsql

import (
	"context"
	"log"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/u2lentaru/billbck/internal/models"
	"github.com/u2lentaru/billbck/pkg/pgclient"
)

//type ReasonStorage struct
type ReasonStorage struct {
	db *pgxpool.Pool
}

//func NewReasonStorage(db *pgxpool.Pool) *ReasonStorage
func NewReasonStorage(db *pgxpool.Pool) *ReasonStorage {
	return &ReasonStorage{db: db}
}

//func (est *ReasonStorage) GetList(ctx context.Context, pg, pgs int, gs1 string, ord int, dsc bool) (models.Reason_count, error)
func (est *ReasonStorage) GetList(ctx context.Context, pg, pgs int, gs1 string, ord int, dsc bool) (models.Reason_count, error) {
	dbpool := pgclient.WDB
	gs := models.Reason{}

	gsc := 0
	err := dbpool.QueryRow(ctx, "SELECT * from func_reasons_cnt($1);", gs1).Scan(&gsc)

	if err != nil {
		log.Println(err.Error(), "func_reasons_cnt")
		return models.Reason_count{Values: []models.Reason{}, Count: gsc, Auth: models.Auth{}}, err
	}

	out_arr := make([]models.Reason, 0,
		func() int {
			if gsc < pgs {
				return gsc
			} else {
				return pgs
			}
		}())

	rows, err := dbpool.Query(ctx, "SELECT * from func_reasons_get($1,$2,$3,$4,$5);", pg, pgs, gs1, ord, dsc)
	if err != nil {
		log.Println(err.Error(), "func_reasons_get")
		return models.Reason_count{Values: []models.Reason{}, Count: gsc, Auth: models.Auth{}}, err
	}

	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&gs.Id, &gs.ReasonName)
		if err != nil {
			log.Println("failed to scan row:", err)
		}

		out_arr = append(out_arr, gs)
	}

	out_count := models.Reason_count{Values: out_arr, Count: gsc, Auth: models.Auth{}}
	if err != nil {
		log.Println(err.Error())
		return models.Reason_count{}, err
	}

	return out_count, nil
}

//func (est *ReasonStorage) Add(ctx context.Context, a models.Reason) (int, error)
func (est *ReasonStorage) Add(ctx context.Context, a models.Reason) (int, error) {
	dbpool := pgclient.WDB
	ai := 0

	err := dbpool.QueryRow(ctx, "SELECT func_reasons_add($1);", a.ReasonName).Scan(&ai)

	if err != nil {
		log.Println("Failed execute func_reasons_add: ", err)
		return 0, err
	}

	return ai, nil
}

//func (est *ReasonStorage) Upd(ctx context.Context, u models.Reason) (int, error)
func (est *ReasonStorage) Upd(ctx context.Context, u models.Reason) (int, error) {
	dbpool := pgclient.WDB
	ui := 0

	err := dbpool.QueryRow(ctx, "SELECT func_reasons_upd($1,$2);", u.Id, u.ReasonName).Scan(&ui)

	if err != nil {
		log.Println("Failed execute func_reasons_upd: ", err)
		return 0, err
	}
	return ui, nil
}

//func (est *ReasonStorage) Del(ctx context.Context, d []int) ([]int, error)
func (est *ReasonStorage) Del(ctx context.Context, d []int) ([]int, error) {
	dbpool := pgclient.WDB
	res := []int{}
	i := 0
	for _, id := range d {
		err := dbpool.QueryRow(ctx, "SELECT func_reasons_del($1);", id).Scan(&i)
		res = append(res, i)

		if err != nil {
			log.Println("Failed execute func_reasons_del: ", err)
		}
	}
	return res, nil
}

//func (est *ReasonStorage) GetOne(ctx context.Context, i int) (models.Reason_count, error)
func (est *ReasonStorage) GetOne(ctx context.Context, i int) (models.Reason_count, error) {
	dbpool := pgclient.WDB
	out_arr := []models.Reason{}
	g := models.Reason{}

	err := dbpool.QueryRow(ctx, "SELECT * from func_reason_get($1);", i).Scan(&g.Id, &g.ReasonName)

	if err != nil && err != pgx.ErrNoRows {
		log.Println("Failed execute func_reason_get: ", err)
		return models.Reason_count{Values: []models.Reason{}, Count: 0, Auth: models.Auth{}}, err
	}

	out_arr = append(out_arr, g)

	out_count := models.Reason_count{Values: out_arr, Count: 1, Auth: models.Auth{}}
	return out_count, nil
}
