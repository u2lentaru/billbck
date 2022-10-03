package pgsql

import (
	"context"
	"log"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/u2lentaru/billbck/internal/models"
	"github.com/u2lentaru/billbck/pkg/pgclient"
)

//type ObjTypeStorage struct
type ObjTypeStorage struct {
	db *pgxpool.Pool
}

//func NewObjTypeStorage(db *pgxpool.Pool) *ObjTypeStorage
func NewObjTypeStorage(db *pgxpool.Pool) *ObjTypeStorage {
	return &ObjTypeStorage{db: db}
}

//func (est *ObjTypeStorage) GetList(ctx context.Context, pg, pgs int, gs1 string, ord int, dsc bool) (models.ObjType_count, error)
func (est *ObjTypeStorage) GetList(ctx context.Context, pg, pgs int, gs1 string, ord int, dsc bool) (models.ObjType_count, error) {
	dbpool := pgclient.WDB
	gs := models.ObjType{}

	gsc := 0
	err := dbpool.QueryRow(ctx, "SELECT * from func_obj_types_cnt($1);", gs1).Scan(&gsc)

	if err != nil {
		log.Println(err.Error(), "func_obj_types_cnt")
		return models.ObjType_count{Values: []models.ObjType{}, Count: gsc, Auth: models.Auth{}}, err
	}

	out_arr := make([]models.ObjType, 0,
		func() int {
			if gsc < pgs {
				return gsc
			} else {
				return pgs
			}
		}())

	rows, err := dbpool.Query(ctx, "SELECT * from func_obj_types_get($1,$2,$3,$4,$5);", pg, pgs, gs1, ord, dsc)
	if err != nil {
		log.Println(err.Error())
		return models.ObjType_count{Values: []models.ObjType{}, Count: gsc, Auth: models.Auth{}}, err
	}

	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&gs.Id, &gs.ObjTypeName)
		if err != nil {
			log.Println("failed to scan row:", err)
		}

		out_arr = append(out_arr, gs)
	}

	out_count := models.ObjType_count{Values: out_arr, Count: gsc, Auth: models.Auth{}}
	if err != nil {
		log.Println(err.Error())
		return models.ObjType_count{}, err
	}

	return out_count, nil
}

//func (est *ObjTypeStorage) Add(ctx context.Context, a models.ObjType) (int, error)
func (est *ObjTypeStorage) Add(ctx context.Context, a models.ObjType) (int, error) {
	dbpool := pgclient.WDB
	ai := 0

	err := dbpool.QueryRow(ctx, "SELECT func_obj_types_add($1);", a.ObjTypeName).Scan(&ai)

	if err != nil {
		log.Println("Failed execute func_obj_types_add: ", err)
		return 0, err
	}

	return ai, nil
}

//func (est *ObjTypeStorage) Upd(ctx context.Context, u models.ObjType) (int, error)
func (est *ObjTypeStorage) Upd(ctx context.Context, u models.ObjType) (int, error) {
	dbpool := pgclient.WDB
	ui := 0

	err := dbpool.QueryRow(ctx, "SELECT func_obj_types_upd($1,$2);", u.Id, u.ObjTypeName).Scan(&ui)

	if err != nil {
		log.Println("Failed execute func_obj_types_upd: ", err)
		return 0, err
	}
	return ui, nil
}

//func (est *ObjTypeStorage) Del(ctx context.Context, d []int) ([]int, error)
func (est *ObjTypeStorage) Del(ctx context.Context, d []int) ([]int, error) {
	dbpool := pgclient.WDB
	res := []int{}
	i := 0
	for _, id := range d {
		err := dbpool.QueryRow(ctx, "SELECT func_obj_types_del($1);", id).Scan(&i)
		res = append(res, i)

		if err != nil {
			log.Println("Failed execute func_obj_types_del: ", err)
		}
	}
	return res, nil
}

//func (est *ObjTypeStorage) GetOne(ctx context.Context, i int) (models.ObjType_count, error)
func (est *ObjTypeStorage) GetOne(ctx context.Context, i int) (models.ObjType_count, error) {
	dbpool := pgclient.WDB
	out_arr := []models.ObjType{}
	g := models.ObjType{}

	err := dbpool.QueryRow(ctx, "SELECT * from func_obj_type_get($1);", i).Scan(&g.Id, &g.ObjTypeName)

	if err != nil && err != pgx.ErrNoRows {
		log.Println("Failed execute from func_obj_type_get: ", err)
		return models.ObjType_count{Values: []models.ObjType{}, Count: 0, Auth: models.Auth{}}, err
	}

	out_arr = append(out_arr, g)

	out_count := models.ObjType_count{Values: out_arr, Count: 1, Auth: models.Auth{}}
	return out_count, nil
}
