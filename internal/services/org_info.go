package services

import (
	"context"
	"log"

	"github.com/u2lentaru/billbck/internal/adapters/db/pgsql"
	"github.com/u2lentaru/billbck/internal/models"
)

type OrgInfoService struct {
	storage pgsql.OrgInfoStorage
}

type ifOrgInfoStorage interface {
	GetList(ctx context.Context, pg, pgs int, gs1, gs2 string, ord int, dsc bool) (models.OrgInfo_count, error)
	Add(ctx context.Context, ea models.OrgInfo) (int, error)
	Upd(ctx context.Context, eu models.OrgInfo) (int, error)
	Del(ctx context.Context, ed []int) ([]int, error)
	GetOne(ctx context.Context, i int) (models.OrgInfo_count, error)
}

//NewOrgInfoService(storage pg.OrgInfoStorage) *OrgInfoService
func NewOrgInfoService(storage pgsql.OrgInfoStorage) *OrgInfoService {
	return &OrgInfoService{storage}
}

//func (esv *OrgInfoService) GetList(ctx context.Context, pg, pgs int, gs1, gs2 string, ord int, dsc bool) (models.OrgInfo_count, error)
func (esv *OrgInfoService) GetList(ctx context.Context, pg, pgs int, gs1, gs2 string, ord int, dsc bool) (models.OrgInfo_count, error) {
	var est ifOrgInfoStorage
	est = &esv.storage

	out_count, err := est.GetList(ctx, pg, pgs, gs1, gs2, ord, dsc)

	if err != nil {
		log.Println("OrgInfoStorage.GetList", err)
		return models.OrgInfo_count{Values: []models.OrgInfo{}, Count: 0, Auth: models.Auth{}}, err
	}

	return out_count, nil
}

//func (esv *OrgInfoService) Add(ctx context.Context, ea models.OrgInfo) (int, error)
func (esv *OrgInfoService) Add(ctx context.Context, ea models.OrgInfo) (int, error) {
	var est ifOrgInfoStorage
	est = &esv.storage

	ai, err := est.Add(ctx, ea)

	if err != nil {
		log.Println("OrgInfoStorage.Add", err)
		return 0, err
	}

	return ai, nil
}

//func (esv *OrgInfoService) Upd(ctx context.Context, eu models.OrgInfo) (int, error)
func (esv *OrgInfoService) Upd(ctx context.Context, eu models.OrgInfo) (int, error) {
	var est ifOrgInfoStorage
	est = &esv.storage

	ui, err := est.Upd(ctx, eu)

	if err != nil {
		log.Println("OrgInfoStorage.Upd", err)
		return 0, err
	}

	return ui, nil
}

//func (esv *OrgInfoService) Del(ctx context.Context, ed []int) ([]int, error)
func (esv *OrgInfoService) Del(ctx context.Context, ed []int) ([]int, error) {
	var est ifOrgInfoStorage
	est = &esv.storage

	res, err := est.Del(ctx, ed)

	if err != nil {
		log.Println("OrgInfoStorage.Del", err)
		return []int{}, err
	}

	return res, nil
}

//func (esv *OrgInfoService) GetOne(ctx context.Context, i int) (models.OrgInfo_count, error)
func (esv *OrgInfoService) GetOne(ctx context.Context, i int) (models.OrgInfo_count, error) {
	var est ifOrgInfoStorage
	est = &esv.storage

	out_count, err := est.GetOne(ctx, i)

	if err != nil {
		log.Println("OrgInfoStorage.GetOne", err)
		return models.OrgInfo_count{Values: []models.OrgInfo{}, Count: 0, Auth: models.Auth{}}, err
	}

	return out_count, nil
}
