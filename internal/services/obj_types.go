package services

import (
	"context"
	"log"

	"github.com/u2lentaru/billbck/internal/adapters/db/pgsql"
	"github.com/u2lentaru/billbck/internal/models"
)

type ObjTypeService struct {
	storage pgsql.ObjTypeStorage
}

type ifObjTypeStorage interface {
	GetList(ctx context.Context, pg, pgs int, gs1 string, ord int, dsc bool) (models.ObjType_count, error)
	Add(ctx context.Context, ea models.ObjType) (int, error)
	Upd(ctx context.Context, eu models.ObjType) (int, error)
	Del(ctx context.Context, ed []int) ([]int, error)
	GetOne(ctx context.Context, i int) (models.ObjType_count, error)
}

//func NewObjTypeService(storage pgsql.ObjTypeStorage) *ObjTypeService
func NewObjTypeService(storage pgsql.ObjTypeStorage) *ObjTypeService {
	return &ObjTypeService{storage}
}

//func (esv *ObjTypeService) GetList(ctx context.Context, pg, pgs int, gs1 string, ord int, dsc bool) (models.ObjType_count, error)
func (esv *ObjTypeService) GetList(ctx context.Context, pg, pgs int, gs1 string, ord int, dsc bool) (models.ObjType_count, error) {
	var est ifObjTypeStorage
	est = &esv.storage

	out_count, err := est.GetList(ctx, pg, pgs, gs1, ord, dsc)

	if err != nil {
		log.Println("ObjTypeStorage.GetList", err)
		return models.ObjType_count{Values: []models.ObjType{}, Count: 0, Auth: models.Auth{}}, err
	}

	return out_count, nil
}

//func (esv *ObjTypeService) Add(ctx context.Context, ea models.ObjType) (int, error)
func (esv *ObjTypeService) Add(ctx context.Context, ea models.ObjType) (int, error) {
	var est ifObjTypeStorage
	est = &esv.storage

	ai, err := est.Add(ctx, ea)

	if err != nil {
		log.Println("ObjTypeStorage.Add", err)
		return 0, err
	}

	return ai, nil
}

//func (esv *ObjTypeService) Upd(ctx context.Context, eu models.ObjType) (int, error)
func (esv *ObjTypeService) Upd(ctx context.Context, eu models.ObjType) (int, error) {
	var est ifObjTypeStorage
	est = &esv.storage

	ui, err := est.Upd(ctx, eu)

	if err != nil {
		log.Println("ObjTypeStorage.Upd", err)
		return 0, err
	}

	return ui, nil
}

//func (esv *ObjTypeService) Del(ctx context.Context, ed []int) ([]int, error) {
func (esv *ObjTypeService) Del(ctx context.Context, ed []int) ([]int, error) {
	var est ifObjTypeStorage
	est = &esv.storage

	res, err := est.Del(ctx, ed)

	if err != nil {
		log.Println("ObjTypeStorage.Del", err)
		return []int{}, err
	}

	return res, nil
}

//func (esv *ObjTypeService) GetOne(ctx context.Context, i int) (models.ObjType_count, error)
func (esv *ObjTypeService) GetOne(ctx context.Context, i int) (models.ObjType_count, error) {
	var est ifObjTypeStorage
	est = &esv.storage

	out_count, err := est.GetOne(ctx, i)

	if err != nil {
		log.Println("ObjTypeStorage.GetOne", err)
		return models.ObjType_count{Values: []models.ObjType{}, Count: 0, Auth: models.Auth{}}, err
	}

	return out_count, nil
}
