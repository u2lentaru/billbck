package pgsql

import (
	"context"
	"log"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/u2lentaru/billbck/internal/models"
	"github.com/u2lentaru/billbck/pkg/pgclient"
)

//type ServiceTypeStorage struct
type ServiceTypeStorage struct {
	db *pgxpool.Pool
}

//func NewServiceTypeStorage(db *pgxpool.Pool) *ServiceTypeStorage
func NewServiceTypeStorage(db *pgxpool.Pool) *ServiceTypeStorage {
	return &ServiceTypeStorage{db: db}
}

//func (est *ServiceTypeStorage) GetList(ctx context.Context, pg, pgs int, gs1 string, ord int, dsc bool) (models.ServiceType_count, error)
func (est *ServiceTypeStorage) GetList(ctx context.Context, pg, pgs int, gs1 string, ord int, dsc bool) (models.ServiceType_count, error) {
	dbpool := pgclient.WDB
	gs := models.ServiceType{}

	gsc := 0
	err := dbpool.QueryRow(ctx, "SELECT * from func_service_types_cnt($1);", gs1).Scan(&gsc)

	if err != nil {
		log.Println(err.Error(), "func_service_types_cnt")
		return models.ServiceType_count{Values: []models.ServiceType{}, Count: gsc, Auth: models.Auth{}}, err
	}

	out_arr := make([]models.ServiceType, 0,
		func() int {
			if gsc < pgs {
				return gsc
			} else {
				return pgs
			}
		}())

	rows, err := dbpool.Query(ctx, "SELECT * from func_service_types_get($1,$2,$3,$4,$5);", pg, pgs, gs1, ord, dsc)
	if err != nil {
		log.Println(err.Error(), "func_service_types_get")
		return models.ServiceType_count{Values: []models.ServiceType{}, Count: gsc, Auth: models.Auth{}}, err
	}

	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&gs.Id, &gs.ServiceTypeName)
		if err != nil {
			log.Println("failed to scan row:", err)
		}

		out_arr = append(out_arr, gs)
	}

	out_count := models.ServiceType_count{Values: out_arr, Count: gsc, Auth: models.Auth{}}
	return out_count, nil
}

//func (est *ServiceTypeStorage) Add(ctx context.Context, a models.ServiceType) (int, error)
func (est *ServiceTypeStorage) Add(ctx context.Context, a models.ServiceType) (int, error) {
	dbpool := pgclient.WDB
	ai := 0

	err := dbpool.QueryRow(ctx, "SELECT func_service_types_add($1);", a.ServiceTypeName).Scan(&ai)

	if err != nil {
		log.Println("Failed execute func_service_types_add: ", err)
		return 0, err
	}

	return ai, nil
}

//func (est *ServiceTypeStorage) Upd(ctx context.Context, u models.ServiceType) (int, error)
func (est *ServiceTypeStorage) Upd(ctx context.Context, u models.ServiceType) (int, error) {
	dbpool := pgclient.WDB
	ui := 0

	err := dbpool.QueryRow(ctx, "SELECT func_service_types_upd($1,$2);", u.Id, u.ServiceTypeName).Scan(&ui)

	if err != nil {
		log.Println("Failed execute func_service_types_upd: ", err)
		return 0, err
	}
	return ui, nil
}

//func (est *ServiceTypeStorage) Del(ctx context.Context, d []int) ([]int, error)
func (est *ServiceTypeStorage) Del(ctx context.Context, d []int) ([]int, error) {
	dbpool := pgclient.WDB
	res := []int{}
	i := 0
	for _, id := range d {
		err := dbpool.QueryRow(ctx, "SELECT func_service_types_del($1);", id).Scan(&i)
		res = append(res, i)

		if err != nil {
			log.Println("Failed execute func_service_types_del: ", err)
		}
	}
	return res, nil
}

//func (est *ServiceTypeStorage) GetOne(ctx context.Context, i int) (models.ServiceType_count, error)
func (est *ServiceTypeStorage) GetOne(ctx context.Context, i int) (models.ServiceType_count, error) {
	dbpool := pgclient.WDB
	out_arr := []models.ServiceType{}
	g := models.ServiceType{}

	err := dbpool.QueryRow(ctx, "SELECT * from func_service_type_get($1);", i).Scan(&g.Id, &g.ServiceTypeName)

	if err != nil && err != pgx.ErrNoRows {
		log.Println("Failed execute func_service_type_get: ", err)
		return models.ServiceType_count{Values: []models.ServiceType{}, Count: 0, Auth: models.Auth{}}, err
	}

	out_arr = append(out_arr, g)

	out_count := models.ServiceType_count{Values: out_arr, Count: 1, Auth: models.Auth{}}
	return out_count, nil
}
