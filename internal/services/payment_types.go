package services

import (
	"context"
	"log"

	"github.com/u2lentaru/billbck/internal/adapters/db/pgsql"
	"github.com/u2lentaru/billbck/internal/models"
)

type PaymentTypeService struct {
	storage pgsql.PaymentTypeStorage
}

type ifPaymentTypeStorage interface {
	GetList(ctx context.Context, pg, pgs int, gs1 string, ord int, dsc bool) (models.PaymentType_count, error)
	Add(ctx context.Context, ea models.PaymentType) (int, error)
	Upd(ctx context.Context, eu models.PaymentType) (int, error)
	Del(ctx context.Context, ed []int) ([]int, error)
	GetOne(ctx context.Context, i int) (models.PaymentType_count, error)
}

//func NewPaymentTypeService(storage pgsql.PaymentTypeStorage) *PaymentTypeService
func NewPaymentTypeService(storage pgsql.PaymentTypeStorage) *PaymentTypeService {
	return &PaymentTypeService{storage}
}

//func (esv *PaymentTypeService) GetList(ctx context.Context, pg, pgs int, gs1 string, ord int, dsc bool) (models.PaymentType_count, error)
func (esv *PaymentTypeService) GetList(ctx context.Context, pg, pgs int, gs1 string, ord int, dsc bool) (models.PaymentType_count, error) {
	var est ifPaymentTypeStorage
	est = &esv.storage

	out_count, err := est.GetList(ctx, pg, pgs, gs1, ord, dsc)

	if err != nil {
		log.Println("PaymentTypeStorage.GetList", err)
		return models.PaymentType_count{Values: []models.PaymentType{}, Count: 0, Auth: models.Auth{}}, err
	}

	return out_count, nil
}

//func (esv *PaymentTypeService) Add(ctx context.Context, ea models.PaymentType) (int, error)
func (esv *PaymentTypeService) Add(ctx context.Context, ea models.PaymentType) (int, error) {
	var est ifPaymentTypeStorage
	est = &esv.storage

	ai, err := est.Add(ctx, ea)

	if err != nil {
		log.Println("PaymentTypeStorage.Add", err)
		return 0, err
	}

	return ai, nil
}

//func (esv *PaymentTypeService) Upd(ctx context.Context, eu models.PaymentType) (int, error)
func (esv *PaymentTypeService) Upd(ctx context.Context, eu models.PaymentType) (int, error) {
	var est ifPaymentTypeStorage
	est = &esv.storage

	ui, err := est.Upd(ctx, eu)

	if err != nil {
		log.Println("PaymentTypeStorage.Upd", err)
		return 0, err
	}

	return ui, nil
}

//func (esv *PaymentTypeService) Del(ctx context.Context, ed []int) ([]int, error) {
func (esv *PaymentTypeService) Del(ctx context.Context, ed []int) ([]int, error) {
	var est ifPaymentTypeStorage
	est = &esv.storage

	res, err := est.Del(ctx, ed)

	if err != nil {
		log.Println("PaymentTypeStorage.Del", err)
		return []int{}, err
	}

	return res, nil
}

//func (esv *PaymentTypeService) GetOne(ctx context.Context, i int) (models.PaymentType_count, error)
func (esv *PaymentTypeService) GetOne(ctx context.Context, i int) (models.PaymentType_count, error) {
	var est ifPaymentTypeStorage
	est = &esv.storage

	out_count, err := est.GetOne(ctx, i)

	if err != nil {
		log.Println("PaymentTypeStorage.GetOne", err)
		return models.PaymentType_count{Values: []models.PaymentType{}, Count: 0, Auth: models.Auth{}}, err
	}

	return out_count, nil
}
