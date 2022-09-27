package services

import (
	"context"
	"log"

	"github.com/u2lentaru/billbck/internal/adapters/db/pgsql"
	"github.com/u2lentaru/billbck/internal/models"
)

type KskService struct {
	storage pgsql.KskStorage
}

type ifKskStorage interface {
	GetList(ctx context.Context, pg, pgs int, gs1, gs2 string, ord int, dsc bool) (models.Ksk_count, error)
	Add(ctx context.Context, ea models.Ksk) (int, error)
	Upd(ctx context.Context, eu models.Ksk) (int, error)
	Del(ctx context.Context, ed []int) ([]int, error)
	GetOne(ctx context.Context, i int) (models.Ksk_count, error)
}

//NewKskService(storage pgsql.KskStorage) *KskService
func NewKskService(storage pgsql.KskStorage) *KskService {
	return &KskService{storage}
}

//func (esv *KskService) GetList(ctx context.Context, pg, pgs int, gs1, gs2 string, ord int, dsc bool) (models.Ksk_count, error)
func (esv *KskService) GetList(ctx context.Context, pg, pgs int, gs1, gs2 string, ord int, dsc bool) (models.Ksk_count, error) {
	var est ifKskStorage
	est = &esv.storage

	out_count, err := est.GetList(ctx, pg, pgs, gs1, gs2, ord, dsc)

	if err != nil {
		log.Println("KskStorage.GetList", err)
		return models.Ksk_count{Values: []models.Ksk{}, Count: 0, Auth: models.Auth{}}, err
	}

	return out_count, nil
}

//func (esv *KskService) Add(ctx context.Context, ea models.Ksk) (int, error)
func (esv *KskService) Add(ctx context.Context, ea models.Ksk) (int, error) {
	var est ifKskStorage
	est = &esv.storage

	ai, err := est.Add(ctx, ea)

	if err != nil {
		log.Println("KskStorage.Add", err)
		return 0, err
	}

	return ai, nil
}

//func (esv *KskService) Upd(ctx context.Context, eu models.Ksk) (int, error)
func (esv *KskService) Upd(ctx context.Context, eu models.Ksk) (int, error) {
	var est ifKskStorage
	est = &esv.storage

	ui, err := est.Upd(ctx, eu)

	if err != nil {
		log.Println("KskStorage.Upd", err)
		return 0, err
	}

	return ui, nil
}

//func (esv *KskService) Del(ctx context.Context, ed []int) ([]int, error)
func (esv *KskService) Del(ctx context.Context, ed []int) ([]int, error) {
	var est ifKskStorage
	est = &esv.storage

	res, err := est.Del(ctx, ed)

	if err != nil {
		log.Println("KskStorage.Del", err)
		return []int{}, err
	}

	return res, nil
}

//func (esv *KskService) GetOne(ctx context.Context, i int) (models.Ksk_count, error)
func (esv *KskService) GetOne(ctx context.Context, i int) (models.Ksk_count, error) {
	var est ifKskStorage
	est = &esv.storage

	out_count, err := est.GetOne(ctx, i)

	if err != nil {
		log.Println("KskStorage.GetOne", err)
		return models.Ksk_count{Values: []models.Ksk{}, Count: 0, Auth: models.Auth{}}, err
	}

	return out_count, nil
}
