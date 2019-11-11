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
	"context"
	"github.com/nalej/device-controller/pkg/entities"
	"github.com/nalej/grpc-common-go"
	"github.com/nalej/grpc-device-controller-go"
)

type Handler struct {
	Manager Manager
}

func NewHandler(manager Manager) *Handler {
	return &Handler{manager}
}

func (h *Handler) Ping(ctx context.Context, in *grpc_common_go.Empty) (*grpc_common_go.Success, error) {
	return h.Manager.Ping()
}

func (h *Handler) RegisterLatency(ctx context.Context, ping *grpc_device_controller_go.RegisterLatencyRequest) (*grpc_device_controller_go.RegisterLatencyResult, error) {
	err := entities.ValidRegisterLatencyRequest(ping)
	if err != nil {
		return nil, err
	}
	return h.Manager.RegisterPing(ping)
}

func (h *Handler) SelectCluster(ctx context.Context, request *grpc_device_controller_go.SelectClusterRequest) (*grpc_device_controller_go.SelectedCluster, error) {
	err := entities.ValidSelectClusterRequest(request)
	if err != nil {
		return nil, err
	}
	return h.Manager.SelectCluster(request)
}
