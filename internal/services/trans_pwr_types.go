package services

import (
	"context"
	"log"

	"github.com/u2lentaru/billbck/internal/adapters/db/pgsql"
	"github.com/u2lentaru/billbck/internal/models"
)

type TransPwrTypeService struct {
	storage pgsql.TransPwrTypeStorage
}

type ifTransPwrTypeStorage interface {
	GetList(ctx context.Context, pg, pgs int, gs1 string, ord int, dsc bool) (models.TransPwrType_count, error)
	Add(ctx context.Context, ea models.TransPwrType) (int, error)
	Upd(ctx context.Context, eu models.TransPwrType) (int, error)
	Del(ctx context.Context, ed []int) ([]int, error)
	GetOne(ctx context.Context, i int) (models.TransPwrType_count, error)
}

//func NewTransPwrTypeService(storage pgsql.TransPwrTypetorage) *TransPwrTypeService
func NewTransPwrTypeService(storage pgsql.TransPwrTypeStorage) *TransPwrTypeService {
	return &TransPwrTypeService{storage}
}

//func (esv *TransPwrTypeService) GetList(ctx context.Context, pg, pgs int, gs1 string, ord int, dsc bool) (models.TransPwrType_count, error)
func (esv *TransPwrTypeService) GetList(ctx context.Context, pg, pgs int, gs1 string, ord int, dsc bool) (models.TransPwrType_count, error) {
	var est ifTransPwrTypeStorage
	est = &esv.storage

	out_count, err := est.GetList(ctx, pg, pgs, gs1, ord, dsc)

	if err != nil {
		log.Println("TransPwrTypeStorage.GetList", err)
		return models.TransPwrType_count{Values: []models.TransPwrType{}, Count: 0, Auth: models.Auth{}}, err
	}

	return out_count, nil
}

//func (esv *TransPwrTypeService) Add(ctx context.Context, ea models.TransPwrType) (int, error)
func (esv *TransPwrTypeService) Add(ctx context.Context, ea models.TransPwrType) (int, error) {
	var est ifTransPwrTypeStorage
	est = &esv.storage

	ai, err := est.Add(ctx, ea)

	if err != nil {
		log.Println("TransPwrTypeStorage.Add", err)
		return 0, err
	}

	return ai, nil
}

//func (esv *TransPwrTypeService) Upd(ctx context.Context, eu models.TransPwrType) (int, error)
func (esv *TransPwrTypeService) Upd(ctx context.Context, eu models.TransPwrType) (int, error) {
	var est ifTransPwrTypeStorage
	est = &esv.storage

	ui, err := est.Upd(ctx, eu)

	if err != nil {
		log.Println("TransPwrTypeStorage.Upd", err)
		return 0, err
	}

	return ui, nil
}

//func (esv *TransPwrTypeService) Del(ctx context.Context, ed []int) ([]int, error) {
func (esv *TransPwrTypeService) Del(ctx context.Context, ed []int) ([]int, error) {
	var est ifTransPwrTypeStorage
	est = &esv.storage

	res, err := est.Del(ctx, ed)

	if err != nil {
		log.Println("TransPwrTypeStorage.Del", err)
		return []int{}, err
	}

	return res, nil
}

//func (esv *TransPwrTypeService) GetOne(ctx context.Context, i int) (models.TransPwrType_count, error)
func (esv *TransPwrTypeService) GetOne(ctx context.Context, i int) (models.TransPwrType_count, error) {
	var est ifTransPwrTypeStorage
	est = &esv.storage

	out_count, err := est.GetOne(ctx, i)

	if err != nil {
		log.Println("TransPwrTypeStorage.GetOne", err)
		return models.TransPwrType_count{Values: []models.TransPwrType{}, Count: 0, Auth: models.Auth{}}, err
	}

	return out_count, nil
}
