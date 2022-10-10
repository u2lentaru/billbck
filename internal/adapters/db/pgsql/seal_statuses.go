package pgsql

import (
	"context"
	"log"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/u2lentaru/billbck/internal/models"
	"github.com/u2lentaru/billbck/pkg/pgclient"
)

//type SealStatusStorage struct
type SealStatusStorage struct {
	db *pgxpool.Pool
}

//func NewSealStatusStorage(db *pgxpool.Pool) *SealStatusStorage
func NewSealStatusStorage(db *pgxpool.Pool) *SealStatusStorage {
	return &SealStatusStorage{db: db}
}

//func (est *SealStatusStorage) GetList(ctx context.Context, pg, pgs int, gs1 string, ord int, dsc bool) (models.SealStatus_count, error)
func (est *SealStatusStorage) GetList(ctx context.Context, pg, pgs int, gs1 string, ord int, dsc bool) (models.SealStatus_count, error) {
	dbpool := pgclient.WDB
	gs := models.SealStatus{}

	gsc := 0
	err := dbpool.QueryRow(ctx, "SELECT * from func_seal_statuses_cnt($1);", gs1).Scan(&gsc)

	if err != nil {
		log.Println(err.Error(), "func_seal_statuses_cnt")
		return models.SealStatus_count{Values: []models.SealStatus{}, Count: gsc, Auth: models.Auth{}}, err
	}

	out_arr := make([]models.SealStatus, 0,
		func() int {
			if gsc < pgs {
				return gsc
			} else {
				return pgs
			}
		}())

	rows, err := dbpool.Query(ctx, "SELECT * from func_seal_statuses_get($1,$2,$3,$4,$5);", pg, pgs, gs1, ord, dsc)
	if err != nil {
		log.Println(err.Error(), "func_seal_statuses_get")
		return models.SealStatus_count{Values: []models.SealStatus{}, Count: gsc, Auth: models.Auth{}}, err
	}

	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&gs.Id, &gs.SealStatusName)
		if err != nil {
			log.Println("failed to scan row:", err)
		}

		out_arr = append(out_arr, gs)
	}

	out_count := models.SealStatus_count{Values: out_arr, Count: gsc, Auth: models.Auth{}}
	if err != nil {
		log.Println(err.Error())
		return models.SealStatus_count{}, err
	}

	return out_count, nil
}

//func (est *SealStatusStorage) Add(ctx context.Context, a models.SealStatus) (int, error)
func (est *SealStatusStorage) Add(ctx context.Context, a models.SealStatus) (int, error) {
	dbpool := pgclient.WDB
	ai := 0

	err := dbpool.QueryRow(ctx, "SELECT func_seal_statuses_add($1);", a.SealStatusName).Scan(&ai)

	if err != nil {
		log.Println("Failed execute func_seal_statuses_add: ", err)
		return 0, err
	}

	return ai, nil
}

//func (est *SealStatusStorage) Upd(ctx context.Context, u models.SealStatus) (int, error)
func (est *SealStatusStorage) Upd(ctx context.Context, u models.SealStatus) (int, error) {
	dbpool := pgclient.WDB
	ui := 0

	err := dbpool.QueryRow(ctx, "SELECT func_seal_statuses_upd($1,$2);", u.Id, u.SealStatusName).Scan(&ui)

	if err != nil {
		log.Println("Failed execute func_seal_statuses_upd: ", err)
		return 0, err
	}
	return ui, nil
}

//func (est *SealStatusStorage) Del(ctx context.Context, d []int) ([]int, error)
func (est *SealStatusStorage) Del(ctx context.Context, d []int) ([]int, error) {
	dbpool := pgclient.WDB
	res := []int{}
	i := 0
	for _, id := range d {
		err := dbpool.QueryRow(ctx, "SELECT func_seal_statuses_del($1);", id).Scan(&i)
		res = append(res, i)

		if err != nil {
			log.Println("Failed execute func_seal_statuses_del: ", err)
		}
	}
	return res, nil
}

//func (est *SealStatusStorage) GetOne(ctx context.Context, i int) (models.SealStatus_count, error)
func (est *SealStatusStorage) GetOne(ctx context.Context, i int) (models.SealStatus_count, error) {
	dbpool := pgclient.WDB
	out_arr := []models.SealStatus{}
	g := models.SealStatus{}

	err := dbpool.QueryRow(ctx, "SELECT * from func_seal_status_get($1);", i).Scan(&g.Id, &g.SealStatusName)

	if err != nil && err != pgx.ErrNoRows {
		log.Println("Failed execute func_seal_status_get: ", err)
		return models.SealStatus_count{Values: []models.SealStatus{}, Count: 0, Auth: models.Auth{}}, err
	}

	out_arr = append(out_arr, g)

	out_count := models.SealStatus_count{Values: out_arr, Count: 1, Auth: models.Auth{}}
	return out_count, nil
}
