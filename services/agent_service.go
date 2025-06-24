package services

import (
	"fmt"
	"strings"
	"web-agent-backend/database"
	"web-agent-backend/models"
	"web-agent-backend/util"
)

// AgentService 代理服务
// 只依赖DatabaseManager，专注代理业务逻辑
type AgentService struct {
	dbManager *database.DatabaseManager
}

// NewAgentService 创建代理服务实例
func NewAgentService(dbManager *database.DatabaseManager) *AgentService {
	return &AgentService{
		dbManager: dbManager,
	}
}

// ===== 私有方法 =====

// GetSubAgents 获取指定代理的所有子代理
func (s *AgentService) GetSubAgents(software, username string) ([]*models.Agent, error) {
	// 获取软件位数据库连接
	softwareDB, err := s.dbManager.GetSoftwareDB(software)
	if err != nil {
		return nil, fmt.Errorf("连接软件位数据库失败: %v", err)
	}

	// 查询所有代理，只包含启用(Stat=0)或禁用(Stat=1)的代理
	var allAgents []*models.Agent
	err = softwareDB.Where("deltm = ? AND Stat IN (0, 1)", 0).Find(&allAgents).Error
	if err != nil {
		return nil, fmt.Errorf("查询代理列表失败: %v", err)
	}

	// 筛选出子代理
	var subAgents []*models.Agent
	for _, agent := range allAgents {
		// 跳过自己
		if agent.User == username {
			continue
		}

		// 检查是否是子代理
		if agent.IsChildOf(username) {
			// 解析卡类型权限
			agent.CardTypeAuthNameArray = util.ParseBracketList(agent.CardTypeAuthName)
			subAgents = append(subAgents, agent)
		}
	}

	return subAgents, nil
}

// GetDirectSubAgents 获取指定代理的直接子代理
func (s *AgentService) GetDirectSubAgents(software, username string) ([]*models.Agent, error) {
	// 获取软件位数据库连接
	softwareDB, err := s.dbManager.GetSoftwareDB(software)
	if err != nil {
		return nil, fmt.Errorf("连接软件位数据库失败: %v", err)
	}

	// 查询所有代理，只包含启用(Stat=0)或禁用(Stat=1)的代理
	var allAgents []*models.Agent
	err = softwareDB.Where("deltm = ? AND Stat IN (0, 1)", 0).Find(&allAgents).Error
	if err != nil {
		return nil, fmt.Errorf("查询代理列表失败: %v", err)
	}

	// 筛选出直接子代理
	var directSubAgents []*models.Agent
	for _, agent := range allAgents {
		// 跳过自己
		if agent.User == username {
			continue
		}

		// 检查是否是直接子代理
		if agent.IsDirectChildOf(username) {
			// 解析卡类型权限
			agent.CardTypeAuthNameArray = util.ParseBracketList(agent.CardTypeAuthName)
			directSubAgents = append(directSubAgents, agent)
		}
	}

	return directSubAgents, nil
}

// CreateSubAgent 创建子代理
func (s *AgentService) CreateSubAgent(software, parentUsername string, newAgent *models.Agent) error {
	// 获取软件位数据库连接
	softwareDB, err := s.dbManager.GetSoftwareDB(software)
	if err != nil {
		return fmt.Errorf("连接软件位数据库失败: %v", err)
	}

	// 1. 检查用户名是否已存在
	var count int64
	err = softwareDB.Model(&models.Agent{}).Where("User = ?", newAgent.User).Count(&count).Error
	if err != nil {
		return fmt.Errorf("检查用户名失败: %v", err)
	}
	if count > 0 {
		return fmt.Errorf("用户名已存在")
	}

	// 2. 获取父代理信息
	var parentAgent models.Agent
	err = softwareDB.Where("User = ?", parentUsername).First(&parentAgent).Error
	if err != nil {
		return fmt.Errorf("获取父代理信息失败: %v", err)
	}

	// 3. 检查父代理权限（Service层安全检查）
	authority, err := parentAgent.GetAuthorityUint64()
	if err != nil {
		return fmt.Errorf("解析父代理权限失败: %v", err)
	}
	if (authority & util.PermManageAgent) == 0 {
		return fmt.Errorf("无权创建子代理")
	}

	// 4. 生成FNode
	newAgent.FNode = util.GenerateChildFNode(parentAgent.FNode, newAgent.User)

	// 5. 计算总利率
	// 子代理的总利率 = 父代理的总利率 * (子代理利率 / 100)
	newAgent.TatalParities = parentAgent.TatalParities * (newAgent.Parities / 100.0)

	// 6. 创建代理
	err = softwareDB.Create(newAgent).Error
	if err != nil {
		return fmt.Errorf("创建代理失败: %v", err)
	}

	return nil
}

// GetDirectSubAgentsWithSearch 获取直接子代理列表（支持搜索）
func (s *AgentService) GetDirectSubAgentsWithSearch(software, parentAgent string, searchType int, keyword string) ([]*models.Agent, error) {
	// 获取直接子代理列表
	agents, err := s.GetDirectSubAgents(software, parentAgent)
	if err != nil {
		return nil, err
	}

	// 如果没有关键词，返回所有结果
	if keyword == "" {
		return agents, nil
	}

	// 根据搜索类型和关键词过滤结果
	var filteredAgents []*models.Agent

	for _, agent := range agents {
		if searchType == 0 { // 精准搜索
			if agent.User == keyword {
				filteredAgents = append(filteredAgents, agent)
			}
		} else { // 模糊搜索
			if strings.Contains(agent.User, keyword) {
				filteredAgents = append(filteredAgents, agent)
			}
		}
	}

	return filteredAgents, nil
}

// GetSubAgentsWithSearch 获取所有子代理列表（支持搜索）
func (s *AgentService) GetSubAgentsWithSearch(software, parentAgent string, searchType int, keyword string) ([]*models.Agent, error) {
	// 获取所有子代理列表
	agents, err := s.GetSubAgents(software, parentAgent)
	if err != nil {
		return nil, err
	}

	// 如果没有关键词，返回所有结果
	if keyword == "" {
		return agents, nil
	}

	// 根据搜索类型和关键词过滤结果
	var filteredAgents []*models.Agent

	for _, agent := range agents {
		if searchType == 0 { // 精准搜索
			if agent.User == keyword {
				filteredAgents = append(filteredAgents, agent)
			}
		} else { // 模糊搜索
			if strings.Contains(agent.User, keyword) {
				filteredAgents = append(filteredAgents, agent)
			}
		}
	}

	return filteredAgents, nil
}

// DisableAgent 禁用代理（支持批量）
func (s *AgentService) DisableAgent(software string, usernames []string) (map[string]bool, error) {
	// 获取软件位数据库连接
	softwareDB, err := s.dbManager.GetSoftwareDB(software)
	if err != nil {
		return nil, fmt.Errorf("连接软件位数据库失败: %v", err)
	}

	// 存储每个代理的操作结果
	results := make(map[string]bool)
	for _, username := range usernames {
		// 更新代理状态
		result := softwareDB.Model(&models.Agent{}).
			Where("User = ?", username).
			Updates(map[string]interface{}{
				"Stat":        1, // 1表示禁用
				"CardsEnable": 0, // 0表示禁用代理生成的卡密
			})

		if result.Error != nil {
			results[username] = false
			continue
		}

		results[username] = result.RowsAffected > 0
	}

	return results, nil
}

// EnableAgent 启用代理（支持批量）
func (s *AgentService) EnableAgent(software string, usernames []string) (map[string]bool, error) {
	// 获取软件位数据库连接
	softwareDB, err := s.dbManager.GetSoftwareDB(software)
	if err != nil {
		return nil, fmt.Errorf("连接软件位数据库失败: %v", err)
	}

	// 存储每个代理的操作结果
	results := make(map[string]bool)
	for _, username := range usernames {
		// 更新代理状态
		result := softwareDB.Model(&models.Agent{}).
			Where("User = ?", username).
			Updates(map[string]interface{}{
				"Stat":        0, // 0表示启用
				"CardsEnable": 1, // 1表示子代理卡密启用
			})

		if result.Error != nil {
			results[username] = false
			continue
		}

		results[username] = result.RowsAffected > 0
	}

	return results, nil
}

// UpdateAgentRemark 修改代理备注
func (s *AgentService) UpdateAgentRemark(software, username, remark string) error {
	// 获取软件位数据库连接
	softwareDB, err := s.dbManager.GetSoftwareDB(software)
	if err != nil {
		return fmt.Errorf("连接软件位数据库失败: %v", err)
	}

	// 更新代理备注
	result := softwareDB.Model(&models.Agent{}).
		Where("User = ?", username).
		Update("Remarks", remark)

	if result.Error != nil {
		return fmt.Errorf("更新备注失败: %v", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("代理不存在或无权限修改")
	}

	return nil
}

// DeleteSubAgent 删除子代理
func (s *AgentService) DeleteSubAgent(software, parentAgentName, subAgentName string) error {
	// 获取软件位数据库连接
	softwareDB, err := s.dbManager.GetSoftwareDB(software)
	if err != nil {
		return fmt.Errorf("连接软件位数据库失败: %v", err)
	}

	// 1. 检查父代理权限
	var parentAgent models.Agent
	err = softwareDB.Where("User = ?", parentAgentName).First(&parentAgent).Error
	if err != nil {
		return fmt.Errorf("获取父代理信息失败: %v", err)
	}

	authority, err := parentAgent.GetAuthorityUint64()
	if err != nil {
		return fmt.Errorf("解析父代理权限失败: %v", err)
	}
	if (authority & util.PermManageAgent) == 0 {
		return fmt.Errorf("无权删除子代理")
	}

	// 2. 检查子代理是否存在
	var subAgent models.Agent
	err = softwareDB.Where("User = ? AND Deltm = 0", subAgentName).First(&subAgent).Error
	if err != nil {
		return fmt.Errorf("子代理不存在或已被删除: %v", err)
	}

	// 3. 检查是否为该父代理的直接子代理
	// 通过FNode字段验证层级关系
	if !util.IsDirectChildAgent(subAgentName, parentAgentName, subAgent.FNode) {
		return fmt.Errorf("只能删除直接下级代理")
	}

	// 4. 检查该子代理下是否还有子代理
	var grandChildCount int64
	err = softwareDB.Model(&models.Agent{}).
		Where("FNode LIKE ? AND Deltm = 0", "%\""+subAgentName+"\"%").
		Count(&grandChildCount).Error
	if err != nil {
		return fmt.Errorf("检查下级代理失败: %v", err)
	}

	if grandChildCount > 0 {
		return fmt.Errorf("该代理下还有子代理，无法删除")
	}

	// 5. 软删除：设置Deltm标志为1
	err = softwareDB.Model(&models.Agent{}).
		Where("User = ?", subAgentName).
		Update("Deltm", 1).Error
	if err != nil {
		return fmt.Errorf("删除子代理失败: %v", err)
	}

	return nil
}

// AddMoney 给子代理充值
func (s *AgentService) AddMoney(software, parentAgentName, targetAgentName string, amount float64, timeHours int) error {
	// 获取软件位数据库连接
	softwareDB, err := s.dbManager.GetSoftwareDB(software)
	if err != nil {
		return fmt.Errorf("连接软件位数据库失败: %v", err)
	}

	// 1. 获取父代理信息
	var parentAgent models.Agent
	err = softwareDB.Where("User = ? AND Deltm = 0", parentAgentName).First(&parentAgent).Error
	if err != nil {
		return fmt.Errorf("获取父代理信息失败: %v", err)
	}

	// 2. 检查父代理权限
	authority, err := parentAgent.GetAuthorityUint64()
	if err != nil {
		return fmt.Errorf("解析父代理权限失败: %v", err)
	}
	if (authority & util.PermManageAgent) == 0 {
		return fmt.Errorf("无权给子代理充值")
	}

	// 3. 获取目标代理信息
	var targetAgent models.Agent
	err = softwareDB.Where("User = ? AND Deltm = 0", targetAgentName).First(&targetAgent).Error
	if err != nil {
		return fmt.Errorf("目标代理不存在或已被删除: %v", err)
	}

	// 4. 检查是否为下级代理（包括间接下级）
	if !util.IsChildAgent(targetAgentName, parentAgentName, targetAgent.FNode) {
		return fmt.Errorf("只能给下级代理充值")
	}

	// 5. 检查父代理余额和时长是否足够
	if amount > 0 && parentAgent.AccountBalance < amount {
		return fmt.Errorf("余额不足，当前余额: %.2f", parentAgent.AccountBalance)
	}

	timeSeconds := timeHours * 3600
	if timeSeconds > 0 && parentAgent.AccountTime < timeSeconds {
		return fmt.Errorf("时长不足，当前时长: %d 秒", parentAgent.AccountTime)
	}

	// 6. 计算子代理实际获得的金额（根据目标代理的充值利率）
	// 使用目标代理的parities字段作为充值利率
	actualAmount := amount * (targetAgent.Parities / 100.0)

	// 7. 开始事务操作
	tx := softwareDB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 8. 扣除父代理的余额和时长
	if amount > 0 {
		err = tx.Model(&models.Agent{}).
			Where("User = ?", parentAgentName).
			Update("AccountBalance", parentAgent.AccountBalance-amount).Error
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("扣除父代理余额失败: %v", err)
		}
	}

	if timeSeconds > 0 {
		err = tx.Model(&models.Agent{}).
			Where("User = ?", parentAgentName).
			Update("AccountTime", parentAgent.AccountTime-timeSeconds).Error
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("扣除父代理时长失败: %v", err)
		}
	}

	// 9. 增加子代理的余额和时长
	if actualAmount > 0 {
		err = tx.Model(&models.Agent{}).
			Where("User = ?", targetAgentName).
			Update("AccountBalance", targetAgent.AccountBalance+actualAmount).Error
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("增加子代理余额失败: %v", err)
		}
	}

	if timeSeconds > 0 {
		err = tx.Model(&models.Agent{}).
			Where("User = ?", targetAgentName).
			Update("AccountTime", targetAgent.AccountTime+timeSeconds).Error
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("增加子代理时长失败: %v", err)
		}
	}

	// 10. 提交事务
	err = tx.Commit().Error
	if err != nil {
		return fmt.Errorf("提交事务失败: %v", err)
	}

	return nil
}

// GetAllCardTypes 获取所有卡类型（复用CardTypeService的逻辑）
func (s *AgentService) GetAllCardTypes(software string) ([]models.CardType, error) {
	// 获取软件位数据库连接
	softwareDB, err := s.dbManager.GetSoftwareDB(software)
	if err != nil {
		return nil, fmt.Errorf("连接软件位数据库失败: %v", err)
	}

	var cardTypes []models.CardType
	err = softwareDB.Find(&cardTypes).Error
	if err != nil {
		return nil, fmt.Errorf("查询卡类型失败: %v", err)
	}

	return cardTypes, nil
}

// GetAgentAuthorizedCardTypes 获取代理已授权的卡类型
// 当parentAgentName为空时，直接获取targetAgentName的权限（用于获取当前用户自己的权限）
// 当parentAgentName不为空时，返回子代理权限与父代理权限的交集（父代理视角）
func (s *AgentService) GetAgentAuthorizedCardTypes(software, parentAgentName, targetAgentName string) ([]string, error) {
	// 获取软件位数据库连接
	softwareDB, err := s.dbManager.GetSoftwareDB(software)
	if err != nil {
		return nil, fmt.Errorf("连接软件位数据库失败: %v", err)
	}

	// 获取目标代理信息
	var targetAgent models.Agent
	err = softwareDB.Where("User = ? AND Deltm = 0", targetAgentName).First(&targetAgent).Error
	if err != nil {
		return nil, fmt.Errorf("目标代理不存在: %v", err)
	}

	// 解析目标代理的CardTypeAuthName字段
	targetAuthorizedCardTypes := util.ParseBracketList(targetAgent.CardTypeAuthName)

	// 如果没有parentAgentName，直接返回目标代理的权限（用于获取自己的权限）
	if parentAgentName == "" {
		return targetAuthorizedCardTypes, nil
	}

	// 如果有parentAgentName，需要验证权限关系并返回交集
	// 检查是否为下级代理
	if !util.IsChildAgent(targetAgentName, parentAgentName, targetAgent.FNode) {
		return nil, fmt.Errorf("只能查询下级代理权限")
	}

	// 获取父代理的权限
	var parentAgent models.Agent
	err = softwareDB.Where("User = ? AND Deltm = 0", parentAgentName).First(&parentAgent).Error
	if err != nil {
		return nil, fmt.Errorf("父代理不存在: %v", err)
	}

	parentAuthorizedCardTypes := util.ParseBracketList(parentAgent.CardTypeAuthName)

	// 返回子代理权限与父代理权限的交集（父代理视角下的子代理权限）
	intersection := make([]string, 0)
	for _, childCardType := range targetAuthorizedCardTypes {
		for _, parentCardType := range parentAuthorizedCardTypes {
			if childCardType == parentCardType {
				intersection = append(intersection, childCardType)
				break
			}
		}
	}

	return intersection, nil
}

// UpdateAgentCardTypePermissions 更新代理卡类型权限
func (s *AgentService) UpdateAgentCardTypePermissions(software, parentAgentName, targetAgentName string, cardTypeNames []string) error {
	// 获取软件位数据库连接
	softwareDB, err := s.dbManager.GetSoftwareDB(software)
	if err != nil {
		return fmt.Errorf("连接软件位数据库失败: %v", err)
	}

	// 1. 检查父代理权限
	var parentAgent models.Agent
	err = softwareDB.Where("User = ? AND Deltm = 0", parentAgentName).First(&parentAgent).Error
	if err != nil {
		return fmt.Errorf("获取父代理信息失败: %v", err)
	}

	parentAuthority, err := parentAgent.GetAuthorityUint64()
	if err != nil {
		return fmt.Errorf("解析父代理权限失败: %v", err)
	}
	if (parentAuthority & util.PermManageAgent) == 0 {
		return fmt.Errorf("无权修改代理权限")
	}

	// 2. 检查目标代理是否存在
	var targetAgent models.Agent
	err = softwareDB.Where("User = ? AND Deltm = 0", targetAgentName).First(&targetAgent).Error
	if err != nil {
		return fmt.Errorf("目标代理不存在或已被删除: %v", err)
	}

	// 3. 检查是否为下级代理
	if !util.IsChildAgent(targetAgentName, parentAgentName, targetAgent.FNode) {
		return fmt.Errorf("只能修改下级代理权限")
	}

	// 4. 构建新的CardTypeAuthName
	newCardTypeAuthName := util.BuildBracketList(cardTypeNames)

	// 5. 更新代理的卡类型权限
	err = softwareDB.Model(&models.Agent{}).
		Where("User = ?", targetAgentName).
		Update("CardTypeAuthName", newCardTypeAuthName).Error
	if err != nil {
		return fmt.Errorf("更新代理卡类型权限失败: %v", err)
	}

	return nil
}
