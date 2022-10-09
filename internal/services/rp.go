package services

import (
	"context"
	"log"

	"github.com/u2lentaru/billbck/internal/adapters/db/pgsql"
	"github.com/u2lentaru/billbck/internal/models"
)

type RpService struct {
	storage pgsql.RpStorage
}

type ifRpStorage interface {
	GetList(ctx context.Context, pg, pgs int, gs1, gs2 string, ord int, dsc bool) (models.Rp_count, error)
	Add(ctx context.Context, ea models.Rp) (int, error)
	Upd(ctx context.Context, eu models.Rp) (int, error)
	Del(ctx context.Context, ed []int) ([]int, error)
	GetOne(ctx context.Context, i int) (models.Rp_count, error)
}

//NewRpService(storage pg.RpStorage) *RpService
func NewRpService(storage pgsql.RpStorage) *RpService {
	return &RpService{storage}
}

//func (esv *RpService) GetList(ctx context.Context, pg, pgs int, gs1, gs2 string, ord int, dsc bool) (models.Rp_count, error)
func (esv *RpService) GetList(ctx context.Context, pg, pgs int, gs1, gs2 string, ord int, dsc bool) (models.Rp_count, error) {
	var est ifRpStorage
	est = &esv.storage

	out_count, err := est.GetList(ctx, pg, pgs, gs1, gs2, ord, dsc)

	if err != nil {
		log.Println("RpStorage.GetList", err)
		return models.Rp_count{Values: []models.Rp{}, Count: 0, Auth: models.Auth{}}, err
	}

	return out_count, nil
}

//func (esv *RpService) Add(ctx context.Context, ea models.Rp) (int, error)
func (esv *RpService) Add(ctx context.Context, ea models.Rp) (int, error) {
	var est ifRpStorage
	est = &esv.storage

	ai, err := est.Add(ctx, ea)

	if err != nil {
		log.Println("RpStorage.Add", err)
		return 0, err
	}

	return ai, nil
}

//func (esv *RpService) Upd(ctx context.Context, eu models.Rp) (int, error)
func (esv *RpService) Upd(ctx context.Context, eu models.Rp) (int, error) {
	var est ifRpStorage
	est = &esv.storage

	ui, err := est.Upd(ctx, eu)

	if err != nil {
		log.Println("RpStorage.Upd", err)
		return 0, err
	}

	return ui, nil
}

//func (esv *RpService) Del(ctx context.Context, ed []int) ([]int, error)
func (esv *RpService) Del(ctx context.Context, ed []int) ([]int, error) {
	var est ifRpStorage
	est = &esv.storage

	res, err := est.Del(ctx, ed)

	if err != nil {
		log.Println("RpStorage.Del", err)
		return []int{}, err
	}

	return res, nil
}

//func (esv *RpService) GetOne(ctx context.Context, i int) (models.Rp_count, error)
func (esv *RpService) GetOne(ctx context.Context, i int) (models.Rp_count, error) {
	var est ifRpStorage
	est = &esv.storage

	out_count, err := est.GetOne(ctx, i)

	if err != nil {
		log.Println("RpStorage.GetOne", err)
		return models.Rp_count{Values: []models.Rp{}, Count: 0, Auth: models.Auth{}}, err
	}

	return out_count, nil
}
