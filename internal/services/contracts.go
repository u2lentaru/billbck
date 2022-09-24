package services

import (
	"context"
	"log"

	"github.com/u2lentaru/billbck/internal/adapters/db/pgsql"
	"github.com/u2lentaru/billbck/internal/models"
)

type ContractService struct {
	storage pgsql.ContractStorage
}

type ifContractStorage interface {
	GetList(ctx context.Context, pg, pgs int, gs1, gs2, gs3, gs4 string, gs5, gs6 int, gs7, gs8, gs9, gs10, gs11, gs12, gs13, gs14 string, ord int, dsc bool) (models.Contract_count, error)
	Add(ctx context.Context, ea models.Contract) (int, error)
	Upd(ctx context.Context, eu models.Contract) (int, error)
	Del(ctx context.Context, ed models.IdClose) (int, error)
	GetOne(ctx context.Context, i int) (models.Contract_count, error)
	GetObj(ctx context.Context, i int, a string) (models.ObjContract, error)
	GetHist(ctx context.Context, i int) (string, error)
}

//NewContractService(storage pg.ContractStorage) *ContractService
func NewContractService(storage pgsql.ContractStorage) *ContractService {
	return &ContractService{storage}
}

//func (esv *ContractService) GetList(ctx context.Context, pg, pgs int, gs1, gs2, gs3, gs4 string, gs5, gs6 int, gs7, gs8, gs9, gs10, gs11, gs12, gs13, gs14 string, ord int, dsc bool) (models.Contract_count, error)
func (esv *ContractService) GetList(ctx context.Context, pg, pgs int, gs1, gs2, gs3, gs4 string, gs5, gs6 int, gs7, gs8, gs9, gs10, gs11, gs12, gs13, gs14 string, ord int, dsc bool) (models.Contract_count, error) {
	var est ifContractStorage
	est = &esv.storage

	out_count, err := est.GetList(ctx, pg, pgs, gs1, gs2, gs3, gs4, gs5, gs6, gs7, gs8, gs9, gs10, gs11, gs12, gs13, gs14, ord, dsc)

	if err != nil {
		log.Println("ContractStorage.GetList", err)
		return models.Contract_count{Values: []models.Contract{}, Count: 0, Auth: models.Auth{}}, err
	}

	return out_count, nil
}

//func (esv *ContractService) Add(ctx context.Context, ea models.Contract) (int, error)
func (esv *ContractService) Add(ctx context.Context, ea models.Contract) (int, error) {
	var est ifContractStorage
	est = &esv.storage

	ai, err := est.Add(ctx, ea)

	if err != nil {
		log.Println("ContractStorage.Add", err)
		return 0, err
	}

	return ai, nil
}

//func (esv *ContractService) Upd(ctx context.Context, eu models.Contract) (int, error)
func (esv *ContractService) Upd(ctx context.Context, eu models.Contract) (int, error) {
	var est ifContractStorage
	est = &esv.storage

	ui, err := est.Upd(ctx, eu)

	if err != nil {
		log.Println("ContractStorage.Upd", err)
		return 0, err
	}

	return ui, nil
}

//func (esv *ContractService) Del(ctx context.Context, ed []int) (int, error)
func (esv *ContractService) Del(ctx context.Context, ed models.IdClose) (int, error) {
	var est ifContractStorage
	est = &esv.storage

	res, err := est.Del(ctx, ed)

	if err != nil {
		log.Println("ContractStorage.Del", err)
		return 0, err
	}

	return res, nil
}

//func (esv *ContractService) GetOne(ctx context.Context, i int) (models.Contract_count, error)
func (esv *ContractService) GetOne(ctx context.Context, i int) (models.Contract_count, error) {
	var est ifContractStorage
	est = &esv.storage

	out_count, err := est.GetOne(ctx, i)

	if err != nil {
		log.Println("ContractStorage.GetOne", err)
		return models.Contract_count{Values: []models.Contract{}, Count: 0, Auth: models.Auth{}}, err
	}

	return out_count, nil
}

//func (esv *ContractService) GetObj(ctx context.Context, i int, s string) (models.Contract_count, error)
func (esv *ContractService) GetObj(ctx context.Context, i int, s string) (models.ObjContract, error) {
	var est ifContractStorage
	est = &esv.storage

	out_count, err := est.GetObj(ctx, i, s)

	if err != nil {
		log.Println("ContractStorage.GetObj", err)
		return models.ObjContract{}, err
	}

	return out_count, nil
}

//func (esv *ContractService) GetHist(ctx context.Context, i int) (string, error)
func (esv *ContractService) GetHist(ctx context.Context, i int) (string, error) {
	var est ifContractStorage
	est = &esv.storage

	out_count, err := est.GetHist(ctx, i)

	if err != nil {
		log.Println("ContractStorage.GetHist", err)
		return "", err
	}

	return out_count, nil
}
