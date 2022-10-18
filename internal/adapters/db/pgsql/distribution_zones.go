package pgsql

import (
	"context"
	"log"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/u2lentaru/billbck/internal/models"
	"github.com/u2lentaru/billbck/pkg/pgclient"
)

//type DistributionZoneStorage struct
type DistributionZoneStorage struct {
	db *pgxpool.Pool
}

//func NewDistributionZoneStorage(db *pgxpool.Pool) *DistributionZoneStorage
func NewDistributionZoneStorage(db *pgxpool.Pool) *DistributionZoneStorage {
	return &DistributionZoneStorage{db: db}
}

//func (est *DistributionZoneStorage) GetList(ctx context.Context, pg, pgs int, gs1 string, ord int, dsc bool) (models.DistributionZone_count, error)
func (est *DistributionZoneStorage) GetList(ctx context.Context, pg, pgs int, gs1 string, ord int, dsc bool) (models.DistributionZone_count, error) {
	dbpool := pgclient.WDB
	gs := models.DistributionZone{}

	gsc := 0
	err := dbpool.QueryRow(ctx, "SELECT * from func_distribution_zones_cnt($1);", gs1).Scan(&gsc)

	if err != nil {
		log.Println(err.Error(), "func_distribution_zones_cnt")
		return models.DistributionZone_count{Values: []models.DistributionZone{}, Count: gsc, Auth: models.Auth{}}, err
	}

	out_arr := make([]models.DistributionZone, 0,
		func() int {
			if gsc < pgs {
				return gsc
			} else {
				return pgs
			}
		}())

	rows, err := dbpool.Query(ctx, "SELECT * from func_distribution_zones_get($1,$2,$3,$4,$5);", pg, pgs, gs1, ord, dsc)
	if err != nil {
		log.Println(err.Error())
		return models.DistributionZone_count{Values: []models.DistributionZone{}, Count: gsc, Auth: models.Auth{}}, err
	}

	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&gs.Id, &gs.DistributionZoneName)
		if err != nil {
			log.Println("failed to scan row:", err)
		}

		out_arr = append(out_arr, gs)
	}

	out_count := models.DistributionZone_count{Values: out_arr, Count: gsc, Auth: models.Auth{}}
	if err != nil {
		log.Println(err.Error())
		return models.DistributionZone_count{}, err
	}

	return out_count, nil
}

//func (est *DistributionZoneStorage) Add(ctx context.Context, a models.DistributionZone) (int, error)
func (est *DistributionZoneStorage) Add(ctx context.Context, a models.DistributionZone) (int, error) {
	dbpool := pgclient.WDB
	ai := 0

	err := dbpool.QueryRow(ctx, "SELECT func_distribution_zones_add($1);", a.DistributionZoneName).Scan(&ai)

	if err != nil {
		log.Println("Failed execute func_distribution_zones_add: ", err)
		return 0, err
	}

	return ai, nil
}

//func (est *DistributionZoneStorage) Upd(ctx context.Context, u models.DistributionZone) (int, error)
func (est *DistributionZoneStorage) Upd(ctx context.Context, u models.DistributionZone) (int, error) {
	dbpool := pgclient.WDB
	ui := 0

	err := dbpool.QueryRow(ctx, "SELECT func_distribution_zones_upd($1,$2);", u.Id, u.DistributionZoneName).Scan(&ui)

	if err != nil {
		log.Println("Failed execute func_distribution_zones_upd: ", err)
		return 0, err
	}
	return ui, nil
}

//func (est *DistributionZoneStorage) Del(ctx context.Context, d []int) ([]int, error)
func (est *DistributionZoneStorage) Del(ctx context.Context, d []int) ([]int, error) {
	dbpool := pgclient.WDB
	res := []int{}
	i := 0
	for _, id := range d {
		err := dbpool.QueryRow(ctx, "SELECT func_distribution_zones_del($1);", id).Scan(&i)
		res = append(res, i)

		if err != nil {
			log.Println("Failed execute func_distribution_zones_del: ", err)
		}
	}
	return res, nil
}

//func (est *DistributionZoneStorage) GetOne(ctx context.Context, i int) (models.DistributionZone_count, error)
func (est *DistributionZoneStorage) GetOne(ctx context.Context, i int) (models.DistributionZone_count, error) {
	dbpool := pgclient.WDB
	out_arr := []models.DistributionZone{}
	g := models.DistributionZone{}

	err := dbpool.QueryRow(ctx, "SELECT * from func_distribution_zone_get($1);", i).Scan(&g.Id, &g.DistributionZoneName)

	if err != nil && err != pgx.ErrNoRows {
		log.Println("Failed execute from func_distribution_zone_get: ", err)
		return models.DistributionZone_count{Values: []models.DistributionZone{}, Count: 0, Auth: models.Auth{}}, err
	}

	out_arr = append(out_arr, g)

	out_count := models.DistributionZone_count{Values: out_arr, Count: 1, Auth: models.Auth{}}
	return out_count, nil
}
