package pgsql

import (
	"context"
	"log"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/u2lentaru/billbck/internal/models"
	"github.com/u2lentaru/billbck/pkg/pgclient"
)

//type PuValueTypeStorage struct
type PuValueTypeStorage struct {
	db *pgxpool.Pool
}

//func NewPuValueTypeStorage(db *pgxpool.Pool) *PuValueTypeStorage
func NewPuValueTypeStorage(db *pgxpool.Pool) *PuValueTypeStorage {
	return &PuValueTypeStorage{db: db}
}

//func (est *PuValueTypeStorage) GetList(ctx context.Context, pg, pgs int, gs1 string, ord int, dsc bool) (models.PuValueType_count, error)
func (est *PuValueTypeStorage) GetList(ctx context.Context, pg, pgs int, gs1 string, ord int, dsc bool) (models.PuValueType_count, error) {
	dbpool := pgclient.WDB
	gs := models.PuValueType{}

	gsc := 0
	err := dbpool.QueryRow(ctx, "SELECT * from func_pu_value_types_cnt($1);", gs1).Scan(&gsc)

	if err != nil {
		log.Println(err.Error(), "func_pu_value_types_cnt")
		return models.PuValueType_count{Values: []models.PuValueType{}, Count: gsc, Auth: models.Auth{}}, err
	}

	out_arr := make([]models.PuValueType, 0,
		func() int {
			if gsc < pgs {
				return gsc
			} else {
				return pgs
			}
		}())

	rows, err := dbpool.Query(ctx, "SELECT * from func_pu_value_types_get($1,$2,$3,$4,$5);", pg, pgs, gs1, ord, dsc)
	if err != nil {
		log.Println(err.Error(), "func_pu_value_types_get")
		return models.PuValueType_count{Values: []models.PuValueType{}, Count: gsc, Auth: models.Auth{}}, err
	}

	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&gs.Id, &gs.PuValueTypeName)
		if err != nil {
			log.Println("failed to scan row:", err)
		}

		out_arr = append(out_arr, gs)
	}

	out_count := models.PuValueType_count{Values: out_arr, Count: gsc, Auth: models.Auth{}}
	if err != nil {
		log.Println(err.Error())
		return models.PuValueType_count{}, err
	}

	return out_count, nil
}

//func (est *PuValueTypeStorage) Add(ctx context.Context, a models.PuValueType) (int, error)
func (est *PuValueTypeStorage) Add(ctx context.Context, a models.PuValueType) (int, error) {
	dbpool := pgclient.WDB
	ai := 0

	err := dbpool.QueryRow(ctx, "SELECT func_pu_value_types_add($1);", a.PuValueTypeName).Scan(&ai)

	if err != nil {
		log.Println("Failed execute func_pu_value_types_add: ", err)
		return 0, err
	}

	return ai, nil
}

//func (est *PuValueTypeStorage) Upd(ctx context.Context, u models.PuValueType) (int, error)
func (est *PuValueTypeStorage) Upd(ctx context.Context, u models.PuValueType) (int, error) {
	dbpool := pgclient.WDB
	ui := 0

	err := dbpool.QueryRow(ctx, "SELECT func_pu_value_types_upd($1,$2);", u.Id, u.PuValueTypeName).Scan(&ui)

	if err != nil {
		log.Println("Failed execute func_pu_value_types_upd: ", err)
		return 0, err
	}
	return ui, nil
}

//func (est *PuValueTypeStorage) Del(ctx context.Context, d []int) ([]int, error)
func (est *PuValueTypeStorage) Del(ctx context.Context, d []int) ([]int, error) {
	dbpool := pgclient.WDB
	res := []int{}
	i := 0
	for _, id := range d {
		err := dbpool.QueryRow(ctx, "SELECT func_pu_value_types_del($1);", id).Scan(&i)
		res = append(res, i)

		if err != nil {
			log.Println("Failed execute func_pu_value_types_del: ", err)
		}
	}
	return res, nil
}

//func (est *PuValueTypeStorage) GetOne(ctx context.Context, i int) (models.PuValueType_count, error)
func (est *PuValueTypeStorage) GetOne(ctx context.Context, i int) (models.PuValueType_count, error) {
	dbpool := pgclient.WDB
	out_arr := []models.PuValueType{}
	g := models.PuValueType{}

	err := dbpool.QueryRow(ctx, "SELECT * from func_pu_value_type_get($1);", i).Scan(&g.Id, &g.PuValueTypeName)

	if err != nil && err != pgx.ErrNoRows {
		log.Println("Failed execute func_pu_value_type_get: ", err)
		return models.PuValueType_count{Values: []models.PuValueType{}, Count: 0, Auth: models.Auth{}}, err
	}

	out_arr = append(out_arr, g)

	out_count := models.PuValueType_count{Values: out_arr, Count: 1, Auth: models.Auth{}}
	return out_count, nil
}
