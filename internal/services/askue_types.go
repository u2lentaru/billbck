package services

import (
	"context"
	"log"

	"github.com/u2lentaru/billbck/internal/adapters/db/pgsql"
	"github.com/u2lentaru/billbck/internal/models"
)

type AskueTypeService struct {
	storage pgsql.AskueTypeStorage
}

type ifAskueTypeStorage interface {
	GetList(ctx context.Context, pg, pgs int, gs1 string, ord int, dsc bool) (models.AskueType_count, error)
	Add(ctx context.Context, ea models.AskueType) (int, error)
	Upd(ctx context.Context, eu models.AskueType) (int, error)
	Del(ctx context.Context, ed []int) ([]int, error)
	GetOne(ctx context.Context, i int) (models.AskueType_count, error)
}

//func NewAskueTypeService(storage pgsql.AskueTypeStorage) *AskueTypeService
func NewAskueTypeService(storage pgsql.AskueTypeStorage) *AskueTypeService {
	return &AskueTypeService{storage}
}

//func (esv *AskueTypeService) GetList(ctx context.Context, pg, pgs int, gs1 string, ord int, dsc bool) (models.AskueType_count, error)
func (esv *AskueTypeService) GetList(ctx context.Context, pg, pgs int, gs1 string, ord int, dsc bool) (models.AskueType_count, error) {
	var est ifAskueTypeStorage
	est = &esv.storage

	out_count, err := est.GetList(ctx, pg, pgs, gs1, ord, dsc)

	if err != nil {
		log.Println("AskueTypeStorage.GetList", err)
		return models.AskueType_count{Values: []models.AskueType{}, Count: 0, Auth: models.Auth{}}, err
	}

	return out_count, nil
}

//func (esv *AskueTypeService) Add(ctx context.Context, ea models.AskueType) (int, error)
func (esv *AskueTypeService) Add(ctx context.Context, ea models.AskueType) (int, error) {
	var est ifAskueTypeStorage
	est = &esv.storage

	ai, err := est.Add(ctx, ea)

	if err != nil {
		log.Println("AskueTypeStorage.Add", err)
		return 0, err
	}

	return ai, nil
}

//func (esv *AskueTypeService) Upd(ctx context.Context, eu models.AskueType) (int, error)
func (esv *AskueTypeService) Upd(ctx context.Context, eu models.AskueType) (int, error) {
	var est ifAskueTypeStorage
	est = &esv.storage

	ui, err := est.Upd(ctx, eu)

	if err != nil {
		log.Println("AskueTypeStorage.Upd", err)
		return 0, err
	}

	return ui, nil
}

//func (esv *AskueTypeService) Del(ctx context.Context, ed []int) ([]int, error) {
func (esv *AskueTypeService) Del(ctx context.Context, ed []int) ([]int, error) {
	var est ifAskueTypeStorage
	est = &esv.storage

	res, err := est.Del(ctx, ed)

	if err != nil {
		log.Println("AskueTypeStorage.Del", err)
		return []int{}, err
	}

	return res, nil
}

//func (esv *AskueTypeService) GetOne(ctx context.Context, i int) (models.AskueType_count, error)
func (esv *AskueTypeService) GetOne(ctx context.Context, i int) (models.AskueType_count, error) {
	var est ifAskueTypeStorage
	est = &esv.storage

	out_count, err := est.GetOne(ctx, i)

	if err != nil {
		log.Println("AskueTypeStorage.GetOne", err)
		return models.AskueType_count{Values: []models.AskueType{}, Count: 0, Auth: models.Auth{}}, err
	}

	return out_count, nil
}
