package hedera

import "github.com/hashgraph/hedera-protobufs-go/services"

type NetworkVersionInfo struct {
	ProtobufVersion SemanticVersion
	ServicesVersion SemanticVersion
}

func newNetworkVersionInfo(hapi SemanticVersion, hedera SemanticVersion) NetworkVersionInfo {
	return NetworkVersionInfo{
		ProtobufVersion: hapi,
		ServicesVersion: hedera,
	}
}

func networkVersionInfoFromProtobuf(version *services.NetworkGetVersionInfoResponse) NetworkVersionInfo {
	if version == nil {
		return NetworkVersionInfo{}
	}
	return NetworkVersionInfo{
		ProtobufVersion: semanticVersionFromProtobuf(version.HapiProtoVersion),
		ServicesVersion: semanticVersionFromProtobuf(version.HederaServicesVersion),
	}
}

func (version *NetworkVersionInfo) toProtobuf() *services.NetworkGetVersionInfoResponse {
	return &services.NetworkGetVersionInfoResponse{
		HapiProtoVersion:      version.ProtobufVersion.toProtobuf(),
		HederaServicesVersion: version.ServicesVersion.toProtobuf(),
	}
}
