package pgsql

import (
	"context"
	"log"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/u2lentaru/billbck/internal/models"
	"github.com/u2lentaru/billbck/pkg/pgclient"
)

//type RpStorage struct
type RpStorage struct {
	db *pgxpool.Pool
}

//func NewRpStorage(db *pgxpool.Pool) *RpStorage
func NewRpStorage(db *pgxpool.Pool) *RpStorage {
	return &RpStorage{db: db}
}

//func (est *RpStorage) GetList(ctx context.Context, pg, pgs int, gs1, gs2 string, ord int, dsc bool) (models.Rp_count, error) {
func (est *RpStorage) GetList(ctx context.Context, pg, pgs int, gs1, gs2 string, ord int, dsc bool) (models.Rp_count, error) {
	dbpool := pgclient.WDB
	gs := models.Rp{}

	gsc := 0
	err := dbpool.QueryRow(ctx, "SELECT * from func_rp_cnt($1,$2);", gs1, gs2).Scan(&gsc)

	if err != nil {
		log.Println(err.Error(), "func_rp_cnt")
		return models.Rp_count{Values: []models.Rp{}, Count: gsc, Auth: models.Auth{}}, err
	}

	out_arr := make([]models.Rp, 0,
		func() int {
			if gsc < pgs {
				return gsc
			} else {
				return pgs
			}
		}())

	rows, err := dbpool.Query(ctx, "SELECT * from func_rp_get($1,$2,$3,$4,$5,$6);", pg, pgs, gs1, gs2, ord, dsc)
	if err != nil {
		log.Println(err.Error(), "func_rp_get")
		return models.Rp_count{Values: []models.Rp{}, Count: gsc, Auth: models.Auth{}}, err
	}

	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&gs.Id, &gs.RpName, &gs.InvNumber, &gs.InputVoltage.Id, &gs.OutputVoltage1.Id, &gs.OutputVoltage2.Id,
			&gs.Tp.Id, &gs.InputVoltage.VoltageName, &gs.InputVoltage.VoltageValue, &gs.OutputVoltage1.VoltageName,
			&gs.OutputVoltage1.VoltageValue, &gs.OutputVoltage2.VoltageName, &gs.OutputVoltage2.VoltageValue,
			&gs.Tp.TpName, &gs.Tp.GRp.Id, &gs.Tp.GRp.GRpName)
		if err != nil {
			log.Println("failed to scan row:", err)
		}

		out_arr = append(out_arr, gs)
	}

	out_count := models.Rp_count{Values: out_arr, Count: gsc, Auth: models.Auth{}}
	return out_count, nil
}

//func (est *RpStorage) Add(ctx context.Context, ea models.Rp) (int, error)
func (est *RpStorage) Add(ctx context.Context, a models.Rp) (int, error) {
	dbpool := pgclient.WDB
	ai := 0

	err := dbpool.QueryRow(ctx, "SELECT func_rp_add($1,$2,$3,$4,$5,$6);", a.RpName, a.InvNumber, a.InputVoltage.Id,
		a.OutputVoltage1.Id, a.OutputVoltage2.Id, a.Tp.Id).Scan(&ai)

	if err != nil {
		log.Println("Failed execute func_rp_add: ", err)
		return 0, err
	}

	return ai, nil
}

//func (est *RpStorage) Upd(ctx context.Context, eu models.Rp) (int, error)
func (est *RpStorage) Upd(ctx context.Context, u models.Rp) (int, error) {
	dbpool := pgclient.WDB
	ui := 0

	err := dbpool.QueryRow(ctx, "SELECT func_rp_upd($1,$2,$3,$4,$5,$6,$7);", u.Id, u.RpName, u.InvNumber,
		u.InputVoltage.Id, u.OutputVoltage1.Id, u.OutputVoltage2.Id, u.Tp.Id).Scan(&ui)

	if err != nil {
		log.Println("Failed execute func_rp_upd: ", err)
		return 0, err
	}
	return ui, nil
}

//func (est *RpStorage) Del(ctx context.Context, ed []int) ([]int, error)
func (est *RpStorage) Del(ctx context.Context, d []int) ([]int, error) {
	dbpool := pgclient.WDB
	res := []int{}
	i := 0
	for _, id := range d {
		err := dbpool.QueryRow(ctx, "SELECT func_rp_del($1);", id).Scan(&i)
		res = append(res, i)

		if err != nil {
			log.Println("Failed execute func_rp_del: ", err)
		}
	}
	return res, nil
}

//func (est *RpStorage) GetOne(ctx context.Context, i int) (models.Rp_count, error)
func (est *RpStorage) GetOne(ctx context.Context, i int) (models.Rp_count, error) {
	dbpool := pgclient.WDB
	out_arr := []models.Rp{}
	g := models.Rp{}

	err := dbpool.QueryRow(ctx, "SELECT * from func_rp_getbyid($1);", i).Scan(&g.Id, &g.RpName, &g.InvNumber,
		&g.InputVoltage.Id, &g.OutputVoltage1.Id, &g.OutputVoltage2.Id, &g.Tp.Id, &g.InputVoltage.VoltageName,
		&g.InputVoltage.VoltageValue, &g.OutputVoltage1.VoltageName, &g.OutputVoltage1.VoltageValue, &g.OutputVoltage2.VoltageName,
		&g.OutputVoltage2.VoltageValue, &g.Tp.TpName, &g.Tp.GRp.Id, &g.Tp.GRp.GRpName)

	if err != nil && err != pgx.ErrNoRows {
		log.Println("Failed execute func_rp_getbyid: ", err)
		return models.Rp_count{Values: []models.Rp{}, Count: 0, Auth: models.Auth{}}, err
	}

	out_arr = append(out_arr, g)

	out_count := models.Rp_count{Values: out_arr, Count: 1, Auth: models.Auth{}}
	return out_count, nil

}
