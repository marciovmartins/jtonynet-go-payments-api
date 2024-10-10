package service

import (
	"fmt"

	"github.com/jtonynet/go-payments-api/internal/adapter/repository"
	"github.com/jtonynet/go-payments-api/internal/core/domain"
	"github.com/jtonynet/go-payments-api/internal/core/port"
)

type Transaction struct {
	repos repository.Repos
}

func NewTransction(repos repository.Repos) Transaction {
	return Transaction{repos}
}

func (t *Transaction) HandleTransactionL1(tDTOHandler port.TransactionDTORequest) (string, error) {
	tDomain, _ := mapTransactionDTOtoDomain(tDTOHandler)

	accountRepo := t.repos.Account
	accountDTORepository, err := accountRepo.FindByUUID(tDomain.AccountUUID)
	if err != nil {
		//TODO implements error in future
		return "todo KCT", fmt.Errorf("dont retrieve account: %v", err)
	}
	aDomain, _ := mapAccountDTOtoDomain(accountDTORepository)

	balanceRepo := t.repos.Balance
	balanceDTORepository, err := balanceRepo.FindByAccountUUID(aDomain.UUID)
	if err != nil {
		//TODO implements error in future
		return "todo KCT", fmt.Errorf("dont retrieve account: %v", err)
	}
	bDomain, _ := mapBalanceDTOtoDomain(balanceDTORepository)

	approvedBalance, _ := bDomain.Approve(tDomain)
	// TODO TRATAR ERROS

	err = balanceRepo.Update(mapBalanceDomainToEntity(approvedBalance))
	if err != nil {
		//TODO implements error in future
		return "todo KCT", fmt.Errorf("dont retrieve account: %v", err)
	}

	transactionRepo := t.repos.Transaction
	err = transactionRepo.Save(mapTransactionDomainToEntity(tDomain))
	if err != nil {
		//TODO implements error in future
		return "todo KCT", fmt.Errorf("dont retrieve account: %v", err)
	}

	fmt.Println(transactionRepo)

	return "00", nil
}

func mapTransactionDomainToEntity(tDomain domain.Transaction) port.TransactionEntity {
	return port.TransactionEntity{}
}

func mapBalanceDomainToEntity(approvedBalance domain.Balance) port.BalanceDTORepository {
	return port.BalanceDTORepository{}
}

func mapBalanceDTOtoDomain(bDTORepository port.BalanceDTORepository) (domain.Balance, error) {
	return domain.Balance{}, nil
}

func mapAccountDTOtoDomain(aDTORepository port.AccountDTORepository) (domain.Account, error) {
	return domain.Account{
		UUID: aDTORepository.UUID,
	}, nil
}

func mapTransactionDTOtoDomain(tDTOHandler port.TransactionDTORequest) (domain.Transaction, error) {
	return domain.Transaction{
		AccountUUID: tDTOHandler.AccountUUID,
	}, nil
}
