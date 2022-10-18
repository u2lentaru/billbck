package pgsql

import (
	"context"
	"log"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/u2lentaru/billbck/internal/models"
	"github.com/u2lentaru/billbck/pkg/pgclient"
)

//type CashdeskStorage struct
type CashdeskStorage struct {
	db *pgxpool.Pool
}

//func NewCashdeskStorage(db *pgxpool.Pool) *CashdeskStorage
func NewCashdeskStorage(db *pgxpool.Pool) *CashdeskStorage {
	return &CashdeskStorage{db: db}
}

//func (est *CashdeskStorage) GetList(ctx context.Context, pg, pgs int, gs1 string, ord int, dsc bool) (models.Cashdesk_count, error)
func (est *CashdeskStorage) GetList(ctx context.Context, pg, pgs int, gs1 string, ord int, dsc bool) (models.Cashdesk_count, error) {
	dbpool := pgclient.WDB
	gs := models.Cashdesk{}

	gsc := 0
	err := dbpool.QueryRow(ctx, "SELECT * from func_cashdesks_cnt($1);", gs1).Scan(&gsc)

	if err != nil {
		log.Println(err.Error(), "func_cashdesks_cnt")
		return models.Cashdesk_count{Values: []models.Cashdesk{}, Count: gsc, Auth: models.Auth{}}, err
	}

	out_arr := make([]models.Cashdesk, 0,
		func() int {
			if gsc < pgs {
				return gsc
			} else {
				return pgs
			}
		}())

	rows, err := dbpool.Query(ctx, "SELECT * from func_cashdesks_get($1,$2,$3,$4,$5);", pg, pgs, gs1, ord, dsc)
	if err != nil {
		log.Println(err.Error())
		return models.Cashdesk_count{Values: []models.Cashdesk{}, Count: gsc, Auth: models.Auth{}}, err
	}

	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&gs.Id, &gs.CashdeskName)
		if err != nil {
			log.Println("failed to scan row:", err)
		}

		out_arr = append(out_arr, gs)
	}

	out_count := models.Cashdesk_count{Values: out_arr, Count: gsc, Auth: models.Auth{}}
	return out_count, nil
}

//func (est *CashdeskStorage) Add(ctx context.Context, a models.Cashdesk) (int, error)
func (est *CashdeskStorage) Add(ctx context.Context, a models.Cashdesk) (int, error) {
	dbpool := pgclient.WDB
	ai := 0

	err := dbpool.QueryRow(ctx, "SELECT func_cashdesks_add($1);", a.CashdeskName).Scan(&ai)

	if err != nil {
		log.Println("Failed execute func_cashdesks_add: ", err)
		return 0, err
	}

	return ai, nil
}

//func (est *CashdeskStorage) Upd(ctx context.Context, u models.Cashdesk) (int, error)
func (est *CashdeskStorage) Upd(ctx context.Context, u models.Cashdesk) (int, error) {
	dbpool := pgclient.WDB
	ui := 0

	err := dbpool.QueryRow(ctx, "SELECT func_cashdesks_upd($1,$2);", u.Id, u.CashdeskName).Scan(&ui)

	if err != nil {
		log.Println("Failed execute func_cashdesks_upd: ", err)
		return 0, err
	}
	return ui, nil
}

//func (est *CashdeskStorage) Del(ctx context.Context, d []int) ([]int, error)
func (est *CashdeskStorage) Del(ctx context.Context, d []int) ([]int, error) {
	dbpool := pgclient.WDB
	res := []int{}
	i := 0
	for _, id := range d {
		err := dbpool.QueryRow(ctx, "SELECT func_cashdesks_del($1);", id).Scan(&i)
		res = append(res, i)

		if err != nil {
			log.Println("Failed execute func_cashdesks_del: ", err)
		}
	}
	return res, nil
}

//func (est *CashdeskStorage) GetOne(ctx context.Context, i int) (models.Cashdesk_count, error)
func (est *CashdeskStorage) GetOne(ctx context.Context, i int) (models.Cashdesk_count, error) {
	dbpool := pgclient.WDB
	out_arr := []models.Cashdesk{}
	g := models.Cashdesk{}

	err := dbpool.QueryRow(ctx, "SELECT * from func_cashdesk_get($1);", i).Scan(&g.Id, &g.CashdeskName)

	if err != nil && err != pgx.ErrNoRows {
		log.Println("Failed execute from func_cashdesk_get: ", err)
		return models.Cashdesk_count{Values: []models.Cashdesk{}, Count: 0, Auth: models.Auth{}}, err
	}

	out_arr = append(out_arr, g)

	out_count := models.Cashdesk_count{Values: out_arr, Count: 1, Auth: models.Auth{}}
	return out_count, nil
}
