package services

import (
	"context"
	"log"

	"github.com/u2lentaru/billbck/internal/adapters/db/pgsql"
	"github.com/u2lentaru/billbck/internal/models"
)

type HouseService struct {
	storage pgsql.HouseStorage
}

type ifHouseStorage interface {
	GetList(ctx context.Context, pg, pgs int, gs1, gs2 string, gs3, ord int, dsc bool) (models.House_count, error)
	Add(ctx context.Context, ea models.House) (int, error)
	Upd(ctx context.Context, eu models.House) (int, error)
	Del(ctx context.Context, ed []int) ([]int, error)
	GetOne(ctx context.Context, i int) (models.House_count, error)
}

//NewHouseService(storage pg.HouseStorage) *HouseService
func NewHouseService(storage pgsql.HouseStorage) *HouseService {
	return &HouseService{storage}
}

//func (esv *HouseService) GetList(ctx context.Context, pg, pgs int, gs1, gs2 string, gs3, ord int, dsc bool) (models.House_count, error)
func (esv *HouseService) GetList(ctx context.Context, pg, pgs int, gs1, gs2 string, gs3, ord int, dsc bool) (models.House_count, error) {
	var est ifHouseStorage
	est = &esv.storage

	out_count, err := est.GetList(ctx, pg, pgs, gs1, gs2, gs3, ord, dsc)

	if err != nil {
		log.Println("HouseStorage.GetList", err)
		return models.House_count{Values: []models.House{}, Count: 0, Auth: models.Auth{}}, err
	}

	return out_count, nil
}

//func (esv *HouseService) Add(ctx context.Context, ea models.House) (int, error)
func (esv *HouseService) Add(ctx context.Context, ea models.House) (int, error) {
	var est ifHouseStorage
	est = &esv.storage

	ai, err := est.Add(ctx, ea)

	if err != nil {
		log.Println("HouseStorage.Add", err)
		return 0, err
	}

	return ai, nil
}

//func (esv *HouseService) Upd(ctx context.Context, eu models.House) (int, error)
func (esv *HouseService) Upd(ctx context.Context, eu models.House) (int, error) {
	var est ifHouseStorage
	est = &esv.storage

	ui, err := est.Upd(ctx, eu)

	if err != nil {
		log.Println("HouseStorage.Upd", err)
		return 0, err
	}

	return ui, nil
}

//func (esv *HouseService) Del(ctx context.Context, ed []int) ([]int, error)
func (esv *HouseService) Del(ctx context.Context, ed []int) ([]int, error) {
	var est ifHouseStorage
	est = &esv.storage

	res, err := est.Del(ctx, ed)

	if err != nil {
		log.Println("HouseStorage.Del", err)
		return []int{}, err
	}

	return res, nil
}

//func (esv *HouseService) GetOne(ctx context.Context, i int) (models.House_count, error)
func (esv *HouseService) GetOne(ctx context.Context, i int) (models.House_count, error) {
	var est ifHouseStorage
	est = &esv.storage

	out_count, err := est.GetOne(ctx, i)

	if err != nil {
		log.Println("HouseStorage.GetOne", err)
		return models.House_count{Values: []models.House{}, Count: 0, Auth: models.Auth{}}, err
	}

	return out_count, nil
}
