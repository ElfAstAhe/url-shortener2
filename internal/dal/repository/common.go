package repository

type DMLResult struct {
	Err    error
	ID     string
	Entity string
}

func NewDMLResult(err error, entity string, ID string) *DMLResult {
	return &DMLResult{
		Err:    err,
		ID:     ID,
		Entity: entity,
	}
}
