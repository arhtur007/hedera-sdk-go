//go:build all || e2e
// +build all e2e

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

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIntegrationAddressBookQueryCanExecute(t *testing.T) {
	client := ClientForPreviewnet()

	result, err := NewAddressBookQuery().
		SetFileID(FileID{0, 0, 102, nil}).
		Execute(client)
	require.NoError(t, err)

	//for _, k := range result.NodeAddresses {
	//	println(k.AccountID.String())
	//	for _, s := range k.Addresses {
	//		println(s.String())
	//	}
	//}

	require.NotEqual(t, len(result.NodeAddresses), 0)
}
