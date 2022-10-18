package pgsql

import (
	"context"
	"log"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/u2lentaru/billbck/internal/models"
	"github.com/u2lentaru/billbck/pkg/pgclient"
)

//type PaymentTypeStorage struct
type PaymentTypeStorage struct {
	db *pgxpool.Pool
}

//func NewPaymentTypeStorage(db *pgxpool.Pool) *PaymentTypeStorage
func NewPaymentTypeStorage(db *pgxpool.Pool) *PaymentTypeStorage {
	return &PaymentTypeStorage{db: db}
}

//func (est *PaymentTypeStorage) GetList(ctx context.Context, pg, pgs int, gs1 string, ord int, dsc bool) (models.PaymentType_count, error)
func (est *PaymentTypeStorage) GetList(ctx context.Context, pg, pgs int, gs1 string, ord int, dsc bool) (models.PaymentType_count, error) {
	dbpool := pgclient.WDB
	gs := models.PaymentType{}

	gsc := 0
	err := dbpool.QueryRow(ctx, "SELECT * from func_payment_types_cnt($1);", gs1).Scan(&gsc)

	if err != nil {
		log.Println(err.Error(), "func_payment_types_cnt")
		return models.PaymentType_count{Values: []models.PaymentType{}, Count: gsc, Auth: models.Auth{}}, err
	}

	out_arr := make([]models.PaymentType, 0,
		func() int {
			if gsc < pgs {
				return gsc
			} else {
				return pgs
			}
		}())

	rows, err := dbpool.Query(ctx, "SELECT * from func_payment_types_get($1,$2,$3,$4,$5);", pg, pgs, gs1, ord, dsc)
	if err != nil {
		log.Println(err.Error())
		return models.PaymentType_count{Values: []models.PaymentType{}, Count: gsc, Auth: models.Auth{}}, err
	}

	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&gs.Id, &gs.PaymentTypeName)
		if err != nil {
			log.Println("failed to scan row:", err)
		}

		out_arr = append(out_arr, gs)
	}

	out_count := models.PaymentType_count{Values: out_arr, Count: gsc, Auth: models.Auth{}}
	return out_count, nil
}

//func (est *PaymentTypeStorage) Add(ctx context.Context, a models.PaymentType) (int, error)
func (est *PaymentTypeStorage) Add(ctx context.Context, a models.PaymentType) (int, error) {
	dbpool := pgclient.WDB
	ai := 0

	err := dbpool.QueryRow(ctx, "SELECT func_payment_types_add($1);", a.PaymentTypeName).Scan(&ai)

	if err != nil {
		log.Println("Failed execute func_payment_types_add: ", err)
		return 0, err
	}

	return ai, nil
}

//func (est *PaymentTypeStorage) Upd(ctx context.Context, u models.PaymentType) (int, error)
func (est *PaymentTypeStorage) Upd(ctx context.Context, u models.PaymentType) (int, error) {
	dbpool := pgclient.WDB
	ui := 0

	err := dbpool.QueryRow(ctx, "SELECT func_payment_types_upd($1,$2);", u.Id, u.PaymentTypeName).Scan(&ui)

	if err != nil {
		log.Println("Failed execute func_payment_types_upd: ", err)
		return 0, err
	}
	return ui, nil
}

//func (est *PaymentTypeStorage) Del(ctx context.Context, d []int) ([]int, error)
func (est *PaymentTypeStorage) Del(ctx context.Context, d []int) ([]int, error) {
	dbpool := pgclient.WDB
	res := []int{}
	i := 0
	for _, id := range d {
		err := dbpool.QueryRow(ctx, "SELECT func_payment_types_del($1);", id).Scan(&i)
		res = append(res, i)

		if err != nil {
			log.Println("Failed execute func_payment_types_del: ", err)
		}
	}
	return res, nil
}

//func (est *PaymentTypeStorage) GetOne(ctx context.Context, i int) (models.PaymentType_count, error)
func (est *PaymentTypeStorage) GetOne(ctx context.Context, i int) (models.PaymentType_count, error) {
	dbpool := pgclient.WDB
	out_arr := []models.PaymentType{}
	g := models.PaymentType{}

	err := dbpool.QueryRow(ctx, "SELECT * from func_payment_type_get($1);", i).Scan(&g.Id, &g.PaymentTypeName)

	if err != nil && err != pgx.ErrNoRows {
		log.Println("Failed execute func_payment_type_get: ", err)
		return models.PaymentType_count{Values: []models.PaymentType{}, Count: 0, Auth: models.Auth{}}, err
	}

	out_arr = append(out_arr, g)

	out_count := models.PaymentType_count{Values: out_arr, Count: 1, Auth: models.Auth{}}
	return out_count, nil
}
