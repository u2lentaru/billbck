package pgsql

import (
	"context"
	"log"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/u2lentaru/billbck/internal/models"
	"github.com/u2lentaru/billbck/pkg/pgclient"
)

//type VoltageStorage struct
type VoltageStorage struct {
	db *pgxpool.Pool
}

//func NewVoltageStorage(db *pgxpool.Pool) *VoltageStorage
func NewVoltageStorage(db *pgxpool.Pool) *VoltageStorage {
	return &VoltageStorage{db: db}
}

//func (est *VoltageStorage) GetList(ctx context.Context, pg, pgs int, gs1 string, ord int, dsc bool) (models.Voltage_count, error)
func (est *VoltageStorage) GetList(ctx context.Context, pg, pgs int, gs1 string, ord int, dsc bool) (models.Voltage_count, error) {
	dbpool := pgclient.WDB
	gs := models.Voltage{}

	gsc := 0
	err := dbpool.QueryRow(ctx, "SELECT * from func_voltages_cnt($1);", gs1).Scan(&gsc)

	if err != nil {
		log.Println(err.Error(), "func_voltages_cnt")
		return models.Voltage_count{Values: []models.Voltage{}, Count: gsc, Auth: models.Auth{}}, err
	}

	out_arr := make([]models.Voltage, 0,
		func() int {
			if gsc < pgs {
				return gsc
			} else {
				return pgs
			}
		}())

	rows, err := dbpool.Query(ctx, "SELECT * from func_voltages_get($1,$2,$3,$4,$5);", pg, pgs, gs1, ord, dsc)
	if err != nil {
		log.Println(err.Error(), "func_voltages_get")
		return models.Voltage_count{Values: []models.Voltage{}, Count: gsc, Auth: models.Auth{}}, err
	}

	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&gs.Id, &gs.VoltageName, &gs.VoltageValue)
		if err != nil {
			log.Println("failed to scan row:", err)
		}

		out_arr = append(out_arr, gs)
	}

	out_count := models.Voltage_count{Values: out_arr, Count: gsc, Auth: models.Auth{}}
	return out_count, nil
}

//func (est *VoltageStorage) Add(ctx context.Context, a models.Voltage) (int, error)
func (est *VoltageStorage) Add(ctx context.Context, a models.Voltage) (int, error) {
	dbpool := pgclient.WDB
	ai := 0

	err := dbpool.QueryRow(ctx, "SELECT func_voltages_add($1, $2);", a.VoltageName, a.VoltageValue).Scan(&ai)

	if err != nil {
		log.Println("Failed execute func_voltages_add: ", err)
		return 0, err
	}

	return ai, nil
}

//func (est *VoltageStorage) Upd(ctx context.Context, u models.Voltage) (int, error)
func (est *VoltageStorage) Upd(ctx context.Context, u models.Voltage) (int, error) {
	dbpool := pgclient.WDB
	ui := 0

	err := dbpool.QueryRow(ctx, "SELECT func_voltages_upd($1,$2,$3);", u.Id, u.VoltageName, u.VoltageValue).Scan(&ui)

	if err != nil {
		log.Println("Failed execute func_voltages_upd: ", err)
		return 0, err
	}
	return ui, nil
}

//func (est *VoltageStorage) Del(ctx context.Context, d []int) ([]int, error)
func (est *VoltageStorage) Del(ctx context.Context, d []int) ([]int, error) {
	dbpool := pgclient.WDB
	res := []int{}
	i := 0
	for _, id := range d {
		err := dbpool.QueryRow(ctx, "SELECT func_voltages_del($1);", id).Scan(&i)
		res = append(res, i)

		if err != nil {
			log.Println("Failed execute func_voltages_del: ", err)
		}
	}
	return res, nil
}

//func (est *VoltageStorage) GetOne(ctx context.Context, i int) (models.Voltage_count, error)
func (est *VoltageStorage) GetOne(ctx context.Context, i int) (models.Voltage_count, error) {
	dbpool := pgclient.WDB
	out_arr := []models.Voltage{}
	g := models.Voltage{}

	err := dbpool.QueryRow(ctx, "SELECT * from func_voltage_get($1);", i).Scan(&g.Id, &g.VoltageName, &g.VoltageValue)

	if err != nil && err != pgx.ErrNoRows {
		log.Println("Failed execute func_voltage_get: ", err)
		return models.Voltage_count{Values: []models.Voltage{}, Count: 0, Auth: models.Auth{}}, err
	}

	out_arr = append(out_arr, g)

	out_count := models.Voltage_count{Values: out_arr, Count: 1, Auth: models.Auth{}}
	return out_count, nil
}
