package services

import (
	"context"
	"log"

	"github.com/u2lentaru/billbck/internal/adapters/db/pgsql"
	"github.com/u2lentaru/billbck/internal/models"
)

type AreaService struct {
	storage pgsql.AreaStorage
}

type ifAreaStorage interface {
	GetList(ctx context.Context, pg, pgs int, gs1, gs2 string, ord int, dsc bool) (models.Area_count, error)
	Add(ctx context.Context, ea models.Area) (int, error)
	Upd(ctx context.Context, eu models.Area) (int, error)
	Del(ctx context.Context, ed []int) ([]int, error)
	GetOne(ctx context.Context, i int) (models.Area_count, error)
}

//NewAreaService(storage pg.AreaStorage) *AreaService
func NewAreaService(storage pgsql.AreaStorage) *AreaService {
	return &AreaService{storage}
}

//func (esv *AreaService) GetList(ctx context.Context, pg, pgs int, gs1, gs2 string, ord int, dsc bool) (models.Area_count, error)
func (esv *AreaService) GetList(ctx context.Context, pg, pgs int, gs1, gs2 string, ord int, dsc bool) (models.Area_count, error) {
	var est ifAreaStorage
	est = &esv.storage

	out_count, err := est.GetList(ctx, pg, pgs, gs1, gs2, ord, dsc)

	if err != nil {
		log.Println("AreaStorage.GetList", err)
		return models.Area_count{Values: []models.Area{}, Count: 0, Auth: models.Auth{}}, err
	}

	return out_count, nil
}

//func (esv *AreaService) Add(ctx context.Context, ea models.Area) (int, error)
func (esv *AreaService) Add(ctx context.Context, ea models.Area) (int, error) {
	var est ifAreaStorage
	est = &esv.storage

	ai, err := est.Add(ctx, ea)

	if err != nil {
		log.Println("AreaStorage.Add", err)
		return 0, err
	}

	return ai, nil
}

//func (esv *AreaService) Upd(ctx context.Context, eu models.Area) (int, error)
func (esv *AreaService) Upd(ctx context.Context, eu models.Area) (int, error) {
	var est ifAreaStorage
	est = &esv.storage

	ui, err := est.Upd(ctx, eu)

	if err != nil {
		log.Println("AreaStorage.Upd", err)
		return 0, err
	}

	return ui, nil
}

//func (esv *AreaService) Del(ctx context.Context, ed []int) ([]int, error)
func (esv *AreaService) Del(ctx context.Context, ed []int) ([]int, error) {
	var est ifAreaStorage
	est = &esv.storage

	res, err := est.Del(ctx, ed)

	if err != nil {
		log.Println("AreaStorage.Del", err)
		return []int{}, err
	}

	return res, nil
}

//func (esv *AreaService) GetOne(ctx context.Context, i int) (models.Area_count, error)
func (esv *AreaService) GetOne(ctx context.Context, i int) (models.Area_count, error) {
	var est ifAreaStorage
	est = &esv.storage

	out_count, err := est.GetOne(ctx, i)

	if err != nil {
		log.Println("AreaStorage.GetOne", err)
		return models.Area_count{Values: []models.Area{}, Count: 0, Auth: models.Auth{}}, err
	}

	return out_count, nil

}
