package pgsql

import (
	"context"
	"log"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/u2lentaru/billbck/internal/models"
	"github.com/u2lentaru/billbck/pkg/pgclient"
)

//type PuTypeStorage struct
type PuTypeStorage struct {
	db *pgxpool.Pool
}

//func NewPuTypeStorage(db *pgxpool.Pool) *PuTypeStorage
func NewPuTypeStorage(db *pgxpool.Pool) *PuTypeStorage {
	return &PuTypeStorage{db: db}
}

//func (est *PuTypeStorage) GetList(ctx context.Context, pg, pgs int, gs1 string, ord int, dsc bool) (models.PuType_count, error)
func (est *PuTypeStorage) GetList(ctx context.Context, pg, pgs int, gs1 string, ord int, dsc bool) (models.PuType_count, error) {
	dbpool := pgclient.WDB
	gs := models.PuType{}

	gsc := 0
	err := dbpool.QueryRow(ctx, "SELECT * from func_pu_types_cnt($1);", gs1).Scan(&gsc)

	if err != nil {
		log.Println(err.Error(), "func_pu_types_cnt")
		return models.PuType_count{Values: []models.PuType{}, Count: gsc, Auth: models.Auth{}}, err
	}

	out_arr := make([]models.PuType, 0,
		func() int {
			if gsc < pgs {
				return gsc
			} else {
				return pgs
			}
		}())

	rows, err := dbpool.Query(ctx, "SELECT * from func_pu_types_get($1,$2,$3,$4,$5);", pg, pgs, gs1, ord, dsc)
	if err != nil {
		log.Println(err.Error(), "func_pu_types_get")
		return models.PuType_count{Values: []models.PuType{}, Count: gsc, Auth: models.Auth{}}, err
	}

	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&gs.Id, &gs.PuTypeName)
		if err != nil {
			log.Println("failed to scan row:", err)
		}

		out_arr = append(out_arr, gs)
	}

	out_count := models.PuType_count{Values: out_arr, Count: gsc, Auth: models.Auth{}}
	return out_count, nil
}

//func (est *PuTypeStorage) Add(ctx context.Context, a models.PuType) (int, error)
func (est *PuTypeStorage) Add(ctx context.Context, a models.PuType) (int, error) {
	dbpool := pgclient.WDB
	ai := 0

	err := dbpool.QueryRow(ctx, "SELECT func_pu_types_add($1);", a.PuTypeName).Scan(&ai)

	if err != nil {
		log.Println("Failed execute func_pu_types_add: ", err)
		return 0, err
	}

	return ai, nil
}

//func (est *PuTypeStorage) Upd(ctx context.Context, u models.PuType) (int, error)
func (est *PuTypeStorage) Upd(ctx context.Context, u models.PuType) (int, error) {
	dbpool := pgclient.WDB
	ui := 0

	err := dbpool.QueryRow(ctx, "SELECT func_pu_types_upd($1,$2);", u.Id, u.PuTypeName).Scan(&ui)

	if err != nil {
		log.Println("Failed execute func_pu_types_upd: ", err)
		return 0, err
	}
	return ui, nil
}

//func (est *PuTypeStorage) Del(ctx context.Context, d []int) ([]int, error)
func (est *PuTypeStorage) Del(ctx context.Context, d []int) ([]int, error) {
	dbpool := pgclient.WDB
	res := []int{}
	i := 0
	for _, id := range d {
		err := dbpool.QueryRow(ctx, "SELECT func_pu_types_del($1);", id).Scan(&i)
		res = append(res, i)

		if err != nil {
			log.Println("Failed execute func_pu_types_del: ", err)
		}
	}
	return res, nil
}

//func (est *PuTypeStorage) GetOne(ctx context.Context, i int) (models.PuType_count, error)
func (est *PuTypeStorage) GetOne(ctx context.Context, i int) (models.PuType_count, error) {
	dbpool := pgclient.WDB
	out_arr := []models.PuType{}
	g := models.PuType{}

	err := dbpool.QueryRow(ctx, "SELECT * from func_pu_type_get($1);", i).Scan(&g.Id, &g.PuTypeName)

	if err != nil && err != pgx.ErrNoRows {
		log.Println("Failed execute from func_pu_type_get: ", err)
		return models.PuType_count{Values: []models.PuType{}, Count: 0, Auth: models.Auth{}}, err
	}

	out_arr = append(out_arr, g)

	out_count := models.PuType_count{Values: out_arr, Count: 1, Auth: models.Auth{}}
	return out_count, nil
}
