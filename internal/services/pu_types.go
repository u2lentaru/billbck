package services

import (
	"context"
	"log"

	"github.com/u2lentaru/billbck/internal/adapters/db/pgsql"
	"github.com/u2lentaru/billbck/internal/models"
)

type PuTypeService struct {
	storage pgsql.PuTypeStorage
}

type ifPuTypeStorage interface {
	GetList(ctx context.Context, pg, pgs int, gs1 string, ord int, dsc bool) (models.PuType_count, error)
	Add(ctx context.Context, ea models.PuType) (int, error)
	Upd(ctx context.Context, eu models.PuType) (int, error)
	Del(ctx context.Context, ed []int) ([]int, error)
	GetOne(ctx context.Context, i int) (models.PuType_count, error)
}

//func NewPuTypeService(storage pgsql.PuTypetorage) *PuTypeService
func NewPuTypeService(storage pgsql.PuTypeStorage) *PuTypeService {
	return &PuTypeService{storage}
}

//func (esv *PuTypeService) GetList(ctx context.Context, pg, pgs int, gs1 string, ord int, dsc bool) (models.PuType_count, error)
func (esv *PuTypeService) GetList(ctx context.Context, pg, pgs int, gs1 string, ord int, dsc bool) (models.PuType_count, error) {
	var est ifPuTypeStorage
	est = &esv.storage

	out_count, err := est.GetList(ctx, pg, pgs, gs1, ord, dsc)

	if err != nil {
		log.Println("PuTypeStorage.GetList", err)
		return models.PuType_count{Values: []models.PuType{}, Count: 0, Auth: models.Auth{}}, err
	}

	return out_count, nil
}

//func (esv *PuTypeService) Add(ctx context.Context, ea models.PuType) (int, error)
func (esv *PuTypeService) Add(ctx context.Context, ea models.PuType) (int, error) {
	var est ifPuTypeStorage
	est = &esv.storage

	ai, err := est.Add(ctx, ea)

	if err != nil {
		log.Println("PuTypeStorage.Add", err)
		return 0, err
	}

	return ai, nil
}

//func (esv *PuTypeService) Upd(ctx context.Context, eu models.PuType) (int, error)
func (esv *PuTypeService) Upd(ctx context.Context, eu models.PuType) (int, error) {
	var est ifPuTypeStorage
	est = &esv.storage

	ui, err := est.Upd(ctx, eu)

	if err != nil {
		log.Println("PuTypeStorage.Upd", err)
		return 0, err
	}

	return ui, nil
}

//func (esv *PuTypeService) Del(ctx context.Context, ed []int) ([]int, error) {
func (esv *PuTypeService) Del(ctx context.Context, ed []int) ([]int, error) {
	var est ifPuTypeStorage
	est = &esv.storage

	res, err := est.Del(ctx, ed)

	if err != nil {
		log.Println("PuTypeStorage.Del", err)
		return []int{}, err
	}

	return res, nil
}

//func (esv *PuTypeService) GetOne(ctx context.Context, i int) (models.PuType_count, error)
func (esv *PuTypeService) GetOne(ctx context.Context, i int) (models.PuType_count, error) {
	var est ifPuTypeStorage
	est = &esv.storage

	out_count, err := est.GetOne(ctx, i)

	if err != nil {
		log.Println("PuTypeStorage.GetOne", err)
		return models.PuType_count{Values: []models.PuType{}, Count: 0, Auth: models.Auth{}}, err
	}

	return out_count, nil
}
