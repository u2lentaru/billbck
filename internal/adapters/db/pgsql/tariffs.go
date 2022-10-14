package pgsql

import (
	"context"
	"log"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/u2lentaru/billbck/internal/models"
	"github.com/u2lentaru/billbck/pkg/pgclient"
)

//type TariffStorage struct
type TariffStorage struct {
	db *pgxpool.Pool
}

//func NewTariffStorage(db *pgxpool.Pool) *TariffStorage
func NewTariffStorage(db *pgxpool.Pool) *TariffStorage {
	return &TariffStorage{db: db}
}

//func (est *TariffStorage) GetList(ctx context.Context, pg, pgs int, gs1 string, ord int, dsc bool) (models.Tariff_count, error)
func (est *TariffStorage) GetList(ctx context.Context, pg, pgs int, gs1 string, ord int, dsc bool) (models.Tariff_count, error) {
	dbpool := pgclient.WDB
	gs := models.Tariff{}

	gsc := 0
	err := dbpool.QueryRow(ctx, "SELECT * from func_tariffs_cnt($1);", gs1).Scan(&gsc)

	if err != nil {
		log.Println(err.Error(), "func_tariffs_cnt")
		return models.Tariff_count{Values: []models.Tariff{}, Count: gsc, Auth: models.Auth{}}, err
	}

	out_arr := make([]models.Tariff, 0,
		func() int {
			if gsc < pgs {
				return gsc
			} else {
				return pgs
			}
		}())

	rows, err := dbpool.Query(ctx, "SELECT * from func_tariffs_get($1,$2,$3,$4,$5);", pg, pgs, gs1, ord, dsc)
	if err != nil {
		log.Println(err.Error(), "func_tariffs_get")
		return models.Tariff_count{Values: []models.Tariff{}, Count: gsc, Auth: models.Auth{}}, err
	}

	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&gs.Id, &gs.TariffName, &gs.TariffGroup.Id, &gs.Norma, &gs.Tariff, &gs.Startdate, &gs.Enddate, &gs.TariffGroup.TariffGroupName)
		if err != nil {
			log.Println("failed to scan row:", err)
		}

		out_arr = append(out_arr, gs)
	}

	out_count := models.Tariff_count{Values: out_arr, Count: gsc, Auth: models.Auth{}}
	if err != nil {
		log.Println(err.Error())
		return models.Tariff_count{}, err
	}

	return out_count, nil
}

//func (est *TariffStorage) Add(ctx context.Context, a models.Tariff) (int, error)
func (est *TariffStorage) Add(ctx context.Context, a models.Tariff) (int, error) {
	dbpool := pgclient.WDB
	ai := 0

	err := dbpool.QueryRow(ctx, "SELECT func_tariffs_add($1,$2,$3,$4,$5);", a.TariffName, a.TariffGroup.Id, a.Norma, a.Tariff, a.Startdate).Scan(&ai)

	if err != nil {
		log.Println("Failed execute func_tariffs_add: ", err)
		return 0, err
	}

	return ai, nil
}

//func (est *TariffStorage) Upd(ctx context.Context, u models.Tariff) (int, error)
func (est *TariffStorage) Upd(ctx context.Context, u models.Tariff) (int, error) {
	dbpool := pgclient.WDB
	ui := 0

	err := dbpool.QueryRow(ctx, "SELECT func_tariffs_upd($1,$2,$3,$4,$5,$6,$7);", u.Id, u.TariffName, u.TariffGroup.Id, u.Norma, u.Tariff,
		u.Startdate, u.Enddate).Scan(&ui)

	if err != nil {
		log.Println("Failed execute func_tariffs_upd: ", err)
		return 0, err
	}
	return ui, nil
}

//func (est *TariffStorage) Del(ctx context.Context, d []int) ([]int, error)
func (est *TariffStorage) Del(ctx context.Context, d []int) ([]int, error) {
	dbpool := pgclient.WDB
	res := []int{}
	i := 0
	for _, id := range d {
		err := dbpool.QueryRow(ctx, "SELECT func_tariffs_del($1);", id).Scan(&i)
		res = append(res, i)

		if err != nil {
			log.Println("Failed execute func_tariffs_del: ", err)
		}
	}
	return res, nil
}

//func (est *TariffStorage) GetOne(ctx context.Context, i int) (models.Tariff_count, error)
func (est *TariffStorage) GetOne(ctx context.Context, i int) (models.Tariff_count, error) {
	dbpool := pgclient.WDB
	out_arr := []models.Tariff{}
	g := models.Tariff{}

	err := dbpool.QueryRow(ctx, "SELECT * from func_tariff_get($1);", i).Scan(&g.Id, &g.TariffName, &g.TariffGroup.Id, &g.Norma, &g.Tariff,
		&g.Startdate, &g.Enddate, &g.TariffGroup.TariffGroupName)

	if err != nil && err != pgx.ErrNoRows {
		log.Println("Failed execute func_tariff_get: ", err)
		return models.Tariff_count{Values: []models.Tariff{}, Count: 0, Auth: models.Auth{}}, err
	}

	out_arr = append(out_arr, g)

	out_count := models.Tariff_count{Values: out_arr, Count: 1, Auth: models.Auth{}}
	return out_count, nil
}
