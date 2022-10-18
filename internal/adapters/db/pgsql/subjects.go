package pgsql

import (
	"context"
	"database/sql"
	"log"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/u2lentaru/billbck/internal/models"
	"github.com/u2lentaru/billbck/internal/utils"
	"github.com/u2lentaru/billbck/pkg/pgclient"
)

//type SubjectStorage struct
type SubjectStorage struct {
	db *pgxpool.Pool
}

//func NewSubjectStorage(db *pgxpool.Pool) *SubjectStorage
func NewSubjectStorage(db *pgxpool.Pool) *SubjectStorage {
	return &SubjectStorage{db: db}
}

//func (est *SubjectStorage) GetList(ctx context.Context, pg, pgs int, gs1, gs2, gs3 string, ord int, dsc bool) (models.Subject_count, error) {
func (est *SubjectStorage) GetList(ctx context.Context, pg, pgs int, gs1, gs2, gs3 string, ord int, dsc bool) (models.Subject_count, error) {
	dbpool := pgclient.WDB
	gs := models.Subject{}

	gsc := 0
	err := dbpool.QueryRow(ctx, "SELECT * from func_sub_details_cnt($1,$2,$3);", gs1, gs2, gs3).Scan(&gsc)

	if err != nil {
		log.Println(err.Error(), "func_sub_details_cnt")
		return models.Subject_count{Values: []models.Subject{}, Count: gsc, Auth: models.Auth{}}, err
	}

	out_arr := make([]models.Subject, 0,
		func() int {
			if gsc < pgs {
				return gsc
			} else {
				return pgs
			}
		}())

	rows, err := dbpool.Query(ctx, "SELECT * from func_sub_details_get($1,$2,$3,$4,$5,$6,$7);", pg, pgs, gs1, gs2, gs3, ord, dsc)
	if err != nil {
		log.Println(err.Error(), "func_sub_details_get")
		return models.Subject_count{Values: []models.Subject{}, Count: gsc, Auth: models.Auth{}}, err
	}

	defer rows.Close()

	var shp, sap sql.NullInt32
	var shn, san sql.NullString

	for rows.Next() {
		err = rows.Scan(&gs.SubId, &gs.SubType.Id, &gs.SubPhys, &gs.SubDescr, &gs.SubName, &gs.SubBin, &shp, &gs.SubHeadName,
			&sap, &gs.SubAccName, &gs.SubAddr, &gs.SubPhone, &gs.SubStart, &gs.SubAccNumber, &gs.SubType.SubTypeName,
			&shn, &san, &gs.Job, &gs.Email, &gs.MobPhone, &gs.JobPhone, &gs.Notes)
		if err != nil {
			log.Println("failed to scan row:", err)
		}

		gs.SubHeadPos.Id = int(shp.Int32)
		gs.SubAccPos.Id = int(sap.Int32)
		gs.SubHeadPos.PositionName = shn.String
		gs.SubAccPos.PositionName = san.String

		out_arr = append(out_arr, gs)
	}

	out_count := models.Subject_count{Values: out_arr, Count: gsc, Auth: models.Auth{}}
	return out_count, nil
}

//func (est *SubjectStorage) Add(ctx context.Context, ea models.Subject) (int, error)
func (est *SubjectStorage) Add(ctx context.Context, a models.Subject) (int, error) {
	dbpool := pgclient.WDB
	ai := 0

	err := dbpool.QueryRow(ctx, "SELECT func_sub_details_add($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18);",
		a.SubType.Id, a.SubPhys, a.SubDescr, a.SubName, a.SubBin, utils.NullableInt(int32(a.SubHeadPos.Id)), a.SubHeadName,
		utils.NullableInt(int32(a.SubAccPos.Id)), a.SubAccName, a.SubAddr, a.SubPhone, a.SubStart, a.SubAccNumber, a.Job,
		a.Email, a.MobPhone, a.JobPhone, a.Notes).Scan(&ai)

	if err != nil {
		log.Println("Failed execute func_sub_details_add: ", err)
		return 0, err
	}

	return ai, nil
}

//func (est *SubjectStorage) Upd(ctx context.Context, eu models.Subject) (int, error)
func (est *SubjectStorage) Upd(ctx context.Context, u models.Subject) (int, error) {
	dbpool := pgclient.WDB
	ui := 0

	err := dbpool.QueryRow(ctx, "SELECT func_sub_details_upd($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19);",
		u.SubId, u.SubType.Id, u.SubPhys, u.SubDescr, u.SubName, u.SubBin, utils.NullableInt(int32(u.SubHeadPos.Id)), u.SubHeadName,
		utils.NullableInt(int32(u.SubAccPos.Id)), u.SubAccName, u.SubAddr, u.SubPhone, u.SubStart, u.SubAccNumber, u.Job, u.Email,
		u.MobPhone, u.JobPhone, u.Notes).Scan(&ui)

	if err != nil {
		log.Println("Failed execute func_sub_details_upd: ", err)
		return 0, err
	}
	return ui, nil
}

//func (est *SubjectStorage) Del(ctx context.Context, d models.SubjectClose) (int, error)
func (est *SubjectStorage) Del(ctx context.Context, d models.SubjectClose) (int, error) {
	dbpool := pgclient.WDB
	i := 0
	err := dbpool.QueryRow(ctx, "SELECT func_subjects_close($1,$2);", d.SubId, d.SubClose).Scan(&i)
	if err != nil {
		log.Println("Failed execute func_subjects_close: ", err)
	}
	return i, nil
}

//func (est *SubjectStorage) GetOne(ctx context.Context, i int) (models.Subject_count, error)
func (est *SubjectStorage) GetOne(ctx context.Context, i int) (models.Subject_count, error) {
	dbpool := pgclient.WDB
	out_arr := []models.Subject{}
	g := models.Subject{}

	var shp, sap sql.NullInt32
	var shn, san sql.NullString

	err := dbpool.QueryRow(ctx, "SELECT * from func_sub_detail_get($1);", i).Scan(&g.SubId, &g.SubType.Id, &g.SubPhys, &g.SubDescr,
		&g.SubName, &g.SubBin, &shp, &g.SubHeadName, &sap, &g.SubAccName, &g.SubAddr, &g.SubPhone, &g.SubStart, &g.SubAccNumber,
		&g.SubType.SubTypeName, &shn, &san, &g.Job, &g.Email, &g.MobPhone, &g.JobPhone, &g.Notes)

	if err != nil && err != pgx.ErrNoRows {
		log.Println("Failed execute func_sub_detail_get: ", err)
		return models.Subject_count{Values: []models.Subject{}, Count: 0, Auth: models.Auth{}}, err
	}

	g.SubHeadPos.Id = int(shp.Int32)
	g.SubAccPos.Id = int(sap.Int32)
	g.SubHeadPos.PositionName = shn.String
	g.SubAccPos.PositionName = san.String

	out_arr = append(out_arr, g)

	out_count := models.Subject_count{Values: out_arr, Count: 1, Auth: models.Auth{}}
	return out_count, nil
}

//func (est *SubjectStorage) GetHist(ctx context.Context, i int) ([]string, error)
func (est *SubjectStorage) GetHist(ctx context.Context, i int) ([]string, error) {
	dbpool := pgclient.WDB

	hist_arr := []string{}

	h := ""
	rows, err := dbpool.Query(ctx, "SELECT * from func_sub_detail_get_hist($1);", i)

	if err != nil && err != pgx.ErrNoRows {
		log.Println("Failed execute func_sub_detail_get_hist: ", err)
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&h)
		if err != nil {
			log.Println("failed to scan row:", err)
		}
		hist_arr = append(hist_arr, h)
	}

	return hist_arr, nil
}
