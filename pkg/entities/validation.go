/*
 * Copyright (C) 2019 Nalej - All Rights Reserved
 */
package entities

import (
	"github.com/nalej/derrors"
	"github.com/nalej/grpc-device-controller-go"
)

const (
	emptyOrganizationId = "organization_id cannot be empty"
	emptyDeviceGroupId = "device_group_id cannot be empty"
	emptyDeviceId = "device_id cannot be empty"
	invalidMeasure = "Measure cannot be zero or less than zero"
	invalidListMeasure= "LatencyList cannot be empty"
)

func validLatency(latency * grpc_device_controller_go.Latency) derrors.Error {
	if latency.Measure <= 0 {
		return derrors.NewInvalidArgumentError(invalidMeasure)
	}
	return nil
}

func validLatencyList(list []*grpc_device_controller_go.Latency) derrors.Error {

	if list == nil || len(list.Latencies) == 0 {
		return derrors.NewInvalidArgumentError(invalidListMeasure)
	}
	for _, latency := range list.Latencies {
		err := validLatency(latency)
		if err != nil {
			return err
		}
	}
	return nil
}
func ValidRegisterLatencyRequest(latency * grpc_device_controller_go.RegisterLatencyRequest) derrors.Error{
	if latency.OrganizationId == "" {
		return derrors.NewInvalidArgumentError(emptyOrganizationId)
	}
	if latency.DeviceGroupId == "" {
		return derrors.NewInvalidArgumentError(emptyDeviceGroupId)
	}
	if latency.DeviceId == "" {
		return derrors.NewInvalidArgumentError(emptyDeviceId)
	}
	return validLatency(latency.Latency)
}

func ValidSelectClusterRequest(request * grpc_device_controller_go.SelectClusterRequest) derrors.Error {
	return validLatencyList(request.Latencies)
}