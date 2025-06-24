package models

// CardInfo 卡密信息模型
// 对应数据库中的CardInfo表，存储卡密详细信息
// CREATE TABLE CardInfo (Prefix_Name NVARCHAR (200) UNIQUE NOT NULL, Whom NVARCHAR (200), CardType NVARCHAR (200), FYI INTEGER, state NVARCHAR (20), Bind INTEGER, OpenNum INTEGER, LoginCount INTEGER, IP NVARCHAR (40), Remarks NVARCHAR (400), CreateData_ INTEGER, ActivateTime_ INTEGER, ExpiredTime_ INTEGER, LastLoginTime_ INTEGER, delstate INTEGER, Price REAL, cty BOOLEAN, ExpiredTime__ INTEGER, UnBindCount INTEGER DEFAULT (0), UnBindDeduct INTEGER DEFAULT (0), Attr_UnBindLimitTime INTEGER DEFAULT (0), Attr_UnBindDeductTime INTEGER DEFAULT (0), Attr_UnBindFreeCount INTEGER DEFAULT (0), Attr_UnBindMaxCount INTEGER DEFAULT (0), BindIP INTEGER DEFAULT (0), BanTime INTEGER DEFAULT (0), Owner TEXT DEFAULT (”), BindUser INTEGER DEFAULT (0), NowBindMachineNum INTEGER DEFAULT (0), BindMachineNum INTEGER DEFAULT (1), PCSign2 TEXT DEFAULT NULL, BanDurationTime INTEGER DEFAULT (0), GiveBackBanTime INTEGER DEFAULT (0), PICXCount INTEGER DEFAULT (0), LockBindPcsign INTEGER DEFAULT (0), 'LastRechargeTime' INTEGER DEFAULT (0), 'UserExtraData' BLOB DEFAULT (NULL))
type CardInfo struct {
	PrefixName           string  `gorm:"column:Prefix_Name;size:200;unique;not null" json:"prefix_name"`        // 卡号（唯一标识）
	Whom                 string  `gorm:"column:Whom;size:200" json:"whom"`                                      // 制卡人
	CardType             string  `gorm:"column:CardType;size:200" json:"card_type"`                             // 卡类型
	FYI                  int     `gorm:"column:FYI" json:"fyi"`                                                 // 点数
	State                string  `gorm:"column:state;size:20" json:"state"`                                     // 卡密状态（启用/禁用等）
	Bind                 int     `gorm:"column:Bind" json:"bind"`                                               // 绑定状态
	OpenNum              int     `gorm:"column:OpenNum" json:"open_num"`                                        // 多开数量
	LoginCount           int     `gorm:"column:LoginCount" json:"login_count"`                                  // 登录次数
	IP                   string  `gorm:"column:IP;size:40" json:"ip"`                                           // 最后登录IP地址
	Remarks              string  `gorm:"column:Remarks;size:400" json:"remarks"`                                // 备注信息
	CreateData_          int64   `gorm:"column:CreateData_" json:"create_data"`                                 // 创建时间戳
	ActivateTime_        int64   `gorm:"column:ActivateTime_" json:"activate_time"`                             // 激活时间戳
	ExpiredTime_         int64   `gorm:"column:ExpiredTime_" json:"expired_time"`                               // 有效期秒数
	LastLoginTime_       int64   `gorm:"column:LastLoginTime_" json:"last_login_time"`                          // 最后登录时间戳
	Delstate             int     `gorm:"column:delstate" json:"-"`                                              // 删除状态标记
	Price                float64 `gorm:"column:Price" json:"price"`                                             // 卡密价格
	Cty                  bool    `gorm:"column:cty" json:"cty"`                                                 // Cty标志
	ExpiredTime__        int64   `gorm:"column:ExpiredTime__" json:"expired_time_2"`                            // 过期时间戳
	UnBindCount          int     `gorm:"column:UnBindCount;default:0" json:"unbind_count"`                      // 解绑次数
	UnBindDeduct         int     `gorm:"column:UnBindDeduct;default:0" json:"unbind_deduct"`                    // 解绑扣除数量
	AttrUnBindLimitTime  int     `gorm:"column:Attr_UnBindLimitTime;default:0" json:"attr_unbind_limit_time"`   // 解绑限制时间
	AttrUnBindDeductTime int     `gorm:"column:Attr_UnBindDeductTime;default:0" json:"attr_unbind_deduct_time"` // 解绑扣除时间
	AttrUnBindFreeCount  int     `gorm:"column:Attr_UnBindFreeCount;default:0" json:"attr_unbind_free_count"`   // 解绑免费次数
	AttrUnBindMaxCount   int     `gorm:"column:Attr_UnBindMaxCount;default:0" json:"attr_unbind_max_count"`     // 解绑最大次数限制
	BindIP               int     `gorm:"column:BindIP;default:0" json:"bind_ip"`                                // 绑定IP状态
	BanTime              int     `gorm:"column:BanTime;default:0" json:"ban_time"`                              // 封禁时间
	Owner                string  `gorm:"column:Owner;type:text;default:''" json:"owner"`                        // 卡密所有者/充值账号
	BindUser             int     `gorm:"column:BindUser;default:0" json:"bind_user"`                            // 绑定用户状态
	NowBindMachineNum    int     `gorm:"column:NowBindMachineNum;default:0" json:"now_bind_machine_num"`        // 当前绑定机器数量
	BindMachineNum       int     `gorm:"column:BindMachineNum;default:1" json:"bind_machine_num"`               // 允许绑定的机器数量上限
	PCSign2              string  `gorm:"column:PCSign2;type:text" json:"pc_sign2"`                              // PC机器码签名2
	BanDurationTime      int     `gorm:"column:BanDurationTime;default:0" json:"ban_duration_time"`             // 封禁持续时间
	GiveBackBanTime      int     `gorm:"column:GiveBackBanTime;default:0" json:"give_back_ban_time"`            // 归还的封禁时间
	PICXCount            int     `gorm:"column:PICXCount;default:0" json:"picx_count"`                          // PICX计数器
	LockBindPcsign       int     `gorm:"column:LockBindPcsign;default:0" json:"lock_bind_pcsign"`               // 锁定绑定PC签名状态
	LastRechargeTime     int64   `gorm:"column:LastRechargeTime;default:0" json:"last_recharge_time"`           // 最后充值时间戳
	UserExtraData        []byte  `gorm:"column:UserExtraData;type:blob" json:"user_extra_data"`                 // 用户扩展数据（二进制存储）
}

// TableName 指定表名
func (CardInfo) TableName() string {
	return "CardInfo"
}
