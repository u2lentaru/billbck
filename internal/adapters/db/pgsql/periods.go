package pgsql

import (
	"context"
	"log"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/u2lentaru/billbck/internal/models"
	"github.com/u2lentaru/billbck/pkg/pgclient"
)

//type PeriodStorage struct
type PeriodStorage struct {
	db *pgxpool.Pool
}

//func NewPeriodStorage(db *pgxpool.Pool) *PeriodStorage
func NewPeriodStorage(db *pgxpool.Pool) *PeriodStorage {
	return &PeriodStorage{db: db}
}

//func (est *PeriodStorage) GetList(ctx context.Context, pg, pgs int, gs1 string, ord int, dsc bool) (models.Period_count, error)
func (est *PeriodStorage) GetList(ctx context.Context, pg, pgs int, gs1 string, ord int, dsc bool) (models.Period_count, error) {
	dbpool := pgclient.WDB
	gs := models.Period{}

	gsc := 0
	err := dbpool.QueryRow(ctx, "SELECT * from func_periods_cnt($1);", gs1).Scan(&gsc)

	if err != nil {
		log.Println(err.Error(), "func_periods_cnt")
		return models.Period_count{Values: []models.Period{}, Count: gsc, Auth: models.Auth{}}, err
	}

	out_arr := make([]models.Period, 0,
		func() int {
			if gsc < pgs {
				return gsc
			} else {
				return pgs
			}
		}())

	rows, err := dbpool.Query(ctx, "SELECT * from func_periods_get($1,$2,$3,$4,$5);", pg, pgs, gs1, ord, dsc)
	if err != nil {
		log.Println(err.Error(), "func_periods_get")
		return models.Period_count{Values: []models.Period{}, Count: gsc, Auth: models.Auth{}}, err
	}

	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&gs.Id, &gs.PeriodName, &gs.Startdate, &gs.Enddate)
		if err != nil {
			log.Println("failed to scan row:", err)
		}

		out_arr = append(out_arr, gs)
	}

	out_count := models.Period_count{Values: out_arr, Count: gsc, Auth: models.Auth{}}
	return out_count, nil
}

//func (est *PeriodStorage) Add(ctx context.Context, a models.Period) (int, error)
func (est *PeriodStorage) Add(ctx context.Context, a models.Period) (int, error) {
	dbpool := pgclient.WDB
	ai := 0

	err := dbpool.QueryRow(ctx, "SELECT func_periods_add($1,$2,$3);", a.PeriodName, a.Startdate, a.Enddate).Scan(&ai)

	if err != nil {
		log.Println("Failed execute func_periods_add: ", err)
		return 0, err
	}

	return ai, nil
}

//func (est *PeriodStorage) Upd(ctx context.Context, u models.Period) (int, error)
func (est *PeriodStorage) Upd(ctx context.Context, u models.Period) (int, error) {
	dbpool := pgclient.WDB
	ui := 0

	err := dbpool.QueryRow(ctx, "SELECT func_periods_upd($1,$2,$3,$4);", u.Id, u.PeriodName, u.Startdate, u.Enddate).Scan(&ui)

	if err != nil {
		log.Println("Failed execute func_periods_upd: ", err)
		return 0, err
	}
	return ui, nil
}

//func (est *PeriodStorage) Del(ctx context.Context, d []int) ([]int, error)
func (est *PeriodStorage) Del(ctx context.Context, d []int) ([]int, error) {
	dbpool := pgclient.WDB
	res := []int{}
	i := 0
	for _, id := range d {
		err := dbpool.QueryRow(ctx, "SELECT func_periods_del($1);", id).Scan(&i)
		res = append(res, i)

		if err != nil {
			log.Println("Failed execute func_periods_del: ", err)
		}
	}
	return res, nil
}

//func (est *PeriodStorage) GetOne(ctx context.Context, i int) (models.Period_count, error)
func (est *PeriodStorage) GetOne(ctx context.Context, i int) (models.Period_count, error) {
	dbpool := pgclient.WDB
	out_arr := []models.Period{}
	g := models.Period{}

	err := dbpool.QueryRow(ctx, "SELECT * from func_period_get($1);", i).Scan(&g.Id, &g.PeriodName, &g.Startdate,
		&g.Enddate)

	if err != nil && err != pgx.ErrNoRows {
		log.Println("Failed execute func_period_get: ", err)
		return models.Period_count{Values: []models.Period{}, Count: 0, Auth: models.Auth{}}, err
	}

	out_arr = append(out_arr, g)

	out_count := models.Period_count{Values: out_arr, Count: 1, Auth: models.Auth{}}
	return out_count, nil
}
