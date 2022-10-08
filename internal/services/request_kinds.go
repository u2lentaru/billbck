package services

import (
	"context"
	"log"

	"github.com/u2lentaru/billbck/internal/adapters/db/pgsql"
	"github.com/u2lentaru/billbck/internal/models"
)

type RequestKindService struct {
	storage pgsql.RequestKindStorage
}

type ifRequestKindStorage interface {
	GetList(ctx context.Context, pg, pgs int, gs1 string, ord int, dsc bool) (models.RequestKind_count, error)
	Add(ctx context.Context, ea models.RequestKind) (int, error)
	Upd(ctx context.Context, eu models.RequestKind) (int, error)
	Del(ctx context.Context, ed []int) ([]int, error)
	GetOne(ctx context.Context, i int) (models.RequestKind_count, error)
}

//func NewRequestKindService(storage pgsql.RequestKindtorage) *RequestKindService
func NewRequestKindService(storage pgsql.RequestKindStorage) *RequestKindService {
	return &RequestKindService{storage}
}

//func (esv *RequestKindService) GetList(ctx context.Context, pg, pgs int, gs1 string, ord int, dsc bool) (models.RequestKind_count, error)
func (esv *RequestKindService) GetList(ctx context.Context, pg, pgs int, gs1 string, ord int, dsc bool) (models.RequestKind_count, error) {
	var est ifRequestKindStorage
	est = &esv.storage

	out_count, err := est.GetList(ctx, pg, pgs, gs1, ord, dsc)

	if err != nil {
		log.Println("RequestKindStorage.GetList", err)
		return models.RequestKind_count{Values: []models.RequestKind{}, Count: 0, Auth: models.Auth{}}, err
	}

	return out_count, nil
}

//func (esv *RequestKindService) Add(ctx context.Context, ea models.RequestKind) (int, error)
func (esv *RequestKindService) Add(ctx context.Context, ea models.RequestKind) (int, error) {
	var est ifRequestKindStorage
	est = &esv.storage

	ai, err := est.Add(ctx, ea)

	if err != nil {
		log.Println("RequestKindStorage.Add", err)
		return 0, err
	}

	return ai, nil
}

//func (esv *RequestKindService) Upd(ctx context.Context, eu models.RequestKind) (int, error)
func (esv *RequestKindService) Upd(ctx context.Context, eu models.RequestKind) (int, error) {
	var est ifRequestKindStorage
	est = &esv.storage

	ui, err := est.Upd(ctx, eu)

	if err != nil {
		log.Println("RequestKindStorage.Upd", err)
		return 0, err
	}

	return ui, nil
}

//func (esv *RequestKindService) Del(ctx context.Context, ed []int) ([]int, error) {
func (esv *RequestKindService) Del(ctx context.Context, ed []int) ([]int, error) {
	var est ifRequestKindStorage
	est = &esv.storage

	res, err := est.Del(ctx, ed)

	if err != nil {
		log.Println("RequestKindStorage.Del", err)
		return []int{}, err
	}

	return res, nil
}

//func (esv *RequestKindService) GetOne(ctx context.Context, i int) (models.RequestKind_count, error)
func (esv *RequestKindService) GetOne(ctx context.Context, i int) (models.RequestKind_count, error) {
	var est ifRequestKindStorage
	est = &esv.storage

	out_count, err := est.GetOne(ctx, i)

	if err != nil {
		log.Println("RequestKindStorage.GetOne", err)
		return models.RequestKind_count{Values: []models.RequestKind{}, Count: 0, Auth: models.Auth{}}, err
	}

	return out_count, nil
}
