package pgsql

import (
	"context"
	"log"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/u2lentaru/billbck/internal/models"
	"github.com/u2lentaru/billbck/pkg/pgclient"
)

//type CityStorage struct
type CityStorage struct {
	db *pgxpool.Pool
}

//func NewCityStorage(db *pgxpool.Pool) *CityStorage
func NewCityStorage(db *pgxpool.Pool) *CityStorage {
	return &CityStorage{db: db}
}

//func (est *CityStorage) GetList(ctx context.Context, pg, pgs int, gs1 string, ord int, dsc bool) (models.City_count, error)
func (est *CityStorage) GetList(ctx context.Context, pg, pgs int, gs1 string, ord int, dsc bool) (models.City_count, error) {
	dbpool := pgclient.WDB
	gs := models.City{}

	gsc := 0
	err := dbpool.QueryRow(ctx, "SELECT * from func_cities_cnt($1);", gs1).Scan(&gsc)

	if err != nil {
		log.Println(err.Error(), "func_cities_cnt")
		return models.City_count{Values: []models.City{}, Count: gsc, Auth: models.Auth{}}, err
	}

	out_arr := make([]models.City, 0,
		func() int {
			if gsc < pgs {
				return gsc
			} else {
				return pgs
			}
		}())

	rows, err := dbpool.Query(ctx, "SELECT * from func_cities_get($1,$2,$3,$4,$5);", pg, pgs, gs1, ord, dsc)
	if err != nil {
		log.Println(err.Error())
		return models.City_count{Values: []models.City{}, Count: gsc, Auth: models.Auth{}}, err
	}

	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&gs.Id, &gs.CityName)
		if err != nil {
			log.Println("failed to scan row:", err)
		}

		out_arr = append(out_arr, gs)
	}

	out_count := models.City_count{Values: out_arr, Count: gsc, Auth: models.Auth{}}
	return out_count, nil
}

//func (est *CityStorage) Add(ctx context.Context, a models.City) (int, error)
func (est *CityStorage) Add(ctx context.Context, a models.City) (int, error) {
	dbpool := pgclient.WDB
	ai := 0

	err := dbpool.QueryRow(ctx, "SELECT func_cities_add($1);", a.CityName).Scan(&ai)

	if err != nil {
		log.Println("Failed execute func_cities_add: ", err)
		return 0, err
	}

	return ai, nil
}

//func (est *CityStorage) Upd(ctx context.Context, u models.City) (int, error)
func (est *CityStorage) Upd(ctx context.Context, u models.City) (int, error) {
	dbpool := pgclient.WDB
	ui := 0

	err := dbpool.QueryRow(ctx, "SELECT func_cities_upd($1,$2);", u.Id, u.CityName).Scan(&ui)

	if err != nil {
		log.Println("Failed execute func_cities_upd: ", err)
		return 0, err
	}
	return ui, nil
}

//func (est *CityStorage) Del(ctx context.Context, d []int) ([]int, error)
func (est *CityStorage) Del(ctx context.Context, d []int) ([]int, error) {
	dbpool := pgclient.WDB
	res := []int{}
	i := 0
	for _, id := range d {
		err := dbpool.QueryRow(ctx, "SELECT func_cities_del($1);", id).Scan(&i)
		res = append(res, i)

		if err != nil {
			log.Println("Failed execute func_cities_del: ", err)
		}
	}
	return res, nil
}

//func (est *CityStorage) GetOne(ctx context.Context, i int) (models.City_count, error)
func (est *CityStorage) GetOne(ctx context.Context, i int) (models.City_count, error) {
	dbpool := pgclient.WDB
	out_arr := []models.City{}
	g := models.City{}

	err := dbpool.QueryRow(ctx, "SELECT * from func_city_get($1);", i).Scan(&g.Id, &g.CityName)

	if err != nil && err != pgx.ErrNoRows {
		log.Println("Failed execute from func_city_get: ", err)
		return models.City_count{Values: []models.City{}, Count: 0, Auth: models.Auth{}}, err
	}

	out_arr = append(out_arr, g)

	out_count := models.City_count{Values: out_arr, Count: 1, Auth: models.Auth{}}
	return out_count, nil
}
