package services

import (
	"context"
	"log"

	"github.com/u2lentaru/billbck/internal/adapters/db/pgsql"
	"github.com/u2lentaru/billbck/internal/models"
)

type SubBankService struct {
	storage pgsql.SubBankStorage
}

type ifSubBankStorage interface {
	GetList(ctx context.Context, pg, pgs int, gs1 string, gs2 int, gs3 string, ord int, dsc bool) (models.SubBank_count, error)
	Add(ctx context.Context, ea models.SubBank) (int, error)
	Upd(ctx context.Context, eu models.SubBank) (int, error)
	Del(ctx context.Context, ed []int) ([]int, error)
	GetOne(ctx context.Context, i int) (models.SubBank_count, error)
	SetActive(ctx context.Context, i int) (int, error)
}

//func NewSubBankService(storage pgsql.SubBanktorage) *SubBankService
func NewSubBankService(storage pgsql.SubBankStorage) *SubBankService {
	return &SubBankService{storage}
}

//func (esv *SubBankService) GetList(ctx context.Context, pg, pgs int, gs1 string, gs2 int, gs3 string, ord int, dsc bool) (models.SubBank_count, error)
func (esv *SubBankService) GetList(ctx context.Context, pg, pgs int, gs1 string, gs2 int, gs3 string, ord int, dsc bool) (models.SubBank_count, error) {
	var est ifSubBankStorage
	est = &esv.storage

	out_count, err := est.GetList(ctx, pg, pgs, gs1, gs2, gs3, ord, dsc)

	if err != nil {
		log.Println("SubBankStorage.GetList", err)
		return models.SubBank_count{Values: []models.SubBank{}, Count: 0, Auth: models.Auth{}}, err
	}

	return out_count, nil
}

//func (esv *SubBankService) Add(ctx context.Context, ea models.SubBank) (int, error)
func (esv *SubBankService) Add(ctx context.Context, ea models.SubBank) (int, error) {
	var est ifSubBankStorage
	est = &esv.storage

	ai, err := est.Add(ctx, ea)

	if err != nil {
		log.Println("SubBankStorage.Add", err)
		return 0, err
	}

	return ai, nil
}

//func (esv *SubBankService) Upd(ctx context.Context, eu models.SubBank) (int, error)
func (esv *SubBankService) Upd(ctx context.Context, eu models.SubBank) (int, error) {
	var est ifSubBankStorage
	est = &esv.storage

	ui, err := est.Upd(ctx, eu)

	if err != nil {
		log.Println("SubBankStorage.Upd", err)
		return 0, err
	}

	return ui, nil
}

//func (esv *SubBankService) Del(ctx context.Context, ed []int) ([]int, error) {
func (esv *SubBankService) Del(ctx context.Context, ed []int) ([]int, error) {
	var est ifSubBankStorage
	est = &esv.storage

	res, err := est.Del(ctx, ed)

	if err != nil {
		log.Println("SubBankStorage.Del", err)
		return []int{}, err
	}

	return res, nil
}

//func (esv *SubBankService) GetOne(ctx context.Context, i int) (models.SubBank_count, error)
func (esv *SubBankService) GetOne(ctx context.Context, i int) (models.SubBank_count, error) {
	var est ifSubBankStorage
	est = &esv.storage

	out_count, err := est.GetOne(ctx, i)

	if err != nil {
		log.Println("SubBankStorage.GetOne", err)
		return models.SubBank_count{Values: []models.SubBank{}, Count: 0, Auth: models.Auth{}}, err
	}

	return out_count, nil
}

//func (esv *SubBankService) SetActive(ctx context.Context, i int) (int, error)
func (esv *SubBankService) SetActive(ctx context.Context, i int) (int, error) {
	var est ifSubBankStorage
	est = &esv.storage

	si, err := est.SetActive(ctx, i)

	if err != nil {
		log.Println("SubBankStorage.SetActive", err)
		return 0, err
	}

	return si, nil
}
