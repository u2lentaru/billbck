package services

import (
	"context"
	"log"

	"github.com/u2lentaru/billbck/internal/adapters/db/pgsql"
	"github.com/u2lentaru/billbck/internal/models"
)

type ReasonService struct {
	storage pgsql.ReasonStorage
}

type ifReasonStorage interface {
	GetList(ctx context.Context, pg, pgs int, gs1 string, ord int, dsc bool) (models.Reason_count, error)
	Add(ctx context.Context, ea models.Reason) (int, error)
	Upd(ctx context.Context, eu models.Reason) (int, error)
	Del(ctx context.Context, ed []int) ([]int, error)
	GetOne(ctx context.Context, i int) (models.Reason_count, error)
}

//func NewReasonService(storage pgsql.Reasontorage) *ReasonService
func NewReasonService(storage pgsql.ReasonStorage) *ReasonService {
	return &ReasonService{storage}
}

//func (esv *ReasonService) GetList(ctx context.Context, pg, pgs int, gs1 string, ord int, dsc bool) (models.Reason_count, error)
func (esv *ReasonService) GetList(ctx context.Context, pg, pgs int, gs1 string, ord int, dsc bool) (models.Reason_count, error) {
	var est ifReasonStorage
	est = &esv.storage

	out_count, err := est.GetList(ctx, pg, pgs, gs1, ord, dsc)

	if err != nil {
		log.Println("ReasonStorage.GetList", err)
		return models.Reason_count{Values: []models.Reason{}, Count: 0, Auth: models.Auth{}}, err
	}

	return out_count, nil
}

//func (esv *ReasonService) Add(ctx context.Context, ea models.Reason) (int, error)
func (esv *ReasonService) Add(ctx context.Context, ea models.Reason) (int, error) {
	var est ifReasonStorage
	est = &esv.storage

	ai, err := est.Add(ctx, ea)

	if err != nil {
		log.Println("ReasonStorage.Add", err)
		return 0, err
	}

	return ai, nil
}

//func (esv *ReasonService) Upd(ctx context.Context, eu models.Reason) (int, error)
func (esv *ReasonService) Upd(ctx context.Context, eu models.Reason) (int, error) {
	var est ifReasonStorage
	est = &esv.storage

	ui, err := est.Upd(ctx, eu)

	if err != nil {
		log.Println("ReasonStorage.Upd", err)
		return 0, err
	}

	return ui, nil
}

//func (esv *ReasonService) Del(ctx context.Context, ed []int) ([]int, error) {
func (esv *ReasonService) Del(ctx context.Context, ed []int) ([]int, error) {
	var est ifReasonStorage
	est = &esv.storage

	res, err := est.Del(ctx, ed)

	if err != nil {
		log.Println("ReasonStorage.Del", err)
		return []int{}, err
	}

	return res, nil
}

//func (esv *ReasonService) GetOne(ctx context.Context, i int) (models.Reason_count, error)
func (esv *ReasonService) GetOne(ctx context.Context, i int) (models.Reason_count, error) {
	var est ifReasonStorage
	est = &esv.storage

	out_count, err := est.GetOne(ctx, i)

	if err != nil {
		log.Println("ReasonStorage.GetOne", err)
		return models.Reason_count{Values: []models.Reason{}, Count: 0, Auth: models.Auth{}}, err
	}

	return out_count, nil
}
