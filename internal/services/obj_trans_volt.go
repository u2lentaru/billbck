package services

import (
	"context"
	"log"

	"github.com/u2lentaru/billbck/internal/adapters/db/pgsql"
	"github.com/u2lentaru/billbck/internal/models"
)

type ObjTransVoltService struct {
	storage pgsql.ObjTransVoltStorage
}

type ifObjTransVoltStorage interface {
	GetList(ctx context.Context, pg, pgs int, gs1, gs2 string, ord int, dsc bool) (models.ObjTransVolt_count, error)
	Add(ctx context.Context, ea models.ObjTransVolt) (int, error)
	Upd(ctx context.Context, eu models.ObjTransVolt) (int, error)
	Del(ctx context.Context, ed []int) ([]int, error)
	GetOne(ctx context.Context, i int) (models.ObjTransVolt_count, error)
	GetObj(ctx context.Context, gs1, gs2 string) (models.ObjTransVolt_count, error)
}

//NewObjTransVoltService(storage pgsql.ObjTransVoltStorage) *ObjTransVoltService
func NewObjTransVoltService(storage pgsql.ObjTransVoltStorage) *ObjTransVoltService {
	return &ObjTransVoltService{storage}
}

//func (esv *ObjTransVoltService) GetList(ctx context.Context, pg, pgs int, gs1, gs2 string, ord int, dsc bool) (models.ObjTransVolt_count, error)
func (esv *ObjTransVoltService) GetList(ctx context.Context, pg, pgs int, gs1, gs2 string, ord int, dsc bool) (models.ObjTransVolt_count, error) {
	var est ifObjTransVoltStorage
	est = &esv.storage

	out_count, err := est.GetList(ctx, pg, pgs, gs1, gs2, ord, dsc)

	if err != nil {
		log.Println("ObjTransVoltStorage.GetList", err)
		return models.ObjTransVolt_count{Values: []models.ObjTransVolt{}, Count: 0, Auth: models.Auth{}}, err
	}

	return out_count, nil
}

//func (esv *ObjTransVoltService) Add(ctx context.Context, ea models.ObjTransVolt) (int, error)
func (esv *ObjTransVoltService) Add(ctx context.Context, ea models.ObjTransVolt) (int, error) {
	var est ifObjTransVoltStorage
	est = &esv.storage

	ai, err := est.Add(ctx, ea)

	if err != nil {
		log.Println("ObjTransVoltStorage.Add", err)
		return 0, err
	}

	return ai, nil
}

//func (esv *ObjTransVoltService) Upd(ctx context.Context, eu models.ObjTransVolt) (int, error)
func (esv *ObjTransVoltService) Upd(ctx context.Context, eu models.ObjTransVolt) (int, error) {
	var est ifObjTransVoltStorage
	est = &esv.storage

	ui, err := est.Upd(ctx, eu)

	if err != nil {
		log.Println("ObjTransVoltStorage.Upd", err)
		return 0, err
	}

	return ui, nil
}

//func (esv *ObjTransVoltService) Del(ctx context.Context, ed []int) ([]int, error)
func (esv *ObjTransVoltService) Del(ctx context.Context, ed []int) ([]int, error) {
	var est ifObjTransVoltStorage
	est = &esv.storage

	res, err := est.Del(ctx, ed)

	if err != nil {
		log.Println("ObjTransVoltStorage.Del", err)
		return []int{}, err
	}

	return res, nil
}

//func (esv *ObjTransVoltService) GetOne(ctx context.Context, i int) (models.ObjTransVolt_count, error)
func (esv *ObjTransVoltService) GetOne(ctx context.Context, i int) (models.ObjTransVolt_count, error) {
	var est ifObjTransVoltStorage
	est = &esv.storage

	out_count, err := est.GetOne(ctx, i)

	if err != nil {
		log.Println("ObjTransVoltStorage.GetOne", err)
		return models.ObjTransVolt_count{Values: []models.ObjTransVolt{}, Count: 0, Auth: models.Auth{}}, err
	}

	return out_count, nil
}

//func (esv *ObjTransVoltService) GetObj(ctx context.Context, gs1, gs2 string) (models.ObjTransVolt_count, error)
func (esv *ObjTransVoltService) GetObj(ctx context.Context, gs1, gs2 string) (models.ObjTransVolt_count, error) {
	var est ifObjTransVoltStorage
	est = &esv.storage

	out_count, err := est.GetObj(ctx, gs1, gs2)

	if err != nil {
		log.Println("ObjTransVoltStorage.GetObj", err)
		return models.ObjTransVolt_count{Values: []models.ObjTransVolt{}, Count: 0, Auth: models.Auth{}}, err
	}

	return out_count, nil
}
