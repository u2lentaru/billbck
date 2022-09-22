package services

import (
	"context"
	"log"

	"github.com/u2lentaru/billbck/internal/adapters/db/pgsql"
	"github.com/u2lentaru/billbck/internal/models"
)

type BankService struct {
	storage pgsql.BankStorage
}

type ifBankStorage interface {
	GetList(ctx context.Context, pg, pgs int, gs1, gs2 string, ord int, dsc bool) (models.Bank_count, error)
	Add(ctx context.Context, ea models.Bank) (int, error)
	Upd(ctx context.Context, eu models.Bank) (int, error)
	Del(ctx context.Context, ed []int) ([]int, error)
	GetOne(ctx context.Context, i int) (models.Bank_count, error)
}

//NewBankService(storage pg.BankStorage) *BankService
func NewBankService(storage pgsql.BankStorage) *BankService {
	return &BankService{storage}
}

//func (esv *BankService) GetList(ctx context.Context, pg, pgs int, gs1, gs2 string, ord int, dsc bool) (models.Bank_count, error)
func (esv *BankService) GetList(ctx context.Context, pg, pgs int, gs1, gs2 string, ord int, dsc bool) (models.Bank_count, error) {
	var est ifBankStorage
	est = &esv.storage

	out_count, err := est.GetList(ctx, pg, pgs, gs1, gs2, ord, dsc)

	if err != nil {
		log.Println("BankStorage.GetList", err)
		return models.Bank_count{Values: []models.Bank{}, Count: 0, Auth: models.Auth{}}, err
	}

	return out_count, nil
}

//func (esv *BankService) Add(ctx context.Context, ea models.Bank) (int, error)
func (esv *BankService) Add(ctx context.Context, ea models.Bank) (int, error) {
	var est ifBankStorage
	est = &esv.storage

	ai, err := est.Add(ctx, ea)

	if err != nil {
		log.Println("BankStorage.Add", err)
		return 0, err
	}

	return ai, nil
}

//func (esv *BankService) Upd(ctx context.Context, eu models.Bank) (int, error)
func (esv *BankService) Upd(ctx context.Context, eu models.Bank) (int, error) {
	var est ifBankStorage
	est = &esv.storage

	ui, err := est.Upd(ctx, eu)

	if err != nil {
		log.Println("BankStorage.Upd", err)
		return 0, err
	}

	return ui, nil
}

//func (esv *BankService) Del(ctx context.Context, ed []int) ([]int, error)
func (esv *BankService) Del(ctx context.Context, ed []int) ([]int, error) {
	var est ifBankStorage
	est = &esv.storage

	res, err := est.Del(ctx, ed)

	if err != nil {
		log.Println("BankStorage.Del", err)
		return []int{}, err
	}

	return res, nil
}

//func (esv *BankService) GetOne(ctx context.Context, i int) (models.Bank_count, error)
func (esv *BankService) GetOne(ctx context.Context, i int) (models.Bank_count, error) {
	var est ifBankStorage
	est = &esv.storage

	out_count, err := est.GetOne(ctx, i)

	if err != nil {
		log.Println("BankStorage.GetOne", err)
		return models.Bank_count{Values: []models.Bank{}, Count: 0, Auth: models.Auth{}}, err
	}

	return out_count, nil

}
