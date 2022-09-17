package pg

import (
	"context"
	"log"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/u2lentaru/billbck/internal/models"
)

//type ActTypeStorage struct
type ActTypeStorage struct {
	db *pgxpool.Pool
}

//func NewActTypeStorage(db *pgxpool.Pool) *ActTypeStorage
func NewActTypeStorage(db *pgxpool.Pool) *ActTypeStorage {
	return &ActTypeStorage{db: db}
}

//func (est *ActTypeStorage) GetActType(ctx context.Context, Dbpool *pgxpool.Pool, pg, pgs int, nm string, ord int, dsc bool) (models.ActType_count, error)
func (est *ActTypeStorage) GetActType(ctx context.Context, Dbpool *pgxpool.Pool, pg, pgs int, nm string, ord int, dsc bool) (models.ActType_count, error) {
	gsc := 0
	err := Dbpool.QueryRow(ctx, "SELECT * from func_act_types_cnt($1);", nm).Scan(&gsc)
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

	rows, err := Dbpool.Query(ctx, "SELECT * from func_act_types_get($1,$2,$3,$4,$5);", pg, pgs, nm, ord, dsc)
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

//func (est *ActTypeStorage) AddActType(ctx context.Context, Dbpool *pgxpool.Pool) (int, error)
func (est *ActTypeStorage) AddActType(ctx context.Context, Dbpool *pgxpool.Pool) (int, error) {
	ai := 0
	e := models.ActType{}

	err := Dbpool.QueryRow(ctx, "SELECT func_act_types_add($1);", e.ActTypeName).Scan(&ai)

	if err != nil {
		log.Println("Failed execute func_act_types_add: ", err)
		return 0, err
	}

	return ai, nil
}

//func (est *ActTypeStorage) UpdActType(ctx context.Context, Dbpool *pgxpool.Pool)
func (est *ActTypeStorage) UpdActType(ctx context.Context, Dbpool *pgxpool.Pool) (int, error) {
	ui := 0
	e := models.ActType{}

	err := Dbpool.QueryRow(context.Background(), "SELECT func_act_types_upd($1,$2);", e.Id, e.ActTypeName).Scan(&ui)

	if err != nil {
		log.Println("Failed execute func_act_types_upd: ", err)
		return 0, err
	}
	return ui, nil
}

//func (est *ActTypeStorage) DelActType(ctx context.Context, Dbpool *pgxpool.Pool, d []int) ([]int, error)
func (est *ActTypeStorage) DelActType(ctx context.Context, Dbpool *pgxpool.Pool, d []int) ([]int, error) {
	res := []int{}
	i := 0
	for _, id := range d {
		err := Dbpool.QueryRow(ctx, "SELECT func_act_types_del($1);", id).Scan(&i)
		res = append(res, i)

		if err != nil {
			log.Println("Failed execute func_act_types_del: ", err)
			return []int{}, err
		}
	}
	return res, nil
}

//func (est *ActTypeStorage) GetActType(ctx context.Context, Dbpool *pgxpool.Pool, i int) (models.ActType_count, error)
func (est *ActTypeStorage) GetDistributionZone(ctx context.Context, Dbpool *pgxpool.Pool, i int) (models.ActType_count, error) {
	out_arr := []models.ActType{}
	e := models.ActType{}
	auth := models.Auth{Create: true, Read: true, Update: true, Delete: true}

	err := Dbpool.QueryRow(context.Background(), "SELECT * from func_act_type_get($1);", i).Scan(&(e.Id), &(e.ActTypeName))

	if err != nil && err != pgx.ErrNoRows {
		log.Println("Failed execute from func_act_type_get: ", err)
		return models.ActType_count{Values: []models.ActType{}, Count: 0, Auth: auth}, err
	}

	out_arr = append(out_arr, e)

	out_count := models.ActType_count{Values: out_arr, Count: 0, Auth: auth}
	return out_count, nil
}
