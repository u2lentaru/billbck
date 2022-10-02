package services

import (
	"context"
	"log"

	"github.com/u2lentaru/billbck/internal/adapters/db/pgsql"
	"github.com/u2lentaru/billbck/internal/models"
)

type ObjTransCurrService struct {
	storage pgsql.ObjTransCurrStorage
}

type ifObjTransCurrStorage interface {
	GetList(ctx context.Context, pg, pgs int, gs1, gs2 string, ord int, dsc bool) (models.ObjTransCurr_count, error)
	Add(ctx context.Context, ea models.ObjTransCurr) (int, error)
	Upd(ctx context.Context, eu models.ObjTransCurr) (int, error)
	Del(ctx context.Context, ed []int) ([]int, error)
	GetOne(ctx context.Context, i int) (models.ObjTransCurr_count, error)
	GetObj(ctx context.Context, gs1, gs2 string) (models.ObjTransCurr_count, error)
}

//NewObjTransCurrService(storage pgsql.ObjTransCurrStorage) *ObjTransCurrService
func NewObjTransCurrService(storage pgsql.ObjTransCurrStorage) *ObjTransCurrService {
	return &ObjTransCurrService{storage}
}

//func (esv *ObjTransCurrService) GetList(ctx context.Context, pg, pgs int, gs1, gs2 string, ord int, dsc bool) (models.ObjTransCurr_count, error)
func (esv *ObjTransCurrService) GetList(ctx context.Context, pg, pgs int, gs1, gs2 string, ord int, dsc bool) (models.ObjTransCurr_count, error) {
	var est ifObjTransCurrStorage
	est = &esv.storage

	out_count, err := est.GetList(ctx, pg, pgs, gs1, gs2, ord, dsc)

	if err != nil {
		log.Println("ObjTransCurrStorage.GetList", err)
		return models.ObjTransCurr_count{Values: []models.ObjTransCurr{}, Count: 0, Auth: models.Auth{}}, err
	}

	return out_count, nil
}

//func (esv *ObjTransCurrService) Add(ctx context.Context, ea models.ObjTransCurr) (int, error)
func (esv *ObjTransCurrService) Add(ctx context.Context, ea models.ObjTransCurr) (int, error) {
	var est ifObjTransCurrStorage
	est = &esv.storage

	ai, err := est.Add(ctx, ea)

	if err != nil {
		log.Println("ObjTransCurrStorage.Add", err)
		return 0, err
	}

	return ai, nil
}

//func (esv *ObjTransCurrService) Upd(ctx context.Context, eu models.ObjTransCurr) (int, error)
func (esv *ObjTransCurrService) Upd(ctx context.Context, eu models.ObjTransCurr) (int, error) {
	var est ifObjTransCurrStorage
	est = &esv.storage

	ui, err := est.Upd(ctx, eu)

	if err != nil {
		log.Println("ObjTransCurrStorage.Upd", err)
		return 0, err
	}

	return ui, nil
}

//func (esv *ObjTransCurrService) Del(ctx context.Context, ed []int) ([]int, error)
func (esv *ObjTransCurrService) Del(ctx context.Context, ed []int) ([]int, error) {
	var est ifObjTransCurrStorage
	est = &esv.storage

	res, err := est.Del(ctx, ed)

	if err != nil {
		log.Println("ObjTransCurrStorage.Del", err)
		return []int{}, err
	}

	return res, nil
}

//func (esv *ObjTransCurrService) GetOne(ctx context.Context, i int) (models.ObjTransCurr_count, error)
func (esv *ObjTransCurrService) GetOne(ctx context.Context, i int) (models.ObjTransCurr_count, error) {
	var est ifObjTransCurrStorage
	est = &esv.storage

	out_count, err := est.GetOne(ctx, i)

	if err != nil {
		log.Println("ObjTransCurrStorage.GetOne", err)
		return models.ObjTransCurr_count{Values: []models.ObjTransCurr{}, Count: 0, Auth: models.Auth{}}, err
	}

	return out_count, nil
}

//func (esv *ObjTransCurrService) GetObj(ctx context.Context, gs1, gs2 string) (models.ObjTransCurr_count, error)
func (esv *ObjTransCurrService) GetObj(ctx context.Context, gs1, gs2 string) (models.ObjTransCurr_count, error) {
	var est ifObjTransCurrStorage
	est = &esv.storage

	out_count, err := est.GetObj(ctx, gs1, gs2)

	if err != nil {
		log.Println("ObjTransCurrStorage.GetObj", err)
		return models.ObjTransCurr_count{Values: []models.ObjTransCurr{}, Count: 0, Auth: models.Auth{}}, err
	}

	return out_count, nil
}
