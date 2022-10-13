package services

import (
	"context"
	"log"

	"github.com/u2lentaru/billbck/internal/adapters/db/pgsql"
	"github.com/u2lentaru/billbck/internal/models"
)

type SubTypeService struct {
	storage pgsql.SubTypeStorage
}

type ifSubTypeStorage interface {
	GetList(ctx context.Context, pg, pgs int, gs1, gs2 string, ord int, dsc bool) (models.SubType_count, error)
	Add(ctx context.Context, ea models.SubType) (int, error)
	Upd(ctx context.Context, eu models.SubType) (int, error)
	Del(ctx context.Context, ed []int) ([]int, error)
	GetOne(ctx context.Context, i int) (models.SubType_count, error)
}

//NewSubTypeService(storage pg.SubTypeStorage) *SubTypeService
func NewSubTypeService(storage pgsql.SubTypeStorage) *SubTypeService {
	return &SubTypeService{storage}
}

//func (esv *SubTypeService) GetList(ctx context.Context, pg, pgs int, gs1, gs2 string, ord int, dsc bool) (models.SubType_count, error)
func (esv *SubTypeService) GetList(ctx context.Context, pg, pgs int, gs1, gs2 string, ord int, dsc bool) (models.SubType_count, error) {
	var est ifSubTypeStorage
	est = &esv.storage

	out_count, err := est.GetList(ctx, pg, pgs, gs1, gs2, ord, dsc)

	if err != nil {
		log.Println("SubTypeStorage.GetList", err)
		return models.SubType_count{Values: []models.SubType{}, Count: 0, Auth: models.Auth{}}, err
	}

	return out_count, nil
}

//func (esv *SubTypeService) Add(ctx context.Context, ea models.SubType) (int, error)
func (esv *SubTypeService) Add(ctx context.Context, ea models.SubType) (int, error) {
	var est ifSubTypeStorage
	est = &esv.storage

	ai, err := est.Add(ctx, ea)

	if err != nil {
		log.Println("SubTypeStorage.Add", err)
		return 0, err
	}

	return ai, nil
}

//func (esv *SubTypeService) Upd(ctx context.Context, eu models.SubType) (int, error)
func (esv *SubTypeService) Upd(ctx context.Context, eu models.SubType) (int, error) {
	var est ifSubTypeStorage
	est = &esv.storage

	ui, err := est.Upd(ctx, eu)

	if err != nil {
		log.Println("SubTypeStorage.Upd", err)
		return 0, err
	}

	return ui, nil
}

//func (esv *SubTypeService) Del(ctx context.Context, ed []int) ([]int, error)
func (esv *SubTypeService) Del(ctx context.Context, ed []int) ([]int, error) {
	var est ifSubTypeStorage
	est = &esv.storage

	res, err := est.Del(ctx, ed)

	if err != nil {
		log.Println("SubTypeStorage.Del", err)
		return []int{}, err
	}

	return res, nil
}

//func (esv *SubTypeService) GetOne(ctx context.Context, i int) (models.SubType_count, error)
func (esv *SubTypeService) GetOne(ctx context.Context, i int) (models.SubType_count, error) {
	var est ifSubTypeStorage
	est = &esv.storage

	out_count, err := est.GetOne(ctx, i)

	if err != nil {
		log.Println("SubTypeStorage.GetOne", err)
		return models.SubType_count{Values: []models.SubType{}, Count: 0, Auth: models.Auth{}}, err
	}

	return out_count, nil
}
