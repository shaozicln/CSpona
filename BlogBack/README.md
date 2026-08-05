# BlogBack

Gin + Gorm + MySQL 后端。

## 目录

| 路径 | 说明 |
|------|------|
| `api/` | 接口与业务 |
| `routes/` | 路由 |
| `middleware/` | CORS 等中间件 |
| `utils/` | 配置、AES 等工具 |
| `database/` | SQL 备份 / schema |
| `config.example.ini` | 配置模板 |
| `config.ini` | 本地真实配置（gitignore，需自行复制模板） |

## 启动

```bash
cp config.example.ini config.ini   # 首次
# 编辑 config.ini
go run .
```

详见仓库根目录 `README.md`。
