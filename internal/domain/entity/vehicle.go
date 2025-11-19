package entity

type Vehicle struct {
	ID        string
	TenantID  string
	Plate     string
	Model     string
	ColdChain bool // 是否冷链
	IceBoxNo  string
	CreatedAt int64
	UpdatedAt int64
}
