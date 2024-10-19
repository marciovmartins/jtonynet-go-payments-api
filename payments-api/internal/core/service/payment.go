package service

import (
	"fmt"

	"github.com/jtonynet/go-payments-api/internal/adapter/repository"
	"github.com/jtonynet/go-payments-api/internal/core/constants"
	"github.com/jtonynet/go-payments-api/internal/core/domain"
	"github.com/jtonynet/go-payments-api/internal/core/port"
)

type Payment struct {
	repos repository.Repos
}

func NewPayment(repos repository.Repos) *Payment {
	return &Payment{repos}
}

func (t *Payment) Execute(tRequest port.TransactionRequest) (string, error) {
	accountRepo := t.repos.Account
	accountEntity, err := accountRepo.FindByUID(tRequest.AccountUID)
	if err != nil {
		return constants.CODE_REJECTED_GENERIC, fmt.Errorf("dont retrieve account entity: %v", err)
	}

	tDomain, err := domain.NewTransaction(
		accountEntity.ID,
		tRequest.AccountUID,
		tRequest.MCCcode,
		tRequest.TotalAmount,
		tRequest.Merchant)
	if err != nil {
		return constants.CODE_REJECTED_GENERIC, fmt.Errorf("dont instantiate transaction domain from request: %v", err)
	}

	accountDomain, err := domain.NewAccountFromEntity(accountEntity)
	if err != nil {
		return constants.CODE_REJECTED_GENERIC, fmt.Errorf("dont instantiate account domain from entity: %v", err)
	}

	balanceRepo := t.repos.Balance
	balanceEntity, err := balanceRepo.FindByAccountID(accountDomain.ID)
	if err != nil {
		return constants.CODE_REJECTED_GENERIC, fmt.Errorf("dont retrieve balance entity: %v", err)
	}

	balanceDomain, err := domain.NewBalanceFromEntity(balanceEntity)
	if err != nil {
		return constants.CODE_REJECTED_GENERIC, fmt.Errorf("dont instantiate balance domain from entity: %v", err)
	}

	approvedBalance, cErr := balanceDomain.Approve(tDomain)
	if cErr != nil {
		return cErr.Code, fmt.Errorf("dont approve balance domain: %v", err)
	}

	err = balanceRepo.UpdateTotalAmount(mapBalanceDomainToEntity(approvedBalance))
	if err != nil {
		return constants.CODE_REJECTED_GENERIC, fmt.Errorf("dont update balance entity: %v", err)
	}

	transactionRepo := t.repos.Transaction
	err = transactionRepo.Save(mapTransactionDomainToEntity(tDomain))
	if err != nil {
		return constants.CODE_REJECTED_GENERIC, fmt.Errorf("dont save transaction entity: %v", err)
	}

	return constants.CODE_APPROVED, nil
}

func mapBalanceDomainToEntity(dBalance domain.Balance) port.BalanceEntity {
	bCategories := make(map[int]port.BalanceByCategoryEntity)
	for _, categoryItem := range dBalance.GetCategories().Itens {
		bCategory := port.BalanceByCategoryEntity{
			ID:        categoryItem.ID,
			AccountID: dBalance.GetAccountID(),
			Amount:    categoryItem.Amount,
			Category:  port.Categories[categoryItem.Name],
		}

		bCategories[categoryItem.Order] = bCategory
	}

	return port.BalanceEntity{
		AccountID:   dBalance.GetAccountID(),
		AmountTotal: dBalance.GetAmountTotal(),
		Categories:  bCategories,
	}
}

func mapTransactionDomainToEntity(tDomain *domain.Transaction) port.TransactionEntity {
	return port.TransactionEntity{
		AccountID:   tDomain.AccountID,
		TotalAmount: tDomain.TotalAmount,
		MCCcode:     tDomain.MCCcode,
		Merchant:    tDomain.Merchant,
	}
}
