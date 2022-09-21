package pgsql

import (
	"context"
	"log"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/u2lentaru/billbck/internal/models"
	"github.com/u2lentaru/billbck/pkg/pgclient"
)

//type BankStorage struct
type BankStorage struct {
	db *pgxpool.Pool
}

//func NewBankStorage(db *pgxpool.Pool) *BankStorage
func NewBankStorage(db *pgxpool.Pool) *BankStorage {
	return &BankStorage{db: db}
}

//func (est *BankStorage) GetList(ctx context.Context, pg, pgs int, gs1, gs2 string, ord int, dsc bool) (models.Bank_count, error)
func (est *BankStorage) GetList(ctx context.Context, pg, pgs int, gs1, gs2 string, ord int, dsc bool) (models.Bank_count, error) {
	dbpool := pgclient.WDB
	auth := models.Auth{Create: true, Read: true, Update: true, Delete: true}
	gs := models.Bank{}

	gsc := 0
	err := dbpool.QueryRow(ctx, "SELECT * from func_banks_cnt($1,$2);", gs1, gs2).Scan(&gsc)

	if err != nil {
		log.Println(err.Error(), "func_banks_cnt")
		return models.Bank_count{Values: []models.Bank{}, Count: gsc, Auth: auth}, err
	}

	out_arr := make([]models.Bank, 0,
		func() int {
			if gsc < pgs {
				return gsc
			} else {
				return pgs
			}
		}())

	rows, err := dbpool.Query(ctx, "SELECT * from func_banks_get($1,$2,$3,$4,$5,$6);", pg, pgs, gs1, gs2, ord, dsc)
	if err != nil {
		log.Println(err.Error())
		return models.Bank_count{Values: []models.Bank{}, Count: gsc, Auth: auth}, err
	}

	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&gs.Id, &gs.BankName, &gs.BankDescr, &gs.Mfo)
		if err != nil {
			log.Println("failed to scan row:", err)
		}

		out_arr = append(out_arr, gs)
	}

	out_count := models.Bank_count{Values: out_arr, Count: gsc, Auth: auth}
	if err != nil {
		log.Println(err.Error())
		return models.Bank_count{}, err
	}

	return out_count, nil
}

//func (est *BankStorage) Add(ctx context.Context, a models.Bank) (int, error)
func (est *BankStorage) Add(ctx context.Context, a models.Bank) (int, error) {
	dbpool := pgclient.WDB
	ai := 0

	err := dbpool.QueryRow(context.Background(), "SELECT func_banks_add($1,$2,$3);", a.BankName, a.BankDescr, a.Mfo).Scan(&ai)

	if err != nil {
		log.Println("Failed execute func_banks_add: ", err)
		return 0, err
	}

	return ai, nil
}

//func (est *BankStorage) Upd(ctx context.Context, u models.Bank) (int, error)
func (est *BankStorage) Upd(ctx context.Context, u models.Bank) (int, error) {
	dbpool := pgclient.WDB
	ui := 0

	err := dbpool.QueryRow(context.Background(), "SELECT func_banks_upd($1,$2,$3,$4);", u.Id, u.BankName, u.BankDescr, u.Mfo).Scan(&ui)

	if err != nil {
		log.Println("Failed execute func_banks_upd: ", err)
		return 0, err
	}
	return ui, nil
}

//func (est *BankStorage) Del(ctx context.Context, d []int) ([]int, error)
func (est *BankStorage) Del(ctx context.Context, d []int) ([]int, error) {
	dbpool := pgclient.WDB
	res := []int{}
	i := 0
	for _, id := range d {
		err := dbpool.QueryRow(context.Background(), "SELECT func_banks_del($1);", id).Scan(&i)
		res = append(res, i)

		if err != nil {
			log.Println("Failed execute func_banks_del: ", err)
		}
	}
	return res, nil
}

//func (est *BankStorage) GetOne(ctx context.Context, i int) (models.Bank_count, error)
func (est *BankStorage) GetOne(ctx context.Context, i int) (models.Bank_count, error) {
	dbpool := pgclient.WDB
	out_arr := []models.Bank{}
	g := models.Bank{}
	auth := models.Auth{Create: true, Read: true, Update: true, Delete: true}

	err := dbpool.QueryRow(context.Background(), "SELECT * from func_bank_get($1);", i).Scan(&g.Id, &g.BankName, &g.BankDescr, &g.Mfo)

	if err != nil && err != pgx.ErrNoRows {
		log.Println("Failed execute from func_bank_get: ", err)
		return models.Bank_count{Values: []models.Bank{}, Count: 0, Auth: auth}, err
	}

	out_arr = append(out_arr, g)

	out_count := models.Bank_count{Values: out_arr, Count: 0, Auth: auth}
	return out_count, nil
}
