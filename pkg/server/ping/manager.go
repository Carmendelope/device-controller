/*
 * Copyright 2019 Nalej
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package ping

import (
	"github.com/nalej/device-controller/pkg/login_helper"
	"github.com/nalej/grpc-cluster-api-go"
	"github.com/nalej/grpc-common-go"
	"github.com/nalej/grpc-device-controller-go"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	grpc_status "google.golang.org/grpc/status"
	"math"
)

type Manager struct {
	Threshold int
	// LoginHelper Helper
	ClusterAPILoginHelper *login_helper.LoginHelper
	// clusterAPIClient
	ClusterAPIClient grpc_cluster_api_go.DeviceManagerClient
}

func NewManager(threshold int, helper *login_helper.LoginHelper, client grpc_cluster_api_go.DeviceManagerClient) Manager {
	return Manager{
		Threshold:             threshold,
		ClusterAPILoginHelper: helper,
		ClusterAPIClient:      client,
	}
}

func (m *Manager) Ping() (*grpc_common_go.Success, error) {
	return &grpc_common_go.Success{}, nil
}

func (m *Manager) sendRegisterPingToClusterAPI(ping *grpc_device_controller_go.RegisterLatencyRequest) error {

	ctx, cancel := m.ClusterAPILoginHelper.GetContext()
	defer cancel()

	_, err := m.ClusterAPIClient.RegisterLatency(ctx, ping)
	if err != nil {
		st := grpc_status.Convert(err).Code()
		if st == codes.Unauthenticated {
			errLogin := m.ClusterAPILoginHelper.RerunAuthentication()
			if errLogin != nil {
				log.Error().Err(errLogin).Msg("error during reauthentication")
			}
			ctx2, cancel2 := m.ClusterAPILoginHelper.GetContext()
			defer cancel2()
			_, err = m.ClusterAPIClient.RegisterLatency(ctx2, ping)
		} else {
			log.Error().Err(err).Str("OrganizationId", ping.OrganizationId).Str("deviceGroupId", ping.DeviceGroupId).Str("deviceId", ping.DeviceId).Msgf("error recording latencies")
		}
	}

	return nil
}

func (m *Manager) RegisterPing(ping *grpc_device_controller_go.RegisterLatencyRequest) (*grpc_device_controller_go.RegisterLatencyResult, error) {
	result := grpc_device_controller_go.RegisterResult_OK
	if int(ping.Latency) > m.Threshold {
		result = grpc_device_controller_go.RegisterResult_LATENCY_CHECK_REQUIRED
	}

	go m.sendRegisterPingToClusterAPI(ping)

	return &grpc_device_controller_go.RegisterLatencyResult{
		Result: result,
	}, nil
}

func (m *Manager) SelectCluster(request *grpc_device_controller_go.SelectClusterRequest) (*grpc_device_controller_go.SelectedCluster, error) {

	pos := 0
	min := math.MaxInt32

	for i, latency := range request.Latencies {
		if int(latency) < min {
			min = int(latency)
			pos = i
		}
	}

	return &grpc_device_controller_go.SelectedCluster{
		ClusterIndex: int32(pos),
	}, nil
}
