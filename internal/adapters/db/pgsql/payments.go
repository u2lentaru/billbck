package pgsql

import (
	"context"
	"log"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/u2lentaru/billbck/internal/models"
	"github.com/u2lentaru/billbck/internal/utils"
	"github.com/u2lentaru/billbck/pkg/pgclient"
)

//type PaymentStorage struct
type PaymentStorage struct {
	db *pgxpool.Pool
}

//func NewPaymentStorage(db *pgxpool.Pool) *PaymentStorage
func NewPaymentStorage(db *pgxpool.Pool) *PaymentStorage {
	return &PaymentStorage{db: db}
}

//func (est *PaymentStorage) GetList(ctx context.Context, pg, pgs int, gs1, gs2 string, ord int, dsc bool) (models.Payment_count, error) {
func (est *PaymentStorage) GetList(ctx context.Context, pg, pgs int, gs1, gs2 string, ord int, dsc bool) (models.Payment_count, error) {
	dbpool := pgclient.WDB
	gs := models.Payment{}

	gsc := 0
	err := dbpool.QueryRow(ctx, "SELECT * from func_payments_cnt($1,$2);", gs1, utils.NullableString(gs2)).Scan(&gsc)

	if err != nil {
		log.Println(err.Error(), "func_payments_cnt")
		return models.Payment_count{Values: []models.Payment{}, Count: gsc, Auth: models.Auth{}}, err
	}

	out_arr := make([]models.Payment, 0,
		func() int {
			if gsc < pgs {
				return gsc
			} else {
				return pgs
			}
		}())

	rows, err := dbpool.Query(ctx, "SELECT * from func_payments_get($1,$2,$3,$4,$5,$6);", pg, pgs, gs1, utils.NullableString(gs2), ord, dsc)
	if err != nil {
		log.Println(err.Error())
		return models.Payment_count{Values: []models.Payment{}, Count: gsc, Auth: models.Auth{}}, err
	}

	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&gs.Id, &gs.PaymentDate, &gs.Contract.Id, &gs.Object.Id, &gs.PaymentType.Id, &gs.ChargeType.Id, &gs.Cashdesk.Id,
			&gs.BundleNumber, &gs.Amount, &gs.Contract.ContractNumber, &gs.Object.ObjectName, &gs.PaymentType.PaymentTypeName,
			&gs.ChargeType.ChargeTypeName, &gs.Cashdesk.CashdeskName)
		if err != nil {
			log.Println("failed to scan row:", err)
		}

		out_arr = append(out_arr, gs)
	}

	out_count := models.Payment_count{Values: out_arr, Count: gsc, Auth: models.Auth{}}
	return out_count, nil
}

//func (est *PaymentStorage) Add(ctx context.Context, ea models.Payment) (int, error)
func (est *PaymentStorage) Add(ctx context.Context, a models.Payment) (int, error) {
	dbpool := pgclient.WDB
	ai := 0

	err := dbpool.QueryRow(ctx, "SELECT func_payments_add($1,$2,$3,$4,$5,$6,$7,$8);", a.PaymentDate, a.Contract.Id,
		a.Object.Id, a.PaymentType.Id, a.ChargeType.Id, a.Cashdesk.Id, a.BundleNumber, a.Amount).Scan(&ai)

	if err != nil {
		log.Println("Failed execute func_payments_add: ", err)
		return 0, err
	}

	return ai, nil
}

//func (est *PaymentStorage) Upd(ctx context.Context, eu models.Payment) (int, error)
func (est *PaymentStorage) Upd(ctx context.Context, u models.Payment) (int, error) {
	dbpool := pgclient.WDB
	ui := 0

	err := dbpool.QueryRow(ctx, "SELECT func_payments_upd($1,$2,$3,$4,$5,$6,$7,$8,$9);", u.Id, u.PaymentDate, u.Contract.Id,
		u.Object.Id, u.PaymentType.Id, u.ChargeType.Id, u.Cashdesk.Id, u.BundleNumber, u.Amount).Scan(&ui)

	if err != nil {
		log.Println("Failed execute func_payments_upd: ", err)
		return 0, err
	}
	return ui, nil
}

//func (est *PaymentStorage) Del(ctx context.Context, ed []int) ([]int, error)
func (est *PaymentStorage) Del(ctx context.Context, d []int) ([]int, error) {
	dbpool := pgclient.WDB
	res := []int{}
	i := 0
	for _, id := range d {
		err := dbpool.QueryRow(ctx, "SELECT func_payments_del($1);", id).Scan(&i)
		res = append(res, i)

		if err != nil {
			log.Println("Failed execute func_payments_del: ", err)
		}
	}
	return res, nil
}

//func (est *PaymentStorage) GetOne(ctx context.Context, i int) (models.Payment_count, error)
func (est *PaymentStorage) GetOne(ctx context.Context, i int) (models.Payment_count, error) {
	dbpool := pgclient.WDB
	out_arr := []models.Payment{}
	g := models.Payment{}

	err := dbpool.QueryRow(ctx, "SELECT * from func_payment_get($1);", i).Scan(&g.Id, &g.PaymentDate, &g.Contract.Id,
		&g.Object.Id, &g.PaymentType.Id, &g.ChargeType.Id, &g.Cashdesk.Id, &g.BundleNumber, &g.Amount, &g.Contract.ContractNumber,
		&g.Object.ObjectName, &g.PaymentType.PaymentTypeName, &g.ChargeType.ChargeTypeName, &g.Cashdesk.CashdeskName)

	if err != nil && err != pgx.ErrNoRows {
		log.Println("Failed execute from func_payment_get: ", err)
		return models.Payment_count{Values: []models.Payment{}, Count: 0, Auth: models.Auth{}}, err
	}

	out_arr = append(out_arr, g)

	out_count := models.Payment_count{Values: out_arr, Count: 1, Auth: models.Auth{}}
	return out_count, nil
}
