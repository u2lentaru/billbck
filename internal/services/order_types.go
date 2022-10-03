package services

import (
	"context"
	"log"

	"github.com/u2lentaru/billbck/internal/adapters/db/pgsql"
	"github.com/u2lentaru/billbck/internal/models"
)

type OrderTypeService struct {
	storage pgsql.OrderTypeStorage
}

type ifOrderTypeStorage interface {
	GetList(ctx context.Context, pg, pgs int, gs1 string, ord int, dsc bool) (models.OrderType_count, error)
	Add(ctx context.Context, ea models.OrderType) (int, error)
	Upd(ctx context.Context, eu models.OrderType) (int, error)
	Del(ctx context.Context, ed []int) ([]int, error)
	GetOne(ctx context.Context, i int) (models.OrderType_count, error)
}

//func NewOrderTypeService(storage pgsql.OrderTypeStorage) *OrderTypeService
func NewOrderTypeService(storage pgsql.OrderTypeStorage) *OrderTypeService {
	return &OrderTypeService{storage}
}

//func (esv *OrderTypeService) GetList(ctx context.Context, pg, pgs int, gs1 string, ord int, dsc bool) (models.OrderType_count, error)
func (esv *OrderTypeService) GetList(ctx context.Context, pg, pgs int, gs1 string, ord int, dsc bool) (models.OrderType_count, error) {
	var est ifOrderTypeStorage
	est = &esv.storage

	out_count, err := est.GetList(ctx, pg, pgs, gs1, ord, dsc)

	if err != nil {
		log.Println("OrderTypeStorage.GetList", err)
		return models.OrderType_count{Values: []models.OrderType{}, Count: 0, Auth: models.Auth{}}, err
	}

	return out_count, nil
}

//func (esv *OrderTypeService) Add(ctx context.Context, ea models.OrderType) (int, error)
func (esv *OrderTypeService) Add(ctx context.Context, ea models.OrderType) (int, error) {
	var est ifOrderTypeStorage
	est = &esv.storage

	ai, err := est.Add(ctx, ea)

	if err != nil {
		log.Println("OrderTypeStorage.Add", err)
		return 0, err
	}

	return ai, nil
}

//func (esv *OrderTypeService) Upd(ctx context.Context, eu models.OrderType) (int, error)
func (esv *OrderTypeService) Upd(ctx context.Context, eu models.OrderType) (int, error) {
	var est ifOrderTypeStorage
	est = &esv.storage

	ui, err := est.Upd(ctx, eu)

	if err != nil {
		log.Println("OrderTypeStorage.Upd", err)
		return 0, err
	}

	return ui, nil
}

//func (esv *OrderTypeService) Del(ctx context.Context, ed []int) ([]int, error) {
func (esv *OrderTypeService) Del(ctx context.Context, ed []int) ([]int, error) {
	var est ifOrderTypeStorage
	est = &esv.storage

	res, err := est.Del(ctx, ed)

	if err != nil {
		log.Println("OrderTypeStorage.Del", err)
		return []int{}, err
	}

	return res, nil
}

//func (esv *OrderTypeService) GetOne(ctx context.Context, i int) (models.OrderType_count, error)
func (esv *OrderTypeService) GetOne(ctx context.Context, i int) (models.OrderType_count, error) {
	var est ifOrderTypeStorage
	est = &esv.storage

	out_count, err := est.GetOne(ctx, i)

	if err != nil {
		log.Println("OrderTypeStorage.GetOne", err)
		return models.OrderType_count{Values: []models.OrderType{}, Count: 0, Auth: models.Auth{}}, err
	}

	return out_count, nil
}
