# API 测试示例

## 使用 curl 测试 API

### 1. 健康检查

```bash
curl http://localhost:8080/health
```

### 2. 微信登录

```bash
curl -X POST http://localhost:8080/api/v1/auth/wechat/login \
  -H "Content-Type: application/json" \
  -d '{
    "code": "微信小程序返回的code"
  }'
```

响应示例：
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "id": 1,
      "open_id": "oxxxxxxxxxxxxxx",
      "nickname": "",
      "avatar_url": "",
      "gender": 0,
      "status": 1
    }
  }
}
```

### 3. 获取用户信息（需要 Token）

```bash
curl http://localhost:8080/api/v1/user/info \
  -H "Authorization: Bearer YOUR_TOKEN_HERE"
```

### 4. 更新用户信息（需要 Token）

```bash
curl -X PUT http://localhost:8080/api/v1/user/info \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" \
  -d '{
    "nickname": "张三",
    "avatar_url": "https://example.com/avatar.jpg",
    "gender": 1,
    "phone": "13800138000",
    "email": "zhangsan@example.com"
  }'
```

## 使用 Postman 测试

1. 导入以下环境变量：
   - `base_url`: http://localhost:8080
   - `token`: (登录后获取的 token)

2. 创建请求集合：
   - 健康检查: GET {{base_url}}/health
   - 微信登录: POST {{base_url}}/api/v1/auth/wechat/login
   - 获取用户信息: GET {{base_url}}/api/v1/user/info
   - 更新用户信息: PUT {{base_url}}/api/v1/user/info

3. 在需要认证的请求中添加 Header：
   - Key: Authorization
   - Value: Bearer {{token}}

## 微信小程序端调用示例

```javascript
// 1. 微信登录
wx.login({
  success: (res) => {
    if (res.code) {
      // 发送 code 到后端
      wx.request({
        url: 'https://your-domain.com/api/v1/auth/wechat/login',
        method: 'POST',
        data: {
          code: res.code
        },
        success: (res) => {
          if (res.data.code === 0) {
            // 保存 token
            wx.setStorageSync('token', res.data.data.token)
            // 保存用户信息
            wx.setStorageSync('userInfo', res.data.data.user)
          }
        }
      })
    }
  }
})

// 2. 获取用户信息
wx.request({
  url: 'https://your-domain.com/api/v1/user/info',
  method: 'GET',
  header: {
    'Authorization': 'Bearer ' + wx.getStorageSync('token')
  },
  success: (res) => {
    console.log(res.data)
  }
})

// 3. 更新用户信息
wx.request({
  url: 'https://your-domain.com/api/v1/user/info',
  method: 'PUT',
  header: {
    'Authorization': 'Bearer ' + wx.getStorageSync('token')
  },
  data: {
    nickname: '张三',
    avatar_url: 'https://example.com/avatar.jpg',
    gender: 1
  },
  success: (res) => {
    console.log(res.data)
  }
})
```

## 错误码说明

| 错误码 | 说明 |
|--------|------|
| 0 | 成功 |
| 400 | 请求参数错误 |
| 401 | 未授权/Token 无效 |
| 403 | 禁止访问 |
| 404 | 资源不存在 |
| 500 | 服务器内部错误 |

