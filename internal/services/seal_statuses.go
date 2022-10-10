package services

import (
	"context"
	"log"

	"github.com/u2lentaru/billbck/internal/adapters/db/pgsql"
	"github.com/u2lentaru/billbck/internal/models"
)

type SealStatusService struct {
	storage pgsql.SealStatusStorage
}

type ifSealStatusStorage interface {
	GetList(ctx context.Context, pg, pgs int, gs1 string, ord int, dsc bool) (models.SealStatus_count, error)
	Add(ctx context.Context, ea models.SealStatus) (int, error)
	Upd(ctx context.Context, eu models.SealStatus) (int, error)
	Del(ctx context.Context, ed []int) ([]int, error)
	GetOne(ctx context.Context, i int) (models.SealStatus_count, error)
}

//func NewSealStatusService(storage pgsql.SealStatustorage) *SealStatusService
func NewSealStatusService(storage pgsql.SealStatusStorage) *SealStatusService {
	return &SealStatusService{storage}
}

//func (esv *SealStatusService) GetList(ctx context.Context, pg, pgs int, gs1 string, ord int, dsc bool) (models.SealStatus_count, error)
func (esv *SealStatusService) GetList(ctx context.Context, pg, pgs int, gs1 string, ord int, dsc bool) (models.SealStatus_count, error) {
	var est ifSealStatusStorage
	est = &esv.storage

	out_count, err := est.GetList(ctx, pg, pgs, gs1, ord, dsc)

	if err != nil {
		log.Println("SealStatusStorage.GetList", err)
		return models.SealStatus_count{Values: []models.SealStatus{}, Count: 0, Auth: models.Auth{}}, err
	}

	return out_count, nil
}

//func (esv *SealStatusService) Add(ctx context.Context, ea models.SealStatus) (int, error)
func (esv *SealStatusService) Add(ctx context.Context, ea models.SealStatus) (int, error) {
	var est ifSealStatusStorage
	est = &esv.storage

	ai, err := est.Add(ctx, ea)

	if err != nil {
		log.Println("SealStatusStorage.Add", err)
		return 0, err
	}

	return ai, nil
}

//func (esv *SealStatusService) Upd(ctx context.Context, eu models.SealStatus) (int, error)
func (esv *SealStatusService) Upd(ctx context.Context, eu models.SealStatus) (int, error) {
	var est ifSealStatusStorage
	est = &esv.storage

	ui, err := est.Upd(ctx, eu)

	if err != nil {
		log.Println("SealStatusStorage.Upd", err)
		return 0, err
	}

	return ui, nil
}

//func (esv *SealStatusService) Del(ctx context.Context, ed []int) ([]int, error) {
func (esv *SealStatusService) Del(ctx context.Context, ed []int) ([]int, error) {
	var est ifSealStatusStorage
	est = &esv.storage

	res, err := est.Del(ctx, ed)

	if err != nil {
		log.Println("SealStatusStorage.Del", err)
		return []int{}, err
	}

	return res, nil
}

//func (esv *SealStatusService) GetOne(ctx context.Context, i int) (models.SealStatus_count, error)
func (esv *SealStatusService) GetOne(ctx context.Context, i int) (models.SealStatus_count, error) {
	var est ifSealStatusStorage
	est = &esv.storage

	out_count, err := est.GetOne(ctx, i)

	if err != nil {
		log.Println("SealStatusStorage.GetOne", err)
		return models.SealStatus_count{Values: []models.SealStatus{}, Count: 0, Auth: models.Auth{}}, err
	}

	return out_count, nil
}
