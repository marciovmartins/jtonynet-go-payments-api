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

func (p *Payment) Execute(tRequest port.TransactionPaymentRequest) (string, error) {
	accountEntity, err := p.AccountRepository.FindByUID(tRequest.AccountUID)
	if err != nil {
		return rejectedCodeError(fmt.Errorf("failed to retrieve account entity: %w", err))
	}

	mDomain := domain.Merchant{}
	merchantEntity, err := p.MerchantRepository.FindByName(tRequest.Merchant)
	if err == nil {
		mDomain = mapMerchantEntityToDomain(merchantEntity)
	}

	tDomain := accountEntity.NewTransaction(
		tRequest.MccCode,
		tRequest.TotalAmount,
		mDomain,
	)

	accountDomain, err := mapAccountEntityToDomain(accountEntity)
	if err != nil {
		return rejectedCodeError(fmt.Errorf("failed to map account domain from entity: %w", err))
	}

	balanceEntity, err := p.BalanceRepository.FindByAccountID(accountDomain.ID)
	if err != nil {
		return rejectedCodeError(fmt.Errorf("failed to retrieve balance entity: %w", err))
	}

	balanceDomain, err := mapBalanceEntityToDomain(balanceEntity)
	if err != nil {
		return rejectedCodeError(fmt.Errorf("failed to map balance domain from entity: %w", err))
	}

	approvedBalance, cErr := balanceDomain.ApproveTransaction(tDomain)
	if cErr != nil {
		log.Println(cErr)
		return cErr.Code, fmt.Errorf("failed to approve balance domain: %s", cErr.Message)
	}

	err = p.BalanceRepository.UpdateTotalAmount(mapBalanceDomainToEntity(approvedBalance))
	if err != nil {
		return rejectedCodeError(fmt.Errorf("failed to update balance entity: %w", err))
	}

	err = p.TransactionRepository.Save(mapTransactionDomainToEntity(tDomain))
	if err != nil {
		return rejectedCodeError(fmt.Errorf("failed to save transaction entity: %w", err))
	}

	return domain.CODE_APPROVED, nil
}

func rejectedCodeError(err error) (string, error) {
	log.Println(err)
	return domain.CODE_REJECTED_GENERIC, err
}
