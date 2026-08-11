# 国美考务系统

当前根目录已替换为 `国美考务-源码+演示库` 中提供的正式功能源码。

## 技术栈

- 前端：Vue 3.4、TypeScript、Vite 5、Element Plus、Pinia
- 后端：Go、Gin、GORM
- 数据库：MySQL 8.0.41

## 快速启动

~~~powershell
cd D:\2026\codex\YWJLDP
.\scripts\start-local.ps1
~~~

打开：`http://127.0.0.1:8848`

登录：

- 用户名：`admin`
- 密码：`Admin123`
- 验证码：填写登录页显示的算式结果

停止服务：

~~~powershell
.\scripts\stop-local.ps1
~~~

详细配置、端口与数据量见 [本地运行说明](./docs/本地运行说明.md)。

## 目录

~~~text
YWJLDP/
├─ frontend/                  正式 Vue 3 前端
├─ backend/                   正式 Go 后端
├─ scripts/                   本地启动、停止和数据库重置脚本
├─ docs/                      本地运行文档
├─ .runtime/                  MySQL、构建缓存和可执行程序（不提交 Git）
└─ 国美考务-源码+演示库/       原始交付源码与 SQL 备份
~~~

