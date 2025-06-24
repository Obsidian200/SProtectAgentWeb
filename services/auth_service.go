package services

import (
	"SProtectAgentWeb/database"
	"SProtectAgentWeb/models"
	"SProtectAgentWeb/util"
	"fmt"
	"log"

	"gorm.io/gorm"
)

// AuthService 认证服务
// 专注认证相关的业务逻辑
type AuthService struct {
	dbManager *database.DatabaseManager
}

// NewAuthService 创建认证服务实例
func NewAuthService(dbManager *database.DatabaseManager) *AuthService {
	return &AuthService{
		dbManager: dbManager,
	}
}

// CreateUserSession 创建用户会话
func (s *AuthService) CreateUserSession(username, password, ipAddress, userAgent string) (*models.UserSession, error) {
	// 验证用户凭据并获取所有可访问的软件位
	softwareAgents, softwareList, err := s.validateAndGetAllSoftware(username, password)
	if err != nil {
		return nil, err
	}

	// 创建用户会话
	userSession := &models.UserSession{
		Username:          username,
		Password:          password,
		IPAddress:         ipAddress,
		SoftwareList:      softwareList,
		SoftwareAgentInfo: softwareAgents,
	}

	return userSession, nil
}

func (s *AuthService) GetUserInfo(username, password string) (*models.UserSession, error) {
	// 验证用户凭据并获取所有可访问的软件位
	softwareAgents, softwareList, err := s.validateAndGetAllSoftware(username, password)
	if err != nil {
		return nil, err
	}

	// 创建用户会话
	userSession := &models.UserSession{
		Username:          username,
		SoftwareList:      softwareList,
		SoftwareAgentInfo: softwareAgents,
	}

	return userSession, nil
}

// ChangePassword 修改密码
func (s *AuthService) ChangePassword(username, software, oldPassword, newPassword string) error {
	// 获取软件位数据库连接
	softwareDB, err := s.dbManager.GetSoftwareDB(software)
	if err != nil {
		return fmt.Errorf("连接软件位数据库失败: %v", err)
	}

	// 验证旧密码
	var agent models.Agent
	err = softwareDB.Where("User = ? AND Password = ?", username, oldPassword).First(&agent).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("旧密码错误")
		}
		return fmt.Errorf("验证用户信息失败: %v", err)
	}

	// 更新密码
	err = softwareDB.Model(&models.Agent{}).Where("User = ?", username).Update("Password", newPassword).Error
	if err != nil {
		return fmt.Errorf("更新密码失败: %v", err)
	}

	return nil
}

// GetAgentInfo 获取指定软件位的代理信息
func (s *AuthService) GetAgentInfo(username, software string) (*models.Agent, error) {
	softwareDB, err := s.dbManager.GetSoftwareDB(software)
	if err != nil {
		return nil, fmt.Errorf("连接软件位数据库失败: %v", err)
	}

	var agent models.Agent
	err = softwareDB.Where("User = ?", username).First(&agent).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("代理不存在")
		}
		return nil, fmt.Errorf("查询代理信息失败: %v", err)
	}

	return &agent, nil
}

// ===== 私有方法 =====

// validateAndGetAllSoftware 验证用户凭据并获取所有软件位信息
func (s *AuthService) validateAndGetAllSoftware(username, password string) (map[string]*models.Agent, []string, error) {
	softwareDBMap, err := s.dbManager.GetAllSoftwareDB()

	if err != nil {
		return nil, nil, fmt.Errorf("获取软件位数据库失败: %v", err)
	}

	log.Println("可用软件位:", softwareDBMap)

	softwareAgentMap := make(map[string]*models.Agent)
	softwareArr := []string{}

	// 在每个软件位中查找匹配的代理
	for softwareName, db := range softwareDBMap {
		var agent models.Agent
		err = db.Where("User = ? AND Password = ?", username, password).First(&agent).Error
		if err != nil {
			continue // 该软件位中未找到匹配的账户
		}

		// 验证代理是否有效
		if !agent.IsValid() {
			continue
		}

		// 直接使用util层函数解析卡类型权限
		agent.CardTypeAuthNameArray = util.ParseBracketList(agent.CardTypeAuthName)

		// 存储该用户在这个软件位中的代理信息
		softwareAgentMap[softwareName] = &agent
		softwareArr = append(softwareArr, softwareName)

	}

	if len(softwareArr) == 0 {
		return nil, nil, fmt.Errorf("用户名或密码错误，或账户已被禁用")
	}

	return softwareAgentMap, softwareArr, nil
}
