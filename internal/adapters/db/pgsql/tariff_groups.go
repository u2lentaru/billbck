package pgsql

import (
	"context"
	"log"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/u2lentaru/billbck/internal/models"
	"github.com/u2lentaru/billbck/pkg/pgclient"
)

//type TariffGroupStorage struct
type TariffGroupStorage struct {
	db *pgxpool.Pool
}

//func NewTariffGroupStorage(db *pgxpool.Pool) *TariffGroupStorage
func NewTariffGroupStorage(db *pgxpool.Pool) *TariffGroupStorage {
	return &TariffGroupStorage{db: db}
}

//func (est *TariffGroupStorage) GetList(ctx context.Context, pg, pgs int, gs1 string, ord int, dsc bool) (models.TariffGroup_count, error)
func (est *TariffGroupStorage) GetList(ctx context.Context, pg, pgs int, gs1 string, ord int, dsc bool) (models.TariffGroup_count, error) {
	dbpool := pgclient.WDB
	gs := models.TariffGroup{}

	gsc := 0
	err := dbpool.QueryRow(ctx, "SELECT * from func_tariff_groups_cnt($1);", gs1).Scan(&gsc)

	if err != nil {
		log.Println(err.Error(), "func_tariff_groups_cnt")
		return models.TariffGroup_count{Values: []models.TariffGroup{}, Count: gsc, Auth: models.Auth{}}, err
	}

	out_arr := make([]models.TariffGroup, 0,
		func() int {
			if gsc < pgs {
				return gsc
			} else {
				return pgs
			}
		}())

	rows, err := dbpool.Query(ctx, "SELECT * from func_tariff_groups_get($1,$2,$3,$4,$5);", pg, pgs, gs1, ord, dsc)
	if err != nil {
		log.Println(err.Error(), "func_tariff_groups_get")
		return models.TariffGroup_count{Values: []models.TariffGroup{}, Count: gsc, Auth: models.Auth{}}, err
	}

	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&gs.Id, &gs.TariffGroupName)
		if err != nil {
			log.Println("failed to scan row:", err)
		}

		out_arr = append(out_arr, gs)
	}

	out_count := models.TariffGroup_count{Values: out_arr, Count: gsc, Auth: models.Auth{}}
	if err != nil {
		log.Println(err.Error())
		return models.TariffGroup_count{}, err
	}

	return out_count, nil
}

//func (est *TariffGroupStorage) Add(ctx context.Context, a models.TariffGroup) (int, error)
func (est *TariffGroupStorage) Add(ctx context.Context, a models.TariffGroup) (int, error) {
	dbpool := pgclient.WDB
	ai := 0

	err := dbpool.QueryRow(ctx, "SELECT func_tariff_groups_add($1);", a.TariffGroupName).Scan(&ai)

	if err != nil {
		log.Println("Failed execute func_tariff_groups_add: ", err)
		return 0, err
	}

	return ai, nil
}

//func (est *TariffGroupStorage) Upd(ctx context.Context, u models.TariffGroup) (int, error)
func (est *TariffGroupStorage) Upd(ctx context.Context, u models.TariffGroup) (int, error) {
	dbpool := pgclient.WDB
	ui := 0

	err := dbpool.QueryRow(ctx, "SELECT func_tariff_groups_upd($1,$2);", u.Id, u.TariffGroupName).Scan(&ui)

	if err != nil {
		log.Println("Failed execute func_tariff_groups_upd: ", err)
		return 0, err
	}
	return ui, nil
}

//func (est *TariffGroupStorage) Del(ctx context.Context, d []int) ([]int, error)
func (est *TariffGroupStorage) Del(ctx context.Context, d []int) ([]int, error) {
	dbpool := pgclient.WDB
	res := []int{}
	i := 0
	for _, id := range d {
		err := dbpool.QueryRow(ctx, "SELECT func_tariff_groups_del($1);", id).Scan(&i)
		res = append(res, i)

		if err != nil {
			log.Println("Failed execute func_tariff_groups_del: ", err)
		}
	}
	return res, nil
}

//func (est *TariffGroupStorage) GetOne(ctx context.Context, i int) (models.TariffGroup_count, error)
func (est *TariffGroupStorage) GetOne(ctx context.Context, i int) (models.TariffGroup_count, error) {
	dbpool := pgclient.WDB
	out_arr := []models.TariffGroup{}
	g := models.TariffGroup{}

	err := dbpool.QueryRow(ctx, "SELECT * from func_tariff_group_get($1);", i).Scan(&g.Id, &g.TariffGroupName)

	if err != nil && err != pgx.ErrNoRows {
		log.Println("Failed execute func_tariff_group_get: ", err)
		return models.TariffGroup_count{Values: []models.TariffGroup{}, Count: 0, Auth: models.Auth{}}, err
	}

	out_arr = append(out_arr, g)

	out_count := models.TariffGroup_count{Values: out_arr, Count: 1, Auth: models.Auth{}}
	return out_count, nil
}
