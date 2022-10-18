package pgsql

import (
	"context"
	"log"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/u2lentaru/billbck/internal/models"
	"github.com/u2lentaru/billbck/pkg/pgclient"
)

//type SealColourStorage struct
type SealColourStorage struct {
	db *pgxpool.Pool
}

//func NewSealColourStorage(db *pgxpool.Pool) *SealColourStorage
func NewSealColourStorage(db *pgxpool.Pool) *SealColourStorage {
	return &SealColourStorage{db: db}
}

//func (est *SealColourStorage) GetList(ctx context.Context, pg, pgs int, gs1 string, ord int, dsc bool) (models.SealColour_count, error)
func (est *SealColourStorage) GetList(ctx context.Context, pg, pgs int, gs1 string, ord int, dsc bool) (models.SealColour_count, error) {
	dbpool := pgclient.WDB
	gs := models.SealColour{}

	gsc := 0
	err := dbpool.QueryRow(ctx, "SELECT * from func_seal_colours_cnt($1);", gs1).Scan(&gsc)

	if err != nil {
		log.Println(err.Error(), "func_seal_colours_cnt")
		return models.SealColour_count{Values: []models.SealColour{}, Count: gsc, Auth: models.Auth{}}, err
	}

	out_arr := make([]models.SealColour, 0,
		func() int {
			if gsc < pgs {
				return gsc
			} else {
				return pgs
			}
		}())

	rows, err := dbpool.Query(ctx, "SELECT * from func_seal_colours_get($1,$2,$3,$4,$5);", pg, pgs, gs1, ord, dsc)
	if err != nil {
		log.Println(err.Error(), "func_seal_colours_get")
		return models.SealColour_count{Values: []models.SealColour{}, Count: gsc, Auth: models.Auth{}}, err
	}

	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&gs.Id, &gs.SealColourName)
		if err != nil {
			log.Println("failed to scan row:", err)
		}

		out_arr = append(out_arr, gs)
	}

	out_count := models.SealColour_count{Values: out_arr, Count: gsc, Auth: models.Auth{}}
	return out_count, nil
}

//func (est *SealColourStorage) Add(ctx context.Context, a models.SealColour) (int, error)
func (est *SealColourStorage) Add(ctx context.Context, a models.SealColour) (int, error) {
	dbpool := pgclient.WDB
	ai := 0

	err := dbpool.QueryRow(ctx, "SELECT func_seal_colours_add($1);", a.SealColourName).Scan(&ai)

	if err != nil {
		log.Println("Failed execute func_seal_colours_add: ", err)
		return 0, err
	}

	return ai, nil
}

//func (est *SealColourStorage) Upd(ctx context.Context, u models.SealColour) (int, error)
func (est *SealColourStorage) Upd(ctx context.Context, u models.SealColour) (int, error) {
	dbpool := pgclient.WDB
	ui := 0

	err := dbpool.QueryRow(ctx, "SELECT func_seal_colours_upd($1,$2);", u.Id, u.SealColourName).Scan(&ui)

	if err != nil {
		log.Println("Failed execute func_seal_colours_upd: ", err)
		return 0, err
	}
	return ui, nil
}

//func (est *SealColourStorage) Del(ctx context.Context, d []int) ([]int, error)
func (est *SealColourStorage) Del(ctx context.Context, d []int) ([]int, error) {
	dbpool := pgclient.WDB
	res := []int{}
	i := 0
	for _, id := range d {
		err := dbpool.QueryRow(ctx, "SELECT func_seal_colours_del($1);", id).Scan(&i)
		res = append(res, i)

		if err != nil {
			log.Println("Failed execute func_seal_colours_del: ", err)
		}
	}
	return res, nil
}

//func (est *SealColourStorage) GetOne(ctx context.Context, i int) (models.SealColour_count, error)
func (est *SealColourStorage) GetOne(ctx context.Context, i int) (models.SealColour_count, error) {
	dbpool := pgclient.WDB
	out_arr := []models.SealColour{}
	g := models.SealColour{}

	err := dbpool.QueryRow(ctx, "SELECT * from func_seal_colour_get($1);", i).Scan(&g.Id, &g.SealColourName)

	if err != nil && err != pgx.ErrNoRows {
		log.Println("Failed execute func_seal_colour_get: ", err)
		return models.SealColour_count{Values: []models.SealColour{}, Count: 0, Auth: models.Auth{}}, err
	}

	out_arr = append(out_arr, g)

	out_count := models.SealColour_count{Values: out_arr, Count: 1, Auth: models.Auth{}}
	return out_count, nil
}
