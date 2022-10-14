package services

import (
	"context"
	"log"

	"github.com/u2lentaru/billbck/internal/adapters/db/pgsql"
	"github.com/u2lentaru/billbck/internal/models"
)

type TariffGroupService struct {
	storage pgsql.TariffGroupStorage
}

type ifTariffGroupStorage interface {
	GetList(ctx context.Context, pg, pgs int, gs1 string, ord int, dsc bool) (models.TariffGroup_count, error)
	Add(ctx context.Context, ea models.TariffGroup) (int, error)
	Upd(ctx context.Context, eu models.TariffGroup) (int, error)
	Del(ctx context.Context, ed []int) ([]int, error)
	GetOne(ctx context.Context, i int) (models.TariffGroup_count, error)
}

//func NewTariffGroupService(storage pgsql.TariffGrouptorage) *TariffGroupService
func NewTariffGroupService(storage pgsql.TariffGroupStorage) *TariffGroupService {
	return &TariffGroupService{storage}
}

//func (esv *TariffGroupService) GetList(ctx context.Context, pg, pgs int, gs1 string, ord int, dsc bool) (models.TariffGroup_count, error)
func (esv *TariffGroupService) GetList(ctx context.Context, pg, pgs int, gs1 string, ord int, dsc bool) (models.TariffGroup_count, error) {
	var est ifTariffGroupStorage
	est = &esv.storage

	out_count, err := est.GetList(ctx, pg, pgs, gs1, ord, dsc)

	if err != nil {
		log.Println("TariffGroupStorage.GetList", err)
		return models.TariffGroup_count{Values: []models.TariffGroup{}, Count: 0, Auth: models.Auth{}}, err
	}

	return out_count, nil
}

//func (esv *TariffGroupService) Add(ctx context.Context, ea models.TariffGroup) (int, error)
func (esv *TariffGroupService) Add(ctx context.Context, ea models.TariffGroup) (int, error) {
	var est ifTariffGroupStorage
	est = &esv.storage

	ai, err := est.Add(ctx, ea)

	if err != nil {
		log.Println("TariffGroupStorage.Add", err)
		return 0, err
	}

	return ai, nil
}

//func (esv *TariffGroupService) Upd(ctx context.Context, eu models.TariffGroup) (int, error)
func (esv *TariffGroupService) Upd(ctx context.Context, eu models.TariffGroup) (int, error) {
	var est ifTariffGroupStorage
	est = &esv.storage

	ui, err := est.Upd(ctx, eu)

	if err != nil {
		log.Println("TariffGroupStorage.Upd", err)
		return 0, err
	}

	return ui, nil
}

//func (esv *TariffGroupService) Del(ctx context.Context, ed []int) ([]int, error) {
func (esv *TariffGroupService) Del(ctx context.Context, ed []int) ([]int, error) {
	var est ifTariffGroupStorage
	est = &esv.storage

	res, err := est.Del(ctx, ed)

	if err != nil {
		log.Println("TariffGroupStorage.Del", err)
		return []int{}, err
	}

	return res, nil
}

//func (esv *TariffGroupService) GetOne(ctx context.Context, i int) (models.TariffGroup_count, error)
func (esv *TariffGroupService) GetOne(ctx context.Context, i int) (models.TariffGroup_count, error) {
	var est ifTariffGroupStorage
	est = &esv.storage

	out_count, err := est.GetOne(ctx, i)

	if err != nil {
		log.Println("TariffGroupStorage.GetOne", err)
		return models.TariffGroup_count{Values: []models.TariffGroup{}, Count: 0, Auth: models.Auth{}}, err
	}

	return out_count, nil
}
