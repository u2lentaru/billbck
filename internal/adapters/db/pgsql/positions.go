package pgsql

import (
	"context"
	"log"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/u2lentaru/billbck/internal/models"
	"github.com/u2lentaru/billbck/pkg/pgclient"
)

//type PositionStorage struct
type PositionStorage struct {
	db *pgxpool.Pool
}

//func NewPositionStorage(db *pgxpool.Pool) *PositionStorage
func NewPositionStorage(db *pgxpool.Pool) *PositionStorage {
	return &PositionStorage{db: db}
}

//func (est *PositionStorage) GetList(ctx context.Context, pg, pgs int, gs1 string, ord int, dsc bool) (models.Position_count, error)
func (est *PositionStorage) GetList(ctx context.Context, pg, pgs int, gs1 string, ord int, dsc bool) (models.Position_count, error) {
	dbpool := pgclient.WDB
	gs := models.Position{}

	gsc := 0
	err := dbpool.QueryRow(ctx, "SELECT * from func_positions_cnt($1);", gs1).Scan(&gsc)

	if err != nil {
		log.Println(err.Error(), "func_positions_cnt")
		return models.Position_count{Values: []models.Position{}, Count: gsc, Auth: models.Auth{}}, err
	}

	out_arr := make([]models.Position, 0,
		func() int {
			if gsc < pgs {
				return gsc
			} else {
				return pgs
			}
		}())

	rows, err := dbpool.Query(ctx, "SELECT * from func_positions_get($1,$2,$3,$4,$5);", pg, pgs, gs1, ord, dsc)
	if err != nil {
		log.Println(err.Error(), "func_positions_get")
		return models.Position_count{Values: []models.Position{}, Count: gsc, Auth: models.Auth{}}, err
	}

	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&gs.Id, &gs.PositionName)
		if err != nil {
			log.Println("failed to scan row:", err)
		}

		out_arr = append(out_arr, gs)
	}

	out_count := models.Position_count{Values: out_arr, Count: gsc, Auth: models.Auth{}}
	if err != nil {
		log.Println(err.Error())
		return models.Position_count{}, err
	}

	return out_count, nil
}

//func (est *PositionStorage) Add(ctx context.Context, a models.Position) (int, error)
func (est *PositionStorage) Add(ctx context.Context, a models.Position) (int, error) {
	dbpool := pgclient.WDB
	ai := 0

	err := dbpool.QueryRow(ctx, "SELECT func_positions_add($1);", a.PositionName).Scan(&ai)

	if err != nil {
		log.Println("Failed execute func_positions_add: ", err)
		return 0, err
	}

	return ai, nil
}

//func (est *PositionStorage) Upd(ctx context.Context, u models.Position) (int, error)
func (est *PositionStorage) Upd(ctx context.Context, u models.Position) (int, error) {
	dbpool := pgclient.WDB
	ui := 0

	err := dbpool.QueryRow(ctx, "SELECT func_positions_upd($1,$2);", u.Id, u.PositionName).Scan(&ui)

	if err != nil {
		log.Println("Failed execute func_positions_upd: ", err)
		return 0, err
	}
	return ui, nil
}

//func (est *PositionStorage) Del(ctx context.Context, d []int) ([]int, error)
func (est *PositionStorage) Del(ctx context.Context, d []int) ([]int, error) {
	dbpool := pgclient.WDB
	res := []int{}
	i := 0
	for _, id := range d {
		err := dbpool.QueryRow(ctx, "SELECT func_positions_del($1);", id).Scan(&i)
		res = append(res, i)

		if err != nil {
			log.Println("Failed execute func_positions_del: ", err)
		}
	}
	return res, nil
}

//func (est *PositionStorage) GetOne(ctx context.Context, i int) (models.Position_count, error)
func (est *PositionStorage) GetOne(ctx context.Context, i int) (models.Position_count, error) {
	dbpool := pgclient.WDB
	out_arr := []models.Position{}
	g := models.Position{}

	err := dbpool.QueryRow(ctx, "SELECT * from func_position_get($1);", i).Scan(&g.Id, &g.PositionName)

	if err != nil && err != pgx.ErrNoRows {
		log.Println("Failed execute from func_position_get: ", err)
		return models.Position_count{Values: []models.Position{}, Count: 0, Auth: models.Auth{}}, err
	}

	out_arr = append(out_arr, g)

	out_count := models.Position_count{Values: out_arr, Count: 1, Auth: models.Auth{}}
	return out_count, nil
}
