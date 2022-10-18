package pgsql

import (
	"context"
	"log"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/u2lentaru/billbck/internal/models"
	"github.com/u2lentaru/billbck/pkg/pgclient"
)

//type ObjTransPwrStorage struct
type ObjTransPwrStorage struct {
	db *pgxpool.Pool
}

//func NewObjTransPwrStorage(db *pgxpool.Pool) *ObjTransPwrStorage
func NewObjTransPwrStorage(db *pgxpool.Pool) *ObjTransPwrStorage {
	return &ObjTransPwrStorage{db: db}
}

//func (est *ObjTransPwrStorage) GetList(ctx context.Context, pg, pgs int, gs1, gs2 string, ord int, dsc bool) (models.ObjTransPwr_count, error) {
func (est *ObjTransPwrStorage) GetList(ctx context.Context, pg, pgs int, gs1, gs2 string, ord int, dsc bool) (models.ObjTransPwr_count, error) {
	dbpool := pgclient.WDB
	gs := models.ObjTransPwr{}

	gsc := 0
	err := dbpool.QueryRow(ctx, "SELECT * from func_obj_trans_pwr_cnt($1,$2);", gs1, gs2).Scan(&gsc)

	if err != nil {
		log.Println(err.Error(), "func_obj_trans_pwr_cnt")
		return models.ObjTransPwr_count{Values: []models.ObjTransPwr{}, Count: gsc, Auth: models.Auth{}}, err
	}

	out_arr := make([]models.ObjTransPwr, 0,
		func() int {
			if gsc < pgs {
				return gsc
			} else {
				return pgs
			}
		}())

	rows, err := dbpool.Query(ctx, "SELECT * from func_obj_trans_pwr_get($1,$2,$3,$4,$5,$6);", pg, pgs, gs1, gs2, ord, dsc)
	if err != nil {
		log.Println(err.Error())
		return models.ObjTransPwr_count{Values: []models.ObjTransPwr{}, Count: gsc, Auth: models.Auth{}}, err
	}

	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&gs.Id, &gs.ObjId, &gs.ObjTypeId, &gs.TransPwr.Id, &gs.Startdate, &gs.Enddate, &gs.ObjName,
			&gs.TransPwr.TransPwrName, &gs.TransPwr.TransPwrType.Id, &gs.TransPwr.TransPwrType.TransPwrTypeName,
			&gs.TransPwr.TransPwrType.ShortCircuitPower, &gs.TransPwr.TransPwrType.IdlingLossPower,
			&gs.TransPwr.TransPwrType.NominalPower)
		if err != nil {
			log.Println("failed to scan row:", err)
		}

		out_arr = append(out_arr, gs)
	}

	out_count := models.ObjTransPwr_count{Values: out_arr, Count: gsc, Auth: models.Auth{}}
	if err != nil {
		log.Println(err.Error())
		return models.ObjTransPwr_count{}, err
	}

	return out_count, nil
}

//func (est *ObjTransPwrStorage) Add(ctx context.Context, ea models.ObjTransPwr) (int, error)
func (est *ObjTransPwrStorage) Add(ctx context.Context, a models.ObjTransPwr) (int, error) {
	dbpool := pgclient.WDB
	ai := 0

	err := dbpool.QueryRow(ctx, "SELECT func_obj_trans_pwr_add($1,$2,$3,$4);", a.ObjId, a.ObjTypeId, a.TransPwr.Id, a.Startdate).Scan(&ai)

	if err != nil {
		log.Println("Failed execute func_obj_trans_pwr_add: ", err)
		return 0, err
	}

	return ai, nil
}

//func (est *ObjTransPwrStorage) Upd(ctx context.Context, u models.ObjTransPwr) (int, error)
func (est *ObjTransPwrStorage) Upd(ctx context.Context, u models.ObjTransPwr) (int, error) {
	dbpool := pgclient.WDB
	ui := 0

	err := dbpool.QueryRow(ctx, "SELECT func_obj_trans_pwr_upd($1,$2,$3,$4,$5,$6);", u.Id, u.ObjId, u.ObjTypeId, u.TransPwr.Id, u.Startdate,
		u.Enddate).Scan(&ui)

	if err != nil {
		log.Println("Failed execute func_obj_trans_pwr_upd: ", err)
		return 0, err
	}
	return ui, nil
}

//func (est *ObjTransPwrStorage) Del(ctx context.Context, d []int) ([]int, error)
func (est *ObjTransPwrStorage) Del(ctx context.Context, d []int) ([]int, error) {
	dbpool := pgclient.WDB
	res := []int{}
	i := 0
	for _, id := range d {
		err := dbpool.QueryRow(ctx, "SELECT func_obj_trans_pwr_del($1);", id).Scan(&i)
		res = append(res, i)

		if err != nil {
			log.Println("Failed execute func_obj_trans_pwr_del: ", err)
		}
	}
	return res, nil
}

//func (est *ObjTransPwrStorage) GetOne(ctx context.Context, i int) (models.ObjTransPwr_count, error)
func (est *ObjTransPwrStorage) GetOne(ctx context.Context, i int) (models.ObjTransPwr_count, error) {
	dbpool := pgclient.WDB
	out_arr := []models.ObjTransPwr{}
	g := models.ObjTransPwr{}

	err := dbpool.QueryRow(ctx, "SELECT * from func_obj_trans_pwr_getbyid($1);", i).Scan(&g.Id, &g.ObjId, &g.ObjTypeId, &g.TransPwr.Id,
		&g.Startdate, &g.Enddate, &g.ObjName, &g.TransPwr.TransPwrName, &g.TransPwr.TransPwrType.Id, &g.TransPwr.TransPwrType.TransPwrTypeName,
		&g.TransPwr.TransPwrType.ShortCircuitPower, &g.TransPwr.TransPwrType.IdlingLossPower, &g.TransPwr.TransPwrType.NominalPower)

	if err != nil && err != pgx.ErrNoRows {
		log.Println("Failed execute func_obj_trans_pwr_getbyid: ", err)
		return models.ObjTransPwr_count{Values: []models.ObjTransPwr{}, Count: 0, Auth: models.Auth{}}, err
	}

	out_arr = append(out_arr, g)

	out_count := models.ObjTransPwr_count{Values: out_arr, Count: 1, Auth: models.Auth{}}
	return out_count, nil
}

//func (est *ObjTransPwrStorage) GetObj(ctx context.Context, gs1, gs2 string) (models.ObjTransPwr_count, error)
func (est *ObjTransPwrStorage) GetObj(ctx context.Context, gs1, gs2 string) (models.ObjTransPwr_count, error) {
	dbpool := pgclient.WDB
	out_arr := []models.ObjTransPwr{}
	g := models.ObjTransPwr{}

	err := dbpool.QueryRow(ctx, "SELECT * from func_obj_trans_pwr_getbyobj($1,$2);", gs1, gs2).Scan(&g.Id, &g.ObjId, &g.ObjTypeId, &g.TransPwr.Id,
		&g.Startdate, &g.Enddate, &g.ObjName, &g.TransPwr.TransPwrName, &g.TransPwr.TransPwrType.Id, &g.TransPwr.TransPwrType.TransPwrTypeName,
		&g.TransPwr.TransPwrType.ShortCircuitPower, &g.TransPwr.TransPwrType.IdlingLossPower, &g.TransPwr.TransPwrType.NominalPower)

	if err != nil && err != pgx.ErrNoRows {
		log.Println("Failed execute func_obj_trans_pwr_getbyobj: ", err)
		return models.ObjTransPwr_count{Values: []models.ObjTransPwr{}, Count: 0, Auth: models.Auth{}}, err
	}

	out_arr = append(out_arr, g)

	out_count := models.ObjTransPwr_count{Values: out_arr, Count: 1, Auth: models.Auth{}}
	return out_count, nil
}
