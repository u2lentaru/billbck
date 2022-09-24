package services

import (
	"context"
	"log"

	"github.com/u2lentaru/billbck/internal/adapters/db/pgsql"
	"github.com/u2lentaru/billbck/internal/models"
)

type ContractMotService struct {
	storage pgsql.ContractMotStorage
}

type ifContractMotStorage interface {
	GetList(ctx context.Context, pg, pgs int, gs1, gs2 string, ord int, dsc bool) (models.ContractMot_count, error)
	Add(ctx context.Context, ea models.ContractMot) (int, error)
	Upd(ctx context.Context, eu models.ContractMot) (int, error)
	Del(ctx context.Context, ed []int) ([]int, error)
	GetOne(ctx context.Context, i int) (models.ContractMot_count, error)
}

//NewContractMotService(storage pg.ContractMotStorage) *ContractMotService
func NewContractMotService(storage pgsql.ContractMotStorage) *ContractMotService {
	return &ContractMotService{storage}
}

//func (esv *ContractMotService) GetList(ctx context.Context, pg, pgs int, gs1, gs2 string, ord int, dsc bool) (models.ContractMot_count, error)
func (esv *ContractMotService) GetList(ctx context.Context, pg, pgs int, gs1, gs2 string, ord int, dsc bool) (models.ContractMot_count, error) {
	var est ifContractMotStorage
	est = &esv.storage

	out_count, err := est.GetList(ctx, pg, pgs, gs1, gs2, ord, dsc)

	if err != nil {
		log.Println("ContractMotStorage.GetList", err)
		return models.ContractMot_count{Values: []models.ContractMot{}, Count: 0, Auth: models.Auth{}}, err
	}

	return out_count, nil
}

//func (esv *ContractMotService) Add(ctx context.Context, ea models.ContractMot) (int, error)
func (esv *ContractMotService) Add(ctx context.Context, ea models.ContractMot) (int, error) {
	var est ifContractMotStorage
	est = &esv.storage

	ai, err := est.Add(ctx, ea)

	if err != nil {
		log.Println("ContractMotStorage.Add", err)
		return 0, err
	}

	return ai, nil
}

//func (esv *ContractMotService) Upd(ctx context.Context, eu models.ContractMot) (int, error)
func (esv *ContractMotService) Upd(ctx context.Context, eu models.ContractMot) (int, error) {
	var est ifContractMotStorage
	est = &esv.storage

	ui, err := est.Upd(ctx, eu)

	if err != nil {
		log.Println("ContractMotStorage.Upd", err)
		return 0, err
	}

	return ui, nil
}

//func (esv *ContractMotService) Del(ctx context.Context, ed []int) ([]int, error)
func (esv *ContractMotService) Del(ctx context.Context, ed []int) ([]int, error) {
	var est ifContractMotStorage
	est = &esv.storage

	res, err := est.Del(ctx, ed)

	if err != nil {
		log.Println("ContractMotStorage.Del", err)
		return []int{}, err
	}

	return res, nil
}

//func (esv *ContractMotService) GetOne(ctx context.Context, i int) (models.ContractMot_count, error)
func (esv *ContractMotService) GetOne(ctx context.Context, i int) (models.ContractMot_count, error) {
	var est ifContractMotStorage
	est = &esv.storage

	out_count, err := est.GetOne(ctx, i)

	if err != nil {
		log.Println("ContractMotStorage.GetOne", err)
		return models.ContractMot_count{Values: []models.ContractMot{}, Count: 0, Auth: models.Auth{}}, err
	}

	return out_count, nil
}
