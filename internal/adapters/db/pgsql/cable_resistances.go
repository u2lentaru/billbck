package pgsql

import (
	"context"
	"log"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/u2lentaru/billbck/internal/models"
	"github.com/u2lentaru/billbck/pkg/pgclient"
)

//type CableResistanceStorage struct
type CableResistanceStorage struct {
	db *pgxpool.Pool
}

//func NewCableResistanceStorage(db *pgxpool.Pool) *CableResistanceStorage
func NewCableResistanceStorage(db *pgxpool.Pool) *CableResistanceStorage {
	return &CableResistanceStorage{db: db}
}

//func (est *CableResistanceStorage) GetList(ctx context.Context, pg, pgs int, gs1 string, ord int, dsc bool) (models.CableResistance_count, error)
func (est *CableResistanceStorage) GetList(ctx context.Context, pg, pgs int, gs1 string, ord int, dsc bool) (models.CableResistance_count, error) {
	dbpool := pgclient.WDB
	auth := models.Auth{Create: true, Read: true, Update: true, Delete: true}
	gs := models.CableResistance{}

	gsc := 0
	err := dbpool.QueryRow(ctx, "SELECT * from func_cable_resistances_cnt($1);", gs1).Scan(&gsc)

	if err != nil {
		log.Println(err.Error(), "func_cable_resistances_cnt")
		return models.CableResistance_count{Values: []models.CableResistance{}, Count: gsc, Auth: auth}, err
	}

	out_arr := make([]models.CableResistance, 0,
		func() int {
			if gsc < pgs {
				return gsc
			} else {
				return pgs
			}
		}())

	rows, err := dbpool.Query(ctx, "SELECT * from func_cable_resistances_get($1,$2,$3,$4,$5);", pg, pgs, gs1, ord, dsc)
	if err != nil {
		log.Println(err.Error())
		return models.CableResistance_count{Values: []models.CableResistance{}, Count: gsc, Auth: auth}, err
	}

	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&gs.Id, &gs.CableResistanceName, &gs.Resistance, &gs.MaterialType)
		if err != nil {
			log.Println("failed to scan row:", err)
		}

		out_arr = append(out_arr, gs)
	}

	out_count := models.CableResistance_count{Values: out_arr, Count: gsc, Auth: auth}
	if err != nil {
		log.Println(err.Error())
		return models.CableResistance_count{}, err
	}

	return out_count, nil
}

//func (est *CableResistanceStorage) Add(ctx context.Context, a models.CableResistance) (int, error)
func (est *CableResistanceStorage) Add(ctx context.Context, a models.CableResistance) (int, error) {
	dbpool := pgclient.WDB
	ai := 0

	err := dbpool.QueryRow(context.Background(), "SELECT func_cable_resistances_add($1,$2,$3);", a.CableResistanceName,
		a.Resistance, a.MaterialType).Scan(&ai)

	if err != nil {
		log.Println("Failed execute func_cable_resistances_add: ", err)
		return 0, err
	}

	return ai, nil
}

//func (est *CableResistanceStorage) Upd(ctx context.Context, u models.CableResistance) (int, error)
func (est *CableResistanceStorage) Upd(ctx context.Context, u models.CableResistance) (int, error) {
	dbpool := pgclient.WDB
	ui := 0

	err := dbpool.QueryRow(context.Background(), "SELECT func_cable_resistances_upd($1,$2,$3,$4);", u.Id, u.CableResistanceName,
		u.Resistance, u.MaterialType).Scan(&ui)

	if err != nil {
		log.Println("Failed execute func_cable_resistances_upd: ", err)
		return 0, err
	}
	return ui, nil
}

//func (est *CableResistanceStorage) Del(ctx context.Context, d []int) ([]int, error)
func (est *CableResistanceStorage) Del(ctx context.Context, d []int) ([]int, error) {
	dbpool := pgclient.WDB
	res := []int{}
	i := 0
	for _, id := range d {
		err := dbpool.QueryRow(context.Background(), "SELECT func_cable_resistances_del($1);", id).Scan(&i)
		res = append(res, i)

		if err != nil {
			log.Println("Failed execute func_cable_resistances_del: ", err)
		}
	}
	return res, nil
}

//func (est *CableResistanceStorage) GetOne(ctx context.Context, i int) (models.CableResistance_count, error)
func (est *CableResistanceStorage) GetOne(ctx context.Context, i int) (models.CableResistance_count, error) {
	dbpool := pgclient.WDB
	out_arr := []models.CableResistance{}
	g := models.CableResistance{}
	auth := models.Auth{Create: true, Read: true, Update: true, Delete: true}

	err := dbpool.QueryRow(context.Background(), "SELECT * from func_cable_resistance_get($1);", i).Scan(&g.Id,
		&g.CableResistanceName, &g.Resistance, &g.MaterialType)

	if err != nil && err != pgx.ErrNoRows {
		log.Println("Failed execute from func_cable_resistance_get: ", err)
		return models.CableResistance_count{Values: []models.CableResistance{}, Count: 0, Auth: auth}, err
	}

	out_arr = append(out_arr, g)

	out_count := models.CableResistance_count{Values: out_arr, Count: 0, Auth: auth}
	return out_count, nil
}
