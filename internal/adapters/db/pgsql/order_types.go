package pgsql

import (
	"context"
	"log"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/u2lentaru/billbck/internal/models"
	"github.com/u2lentaru/billbck/pkg/pgclient"
)

//type OrderTypeStorage struct
type OrderTypeStorage struct {
	db *pgxpool.Pool
}

//func NewOrderTypeStorage(db *pgxpool.Pool) *OrderTypeStorage
func NewOrderTypeStorage(db *pgxpool.Pool) *OrderTypeStorage {
	return &OrderTypeStorage{db: db}
}

//func (est *OrderTypeStorage) GetList(ctx context.Context, pg, pgs int, gs1 string, ord int, dsc bool) (models.OrderType_count, error)
func (est *OrderTypeStorage) GetList(ctx context.Context, pg, pgs int, gs1 string, ord int, dsc bool) (models.OrderType_count, error) {
	dbpool := pgclient.WDB
	gs := models.OrderType{}

	gsc := 0
	err := dbpool.QueryRow(ctx, "SELECT * from func_order_types_cnt($1);", gs1).Scan(&gsc)

	if err != nil {
		log.Println(err.Error(), "func_order_types_cnt")
		return models.OrderType_count{Values: []models.OrderType{}, Count: gsc, Auth: models.Auth{}}, err
	}

	out_arr := make([]models.OrderType, 0,
		func() int {
			if gsc < pgs {
				return gsc
			} else {
				return pgs
			}
		}())

	rows, err := dbpool.Query(ctx, "SELECT * from func_order_types_get($1,$2,$3,$4,$5);", pg, pgs, gs1, ord, dsc)
	if err != nil {
		log.Println(err.Error(), "func_order_types_get")
		return models.OrderType_count{Values: []models.OrderType{}, Count: gsc, Auth: models.Auth{}}, err
	}

	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&gs.Id, &gs.OrderTypeName)
		if err != nil {
			log.Println("failed to scan row:", err)
		}

		out_arr = append(out_arr, gs)
	}

	out_count := models.OrderType_count{Values: out_arr, Count: gsc, Auth: models.Auth{}}
	return out_count, nil
}

//func (est *OrderTypeStorage) Add(ctx context.Context, a models.OrderType) (int, error)
func (est *OrderTypeStorage) Add(ctx context.Context, a models.OrderType) (int, error) {
	dbpool := pgclient.WDB
	ai := 0

	err := dbpool.QueryRow(ctx, "SELECT func_order_types_add($1);", a.OrderTypeName).Scan(&ai)

	if err != nil {
		log.Println("Failed execute func_order_types_add: ", err)
		return 0, err
	}

	return ai, nil
}

//func (est *OrderTypeStorage) Upd(ctx context.Context, u models.OrderType) (int, error)
func (est *OrderTypeStorage) Upd(ctx context.Context, u models.OrderType) (int, error) {
	dbpool := pgclient.WDB
	ui := 0

	err := dbpool.QueryRow(ctx, "SELECT func_order_types_upd($1,$2);", u.Id, u.OrderTypeName).Scan(&ui)

	if err != nil {
		log.Println("Failed execute func_order_types_upd: ", err)
		return 0, err
	}
	return ui, nil
}

//func (est *OrderTypeStorage) Del(ctx context.Context, d []int) ([]int, error)
func (est *OrderTypeStorage) Del(ctx context.Context, d []int) ([]int, error) {
	dbpool := pgclient.WDB
	res := []int{}
	i := 0
	for _, id := range d {
		err := dbpool.QueryRow(ctx, "SELECT func_order_types_del($1);", id).Scan(&i)
		res = append(res, i)

		if err != nil {
			log.Println("Failed execute func_order_types_del: ", err)
		}
	}
	return res, nil
}

//func (est *OrderTypeStorage) GetOne(ctx context.Context, i int) (models.OrderType_count, error)
func (est *OrderTypeStorage) GetOne(ctx context.Context, i int) (models.OrderType_count, error) {
	dbpool := pgclient.WDB
	out_arr := []models.OrderType{}
	g := models.OrderType{}

	err := dbpool.QueryRow(ctx, "SELECT * from func_order_type_get($1);", i).Scan(&g.Id, &g.OrderTypeName)

	if err != nil && err != pgx.ErrNoRows {
		log.Println("Failed execute func_order_type_get: ", err)
		return models.OrderType_count{Values: []models.OrderType{}, Count: 0, Auth: models.Auth{}}, err
	}

	out_arr = append(out_arr, g)

	out_count := models.OrderType_count{Values: out_arr, Count: 1, Auth: models.Auth{}}
	return out_count, nil
}
