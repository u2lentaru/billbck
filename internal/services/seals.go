package services

import (
	"context"
	"log"

	"github.com/u2lentaru/billbck/internal/adapters/db/pgsql"
	"github.com/u2lentaru/billbck/internal/models"
)

type SealService struct {
	storage pgsql.SealStorage
}

type ifSealStorage interface {
	GetList(ctx context.Context, pg, pgs int, gs1 string, ord int, dsc bool) (models.Seal_count, error)
	Add(ctx context.Context, a models.Seal) (int, error)
	Upd(ctx context.Context, u models.Seal) (int, error)
	Del(ctx context.Context, d []int) ([]int, error)
	GetOne(ctx context.Context, i int) (models.Seal_count, error)
}

//func NewSealService(storage pgsql.Sealtorage) *SealService
func NewSealService(storage pgsql.SealStorage) *SealService {
	return &SealService{storage}
}

//func (esv *SealService) GetList(ctx context.Context, pg, pgs int, gs1 string, ord int, dsc bool) (models.Seal_count, error)
func (esv *SealService) GetList(ctx context.Context, pg, pgs int, gs1 string, ord int, dsc bool) (models.Seal_count, error) {
	var est ifSealStorage
	est = &esv.storage

	out_count, err := est.GetList(ctx, pg, pgs, gs1, ord, dsc)

	if err != nil {
		log.Println("SealStorage.GetList", err)
		return models.Seal_count{Values: []models.Seal{}, Count: 0, Auth: models.Auth{}}, err
	}

	return out_count, nil
}

//func (esv *SealService) Add(ctx context.Context, ea models.Seal) (int, error)
func (esv *SealService) Add(ctx context.Context, ea models.Seal) (int, error) {
	var est ifSealStorage
	est = &esv.storage

	ai, err := est.Add(ctx, ea)

	if err != nil {
		log.Println("SealStorage.Add", err)
		return 0, err
	}

	return ai, nil
}

//func (esv *SealService) Upd(ctx context.Context, eu models.Seal) (int, error)
func (esv *SealService) Upd(ctx context.Context, eu models.Seal) (int, error) {
	var est ifSealStorage
	est = &esv.storage

	ui, err := est.Upd(ctx, eu)

	if err != nil {
		log.Println("SealStorage.Upd", err)
		return 0, err
	}

	return ui, nil
}

//func (esv *SealService) Del(ctx context.Context, ed []int) ([]int, error) {
func (esv *SealService) Del(ctx context.Context, ed []int) ([]int, error) {
	var est ifSealStorage
	est = &esv.storage

	res, err := est.Del(ctx, ed)

	if err != nil {
		log.Println("SealStorage.Del", err)
		return []int{}, err
	}

	return res, nil
}

//func (esv *SealService) GetOne(ctx context.Context, i int) (models.Seal_count, error)
func (esv *SealService) GetOne(ctx context.Context, i int) (models.Seal_count, error) {
	var est ifSealStorage
	est = &esv.storage

	out_count, err := est.GetOne(ctx, i)

	if err != nil {
		log.Println("SealStorage.GetOne", err)
		return models.Seal_count{Values: []models.Seal{}, Count: 0, Auth: models.Auth{}}, err
	}

	return out_count, nil
}
