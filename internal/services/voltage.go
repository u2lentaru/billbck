package services

import (
	"context"
	"log"

	"github.com/u2lentaru/billbck/internal/adapters/db/pgsql"
	"github.com/u2lentaru/billbck/internal/models"
)

type VoltageService struct {
	storage pgsql.VoltageStorage
}

type ifVoltageStorage interface {
	GetList(ctx context.Context, pg, pgs int, gs1 string, ord int, dsc bool) (models.Voltage_count, error)
	Add(ctx context.Context, ea models.Voltage) (int, error)
	Upd(ctx context.Context, eu models.Voltage) (int, error)
	Del(ctx context.Context, ed []int) ([]int, error)
	GetOne(ctx context.Context, i int) (models.Voltage_count, error)
}

//func NewVoltageService(storage pgsql.Voltagetorage) *VoltageService
func NewVoltageService(storage pgsql.VoltageStorage) *VoltageService {
	return &VoltageService{storage}
}

//func (esv *VoltageService) GetList(ctx context.Context, pg, pgs int, gs1 string, ord int, dsc bool) (models.Voltage_count, error)
func (esv *VoltageService) GetList(ctx context.Context, pg, pgs int, gs1 string, ord int, dsc bool) (models.Voltage_count, error) {
	var est ifVoltageStorage
	est = &esv.storage

	out_count, err := est.GetList(ctx, pg, pgs, gs1, ord, dsc)

	if err != nil {
		log.Println("VoltageStorage.GetList", err)
		return models.Voltage_count{Values: []models.Voltage{}, Count: 0, Auth: models.Auth{}}, err
	}

	return out_count, nil
}

//func (esv *VoltageService) Add(ctx context.Context, ea models.Voltage) (int, error)
func (esv *VoltageService) Add(ctx context.Context, ea models.Voltage) (int, error) {
	var est ifVoltageStorage
	est = &esv.storage

	ai, err := est.Add(ctx, ea)

	if err != nil {
		log.Println("VoltageStorage.Add", err)
		return 0, err
	}

	return ai, nil
}

//func (esv *VoltageService) Upd(ctx context.Context, eu models.Voltage) (int, error)
func (esv *VoltageService) Upd(ctx context.Context, eu models.Voltage) (int, error) {
	var est ifVoltageStorage
	est = &esv.storage

	ui, err := est.Upd(ctx, eu)

	if err != nil {
		log.Println("VoltageStorage.Upd", err)
		return 0, err
	}

	return ui, nil
}

//func (esv *VoltageService) Del(ctx context.Context, ed []int) ([]int, error) {
func (esv *VoltageService) Del(ctx context.Context, ed []int) ([]int, error) {
	var est ifVoltageStorage
	est = &esv.storage

	res, err := est.Del(ctx, ed)

	if err != nil {
		log.Println("VoltageStorage.Del", err)
		return []int{}, err
	}

	return res, nil
}

//func (esv *VoltageService) GetOne(ctx context.Context, i int) (models.Voltage_count, error)
func (esv *VoltageService) GetOne(ctx context.Context, i int) (models.Voltage_count, error) {
	var est ifVoltageStorage
	est = &esv.storage

	out_count, err := est.GetOne(ctx, i)

	if err != nil {
		log.Println("VoltageStorage.GetOne", err)
		return models.Voltage_count{Values: []models.Voltage{}, Count: 0, Auth: models.Auth{}}, err
	}

	return out_count, nil
}
