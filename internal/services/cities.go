package services

import (
	"context"
	"log"

	"github.com/u2lentaru/billbck/internal/adapters/db/pgsql"
	"github.com/u2lentaru/billbck/internal/models"
)

type CityService struct {
	storage pgsql.CityStorage
}

type ifCityStorage interface {
	GetList(ctx context.Context, pg, pgs int, gs1 string, ord int, dsc bool) (models.City_count, error)
	Add(ctx context.Context, ea models.City) (int, error)
	Upd(ctx context.Context, eu models.City) (int, error)
	Del(ctx context.Context, ed []int) ([]int, error)
	GetOne(ctx context.Context, i int) (models.City_count, error)
}

//func NewCityService(storage pgsql.CityStorage) *CityService
func NewCityService(storage pgsql.CityStorage) *CityService {
	return &CityService{storage}
}

//func (esv *CityService) GetList(ctx context.Context, pg, pgs int, gs1 string, ord int, dsc bool) (models.City_count, error)
func (esv *CityService) GetList(ctx context.Context, pg, pgs int, gs1 string, ord int, dsc bool) (models.City_count, error) {
	var est ifCityStorage
	est = &esv.storage

	out_count, err := est.GetList(ctx, pg, pgs, gs1, ord, dsc)

	if err != nil {
		log.Println("CityStorage.GetList", err)
		return models.City_count{Values: []models.City{}, Count: 0, Auth: models.Auth{}}, err
	}

	return out_count, nil
}

//func (esv *CityService) Add(ctx context.Context, ea models.City) (int, error)
func (esv *CityService) Add(ctx context.Context, ea models.City) (int, error) {
	var est ifCityStorage
	est = &esv.storage

	ai, err := est.Add(ctx, ea)

	if err != nil {
		log.Println("CityStorage.Add", err)
		return 0, err
	}

	return ai, nil
}

//func (esv *CityService) Upd(ctx context.Context, eu models.City) (int, error)
func (esv *CityService) Upd(ctx context.Context, eu models.City) (int, error) {
	var est ifCityStorage
	est = &esv.storage

	ui, err := est.Upd(ctx, eu)

	if err != nil {
		log.Println("CityStorage.Upd", err)
		return 0, err
	}

	return ui, nil
}

//func (esv *CityService) Del(ctx context.Context, ed []int) ([]int, error) {
func (esv *CityService) Del(ctx context.Context, ed []int) ([]int, error) {
	var est ifCityStorage
	est = &esv.storage

	res, err := est.Del(ctx, ed)

	if err != nil {
		log.Println("CityStorage.Del", err)
		return []int{}, err
	}

	return res, nil
}

//func (esv *CityService) GetOne(ctx context.Context, i int) (models.City_count, error)
func (esv *CityService) GetOne(ctx context.Context, i int) (models.City_count, error) {
	var est ifCityStorage
	est = &esv.storage

	out_count, err := est.GetOne(ctx, i)

	if err != nil {
		log.Println("CityStorage.GetOne", err)
		return models.City_count{Values: []models.City{}, Count: 0, Auth: models.Auth{}}, err
	}

	return out_count, nil
}
