package pgsql

import (
	"context"
	"log"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/u2lentaru/billbck/internal/models"
	"github.com/u2lentaru/billbck/pkg/pgclient"
)

//type SealTypeStorage struct
type SealTypeStorage struct {
	db *pgxpool.Pool
}

//func NewSealTypeStorage(db *pgxpool.Pool) *SealTypeStorage
func NewSealTypeStorage(db *pgxpool.Pool) *SealTypeStorage {
	return &SealTypeStorage{db: db}
}

//func (est *SealTypeStorage) GetList(ctx context.Context, pg, pgs int, gs1 string, ord int, dsc bool) (models.SealType_count, error)
func (est *SealTypeStorage) GetList(ctx context.Context, pg, pgs int, gs1 string, ord int, dsc bool) (models.SealType_count, error) {
	dbpool := pgclient.WDB
	gs := models.SealType{}

	gsc := 0
	err := dbpool.QueryRow(ctx, "SELECT * from func_seal_types_cnt($1);", gs1).Scan(&gsc)

	if err != nil {
		log.Println(err.Error(), "func_seal_types_cnt")
		return models.SealType_count{Values: []models.SealType{}, Count: gsc, Auth: models.Auth{}}, err
	}

	out_arr := make([]models.SealType, 0,
		func() int {
			if gsc < pgs {
				return gsc
			} else {
				return pgs
			}
		}())

	rows, err := dbpool.Query(ctx, "SELECT * from func_seal_types_get($1,$2,$3,$4,$5);", pg, pgs, gs1, ord, dsc)
	if err != nil {
		log.Println(err.Error(), "func_seal_types_get")
		return models.SealType_count{Values: []models.SealType{}, Count: gsc, Auth: models.Auth{}}, err
	}

	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&gs.Id, &gs.SealTypeName)
		if err != nil {
			log.Println("failed to scan row:", err)
		}

		out_arr = append(out_arr, gs)
	}

	out_count := models.SealType_count{Values: out_arr, Count: gsc, Auth: models.Auth{}}
	if err != nil {
		log.Println(err.Error())
		return models.SealType_count{}, err
	}

	return out_count, nil
}

//func (est *SealTypeStorage) Add(ctx context.Context, a models.SealType) (int, error)
func (est *SealTypeStorage) Add(ctx context.Context, a models.SealType) (int, error) {
	dbpool := pgclient.WDB
	ai := 0

	err := dbpool.QueryRow(ctx, "SELECT func_seal_types_add($1);", a.SealTypeName).Scan(&ai)

	if err != nil {
		log.Println("Failed execute func_seal_types_add: ", err)
		return 0, err
	}

	return ai, nil
}

//func (est *SealTypeStorage) Upd(ctx context.Context, u models.SealType) (int, error)
func (est *SealTypeStorage) Upd(ctx context.Context, u models.SealType) (int, error) {
	dbpool := pgclient.WDB
	ui := 0

	err := dbpool.QueryRow(ctx, "SELECT func_seal_types_upd($1,$2);", u.Id, u.SealTypeName).Scan(&ui)

	if err != nil {
		log.Println("Failed execute func_seal_types_upd: ", err)
		return 0, err
	}
	return ui, nil
}

//func (est *SealTypeStorage) Del(ctx context.Context, d []int) ([]int, error)
func (est *SealTypeStorage) Del(ctx context.Context, d []int) ([]int, error) {
	dbpool := pgclient.WDB
	res := []int{}
	i := 0
	for _, id := range d {
		err := dbpool.QueryRow(ctx, "SELECT func_seal_types_del($1);", id).Scan(&i)
		res = append(res, i)

		if err != nil {
			log.Println("Failed execute func_seal_types_del: ", err)
		}
	}
	return res, nil
}

//func (est *SealTypeStorage) GetOne(ctx context.Context, i int) (models.SealType_count, error)
func (est *SealTypeStorage) GetOne(ctx context.Context, i int) (models.SealType_count, error) {
	dbpool := pgclient.WDB
	out_arr := []models.SealType{}
	g := models.SealType{}

	err := dbpool.QueryRow(ctx, "SELECT * from func_seal_type_get($1);", i).Scan(&g.Id, &g.SealTypeName)

	if err != nil && err != pgx.ErrNoRows {
		log.Println("Failed execute func_seal_type_get: ", err)
		return models.SealType_count{Values: []models.SealType{}, Count: 0, Auth: models.Auth{}}, err
	}

	out_arr = append(out_arr, g)

	out_count := models.SealType_count{Values: out_arr, Count: 1, Auth: models.Auth{}}
	return out_count, nil
}
