package services

import (
	"context"
	"log"

	"github.com/u2lentaru/billbck/internal/adapters/db/pgsql"
	"github.com/u2lentaru/billbck/internal/models"
)

type SectorService struct {
	storage pgsql.SectorStorage
}

type ifSectorStorage interface {
	GetList(ctx context.Context, pg, pgs int, gs1 string, ord int, dsc bool) (models.Sector_count, error)
	Add(ctx context.Context, ea models.Sector) (int, error)
	Upd(ctx context.Context, eu models.Sector) (int, error)
	Del(ctx context.Context, ed []int) ([]int, error)
	GetOne(ctx context.Context, i int) (models.Sector_count, error)
}

//func NewSectorService(storage pgsql.Sectortorage) *SectorService
func NewSectorService(storage pgsql.SectorStorage) *SectorService {
	return &SectorService{storage}
}

//func (esv *SectorService) GetList(ctx context.Context, pg, pgs int, gs1 string, ord int, dsc bool) (models.Sector_count, error)
func (esv *SectorService) GetList(ctx context.Context, pg, pgs int, gs1 string, ord int, dsc bool) (models.Sector_count, error) {
	var est ifSectorStorage
	est = &esv.storage

	out_count, err := est.GetList(ctx, pg, pgs, gs1, ord, dsc)

	if err != nil {
		log.Println("SectorStorage.GetList", err)
		return models.Sector_count{Values: []models.Sector{}, Count: 0, Auth: models.Auth{}}, err
	}

	return out_count, nil
}

//func (esv *SectorService) Add(ctx context.Context, ea models.Sector) (int, error)
func (esv *SectorService) Add(ctx context.Context, ea models.Sector) (int, error) {
	var est ifSectorStorage
	est = &esv.storage

	ai, err := est.Add(ctx, ea)

	if err != nil {
		log.Println("SectorStorage.Add", err)
		return 0, err
	}

	return ai, nil
}

//func (esv *SectorService) Upd(ctx context.Context, eu models.Sector) (int, error)
func (esv *SectorService) Upd(ctx context.Context, eu models.Sector) (int, error) {
	var est ifSectorStorage
	est = &esv.storage

	ui, err := est.Upd(ctx, eu)

	if err != nil {
		log.Println("SectorStorage.Upd", err)
		return 0, err
	}

	return ui, nil
}

//func (esv *SectorService) Del(ctx context.Context, ed []int) ([]int, error) {
func (esv *SectorService) Del(ctx context.Context, ed []int) ([]int, error) {
	var est ifSectorStorage
	est = &esv.storage

	res, err := est.Del(ctx, ed)

	if err != nil {
		log.Println("SectorStorage.Del", err)
		return []int{}, err
	}

	return res, nil
}

//func (esv *SectorService) GetOne(ctx context.Context, i int) (models.Sector_count, error)
func (esv *SectorService) GetOne(ctx context.Context, i int) (models.Sector_count, error) {
	var est ifSectorStorage
	est = &esv.storage

	out_count, err := est.GetOne(ctx, i)

	if err != nil {
		log.Println("SectorStorage.GetOne", err)
		return models.Sector_count{Values: []models.Sector{}, Count: 0, Auth: models.Auth{}}, err
	}

	return out_count, nil
}
