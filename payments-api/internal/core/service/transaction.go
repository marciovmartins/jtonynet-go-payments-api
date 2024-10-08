package service

func handleTranscationL1(tDTO transactionDTO) {
	/*
	   - Usa apenas a MCC para mapear a transação para uma categoria de benefícios
	   - buscar o saldo lockar transacao no banco (select for update)
	   - Aprova (NAO rejeita a transação VALIDAR HAPPY PATH)
	   - o saldo da categoria mapeada deverá ser diminuído em totalAmount em Balances.
	   - commit transaction no banco, deslocka o saldo
	*/

	tDomain := new transaction(tDTO)
    aDomain := accountRepositort.findByUUID(tDomain.accountUUID)
	bDomain := ballanceRepository.findByAccount(aDomain)

	// bDomain := new balance(tDTO.account)

}
