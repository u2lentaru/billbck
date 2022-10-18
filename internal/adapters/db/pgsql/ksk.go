package pgsql

import (
	"context"
	"log"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/u2lentaru/billbck/internal/models"
	"github.com/u2lentaru/billbck/pkg/pgclient"
)

//type KskStorage struct
type KskStorage struct {
	db *pgxpool.Pool
}

//func NewKskStorage(db *pgxpool.Pool) *KskStorage
func NewKskStorage(db *pgxpool.Pool) *KskStorage {
	return &KskStorage{db: db}
}

//func (est *KskStorage) GetList(ctx context.Context, pg, pgs int, gs1, gs2 string, ord int, dsc bool) (models.Ksk_count, error) {
func (est *KskStorage) GetList(ctx context.Context, pg, pgs int, gs1, gs2 string, ord int, dsc bool) (models.Ksk_count, error) {
	dbpool := pgclient.WDB
	gs := models.Ksk{}

	gsc := 0
	err := dbpool.QueryRow(ctx, "SELECT * from func_ksks_cnt($1,$2);", gs1, gs2).Scan(&gsc)

	if err != nil {
		log.Println(err.Error(), "func_ksks_cnt")
		return models.Ksk_count{Values: []models.Ksk{}, Count: gsc, Auth: models.Auth{}}, err
	}

	out_arr := make([]models.Ksk, 0,
		func() int {
			if gsc < pgs {
				return gsc
			} else {
				return pgs
			}
		}())

	rows, err := dbpool.Query(ctx, "SELECT * from func_ksks_get($1,$2,$3,$4,$5,$6);", pg, pgs, gs1, gs2, ord, dsc)
	if err != nil {
		log.Println(err.Error())
		return models.Ksk_count{Values: []models.Ksk{}, Count: gsc, Auth: models.Auth{}}, err
	}

	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&gs.Id, &gs.KskName, &gs.KskAddress, &gs.KskHead, &gs.KskPhone)
		if err != nil {
			log.Println("failed to scan row:", err)
		}

		out_arr = append(out_arr, gs)
	}

	out_count := models.Ksk_count{Values: out_arr, Count: gsc, Auth: models.Auth{}}
	if err != nil {
		log.Println(err.Error())
		return models.Ksk_count{}, err
	}

	return out_count, nil
}

//func (est *KskStorage) Add(ctx context.Context, ea models.Ksk) (int, error)
func (est *KskStorage) Add(ctx context.Context, a models.Ksk) (int, error) {
	dbpool := pgclient.WDB
	ai := 0

	err := dbpool.QueryRow(ctx, "SELECT func_ksks_add($1,$2,$3,$4);", a.KskName, a.KskAddress, a.KskHead, a.KskPhone).Scan(&ai)

	if err != nil {
		log.Println("Failed execute func_ksks_add: ", err)
		return 0, err
	}

	return ai, nil
}

//func (est *KskStorage) Upd(ctx context.Context, eu models.Ksk) (int, error)
func (est *KskStorage) Upd(ctx context.Context, u models.Ksk) (int, error) {
	dbpool := pgclient.WDB
	ui := 0

	err := dbpool.QueryRow(ctx, "SELECT func_ksks_upd($1,$2,$3,$4,$5);", u.Id, u.KskName, u.KskAddress, u.KskHead, u.KskPhone).Scan(&ui)

	if err != nil {
		log.Println("Failed execute func_ksks_upd: ", err)
		return 0, err
	}
	return ui, nil
}

//func (est *KskStorage) Del(ctx context.Context, ed []int) ([]int, error)
func (est *KskStorage) Del(ctx context.Context, d []int) ([]int, error) {
	dbpool := pgclient.WDB
	res := []int{}
	i := 0
	for _, id := range d {
		err := dbpool.QueryRow(ctx, "SELECT func_ksks_del($1);", id).Scan(&i)
		res = append(res, i)

		if err != nil {
			log.Println("Failed execute func_ksks_del: ", err)
		}
	}
	return res, nil
}

//func (est *KskStorage) GetOne(ctx context.Context, i int) (models.Ksk_count, error)
func (est *KskStorage) GetOne(ctx context.Context, i int) (models.Ksk_count, error) {
	dbpool := pgclient.WDB
	out_arr := []models.Ksk{}
	g := models.Ksk{}

	err := dbpool.QueryRow(ctx, "SELECT * from func_ksk_get($1);", i).Scan(&g.Id, &g.KskName, &g.KskAddress, &g.KskHead, &g.KskPhone)

	if err != nil && err != pgx.ErrNoRows {
		log.Println("Failed execute from func_ksk_get: ", err)
		return models.Ksk_count{Values: []models.Ksk{}, Count: 0, Auth: models.Auth{}}, err
	}

	out_arr = append(out_arr, g)

	out_count := models.Ksk_count{Values: out_arr, Count: 1, Auth: models.Auth{}}
	return out_count, nil
}
