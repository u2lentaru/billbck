package services

import (
	"context"
	"log"

	"github.com/u2lentaru/billbck/internal/adapters/db/pgsql"
	"github.com/u2lentaru/billbck/internal/models"
)

type ResultService struct {
	storage pgsql.ResultStorage
}

type ifResultStorage interface {
	GetList(ctx context.Context, pg, pgs int, gs1 string, ord int, dsc bool) (models.Result_count, error)
	Add(ctx context.Context, ea models.Result) (int, error)
	Upd(ctx context.Context, eu models.Result) (int, error)
	Del(ctx context.Context, ed []int) ([]int, error)
	GetOne(ctx context.Context, i int) (models.Result_count, error)
}

//func NewResultService(storage pgsql.Resulttorage) *ResultService
func NewResultService(storage pgsql.ResultStorage) *ResultService {
	return &ResultService{storage}
}

//func (esv *ResultService) GetList(ctx context.Context, pg, pgs int, gs1 string, ord int, dsc bool) (models.Result_count, error)
func (esv *ResultService) GetList(ctx context.Context, pg, pgs int, gs1 string, ord int, dsc bool) (models.Result_count, error) {
	var est ifResultStorage
	est = &esv.storage

	out_count, err := est.GetList(ctx, pg, pgs, gs1, ord, dsc)

	if err != nil {
		log.Println("ResultStorage.GetList", err)
		return models.Result_count{Values: []models.Result{}, Count: 0, Auth: models.Auth{}}, err
	}

	return out_count, nil
}

//func (esv *ResultService) Add(ctx context.Context, ea models.Result) (int, error)
func (esv *ResultService) Add(ctx context.Context, ea models.Result) (int, error) {
	var est ifResultStorage
	est = &esv.storage

	ai, err := est.Add(ctx, ea)

	if err != nil {
		log.Println("ResultStorage.Add", err)
		return 0, err
	}

	return ai, nil
}

//func (esv *ResultService) Upd(ctx context.Context, eu models.Result) (int, error)
func (esv *ResultService) Upd(ctx context.Context, eu models.Result) (int, error) {
	var est ifResultStorage
	est = &esv.storage

	ui, err := est.Upd(ctx, eu)

	if err != nil {
		log.Println("ResultStorage.Upd", err)
		return 0, err
	}

	return ui, nil
}

//func (esv *ResultService) Del(ctx context.Context, ed []int) ([]int, error) {
func (esv *ResultService) Del(ctx context.Context, ed []int) ([]int, error) {
	var est ifResultStorage
	est = &esv.storage

	res, err := est.Del(ctx, ed)

	if err != nil {
		log.Println("ResultStorage.Del", err)
		return []int{}, err
	}

	return res, nil
}

//func (esv *ResultService) GetOne(ctx context.Context, i int) (models.Result_count, error)
func (esv *ResultService) GetOne(ctx context.Context, i int) (models.Result_count, error) {
	var est ifResultStorage
	est = &esv.storage

	out_count, err := est.GetOne(ctx, i)

	if err != nil {
		log.Println("ResultStorage.GetOne", err)
		return models.Result_count{Values: []models.Result{}, Count: 0, Auth: models.Auth{}}, err
	}

	return out_count, nil
}
