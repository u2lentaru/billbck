package services

import (
	"context"
	"log"

	"github.com/u2lentaru/billbck/internal/adapters/db/pgsql"
	"github.com/u2lentaru/billbck/internal/models"
)

type ReliabilityService struct {
	storage pgsql.ReliabilityStorage
}

type ifReliabilityStorage interface {
	GetList(ctx context.Context, pg, pgs int, gs1 string, ord int, dsc bool) (models.Reliability_count, error)
	Add(ctx context.Context, ea models.Reliability) (int, error)
	Upd(ctx context.Context, eu models.Reliability) (int, error)
	Del(ctx context.Context, ed []int) ([]int, error)
	GetOne(ctx context.Context, i int) (models.Reliability_count, error)
}

//func NewReliabilityService(storage pgsql.Reliabilitytorage) *ReliabilityService
func NewReliabilityService(storage pgsql.ReliabilityStorage) *ReliabilityService {
	return &ReliabilityService{storage}
}

//func (esv *ReliabilityService) GetList(ctx context.Context, pg, pgs int, gs1 string, ord int, dsc bool) (models.Reliability_count, error)
func (esv *ReliabilityService) GetList(ctx context.Context, pg, pgs int, gs1 string, ord int, dsc bool) (models.Reliability_count, error) {
	var est ifReliabilityStorage
	est = &esv.storage

	out_count, err := est.GetList(ctx, pg, pgs, gs1, ord, dsc)

	if err != nil {
		log.Println("ReliabilityStorage.GetList", err)
		return models.Reliability_count{Values: []models.Reliability{}, Count: 0, Auth: models.Auth{}}, err
	}

	return out_count, nil
}

//func (esv *ReliabilityService) Add(ctx context.Context, ea models.Reliability) (int, error)
func (esv *ReliabilityService) Add(ctx context.Context, ea models.Reliability) (int, error) {
	var est ifReliabilityStorage
	est = &esv.storage

	ai, err := est.Add(ctx, ea)

	if err != nil {
		log.Println("ReliabilityStorage.Add", err)
		return 0, err
	}

	return ai, nil
}

//func (esv *ReliabilityService) Upd(ctx context.Context, eu models.Reliability) (int, error)
func (esv *ReliabilityService) Upd(ctx context.Context, eu models.Reliability) (int, error) {
	var est ifReliabilityStorage
	est = &esv.storage

	ui, err := est.Upd(ctx, eu)

	if err != nil {
		log.Println("ReliabilityStorage.Upd", err)
		return 0, err
	}

	return ui, nil
}

//func (esv *ReliabilityService) Del(ctx context.Context, ed []int) ([]int, error) {
func (esv *ReliabilityService) Del(ctx context.Context, ed []int) ([]int, error) {
	var est ifReliabilityStorage
	est = &esv.storage

	res, err := est.Del(ctx, ed)

	if err != nil {
		log.Println("ReliabilityStorage.Del", err)
		return []int{}, err
	}

	return res, nil
}

//func (esv *ReliabilityService) GetOne(ctx context.Context, i int) (models.Reliability_count, error)
func (esv *ReliabilityService) GetOne(ctx context.Context, i int) (models.Reliability_count, error) {
	var est ifReliabilityStorage
	est = &esv.storage

	out_count, err := est.GetOne(ctx, i)

	if err != nil {
		log.Println("ReliabilityStorage.GetOne", err)
		return models.Reliability_count{Values: []models.Reliability{}, Count: 0, Auth: models.Auth{}}, err
	}

	return out_count, nil
}
