package pgsql

import (
	"context"
	"log"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/u2lentaru/billbck/internal/models"
	"github.com/u2lentaru/billbck/pkg/pgclient"
)

//type GRpStorage struct
type GRpStorage struct {
	db *pgxpool.Pool
}

//func NewGRpStorage(db *pgxpool.Pool) *GRpStorage
func NewGRpStorage(db *pgxpool.Pool) *GRpStorage {
	return &GRpStorage{db: db}
}

//func (est *GRpStorage) GetList(ctx context.Context, pg, pgs int, gs1 string, ord int, dsc bool) (models.GRp_count, error)
func (est *GRpStorage) GetList(ctx context.Context, pg, pgs int, gs1 string, ord int, dsc bool) (models.GRp_count, error) {
	dbpool := pgclient.WDB
	gs := models.GRp{}

	gsc := 0
	err := dbpool.QueryRow(ctx, "SELECT * from func_grp_cnt($1);", gs1).Scan(&gsc)

	if err != nil {
		log.Println(err.Error(), "func_grp_cnt")
		return models.GRp_count{Values: []models.GRp{}, Count: gsc, Auth: models.Auth{}}, err
	}

	out_arr := make([]models.GRp, 0,
		func() int {
			if gsc < pgs {
				return gsc
			} else {
				return pgs
			}
		}())

	rows, err := dbpool.Query(ctx, "SELECT * from func_grp_get($1,$2,$3,$4,$5);", pg, pgs, gs1, ord, dsc)
	if err != nil {
		log.Println(err.Error())
		return models.GRp_count{Values: []models.GRp{}, Count: gsc, Auth: models.Auth{}}, err
	}

	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&gs.Id, &gs.GRpName)
		if err != nil {
			log.Println("failed to scan row:", err)
		}

		out_arr = append(out_arr, gs)
	}

	out_count := models.GRp_count{Values: out_arr, Count: gsc, Auth: models.Auth{}}
	return out_count, nil
}

//func (est *GRpStorage) Add(ctx context.Context, a models.GRp) (int, error)
func (est *GRpStorage) Add(ctx context.Context, a models.GRp) (int, error) {
	dbpool := pgclient.WDB
	ai := 0

	err := dbpool.QueryRow(ctx, "SELECT func_grp_add($1);", a.GRpName).Scan(&ai)

	if err != nil {
		log.Println("Failed execute func_grp_add: ", err)
		return 0, err
	}

	return ai, nil
}

//func (est *GRpStorage) Upd(ctx context.Context, u models.GRp) (int, error)
func (est *GRpStorage) Upd(ctx context.Context, u models.GRp) (int, error) {
	dbpool := pgclient.WDB
	ui := 0

	err := dbpool.QueryRow(ctx, "SELECT func_grp_upd($1,$2);", u.Id, u.GRpName).Scan(&ui)

	if err != nil {
		log.Println("Failed execute func_grp_upd: ", err)
		return 0, err
	}
	return ui, nil
}

//func (est *GRpStorage) Del(ctx context.Context, d []int) ([]int, error)
func (est *GRpStorage) Del(ctx context.Context, d []int) ([]int, error) {
	dbpool := pgclient.WDB
	res := []int{}
	i := 0
	for _, id := range d {
		err := dbpool.QueryRow(ctx, "SELECT func_grp_del($1);", id).Scan(&i)
		res = append(res, i)

		if err != nil {
			log.Println("Failed execute func_grp_del: ", err)
		}
	}
	return res, nil
}

//func (est *GRpStorage) GetOne(ctx context.Context, i int) (models.GRp_count, error)
func (est *GRpStorage) GetOne(ctx context.Context, i int) (models.GRp_count, error) {
	dbpool := pgclient.WDB
	out_arr := []models.GRp{}
	g := models.GRp{}

	err := dbpool.QueryRow(ctx, "SELECT * from func_grp_getbyid($1);", i).Scan(&g.Id, &g.GRpName)

	if err != nil && err != pgx.ErrNoRows {
		log.Println("Failed execute from func_grp_getbyid: ", err)
		return models.GRp_count{Values: []models.GRp{}, Count: 0, Auth: models.Auth{}}, err
	}

	out_arr = append(out_arr, g)

	out_count := models.GRp_count{Values: out_arr, Count: 1, Auth: models.Auth{}}
	return out_count, nil
}
