package service

import (
	"fmt"

	"github.com/jtonynet/go-payments-api/internal/core/constants"
	"github.com/jtonynet/go-payments-api/internal/core/domain"
	"github.com/jtonynet/go-payments-api/internal/core/port"
)

type Payment struct {
	AccountRepository     port.AccountRepository
	BalanceRepository     port.BalanceRepository
	TransactionRepository port.TransactionRepository
}

func NewPayment(
	aRepository port.AccountRepository,
	bRepository port.BalanceRepository,
	tRepository port.TransactionRepository,
) *Payment {
	return &Payment{
		AccountRepository:     aRepository,
		BalanceRepository:     bRepository,
		TransactionRepository: tRepository,
	}
}

func (t *Payment) Execute(tRequest port.TransactionPaymentRequest) (string, error) {
	accountRepo := t.AccountRepository
	accountEntity, err := accountRepo.FindByUID(tRequest.AccountUID)
	if err != nil {
		return rejectedCodeError(fmt.Errorf("dont retrieve account entity: %w", err))
	}

	tDomain, err := domain.NewTransaction(
		accountEntity.ID,
		tRequest.AccountUID,
		tRequest.MCCcode,
		tRequest.TotalAmount,
		tRequest.Merchant)
	if err != nil {
		return rejectedCodeError(fmt.Errorf("dont instantiate transaction domain from request: %w", err))
	}

	accountDomain, err := domain.NewAccountFromEntity(accountEntity)
	if err != nil {
		return rejectedCodeError(fmt.Errorf("dont instantiate account domain from entity: %w", err))
	}

	balanceRepo := t.BalanceRepository
	balanceEntity, err := balanceRepo.FindByAccountID(accountDomain.ID)
	if err != nil {
		return rejectedCodeError(fmt.Errorf("dont retrieve balance entity: %w", err))
	}

	balanceDomain, err := domain.NewBalanceFromEntity(balanceEntity)
	if err != nil {
		return rejectedCodeError(fmt.Errorf("dont instantiate balance domain from entity: %w", err))
	}

	approvedBalance, cErr := balanceDomain.ApproveTransaction(tDomain)
	if cErr != nil {
		return cErr.Code, fmt.Errorf("dont approve balance domain: %w", err)
	}

	err = balanceRepo.UpdateTotalAmount(mapBalanceDomainToEntity(approvedBalance))
	if err != nil {
		return rejectedCodeError(fmt.Errorf("dont update balance entity: %w", err))
	}

	transactionRepo := t.TransactionRepository
	err = transactionRepo.Save(mapTransactionDomainToEntity(tDomain))
	if err != nil {
		return rejectedCodeError(fmt.Errorf("dont save transaction entity: %w", err))
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

func rejectedCodeError(err error) (string, error) {
	return constants.CODE_REJECTED_GENERIC, err
}
