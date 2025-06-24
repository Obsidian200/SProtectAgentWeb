package services

import (
	"SProtectAgentWeb/database"
	"SProtectAgentWeb/models"
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"
)

type CardService struct {
	dbManager *database.DatabaseManager
}

func NewCardService(dbManager *database.DatabaseManager) *CardService {
	return &CardService{dbManager: dbManager}
}

// DisableCard 禁用卡密
func (s *CardService) DisableCard(software string, cardKey string) (bool, error) {
	db, err := s.dbManager.GetSoftwareDB(software)
	if err != nil {
		return false, fmt.Errorf("获取数据库连接失败: %v", err)
	}

	currentTime := time.Now().Unix()

	// 更新卡密状态
	result := db.Model(&models.CardInfo{}).
		Where("Prefix_Name = ?", cardKey).
		Updates(map[string]interface{}{
			"state":   "禁用", // 使用常量替代硬编码
			"BanTime": currentTime,
		})

	if result.Error != nil {
		return false, fmt.Errorf("禁用卡密失败: %v", result.Error)
	}

	if result.RowsAffected == 0 {
		return false, fmt.Errorf("卡密不存在或已经是禁用状态")
	}

	// 记录操作日志
	// 如果有现成的日志函数，可以调用
	// s.auditService.LogOperation(software, operator, "禁用卡密", cardKey)

	return true, nil
}

// EnableCard 启用卡密
func (s *CardService) EnableCard(software string, cardKey string) (bool, error) {
	db, err := s.dbManager.GetSoftwareDB(software)
	if err != nil {
		return false, fmt.Errorf("获取数据库连接失败: %v", err)
	}

	// 更新卡密状态
	result := db.Model(&models.CardInfo{}).
		Where("Prefix_Name = ?", cardKey).
		Updates(map[string]interface{}{
			"state":   "启用", // 使用常量替代硬编码
			"BanTime": 0,
		})

	if result.Error != nil {
		return false, fmt.Errorf("启用卡密失败: %v", result.Error)
	}

	if result.RowsAffected == 0 {
		return false, fmt.Errorf("卡密不存在或已经是启用状态")
	}

	// 记录操作日志
	// 如果有现成的日志函数，可以调用
	// s.auditService.LogOperation(software, operator, "启用卡密", cardKey)

	return true, nil
}

// EnableCardWithBanTimeReturn 启用卡密并归还封禁时间
func (s *CardService) EnableCardWithBanTimeReturn(software string, cardKey string) (bool, error) {
	db, err := s.dbManager.GetSoftwareDB(software)
	if err != nil {
		return false, fmt.Errorf("获取数据库连接失败: %v", err)
	}

	// 首先获取卡密信息，用于计算封禁时间
	var card models.CardInfo
	if err := db.Where("Prefix_Name = ?", cardKey).First(&card).Error; err != nil {
		return false, fmt.Errorf("卡密不存在: %v", err)
	}

	// 如果卡密不是禁用状态或没有封禁时间，则无需归还
	if card.State != "禁用" || card.BanTime == 0 {
		return false, fmt.Errorf("卡密不是禁用状态或没有封禁时间")
	}

	// 计算封禁时长（秒）
	currentTime := time.Now().Unix()
	banDuration := currentTime - int64(card.BanTime)

	// 如果封禁时间为负数（可能是系统时间被调整），则设为0
	if banDuration < 0 {
		banDuration = 0
	}

	// 使用事务确保操作的原子性
	err = db.Transaction(func(tx *gorm.DB) error {
		// 1. 更新卡密状态为启用，清除封禁时间
		if err := tx.Model(&models.CardInfo{}).
			Where("Prefix_Name = ?", cardKey).
			Updates(map[string]interface{}{
				"State":   "启用",
				"BanTime": 0,
			}).Error; err != nil {
			return err
		}

		// 2. 如果有到期时间，延长到期时间
		if card.ExpiredTime__ > 0 {
			newExpireTime := card.ExpiredTime__ + banDuration
			if err := tx.Model(&models.CardInfo{}).
				Where("Prefix_Name = ?", cardKey).
				Update("ExpiredTime__", newExpireTime).Error; err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return false, fmt.Errorf("启用卡密并归还封禁时间失败: %v", err)
	}

	// 记录操作日志
	// 如果有现成的日志函数，可以调用
	// s.auditService.LogOperation(software, operator, "启用卡密(归还封禁时间)", cardKey)

	return true, nil
}

// CardQueryParams 卡密查询参数
type CardQueryParams struct {
	Software     string   // 软件名称
	Status       string   // 状态筛选：0-全部，1-启用，2-禁用
	SearchType   int      // 搜索类型：0-精准搜索，1-模糊搜索
	Keywords     []string // 搜索关键词数组，支持多个卡密同时搜索
	Page         int      // 页码
	PageSize     int      // 每页数量
	CurrentAgent string   // 当前代理用户名
}

// GetCardList 获取卡密列表
func (s *CardService) GetCardList(params *CardQueryParams) ([]models.CardInfo, int64, error) {
	// 获取软件位数据库连接
	db, err := s.dbManager.GetSoftwareDB(params.Software)
	if err != nil {
		return nil, 0, fmt.Errorf("获取数据库连接失败: %v", err)
	}

	// 构建基础查询
	query := db.Table("CardInfo")

	// 只查询当前代理的卡密
	query = query.Where("Whom = ?", params.CurrentAgent)

	// 过滤掉已删除的卡密
	query = query.Where("delstate = 0")

	// 添加状态筛选
	if params.Status == "1" {
		query = query.Where("state = ?", "启用")
	} else if params.Status == "2" {
		query = query.Where("state = ?", "禁用")
	}

	// 添加关键词搜索
	if len(params.Keywords) > 0 {
		if params.SearchType == 0 {
			// 精准搜索 - 使用 IN 查询
			query = query.Where("Prefix_Name IN ?", params.Keywords)
		} else {
			// 模糊搜索 - 构建 OR 条件
			orConditions := make([]string, 0, len(params.Keywords))
			orArgs := make([]interface{}, 0, len(params.Keywords))

			for _, keyword := range params.Keywords {
				orConditions = append(orConditions, "Prefix_Name LIKE ?")
				orArgs = append(orArgs, "%"+keyword+"%")
			}

			// 将所有条件用 OR 连接
			query = query.Where(
				"("+strings.Join(orConditions, " OR ")+")",
				orArgs...,
			)
		}
	}

	// 获取总数
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("获取卡密总数失败: %v", err)
	}

	// 分页查询
	offset := (params.Page - 1) * params.PageSize
	query = query.Offset(offset).Limit(params.PageSize)

	// 排序
	query = query.Order("CreateData_ DESC")

	// 执行查询
	var results []models.CardInfo
	if err := query.Find(&results).Error; err != nil {
		return nil, 0, fmt.Errorf("查询卡密列表失败: %v", err)
	}

	return results, total, nil
}
