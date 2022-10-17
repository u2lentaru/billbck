package pgsql

import (
	"context"
	"log"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/u2lentaru/billbck/internal/models"
	"github.com/u2lentaru/billbck/pkg/pgclient"
)

//type UserStorage struct
type UserStorage struct {
	db *pgxpool.Pool
}

//func NewUserStorage(db *pgxpool.Pool) *UserStorage
func NewUserStorage(db *pgxpool.Pool) *UserStorage {
	return &UserStorage{db: db}
}

//func (est *UserStorage) GetList(ctx context.Context, pg, pgs int, gs1 string, ord int, dsc bool) (models.User_count, error)
func (est *UserStorage) GetList(ctx context.Context, pg, pgs int, gs1 string, ord int, dsc bool) (models.User_count, error) {
	dbpool := pgclient.WDB
	gs := models.User{}

	gsc := 0
	err := dbpool.QueryRow(ctx, "SELECT * from func_users_cnt($1);", gs1).Scan(&gsc)

	if err != nil {
		log.Println(err.Error(), "func_users_cnt")
		return models.User_count{Values: []models.User{}, Count: gsc, Auth: models.Auth{}}, err
	}

	out_arr := make([]models.User, 0,
		func() int {
			if gsc < pgs {
				return gsc
			} else {
				return pgs
			}
		}())

	rows, err := dbpool.Query(ctx, "SELECT * from func_users_get($1,$2,$3,$4,$5);", pg, pgs, gs1, ord, dsc)
	if err != nil {
		log.Println(err.Error(), "func_users_get")
		return models.User_count{Values: []models.User{}, Count: gsc, Auth: models.Auth{}}, err
	}

	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&gs.Id, &gs.UserName, &gs.OrgInfo.Id, &gs.Lang.Id, &gs.ChangePass, &gs.Position.Id, &gs.UserFullName,
			&gs.Created, &gs.Closed, &gs.OrgInfo.OIName, &gs.Lang.LangName, &gs.Position.PositionName)
		if err != nil {
			log.Println("failed to scan row:", err)
		}

		out_arr = append(out_arr, gs)
	}

	out_count := models.User_count{Values: out_arr, Count: gsc, Auth: models.Auth{}}

	return out_count, nil
}

//func (est *UserStorage) Add(ctx context.Context, a models.User) (int, error)
func (est *UserStorage) Add(ctx context.Context, a models.User) (int, error) {
	dbpool := pgclient.WDB
	ai := 0

	err := dbpool.QueryRow(ctx, "SELECT func_users_add($1,$2,$3,$4,$5,$6,$7);", a.UserName, a.OrgInfo.Id, a.Lang.Id, a.ChangePass,
		a.Position.Id, a.UserFullName, a.Created).Scan(&ai)

	if err != nil {
		log.Println("Failed execute func_users_add: ", err)
		return 0, err
	}

	return ai, nil
}

//func (est *UserStorage) Upd(ctx context.Context, u models.User) (int, error)
func (est *UserStorage) Upd(ctx context.Context, u models.User) (int, error) {
	dbpool := pgclient.WDB
	ui := 0

	err := dbpool.QueryRow(ctx, "SELECT func_users_upd($1,$2,$3,$4,$5,$6,$7,$8,$9);", u.Id, u.UserName, u.OrgInfo.Id, u.Lang.Id,
		u.ChangePass, u.Position.Id, u.UserFullName, u.Created, u.Closed).Scan(&ui)

	if err != nil {
		log.Println("Failed execute func_users_upd: ", err)
		return 0, err
	}
	return ui, nil
}

//func (est *UserStorage) Del(ctx context.Context, d []int) ([]int, error)
func (est *UserStorage) Del(ctx context.Context, d []int) ([]int, error) {
	dbpool := pgclient.WDB
	res := []int{}
	i := 0
	for _, id := range d {
		err := dbpool.QueryRow(ctx, "SELECT func_users_del($1);", id).Scan(&i)
		res = append(res, i)

		if err != nil {
			log.Println("Failed execute func_users_del: ", err)
		}
	}
	return res, nil
}

//func (est *UserStorage) GetOne(ctx context.Context, i int) (models.User_count, error)
func (est *UserStorage) GetOne(ctx context.Context, i int) (models.User_count, error) {
	dbpool := pgclient.WDB
	out_arr := []models.User{}
	g := models.User{}

	err := dbpool.QueryRow(ctx, "SELECT * from func_user_get($1);", i).Scan(&g.Id, &g.UserName, &g.OrgInfo.Id, &g.Lang.Id, &g.ChangePass,
		&g.Position.Id, &g.UserFullName, &g.Created, &g.Closed, &g.OrgInfo.OIName, &g.Lang.LangName, &g.Position.PositionName)

	if err != nil && err != pgx.ErrNoRows {
		log.Println("Failed execute func_user_get: ", err)
		return models.User_count{Values: []models.User{}, Count: 0, Auth: models.Auth{}}, err
	}

	out_arr = append(out_arr, g)

	out_count := models.User_count{Values: out_arr, Count: 1, Auth: models.Auth{}}
	return out_count, nil
}
