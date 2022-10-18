package pgsql

import (
	"context"
	"log"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/u2lentaru/billbck/internal/models"
	"github.com/u2lentaru/billbck/pkg/pgclient"
)

//type ObjTransCurrStorage struct
type ObjTransCurrStorage struct {
	db *pgxpool.Pool
}

//func NewObjTransCurrStorage(db *pgxpool.Pool) *ObjTransCurrStorage
func NewObjTransCurrStorage(db *pgxpool.Pool) *ObjTransCurrStorage {
	return &ObjTransCurrStorage{db: db}
}

//func (est *ObjTransCurrStorage) GetList(ctx context.Context, pg, pgs int, gs1, gs2 string, ord int, dsc bool) (models.ObjTransCurr_count, error) {
func (est *ObjTransCurrStorage) GetList(ctx context.Context, pg, pgs int, gs1, gs2 string, ord int, dsc bool) (models.ObjTransCurr_count, error) {
	dbpool := pgclient.WDB
	gs := models.ObjTransCurr{}

	gsc := 0
	err := dbpool.QueryRow(ctx, "SELECT * from func_obj_trans_curr_cnt($1,$2);", gs1, gs2).Scan(&gsc)

	if err != nil {
		log.Println(err.Error(), "func_obj_trans_curr_cnt")
		return models.ObjTransCurr_count{Values: []models.ObjTransCurr{}, Count: gsc, Auth: models.Auth{}}, err
	}

	out_arr := make([]models.ObjTransCurr, 0,
		func() int {
			if gsc < pgs {
				return gsc
			} else {
				return pgs
			}
		}())

	rows, err := dbpool.Query(ctx, "SELECT * from func_obj_trans_curr_get($1,$2,$3,$4,$5,$6);", pg, pgs, gs1, gs2, ord, dsc)
	if err != nil {
		log.Println(err.Error())
		return models.ObjTransCurr_count{Values: []models.ObjTransCurr{}, Count: gsc, Auth: models.Auth{}}, err
	}

	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&gs.Id, &gs.ObjId, &gs.ObjTypeId, &gs.TransCurr.Id, &gs.Startdate, &gs.Enddate, &gs.ObjName,
			&gs.TransCurr.TransCurrName, &gs.TransCurr.TransType.Id, &gs.TransCurr.CheckDate, &gs.TransCurr.NextCheckDate,
			&gs.TransCurr.ProdDate, &gs.TransCurr.Serial1, &gs.TransCurr.Serial2, &gs.TransCurr.Serial3,
			&gs.TransCurr.TransType.TransTypeName, &gs.TransCurr.TransType.Ratio, &gs.TransCurr.TransType.Class,
			&gs.TransCurr.TransType.MaxCurr, &gs.TransCurr.TransType.NomCurr)
		if err != nil {
			log.Println("failed to scan row:", err)
		}

		out_arr = append(out_arr, gs)
	}

	out_count := models.ObjTransCurr_count{Values: out_arr, Count: gsc, Auth: models.Auth{}}
	if err != nil {
		log.Println(err.Error())
		return models.ObjTransCurr_count{}, err
	}

	return out_count, nil
}

//func (est *ObjTransCurrStorage) Add(ctx context.Context, ea models.ObjTransCurr) (int, error)
func (est *ObjTransCurrStorage) Add(ctx context.Context, a models.ObjTransCurr) (int, error) {
	dbpool := pgclient.WDB
	ai := 0

	err := dbpool.QueryRow(ctx, "SELECT func_obj_trans_curr_add($1,$2,$3,$4);", a.ObjId, a.ObjTypeId, a.TransCurr.Id, a.Startdate).Scan(&ai)

	if err != nil {
		log.Println("Failed execute func_obj_trans_curr_add: ", err)
		return 0, err
	}

	return ai, nil
}

//func (est *ObjTransCurrStorage) Upd(ctx context.Context, u models.ObjTransCurr) (int, error)
func (est *ObjTransCurrStorage) Upd(ctx context.Context, u models.ObjTransCurr) (int, error) {
	dbpool := pgclient.WDB
	ui := 0

	err := dbpool.QueryRow(ctx, "SELECT func_obj_trans_curr_upd($1,$2,$3,$4,$5,$6);", u.Id, u.ObjId, u.ObjTypeId, u.TransCurr.Id,
		u.Startdate, u.Enddate).Scan(&ui)

	if err != nil {
		log.Println("Failed execute func_obj_trans_curr_upd: ", err)
		return 0, err
	}
	return ui, nil
}

//func (est *ObjTransCurrStorage) Del(ctx context.Context, d []int) ([]int, error)
func (est *ObjTransCurrStorage) Del(ctx context.Context, d []int) ([]int, error) {
	dbpool := pgclient.WDB
	res := []int{}
	i := 0
	for _, id := range d {
		err := dbpool.QueryRow(ctx, "SELECT func_obj_trans_curr_del($1);", id).Scan(&i)
		res = append(res, i)

		if err != nil {
			log.Println("Failed execute func_obj_trans_curr_del: ", err)
		}
	}
	return res, nil
}

//func (est *ObjTransCurrStorage) GetOne(ctx context.Context, i int) (models.ObjTransCurr_count, error)
func (est *ObjTransCurrStorage) GetOne(ctx context.Context, i int) (models.ObjTransCurr_count, error) {
	dbpool := pgclient.WDB
	out_arr := []models.ObjTransCurr{}
	g := models.ObjTransCurr{}

	err := dbpool.QueryRow(ctx, "SELECT * from func_obj_trans_curr_getbyid($1);", i).Scan(&g.Id, &g.ObjId, &g.ObjTypeId, &g.TransCurr.Id,
		&g.Startdate, &g.Enddate, &g.ObjName, &g.TransCurr.TransCurrName, &g.TransCurr.TransType.Id, &g.TransCurr.CheckDate,
		&g.TransCurr.NextCheckDate, &g.TransCurr.ProdDate, &g.TransCurr.Serial1, &g.TransCurr.Serial2, &g.TransCurr.Serial3,
		&g.TransCurr.TransType.TransTypeName, &g.TransCurr.TransType.Ratio, &g.TransCurr.TransType.Class, &g.TransCurr.TransType.MaxCurr,
		&g.TransCurr.TransType.NomCurr)

	if err != nil && err != pgx.ErrNoRows {
		log.Println("Failed execute func_obj_trans_curr_getbyid: ", err)
		return models.ObjTransCurr_count{Values: []models.ObjTransCurr{}, Count: 0, Auth: models.Auth{}}, err
	}

	out_arr = append(out_arr, g)

	out_count := models.ObjTransCurr_count{Values: out_arr, Count: 1, Auth: models.Auth{}}
	return out_count, nil
}

//func (est *ObjTransCurrStorage) GetObj(ctx context.Context, gs1, gs2 string) (models.ObjTransCurr_count, error)
func (est *ObjTransCurrStorage) GetObj(ctx context.Context, gs1, gs2 string) (models.ObjTransCurr_count, error) {
	dbpool := pgclient.WDB
	out_arr := []models.ObjTransCurr{}
	g := models.ObjTransCurr{}

	err := dbpool.QueryRow(ctx, "SELECT * from func_obj_trans_curr_getbyobj($1,$2);", gs1, gs2).Scan(&g.Id, &g.ObjId,
		&g.ObjTypeId, &g.TransCurr.Id, &g.Startdate, &g.Enddate, &g.ObjName, &g.TransCurr.TransCurrName, &g.TransCurr.TransType.Id,
		&g.TransCurr.CheckDate, &g.TransCurr.NextCheckDate, &g.TransCurr.ProdDate, &g.TransCurr.Serial1, &g.TransCurr.Serial2,
		&g.TransCurr.Serial3, &g.TransCurr.TransType.TransTypeName, &g.TransCurr.TransType.Ratio, &g.TransCurr.TransType.Class,
		&g.TransCurr.TransType.MaxCurr, &g.TransCurr.TransType.NomCurr)

	if err != nil && err != pgx.ErrNoRows {
		log.Println("Failed execute func_obj_trans_curr_getbyobj: ", err)
		return models.ObjTransCurr_count{Values: []models.ObjTransCurr{}, Count: 0, Auth: models.Auth{}}, err
	}

	out_arr = append(out_arr, g)

	out_count := models.ObjTransCurr_count{Values: out_arr, Count: 1, Auth: models.Auth{}}
	return out_count, nil
}
