package pgsql

import (
	"context"
	"log"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/u2lentaru/billbck/internal/models"
	"github.com/u2lentaru/billbck/pkg/pgclient"
)

//type BuildingTypeStorage struct
type BuildingTypeStorage struct {
	db *pgxpool.Pool
}

//func NewBuildingTypeStorage(db *pgxpool.Pool) *BuildingTypeStorage
func NewBuildingTypeStorage(db *pgxpool.Pool) *BuildingTypeStorage {
	return &BuildingTypeStorage{db: db}
}

//func (est *BuildingTypeStorage) GetList(ctx context.Context, pg, pgs int, gs1 string, ord int, dsc bool) (models.BuildingType_count, error)
func (est *BuildingTypeStorage) GetList(ctx context.Context, pg, pgs int, gs1 string, ord int, dsc bool) (models.BuildingType_count, error) {
	dbpool := pgclient.WDB
	auth := models.Auth{Create: true, Read: true, Update: true, Delete: true}
	gs := models.BuildingType{}

	gsc := 0
	err := dbpool.QueryRow(ctx, "SELECT * from func_building_types_cnt($1);", gs1).Scan(&gsc)

	if err != nil {
		log.Println(err.Error(), "func_building_types_cnt")
		return models.BuildingType_count{Values: []models.BuildingType{}, Count: gsc, Auth: auth}, err
	}

	out_arr := make([]models.BuildingType, 0,
		func() int {
			if gsc < pgs {
				return gsc
			} else {
				return pgs
			}
		}())

	rows, err := dbpool.Query(ctx, "SELECT * from func_building_types_get($1,$2,$3,$4,$5);", pg, pgs, gs1, ord, dsc)
	if err != nil {
		log.Println(err.Error())
		return models.BuildingType_count{Values: []models.BuildingType{}, Count: gsc, Auth: auth}, err
	}

	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&gs.Id, &gs.BuildingTypeName)
		if err != nil {
			log.Println("failed to scan row:", err)
		}

		out_arr = append(out_arr, gs)
	}

	out_count := models.BuildingType_count{Values: out_arr, Count: gsc, Auth: auth}
	if err != nil {
		log.Println(err.Error())
		return models.BuildingType_count{}, err
	}

	return out_count, nil
}

//func (est *BuildingTypeStorage) Add(ctx context.Context, a models.BuildingType) (int, error)
func (est *BuildingTypeStorage) Add(ctx context.Context, a models.BuildingType) (int, error) {
	dbpool := pgclient.WDB
	ai := 0

	err := dbpool.QueryRow(context.Background(), "SELECT func_building_types_add($1);", a.BuildingTypeName).Scan(&ai)

	if err != nil {
		log.Println("Failed execute func_building_types_add: ", err)
		return 0, err
	}

	return ai, nil
}

//func (est *BuildingTypeStorage) Upd(ctx context.Context, u models.BuildingType) (int, error)
func (est *BuildingTypeStorage) Upd(ctx context.Context, u models.BuildingType) (int, error) {
	dbpool := pgclient.WDB
	ui := 0

	err := dbpool.QueryRow(context.Background(), "SELECT func_building_types_upd($1,$2);", u.Id, u.BuildingTypeName).Scan(&ui)

	if err != nil {
		log.Println("Failed execute func_building_types_upd: ", err)
		return 0, err
	}
	return ui, nil
}

//func (est *BuildingTypeStorage) Del(ctx context.Context, d []int) ([]int, error)
func (est *BuildingTypeStorage) Del(ctx context.Context, d []int) ([]int, error) {
	dbpool := pgclient.WDB
	res := []int{}
	i := 0
	for _, id := range d {
		err := dbpool.QueryRow(context.Background(), "SELECT func_building_types_del($1);", id).Scan(&i)
		res = append(res, i)

		if err != nil {
			log.Println("Failed execute func_building_types_del: ", err)
		}
	}
	return res, nil
}

//func (est *BuildingTypeStorage) GetOne(ctx context.Context, i int) (models.BuildingType_count, error)
func (est *BuildingTypeStorage) GetOne(ctx context.Context, i int) (models.BuildingType_count, error) {
	dbpool := pgclient.WDB
	out_arr := []models.BuildingType{}
	g := models.BuildingType{}
	auth := models.Auth{Create: true, Read: true, Update: true, Delete: true}

	err := dbpool.QueryRow(context.Background(), "SELECT * from func_building_type_get($1);", i).Scan(&g.Id, &g.BuildingTypeName)

	if err != nil && err != pgx.ErrNoRows {
		log.Println("Failed execute from func_building_type_get: ", err)
		return models.BuildingType_count{Values: []models.BuildingType{}, Count: 0, Auth: auth}, err
	}

	out_arr = append(out_arr, g)

	out_count := models.BuildingType_count{Values: out_arr, Count: 0, Auth: auth}
	return out_count, nil
}
