package pgsql

import (
	"context"
	"log"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/u2lentaru/billbck/internal/models"
	"github.com/u2lentaru/billbck/pkg/pgclient"
)

//type ObjLineStorage struct
type ObjLineStorage struct {
	db *pgxpool.Pool
}

//func NewObjLineStorage(db *pgxpool.Pool) *ObjLineStorage
func NewObjLineStorage(db *pgxpool.Pool) *ObjLineStorage {
	return &ObjLineStorage{db: db}
}

//func (est *ObjLineStorage) GetList(ctx context.Context, pg, pgs int, gs1, gs2 string, ord int, dsc bool) (models.ObjLine_count, error) {
func (est *ObjLineStorage) GetList(ctx context.Context, pg, pgs int, gs1, gs2 string, ord int, dsc bool) (models.ObjLine_count, error) {
	dbpool := pgclient.WDB
	gs := models.ObjLine{}

	gsc := 0
	err := dbpool.QueryRow(ctx, "SELECT * from func_obj_lines_cnt($1,$2);", gs1, gs2).Scan(&gsc)

	if err != nil {
		log.Println(err.Error(), "func_obj_lines_cnt")
		return models.ObjLine_count{Values: []models.ObjLine{}, Count: gsc, Auth: models.Auth{}}, err
	}

	out_arr := make([]models.ObjLine, 0,
		func() int {
			if gsc < pgs {
				return gsc
			} else {
				return pgs
			}
		}())

	rows, err := dbpool.Query(ctx, "SELECT * from func_obj_lines_get($1,$2,$3,$4,$5,$6);", pg, pgs, gs1, gs2, ord, dsc)
	if err != nil {
		log.Println(err.Error())
		return models.ObjLine_count{Values: []models.ObjLine{}, Count: gsc, Auth: models.Auth{}}, err
	}

	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&gs.Id, &gs.ObjId, &gs.ObjTypeId, &gs.CableResistance.Id, &gs.LineLength, &gs.Startdate, &gs.Enddate,
			&gs.ObjName, &gs.CableResistance.CableResistanceName, &gs.CableResistance.Resistance, &gs.CableResistance.MaterialType)
		if err != nil {
			log.Println("failed to scan row:", err)
		}

		out_arr = append(out_arr, gs)
	}

	out_count := models.ObjLine_count{Values: out_arr, Count: gsc, Auth: models.Auth{}}
	return out_count, nil
}

//func (est *ObjLineStorage) Add(ctx context.Context, ea models.ObjLine) (int, error)
func (est *ObjLineStorage) Add(ctx context.Context, a models.ObjLine) (int, error) {
	dbpool := pgclient.WDB
	ai := 0

	err := dbpool.QueryRow(ctx, "SELECT func_obj_lines_add($1,$2,$3,$4,$5);", a.ObjId, a.ObjTypeId, a.CableResistance.Id, a.LineLength,
		a.Startdate).Scan(&ai)

	if err != nil {
		log.Println("Failed execute func_obj_lines_add: ", err)
		return 0, err
	}

	return ai, nil
}

//func (est *ObjLineStorage) Upd(ctx context.Context, u models.ObjLine) (int, error)
func (est *ObjLineStorage) Upd(ctx context.Context, u models.ObjLine) (int, error) {
	dbpool := pgclient.WDB
	ui := 0

	err := dbpool.QueryRow(ctx, "SELECT func_obj_lines_upd($1,$2,$3,$4,$5,$6,$7);", u.Id, u.ObjId, u.ObjTypeId, u.CableResistance.Id,
		u.LineLength, u.Startdate, u.Enddate).Scan(&ui)

	if err != nil {
		log.Println("Failed execute func_obj_lines_upd: ", err)
		return 0, err
	}
	return ui, nil
}

//func (est *ObjLineStorage) Del(ctx context.Context, d []int) ([]int, error)
func (est *ObjLineStorage) Del(ctx context.Context, d []int) ([]int, error) {
	dbpool := pgclient.WDB
	res := []int{}
	i := 0
	for _, id := range d {
		err := dbpool.QueryRow(ctx, "SELECT func_obj_lines_del($1);", id).Scan(&i)
		res = append(res, i)

		if err != nil {
			log.Println("Failed execute func_obj_lines_del: ", err)
		}
	}
	return res, nil
}

//func (est *ObjLineStorage) GetOne(ctx context.Context, i int) (models.ObjLine_count, error)
func (est *ObjLineStorage) GetOne(ctx context.Context, i int) (models.ObjLine_count, error) {
	dbpool := pgclient.WDB
	out_arr := []models.ObjLine{}
	g := models.ObjLine{}

	err := dbpool.QueryRow(ctx, "SELECT * from func_obj_line_get($1);", i).Scan(&g.Id, &g.ObjId, &g.ObjTypeId, &g.CableResistance.Id,
		&g.LineLength, &g.Startdate, &g.Enddate, &g.ObjName, &g.CableResistance.CableResistanceName, &g.CableResistance.Resistance,
		&g.CableResistance.MaterialType)

	if err != nil && err != pgx.ErrNoRows {
		log.Println("Failed execute func_obj_line_get: ", err)
		return models.ObjLine_count{Values: []models.ObjLine{}, Count: 0, Auth: models.Auth{}}, err
	}

	out_arr = append(out_arr, g)

	out_count := models.ObjLine_count{Values: out_arr, Count: 1, Auth: models.Auth{}}
	return out_count, nil
}

//func (est *ObjLineStorage) GetObj(ctx context.Context, gs1, gs2 string) (models.ObjLine_count, error)
func (est *ObjLineStorage) GetObj(ctx context.Context, gs1, gs2 string) (models.ObjLine_count, error) {
	dbpool := pgclient.WDB
	gs := models.ObjLine{}

	gsc := 0
	err := dbpool.QueryRow(ctx, "SELECT * from func_obj_lines_obj_cnt($1,$2);", gs1, gs2).Scan(&gsc)

	if err != nil {
		log.Println(err.Error(), "func_obj_lines_obj_cnt")
		return models.ObjLine_count{Values: []models.ObjLine{}, Count: gsc, Auth: models.Auth{}}, err
	}

	out_arr := make([]models.ObjLine, 0, gsc)

	rows, err := dbpool.Query(ctx, "SELECT * from func_obj_lines_obj($1,$2);", gs1, gs2)
	if err != nil {
		log.Println(err.Error())
		return models.ObjLine_count{Values: []models.ObjLine{}, Count: gsc, Auth: models.Auth{}}, err
	}

	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&gs.Id, &gs.ObjId, &gs.ObjTypeId, &gs.CableResistance.Id, &gs.LineLength, &gs.Startdate, &gs.Enddate,
			&gs.ObjName, &gs.CableResistance.CableResistanceName, &gs.CableResistance.Resistance, &gs.CableResistance.MaterialType)
		if err != nil {
			log.Println("failed to scan row:", err)
		}

		out_arr = append(out_arr, gs)
	}

	out_count := models.ObjLine_count{Values: out_arr, Count: gsc, Auth: models.Auth{}}
	if err != nil {
		log.Println(err.Error())
		return models.ObjLine_count{}, err
	}

	return out_count, nil
}
