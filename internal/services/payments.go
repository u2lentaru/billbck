package services

import (
	"context"
	"log"

	"github.com/u2lentaru/billbck/internal/adapters/db/pgsql"
	"github.com/u2lentaru/billbck/internal/models"
)

type PaymentService struct {
	storage pgsql.PaymentStorage
}

type ifPaymentStorage interface {
	GetList(ctx context.Context, pg, pgs int, gs1, gs2 string, ord int, dsc bool) (models.Payment_count, error)
	Add(ctx context.Context, ea models.Payment) (int, error)
	Upd(ctx context.Context, eu models.Payment) (int, error)
	Del(ctx context.Context, ed []int) ([]int, error)
	GetOne(ctx context.Context, i int) (models.Payment_count, error)
}

//NewPaymentService(storage pg.PaymentStorage) *PaymentService
func NewPaymentService(storage pgsql.PaymentStorage) *PaymentService {
	return &PaymentService{storage}
}

//func (esv *PaymentService) GetList(ctx context.Context, pg, pgs int, gs1, gs2 string, ord int, dsc bool) (models.Payment_count, error)
func (esv *PaymentService) GetList(ctx context.Context, pg, pgs int, gs1, gs2 string, ord int, dsc bool) (models.Payment_count, error) {
	var est ifPaymentStorage
	est = &esv.storage

	out_count, err := est.GetList(ctx, pg, pgs, gs1, gs2, ord, dsc)

	if err != nil {
		log.Println("PaymentStorage.GetList", err)
		return models.Payment_count{Values: []models.Payment{}, Count: 0, Auth: models.Auth{}}, err
	}

	return out_count, nil
}

//func (esv *PaymentService) Add(ctx context.Context, ea models.Payment) (int, error)
func (esv *PaymentService) Add(ctx context.Context, ea models.Payment) (int, error) {
	var est ifPaymentStorage
	est = &esv.storage

	ai, err := est.Add(ctx, ea)

	if err != nil {
		log.Println("PaymentStorage.Add", err)
		return 0, err
	}

	return ai, nil
}

//func (esv *PaymentService) Upd(ctx context.Context, eu models.Payment) (int, error)
func (esv *PaymentService) Upd(ctx context.Context, eu models.Payment) (int, error) {
	var est ifPaymentStorage
	est = &esv.storage

	ui, err := est.Upd(ctx, eu)

	if err != nil {
		log.Println("PaymentStorage.Upd", err)
		return 0, err
	}

	return ui, nil
}

//func (esv *PaymentService) Del(ctx context.Context, ed []int) ([]int, error)
func (esv *PaymentService) Del(ctx context.Context, ed []int) ([]int, error) {
	var est ifPaymentStorage
	est = &esv.storage

	res, err := est.Del(ctx, ed)

	if err != nil {
		log.Println("PaymentStorage.Del", err)
		return []int{}, err
	}

	return res, nil
}

//func (esv *PaymentService) GetOne(ctx context.Context, i int) (models.Payment_count, error)
func (esv *PaymentService) GetOne(ctx context.Context, i int) (models.Payment_count, error) {
	var est ifPaymentStorage
	est = &esv.storage

	out_count, err := est.GetOne(ctx, i)

	if err != nil {
		log.Println("PaymentStorage.GetOne", err)
		return models.Payment_count{Values: []models.Payment{}, Count: 0, Auth: models.Auth{}}, err
	}

	return out_count, nil
}
