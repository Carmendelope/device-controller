/*
 * Copyright (C) 2019 Nalej - All Rights Reserved
 */

package ping

import (
	"github.com/nalej/device-controller/pkg/login_helper"
	"github.com/nalej/grpc-cluster-api-go"
	"github.com/nalej/grpc-common-go"
	"github.com/nalej/grpc-device-controller-go"
	"math"
)

type Manager struct {
	Threshold int
	// LoginHelper Helper
	ClusterAPILoginHelper *login_helper.LoginHelper
	// clusterAPIClient
	ClusterAPIClient grpc_cluster_api_go.DeviceManagerClient
}

func NewManager(threshold int, helper *login_helper.LoginHelper, client grpc_cluster_api_go.DeviceManagerClient ) Manager {
	return Manager{
		Threshold: threshold,
		ClusterAPILoginHelper: helper,
		ClusterAPIClient: client,
	}
}

func (m * Manager) Ping () (*grpc_common_go.Success, error) {
	return &grpc_common_go.Success{}, nil
}

func (m * Manager) SendRegisterPingToClusterAPI(ping * grpc_device_controller_go.RegisterLatencyRequest) error {
	return nil
}

func (m * Manager) RegisterPing (ping * grpc_device_controller_go.RegisterLatencyRequest) (* grpc_device_controller_go.RegisterLatencyResult, error) {
	result := grpc_device_controller_go.RegisterResult_OK
	if int(ping.Latency) > m.Threshold {
		result = grpc_device_controller_go.RegisterResult_LATENCY_CHECK_REQUIRED
	}
	return &grpc_device_controller_go.RegisterLatencyResult{
		Result: result,
	}, nil
}

func (m * Manager) SelectCluster(request * grpc_device_controller_go.SelectClusterRequest) (* grpc_device_controller_go.SelectedCluster, error) {

	pos := 0
	min := math.MaxInt32

	for i,latency := range request.Latencies {
		if int(latency) < min {
			min = int(latency)
			pos = i
		}
	}

	return &grpc_device_controller_go.SelectedCluster{
		ClusterIndex: int32(pos),
	}, nil
}