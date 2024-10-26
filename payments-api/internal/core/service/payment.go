package service

import (
	"fmt"
	"log"

	"github.com/jtonynet/go-payments-api/internal/core/domain"
	"github.com/jtonynet/go-payments-api/internal/core/port"
)

type Payment struct {
	AccountRepository     port.AccountRepository
	BalanceRepository     port.BalanceRepository
	TransactionRepository port.TransactionRepository
	MerchantRepository    port.MerchantRepository
}

func NewPayment(
	aRepository port.AccountRepository,
	bRepository port.BalanceRepository,
	tRepository port.TransactionRepository,
	mRepository port.MerchantRepository,
) *Payment {
	return &Payment{
		AccountRepository:     aRepository,
		BalanceRepository:     bRepository,
		TransactionRepository: tRepository,
		MerchantRepository:    mRepository,
	}
}

func (p *Payment) Execute(tpr port.TransactionPaymentRequest) (string, error) {
	accountEntity, err := p.AccountRepository.FindByUID(tpr.AccountUID)
	if err != nil {
		return rejectedGenericErr(fmt.Errorf("failed to retrieve account entity: %w", err))
	}

	account, err := mapAccountEntityToDomain(accountEntity)
	if err != nil {
		return rejectedGenericErr(fmt.Errorf("failed to map account domain from entity: %w", err))
	}

	var merchant domain.Merchant
	merchantEntity, err := p.MerchantRepository.FindByName(tpr.Merchant)
	if err != nil {
		return rejectedGenericErr(fmt.Errorf("failed to retrieve merchant entity with name %s", tpr.Merchant))
	}

	if merchantEntity != nil {
		merchant = mapMerchantEntityToDomain(merchantEntity)
	}

	transaction := account.NewTransaction(
		tpr.MccCode,
		tpr.TotalAmount,
		merchant,
	)

	balanceEntity, err := p.BalanceRepository.FindByAccountID(account.ID)
	if err != nil {
		return rejectedGenericErr(fmt.Errorf("failed to retrieve balance entity: %w", err))
	}

	balance, err := mapBalanceEntityToDomain(balanceEntity)
	if err != nil {
		return rejectedGenericErr(fmt.Errorf("failed to map balance domain from entity: %w", err))
	}

	approvedBalance, cErr := balance.ApproveTransaction(transaction)
	if cErr != nil {
		return rejectedCustomErr(cErr)
	}

	err = p.BalanceRepository.UpdateTotalAmount(mapBalanceDomainToEntity(approvedBalance))
	if err != nil {
		return rejectedGenericErr(fmt.Errorf("failed to update balance entity: %w", err))
	}

	err = p.TransactionRepository.Save(mapTransactionDomainToEntity(transaction))
	if err != nil {
		return rejectedGenericErr(fmt.Errorf("failed to save transaction entity: %w", err))
	}

	return domain.CODE_APPROVED, nil
}

func rejectedGenericErr(err error) (string, error) {
	log.Println(err)
	return domain.CODE_REJECTED_GENERIC, err
}

func rejectedCustomErr(cErr *domain.CustomError) (string, error) {
	log.Println(cErr)
	return cErr.Code, fmt.Errorf("failed to approve balance domain: %s", cErr.Message)
}
