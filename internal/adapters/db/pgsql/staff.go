package pgsql

import (
	"context"
	"log"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/u2lentaru/billbck/internal/models"
	"github.com/u2lentaru/billbck/pkg/pgclient"
)

//type StaffStorage struct
type StaffStorage struct {
	db *pgxpool.Pool
}

//func NewStaffStorage(db *pgxpool.Pool) *StaffStorage
func NewStaffStorage(db *pgxpool.Pool) *StaffStorage {
	return &StaffStorage{db: db}
}

//func (est *StaffStorage) GetList(ctx context.Context, pg, pgs int, gs1 string, ord int, dsc bool) (models.Staff_count, error)
func (est *StaffStorage) GetList(ctx context.Context, pg, pgs int, gs1 string, ord int, dsc bool) (models.Staff_count, error) {
	dbpool := pgclient.WDB
	gs := models.Staff{}

	gsc := 0
	err := dbpool.QueryRow(ctx, "SELECT * from func_staff_cnt($1);", gs1).Scan(&gsc)

	if err != nil {
		log.Println(err.Error(), "func_staff_cnt")
		return models.Staff_count{Values: []models.Staff{}, Count: gsc, Auth: models.Auth{}}, err
	}

	out_arr := make([]models.Staff, 0,
		func() int {
			if gsc < pgs {
				return gsc
			} else {
				return pgs
			}
		}())

	rows, err := dbpool.Query(ctx, "SELECT * from func_staff_get($1,$2,$3,$4,$5);", pg, pgs, gs1, ord, dsc)
	if err != nil {
		log.Println(err.Error(), "func_staff_get")
		return models.Staff_count{Values: []models.Staff{}, Count: gsc, Auth: models.Auth{}}, err
	}

	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&gs.Id, &gs.StaffName, &gs.OrgInfo.Id, &gs.Phone, &gs.Notes, &gs.OrgInfo.OIName)
		if err != nil {
			log.Println("failed to scan row:", err)
		}

		out_arr = append(out_arr, gs)
	}

	out_count := models.Staff_count{Values: out_arr, Count: gsc, Auth: models.Auth{}}
	return out_count, nil
}

//func (est *StaffStorage) Add(ctx context.Context, a models.Staff) (int, error)
func (est *StaffStorage) Add(ctx context.Context, a models.Staff) (int, error) {
	dbpool := pgclient.WDB
	ai := 0

	err := dbpool.QueryRow(ctx, "SELECT func_staff_add($1,$2,$3,$4);", a.StaffName, a.OrgInfo.Id, a.Phone, a.Notes).Scan(&ai)

	if err != nil {
		log.Println("Failed execute func_staff_add: ", err)
		return 0, err
	}

	return ai, nil
}

//func (est *StaffStorage) Upd(ctx context.Context, u models.Staff) (int, error)
func (est *StaffStorage) Upd(ctx context.Context, u models.Staff) (int, error) {
	dbpool := pgclient.WDB
	ui := 0

	err := dbpool.QueryRow(ctx, "SELECT func_staff_upd($1,$2,$3,$4,$5);", u.Id, u.StaffName, u.OrgInfo.Id, u.Phone,
		u.Notes).Scan(&ui)

	if err != nil {
		log.Println("Failed execute func_staff_upd: ", err)
		return 0, err
	}
	return ui, nil
}

//func (est *StaffStorage) Del(ctx context.Context, d []int) ([]int, error)
func (est *StaffStorage) Del(ctx context.Context, d []int) ([]int, error) {
	dbpool := pgclient.WDB
	res := []int{}
	i := 0
	for _, id := range d {
		err := dbpool.QueryRow(ctx, "SELECT func_staff_del($1);", id).Scan(&i)
		res = append(res, i)

		if err != nil {
			log.Println("Failed execute func_staff_del: ", err)
		}
	}
	return res, nil
}

//func (est *StaffStorage) GetOne(ctx context.Context, i int) (models.Staff_count, error)
func (est *StaffStorage) GetOne(ctx context.Context, i int) (models.Staff_count, error) {
	dbpool := pgclient.WDB
	out_arr := []models.Staff{}
	g := models.Staff{}

	err := dbpool.QueryRow(ctx, "SELECT * from func_staff_getbyid($1);", i).Scan(&g.Id, &g.StaffName, &g.OrgInfo.Id,
		&g.Phone, &g.Notes, &g.OrgInfo.OIName)

	if err != nil && err != pgx.ErrNoRows {
		log.Println("Failed execute func_staff_getbyid: ", err)
		return models.Staff_count{Values: []models.Staff{}, Count: 0, Auth: models.Auth{}}, err
	}

	out_arr = append(out_arr, g)

	out_count := models.Staff_count{Values: out_arr, Count: 1, Auth: models.Auth{}}
	return out_count, nil
}
