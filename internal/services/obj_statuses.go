package services

import (
	"context"
	"log"

	"github.com/u2lentaru/billbck/internal/adapters/db/pgsql"
	"github.com/u2lentaru/billbck/internal/models"
)

type ObjStatusService struct {
	storage pgsql.ObjStatusStorage
}

type ifObjStatusStorage interface {
	GetList(ctx context.Context, pg, pgs int, gs1 string, ord int, dsc bool) (models.ObjStatus_count, error)
	Add(ctx context.Context, ea models.ObjStatus) (int, error)
	Upd(ctx context.Context, eu models.ObjStatus) (int, error)
	Del(ctx context.Context, ed []int) ([]int, error)
	GetOne(ctx context.Context, i int) (models.ObjStatus_count, error)
}

//func NewObjStatusService(storage pgsql.ObjStatusStorage) *ObjStatusService
func NewObjStatusService(storage pgsql.ObjStatusStorage) *ObjStatusService {
	return &ObjStatusService{storage}
}

//func (esv *ObjStatusService) GetList(ctx context.Context, pg, pgs int, gs1 string, ord int, dsc bool) (models.ObjStatus_count, error)
func (esv *ObjStatusService) GetList(ctx context.Context, pg, pgs int, gs1 string, ord int, dsc bool) (models.ObjStatus_count, error) {
	var est ifObjStatusStorage
	est = &esv.storage

	out_count, err := est.GetList(ctx, pg, pgs, gs1, ord, dsc)

	if err != nil {
		log.Println("ObjStatusStorage.GetList", err)
		return models.ObjStatus_count{Values: []models.ObjStatus{}, Count: 0, Auth: models.Auth{}}, err
	}

	return out_count, nil
}

//func (esv *ObjStatusService) Add(ctx context.Context, ea models.ObjStatus) (int, error)
func (esv *ObjStatusService) Add(ctx context.Context, ea models.ObjStatus) (int, error) {
	var est ifObjStatusStorage
	est = &esv.storage

	ai, err := est.Add(ctx, ea)

	if err != nil {
		log.Println("ObjStatusStorage.Add", err)
		return 0, err
	}

	return ai, nil
}

//func (esv *ObjStatusService) Upd(ctx context.Context, eu models.ObjStatus) (int, error)
func (esv *ObjStatusService) Upd(ctx context.Context, eu models.ObjStatus) (int, error) {
	var est ifObjStatusStorage
	est = &esv.storage

	ui, err := est.Upd(ctx, eu)

	if err != nil {
		log.Println("ObjStatusStorage.Upd", err)
		return 0, err
	}

	return ui, nil
}

//func (esv *ObjStatusService) Del(ctx context.Context, ed []int) ([]int, error) {
func (esv *ObjStatusService) Del(ctx context.Context, ed []int) ([]int, error) {
	var est ifObjStatusStorage
	est = &esv.storage

	res, err := est.Del(ctx, ed)

	if err != nil {
		log.Println("ObjStatusStorage.Del", err)
		return []int{}, err
	}

	return res, nil
}

//func (esv *ObjStatusService) GetOne(ctx context.Context, i int) (models.ObjStatus_count, error)
func (esv *ObjStatusService) GetOne(ctx context.Context, i int) (models.ObjStatus_count, error) {
	var est ifObjStatusStorage
	est = &esv.storage

	out_count, err := est.GetOne(ctx, i)

	if err != nil {
		log.Println("ObjStatusStorage.GetOne", err)
		return models.ObjStatus_count{Values: []models.ObjStatus{}, Count: 0, Auth: models.Auth{}}, err
	}

	return out_count, nil
}
