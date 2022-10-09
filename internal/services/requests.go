package services

import (
	"context"
	"log"

	"github.com/u2lentaru/billbck/internal/adapters/db/pgsql"
	"github.com/u2lentaru/billbck/internal/models"
)

type RequestService struct {
	storage pgsql.RequestStorage
}

type ifRequestStorage interface {
	GetList(ctx context.Context, pg, pgs int, gs1 string, ord int, dsc bool) (models.Request_count, error)
	Add(ctx context.Context, ea models.Request) (int, error)
	Upd(ctx context.Context, eu models.Request) (int, error)
	Del(ctx context.Context, ed []int) ([]int, error)
	GetOne(ctx context.Context, i int) (models.Request_count, error)
}

//func NewRequestService(storage pgsql.Requesttorage) *RequestService
func NewRequestService(storage pgsql.RequestStorage) *RequestService {
	return &RequestService{storage}
}

//func (esv *RequestService) GetList(ctx context.Context, pg, pgs int, gs1 string, ord int, dsc bool) (models.Request_count, error)
func (esv *RequestService) GetList(ctx context.Context, pg, pgs int, gs1 string, ord int, dsc bool) (models.Request_count, error) {
	var est ifRequestStorage
	est = &esv.storage

	out_count, err := est.GetList(ctx, pg, pgs, gs1, ord, dsc)

	if err != nil {
		log.Println("RequestStorage.GetList", err)
		return models.Request_count{Values: []models.Request{}, Count: 0, Auth: models.Auth{}}, err
	}

	return out_count, nil
}

//func (esv *RequestService) Add(ctx context.Context, ea models.Request) (int, error)
func (esv *RequestService) Add(ctx context.Context, ea models.Request) (int, error) {
	var est ifRequestStorage
	est = &esv.storage

	ai, err := est.Add(ctx, ea)

	if err != nil {
		log.Println("RequestStorage.Add", err)
		return 0, err
	}

	return ai, nil
}

//func (esv *RequestService) Upd(ctx context.Context, eu models.Request) (int, error)
func (esv *RequestService) Upd(ctx context.Context, eu models.Request) (int, error) {
	var est ifRequestStorage
	est = &esv.storage

	ui, err := est.Upd(ctx, eu)

	if err != nil {
		log.Println("RequestStorage.Upd", err)
		return 0, err
	}

	return ui, nil
}

//func (esv *RequestService) Del(ctx context.Context, ed []int) ([]int, error) {
func (esv *RequestService) Del(ctx context.Context, ed []int) ([]int, error) {
	var est ifRequestStorage
	est = &esv.storage

	res, err := est.Del(ctx, ed)

	if err != nil {
		log.Println("RequestStorage.Del", err)
		return []int{}, err
	}

	return res, nil
}

//func (esv *RequestService) GetOne(ctx context.Context, i int) (models.Request_count, error)
func (esv *RequestService) GetOne(ctx context.Context, i int) (models.Request_count, error) {
	var est ifRequestStorage
	est = &esv.storage

	out_count, err := est.GetOne(ctx, i)

	if err != nil {
		log.Println("RequestStorage.GetOne", err)
		return models.Request_count{Values: []models.Request{}, Count: 0, Auth: models.Auth{}}, err
	}

	return out_count, nil
}
