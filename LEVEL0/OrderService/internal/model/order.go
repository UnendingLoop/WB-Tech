package model

type Order struct {
	OrderUID    string `gorm:"primaryKey" json:"order_uid"`
	TrackNumber string `gorm:"not null" json:"track_number"`
	Entry       string `gorm:"not null" json:"entry"`

	DeliveryID uint     `json:"-"`
	Delivery   Delivery `gorm:"foreignKey:DeliveryID" json:"delivery"`

	PaymentID uint    `json:"-"`
	Payment   Payment `gorm:"foreignKey:PaymentID" json:"payment"`

	Items []Item `gorm:"foreignKey:OrderID" json:"items"`

	Locale            string `gorm:"not null" json:"locale"`
	InternalSignature string `gorm:"not null" json:"internal_signature"`
	CustomerID        string `gorm:"not null" json:"customer_id"`
	DeliveryService   string `gorm:"not null" json:"delivery_service"`
	ShardKey          string `gorm:"not null" json:"shardkey"`
	SMID              int    `gorm:"not null" json:"sm_id"`
	DateCreated       string `gorm:"not null" json:"date_created"`
	OofShard          string `gorm:"not null" json:"oof_shard"`
}

type Delivery struct {
	ID      uint   `gorm:"primaryKey" json:"-"`
	Name    string `gorm:"not null" json:"name"`
	Phone   string `gorm:"not null" json:"phone"`
	Zip     string `gorm:"not null" json:"zip"`
	City    string `gorm:"not null" json:"city"`
	Address string `gorm:"not null" json:"address"`
	Region  string `gorm:"not null" json:"region"`
	Email   string `gorm:"not null" json:"email"`
}

type Payment struct {
	ID           uint   `gorm:"primaryKey" json:"-"`
	Transaction  string `gorm:"not null" json:"transaction"`
	RequestID    string `gorm:"not null" json:"request_id"`
	Currency     string `gorm:"not null" json:"currency"`
	Provider     string `gorm:"not null" json:"provider"`
	Amount       uint   `gorm:"not null" json:"amount"`
	PaymentDT    int64  `gorm:"not null" json:"payment_dt"` // UNIX timestamp
	Bank         string `gorm:"not null" json:"bank"`
	DeliveryCost uint   `gorm:"not null" json:"delivery_cost"`
	GoodsTotal   uint   `gorm:"not null" json:"goods_total"`
	CustomFee    uint   `gorm:"not null" json:"custom_fee"`
}

type Item struct {
	ID      uint   `gorm:"primaryKey" json:"-"`
	OrderID string `gorm:"not null" json:"-"` // foreign key to Order.OrderUID

	ChrtID      uint   `gorm:"not null" json:"chrt_id"`
	TrackNumber string `gorm:"not null" json:"track_number"`
	Price       uint   `gorm:"not null" json:"price"`
	RID         string `gorm:"not null" json:"rid"`
	Name        string `gorm:"not null" json:"name"`
	Sale        uint   `gorm:"not null" json:"sale"`
	Size        string `gorm:"not null" json:"size"`
	TotalPrice  uint   `gorm:"not null" json:"total_price"`
	NMID        uint   `gorm:"not null" json:"nm_id"`
	Brand       string `gorm:"not null" json:"brand"`
	Status      int    `gorm:"not null" json:"status"`
}
