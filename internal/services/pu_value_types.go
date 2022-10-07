package services

import (
	"context"
	"log"

	"github.com/u2lentaru/billbck/internal/adapters/db/pgsql"
	"github.com/u2lentaru/billbck/internal/models"
)

type PuValueTypeService struct {
	storage pgsql.PuValueTypeStorage
}

type ifPuValueTypeStorage interface {
	GetList(ctx context.Context, pg, pgs int, gs1 string, ord int, dsc bool) (models.PuValueType_count, error)
	Add(ctx context.Context, ea models.PuValueType) (int, error)
	Upd(ctx context.Context, eu models.PuValueType) (int, error)
	Del(ctx context.Context, ed []int) ([]int, error)
	GetOne(ctx context.Context, i int) (models.PuValueType_count, error)
}

//func NewPuValueTypeService(storage pgsql.PuValueTypetorage) *PuValueTypeService
func NewPuValueTypeService(storage pgsql.PuValueTypeStorage) *PuValueTypeService {
	return &PuValueTypeService{storage}
}

//func (esv *PuValueTypeService) GetList(ctx context.Context, pg, pgs int, gs1 string, ord int, dsc bool) (models.PuValueType_count, error)
func (esv *PuValueTypeService) GetList(ctx context.Context, pg, pgs int, gs1 string, ord int, dsc bool) (models.PuValueType_count, error) {
	var est ifPuValueTypeStorage
	est = &esv.storage

	out_count, err := est.GetList(ctx, pg, pgs, gs1, ord, dsc)

	if err != nil {
		log.Println("PuValueTypeStorage.GetList", err)
		return models.PuValueType_count{Values: []models.PuValueType{}, Count: 0, Auth: models.Auth{}}, err
	}

	return out_count, nil
}

//func (esv *PuValueTypeService) Add(ctx context.Context, ea models.PuValueType) (int, error)
func (esv *PuValueTypeService) Add(ctx context.Context, ea models.PuValueType) (int, error) {
	var est ifPuValueTypeStorage
	est = &esv.storage

	ai, err := est.Add(ctx, ea)

	if err != nil {
		log.Println("PuValueTypeStorage.Add", err)
		return 0, err
	}

	return ai, nil
}

//func (esv *PuValueTypeService) Upd(ctx context.Context, eu models.PuValueType) (int, error)
func (esv *PuValueTypeService) Upd(ctx context.Context, eu models.PuValueType) (int, error) {
	var est ifPuValueTypeStorage
	est = &esv.storage

	ui, err := est.Upd(ctx, eu)

	if err != nil {
		log.Println("PuValueTypeStorage.Upd", err)
		return 0, err
	}

	return ui, nil
}

//func (esv *PuValueTypeService) Del(ctx context.Context, ed []int) ([]int, error) {
func (esv *PuValueTypeService) Del(ctx context.Context, ed []int) ([]int, error) {
	var est ifPuValueTypeStorage
	est = &esv.storage

	res, err := est.Del(ctx, ed)

	if err != nil {
		log.Println("PuValueTypeStorage.Del", err)
		return []int{}, err
	}

	return res, nil
}

//func (esv *PuValueTypeService) GetOne(ctx context.Context, i int) (models.PuValueType_count, error)
func (esv *PuValueTypeService) GetOne(ctx context.Context, i int) (models.PuValueType_count, error) {
	var est ifPuValueTypeStorage
	est = &esv.storage

	out_count, err := est.GetOne(ctx, i)

	if err != nil {
		log.Println("PuValueTypeStorage.GetOne", err)
		return models.PuValueType_count{Values: []models.PuValueType{}, Count: 0, Auth: models.Auth{}}, err
	}

	return out_count, nil
}
