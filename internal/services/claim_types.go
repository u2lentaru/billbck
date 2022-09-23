package services

import (
	"context"
	"log"

	"github.com/u2lentaru/billbck/internal/adapters/db/pgsql"
	"github.com/u2lentaru/billbck/internal/models"
)

type ClaimTypeService struct {
	storage pgsql.ClaimTypeStorage
}

type ifClaimTypeStorage interface {
	GetList(ctx context.Context, pg, pgs int, gs1 string, ord int, dsc bool) (models.ClaimType_count, error)
	Add(ctx context.Context, ea models.ClaimType) (int, error)
	Upd(ctx context.Context, eu models.ClaimType) (int, error)
	Del(ctx context.Context, ed []int) ([]int, error)
	GetOne(ctx context.Context, i int) (models.ClaimType_count, error)
}

//func NewClaimTypeService(storage pgsql.ClaimTypeStorage) *ClaimTypeService
func NewClaimTypeService(storage pgsql.ClaimTypeStorage) *ClaimTypeService {
	return &ClaimTypeService{storage}
}

//func (esv *ClaimTypeService) GetList(ctx context.Context, pg, pgs int, gs1 string, ord int, dsc bool) (models.ClaimType_count, error)
func (esv *ClaimTypeService) GetList(ctx context.Context, pg, pgs int, gs1 string, ord int, dsc bool) (models.ClaimType_count, error) {
	var est ifClaimTypeStorage
	est = &esv.storage

	out_count, err := est.GetList(ctx, pg, pgs, gs1, ord, dsc)

	if err != nil {
		log.Println("ClaimTypeStorage.GetList", err)
		return models.ClaimType_count{Values: []models.ClaimType{}, Count: 0, Auth: models.Auth{}}, err
	}

	return out_count, nil
}

//func (esv *ClaimTypeService) Add(ctx context.Context, ea models.ClaimType) (int, error)
func (esv *ClaimTypeService) Add(ctx context.Context, ea models.ClaimType) (int, error) {
	var est ifClaimTypeStorage
	est = &esv.storage

	ai, err := est.Add(ctx, ea)

	if err != nil {
		log.Println("ClaimTypeStorage.Add", err)
		return 0, err
	}

	return ai, nil
}

//func (esv *ClaimTypeService) Upd(ctx context.Context, eu models.ClaimType) (int, error)
func (esv *ClaimTypeService) Upd(ctx context.Context, eu models.ClaimType) (int, error) {
	var est ifClaimTypeStorage
	est = &esv.storage

	ui, err := est.Upd(ctx, eu)

	if err != nil {
		log.Println("ClaimTypeStorage.Upd", err)
		return 0, err
	}

	return ui, nil
}

//func (esv *ClaimTypeService) Del(ctx context.Context, ed []int) ([]int, error) {
func (esv *ClaimTypeService) Del(ctx context.Context, ed []int) ([]int, error) {
	var est ifClaimTypeStorage
	est = &esv.storage

	res, err := est.Del(ctx, ed)

	if err != nil {
		log.Println("ClaimTypeStorage.Del", err)
		return []int{}, err
	}

	return res, nil
}

//func (esv *ClaimTypeService) GetOne(ctx context.Context, i int) (models.ClaimType_count, error)
func (esv *ClaimTypeService) GetOne(ctx context.Context, i int) (models.ClaimType_count, error) {
	var est ifClaimTypeStorage
	est = &esv.storage

	out_count, err := est.GetOne(ctx, i)

	if err != nil {
		log.Println("ClaimTypeStorage.GetOne", err)
		return models.ClaimType_count{Values: []models.ClaimType{}, Count: 0, Auth: models.Auth{}}, err
	}

	return out_count, nil
}
