package models

// CardType 卡密类型模型
type CardType struct {
	Name                 string  `gorm:"column:Name;type:NVARCHAR(200);unique;not null" json:"name"`                         // 卡密类型名称
	Prefix               string  `gorm:"column:Prefix;type:NVARCHAR(50)" json:"prefix"`                                      // 前缀
	Duration             int     `gorm:"column:Duration;type:INTEGER" json:"duration"`                                       // 时长
	FYI                  int     `gorm:"column:FYI;type:INTEGER" json:"fyi"`                                                 // FYI字段
	Price                float64 `gorm:"column:Price;type:REAL" json:"price"`                                                // 价格
	Param                string  `gorm:"column:Param;type:NVARCHAR(400)" json:"param"`                                       // 参数
	Bind                 int     `gorm:"column:Bind;type:INTEGER" json:"bind"`                                               // 绑定
	OpenNum              int     `gorm:"column:OpenNum;type:INTEGER" json:"open_num"`                                        // 开通数量
	Remarks              string  `gorm:"column:Remarks;type:NVARCHAR(400)" json:"remarks"`                                   // 备注
	CannotBeChanged      bool    `gorm:"column:CannotBeChanged;type:BOOLEAN" json:"cannot_be_changed"`                       // 不可更改
	AttrUnBindLimitTime  int     `gorm:"column:Attr_UnBindLimitTime;type:INTEGER;default:0" json:"attr_unbind_limit_time"`   // 解绑限制时间
	AttrUnBindDeductTime int     `gorm:"column:Attr_UnBindDeductTime;type:INTEGER;default:0" json:"attr_unbind_deduct_time"` // 解绑扣除时间
	AttrUnBindFreeCount  int     `gorm:"column:Attr_UnBindFreeCount;type:INTEGER;default:0" json:"attr_unbind_free_count"`   // 解绑免费次数
	AttrUnBindMaxCount   int     `gorm:"column:Attr_UnBindMaxCount;type:INTEGER;default:0" json:"attr_unbind_max_count"`     // 解绑最大次数
	BindIP               int     `gorm:"column:BindIP;type:INTEGER;default:0" json:"bind_ip"`                                // 绑定IP
	BindMachineNum       int     `gorm:"column:BindMachineNum;type:INTEGER;default:1" json:"bind_machine_num"`               // 绑定机器数量
	LockBindPcsign       int     `gorm:"column:LockBindPcsign;type:INTEGER;default:0" json:"lock_bind_pcsign"`               // 锁定绑定PC签名
}

// TableName 指定表名
func (CardType) TableName() string {
	return "CardType"
}

// CalculatePrice 根据利率计算价格
func (ct *CardType) CalculatePrice(parities float64) float64 {
	if ct.Price <= 0 {
		return 0.0
	}

	// 计算实际价格：原价 × (利率 / 100)
	actualPrice := ct.Price * (parities / 100.0)

	// 保留2位小数（截断，不四舍五入）
	return float64(int(actualPrice*100)) / 100.0
}
