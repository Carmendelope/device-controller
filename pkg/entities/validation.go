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
package entities

import (
	"github.com/nalej/derrors"
	"github.com/nalej/grpc-device-controller-go"
)

const (
	emptyOrganizationId = "organization_id cannot be empty"
	emptyDeviceGroupId  = "device_group_id cannot be empty"
	emptyDeviceId       = "device_id cannot be empty"
	invalidMeasure      = "Measure cannot be zero or less than zero"
	invalidListMeasure  = "LatencyList cannot be empty"
)

func ValidRegisterLatencyRequest(latency *grpc_device_controller_go.RegisterLatencyRequest) derrors.Error {
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

func ValidSelectClusterRequest(request *grpc_device_controller_go.SelectClusterRequest) derrors.Error {
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
