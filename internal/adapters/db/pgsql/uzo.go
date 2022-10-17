package pgsql

import (
	"context"
	"log"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/u2lentaru/billbck/internal/models"
	"github.com/u2lentaru/billbck/pkg/pgclient"
)

//type UzoStorage struct
type UzoStorage struct {
	db *pgxpool.Pool
}

//func NewUzoStorage(db *pgxpool.Pool) *UzoStorage
func NewUzoStorage(db *pgxpool.Pool) *UzoStorage {
	return &UzoStorage{db: db}
}

//func (est *UzoStorage) GetList(ctx context.Context, pg, pgs int, gs1 string, ord int, dsc bool) (models.Uzo_count, error)
func (est *UzoStorage) GetList(ctx context.Context, pg, pgs int, gs1 string, ord int, dsc bool) (models.Uzo_count, error) {
	dbpool := pgclient.WDB
	gs := models.Uzo{}

	gsc := 0
	err := dbpool.QueryRow(ctx, "SELECT * from func_uzo_cnt($1);", gs1).Scan(&gsc)

	if err != nil {
		log.Println(err.Error(), "func_uzo_cnt")
		return models.Uzo_count{Values: []models.Uzo{}, Count: gsc, Auth: models.Auth{}}, err
	}

	out_arr := make([]models.Uzo, 0,
		func() int {
			if gsc < pgs {
				return gsc
			} else {
				return pgs
			}
		}())

	rows, err := dbpool.Query(ctx, "SELECT * from func_uzo_get($1,$2,$3,$4,$5);", pg, pgs, gs1, ord, dsc)
	if err != nil {
		log.Println(err.Error(), "func_uzo_get")
		return models.Uzo_count{Values: []models.Uzo{}, Count: gsc, Auth: models.Auth{}}, err
	}

	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&gs.Id, &gs.UzoName, &gs.UzoValue)
		if err != nil {
			log.Println("failed to scan row:", err)
		}

		out_arr = append(out_arr, gs)
	}

	out_count := models.Uzo_count{Values: out_arr, Count: gsc, Auth: models.Auth{}}

	return out_count, nil
}

//func (est *UzoStorage) Add(ctx context.Context, a models.Uzo) (int, error)
func (est *UzoStorage) Add(ctx context.Context, a models.Uzo) (int, error) {
	dbpool := pgclient.WDB
	ai := 0

	err := dbpool.QueryRow(ctx, "SELECT func_uzo_add($1, $2);", a.UzoName, a.UzoValue).Scan(&ai)

	if err != nil {
		log.Println("Failed execute func_uzo_add: ", err)
		return 0, err
	}

	return ai, nil
}

//func (est *UzoStorage) Upd(ctx context.Context, u models.Uzo) (int, error)
func (est *UzoStorage) Upd(ctx context.Context, u models.Uzo) (int, error) {
	dbpool := pgclient.WDB
	ui := 0

	err := dbpool.QueryRow(ctx, "SELECT func_uzo_upd($1,$2,$3);", u.Id, u.UzoName, u.UzoValue).Scan(&ui)

	if err != nil {
		log.Println("Failed execute func_uzo_upd: ", err)
		return 0, err
	}
	return ui, nil
}

//func (est *UzoStorage) Del(ctx context.Context, d []int) ([]int, error)
func (est *UzoStorage) Del(ctx context.Context, d []int) ([]int, error) {
	dbpool := pgclient.WDB
	res := []int{}
	i := 0
	for _, id := range d {
		err := dbpool.QueryRow(ctx, "SELECT func_uzo_del($1);", id).Scan(&i)
		res = append(res, i)

		if err != nil {
			log.Println("Failed execute func_uzo_del: ", err)
		}
	}
	return res, nil
}

//func (est *UzoStorage) GetOne(ctx context.Context, i int) (models.Uzo_count, error)
func (est *UzoStorage) GetOne(ctx context.Context, i int) (models.Uzo_count, error) {
	dbpool := pgclient.WDB
	out_arr := []models.Uzo{}
	g := models.Uzo{}

	err := dbpool.QueryRow(ctx, "SELECT * from func_uzo_getbyid($1);", i).Scan(&g.Id, &g.UzoName, &g.UzoValue)

	if err != nil && err != pgx.ErrNoRows {
		log.Println("Failed execute func_uzo_getbyid: ", err)
		return models.Uzo_count{Values: []models.Uzo{}, Count: 0, Auth: models.Auth{}}, err
	}

	out_arr = append(out_arr, g)

	out_count := models.Uzo_count{Values: out_arr, Count: 1, Auth: models.Auth{}}
	return out_count, nil
}
