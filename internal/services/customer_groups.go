package services

import (
	"context"
	"log"

	"github.com/u2lentaru/billbck/internal/adapters/db/pgsql"
	"github.com/u2lentaru/billbck/internal/models"
)

type CustomerGroupService struct {
	storage pgsql.CustomerGroupStorage
}

type ifCustomerGroupStorage interface {
	GetList(ctx context.Context, pg, pgs int, gs1 string, ord int, dsc bool) (models.CustomerGroup_count, error)
	Add(ctx context.Context, ea models.CustomerGroup) (int, error)
	Upd(ctx context.Context, eu models.CustomerGroup) (int, error)
	Del(ctx context.Context, ed []int) ([]int, error)
	GetOne(ctx context.Context, i int) (models.CustomerGroup_count, error)
}

//func NewCustomerGroupService(storage pgsql.CustomerGroupStorage) *CustomerGroupService
func NewCustomerGroupService(storage pgsql.CustomerGroupStorage) *CustomerGroupService {
	return &CustomerGroupService{storage}
}

//func (esv *CustomerGroupService) GetList(ctx context.Context, pg, pgs int, gs1 string, ord int, dsc bool) (models.CustomerGroup_count, error)
func (esv *CustomerGroupService) GetList(ctx context.Context, pg, pgs int, gs1 string, ord int, dsc bool) (models.CustomerGroup_count, error) {
	var est ifCustomerGroupStorage
	est = &esv.storage

	out_count, err := est.GetList(ctx, pg, pgs, gs1, ord, dsc)

	if err != nil {
		log.Println("CustomerGroupStorage.GetList", err)
		return models.CustomerGroup_count{Values: []models.CustomerGroup{}, Count: 0, Auth: models.Auth{}}, err
	}

	return out_count, nil
}

//func (esv *CustomerGroupService) Add(ctx context.Context, ea models.CustomerGroup) (int, error)
func (esv *CustomerGroupService) Add(ctx context.Context, ea models.CustomerGroup) (int, error) {
	var est ifCustomerGroupStorage
	est = &esv.storage

	ai, err := est.Add(ctx, ea)

	if err != nil {
		log.Println("CustomerGroupStorage.Add", err)
		return 0, err
	}

	return ai, nil
}

//func (esv *CustomerGroupService) Upd(ctx context.Context, eu models.CustomerGroup) (int, error)
func (esv *CustomerGroupService) Upd(ctx context.Context, eu models.CustomerGroup) (int, error) {
	var est ifCustomerGroupStorage
	est = &esv.storage

	ui, err := est.Upd(ctx, eu)

	if err != nil {
		log.Println("CustomerGroupStorage.Upd", err)
		return 0, err
	}

	return ui, nil
}

//func (esv *CustomerGroupService) Del(ctx context.Context, ed []int) ([]int, error) {
func (esv *CustomerGroupService) Del(ctx context.Context, ed []int) ([]int, error) {
	var est ifCustomerGroupStorage
	est = &esv.storage

	res, err := est.Del(ctx, ed)

	if err != nil {
		log.Println("CustomerGroupStorage.Del", err)
		return []int{}, err
	}

	return res, nil
}

//func (esv *CustomerGroupService) GetOne(ctx context.Context, i int) (models.CustomerGroup_count, error)
func (esv *CustomerGroupService) GetOne(ctx context.Context, i int) (models.CustomerGroup_count, error) {
	var est ifCustomerGroupStorage
	est = &esv.storage

	out_count, err := est.GetOne(ctx, i)

	if err != nil {
		log.Println("CustomerGroupStorage.GetOne", err)
		return models.CustomerGroup_count{Values: []models.CustomerGroup{}, Count: 0, Auth: models.Auth{}}, err
	}

	return out_count, nil
}
