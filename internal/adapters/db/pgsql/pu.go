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

//type PuStorage struct
type PuStorage struct {
	db *pgxpool.Pool
}

//func NewPuStorage(db *pgxpool.Pool) *PuStorage
func NewPuStorage(db *pgxpool.Pool) *PuStorage {
	return &PuStorage{db: db}
}

//func (est *PuStorage) GetList(ctx context.Context, pg, pgs int, gs1, gs2, gs3, gs4, gs5, gs6, gs7 string, ord int, dsc bool) (models.Pu_count, error)
func (est *PuStorage) GetList(ctx context.Context, pg, pgs int, gs1, gs2, gs3, gs4, gs5, gs6, gs7 string, ord int, dsc bool) (models.Pu_count, error) {
	dbpool := pgclient.WDB
	gs := models.Pu{}

	gsc := 0
	err := dbpool.QueryRow(ctx, "SELECT * from func_pu_cnt($1,$2,$3,$4,$5,$6,$7);", gs1, gs2, gs3, utils.NullableString(gs4), utils.NullableString(gs5),
		utils.NullableString(gs6), utils.NullableString(gs7)).Scan(&gsc)

	if err != nil {
		log.Println(err.Error(), "func_pu_cnt")
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

	rows, err := dbpool.Query(ctx, "SELECT * from func_pu_get($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11);", pg, pgs, gs1, gs2, gs3, utils.NullableString(gs4),
		utils.NullableString(gs5), utils.NullableString(gs6), utils.NullableString(gs7), ord, dsc)
	if err != nil {
		log.Println(err.Error())
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

//func (est *PuStorage) Add(ctx context.Context, ea models.Pu) (int, error)
func (est *PuStorage) Add(ctx context.Context, a models.Pu) (int, error) {
	dbpool := pgclient.WDB
	ai := 0

	err := dbpool.QueryRow(ctx, "SELECT func_pu_add($1,$2,$3,$4,$5,$6,$7,$8,$9,$10);",
		a.Object.Id, a.PuObjectType, a.PuType.Id, a.PuNumber, a.InstallDate, a.CheckInterval, a.InitialValue, a.DevStopped, a.Startdate,
		a.Pid).Scan(&ai)

	if err != nil {
		log.Println("Failed execute func_pu_add: ", err)
		return 0, err
	}

	return ai, nil
}

//func (est *PuStorage) Upd(ctx context.Context, eu models.Pu) (int, error)
func (est *PuStorage) Upd(ctx context.Context, u models.Pu) (int, error) {
	dbpool := pgclient.WDB
	ui := 0

	err := dbpool.QueryRow(ctx, "SELECT func_pu_upd($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12);", u.Id, u.Object.Id, u.PuObjectType,
		u.PuType.Id, u.PuNumber, u.InstallDate, u.CheckInterval, u.InitialValue, u.DevStopped, u.Startdate, u.Enddate, u.Pid).Scan(&ui)

	if err != nil {
		log.Println("Failed execute func_pu_upd: ", err)
		return 0, err
	}
	return ui, nil
}

//func (est *PuStorage) Del(ctx context.Context, d []int) ([]int, error)
func (est *PuStorage) Del(ctx context.Context, d []int) ([]int, error) {
	dbpool := pgclient.WDB
	res := []int{}
	i := 0
	for _, id := range d {
		err := dbpool.QueryRow(ctx, "SELECT func_pu_del($1);", id).Scan(&i)
		res = append(res, i)

		if err != nil {
			log.Println("Failed execute func_pu_del: ", err)
		}
	}
	return res, nil
}

//func (est *PuStorage) GetOne(ctx context.Context, i int) (models.Pu_count, error)
func (est *PuStorage) GetOne(ctx context.Context, i int) (models.Pu_count, error) {
	dbpool := pgclient.WDB
	out_arr := []models.Pu{}
	g := models.Pu{}

	err := dbpool.QueryRow(ctx, "SELECT * from func_pu_getbyid($1);", i).Scan(&g.Id, &g.Startdate, &g.Enddate,
		&g.PuType.Id, &g.PuType.PuTypeName, &g.PuNumber, &g.InstallDate, &g.CheckInterval, &g.InitialValue, &g.DevStopped, &g.Object.Id,
		&g.PuObjectType, &g.Object.ObjectName, &g.Object.House.Id, &g.Object.House.HouseNumber, &g.Object.FlatNumber,
		&g.Object.House.BuildingNumber, &g.Object.RegQty, &g.Object.House.Street.Id, &g.Object.House.Street.StreetName,
		&g.Object.House.Street.City.CityName, &g.Object.House.BuildingType.BuildingTypeName, &g.Object.House.Street.City.Id,
		&g.Object.House.Street.Created, &g.Object.House.Street.Closed, &g.Pid)

	if err != nil && err != pgx.ErrNoRows {
		log.Println("Failed execute from func_pu_getbyid: ", err)
		return models.Pu_count{Values: []models.Pu{}, Count: 0, Auth: models.Auth{}}, err
	}

	out_arr = append(out_arr, g)

	out_count := models.Pu_count{Values: out_arr, Count: 1, Auth: models.Auth{}}
	return out_count, nil
}

//func (est *PuStorage) GetObj(ctx context.Context,  gs1, gs2 string) (models.Pu_count, error)
func (est *PuStorage) GetObj(ctx context.Context, gs1, gs2 string) (models.Pu_count, error) {
	dbpool := pgclient.WDB
	out_arr := []models.Pu{}
	gs := models.Pu{}

	rows, err := dbpool.Query(ctx, "SELECT * from func_pu_obj($1,$2);", gs1, gs2)
	if err != nil {
		log.Println("Failed execute from func_pu_obj: ", err)
		return models.Pu_count{Values: []models.Pu{}, Count: 0, Auth: models.Auth{}}, err
	}

	defer rows.Close()

	gsc := 0

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
		gsc++
	}

	out_count := models.Pu_count{Values: out_arr, Count: 1, Auth: models.Auth{}}
	return out_count, nil
}
