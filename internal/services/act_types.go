package services

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/u2lentaru/billbck/internal/adapters/db/pg"
	"github.com/u2lentaru/billbck/internal/models"
)

type ActTypeService struct {
	storage pg.ActTypeStorage
}

type ActTypeStorage interface {
	GetList(ctx context.Context, Dbpool *pgxpool.Pool, pg, pgs int, nm string, ord int, dsc bool) (models.ActType_count, error)
	Add(ctx context.Context, Dbpool *pgxpool.Pool, ea models.ActType) (int, error)
	Upd(ctx context.Context, Dbpool *pgxpool.Pool, eu models.ActType) (int, error)
	Del(ctx context.Context, Dbpool *pgxpool.Pool, ed []int) ([]int, error)
	GetOne(ctx context.Context, Dbpool *pgxpool.Pool, i int) (models.ActType_count, error)
}

//NewActTypeService(storage pg.ActTypeStorage) *ActTypeService
func NewActTypeService(storage pg.ActTypeStorage) *ActTypeService {
	return &ActTypeService{storage}
}

/*
//(dz *srvActType) GetActTypes(ctx context.Context, Dbpool *pgxpool.Pool, pg, pgs int, nm string, ord int, dsc bool) (srvActType_count, error)
func (dz srvActType) GetDistributionZones(ctx context.Context, Dbpool *pgxpool.Pool, pg, pgs int, nm string, ord int, dsc bool) (srvActType_count, error) {
	gsc := 0
	err := Dbpool.QueryRow(ctx, "SELECT * from func_distribution_zones_cnt($1);", nm).Scan(&gsc)
	auth := models.Auth{Create: true, Read: true, Update: true, Delete: true}

	if err != nil {
		log.Println(err.Error(), "func_distribution_zones_cnt")
		return srvActType_count{Values: []models.ActType{}, Count: gsc, Auth: auth}, err
	}

	out_arr := make([]srvActType, 0,
		func() int {
			if gsc < pgs {
				return gsc
			} else {
				return pgs
			}
		}())

	rows, err := Dbpool.Query(ctx, "SELECT * from func_distribution_zones_get($1,$2,$3,$4,$5);", pg, pgs, nm, ord, dsc)
	if err != nil {
		log.Println(err.Error())
		return srvActType_count{Values: []models.ActType{}, Count: gsc, Auth: auth}, err
	}

	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&(dz.Id), &(dz.ActTypeName))
		if err != nil {
			log.Println("failed to scan row:", err)
		}

		out_arr = append(out_arr, *dz)
	}

	out_count := models.ActType_count{Values: out_arr, Count: gsc, Auth: auth}
	if err != nil {
		log.Println(err.Error())
		return models.ActType_count{}, err
	}

	return out_count, nil
}

//func (dz *srvActType) AddActType(ctx context.Context, Dbpool *pgxpool.Pool) (int, error)
func (dz *srvActType) AddActType(ctx context.Context, Dbpool *pgxpool.Pool) (int, error) {
	ai := 0
	err := Dbpool.QueryRow(ctx, "SELECT func_distribution_zones_add($1);", dz.ActTypeName).Scan(&ai)

	if err != nil {
		log.Println("Failed execute func_distribution_zones_add: ", err)
		return 0, err
	}

	return ai, nil
}

//func (dz *srvActType) UpdActType(ctx context.Context, Dbpool *pgxpool.Pool)
func (dz *srvActType) UpdActType(ctx context.Context, Dbpool *pgxpool.Pool) (int, error) {
	ui := 0
	err := Dbpool.QueryRow(context.Background(), "SELECT func_distribution_zones_upd($1,$2);", dz.Id, dz.ActTypeName).Scan(&ui)

	if err != nil {
		log.Println("Failed execute func_distribution_zones_upd: ", err)
		return 0, err
	}
	return ui, nil
}

//func (dz *srvActType) DelActType(ctx context.Context, Dbpool *pgxpool.Pool, d []int) ([]int, error)
func (dz *srvActType) DelActType(ctx context.Context, Dbpool *pgxpool.Pool, d []int) ([]int, error) {
	res := []int{}
	i := 0
	for _, id := range d {
		err := Dbpool.QueryRow(ctx, "SELECT func_distribution_zones_del($1);", id).Scan(&i)
		res = append(res, i)

		if err != nil {
			log.Println("Failed execute func_distribution_zones_del: ", err)
			return []int{}, err
		}
	}
	return res, nil
}

//func (dz *srvActType) GetActType(ctx context.Context, Dbpool *pgxpool.Pool, i int) (models.ActType_count, error)
func (dz *srvActType) GetDistributionZone(ctx context.Context, Dbpool *pgxpool.Pool, i int) (models.ActType_count, error) {
	out_arr := []models.ActType{}
	auth := models.Auth{Create: true, Read: true, Update: true, Delete: true}

	err := Dbpool.QueryRow(context.Background(), "SELECT * from func_distribution_zone_get($1);", i).Scan(&(dz.Id), &(dz.ActTypeName))

	if err != nil && err != pgx.ErrNoRows {
		log.Println("Failed execute from func_distribution_zone_get: ", err)
		return models.ActType_count{Values: []models.ActType{}, Count: 0, Auth: auth}, err
	}

	out_arr = append(out_arr, dz)

	out_count := models.ActType_count{Values: out_arr, Count: 0, Auth: auth}
	return out_count, nil
}*/
