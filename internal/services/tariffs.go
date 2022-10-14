package services

import (
	"context"
	"log"

	"github.com/u2lentaru/billbck/internal/adapters/db/pgsql"
	"github.com/u2lentaru/billbck/internal/models"
)

type TariffService struct {
	storage pgsql.TariffStorage
}

type ifTariffStorage interface {
	GetList(ctx context.Context, pg, pgs int, gs1 string, ord int, dsc bool) (models.Tariff_count, error)
	Add(ctx context.Context, ea models.Tariff) (int, error)
	Upd(ctx context.Context, eu models.Tariff) (int, error)
	Del(ctx context.Context, ed []int) ([]int, error)
	GetOne(ctx context.Context, i int) (models.Tariff_count, error)
}

//func NewTariffService(storage pgsql.Tarifftorage) *TariffService
func NewTariffService(storage pgsql.TariffStorage) *TariffService {
	return &TariffService{storage}
}

//func (esv *TariffService) GetList(ctx context.Context, pg, pgs int, gs1 string, ord int, dsc bool) (models.Tariff_count, error)
func (esv *TariffService) GetList(ctx context.Context, pg, pgs int, gs1 string, ord int, dsc bool) (models.Tariff_count, error) {
	var est ifTariffStorage
	est = &esv.storage

	out_count, err := est.GetList(ctx, pg, pgs, gs1, ord, dsc)

	if err != nil {
		log.Println("TariffStorage.GetList", err)
		return models.Tariff_count{Values: []models.Tariff{}, Count: 0, Auth: models.Auth{}}, err
	}

	return out_count, nil
}

//func (esv *TariffService) Add(ctx context.Context, ea models.Tariff) (int, error)
func (esv *TariffService) Add(ctx context.Context, ea models.Tariff) (int, error) {
	var est ifTariffStorage
	est = &esv.storage

	ai, err := est.Add(ctx, ea)

	if err != nil {
		log.Println("TariffStorage.Add", err)
		return 0, err
	}

	return ai, nil
}

//func (esv *TariffService) Upd(ctx context.Context, eu models.Tariff) (int, error)
func (esv *TariffService) Upd(ctx context.Context, eu models.Tariff) (int, error) {
	var est ifTariffStorage
	est = &esv.storage

	ui, err := est.Upd(ctx, eu)

	if err != nil {
		log.Println("TariffStorage.Upd", err)
		return 0, err
	}

	return ui, nil
}

//func (esv *TariffService) Del(ctx context.Context, ed []int) ([]int, error) {
func (esv *TariffService) Del(ctx context.Context, ed []int) ([]int, error) {
	var est ifTariffStorage
	est = &esv.storage

	res, err := est.Del(ctx, ed)

	if err != nil {
		log.Println("TariffStorage.Del", err)
		return []int{}, err
	}

	return res, nil
}

//func (esv *TariffService) GetOne(ctx context.Context, i int) (models.Tariff_count, error)
func (esv *TariffService) GetOne(ctx context.Context, i int) (models.Tariff_count, error) {
	var est ifTariffStorage
	est = &esv.storage

	out_count, err := est.GetOne(ctx, i)

	if err != nil {
		log.Println("TariffStorage.GetOne", err)
		return models.Tariff_count{Values: []models.Tariff{}, Count: 0, Auth: models.Auth{}}, err
	}

	return out_count, nil
}
