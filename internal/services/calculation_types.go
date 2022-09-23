package services

import (
	"context"
	"log"

	"github.com/u2lentaru/billbck/internal/adapters/db/pgsql"
	"github.com/u2lentaru/billbck/internal/models"
)

type CalculationTypeService struct {
	storage pgsql.CalculationTypeStorage
}

type ifCalculationTypeStorage interface {
	GetList(ctx context.Context, pg, pgs int, gs1 string, ord int, dsc bool) (models.CalculationType_count, error)
	Add(ctx context.Context, ea models.CalculationType) (int, error)
	Upd(ctx context.Context, eu models.CalculationType) (int, error)
	Del(ctx context.Context, ed []int) ([]int, error)
	GetOne(ctx context.Context, i int) (models.CalculationType_count, error)
}

//func NewCalculationTypeService(storage pgsql.CalculationTypeStorage) *CalculationTypeService
func NewCalculationTypeService(storage pgsql.CalculationTypeStorage) *CalculationTypeService {
	return &CalculationTypeService{storage}
}

//func (esv *CalculationTypeService) GetList(ctx context.Context, pg, pgs int, gs1 string, ord int, dsc bool) (models.CalculationType_count, error)
func (esv *CalculationTypeService) GetList(ctx context.Context, pg, pgs int, gs1 string, ord int, dsc bool) (models.CalculationType_count, error) {
	var est ifCalculationTypeStorage
	est = &esv.storage

	out_count, err := est.GetList(ctx, pg, pgs, gs1, ord, dsc)

	if err != nil {
		log.Println("CalculationTypeStorage.GetList", err)
		return models.CalculationType_count{Values: []models.CalculationType{}, Count: 0, Auth: models.Auth{}}, err
	}

	return out_count, nil
}

//func (esv *CalculationTypeService) Add(ctx context.Context, ea models.CalculationType) (int, error)
func (esv *CalculationTypeService) Add(ctx context.Context, ea models.CalculationType) (int, error) {
	var est ifCalculationTypeStorage
	est = &esv.storage

	ai, err := est.Add(ctx, ea)

	if err != nil {
		log.Println("CalculationTypeStorage.Add", err)
		return 0, err
	}

	return ai, nil
}

//func (esv *CalculationTypeService) Upd(ctx context.Context, eu models.CalculationType) (int, error)
func (esv *CalculationTypeService) Upd(ctx context.Context, eu models.CalculationType) (int, error) {
	var est ifCalculationTypeStorage
	est = &esv.storage

	ui, err := est.Upd(ctx, eu)

	if err != nil {
		log.Println("CalculationTypeStorage.Upd", err)
		return 0, err
	}

	return ui, nil
}

//func (esv *CalculationTypeService) Del(ctx context.Context, ed []int) ([]int, error) {
func (esv *CalculationTypeService) Del(ctx context.Context, ed []int) ([]int, error) {
	var est ifCalculationTypeStorage
	est = &esv.storage

	res, err := est.Del(ctx, ed)

	if err != nil {
		log.Println("CalculationTypeStorage.Del", err)
		return []int{}, err
	}

	return res, nil
}

//func (esv *CalculationTypeService) GetOne(ctx context.Context, i int) (models.CalculationType_count, error)
func (esv *CalculationTypeService) GetOne(ctx context.Context, i int) (models.CalculationType_count, error) {
	var est ifCalculationTypeStorage
	est = &esv.storage

	out_count, err := est.GetOne(ctx, i)

	if err != nil {
		log.Println("CalculationTypeStorage.GetOne", err)
		return models.CalculationType_count{Values: []models.CalculationType{}, Count: 0, Auth: models.Auth{}}, err
	}

	return out_count, nil
}
