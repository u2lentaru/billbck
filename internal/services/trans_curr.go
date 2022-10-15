package services

import (
	"context"
	"log"

	"github.com/u2lentaru/billbck/internal/adapters/db/pgsql"
	"github.com/u2lentaru/billbck/internal/models"
)

type TransCurrService struct {
	storage pgsql.TransCurrStorage
}

type ifTransCurrStorage interface {
	GetList(ctx context.Context, pg, pgs int, gs1 string, ord int, dsc bool) (models.TransCurr_count, error)
	Add(ctx context.Context, ea models.TransCurr) (int, error)
	Upd(ctx context.Context, eu models.TransCurr) (int, error)
	Del(ctx context.Context, ed []int) ([]int, error)
	GetOne(ctx context.Context, i int) (models.TransCurr_count, error)
}

//func NewTransCurrService(storage pgsql.TransCurrtorage) *TransCurrService
func NewTransCurrService(storage pgsql.TransCurrStorage) *TransCurrService {
	return &TransCurrService{storage}
}

//func (esv *TransCurrService) GetList(ctx context.Context, pg, pgs int, gs1 string, ord int, dsc bool) (models.TransCurr_count, error)
func (esv *TransCurrService) GetList(ctx context.Context, pg, pgs int, gs1 string, ord int, dsc bool) (models.TransCurr_count, error) {
	var est ifTransCurrStorage
	est = &esv.storage

	out_count, err := est.GetList(ctx, pg, pgs, gs1, ord, dsc)

	if err != nil {
		log.Println("TransCurrStorage.GetList", err)
		return models.TransCurr_count{Values: []models.TransCurr{}, Count: 0, Auth: models.Auth{}}, err
	}

	return out_count, nil
}

//func (esv *TransCurrService) Add(ctx context.Context, ea models.TransCurr) (int, error)
func (esv *TransCurrService) Add(ctx context.Context, ea models.TransCurr) (int, error) {
	var est ifTransCurrStorage
	est = &esv.storage

	ai, err := est.Add(ctx, ea)

	if err != nil {
		log.Println("TransCurrStorage.Add", err)
		return 0, err
	}

	return ai, nil
}

//func (esv *TransCurrService) Upd(ctx context.Context, eu models.TransCurr) (int, error)
func (esv *TransCurrService) Upd(ctx context.Context, eu models.TransCurr) (int, error) {
	var est ifTransCurrStorage
	est = &esv.storage

	ui, err := est.Upd(ctx, eu)

	if err != nil {
		log.Println("TransCurrStorage.Upd", err)
		return 0, err
	}

	return ui, nil
}

//func (esv *TransCurrService) Del(ctx context.Context, ed []int) ([]int, error) {
func (esv *TransCurrService) Del(ctx context.Context, ed []int) ([]int, error) {
	var est ifTransCurrStorage
	est = &esv.storage

	res, err := est.Del(ctx, ed)

	if err != nil {
		log.Println("TransCurrStorage.Del", err)
		return []int{}, err
	}

	return res, nil
}

//func (esv *TransCurrService) GetOne(ctx context.Context, i int) (models.TransCurr_count, error)
func (esv *TransCurrService) GetOne(ctx context.Context, i int) (models.TransCurr_count, error) {
	var est ifTransCurrStorage
	est = &esv.storage

	out_count, err := est.GetOne(ctx, i)

	if err != nil {
		log.Println("TransCurrStorage.GetOne", err)
		return models.TransCurr_count{Values: []models.TransCurr{}, Count: 0, Auth: models.Auth{}}, err
	}

	return out_count, nil
}
