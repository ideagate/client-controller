package model

import v1 "github.com/bayu-aditya/ideagate/backend/model/gen-go/client/controller/v1"

func ConvertPodStatus(status string) v1.PodStatus {
	switch status {
	case "Pending":
		return v1.PodStatus_POD_STATUS_PENDING
	case "Running":
		return v1.PodStatus_POD_STATUS_RUNNING
	case "Succeeded":
		return v1.PodStatus_POD_STATUS_SUCCEEDED
	case "Failed":
		return v1.PodStatus_POD_STATUS_FAILED
	case "Unknown":
		return v1.PodStatus_POD_STATUS_UNSPECIFIED
	}

	return v1.PodStatus_POD_STATUS_UNSPECIFIED
}
