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

//type ActStorage struct
type ActStorage struct {
	db *pgxpool.Pool
}

//func NewActStorage(db *pgxpool.Pool) *ActStorage
func NewActStorage(db *pgxpool.Pool) *ActStorage {
	return &ActStorage{db: db}
}

//func (est *ActStorage) GetList(ctx context.Context, pg, pgs int, gs1, gs2 string, gs3, ord int, dsc bool) (models.Act_count, error)
func (est *ActStorage) GetList(ctx context.Context, pg, pgs int, gs1, gs2 string, gs3, ord int, dsc bool) (models.Act_count, error) {
	dbpool := pgclient.WDB
	gs := models.Act{}

	gsc := 0
	err := dbpool.QueryRow(ctx, "SELECT * from func_acts_cnt($1,$2,$3);", gs1, gs2, gs3).Scan(&gsc)

	if err != nil {
		log.Println(err.Error(), "func_acts_cnt")
		return models.Act_count{Values: []models.Act{}, Count: gsc, Auth: models.Auth{}}, err
	}

	out_arr := make([]models.Act, 0,
		func() int {
			if gsc < pgs {
				return gsc
			} else {
				return pgs
			}
		}())

	rows, err := dbpool.Query(ctx, "SELECT * from func_acts_get($1,$2,$3,$4,$5,$6,$7);", pg, pgs, gs1, gs2, gs3, ord, dsc)
	if err != nil {
		log.Println(err.Error())
		return models.Act_count{Values: []models.Act{}, Count: gsc, Auth: models.Auth{}}, err
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&gs.Id, &gs.ActType.Id, &gs.ActNumber, &gs.ActDate, &gs.Object.Id, &gs.Staff.Id, &gs.Notes, &gs.Activated,
			&gs.ActType.ActTypeName, &gs.Object.ObjectName, &gs.Object.FlatNumber, &gs.Object.RegQty, &gs.Object.House.Street.StreetName,
			&gs.Object.House.Street.City.CityName, &gs.Object.House.HouseNumber, &gs.Object.House.BuildingNumber,
			&gs.Object.TariffGroup.TariffGroupName, &gs.Staff.StaffName)

		if err != nil {
			log.Println("failed to scan row:", err)
		}

		out_arr = append(out_arr, gs)
	}

	out_count := models.Act_count{Values: out_arr, Count: gsc, Auth: models.Auth{}}
	if err != nil {
		log.Println(err.Error())
		return models.Act_count{}, err
	}

	return out_count, nil
}

//func (est *ActStorage) Add(ctx context.Context, a models.Act) (int, error)
func (est *ActStorage) Add(ctx context.Context, a models.Act) (int, error) {
	dbpool := pgclient.WDB
	ai := 0

	err := dbpool.QueryRow(context.Background(), "SELECT func_acts_add($1,$2,$3,$4,$5,$6);",
		a.ActType.Id, a.ActNumber, a.ActDate, a.Object.Id, a.Staff.Id, a.Notes).Scan(&ai)

	if err != nil {
		log.Println("Failed execute func_acts_add: ", err)
		return 0, err
	}

	return ai, nil
}

//func (est *ActStorage) Upd(ctx context.Context, u models.Act) (int, error)
func (est *ActStorage) Upd(ctx context.Context, u models.Act) (int, error) {
	dbpool := pgclient.WDB
	ui := 0

	err := dbpool.QueryRow(context.Background(), "SELECT func_acts_upd($1,$2,$3,$4,$5,$6,$7);", u.Id, u.ActType.Id, u.ActNumber,
		u.ActDate, u.Object.Id, u.Staff.Id, u.Notes).Scan(&ui)

	if err != nil {
		log.Println("Failed execute func_acts_upd: ", err)
		return 0, err
	}
	return ui, nil
}

//func (est *ActStorage) Del(ctx context.Context, d []int) ([]int, error)
func (est *ActStorage) Del(ctx context.Context, d []int) ([]int, error) {
	dbpool := pgclient.WDB
	res := []int{}
	i := 0
	for _, id := range d {
		err := dbpool.QueryRow(context.Background(), "SELECT func_acts_del($1);", id).Scan(&i)
		res = append(res, i)

		if err != nil {
			log.Println("Failed execute func_acts_del: ", err)
		}
	}
	return res, nil
}

//func (est *ActStorage) GetOne(ctx context.Context, i int) (models.Act_count, error)
func (est *ActStorage) GetOne(ctx context.Context, i int) (models.Act_count, error) {
	dbpool := pgclient.WDB
	out_arr := []models.Act{}
	g := models.Act{}

	err := dbpool.QueryRow(context.Background(), "SELECT * from func_act_get($1);", i).Scan(&g.Id, &g.ActType.Id, &g.ActNumber,
		&g.ActDate, &g.Object.Id, &g.Staff.Id, &g.Notes, &g.Activated, &g.ActType.ActTypeName, &g.Object.ObjectName,
		&g.Object.FlatNumber, &g.Object.RegQty, &g.Object.House.Street.StreetName, &g.Object.House.Street.City.CityName,
		&g.Object.House.HouseNumber, &g.Object.House.BuildingNumber, &g.Object.TariffGroup.TariffGroupName, &g.Staff.StaffName)

	if err != nil && err != pgx.ErrNoRows {
		log.Println("Failed execute from func_act_get: ", err)
	}

	out_arr = append(out_arr, g)

	out_count := models.Act_count{Values: out_arr, Count: 1, Auth: models.Auth{}}
	return out_count, nil
}

//func (est *ActStorage) Activate(ctx context.Context, i int, d string) (int, error)
func (est *ActStorage) Activate(ctx context.Context, i int, d string) (int, error) {
	dbpool := pgclient.WDB
	ai := 0

	err := dbpool.QueryRow(ctx, "SELECT * from func_acts_activate($1,$2);", utils.NullableInt(int32(i)), utils.NullableString(d)).Scan(&ai)

	if err != nil {
		log.Println("Failed execute func_acts_activate: ", err)
		return 0, err
	}
	return ai, nil
}
