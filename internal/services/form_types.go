package services

import (
	"context"
	"log"

	"github.com/u2lentaru/billbck/internal/adapters/db/pgsql"
	"github.com/u2lentaru/billbck/internal/models"
)

type FormTypeService struct {
	storage pgsql.FormTypeStorage
}

type ifFormTypeStorage interface {
	GetList(ctx context.Context, pg, pgs int, gs1, gs2 string, ord int, dsc bool) (models.FormType_count, error)
	Add(ctx context.Context, ea models.FormType) (int, error)
	Upd(ctx context.Context, eu models.FormType) (int, error)
	Del(ctx context.Context, ed []int) ([]int, error)
	GetOne(ctx context.Context, i int) (models.FormType_count, error)
}

//NewFormTypeService(storage pg.FormTypeStorage) *FormTypeService
func NewFormTypeService(storage pgsql.FormTypeStorage) *FormTypeService {
	return &FormTypeService{storage}
}

//func (esv *FormTypeService) GetList(ctx context.Context, pg, pgs int, gs1, gs2 string, ord int, dsc bool) (models.FormType_count, error)
func (esv *FormTypeService) GetList(ctx context.Context, pg, pgs int, gs1, gs2 string, ord int, dsc bool) (models.FormType_count, error) {
	var est ifFormTypeStorage
	est = &esv.storage

	out_count, err := est.GetList(ctx, pg, pgs, gs1, gs2, ord, dsc)

	if err != nil {
		log.Println("FormTypeStorage.GetList", err)
		return models.FormType_count{Values: []models.FormType{}, Count: 0, Auth: models.Auth{}}, err
	}

	return out_count, nil
}

//func (esv *FormTypeService) Add(ctx context.Context, ea models.FormType) (int, error)
func (esv *FormTypeService) Add(ctx context.Context, ea models.FormType) (int, error) {
	var est ifFormTypeStorage
	est = &esv.storage

	ai, err := est.Add(ctx, ea)

	if err != nil {
		log.Println("FormTypeStorage.Add", err)
		return 0, err
	}

	return ai, nil
}

//func (esv *FormTypeService) Upd(ctx context.Context, eu models.FormType) (int, error)
func (esv *FormTypeService) Upd(ctx context.Context, eu models.FormType) (int, error) {
	var est ifFormTypeStorage
	est = &esv.storage

	ui, err := est.Upd(ctx, eu)

	if err != nil {
		log.Println("FormTypeStorage.Upd", err)
		return 0, err
	}

	return ui, nil
}

//func (esv *FormTypeService) Del(ctx context.Context, ed []int) ([]int, error)
func (esv *FormTypeService) Del(ctx context.Context, ed []int) ([]int, error) {
	var est ifFormTypeStorage
	est = &esv.storage

	res, err := est.Del(ctx, ed)

	if err != nil {
		log.Println("FormTypeStorage.Del", err)
		return []int{}, err
	}

	return res, nil
}

//func (esv *FormTypeService) GetOne(ctx context.Context, i int) (models.FormType_count, error)
func (esv *FormTypeService) GetOne(ctx context.Context, i int) (models.FormType_count, error) {
	var est ifFormTypeStorage
	est = &esv.storage

	out_count, err := est.GetOne(ctx, i)

	if err != nil {
		log.Println("FormTypeStorage.GetOne", err)
		return models.FormType_count{Values: []models.FormType{}, Count: 0, Auth: models.Auth{}}, err
	}

	return out_count, nil
}
