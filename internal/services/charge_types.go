package services

import (
	"context"
	"log"

	"github.com/u2lentaru/billbck/internal/adapters/db/pgsql"
	"github.com/u2lentaru/billbck/internal/models"
)

type ChargeTypeService struct {
	storage pgsql.ChargeTypeStorage
}

type ifChargeTypeStorage interface {
	GetList(ctx context.Context, pg, pgs int, gs1 string, ord int, dsc bool) (models.ChargeType_count, error)
	Add(ctx context.Context, ea models.ChargeType) (int, error)
	Upd(ctx context.Context, eu models.ChargeType) (int, error)
	Del(ctx context.Context, ed []int) ([]int, error)
	GetOne(ctx context.Context, i int) (models.ChargeType_count, error)
}

//func NewChargeTypeService(storage pgsql.ChargeTypeStorage) *ChargeTypeService
func NewChargeTypeService(storage pgsql.ChargeTypeStorage) *ChargeTypeService {
	return &ChargeTypeService{storage}
}

//func (esv *ChargeTypeService) GetList(ctx context.Context, pg, pgs int, gs1 string, ord int, dsc bool) (models.ChargeType_count, error)
func (esv *ChargeTypeService) GetList(ctx context.Context, pg, pgs int, gs1 string, ord int, dsc bool) (models.ChargeType_count, error) {
	var est ifChargeTypeStorage
	est = &esv.storage

	out_count, err := est.GetList(ctx, pg, pgs, gs1, ord, dsc)

	if err != nil {
		log.Println("ChargeTypeStorage.GetList", err)
		return models.ChargeType_count{Values: []models.ChargeType{}, Count: 0, Auth: models.Auth{}}, err
	}

	return out_count, nil
}

//func (esv *ChargeTypeService) Add(ctx context.Context, ea models.ChargeType) (int, error)
func (esv *ChargeTypeService) Add(ctx context.Context, ea models.ChargeType) (int, error) {
	var est ifChargeTypeStorage
	est = &esv.storage

	ai, err := est.Add(ctx, ea)

	if err != nil {
		log.Println("ChargeTypeStorage.Add", err)
		return 0, err
	}

	return ai, nil
}

//func (esv *ChargeTypeService) Upd(ctx context.Context, eu models.ChargeType) (int, error)
func (esv *ChargeTypeService) Upd(ctx context.Context, eu models.ChargeType) (int, error) {
	var est ifChargeTypeStorage
	est = &esv.storage

	ui, err := est.Upd(ctx, eu)

	if err != nil {
		log.Println("ChargeTypeStorage.Upd", err)
		return 0, err
	}

	return ui, nil
}

//func (esv *ChargeTypeService) Del(ctx context.Context, ed []int) ([]int, error) {
func (esv *ChargeTypeService) Del(ctx context.Context, ed []int) ([]int, error) {
	var est ifChargeTypeStorage
	est = &esv.storage

	res, err := est.Del(ctx, ed)

	if err != nil {
		log.Println("ChargeTypeStorage.Del", err)
		return []int{}, err
	}

	return res, nil
}

//func (esv *ChargeTypeService) GetOne(ctx context.Context, i int) (models.ChargeType_count, error)
func (esv *ChargeTypeService) GetOne(ctx context.Context, i int) (models.ChargeType_count, error) {
	var est ifChargeTypeStorage
	est = &esv.storage

	out_count, err := est.GetOne(ctx, i)

	if err != nil {
		log.Println("ChargeTypeStorage.GetOne", err)
		return models.ChargeType_count{Values: []models.ChargeType{}, Count: 0, Auth: models.Auth{}}, err
	}

	return out_count, nil
}
