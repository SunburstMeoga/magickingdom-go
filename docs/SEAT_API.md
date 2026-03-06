# 座位管理接口测试文档

## 接口列表

### 1. 获取座位布局
**接口**: `GET /api/v1/seats/layout`  
**认证**: 不需要  
**说明**: 获取所有座位的布局信息

```bash
curl http://localhost:8080/api/v1/seats/layout
```

---

### 2. 获取用户当前座位
**接口**: `GET /api/v1/seats/my-seat`  
**认证**: 需要 JWT Token  
**说明**: 获取当前登录用户的座位信息

```bash
curl -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  http://localhost:8080/api/v1/seats/my-seat
```

**响应示例**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "has_seat": true,
    "seat_id": "C09+",
    "seat_type": "card"
  }
}
```

---

### 3. 用户入座
**接口**: `POST /api/v1/seats/join`  
**认证**: 需要 JWT Token  
**说明**: 用户加入某个座位

```bash
curl -X POST \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "seat_id": "C09+",
    "seat_type": "card"
  }' \
  http://localhost:8080/api/v1/seats/join
```

**请求参数**:
- `seat_id`: 座位ID（必填）
- `seat_type`: 座位类型（必填），可选值：`card`、`vip`、`table`、`first-class`

**业务规则**:
- 用户只能同时占用一个座位
- 如果用户已经在其他座位上，需要先离座
- 如果用户已经在目标座位上，直接返回成功

---

### 4. 用户离座
**接口**: `POST /api/v1/seats/leave`  
**认证**: 需要 JWT Token  
**说明**: 用户离开当前座位

```bash
curl -X POST \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  http://localhost:8080/api/v1/seats/leave
```

---

### 5. 获取座位占用信息
**接口**: `GET /api/v1/seats/occupancy?seat_id=C09+`  
**认证**: 不需要  
**说明**: 获取某个座位的入座人数和用户信息

```bash
curl "http://localhost:8080/api/v1/seats/occupancy?seat_id=C09%2B"
```

**响应示例**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "seat_id": "C09+",
    "seat_type": "card",
    "occupied_num": 2,
    "users": [
      {
        "user_id": 1,
        "nickname": "张三",
        "avatar_url": "https://example.com/avatar1.jpg",
        "joined_at": "2026-03-06T22:00:00Z"
      },
      {
        "user_id": 2,
        "nickname": "李四",
        "avatar_url": "https://example.com/avatar2.jpg",
        "joined_at": "2026-03-06T22:05:00Z"
      }
    ]
  }
}
```

---

## 前端使用流程

### 1. 页面加载时
```javascript
// 1. 获取座位布局
const layoutRes = await fetch('/api/v1/seats/layout');
const layout = await layoutRes.json();

// 2. 检查用户是否已入座（需要登录）
const mySeatRes = await fetch('/api/v1/seats/my-seat', {
  headers: { 'Authorization': `Bearer ${token}` }
});
const mySeat = await mySeatRes.json();

if (mySeat.data.has_seat) {
  console.log('用户当前在座位:', mySeat.data.seat_id);
}
```

### 2. 用户点击座位时
```javascript
// 获取该座位的占用信息
const occupancyRes = await fetch(`/api/v1/seats/occupancy?seat_id=${seatId}`);
const occupancy = await occupancyRes.json();

console.log('当前入座人数:', occupancy.data.occupied_num);
console.log('入座用户:', occupancy.data.users);
```

### 3. 用户入座
```javascript
const joinRes = await fetch('/api/v1/seats/join', {
  method: 'POST',
  headers: {
    'Authorization': `Bearer ${token}`,
    'Content-Type': 'application/json'
  },
  body: JSON.stringify({
    seat_id: 'C09+',
    seat_type: 'card'
  })
});

if (joinRes.ok) {
  console.log('入座成功');
}
```

### 4. 用户离座
```javascript
const leaveRes = await fetch('/api/v1/seats/leave', {
  method: 'POST',
  headers: { 'Authorization': `Bearer ${token}` }
});

if (leaveRes.ok) {
  console.log('离座成功');
}
```

---

## 数据库表结构

### seat_occupancies 表
| 字段 | 类型 | 说明 |
|------|------|------|
| id | bigint | 主键 |
| user_id | bigint | 用户ID（外键关联users表） |
| seat_id | varchar(50) | 座位ID |
| seat_type | varchar(20) | 座位类型 |
| status | int | 状态：1-占用中，0-已离座 |
| created_at | datetime | 创建时间（入座时间） |
| updated_at | datetime | 更新时间 |
| deleted_at | datetime | 软删除时间 |

---

## 注意事项

1. **并发控制**: 当前实现没有处理并发入座的情况，如果需要严格控制座位容量，建议添加数据库锁或使用Redis分布式锁

2. **座位容量**: 当前实现允许多人入座同一个座位，如果需要限制座位容量，需要在入座时检查当前人数

3. **实时更新**: 建议使用WebSocket或轮询来实时更新座位占用状态

4. **历史记录**: 当前使用status字段标记离座，保留了历史记录。如果需要完全删除记录，可以使用硬删除

