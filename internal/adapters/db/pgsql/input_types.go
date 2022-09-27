package pgsql

import (
	"context"
	"log"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/u2lentaru/billbck/internal/models"
	"github.com/u2lentaru/billbck/pkg/pgclient"
)

//type InputTypeStorage struct
type InputTypeStorage struct {
	db *pgxpool.Pool
}

//func NewInputTypeStorage(db *pgxpool.Pool) *InputTypeStorage
func NewInputTypeStorage(db *pgxpool.Pool) *InputTypeStorage {
	return &InputTypeStorage{db: db}
}

//func (est *InputTypeStorage) GetList(ctx context.Context, pg, pgs int, gs1 string, ord int, dsc bool) (models.InputType_count, error)
func (est *InputTypeStorage) GetList(ctx context.Context, pg, pgs int, gs1 string, ord int, dsc bool) (models.InputType_count, error) {
	dbpool := pgclient.WDB
	gs := models.InputType{}

	gsc := 0
	err := dbpool.QueryRow(ctx, "SELECT * from func_input_types_cnt($1);", gs1).Scan(&gsc)

	if err != nil {
		log.Println(err.Error(), "func_input_types_cnt")
		return models.InputType_count{Values: []models.InputType{}, Count: gsc, Auth: models.Auth{}}, err
	}

	out_arr := make([]models.InputType, 0,
		func() int {
			if gsc < pgs {
				return gsc
			} else {
				return pgs
			}
		}())

	rows, err := dbpool.Query(ctx, "SELECT * from func_input_types_get($1,$2,$3,$4,$5);", pg, pgs, gs1, ord, dsc)
	if err != nil {
		log.Println(err.Error())
		return models.InputType_count{Values: []models.InputType{}, Count: gsc, Auth: models.Auth{}}, err
	}

	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&gs.Id, &gs.InputTypeName)
		if err != nil {
			log.Println("failed to scan row:", err)
		}

		out_arr = append(out_arr, gs)
	}

	out_count := models.InputType_count{Values: out_arr, Count: gsc, Auth: models.Auth{}}
	if err != nil {
		log.Println(err.Error())
		return models.InputType_count{}, err
	}

	return out_count, nil
}

//func (est *InputTypeStorage) Add(ctx context.Context, a models.InputType) (int, error)
func (est *InputTypeStorage) Add(ctx context.Context, a models.InputType) (int, error) {
	dbpool := pgclient.WDB
	ai := 0

	err := dbpool.QueryRow(context.Background(), "SELECT func_input_types_add($1);", a.InputTypeName).Scan(&ai)

	if err != nil {
		log.Println("Failed execute func_input_types_add: ", err)
		return 0, err
	}

	return ai, nil
}

//func (est *InputTypeStorage) Upd(ctx context.Context, u models.InputType) (int, error)
func (est *InputTypeStorage) Upd(ctx context.Context, u models.InputType) (int, error) {
	dbpool := pgclient.WDB
	ui := 0

	err := dbpool.QueryRow(context.Background(), "SELECT func_input_types_upd($1,$2);", u.Id, u.InputTypeName).Scan(&ui)

	if err != nil {
		log.Println("Failed execute func_input_types_upd: ", err)
		return 0, err
	}
	return ui, nil
}

//func (est *InputTypeStorage) Del(ctx context.Context, d []int) ([]int, error)
func (est *InputTypeStorage) Del(ctx context.Context, d []int) ([]int, error) {
	dbpool := pgclient.WDB
	res := []int{}
	i := 0
	for _, id := range d {
		err := dbpool.QueryRow(context.Background(), "SELECT func_input_types_del($1);", id).Scan(&i)
		res = append(res, i)

		if err != nil {
			log.Println("Failed execute func_input_types_del: ", err)
		}
	}
	return res, nil
}

//func (est *InputTypeStorage) GetOne(ctx context.Context, i int) (models.InputType_count, error)
func (est *InputTypeStorage) GetOne(ctx context.Context, i int) (models.InputType_count, error) {
	dbpool := pgclient.WDB
	out_arr := []models.InputType{}
	g := models.InputType{}

	err := dbpool.QueryRow(context.Background(), "SELECT * from func_input_type_get($1);", i).Scan(&g.Id, &g.InputTypeName)

	if err != nil && err != pgx.ErrNoRows {
		log.Println("Failed execute from func_input_type_get: ", err)
		return models.InputType_count{Values: []models.InputType{}, Count: 0, Auth: models.Auth{}}, err
	}

	out_arr = append(out_arr, g)

	out_count := models.InputType_count{Values: out_arr, Count: 1, Auth: models.Auth{}}
	return out_count, nil
}
