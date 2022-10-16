package pgsql

import (
	"context"
	"log"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/u2lentaru/billbck/internal/models"
	"github.com/u2lentaru/billbck/pkg/pgclient"
)

//type BalanceStorage struct
type BalanceStorage struct {
	db *pgxpool.Pool
}

//func NewBalanceStorage(db *pgxpool.Pool) *BalanceStorage
func NewBalanceStorage(db *pgxpool.Pool) *BalanceStorage {
	return &BalanceStorage{db: db}
}

//func (est *BalanceStorage) GetList(ctx context.Context, pg, pgs, gs1, gs2 int) (models.Balance_count, error) {
func (est *BalanceStorage) GetList(ctx context.Context, pg, pgs, gs1, gs2 int) (models.Balance_count, error) {
	dbpool := pgclient.WDB
	gs := models.Balance{}

	gsc := 0
	err := dbpool.QueryRow(ctx, "SELECT * from func_balance_cnt($1,$2);", gs1, gs2).Scan(&gsc)

	if err != nil {
		log.Println(err.Error(), "func_balance_cnt")
		return models.Balance_count{Values: []models.Balance{}, Count: gsc}, err
	}

	out_arr := make([]models.Balance, 0,
		func() int {
			if gsc < pgs {
				return gsc
			} else {
				return pgs
			}
		}())

	rows, err := dbpool.Query(ctx, "SELECT * from func_balance_get($1,$2,$3,$4);", pg, pgs, gs1, gs2)
	if err != nil {
		log.Println(err.Error(), "func_balance_get")
		return models.Balance_count{Values: []models.Balance{}, Count: gsc}, err
	}

	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&gs.Id, &gs.PId, &gs.BName, &gs.BTypeId, &gs.BTypeName, &gs.ChildCount, &gs.ReqId)
		if err != nil {
			log.Println("failed to scan row:", err)
		}

		out_arr = append(out_arr, gs)
	}

	out_count := models.Balance_count{Values: out_arr, Count: gsc}

	return out_count, nil
}

//func (est *BalanceStorage) GetNode(ctx context.Context, i, tid string) (models.Balance, error)
func (est *BalanceStorage) GetNode(ctx context.Context, i, tid string) (models.Balance, error) {
	dbpool := pgclient.WDB
	g := models.Balance{}

	err := dbpool.QueryRow(ctx, "SELECT * from func_balance_getbyid($1,$2);", i, tid).Scan(&g.Id, &g.PId, &g.BName, &g.BTypeId,
		&g.BTypeName, &g.ChildCount, &g.ReqId)

	if err != nil && err != pgx.ErrNoRows {
		log.Println("Failed execute func_balance_getbyid: ", err)
		return models.Balance{}, err
	}

	return g, nil
}

//func (est *BalanceStorage) GetNodeSum(ctx context.Context, gs1, gs2 int, gs3, gs4 string) (models.Json_sum, error)
func (est *BalanceStorage) GetNodeSum(ctx context.Context, gs1, gs2 int, gs3, gs4 string) (models.Json_sum, error) {
	dbpool := pgclient.WDB

	sum := 0.0
	err := dbpool.QueryRow(ctx, "SELECT * from func_balance_sum($1,$2,$3,$4);", gs1, gs2, gs3, gs4).Scan(&sum)
	if err != nil && err != pgx.ErrNoRows {
		log.Println("Failed execute func_balance_sum: ", err)
		return models.Json_sum{}, err
	}

	out_count := models.Json_sum{Sum: float32(sum)}

	return out_count, nil
}

//func (est *BalanceStorage) GetNodeSumL1(ctx context.Context, gs1, gs2 int, gs3, gs4 string) (models.Json_sum, error)
func (est *BalanceStorage) GetNodeSumL1(ctx context.Context, gs1, gs2 int, gs3, gs4 string) (models.Json_sum, error) {
	dbpool := pgclient.WDB

	sum := 0.0
	err := dbpool.QueryRow(ctx, "SELECT * from func_balance_sum1($1,$2,$3,$4);", gs1, gs2, gs3, gs4).Scan(&sum)
	if err != nil && err != pgx.ErrNoRows {
		log.Println("Failed execute func_balance_sum1: ", err)
		return models.Json_sum{}, err
	}

	out_count := models.Json_sum{Sum: float32(sum)}

	return out_count, nil
}

//func (est *BalanceStorage) GetNodeSumL0(ctx context.Context, gs1, gs2 int, gs3, gs4 string) (models.Json_sum, error)
func (est *BalanceStorage) GetNodeSumL0(ctx context.Context, gs1, gs2 int, gs3, gs4 string) (models.Json_sum, error) {
	dbpool := pgclient.WDB

	sum := 0.0
	err := dbpool.QueryRow(ctx, "SELECT * from func_balance_sum0($1,$2,$3,$4);", gs1, gs2, gs3, gs4).Scan(&sum)
	if err != nil && err != pgx.ErrNoRows {
		log.Println("Failed execute func_balance_sum0: ", err)
		return models.Json_sum{}, err
	}

	out_count := models.Json_sum{Sum: float32(sum)}

	return out_count, nil
}

//func (est *BalanceStorage) GetTabL1(ctx context.Context, pg, pgs, gs1, gs2 int, gs3, gs4 string) (models.BalanceTab_sum, error) {
func (est *BalanceStorage) GetTabL1(ctx context.Context, pg, pgs, gs1, gs2 int, gs3, gs4 string) (models.BalanceTab_sum, error) {
	dbpool := pgclient.WDB
	gs := models.BalanceTab{}
	out_arr := []models.BalanceTab{}

	rows, err := dbpool.Query(ctx, "SELECT * from func_balance_tab1($1,$2,$3,$4,$5,$6);", pg, pgs, gs1, gs2, gs3, gs4)
	if err != nil {
		log.Println(err.Error(), "func_balance_tab1")
		return models.BalanceTab_sum{}, err
	}

	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&gs.Id, &gs.PId, &gs.BName, &gs.BTypeId, &gs.BTypeName, &gs.Sum, &gs.ReqId)
		if err != nil {
			log.Println("failed to scan row:", err)
		}

		out_arr = append(out_arr, gs)
	}

	ogsc := -1.0

	//input node sum

	igsc := 0.0
	err = dbpool.QueryRow(ctx, "SELECT * from func_balance_sum($1,$2,$3,$4);", gs1, gs2, gs3, gs4).Scan(&igsc)
	if err != nil {
		log.Println(err.Error(), "func_balance_sum")
		return models.BalanceTab_sum{}, err
	}

	tab1cnt := 0
	err = dbpool.QueryRow(ctx, "SELECT * from func_balance_tab1_cnt($1,$2);", gs1, gs2).Scan(&tab1cnt)
	if err != nil {
		log.Println(err.Error(), "func_balance_tab1_cnt")
		return models.BalanceTab_sum{}, err
	}

	out_count := models.BalanceTab_sum{Values: out_arr, InSum: igsc, OutSum: ogsc, Count: tab1cnt}

	return out_count, nil
}

//func (est *BalanceStorage) GetTabL0(ctx context.Context, pg, pgs, gs1, gs2 int, gs3, gs4 string) (models.BalanceTab_sum, error) {
func (est *BalanceStorage) GetTabL0(ctx context.Context, pg, pgs, gs1, gs2 int, gs3, gs4 string) (models.BalanceTab_sum, error) {
	dbpool := pgclient.WDB
	gs := models.BalanceTab{}
	out_arr := []models.BalanceTab{}

	rows, err := dbpool.Query(ctx, "SELECT * from func_balance_tab0($1,$2,$3,$4,$5,$6);", pg, pgs, gs1, gs2, gs3, gs4)
	if err != nil {
		log.Println(err.Error(), "func_balance_tab0")
		return models.BalanceTab_sum{}, err
	}

	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&gs.Id, &gs.PId, &gs.BName, &gs.BTypeId, &gs.BTypeName, &gs.Sum, &gs.ReqId)
		if err != nil {
			log.Println("failed to scan row:", err)
		}

		out_arr = append(out_arr, gs)
	}

	ogsc := -1.0

	//input node sum

	igsc := 0.0
	err = dbpool.QueryRow(ctx, "SELECT * from func_balance_sum($1,$2,$3,$4);", gs1, gs2, gs3, gs4).Scan(&igsc)
	if err != nil {
		log.Println(err.Error(), "func_balance_sum")
		return models.BalanceTab_sum{}, err
	}

	tab0cnt := 0
	err = dbpool.QueryRow(ctx, "SELECT * from func_balance_tab0_cnt($1,$2);", gs1, gs2).Scan(&tab0cnt)
	if err != nil {
		log.Println(err.Error(), "func_balance_tab1_cnt")
		return models.BalanceTab_sum{}, err
	}

	out_count := models.BalanceTab_sum{Values: out_arr, InSum: igsc, OutSum: ogsc, Count: tab0cnt}

	return out_count, nil
}

//func (est *BalanceStorage) GetBranch(ctx context.Context, gs1, gs2 int) (models.BalanceTab_sum, error)
func (est *BalanceStorage) GetBranch(ctx context.Context, gs1, gs2 int) (models.BalanceTab_sum, error) {
	dbpool := pgclient.WDB
	gs := models.BalanceTab{}
	out_arr := []models.BalanceTab{}

	rows, err := dbpool.Query(ctx, "SELECT * from func_balance_getbranch($1,$2);", gs1, gs2)
	if err != nil {
		log.Println(err.Error(), "func_balance_getbranch")
		return models.BalanceTab_sum{}, err
	}

	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&gs.Id, &gs.PId, &gs.BName, &gs.BTypeId, &gs.BTypeName, &gs.Sum, &gs.ReqId)
		if err != nil {
			log.Println("failed to scan row:", err)
		}

		out_arr = append(out_arr, gs)
	}

	out_count := models.BalanceTab_sum{Values: out_arr, InSum: 0, OutSum: 0, Count: 0}

	return out_count, nil

}
