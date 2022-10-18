package pgsql

import (
	"context"
	"log"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/u2lentaru/billbck/internal/models"
	"github.com/u2lentaru/billbck/pkg/pgclient"
)

//type EquipmentTypeStorage struct
type EquipmentTypeStorage struct {
	db *pgxpool.Pool
}

//func NewEquipmentTypeStorage(db *pgxpool.Pool) *EquipmentTypeStorage
func NewEquipmentTypeStorage(db *pgxpool.Pool) *EquipmentTypeStorage {
	return &EquipmentTypeStorage{db: db}
}

//func (est *EquipmentTypeStorage) GetList(ctx context.Context, pg, pgs int, gs1 string, ord int, dsc bool) (models.EquipmentType_count, error)
func (est *EquipmentTypeStorage) GetList(ctx context.Context, pg, pgs int, gs1 string, ord int, dsc bool) (models.EquipmentType_count, error) {
	dbpool := pgclient.WDB
	gs := models.EquipmentType{}

	gsc := 0
	err := dbpool.QueryRow(ctx, "SELECT * from func_equipment_types_cnt($1);", gs1).Scan(&gsc)

	if err != nil {
		log.Println(err.Error(), "func_equipment_types_cnt")
		return models.EquipmentType_count{Values: []models.EquipmentType{}, Count: gsc, Auth: models.Auth{}}, err
	}

	out_arr := make([]models.EquipmentType, 0,
		func() int {
			if gsc < pgs {
				return gsc
			} else {
				return pgs
			}
		}())

	rows, err := dbpool.Query(ctx, "SELECT * from func_equipment_types_get($1,$2,$3,$4,$5);", pg, pgs, gs1, ord, dsc)
	if err != nil {
		log.Println(err.Error())
		return models.EquipmentType_count{Values: []models.EquipmentType{}, Count: gsc, Auth: models.Auth{}}, err
	}

	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&gs.Id, &gs.EquipmentTypeName, &gs.EquipmentTypePower)
		if err != nil {
			log.Println("failed to scan row:", err)
		}

		out_arr = append(out_arr, gs)
	}

	out_count := models.EquipmentType_count{Values: out_arr, Count: gsc, Auth: models.Auth{}}
	return out_count, nil
}

//func (est *EquipmentTypeStorage) Add(ctx context.Context, a models.EquipmentType) (int, error)
func (est *EquipmentTypeStorage) Add(ctx context.Context, a models.EquipmentType) (int, error) {
	dbpool := pgclient.WDB
	ai := 0

	err := dbpool.QueryRow(ctx, "SELECT func_equipment_types_add($1,$2);", a.EquipmentTypeName,
		a.EquipmentTypePower).Scan(&ai)

	if err != nil {
		log.Println("Failed execute func_equipment_types_add: ", err)
		return 0, err
	}

	return ai, nil
}

//func (est *EquipmentTypeStorage) Upd(ctx context.Context, u models.EquipmentType) (int, error)
func (est *EquipmentTypeStorage) Upd(ctx context.Context, u models.EquipmentType) (int, error) {
	dbpool := pgclient.WDB
	ui := 0

	err := dbpool.QueryRow(ctx, "SELECT func_equipment_types_upd($1,$2,$3);", u.Id, u.EquipmentTypeName,
		u.EquipmentTypePower).Scan(&ui)

	if err != nil {
		log.Println("Failed execute func_equipment_types_upd: ", err)
		return 0, err
	}
	return ui, nil
}

//func (est *EquipmentTypeStorage) Del(ctx context.Context, d []int) ([]int, error)
func (est *EquipmentTypeStorage) Del(ctx context.Context, d []int) ([]int, error) {
	dbpool := pgclient.WDB
	res := []int{}
	i := 0
	for _, id := range d {
		err := dbpool.QueryRow(ctx, "SELECT func_equipment_types_del($1);", id).Scan(&i)
		res = append(res, i)

		if err != nil {
			log.Println("Failed execute func_equipment_types_del: ", err)
		}
	}
	return res, nil
}

//func (est *EquipmentTypeStorage) GetOne(ctx context.Context, i int) (models.EquipmentType_count, error)
func (est *EquipmentTypeStorage) GetOne(ctx context.Context, i int) (models.EquipmentType_count, error) {
	dbpool := pgclient.WDB
	out_arr := []models.EquipmentType{}
	g := models.EquipmentType{}

	err := dbpool.QueryRow(ctx, "SELECT * from func_equipment_type_get($1);", i).Scan(&g.Id, &g.EquipmentTypeName,
		&g.EquipmentTypePower)

	if err != nil && err != pgx.ErrNoRows {
		log.Println("Failed execute from func_equipment_type_get: ", err)
		return models.EquipmentType_count{Values: []models.EquipmentType{}, Count: 0, Auth: models.Auth{}}, err
	}

	out_arr = append(out_arr, g)

	out_count := models.EquipmentType_count{Values: out_arr, Count: 1, Auth: models.Auth{}}
	return out_count, nil
}
