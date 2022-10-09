package services

import (
	"context"
	"log"

	"github.com/u2lentaru/billbck/internal/adapters/db/pgsql"
	"github.com/u2lentaru/billbck/internal/models"
)

type RequestTypeService struct {
	storage pgsql.RequestTypeStorage
}

type ifRequestTypeStorage interface {
	GetList(ctx context.Context, pg, pgs int, gs1, gs2 string, ord int, dsc bool) (models.RequestType_count, error)
	Add(ctx context.Context, ea models.RequestType) (int, error)
	Upd(ctx context.Context, eu models.RequestType) (int, error)
	Del(ctx context.Context, ed []int) ([]int, error)
	GetOne(ctx context.Context, i int) (models.RequestType_count, error)
}

//NewRequestTypeService(storage pg.RequestTypeStorage) *RequestTypeService
func NewRequestTypeService(storage pgsql.RequestTypeStorage) *RequestTypeService {
	return &RequestTypeService{storage}
}

//func (esv *RequestTypeService) GetList(ctx context.Context, pg, pgs int, gs1, gs2 string, ord int, dsc bool) (models.RequestType_count, error)
func (esv *RequestTypeService) GetList(ctx context.Context, pg, pgs int, gs1, gs2 string, ord int, dsc bool) (models.RequestType_count, error) {
	var est ifRequestTypeStorage
	est = &esv.storage

	out_count, err := est.GetList(ctx, pg, pgs, gs1, gs2, ord, dsc)

	if err != nil {
		log.Println("RequestTypeStorage.GetList", err)
		return models.RequestType_count{Values: []models.RequestType{}, Count: 0, Auth: models.Auth{}}, err
	}

	return out_count, nil
}

//func (esv *RequestTypeService) Add(ctx context.Context, ea models.RequestType) (int, error)
func (esv *RequestTypeService) Add(ctx context.Context, ea models.RequestType) (int, error) {
	var est ifRequestTypeStorage
	est = &esv.storage

	ai, err := est.Add(ctx, ea)

	if err != nil {
		log.Println("RequestTypeStorage.Add", err)
		return 0, err
	}

	return ai, nil
}

//func (esv *RequestTypeService) Upd(ctx context.Context, eu models.RequestType) (int, error)
func (esv *RequestTypeService) Upd(ctx context.Context, eu models.RequestType) (int, error) {
	var est ifRequestTypeStorage
	est = &esv.storage

	ui, err := est.Upd(ctx, eu)

	if err != nil {
		log.Println("RequestTypeStorage.Upd", err)
		return 0, err
	}

	return ui, nil
}

//func (esv *RequestTypeService) Del(ctx context.Context, ed []int) ([]int, error)
func (esv *RequestTypeService) Del(ctx context.Context, ed []int) ([]int, error) {
	var est ifRequestTypeStorage
	est = &esv.storage

	res, err := est.Del(ctx, ed)

	if err != nil {
		log.Println("RequestTypeStorage.Del", err)
		return []int{}, err
	}

	return res, nil
}

//func (esv *RequestTypeService) GetOne(ctx context.Context, i int) (models.RequestType_count, error)
func (esv *RequestTypeService) GetOne(ctx context.Context, i int) (models.RequestType_count, error) {
	var est ifRequestTypeStorage
	est = &esv.storage

	out_count, err := est.GetOne(ctx, i)

	if err != nil {
		log.Println("RequestTypeStorage.GetOne", err)
		return models.RequestType_count{Values: []models.RequestType{}, Count: 0, Auth: models.Auth{}}, err
	}

	return out_count, nil
}
