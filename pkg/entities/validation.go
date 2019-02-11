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
	if latency.Latency <= 0 {
		return derrors.NewInvalidArgumentError(invalidMeasure)
	}
	return nil
}

func ValidSelectClusterRequest(request * grpc_device_controller_go.SelectClusterRequest) derrors.Error {
	if request.OrganizationId == "" {
		return derrors.NewInvalidArgumentError(emptyOrganizationId)
	}
	if request.DeviceGroupId == "" {
		return derrors.NewInvalidArgumentError(emptyDeviceGroupId)
	}
	if request.DeviceId == "" {
		return derrors.NewInvalidArgumentError(emptyDeviceId)
	}
	for _, latency := range request.Latencies {
		if latency <= 0 {
			return derrors.NewInvalidArgumentError(invalidMeasure)
		}
	}
	return nil

}