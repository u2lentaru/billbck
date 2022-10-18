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

//type SubBankStorage struct
type SubBankStorage struct {
	db *pgxpool.Pool
}

//func NewSubBankStorage(db *pgxpool.Pool) *SubBankStorage
func NewSubBankStorage(db *pgxpool.Pool) *SubBankStorage {
	return &SubBankStorage{db: db}
}

//func (est *SubBankStorage) GetList(ctx context.Context, pg, pgs int, gs1 string, gs2  int, gs3 string, ord int, dsc bool) (models.SubBank_count, error)
func (est *SubBankStorage) GetList(ctx context.Context, pg, pgs int, gs1 string, gs2 int, gs3 string, ord int, dsc bool) (models.SubBank_count, error) {
	dbpool := pgclient.WDB
	gs := models.SubBank{}

	gsc := 0
	err := dbpool.QueryRow(ctx, "SELECT * from func_sub_banks_cnt($1,$2,$3);", gs1, utils.NullableInt(int32(gs2)), gs3).Scan(&gsc)

	if err != nil {
		log.Println(err.Error(), "func_sub_banks_cnt")
		return models.SubBank_count{Values: []models.SubBank{}, Count: gsc, Auth: models.Auth{}}, err
	}

	out_arr := make([]models.SubBank, 0,
		func() int {
			if gsc < pgs {
				return gsc
			} else {
				return pgs
			}
		}())

	rows, err := dbpool.Query(ctx, "SELECT * from func_sub_banks_get($1,$2,$3,$4,$5,$6,$7);", pg, pgs, gs1, utils.NullableInt(int32(gs2)), gs3, ord, dsc)
	if err != nil {
		log.Println(err.Error(), "func_sub_banks_get")
		return models.SubBank_count{Values: []models.SubBank{}, Count: gsc, Auth: models.Auth{}}, err
	}

	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&gs.Id, &gs.Sub.Id, &gs.Bank.Id, &gs.AccNumber, &gs.Active, &gs.Sub.SBName, &gs.Bank.BankName)
		if err != nil {
			log.Println("failed to scan row:", err)
		}

		out_arr = append(out_arr, gs)
	}

	out_count := models.SubBank_count{Values: out_arr, Count: gsc, Auth: models.Auth{}}
	return out_count, nil
}

//func (est *SubBankStorage) Add(ctx context.Context, ea models.SubBank) (int, error)
func (est *SubBankStorage) Add(ctx context.Context, a models.SubBank) (int, error) {
	dbpool := pgclient.WDB
	ai := 0

	err := dbpool.QueryRow(ctx, "SELECT func_sub_banks_add($1,$2,$3);", a.Sub.Id, a.Bank.Id, a.AccNumber).Scan(&ai)

	if err != nil {
		log.Println("Failed execute func_sub_banks_add: ", err)
		return 0, err
	}

	return ai, nil
}

//func (est *SubBankStorage) Upd(ctx context.Context, eu models.SubBank) (int, error)
func (est *SubBankStorage) Upd(ctx context.Context, u models.SubBank) (int, error) {
	dbpool := pgclient.WDB
	ui := 0

	err := dbpool.QueryRow(ctx, "SELECT func_sub_banks_upd($1,$2,$3,$4);", u.Id, u.Sub.Id, u.Bank.Id, u.AccNumber).Scan(&ui)

	if err != nil {
		log.Println("Failed execute func_sub_banks_upd: ", err)
		return 0, err
	}
	return ui, nil
}

//func (est *SubBankStorage) Del(ctx context.Context, d []int) ([]int, error)
func (est *SubBankStorage) Del(ctx context.Context, d []int) ([]int, error) {
	dbpool := pgclient.WDB
	res := []int{}
	i := 0
	for _, id := range d {
		err := dbpool.QueryRow(ctx, "SELECT func_sub_banks_del($1);", id).Scan(&i)
		res = append(res, i)

		if err != nil {
			log.Println("Failed execute func_sub_banks_del: ", err)
		}
	}
	return res, nil
}

//func (est *SubBankStorage) GetOne(ctx context.Context, i int) (models.SubBank_count, error)
func (est *SubBankStorage) GetOne(ctx context.Context, i int) (models.SubBank_count, error) {
	dbpool := pgclient.WDB
	out_arr := []models.SubBank{}
	g := models.SubBank{}

	err := dbpool.QueryRow(ctx, "SELECT * from func_sub_bank_get($1);", i).Scan(&g.Id, &g.Sub.Id, &g.Bank.Id,
		&g.AccNumber, &g.Active, &g.Sub.SBName, &g.Bank.BankName)

	if err != nil && err != pgx.ErrNoRows {
		log.Println("Failed execute func_sub_bank_get: ", err)
		return models.SubBank_count{Values: []models.SubBank{}, Count: 0, Auth: models.Auth{}}, err
	}

	out_arr = append(out_arr, g)

	out_count := models.SubBank_count{Values: out_arr, Count: 1, Auth: models.Auth{}}
	return out_count, nil
}

//func (est *SubBankStorage) SetActive(ctx context.Context, i int) (int, error)
func (est *SubBankStorage) SetActive(ctx context.Context, i int) (int, error) {
	dbpool := pgclient.WDB
	si := 0

	err := dbpool.QueryRow(ctx, "SELECT func_sub_banks_set_active($1);", i).Scan(&si)

	if err != nil {
		log.Println("Failed execute func_sub_banks_set_active: ", err)
		return 0, err
	}
	return si, nil
}
