# 📖 StrShelf

一个简单优雅的书签管理应用，帮助你组织和分享你的链接集合。

## ✨ 功能特性

- 🔐 **用户认证** - JWT 令牌认证，支持用户登录和注册
- 📝 **链接管理** - 创建、编辑、删除书签
- 🔍 **搜索功能(TODO)** - 快速搜索你的书签
- 📅 **时间分组** - 按创建日期自动组织书签
- 🎨 **现代化界面** - Vue 3 构建的响应式前端
- 💾 **数据持久化** - PostgreSQL 数据库支持
- 🚀 **部署便捷** - Docker 友好的构建配置

## 🛠️ 技术栈

### 前端

- **Vue 3** - 现代化 JavaScript 框架
- **TypeScript** - 类型安全的开发体验
- **Vite** - 快速的前端构建工具
- **CSS3** - 丰富的动画和样式效果

### 后端

- **Go** - 高性能后端语言
- **Gin** - 轻量级 Web 框架
- **GORM** - Go ORM 库
- **PostgreSQL** - 可靠的关系数据库
- **JWT** - 安全的身份认证
- **Zap** - 高性能日志库

## 🚀 快速开始

### 前置要求

- Node.js >= 18
- Go >= 1.21
- PostgreSQL >= 13

### 安装与构建

1. **克隆项目**

```bash
git clone https://github.com/swim233/StrShelf.git
cd StrShelf
```

2. **安装依赖**

```bash
make install
```

3. **初始化数据库**

```bash
# 创建数据库并执行初始化脚本
psql -U postgres -f strshelf.sql
```

4. **配置应用**

```bash
# 复制配置示例
cp packages/api/config/config.json.example packages/api/config/config.json

# 编辑配置文件，设置：
# - secret_key: JWT 密钥
# - dsn: 数据库连接字符串
# - allow_signup: 是否允许注册
```

5. **编译构建**

```bash
make build
```

### 开发模式

**启动后端服务**

```bash
cd packages/api
go run main.go
```

**启动前端开发服务器**

```bash
npm run dev -w strshelf-web
```

访问 `http://localhost:5173` 查看应用。

## 📚 API 文档

### 用户认证

#### 用户登录

```
POST /v1/user.login
Content-Type: application/json

{
  "account": "username",
  "password": "password"
}

响应:
{
  "token": "jwt_token_string"
}
```

#### 验证令牌

```
POST /v1/user.verify
Content-Type: application/json

{
  "token": "jwt_token_string"
}

响应:
{
  "msg": "successful"
}
```

#### 用户注册 (需启用)

```
POST /v1/user.signup
Content-Type: application/json

{
  "account": "newuser",
  "password": "password"
}
```

### 书签管理

#### 获取所有书签

```
POST /v1/item.get

响应:
{
  "code": 200,
  "data": [
    {
      "id": 1,
      "title": "书签标题",
      "link": "https://example.com",
      "comment": "书签描述",
      "gmt_created": 1234567890,
      "gmt_modified": 1234567890,
      "deleted": false
    }
  ],
  "msg": ""
}
```

#### 创建书签

```
POST /v1/item.post
Authorization: Bearer {token}
Content-Type: application/json

{
  "title": "新书签",
  "link": "https://example.com",
  "comment": "描述"
}
```

#### 编辑书签

```
POST /v1/item.edit
Authorization: Bearer {token}
Content-Type: application/json

{
  "id": 1,
  "new_title": "更新标题",
  "new_link": "https://newurl.com",
  "new_comment": "新描述"
}
```

#### 删除书签

```
POST /v1/item.delete
Authorization: Bearer {token}
Content-Type: application/json

{
  "id": 1
}
```

## ⚙️ 配置说明

编辑 `packages/api/config/config.json`:

```json
{
  "port": "1111",
  "secret_key": "你的JWT密钥",
  "dsn": "host=localhost user=postgres dbname=strshelf port=5432 sslmode=disable TimeZone=Asia/Shanghai",
  "allow_signup": false
}
```

**配置字段说明：**

- `port` - API 服务端口 (默认: 1111)
- `secret_key` - JWT 签名密钥 (必需，修改为强密钥)
- `dsn` - PostgreSQL 数据库连接字符串
- `allow_signup` - 是否允许新用户注册

## 📝 数据库架构

### shelf_item_v1 表

存储书签信息

| 字段         | 类型      | 说明        |
| ------------ | --------- | ----------- |
| id           | bigint    | 主键 (自增) |
| title        | text      | 书签标题    |
| link         | text      | 链接地址    |
| comment      | text      | 书签备注    |
| gmt_created  | timestamp | 创建时间    |
| gmt_modified | timestamp | 修改时间    |
| gmt_deleted  | timestamp | 删除时间    |
| deleted      | boolean   | 软删除标记  |

### shelf_user_v1 表

存储用户信息

| 字段        | 类型      | 说明        |
| ----------- | --------- | ----------- |
| id          | bigint    | 主键 (自增) |
| username    | text      | 用户名      |
| password    | text      | 密码哈希值  |
| gmt_created | timestamp | 创建时间    |

## 🐳 Docker 部署

```bash
# 构建 Docker 镜像
docker build -t strshelf .

# 运行容器
docker run -p 1111:1111 \
  -e DB_DSN="host=postgres user=postgres dbname=strshelf port=5432 sslmode=disable" \
  -e SECRET_KEY="your-secret-key" \
  strshelf
```

## 📄 许可证

MIT License

## 🤝 贡献

欢迎提交 Issue 和 Pull Request！

## 📮 联系方式

如有问题或建议，请通过以下方式联系：

- 提交 GitHub Issue
- 发送邮件至 [swim853279614@163.com](swim853279614@163.com)

---
