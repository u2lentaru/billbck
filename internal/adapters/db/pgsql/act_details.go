package pgsql

import (
	"context"
	"database/sql"
	"log"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/u2lentaru/billbck/internal/models"
	"github.com/u2lentaru/billbck/pkg/pgclient"
)

//type ActDetailStorage struct
type ActDetailStorage struct {
	db *pgxpool.Pool
}

//func NewActDetailStorage(db *pgxpool.Pool) *ActTypeStorage
func NewActDetailStorage(db *pgxpool.Pool) *ActDetailStorage {
	return &ActDetailStorage{db: db}
}

//func (est *ActDetailStorage) GetList(ctx context.Context, pg, pgs, nm, ord int, dsc bool) (models.ActDetail_count, error)
func (est *ActDetailStorage) GetList(ctx context.Context, pg, pgs, nm, ord int, dsc bool) (models.ActDetail_count, error) {
	dbpool := pgclient.WDB
	auth := models.Auth{Create: true, Read: true, Update: true, Delete: true}
	gs := models.ActDetail{}

	gsc := 0
	err := dbpool.QueryRow(ctx, "SELECT * from func_act_details_cnt($1);", nm).Scan(&gsc)

	if err != nil {
		log.Println(err.Error(), "func_act_details_cnt")
		return models.ActDetail_count{Values: []models.ActDetail{}, Count: gsc, Auth: auth}, err
	}

	out_arr := make([]models.ActDetail, 0,
		func() int {
			if gsc < pgs {
				return gsc
			} else {
				return pgs
			}
		}())

	rows, err := dbpool.Query(ctx, "SELECT * from func_act_details_get($1,$2,$3,$4,$5);", pg, pgs, nm, ord, dsc)
	if err != nil {
		log.Println(err.Error(), "func_act_details_get")
		return models.ActDetail_count{Values: []models.ActDetail{}, Count: gsc, Auth: auth}, err
	}

	defer rows.Close()

	var pti, si, ci, sti, ri, vi sql.NullInt32
	var ccn, ptn, stn, rn, vn sql.NullString

	for rows.Next() {
		err = rows.Scan(&gs.Act.Id, &gs.Act.ActType.Id, &gs.Act.ActNumber, &gs.Act.ActDate, &gs.Act.Object.Id, &gs.Act.Staff.Id,
			&gs.Act.Notes, &gs.Act.Activated, &gs.Act.ActType.ActTypeName, &gs.Act.Object.ObjectName, &gs.Act.Object.FlatNumber,
			&gs.Act.Object.RegQty, &gs.Act.Object.House.Street.StreetName, &gs.Act.Object.House.Street.City.CityName,
			&gs.Act.Object.House.HouseNumber, &gs.Act.Object.House.BuildingNumber, &gs.Act.Object.TariffGroup.TariffGroupName,
			&gs.Act.Staff.StaffName, &gs.Id, &gs.ActDetailDate, &gs.PuId, &pti, &gs.PuNumber, &gs.InstallDate, &gs.CheckInterval,
			&gs.InitialValue, &gs.DevStopped, &gs.Startdate, &gs.Enddate, &gs.Pid, &gs.AdPuValue, &si, &gs.SealNumber,
			&gs.SealDate, &ci, &gs.ConclusionNumber, &sti, &ri, &vi, &gs.Customer, &gs.CustomerPhone, &gs.CustomerPos,
			&gs.Notes, &ptn, &ccn, &stn, &rn, &vn)

		if err != nil {
			log.Println("failed to scan row:", err)
		}

		gs.Reason.Id = int(ri.Int32)
		gs.Violation.Id = int(vi.Int32)
		gs.Reason.ReasonName = rn.String
		gs.Violation.ViolationName = vn.String
		gs.PuType.Id = int(pti.Int32)
		gs.Seal.Id = int(si.Int32)
		// gs.Seal.SealNumber = ssn.String
		gs.Conclusion.Id = int(ci.Int32)
		gs.Conclusion.ConclusionName = ccn.String
		gs.ShutdownType.Id = int(sti.Int32)
		gs.PuType.PuTypeName = ptn.String
		gs.ShutdownType.ShutdownTypeName = stn.String

		out_arr = append(out_arr, gs)
	}

	out_count := models.ActDetail_count{Values: out_arr, Count: gsc, Auth: auth}

	if err != nil {
		log.Println(err.Error())
		return models.ActDetail_count{}, err
	}

	return out_count, nil
}

//func (est *ActDetailStorage) Add(ctx context.Context, ea models.ActDetail) (int, error)
func (est *ActDetailStorage) Add(ctx context.Context, a models.ActDetail) (int, error) {
	dbpool := pgclient.WDB
	ai := 0

	err := dbpool.QueryRow(context.Background(),
		"SELECT func_act_details_add($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19,$20,$21,$22,$23,$24,$25);",
		a.Act.Id, a.PuId, a.SealNumber, a.AdPuValue, a.ActDetailDate, a.PuNumber, a.InstallDate, a.CheckInterval, a.InitialValue,
		a.DevStopped, a.Startdate, a.Enddate, a.Pid, a.Seal.Id, a.SealDate, a.Notes, a.PuType.Id, a.Conclusion.Id,
		a.ConclusionNumber, a.ShutdownType.Id, a.CustomerPhone, a.CustomerPos, a.Reason.Id, a.Violation.Id, a.Customer).Scan(&ai)

	if err != nil {
		log.Println("Failed execute func_act_details_add: ", err)
		return 0, err
	}

	return ai, nil
}

//func (est *ActDetailStorage) Upd(ctx context.Context, eu models.ActDetail) (int, error)
func (est *ActDetailStorage) Upd(ctx context.Context, u models.ActDetail) (int, error) {
	dbpool := pgclient.WDB
	ui := 0

	err := dbpool.QueryRow(context.Background(),
		"SELECT func_act_details_upd($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19,$20,$21,$22,$23,$24,$25,$26);",
		u.Id, u.Act.Id, u.PuId, u.SealNumber, u.AdPuValue, u.ActDetailDate, u.PuNumber, u.InstallDate, u.CheckInterval, u.InitialValue,
		u.DevStopped, u.Startdate, u.Enddate, u.Pid, u.Seal.Id, u.SealDate, u.Notes, u.PuType.Id, u.Conclusion.Id,
		u.ConclusionNumber, u.ShutdownType.Id, u.CustomerPhone, u.CustomerPos, u.Reason.Id, u.Violation.Id, u.Customer).Scan(&ui)

	if err != nil {
		log.Println("Failed execute func_act_types_upd: ", err)
		return 0, err
	}
	return ui, nil
}

//func (est *ActDetailStorage) Del(ctx context.Context, ed []int) ([]int, error)
func (est *ActDetailStorage) Del(ctx context.Context, ed []int) ([]int, error) {
	dbpool := pgclient.WDB
	res := []int{}
	i := 0
	for _, id := range ed {
		err := dbpool.QueryRow(context.Background(), "SELECT func_act_details_del($1);", id).Scan(&i)
		res = append(res, i)

		if err != nil {
			log.Println("Failed execute func_act_details_del: ", err)
			return []int{}, err
		}
	}
	return res, nil
}

//func (est *ActDetailStorage) GetOne(ctx context.Context, i int) (models.ActDetail_count, error)
func (est *ActDetailStorage) GetOne(ctx context.Context, i int) (models.ActDetail_count, error) {
	dbpool := pgclient.WDB
	out_arr := []models.ActDetail{}
	g := models.ActDetail{}
	auth := models.Auth{Create: true, Read: true, Update: true, Delete: true}

	var pti, si, ci, sti, ri, vi sql.NullInt32
	var ccn, ptn, stn, rn, vn sql.NullString

	err := dbpool.QueryRow(context.Background(), "SELECT * from func_act_detail_get($1);", i).Scan(&g.Act.Id, &g.Act.ActType.Id,
		&g.Act.ActNumber, &g.Act.ActDate, &g.Act.Object.Id, &g.Act.Staff.Id, &g.Act.Notes, &g.Act.Activated, &g.Act.ActType.ActTypeName,
		&g.Act.Object.ObjectName, &g.Act.Object.FlatNumber, &g.Act.Object.RegQty, &g.Act.Object.House.Street.StreetName,
		&g.Act.Object.House.Street.City.CityName, &g.Act.Object.House.HouseNumber, &g.Act.Object.House.BuildingNumber,
		&g.Act.Object.TariffGroup.TariffGroupName, &g.Act.Staff.StaffName, &g.Id, &g.ActDetailDate, &g.PuId, &pti, &g.PuNumber,
		&g.InstallDate, &g.CheckInterval, &g.InitialValue, &g.DevStopped, &g.Startdate, &g.Enddate, &g.Pid, &g.AdPuValue, &si,
		&g.SealNumber, &g.SealDate, &ci, &g.ConclusionNumber, &sti, &ri, &vi, &g.Customer, &g.CustomerPhone, &g.CustomerPos,
		&g.Notes, &ptn, &ccn, &stn, &rn, &vn)

	if err != nil && err != pgx.ErrNoRows {
		log.Println("Failed execute from func_act_type_get: ", err)
		return models.ActDetail_count{Values: []models.ActDetail{}, Count: 0, Auth: auth}, err
	}

	g.Reason.Id = int(ri.Int32)
	g.Violation.Id = int(vi.Int32)
	g.Reason.ReasonName = rn.String
	g.Violation.ViolationName = vn.String
	g.PuType.Id = int(pti.Int32)
	g.Seal.Id = int(si.Int32)
	// g.Seal.SealNumber = ssn.String
	g.Conclusion.Id = int(ci.Int32)
	g.Conclusion.ConclusionName = ccn.String
	g.ShutdownType.Id = int(sti.Int32)
	g.PuType.PuTypeName = ptn.String
	g.ShutdownType.ShutdownTypeName = stn.String

	out_arr = append(out_arr, g)

	out_count := models.ActDetail_count{Values: out_arr, Count: 0, Auth: auth}
	return out_count, nil
}
