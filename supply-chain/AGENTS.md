# AGENTS.md - 供应链系统开发指南

## 项目概述

供应链管理系统，包含供应商管理、采购管理、库存管理、销售订单、物流跟踪等模块。

**项目结构：**
- `supply-chain/` - 前端项目 (Vue 3)
- `supply-chain-server/` - 后端项目 (Go)

---

## 前端 (Vue 3)

### 构建命令
```bash
cd supply-chain
npm run dev          # 启动开发服务器 (端口 3000)
npm run build        # TypeScript 检查 + 构建生产版本
npm run preview      # 预览生产构建
npx vue-tsc --noEmit # 仅运行 TypeScript 类型检查
```

### 技术栈
- Vue 3.4 + TypeScript 5.3 + Vite 5
- Pinia 状态管理 + Vue Router 4
- PrimeVue 4.x + SCSS
- Chart.js + ECharts
- Vue I18n (中英文)

### 文件命名
- 页面组件: PascalCase + Page 后缀 (SupplierPage.vue)
- Store: camelCase.ts (supply.ts)
- 类型定义: camelCase.ts (index.ts)

### 路径导入
使用 @ 别名导入模块:
```ts
import { useSupplyStore } from '@/stores/supply'
import type { Supplier } from '@/types'
```

### Vue 组件结构
```vue
<template>
  <!-- 模板内容 -->
</template>

<script setup lang="ts">
// 1. 导入
import { ref, computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'

// 2. 类型定义
interface Form { name: string }

// 3. Composables
const { t } = useI18n()

// 4. 响应式状态
const loading = ref(false)

// 5. 计算属性
const filteredList = computed(() => [])

// 6. 方法
function handleSubmit() {}

// 7. 生命周期
onMounted(() => {})
</script>

<style lang="scss" scoped>
</style>
```

### TypeScript 规范
- 所有变量和函数必须有明确类型
- 接口使用 PascalCase 命名
- 使用联合类型定义状态字面量
- 避免使用 any，导入类型使用 `import type`

---

## 后端 (Go + Gin)

### 构建命令
```bash
cd supply-chain-server
go mod tidy                    # 安装依赖
go build -o bin/server ./cmd/server/main.go  # 编译
./bin/server                   # 运行服务 (端口 8080)
```

### 测试命令
```bash
go test ./...                  # 运行所有测试
go test ./internal/service -run TestName  # 运行单个测试
go test -v ./...               # 详细输出
```

### 技术栈
- Go 1.24 + Gin Web 框架
- MySQL + GORM ORM
- JWT 认证 (golang-jwt/jwt)
- Viper 配置管理

### 项目结构
```
supply-chain-server/
├── cmd/server/main.go       # 程序入口
├── configs/config.yaml      # 配置文件
├── internal/
│   ├── config/              # 配置管理
│   ├── model/               # 数据模型
│   ├── repository/          # 数据访问层
│   ├── service/             # 业务逻辑层
│   ├── handler/             # HTTP 处理器
│   ├── middleware/          # 中间件 (JWT, CORS)
│   └── router/              # 路由配置
├── pkg/
│   ├── database/            # 数据库连接
│   ├── response/            # 统一响应格式
│   └── utils/               # 工具函数
└── docs/sql/schema.sql      # 数据库脚本
```

### Go 代码规范

#### 命名约定
- 包名: 小写单词 (repository, service, handler)
- 类型: PascalCase (UserService, SupplierRepository)
- 方法: PascalCase 导出, camelCase 私有
- 常量: UPPER_SNAKE_CASE 或 PascalCase

#### 错误处理
```go
// 返回错误
if err != nil {
    return fmt.Errorf("操作失败: %w", err)
}

// 处理错误
if err := service.Create(user); err != nil {
    response.ServerError(c, err.Error())
    return
}
```

#### Handler 示例
```go
func (h *SupplierHandler) List(c *gin.Context) {
    page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
    pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
    
    items, total, err := h.service.List(page, pageSize, "")
    if err != nil {
        response.ServerError(c, err.Error())
        return
    }
    
    response.PageSuccess(c, items, total, page, pageSize)
}
```

#### Repository 示例
```go
func (r *SupplierRepository) List(page, pageSize int, keyword string) ([]model.Supplier, int64, error) {
    var suppliers []model.Supplier
    var total int64

    query := r.db.Model(&model.Supplier{})
    if keyword != "" {
        query = query.Where("name LIKE ?", "%"+keyword+"%")
    }

    query.Count(&total)
    err := query.Offset((page - 1) * pageSize).Limit(pageSize).Find(&suppliers).Error
    return suppliers, total, err
}
```

### API 路由

```
POST   /api/auth/login           # 登录
POST   /api/auth/register        # 注册
GET    /api/auth/profile         # 用户信息 (需认证)

GET    /api/suppliers            # 供应商列表
POST   /api/suppliers            # 创建供应商
GET    /api/suppliers/:id        # 供应商详情
PUT    /api/suppliers/:id        # 更新供应商
DELETE /api/suppliers/:id        # 删除供应商

GET    /api/inventory            # 库存列表
POST   /api/inventory/stock-in   # 入库
POST   /api/inventory/stock-out  # 出库
GET    /api/inventory/stats      # 库存统计

GET    /api/dashboard/stats      # 仪表盘统计
```

### 数据库配置
修改 `configs/config.yaml`:
```yaml
database:
  host: "localhost"
  port: 3306
  user: "root"
  password: "your_password"
  dbname: "supply_chain"
```

### 初始化数据库
```bash
mysql -u root -p < docs/sql/schema.sql
```

---

## 开发流程

1. **启动后端**: `cd supply-chain-server && go run ./cmd/server/main.go`
2. **启动前端**: `cd supply-chain && npm run dev`
3. **访问应用**: http://localhost:3000
4. **API 文档**: http://localhost:8080/api