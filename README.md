# SProtectAgentWeb

SProtectAgentWeb 是为 SProtect 网络验证系统开发的 Web 代理端，采用 Go + Layui 技术栈，提供简洁易用的代理管理界面。

SProtectAgentWeb不依赖SProtect云计算和服务端，程序独立运行，功能开发和逻辑设计都是通过在PC代理端反复操作和通过对数据库变动对比分析得出，可能存在逻辑遗漏以及功能和官方代理端效果不一致的问题，可向我反馈修复

## 使用说明
在[Realease](https://github.com/Obsidian200/SProtectAgentWeb/releases)页面下载最新版本的zip解压出来

SProtectAgentWeb.exe和SProtectAgentWeb.ini放到服务端目录，运行即可

SProtectAgentWeb.ini是配置文件，修改端口号和数据库路径，默认路径是当前目录

## 主要功能

### 子代理管理
- 创建和管理下级代理
- 代理权限分配(开发中)
- 余额和时长充值
- 代理状态控制

### 卡密管理
- 批量生成卡密
- 卡密启用/禁用
- 卡密状态查询
- 支持自定义前缀

### 卡类型管理
- 卡类型配置
- 价格和时长设置
- 权限分配管理
<!-- 
## 快速开始

### 安装运行

```bash
# 下载依赖
go mod download

# 编译运行
go build -o SProtectAgentWeb
./SProtectAgentWeb
```

### 配置文件

编辑 `config/SProtectAgentWeb.ini`： -->


## 技术栈

- **后端**: Golang
- **前端**: Layui
