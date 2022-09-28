package services

import (
	"context"
	"log"

	"github.com/u2lentaru/billbck/internal/adapters/db/pgsql"
	"github.com/u2lentaru/billbck/internal/models"
)

type ObjLineService struct {
	storage pgsql.ObjLineStorage
}

type ifObjLineStorage interface {
	GetList(ctx context.Context, pg, pgs int, gs1, gs2 string, ord int, dsc bool) (models.ObjLine_count, error)
	Add(ctx context.Context, ea models.ObjLine) (int, error)
	Upd(ctx context.Context, eu models.ObjLine) (int, error)
	Del(ctx context.Context, ed []int) ([]int, error)
	GetOne(ctx context.Context, i int) (models.ObjLine_count, error)
	GetObj(ctx context.Context, gs1, gs2 string) (models.ObjLine_count, error)
}

//NewObjLineService(storage pgsql.ObjLineStorage) *ObjLineService
func NewObjLineService(storage pgsql.ObjLineStorage) *ObjLineService {
	return &ObjLineService{storage}
}

//func (esv *ObjLineService) GetList(ctx context.Context, pg, pgs int, gs1, gs2 string, ord int, dsc bool) (models.ObjLine_count, error)
func (esv *ObjLineService) GetList(ctx context.Context, pg, pgs int, gs1, gs2 string, ord int, dsc bool) (models.ObjLine_count, error) {
	var est ifObjLineStorage
	est = &esv.storage

	out_count, err := est.GetList(ctx, pg, pgs, gs1, gs2, ord, dsc)

	if err != nil {
		log.Println("ObjLineStorage.GetList", err)
		return models.ObjLine_count{Values: []models.ObjLine{}, Count: 0, Auth: models.Auth{}}, err
	}

	return out_count, nil
}

//func (esv *ObjLineService) Add(ctx context.Context, ea models.ObjLine) (int, error)
func (esv *ObjLineService) Add(ctx context.Context, ea models.ObjLine) (int, error) {
	var est ifObjLineStorage
	est = &esv.storage

	ai, err := est.Add(ctx, ea)

	if err != nil {
		log.Println("ObjLineStorage.Add", err)
		return 0, err
	}

	return ai, nil
}

//func (esv *ObjLineService) Upd(ctx context.Context, eu models.ObjLine) (int, error)
func (esv *ObjLineService) Upd(ctx context.Context, eu models.ObjLine) (int, error) {
	var est ifObjLineStorage
	est = &esv.storage

	ui, err := est.Upd(ctx, eu)

	if err != nil {
		log.Println("ObjLineStorage.Upd", err)
		return 0, err
	}

	return ui, nil
}

//func (esv *ObjLineService) Del(ctx context.Context, ed []int) ([]int, error)
func (esv *ObjLineService) Del(ctx context.Context, ed []int) ([]int, error) {
	var est ifObjLineStorage
	est = &esv.storage

	res, err := est.Del(ctx, ed)

	if err != nil {
		log.Println("ObjLineStorage.Del", err)
		return []int{}, err
	}

	return res, nil
}

//func (esv *ObjLineService) GetOne(ctx context.Context, i int) (models.ObjLine_count, error)
func (esv *ObjLineService) GetOne(ctx context.Context, i int) (models.ObjLine_count, error) {
	var est ifObjLineStorage
	est = &esv.storage

	out_count, err := est.GetOne(ctx, i)

	if err != nil {
		log.Println("ObjLineStorage.GetOne", err)
		return models.ObjLine_count{Values: []models.ObjLine{}, Count: 0, Auth: models.Auth{}}, err
	}

	return out_count, nil
}

//func (esv *ObjLineService) GetObj(ctx context.Context, gs1, gs2 string) (models.ObjLine_count, error)
func (esv *ObjLineService) GetObj(ctx context.Context, gs1, gs2 string) (models.ObjLine_count, error) {
	var est ifObjLineStorage
	est = &esv.storage

	out_count, err := est.GetObj(ctx, gs1, gs2)

	if err != nil {
		log.Println("ObjLineStorage.GetObj", err)
		return models.ObjLine_count{Values: []models.ObjLine{}, Count: 0, Auth: models.Auth{}}, err
	}

	return out_count, nil
}
