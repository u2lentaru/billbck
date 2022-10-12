package services

import (
	"context"
	"log"

	"github.com/u2lentaru/billbck/internal/adapters/db/pgsql"
	"github.com/u2lentaru/billbck/internal/models"
)

type ShutdownTypeService struct {
	storage pgsql.ShutdownTypeStorage
}

type ifShutdownTypeStorage interface {
	GetList(ctx context.Context, pg, pgs int, gs1 string, ord int, dsc bool) (models.ShutdownType_count, error)
	Add(ctx context.Context, ea models.ShutdownType) (int, error)
	Upd(ctx context.Context, eu models.ShutdownType) (int, error)
	Del(ctx context.Context, ed []int) ([]int, error)
	GetOne(ctx context.Context, i int) (models.ShutdownType_count, error)
}

//func NewShutdownTypeService(storage pgsql.ShutdownTypetorage) *ShutdownTypeService
func NewShutdownTypeService(storage pgsql.ShutdownTypeStorage) *ShutdownTypeService {
	return &ShutdownTypeService{storage}
}

//func (esv *ShutdownTypeService) GetList(ctx context.Context, pg, pgs int, gs1 string, ord int, dsc bool) (models.ShutdownType_count, error)
func (esv *ShutdownTypeService) GetList(ctx context.Context, pg, pgs int, gs1 string, ord int, dsc bool) (models.ShutdownType_count, error) {
	var est ifShutdownTypeStorage
	est = &esv.storage

	out_count, err := est.GetList(ctx, pg, pgs, gs1, ord, dsc)

	if err != nil {
		log.Println("ShutdownTypeStorage.GetList", err)
		return models.ShutdownType_count{Values: []models.ShutdownType{}, Count: 0, Auth: models.Auth{}}, err
	}

	return out_count, nil
}

//func (esv *ShutdownTypeService) Add(ctx context.Context, ea models.ShutdownType) (int, error)
func (esv *ShutdownTypeService) Add(ctx context.Context, ea models.ShutdownType) (int, error) {
	var est ifShutdownTypeStorage
	est = &esv.storage

	ai, err := est.Add(ctx, ea)

	if err != nil {
		log.Println("ShutdownTypeStorage.Add", err)
		return 0, err
	}

	return ai, nil
}

//func (esv *ShutdownTypeService) Upd(ctx context.Context, eu models.ShutdownType) (int, error)
func (esv *ShutdownTypeService) Upd(ctx context.Context, eu models.ShutdownType) (int, error) {
	var est ifShutdownTypeStorage
	est = &esv.storage

	ui, err := est.Upd(ctx, eu)

	if err != nil {
		log.Println("ShutdownTypeStorage.Upd", err)
		return 0, err
	}

	return ui, nil
}

//func (esv *ShutdownTypeService) Del(ctx context.Context, ed []int) ([]int, error) {
func (esv *ShutdownTypeService) Del(ctx context.Context, ed []int) ([]int, error) {
	var est ifShutdownTypeStorage
	est = &esv.storage

	res, err := est.Del(ctx, ed)

	if err != nil {
		log.Println("ShutdownTypeStorage.Del", err)
		return []int{}, err
	}

	return res, nil
}

//func (esv *ShutdownTypeService) GetOne(ctx context.Context, i int) (models.ShutdownType_count, error)
func (esv *ShutdownTypeService) GetOne(ctx context.Context, i int) (models.ShutdownType_count, error) {
	var est ifShutdownTypeStorage
	est = &esv.storage

	out_count, err := est.GetOne(ctx, i)

	if err != nil {
		log.Println("ShutdownTypeStorage.GetOne", err)
		return models.ShutdownType_count{Values: []models.ShutdownType{}, Count: 0, Auth: models.Auth{}}, err
	}

	return out_count, nil
}
