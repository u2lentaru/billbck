package services

import (
	"context"
	"log"

	"github.com/u2lentaru/billbck/internal/adapters/db/pgsql"
	"github.com/u2lentaru/billbck/internal/models"
)

type PositionService struct {
	storage pgsql.PositionStorage
}

type ifPositionStorage interface {
	GetList(ctx context.Context, pg, pgs int, gs1 string, ord int, dsc bool) (models.Position_count, error)
	Add(ctx context.Context, ea models.Position) (int, error)
	Upd(ctx context.Context, eu models.Position) (int, error)
	Del(ctx context.Context, ed []int) ([]int, error)
	GetOne(ctx context.Context, i int) (models.Position_count, error)
}

//func NewPositionService(storage pgsql.Positiontorage) *PositionService
func NewPositionService(storage pgsql.PositionStorage) *PositionService {
	return &PositionService{storage}
}

//func (esv *PositionService) GetList(ctx context.Context, pg, pgs int, gs1 string, ord int, dsc bool) (models.Position_count, error)
func (esv *PositionService) GetList(ctx context.Context, pg, pgs int, gs1 string, ord int, dsc bool) (models.Position_count, error) {
	var est ifPositionStorage
	est = &esv.storage

	out_count, err := est.GetList(ctx, pg, pgs, gs1, ord, dsc)

	if err != nil {
		log.Println("PositionStorage.GetList", err)
		return models.Position_count{Values: []models.Position{}, Count: 0, Auth: models.Auth{}}, err
	}

	return out_count, nil
}

//func (esv *PositionService) Add(ctx context.Context, ea models.Position) (int, error)
func (esv *PositionService) Add(ctx context.Context, ea models.Position) (int, error) {
	var est ifPositionStorage
	est = &esv.storage

	ai, err := est.Add(ctx, ea)

	if err != nil {
		log.Println("PositionStorage.Add", err)
		return 0, err
	}

	return ai, nil
}

//func (esv *PositionService) Upd(ctx context.Context, eu models.Position) (int, error)
func (esv *PositionService) Upd(ctx context.Context, eu models.Position) (int, error) {
	var est ifPositionStorage
	est = &esv.storage

	ui, err := est.Upd(ctx, eu)

	if err != nil {
		log.Println("PositionStorage.Upd", err)
		return 0, err
	}

	return ui, nil
}

//func (esv *PositionService) Del(ctx context.Context, ed []int) ([]int, error) {
func (esv *PositionService) Del(ctx context.Context, ed []int) ([]int, error) {
	var est ifPositionStorage
	est = &esv.storage

	res, err := est.Del(ctx, ed)

	if err != nil {
		log.Println("PositionStorage.Del", err)
		return []int{}, err
	}

	return res, nil
}

//func (esv *PositionService) GetOne(ctx context.Context, i int) (models.Position_count, error)
func (esv *PositionService) GetOne(ctx context.Context, i int) (models.Position_count, error) {
	var est ifPositionStorage
	est = &esv.storage

	out_count, err := est.GetOne(ctx, i)

	if err != nil {
		log.Println("PositionStorage.GetOne", err)
		return models.Position_count{Values: []models.Position{}, Count: 0, Auth: models.Auth{}}, err
	}

	return out_count, nil
}
