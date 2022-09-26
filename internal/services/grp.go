package services

import (
	"context"
	"log"

	"github.com/u2lentaru/billbck/internal/adapters/db/pgsql"
	"github.com/u2lentaru/billbck/internal/models"
)

type GRpService struct {
	storage pgsql.GRpStorage
}

type ifGRpStorage interface {
	GetList(ctx context.Context, pg, pgs int, gs1 string, ord int, dsc bool) (models.GRp_count, error)
	Add(ctx context.Context, ea models.GRp) (int, error)
	Upd(ctx context.Context, eu models.GRp) (int, error)
	Del(ctx context.Context, ed []int) ([]int, error)
	GetOne(ctx context.Context, i int) (models.GRp_count, error)
}

//func NewGRpService(storage pgsql.GRpStorage) *GRpService
func NewGRpService(storage pgsql.GRpStorage) *GRpService {
	return &GRpService{storage}
}

//func (esv *GRpService) GetList(ctx context.Context, pg, pgs int, gs1 string, ord int, dsc bool) (models.GRp_count, error)
func (esv *GRpService) GetList(ctx context.Context, pg, pgs int, gs1 string, ord int, dsc bool) (models.GRp_count, error) {
	var est ifGRpStorage
	est = &esv.storage

	out_count, err := est.GetList(ctx, pg, pgs, gs1, ord, dsc)

	if err != nil {
		log.Println("GRpStorage.GetList", err)
		return models.GRp_count{Values: []models.GRp{}, Count: 0, Auth: models.Auth{}}, err
	}

	return out_count, nil
}

//func (esv *GRpService) Add(ctx context.Context, ea models.GRp) (int, error)
func (esv *GRpService) Add(ctx context.Context, ea models.GRp) (int, error) {
	var est ifGRpStorage
	est = &esv.storage

	ai, err := est.Add(ctx, ea)

	if err != nil {
		log.Println("GRpStorage.Add", err)
		return 0, err
	}

	return ai, nil
}

//func (esv *GRpService) Upd(ctx context.Context, eu models.GRp) (int, error)
func (esv *GRpService) Upd(ctx context.Context, eu models.GRp) (int, error) {
	var est ifGRpStorage
	est = &esv.storage

	ui, err := est.Upd(ctx, eu)

	if err != nil {
		log.Println("GRpStorage.Upd", err)
		return 0, err
	}

	return ui, nil
}

//func (esv *GRpService) Del(ctx context.Context, ed []int) ([]int, error) {
func (esv *GRpService) Del(ctx context.Context, ed []int) ([]int, error) {
	var est ifGRpStorage
	est = &esv.storage

	res, err := est.Del(ctx, ed)

	if err != nil {
		log.Println("GRpStorage.Del", err)
		return []int{}, err
	}

	return res, nil
}

//func (esv *GRpService) GetOne(ctx context.Context, i int) (models.GRp_count, error)
func (esv *GRpService) GetOne(ctx context.Context, i int) (models.GRp_count, error) {
	var est ifGRpStorage
	est = &esv.storage

	out_count, err := est.GetOne(ctx, i)

	if err != nil {
		log.Println("GRpStorage.GetOne", err)
		return models.GRp_count{Values: []models.GRp{}, Count: 0, Auth: models.Auth{}}, err
	}

	return out_count, nil
}
