package pgsql

import (
	"context"
	"log"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/u2lentaru/billbck/internal/models"
	"github.com/u2lentaru/billbck/pkg/pgclient"
)

//type ReliabilityStorage struct
type ReliabilityStorage struct {
	db *pgxpool.Pool
}

//func NewReliabilityStorage(db *pgxpool.Pool) *ReliabilityStorage
func NewReliabilityStorage(db *pgxpool.Pool) *ReliabilityStorage {
	return &ReliabilityStorage{db: db}
}

//func (est *ReliabilityStorage) GetList(ctx context.Context, pg, pgs int, gs1 string, ord int, dsc bool) (models.Reliability_count, error)
func (est *ReliabilityStorage) GetList(ctx context.Context, pg, pgs int, gs1 string, ord int, dsc bool) (models.Reliability_count, error) {
	dbpool := pgclient.WDB
	gs := models.Reliability{}

	gsc := 0
	err := dbpool.QueryRow(ctx, "SELECT * from func_reliabilities_cnt($1);", gs1).Scan(&gsc)

	if err != nil {
		log.Println(err.Error(), "func_reliabilities_cnt")
		return models.Reliability_count{Values: []models.Reliability{}, Count: gsc, Auth: models.Auth{}}, err
	}

	out_arr := make([]models.Reliability, 0,
		func() int {
			if gsc < pgs {
				return gsc
			} else {
				return pgs
			}
		}())

	rows, err := dbpool.Query(ctx, "SELECT * from func_reliabilities_get($1,$2,$3,$4,$5);", pg, pgs, gs1, ord, dsc)
	if err != nil {
		log.Println(err.Error(), "func_reliabilities_get")
		return models.Reliability_count{Values: []models.Reliability{}, Count: gsc, Auth: models.Auth{}}, err
	}

	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&gs.Id, &gs.ReliabilityName)
		if err != nil {
			log.Println("failed to scan row:", err)
		}

		out_arr = append(out_arr, gs)
	}

	out_count := models.Reliability_count{Values: out_arr, Count: gsc, Auth: models.Auth{}}
	return out_count, nil
}

//func (est *ReliabilityStorage) Add(ctx context.Context, a models.Reliability) (int, error)
func (est *ReliabilityStorage) Add(ctx context.Context, a models.Reliability) (int, error) {
	dbpool := pgclient.WDB
	ai := 0

	err := dbpool.QueryRow(ctx, "SELECT func_reliabilities_add($1);", a.ReliabilityName).Scan(&ai)

	if err != nil {
		log.Println("Failed execute func_reliabilities_add: ", err)
		return 0, err
	}

	return ai, nil
}

//func (est *ReliabilityStorage) Upd(ctx context.Context, u models.Reliability) (int, error)
func (est *ReliabilityStorage) Upd(ctx context.Context, u models.Reliability) (int, error) {
	dbpool := pgclient.WDB
	ui := 0

	err := dbpool.QueryRow(ctx, "SELECT func_reliabilities_upd($1,$2);", u.Id, u.ReliabilityName).Scan(&ui)

	if err != nil {
		log.Println("Failed execute func_reliabilities_upd: ", err)
		return 0, err
	}
	return ui, nil
}

//func (est *ReliabilityStorage) Del(ctx context.Context, d []int) ([]int, error)
func (est *ReliabilityStorage) Del(ctx context.Context, d []int) ([]int, error) {
	dbpool := pgclient.WDB
	res := []int{}
	i := 0
	for _, id := range d {
		err := dbpool.QueryRow(ctx, "SELECT func_reliabilities_del($1);", id).Scan(&i)
		res = append(res, i)

		if err != nil {
			log.Println("Failed execute func_reliabilities_del: ", err)
		}
	}
	return res, nil
}

//func (est *ReliabilityStorage) GetOne(ctx context.Context, i int) (models.Reliability_count, error)
func (est *ReliabilityStorage) GetOne(ctx context.Context, i int) (models.Reliability_count, error) {
	dbpool := pgclient.WDB
	out_arr := []models.Reliability{}
	g := models.Reliability{}

	err := dbpool.QueryRow(ctx, "SELECT * from func_reliability_get($1);", i).Scan(&g.Id, &g.ReliabilityName)

	if err != nil && err != pgx.ErrNoRows {
		log.Println("Failed execute func_reliability_get: ", err)
		return models.Reliability_count{Values: []models.Reliability{}, Count: 0, Auth: models.Auth{}}, err
	}

	out_arr = append(out_arr, g)

	out_count := models.Reliability_count{Values: out_arr, Count: 1, Auth: models.Auth{}}
	return out_count, nil
}
