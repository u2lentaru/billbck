package models

import (
	"context"
	"log"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

//DistributionZone struct
type DistributionZone struct {
	Id                   int    `json:"id"`
	DistributionZoneName string `json:"distributionzonename"`
}

//AddDistributionZone struct
type AddDistributionZone struct {
	DistributionZoneName string `json:"distributionzonename"`
}

//DistributionZone_count  struct
type DistributionZone_count struct {
	Values []DistributionZone `json:"values"`
	Count  int                `json:"count"`
	Auth   Auth               `json:"auth"`
}

//NewDistributionZone() *DistributionZone
func NewDistributionZone() *DistributionZone {
	return &DistributionZone{}
}

//(dz *DistributionZone) GetDistributionZones(ctx context.Context, Dbpool *pgxpool.Pool, pg, pgs int, nm string, ord int, dsc bool) (DistributionZone_count, error)
func (dz *DistributionZone) GetDistributionZones(ctx context.Context, Dbpool *pgxpool.Pool, pg, pgs int, nm string, ord int, dsc bool) (DistributionZone_count, error) {
	gsc := 0
	err := Dbpool.QueryRow(ctx, "SELECT * from func_distribution_zones_cnt($1);", nm).Scan(&gsc)
	auth := Auth{Create: true, Read: true, Update: true, Delete: true}

	if err != nil {
		log.Println(err.Error(), "func_distribution_zones_cnt")
		return DistributionZone_count{Values: []DistributionZone{}, Count: gsc, Auth: auth}, err
	}

	out_arr := make([]DistributionZone, 0,
		func() int {
			if gsc < pgs {
				return gsc
			} else {
				return pgs
			}
		}())

	rows, err := Dbpool.Query(ctx, "SELECT * from func_distribution_zones_get($1,$2,$3,$4,$5);", pg, pgs, nm, ord, dsc)
	if err != nil {
		log.Println(err.Error())
		return DistributionZone_count{Values: []DistributionZone{}, Count: gsc, Auth: auth}, err
	}

	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&(dz.Id), &(dz.DistributionZoneName))
		if err != nil {
			log.Println("failed to scan row:", err)
		}

		out_arr = append(out_arr, *dz)
	}

	out_count := DistributionZone_count{Values: out_arr, Count: gsc, Auth: auth}
	if err != nil {
		log.Println(err.Error())
		return DistributionZone_count{}, err
	}

	return out_count, nil
}

//func (dz *DistributionZone) AddDistributionZone(ctx context.Context, Dbpool *pgxpool.Pool) (int, error)
func (dz *DistributionZone) AddDistributionZone(ctx context.Context, Dbpool *pgxpool.Pool) (int, error) {
	ai := 0
	err := Dbpool.QueryRow(ctx, "SELECT func_distribution_zones_add($1);", dz.DistributionZoneName).Scan(&ai)

	if err != nil {
		log.Println("Failed execute func_distribution_zones_add: ", err)
		return 0, err
	}

	return ai, nil
}

//func (dz *DistributionZone) UpdDistributionZone(ctx context.Context, Dbpool *pgxpool.Pool)
func (dz *DistributionZone) UpdDistributionZone(ctx context.Context, Dbpool *pgxpool.Pool) (int, error) {
	ui := 0
	err := Dbpool.QueryRow(context.Background(), "SELECT func_distribution_zones_upd($1,$2);", dz.Id, dz.DistributionZoneName).Scan(&ui)

	if err != nil {
		log.Println("Failed execute func_distribution_zones_upd: ", err)
		return 0, err
	}
	return ui, nil
}

//func (dz *DistributionZone) DelDistributionZone(ctx context.Context, Dbpool *pgxpool.Pool, d []int) ([]int, error)
func (dz *DistributionZone) DelDistributionZone(ctx context.Context, Dbpool *pgxpool.Pool, d []int) ([]int, error) {
	res := []int{}
	i := 0
	for _, id := range d {
		err := Dbpool.QueryRow(ctx, "SELECT func_distribution_zones_del($1);", id).Scan(&i)
		res = append(res, i)

		if err != nil {
			log.Println("Failed execute func_distribution_zones_del: ", err)
			return []int{}, err
		}
	}
	return res, nil
}

//func (dz *DistributionZone) GetDistributionZone(ctx context.Context, Dbpool *pgxpool.Pool, i int) (DistributionZone_count, error)
func (dz *DistributionZone) GetDistributionZone(ctx context.Context, Dbpool *pgxpool.Pool, i int) (DistributionZone_count, error) {
	out_arr := []DistributionZone{}
	auth := Auth{Create: true, Read: true, Update: true, Delete: true}

	err := Dbpool.QueryRow(context.Background(), "SELECT * from func_distribution_zone_get($1);", i).Scan(&(dz.Id), &(dz.DistributionZoneName))

	if err != nil && err != pgx.ErrNoRows {
		log.Println("Failed execute from func_distribution_zone_get: ", err)
		return DistributionZone_count{Values: []DistributionZone{}, Count: 0, Auth: auth}, err
	}

	out_arr = append(out_arr, *dz)

	out_count := DistributionZone_count{Values: out_arr, Count: 0, Auth: auth}
	return out_count, nil
}
