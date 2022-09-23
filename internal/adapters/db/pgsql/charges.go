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

//type ChargeStorage struct
type ChargeStorage struct {
	db *pgxpool.Pool
}

//func NewChargeStorage(db *pgxpool.Pool) *ChargeStorage
func NewChargeStorage(db *pgxpool.Pool) *ChargeStorage {
	return &ChargeStorage{db: db}
}

//func (est *ChargeStorage) GetList(ctx context.Context, pg, pgs int, gs1, gs2 string, ord int, dsc bool) (models.Charge_count, error) {
func (est *ChargeStorage) GetList(ctx context.Context, pg, pgs int, gs1, gs2 string, ord int, dsc bool) (models.Charge_count, error) {
	dbpool := pgclient.WDB
	gs := models.Charge{}

	gsc := 0
	err := dbpool.QueryRow(ctx, "SELECT * from func_charges_cnt($1,$2);", gs1, utils.NullableString(gs2)).Scan(&gsc)

	if err != nil {
		log.Println(err.Error(), "func_charges_cnt")
		return models.Charge_count{Values: []models.Charge{}, Count: gsc, Auth: models.Auth{}}, err
	}

	out_arr := make([]models.Charge, 0,
		func() int {
			if gsc < pgs {
				return gsc
			} else {
				return pgs
			}
		}())

	rows, err := dbpool.Query(ctx, "SELECT * from func_charges_get($1,$2,$3,$4,$5,$6);", pg, pgs, gs1, utils.NullableString(gs2), ord, dsc)
	if err != nil {
		log.Println(err.Error())
		return models.Charge_count{Values: []models.Charge{}, Count: gsc, Auth: models.Auth{}}, err
	}

	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&gs.Id, &gs.ChargeDate, &gs.Contract.Id, &gs.Object.Id, &gs.ObjTypeId, &gs.Pu.Id, &gs.ChargeType.Id, &gs.Qty,
			&gs.TransLoss, &gs.Lineloss, &gs.Startdate, &gs.Enddate, &gs.Contract.ContractNumber, &gs.Object.ObjectName,
			&gs.Pu.PuNumber, &gs.ChargeType.ChargeTypeName)
		if err != nil {
			log.Println("failed to scan row:", err)
		}

		out_arr = append(out_arr, gs)
	}

	out_count := models.Charge_count{Values: out_arr, Count: gsc, Auth: models.Auth{}}
	if err != nil {
		log.Println(err.Error())
		return models.Charge_count{}, err
	}

	return out_count, nil
}

//func (est *ChargeStorage) Add(ctx context.Context, ea models.Charge) (int, error)
func (est *ChargeStorage) Add(ctx context.Context, a models.Charge) (int, error) {
	dbpool := pgclient.WDB
	ai := 0

	err := dbpool.QueryRow(context.Background(), "SELECT func_charges_add($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11);", a.ChargeDate, a.Contract.Id,
		a.Object.Id, a.ObjTypeId, a.Pu.Id, a.ChargeType.Id, a.Qty, a.TransLoss, a.Lineloss, a.Startdate, a.Enddate).Scan(&ai)

	if err != nil {
		log.Println("Failed execute func_charges_add: ", err)
		return 0, err
	}

	return ai, nil
}

//func (est *ChargeStorage) Upd(ctx context.Context, eu models.Charge) (int, error)
func (est *ChargeStorage) Upd(ctx context.Context, u models.Charge) (int, error) {
	dbpool := pgclient.WDB
	ui := 0

	err := dbpool.QueryRow(context.Background(), "SELECT func_charges_upd($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12);", u.Id, u.ChargeDate,
		u.Contract.Id, u.Object.Id, u.ObjTypeId, u.Pu.Id, u.ChargeType.Id, u.Qty, u.TransLoss, u.Lineloss, u.Startdate, u.Enddate).Scan(&ui)

	if err != nil {
		log.Println("Failed execute func_charges_upd: ", err)
		return 0, err
	}
	return ui, nil
}

//func (est *ChargeStorage) Del(ctx context.Context, ed []int) ([]int, error)
func (est *ChargeStorage) Del(ctx context.Context, d []int) ([]int, error) {
	dbpool := pgclient.WDB
	res := []int{}
	i := 0
	for _, id := range d {
		err := dbpool.QueryRow(context.Background(), "SELECT func_charges_del($1);", id).Scan(&i)
		res = append(res, i)

		if err != nil {
			log.Println("Failed execute func_charges_del: ", err)
		}
	}
	return res, nil
}

//func (est *ChargeStorage) GetOne(ctx context.Context, i int) (models.Charge_count, error)
func (est *ChargeStorage) GetOne(ctx context.Context, i int) (models.Charge_count, error) {
	dbpool := pgclient.WDB
	out_arr := []models.Charge{}
	g := models.Charge{}

	err := dbpool.QueryRow(context.Background(), "SELECT * from func_charge_get($1);", i).Scan(&g.Id, &g.ChargeDate, &g.Contract.Id,
		&g.Object.Id, &g.ObjTypeId, &g.Pu.Id, &g.ChargeType.Id, &g.Qty, &g.TransLoss, &g.Lineloss, &g.Startdate, &g.Enddate,
		&g.Contract.ContractNumber, &g.Object.ObjectName, &g.Pu.PuNumber, &g.ChargeType.ChargeTypeName)

	if err != nil && err != pgx.ErrNoRows {
		log.Println("Failed execute from func_charge_get: ", err)
		return models.Charge_count{Values: []models.Charge{}, Count: 0, Auth: models.Auth{}}, err
	}

	out_arr = append(out_arr, g)

	out_count := models.Charge_count{Values: out_arr, Count: 1, Auth: models.Auth{}}
	return out_count, nil
}

//func (est *ChargeStorage) ChargeRun(ctx context.Context, i int) (int, error)
func (est *ChargeStorage) ChargeRun(ctx context.Context, i int) (int, error) {
	dbpool := pgclient.WDB
	pr := 0

	err := dbpool.QueryRow(context.Background(), "SELECT * from func_charges_run($1);", i).Scan(&pr)

	if err != nil {
		log.Println("Failed execute func_charges_run: ", err)
	}

	return pr, nil
}
