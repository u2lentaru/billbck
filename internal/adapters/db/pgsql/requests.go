package pgsql

import (
	"context"
	"database/sql"
	"log"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/u2lentaru/billbck/internal/models"
	"github.com/u2lentaru/billbck/internal/utils"
	"github.com/u2lentaru/billbck/pkg/pgclient"
)

//type RequestStorage struct
type RequestStorage struct {
	db *pgxpool.Pool
}

//func NewRequestStorage(db *pgxpool.Pool) *RequestStorage
func NewRequestStorage(db *pgxpool.Pool) *RequestStorage {
	return &RequestStorage{db: db}
}

//func (est *RequestStorage) GetList(ctx context.Context, pg, pgs int, gs1 string, ord int, dsc bool) (models.Request_count, error)
func (est *RequestStorage) GetList(ctx context.Context, pg, pgs int, gs1 string, ord int, dsc bool) (models.Request_count, error) {
	dbpool := pgclient.WDB
	gs := models.Request{}

	gsc := 0
	err := dbpool.QueryRow(ctx, "SELECT * from func_requests_cnt($1);", gs1).Scan(&gsc)

	if err != nil {
		log.Println(err.Error(), "func_requests_cnt")
		return models.Request_count{Values: []models.Request{}, Count: gsc, Auth: models.Auth{}}, err
	}

	out_arr := make([]models.Request, 0,
		func() int {
			if gsc < pgs {
				return gsc
			} else {
				return pgs
			}
		}())

	rows, err := dbpool.Query(ctx, "SELECT * from func_requests_get($1,$2,$3,$4,$5);", pg, pgs, gs1, ord, dsc)
	if err != nil {
		log.Println(err.Error(), "func_requests_get")
		return models.Request_count{Values: []models.Request{}, Count: gsc, Auth: models.Auth{}}, err
	}

	defer rows.Close()

	var ci, cli, ai, oi sql.NullInt32
	var cn, cln, an, on sql.NullString

	for rows.Next() {
		err = rows.Scan(&gs.Id, &gs.RequestNumber, &gs.RequestDate, &ci, &gs.ServiceType.Id,
			&gs.RequestType.Id, &gs.RequestKind.Id, &cli, &gs.TermDate, &gs.Executive, &gs.Accept, &gs.Notes,
			&gs.Result.Id, &ai, &oi, &cn, &gs.ServiceType.ServiceTypeName, &gs.RequestType.RequestTypeName,
			&gs.RequestKind.RequestKindName, &cln, &gs.Result.ResultName, &an, &on)

		gs.Contract.Id = int(ci.Int32)
		gs.Contract.ContractNumber = cn.String
		gs.ClaimType.Id = int(cli.Int32)
		gs.ClaimType.ClaimTypeName = cln.String
		gs.Act.Id = int(ai.Int32)
		gs.Act.ActNumber = an.String
		gs.Object.Id = int(oi.Int32)
		gs.Object.ObjectName = on.String

		if err != nil {
			log.Println("failed to scan row:", err)
		}

		out_arr = append(out_arr, gs)
	}

	out_count := models.Request_count{Values: out_arr, Count: gsc, Auth: models.Auth{}}
	if err != nil {
		log.Println(err.Error())
		return models.Request_count{}, err
	}

	return out_count, nil
}

//func (est *RequestStorage) Add(ctx context.Context, a models.Request) (int, error)
func (est *RequestStorage) Add(ctx context.Context, a models.Request) (int, error) {
	dbpool := pgclient.WDB
	ai := 0

	err := dbpool.QueryRow(ctx, "SELECT func_requests_add($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14);",
		a.RequestNumber, a.RequestDate, utils.NullableInt(int32(a.Contract.Id)), a.ServiceType.Id, a.RequestType.Id, a.RequestKind.Id,
		utils.NullableInt(int32(a.ClaimType.Id)), a.TermDate, a.Executive, a.Accept, a.Notes, a.Result.Id,
		utils.NullableInt(int32(a.Act.Id)), utils.NullableInt(int32(a.Object.Id))).Scan(&ai)

	if err != nil {
		log.Println("Failed execute func_requests_add: ", err)
		return 0, err
	}

	return ai, nil
}

//func (est *RequestStorage) Upd(ctx context.Context, u models.Request) (int, error)
func (est *RequestStorage) Upd(ctx context.Context, u models.Request) (int, error) {
	dbpool := pgclient.WDB
	ui := 0

	err := dbpool.QueryRow(ctx, "SELECT func_requests_upd($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15);",
		u.Id, u.RequestNumber, u.RequestDate, utils.NullableInt(int32(u.Contract.Id)), u.ServiceType.Id,
		u.RequestType.Id, u.RequestKind.Id, utils.NullableInt(int32(u.ClaimType.Id)), u.TermDate, u.Executive, u.Accept,
		u.Notes, u.Result.Id, utils.NullableInt(int32(u.Act.Id)), utils.NullableInt(int32(u.Object.Id))).Scan(&ui)

	if err != nil {
		log.Println("Failed execute func_requests_upd: ", err)
		return 0, err
	}
	return ui, nil
}

//func (est *RequestStorage) Del(ctx context.Context, d []int) ([]int, error)
func (est *RequestStorage) Del(ctx context.Context, d []int) ([]int, error) {
	dbpool := pgclient.WDB
	res := []int{}
	i := 0
	for _, id := range d {
		err := dbpool.QueryRow(ctx, "SELECT func_requests_del($1);", id).Scan(&i)
		res = append(res, i)

		if err != nil {
			log.Println("Failed execute func_requests_del: ", err)
		}
	}
	return res, nil
}

//func (est *RequestStorage) GetOne(ctx context.Context, i int) (models.Request_count, error)
func (est *RequestStorage) GetOne(ctx context.Context, i int) (models.Request_count, error) {
	dbpool := pgclient.WDB
	out_arr := []models.Request{}
	g := models.Request{}

	var ci, cli, ai, oi sql.NullInt32
	var cn, cln, an, on sql.NullString

	err := dbpool.QueryRow(ctx, "SELECT * from func_request_get($1);", i).Scan(&g.Id,
		&g.RequestNumber, &g.RequestDate, &ci, &g.ServiceType.Id, &g.RequestType.Id, &g.RequestKind.Id, &cli, &g.TermDate,
		&g.Executive, &g.Accept, &g.Notes, &g.Result.Id, &ai, &oi, &cn, &g.ServiceType.ServiceTypeName,
		&g.RequestType.RequestTypeName, &g.RequestKind.RequestKindName, &cln, &g.Result.ResultName, &an, &on)

	g.Contract.Id = int(ci.Int32)
	g.Contract.ContractNumber = cn.String
	g.ClaimType.Id = int(cli.Int32)
	g.ClaimType.ClaimTypeName = cln.String
	g.Act.Id = int(ai.Int32)
	g.Act.ActNumber = an.String
	g.Object.Id = int(oi.Int32)
	g.Object.ObjectName = on.String

	if err != nil && err != pgx.ErrNoRows {
		log.Println("Failed execute func_request_get: ", err)
		return models.Request_count{Values: []models.Request{}, Count: 0, Auth: models.Auth{}}, err
	}

	out_arr = append(out_arr, g)

	out_count := models.Request_count{Values: out_arr, Count: 1, Auth: models.Auth{}}
	return out_count, nil
}
