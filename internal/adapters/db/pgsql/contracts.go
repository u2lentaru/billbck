package pgsql

import (
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/u2lentaru/billbck/internal/models"
	"github.com/u2lentaru/billbck/internal/utils"
	"github.com/u2lentaru/billbck/pkg/pgclient"
)

//type ContractStorage struct
type ContractStorage struct {
	db *pgxpool.Pool
}

//func NewContractStorage(db *pgxpool.Pool) *ContractStorage
func NewContractStorage(db *pgxpool.Pool) *ContractStorage {
	return &ContractStorage{db: db}
}

//func (est *ContractStorage) GetList(ctx context.Context, pg, pgs int, gs1, gs2, gs3, gs4 string, gs5, gs6 int, gs7, gs8, gs9, gs10, gs11, gs12, gs13, gs14 string, ord int, dsc bool) (models.Contract_count, error)
func (est *ContractStorage) GetList(ctx context.Context, pg, pgs int, gs1, gs2, gs3, gs4 string, gs5, gs6 int, gs7, gs8, gs9, gs10, gs11, gs12, gs13, gs14 string, ord int, dsc bool) (models.Contract_count, error) {
	dbpool := pgclient.WDB
	gs := models.Contract{}

	gsc := 0
	err := dbpool.QueryRow(ctx, "SELECT * from func_contracts_cnt($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14);", gs1, gs2, utils.NullableString(gs3),
		gs4, utils.NullableInt(int32(gs5)), utils.NullableInt(int32(gs6)), gs7, gs8, gs9, gs10, gs11, gs12, utils.NullableString(gs13),
		utils.NullableString(gs14)).Scan(&gsc)

	if err != nil {
		log.Println(err.Error(), "func_contracts_cnt")
		return models.Contract_count{Values: []models.Contract{}, Count: gsc, Auth: models.Auth{}}, err
	}

	out_arr := make([]models.Contract, 0,
		func() int {
			if gsc < pgs {
				return gsc
			} else {
				return pgs
			}
		}())

	rows, err := dbpool.Query(ctx, "SELECT * from func_contracts_get($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18);",
		pg, pgs, gs1, gs2, utils.NullableString(gs3), gs4, utils.NullableInt(int32(gs5)), utils.NullableInt(int32(gs6)), gs7, gs8, gs9,
		gs10, gs11, gs12, utils.NullableString(gs13), utils.NullableString(gs14), ord, dsc)
	if err != nil {
		log.Println(err.Error())
		return models.Contract_count{Values: []models.Contract{}, Count: gsc, Auth: models.Auth{}}, err
	}

	defer rows.Close()

	var rmi sql.NullInt32
	var rmr, rmk sql.NullString

	for rows.Next() {
		err = rows.Scan(&gs.Id, &gs.BarCode, &gs.PersonalAccount, &gs.Id, &gs.ContractNumber, &gs.Startdate, &gs.Enddate, &gs.Customer.SubId,
			&gs.Consignee.SubId, &gs.EsoContractNumber, &gs.Eso.Id, &gs.Area.Id, &gs.CustomerGroup.Id, &rmi, &gs.Notes, &gs.MotNotes,
			&gs.Customer.SubName, &gs.Customer.SubAddr, &gs.Consignee.SubName, &gs.Eso.EsoName, &gs.Area.AreaName,
			&gs.CustomerGroup.CustomerGroupName, &rmr, &rmk)

		gs.ContractMot.Id = int(rmi.Int32)
		gs.ContractMot.ContractMotNameRu = rmr.String
		gs.ContractMot.ContractMotNameKz = rmk.String

		if err != nil {
			log.Println("failed to scan row:", err)
		}

		out_arr = append(out_arr, gs)
	}

	out_count := models.Contract_count{Values: out_arr, Count: gsc, Auth: models.Auth{}}
	if err != nil {
		log.Println(err.Error())
		return models.Contract_count{}, err
	}

	return out_count, nil
}

//func (est *ContractStorage) Add(ctx context.Context, ea models.Contract) (int, error)
func (est *ContractStorage) Add(ctx context.Context, a models.Contract) (int, error) {
	dbpool := pgclient.WDB
	ai := 0

	err := dbpool.QueryRow(context.Background(), "SELECT func_contracts_add($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11);",
		a.PersonalAccount, a.BarCode, a.ContractNumber, a.Startdate, a.Customer.SubId, a.Consignee.SubId, a.EsoContractNumber, a.Eso.Id,
		a.Area.Id, a.CustomerGroup.Id, a.Notes).Scan(&ai)

	if err != nil {
		log.Println("Failed execute func_contracts_add: ", err)
		return 0, err
	}

	return ai, nil
}

//func (est *ContractStorage) Upd(ctx context.Context, eu models.Contract) (int, error)
func (est *ContractStorage) Upd(ctx context.Context, u models.Contract) (int, error) {
	dbpool := pgclient.WDB
	ui := 0

	err := dbpool.QueryRow(context.Background(), "SELECT func_contracts_upd($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13);", u.Id,
		u.ContractNumber, u.Startdate, u.Enddate, u.Customer.SubId, u.Consignee.SubId, u.EsoContractNumber, u.Eso.Id, u.Area.Id,
		u.CustomerGroup.Id, utils.NullableInt(int32(u.ContractMot.Id)), u.Notes, u.MotNotes).Scan(&ui)

	if err != nil {
		log.Println("Failed execute func_contracts_upd: ", err)
		return 0, err
	}
	return ui, nil
}

//func (est *ContractStorage) Del(ctx context.Context, d models.IdClose) (int, error)
func (est *ContractStorage) Del(ctx context.Context, d models.IdClose) (int, error) {
	dbpool := pgclient.WDB

	i := 0
	_, err := time.Parse("2006-01-02", d.CloseDate)
	if err != nil {
		d.CloseDate = time.Now().Format("2006-01-02")
	}
	err = dbpool.QueryRow(ctx, "SELECT func_contracts_del($1,$2,$3,$4);", d.Id, d.CloseDate, d.ContractMot.Id,
		utils.NullableString(d.MotNotes)).Scan(&i)

	if err != nil {
		log.Println("Failed execute func_contracts_del: ", err)
		return 0, err
	}

	return i, nil
}

//func (est *ContractStorage) GetOne(ctx context.Context, i int) (models.Contract_count, error)
func (est *ContractStorage) GetOne(ctx context.Context, i int) (models.Contract_count, error) {
	dbpool := pgclient.WDB
	out_arr := []models.Contract{}
	g := models.Contract{}

	var rmi sql.NullInt32
	var rmr, rmk sql.NullString

	err := dbpool.QueryRow(context.Background(), "SELECT * from func_contract_get($1);", i).Scan(&g.Id, &g.BarCode, &g.PersonalAccount,
		&g.Id, &g.ContractNumber, &g.Startdate, &g.Enddate, &g.Customer.SubId, &g.Consignee.SubId, &g.EsoContractNumber, &g.Eso.Id,
		&g.Area.Id, &g.CustomerGroup.Id, &rmi, &g.Notes, &g.MotNotes, &g.Customer.SubName, &g.Customer.SubAddr, &g.Consignee.SubName,
		&g.Eso.EsoName, &g.Area.AreaName, &g.CustomerGroup.CustomerGroupName, &rmr, &rmk)

	g.ContractMot.Id = int(rmi.Int32)
	g.ContractMot.ContractMotNameRu = rmr.String
	g.ContractMot.ContractMotNameKz = rmk.String

	if err != nil && err != pgx.ErrNoRows {
		log.Println("Failed execute from func_contract_get: ", err)
		return models.Contract_count{Values: []models.Contract{}, Count: 0, Auth: models.Auth{}}, err
	}

	out_arr = append(out_arr, g)

	out_count := models.Contract_count{Values: out_arr, Count: 1, Auth: models.Auth{}}
	return out_count, nil
}

//func (est *ContractStorage) GetObj(ctx context.Context, i int, a string) (models.ObjContract, error)
func (est *ContractStorage) GetObj(ctx context.Context, i int, a string) (models.ObjContract_count, error) {
	dbpool := pgclient.WDB
	out_arr := []models.ObjContract{}
	gs := models.ObjContract{}

	rows, err := dbpool.Query(ctx, "SELECT * from func_contract_getobject($1,$2);", i, utils.NullableString(a))

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
	err = dbpool.QueryRow(ctx, "SELECT * from func_contract_getobject_cnt($1,$2);", i, utils.NullableString(a)).Scan(&gsc)

	if err != nil {
		log.Println("Failed execute from func_contract_getobject_cnt: ", err)
		return models.ObjContract_count{}, err
	}

	out_count := models.ObjContract_count{Values: out_arr, Count: gsc}

	return out_count, nil
}

//func (est *ContractStorage) GetHist(ctx context.Context, i int) (string, error)
func (est *ContractStorage) GetHist(ctx context.Context, i int) (string, error) {
	dbpool := pgclient.WDB

	h := ""
	rows, err := dbpool.Query(ctx, "SELECT * from func_contracts_hist($1);", i)

	if err != nil {
		log.Println("Failed execute from func_contract_getobject: ", err)
		return "", err
	}

	defer rows.Close()

	qa := false
	hist_arr := "["

	for rows.Next() {

		if qa {
			hist_arr += ","
		}
		qa = true

		err = rows.Scan(&h)
		hist_arr += h

		if err != nil {
			log.Println("failed to scan row:", err)
		}

	}
	hist_arr += "]"
	return hist_arr, nil
}
