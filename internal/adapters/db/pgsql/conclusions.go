package pgsql

import (
	"context"
	"log"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/u2lentaru/billbck/internal/models"
	"github.com/u2lentaru/billbck/pkg/pgclient"
)

//type ConclusionStorage struct
type ConclusionStorage struct {
	db *pgxpool.Pool
}

//func NewConclusionStorage(db *pgxpool.Pool) *ConclusionStorage
func NewConclusionStorage(db *pgxpool.Pool) *ConclusionStorage {
	return &ConclusionStorage{db: db}
}

//func (est *ConclusionStorage) GetList(ctx context.Context, pg, pgs int, gs1 string, ord int, dsc bool) (models.Conclusion_count, error)
func (est *ConclusionStorage) GetList(ctx context.Context, pg, pgs int, gs1 string, ord int, dsc bool) (models.Conclusion_count, error) {
	dbpool := pgclient.WDB
	gs := models.Conclusion{}

	gsc := 0
	err := dbpool.QueryRow(ctx, "SELECT * from func_conclusions_cnt($1);", gs1).Scan(&gsc)

	if err != nil {
		log.Println(err.Error(), "func_conclusions_cnt")
		return models.Conclusion_count{Values: []models.Conclusion{}, Count: gsc, Auth: models.Auth{}}, err
	}

	out_arr := make([]models.Conclusion, 0,
		func() int {
			if gsc < pgs {
				return gsc
			} else {
				return pgs
			}
		}())

	rows, err := dbpool.Query(ctx, "SELECT * from func_conclusions_get($1,$2,$3,$4,$5);", pg, pgs, gs1, ord, dsc)
	if err != nil {
		log.Println(err.Error())
		return models.Conclusion_count{Values: []models.Conclusion{}, Count: gsc, Auth: models.Auth{}}, err
	}

	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&gs.Id, &gs.ConclusionName)
		if err != nil {
			log.Println("failed to scan row:", err)
		}

		out_arr = append(out_arr, gs)
	}

	out_count := models.Conclusion_count{Values: out_arr, Count: gsc, Auth: models.Auth{}}
	if err != nil {
		log.Println(err.Error())
		return models.Conclusion_count{}, err
	}

	return out_count, nil
}

//func (est *ConclusionStorage) Add(ctx context.Context, a models.Conclusion) (int, error)
func (est *ConclusionStorage) Add(ctx context.Context, a models.Conclusion) (int, error) {
	dbpool := pgclient.WDB
	ai := 0

	err := dbpool.QueryRow(ctx, "SELECT func_conclusions_add($1);", a.ConclusionName).Scan(&ai)

	if err != nil {
		log.Println("Failed execute func_conclusions_add: ", err)
		return 0, err
	}

	return ai, nil
}

//func (est *ConclusionStorage) Upd(ctx context.Context, u models.Conclusion) (int, error)
func (est *ConclusionStorage) Upd(ctx context.Context, u models.Conclusion) (int, error) {
	dbpool := pgclient.WDB
	ui := 0

	err := dbpool.QueryRow(ctx, "SELECT func_conclusions_upd($1,$2);", u.Id, u.ConclusionName).Scan(&ui)

	if err != nil {
		log.Println("Failed execute func_conclusions_upd: ", err)
		return 0, err
	}
	return ui, nil
}

//func (est *ConclusionStorage) Del(ctx context.Context, d []int) ([]int, error)
func (est *ConclusionStorage) Del(ctx context.Context, d []int) ([]int, error) {
	dbpool := pgclient.WDB
	res := []int{}
	i := 0
	for _, id := range d {
		err := dbpool.QueryRow(ctx, "SELECT func_conclusions_del($1);", id).Scan(&i)
		res = append(res, i)

		if err != nil {
			log.Println("Failed execute func_conclusions_del: ", err)
		}
	}
	return res, nil
}

//func (est *ConclusionStorage) GetOne(ctx context.Context, i int) (models.Conclusion_count, error)
func (est *ConclusionStorage) GetOne(ctx context.Context, i int) (models.Conclusion_count, error) {
	dbpool := pgclient.WDB
	out_arr := []models.Conclusion{}
	g := models.Conclusion{}

	err := dbpool.QueryRow(ctx, "SELECT * from func_conclusion_get($1);", i).Scan(&g.Id, &g.ConclusionName)

	if err != nil && err != pgx.ErrNoRows {
		log.Println("Failed execute from func_conclusion_get: ", err)
		return models.Conclusion_count{Values: []models.Conclusion{}, Count: 0, Auth: models.Auth{}}, err
	}

	out_arr = append(out_arr, g)

	out_count := models.Conclusion_count{Values: out_arr, Count: 1, Auth: models.Auth{}}
	return out_count, nil
}
