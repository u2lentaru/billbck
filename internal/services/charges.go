package services

import (
	"context"
	"log"

	"github.com/u2lentaru/billbck/internal/adapters/db/pgsql"
	"github.com/u2lentaru/billbck/internal/models"
)

type ChargeService struct {
	storage pgsql.ChargeStorage
}

type ifChargeStorage interface {
	GetList(ctx context.Context, pg, pgs int, gs1, gs2 string, ord int, dsc bool) (models.Charge_count, error)
	Add(ctx context.Context, ea models.Charge) (int, error)
	Upd(ctx context.Context, eu models.Charge) (int, error)
	Del(ctx context.Context, ed []int) ([]int, error)
	GetOne(ctx context.Context, i int) (models.Charge_count, error)
	ChargeRun(ctx context.Context, i int) (int, error)
}

//NewChargeService(storage pg.ChargeStorage) *ChargeService
func NewChargeService(storage pgsql.ChargeStorage) *ChargeService {
	return &ChargeService{storage}
}

//func (esv *ChargeService) GetList(ctx context.Context, pg, pgs int, gs1, gs2 string, ord int, dsc bool) (models.Charge_count, error)
func (esv *ChargeService) GetList(ctx context.Context, pg, pgs int, gs1, gs2 string, ord int, dsc bool) (models.Charge_count, error) {
	var est ifChargeStorage
	est = &esv.storage

	out_count, err := est.GetList(ctx, pg, pgs, gs1, gs2, ord, dsc)

	if err != nil {
		log.Println("ChargeStorage.GetList", err)
		return models.Charge_count{Values: []models.Charge{}, Count: 0, Auth: models.Auth{}}, err
	}

	return out_count, nil
}

//func (esv *ChargeService) Add(ctx context.Context, ea models.Charge) (int, error)
func (esv *ChargeService) Add(ctx context.Context, ea models.Charge) (int, error) {
	var est ifChargeStorage
	est = &esv.storage

	ai, err := est.Add(ctx, ea)

	if err != nil {
		log.Println("ChargeStorage.Add", err)
		return 0, err
	}

	return ai, nil
}

//func (esv *ChargeService) Upd(ctx context.Context, eu models.Charge) (int, error)
func (esv *ChargeService) Upd(ctx context.Context, eu models.Charge) (int, error) {
	var est ifChargeStorage
	est = &esv.storage

	ui, err := est.Upd(ctx, eu)

	if err != nil {
		log.Println("ChargeStorage.Upd", err)
		return 0, err
	}

	return ui, nil
}

//func (esv *ChargeService) Del(ctx context.Context, ed []int) ([]int, error)
func (esv *ChargeService) Del(ctx context.Context, ed []int) ([]int, error) {
	var est ifChargeStorage
	est = &esv.storage

	res, err := est.Del(ctx, ed)

	if err != nil {
		log.Println("ChargeStorage.Del", err)
		return []int{}, err
	}

	return res, nil
}

//func (esv *ChargeService) GetOne(ctx context.Context, i int) (models.Charge_count, error)
func (esv *ChargeService) GetOne(ctx context.Context, i int) (models.Charge_count, error) {
	var est ifChargeStorage
	est = &esv.storage

	out_count, err := est.GetOne(ctx, i)

	if err != nil {
		log.Println("ChargeStorage.GetOne", err)
		return models.Charge_count{Values: []models.Charge{}, Count: 0, Auth: models.Auth{}}, err
	}

	return out_count, nil
}

//func (esv *ChargeService) ChargeRun(ctx context.Context, i int) (int, error)
func (esv *ChargeService) ChargeRun(ctx context.Context, i int) (int, error) {
	var est ifChargeStorage
	est = &esv.storage

	pr, err := est.ChargeRun(ctx, i)

	if err != nil {
		log.Println("ChargeStorage.ChargeRun", err)
		return 0, err
	}

	return pr, nil
}
