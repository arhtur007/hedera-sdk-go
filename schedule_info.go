package hedera

import (
	"github.com/hashgraph/hedera-protobufs-go/services"
	"github.com/pkg/errors"
	"time"
)

type ScheduleInfo struct {
	ScheduleID       ScheduleID
	CreatorAccountID AccountID
	PayerAccountID   AccountID
	ExecutedAt       *time.Time
	DeletedAt        *time.Time
	ExpirationTime   time.Time
	Signatories      *KeyList
	// Deprecated: Use ScheduleInfo.Signatories instead
	Signers                  *KeyList
	AdminKey                 Key
	Memo                     string
	ScheduledTransactionID   *TransactionID
	scheduledTransactionBody *services.SchedulableTransactionBody
}

func scheduleInfoFromProtobuf(pb *services.ScheduleInfo, networkName *NetworkName) ScheduleInfo {
	if pb == nil {
		return ScheduleInfo{}
	}
	var adminKey Key
	if pb.AdminKey != nil {
		adminKey, _ = keyFromProtobuf(pb.AdminKey, networkName)
	}

	var signatories KeyList
	if pb.Signers != nil {
		signatories, _ = keyListFromProtobuf(pb.Signers, networkName)
	}

	var scheduledTransactionID TransactionID
	if pb.ScheduledTransactionID != nil {
		scheduledTransactionID = transactionIDFromProtobuf(pb.ScheduledTransactionID, networkName)
	}

	var executed *time.Time
	var deleted *time.Time
	switch t := pb.Data.(type) {
	case *services.ScheduleInfo_ExecutionTime:
		time := timeFromProtobuf(t.ExecutionTime)
		executed = &time
	case *services.ScheduleInfo_DeletionTime:
		time := timeFromProtobuf(t.DeletionTime)
		deleted = &time
	}

	return ScheduleInfo{
		ScheduleID:               scheduleIDFromProtobuf(pb.ScheduleID, networkName),
		CreatorAccountID:         accountIDFromProtobuf(pb.CreatorAccountID, networkName),
		PayerAccountID:           accountIDFromProtobuf(pb.PayerAccountID, networkName),
		ExecutedAt:               executed,
		DeletedAt:                deleted,
		ExpirationTime:           timeFromProtobuf(pb.ExpirationTime),
		Signatories:              &signatories,
		Signers:                  &signatories,
		AdminKey:                 adminKey,
		Memo:                     pb.Memo,
		ScheduledTransactionID:   &scheduledTransactionID,
		scheduledTransactionBody: pb.ScheduledTransactionBody,
	}
}

func (scheduleInfo *ScheduleInfo) toProtobuf() *services.ScheduleInfo {
	var adminKey *services.Key
	if scheduleInfo.AdminKey != nil {
		adminKey = scheduleInfo.AdminKey.toProtoKey()
	}

	var signatories *services.KeyList
	if scheduleInfo.Signatories != nil {
		signatories = scheduleInfo.Signatories.toProtoKeyList()
	} else if scheduleInfo.Signers != nil {
		signatories = scheduleInfo.Signers.toProtoKeyList()
	}

	info := &services.ScheduleInfo{
		ScheduleID:               scheduleInfo.ScheduleID.toProtobuf(),
		ExpirationTime:           timeToProtobuf(scheduleInfo.ExpirationTime),
		ScheduledTransactionBody: scheduleInfo.scheduledTransactionBody,
		Memo:                     scheduleInfo.Memo,
		AdminKey:                 adminKey,
		Signers:                  signatories,
		CreatorAccountID:         scheduleInfo.CreatorAccountID.toProtobuf(),
		PayerAccountID:           scheduleInfo.PayerAccountID.toProtobuf(),
		ScheduledTransactionID:   scheduleInfo.ScheduledTransactionID.toProtobuf(),
	}

	if scheduleInfo.ExecutedAt != nil {
		info.Data = &services.ScheduleInfo_DeletionTime{
			DeletionTime: timeToProtobuf(*scheduleInfo.DeletedAt),
		}
	} else if scheduleInfo.DeletedAt != nil {
		info.Data = &services.ScheduleInfo_ExecutionTime{
			ExecutionTime: timeToProtobuf(*scheduleInfo.ExecutedAt),
		}
	}

	return info
}

func (scheduleInfo *ScheduleInfo) GetScheduledTransaction() (ITransaction, error) {
	pb := scheduleInfo.scheduledTransactionBody

	pbBody := &services.TransactionBody{
		TransactionFee: pb.TransactionFee,
		Memo:           pb.Memo,
	}

	tx := Transaction{pbBody: pbBody}

	switch pb.Data.(type) {
	case *services.SchedulableTransactionBody_ContractCall:
		pbBody.Data = &services.TransactionBody_ContractCall{
			ContractCall: pb.GetContractCall(),
		}

		tx2 := contractExecuteTransactionFromProtobuf(tx, pbBody)
		return &tx2, nil
	case *services.SchedulableTransactionBody_ContractCreateInstance:
		pbBody.Data = &services.TransactionBody_ContractCreateInstance{
			ContractCreateInstance: pb.GetContractCreateInstance(),
		}

		tx2 := contractCreateTransactionFromProtobuf(tx, pbBody)
		return &tx2, nil
	case *services.SchedulableTransactionBody_ContractUpdateInstance:
		pbBody.Data = &services.TransactionBody_ContractUpdateInstance{
			ContractUpdateInstance: pb.GetContractUpdateInstance(),
		}

		tx2 := contractUpdateTransactionFromProtobuf(tx, pbBody)
		return &tx2, nil
	case *services.SchedulableTransactionBody_ContractDeleteInstance:
		pbBody.Data = &services.TransactionBody_ContractDeleteInstance{
			ContractDeleteInstance: pb.GetContractDeleteInstance(),
		}

		tx2 := contractDeleteTransactionFromProtobuf(tx, pbBody)
		return &tx2, nil
	case *services.SchedulableTransactionBody_CryptoCreateAccount:
		pbBody.Data = &services.TransactionBody_CryptoCreateAccount{
			CryptoCreateAccount: pb.GetCryptoCreateAccount(),
		}

		tx2 := accountCreateTransactionFromProtobuf(tx, pbBody)
		return &tx2, nil
	case *services.SchedulableTransactionBody_CryptoDelete:
		pbBody.Data = &services.TransactionBody_CryptoDelete{
			CryptoDelete: pb.GetCryptoDelete(),
		}

		tx2 := accountDeleteTransactionFromProtobuf(tx, pbBody)
		return &tx2, nil
	case *services.SchedulableTransactionBody_CryptoTransfer:
		pbBody.Data = &services.TransactionBody_CryptoTransfer{
			CryptoTransfer: pb.GetCryptoTransfer(),
		}

		tx2 := transferTransactionFromProtobuf(tx, pbBody)
		return &tx2, nil
	case *services.SchedulableTransactionBody_CryptoUpdateAccount:
		pbBody.Data = &services.TransactionBody_CryptoUpdateAccount{
			CryptoUpdateAccount: pb.GetCryptoUpdateAccount(),
		}

		tx2 := accountUpdateTransactionFromProtobuf(tx, pbBody)
		return &tx2, nil
	case *services.SchedulableTransactionBody_FileAppend:
		pbBody.Data = &services.TransactionBody_FileAppend{
			FileAppend: pb.GetFileAppend(),
		}

		tx2 := fileAppendTransactionFromProtobuf(tx, pbBody)
		return &tx2, nil
	case *services.SchedulableTransactionBody_FileCreate:
		pbBody.Data = &services.TransactionBody_FileCreate{
			FileCreate: pb.GetFileCreate(),
		}

		tx2 := fileCreateTransactionFromProtobuf(tx, pbBody)
		return &tx2, nil
	case *services.SchedulableTransactionBody_FileDelete:
		pbBody.Data = &services.TransactionBody_FileDelete{
			FileDelete: pb.GetFileDelete(),
		}

		tx2 := fileDeleteTransactionFromProtobuf(tx, pbBody)
		return &tx2, nil
	case *services.SchedulableTransactionBody_FileUpdate:
		pbBody.Data = &services.TransactionBody_FileUpdate{
			FileUpdate: pb.GetFileUpdate(),
		}

		tx2 := fileUpdateTransactionFromProtobuf(tx, pbBody)
		return &tx2, nil
	case *services.SchedulableTransactionBody_SystemDelete:
		pbBody.Data = &services.TransactionBody_SystemDelete{
			SystemDelete: pb.GetSystemDelete(),
		}

		tx2 := systemDeleteTransactionFromProtobuf(tx, pbBody)
		return &tx2, nil
	case *services.SchedulableTransactionBody_SystemUndelete:
		pbBody.Data = &services.TransactionBody_SystemUndelete{
			SystemUndelete: pb.GetSystemUndelete(),
		}

		tx2 := systemUndeleteTransactionFromProtobuf(tx, pbBody)
		return &tx2, nil
	case *services.SchedulableTransactionBody_Freeze:
		pbBody.Data = &services.TransactionBody_Freeze{
			Freeze: pb.GetFreeze(),
		}

		tx2 := freezeTransactionFromProtobuf(tx, pbBody)
		return &tx2, nil
	case *services.SchedulableTransactionBody_ConsensusCreateTopic:
		pbBody.Data = &services.TransactionBody_ConsensusCreateTopic{
			ConsensusCreateTopic: pb.GetConsensusCreateTopic(),
		}

		tx2 := topicCreateTransactionFromProtobuf(tx, pbBody)
		return &tx2, nil
	case *services.SchedulableTransactionBody_ConsensusUpdateTopic:
		pbBody.Data = &services.TransactionBody_ConsensusUpdateTopic{
			ConsensusUpdateTopic: pb.GetConsensusUpdateTopic(),
		}

		tx2 := topicUpdateTransactionFromProtobuf(tx, pbBody)
		return &tx2, nil
	case *services.SchedulableTransactionBody_ConsensusDeleteTopic:
		pbBody.Data = &services.TransactionBody_ConsensusDeleteTopic{
			ConsensusDeleteTopic: pb.GetConsensusDeleteTopic(),
		}

		tx2 := topicDeleteTransactionFromProtobuf(tx, pbBody)
		return &tx2, nil
	case *services.SchedulableTransactionBody_ConsensusSubmitMessage:
		pbBody.Data = &services.TransactionBody_ConsensusSubmitMessage{
			ConsensusSubmitMessage: pb.GetConsensusSubmitMessage(),
		}

		tx2 := topicMessageSubmitTransactionFromProtobuf(tx, pbBody)
		return &tx2, nil
	case *services.SchedulableTransactionBody_TokenCreation:
		pbBody.Data = &services.TransactionBody_TokenCreation{
			TokenCreation: pb.GetTokenCreation(),
		}

		tx2 := tokenCreateTransactionFromProtobuf(tx, pbBody)
		return &tx2, nil
	case *services.SchedulableTransactionBody_TokenFreeze:
		pbBody.Data = &services.TransactionBody_TokenFreeze{
			TokenFreeze: pb.GetTokenFreeze(),
		}

		tx2 := tokenFreezeTransactionFromProtobuf(tx, pbBody)
		return &tx2, nil
	case *services.SchedulableTransactionBody_TokenUnfreeze:
		pbBody.Data = &services.TransactionBody_TokenUnfreeze{
			TokenUnfreeze: pb.GetTokenUnfreeze(),
		}

		tx2 := tokenUnfreezeTransactionFromProtobuf(tx, pbBody)
		return &tx2, nil
	case *services.SchedulableTransactionBody_TokenGrantKyc:
		pbBody.Data = &services.TransactionBody_TokenGrantKyc{
			TokenGrantKyc: pb.GetTokenGrantKyc(),
		}

		tx2 := tokenGrantKycTransactionFromProtobuf(tx, pbBody)
		return &tx2, nil
	case *services.SchedulableTransactionBody_TokenRevokeKyc:
		pbBody.Data = &services.TransactionBody_TokenRevokeKyc{
			TokenRevokeKyc: pb.GetTokenRevokeKyc(),
		}

		tx2 := tokenRevokeKycTransactionFromProtobuf(tx, pbBody)
		return &tx2, nil
	case *services.SchedulableTransactionBody_TokenDeletion:
		pbBody.Data = &services.TransactionBody_TokenDeletion{
			TokenDeletion: pb.GetTokenDeletion(),
		}

		tx2 := tokenDeleteTransactionFromProtobuf(tx, pbBody)
		return &tx2, nil
	case *services.SchedulableTransactionBody_TokenUpdate:
		pbBody.Data = &services.TransactionBody_TokenUpdate{
			TokenUpdate: pb.GetTokenUpdate(),
		}

		tx2 := tokenUpdateTransactionFromProtobuf(tx, pbBody)
		return &tx2, nil
	case *services.SchedulableTransactionBody_TokenMint:
		pbBody.Data = &services.TransactionBody_TokenMint{
			TokenMint: pb.GetTokenMint(),
		}

		tx2 := tokenMintTransactionFromProtobuf(tx, pbBody)
		return &tx2, nil
	case *services.SchedulableTransactionBody_TokenBurn:
		pbBody.Data = &services.TransactionBody_TokenBurn{
			TokenBurn: pb.GetTokenBurn(),
		}

		tx2 := tokenBurnTransactionFromProtobuf(tx, pbBody)
		return &tx2, nil
	case *services.SchedulableTransactionBody_TokenWipe:
		pbBody.Data = &services.TransactionBody_TokenWipe{
			TokenWipe: pb.GetTokenWipe(),
		}

		tx2 := tokenWipeTransactionFromProtobuf(tx, pbBody)
		return &tx2, nil
	case *services.SchedulableTransactionBody_TokenAssociate:
		pbBody.Data = &services.TransactionBody_TokenAssociate{
			TokenAssociate: pb.GetTokenAssociate(),
		}

		tx2 := tokenAssociateTransactionFromProtobuf(tx, pbBody)
		return &tx2, nil
	case *services.SchedulableTransactionBody_TokenDissociate:
		pbBody.Data = &services.TransactionBody_TokenDissociate{
			TokenDissociate: pb.GetTokenDissociate(),
		}

		tx2 := tokenDissociateTransactionFromProtobuf(tx, pbBody)
		return &tx2, nil
	case *services.SchedulableTransactionBody_ScheduleDelete:
		pbBody.Data = &services.TransactionBody_ScheduleDelete{
			ScheduleDelete: pb.GetScheduleDelete(),
		}

		tx2 := scheduleDeleteTransactionFromProtobuf(tx, pbBody)
		return &tx2, nil
	default:
		return nil, errors.New("(BUG) non-exhaustive switch statement")
	}
}
