package services

import (
	"context"
	"log"

	"github.com/u2lentaru/billbck/internal/adapters/db/pgsql"
	"github.com/u2lentaru/billbck/internal/models"
)

type TpService struct {
	storage pgsql.TpStorage
}

type ifTpStorage interface {
	GetList(ctx context.Context, pg, pgs int, gs1 string, ord int, dsc bool) (models.Tp_count, error)
	Add(ctx context.Context, ea models.Tp) (int, error)
	Upd(ctx context.Context, eu models.Tp) (int, error)
	Del(ctx context.Context, ed []int) ([]int, error)
	GetOne(ctx context.Context, i int) (models.Tp_count, error)
}

//func NewTpService(storage pgsql.TpStorage) *TpService
func NewTpService(storage pgsql.TpStorage) *TpService {
	return &TpService{storage}
}

//func (esv *TpService) GetList(ctx context.Context, pg, pgs int, gs1 string, ord int, dsc bool) (models.Tp_count, error)
func (esv *TpService) GetList(ctx context.Context, pg, pgs int, gs1 string, ord int, dsc bool) (models.Tp_count, error) {
	var est ifTpStorage
	est = &esv.storage

	out_count, err := est.GetList(ctx, pg, pgs, gs1, ord, dsc)

	if err != nil {
		log.Println("TpStorage.GetList", err)
		return models.Tp_count{Values: []models.Tp{}, Count: 0, Auth: models.Auth{}}, err
	}

	return out_count, nil
}

//func (esv *TpService) Add(ctx context.Context, ea models.Tp) (int, error)
func (esv *TpService) Add(ctx context.Context, ea models.Tp) (int, error) {
	var est ifTpStorage
	est = &esv.storage

	ai, err := est.Add(ctx, ea)

	if err != nil {
		log.Println("TpStorage.Add", err)
		return 0, err
	}

	return ai, nil
}

//func (esv *TpService) Upd(ctx context.Context, eu models.Tp) (int, error)
func (esv *TpService) Upd(ctx context.Context, eu models.Tp) (int, error) {
	var est ifTpStorage
	est = &esv.storage

	ui, err := est.Upd(ctx, eu)

	if err != nil {
		log.Println("TpStorage.Upd", err)
		return 0, err
	}

	return ui, nil
}

//func (esv *TpService) Del(ctx context.Context, ed []int) ([]int, error) {
func (esv *TpService) Del(ctx context.Context, ed []int) ([]int, error) {
	var est ifTpStorage
	est = &esv.storage

	res, err := est.Del(ctx, ed)

	if err != nil {
		log.Println("TpStorage.Del", err)
		return []int{}, err
	}

	return res, nil
}

//func (esv *TpService) GetOne(ctx context.Context, i int) (models.Tp_count, error)
func (esv *TpService) GetOne(ctx context.Context, i int) (models.Tp_count, error) {
	var est ifTpStorage
	est = &esv.storage

	out_count, err := est.GetOne(ctx, i)

	if err != nil {
		log.Println("TpStorage.GetOne", err)
		return models.Tp_count{Values: []models.Tp{}, Count: 0, Auth: models.Auth{}}, err
	}

	return out_count, nil
}
