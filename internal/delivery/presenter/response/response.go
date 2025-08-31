package response

import "fmt"

type Meta struct {
	TotalData int64 `json:"total_data"`
	TotalPage int64 `json:"total_page"`
	Page      int64 `json:"page"`
	PageSize  int64 `json:"page_size"`
}

type TxError struct {
	Op  string // "commit", "rollback", "operation"
	Err error
}

func (e *TxError) Error() string {
	return fmt.Sprintf("transaction %s error: %v", e.Op, e.Err)
}

func (e *TxError) Unwrap() error {
	return e.Err
}
