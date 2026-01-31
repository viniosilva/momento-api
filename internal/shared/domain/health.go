package domain

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
