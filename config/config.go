package config

import (
	"fmt"
	"path/filepath"

	"gopkg.in/ini.v1"
)

// 服务器配置结构体
type ServerConfig struct {
	Host     string `ini:"服务器地址"` // 服务器监听地址
	Port     int    `ini:"端口"`    // 服务器监听端口
	DataPath string `ini:"数据库路径"` // 数据库路径
}

// JWT配置结构体
type JWTConfig struct {
	Secret string `ini:"JWT密钥"` // JWT签名密钥
}

// 应用程序配置结构体
type Config struct {
	Server ServerConfig `ini:"服务器设置"`
	JWT    JWTConfig    `ini:"认证设置"`
}

// 全局配置实例
var AppConfig *Config

// LoadConfig 加载配置文件
func LoadConfig(configPath string) (*Config, error) {
	cfg, err := ini.Load(configPath)
	if err != nil {
		return nil, fmt.Errorf("加载配置文件失败: %v", err)
	}

	config := &Config{}
	err = cfg.MapTo(config)
	if err != nil {
		return nil, fmt.Errorf("解析配置文件失败: %v", err)
	}

	// 设置默认值
	setDefaults(config)

	// 验证配置
	if err := validateConfig(config); err != nil {
		return nil, fmt.Errorf("配置验证失败: %v", err)
	}

	return config, nil
}

// InitConfig 初始化全局配置
func InitConfig() error {

	config, err := LoadConfig("config/SProtectAgentWeb.ini")
	if err != nil {
		return err
	}

	AppConfig = config
	return nil
}

// setDefaults 设置默认配置值
func setDefaults(config *Config) {
	// 服务器默认配置
	if config.Server.Host == "" {
		config.Server.Host = "0.0.0.0"
	}
	if config.Server.Port == 0 {
		config.Server.Port = 8080
	}

	// JWT默认配置
	if config.JWT.Secret == "" {
		config.JWT.Secret = "default-secret-key-please-change-in-production"
	}
}

// validateConfig 验证配置的有效性
func validateConfig(config *Config) error {
	// 验证端口范围
	if config.Server.Port < 1 || config.Server.Port > 65535 {
		return fmt.Errorf("服务器端口必须在1-65535范围内")
	}

	// 验证JWT密钥长度
	if len(config.JWT.Secret) < 16 {
		return fmt.Errorf("JWT密钥长度至少需要16个字符")
	}

	// 验证服务器地址
	if config.Server.Host == "" {
		return fmt.Errorf("服务器地址不能为空")
	}

	return nil
}

// GetServerAddress 获取完整的服务器地址
func GetServerAddress() string {
	return fmt.Sprintf("%s:%d", AppConfig.Server.Host, AppConfig.Server.Port)
}

// 以下为硬编码的系统配置，用户无需修改

// GetRunMode 获取运行模式
func GetRunMode() string {
	return "release" // 生产环境默认使用release模式
}

// GetDataPath 获取数据库路径
func GetDataPath() string {
	//判断路径最后字符串是不是\，如果不是，则添加\
	if AppConfig.Server.DataPath[len(AppConfig.Server.DataPath)-1:] != "\\" {
		AppConfig.Server.DataPath += "\\"
	}

	return AppConfig.Server.DataPath
}

func GetJWTSecret() string {
	return AppConfig.JWT.Secret
}

// GetJWTExpireTime 获取JWT过期时间(秒)
func GetJWTExpireTime() int {
	return 28800 // 8小时
}

// GetAppName 获取应用名称
func GetAppName() string {
	return "SProtectAgentWeb"
}

// GetAppVersion 获取应用版本
func GetAppVersion() string {
	return "1.0.0"
}

// GetDatabasePath 获取指定数据库的完整路径
func GetDatabasePath(dbName string) string {
	return filepath.Join(GetDataPath(), dbName)
}

// GetMainDatabasePath 获取主数据库路径
func GetMainDatabasePath() string {
	return GetDatabasePath("idc.db")
}

// GetSoftwareDatabasePath 获取软件位数据库路径
func GetSoftwareDatabasePath(softwareName string) string {
	return GetDatabasePath(fmt.Sprintf("idc_%s.db", softwareName))
}

// GetAuditDatabasePath 获取审计日志数据库路径
func GetAuditDatabasePath() string {
	return GetDatabasePath("audit_logs.db")
}

// IsAuditEnabled 是否启用审计日志
func IsAuditEnabled() bool {
	return true
}

// GetAuditCleanDays 获取审计日志保留天数
func GetAuditCleanDays() int {
	return 90
}
