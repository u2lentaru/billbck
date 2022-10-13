package pgsql

import (
	"context"
	"log"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/u2lentaru/billbck/internal/models"
	"github.com/u2lentaru/billbck/pkg/pgclient"
)

//type SubPuStorage struct
type SubPuStorage struct {
	db *pgxpool.Pool
}

//func NewSubPuStorage(db *pgxpool.Pool) *SubPuStorage
func NewSubPuStorage(db *pgxpool.Pool) *SubPuStorage {
	return &SubPuStorage{db: db}
}

//func (est *SubPuStorage) GetList(ctx context.Context, pg, pgs, gs1, ord int, dsc bool) (models.SubPu_count, error)
func (est *SubPuStorage) GetList(ctx context.Context, pg, pgs, gs1, ord int, dsc bool) (models.Pu_count, error) {
	dbpool := pgclient.WDB
	gs := models.Pu{}

	gsc := 0
	err := dbpool.QueryRow(ctx, "SELECT * from func_sub_pu_cnt($1);", gs1).Scan(&gsc)

	if err != nil {
		log.Println(err.Error(), "func_sub_pu_cnt")
		return models.Pu_count{Values: []models.Pu{}, Count: gsc, Auth: models.Auth{}}, err
	}

	out_arr := make([]models.Pu, 0,
		func() int {
			if gsc < pgs {
				return gsc
			} else {
				return pgs
			}
		}())

	rows, err := dbpool.Query(ctx, "SELECT * from func_sub_pu_get($1,$2,$3::int,$4,$5);", pg, pgs, gs1, ord, dsc)
	if err != nil {
		log.Println(err.Error(), "func_sub_pu_get")
		return models.Pu_count{Values: []models.Pu{}, Count: gsc, Auth: models.Auth{}}, err
	}

	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&gs.Id, &gs.Startdate, &gs.Enddate, &gs.PuType.Id, &gs.PuType.PuTypeName, &gs.PuNumber, &gs.InstallDate,
			&gs.CheckInterval, &gs.InitialValue, &gs.DevStopped, &gs.Object.Id, &gs.PuObjectType, &gs.Object.ObjectName, &gs.Object.House.Id,
			&gs.Object.House.HouseNumber, &gs.Object.FlatNumber, &gs.Object.House.BuildingNumber, &gs.Object.RegQty,
			&gs.Object.House.Street.Id, &gs.Object.House.Street.StreetName, &gs.Object.House.Street.City.CityName,
			&gs.Object.House.BuildingType.BuildingTypeName, &gs.Object.House.Street.City.Id, &gs.Object.House.Street.Created,
			&gs.Object.House.Street.Closed, &gs.Pid)
		if err != nil {
			log.Println("failed to scan row:", err)
		}

		out_arr = append(out_arr, gs)
	}

	out_count := models.Pu_count{Values: out_arr, Count: gsc, Auth: models.Auth{}}
	if err != nil {
		log.Println(err.Error())
		return models.Pu_count{}, err
	}

	return out_count, nil
}

//func (est *SubPuStorage) Add(ctx context.Context, a models.SubPu) (int, error)
func (est *SubPuStorage) Add(ctx context.Context, a models.SubPu) (int, error) {
	dbpool := pgclient.WDB
	ai := 0

	err := dbpool.QueryRow(ctx, "SELECT func_sub_pu_add($1,$2);", a.ParId, a.SubId).Scan(&ai)

	if err != nil {
		log.Println("Failed execute func_sub_pu_add: ", err)
		return 0, err
	}

	return ai, nil
}

//func (est *SubPuStorage) Upd(ctx context.Context, u models.SubPu) (int, error)
func (est *SubPuStorage) Upd(ctx context.Context, u models.SubPu) (int, error) {
	dbpool := pgclient.WDB
	ui := 0

	err := dbpool.QueryRow(ctx, "SELECT func_sub_pu_upd($1,$2,$3);", u.Id, u.ParId, u.SubId).Scan(&ui)

	if err != nil {
		log.Println("Failed execute func_sub_pu_upd: ", err)
		return 0, err
	}
	return ui, nil
}

//func (est *SubPuStorage) Del(ctx context.Context, d []int) ([]int, error)
func (est *SubPuStorage) Del(ctx context.Context, d []int) ([]int, error) {
	dbpool := pgclient.WDB
	res := []int{}
	i := 0
	for _, id := range d {
		err := dbpool.QueryRow(ctx, "SELECT func_sub_pu_del($1);", id).Scan(&i)
		res = append(res, i)

		if err != nil {
			log.Println("Failed execute func_sub_pu_del: ", err)
		}
	}
	return res, nil
}

//func (est *SubPuStorage) GetOne(ctx context.Context, i int) (models.SubPu_count, error)
func (est *SubPuStorage) GetOne(ctx context.Context, i int) (models.SubPu_count, error) {
	dbpool := pgclient.WDB
	out_arr := []models.SubPu{}
	g := models.SubPu{}

	err := dbpool.QueryRow(ctx, "SELECT * from func_sub_pu_getbyid($1);", i).Scan(&g.Id, &g.ParId, &g.SubId)

	if err != nil && err != pgx.ErrNoRows {
		log.Println("Failed execute func_sub_pu_getbyid: ", err)
		return models.SubPu_count{Values: []models.SubPu{}, Count: 0, Auth: models.Auth{}}, err
	}

	out_arr = append(out_arr, g)

	out_count := models.SubPu_count{Values: out_arr, Count: 1, Auth: models.Auth{}}
	return out_count, nil
}

//func (est *SubPuStorage) GetPrl(ctx context.Context, pg, pgs, gs1, gs2, ord int, dsc bool) (models.SubPu_count, error)
func (est *SubPuStorage) GetPrl(ctx context.Context, pg, pgs, gs1, gs2, ord int, dsc bool) (models.Pu_count, error) {
	dbpool := pgclient.WDB
	gs := models.Pu{}

	gsc := 0
	err := dbpool.QueryRow(ctx, "SELECT * from func_sub_pu_prl_cnt($1,$2);", gs1, gs2).Scan(&gsc)

	if err != nil {
		log.Println(err.Error(), "func_sub_pu_prl_cnt")
		return models.Pu_count{Values: []models.Pu{}, Count: gsc, Auth: models.Auth{}}, err
	}

	out_arr := make([]models.Pu, 0,
		func() int {
			if gsc < pgs {
				return gsc
			} else {
				return pgs
			}
		}())

	rows, err := dbpool.Query(ctx, "SELECT * from func_sub_pu_prl($1,$2,$3::int,$4::int,$5,$6);", pg, pgs, gs1, gs2, ord, dsc)
	if err != nil {
		log.Println(err.Error(), "func_sub_pu_prl")
		return models.Pu_count{Values: []models.Pu{}, Count: gsc, Auth: models.Auth{}}, err
	}

	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&gs.Id, &gs.Startdate, &gs.Enddate, &gs.PuType.Id, &gs.PuType.PuTypeName, &gs.PuNumber, &gs.InstallDate,
			&gs.CheckInterval, &gs.InitialValue, &gs.DevStopped, &gs.Object.Id, &gs.PuObjectType, &gs.Object.ObjectName, &gs.Object.House.Id,
			&gs.Object.House.HouseNumber, &gs.Object.FlatNumber, &gs.Object.House.BuildingNumber, &gs.Object.RegQty,
			&gs.Object.House.Street.Id, &gs.Object.House.Street.StreetName, &gs.Object.House.Street.City.CityName,
			&gs.Object.House.BuildingType.BuildingTypeName, &gs.Object.House.Street.City.Id, &gs.Object.House.Street.Created,
			&gs.Object.House.Street.Closed, &gs.Pid)
		if err != nil {
			log.Println("failed to scan row:", err)
		}

		out_arr = append(out_arr, gs)
	}

	out_count := models.Pu_count{Values: out_arr, Count: gsc, Auth: models.Auth{}}
	if err != nil {
		log.Println(err.Error())
		return models.Pu_count{}, err
	}

	return out_count, nil
}
