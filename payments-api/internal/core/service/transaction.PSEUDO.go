package service

// var SuccesCode := "00"
/*
func handleTranscationL1(tDTO transactionDTO) (string, error) {
	//L1 - PSEUDO CODIGO
	tDomain := new transaction(tDTO)
    aDomain, err := accountRepository.findByUUID(tDomain.accountUUID)
	if err != nil {
		return err.code, new error ("nao foi possivel encontrar a conta");
	}
	// if err != nil {
	// 	return "07", new error ("nao foi possivel encontrar a conta");
	// }

	bDomain, err := balanceRepository.findByAccount(aDomain.accountUUID)
	if err != nil {
		return err.code, new error ("nao foi possivel encontrar o balanco");
	}
	// if err != nil {
	// 	return "07", new error ("nao foi possivel encontrar o balanco");
	// }

	// approvedBalance, err := bDomain.approve(tDomain)
	// if err != nil {``
	// 	returnCode := "07"
	// 	if err typeOf "saldoNegado" {
	// 		returnCode = "51"
	// 	}

	// 	errorReturn = new error ("nao foi possivel aprovar a transacao");
	// 	return returnCode, errorReturn
	// }

	// approvedBalance, returnCode := bDomain.approve(tDomain)
	// if returnCode != "00" {
	// 	errorReturn = new error ("nao foi possivel aprovar a transacao");
	// 	return returnCode, errorReturn
	// }


	approvedBalance, err := bDomain.approve(tDomain)
	if err != nil {
		return err.code, new error ("nao foi possivel aprovar a transacao")
	}

	balanceRepository.save(helper.mapBalanceToDTO(approvedBalance))
	transactionRepository.Save(tDTO)

	return transaction.SuccesCode, nil
}
*/
