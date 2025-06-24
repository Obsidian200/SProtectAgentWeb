package services

import (
	"fmt"
	"web-agent-backend/database"
	"web-agent-backend/models"
)

// CardTypeService 卡密类型服务
type CardTypeService struct {
	dbManager       *database.DatabaseManager
	softwareService *SoftwareService
}

// NewCardTypeService 创建卡密类型服务实例
func NewCardTypeService(dbManager *database.DatabaseManager, softwareService *SoftwareService) *CardTypeService {
	return &CardTypeService{
		dbManager:       dbManager,
		softwareService: softwareService,
	}
}

// GetCardTypeList 获取卡密类型列表（根据代理权限过滤）
func (s *CardTypeService) GetCardTypeList(software, agentName string) ([]models.CardType, error) {
	// 1. 通过SoftwareService获取代理信息
	agent, err := s.softwareService.GetSoftwareAgentByName(software, agentName)
	if err != nil {
		return nil, fmt.Errorf("获取代理信息失败: %v", err)
	}

	// 2. 获取允许的卡密类型名称列表（已在SoftwareService中处理）
	allowedTypes := agent.CardTypeAuthNameArray

	// 如果没有权限，返回空列表
	if len(allowedTypes) == 0 {
		return []models.CardType{}, nil
	}

	// 3. 获取软件位数据库连接
	db, err := s.dbManager.GetSoftwareDB(software)
	if err != nil {
		return nil, fmt.Errorf("获取数据库连接失败: %v", err)
	}

	// 4. 根据权限查询卡密类型详情
	var cardTypes []models.CardType
	err = db.Where("Name IN ?", allowedTypes).Find(&cardTypes).Error
	if err != nil {
		return nil, fmt.Errorf("查询卡密类型失败: %v", err)
	}

	// 5. 根据代理的总利率计算实际价格
	for i := range cardTypes {
		// 使用CardType模型的CalculatePrice方法计算实际价格
		// 例如：原价10元，总利率132% -> 10 × 1.32 = 13.20元（保留2位小数，截断）
		cardTypes[i].Price = cardTypes[i].CalculatePrice(agent.TatalParities)
	}

	return cardTypes, nil
}

// GetCardTypeByName 根据名称获取卡密类型
func (s *CardTypeService) GetCardTypeByName(software, name string) (*models.CardType, error) {
	// 获取软件位数据库连接
	db, err := s.dbManager.GetSoftwareDB(software)
	if err != nil {
		return nil, fmt.Errorf("获取数据库连接失败: %v", err)
	}

	// 查询指定名称的卡密类型
	var cardType models.CardType
	err = db.Where("Name = ?", name).First(&cardType).Error
	if err != nil {
		return nil, fmt.Errorf("查询卡密类型失败: %v", err)
	}

	return &cardType, nil
}

// CreateCardType 创建卡密类型
func (s *CardTypeService) CreateCardType(software string, cardType *models.CardType) error {
	// 获取软件位数据库连接
	db, err := s.dbManager.GetSoftwareDB(software)
	if err != nil {
		return fmt.Errorf("获取数据库连接失败: %v", err)
	}

	// 创建卡密类型
	err = db.Create(cardType).Error
	if err != nil {
		return fmt.Errorf("创建卡密类型失败: %v", err)
	}

	return nil
}

// UpdateCardType 更新卡密类型
func (s *CardTypeService) UpdateCardType(software, name string, updates map[string]interface{}) error {
	// 获取软件位数据库连接
	db, err := s.dbManager.GetSoftwareDB(software)
	if err != nil {
		return fmt.Errorf("获取数据库连接失败: %v", err)
	}

	// 更新卡密类型
	result := db.Model(&models.CardType{}).
		Where("Name = ?", name).
		Updates(updates)

	if result.Error != nil {
		return fmt.Errorf("更新卡密类型失败: %v", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("卡密类型不存在")
	}

	return nil
}

// DeleteCardType 删除卡密类型
func (s *CardTypeService) DeleteCardType(software, name string) error {
	// 获取软件位数据库连接
	db, err := s.dbManager.GetSoftwareDB(software)
	if err != nil {
		return fmt.Errorf("获取数据库连接失败: %v", err)
	}

	// 检查是否为系统卡密类型（不可删除）
	var cardType models.CardType
	err = db.Where("Name = ?", name).First(&cardType).Error
	if err != nil {
		return fmt.Errorf("查询卡密类型失败: %v", err)
	}

	if cardType.CannotBeChanged {
		return fmt.Errorf("系统卡密类型不可删除")
	}

	// 删除卡密类型
	result := db.Where("Name = ?", name).Delete(&models.CardType{})
	if result.Error != nil {
		return fmt.Errorf("删除卡密类型失败: %v", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("卡密类型不存在")
	}

	return nil
}
