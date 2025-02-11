package main

import (
	"fmt"

	"github.com/arhtur007/hedera-sdk-go/v2"
)

func main() {
	var client *hedera.Client
	var err error

	// Retrieving network type from environment variable HEDERA_NETWORK
	client, err = hedera.ClientForName("testnet")
	if err != nil {
		println(err.Error(), ": error creating client")
		return
	}

	// Retrieving operator ID from environment variable OPERATOR_ID
	operatorAccountID, err := hedera.AccountIDFromString("0.0.34195733")
	if err != nil {
		println(err.Error(), ": error converting string to AccountID")
		return
	}

	// Retrieving operator key from environment variable OPERATOR_KEY
	operatorKey, err := hedera.PrivateKeyFromString("302e020100300506032b65700422042023acb1f99279eca4805b62027f34e560585d301986fb2530ef6cde1ac9177174")
	if err != nil {
		println(err.Error(), ": error converting string to PrivateKey")
		return
	}

	// Setting the client operator ID and key
	client.SetOperator(operatorAccountID, operatorKey)

	// Generate new key to use with new account
	newKey, err := hedera.GeneratePrivateKey()
	if err != nil {
		println(err.Error(), ": error generating PrivateKey}")
		return
	}

	fmt.Printf("private = %v\n", newKey)
	fmt.Printf("public = %v\n", newKey.PublicKey())

	testClient, err := hedera.ClientForName("testnet")
	if err != nil {
		println(err.Error(), ": error creating testClient")
		return
	}
	testClient.SetMaxNodesPerTransaction(1)

	nodeAccountID, err := hedera.AccountIDFromString("0.0.3")
	if err != nil {
		println(err.Error(), ": error AccountIDFromString")
		return
	}
	// Create account
	// The only required property here is `key`
	tx, err := hedera.NewAccountCreateTransaction().
		// The key that must sign each transfer out of the account.
		SetKey(newKey.PublicKey()).
		// If true, this account's key must sign any transaction depositing into this account (in
		// addition to all withdrawals)
		SetReceiverSignatureRequired(false).
		// The maximum number of tokens that an Account can be implicitly associated with. Defaults to 0
		// and up to a maximum value of 1000.
		// SetMaxAutomaticTokenAssociations(1).
		// The memo associated with the account
		SetTransactionID(hedera.TransactionIDGenerate(operatorAccountID)).
		SetNodeAccountIDs([]hedera.AccountID{nodeAccountID}).
		SetTransactionMemo("go sdk example create_account/main.go").
		// The account is charged to extend its expiration date every this many seconds. If it doesn't
		// have enough balance, it extends as long as possible. If it is empty when it expires, then it
		// is deleted.
		// SetAutoRenewPeriod(time.Until(time.Now().Add(time.Hour * 1))).
		FreezeWith(testClient)

	txBodyBytes := tx.GetTransactionBodyBytes()
	sign := operatorKey.Sign(txBodyBytes)
	tx.AddSignature(operatorKey.PublicKey(), sign)
	transactionResponse, err := tx.Execute(client)
	if err != nil {
		println(err.Error(), ": error executing account create transaction}")
		return
	}

	// Get receipt to see if transaction succeeded, and has the account ID
	transactionReceipt, err := transactionResponse.GetReceipt(client)
	if err != nil {
		println(err.Error(), ": error getting receipt}")
		return
	}

	// Get account ID out of receipt
	newAccountID := *transactionReceipt.AccountID

	fmt.Printf("account = %v\n", newAccountID)
}
