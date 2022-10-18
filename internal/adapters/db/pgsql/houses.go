package pgsql

import (
	"context"
	"log"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/u2lentaru/billbck/internal/models"
	"github.com/u2lentaru/billbck/pkg/pgclient"
)

//type HouseStorage struct
type HouseStorage struct {
	db *pgxpool.Pool
}

//func NewHouseStorage(db *pgxpool.Pool) *HouseStorage
func NewHouseStorage(db *pgxpool.Pool) *HouseStorage {
	return &HouseStorage{db: db}
}

//func (est *HouseStorage) GetList(ctx context.Context, pg, pgs int, gs1, gs2 string, gs3, ord int, dsc bool) (models.House_count, error) {
func (est *HouseStorage) GetList(ctx context.Context, pg, pgs int, gs1, gs2 string, gs3, ord int, dsc bool) (models.House_count, error) {
	dbpool := pgclient.WDB
	gs := models.House{}

	gsc := 0
	err := dbpool.QueryRow(ctx, "SELECT * from func_houses_cnt($1,$2,$3);", gs1, gs2, gs3).Scan(&gsc)

	if err != nil {
		log.Println(err.Error(), "func_houses_cnt")
		return models.House_count{Values: []models.House{}, Count: gsc, Auth: models.Auth{}}, err
	}

	out_arr := make([]models.House, 0,
		func() int {
			if gsc < pgs {
				return gsc
			} else {
				return pgs
			}
		}())

	rows, err := dbpool.Query(ctx, "SELECT * from func_houses_get($1,$2,$3,$4,$5,$6,$7);", pg, pgs, gs1, gs2, gs3, ord, dsc)
	if err != nil {
		log.Println(err.Error())
		return models.House_count{Values: []models.House{}, Count: gsc, Auth: models.Auth{}}, err
	}

	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&gs.Id, &gs.BuildingType.Id, &gs.Street.Id, &gs.HouseNumber, &gs.BuildingNumber, &gs.RP.Id, &gs.Area.Id, &gs.Ksk.Id,
			&gs.Sector.Id, &gs.Connector.Id, &gs.InputType.Id, &gs.Reliability.Id, &gs.Voltage.Id, &gs.Notes, &gs.BuildingType.BuildingTypeName,
			&gs.Street.StreetName, &gs.Street.Created, &gs.Street.City.CityName, &gs.RP.RpName, &gs.Area.AreaName, &gs.Area.AreaNumber, &gs.Ksk.KskName,
			&gs.Sector.SectorName, &gs.Connector.ConnectorName, &gs.InputType.InputTypeName, &gs.Reliability.ReliabilityName,
			&gs.Voltage.VoltageName, &gs.Voltage.VoltageValue)
		if err != nil {
			log.Println("failed to scan row:", err)
		}

		out_arr = append(out_arr, gs)
	}

	out_count := models.House_count{Values: out_arr, Count: gsc, Auth: models.Auth{}}
	if err != nil {
		log.Println(err.Error())
		return models.House_count{}, err
	}

	return out_count, nil
}

//func (est *HouseStorage) Add(ctx context.Context, ea models.House) (int, error)
func (est *HouseStorage) Add(ctx context.Context, a models.House) (int, error) {
	dbpool := pgclient.WDB
	ai := 0

	err := dbpool.QueryRow(ctx, "SELECT func_houses_add($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13);", a.BuildingType.Id, a.Street.Id,
		a.HouseNumber, a.BuildingNumber, a.RP.Id, a.Area.Id, a.Ksk.Id, a.Sector.Id, a.Connector.Id, a.InputType.Id, a.Reliability.Id,
		a.Voltage.Id, a.Notes).Scan(&ai)

	if err != nil {
		log.Println("Failed execute func_houses_add: ", err)
		return 0, err
	}

	return ai, nil
}

//func (est *HouseStorage) Upd(ctx context.Context, eu models.House) (int, error)
func (est *HouseStorage) Upd(ctx context.Context, u models.House) (int, error) {
	dbpool := pgclient.WDB
	ui := 0

	err := dbpool.QueryRow(ctx, "SELECT func_houses_upd($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14);", u.Id, u.BuildingType.Id,
		u.Street.Id, u.HouseNumber, u.BuildingNumber, u.RP.Id, u.Area.Id, u.Ksk.Id, u.Sector.Id, u.Connector.Id, u.InputType.Id,
		u.Reliability.Id, u.Voltage.Id, u.Notes).Scan(&ui)

	if err != nil {
		log.Println("Failed execute func_houses_upd: ", err)
		return 0, err
	}
	return ui, nil
}

//func (est *HouseStorage) Del(ctx context.Context, ed []int) ([]int, error)
func (est *HouseStorage) Del(ctx context.Context, d []int) ([]int, error) {
	dbpool := pgclient.WDB
	res := []int{}
	i := 0
	for _, id := range d {
		err := dbpool.QueryRow(ctx, "SELECT func_houses_del($1);", id).Scan(&i)
		res = append(res, i)

		if err != nil {
			log.Println("Failed execute func_houses_del: ", err)
		}
	}
	return res, nil
}

//func (est *HouseStorage) GetOne(ctx context.Context, i int) (models.House_count, error)
func (est *HouseStorage) GetOne(ctx context.Context, i int) (models.House_count, error) {
	dbpool := pgclient.WDB
	out_arr := []models.House{}
	g := models.House{}

	err := dbpool.QueryRow(ctx, "SELECT * from func_house_get($1);", i).Scan(&g.Id, &g.BuildingType.Id, &g.Street.Id, &g.HouseNumber,
		&g.BuildingNumber, &g.RP.Id, &g.Area.Id, &g.Ksk.Id, &g.Sector.Id, &g.Connector.Id, &g.InputType.Id, &g.Reliability.Id, &g.Voltage.Id,
		&g.Notes, &g.BuildingType.BuildingTypeName, &g.Street.StreetName, &g.Street.Created, &g.Street.City.CityName, &g.RP.RpName,
		&g.Area.AreaName, &g.Area.AreaNumber, &g.Ksk.KskName, &g.Sector.SectorName, &g.Connector.ConnectorName, &g.InputType.InputTypeName,
		&g.Reliability.ReliabilityName, &g.Voltage.VoltageName, &g.Voltage.VoltageValue)

	if err != nil && err != pgx.ErrNoRows {
		log.Println("Failed execute from func_house_get: ", err)
		return models.House_count{Values: []models.House{}, Count: 0, Auth: models.Auth{}}, err
	}

	out_arr = append(out_arr, g)

	out_count := models.House_count{Values: out_arr, Count: 1, Auth: models.Auth{}}
	return out_count, nil
}
