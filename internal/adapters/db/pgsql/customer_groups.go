package pgsql

import (
	"context"
	"log"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/u2lentaru/billbck/internal/models"
	"github.com/u2lentaru/billbck/pkg/pgclient"
)

//type CustomerGroupStorage struct
type CustomerGroupStorage struct {
	db *pgxpool.Pool
}

//func NewCustomerGroupStorage(db *pgxpool.Pool) *CustomerGroupStorage
func NewCustomerGroupStorage(db *pgxpool.Pool) *CustomerGroupStorage {
	return &CustomerGroupStorage{db: db}
}

//func (est *CustomerGroupStorage) GetList(ctx context.Context, pg, pgs int, gs1 string, ord int, dsc bool) (models.CustomerGroup_count, error)
func (est *CustomerGroupStorage) GetList(ctx context.Context, pg, pgs int, gs1 string, ord int, dsc bool) (models.CustomerGroup_count, error) {
	dbpool := pgclient.WDB
	gs := models.CustomerGroup{}

	gsc := 0
	err := dbpool.QueryRow(ctx, "SELECT * from func_customer_groups_cnt($1);", gs1).Scan(&gsc)

	if err != nil {
		log.Println(err.Error(), "func_customer_groups_cnt")
		return models.CustomerGroup_count{Values: []models.CustomerGroup{}, Count: gsc, Auth: models.Auth{}}, err
	}

	out_arr := make([]models.CustomerGroup, 0,
		func() int {
			if gsc < pgs {
				return gsc
			} else {
				return pgs
			}
		}())

	rows, err := dbpool.Query(ctx, "SELECT * from func_customer_groups_get($1,$2,$3,$4,$5);", pg, pgs, gs1, ord, dsc)
	if err != nil {
		log.Println(err.Error())
		return models.CustomerGroup_count{Values: []models.CustomerGroup{}, Count: gsc, Auth: models.Auth{}}, err
	}

	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&gs.Id, &gs.CustomerGroupName)
		if err != nil {
			log.Println("failed to scan row:", err)
		}

		out_arr = append(out_arr, gs)
	}

	out_count := models.CustomerGroup_count{Values: out_arr, Count: gsc, Auth: models.Auth{}}
	if err != nil {
		log.Println(err.Error())
		return models.CustomerGroup_count{}, err
	}

	return out_count, nil
}

//func (est *CustomerGroupStorage) Add(ctx context.Context, a models.CustomerGroup) (int, error)
func (est *CustomerGroupStorage) Add(ctx context.Context, a models.CustomerGroup) (int, error) {
	dbpool := pgclient.WDB
	ai := 0

	err := dbpool.QueryRow(ctx, "SELECT func_customer_groups_add($1);", a.CustomerGroupName).Scan(&ai)

	if err != nil {
		log.Println("Failed execute func_customer_groups_add: ", err)
		return 0, err
	}

	return ai, nil
}

//func (est *CustomerGroupStorage) Upd(ctx context.Context, u models.CustomerGroup) (int, error)
func (est *CustomerGroupStorage) Upd(ctx context.Context, u models.CustomerGroup) (int, error) {
	dbpool := pgclient.WDB
	ui := 0

	err := dbpool.QueryRow(ctx, "SELECT func_customer_groups_upd($1,$2);", u.Id, u.CustomerGroupName).Scan(&ui)

	if err != nil {
		log.Println("Failed execute func_customer_groups_upd: ", err)
		return 0, err
	}
	return ui, nil
}

//func (est *CustomerGroupStorage) Del(ctx context.Context, d []int) ([]int, error)
func (est *CustomerGroupStorage) Del(ctx context.Context, d []int) ([]int, error) {
	dbpool := pgclient.WDB
	res := []int{}
	i := 0
	for _, id := range d {
		err := dbpool.QueryRow(ctx, "SELECT func_customer_groups_del($1);", id).Scan(&i)
		res = append(res, i)

		if err != nil {
			log.Println("Failed execute func_customer_groups_del: ", err)
		}
	}
	return res, nil
}

//func (est *CustomerGroupStorage) GetOne(ctx context.Context, i int) (models.CustomerGroup_count, error)
func (est *CustomerGroupStorage) GetOne(ctx context.Context, i int) (models.CustomerGroup_count, error) {
	dbpool := pgclient.WDB
	out_arr := []models.CustomerGroup{}
	g := models.CustomerGroup{}

	err := dbpool.QueryRow(ctx, "SELECT * from func_customer_group_get($1);", i).Scan(&g.Id, &g.CustomerGroupName)

	if err != nil && err != pgx.ErrNoRows {
		log.Println("Failed execute from func_customer_group_get: ", err)
		return models.CustomerGroup_count{Values: []models.CustomerGroup{}, Count: 0, Auth: models.Auth{}}, err
	}

	out_arr = append(out_arr, g)

	out_count := models.CustomerGroup_count{Values: out_arr, Count: 1, Auth: models.Auth{}}
	return out_count, nil
}
