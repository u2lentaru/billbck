package pgsql

import (
	"context"
	"log"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/u2lentaru/billbck/internal/models"
	"github.com/u2lentaru/billbck/pkg/pgclient"
)

//type ClaimTypeStorage struct
type ClaimTypeStorage struct {
	db *pgxpool.Pool
}

//func NewClaimTypeStorage(db *pgxpool.Pool) *ClaimTypeStorage
func NewClaimTypeStorage(db *pgxpool.Pool) *ClaimTypeStorage {
	return &ClaimTypeStorage{db: db}
}

//func (est *ClaimTypeStorage) GetList(ctx context.Context, pg, pgs int, gs1 string, ord int, dsc bool) (models.ClaimType_count, error)
func (est *ClaimTypeStorage) GetList(ctx context.Context, pg, pgs int, gs1 string, ord int, dsc bool) (models.ClaimType_count, error) {
	dbpool := pgclient.WDB
	gs := models.ClaimType{}

	gsc := 0
	err := dbpool.QueryRow(ctx, "SELECT * from func_claim_types_cnt($1);", gs1).Scan(&gsc)

	if err != nil {
		log.Println(err.Error(), "func_claim_types_cnt")
		return models.ClaimType_count{Values: []models.ClaimType{}, Count: gsc, Auth: models.Auth{}}, err
	}

	out_arr := make([]models.ClaimType, 0,
		func() int {
			if gsc < pgs {
				return gsc
			} else {
				return pgs
			}
		}())

	rows, err := dbpool.Query(ctx, "SELECT * from func_claim_types_get($1,$2,$3,$4,$5);", pg, pgs, gs1, ord, dsc)
	if err != nil {
		log.Println(err.Error())
		return models.ClaimType_count{Values: []models.ClaimType{}, Count: gsc, Auth: models.Auth{}}, err
	}

	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&gs.Id, &gs.ClaimTypeName)
		if err != nil {
			log.Println("failed to scan row:", err)
		}

		out_arr = append(out_arr, gs)
	}

	out_count := models.ClaimType_count{Values: out_arr, Count: gsc, Auth: models.Auth{}}
	if err != nil {
		log.Println(err.Error())
		return models.ClaimType_count{}, err
	}

	return out_count, nil
}

//func (est *ClaimTypeStorage) Add(ctx context.Context, a models.ClaimType) (int, error)
func (est *ClaimTypeStorage) Add(ctx context.Context, a models.ClaimType) (int, error) {
	dbpool := pgclient.WDB
	ai := 0

	err := dbpool.QueryRow(context.Background(), "SELECT func_claim_types_add($1);", a.ClaimTypeName).Scan(&ai)

	if err != nil {
		log.Println("Failed execute func_claim_types_add: ", err)
		return 0, err
	}

	return ai, nil
}

//func (est *ClaimTypeStorage) Upd(ctx context.Context, u models.ClaimType) (int, error)
func (est *ClaimTypeStorage) Upd(ctx context.Context, u models.ClaimType) (int, error) {
	dbpool := pgclient.WDB
	ui := 0

	err := dbpool.QueryRow(context.Background(), "SELECT func_claim_types_upd($1,$2);", u.Id, u.ClaimTypeName).Scan(&ui)

	if err != nil {
		log.Println("Failed execute func_claim_types_upd: ", err)
		return 0, err
	}
	return ui, nil
}

//func (est *ClaimTypeStorage) Del(ctx context.Context, d []int) ([]int, error)
func (est *ClaimTypeStorage) Del(ctx context.Context, d []int) ([]int, error) {
	dbpool := pgclient.WDB
	res := []int{}
	i := 0
	for _, id := range d {
		err := dbpool.QueryRow(context.Background(), "SELECT func_claim_types_del($1);", id).Scan(&i)
		res = append(res, i)

		if err != nil {
			log.Println("Failed execute func_claim_types_del: ", err)
		}
	}
	return res, nil
}

//func (est *ClaimTypeStorage) GetOne(ctx context.Context, i int) (models.ClaimType_count, error)
func (est *ClaimTypeStorage) GetOne(ctx context.Context, i int) (models.ClaimType_count, error) {
	dbpool := pgclient.WDB
	out_arr := []models.ClaimType{}
	g := models.ClaimType{}

	err := dbpool.QueryRow(context.Background(), "SELECT * from func_claim_type_get($1);", i).Scan(&g.Id, &g.ClaimTypeName)

	if err != nil && err != pgx.ErrNoRows {
		log.Println("Failed execute from func_claim_type_get: ", err)
		return models.ClaimType_count{Values: []models.ClaimType{}, Count: 0, Auth: models.Auth{}}, err
	}

	out_arr = append(out_arr, g)

	out_count := models.ClaimType_count{Values: out_arr, Count: 1, Auth: models.Auth{}}
	return out_count, nil
}
