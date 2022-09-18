package pgsql

import (
	"context"
	"log"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/u2lentaru/billbck/internal/models"
	"github.com/u2lentaru/billbck/pkg/pgclient"
)

//type ActTypeStorage struct
type ActTypeStorage struct {
	db *pgxpool.Pool
}

//func NewActTypeStorage(db *pgxpool.Pool) *ActTypeStorage
func NewActTypeStorage(db *pgxpool.Pool) *ActTypeStorage {
	return &ActTypeStorage{db: db}
}

//func (est *ActTypeStorage) GetList(ctx context.Context, pg, pgs int, nm string, ord int, dsc bool) (models.ActType_count, error)
func (est *ActTypeStorage) GetList(ctx context.Context, pg, pgs int, nm string, ord int, dsc bool) (models.ActType_count, error) {
	dbpool := pgclient.WDB
	gsc := 0
	err := dbpool.QueryRow(ctx, "SELECT * from func_act_types_cnt($1);", nm).Scan(&gsc)
	auth := models.Auth{Create: true, Read: true, Update: true, Delete: true}
	e := models.ActType{}

	if err != nil {
		log.Println(err.Error(), "func_act_types_cnt")
		return models.ActType_count{Values: []models.ActType{}, Count: gsc, Auth: auth}, err
	}

	out_arr := make([]models.ActType, 0,
		func() int {
			if gsc < pgs {
				return gsc
			} else {
				return pgs
			}
		}())

	rows, err := dbpool.Query(ctx, "SELECT * from func_act_types_get($1,$2,$3,$4,$5);", pg, pgs, nm, ord, dsc)
	if err != nil {
		log.Println(err.Error())
		return models.ActType_count{Values: []models.ActType{}, Count: gsc, Auth: auth}, err
	}

	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&(e.Id), &(e.ActTypeName))
		if err != nil {
			log.Println("failed to scan row:", err)
		}

		out_arr = append(out_arr, e)
	}

	out_count := models.ActType_count{Values: out_arr, Count: gsc, Auth: auth}
	if err != nil {
		log.Println(err.Error())
		return models.ActType_count{}, err
	}

	return out_count, nil
}

//func (est *ActTypeStorage) Add(ctx context.Context, ea models.ActType) (int, error)
func (est *ActTypeStorage) Add(ctx context.Context, ea models.ActType) (int, error) {
	dbpool := pgclient.WDB
	ai := 0

	err := dbpool.QueryRow(ctx, "SELECT func_act_types_add($1);", ea.ActTypeName).Scan(&ai)

	if err != nil {
		log.Println("Failed execute func_act_types_add: ", err)
		return 0, err
	}

	return ai, nil
}

//func (est *ActTypeStorage) Upd(ctx context.Context, eu models.ActType) (int, error)
func (est *ActTypeStorage) Upd(ctx context.Context, eu models.ActType) (int, error) {
	dbpool := pgclient.WDB
	ui := 0

	err := dbpool.QueryRow(context.Background(), "SELECT func_act_types_upd($1,$2);", eu.Id, eu.ActTypeName).Scan(&ui)

	if err != nil {
		log.Println("Failed execute func_act_types_upd: ", err)
		return 0, err
	}
	return ui, nil
}

//func (est *ActTypeStorage) Del(ctx context.Context, ed []int) ([]int, error)
func (est *ActTypeStorage) Del(ctx context.Context, ed []int) ([]int, error) {
	dbpool := pgclient.WDB
	res := []int{}
	i := 0
	for _, id := range ed {
		err := dbpool.QueryRow(ctx, "SELECT func_act_types_del($1);", id).Scan(&i)
		res = append(res, i)

		if err != nil {
			log.Println("Failed execute func_act_types_del: ", err)
			return []int{}, err
		}
	}
	return res, nil
}

//func (est *ActTypeStorage) GetOne(ctx context.Context, Dbpool *pgxpool.Pool, i int) (models.ActType_count, error)
func (est *ActTypeStorage) GetOne(ctx context.Context, i int) (models.ActType_count, error) {
	dbpool := pgclient.WDB
	out_arr := []models.ActType{}
	e := models.ActType{}
	auth := models.Auth{Create: true, Read: true, Update: true, Delete: true}

	err := dbpool.QueryRow(context.Background(), "SELECT * from func_act_type_get($1);", i).Scan(&(e.Id), &(e.ActTypeName))

	if err != nil && err != pgx.ErrNoRows {
		log.Println("Failed execute from func_act_type_get: ", err)
		return models.ActType_count{Values: []models.ActType{}, Count: 0, Auth: auth}, err
	}

	out_arr = append(out_arr, e)

	out_count := models.ActType_count{Values: out_arr, Count: 0, Auth: auth}
	return out_count, nil
}
