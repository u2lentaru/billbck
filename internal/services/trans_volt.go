package services

import (
	"context"
	"log"

	"github.com/u2lentaru/billbck/internal/adapters/db/pgsql"
	"github.com/u2lentaru/billbck/internal/models"
)

type TransVoltService struct {
	storage pgsql.TransVoltStorage
}

type ifTransVoltStorage interface {
	GetList(ctx context.Context, pg, pgs int, gs1 string, ord int, dsc bool) (models.TransVolt_count, error)
	Add(ctx context.Context, ea models.TransVolt) (int, error)
	Upd(ctx context.Context, eu models.TransVolt) (int, error)
	Del(ctx context.Context, ed []int) ([]int, error)
	GetOne(ctx context.Context, i int) (models.TransVolt_count, error)
}

//func NewTransVoltService(storage pgsql.TransVolttorage) *TransVoltService
func NewTransVoltService(storage pgsql.TransVoltStorage) *TransVoltService {
	return &TransVoltService{storage}
}

//func (esv *TransVoltService) GetList(ctx context.Context, pg, pgs int, gs1 string, ord int, dsc bool) (models.TransVolt_count, error)
func (esv *TransVoltService) GetList(ctx context.Context, pg, pgs int, gs1 string, ord int, dsc bool) (models.TransVolt_count, error) {
	var est ifTransVoltStorage
	est = &esv.storage

	out_count, err := est.GetList(ctx, pg, pgs, gs1, ord, dsc)

	if err != nil {
		log.Println("TransVoltStorage.GetList", err)
		return models.TransVolt_count{Values: []models.TransVolt{}, Count: 0, Auth: models.Auth{}}, err
	}

	return out_count, nil
}

//func (esv *TransVoltService) Add(ctx context.Context, ea models.TransVolt) (int, error)
func (esv *TransVoltService) Add(ctx context.Context, ea models.TransVolt) (int, error) {
	var est ifTransVoltStorage
	est = &esv.storage

	ai, err := est.Add(ctx, ea)

	if err != nil {
		log.Println("TransVoltStorage.Add", err)
		return 0, err
	}

	return ai, nil
}

//func (esv *TransVoltService) Upd(ctx context.Context, eu models.TransVolt) (int, error)
func (esv *TransVoltService) Upd(ctx context.Context, eu models.TransVolt) (int, error) {
	var est ifTransVoltStorage
	est = &esv.storage

	ui, err := est.Upd(ctx, eu)

	if err != nil {
		log.Println("TransVoltStorage.Upd", err)
		return 0, err
	}

	return ui, nil
}

//func (esv *TransVoltService) Del(ctx context.Context, ed []int) ([]int, error) {
func (esv *TransVoltService) Del(ctx context.Context, ed []int) ([]int, error) {
	var est ifTransVoltStorage
	est = &esv.storage

	res, err := est.Del(ctx, ed)

	if err != nil {
		log.Println("TransVoltStorage.Del", err)
		return []int{}, err
	}

	return res, nil
}

//func (esv *TransVoltService) GetOne(ctx context.Context, i int) (models.TransVolt_count, error)
func (esv *TransVoltService) GetOne(ctx context.Context, i int) (models.TransVolt_count, error) {
	var est ifTransVoltStorage
	est = &esv.storage

	out_count, err := est.GetOne(ctx, i)

	if err != nil {
		log.Println("TransVoltStorage.GetOne", err)
		return models.TransVolt_count{Values: []models.TransVolt{}, Count: 0, Auth: models.Auth{}}, err
	}

	return out_count, nil
}
