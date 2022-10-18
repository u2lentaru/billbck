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

//type ObjectStorage struct
type ObjectStorage struct {
	db *pgxpool.Pool
}

//func NewObjectStorage(db *pgxpool.Pool) *ObjectStorage
func NewObjectStorage(db *pgxpool.Pool) *ObjectStorage {
	return &ObjectStorage{db: db}
}

//func (est *ObjectStorage) GetList(ctx context.Context, pg, pgs int, gs1, gs2 string, gs3, gs3f bool, ord int, dsc bool) (models.Object_count, error)
func (est *ObjectStorage) GetList(ctx context.Context, pg, pgs int, gs1, gs2 string, gs3, gs3f bool, ord int, dsc bool) (models.Object_count, error) {
	dbpool := pgclient.WDB
	gs := models.Object{}

	gsc := 0
	err := dbpool.QueryRow(ctx, "SELECT * from func_objects_cnt($1,$2,$3);", gs1, utils.NullableString(gs2), utils.NullableBool(gs3, gs3f)).Scan(&gsc)

	if err != nil {
		log.Println(err.Error(), "func_objects_cnt")
		return models.Object_count{Values: []models.Object{}, Count: gsc, Auth: models.Auth{}}, err
	}

	out_arr := make([]models.Object, 0,
		func() int {
			if gsc < pgs {
				return gsc
			} else {
				return pgs
			}
		}())

	rows, err := dbpool.Query(ctx, "SELECT * from func_objects_get($1,$2,$3,$4,$5,$6,$7);", pg, pgs, gs1, utils.NullableString(gs2),
		utils.NullableBool(gs3, gs3f), ord, dsc)
	if err != nil {
		log.Println(err.Error(), "func_objects_get")
		return models.Object_count{Values: []models.Object{}, Count: gsc, Auth: models.Auth{}}, err
	}

	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&gs.Id, &gs.ObjectName, &gs.House.Id, &gs.FlatNumber, &gs.ObjType.Id, &gs.RegQty, &gs.Uzo.Id, &gs.TariffGroup.Id,
			&gs.Notes, &gs.CalculationType.Id, &gs.ObjStatus.Id, &gs.MffId, &gs.House.BuildingType.Id, &gs.House.Street.Id, &gs.House.HouseNumber,
			&gs.House.BuildingNumber, &gs.House.RP.Id, &gs.House.Area.Id, &gs.House.Ksk.Id, &gs.House.Sector.Id, &gs.House.Connector.Id,
			&gs.House.InputType.Id, &gs.House.Reliability.Id, &gs.House.Voltage.Id, &gs.House.BuildingType.BuildingTypeName,
			&gs.House.Street.StreetName, &gs.House.Street.Created, &gs.House.Street.City.CityName, &gs.House.RP.RpName, &gs.House.Area.AreaName,
			&gs.House.Area.AreaNumber, &gs.House.Ksk.KskName, &gs.House.Sector.SectorName, &gs.House.Connector.ConnectorName,
			&gs.House.InputType.InputTypeName, &gs.House.Reliability.ReliabilityName, &gs.House.Voltage.VoltageName, &gs.House.Voltage.VoltageValue,
			&gs.ObjType.ObjTypeName, &gs.Uzo.UzoName, &gs.Uzo.UzoValue, &gs.TariffGroup.TariffGroupName, &gs.CalculationType.CalculationTypeName,
			&gs.ObjStatus.ObjStatusName)

		if err != nil {
			log.Println("failed to scan row:", err)
		}

		out_arr = append(out_arr, gs)
	}

	out_count := models.Object_count{Values: out_arr, Count: gsc, Auth: models.Auth{}}
	return out_count, nil
}

//func (est *ObjectStorage) Add(ctx context.Context, ea models.Object) (int, error)
func (est *ObjectStorage) Add(ctx context.Context, a models.Object) (int, error) {
	dbpool := pgclient.WDB
	ai := 0

	err := dbpool.QueryRow(ctx, "SELECT func_objects_add($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11);", a.ObjectName, a.House.Id,
		a.FlatNumber, a.ObjType.Id, a.RegQty, a.Uzo.Id, a.TariffGroup.Id, a.CalculationType.Id, a.ObjStatus.Id, a.Notes, a.MffId).Scan(&ai)

	if err != nil {
		log.Println("Failed execute func_objects_add: ", err)
		return 0, err
	}

	return ai, nil
}

//func (est *ObjectStorage) Upd(ctx context.Context, eu models.Object) (int, error)
func (est *ObjectStorage) Upd(ctx context.Context, u models.Object) (int, error) {
	dbpool := pgclient.WDB
	ui := 0

	err := dbpool.QueryRow(ctx, "SELECT func_objects_upd($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12);", u.Id, u.ObjectName,
		u.House.Id, u.FlatNumber, u.ObjType.Id, u.RegQty, u.Uzo.Id, u.TariffGroup.Id, u.CalculationType.Id, u.ObjStatus.Id, u.Notes,
		u.MffId).Scan(&ui)

	if err != nil {
		log.Println("Failed execute func_objects_upd: ", err)
		return 0, err
	}
	return ui, nil
}

//func (est *ObjectStorage) Del(ctx context.Context, d []int) ([]int, error)
func (est *ObjectStorage) Del(ctx context.Context, d []int) ([]int, error) {
	dbpool := pgclient.WDB
	res := []int{}
	i := 0
	for _, id := range d {
		err := dbpool.QueryRow(ctx, "SELECT func_objects_del($1);", id).Scan(&i)
		res = append(res, i)

		if err != nil {
			log.Println("Failed execute func_objects_del: ", err)
		}
	}
	return res, nil
}

//func (est *ObjectStorage) GetOne(ctx context.Context, i int) (models.Object_count, error)
func (est *ObjectStorage) GetOne(ctx context.Context, i int) (models.Object_count, error) {
	dbpool := pgclient.WDB
	out_arr := []models.Object{}
	g := models.Object{}

	err := dbpool.QueryRow(ctx, "SELECT * from func_object_get($1);", i).Scan(&g.Id, &g.ObjectName, &g.House.Id, &g.FlatNumber, &g.ObjType.Id,
		&g.RegQty, &g.Uzo.Id, &g.TariffGroup.Id, &g.Notes, &g.CalculationType.Id, &g.ObjStatus.Id, &g.MffId, &g.House.BuildingType.Id,
		&g.House.Street.Id, &g.House.HouseNumber, &g.House.BuildingNumber, &g.House.RP.Id, &g.House.Area.Id, &g.House.Ksk.Id, &g.House.Sector.Id,
		&g.House.Connector.Id, &g.House.InputType.Id, &g.House.Reliability.Id, &g.House.Voltage.Id, &g.House.BuildingType.BuildingTypeName,
		&g.House.Street.StreetName, &g.House.Street.Created, &g.House.Street.City.CityName, &g.House.RP.RpName, &g.House.Area.AreaName,
		&g.House.Area.AreaNumber, &g.House.Ksk.KskName, &g.House.Sector.SectorName, &g.House.Connector.ConnectorName, &g.House.InputType.InputTypeName,
		&g.House.Reliability.ReliabilityName, &g.House.Voltage.VoltageName, &g.House.Voltage.VoltageValue, &g.ObjType.ObjTypeName, &g.Uzo.UzoName,
		&g.Uzo.UzoValue, &g.TariffGroup.TariffGroupName, &g.CalculationType.CalculationTypeName, &g.ObjStatus.ObjStatusName)

	if err != nil && err != pgx.ErrNoRows {
		log.Println("Failed execute from func_object_get: ", err)
		return models.Object_count{Values: []models.Object{}, Count: 0, Auth: models.Auth{}}, err
	}

	out_arr = append(out_arr, g)

	out_count := models.Object_count{Values: out_arr, Count: 1, Auth: models.Auth{}}
	return out_count, nil
}

//func (est *ObjectStorage) GetObjContract(ctx context.Context, i int, a string) (models.ObjContract_count, error)
func (est *ObjectStorage) GetObjContract(ctx context.Context, i int, a string) (models.ObjContract_count, error) {
	dbpool := pgclient.WDB
	out_arr := []models.ObjContract{}
	gs := models.ObjContract{}

	rows, err := dbpool.Query(ctx, "SELECT * from func_object_getcontract($1,$2);", i, utils.NullableString(a))

	if err != nil && err != pgx.ErrNoRows {
		log.Println("Failed execute from func_contract_getobject: ", err)
		return models.ObjContract_count{}, err
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&gs.Id, &gs.Contract.Id, &gs.Object.Id, &gs.Startdate, &gs.Enddate, &gs.Object.ObjectName, &gs.Object.RegQty,
			&gs.Object.FlatNumber, &gs.Object.House.Id, &gs.Object.House.HouseNumber, &gs.Object.House.BuildingNumber,
			&gs.Object.House.Street.City.CityName, &gs.Object.House.Street.Id, &gs.Object.House.Street.StreetName, &gs.Object.TariffGroup.Id,
			&gs.Object.TariffGroup.TariffGroupName, &gs.Contract.ContractNumber, &gs.Contract.Startdate, &gs.Contract.Enddate,
			&gs.Contract.Customer.SubId, &gs.Contract.Customer.SubName, &gs.Contract.Consignee.SubId, &gs.Contract.Consignee.SubName)
		if err != nil {
			log.Println("failed to scan row:", err)
		}

		out_arr = append(out_arr, gs)
	}

	gsc := 0
	err = dbpool.QueryRow(ctx, "SELECT * from func_object_getcontract_cnt($1,$2);", i, utils.NullableString(a)).Scan(&gsc)

	if err != nil {
		log.Println("Failed execute from func_contract_getobject_cnt: ", err)
		return models.ObjContract_count{}, err
	}

	out_count := models.ObjContract_count{Values: out_arr, Count: gsc, Auth: models.Auth{}}

	return out_count, nil
}

//func (est *ObjectStorage) GetMff(ctx context.Context, i int) (models.Object_count, error)
func (est *ObjectStorage) GetMff(ctx context.Context, i int) (models.Object_count, error) {
	dbpool := pgclient.WDB
	g := models.Object{}
	out_arr := []models.Object{}

	err := dbpool.QueryRow(ctx, "SELECT * from func_objects_mff($1);", i).Scan(&g.Id, &g.ObjectName, &g.House.Id, &g.FlatNumber, &g.ObjType.Id,
		&g.RegQty, &g.Uzo.Id, &g.TariffGroup.Id, &g.Notes, &g.CalculationType.Id, &g.ObjStatus.Id, &g.MffId)

	if err != nil && err != pgx.ErrNoRows {
		log.Println("Failed execute from func_objects_mff: ", err)
	}

	out_arr = append(out_arr, g)

	out_count := models.Object_count{Values: out_arr, Count: 1, Auth: models.Auth{}}
	return out_count, nil
}
