package services

import (
	"context"
	"log"

	"github.com/u2lentaru/billbck/internal/adapters/db/pgsql"
	"github.com/u2lentaru/billbck/internal/models"
)

type DistributionZoneService struct {
	storage pgsql.DistributionZoneStorage
}

type ifDistributionZoneStorage interface {
	GetList(ctx context.Context, pg, pgs int, gs1 string, ord int, dsc bool) (models.DistributionZone_count, error)
	Add(ctx context.Context, ea models.DistributionZone) (int, error)
	Upd(ctx context.Context, eu models.DistributionZone) (int, error)
	Del(ctx context.Context, ed []int) ([]int, error)
	GetOne(ctx context.Context, i int) (models.DistributionZone_count, error)
}

//func NewDistributionZoneService(storage pgsql.DistributionZoneStorage) *DistributionZoneService
func NewDistributionZoneService(storage pgsql.DistributionZoneStorage) *DistributionZoneService {
	return &DistributionZoneService{storage}
}

//func (esv *DistributionZoneService) GetList(ctx context.Context, pg, pgs int, gs1 string, ord int, dsc bool) (models.DistributionZone_count, error)
func (esv *DistributionZoneService) GetList(ctx context.Context, pg, pgs int, gs1 string, ord int, dsc bool) (models.DistributionZone_count, error) {
	var est ifDistributionZoneStorage
	est = &esv.storage

	out_count, err := est.GetList(ctx, pg, pgs, gs1, ord, dsc)

	if err != nil {
		log.Println("DistributionZoneStorage.GetList", err)
		return models.DistributionZone_count{Values: []models.DistributionZone{}, Count: 0, Auth: models.Auth{}}, err
	}

	return out_count, nil
}

//func (esv *DistributionZoneService) Add(ctx context.Context, ea models.DistributionZone) (int, error)
func (esv *DistributionZoneService) Add(ctx context.Context, ea models.DistributionZone) (int, error) {
	var est ifDistributionZoneStorage
	est = &esv.storage

	ai, err := est.Add(ctx, ea)

	if err != nil {
		log.Println("DistributionZoneStorage.Add", err)
		return 0, err
	}

	return ai, nil
}

//func (esv *DistributionZoneService) Upd(ctx context.Context, eu models.DistributionZone) (int, error)
func (esv *DistributionZoneService) Upd(ctx context.Context, eu models.DistributionZone) (int, error) {
	var est ifDistributionZoneStorage
	est = &esv.storage

	ui, err := est.Upd(ctx, eu)

	if err != nil {
		log.Println("DistributionZoneStorage.Upd", err)
		return 0, err
	}

	return ui, nil
}

//func (esv *DistributionZoneService) Del(ctx context.Context, ed []int) ([]int, error) {
func (esv *DistributionZoneService) Del(ctx context.Context, ed []int) ([]int, error) {
	var est ifDistributionZoneStorage
	est = &esv.storage

	res, err := est.Del(ctx, ed)

	if err != nil {
		log.Println("DistributionZoneStorage.Del", err)
		return []int{}, err
	}

	return res, nil
}

//func (esv *DistributionZoneService) GetOne(ctx context.Context, i int) (models.DistributionZone_count, error)
func (esv *DistributionZoneService) GetOne(ctx context.Context, i int) (models.DistributionZone_count, error) {
	var est ifDistributionZoneStorage
	est = &esv.storage

	out_count, err := est.GetOne(ctx, i)

	if err != nil {
		log.Println("DistributionZoneStorage.GetOne", err)
		return models.DistributionZone_count{Values: []models.DistributionZone{}, Count: 0, Auth: models.Auth{}}, err
	}

	return out_count, nil
}
