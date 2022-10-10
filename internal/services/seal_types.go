package services

import (
	"context"
	"log"

	"github.com/u2lentaru/billbck/internal/adapters/db/pgsql"
	"github.com/u2lentaru/billbck/internal/models"
)

type SealTypeService struct {
	storage pgsql.SealTypeStorage
}

type ifSealTypeStorage interface {
	GetList(ctx context.Context, pg, pgs int, gs1 string, ord int, dsc bool) (models.SealType_count, error)
	Add(ctx context.Context, ea models.SealType) (int, error)
	Upd(ctx context.Context, eu models.SealType) (int, error)
	Del(ctx context.Context, ed []int) ([]int, error)
	GetOne(ctx context.Context, i int) (models.SealType_count, error)
}

//func NewSealTypeService(storage pgsql.SealTypetorage) *SealTypeService
func NewSealTypeService(storage pgsql.SealTypeStorage) *SealTypeService {
	return &SealTypeService{storage}
}

//func (esv *SealTypeService) GetList(ctx context.Context, pg, pgs int, gs1 string, ord int, dsc bool) (models.SealType_count, error)
func (esv *SealTypeService) GetList(ctx context.Context, pg, pgs int, gs1 string, ord int, dsc bool) (models.SealType_count, error) {
	var est ifSealTypeStorage
	est = &esv.storage

	out_count, err := est.GetList(ctx, pg, pgs, gs1, ord, dsc)

	if err != nil {
		log.Println("SealTypeStorage.GetList", err)
		return models.SealType_count{Values: []models.SealType{}, Count: 0, Auth: models.Auth{}}, err
	}

	return out_count, nil
}

//func (esv *SealTypeService) Add(ctx context.Context, ea models.SealType) (int, error)
func (esv *SealTypeService) Add(ctx context.Context, ea models.SealType) (int, error) {
	var est ifSealTypeStorage
	est = &esv.storage

	ai, err := est.Add(ctx, ea)

	if err != nil {
		log.Println("SealTypeStorage.Add", err)
		return 0, err
	}

	return ai, nil
}

//func (esv *SealTypeService) Upd(ctx context.Context, eu models.SealType) (int, error)
func (esv *SealTypeService) Upd(ctx context.Context, eu models.SealType) (int, error) {
	var est ifSealTypeStorage
	est = &esv.storage

	ui, err := est.Upd(ctx, eu)

	if err != nil {
		log.Println("SealTypeStorage.Upd", err)
		return 0, err
	}

	return ui, nil
}

//func (esv *SealTypeService) Del(ctx context.Context, ed []int) ([]int, error) {
func (esv *SealTypeService) Del(ctx context.Context, ed []int) ([]int, error) {
	var est ifSealTypeStorage
	est = &esv.storage

	res, err := est.Del(ctx, ed)

	if err != nil {
		log.Println("SealTypeStorage.Del", err)
		return []int{}, err
	}

	return res, nil
}

//func (esv *SealTypeService) GetOne(ctx context.Context, i int) (models.SealType_count, error)
func (esv *SealTypeService) GetOne(ctx context.Context, i int) (models.SealType_count, error) {
	var est ifSealTypeStorage
	est = &esv.storage

	out_count, err := est.GetOne(ctx, i)

	if err != nil {
		log.Println("SealTypeStorage.GetOne", err)
		return models.SealType_count{Values: []models.SealType{}, Count: 0, Auth: models.Auth{}}, err
	}

	return out_count, nil
}
