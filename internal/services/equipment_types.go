package services

import (
	"context"
	"log"

	"github.com/u2lentaru/billbck/internal/adapters/db/pgsql"
	"github.com/u2lentaru/billbck/internal/models"
)

type EquipmentTypeService struct {
	storage pgsql.EquipmentTypeStorage
}

type ifEquipmentTypeStorage interface {
	GetList(ctx context.Context, pg, pgs int, gs1 string, ord int, dsc bool) (models.EquipmentType_count, error)
	Add(ctx context.Context, ea models.EquipmentType) (int, error)
	Upd(ctx context.Context, eu models.EquipmentType) (int, error)
	Del(ctx context.Context, ed []int) ([]int, error)
	GetOne(ctx context.Context, i int) (models.EquipmentType_count, error)
}

//func NewEquipmentTypeService(storage pgsql.EquipmentTypeStorage) *EquipmentTypeService
func NewEquipmentTypeService(storage pgsql.EquipmentTypeStorage) *EquipmentTypeService {
	return &EquipmentTypeService{storage}
}

//func (esv *EquipmentTypeService) GetList(ctx context.Context, pg, pgs int, gs1 string, ord int, dsc bool) (models.EquipmentType_count, error)
func (esv *EquipmentTypeService) GetList(ctx context.Context, pg, pgs int, gs1 string, ord int, dsc bool) (models.EquipmentType_count, error) {
	var est ifEquipmentTypeStorage
	est = &esv.storage

	out_count, err := est.GetList(ctx, pg, pgs, gs1, ord, dsc)

	if err != nil {
		log.Println("EquipmentTypeStorage.GetList", err)
		return models.EquipmentType_count{Values: []models.EquipmentType{}, Count: 0, Auth: models.Auth{}}, err
	}

	return out_count, nil
}

//func (esv *EquipmentTypeService) Add(ctx context.Context, ea models.EquipmentType) (int, error)
func (esv *EquipmentTypeService) Add(ctx context.Context, ea models.EquipmentType) (int, error) {
	var est ifEquipmentTypeStorage
	est = &esv.storage

	ai, err := est.Add(ctx, ea)

	if err != nil {
		log.Println("EquipmentTypeStorage.Add", err)
		return 0, err
	}

	return ai, nil
}

//func (esv *EquipmentTypeService) Upd(ctx context.Context, eu models.EquipmentType) (int, error)
func (esv *EquipmentTypeService) Upd(ctx context.Context, eu models.EquipmentType) (int, error) {
	var est ifEquipmentTypeStorage
	est = &esv.storage

	ui, err := est.Upd(ctx, eu)

	if err != nil {
		log.Println("EquipmentTypeStorage.Upd", err)
		return 0, err
	}

	return ui, nil
}

//func (esv *EquipmentTypeService) Del(ctx context.Context, ed []int) ([]int, error) {
func (esv *EquipmentTypeService) Del(ctx context.Context, ed []int) ([]int, error) {
	var est ifEquipmentTypeStorage
	est = &esv.storage

	res, err := est.Del(ctx, ed)

	if err != nil {
		log.Println("EquipmentTypeStorage.Del", err)
		return []int{}, err
	}

	return res, nil
}

//func (esv *EquipmentTypeService) GetOne(ctx context.Context, i int) (models.EquipmentType_count, error)
func (esv *EquipmentTypeService) GetOne(ctx context.Context, i int) (models.EquipmentType_count, error) {
	var est ifEquipmentTypeStorage
	est = &esv.storage

	out_count, err := est.GetOne(ctx, i)

	if err != nil {
		log.Println("EquipmentTypeStorage.GetOne", err)
		return models.EquipmentType_count{Values: []models.EquipmentType{}, Count: 0, Auth: models.Auth{}}, err
	}

	return out_count, nil
}
