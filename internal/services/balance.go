package services

import (
	"context"
	"log"

	"github.com/u2lentaru/billbck/internal/adapters/db/pgsql"
	"github.com/u2lentaru/billbck/internal/models"
)

type BalanceService struct {
	storage pgsql.BalanceStorage
}

type ifBalanceStorage interface {
	GetList(ctx context.Context, pg, pgs, gs1, gs2 int) (models.Balance_count, error)
	GetNode(ctx context.Context, gs1, gs2 string) (models.Balance, error)
	GetNodeSum(ctx context.Context, gs1, gs2 int, gs3, gs4 string) (models.Json_sum, error)
	GetNodeSumL1(ctx context.Context, gs1, gs2 int, gs3, gs4 string) (models.Json_sum, error)
	GetNodeSumL0(ctx context.Context, gs1, gs2 int, gs3, gs4 string) (models.Json_sum, error)
	GetTabL1(ctx context.Context, pg, pgs, gs1, gs2 int, gs3, gs4 string) (models.BalanceTab_sum, error)
	GetTabL0(ctx context.Context, pg, pgs, gs1, gs2 int, gs3, gs4 string) (models.BalanceTab_sum, error)
	GetBranch(ctx context.Context, gs1, gs2 int) (models.BalanceTab_sum, error)
}

//NewBalanceService(storage pg.BalanceStorage) *BalanceService
func NewBalanceService(storage pgsql.BalanceStorage) *BalanceService {
	return &BalanceService{storage}
}

//func (esv *BalanceService) GetList(ctx context.Context, pg, pgs, gs1, gs2 int) (models.Balance_count, error)
func (esv *BalanceService) GetList(ctx context.Context, pg, pgs, gs1, gs2 int) (models.Balance_count, error) {
	var est ifBalanceStorage
	est = &esv.storage

	out_count, err := est.GetList(ctx, pg, pgs, gs1, gs2)

	if err != nil {
		log.Println("BalanceStorage.GetList", err)
		return models.Balance_count{Values: []models.Balance{}, Count: 0}, err
	}

	return out_count, nil
}

//func (esv *BalanceService) GetNode(ctx context.Context, gs1, gs2 string) (models.Balance, error)
func (esv *BalanceService) GetNode(ctx context.Context, gs1, gs2 string) (models.Balance, error) {
	var est ifBalanceStorage
	est = &esv.storage

	out_count, err := est.GetNode(ctx, gs1, gs2)

	if err != nil {
		log.Println("BalanceStorage.GetNode", err)
		return models.Balance{}, err
	}

	return out_count, nil
}

//func (esv *BalanceService) GetNodeSum(ctx context.Context, gs1, gs2 int, gs3, gs4 string) (models.Json_sum, error)
func (esv *BalanceService) GetNodeSum(ctx context.Context, gs1, gs2 int, gs3, gs4 string) (models.Json_sum, error) {
	var est ifBalanceStorage
	est = &esv.storage

	out_count, err := est.GetNodeSum(ctx, gs1, gs2, gs3, gs4)

	if err != nil {
		log.Println("BalanceStorage.GetNodeSum", err)
		return models.Json_sum{}, err
	}

	return out_count, nil
}

//func (esv *BalanceService) GetNodeSumL1(ctx context.Context, gs1, gs2 int, gs3, gs4 string) (models.Json_sum, error)
func (esv *BalanceService) GetNodeSumL1(ctx context.Context, gs1, gs2 int, gs3, gs4 string) (models.Json_sum, error) {
	var est ifBalanceStorage
	est = &esv.storage

	out_count, err := est.GetNodeSumL1(ctx, gs1, gs2, gs3, gs4)

	if err != nil {
		log.Println("BalanceStorage.GetNodeSumL1", err)
		return models.Json_sum{}, err
	}

	return out_count, nil
}

//func (esv *BalanceService) GetNodeSumL0(ctx context.Context, gs1, gs2 int, gs3, gs4 string) (models.Json_sum, error)
func (esv *BalanceService) GetNodeSumL0(ctx context.Context, gs1, gs2 int, gs3, gs4 string) (models.Json_sum, error) {
	var est ifBalanceStorage
	est = &esv.storage

	out_count, err := est.GetNodeSumL0(ctx, gs1, gs2, gs3, gs4)

	if err != nil {
		log.Println("BalanceStorage.GetNodeSumL0", err)
		return models.Json_sum{}, err
	}

	return out_count, nil
}

//func (esv *BalanceService) GetTabL1(ctx context.Context, pg, pgs, gs1, gs2 int, gs3, gs4 string) (models.BalanceTab_sum, error)
func (esv *BalanceService) GetTabL1(ctx context.Context, pg, pgs, gs1, gs2 int, gs3, gs4 string) (models.BalanceTab_sum, error) {
	var est ifBalanceStorage
	est = &esv.storage

	out_count, err := est.GetTabL1(ctx, pg, pgs, gs1, gs2, gs3, gs4)

	if err != nil {
		log.Println("BalanceStorage.GetTabL1", err)
		return models.BalanceTab_sum{}, err
	}

	return out_count, nil
}

//func (esv *BalanceService) GetTabL0(ctx context.Context, pg, pgs, gs1, gs2 int, gs3, gs4 string) (models.BalanceTab_sum, error)
func (esv *BalanceService) GetTabL0(ctx context.Context, pg, pgs, gs1, gs2 int, gs3, gs4 string) (models.BalanceTab_sum, error) {
	var est ifBalanceStorage
	est = &esv.storage

	out_count, err := est.GetTabL0(ctx, pg, pgs, gs1, gs2, gs3, gs4)

	if err != nil {
		log.Println("BalanceStorage.GetTabL0", err)
		return models.BalanceTab_sum{}, err
	}

	return out_count, nil
}

//func (esv *BalanceService) GetBranch(ctx context.Context, gs1, gs2 int) (models.BalanceTab_sum, error)
func (esv *BalanceService) GetBranch(ctx context.Context, gs1, gs2 int) (models.BalanceTab_sum, error) {
	var est ifBalanceStorage
	est = &esv.storage

	out_count, err := est.GetBranch(ctx, gs1, gs2)

	if err != nil {
		log.Println("BalanceStorage.GetBranch", err)
		return models.BalanceTab_sum{}, err
	}

	return out_count, nil
}
