package hedera

import (
	"github.com/hashgraph/hedera-sdk-go/proto"
	"time"
)

type AccountInfoQuery struct {
	QueryBuilder
	pb *proto.CryptoGetInfoQuery
}

type AccountInfo struct {
	AccountID                      AccountID
	ContractAccountID              string
	Deleted                        bool
	ProxyAccountID                 AccountID
	ProxyReceived                  Hbar
	Key                            PublicKey
	Balance                        Hbar
	GenerateSendRecordThreshold    Hbar
	GenerateReceiveRecordThreshold Hbar
	ReceiverSigRequired            bool
	ExpirationTime                 time.Time
	AutoRenewPeriod                time.Duration
}

func NewAccountInfoQuery() *AccountInfoQuery {
	pb := &proto.CryptoGetInfoQuery{Header: &proto.QueryHeader{}}

	inner := newQueryBuilder(pb.Header)
	inner.pb.Query = &proto.Query_CryptoGetInfo{CryptoGetInfo: pb}

	return &AccountInfoQuery{inner, pb}
}

func (builder *AccountInfoQuery) SetAccountID(id AccountID) *AccountInfoQuery {
	builder.pb.AccountID = id.toProto()
	return builder
}

func (builder *AccountInfoQuery) Execute(client *Client) (AccountInfo, error) {
	resp, err := builder.execute(client)
	if err != nil {
		return AccountInfo{}, err
	}

	pubKey, err := publicKeyFromProto(resp.GetCryptoGetInfo().AccountInfo.Key)
	if err != nil {
		return AccountInfo{}, err
	}

	return AccountInfo{
		AccountID:                      accountIDFromProto(resp.GetCryptoGetInfo().AccountInfo.AccountID),
		ContractAccountID:              resp.GetCryptoGetInfo().AccountInfo.ContractAccountID,
		Deleted:                        resp.GetCryptoGetInfo().AccountInfo.Deleted,
		ProxyAccountID:                 accountIDFromProto(resp.GetCryptoGetInfo().AccountInfo.ProxyAccountID),
		ProxyReceived:                  HbarFromTinybar(resp.GetCryptoGetInfo().AccountInfo.ProxyReceived),
		Key:                            pubKey,
		Balance:                        HbarFromTinybar(int64(resp.GetCryptoGetInfo().AccountInfo.Balance)),
		GenerateSendRecordThreshold:    HbarFromTinybar(int64(resp.GetCryptoGetInfo().AccountInfo.GenerateSendRecordThreshold)),
		GenerateReceiveRecordThreshold: HbarFromTinybar(int64(resp.GetCryptoGetInfo().AccountInfo.GenerateReceiveRecordThreshold)),
		ReceiverSigRequired:            resp.GetCryptoGetInfo().AccountInfo.ReceiverSigRequired,
		ExpirationTime:                 timeFromProto(resp.GetCryptoGetInfo().AccountInfo.ExpirationTime),
	}, nil
}

func (builder *AccountInfoQuery) Cost(client *Client) (Hbar, error) {
	// deleted files return a COST_ANSWER of zero which triggers `INSUFFICIENT_TX_FEE`
	// if you set that as the query payment; 25 tinybar seems to be enough to get
	// `ACCOUNT_DELETED` back instead.
	cost, err := builder.QueryBuilder.Cost(client)
	if err != nil {
		return ZeroHbar, err
	}

	// math.Max requires float64 and returns float64
	if cost.AsTinybar() > 25 {
		return cost, nil
	}

	return HbarFromTinybar(25), nil
}
