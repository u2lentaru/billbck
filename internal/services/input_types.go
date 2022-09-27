package services

import (
	"context"
	"log"

	"github.com/u2lentaru/billbck/internal/adapters/db/pgsql"
	"github.com/u2lentaru/billbck/internal/models"
)

type InputTypeService struct {
	storage pgsql.InputTypeStorage
}

type ifInputTypeStorage interface {
	GetList(ctx context.Context, pg, pgs int, gs1 string, ord int, dsc bool) (models.InputType_count, error)
	Add(ctx context.Context, ea models.InputType) (int, error)
	Upd(ctx context.Context, eu models.InputType) (int, error)
	Del(ctx context.Context, ed []int) ([]int, error)
	GetOne(ctx context.Context, i int) (models.InputType_count, error)
}

//func NewInputTypeService(storage pgsql.InputTypeStorage) *InputTypeService
func NewInputTypeService(storage pgsql.InputTypeStorage) *InputTypeService {
	return &InputTypeService{storage}
}

//func (esv *InputTypeService) GetList(ctx context.Context, pg, pgs int, gs1 string, ord int, dsc bool) (models.InputType_count, error)
func (esv *InputTypeService) GetList(ctx context.Context, pg, pgs int, gs1 string, ord int, dsc bool) (models.InputType_count, error) {
	var est ifInputTypeStorage
	est = &esv.storage

	out_count, err := est.GetList(ctx, pg, pgs, gs1, ord, dsc)

	if err != nil {
		log.Println("InputTypeStorage.GetList", err)
		return models.InputType_count{Values: []models.InputType{}, Count: 0, Auth: models.Auth{}}, err
	}

	return out_count, nil
}

//func (esv *InputTypeService) Add(ctx context.Context, ea models.InputType) (int, error)
func (esv *InputTypeService) Add(ctx context.Context, ea models.InputType) (int, error) {
	var est ifInputTypeStorage
	est = &esv.storage

	ai, err := est.Add(ctx, ea)

	if err != nil {
		log.Println("InputTypeStorage.Add", err)
		return 0, err
	}

	return ai, nil
}

//func (esv *InputTypeService) Upd(ctx context.Context, eu models.InputType) (int, error)
func (esv *InputTypeService) Upd(ctx context.Context, eu models.InputType) (int, error) {
	var est ifInputTypeStorage
	est = &esv.storage

	ui, err := est.Upd(ctx, eu)

	if err != nil {
		log.Println("InputTypeStorage.Upd", err)
		return 0, err
	}

	return ui, nil
}

//func (esv *InputTypeService) Del(ctx context.Context, ed []int) ([]int, error) {
func (esv *InputTypeService) Del(ctx context.Context, ed []int) ([]int, error) {
	var est ifInputTypeStorage
	est = &esv.storage

	res, err := est.Del(ctx, ed)

	if err != nil {
		log.Println("InputTypeStorage.Del", err)
		return []int{}, err
	}

	return res, nil
}

//func (esv *InputTypeService) GetOne(ctx context.Context, i int) (models.InputType_count, error)
func (esv *InputTypeService) GetOne(ctx context.Context, i int) (models.InputType_count, error) {
	var est ifInputTypeStorage
	est = &esv.storage

	out_count, err := est.GetOne(ctx, i)

	if err != nil {
		log.Println("InputTypeStorage.GetOne", err)
		return models.InputType_count{Values: []models.InputType{}, Count: 0, Auth: models.Auth{}}, err
	}

	return out_count, nil
}
