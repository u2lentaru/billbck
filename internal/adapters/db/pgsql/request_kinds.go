package pgsql

import (
	"context"
	"log"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/u2lentaru/billbck/internal/models"
	"github.com/u2lentaru/billbck/pkg/pgclient"
)

//type RequestKindStorage struct
type RequestKindStorage struct {
	db *pgxpool.Pool
}

//func NewRequestKindStorage(db *pgxpool.Pool) *RequestKindStorage
func NewRequestKindStorage(db *pgxpool.Pool) *RequestKindStorage {
	return &RequestKindStorage{db: db}
}

//func (est *RequestKindStorage) GetList(ctx context.Context, pg, pgs int, gs1 string, ord int, dsc bool) (models.RequestKind_count, error)
func (est *RequestKindStorage) GetList(ctx context.Context, pg, pgs int, gs1 string, ord int, dsc bool) (models.RequestKind_count, error) {
	dbpool := pgclient.WDB
	gs := models.RequestKind{}

	gsc := 0
	err := dbpool.QueryRow(ctx, "SELECT * from func_request_kinds_cnt($1);", gs1).Scan(&gsc)

	if err != nil {
		log.Println(err.Error(), "func_request_kinds_cnt")
		return models.RequestKind_count{Values: []models.RequestKind{}, Count: gsc, Auth: models.Auth{}}, err
	}

	out_arr := make([]models.RequestKind, 0,
		func() int {
			if gsc < pgs {
				return gsc
			} else {
				return pgs
			}
		}())

	rows, err := dbpool.Query(ctx, "SELECT * from func_request_kinds_get($1,$2,$3,$4,$5);", pg, pgs, gs1, ord, dsc)
	if err != nil {
		log.Println(err.Error(), "func_request_kinds_get")
		return models.RequestKind_count{Values: []models.RequestKind{}, Count: gsc, Auth: models.Auth{}}, err
	}

	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&gs.Id, &gs.RequestKindName)
		if err != nil {
			log.Println("failed to scan row:", err)
		}

		out_arr = append(out_arr, gs)
	}

	out_count := models.RequestKind_count{Values: out_arr, Count: gsc, Auth: models.Auth{}}
	return out_count, nil
}

//func (est *RequestKindStorage) Add(ctx context.Context, a models.RequestKind) (int, error)
func (est *RequestKindStorage) Add(ctx context.Context, a models.RequestKind) (int, error) {
	dbpool := pgclient.WDB
	ai := 0

	err := dbpool.QueryRow(ctx, "SELECT func_request_kinds_add($1);", a.RequestKindName).Scan(&ai)

	if err != nil {
		log.Println("Failed execute func_request_kinds_add: ", err)
		return 0, err
	}

	return ai, nil
}

//func (est *RequestKindStorage) Upd(ctx context.Context, u models.RequestKind) (int, error)
func (est *RequestKindStorage) Upd(ctx context.Context, u models.RequestKind) (int, error) {
	dbpool := pgclient.WDB
	ui := 0

	err := dbpool.QueryRow(ctx, "SELECT func_request_kinds_upd($1,$2);", u.Id, u.RequestKindName).Scan(&ui)

	if err != nil {
		log.Println("Failed execute func_request_kinds_upd: ", err)
		return 0, err
	}
	return ui, nil
}

//func (est *RequestKindStorage) Del(ctx context.Context, d []int) ([]int, error)
func (est *RequestKindStorage) Del(ctx context.Context, d []int) ([]int, error) {
	dbpool := pgclient.WDB
	res := []int{}
	i := 0
	for _, id := range d {
		err := dbpool.QueryRow(ctx, "SELECT func_request_kinds_del($1);", id).Scan(&i)
		res = append(res, i)

		if err != nil {
			log.Println("Failed execute func_request_kinds_del: ", err)
		}
	}
	return res, nil
}

//func (est *RequestKindStorage) GetOne(ctx context.Context, i int) (models.RequestKind_count, error)
func (est *RequestKindStorage) GetOne(ctx context.Context, i int) (models.RequestKind_count, error) {
	dbpool := pgclient.WDB
	out_arr := []models.RequestKind{}
	g := models.RequestKind{}

	err := dbpool.QueryRow(ctx, "SELECT * from func_request_kind_get($1);", i).Scan(&g.Id, &g.RequestKindName)

	if err != nil && err != pgx.ErrNoRows {
		log.Println("Failed execute func_request_kind_get: ", err)
		return models.RequestKind_count{Values: []models.RequestKind{}, Count: 0, Auth: models.Auth{}}, err
	}

	out_arr = append(out_arr, g)

	out_count := models.RequestKind_count{Values: out_arr, Count: 1, Auth: models.Auth{}}
	return out_count, nil
}
