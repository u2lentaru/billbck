package services

import (
	"context"
	"log"

	"github.com/u2lentaru/billbck/internal/adapters/db/pgsql"
	"github.com/u2lentaru/billbck/internal/models"
)

type TransTypeService struct {
	storage pgsql.TransTypeStorage
}

type ifTransTypeStorage interface {
	GetList(ctx context.Context, pg, pgs int, gs1 string, ord int, dsc bool) (models.TransType_count, error)
	Add(ctx context.Context, ea models.TransType) (int, error)
	Upd(ctx context.Context, eu models.TransType) (int, error)
	Del(ctx context.Context, ed []int) ([]int, error)
	GetOne(ctx context.Context, i int) (models.TransType_count, error)
}

//func NewTransTypeService(storage pgsql.TransTypetorage) *TransTypeService
func NewTransTypeService(storage pgsql.TransTypeStorage) *TransTypeService {
	return &TransTypeService{storage}
}

//func (esv *TransTypeService) GetList(ctx context.Context, pg, pgs int, gs1 string, ord int, dsc bool) (models.TransType_count, error)
func (esv *TransTypeService) GetList(ctx context.Context, pg, pgs int, gs1 string, ord int, dsc bool) (models.TransType_count, error) {
	var est ifTransTypeStorage
	est = &esv.storage

	out_count, err := est.GetList(ctx, pg, pgs, gs1, ord, dsc)

	if err != nil {
		log.Println("TransTypeStorage.GetList", err)
		return models.TransType_count{Values: []models.TransType{}, Count: 0, Auth: models.Auth{}}, err
	}

	return out_count, nil
}

//func (esv *TransTypeService) Add(ctx context.Context, ea models.TransType) (int, error)
func (esv *TransTypeService) Add(ctx context.Context, ea models.TransType) (int, error) {
	var est ifTransTypeStorage
	est = &esv.storage

	ai, err := est.Add(ctx, ea)

	if err != nil {
		log.Println("TransTypeStorage.Add", err)
		return 0, err
	}

	return ai, nil
}

//func (esv *TransTypeService) Upd(ctx context.Context, eu models.TransType) (int, error)
func (esv *TransTypeService) Upd(ctx context.Context, eu models.TransType) (int, error) {
	var est ifTransTypeStorage
	est = &esv.storage

	ui, err := est.Upd(ctx, eu)

	if err != nil {
		log.Println("TransTypeStorage.Upd", err)
		return 0, err
	}

	return ui, nil
}

//func (esv *TransTypeService) Del(ctx context.Context, ed []int) ([]int, error) {
func (esv *TransTypeService) Del(ctx context.Context, ed []int) ([]int, error) {
	var est ifTransTypeStorage
	est = &esv.storage

	res, err := est.Del(ctx, ed)

	if err != nil {
		log.Println("TransTypeStorage.Del", err)
		return []int{}, err
	}

	return res, nil
}

//func (esv *TransTypeService) GetOne(ctx context.Context, i int) (models.TransType_count, error)
func (esv *TransTypeService) GetOne(ctx context.Context, i int) (models.TransType_count, error) {
	var est ifTransTypeStorage
	est = &esv.storage

	out_count, err := est.GetOne(ctx, i)

	if err != nil {
		log.Println("TransTypeStorage.GetOne", err)
		return models.TransType_count{Values: []models.TransType{}, Count: 0, Auth: models.Auth{}}, err
	}

	return out_count, nil
}
