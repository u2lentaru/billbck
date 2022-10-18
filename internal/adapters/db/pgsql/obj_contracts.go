package pgsql

import (
	"context"
	"log"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/u2lentaru/billbck/internal/models"
	"github.com/u2lentaru/billbck/internal/utils"
	"github.com/u2lentaru/billbck/pkg/pgclient"
)

//type ObjContractStorage struct
type ObjContractStorage struct {
	db *pgxpool.Pool
}

//func NewObjContractStorage(db *pgxpool.Pool) *ObjContractStorage
func NewObjContractStorage(db *pgxpool.Pool) *ObjContractStorage {
	return &ObjContractStorage{db: db}
}

//func (est *ObjContractStorage) GetList(ctx context.Context, pg, pgs, gs1, gs2, gs3  int, gs4, gs4f bool, ord int, dsc bool) (models.ObjContract_count, error) {
func (est *ObjContractStorage) GetList(ctx context.Context, pg, pgs, gs1, gs2, gs3 int, gs4, gs4f bool, ord int, dsc bool) (models.ObjContract_count, error) {
	dbpool := pgclient.WDB
	gs := models.ObjContract{}

	gsc := 0
	err := dbpool.QueryRow(ctx, "SELECT * from func_obj_contracts_cnt($1,$2,$3,$4);", utils.NullableInt(int32(gs1)),
		utils.NullableInt(int32(gs2)), utils.NullableInt(int32(gs3)), utils.NullableBool(gs4, gs4f)).Scan(&gsc)

	if err != nil {
		log.Println(err.Error(), "func_obj_contracts_cnt")
		return models.ObjContract_count{Values: []models.ObjContract{}, Count: gsc, Auth: models.Auth{}}, err
	}

	out_arr := make([]models.ObjContract, 0,
		func() int {
			if gsc < pgs {
				return gsc
			} else {
				return pgs
			}
		}())

	rows, err := dbpool.Query(ctx, "SELECT * from func_obj_contracts_get($1,$2,$3,$4,$5,$6,$7,$8);", pg, pgs, utils.NullableInt(int32(gs1)),
		utils.NullableInt(int32(gs2)), utils.NullableInt(int32(gs3)), utils.NullableBool(gs4, gs4f), ord, dsc)
	if err != nil {
		log.Println(err.Error())
		return models.ObjContract_count{Values: []models.ObjContract{}, Count: gsc, Auth: models.Auth{}}, err
	}

	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&gs.Id, &gs.Contract.Id, &gs.Object.Id, &gs.ObjTypeId, &gs.Startdate, &gs.Enddate, &gs.Object.ObjectName, &gs.Object.RegQty,
			&gs.Object.FlatNumber, &gs.Object.House.Id, &gs.Object.House.HouseNumber, &gs.Object.House.BuildingNumber,
			&gs.Object.House.Street.City.CityName, &gs.Object.House.Street.Id, &gs.Object.House.Street.StreetName, &gs.Object.TariffGroup.Id,
			&gs.Object.TariffGroup.TariffGroupName, &gs.Contract.ContractNumber, &gs.Contract.Startdate, &gs.Contract.Enddate,
			&gs.Contract.Customer.SubId, &gs.Contract.Customer.SubName, &gs.Contract.Consignee.SubId, &gs.Contract.Consignee.SubName)
		if err != nil {
			log.Println("failed to scan row:", err)
		}

		out_arr = append(out_arr, gs)
	}

	out_count := models.ObjContract_count{Values: out_arr, Count: gsc, Auth: models.Auth{}}
	if err != nil {
		log.Println(err.Error())
		return models.ObjContract_count{}, err
	}

	return out_count, nil
}

//func (est *ObjContractStorage) Add(ctx context.Context, a models.ObjContract) (int, error)
func (est *ObjContractStorage) Add(ctx context.Context, a models.ObjContract) (int, error) {
	dbpool := pgclient.WDB
	ai := 0

	err := dbpool.QueryRow(ctx, "SELECT func_obj_contracts_add($1,$2,$3,$4);", a.Contract.Id, a.Object.Id, a.ObjTypeId, a.Startdate).Scan(&ai)

	if err != nil {
		log.Println("Failed execute func_obj_contracts_add: ", err)
		return 0, err
	}

	return ai, nil
}

//func (est *ObjContractStorage) Upd(ctx context.Context, u models.ObjContract) (int, error)
func (est *ObjContractStorage) Upd(ctx context.Context, u models.ObjContract) (int, error) {
	dbpool := pgclient.WDB
	ui := 0

	err := dbpool.QueryRow(ctx, "SELECT func_obj_contracts_upd($1,$2,$3,$4,$5,$6);", u.Id, u.Contract.Id, u.Object.Id, u.ObjTypeId,
		u.Startdate, u.Enddate).Scan(&ui)

	if err != nil {
		log.Println("Failed execute func_obj_contracts_upd: ", err)
		return 0, err
	}
	return ui, nil
}

//func (est *ObjContractStorage) Del(ctx context.Context, d models.IdClose) (int, error)
func (est *ObjContractStorage) Del(ctx context.Context, d models.IdClose) (int, error) {
	dbpool := pgclient.WDB
	i := 0
	_, err := time.Parse("2006-01-02", d.CloseDate)
	if err != nil {
		d.CloseDate = time.Now().Format("2006-01-02")
	}

	err = dbpool.QueryRow(ctx, "SELECT func_obj_contracts_del($1,$2);", d.Id, d.CloseDate).Scan(&i)

	if err != nil {
		log.Println("Failed execute func_obj_contracts_del: ", err)
	}
	return i, nil
}

//func (est *ObjContractStorage) GetOne(ctx context.Context, i int) (models.ObjContract_count, error)
func (est *ObjContractStorage) GetOne(ctx context.Context, i int, d string) (models.ObjContract_count, error) {
	dbpool := pgclient.WDB
	out_arr := []models.ObjContract{}
	g := models.ObjContract{}

	err := dbpool.QueryRow(ctx, "SELECT * from func_obj_contract_get($1,$2);", i, d).Scan(&g.Id, &g.Contract.Id, &g.Object.Id, &g.ObjTypeId,
		&g.Startdate, &g.Enddate, &g.Object.ObjectName, &g.Object.RegQty, &g.Object.FlatNumber, &g.Object.House.Id, &g.Object.House.HouseNumber,
		&g.Object.House.BuildingNumber, &g.Object.House.Street.City.CityName, &g.Object.House.Street.Id, &g.Object.House.Street.StreetName,
		&g.Object.TariffGroup.Id, &g.Object.TariffGroup.TariffGroupName, &g.Contract.ContractNumber, &g.Contract.Startdate, &g.Contract.Enddate,
		&g.Contract.Customer.SubId, &g.Contract.Customer.SubName, &g.Contract.Consignee.SubId, &g.Contract.Consignee.SubName)

	if err != nil && err != pgx.ErrNoRows {
		log.Println("Failed execute from func_obj_contract_get: ", err)
		return models.ObjContract_count{Values: []models.ObjContract{}, Count: 0, Auth: models.Auth{}}, err
	}

	out_arr = append(out_arr, g)

	out_count := models.ObjContract_count{Values: out_arr, Count: 1, Auth: models.Auth{}}
	return out_count, nil
}
