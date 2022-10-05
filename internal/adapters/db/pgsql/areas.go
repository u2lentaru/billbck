package pgsql

import (
	"context"
	"log"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/u2lentaru/billbck/internal/models"
	"github.com/u2lentaru/billbck/pkg/pgclient"
)

//type AreaStorage struct
type AreaStorage struct {
	db *pgxpool.Pool
}

//func NewAreaStorage(db *pgxpool.Pool) *AreaStorage
func NewAreaStorage(db *pgxpool.Pool) *AreaStorage {
	return &AreaStorage{db: db}
}

//func (est *AreaStorage) GetList(ctx context.Context, pg, pgs int, gs1, gs2 string, ord int, dsc bool) (models.Area_count, error) {
func (est *AreaStorage) GetList(ctx context.Context, pg, pgs int, gs1, gs2 string, ord int, dsc bool) (models.Area_count, error) {
	dbpool := pgclient.WDB
	gs := models.Area{}

	gsc := 0
	err := dbpool.QueryRow(ctx, "SELECT * from func_areas_cnt($1,$2);", gs1, gs2).Scan(&gsc)

	if err != nil {
		log.Println(err.Error(), "func_areas_cnt")
		return models.Area_count{Values: []models.Area{}, Count: gsc, Auth: models.Auth{}}, err
	}

	out_arr := make([]models.Area, 0,
		func() int {
			if gsc < pgs {
				return gsc
			} else {
				return pgs
			}
		}())

	rows, err := dbpool.Query(ctx, "SELECT * from func_areas_get($1,$2,$3,$4,$5,$6);", pg, pgs, gs1, gs2, ord, dsc)
	if err != nil {
		log.Println(err.Error(), "func_areas_get")
		return models.Area_count{Values: []models.Area{}, Count: gsc, Auth: models.Auth{}}, err
	}

	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&gs.Id, &gs.AreaNumber, &gs.AreaName)
		if err != nil {
			log.Println("failed to scan row:", err)
		}

		out_arr = append(out_arr, gs)
	}

	out_count := models.Area_count{Values: out_arr, Count: gsc, Auth: models.Auth{}}
	if err != nil {
		log.Println(err.Error())
		return models.Area_count{}, err
	}

	return out_count, nil
}

//func (est *AreaStorage) Add(ctx context.Context, ea models.Area) (int, error)
func (est *AreaStorage) Add(ctx context.Context, a models.Area) (int, error) {
	dbpool := pgclient.WDB
	ai := 0

	err := dbpool.QueryRow(context.Background(), "SELECT func_areas_add($1,$2);", a.AreaNumber, a.AreaName).Scan(&ai)

	if err != nil {
		log.Println("Failed execute func_areas_add: ", err)
		return 0, err
	}

	return ai, nil
}

//func (est *AreaStorage) Upd(ctx context.Context, eu models.Area) (int, error)
func (est *AreaStorage) Upd(ctx context.Context, u models.Area) (int, error) {
	dbpool := pgclient.WDB
	ui := 0

	err := dbpool.QueryRow(context.Background(), "SELECT func_areas_upd($1,$2,$3);", u.Id, u.AreaNumber, u.AreaName).Scan(&ui)

	if err != nil {
		log.Println("Failed execute func_areas_upd: ", err)
		return 0, err
	}
	return ui, nil
}

//func (est *AreaStorage) Del(ctx context.Context, ed []int) ([]int, error)
func (est *AreaStorage) Del(ctx context.Context, d []int) ([]int, error) {
	dbpool := pgclient.WDB
	res := []int{}
	i := 0
	for _, id := range d {
		err := dbpool.QueryRow(context.Background(), "SELECT func_areas_del($1);", id).Scan(&i)
		res = append(res, i)

		if err != nil {
			log.Println("Failed execute func_areas_del: ", err)
		}
	}
	return res, nil
}

//func (est *AreaStorage) GetOne(ctx context.Context, i int) (models.Area_count, error)
func (est *AreaStorage) GetOne(ctx context.Context, i int) (models.Area_count, error) {
	dbpool := pgclient.WDB
	out_arr := []models.Area{}
	g := models.Area{}

	err := dbpool.QueryRow(context.Background(), "SELECT * from func_area_get($1);", i).Scan(&g.Id, &g.AreaNumber, &g.AreaName)

	if err != nil && err != pgx.ErrNoRows {
		log.Println("Failed execute from func_area_get: ", err)
		return models.Area_count{Values: []models.Area{}, Count: 0, Auth: models.Auth{}}, err
	}

	out_arr = append(out_arr, g)

	out_count := models.Area_count{Values: out_arr, Count: 1, Auth: models.Auth{}}
	return out_count, nil
}
