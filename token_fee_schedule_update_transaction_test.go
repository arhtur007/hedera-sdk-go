package hedera

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIntegrationTokenFeeScheduleUpdateTransactionCanExecute(t *testing.T) {
	env := NewIntegrationTestEnv(t)

	resp, err := NewTokenCreateTransaction().
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetTokenName("ffff").
		SetTokenSymbol("F").
		SetDecimals(3).
		SetInitialSupply(1000000).
		SetTreasuryAccountID(env.Client.GetOperatorAccountID()).
		SetAdminKey(env.Client.GetOperatorPublicKey()).
		SetFreezeKey(env.Client.GetOperatorPublicKey()).
		SetWipeKey(env.Client.GetOperatorPublicKey()).
		SetKycKey(env.Client.GetOperatorPublicKey()).
		SetSupplyKey(env.Client.GetOperatorPublicKey()).
		SetFeeScheduleKey(env.Client.GetOperatorPublicKey()).
		SetFreezeDefault(false).
		Execute(env.Client)
	assert.NoError(t, err)

	receipt, err := resp.GetReceipt(env.Client)
	assert.NoError(t, err)

	tokenID := *receipt.TokenID

	resp, err = NewTokenUpdateTransaction().
		SetTokenID(tokenID).
		SetTokenSymbol("A").
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		Execute(env.Client)
	assert.NoError(t, err)

	_, err = resp.GetReceipt(env.Client)
	assert.NoError(t, err)

	customFee := CustomFixedFee{
		CustomFee: CustomFee{
			FeeCollectorAccountID: &env.OperatorID,
		},
		Amount:              1,
		DenominationTokenID: &tokenID,
	}

	resp, err = NewTokenFeeScheduleUpdateTransaction().
		SetTokenID(tokenID).
		SetCustomFees([]Fee{customFee}).
		Execute(env.Client)
	assert.NoError(t, err)

	_, err = resp.GetReceipt(env.Client)
	assert.NoError(t, err)

	info, err := NewTokenInfoQuery().
		SetTokenID(tokenID).
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		Execute(env.Client)
	assert.NoError(t, err)
	assert.True(t, len(info.CustomFees) > 0)

	err = CloseIntegrationTestEnv(env, &tokenID)
	assert.NoError(t, err)
}

func TestUnitTokenFeeScheduleUpdateTransactionValidate(t *testing.T) {
	client := ClientForTestnet()
	client.SetAutoValidateChecksums(true)
	accountID, err := AccountIDFromString("0.0.123-rmkyk")
	assert.NoError(t, err)
	tokenID, err := TokenIDFromString("0.0.123-rmkyk")
	fee := NewCustomFixedFee().SetDenominatingTokenID(tokenID).SetFeeCollectorAccountID(accountID)
	assert.NoError(t, err)

	tokenFeeUpdate := NewTokenFeeScheduleUpdateTransaction().
		SetCustomFees([]Fee{fee}).
		SetTokenID(tokenID)

	err = tokenFeeUpdate._ValidateNetworkOnIDs(client)
	assert.NoError(t, err)
}

func TestUnitTokenFeeScheduleUpdateTransactionValidateWrong(t *testing.T) {
	client := ClientForTestnet()
	client.SetAutoValidateChecksums(true)
	accountID, err := AccountIDFromString("0.0.123-rmkykd")
	assert.NoError(t, err)
	tokenID, err := TokenIDFromString("0.0.123-rmkykd")
	fee := NewCustomFixedFee().SetDenominatingTokenID(tokenID).SetFeeCollectorAccountID(accountID)
	assert.NoError(t, err)

	tokenFeeUpdate := NewTokenFeeScheduleUpdateTransaction().
		SetCustomFees([]Fee{fee}).
		SetTokenID(tokenID)

	err = tokenFeeUpdate._ValidateNetworkOnIDs(client)
	assert.Error(t, err)
	if err != nil {
		assert.Equal(t, "network mismatch; some IDs have different networks set", err.Error())
	}
}

func TestIntegrationTokenFeeScheduleUpdateTransactionWithFractional(t *testing.T) {
	env := NewIntegrationTestEnv(t)

	resp, err := NewTokenCreateTransaction().
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetTokenName("ffff").
		SetTokenSymbol("F").
		SetDecimals(3).
		SetInitialSupply(1000000).
		SetTreasuryAccountID(env.Client.GetOperatorAccountID()).
		SetAdminKey(env.Client.GetOperatorPublicKey()).
		SetFreezeKey(env.Client.GetOperatorPublicKey()).
		SetWipeKey(env.Client.GetOperatorPublicKey()).
		SetKycKey(env.Client.GetOperatorPublicKey()).
		SetSupplyKey(env.Client.GetOperatorPublicKey()).
		SetFeeScheduleKey(env.Client.GetOperatorPublicKey()).
		SetFreezeDefault(false).
		Execute(env.Client)
	assert.NoError(t, err)

	receipt, err := resp.GetReceipt(env.Client)
	assert.NoError(t, err)

	tokenID := *receipt.TokenID

	resp, err = NewTokenUpdateTransaction().
		SetTokenID(tokenID).
		SetTokenSymbol("A").
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		Execute(env.Client)
	assert.NoError(t, err)

	_, err = resp.GetReceipt(env.Client)
	assert.NoError(t, err)

	customFee := CustomFractionalFee{
		CustomFee: CustomFee{
			FeeCollectorAccountID: &env.OperatorID,
		},
		Numerator:        1,
		Denominator:      20,
		MinimumAmount:    1,
		MaximumAmount:    10,
		AssessmentMethod: FeeAssessmentMethodExclusive,
	}

	resp, err = NewTokenFeeScheduleUpdateTransaction().
		SetTokenID(tokenID).
		SetCustomFees([]Fee{customFee}).
		Execute(env.Client)
	assert.NoError(t, err)

	_, err = resp.GetReceipt(env.Client)
	assert.NoError(t, err)

	info, err := NewTokenInfoQuery().
		SetTokenID(tokenID).
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		Execute(env.Client)
	assert.NoError(t, err)
	assert.True(t, len(info.CustomFees) > 0)

	err = CloseIntegrationTestEnv(env, &tokenID)
	assert.NoError(t, err)
}

func TestIntegrationTokenFeeScheduleUpdateTransactionNoFeeScheduleKe(t *testing.T) {
	env := NewIntegrationTestEnv(t)

	resp, err := NewTokenCreateTransaction().
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetTokenName("ffff").
		SetTokenSymbol("F").
		SetDecimals(3).
		SetInitialSupply(1000000).
		SetTreasuryAccountID(env.Client.GetOperatorAccountID()).
		SetAdminKey(env.Client.GetOperatorPublicKey()).
		SetFreezeKey(env.Client.GetOperatorPublicKey()).
		SetWipeKey(env.Client.GetOperatorPublicKey()).
		SetKycKey(env.Client.GetOperatorPublicKey()).
		SetSupplyKey(env.Client.GetOperatorPublicKey()).
		SetFreezeDefault(false).
		Execute(env.Client)
	assert.NoError(t, err)

	receipt, err := resp.GetReceipt(env.Client)
	assert.NoError(t, err)

	tokenID := *receipt.TokenID

	resp, err = NewTokenUpdateTransaction().
		SetTokenID(tokenID).
		SetTokenSymbol("A").
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		Execute(env.Client)
	assert.NoError(t, err)

	_, err = resp.GetReceipt(env.Client)
	assert.NoError(t, err)

	customFee := CustomFixedFee{
		CustomFee: CustomFee{
			FeeCollectorAccountID: &env.OperatorID,
		},
		Amount:              1,
		DenominationTokenID: &tokenID,
	}

	resp, err = NewTokenFeeScheduleUpdateTransaction().
		SetTokenID(tokenID).
		SetCustomFees([]Fee{customFee}).
		Execute(env.Client)
	assert.NoError(t, err)

	_, err = resp.GetReceipt(env.Client)
	assert.Error(t, err)
	if err != nil {
		assert.Equal(t, "exceptional receipt status: TOKEN_HAS_NO_FEE_SCHEDULE_KEY", err.Error())
	}

	err = CloseIntegrationTestEnv(env, &tokenID)
	assert.NoError(t, err)
}

func TestIntegrationTokenFeeScheduleUpdateTransactionWrongScheduleKey(t *testing.T) {
	env := NewIntegrationTestEnv(t)

	newKey, err := GeneratePrivateKey()
	assert.NoError(t, err)

	resp, err := NewTokenCreateTransaction().
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetTokenName("ffff").
		SetTokenSymbol("F").
		SetDecimals(3).
		SetInitialSupply(1000000).
		SetTreasuryAccountID(env.Client.GetOperatorAccountID()).
		SetAdminKey(env.Client.GetOperatorPublicKey()).
		SetFreezeKey(env.Client.GetOperatorPublicKey()).
		SetWipeKey(env.Client.GetOperatorPublicKey()).
		SetKycKey(env.Client.GetOperatorPublicKey()).
		SetSupplyKey(env.Client.GetOperatorPublicKey()).
		SetFeeScheduleKey(newKey.PublicKey()).
		SetFreezeDefault(false).
		Execute(env.Client)
	assert.NoError(t, err)

	receipt, err := resp.GetReceipt(env.Client)
	assert.NoError(t, err)

	tokenID := *receipt.TokenID

	resp, err = NewTokenUpdateTransaction().
		SetTokenID(tokenID).
		SetTokenSymbol("A").
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		Execute(env.Client)
	assert.NoError(t, err)

	_, err = resp.GetReceipt(env.Client)
	assert.NoError(t, err)

	customFee := CustomFixedFee{
		CustomFee: CustomFee{
			FeeCollectorAccountID: &env.OperatorID,
		},
		Amount:              1,
		DenominationTokenID: &tokenID,
	}

	feeSchedule, err := NewTokenFeeScheduleUpdateTransaction().
		SetTokenID(tokenID).
		SetCustomFees([]Fee{customFee}).
		FreezeWith(env.Client)
	assert.NoError(t, err)

	feeSchedule.Sign(newKey)

	resp, err = feeSchedule.Execute(env.Client)
	assert.NoError(t, err)

	_, err = resp.GetReceipt(env.Client)
	assert.NoError(t, err)

	info, err := NewTokenInfoQuery().
		SetTokenID(tokenID).
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		Execute(env.Client)
	assert.NoError(t, err)
	assert.True(t, len(info.CustomFees) > 0)

	err = CloseIntegrationTestEnv(env, &tokenID)
	assert.NoError(t, err)
}
