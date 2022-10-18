package pgsql

import (
	"context"
	"log"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/u2lentaru/billbck/internal/models"
	"github.com/u2lentaru/billbck/pkg/pgclient"
)

//type ObjStatusStorage struct
type ObjStatusStorage struct {
	db *pgxpool.Pool
}

//func NewObjStatusStorage(db *pgxpool.Pool) *ObjStatusStorage
func NewObjStatusStorage(db *pgxpool.Pool) *ObjStatusStorage {
	return &ObjStatusStorage{db: db}
}

//func (est *ObjStatusStorage) GetList(ctx context.Context, pg, pgs int, gs1 string, ord int, dsc bool) (models.ObjStatus_count, error)
func (est *ObjStatusStorage) GetList(ctx context.Context, pg, pgs int, gs1 string, ord int, dsc bool) (models.ObjStatus_count, error) {
	dbpool := pgclient.WDB
	gs := models.ObjStatus{}

	gsc := 0
	err := dbpool.QueryRow(ctx, "SELECT * from func_obj_statuses_cnt($1);", gs1).Scan(&gsc)

	if err != nil {
		log.Println(err.Error(), "func_obj_statuses_cnt")
		return models.ObjStatus_count{Values: []models.ObjStatus{}, Count: gsc, Auth: models.Auth{}}, err
	}

	out_arr := make([]models.ObjStatus, 0,
		func() int {
			if gsc < pgs {
				return gsc
			} else {
				return pgs
			}
		}())

	rows, err := dbpool.Query(ctx, "SELECT * from func_obj_statuses_get($1,$2,$3,$4,$5);", pg, pgs, gs1, ord, dsc)
	if err != nil {
		log.Println(err.Error())
		return models.ObjStatus_count{Values: []models.ObjStatus{}, Count: gsc, Auth: models.Auth{}}, err
	}

	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&gs.Id, &gs.ObjStatusName)
		if err != nil {
			log.Println("failed to scan row:", err)
		}

		out_arr = append(out_arr, gs)
	}

	out_count := models.ObjStatus_count{Values: out_arr, Count: gsc, Auth: models.Auth{}}
	if err != nil {
		log.Println(err.Error())
		return models.ObjStatus_count{}, err
	}

	return out_count, nil
}

//func (est *ObjStatusStorage) Add(ctx context.Context, a models.ObjStatus) (int, error)
func (est *ObjStatusStorage) Add(ctx context.Context, a models.ObjStatus) (int, error) {
	dbpool := pgclient.WDB
	ai := 0

	err := dbpool.QueryRow(ctx, "SELECT func_obj_statuses_add($1);", a.ObjStatusName).Scan(&ai)

	if err != nil {
		log.Println("Failed execute func_obj_statuses_add: ", err)
		return 0, err
	}

	return ai, nil
}

//func (est *ObjStatusStorage) Upd(ctx context.Context, u models.ObjStatus) (int, error)
func (est *ObjStatusStorage) Upd(ctx context.Context, u models.ObjStatus) (int, error) {
	dbpool := pgclient.WDB
	ui := 0

	err := dbpool.QueryRow(ctx, "SELECT func_obj_statuses_upd($1,$2);", u.Id, u.ObjStatusName).Scan(&ui)

	if err != nil {
		log.Println("Failed execute func_obj_statuses_upd: ", err)
		return 0, err
	}
	return ui, nil
}

//func (est *ObjStatusStorage) Del(ctx context.Context, d []int) ([]int, error)
func (est *ObjStatusStorage) Del(ctx context.Context, d []int) ([]int, error) {
	dbpool := pgclient.WDB
	res := []int{}
	i := 0
	for _, id := range d {
		err := dbpool.QueryRow(ctx, "SELECT func_obj_statuses_del($1);", id).Scan(&i)
		res = append(res, i)

		if err != nil {
			log.Println("Failed execute func_obj_statuses_del: ", err)
		}
	}
	return res, nil
}

//func (est *ObjStatusStorage) GetOne(ctx context.Context, i int) (models.ObjStatus_count, error)
func (est *ObjStatusStorage) GetOne(ctx context.Context, i int) (models.ObjStatus_count, error) {
	dbpool := pgclient.WDB
	out_arr := []models.ObjStatus{}
	g := models.ObjStatus{}

	err := dbpool.QueryRow(ctx, "SELECT * from func_obj_status_get($1);", i).Scan(&g.Id, &g.ObjStatusName)

	if err != nil && err != pgx.ErrNoRows {
		log.Println("Failed execute func_obj_status_get: ", err)
		return models.ObjStatus_count{Values: []models.ObjStatus{}, Count: 0, Auth: models.Auth{}}, err
	}

	out_arr = append(out_arr, g)

	out_count := models.ObjStatus_count{Values: out_arr, Count: 1, Auth: models.Auth{}}
	return out_count, nil
}
