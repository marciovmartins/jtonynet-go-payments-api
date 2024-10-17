package service

import (
	"fmt"

	"github.com/jtonynet/go-payments-api/internal/adapter/repository"
	"github.com/jtonynet/go-payments-api/internal/core/constants"
	"github.com/jtonynet/go-payments-api/internal/core/customError"
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
		msg := fmt.Sprintf("dont retrieve account entity: %v", err)
		return constants.CODE_REJECTED_GENERIC,
			customError.New(constants.CODE_REJECTED_GENERIC, msg)
	}

	tDomain, err := domain.NewTransaction(
		accountEntity.ID,
		tRequest.AccountUID,
		tRequest.MCCcode,
		tRequest.TotalAmount,
		tRequest.Merchant)
	if err != nil {
		msg := fmt.Sprintf("dont instantiate transaction domain from request: %v", err)
		return constants.CODE_REJECTED_GENERIC,
			customError.New(constants.CODE_REJECTED_GENERIC, msg)
	}

	accountDomain, err := domain.NewAccountFromEntity(accountEntity)
	if err != nil {
		msg := fmt.Sprintf("dont instantiate account domain from entity: %v", err)
		return constants.CODE_REJECTED_GENERIC,
			customError.New(constants.CODE_REJECTED_GENERIC, msg)
	}

	balanceRepo := t.repos.Balance
	balanceEntity, err := balanceRepo.FindByAccountID(accountDomain.ID)
	if err != nil {
		msg := fmt.Sprintf("dont retrieve balance entity: %v", err)
		return constants.CODE_REJECTED_GENERIC,
			customError.New(constants.CODE_REJECTED_GENERIC, msg)
	}

	balanceDomain, err := domain.NewBalanceFromEntity(balanceEntity)
	if err != nil {
		msg := fmt.Sprintf("dont instantiate balance domain from entity: %v", err)
		return constants.CODE_REJECTED_GENERIC,
			customError.New(constants.CODE_REJECTED_GENERIC, msg)
	}

	approvedBalance, cErr := balanceDomain.Approve(tDomain)
	if cErr != nil {
		msg := fmt.Sprintf("dont approve balance domain: %v", err)
		return cErr.Code,
			customError.New(cErr.Code, msg)
	}

	err = balanceRepo.Update(mapBalanceDomainToEntity(approvedBalance))
	if err != nil {
		msg := fmt.Sprintf("dont update balance entity: %v", err)
		return constants.CODE_REJECTED_GENERIC,
			customError.New(constants.CODE_REJECTED_GENERIC, msg)
	}

	transactionRepo := t.repos.Transaction
	err = transactionRepo.Save(mapTransactionDomainToEntity(tDomain))
	if err != nil {
		msg := fmt.Sprintf("dont save transaction entity: %v", err)
		return constants.CODE_REJECTED_GENERIC,
			customError.New(constants.CODE_REJECTED_GENERIC, msg)
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
