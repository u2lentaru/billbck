package pgsql

import (
	"context"
	"log"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/u2lentaru/billbck/internal/models"
	"github.com/u2lentaru/billbck/pkg/pgclient"
)

//type ContractTypeStorage struct
type ContractTypeStorage struct {
	db *pgxpool.Pool
}

//func NewContractTypeStorage(db *pgxpool.Pool) *ContractTypeStorage
func NewContractTypeStorage(db *pgxpool.Pool) *ContractTypeStorage {
	return &ContractTypeStorage{db: db}
}

//func (est *ContractTypeStorage) GetList(ctx context.Context, pg, pgs int, gs1 string, ord int, dsc bool) (models.ContractType_count, error)
func (est *ContractTypeStorage) GetList(ctx context.Context, pg, pgs int, gs1 string, ord int, dsc bool) (models.ContractType_count, error) {
	dbpool := pgclient.WDB
	gs := models.ContractType{}

	gsc := 0
	err := dbpool.QueryRow(ctx, "SELECT * from func_contract_types_cnt($1);", gs1).Scan(&gsc)

	if err != nil {
		log.Println(err.Error(), "func_contract_types_cnt")
		return models.ContractType_count{Values: []models.ContractType{}, Count: gsc, Auth: models.Auth{}}, err
	}

	out_arr := make([]models.ContractType, 0,
		func() int {
			if gsc < pgs {
				return gsc
			} else {
				return pgs
			}
		}())

	rows, err := dbpool.Query(ctx, "SELECT * from func_contract_types_get($1,$2,$3,$4,$5);", pg, pgs, gs1, ord, dsc)
	if err != nil {
		log.Println(err.Error())
		return models.ContractType_count{Values: []models.ContractType{}, Count: gsc, Auth: models.Auth{}}, err
	}

	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&gs.Id, &gs.ContractTypeName)
		if err != nil {
			log.Println("failed to scan row:", err)
		}

		out_arr = append(out_arr, gs)
	}

	out_count := models.ContractType_count{Values: out_arr, Count: gsc, Auth: models.Auth{}}
	if err != nil {
		log.Println(err.Error())
		return models.ContractType_count{}, err
	}

	return out_count, nil
}

//func (est *ContractTypeStorage) Add(ctx context.Context, a models.ContractType) (int, error)
func (est *ContractTypeStorage) Add(ctx context.Context, a models.ContractType) (int, error) {
	dbpool := pgclient.WDB
	ai := 0

	err := dbpool.QueryRow(ctx, "SELECT func_contract_types_add($1);", a.ContractTypeName).Scan(&ai)

	if err != nil {
		log.Println("Failed execute func_contract_types_add: ", err)
		return 0, err
	}

	return ai, nil
}

//func (est *ContractTypeStorage) Upd(ctx context.Context, u models.ContractType) (int, error)
func (est *ContractTypeStorage) Upd(ctx context.Context, u models.ContractType) (int, error) {
	dbpool := pgclient.WDB
	ui := 0

	err := dbpool.QueryRow(ctx, "SELECT func_contract_types_upd($1,$2);", u.Id, u.ContractTypeName).Scan(&ui)

	if err != nil {
		log.Println("Failed execute func_contract_types_upd: ", err)
		return 0, err
	}
	return ui, nil
}

//func (est *ContractTypeStorage) Del(ctx context.Context, d []int) ([]int, error)
func (est *ContractTypeStorage) Del(ctx context.Context, d []int) ([]int, error) {
	dbpool := pgclient.WDB
	res := []int{}
	i := 0
	for _, id := range d {
		err := dbpool.QueryRow(ctx, "SELECT func_contract_types_del($1);", id).Scan(&i)
		res = append(res, i)

		if err != nil {
			log.Println("Failed execute func_contract_types_del: ", err)
		}
	}
	return res, nil
}

//func (est *ContractTypeStorage) GetOne(ctx context.Context, i int) (models.ContractType_count, error)
func (est *ContractTypeStorage) GetOne(ctx context.Context, i int) (models.ContractType_count, error) {
	dbpool := pgclient.WDB
	out_arr := []models.ContractType{}
	g := models.ContractType{}

	err := dbpool.QueryRow(ctx, "SELECT * from func_contract_type_get($1);", i).Scan(&g.Id, &g.ContractTypeName)

	if err != nil && err != pgx.ErrNoRows {
		log.Println("Failed execute from func_contract_type_get: ", err)
		return models.ContractType_count{Values: []models.ContractType{}, Count: 0, Auth: models.Auth{}}, err
	}

	out_arr = append(out_arr, g)

	out_count := models.ContractType_count{Values: out_arr, Count: 1, Auth: models.Auth{}}
	return out_count, nil
}
