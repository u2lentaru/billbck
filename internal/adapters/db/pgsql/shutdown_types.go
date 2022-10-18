package pgsql

import (
	"context"
	"log"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/u2lentaru/billbck/internal/models"
	"github.com/u2lentaru/billbck/pkg/pgclient"
)

//type ShutdownTypeStorage struct
type ShutdownTypeStorage struct {
	db *pgxpool.Pool
}

//func NewShutdownTypeStorage(db *pgxpool.Pool) *ShutdownTypeStorage
func NewShutdownTypeStorage(db *pgxpool.Pool) *ShutdownTypeStorage {
	return &ShutdownTypeStorage{db: db}
}

//func (est *ShutdownTypeStorage) GetList(ctx context.Context, pg, pgs int, gs1 string, ord int, dsc bool) (models.ShutdownType_count, error)
func (est *ShutdownTypeStorage) GetList(ctx context.Context, pg, pgs int, gs1 string, ord int, dsc bool) (models.ShutdownType_count, error) {
	dbpool := pgclient.WDB
	gs := models.ShutdownType{}

	gsc := 0
	err := dbpool.QueryRow(ctx, "SELECT * from func_shutdown_types_cnt($1);", gs1).Scan(&gsc)

	if err != nil {
		log.Println(err.Error(), "func_shutdown_types_cnt")
		return models.ShutdownType_count{Values: []models.ShutdownType{}, Count: gsc, Auth: models.Auth{}}, err
	}

	out_arr := make([]models.ShutdownType, 0,
		func() int {
			if gsc < pgs {
				return gsc
			} else {
				return pgs
			}
		}())

	rows, err := dbpool.Query(ctx, "SELECT * from func_shutdown_types_get($1,$2,$3,$4,$5);", pg, pgs, gs1, ord, dsc)
	if err != nil {
		log.Println(err.Error(), "func_shutdown_types_get")
		return models.ShutdownType_count{Values: []models.ShutdownType{}, Count: gsc, Auth: models.Auth{}}, err
	}

	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&gs.Id, &gs.ShutdownTypeName)
		if err != nil {
			log.Println("failed to scan row:", err)
		}

		out_arr = append(out_arr, gs)
	}

	out_count := models.ShutdownType_count{Values: out_arr, Count: gsc, Auth: models.Auth{}}
	return out_count, nil
}

//func (est *ShutdownTypeStorage) Add(ctx context.Context, a models.ShutdownType) (int, error)
func (est *ShutdownTypeStorage) Add(ctx context.Context, a models.ShutdownType) (int, error) {
	dbpool := pgclient.WDB
	ai := 0

	err := dbpool.QueryRow(ctx, "SELECT func_shutdown_types_add($1);", a.ShutdownTypeName).Scan(&ai)

	if err != nil {
		log.Println("Failed execute func_shutdown_types_add: ", err)
		return 0, err
	}

	return ai, nil
}

//func (est *ShutdownTypeStorage) Upd(ctx context.Context, u models.ShutdownType) (int, error)
func (est *ShutdownTypeStorage) Upd(ctx context.Context, u models.ShutdownType) (int, error) {
	dbpool := pgclient.WDB
	ui := 0

	err := dbpool.QueryRow(ctx, "SELECT func_shutdown_types_upd($1,$2);", u.Id, u.ShutdownTypeName).Scan(&ui)

	if err != nil {
		log.Println("Failed execute func_shutdown_types_upd: ", err)
		return 0, err
	}
	return ui, nil
}

//func (est *ShutdownTypeStorage) Del(ctx context.Context, d []int) ([]int, error)
func (est *ShutdownTypeStorage) Del(ctx context.Context, d []int) ([]int, error) {
	dbpool := pgclient.WDB
	res := []int{}
	i := 0
	for _, id := range d {
		err := dbpool.QueryRow(ctx, "SELECT func_shutdown_types_del($1);", id).Scan(&i)
		res = append(res, i)

		if err != nil {
			log.Println("Failed execute func_shutdown_types_del: ", err)
		}
	}
	return res, nil
}

//func (est *ShutdownTypeStorage) GetOne(ctx context.Context, i int) (models.ShutdownType_count, error)
func (est *ShutdownTypeStorage) GetOne(ctx context.Context, i int) (models.ShutdownType_count, error) {
	dbpool := pgclient.WDB
	out_arr := []models.ShutdownType{}
	g := models.ShutdownType{}

	err := dbpool.QueryRow(ctx, "SELECT * from func_shutdown_type_get($1);", i).Scan(&g.Id, &g.ShutdownTypeName)

	if err != nil && err != pgx.ErrNoRows {
		log.Println("Failed execute func_shutdown_type_get: ", err)
		return models.ShutdownType_count{Values: []models.ShutdownType{}, Count: 0, Auth: models.Auth{}}, err
	}

	out_arr = append(out_arr, g)

	out_count := models.ShutdownType_count{Values: out_arr, Count: 1, Auth: models.Auth{}}
	return out_count, nil
}
