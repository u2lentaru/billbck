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

//type RequestTypeStorage struct
type RequestTypeStorage struct {
	db *pgxpool.Pool
}

//func NewRequestTypeStorage(db *pgxpool.Pool) *RequestTypeStorage
func NewRequestTypeStorage(db *pgxpool.Pool) *RequestTypeStorage {
	return &RequestTypeStorage{db: db}
}

//func (est *RequestTypeStorage) GetList(ctx context.Context, pg, pgs int, gs1, gs2 string, ord int, dsc bool) (models.RequestType_count, error)
func (est *RequestTypeStorage) GetList(ctx context.Context, pg, pgs int, gs1, gs2 string, ord int, dsc bool) (models.RequestType_count, error) {
	dbpool := pgclient.WDB
	gs := models.RequestType{}

	gsc := 0
	err := dbpool.QueryRow(ctx, "SELECT * from func_request_types_cnt($1,$2);", gs1, utils.NullableString(gs2)).Scan(&gsc)

	if err != nil {
		log.Println(err.Error(), "func_request_types_cnt")
		return models.RequestType_count{Values: []models.RequestType{}, Count: gsc, Auth: models.Auth{}}, err
	}

	out_arr := make([]models.RequestType, 0,
		func() int {
			if gsc < pgs {
				return gsc
			} else {
				return pgs
			}
		}())

	rows, err := dbpool.Query(ctx, "SELECT * from func_request_types_get($1,$2,$3,$4,$5,$6);", pg, pgs, gs1, utils.NullableString(gs2), ord, dsc)
	if err != nil {
		log.Println(err.Error(), "func_request_types_get")
		return models.RequestType_count{Values: []models.RequestType{}, Count: gsc, Auth: models.Auth{}}, err
	}

	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&gs.Id, &gs.RequestTypeName, &gs.RequestKind.Id, &gs.RequestKind.RequestKindName)
		if err != nil {
			log.Println("failed to scan row:", err)
		}

		out_arr = append(out_arr, gs)
	}

	out_count := models.RequestType_count{Values: out_arr, Count: gsc, Auth: models.Auth{}}
	if err != nil {
		log.Println(err.Error())
		return models.RequestType_count{}, err
	}

	return out_count, nil
}

//func (est *RequestTypeStorage) Add(ctx context.Context, a models.RequestType) (int, error)
func (est *RequestTypeStorage) Add(ctx context.Context, a models.RequestType) (int, error) {
	dbpool := pgclient.WDB
	ai := 0

	err := dbpool.QueryRow(ctx, "SELECT func_request_types_add($1,$2);", a.RequestTypeName, a.RequestKind.Id).Scan(&ai)

	if err != nil {
		log.Println("Failed execute func_request_types_add: ", err)
		return 0, err
	}

	return ai, nil
}

//func (est *RequestTypeStorage) Upd(ctx context.Context, u models.RequestType) (int, error)
func (est *RequestTypeStorage) Upd(ctx context.Context, u models.RequestType) (int, error) {
	dbpool := pgclient.WDB
	ui := 0

	err := dbpool.QueryRow(ctx, "SELECT func_request_types_upd($1,$2,$3);", u.Id, u.RequestTypeName, u.RequestKind.Id).Scan(&ui)

	if err != nil {
		log.Println("Failed execute func_request_types_upd: ", err)
		return 0, err
	}
	return ui, nil
}

//func (est *RequestTypeStorage) Del(ctx context.Context, d []int) ([]int, error)
func (est *RequestTypeStorage) Del(ctx context.Context, d []int) ([]int, error) {
	dbpool := pgclient.WDB
	res := []int{}
	i := 0
	for _, id := range d {
		err := dbpool.QueryRow(ctx, "SELECT func_request_types_del($1);", id).Scan(&i)
		res = append(res, i)

		if err != nil {
			log.Println("Failed execute func_request_types_del: ", err)
		}
	}
	return res, nil
}

//func (est *RequestTypeStorage) GetOne(ctx context.Context, i int) (models.RequestType_count, error)
func (est *RequestTypeStorage) GetOne(ctx context.Context, i int) (models.RequestType_count, error) {
	dbpool := pgclient.WDB
	out_arr := []models.RequestType{}
	g := models.RequestType{}

	err := dbpool.QueryRow(ctx, "SELECT * from func_request_type_get($1);", i).Scan(&g.Id, &g.RequestTypeName,
		&g.RequestKind.Id, &g.RequestKind.RequestKindName)

	if err != nil && err != pgx.ErrNoRows {
		log.Println("Failed execute func_request_type_get: ", err)
		return models.RequestType_count{Values: []models.RequestType{}, Count: 0, Auth: models.Auth{}}, err
	}

	out_arr = append(out_arr, g)

	out_count := models.RequestType_count{Values: out_arr, Count: 1, Auth: models.Auth{}}
	return out_count, nil
}
