package services

import (
	"context"
	"log"

	"github.com/u2lentaru/billbck/internal/adapters/db/pgsql"
	"github.com/u2lentaru/billbck/internal/models"
)

type PuValueService struct {
	storage pgsql.PuValueStorage
}

type ifPuValueStorage interface {
	GetList(ctx context.Context, pg, pgs, gs1, ord int, dsc bool) (models.PuValue_count, error)
	Add(ctx context.Context, ea models.PuValue) (int, error)
	Upd(ctx context.Context, eu models.PuValue) (int, error)
	Del(ctx context.Context, ed []int) ([]int, error)
	GetOne(ctx context.Context, i int) (models.PuValue_count, error)
	AskuePrev(ctx context.Context, af models.AskueFile) (models.PuValueAskue_count, error)
	AskueLoad(ctx context.Context, af models.AskueFile) (models.AskueLoadRes, error)
}

//func NewPuValueService(storage pgsql.PuValueStorage) *PuValueService
func NewPuValueService(storage pgsql.PuValueStorage) *PuValueService {
	return &PuValueService{storage}
}

//func (esv *PuValueService) GetList(ctx context.Context, pg, pgs, gs1, ord int, dsc bool) (models.PuValue_count, error)
func (esv *PuValueService) GetList(ctx context.Context, pg, pgs, gs1, ord int, dsc bool) (models.PuValue_count, error) {
	var est ifPuValueStorage
	est = &esv.storage

	out_count, err := est.GetList(ctx, pg, pgs, gs1, ord, dsc)

	if err != nil {
		log.Println("PuValueStorage.GetList", err)
		return models.PuValue_count{Values: []models.PuValue{}, Count: 0, Auth: models.Auth{}}, err
	}

	return out_count, nil
}

//func (esv *PuValueService) Add(ctx context.Context, ea models.PuValue) (int, error)
func (esv *PuValueService) Add(ctx context.Context, ea models.PuValue) (int, error) {
	var est ifPuValueStorage
	est = &esv.storage

	ai, err := est.Add(ctx, ea)

	if err != nil {
		log.Println("PuValueStorage.Add", err)
		return 0, err
	}

	return ai, nil
}

//func (esv *PuValueService) Upd(ctx context.Context, eu models.PuValue) (int, error)
func (esv *PuValueService) Upd(ctx context.Context, eu models.PuValue) (int, error) {
	var est ifPuValueStorage
	est = &esv.storage

	ui, err := est.Upd(ctx, eu)

	if err != nil {
		log.Println("PuValueStorage.Upd", err)
		return 0, err
	}

	return ui, nil
}

//func (esv *PuValueService) Del(ctx context.Context, ed []int) ([]int, error) {
func (esv *PuValueService) Del(ctx context.Context, ed []int) ([]int, error) {
	var est ifPuValueStorage
	est = &esv.storage

	res, err := est.Del(ctx, ed)

	if err != nil {
		log.Println("PuValueStorage.Del", err)
		return []int{}, err
	}

	return res, nil
}

//func (esv *PuValueService) GetOne(ctx context.Context, i int) (models.PuValue_count, error)
func (esv *PuValueService) GetOne(ctx context.Context, i int) (models.PuValue_count, error) {
	var est ifPuValueStorage
	est = &esv.storage

	out_count, err := est.GetOne(ctx, i)

	if err != nil {
		log.Println("PuValueStorage.GetOne", err)
		return models.PuValue_count{Values: []models.PuValue{}, Count: 0, Auth: models.Auth{}}, err
	}

	return out_count, nil
}

//func (esv *PuValueService) AskuePrev(ctx context.Context, af models.AskueFile) (models.PuValueAskue_count, error)
func (esv *PuValueService) AskuePrev(ctx context.Context, af models.AskueFile) (models.PuValueAskue_count, error) {
	var est ifPuValueStorage
	est = &esv.storage

	p, err := est.AskuePrev(ctx, af)

	if err != nil {
		log.Println("PuValueStorage.AskuePrev", err)
		return models.PuValueAskue_count{}, err
	}

	return p, nil
}

//func (esv *PuValueService) AskueLoad(ctx context.Context, af models.AskueFile) (models.AskueLoadRes, error)
func (esv *PuValueService) AskueLoad(ctx context.Context, af models.AskueFile) (models.AskueLoadRes, error) {
	var est ifPuValueStorage
	est = &esv.storage

	p, err := est.AskueLoad(ctx, af)

	if err != nil {
		log.Println("PuValueStorage.AskueLoad", err)
		return models.AskueLoadRes{}, err
	}

	return p, nil
}
