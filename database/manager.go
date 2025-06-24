package database

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DatabaseManager 数据库连接管理器
// 负责管理所有数据库连接，包括主数据库、软件位数据库和审计日志数据库
type DatabaseManager struct {
	connections map[string]*gorm.DB // 存储所有数据库连接
	mutex       sync.RWMutex        // 读写锁，保证并发安全
	config      *gorm.Config        // GORM配置
	dataPath    string              // 数据库文件存储路径
}

// NewDatabaseManager 创建数据库管理器实例
// dataPath: 数据库文件存储路径
// 返回: 数据库管理器实例
func NewDatabaseManager(dataPath string) *DatabaseManager {
	return &DatabaseManager{
		connections: make(map[string]*gorm.DB),
		config: &gorm.Config{
			Logger: logger.Default.LogMode(logger.Silent), // 设置日志级别
		},
		dataPath: dataPath,
	}
}

// GetDafaultDB 获取主数据库连接
// 主数据库文件：idc.db
// 该数据库包含：
// 1. MultiSoftware 表：所有软件位的配置信息
// 2. 默认软件位（"默认软件"）的业务数据：Agent、CardInfo、CardType等表
// 返回: GORM数据库连接和可能的错误
func (dm *DatabaseManager) GetDafaultDB() (*gorm.DB, error) {
	return dm.getConnection("默认软件", "idc.db")
}

// GetSoftwareDB 获取指定软件位的数据库连接
// software: 软件位名称
// 如果是默认软件位（"默认软件"），返回主数据库 idc.db
// 否则返回对应的 idc_软件位名称.db
// 返回: GORM数据库连接和可能的错误
func (dm *DatabaseManager) GetSoftwareDB(software string) (*gorm.DB, error) {
	// 如果是默认软件位，使用主数据库 idc.db
	if software == "默认软件" {
		return dm.GetDafaultDB()
	}

	// 其他软件位使用独立的数据库文件
	dbName := fmt.Sprintf("idc_%s.db", software)
	return dm.getConnection(software, dbName)
}

// GetAllSoftwareDB 获取所有软件位数据库连接列表
// 返回: 软件位名称到数据库连接的映射
func (dm *DatabaseManager) GetAllSoftwareDB() (map[string]*gorm.DB, error) {
	// 首先获取主数据库以读取所有软件位信息
	mainDB, err := dm.GetDafaultDB()
	if err != nil {
		return nil, fmt.Errorf("获取主数据库失败: %v", err)
	}

	// 查询所有启用的软件位
	var softwares []struct {
		SoftwareName string `gorm:"column:SoftwareName"`
	}

	err = mainDB.Table("MultiSoftware").Where("State = ?", 1).Find(&softwares).Error
	if err != nil {
		return nil, fmt.Errorf("查询软件位列表失败: %v", err)
	}

	// 为每个软件位建立连接
	result := make(map[string]*gorm.DB)
	for _, software := range softwares {
		db, err := dm.GetSoftwareDB(software.SoftwareName)
		if err != nil {
			// 记录错误但继续处理其他软件位
			fmt.Printf("警告: 无法连接软件位 [%s] 的数据库: %v\n", software.SoftwareName, err)
			continue
		}
		result[software.SoftwareName] = db
	}

	return result, nil
}

// getConnection 获取或创建数据库连接
// key: 连接标识符
// filename: 数据库文件名
// 返回: GORM数据库连接和可能的错误
func (dm *DatabaseManager) getConnection(key, filename string) (*gorm.DB, error) {
	// 先尝试读锁获取现有连接
	dm.mutex.RLock()
	if db, exists := dm.connections[key]; exists {
		dm.mutex.RUnlock()
		return db, nil
	}
	dm.mutex.RUnlock()

	// 需要创建新连接，使用写锁
	dm.mutex.Lock()
	defer dm.mutex.Unlock()

	// 再次检查（双重检查锁定模式）
	if db, exists := dm.connections[key]; exists {
		return db, nil
	}

	// 构建数据库文件路径
	dbPath := filepath.Join(dm.dataPath, filename)

	// 检查文件是否存在且非空
	if fileInfo, err := os.Stat(dbPath); err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("数据库文件不存在，路径: %s", dbPath)
		}
		return nil, err
	} else if fileInfo.Size() == 0 {
		return nil, fmt.Errorf("数据库文件为空")
	}

	// 创建SQLite连接
	db, err := gorm.Open(sqlite.Open(dbPath), dm.config)
	if err != nil {
		return nil, fmt.Errorf("连接数据库[%s] 失败: %v", filename, err)
	}

	// 配置SQLite连接参数
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("获取数据库实例失败: %v", err)
	}

	// SQLite连接池配置
	sqlDB.SetMaxOpenConns(1)    // SQLite建议使用单连接
	sqlDB.SetMaxIdleConns(1)    // 最大空闲连接数
	sqlDB.SetConnMaxLifetime(0) // 连接最大生存时间

	// 存储连接
	dm.connections[key] = db

	return db, nil
}

// CloseAll 关闭所有数据库连接
// 在程序退出时调用，确保所有连接被正确关闭
// 返回: 可能的错误
func (dm *DatabaseManager) CloseAll() error {
	dm.mutex.Lock()
	defer dm.mutex.Unlock()

	var lastErr error
	for key, db := range dm.connections {
		if sqlDB, err := db.DB(); err == nil {
			if err := sqlDB.Close(); err != nil {
				lastErr = fmt.Errorf("关闭数据库连接[%s] 失败: %v", key, err)
			}
		}
	}

	// 清空连接映射
	dm.connections = make(map[string]*gorm.DB)

	return lastErr
}

// GetConnectionCount 获取当前连接数量
// 用于监控和调试
// 返回: 当前连接数量
func (dm *DatabaseManager) GetConnectionCount() int {
	dm.mutex.RLock()
	defer dm.mutex.RUnlock()
	return len(dm.connections)
}
