package services

import (
	"context"
	"log"

	"github.com/u2lentaru/billbck/internal/adapters/db/pgsql"
	"github.com/u2lentaru/billbck/internal/models"
)

type ObjTransPwrService struct {
	storage pgsql.ObjTransPwrStorage
}

type ifObjTransPwrStorage interface {
	GetList(ctx context.Context, pg, pgs int, gs1, gs2 string, ord int, dsc bool) (models.ObjTransPwr_count, error)
	Add(ctx context.Context, ea models.ObjTransPwr) (int, error)
	Upd(ctx context.Context, eu models.ObjTransPwr) (int, error)
	Del(ctx context.Context, ed []int) ([]int, error)
	GetOne(ctx context.Context, i int) (models.ObjTransPwr_count, error)
	GetObj(ctx context.Context, gs1, gs2 string) (models.ObjTransPwr_count, error)
}

//NewObjTransPwrService(storage pgsql.ObjTransPwrStorage) *ObjTransPwrService
func NewObjTransPwrService(storage pgsql.ObjTransPwrStorage) *ObjTransPwrService {
	return &ObjTransPwrService{storage}
}

//func (esv *ObjTransPwrService) GetList(ctx context.Context, pg, pgs int, gs1, gs2 string, ord int, dsc bool) (models.ObjTransPwr_count, error)
func (esv *ObjTransPwrService) GetList(ctx context.Context, pg, pgs int, gs1, gs2 string, ord int, dsc bool) (models.ObjTransPwr_count, error) {
	var est ifObjTransPwrStorage
	est = &esv.storage

	out_count, err := est.GetList(ctx, pg, pgs, gs1, gs2, ord, dsc)

	if err != nil {
		log.Println("ObjTransPwrStorage.GetList", err)
		return models.ObjTransPwr_count{Values: []models.ObjTransPwr{}, Count: 0, Auth: models.Auth{}}, err
	}

	return out_count, nil
}

//func (esv *ObjTransPwrService) Add(ctx context.Context, ea models.ObjTransPwr) (int, error)
func (esv *ObjTransPwrService) Add(ctx context.Context, ea models.ObjTransPwr) (int, error) {
	var est ifObjTransPwrStorage
	est = &esv.storage

	ai, err := est.Add(ctx, ea)

	if err != nil {
		log.Println("ObjTransPwrStorage.Add", err)
		return 0, err
	}

	return ai, nil
}

//func (esv *ObjTransPwrService) Upd(ctx context.Context, eu models.ObjTransPwr) (int, error)
func (esv *ObjTransPwrService) Upd(ctx context.Context, eu models.ObjTransPwr) (int, error) {
	var est ifObjTransPwrStorage
	est = &esv.storage

	ui, err := est.Upd(ctx, eu)

	if err != nil {
		log.Println("ObjTransPwrStorage.Upd", err)
		return 0, err
	}

	return ui, nil
}

//func (esv *ObjTransPwrService) Del(ctx context.Context, ed []int) ([]int, error)
func (esv *ObjTransPwrService) Del(ctx context.Context, ed []int) ([]int, error) {
	var est ifObjTransPwrStorage
	est = &esv.storage

	res, err := est.Del(ctx, ed)

	if err != nil {
		log.Println("ObjTransPwrStorage.Del", err)
		return []int{}, err
	}

	return res, nil
}

//func (esv *ObjTransPwrService) GetOne(ctx context.Context, i int) (models.ObjTransPwr_count, error)
func (esv *ObjTransPwrService) GetOne(ctx context.Context, i int) (models.ObjTransPwr_count, error) {
	var est ifObjTransPwrStorage
	est = &esv.storage

	out_count, err := est.GetOne(ctx, i)

	if err != nil {
		log.Println("ObjTransPwrStorage.GetOne", err)
		return models.ObjTransPwr_count{Values: []models.ObjTransPwr{}, Count: 0, Auth: models.Auth{}}, err
	}

	return out_count, nil
}

//func (esv *ObjTransPwrService) GetObj(ctx context.Context, gs1, gs2 string) (models.ObjTransPwr_count, error)
func (esv *ObjTransPwrService) GetObj(ctx context.Context, gs1, gs2 string) (models.ObjTransPwr_count, error) {
	var est ifObjTransPwrStorage
	est = &esv.storage

	out_count, err := est.GetObj(ctx, gs1, gs2)

	if err != nil {
		log.Println("ObjTransPwrStorage.GetObj", err)
		return models.ObjTransPwr_count{Values: []models.ObjTransPwr{}, Count: 0, Auth: models.Auth{}}, err
	}

	return out_count, nil
}
