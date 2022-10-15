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

//type TguStorage struct
type TguStorage struct {
	db *pgxpool.Pool
}

//func NewTguStorage(db *pgxpool.Pool) *TguStorage
func NewTguStorage(db *pgxpool.Pool) *TguStorage {
	return &TguStorage{db: db}
}

//func (est *TguStorage) GetList(ctx context.Context, pg, pgs int, gs1, gs2 string, ord int, dsc bool) (models.Tgu_count, error) {
func (est *TguStorage) GetList(ctx context.Context, pg, pgs int, gs1, gs2 string, ord int, dsc bool) (models.Tgu_count, error) {
	dbpool := pgclient.WDB
	gs := models.Tgu{}

	gsc := 0
	err := dbpool.QueryRow(ctx, "SELECT * from func_tgu_cnt($1,$2);", gs1, utils.NullableString(gs2)).Scan(&gsc)

	if err != nil {
		log.Println(err.Error(), "func_tgu_cnt")
		return models.Tgu_count{Values: []models.Tgu{}, Count: gsc, Auth: models.Auth{}}, err
	}

	out_arr := make([]models.Tgu, 0,
		func() int {
			if gsc < pgs {
				return gsc
			} else {
				return pgs
			}
		}())

	rows, err := dbpool.Query(ctx, "SELECT * from func_tgu_get($1,$2,$3,$4,$5,$6);", pg, pgs, gs1, utils.NullableString(gs2), ord, dsc)
	if err != nil {
		log.Println(err.Error(), "func_tgu_get")
		return models.Tgu_count{Values: []models.Tgu{}, Count: gsc, Auth: models.Auth{}}, err
	}

	defer rows.Close()

	var iid, o1id, o2id, iv, o1v, o2v sql.NullInt32
	var ivn, ov1n, ov2n sql.NullString

	for rows.Next() {
		err = rows.Scan(&gs.Id, &gs.PId, &gs.TguName, &gs.TguType.Id, &gs.InvNumber, &iid, &o1id,
			&o2id, &gs.TguType.TguTypeName, &iv, &o1v, &o2v, &ivn, &ov1n, &ov2n)

		gs.InputVoltage.Id = int(iid.Int32)
		gs.InputVoltage.VoltageName = ivn.String
		gs.InputVoltage.VoltageValue = int(iv.Int32)
		gs.OutputVoltage1.Id = int(o1id.Int32)
		gs.OutputVoltage1.VoltageName = ov1n.String
		gs.OutputVoltage1.VoltageValue = int(o1v.Int32)
		gs.OutputVoltage2.Id = int(o2id.Int32)
		gs.OutputVoltage2.VoltageName = ov2n.String
		gs.OutputVoltage2.VoltageValue = int(o2v.Int32)

		if err != nil {
			log.Println("failed to scan row:", err)
		}

		out_arr = append(out_arr, gs)
	}

	out_count := models.Tgu_count{Values: out_arr, Count: gsc, Auth: models.Auth{}}
	if err != nil {
		log.Println(err.Error())
		return models.Tgu_count{}, err
	}

	return out_count, nil
}

//func (est *TguStorage) Add(ctx context.Context, ea models.Tgu) (int, error)
func (est *TguStorage) Add(ctx context.Context, a models.Tgu) (int, error) {
	dbpool := pgclient.WDB
	ai := 0

	err := dbpool.QueryRow(ctx, "SELECT func_tgu_add($1,$2,$3,$4,$5,$6,$7);", a.PId, a.TguName, a.TguType.Id,
		a.InvNumber, a.InputVoltage.Id, a.OutputVoltage1.Id, a.OutputVoltage2.Id).Scan(&ai)

	if err != nil {
		log.Println("Failed execute func_tgu_add: ", err)
		return 0, err
	}

	return ai, nil
}

//func (est *TguStorage) Upd(ctx context.Context, eu models.Tgu) (int, error)
func (est *TguStorage) Upd(ctx context.Context, u models.Tgu) (int, error) {
	dbpool := pgclient.WDB
	ui := 0

	err := dbpool.QueryRow(ctx, "SELECT func_tgu_upd($1,$2,$3,$4,$5,$6,$7,$8);", u.Id, u.PId, u.TguName, u.TguType.Id,
		u.InvNumber, u.InputVoltage.Id, u.OutputVoltage1.Id, u.OutputVoltage2.Id).Scan(&ui)

	if err != nil {
		log.Println("Failed execute func_tgu_upd: ", err)
		return 0, err
	}
	return ui, nil
}

//func (est *TguStorage) Del(ctx context.Context, ed []int) ([]int, error)
func (est *TguStorage) Del(ctx context.Context, d []int) ([]int, error) {
	dbpool := pgclient.WDB
	res := []int{}
	i := 0
	for _, id := range d {
		err := dbpool.QueryRow(ctx, "SELECT func_tgu_del($1);", id).Scan(&i)
		res = append(res, i)

		if err != nil {
			log.Println("Failed execute func_tgu_del: ", err)
		}
	}
	return res, nil
}

//func (est *TguStorage) GetOne(ctx context.Context, i int) (models.Tgu_count, error)
func (est *TguStorage) GetOne(ctx context.Context, i int) (models.Tgu_count, error) {
	dbpool := pgclient.WDB
	out_arr := []models.Tgu{}
	g := models.Tgu{}

	var iid, o1id, o2id, iv, o1v, o2v sql.NullInt32
	var ivn, ov1n, ov2n sql.NullString

	err := dbpool.QueryRow(ctx, "SELECT * from func_tgu_getbyid($1);", i).Scan(&g.Id, &g.PId, &g.TguName, &g.TguType.Id,
		&g.InvNumber, &iid, &o1id, &o2id, &g.TguType.TguTypeName, &iv, &o1v, &o2v, &ivn, &ov1n, &ov2n)

	g.InputVoltage.Id = int(iid.Int32)
	g.InputVoltage.VoltageName = ivn.String
	g.InputVoltage.VoltageValue = int(iv.Int32)
	g.OutputVoltage1.Id = int(o1id.Int32)
	g.OutputVoltage1.VoltageName = ov1n.String
	g.OutputVoltage1.VoltageValue = int(o1v.Int32)
	g.OutputVoltage2.Id = int(o2id.Int32)
	g.OutputVoltage2.VoltageName = ov2n.String
	g.OutputVoltage2.VoltageValue = int(o2v.Int32)

	if err != nil && err != pgx.ErrNoRows {
		log.Println("Failed execute func_tgu_getbyid: ", err)
		return models.Tgu_count{Values: []models.Tgu{}, Count: 0, Auth: models.Auth{}}, err
	}

	out_arr = append(out_arr, g)

	out_count := models.Tgu_count{Values: out_arr, Count: 1, Auth: models.Auth{}}
	return out_count, nil
}
