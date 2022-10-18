package pgsql

import (
	"context"
	"log"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/u2lentaru/billbck/internal/models"
	"github.com/u2lentaru/billbck/pkg/pgclient"
)

//type FormTypeStorage struct
type FormTypeStorage struct {
	db *pgxpool.Pool
}

//func NewFormTypeStorage(db *pgxpool.Pool) *FormTypeStorage
func NewFormTypeStorage(db *pgxpool.Pool) *FormTypeStorage {
	return &FormTypeStorage{db: db}
}

//func (est *FormTypeStorage) GetList(ctx context.Context, pg, pgs int, gs1, gs2 string, ord int, dsc bool) (models.FormType_count, error) {
func (est *FormTypeStorage) GetList(ctx context.Context, pg, pgs int, gs1, gs2 string, ord int, dsc bool) (models.FormType_count, error) {
	dbpool := pgclient.WDB
	gs := models.FormType{}

	gsc := 0
	err := dbpool.QueryRow(ctx, "SELECT * from func_cnt_form_types_flt($1,$2);", gs1, gs2).Scan(&gsc)

	if err != nil {
		log.Println(err.Error(), "func_cnt_form_types_flt")
		return models.FormType_count{Values: []models.FormType{}, Count: gsc, Auth: models.Auth{}}, err
	}

	out_arr := make([]models.FormType, 0,
		func() int {
			if gsc < pgs {
				return gsc
			} else {
				return pgs
			}
		}())

	rows, err := dbpool.Query(ctx, "SELECT * from func_get_form_types_flt($1,$2,$3,$4,$5,$6);", pg, pgs, gs1, gs2, ord, dsc)
	if err != nil {
		log.Println(err.Error())
		return models.FormType_count{Values: []models.FormType{}, Count: gsc, Auth: models.Auth{}}, err
	}

	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&gs.Id, &gs.FormTypeName, &gs.FormTypeDescr)
		if err != nil {
			log.Println("failed to scan row:", err)
		}

		out_arr = append(out_arr, gs)
	}

	out_count := models.FormType_count{Values: out_arr, Count: gsc, Auth: models.Auth{}}
	if err != nil {
		log.Println(err.Error())
		return models.FormType_count{}, err
	}

	return out_count, nil
}

//func (est *FormTypeStorage) Add(ctx context.Context, ea models.FormType) (int, error)
func (est *FormTypeStorage) Add(ctx context.Context, a models.FormType) (int, error) {
	dbpool := pgclient.WDB
	ai := 0

	err := dbpool.QueryRow(ctx, "SELECT func_add_form_type($1,$2);", a.FormTypeName, a.FormTypeDescr).Scan(&ai)

	if err != nil {
		log.Println("Failed execute func_add_form_type: ", err)
		return 0, err
	}

	return ai, nil
}

//func (est *FormTypeStorage) Upd(ctx context.Context, eu models.FormType) (int, error)
func (est *FormTypeStorage) Upd(ctx context.Context, u models.FormType) (int, error) {
	dbpool := pgclient.WDB
	ui := 0

	err := dbpool.QueryRow(ctx, "SELECT func_upd_form_type($1,$2,$3);", u.Id, u.FormTypeName, u.FormTypeDescr).Scan(&ui)

	if err != nil {
		log.Println("Failed execute func_upd_form_type: ", err)
		return 0, err
	}
	return ui, nil
}

//func (est *FormTypeStorage) Del(ctx context.Context, ed []int) ([]int, error)
func (est *FormTypeStorage) Del(ctx context.Context, d []int) ([]int, error) {
	dbpool := pgclient.WDB
	res := []int{}
	i := 0
	for _, id := range d {
		err := dbpool.QueryRow(ctx, "SELECT func_del_form_type($1);", id).Scan(&i)
		res = append(res, i)

		if err != nil {
			log.Println("Failed execute func_del_form_type: ", err)
		}
	}
	return res, nil
}

//func (est *FormTypeStorage) GetOne(ctx context.Context, i int) (models.FormType_count, error)
func (est *FormTypeStorage) GetOne(ctx context.Context, i int) (models.FormType_count, error) {
	dbpool := pgclient.WDB
	out_arr := []models.FormType{}
	g := models.FormType{}

	err := dbpool.QueryRow(ctx, "SELECT * from func_get_form_type($1);", i).Scan(&g.Id, &g.FormTypeName, &g.FormTypeDescr)

	if err != nil && err != pgx.ErrNoRows {
		log.Println("Failed execute from func_get_form_type: ", err)
		return models.FormType_count{Values: []models.FormType{}, Count: 0, Auth: models.Auth{}}, err
	}

	out_arr = append(out_arr, g)

	out_count := models.FormType_count{Values: out_arr, Count: 1, Auth: models.Auth{}}
	return out_count, nil
}
