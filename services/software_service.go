package services

import (
	"fmt"
	"math"
	"web-agent-backend/database"
	"web-agent-backend/models"
	"web-agent-backend/types"
	"web-agent-backend/util"
)

// SoftwareService 软件位服务
// 只依赖DatabaseManager，专注软件位业务逻辑
type SoftwareService struct {
	dbManager *database.DatabaseManager
}

// NewSoftwareService 创建软件位服务实例
func NewSoftwareService(dbManager *database.DatabaseManager) *SoftwareService {
	return &SoftwareService{
		dbManager: dbManager,
	}
}

// GetSoftwares 获取所有软件位列表
func (s *SoftwareService) GetSoftwares() (*types.GetSoftwaresResponse, error) {
	// 获取主数据库连接
	mainDB, err := s.dbManager.GetDafaultDB()
	if err != nil {
		return nil, fmt.Errorf("连接主数据库失败: %v", err)
	}

	// 查询所有软件位
	var softwares []*models.MultiSoftware
	err = mainDB.Find(&softwares).Error
	if err != nil {
		return nil, fmt.Errorf("查询软件位列表失败: %v", err)
	}

	return &types.GetSoftwaresResponse{
		Softwares: softwares,
	}, nil
}

// GetEnabledSoftwares 获取启用的软件位列表
func (s *SoftwareService) GetEnabledSoftwares() (*types.GetSoftwaresResponse, error) {
	// 获取主数据库连接
	mainDB, err := s.dbManager.GetDafaultDB()
	if err != nil {
		return nil, fmt.Errorf("连接主数据库失败: %v", err)
	}

	// 查询启用的软件位
	var softwares []*models.MultiSoftware
	err = mainDB.Where("State = ?", 1).Find(&softwares).Error
	if err != nil {
		return nil, fmt.Errorf("查询启用软件位列表失败: %v", err)
	}

	return &types.GetSoftwaresResponse{
		Softwares: softwares,
	}, nil
}

// GetSoftwareList 获取软件位列表（分页）
func (s *SoftwareService) GetSoftwareList(req *types.GetSoftwareListRequest) (*types.GetSoftwareListResponse, error) {
	// 获取主数据库连接
	mainDB, err := s.dbManager.GetDafaultDB()
	if err != nil {
		return nil, fmt.Errorf("连接主数据库失败: %v", err)
	}

	// 构建查询条件
	query := mainDB.Model(&models.MultiSoftware{})

	// 添加筛选条件
	if req.State != nil {
		query = query.Where("State = ?", *req.State)
	}
	if req.Search != "" {
		query = query.Where("SoftwareName LIKE ? OR IDC LIKE ?",
			"%"+req.Search+"%", "%"+req.Search+"%")
	}

	// 获取总数
	var total int64
	err = query.Count(&total).Error
	if err != nil {
		return nil, fmt.Errorf("获取软件位总数失败: %v", err)
	}

	// 设置默认分页参数
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 20
	}

	// 设置排序
	orderBy := "ID DESC" // 默认按ID降序
	if req.SortField != "" {
		order := "ASC"
		if req.SortOrder == "desc" {
			order = "DESC"
		}
		orderBy = fmt.Sprintf("%s %s", req.SortField, order)
	}

	// 查询软件位列表
	var softwares []*models.MultiSoftware
	offset := (req.Page - 1) * req.PageSize
	err = query.Order(orderBy).Offset(offset).Limit(req.PageSize).Find(&softwares).Error
	if err != nil {
		return nil, fmt.Errorf("查询软件位列表失败: %v", err)
	}

	// 计算总页数
	totalPages := int(math.Ceil(float64(total) / float64(req.PageSize)))

	return &types.GetSoftwareListResponse{
		Items: softwares,
		Pagination: &types.Pagination{
			Total:      total,
			Page:       req.Page,
			PageSize:   req.PageSize,
			TotalPages: totalPages,
		},
	}, nil
}

// GetSoftwareInfo 获取单个软件位详细信息
func (s *SoftwareService) GetSoftwareInfo(softwareName string) (*types.GetSoftwareInfoResponse, error) {
	// 获取主数据库连接
	mainDB, err := s.dbManager.GetDafaultDB()
	if err != nil {
		return nil, fmt.Errorf("连接主数据库失败: %v", err)
	}

	// 查询软件位信息
	var software models.MultiSoftware
	err = mainDB.Where("SoftwareName = ?", softwareName).First(&software).Error
	if err != nil {
		return nil, fmt.Errorf("查询软件位信息失败: %v", err)
	}

	// 获取统计信息
	statistics, err := s.getSoftwareStatistics(softwareName)
	if err != nil {
		// 统计信息获取失败不影响主要功能
		statistics = &types.SoftwareStats{}
	}

	return &types.GetSoftwareInfoResponse{
		Software:   &software,
		Statistics: statistics,
	}, nil
}

// ValidateSoftware 验证软件位是否存在且可用
func (s *SoftwareService) ValidateSoftware(softwareName string) error {
	// 获取主数据库连接
	mainDB, err := s.dbManager.GetDafaultDB()
	if err != nil {
		return fmt.Errorf("连接主数据库失败: %v", err)
	}

	// 查询软件位
	var software models.MultiSoftware
	err = mainDB.Where("SoftwareName = ? AND State = ?", softwareName, 1).First(&software).Error
	if err != nil {
		return fmt.Errorf("软件位不存在或已禁用: %s", softwareName)
	}

	// 测试软件位数据库连接
	_, err = s.dbManager.GetSoftwareDB(softwareName)
	if err != nil {
		return fmt.Errorf("无法连接软件位数据库: %v", err)
	}

	return nil
}

// ===== 私有方法 =====

// getSoftwareStatistics 获取软件位统计信息
func (s *SoftwareService) getSoftwareStatistics(softwareName string) (*types.SoftwareStats, error) {
	// 尝试连接软件位数据库
	softwareDB, err := s.dbManager.GetSoftwareDB(softwareName)
	if err != nil {
		return nil, fmt.Errorf("连接软件位数据库失败: %v", err)
	}

	stats := &types.SoftwareStats{}

	// 统计代理数量
	softwareDB.Model(&models.Agent{}).Count(&stats.TotalAgents)
	softwareDB.Model(&models.Agent{}).Where("Stat = 1 AND Deltm = 0").Count(&stats.ActiveAgents)

	// 这里可以根据需要添加更多统计逻辑
	// 目前只实现基础的代理统计

	return stats, nil
}

// GetSoftwareAgentByName 根据软件位和代理名称获取代理信息
func (s *SoftwareService) GetSoftwareAgentByName(software, agentName string) (*models.Agent, error) {
	// 获取软件位数据库连接
	db, err := s.dbManager.GetSoftwareDB(software)
	if err != nil {
		return nil, fmt.Errorf("获取数据库连接失败: %v", err)
	}

	// 查询代理信息
	var agent models.Agent
	err = db.Where("User = ?", agentName).First(&agent).Error
	if err != nil {
		return nil, fmt.Errorf("查询代理信息失败: %v", err)
	}

	// 处理CardTypeAuthName，解析为数组格式
	agent.CardTypeAuthNameArray = util.ParseBracketList(agent.CardTypeAuthName)

	return &agent, nil
}
