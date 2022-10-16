package services

import (
	"context"
	"log"

	"github.com/u2lentaru/billbck/internal/adapters/db/pgsql"
	"github.com/u2lentaru/billbck/internal/models"
)

type ViolationService struct {
	storage pgsql.ViolationStorage
}

type ifViolationStorage interface {
	GetList(ctx context.Context, pg, pgs int, gs1 string, ord int, dsc bool) (models.Violation_count, error)
	Add(ctx context.Context, ea models.Violation) (int, error)
	Upd(ctx context.Context, eu models.Violation) (int, error)
	Del(ctx context.Context, ed []int) ([]int, error)
	GetOne(ctx context.Context, i int) (models.Violation_count, error)
}

//func NewViolationService(storage pgsql.Violationtorage) *ViolationService
func NewViolationService(storage pgsql.ViolationStorage) *ViolationService {
	return &ViolationService{storage}
}

//func (esv *ViolationService) GetList(ctx context.Context, pg, pgs int, gs1 string, ord int, dsc bool) (models.Violation_count, error)
func (esv *ViolationService) GetList(ctx context.Context, pg, pgs int, gs1 string, ord int, dsc bool) (models.Violation_count, error) {
	var est ifViolationStorage
	est = &esv.storage

	out_count, err := est.GetList(ctx, pg, pgs, gs1, ord, dsc)

	if err != nil {
		log.Println("ViolationStorage.GetList", err)
		return models.Violation_count{Values: []models.Violation{}, Count: 0, Auth: models.Auth{}}, err
	}

	return out_count, nil
}

//func (esv *ViolationService) Add(ctx context.Context, ea models.Violation) (int, error)
func (esv *ViolationService) Add(ctx context.Context, ea models.Violation) (int, error) {
	var est ifViolationStorage
	est = &esv.storage

	ai, err := est.Add(ctx, ea)

	if err != nil {
		log.Println("ViolationStorage.Add", err)
		return 0, err
	}

	return ai, nil
}

//func (esv *ViolationService) Upd(ctx context.Context, eu models.Violation) (int, error)
func (esv *ViolationService) Upd(ctx context.Context, eu models.Violation) (int, error) {
	var est ifViolationStorage
	est = &esv.storage

	ui, err := est.Upd(ctx, eu)

	if err != nil {
		log.Println("ViolationStorage.Upd", err)
		return 0, err
	}

	return ui, nil
}

//func (esv *ViolationService) Del(ctx context.Context, ed []int) ([]int, error) {
func (esv *ViolationService) Del(ctx context.Context, ed []int) ([]int, error) {
	var est ifViolationStorage
	est = &esv.storage

	res, err := est.Del(ctx, ed)

	if err != nil {
		log.Println("ViolationStorage.Del", err)
		return []int{}, err
	}

	return res, nil
}

//func (esv *ViolationService) GetOne(ctx context.Context, i int) (models.Violation_count, error)
func (esv *ViolationService) GetOne(ctx context.Context, i int) (models.Violation_count, error) {
	var est ifViolationStorage
	est = &esv.storage

	out_count, err := est.GetOne(ctx, i)

	if err != nil {
		log.Println("ViolationStorage.GetOne", err)
		return models.Violation_count{Values: []models.Violation{}, Count: 0, Auth: models.Auth{}}, err
	}

	return out_count, nil
}
