package domain

// HealthStatusEnum represents the possible health status values
// @enum ok error
type HealthStatusEnum string

const (
	HealthStatusValueOk    HealthStatusEnum = "ok"
	HealthStatusValueError HealthStatusEnum = "error"
)

type HealthStatus struct {
	Status HealthStatusEnum
}

func HealthStatusOk() HealthStatus {
	return HealthStatus{
		Status: HealthStatusValueOk,
	}
}

func HealthStatusError() HealthStatus {
	return HealthStatus{
		Status: HealthStatusValueError,
	}
}
