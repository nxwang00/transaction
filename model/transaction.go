package model

type TransactionReq struct {
	Origin        string `json:"origin,omitempty"`
	User_ID       int    `json:"user_id,omitempty"`
	Amount        string `json:"amount,omitempty"`
	Op_Type       string `json:"op_type,omitempty"`
	Registered_At string `json:"registered_at,omitempty"`
}

type Transaction struct {
	ID            int    `json:"id,omitempty"`
	Origin        string `json:"origin,omitempty"`
	User_ID       int    `json:"user_id,omitempty"`
	Amount        string `json:"amount,omitempty"`
	Op_Type       string `json:"op_type,omitempty"`
	Registered_At string `json:"registered_at,omitempty"`
}

type PageInfo struct {
	Page_Number int
	Page_Size   int
}