package pgsql

import (
	"context"
	"log"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/u2lentaru/billbck/internal/models"
	"github.com/u2lentaru/billbck/pkg/pgclient"
)

//type SectorStorage struct
type SectorStorage struct {
	db *pgxpool.Pool
}

//func NewSectorStorage(db *pgxpool.Pool) *SectorStorage
func NewSectorStorage(db *pgxpool.Pool) *SectorStorage {
	return &SectorStorage{db: db}
}

//func (est *SectorStorage) GetList(ctx context.Context, pg, pgs int, gs1 string, ord int, dsc bool) (models.Sector_count, error)
func (est *SectorStorage) GetList(ctx context.Context, pg, pgs int, gs1 string, ord int, dsc bool) (models.Sector_count, error) {
	dbpool := pgclient.WDB
	gs := models.Sector{}

	gsc := 0
	err := dbpool.QueryRow(ctx, "SELECT * from func_sectors_cnt($1);", gs1).Scan(&gsc)

	if err != nil {
		log.Println(err.Error(), "func_sectors_cnt")
		return models.Sector_count{Values: []models.Sector{}, Count: gsc, Auth: models.Auth{}}, err
	}

	out_arr := make([]models.Sector, 0,
		func() int {
			if gsc < pgs {
				return gsc
			} else {
				return pgs
			}
		}())

	rows, err := dbpool.Query(ctx, "SELECT * from func_sectors_get($1,$2,$3,$4,$5);", pg, pgs, gs1, ord, dsc)
	if err != nil {
		log.Println(err.Error(), "func_sectors_get")
		return models.Sector_count{Values: []models.Sector{}, Count: gsc, Auth: models.Auth{}}, err
	}

	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&gs.Id, &gs.SectorName)
		if err != nil {
			log.Println("failed to scan row:", err)
		}

		out_arr = append(out_arr, gs)
	}

	out_count := models.Sector_count{Values: out_arr, Count: gsc, Auth: models.Auth{}}
	return out_count, nil
}

//func (est *SectorStorage) Add(ctx context.Context, a models.Sector) (int, error)
func (est *SectorStorage) Add(ctx context.Context, a models.Sector) (int, error) {
	dbpool := pgclient.WDB
	ai := 0

	err := dbpool.QueryRow(ctx, "SELECT func_sectors_add($1);", a.SectorName).Scan(&ai)

	if err != nil {
		log.Println("Failed execute func_sectors_add: ", err)
		return 0, err
	}

	return ai, nil
}

//func (est *SectorStorage) Upd(ctx context.Context, u models.Sector) (int, error)
func (est *SectorStorage) Upd(ctx context.Context, u models.Sector) (int, error) {
	dbpool := pgclient.WDB
	ui := 0

	err := dbpool.QueryRow(ctx, "SELECT func_sectors_upd($1,$2);", u.Id, u.SectorName).Scan(&ui)

	if err != nil {
		log.Println("Failed execute func_sectors_upd: ", err)
		return 0, err
	}
	return ui, nil
}

//func (est *SectorStorage) Del(ctx context.Context, d []int) ([]int, error)
func (est *SectorStorage) Del(ctx context.Context, d []int) ([]int, error) {
	dbpool := pgclient.WDB
	res := []int{}
	i := 0
	for _, id := range d {
		err := dbpool.QueryRow(ctx, "SELECT func_sectors_del($1);", id).Scan(&i)
		res = append(res, i)

		if err != nil {
			log.Println("Failed execute func_sectors_del: ", err)
		}
	}
	return res, nil
}

//func (est *SectorStorage) GetOne(ctx context.Context, i int) (models.Sector_count, error)
func (est *SectorStorage) GetOne(ctx context.Context, i int) (models.Sector_count, error) {
	dbpool := pgclient.WDB
	out_arr := []models.Sector{}
	g := models.Sector{}

	err := dbpool.QueryRow(ctx, "SELECT * from func_sector_get($1);", i).Scan(&g.Id, &g.SectorName)

	if err != nil && err != pgx.ErrNoRows {
		log.Println("Failed execute func_sector_get: ", err)
		return models.Sector_count{Values: []models.Sector{}, Count: 0, Auth: models.Auth{}}, err
	}

	out_arr = append(out_arr, g)

	out_count := models.Sector_count{Values: out_arr, Count: 1, Auth: models.Auth{}}
	return out_count, nil
}
