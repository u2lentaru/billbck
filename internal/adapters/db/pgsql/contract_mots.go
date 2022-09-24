package pgsql

import (
	"context"
	"log"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/u2lentaru/billbck/internal/models"
	"github.com/u2lentaru/billbck/pkg/pgclient"
)

//type ContractMotStorage struct
type ContractMotStorage struct {
	db *pgxpool.Pool
}

//func NewContractMotStorage(db *pgxpool.Pool) *ContractMotStorage
func NewContractMotStorage(db *pgxpool.Pool) *ContractMotStorage {
	return &ContractMotStorage{db: db}
}

//func (est *ContractMotStorage) GetList(ctx context.Context, pg, pgs int, gs1, gs2 string, ord int, dsc bool) (models.ContractMot_count, error) {
func (est *ContractMotStorage) GetList(ctx context.Context, pg, pgs int, gs1, gs2 string, ord int, dsc bool) (models.ContractMot_count, error) {
	dbpool := pgclient.WDB
	gs := models.ContractMot{}

	gsc := 0
	err := dbpool.QueryRow(ctx, "SELECT * from func_contract_mots_cnt($1,$2);", gs1, gs2).Scan(&gsc)

	if err != nil {
		log.Println(err.Error(), "func_contract_mots_cnt")
		return models.ContractMot_count{Values: []models.ContractMot{}, Count: gsc, Auth: models.Auth{}}, err
	}

	out_arr := make([]models.ContractMot, 0,
		func() int {
			if gsc < pgs {
				return gsc
			} else {
				return pgs
			}
		}())

	rows, err := dbpool.Query(ctx, "SELECT * from func_contract_mots_get($1,$2,$3,$4,$5,$6);", pg, pgs, gs1, gs2, ord, dsc)
	if err != nil {
		log.Println(err.Error())
		return models.ContractMot_count{Values: []models.ContractMot{}, Count: gsc, Auth: models.Auth{}}, err
	}

	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&gs.Id, &gs.ContractMotNameRu, &gs.ContractMotNameKz)
		if err != nil {
			log.Println("failed to scan row:", err)
		}

		out_arr = append(out_arr, gs)
	}

	out_count := models.ContractMot_count{Values: out_arr, Count: gsc, Auth: models.Auth{}}
	if err != nil {
		log.Println(err.Error())
		return models.ContractMot_count{}, err
	}

	return out_count, nil
}

//func (est *ContractMotStorage) Add(ctx context.Context, ea models.ContractMot) (int, error)
func (est *ContractMotStorage) Add(ctx context.Context, a models.ContractMot) (int, error) {
	dbpool := pgclient.WDB
	ai := 0

	err := dbpool.QueryRow(context.Background(), "SELECT func_contract_mots_add($1,$2);", a.ContractMotNameRu, a.ContractMotNameKz).Scan(&ai)

	if err != nil {
		log.Println("Failed execute func_contract_mots_add: ", err)
		return 0, err
	}

	return ai, nil
}

//func (est *ContractMotStorage) Upd(ctx context.Context, eu models.ContractMot) (int, error)
func (est *ContractMotStorage) Upd(ctx context.Context, u models.ContractMot) (int, error) {
	dbpool := pgclient.WDB
	ui := 0

	err := dbpool.QueryRow(context.Background(), "SELECT func_contract_mots_upd($1,$2,$3);", u.Id, u.ContractMotNameRu, u.ContractMotNameKz).Scan(&ui)

	if err != nil {
		log.Println("Failed execute func_contract_mots_upd: ", err)
		return 0, err
	}
	return ui, nil
}

//func (est *ContractMotStorage) Del(ctx context.Context, ed []int) ([]int, error)
func (est *ContractMotStorage) Del(ctx context.Context, d []int) ([]int, error) {
	dbpool := pgclient.WDB
	res := []int{}
	i := 0
	for _, id := range d {
		err := dbpool.QueryRow(context.Background(), "SELECT func_contract_mots_del($1);", id).Scan(&i)
		res = append(res, i)

		if err != nil {
			log.Println("Failed execute func_contract_mots_del: ", err)
		}
	}
	return res, nil
}

//func (est *ContractMotStorage) GetOne(ctx context.Context, i int) (models.ContractMot_count, error)
func (est *ContractMotStorage) GetOne(ctx context.Context, i int) (models.ContractMot_count, error) {
	dbpool := pgclient.WDB
	out_arr := []models.ContractMot{}
	g := models.ContractMot{}

	err := dbpool.QueryRow(context.Background(), "SELECT * from func_contract_mot_get($1);", i).Scan(&g.Id, &g.ContractMotNameRu, &g.ContractMotNameKz)

	if err != nil && err != pgx.ErrNoRows {
		log.Println("Failed execute from func_contract_mot_get: ", err)
		return models.ContractMot_count{Values: []models.ContractMot{}, Count: 0, Auth: models.Auth{}}, err
	}

	out_arr = append(out_arr, g)

	out_count := models.ContractMot_count{Values: out_arr, Count: 1, Auth: models.Auth{}}
	return out_count, nil
}
