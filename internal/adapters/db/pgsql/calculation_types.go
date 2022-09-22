package pgsql

import (
	"context"
	"log"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/u2lentaru/billbck/internal/models"
	"github.com/u2lentaru/billbck/pkg/pgclient"
)

//type CalculationTypeStorage struct
type CalculationTypeStorage struct {
	db *pgxpool.Pool
}

//func NewCalculationTypeStorage(db *pgxpool.Pool) *CalculationTypeStorage
func NewCalculationTypeStorage(db *pgxpool.Pool) *CalculationTypeStorage {
	return &CalculationTypeStorage{db: db}
}

//func (est *CalculationTypeStorage) GetList(ctx context.Context, pg, pgs int, gs1 string, ord int, dsc bool) (models.CalculationType_count, error)
func (est *CalculationTypeStorage) GetList(ctx context.Context, pg, pgs int, gs1 string, ord int, dsc bool) (models.CalculationType_count, error) {
	dbpool := pgclient.WDB
	gs := models.CalculationType{}

	gsc := 0
	err := dbpool.QueryRow(ctx, "SELECT * from func_calculation_types_cnt($1);", gs1).Scan(&gsc)

	if err != nil {
		log.Println(err.Error(), "func_calculation_types_cnt")
		return models.CalculationType_count{Values: []models.CalculationType{}, Count: gsc, Auth: models.Auth{}}, err
	}

	out_arr := make([]models.CalculationType, 0,
		func() int {
			if gsc < pgs {
				return gsc
			} else {
				return pgs
			}
		}())

	rows, err := dbpool.Query(ctx, "SELECT * from func_calculation_types_get($1,$2,$3,$4,$5);", pg, pgs, gs1, ord, dsc)
	if err != nil {
		log.Println(err.Error())
		return models.CalculationType_count{Values: []models.CalculationType{}, Count: gsc, Auth: models.Auth{}}, err
	}

	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&gs.Id, &gs.CalculationTypeName)
		if err != nil {
			log.Println("failed to scan row:", err)
		}

		out_arr = append(out_arr, gs)
	}

	out_count := models.CalculationType_count{Values: out_arr, Count: gsc, Auth: models.Auth{}}
	if err != nil {
		log.Println(err.Error())
		return models.CalculationType_count{}, err
	}

	return out_count, nil
}

//func (est *CalculationTypeStorage) Add(ctx context.Context, a models.CalculationType) (int, error)
func (est *CalculationTypeStorage) Add(ctx context.Context, a models.CalculationType) (int, error) {
	dbpool := pgclient.WDB
	ai := 0

	err := dbpool.QueryRow(context.Background(), "SELECT func_calculation_types_add($1);", a.CalculationTypeName).Scan(&ai)

	if err != nil {
		log.Println("Failed execute func_calculation_types_add: ", err)
		return 0, err
	}

	return ai, nil
}

//func (est *CalculationTypeStorage) Upd(ctx context.Context, u models.CalculationType) (int, error)
func (est *CalculationTypeStorage) Upd(ctx context.Context, u models.CalculationType) (int, error) {
	dbpool := pgclient.WDB
	ui := 0

	err := dbpool.QueryRow(context.Background(), "SELECT func_calculation_types_upd($1,$2);", u.Id, u.CalculationTypeName).Scan(&ui)

	if err != nil {
		log.Println("Failed execute func_calculation_types_upd: ", err)
		return 0, err
	}
	return ui, nil
}

//func (est *CalculationTypeStorage) Del(ctx context.Context, d []int) ([]int, error)
func (est *CalculationTypeStorage) Del(ctx context.Context, d []int) ([]int, error) {
	dbpool := pgclient.WDB
	res := []int{}
	i := 0
	for _, id := range d {
		err := dbpool.QueryRow(context.Background(), "SELECT func_calculation_types_del($1);", id).Scan(&i)
		res = append(res, i)

		if err != nil {
			log.Println("Failed execute func_calculation_types_del: ", err)
		}
	}
	return res, nil
}

//func (est *CalculationTypeStorage) GetOne(ctx context.Context, i int) (models.CalculationType_count, error)
func (est *CalculationTypeStorage) GetOne(ctx context.Context, i int) (models.CalculationType_count, error) {
	dbpool := pgclient.WDB
	out_arr := []models.CalculationType{}
	g := models.CalculationType{}

	err := dbpool.QueryRow(context.Background(), "SELECT * from func_calculation_type_get($1);", i).Scan(&g.Id, &g.CalculationTypeName)

	if err != nil && err != pgx.ErrNoRows {
		log.Println("Failed execute from func_calculation_type_get: ", err)
		return models.CalculationType_count{Values: []models.CalculationType{}, Count: 0, Auth: models.Auth{}}, err
	}

	out_arr = append(out_arr, g)

	out_count := models.CalculationType_count{Values: out_arr, Count: 0, Auth: models.Auth{}}
	return out_count, nil
}
