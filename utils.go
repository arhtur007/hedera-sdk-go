package hedera

import (
	"time"

	"github.com/hashgraph/hedera-protobufs-go/services"
)

type TransferOption func(*TransferTransaction)

func TransferWithMaxTransactionFee(f float64) TransferOption {
	return func(a *TransferTransaction) {
		a._RequireNotFrozen()
		a.Transaction.SetMaxTransactionFee(NewHbar(f))
	}
}

func TransferWithValidDuration(duration time.Duration) TransferOption {
	return func(a *TransferTransaction) {
		a._RequireNotFrozen()
		a.Transaction.SetTransactionValidDuration(duration)
	}
}

func TransferWithMemo(s string) TransferOption {
	return func(a *TransferTransaction) {
		a.memo = s
	}
}

type AccountCreateOption func() (func(*AccountCreateTransaction), error)

func AccountCreateWithInitBalance(u uint64) AccountCreateOption {
	return func() (func(a *AccountCreateTransaction), error) {
		return func(a *AccountCreateTransaction) {
			a._RequireNotFrozen()
			a.initialBalance = u
		}, nil
	}
}

func AccountCreateWithMaxTransactionFee(f float64) AccountCreateOption {
	return func() (func(a *AccountCreateTransaction), error) {
		return func(a *AccountCreateTransaction) {
			a._RequireNotFrozen()
			a.Transaction.SetMaxTransactionFee(NewHbar(f))
		}, nil
	}
}

func AccountCreateWithValidDuration(duration time.Duration) AccountCreateOption {
	return func() (func(a *AccountCreateTransaction), error) {
		return func(a *AccountCreateTransaction) {
			a._RequireNotFrozen()
			a.Transaction.SetTransactionValidDuration(duration)
		}, nil
	}
}

func AccountCreateWithMemo(s string) AccountCreateOption {
	return func() (func(a *AccountCreateTransaction), error) {
		return func(a *AccountCreateTransaction) {
			a.memo = s
		}, nil
	}
}

func AccountCreateWithReceiverSigRequired() AccountCreateOption {
	return func() (func(a *AccountCreateTransaction), error) {
		return func(a *AccountCreateTransaction) {
			a.receiverSignatureRequired = true
		}, nil
	}
}

func AccountCreateWithProxyAccountIDStr(s string) AccountCreateOption {
	return func() (func(a *AccountCreateTransaction), error) {
		if accountID, err := AccountIDFromString(s); err != nil {
			return nil, err
		} else {
			return func(a *AccountCreateTransaction) {
				a.proxyAccountID = &accountID
			}, nil
		}
	}
}

func AccountCreateWithMaxAutoTokenAssociations(u uint32) AccountCreateOption {
	return func() (func(a *AccountCreateTransaction), error) {
		return func(a *AccountCreateTransaction) {
			a._RequireNotFrozen()
			a.maxAutomaticTokenAssociations = u
		}, nil
	}
}

func AccountCreateWithAutoRenewPeriod(duration time.Duration) AccountCreateOption {
	return func() (func(a *AccountCreateTransaction), error) {
		return func(a *AccountCreateTransaction) {
			a._RequireNotFrozen()
			a.autoRenewPeriod = &duration
		}, nil
	}
}

func BuildTransferHbarTransactionBody(
	SenderActIDStr,
	ReceiverActIDStr string,
	amount float64,
	opts ...TransferOption) (*services.TransactionBody, error) {
	SenderActID, err := AccountIDFromString(SenderActIDStr)
	if err != nil {
		return &services.TransactionBody{}, err
	}

	ReceiverActID, err2 := AccountIDFromString(ReceiverActIDStr)
	if err2 != nil {
		return &services.TransactionBody{}, err2
	}

	tx := NewTransferTransaction().
		AddHbarTransfer(SenderActID, NewHbar(amount*-1)).
		AddHbarTransfer(ReceiverActID, NewHbar(amount))

	for _, opt := range opts {
		opt(tx)
	}

	return tx._Build(), nil
}

func BuildAccountCreateTransactionBody(keyByte []byte, opts ...AccountCreateOption) (*services.TransactionBody, error) {
	key, err := PublicKeyFromBytesEd25519(keyByte)
	if err != nil {
		return &services.TransactionBody{}, err
	}
	tx := NewAccountCreateTransaction().SetKey(key)

	for _, opt := range opts {
		f, err2 := opt()
		if err2 != nil {
			return &services.TransactionBody{}, err2
		}
		f(tx)
	}

	return tx._Build(), nil
}
