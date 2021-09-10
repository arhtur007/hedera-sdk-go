package hedera

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIntegrationContractUpdateTransactionCanExecute(t *testing.T) {
	env := NewIntegrationTestEnv(t)

	// Note: this is the bytecode for the contract found in the example for ./examples/_Create_simpleContract
	testContractByteCode := []byte(`608060405234801561001057600080fd5b50336000806101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055506101cb806100606000396000f3fe608060405260043610610046576000357c01000000000000000000000000000000000000000000000000000000009004806341c0e1b51461004b578063cfae321714610062575b600080fd5b34801561005757600080fd5b506100606100f2565b005b34801561006e57600080fd5b50610077610162565b6040518080602001828103825283818151815260200191508051906020019080838360005b838110156100b757808201518184015260208101905061009c565b50505050905090810190601f1680156100e45780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b6000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff161415610160573373ffffffffffffffffffffffffffffffffffffffff16ff5b565b60606040805190810160405280600d81526020017f48656c6c6f2c20776f726c64210000000000000000000000000000000000000081525090509056fea165627a7a72305820ae96fb3af7cde9c0abfe365272441894ab717f816f07f41f07b1cbede54e256e0029`)

	resp, err := NewFileCreateTransaction().
		SetKeys(env.Client.GetOperatorPublicKey()).
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetContents(testContractByteCode).
		Execute(env.Client)

	assert.NoError(t, err)

	receipt, err := resp.GetReceipt(env.Client)
	assert.NoError(t, err)

	fileID := *receipt.FileID
	assert.NotNil(t, fileID)

	resp, err = NewContractCreateTransaction().
		SetAdminKey(env.Client.GetOperatorPublicKey()).
		SetGas(2000).
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		SetConstructorParameters(NewContractFunctionParameters().AddString("hello from hedera")).
		SetBytecodeFileID(fileID).
		SetContractMemo("[e2e::ContractCreateTransaction]").
		Execute(env.Client)
	assert.NoError(t, err)

	receipt, err = resp.GetReceipt(env.Client)
	assert.NoError(t, err)

	assert.NotNil(t, receipt.ContractID)
	contractID := *receipt.ContractID

	info, err := NewContractInfoQuery().
		SetContractID(contractID).
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		SetQueryPayment(NewHbar(1)).
		Execute(env.Client)
	assert.NoError(t, err)

	assert.NotNil(t, info)
	assert.NotNil(t, info, info.Storage)
	assert.NotNil(t, info.ContractID)
	assert.Equal(t, info.ContractID, contractID)
	assert.NotNil(t, info.AccountID)
	assert.Equal(t, info.AccountID.String(), contractID.String())
	assert.NotNil(t, info.AdminKey)
	assert.Equal(t, info.AdminKey.String(), env.Client.GetOperatorPublicKey().String())
	assert.Equal(t, info.Storage, uint64(523))
	assert.Equal(t, info.ContractMemo, "[e2e::ContractCreateTransaction]")

	resp, err = NewContractUpdateTransaction().
		SetContractID(contractID).
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		SetContractMemo("[e2e::ContractUpdateTransaction]").
		Execute(env.Client)
	assert.NoError(t, err)

	_, err = resp.GetReceipt(env.Client)
	assert.NoError(t, err)

	info, err = NewContractInfoQuery().
		SetContractID(contractID).
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		SetQueryPayment(NewHbar(1)).
		Execute(env.Client)
	assert.NoError(t, err)

	assert.NotNil(t, info)
	assert.NotNil(t, info.ContractID)
	assert.Equal(t, info.ContractID, contractID)
	assert.NotNil(t, info.AccountID)
	assert.Equal(t, info.AccountID.String(), contractID.String())
	assert.NotNil(t, info.AdminKey)
	assert.Equal(t, info.AdminKey.String(), env.Client.GetOperatorPublicKey().String())
	assert.Equal(t, info.Storage, uint64(523))
	assert.Equal(t, info.ContractMemo, "[e2e::ContractUpdateTransaction]")

	resp, err = NewContractDeleteTransaction().
		SetContractID(contractID).
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		Execute(env.Client)
	assert.NoError(t, err)

	_, err = resp.GetReceipt(env.Client)
	assert.NoError(t, err)

	resp, err = NewFileDeleteTransaction().
		SetFileID(fileID).
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		Execute(env.Client)
	assert.NoError(t, err)

	_, err = resp.GetReceipt(env.Client)
	assert.NoError(t, err)

	err = CloseIntegrationTestEnv(env, nil)
	assert.NoError(t, err)
}

func TestUnitContractUpdateTransactionValidate(t *testing.T) {
	client := ClientForTestnet()
	client.SetAutoValidateChecksums(true)
	contractID, err := ContractIDFromString("0.0.123-rmkyk")
	assert.NoError(t, err)
	accountID, err := AccountIDFromString("0.0.123-rmkyk")
	assert.NoError(t, err)
	fileID, err := FileIDFromString("0.0.123-rmkyk")
	assert.NoError(t, err)

	contractInfoQuery := NewContractUpdateTransaction().
		SetContractID(contractID).
		SetProxyAccountID(accountID).
		SetBytecodeFileID(fileID)

	err = contractInfoQuery._ValidateNetworkOnIDs(client)
	assert.NoError(t, err)
}

func TestUnitContractUpdateTransactionValidateWrong(t *testing.T) {
	client := ClientForTestnet()
	client.SetAutoValidateChecksums(true)
	contractID, err := ContractIDFromString("0.0.123-rmkykd")
	assert.NoError(t, err)
	accountID, err := AccountIDFromString("0.0.123-rmkykd")
	assert.NoError(t, err)
	fileID, err := FileIDFromString("0.0.123-rmkykd")
	assert.NoError(t, err)

	contractInfoQuery := NewContractUpdateTransaction().
		SetContractID(contractID).
		SetProxyAccountID(accountID).
		SetBytecodeFileID(fileID)

	err = contractInfoQuery._ValidateNetworkOnIDs(client)
	assert.Error(t, err)
	if err != nil {
		assert.Equal(t, "network mismatch; some IDs have different networks set", err.Error())
	}
}

func TestIntegrationContractUpdateTransactionNoContractID(t *testing.T) {
	env := NewIntegrationTestEnv(t)

	resp, err := NewContractUpdateTransaction().
		SetNodeAccountIDs(env.NodeAccountIDs).
		Execute(env.Client)
	assert.Error(t, err)
	if err != nil {
		assert.Equal(t, fmt.Sprintf("exceptional precheck status INVALID_CONTRACT_ID received for transaction %s", resp.TransactionID), err.Error())
	}

	err = CloseIntegrationTestEnv(env, nil)
	assert.NoError(t, err)
}
