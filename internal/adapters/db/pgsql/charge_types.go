package pgsql

import (
	"context"
	"log"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/u2lentaru/billbck/internal/models"
	"github.com/u2lentaru/billbck/pkg/pgclient"
)

//type ChargeTypeStorage struct
type ChargeTypeStorage struct {
	db *pgxpool.Pool
}

//func NewChargeTypeStorage(db *pgxpool.Pool) *ChargeTypeStorage
func NewChargeTypeStorage(db *pgxpool.Pool) *ChargeTypeStorage {
	return &ChargeTypeStorage{db: db}
}

//func (est *ChargeTypeStorage) GetList(ctx context.Context, pg, pgs int, gs1 string, ord int, dsc bool) (models.ChargeType_count, error)
func (est *ChargeTypeStorage) GetList(ctx context.Context, pg, pgs int, gs1 string, ord int, dsc bool) (models.ChargeType_count, error) {
	dbpool := pgclient.WDB
	gs := models.ChargeType{}

	gsc := 0
	err := dbpool.QueryRow(ctx, "SELECT * from func_charge_types_cnt($1);", gs1).Scan(&gsc)

	if err != nil {
		log.Println(err.Error(), "func_charge_types_cnt")
		return models.ChargeType_count{Values: []models.ChargeType{}, Count: gsc, Auth: models.Auth{}}, err
	}

	out_arr := make([]models.ChargeType, 0,
		func() int {
			if gsc < pgs {
				return gsc
			} else {
				return pgs
			}
		}())

	rows, err := dbpool.Query(ctx, "SELECT * from func_charge_types_get($1,$2,$3,$4,$5);", pg, pgs, gs1, ord, dsc)
	if err != nil {
		log.Println(err.Error())
		return models.ChargeType_count{Values: []models.ChargeType{}, Count: gsc, Auth: models.Auth{}}, err
	}

	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&gs.Id, &gs.ChargeTypeName)
		if err != nil {
			log.Println("failed to scan row:", err)
		}

		out_arr = append(out_arr, gs)
	}

	out_count := models.ChargeType_count{Values: out_arr, Count: gsc, Auth: models.Auth{}}
	if err != nil {
		log.Println(err.Error())
		return models.ChargeType_count{}, err
	}

	return out_count, nil
}

//func (est *ChargeTypeStorage) Add(ctx context.Context, a models.ChargeType) (int, error)
func (est *ChargeTypeStorage) Add(ctx context.Context, a models.ChargeType) (int, error) {
	dbpool := pgclient.WDB
	ai := 0

	err := dbpool.QueryRow(ctx, "SELECT func_charge_types_add($1);", a.ChargeTypeName).Scan(&ai)

	if err != nil {
		log.Println("Failed execute func_charge_types_add: ", err)
		return 0, err
	}

	return ai, nil
}

//func (est *ChargeTypeStorage) Upd(ctx context.Context, u models.ChargeType) (int, error)
func (est *ChargeTypeStorage) Upd(ctx context.Context, u models.ChargeType) (int, error) {
	dbpool := pgclient.WDB
	ui := 0

	err := dbpool.QueryRow(ctx, "SELECT func_charge_types_upd($1,$2);", u.Id, u.ChargeTypeName).Scan(&ui)

	if err != nil {
		log.Println("Failed execute func_charge_types_upd: ", err)
		return 0, err
	}
	return ui, nil
}

//func (est *ChargeTypeStorage) Del(ctx context.Context, d []int) ([]int, error)
func (est *ChargeTypeStorage) Del(ctx context.Context, d []int) ([]int, error) {
	dbpool := pgclient.WDB
	res := []int{}
	i := 0
	for _, id := range d {
		err := dbpool.QueryRow(ctx, "SELECT func_charge_types_del($1);", id).Scan(&i)
		res = append(res, i)

		if err != nil {
			log.Println("Failed execute func_charge_types_del: ", err)
		}
	}
	return res, nil
}

//func (est *ChargeTypeStorage) GetOne(ctx context.Context, i int) (models.ChargeType_count, error)
func (est *ChargeTypeStorage) GetOne(ctx context.Context, i int) (models.ChargeType_count, error) {
	dbpool := pgclient.WDB
	out_arr := []models.ChargeType{}
	g := models.ChargeType{}

	err := dbpool.QueryRow(ctx, "SELECT * from func_charge_type_get($1);", i).Scan(&g.Id, &g.ChargeTypeName)

	if err != nil && err != pgx.ErrNoRows {
		log.Println("Failed execute from func_charge_type_get: ", err)
		return models.ChargeType_count{Values: []models.ChargeType{}, Count: 0, Auth: models.Auth{}}, err
	}

	out_arr = append(out_arr, g)

	out_count := models.ChargeType_count{Values: out_arr, Count: 1, Auth: models.Auth{}}
	return out_count, nil
}
