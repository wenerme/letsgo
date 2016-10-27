package bg

type HealthCheckerFunc func() error

func (c HealthCheckerFunc) Check() error {
	return c()
}
