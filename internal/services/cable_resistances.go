package services

import (
	"context"
	"log"

	"github.com/u2lentaru/billbck/internal/adapters/db/pgsql"
	"github.com/u2lentaru/billbck/internal/models"
)

type CableResistanceService struct {
	storage pgsql.CableResistanceStorage
}

type ifCableResistanceStorage interface {
	GetList(ctx context.Context, pg, pgs int, gs1 string, ord int, dsc bool) (models.CableResistance_count, error)
	Add(ctx context.Context, ea models.CableResistance) (int, error)
	Upd(ctx context.Context, eu models.CableResistance) (int, error)
	Del(ctx context.Context, ed []int) ([]int, error)
	GetOne(ctx context.Context, i int) (models.CableResistance_count, error)
}

//func NewCableResistanceService(storage pgsql.CableResistanceStorage) *CableResistanceService
func NewCableResistanceService(storage pgsql.CableResistanceStorage) *CableResistanceService {
	return &CableResistanceService{storage}
}

//func (esv *CableResistanceService) GetList(ctx context.Context, pg, pgs int, gs1 string, ord int, dsc bool) (models.CableResistance_count, error)
func (esv *CableResistanceService) GetList(ctx context.Context, pg, pgs int, gs1 string, ord int, dsc bool) (models.CableResistance_count, error) {
	var est ifCableResistanceStorage
	est = &esv.storage
	auth := models.Auth{Create: true, Read: true, Update: true, Delete: true}

	out_count, err := est.GetList(ctx, pg, pgs, gs1, ord, dsc)

	if err != nil {
		log.Println("CableResistanceStorage.GetList", err)
		return models.CableResistance_count{Values: []models.CableResistance{}, Count: 0, Auth: auth}, err
	}

	return out_count, nil
}

//func (esv *CableResistanceService) Add(ctx context.Context, ea models.CableResistance) (int, error)
func (esv *CableResistanceService) Add(ctx context.Context, ea models.CableResistance) (int, error) {
	var est ifCableResistanceStorage
	est = &esv.storage

	ai, err := est.Add(ctx, ea)

	if err != nil {
		log.Println("CableResistanceStorage.Add", err)
		return 0, err
	}

	return ai, nil
}

//func (esv *CableResistanceService) Upd(ctx context.Context, eu models.CableResistance) (int, error)
func (esv *CableResistanceService) Upd(ctx context.Context, eu models.CableResistance) (int, error) {
	var est ifCableResistanceStorage
	est = &esv.storage

	ui, err := est.Upd(ctx, eu)

	if err != nil {
		log.Println("CableResistanceStorage.Upd", err)
		return 0, err
	}

	return ui, nil
}

//func (esv *CableResistanceService) Del(ctx context.Context, ed []int) ([]int, error) {
func (esv *CableResistanceService) Del(ctx context.Context, ed []int) ([]int, error) {
	var est ifCableResistanceStorage
	est = &esv.storage

	res, err := est.Del(ctx, ed)

	if err != nil {
		log.Println("CableResistanceStorage.Del", err)
		return []int{}, err
	}

	return res, nil
}

//func (esv *CableResistanceService) GetOne(ctx context.Context, i int) (models.CableResistance_count, error)
func (esv *CableResistanceService) GetOne(ctx context.Context, i int) (models.CableResistance_count, error) {
	var est ifCableResistanceStorage
	est = &esv.storage
	auth := models.Auth{Create: true, Read: true, Update: true, Delete: true}

	out_count, err := est.GetOne(ctx, i)

	if err != nil {
		log.Println("CableResistanceStorage.GetOne", err)
		return models.CableResistance_count{Values: []models.CableResistance{}, Count: 0, Auth: auth}, err
	}

	return out_count, nil

}
