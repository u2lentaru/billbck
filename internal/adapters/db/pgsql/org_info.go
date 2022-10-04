package pgsql

import (
	"context"
	"log"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/u2lentaru/billbck/internal/models"
	"github.com/u2lentaru/billbck/pkg/pgclient"
)

//type OrgInfoStorage struct
type OrgInfoStorage struct {
	db *pgxpool.Pool
}

//func NewOrgInfoStorage(db *pgxpool.Pool) *OrgInfoStorage
func NewOrgInfoStorage(db *pgxpool.Pool) *OrgInfoStorage {
	return &OrgInfoStorage{db: db}
}

//func (est *OrgInfoStorage) GetList(ctx context.Context, pg, pgs int, gs1, gs2 string, ord int, dsc bool) (models.OrgInfo_count, error) {
func (est *OrgInfoStorage) GetList(ctx context.Context, pg, pgs int, gs1, gs2 string, ord int, dsc bool) (models.OrgInfo_count, error) {
	dbpool := pgclient.WDB
	gs := models.OrgInfo{}

	gsc := 0
	err := dbpool.QueryRow(ctx, "SELECT * from func_orgs_info_cnt($1,$2);", gs1, gs2).Scan(&gsc)

	if err != nil {
		log.Println(err.Error(), "func_orgs_info_cnt")
		return models.OrgInfo_count{Values: []models.OrgInfo{}, Count: gsc, Auth: models.Auth{}}, err
	}

	out_arr := make([]models.OrgInfo, 0,
		func() int {
			if gsc < pgs {
				return gsc
			} else {
				return pgs
			}
		}())

	rows, err := dbpool.Query(ctx, "SELECT * from func_orgs_info_get($1,$2,$3,$4,$5,$6);", pg, pgs, gs1, gs2, ord, dsc)
	if err != nil {
		log.Println(err.Error())
		return models.OrgInfo_count{Values: []models.OrgInfo{}, Count: gsc, Auth: models.Auth{}}, err
	}

	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&gs.Id, &gs.OIName, &gs.OIBin, &gs.OIAddr, &gs.OIBank.Id, &gs.OIAccNumber, &gs.OIFName, &gs.OIBank.BankName)
		if err != nil {
			log.Println("failed to scan row:", err)
		}

		out_arr = append(out_arr, gs)
	}

	out_count := models.OrgInfo_count{Values: out_arr, Count: gsc, Auth: models.Auth{}}
	if err != nil {
		log.Println(err.Error())
		return models.OrgInfo_count{}, err
	}

	return out_count, nil
}

//func (est *OrgInfoStorage) Add(ctx context.Context, ea models.OrgInfo) (int, error)
func (est *OrgInfoStorage) Add(ctx context.Context, a models.OrgInfo) (int, error) {
	dbpool := pgclient.WDB
	ai := 0

	err := dbpool.QueryRow(ctx, "SELECT func_orgs_info_add($1,$2,$3,$4,$5,$6);", a.OIName, a.OIBin, a.OIAddr, a.OIBank.Id, a.OIAccNumber, a.OIFName).Scan(&ai)

	if err != nil {
		log.Println("Failed execute func_orgs_info_add: ", err)
		return 0, err
	}

	return ai, nil
}

//func (est *OrgInfoStorage) Upd(ctx context.Context, eu models.OrgInfo) (int, error)
func (est *OrgInfoStorage) Upd(ctx context.Context, u models.OrgInfo) (int, error) {
	dbpool := pgclient.WDB
	ui := 0

	err := dbpool.QueryRow(ctx, "SELECT func_orgs_info_upd($1,$2,$3,$4,$5,$6,$7);", u.Id, u.OIName, u.OIBin, u.OIAddr, u.OIBank.Id, u.OIAccNumber, u.OIFName).Scan(&ui)

	if err != nil {
		log.Println("Failed execute func_orgs_info_upd: ", err)
		return 0, err
	}
	return ui, nil
}

//func (est *OrgInfoStorage) Del(ctx context.Context, ed []int) ([]int, error)
func (est *OrgInfoStorage) Del(ctx context.Context, d []int) ([]int, error) {
	dbpool := pgclient.WDB
	res := []int{}
	i := 0
	for _, id := range d {
		err := dbpool.QueryRow(ctx, "SELECT func_orgs_info_del($1);", id).Scan(&i)
		res = append(res, i)

		if err != nil {
			log.Println("Failed execute func_orgs_info_del: ", err)
		}
	}
	return res, nil
}

//func (est *OrgInfoStorage) GetOne(ctx context.Context, i int) (models.OrgInfo_count, error)
func (est *OrgInfoStorage) GetOne(ctx context.Context, i int) (models.OrgInfo_count, error) {
	dbpool := pgclient.WDB
	out_arr := []models.OrgInfo{}
	g := models.OrgInfo{}

	err := dbpool.QueryRow(ctx, "SELECT * from func_org_info_get($1);", i).Scan(&g.Id, &g.OIName, &g.OIBin, &g.OIAddr, &g.OIBank.Id, &g.OIAccNumber, &g.OIFName, &g.OIBank.BankName)

	if err != nil && err != pgx.ErrNoRows {
		log.Println("Failed execute from func_org_info_get: ", err)
		return models.OrgInfo_count{Values: []models.OrgInfo{}, Count: 0, Auth: models.Auth{}}, err
	}

	out_arr = append(out_arr, g)

	out_count := models.OrgInfo_count{Values: out_arr, Count: 1, Auth: models.Auth{}}
	return out_count, nil
}
