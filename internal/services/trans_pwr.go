package services

import (
	"context"
	"log"

	"github.com/u2lentaru/billbck/internal/adapters/db/pgsql"
	"github.com/u2lentaru/billbck/internal/models"
)

type TransPwrService struct {
	storage pgsql.TransPwrStorage
}

type ifTransPwrStorage interface {
	GetList(ctx context.Context, pg, pgs int, gs1 string, ord int, dsc bool) (models.TransPwr_count, error)
	Add(ctx context.Context, ea models.TransPwr) (int, error)
	Upd(ctx context.Context, eu models.TransPwr) (int, error)
	Del(ctx context.Context, ed []int) ([]int, error)
	GetOne(ctx context.Context, i int) (models.TransPwr_count, error)
}

//func NewTransPwrService(storage pgsql.TransPwrtorage) *TransPwrService
func NewTransPwrService(storage pgsql.TransPwrStorage) *TransPwrService {
	return &TransPwrService{storage}
}

//func (esv *TransPwrService) GetList(ctx context.Context, pg, pgs int, gs1 string, ord int, dsc bool) (models.TransPwr_count, error)
func (esv *TransPwrService) GetList(ctx context.Context, pg, pgs int, gs1 string, ord int, dsc bool) (models.TransPwr_count, error) {
	var est ifTransPwrStorage
	est = &esv.storage

	out_count, err := est.GetList(ctx, pg, pgs, gs1, ord, dsc)

	if err != nil {
		log.Println("TransPwrStorage.GetList", err)
		return models.TransPwr_count{Values: []models.TransPwr{}, Count: 0, Auth: models.Auth{}}, err
	}

	return out_count, nil
}

//func (esv *TransPwrService) Add(ctx context.Context, ea models.TransPwr) (int, error)
func (esv *TransPwrService) Add(ctx context.Context, ea models.TransPwr) (int, error) {
	var est ifTransPwrStorage
	est = &esv.storage

	ai, err := est.Add(ctx, ea)

	if err != nil {
		log.Println("TransPwrStorage.Add", err)
		return 0, err
	}

	return ai, nil
}

//func (esv *TransPwrService) Upd(ctx context.Context, eu models.TransPwr) (int, error)
func (esv *TransPwrService) Upd(ctx context.Context, eu models.TransPwr) (int, error) {
	var est ifTransPwrStorage
	est = &esv.storage

	ui, err := est.Upd(ctx, eu)

	if err != nil {
		log.Println("TransPwrStorage.Upd", err)
		return 0, err
	}

	return ui, nil
}

//func (esv *TransPwrService) Del(ctx context.Context, ed []int) ([]int, error) {
func (esv *TransPwrService) Del(ctx context.Context, ed []int) ([]int, error) {
	var est ifTransPwrStorage
	est = &esv.storage

	res, err := est.Del(ctx, ed)

	if err != nil {
		log.Println("TransPwrStorage.Del", err)
		return []int{}, err
	}

	return res, nil
}

//func (esv *TransPwrService) GetOne(ctx context.Context, i int) (models.TransPwr_count, error)
func (esv *TransPwrService) GetOne(ctx context.Context, i int) (models.TransPwr_count, error) {
	var est ifTransPwrStorage
	est = &esv.storage

	out_count, err := est.GetOne(ctx, i)

	if err != nil {
		log.Println("TransPwrStorage.GetOne", err)
		return models.TransPwr_count{Values: []models.TransPwr{}, Count: 0, Auth: models.Auth{}}, err
	}

	return out_count, nil
}
