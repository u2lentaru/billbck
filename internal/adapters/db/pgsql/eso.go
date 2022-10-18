package pgsql

import (
	"context"
	"log"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/u2lentaru/billbck/internal/models"
	"github.com/u2lentaru/billbck/pkg/pgclient"
)

//type EsoStorage struct
type EsoStorage struct {
	db *pgxpool.Pool
}

//func NewEsoStorage(db *pgxpool.Pool) *EsoStorage
func NewEsoStorage(db *pgxpool.Pool) *EsoStorage {
	return &EsoStorage{db: db}
}

//func (est *EsoStorage) GetList(ctx context.Context, pg, pgs int, gs1 string, ord int, dsc bool) (models.Eso_count, error)
func (est *EsoStorage) GetList(ctx context.Context, pg, pgs int, gs1 string, ord int, dsc bool) (models.Eso_count, error) {
	dbpool := pgclient.WDB
	gs := models.Eso{}

	gsc := 0
	err := dbpool.QueryRow(ctx, "SELECT * from func_eso_cnt($1);", gs1).Scan(&gsc)

	if err != nil {
		log.Println(err.Error(), "func_eso_cnt")
		return models.Eso_count{Values: []models.Eso{}, Count: gsc, Auth: models.Auth{}}, err
	}

	out_arr := make([]models.Eso, 0,
		func() int {
			if gsc < pgs {
				return gsc
			} else {
				return pgs
			}
		}())

	rows, err := dbpool.Query(ctx, "SELECT * from func_eso_get($1,$2,$3,$4,$5);", pg, pgs, gs1, ord, dsc)
	if err != nil {
		log.Println(err.Error(), "func_eso_get")
		return models.Eso_count{Values: []models.Eso{}, Count: gsc, Auth: models.Auth{}}, err
	}

	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&gs.Id, &gs.EsoName)
		if err != nil {
			log.Println("failed to scan row:", err)
		}

		out_arr = append(out_arr, gs)
	}

	out_count := models.Eso_count{Values: out_arr, Count: gsc, Auth: models.Auth{}}
	return out_count, nil
}

//func (est *EsoStorage) Add(ctx context.Context, a models.Eso) (int, error)
func (est *EsoStorage) Add(ctx context.Context, a models.Eso) (int, error) {
	dbpool := pgclient.WDB
	ai := 0

	err := dbpool.QueryRow(ctx, "SELECT func_eso_add($1);", a.EsoName).Scan(&ai)

	if err != nil {
		log.Println("Failed execute func_eso_add: ", err)
		return 0, err
	}

	return ai, nil
}

//func (est *EsoStorage) Upd(ctx context.Context, u models.Eso) (int, error)
func (est *EsoStorage) Upd(ctx context.Context, u models.Eso) (int, error) {
	dbpool := pgclient.WDB
	ui := 0

	err := dbpool.QueryRow(ctx, "SELECT func_eso_upd($1,$2);", u.Id, u.EsoName).Scan(&ui)

	if err != nil {
		log.Println("Failed execute func_eso_upd: ", err)
		return 0, err
	}
	return ui, nil
}

//func (est *EsoStorage) Del(ctx context.Context, d []int) ([]int, error)
func (est *EsoStorage) Del(ctx context.Context, d []int) ([]int, error) {
	dbpool := pgclient.WDB
	res := []int{}
	i := 0
	for _, id := range d {
		err := dbpool.QueryRow(ctx, "SELECT func_eso_del($1);", id).Scan(&i)
		res = append(res, i)

		if err != nil {
			log.Println("Failed execute func_eso_del: ", err)
		}
	}
	return res, nil
}

//func (est *EsoStorage) GetOne(ctx context.Context, i int) (models.Eso_count, error)
func (est *EsoStorage) GetOne(ctx context.Context, i int) (models.Eso_count, error) {
	dbpool := pgclient.WDB
	out_arr := []models.Eso{}
	g := models.Eso{}

	err := dbpool.QueryRow(ctx, "SELECT * from func_eso_getbyid($1);", i).Scan(&g.Id, &g.EsoName)

	if err != nil && err != pgx.ErrNoRows {
		log.Println("Failed execute from func_eso_getbyid: ", err)
		return models.Eso_count{Values: []models.Eso{}, Count: 0, Auth: models.Auth{}}, err
	}

	out_arr = append(out_arr, g)

	out_count := models.Eso_count{Values: out_arr, Count: 1, Auth: models.Auth{}}
	return out_count, nil
}
