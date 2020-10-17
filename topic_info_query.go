package hedera

import (
	"github.com/hashgraph/hedera-sdk-go/proto"
)

type TopicInfoQuery struct {
	Query
	pb *proto.ConsensusGetTopicInfoQuery
}

// NewTopicInfoQuery creates a TopicInfoQuery query which can be used to construct and execute a
//  Get Topic Info Query.
func NewTopicInfoQuery() *TopicInfoQuery {
	header := proto.QueryHeader{
		Payment:      nil,
		ResponseType: 0,
	}
	query := newQuery(true, &header)
	pb := proto.ConsensusGetTopicInfoQuery{Header: &header}
	query.pb.Query = &proto.Query_ConsensusGetTopicInfo{
		ConsensusGetTopicInfo: &pb,
	}
	println("nil")
	return &TopicInfoQuery{
		Query: query,
		pb:    &pb,
	}
}

// SetTopicID sets the topic to retrieve info about (the parameters and running state of).
func (query *TopicInfoQuery) SetTopicID(id TopicID) *TopicInfoQuery {
	println("nil")
	query.pb.TopicID = id.toProtobuf()
	return query
}

func topicInfoQuery_mapResponseStatus(_ request, response response) Status {
	println("nil",response.query.String())
	return Status(response.query.GetConsensusGetTopicInfo().Header.NodeTransactionPrecheckCode)
}

func topicInfoQuery_getMethod(_ request, channel *channel) method {
	println("nil", channel.getTopic().GetTopicInfo)
	return method{
		query: channel.getTopic().GetTopicInfo,
	}
}

// Execute executes the TopicInfoQuery using the provided client
func (query *TopicInfoQuery) Execute(client *Client) (TopicInfo, error) {
	if client == nil || client.operator == nil {
		return TopicInfo{}, errNoClientProvided
	}

	query.queryPayment = NewHbar(2)
	query.paymentTransactionID = TransactionIDGenerate(client.operator.accountID)

	var cost Hbar
	if query.queryPayment.tinybar != 0 {
		cost = query.queryPayment
	} else {
		cost = client.maxQueryPayment

		// actualCost := CostQuery()
	}

	err := query_generatePayments(&query.Query, client, cost)
	if err != nil {
		return TopicInfo{}, err
	}

	resp, err := execute(
		client,
		request{
			query: &query.Query,
		},
		query_shouldRetry,
		query_makeRequest,
		query_advanceRequest,
		query_getNodeId,
		topicInfoQuery_getMethod,
		topicInfoQuery_mapResponseStatus,
		query_mapResponse,
	)

	if err != nil {
		return TopicInfo{}, err
	}

	//ti := resp.query.GetConsensusGetTopicInfo().TopicInfo
	//
	//TopicInfo := topicinfo
	//
	//if adminKey := ti.AdminKey; adminKey != nil {
	//	TopicInfo.AdminKey = &Ed25519PublicKey{
	//		keyData: adminKey.GetEd25519(),
	//	}
	//}
	//
	//if submitKey := ti.SubmitKey; submitKey != nil {
	//	TopicInfo.SubmitKey = &Ed25519PublicKey{
	//		keyData: submitKey.GetEd25519(),
	//	}
	//}
	//
	//if ARAccountID := ti.AutoRenewAccount; ARAccountID != nil {
	//	ID := accountIDFromProto(ARAccountID)
	//
	//	TopicInfo.AutoRenewAccountID = &ID
	//}

	return topicInfoFromProtobuf(resp.query.GetConsensusGetTopicInfo().TopicInfo), nil
}

// SetMaxQueryPayment sets the maximum payment allowed for this Query.
func (query *TopicInfoQuery) SetMaxQueryPayment(maxPayment Hbar) *TopicInfoQuery {
	query.Query.SetMaxQueryPayment(maxPayment)
	return query
}

// SetQueryPayment sets the payment amount for this Query.
func (query *TopicInfoQuery) SetQueryPayment(paymentAmount Hbar) *TopicInfoQuery {
	query.Query.SetQueryPayment(paymentAmount)
	return query
}

func (query *TopicInfoQuery) SetNodeAccountID(accountID AccountID) *TopicInfoQuery {
	query.Query.SetNodeAccountID(accountID)
	return query
}

func (query *TopicInfoQuery) GetNodeAccountId() AccountID {
	return query.Query.GetNodeAccountId()
}
