package pgsql

import (
	"context"
	"log"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/u2lentaru/billbck/internal/models"
	"github.com/u2lentaru/billbck/pkg/pgclient"
)

//type SealStorage struct
type SealStorage struct {
	db *pgxpool.Pool
}

//func NewSealStorage(db *pgxpool.Pool) *SealStorage
func NewSealStorage(db *pgxpool.Pool) *SealStorage {
	return &SealStorage{db: db}
}

//func (est *SealStorage) GetList(ctx context.Context, pg, pgs int, gs1 string, ord int, dsc bool) (models.Seal_count, error)
func (est *SealStorage) GetList(ctx context.Context, pg, pgs int, gs1 string, ord int, dsc bool) (models.Seal_count, error) {
	dbpool := pgclient.WDB
	gs := models.Seal{}

	gsc := 0
	err := dbpool.QueryRow(ctx, "SELECT * from func_seals_cnt($1);", gs1).Scan(&gsc)

	if err != nil {
		log.Println(err.Error(), "func_seals_cnt")
		return models.Seal_count{Values: []models.Seal{}, Count: gsc, Auth: models.Auth{}}, err
	}

	out_arr := make([]models.Seal, 0,
		func() int {
			if gsc < pgs {
				return gsc
			} else {
				return pgs
			}
		}())

	rows, err := dbpool.Query(ctx, "SELECT * from func_seals_get($1,$2,$3,$4,$5);", pg, pgs, gs1, ord, dsc)
	if err != nil {
		log.Println(err.Error(), "func_seals_get")
		return models.Seal_count{Values: []models.Seal{}, Count: gsc, Auth: models.Auth{}}, err
	}

	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&gs.Id, &gs.PacketNumber, &gs.Area.Id, &gs.Staff.Id, &gs.SealType.Id, &gs.SealColour.Id, &gs.SealStatus.Id,
			&gs.IssueDate, &gs.ReportDate, &gs.Area.AreaName, &gs.Staff.StaffName, &gs.SealType.SealTypeName, &gs.SealColour.SealColourName,
			&gs.SealStatus.SealStatusName)
		if err != nil {
			log.Println("failed to scan row:", err)
		}

		out_arr = append(out_arr, gs)
	}

	out_count := models.Seal_count{Values: out_arr, Count: gsc, Auth: models.Auth{}}
	if err != nil {
		log.Println(err.Error())
		return models.Seal_count{}, err
	}

	return out_count, nil
}

//func (est *SealStorage) Add(ctx context.Context, a models.Seal) (int, error)
func (est *SealStorage) Add(ctx context.Context, a models.Seal) (int, error) {
	dbpool := pgclient.WDB
	ai := 0

	err := dbpool.QueryRow(ctx, "SELECT func_seals_add($1,$2,$3,$4,$5,$6,$7,$8);", a.PacketNumber, a.Area.Id,
		a.Staff.Id, a.SealType.Id, a.SealColour.Id, a.SealStatus.Id, a.IssueDate, a.ReportDate).Scan(&ai)

	if err != nil {
		log.Println("Failed execute func_seals_add: ", err)
		return 0, err
	}

	return ai, nil
}

//func (est *SealStorage) Upd(ctx context.Context, u models.Seal) (int, error)
func (est *SealStorage) Upd(ctx context.Context, u models.Seal) (int, error) {
	dbpool := pgclient.WDB
	ui := 0

	err := dbpool.QueryRow(ctx, "SELECT func_seals_upd($1,$2,$3,$4,$5,$6,$7,$8,$9);", u.Id, u.PacketNumber, u.Area.Id,
		u.Staff.Id, u.SealType.Id, u.SealColour.Id, u.SealStatus.Id, u.IssueDate, u.ReportDate).Scan(&ui)

	if err != nil {
		log.Println("Failed execute func_seals_upd: ", err)
		return 0, err
	}
	return ui, nil
}

//func (est *SealStorage) Del(ctx context.Context, d []int) ([]int, error)
func (est *SealStorage) Del(ctx context.Context, d []int) ([]int, error) {
	dbpool := pgclient.WDB
	res := []int{}
	i := 0
	for _, id := range d {
		err := dbpool.QueryRow(ctx, "SELECT func_seals_del($1);", id).Scan(&i)
		res = append(res, i)

		if err != nil {
			log.Println("Failed execute func_seals_del: ", err)
		}
	}
	return res, nil
}

//func (est *Storage) GetOne(ctx context.Context, i int) (models.Seal_count, error)
func (est *SealStorage) GetOne(ctx context.Context, i int) (models.Seal_count, error) {
	dbpool := pgclient.WDB
	out_arr := []models.Seal{}
	g := models.Seal{}

	err := dbpool.QueryRow(ctx, "SELECT * from func_seal_get($1);", i).Scan(&g.Id, &g.PacketNumber, &g.Area.Id,
		&g.Staff.Id, &g.SealType.Id, &g.SealColour.Id, &g.SealStatus.Id, &g.IssueDate, &g.ReportDate, &g.Area.AreaName, &g.Staff.StaffName,
		&g.SealType.SealTypeName, &g.SealColour.SealColourName, &g.SealStatus.SealStatusName)

	if err != nil && err != pgx.ErrNoRows {
		log.Println("Failed execute func_seal_get: ", err)
		return models.Seal_count{Values: []models.Seal{}, Count: 0, Auth: models.Auth{}}, err
	}

	out_arr = append(out_arr, g)

	out_count := models.Seal_count{Values: out_arr, Count: 1, Auth: models.Auth{}}
	return out_count, nil
}
