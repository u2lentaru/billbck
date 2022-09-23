package pgsql

import (
	"context"
	"log"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/u2lentaru/billbck/internal/models"
	"github.com/u2lentaru/billbck/pkg/pgclient"
)

//type AskueTypeStorage struct
type AskueTypeStorage struct {
	db *pgxpool.Pool
}

//func NewAskueTypeStorage(db *pgxpool.Pool) *AskueTypeStorage
func NewAskueTypeStorage(db *pgxpool.Pool) *AskueTypeStorage {
	return &AskueTypeStorage{db: db}
}

//func (est *AskueTypeStorage) GetList(ctx context.Context, pg, pgs int, gs1 string, ord int, dsc bool) (models.AskueType_count, error)
func (est *AskueTypeStorage) GetList(ctx context.Context, pg, pgs int, gs1 string, ord int, dsc bool) (models.AskueType_count, error) {
	dbpool := pgclient.WDB
	gs := models.AskueType{}

	gsc := 0
	err := dbpool.QueryRow(ctx, "SELECT * from func_askue_types_cnt($1);", gs1).Scan(&gsc)

	if err != nil {
		log.Println(err.Error(), "func_askue_types_cnt")
		return models.AskueType_count{Values: []models.AskueType{}, Count: gsc, Auth: models.Auth{}}, err
	}

	out_arr := make([]models.AskueType, 0,
		func() int {
			if gsc < pgs {
				return gsc
			} else {
				return pgs
			}
		}())

	rows, err := dbpool.Query(ctx, "SELECT * from func_askue_types_get($1,$2,$3,$4,$5);", pg, pgs, gs1, ord, dsc)
	if err != nil {
		log.Println(err.Error())
		return models.AskueType_count{Values: []models.AskueType{}, Count: gsc, Auth: models.Auth{}}, err
	}

	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&gs.Id, &gs.AskueTypeName, &gs.StartLine, &gs.PuColumn, &gs.ValueColumn, &gs.DateColumn, &gs.DateColumnArray)
		if err != nil {
			log.Println("failed to scan row:", err)
		}

		out_arr = append(out_arr, gs)
	}

	out_count := models.AskueType_count{Values: out_arr, Count: gsc, Auth: models.Auth{}}
	if err != nil {
		log.Println(err.Error())
		return models.AskueType_count{}, err
	}

	return out_count, nil
}

//func (est *AskueTypeStorage) Add(ctx context.Context, a models.AskueType) (int, error)
func (est *AskueTypeStorage) Add(ctx context.Context, a models.AskueType) (int, error) {
	dbpool := pgclient.WDB
	ai := 0

	err := dbpool.QueryRow(context.Background(), "SELECT func_askue_types_add($1,$2,$3,$4,$5,$6);", a.AskueTypeName, a.StartLine,
		a.PuColumn, a.ValueColumn, a.DateColumn, a.DateColumnArray).Scan(&ai)

	if err != nil {
		log.Println("Failed execute func_askue_types_add: ", err)
		return 0, err
	}

	return ai, nil
}

//func (est *AskueTypeStorage) Upd(ctx context.Context, u models.AskueType) (int, error)
func (est *AskueTypeStorage) Upd(ctx context.Context, u models.AskueType) (int, error) {
	dbpool := pgclient.WDB
	ui := 0

	err := dbpool.QueryRow(context.Background(), "SELECT func_askue_types_upd($1,$2,$3,$4,$5,$6,$7);", u.Id, u.AskueTypeName,
		u.StartLine, u.PuColumn, u.ValueColumn, u.DateColumn, u.DateColumnArray).Scan(&ui)

	if err != nil {
		log.Println("Failed execute func_askue_types_upd: ", err)
		return 0, err
	}
	return ui, nil
}

//func (est *AskueTypeStorage) Del(ctx context.Context, d []int) ([]int, error)
func (est *AskueTypeStorage) Del(ctx context.Context, d []int) ([]int, error) {
	dbpool := pgclient.WDB
	res := []int{}
	i := 0
	for _, id := range d {
		err := dbpool.QueryRow(context.Background(), "SELECT func_askue_types_del($1);", id).Scan(&i)
		res = append(res, i)

		if err != nil {
			log.Println("Failed execute func_askue_types_del: ", err)
		}
	}
	return res, nil
}

//func (est *AskueTypeStorage) GetOne(ctx context.Context, i int) (models.AskueType_count, error)
func (est *AskueTypeStorage) GetOne(ctx context.Context, i int) (models.AskueType_count, error) {
	dbpool := pgclient.WDB
	out_arr := []models.AskueType{}
	g := models.AskueType{}

	err := dbpool.QueryRow(context.Background(), "SELECT * from func_askue_type_get($1);", i).Scan(&g.Id, &g.AskueTypeName,
		&g.StartLine, &g.PuColumn, &g.ValueColumn, &g.DateColumn, &g.DateColumnArray)

	if err != nil && err != pgx.ErrNoRows {
		log.Println("Failed execute from func_askue_type_get: ", err)
		return models.AskueType_count{Values: []models.AskueType{}, Count: 0, Auth: models.Auth{}}, err
	}

	out_arr = append(out_arr, g)

	out_count := models.AskueType_count{Values: out_arr, Count: 1, Auth: models.Auth{}}
	return out_count, nil
}
