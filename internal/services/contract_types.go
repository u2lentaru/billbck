package services

import (
	"context"
	"log"

	"github.com/u2lentaru/billbck/internal/adapters/db/pgsql"
	"github.com/u2lentaru/billbck/internal/models"
)

type ContractTypeService struct {
	storage pgsql.ContractTypeStorage
}

type ifContractTypeStorage interface {
	GetList(ctx context.Context, pg, pgs int, gs1 string, ord int, dsc bool) (models.ContractType_count, error)
	Add(ctx context.Context, ea models.ContractType) (int, error)
	Upd(ctx context.Context, eu models.ContractType) (int, error)
	Del(ctx context.Context, ed []int) ([]int, error)
	GetOne(ctx context.Context, i int) (models.ContractType_count, error)
}

//func NewContractTypeService(storage pgsql.ContractTypeStorage) *ContractTypeService
func NewContractTypeService(storage pgsql.ContractTypeStorage) *ContractTypeService {
	return &ContractTypeService{storage}
}

//func (esv *ContractTypeService) GetList(ctx context.Context, pg, pgs int, gs1 string, ord int, dsc bool) (models.ContractType_count, error)
func (esv *ContractTypeService) GetList(ctx context.Context, pg, pgs int, gs1 string, ord int, dsc bool) (models.ContractType_count, error) {
	var est ifContractTypeStorage
	est = &esv.storage

	out_count, err := est.GetList(ctx, pg, pgs, gs1, ord, dsc)

	if err != nil {
		log.Println("ContractTypeStorage.GetList", err)
		return models.ContractType_count{Values: []models.ContractType{}, Count: 0, Auth: models.Auth{}}, err
	}

	return out_count, nil
}

//func (esv *ContractTypeService) Add(ctx context.Context, ea models.ContractType) (int, error)
func (esv *ContractTypeService) Add(ctx context.Context, ea models.ContractType) (int, error) {
	var est ifContractTypeStorage
	est = &esv.storage

	ai, err := est.Add(ctx, ea)

	if err != nil {
		log.Println("ContractTypeStorage.Add", err)
		return 0, err
	}

	return ai, nil
}

//func (esv *ContractTypeService) Upd(ctx context.Context, eu models.ContractType) (int, error)
func (esv *ContractTypeService) Upd(ctx context.Context, eu models.ContractType) (int, error) {
	var est ifContractTypeStorage
	est = &esv.storage

	ui, err := est.Upd(ctx, eu)

	if err != nil {
		log.Println("ContractTypeStorage.Upd", err)
		return 0, err
	}

	return ui, nil
}

//func (esv *ContractTypeService) Del(ctx context.Context, ed []int) ([]int, error) {
func (esv *ContractTypeService) Del(ctx context.Context, ed []int) ([]int, error) {
	var est ifContractTypeStorage
	est = &esv.storage

	res, err := est.Del(ctx, ed)

	if err != nil {
		log.Println("ContractTypeStorage.Del", err)
		return []int{}, err
	}

	return res, nil
}

//func (esv *ContractTypeService) GetOne(ctx context.Context, i int) (models.ContractType_count, error)
func (esv *ContractTypeService) GetOne(ctx context.Context, i int) (models.ContractType_count, error) {
	var est ifContractTypeStorage
	est = &esv.storage

	out_count, err := est.GetOne(ctx, i)

	if err != nil {
		log.Println("ContractTypeStorage.GetOne", err)
		return models.ContractType_count{Values: []models.ContractType{}, Count: 0, Auth: models.Auth{}}, err
	}

	return out_count, nil
}
