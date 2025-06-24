package models

// MultiSoftware 软件位配置模型
// 对应主数据库中的MultiSoftware表，完全兼容现有数据库结构
type MultiSoftware struct {
	SoftwareName        string `gorm:"column:SoftwareName;primaryKey;unique;not null;default:''" json:"software_name"` // 软件位名称（主键）
	State               int    `gorm:"column:State;default:1" json:"state"`                                            // 状态：1=启用，0=禁用
	IDC                 string `gorm:"column:idc;default:''" json:"idc"`                                               // IDC信息/对应数据库文件
	Version             int    `gorm:"column:Version;default:0" json:"version"`                                        // 版本号
	ForceUpdate         int    `gorm:"column:ForceUpdate;default:0" json:"force_update"`                               // 强制更新
	DirectUrl           int    `gorm:"column:DirectUrl;default:0" json:"direct_url"`                                   // 直连URL
	Url                 string `gorm:"column:Url;default:''" json:"url"`                                               // URL地址
	RunExe              string `gorm:"column:RunExe;default:''" json:"run_exe"`                                        // 运行程序
	RunCmd              string `gorm:"column:RunCmd;default:''" json:"run_cmd"`                                        // 运行命令
	Notice              string `gorm:"column:Notice;default:''" json:"notice"`                                         // 公告
	HeartBeat           int    `gorm:"column:HeartBeat;default:0" json:"heart_beat"`                                   // 心跳
	HeartBeatDelayMs    int    `gorm:"column:HeartBeatDelayMs;default:0" json:"heart_beat_delay_ms"`                   // 心跳延迟毫秒
	HeartbeatBroken     int    `gorm:"column:HeartbeatBroken;default:0" json:"heartbeat_broken"`                       // 心跳断开
	BlacklistSec        int    `gorm:"column:BlacklistSec;default:0" json:"blacklist_sec"`                             // 黑名单秒数
	BlacklistCount      int    `gorm:"column:BlacklistCount;default:0" json:"blacklist_count"`                         // 黑名单计数
	DbgBanIP            int    `gorm:"column:DbgBanIP;default:0" json:"dbg_ban_ip"`                                    // 调试封IP
	DbgBanCard          int    `gorm:"column:DbgBanCard;default:0" json:"dbg_ban_card"`                                // 调试封卡
	DbgBanPCSign        int    `gorm:"column:DbgBanPCSign;default:0" json:"dbg_ban_pc_sign"`                           // 调试封机器码
	BeatSilenceSec      int    `gorm:"column:BeatSilenceSec;default:0" json:"beat_silence_sec"`                        // 心跳静默秒数
	BeatKeepOnline      int    `gorm:"column:BeatKeepOnline;default:0" json:"beat_keep_online"`                        // 心跳保持在线
	AutoBanSecond       int    `gorm:"column:AutoBanSecond;default:0" json:"auto_ban_second"`                          // 自动封禁秒数
	AutoBanCount        int    `gorm:"column:AutoBanCount;default:0" json:"auto_ban_count"`                            // 自动封禁计数
	AutoBanRule         int    `gorm:"column:AutoBanRule;default:0" json:"auto_ban_rule"`                              // 自动封禁规则
	ForbidLogin         int    `gorm:"column:ForbidLogin;default:0" json:"forbid_login"`                               // 禁止登录
	ForbidRegister      int    `gorm:"column:ForbidRegister;default:0" json:"forbid_register"`                         // 禁止注册
	ForbidTopUp         int    `gorm:"column:ForbidTopUp;default:0" json:"forbid_top_up"`                              // 禁止充值
	TrialCard           int    `gorm:"column:TrialCard;default:0" json:"trial_card"`                                   // 试用卡
	TrialTime           int    `gorm:"column:TrialTime;default:0" json:"trial_time"`                                   // 试用时间
	PCSignPlan          int    `gorm:"column:PCSignPlan;default:16" json:"pc_sign_plan"`                               // 机器码方案
	MixtureCardRecharge int    `gorm:"column:MixtureCardRecharge;default:0" json:"mixture_card_recharge"`              // 混合卡充值
	HeartBeatTimeoutMs  int    `gorm:"column:HeartBeatTimeoutMs;default:0" json:"heart_beat_timeout_ms"`               // 心跳超时毫秒
	CheckStrength       int    `gorm:"column:CheckStrength;default:1" json:"check_strength"`                           // 检查强度
	ShareSoftware       string `gorm:"column:ShareSoftware;default:''" json:"share_software"`                          // 共享软件
	RestrictionEx       int    `gorm:"column:RestrictionEx;default:0" json:"restriction_ex"`                           // 限制扩展
	InheritedAttribute  int64  `gorm:"column:InheritedAttribute;default:78519735425" json:"inherited_attribute"`       // 继承属性
}

// TableName 指定表名
func (MultiSoftware) TableName() string {
	return "MultiSoftware"
}

// IsActive 检查软件位是否激活
func (s *MultiSoftware) IsActive() bool {
	return s.State == 1
}

// GetStatusText 获取状态文本
func (ms *MultiSoftware) GetStatusText() string {
	if ms.IsActive() {
		return "启用"
	}
	return "禁用"
}

// SoftwareListResponse 软件列表响应
type SoftwareListResponse struct {
	Softwares []SoftwareAgentInfo `json:"softwares"` // 软件位列表
}

// SoftwareInfo 软件位信息
type SoftwareInfo struct {
	Name        string          `json:"name"`        // 软件位名称
	IDC         string          `json:"idc"`         // IDC信息
	State       int             `json:"state"`       // 状态
	AgentInfo   *SoftwareAgent  `json:"agent_info"`  // 代理信息
	Permissions map[string]bool `json:"permissions"` // 权限信息
}

// SoftwareAgent 软件位中的代理信息
type SoftwareAgent struct {
	Username    string          `json:"username"`    // 代理用户名
	Balance     float64         `json:"balance"`     // 账户余额
	TimeStock   int64           `json:"time_stock"`  // 库存时长
	Permissions map[string]bool `json:"permissions"` // 权限信息
	CardTypes   []string        `json:"card_types"`  // 可制作的卡类型
	Status      string          `json:"status"`      // 状态
	Expiration  string          `json:"expiration"`  // 到期时间
}

// getAgentStatus 获取代理状态文本
func GetAgentStatus(agent *Agent) string {
	if agent.Deltm != 0 {
		return "deleted"
	}
	if agent.Stat != 1 {
		return "disabled"
	}
	if agent.IsExpired() {
		return "expired"
	}
	return "active"
}

// GenerationCost 生成成本
type GenerationCost struct {
	BalanceDeducted float64 `json:"balance_deducted"` // 扣除的余额
	TimeDeducted    int64   `json:"time_deducted"`    // 扣除的时长
}

// SoftwareAgentInfo 软件位代理信息
type SoftwareAgentInfo struct {
	SoftwareName string          `json:"software_name"` // 软件位名称
	IDC          string          `json:"idc"`           // IDC信息
	State        int             `json:"state"`         // 状态
	AgentInfo    *SoftwareAgent  `json:"agent_info"`    // 代理信息
	Permissions  map[string]bool `json:"permissions"`   // 权限信息
}

// SoftwareStatistics 软件位统计信息
type SoftwareStatistics struct {
	TotalAgents    int64 `json:"total_agents"`    // 总代理数
	ActiveAgents   int64 `json:"active_agents"`   // 活跃代理数
	DisabledAgents int64 `json:"disabled_agents"` // 禁用代理数
	TotalCards     int64 `json:"total_cards"`     // 总卡密数
	ActiveCards    int64 `json:"active_cards"`    // 激活卡密数
	UsedCards      int64 `json:"used_cards"`      // 已使用卡密数
	CardTypes      int64 `json:"card_types"`      // 卡类型数量
}
