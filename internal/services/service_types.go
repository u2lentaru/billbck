package services

import (
	"context"
	"log"

	"github.com/u2lentaru/billbck/internal/adapters/db/pgsql"
	"github.com/u2lentaru/billbck/internal/models"
)

type ServiceTypeService struct {
	storage pgsql.ServiceTypeStorage
}

type ifServiceTypeStorage interface {
	GetList(ctx context.Context, pg, pgs int, gs1 string, ord int, dsc bool) (models.ServiceType_count, error)
	Add(ctx context.Context, ea models.ServiceType) (int, error)
	Upd(ctx context.Context, eu models.ServiceType) (int, error)
	Del(ctx context.Context, ed []int) ([]int, error)
	GetOne(ctx context.Context, i int) (models.ServiceType_count, error)
}

//func NewServiceTypeService(storage pgsql.ServiceTypeStorage) *ServiceTypeService
func NewServiceTypeService(storage pgsql.ServiceTypeStorage) *ServiceTypeService {
	return &ServiceTypeService{storage}
}

//func (esv *ServiceTypeService) GetList(ctx context.Context, pg, pgs int, gs1 string, ord int, dsc bool) (models.ServiceType_count, error)
func (esv *ServiceTypeService) GetList(ctx context.Context, pg, pgs int, gs1 string, ord int, dsc bool) (models.ServiceType_count, error) {
	var est ifServiceTypeStorage
	est = &esv.storage

	out_count, err := est.GetList(ctx, pg, pgs, gs1, ord, dsc)

	if err != nil {
		log.Println("ServiceTypeStorage.GetList", err)
		return models.ServiceType_count{Values: []models.ServiceType{}, Count: 0, Auth: models.Auth{}}, err
	}

	return out_count, nil
}

//func (esv *ServiceTypeService) Add(ctx context.Context, ea models.ServiceType) (int, error)
func (esv *ServiceTypeService) Add(ctx context.Context, ea models.ServiceType) (int, error) {
	var est ifServiceTypeStorage
	est = &esv.storage

	ai, err := est.Add(ctx, ea)

	if err != nil {
		log.Println("ServiceTypeStorage.Add", err)
		return 0, err
	}

	return ai, nil
}

//func (esv *ServiceTypeService) Upd(ctx context.Context, eu models.ServiceType) (int, error)
func (esv *ServiceTypeService) Upd(ctx context.Context, eu models.ServiceType) (int, error) {
	var est ifServiceTypeStorage
	est = &esv.storage

	ui, err := est.Upd(ctx, eu)

	if err != nil {
		log.Println("ServiceTypeStorage.Upd", err)
		return 0, err
	}

	return ui, nil
}

//func (esv *ServiceTypeService) Del(ctx context.Context, ed []int) ([]int, error) {
func (esv *ServiceTypeService) Del(ctx context.Context, ed []int) ([]int, error) {
	var est ifServiceTypeStorage
	est = &esv.storage

	res, err := est.Del(ctx, ed)

	if err != nil {
		log.Println("ServiceTypeStorage.Del", err)
		return []int{}, err
	}

	return res, nil
}

//func (esv *ServiceTypeService) GetOne(ctx context.Context, i int) (models.ServiceType_count, error)
func (esv *ServiceTypeService) GetOne(ctx context.Context, i int) (models.ServiceType_count, error) {
	var est ifServiceTypeStorage
	est = &esv.storage

	out_count, err := est.GetOne(ctx, i)

	if err != nil {
		log.Println("ServiceTypeStorage.GetOne", err)
		return models.ServiceType_count{Values: []models.ServiceType{}, Count: 0, Auth: models.Auth{}}, err
	}

	return out_count, nil
}
