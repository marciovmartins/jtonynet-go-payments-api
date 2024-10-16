package service

import (
	"fmt"

	"github.com/jtonynet/go-payments-api/internal/adapter/repository"
	"github.com/jtonynet/go-payments-api/internal/core/constants"
	"github.com/jtonynet/go-payments-api/internal/core/customError"
	"github.com/jtonynet/go-payments-api/internal/core/domain"
	"github.com/jtonynet/go-payments-api/internal/core/port"
)

type Transaction struct {
	repos repository.Repos
}

func NewTransaction(repos repository.Repos) *Transaction {
	return &Transaction{repos}
}

func (t *Transaction) HandleTransaction(tRequest port.TransactionRequest) (*Transaction, error) {
	tDomain, err := domain.NewTransactionFromRequest(tRequest)
	if err != nil {
		msg := fmt.Sprintf("dont instantiate transaction domain from request: %v", err)
		return &Transaction{}, customError.New(constants.CODE_REJECTED_GENERIC, msg)
	}

	accountRepo := t.repos.Account
	accountEntity, err := accountRepo.FindByUID(tDomain.AccountUID)
	if err != nil {
		msg := fmt.Sprintf("dont retrieve account entity: %v", err)
		return &Transaction{}, customError.New(constants.CODE_REJECTED_GENERIC, msg)
	}

	accountDomain, err := domain.NewAccountFromEntity(accountEntity)
	if err != nil {
		msg := fmt.Sprintf("dont instantiate account domain from entity: %v", err)
		return &Transaction{}, customError.New(constants.CODE_REJECTED_GENERIC, msg)
	}

	balanceRepo := t.repos.Balance
	balanceEntity, err := balanceRepo.FindByAccountID(accountDomain.ID)
	if err != nil {
		msg := fmt.Sprintf("dont retrieve balance entity: %v", err)
		return &Transaction{}, customError.New(constants.CODE_REJECTED_GENERIC, msg)
	}

	balanceDomain, err := domain.NewBalanceFromEntity(balanceEntity)
	if err != nil {
		msg := fmt.Sprintf("dont instantiate balance domain from entity: %v", err)
		return &Transaction{}, customError.New(constants.CODE_REJECTED_GENERIC, msg)
	}

	approvedBalance, err := balanceDomain.Approve(tDomain)
	if err != nil {
		msg := fmt.Sprintf("dont approve balance domain: %v", err)
		return &Transaction{}, customError.New(constants.CODE_REJECTED_GENERIC, msg)
	}

	err = balanceRepo.Update(mapBalanceDomainToEntity(approvedBalance))
	if err != nil {
		msg := fmt.Sprintf("dont update balance entity: %v", err)
		return &Transaction{}, customError.New(constants.CODE_REJECTED_GENERIC, msg)
	}

	transactionRepo := t.repos.Transaction
	err = transactionRepo.Save(mapTransactionDomainToEntity(tDomain))
	if err != nil {
		msg := fmt.Sprintf("dont save transaction entity: %v", err)
		return &Transaction{}, customError.New(constants.CODE_REJECTED_GENERIC, msg)
	}

	return t, nil
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
		AccountUID:  tDomain.AccountUID,
		TotalAmount: tDomain.TotalAmount,
		MCCcode:     tDomain.MCCcode,
		Merchant:    tDomain.Merchant,
	}
}
