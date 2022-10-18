package pgsql

import (
	"context"
	"log"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/u2lentaru/billbck/internal/models"
	"github.com/u2lentaru/billbck/pkg/pgclient"
)

//type ViolationStorage struct
type ViolationStorage struct {
	db *pgxpool.Pool
}

//func NewViolationStorage(db *pgxpool.Pool) *ViolationStorage
func NewViolationStorage(db *pgxpool.Pool) *ViolationStorage {
	return &ViolationStorage{db: db}
}

//func (est *ViolationStorage) GetList(ctx context.Context, pg, pgs int, gs1 string, ord int, dsc bool) (models.Violation_count, error)
func (est *ViolationStorage) GetList(ctx context.Context, pg, pgs int, gs1 string, ord int, dsc bool) (models.Violation_count, error) {
	dbpool := pgclient.WDB
	gs := models.Violation{}

	gsc := 0
	err := dbpool.QueryRow(ctx, "SELECT * from func_violations_cnt($1);", gs1).Scan(&gsc)

	if err != nil {
		log.Println(err.Error(), "func_violations_cnt")
		return models.Violation_count{Values: []models.Violation{}, Count: gsc, Auth: models.Auth{}}, err
	}

	out_arr := make([]models.Violation, 0,
		func() int {
			if gsc < pgs {
				return gsc
			} else {
				return pgs
			}
		}())

	rows, err := dbpool.Query(ctx, "SELECT * from func_violations_get($1,$2,$3,$4,$5);", pg, pgs, gs1, ord, dsc)
	if err != nil {
		log.Println(err.Error(), "func_violations_get")
		return models.Violation_count{Values: []models.Violation{}, Count: gsc, Auth: models.Auth{}}, err
	}

	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&gs.Id, &gs.ViolationName)
		if err != nil {
			log.Println("failed to scan row:", err)
		}

		out_arr = append(out_arr, gs)
	}

	out_count := models.Violation_count{Values: out_arr, Count: gsc, Auth: models.Auth{}}
	return out_count, nil
}

//func (est *ViolationStorage) Add(ctx context.Context, a models.Violation) (int, error)
func (est *ViolationStorage) Add(ctx context.Context, a models.Violation) (int, error) {
	dbpool := pgclient.WDB
	ai := 0

	err := dbpool.QueryRow(ctx, "SELECT func_violations_add($1);", a.ViolationName).Scan(&ai)

	if err != nil {
		log.Println("Failed execute func_violations_add: ", err)
		return 0, err
	}

	return ai, nil
}

//func (est *ViolationStorage) Upd(ctx context.Context, u models.Violation) (int, error)
func (est *ViolationStorage) Upd(ctx context.Context, u models.Violation) (int, error) {
	dbpool := pgclient.WDB
	ui := 0

	err := dbpool.QueryRow(ctx, "SELECT func_violations_upd($1,$2);", u.Id, u.ViolationName).Scan(&ui)

	if err != nil {
		log.Println("Failed execute func_violations_upd: ", err)
		return 0, err
	}
	return ui, nil
}

//func (est *ViolationStorage) Del(ctx context.Context, d []int) ([]int, error)
func (est *ViolationStorage) Del(ctx context.Context, d []int) ([]int, error) {
	dbpool := pgclient.WDB
	res := []int{}
	i := 0
	for _, id := range d {
		err := dbpool.QueryRow(ctx, "SELECT func_violations_del($1);", id).Scan(&i)
		res = append(res, i)

		if err != nil {
			log.Println("Failed execute func_violations_del: ", err)
		}
	}
	return res, nil
}

//func (est *ViolationStorage) GetOne(ctx context.Context, i int) (models.Violation_count, error)
func (est *ViolationStorage) GetOne(ctx context.Context, i int) (models.Violation_count, error) {
	dbpool := pgclient.WDB
	out_arr := []models.Violation{}
	g := models.Violation{}

	err := dbpool.QueryRow(ctx, "SELECT * from func_violation_get($1);", i).Scan(&g.Id, &g.ViolationName)

	if err != nil && err != pgx.ErrNoRows {
		log.Println("Failed execute func_violation_get: ", err)
		return models.Violation_count{Values: []models.Violation{}, Count: 0, Auth: models.Auth{}}, err
	}

	out_arr = append(out_arr, g)

	out_count := models.Violation_count{Values: out_arr, Count: 1, Auth: models.Auth{}}
	return out_count, nil
}
