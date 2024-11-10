package service

import (
	"fmt"

	"github.com/jtonynet/go-payments-api/internal/core/domain"
	"github.com/jtonynet/go-payments-api/internal/core/port"

	"github.com/jtonynet/go-payments-api/internal/support"
)

type Payment struct {
	accountRepository     port.AccountRepository
	balanceRepository     port.BalanceRepository
	transactionRepository port.TransactionRepository
	merchantRepository    port.MerchantRepository

	logger support.Logger
}

func NewPayment(
	aRepository port.AccountRepository,
	bRepository port.BalanceRepository,
	tRepository port.TransactionRepository,
	mRepository port.MerchantRepository,

	logger support.Logger,
) *Payment {
	return &Payment{
		accountRepository:     aRepository,
		balanceRepository:     bRepository,
		transactionRepository: tRepository,
		merchantRepository:    mRepository,

		logger: logger,
	}
}

func (p *Payment) Execute(tpr port.TransactionPaymentRequest) (string, error) {
	accountEntity, err := p.accountRepository.FindByUID(tpr.AccountUID)
	if err != nil {
		return p.rejectedGenericErr(fmt.Errorf("failed to retrieve account entity: %w", err))
	}

	account, err := mapAccountEntityToDomain(accountEntity)
	if err != nil {
		return p.rejectedGenericErr(fmt.Errorf("failed to map account domain from entity: %w", err))
	}

	var merchant domain.Merchant
	merchantEntity, err := p.merchantRepository.FindByName(tpr.Merchant)
	if err != nil {
		return p.rejectedGenericErr(fmt.Errorf("failed to retrieve merchant entity with name %s", tpr.Merchant))
	}

	if merchantEntity != nil {
		merchant = mapMerchantEntityToDomain(merchantEntity)
	}

	transaction := merchant.NewTransaction(
		tpr.MCC,
		tpr.TotalAmount,
		tpr.Merchant,
		account,
	)

	balanceEntity, err := p.balanceRepository.FindByAccountID(account.ID)
	if err != nil {
		return p.rejectedGenericErr(fmt.Errorf("failed to retrieve balance entity: %w", err))
	}

	balance, err := mapBalanceEntityToDomain(balanceEntity)
	if err != nil {
		return p.rejectedGenericErr(fmt.Errorf("failed to map balance domain from entity: %w", err))
	}
	balance.Logger = p.logger

	approvedBalance, cErr := balance.ApproveTransaction(transaction)
	if cErr != nil {
		return p.rejectedCustomErr(cErr)
	}

	err = p.balanceRepository.UpdateTotalAmount(mapBalanceDomainToEntity(approvedBalance))
	if err != nil {
		return p.rejectedGenericErr(fmt.Errorf("failed to update balance entity: %w", err))
	}

	err = p.transactionRepository.Save(mapTransactionDomainToEntity(transaction))
	if err != nil {
		return p.rejectedGenericErr(fmt.Errorf("failed to save transaction entity: %w", err))
	}

	return domain.CODE_APPROVED, nil
}

func (p *Payment) rejectedGenericErr(err error) (string, error) {
	if p.logger != nil {
		p.logger.Debug(err.Error())
	}

	return domain.CODE_REJECTED_GENERIC, err
}

func (p *Payment) rejectedCustomErr(cErr *domain.CustomError) (string, error) {
	if p.logger != nil {
		p.logger.Debug(cErr.Error())
	}

	return cErr.Code, fmt.Errorf("failed to approve balance domain: %s", cErr.Message)
}
