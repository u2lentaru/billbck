package services

import (
	"context"
	"log"

	"github.com/u2lentaru/billbck/internal/adapters/db/pgsql"
	"github.com/u2lentaru/billbck/internal/models"
)

type TguService struct {
	storage pgsql.TguStorage
}

type ifTguStorage interface {
	GetList(ctx context.Context, pg, pgs int, gs1, gs2 string, ord int, dsc bool) (models.Tgu_count, error)
	Add(ctx context.Context, ea models.Tgu) (int, error)
	Upd(ctx context.Context, eu models.Tgu) (int, error)
	Del(ctx context.Context, ed []int) ([]int, error)
	GetOne(ctx context.Context, i int) (models.Tgu_count, error)
}

//NewTguService(storage pg.TguStorage) *TguService
func NewTguService(storage pgsql.TguStorage) *TguService {
	return &TguService{storage}
}

//func (esv *TguService) GetList(ctx context.Context, pg, pgs int, gs1, gs2 string, ord int, dsc bool) (models.Tgu_count, error)
func (esv *TguService) GetList(ctx context.Context, pg, pgs int, gs1, gs2 string, ord int, dsc bool) (models.Tgu_count, error) {
	var est ifTguStorage
	est = &esv.storage

	out_count, err := est.GetList(ctx, pg, pgs, gs1, gs2, ord, dsc)

	if err != nil {
		log.Println("TguStorage.GetList", err)
		return models.Tgu_count{Values: []models.Tgu{}, Count: 0, Auth: models.Auth{}}, err
	}

	return out_count, nil
}

//func (esv *TguService) Add(ctx context.Context, ea models.Tgu) (int, error)
func (esv *TguService) Add(ctx context.Context, ea models.Tgu) (int, error) {
	var est ifTguStorage
	est = &esv.storage

	ai, err := est.Add(ctx, ea)

	if err != nil {
		log.Println("TguStorage.Add", err)
		return 0, err
	}

	return ai, nil
}

//func (esv *TguService) Upd(ctx context.Context, eu models.Tgu) (int, error)
func (esv *TguService) Upd(ctx context.Context, eu models.Tgu) (int, error) {
	var est ifTguStorage
	est = &esv.storage

	ui, err := est.Upd(ctx, eu)

	if err != nil {
		log.Println("TguStorage.Upd", err)
		return 0, err
	}

	return ui, nil
}

//func (esv *TguService) Del(ctx context.Context, ed []int) ([]int, error)
func (esv *TguService) Del(ctx context.Context, ed []int) ([]int, error) {
	var est ifTguStorage
	est = &esv.storage

	res, err := est.Del(ctx, ed)

	if err != nil {
		log.Println("TguStorage.Del", err)
		return []int{}, err
	}

	return res, nil
}

//func (esv *TguService) GetOne(ctx context.Context, i int) (models.Tgu_count, error)
func (esv *TguService) GetOne(ctx context.Context, i int) (models.Tgu_count, error) {
	var est ifTguStorage
	est = &esv.storage

	out_count, err := est.GetOne(ctx, i)

	if err != nil {
		log.Println("TguStorage.GetOne", err)
		return models.Tgu_count{Values: []models.Tgu{}, Count: 0, Auth: models.Auth{}}, err
	}

	return out_count, nil
}
