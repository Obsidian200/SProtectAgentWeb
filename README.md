# 🛡️ SProtectAgentWeb

---

## 📖 简介

SProtectAgentWeb 是为 SProtect 网络验证系统开发的 Web 代理端，采用 Go + Layui 技术栈，提供简洁易用的代理管理界面。

## ✨ 主要功能

### 👥 子代理管理
- 创建和管理下级代理
- 代理权限分配(开发中)
- 余额和时长充值
- 代理状态控制

### 💳 卡密管理
- 批量生成卡密
- 卡密启用/禁用
- 卡密状态查询
- 支持自定义前缀

### 🏷️ 卡类型管理
- 卡类型配置
- 价格和时长设置
- 权限分配管理

## 🚀 快速开始

### 安装运行

```bash
# 下载依赖
go mod download

# 编译运行
go build -o SProtectAgentWeb
./SProtectAgentWeb
```

### 配置文件

编辑 `config/SProtectAgentWeb.ini`：


## 🛠️ 技术栈

- **后端**: Go + Gin + GORM + SQLite
- **前端**: Layui
