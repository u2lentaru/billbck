package services

import (
	"context"
	"log"

	"github.com/u2lentaru/billbck/internal/adapters/db/pgsql"
	"github.com/u2lentaru/billbck/internal/models"
)

type PeriodService struct {
	storage pgsql.PeriodStorage
}

type ifPeriodStorage interface {
	GetList(ctx context.Context, pg, pgs int, gs1 string, ord int, dsc bool) (models.Period_count, error)
	Add(ctx context.Context, ea models.Period) (int, error)
	Upd(ctx context.Context, eu models.Period) (int, error)
	Del(ctx context.Context, ed []int) ([]int, error)
	GetOne(ctx context.Context, i int) (models.Period_count, error)
}

//func NewPeriodService(storage pgsql.Periodtorage) *PeriodService
func NewPeriodService(storage pgsql.PeriodStorage) *PeriodService {
	return &PeriodService{storage}
}

//func (esv *PeriodService) GetList(ctx context.Context, pg, pgs int, gs1 string, ord int, dsc bool) (models.Period_count, error)
func (esv *PeriodService) GetList(ctx context.Context, pg, pgs int, gs1 string, ord int, dsc bool) (models.Period_count, error) {
	var est ifPeriodStorage
	est = &esv.storage

	out_count, err := est.GetList(ctx, pg, pgs, gs1, ord, dsc)

	if err != nil {
		log.Println("PeriodStorage.GetList", err)
		return models.Period_count{Values: []models.Period{}, Count: 0, Auth: models.Auth{}}, err
	}

	return out_count, nil
}

//func (esv *PeriodService) Add(ctx context.Context, ea models.Period) (int, error)
func (esv *PeriodService) Add(ctx context.Context, ea models.Period) (int, error) {
	var est ifPeriodStorage
	est = &esv.storage

	ai, err := est.Add(ctx, ea)

	if err != nil {
		log.Println("PeriodStorage.Add", err)
		return 0, err
	}

	return ai, nil
}

//func (esv *PeriodService) Upd(ctx context.Context, eu models.Period) (int, error)
func (esv *PeriodService) Upd(ctx context.Context, eu models.Period) (int, error) {
	var est ifPeriodStorage
	est = &esv.storage

	ui, err := est.Upd(ctx, eu)

	if err != nil {
		log.Println("PeriodStorage.Upd", err)
		return 0, err
	}

	return ui, nil
}

//func (esv *PeriodService) Del(ctx context.Context, ed []int) ([]int, error) {
func (esv *PeriodService) Del(ctx context.Context, ed []int) ([]int, error) {
	var est ifPeriodStorage
	est = &esv.storage

	res, err := est.Del(ctx, ed)

	if err != nil {
		log.Println("PeriodStorage.Del", err)
		return []int{}, err
	}

	return res, nil
}

//func (esv *PeriodService) GetOne(ctx context.Context, i int) (models.Period_count, error)
func (esv *PeriodService) GetOne(ctx context.Context, i int) (models.Period_count, error) {
	var est ifPeriodStorage
	est = &esv.storage

	out_count, err := est.GetOne(ctx, i)

	if err != nil {
		log.Println("PeriodStorage.GetOne", err)
		return models.Period_count{Values: []models.Period{}, Count: 0, Auth: models.Auth{}}, err
	}

	return out_count, nil
}
