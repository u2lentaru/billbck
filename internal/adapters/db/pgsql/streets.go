package pgsql

import (
	"context"
	"log"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/u2lentaru/billbck/internal/models"
	"github.com/u2lentaru/billbck/pkg/pgclient"
)

//type StreetStorage struct
type StreetStorage struct {
	db *pgxpool.Pool
}

//func NewStreetStorage(db *pgxpool.Pool) *StreetStorage
func NewStreetStorage(db *pgxpool.Pool) *StreetStorage {
	return &StreetStorage{db: db}
}

//func (est *StreetStorage) GetList(ctx context.Context, pg, pgs int, gs1 string, gs2, ord int, dsc bool) (models.Street_count, error) {
func (est *StreetStorage) GetList(ctx context.Context, pg, pgs int, gs1 string, gs2, ord int, dsc bool) (models.Street_count, error) {
	dbpool := pgclient.WDB
	gs := models.Street{}

	gsc := 0
	err := dbpool.QueryRow(ctx, "SELECT * from func_streets_cnt($1,$2);", gs1, gs2).Scan(&gsc)

	if err != nil {
		log.Println(err.Error(), "func_streets_cnt")
		return models.Street_count{Values: []models.Street{}, Count: gsc, Auth: models.Auth{}}, err
	}

	out_arr := make([]models.Street, 0,
		func() int {
			if gsc < pgs {
				return gsc
			} else {
				return pgs
			}
		}())

	rows, err := dbpool.Query(ctx, "SELECT * from func_streets_get($1,$2,$3,$4,$5,$6);", pg, pgs, gs1, gs2, ord, dsc)
	if err != nil {
		log.Println(err.Error(), "func_streets_get")
		return models.Street_count{Values: []models.Street{}, Count: gsc, Auth: models.Auth{}}, err
	}

	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&gs.Id, &gs.StreetName, &gs.Created, &gs.Closed, &gs.City.CityName, &gs.City.Id)
		if err != nil {
			log.Println("failed to scan row:", err)
		}

		out_arr = append(out_arr, gs)
	}

	out_count := models.Street_count{Values: out_arr, Count: gsc, Auth: models.Auth{}}
	return out_count, nil
}

//func (est *StreetStorage) Add(ctx context.Context, ea models.Street) (int, error)
func (est *StreetStorage) Add(ctx context.Context, a models.Street) (int, error) {
	dbpool := pgclient.WDB
	ai := 0

	err := dbpool.QueryRow(ctx, "SELECT func_streets_add($1,$2,$3);", a.StreetName, a.City.Id, a.Created).Scan(&ai)

	if err != nil {
		log.Println("Failed execute func_streets_add: ", err)
		return 0, err
	}

	return ai, nil
}

//func (est *StreetStorage) Upd(ctx context.Context, eu models.Street) (int, error)
func (est *StreetStorage) Upd(ctx context.Context, u models.Street) (int, error) {
	dbpool := pgclient.WDB
	ui := 0

	err := dbpool.QueryRow(ctx, "SELECT func_streets_upd($1,$2,$3);", u.Id, u.StreetName, u.Created).Scan(&ui)

	if err != nil {
		log.Println("Failed execute func_streets_upd: ", err)
		return 0, err
	}
	return ui, nil
}

//func (est *StreetStorage) Del(ctx context.Context, d models.StreetClose) (models.Json_id, error)
func (est *StreetStorage) Del(ctx context.Context, d models.StreetClose) (models.Json_id, error) {
	dbpool := pgclient.WDB
	i := 0
	err := dbpool.QueryRow(ctx, "SELECT func_streets_del($1,$2);", d.Id, d.Close).Scan(&i)
	if err != nil {
		log.Println("Failed execute func_streets_del: ", err)
	}

	res := models.Json_id{Id: i}

	return res, nil
}

//func (est *StreetStorage) GetOne(ctx context.Context, i int) (models.Street_count, error)
func (est *StreetStorage) GetOne(ctx context.Context, i int) (models.Street_count, error) {
	dbpool := pgclient.WDB
	out_arr := []models.Street{}
	g := models.Street{}

	err := dbpool.QueryRow(ctx, "SELECT * from func_street_get($1);", i).Scan(&g.Id, &g.StreetName, &g.Created, &g.Closed,
		&g.City.CityName, &g.City.Id)

	if err != nil && err != pgx.ErrNoRows {
		log.Println("Failed execute func_street_get: ", err)
		return models.Street_count{Values: []models.Street{}, Count: 0, Auth: models.Auth{}}, err
	}

	out_arr = append(out_arr, g)

	out_count := models.Street_count{Values: out_arr, Count: 1, Auth: models.Auth{}}
	return out_count, nil
}
