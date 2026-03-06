# Postman 测试指南

## 步骤 1: 获取测试 Token

### 请求配置
- **方法**: POST
- **URL**: `http://localhost:8080/api/v1/auth/test-token`
- **Headers**:
  ```
  Content-Type: application/json
  ```
- **Body** (raw JSON):
  ```json
  {
    "user_id": 2
  }
  ```

### 响应示例
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "id": 2,
      "open_id": "test_openid_001",
      "nickname": "张三",
      "avatar_url": "https://thirdwx.qlogo.cn/mmopen/vi_32/test001.png",
      "gender": 1,
      "country": "中国",
      "province": "广东省",
      "city": "深圳市",
      "language": "zh_CN",
      "phone": "13800138001",
      "status": 1
    }
  }
}
```

**复制响应中的 `token` 值，后续请求需要使用。**

---

## 步骤 2: 使用 Token 访问需要认证的接口

### 方法 A: 在 Headers 中手动添加

#### 请求配置
- **方法**: GET
- **URL**: `http://localhost:8080/api/v1/user/info`
- **Headers**:
  ```
  Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
  ```
  
  **注意**: `Bearer` 和 token 之间有一个空格！

### 方法 B: 使用 Postman 的 Authorization 功能（推荐）

1. 在请求页面选择 **Authorization** 标签
2. **Type** 选择 `Bearer Token`
3. 在 **Token** 输入框中粘贴你的 token（不需要加 "Bearer " 前缀）
4. Postman 会自动在请求头中添加 `Authorization: Bearer <token>`

---

## 可用的测试用户 ID

根据之前添加的测试数据，你可以使用以下用户 ID 生成 Token：

| user_id | 昵称 | 城市 | 手机号 |
|---------|------|------|--------|
| 2 | 张三 | 深圳市 | 13800138001 |
| 3 | 李四 | 北京市 | 13800138002 |
| 4 | 王五 | 上海市 | 13800138003 |
| 5 | 赵六 | 杭州市 | 13800138004 |
| 6 | 孙七 | 南京市 | 13800138005 |
| 7 | 周八 | 成都市 | 13800138006 |
| 8 | 吴九 | 武汉市 | 13800138007 |
| 9 | 郑十 | 西安市 | 13800138008 |
| 10 | 钱十一 | 长沙市 | 13800138009 |
| 11 | 冯十二 | 厦门市 | 13800138010 |

---

## 完整测试流程示例

### 1. 生成 Token
```bash
POST http://localhost:8080/api/v1/auth/test-token
Content-Type: application/json

{
  "user_id": 2
}
```

### 2. 获取用户信息
```bash
GET http://localhost:8080/api/v1/user/info
Authorization: Bearer <your_token>
```

### 3. 更新用户信息
```bash
PUT http://localhost:8080/api/v1/user/info
Authorization: Bearer <your_token>
Content-Type: application/json

{
  "nickname": "新昵称",
  "phone": "13900139000"
}
```

### 4. 获取用户当前座位
```bash
GET http://localhost:8080/api/v1/seats/my-seat
Authorization: Bearer <your_token>
```

### 5. 用户入座
```bash
POST http://localhost:8080/api/v1/seats/join
Authorization: Bearer <your_token>
Content-Type: application/json

{
  "seat_id": "C09+",
  "seat_type": "card"
}
```

### 6. 用户离座
```bash
POST http://localhost:8080/api/v1/seats/leave
Authorization: Bearer <your_token>
```

### 7. 查询座位占用信息（无需认证）
```bash
GET http://localhost:8080/api/v1/seats/occupancy?seat_id=C09+
```

---

## 常见错误

### 错误 1: "未提供认证令牌"
**原因**: 没有在请求头中添加 Authorization  
**解决**: 确保添加了 `Authorization: Bearer <token>` 请求头

### 错误 2: "认证令牌格式错误"
**原因**: Authorization 格式不正确  
**解决**: 确保格式为 `Bearer <token>`，注意 Bearer 和 token 之间有空格

### 错误 3: "认证令牌无效或已过期"
**原因**: Token 已过期（默认 7 天）或无效  
**解决**: 重新调用 `/api/v1/auth/test-token` 生成新的 token

---

## Postman Collection 导入

你可以创建一个 Postman Collection 来保存这些请求：

1. 点击 Postman 左侧的 **Collections**
2. 点击 **+** 创建新 Collection
3. 命名为 "MagicKingdom API"
4. 在 Collection 的 **Variables** 标签中添加变量：
   - Variable: `base_url`, Initial Value: `http://localhost:8080`
   - Variable: `token`, Initial Value: (留空，测试时填入)
5. 在 Collection 的 **Authorization** 标签中：
   - Type: `Bearer Token`
   - Token: `{{token}}`
6. 添加请求时使用 `{{base_url}}` 代替完整 URL

这样所有继承 Collection 认证的请求都会自动使用 token 变量。

---

## 注意事项

⚠️ **重要**: `/api/v1/auth/test-token` 接口仅用于开发测试，生产环境必须删除此接口！

生产环境应该只使用微信登录接口 `/api/v1/auth/wechat/login` 来获取 Token。

