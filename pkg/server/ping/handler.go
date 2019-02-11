/*
 * Copyright (C) 2019 Nalej - All Rights Reserved
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

func NewHandler (manager Manager) *Handler {
	return &Handler {manager}
}

func (h * Handler) Ping (ctx context.Context,in *grpc_common_go.Empty) (*grpc_common_go.Success, error) {
	return h.Manager.Ping()
}

	func (h * Handler) RegisterLatency (ctx context.Context, ping * grpc_device_controller_go.RegisterLatencyRequest) (* grpc_device_controller_go.RegisterLatencyResult, error)  {
	err := entities.ValidRegisterLatencyRequest(ping)
	if err != nil {
		return nil, err
	}
	return h.Manager.RegisterPing(ping)
}

func (h * Handler) SelectCluster(ctx context.Context, request * grpc_device_controller_go.SelectClusterRequest) (* grpc_device_controller_go.SelectedCluster, error) {
	err := entities.ValidSelectClusterRequest(request)
	if err != nil {
		return nil, err
	}
	return h.Manager.SelectCluster(request)
}