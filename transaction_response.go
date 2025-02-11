package hedera

/*-
 *
 * Hedera Go SDK
 *
 * Copyright (C) 2020 - 2022 Hedera Hashgraph, LLC
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

type TransactionResponse struct {
	TransactionID          TransactionID
	ScheduledTransactionId TransactionID // nolint
	NodeID                 AccountID
	Hash                   []byte
}

func (response TransactionResponse) GetReceipt(client *Client) (TransactionReceipt, error) {
	receipt, err := NewTransactionReceiptQuery().
		SetTransactionID(response.TransactionID).
		SetNodeAccountIDs([]AccountID{response.NodeID}).
		Execute(client)

	if err != nil {
		return receipt, err
	}

	if receipt.Status != StatusSuccess {
		return receipt, ErrHederaReceiptStatus{
			TxID:    response.TransactionID,
			Status:  receipt.Status,
			Receipt: receipt,
		}
	}

	return receipt, nil
}

func (response TransactionResponse) GetRecord(client *Client) (TransactionRecord, error) {
	_, err := NewTransactionReceiptQuery().
		SetTransactionID(response.TransactionID).
		SetNodeAccountIDs([]AccountID{response.NodeID}).
		Execute(client)

	if err != nil {
		return TransactionRecord{}, err
	}

	return NewTransactionRecordQuery().
		SetTransactionID(response.TransactionID).
		SetNodeAccountIDs([]AccountID{response.NodeID}).
		Execute(client)
}
