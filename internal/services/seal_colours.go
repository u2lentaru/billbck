package services

import (
	"context"
	"log"

	"github.com/u2lentaru/billbck/internal/adapters/db/pgsql"
	"github.com/u2lentaru/billbck/internal/models"
)

type SealColourService struct {
	storage pgsql.SealColourStorage
}

type ifSealColourStorage interface {
	GetList(ctx context.Context, pg, pgs int, gs1 string, ord int, dsc bool) (models.SealColour_count, error)
	Add(ctx context.Context, ea models.SealColour) (int, error)
	Upd(ctx context.Context, eu models.SealColour) (int, error)
	Del(ctx context.Context, ed []int) ([]int, error)
	GetOne(ctx context.Context, i int) (models.SealColour_count, error)
}

//func NewSealColourService(storage pgsql.SealColourtorage) *SealColourService
func NewSealColourService(storage pgsql.SealColourStorage) *SealColourService {
	return &SealColourService{storage}
}

//func (esv *SealColourService) GetList(ctx context.Context, pg, pgs int, gs1 string, ord int, dsc bool) (models.SealColour_count, error)
func (esv *SealColourService) GetList(ctx context.Context, pg, pgs int, gs1 string, ord int, dsc bool) (models.SealColour_count, error) {
	var est ifSealColourStorage
	est = &esv.storage

	out_count, err := est.GetList(ctx, pg, pgs, gs1, ord, dsc)

	if err != nil {
		log.Println("SealColourStorage.GetList", err)
		return models.SealColour_count{Values: []models.SealColour{}, Count: 0, Auth: models.Auth{}}, err
	}

	return out_count, nil
}

//func (esv *SealColourService) Add(ctx context.Context, ea models.SealColour) (int, error)
func (esv *SealColourService) Add(ctx context.Context, ea models.SealColour) (int, error) {
	var est ifSealColourStorage
	est = &esv.storage

	ai, err := est.Add(ctx, ea)

	if err != nil {
		log.Println("SealColourStorage.Add", err)
		return 0, err
	}

	return ai, nil
}

//func (esv *SealColourService) Upd(ctx context.Context, eu models.SealColour) (int, error)
func (esv *SealColourService) Upd(ctx context.Context, eu models.SealColour) (int, error) {
	var est ifSealColourStorage
	est = &esv.storage

	ui, err := est.Upd(ctx, eu)

	if err != nil {
		log.Println("SealColourStorage.Upd", err)
		return 0, err
	}

	return ui, nil
}

//func (esv *SealColourService) Del(ctx context.Context, ed []int) ([]int, error) {
func (esv *SealColourService) Del(ctx context.Context, ed []int) ([]int, error) {
	var est ifSealColourStorage
	est = &esv.storage

	res, err := est.Del(ctx, ed)

	if err != nil {
		log.Println("SealColourStorage.Del", err)
		return []int{}, err
	}

	return res, nil
}

//func (esv *SealColourService) GetOne(ctx context.Context, i int) (models.SealColour_count, error)
func (esv *SealColourService) GetOne(ctx context.Context, i int) (models.SealColour_count, error) {
	var est ifSealColourStorage
	est = &esv.storage

	out_count, err := est.GetOne(ctx, i)

	if err != nil {
		log.Println("SealColourStorage.GetOne", err)
		return models.SealColour_count{Values: []models.SealColour{}, Count: 0, Auth: models.Auth{}}, err
	}

	return out_count, nil
}
