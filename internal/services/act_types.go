package services

import (
	"context"
	"log"

	"github.com/u2lentaru/billbck/internal/adapters/db/pgsql"
	"github.com/u2lentaru/billbck/internal/models"
)

type ActTypeService struct {
	storage pgsql.ActTypeStorage
}

type ifActTypeStorage interface {
	GetList(ctx context.Context, pg, pgs int, nm string, ord int, dsc bool) (models.ActType_count, error)
	Add(ctx context.Context, ea models.ActType) (int, error)
	Upd(ctx context.Context, eu models.ActType) (int, error)
	Del(ctx context.Context, ed []int) ([]int, error)
	GetOne(ctx context.Context, i int) (models.ActType_count, error)
}

//NewActTypeService(storage pg.ActTypeStorage) *ActTypeService
func NewActTypeService(storage pgsql.ActTypeStorage) *ActTypeService {
	return &ActTypeService{storage}
}

//func (esv *ActTypeService) GetList(ctx context.Context, pg, pgs int, nm string, ord int, dsc bool) (models.ActType_count, error)
func (esv *ActTypeService) GetList(ctx context.Context, pg, pgs int, nm string, ord int, dsc bool) (models.ActType_count, error) {
	var est ifActTypeStorage
	est = &esv.storage
	// est = pgsql.NewActTypeStorage(nil)

	out_count, err := est.GetList(ctx, pg, pgs, nm, ord, dsc)

	if err != nil {
		log.Println("ActTypeStorage.GetList", err)
		return models.ActType_count{Values: []models.ActType{}, Count: 0, Auth: models.Auth{}}, err
	}

	return out_count, nil
}

//func (esv *ActTypeService) Add(ctx context.Context, ea models.ActType) (int, error)
func (esv *ActTypeService) Add(ctx context.Context, ea models.ActType) (int, error) {
	var est ifActTypeStorage
	est = &esv.storage
	// est = pgsql.NewActTypeStorage(nil)

	ai, err := est.Add(ctx, ea)

	if err != nil {
		log.Println("ActTypeStorage.Add", err)
		return 0, err
	}

	return ai, nil
}

//func (esv *ActTypeService) Upd(ctx context.Context, eu models.ActType)
func (esv *ActTypeService) Upd(ctx context.Context, eu models.ActType) (int, error) {
	var est ifActTypeStorage
	est = &esv.storage
	// est = pgsql.NewActTypeStorage(nil)

	ui, err := est.Upd(ctx, eu)

	if err != nil {
		log.Println("ActTypeStorage.Upd", err)
		return 0, err
	}

	return ui, nil
}

//func (esv *ActTypeService) Del(ctx context.Context, ed []int) ([]int, error)
func (esv *ActTypeService) Del(ctx context.Context, ed []int) ([]int, error) {
	var est ifActTypeStorage
	est = &esv.storage
	// est = pgsql.NewActTypeStorage(nil)

	res, err := est.Del(ctx, ed)

	if err != nil {
		log.Println("ActTypeStorage.Del", err)
		return []int{}, err
	}

	return res, nil
}

//func (esv *ActTypeService) GetOne(ctx context.Context, i int) (models.ActType_count, error)
func (esv *ActTypeService) GetOne(ctx context.Context, i int) (models.ActType_count, error) {
	var est ifActTypeStorage
	est = &esv.storage
	// est = pgsql.NewActTypeStorage(nil)

	out_count, err := est.GetOne(ctx, i)

	if err != nil {
		log.Println("ActTypeStorage.GetOne", err)
		return models.ActType_count{Values: []models.ActType{}, Count: 0, Auth: models.Auth{}}, err
	}

	return out_count, nil

}
