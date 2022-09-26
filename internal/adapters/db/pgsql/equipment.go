package pgsql

import (
	"context"
	"log"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/u2lentaru/billbck/internal/models"
	"github.com/u2lentaru/billbck/internal/utils"
	"github.com/u2lentaru/billbck/pkg/pgclient"
)

//type EquipmentStorage struct
type EquipmentStorage struct {
	db *pgxpool.Pool
}

//func NewEquipmentStorage(db *pgxpool.Pool) *EquipmentStorage
func NewEquipmentStorage(db *pgxpool.Pool) *EquipmentStorage {
	return &EquipmentStorage{db: db}
}

//func (est *EquipmentStorage) GetList(ctx context.Context, pg, pgs, gs1 int, gs2 string, ord int, dsc bool) (models.Equipment_count, error) {
func (est *EquipmentStorage) GetList(ctx context.Context, pg, pgs, gs1 int, gs2 string, ord int, dsc bool) (models.Equipment_count, error) {
	dbpool := pgclient.WDB
	gs := models.Equipment{}

	gsc := 0
	err := dbpool.QueryRow(ctx, "SELECT * from func_equipment_cnt($1,$2);", gs2, utils.NullableInt(int32(gs1))).Scan(&gsc)

	if err != nil {
		log.Println(err.Error(), "func_equipment_cnt")
		return models.Equipment_count{Values: []models.Equipment{}, Count: gsc, Auth: models.Auth{}}, err
	}

	out_arr := make([]models.Equipment, 0,
		func() int {
			if gsc < pgs {
				return gsc
			} else {
				return pgs
			}
		}())

	rows, err := dbpool.Query(ctx, "SELECT * from func_equipment_get($1,$2,$3,$4,$5,$6);", pg, pgs, gs2, utils.NullableInt(int32(gs1)), ord, dsc)
	if err != nil {
		log.Println(err.Error())
		return models.Equipment_count{Values: []models.Equipment{}, Count: gsc, Auth: models.Auth{}}, err
	}

	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&gs.Id, &gs.EquipmentType.Id, &gs.Object.Id, &gs.Qty, &gs.WorkingHours, &gs.EquipmentType.EquipmentTypeName,
			&gs.EquipmentType.EquipmentTypePower, &gs.Object.ObjectName)
		if err != nil {
			log.Println("failed to scan row:", err)
		}

		out_arr = append(out_arr, gs)
	}

	out_count := models.Equipment_count{Values: out_arr, Count: gsc, Auth: models.Auth{}}
	if err != nil {
		log.Println(err.Error())
		return models.Equipment_count{}, err
	}

	return out_count, nil
}

//func (est *EquipmentStorage) Add(ctx context.Context, ea models.Equipment) (int, error)
func (est *EquipmentStorage) Add(ctx context.Context, a models.Equipment) (int, error) {
	dbpool := pgclient.WDB
	ai := 0

	err := dbpool.QueryRow(context.Background(), "SELECT func_equipment_add($1,$2,$3,$4);", a.EquipmentType.Id, a.Object.Id,
		a.Qty, a.WorkingHours).Scan(&ai)

	if err != nil {
		log.Println("Failed execute func_equipment_add: ", err)
		return 0, err
	}

	return ai, nil
}

//func (est *EquipmentStorage) Upd(ctx context.Context, eu models.Equipment) (int, error)
func (est *EquipmentStorage) Upd(ctx context.Context, u models.Equipment) (int, error) {
	dbpool := pgclient.WDB
	ui := 0

	err := dbpool.QueryRow(context.Background(), "SELECT func_equipment_upd($1,$2,$3,$4,$5);", u.Id, u.EquipmentType.Id, u.Object.Id,
		u.Qty, u.WorkingHours).Scan(&ui)

	if err != nil {
		log.Println("Failed execute func_equipment_upd: ", err)
		return 0, err
	}
	return ui, nil
}

//func (est *EquipmentStorage) Del(ctx context.Context, ed []int) ([]int, error)
func (est *EquipmentStorage) Del(ctx context.Context, d []int) ([]int, error) {
	dbpool := pgclient.WDB
	res := []int{}
	i := 0
	for _, id := range d {
		err := dbpool.QueryRow(context.Background(), "SELECT func_equipment_del($1);", id).Scan(&i)
		res = append(res, i)

		if err != nil {
			log.Println("Failed execute func_equipment_del: ", err)
		}
	}
	return res, nil
}

//func (est *EquipmentStorage) GetOne(ctx context.Context, i int) (models.Equipment_count, error)
func (est *EquipmentStorage) GetOne(ctx context.Context, i int) (models.Equipment_count, error) {
	dbpool := pgclient.WDB
	out_arr := []models.Equipment{}
	g := models.Equipment{}

	err := dbpool.QueryRow(context.Background(), "SELECT * from func_equipment_getbyid($1);", i).Scan(&g.Id, &g.EquipmentType.Id, &g.Object.Id,
		&g.Qty, &g.WorkingHours, &g.EquipmentType.EquipmentTypeName, &g.EquipmentType.EquipmentTypePower, &g.Object.ObjectName)

	if err != nil && err != pgx.ErrNoRows {
		log.Println("Failed execute func_equipment_getbyid: ", err)
		return models.Equipment_count{Values: []models.Equipment{}, Count: 0, Auth: models.Auth{}}, err
	}

	out_arr = append(out_arr, g)

	out_count := models.Equipment_count{Values: out_arr, Count: 1, Auth: models.Auth{}}
	return out_count, nil
}

//func (est *EquipmentStorage) AddList(ctx context.Context, al models.Equipment_count) ([]int, error)
func (est *EquipmentStorage) AddList(ctx context.Context, al models.Equipment_count) ([]int, error) {
	dbpool := pgclient.WDB

	res := []int{}
	i := 0
	first_value := true

	for _, a := range al.Values {

		if first_value {
			err := dbpool.QueryRow(context.Background(), "SELECT func_equipment_delbyobj($1);", a.Object.Id).Scan(&i)

			if err != nil {
				log.Println("Failed execute func_equipment_delbyobj: ", err)
			}

			first_value = false
			i = 0
		}

		err := dbpool.QueryRow(context.Background(), "SELECT func_equipment_add($1,$2,$3,$4);", a.EquipmentType.Id, a.Object.Id,
			a.Qty, a.WorkingHours).Scan(&i)
		res = append(res, i)

		if err != nil {
			log.Println("Failed execute func_equipment_add: ", err)
		}
	}
	return res, nil
}

//func (est *EquipmentStorage) DelObj(ctx context.Context, d []int) ([]int, error)
func (est *EquipmentStorage) DelObj(ctx context.Context, d []int) ([]int, error) {
	dbpool := pgclient.WDB
	res := []int{}
	i := 0
	for _, id := range d {
		err := dbpool.QueryRow(context.Background(), "SELECT func_equipment_delbyobj($1);", id).Scan(&i)
		res = append(res, i)

		if err != nil {
			log.Println("Failed execute func_equipment_delbyobj: ", err)
		}
	}
	return res, nil
}
