package services

import (
	"context"
	"log"

	"github.com/u2lentaru/billbck/internal/adapters/db/pgsql"
	"github.com/u2lentaru/billbck/internal/models"
)

type ObjContractService struct {
	storage pgsql.ObjContractStorage
}

type ifObjContractStorage interface {
	GetList(ctx context.Context, pg, pgs, gs1, gs2, gs3 int, gs4, gs4f bool, ord int, dsc bool) (models.ObjContract_count, error)
	Add(ctx context.Context, ea models.ObjContract) (int, error)
	Upd(ctx context.Context, eu models.ObjContract) (int, error)
	Del(ctx context.Context, ed models.IdClose) (int, error)
	GetOne(ctx context.Context, i int, d string) (models.ObjContract_count, error)
}

//NewObjContractService(storage pg.ObjContractStorage) *ObjContractService
func NewObjContractService(storage pgsql.ObjContractStorage) *ObjContractService {
	return &ObjContractService{storage}
}

//func (esv *ObjContractService) GetList(ctx context.Context, pg, pgs, gs1, gs2, gs3 int, gs4, gs4f bool, ord int, dsc bool) (models.ObjContract_count, error)
func (esv *ObjContractService) GetList(ctx context.Context, pg, pgs, gs1, gs2, gs3 int, gs4, gs4f bool, ord int, dsc bool) (models.ObjContract_count, error) {
	var est ifObjContractStorage
	est = &esv.storage

	out_count, err := est.GetList(ctx, pg, pgs, gs1, gs2, gs3, gs4, gs4f, ord, dsc)

	if err != nil {
		log.Println("ObjContractStorage.GetList", err)
		return models.ObjContract_count{Values: []models.ObjContract{}, Count: 0, Auth: models.Auth{}}, err
	}

	return out_count, nil
}

//func (esv *ObjContractService) Add(ctx context.Context, ea models.ObjContract) (int, error)
func (esv *ObjContractService) Add(ctx context.Context, ea models.ObjContract) (int, error) {
	var est ifObjContractStorage
	est = &esv.storage

	ai, err := est.Add(ctx, ea)

	if err != nil {
		log.Println("ObjContractStorage.Add", err)
		return 0, err
	}

	return ai, nil
}

//func (esv *ObjContractService) Upd(ctx context.Context, eu models.ObjContract) (int, error)
func (esv *ObjContractService) Upd(ctx context.Context, eu models.ObjContract) (int, error) {
	var est ifObjContractStorage
	est = &esv.storage

	ui, err := est.Upd(ctx, eu)

	if err != nil {
		log.Println("ObjContractStorage.Upd", err)
		return 0, err
	}

	return ui, nil
}

//func (esv *ObjContractService) Del(ctx context.Context, ed models.IdClose) (int, error)
func (esv *ObjContractService) Del(ctx context.Context, ed models.IdClose) (int, error) {
	var est ifObjContractStorage
	est = &esv.storage

	res, err := est.Del(ctx, ed)

	if err != nil {
		log.Println("ObjContractStorage.Del", err)
		return 0, err
	}

	return res, nil
}

//func (esv *ObjContractService) GetOne(ctx context.Context, i int) (models.ObjContract_count, error)
func (esv *ObjContractService) GetOne(ctx context.Context, i int, d string) (models.ObjContract_count, error) {
	var est ifObjContractStorage
	est = &esv.storage

	out_count, err := est.GetOne(ctx, i, d)

	if err != nil {
		log.Println("ObjContractStorage.GetOne", err)
		return models.ObjContract_count{Values: []models.ObjContract{}, Count: 0, Auth: models.Auth{}}, err
	}

	return out_count, nil
}
