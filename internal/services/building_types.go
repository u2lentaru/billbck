package services

import (
	"context"
	"log"

	"github.com/u2lentaru/billbck/internal/adapters/db/pgsql"
	"github.com/u2lentaru/billbck/internal/models"
)

type BuildingTypeService struct {
	storage pgsql.BuildingTypeStorage
}

type ifBuildingTypeStorage interface {
	GetList(ctx context.Context, pg, pgs int, gs1 string, ord int, dsc bool) (models.BuildingType_count, error)
	Add(ctx context.Context, ea models.BuildingType) (int, error)
	Upd(ctx context.Context, eu models.BuildingType) (int, error)
	Del(ctx context.Context, ed []int) ([]int, error)
	GetOne(ctx context.Context, i int) (models.BuildingType_count, error)
}

//func NewBuildingTypeService(storage pgsql.BuildingTypeStorage) *BuildingTypeService
func NewBuildingTypeService(storage pgsql.BuildingTypeStorage) *BuildingTypeService {
	return &BuildingTypeService{storage}
}

//func (esv *BuildingTypeService) GetList(ctx context.Context, pg, pgs int, gs1 string, ord int, dsc bool) (models.BuildingType_count, error)
func (esv *BuildingTypeService) GetList(ctx context.Context, pg, pgs int, gs1 string, ord int, dsc bool) (models.BuildingType_count, error) {
	var est ifBuildingTypeStorage
	est = &esv.storage

	out_count, err := est.GetList(ctx, pg, pgs, gs1, ord, dsc)

	if err != nil {
		log.Println("BuildingTypeStorage.GetList", err)
		return models.BuildingType_count{Values: []models.BuildingType{}, Count: 0, Auth: models.Auth{}}, err
	}

	return out_count, nil
}

//func (esv *BuildingTypeService) Add(ctx context.Context, ea models.BuildingType) (int, error)
func (esv *BuildingTypeService) Add(ctx context.Context, ea models.BuildingType) (int, error) {
	var est ifBuildingTypeStorage
	est = &esv.storage

	ai, err := est.Add(ctx, ea)

	if err != nil {
		log.Println("BuildingTypeStorage.Add", err)
		return 0, err
	}

	return ai, nil
}

//func (esv *BuildingTypeService) Upd(ctx context.Context, eu models.BuildingType) (int, error)
func (esv *BuildingTypeService) Upd(ctx context.Context, eu models.BuildingType) (int, error) {
	var est ifBuildingTypeStorage
	est = &esv.storage

	ui, err := est.Upd(ctx, eu)

	if err != nil {
		log.Println("BuildingTypeStorage.Upd", err)
		return 0, err
	}

	return ui, nil
}

//func (esv *BuildingTypeService) Del(ctx context.Context, ed []int) ([]int, error) {
func (esv *BuildingTypeService) Del(ctx context.Context, ed []int) ([]int, error) {
	var est ifBuildingTypeStorage
	est = &esv.storage

	res, err := est.Del(ctx, ed)

	if err != nil {
		log.Println("BuildingTypeStorage.Del", err)
		return []int{}, err
	}

	return res, nil
}

//func (esv *BuildingTypeService) GetOne(ctx context.Context, i int) (models.BuildingType_count, error)
func (esv *BuildingTypeService) GetOne(ctx context.Context, i int) (models.BuildingType_count, error) {
	var est ifBuildingTypeStorage
	est = &esv.storage

	out_count, err := est.GetOne(ctx, i)

	if err != nil {
		log.Println("BuildingTypeStorage.GetOne", err)
		return models.BuildingType_count{Values: []models.BuildingType{}, Count: 0, Auth: models.Auth{}}, err
	}

	return out_count, nil
}
