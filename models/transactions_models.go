package models

type Record struct {
	Key   []byte
	Value []byte
	Topic string
}

type Transaction struct {
	TxID            string  `json:"transaction_id"`
	UserID          string  `json:"user_id"`
	Amount          float32 `json:"amount"`
	Currency        string  `json:"currency"`
	TransactionType string  `json:"transaction_type"`
	Status          string  `json:"status"`
	Timestamp       string  `json:"timestamp"`
	PaymentMethod   string  `json:"payment_method"`
	CardNumber      string  `json:"card_number"`
	BankName        string  `json:"bank_name"`
	MerchantID      string  `json:"merchant_id"`
	MerchantName    string  `json:"merchant_name"`
	Location        string  `json:"location"`
	Description     string  `json:"description"`
	Category        string  `json:"category"`
	InvoiceNumber   string  `json:"invoice_number"`
	ReferenceID     string  `json:"reference_id"`
	TaxAmount       float64 `json:"tax_amount"`
	Discount        float64 `json:"discount"`
	NetAmount       float64 `json:"net_amount"`
	IPAddress       string  `json:"ip_address"`
	DeviceID        string  `json:"device_id"`
}
