package pgsql

import (
	"context"
	"log"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/u2lentaru/billbck/internal/models"
	"github.com/u2lentaru/billbck/pkg/pgclient"
)

//type SubTypeStorage struct
type SubTypeStorage struct {
	db *pgxpool.Pool
}

//func NewSubTypeStorage(db *pgxpool.Pool) *SubTypeStorage
func NewSubTypeStorage(db *pgxpool.Pool) *SubTypeStorage {
	return &SubTypeStorage{db: db}
}

//func (est *SubTypeStorage) GetList(ctx context.Context, pg, pgs int, gs1, gs2 string, ord int, dsc bool) (models.SubType_count, error) {
func (est *SubTypeStorage) GetList(ctx context.Context, pg, pgs int, gs1, gs2 string, ord int, dsc bool) (models.SubType_count, error) {
	dbpool := pgclient.WDB
	gs := models.SubType{}

	gsc := 0
	err := dbpool.QueryRow(ctx, "SELECT * from func_sub_types_cnt($1,$2);", gs1, gs2).Scan(&gsc)

	if err != nil {
		log.Println(err.Error(), "func_sub_types_cnt")
		return models.SubType_count{Values: []models.SubType{}, Count: gsc, Auth: models.Auth{}}, err
	}

	out_arr := make([]models.SubType, 0,
		func() int {
			if gsc < pgs {
				return gsc
			} else {
				return pgs
			}
		}())

	rows, err := dbpool.Query(ctx, "SELECT * from func_sub_types_get($1,$2,$3,$4,$5,$6);", pg, pgs, gs1, gs2, ord, dsc)
	if err != nil {
		log.Println(err.Error(), "func_sub_types_get")
		return models.SubType_count{Values: []models.SubType{}, Count: gsc, Auth: models.Auth{}}, err
	}

	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&gs.Id, &gs.SubTypeName, &gs.SubTypeDescr)
		if err != nil {
			log.Println("failed to scan row:", err)
		}

		out_arr = append(out_arr, gs)
	}

	out_count := models.SubType_count{Values: out_arr, Count: gsc, Auth: models.Auth{}}
	return out_count, nil
}

//func (est *SubTypeStorage) Add(ctx context.Context, ea models.SubType) (int, error)
func (est *SubTypeStorage) Add(ctx context.Context, a models.SubType) (int, error) {
	dbpool := pgclient.WDB
	ai := 0

	err := dbpool.QueryRow(ctx, "SELECT func_sub_types_add($1,$2);", a.SubTypeDescr, a.SubTypeName).Scan(&ai)

	if err != nil {
		log.Println("Failed execute func_sub_types_add: ", err)
		return 0, err
	}

	return ai, nil
}

//func (est *SubTypeStorage) Upd(ctx context.Context, eu models.SubType) (int, error)
func (est *SubTypeStorage) Upd(ctx context.Context, u models.SubType) (int, error) {
	dbpool := pgclient.WDB
	ui := 0

	err := dbpool.QueryRow(ctx, "SELECT func_sub_types_upd($1,$2,$3);", u.Id, u.SubTypeName, u.SubTypeDescr).Scan(&ui)

	if err != nil {
		log.Println("Failed execute func_sub_types_upd: ", err)
		return 0, err
	}
	return ui, nil
}

//func (est *SubTypeStorage) Del(ctx context.Context, ed []int) ([]int, error)
func (est *SubTypeStorage) Del(ctx context.Context, d []int) ([]int, error) {
	dbpool := pgclient.WDB
	res := []int{}
	i := 0
	for _, id := range d {
		err := dbpool.QueryRow(ctx, "SELECT func_sub_types_del($1);", id).Scan(&i)
		res = append(res, i)

		if err != nil {
			log.Println("Failed execute func_sub_types_del: ", err)
		}
	}
	return res, nil
}

//func (est *SubTypeStorage) GetOne(ctx context.Context, i int) (models.SubType_count, error)
func (est *SubTypeStorage) GetOne(ctx context.Context, i int) (models.SubType_count, error) {
	dbpool := pgclient.WDB
	out_arr := []models.SubType{}
	g := models.SubType{}

	err := dbpool.QueryRow(ctx, "SELECT * from func_sub_type_get($1);", i).Scan(&g.Id, &g.SubTypeName, &g.SubTypeDescr)

	if err != nil && err != pgx.ErrNoRows {
		log.Println("Failed execute func_sub_type_get: ", err)
		return models.SubType_count{Values: []models.SubType{}, Count: 0, Auth: models.Auth{}}, err
	}

	out_arr = append(out_arr, g)

	out_count := models.SubType_count{Values: out_arr, Count: 1, Auth: models.Auth{}}
	return out_count, nil
}
