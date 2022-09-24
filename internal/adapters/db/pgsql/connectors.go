package pgsql

import (
	"context"
	"log"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/u2lentaru/billbck/internal/models"
	"github.com/u2lentaru/billbck/pkg/pgclient"
)

//type ConnectorStorage struct
type ConnectorStorage struct {
	db *pgxpool.Pool
}

//func NewConnectorStorage(db *pgxpool.Pool) *ConnectorStorage
func NewConnectorStorage(db *pgxpool.Pool) *ConnectorStorage {
	return &ConnectorStorage{db: db}
}

//func (est *ConnectorStorage) GetList(ctx context.Context, pg, pgs int, gs1 string, ord int, dsc bool) (models.Connector_count, error)
func (est *ConnectorStorage) GetList(ctx context.Context, pg, pgs int, gs1 string, ord int, dsc bool) (models.Connector_count, error) {
	dbpool := pgclient.WDB
	gs := models.Connector{}

	gsc := 0
	err := dbpool.QueryRow(ctx, "SELECT * from func_connectors_cnt($1);", gs1).Scan(&gsc)

	if err != nil {
		log.Println(err.Error(), "func_connectors_cnt")
		return models.Connector_count{Values: []models.Connector{}, Count: gsc, Auth: models.Auth{}}, err
	}

	out_arr := make([]models.Connector, 0,
		func() int {
			if gsc < pgs {
				return gsc
			} else {
				return pgs
			}
		}())

	rows, err := dbpool.Query(ctx, "SELECT * from func_connectors_get($1,$2,$3,$4,$5);", pg, pgs, gs1, ord, dsc)
	if err != nil {
		log.Println(err.Error())
		return models.Connector_count{Values: []models.Connector{}, Count: gsc, Auth: models.Auth{}}, err
	}

	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&gs.Id, &gs.ConnectorName)
		if err != nil {
			log.Println("failed to scan row:", err)
		}

		out_arr = append(out_arr, gs)
	}

	out_count := models.Connector_count{Values: out_arr, Count: gsc, Auth: models.Auth{}}
	if err != nil {
		log.Println(err.Error())
		return models.Connector_count{}, err
	}

	return out_count, nil
}

//func (est *ConnectorStorage) Add(ctx context.Context, a models.Connector) (int, error)
func (est *ConnectorStorage) Add(ctx context.Context, a models.Connector) (int, error) {
	dbpool := pgclient.WDB
	ai := 0

	err := dbpool.QueryRow(context.Background(), "SELECT func_connectors_add($1);", a.ConnectorName).Scan(&ai)

	if err != nil {
		log.Println("Failed execute func_connectors_add: ", err)
		return 0, err
	}

	return ai, nil
}

//func (est *ConnectorStorage) Upd(ctx context.Context, u models.Connector) (int, error)
func (est *ConnectorStorage) Upd(ctx context.Context, u models.Connector) (int, error) {
	dbpool := pgclient.WDB
	ui := 0

	err := dbpool.QueryRow(context.Background(), "SELECT func_connectors_upd($1,$2);", u.Id, u.ConnectorName).Scan(&ui)

	if err != nil {
		log.Println("Failed execute func_connectors_upd: ", err)
		return 0, err
	}
	return ui, nil
}

//func (est *ConnectorStorage) Del(ctx context.Context, d []int) ([]int, error)
func (est *ConnectorStorage) Del(ctx context.Context, d []int) ([]int, error) {
	dbpool := pgclient.WDB
	res := []int{}
	i := 0
	for _, id := range d {
		err := dbpool.QueryRow(context.Background(), "SELECT func_connectors_del($1);", id).Scan(&i)
		res = append(res, i)

		if err != nil {
			log.Println("Failed execute func_connectors_del: ", err)
		}
	}
	return res, nil
}

//func (est *ConnectorStorage) GetOne(ctx context.Context, i int) (models.Connector_count, error)
func (est *ConnectorStorage) GetOne(ctx context.Context, i int) (models.Connector_count, error) {
	dbpool := pgclient.WDB
	out_arr := []models.Connector{}
	g := models.Connector{}

	err := dbpool.QueryRow(context.Background(), "SELECT * from func_connector_get($1);", i).Scan(&g.Id, &g.ConnectorName)

	if err != nil && err != pgx.ErrNoRows {
		log.Println("Failed execute from func_connector_get: ", err)
		return models.Connector_count{Values: []models.Connector{}, Count: 0, Auth: models.Auth{}}, err
	}

	out_arr = append(out_arr, g)

	out_count := models.Connector_count{Values: out_arr, Count: 1, Auth: models.Auth{}}
	return out_count, nil
}
