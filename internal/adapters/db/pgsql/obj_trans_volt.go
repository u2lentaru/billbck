package pgsql

import (
	"context"
	"log"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/u2lentaru/billbck/internal/models"
	"github.com/u2lentaru/billbck/pkg/pgclient"
)

//type ObjTransVoltStorage struct
type ObjTransVoltStorage struct {
	db *pgxpool.Pool
}

//func NewObjTransVoltStorage(db *pgxpool.Pool) *ObjTransVoltStorage
func NewObjTransVoltStorage(db *pgxpool.Pool) *ObjTransVoltStorage {
	return &ObjTransVoltStorage{db: db}
}

//func (est *ObjTransVoltStorage) GetList(ctx context.Context, pg, pgs int, gs1, gs2 string, ord int, dsc bool) (models.ObjTransVolt_count, error) {
func (est *ObjTransVoltStorage) GetList(ctx context.Context, pg, pgs int, gs1, gs2 string, ord int, dsc bool) (models.ObjTransVolt_count, error) {
	dbpool := pgclient.WDB
	gs := models.ObjTransVolt{}

	gsc := 0
	err := dbpool.QueryRow(ctx, "SELECT * from func_obj_trans_volt_cnt($1,$2);", gs1, gs2).Scan(&gsc)

	if err != nil {
		log.Println(err.Error(), "func_obj_trans_volt_cnt")
		return models.ObjTransVolt_count{Values: []models.ObjTransVolt{}, Count: gsc, Auth: models.Auth{}}, err
	}

	out_arr := make([]models.ObjTransVolt, 0,
		func() int {
			if gsc < pgs {
				return gsc
			} else {
				return pgs
			}
		}())

	rows, err := dbpool.Query(ctx, "SELECT * from func_obj_trans_volt_get($1,$2,$3,$4,$5,$6);", pg, pgs, gs1, gs2, ord, dsc)
	if err != nil {
		log.Println(err.Error())
		return models.ObjTransVolt_count{Values: []models.ObjTransVolt{}, Count: gsc, Auth: models.Auth{}}, err
	}

	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&gs.Id, &gs.ObjId, &gs.ObjTypeId, &gs.TransVolt.Id, &gs.Startdate, &gs.Enddate, &gs.ObjName,
			&gs.TransVolt.TransVoltName, &gs.TransVolt.TransType.Id, &gs.TransVolt.CheckDate, &gs.TransVolt.NextCheckDate,
			&gs.TransVolt.ProdDate, &gs.TransVolt.Serial1, &gs.TransVolt.Serial2, &gs.TransVolt.Serial3,
			&gs.TransVolt.TransType.TransTypeName, &gs.TransVolt.TransType.Ratio, &gs.TransVolt.TransType.Class,
			&gs.TransVolt.TransType.MaxCurr, &gs.TransVolt.TransType.NomCurr)
		if err != nil {
			log.Println("failed to scan row:", err)
		}

		out_arr = append(out_arr, gs)
	}

	out_count := models.ObjTransVolt_count{Values: out_arr, Count: gsc, Auth: models.Auth{}}
	if err != nil {
		log.Println(err.Error())
		return models.ObjTransVolt_count{}, err
	}

	return out_count, nil
}

//func (est *ObjTransVoltStorage) Add(ctx context.Context, ea models.ObjTransVolt) (int, error)
func (est *ObjTransVoltStorage) Add(ctx context.Context, a models.ObjTransVolt) (int, error) {
	dbpool := pgclient.WDB
	ai := 0

	err := dbpool.QueryRow(context.Background(), "SELECT func_obj_trans_volt_add($1,$2,$3,$4);", a.ObjId, a.ObjTypeId,
		a.TransVolt.Id, a.Startdate).Scan(&ai)

	if err != nil {
		log.Println("Failed execute func_obj_trans_volt_add: ", err)
		return 0, err
	}

	return ai, nil
}

//func (est *ObjTransVoltStorage) Upd(ctx context.Context, u models.ObjTransVolt) (int, error)
func (est *ObjTransVoltStorage) Upd(ctx context.Context, u models.ObjTransVolt) (int, error) {
	dbpool := pgclient.WDB
	ui := 0

	err := dbpool.QueryRow(context.Background(), "SELECT func_obj_trans_volt_upd($1,$2,$3,$4,$5,$6);", u.Id, u.ObjId, u.ObjTypeId,
		u.TransVolt.Id, u.Startdate, u.Enddate).Scan(&ui)

	if err != nil {
		log.Println("Failed execute func_obj_trans_volt_upd: ", err)
		return 0, err
	}
	return ui, nil
}

//func (est *ObjTransVoltStorage) Del(ctx context.Context, d []int) ([]int, error)
func (est *ObjTransVoltStorage) Del(ctx context.Context, d []int) ([]int, error) {
	dbpool := pgclient.WDB
	res := []int{}
	i := 0
	for _, id := range d {
		err := dbpool.QueryRow(context.Background(), "SELECT func_obj_trans_volt_del($1);", id).Scan(&i)
		res = append(res, i)

		if err != nil {
			log.Println("Failed execute func_obj_trans_volt_del: ", err)
		}
	}
	return res, nil
}

//func (est *ObjTransVoltStorage) GetOne(ctx context.Context, i int) (models.ObjTransVolt_count, error)
func (est *ObjTransVoltStorage) GetOne(ctx context.Context, i int) (models.ObjTransVolt_count, error) {
	dbpool := pgclient.WDB
	out_arr := []models.ObjTransVolt{}
	g := models.ObjTransVolt{}

	err := dbpool.QueryRow(context.Background(), "SELECT * from func_obj_trans_volt_getbyid($1);", i).Scan(&g.Id, &g.ObjId,
		&g.ObjTypeId, &g.TransVolt.Id, &g.Startdate, &g.Enddate, &g.ObjName, &g.TransVolt.TransVoltName, &g.TransVolt.TransType.Id,
		&g.TransVolt.CheckDate, &g.TransVolt.NextCheckDate, &g.TransVolt.ProdDate, &g.TransVolt.Serial1, &g.TransVolt.Serial2,
		&g.TransVolt.Serial3, &g.TransVolt.TransType.TransTypeName, &g.TransVolt.TransType.Ratio, &g.TransVolt.TransType.Class,
		&g.TransVolt.TransType.MaxCurr, &g.TransVolt.TransType.NomCurr)

	if err != nil && err != pgx.ErrNoRows {
		log.Println("Failed execute func_obj_trans_volt_getbyid: ", err)
		return models.ObjTransVolt_count{Values: []models.ObjTransVolt{}, Count: 0, Auth: models.Auth{}}, err
	}

	out_arr = append(out_arr, g)

	out_count := models.ObjTransVolt_count{Values: out_arr, Count: 1, Auth: models.Auth{}}
	return out_count, nil
}

//func (est *ObjTransVoltStorage) GetObj(ctx context.Context, gs1, gs2 string) (models.ObjTransVolt_count, error)
func (est *ObjTransVoltStorage) GetObj(ctx context.Context, gs1, gs2 string) (models.ObjTransVolt_count, error) {
	dbpool := pgclient.WDB
	out_arr := []models.ObjTransVolt{}
	g := models.ObjTransVolt{}

	err := dbpool.QueryRow(context.Background(), "SELECT * from func_obj_trans_volt_getbyobj($1,$2);", gs1, gs2).Scan(&g.Id, &g.ObjId,
		&g.ObjTypeId, &g.TransVolt.Id, &g.Startdate, &g.Enddate, &g.ObjName, &g.TransVolt.TransVoltName, &g.TransVolt.TransType.Id,
		&g.TransVolt.CheckDate, &g.TransVolt.NextCheckDate, &g.TransVolt.ProdDate, &g.TransVolt.Serial1, &g.TransVolt.Serial2,
		&g.TransVolt.Serial3, &g.TransVolt.TransType.TransTypeName, &g.TransVolt.TransType.Ratio, &g.TransVolt.TransType.Class,
		&g.TransVolt.TransType.MaxCurr, &g.TransVolt.TransType.NomCurr)

	if err != nil && err != pgx.ErrNoRows {
		log.Println("Failed execute func_obj_trans_volt_getbyobj: ", err)
		return models.ObjTransVolt_count{Values: []models.ObjTransVolt{}, Count: 0, Auth: models.Auth{}}, err
	}

	out_arr = append(out_arr, g)

	out_count := models.ObjTransVolt_count{Values: out_arr, Count: 1, Auth: models.Auth{}}
	return out_count, nil
}
