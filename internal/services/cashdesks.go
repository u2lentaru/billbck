package services

import (
	"context"
	"log"

	"github.com/u2lentaru/billbck/internal/adapters/db/pgsql"
	"github.com/u2lentaru/billbck/internal/models"
)

type CashdeskService struct {
	storage pgsql.CashdeskStorage
}

type ifCashdeskStorage interface {
	GetList(ctx context.Context, pg, pgs int, gs1 string, ord int, dsc bool) (models.Cashdesk_count, error)
	Add(ctx context.Context, ea models.Cashdesk) (int, error)
	Upd(ctx context.Context, eu models.Cashdesk) (int, error)
	Del(ctx context.Context, ed []int) ([]int, error)
	GetOne(ctx context.Context, i int) (models.Cashdesk_count, error)
}

//func NewCashdeskService(storage pgsql.CashdeskStorage) *CashdeskService
func NewCashdeskService(storage pgsql.CashdeskStorage) *CashdeskService {
	return &CashdeskService{storage}
}

//func (esv *CashdeskService) GetList(ctx context.Context, pg, pgs int, gs1 string, ord int, dsc bool) (models.Cashdesk_count, error)
func (esv *CashdeskService) GetList(ctx context.Context, pg, pgs int, gs1 string, ord int, dsc bool) (models.Cashdesk_count, error) {
	var est ifCashdeskStorage
	est = &esv.storage

	out_count, err := est.GetList(ctx, pg, pgs, gs1, ord, dsc)

	if err != nil {
		log.Println("CashdeskStorage.GetList", err)
		return models.Cashdesk_count{Values: []models.Cashdesk{}, Count: 0, Auth: models.Auth{}}, err
	}

	return out_count, nil
}

//func (esv *CashdeskService) Add(ctx context.Context, ea models.Cashdesk) (int, error)
func (esv *CashdeskService) Add(ctx context.Context, ea models.Cashdesk) (int, error) {
	var est ifCashdeskStorage
	est = &esv.storage

	ai, err := est.Add(ctx, ea)

	if err != nil {
		log.Println("CashdeskStorage.Add", err)
		return 0, err
	}

	return ai, nil
}

//func (esv *CashdeskService) Upd(ctx context.Context, eu models.Cashdesk) (int, error)
func (esv *CashdeskService) Upd(ctx context.Context, eu models.Cashdesk) (int, error) {
	var est ifCashdeskStorage
	est = &esv.storage

	ui, err := est.Upd(ctx, eu)

	if err != nil {
		log.Println("CashdeskStorage.Upd", err)
		return 0, err
	}

	return ui, nil
}

//func (esv *CashdeskService) Del(ctx context.Context, ed []int) ([]int, error) {
func (esv *CashdeskService) Del(ctx context.Context, ed []int) ([]int, error) {
	var est ifCashdeskStorage
	est = &esv.storage

	res, err := est.Del(ctx, ed)

	if err != nil {
		log.Println("CashdeskStorage.Del", err)
		return []int{}, err
	}

	return res, nil
}

//func (esv *CashdeskService) GetOne(ctx context.Context, i int) (models.Cashdesk_count, error)
func (esv *CashdeskService) GetOne(ctx context.Context, i int) (models.Cashdesk_count, error) {
	var est ifCashdeskStorage
	est = &esv.storage

	out_count, err := est.GetOne(ctx, i)

	if err != nil {
		log.Println("CashdeskStorage.GetOne", err)
		return models.Cashdesk_count{Values: []models.Cashdesk{}, Count: 0, Auth: models.Auth{}}, err
	}

	return out_count, nil

}
