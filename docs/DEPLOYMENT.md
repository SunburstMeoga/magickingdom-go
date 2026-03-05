# 部署指南

## 本地开发环境

### 1. 使用 Docker Compose 启动 MySQL

```bash
# 启动 MySQL
docker-compose up -d mysql

# 查看日志
docker-compose logs -f mysql

# 停止
docker-compose down
```

### 2. 配置数据库

修改 `configs/config.yaml`：

```yaml
database:
  host: localhost
  port: 3306
  user: root
  password: root123456
  dbname: magickingdom
```

### 3. 运行应用

```bash
# 开发模式（热重载）
make dev

# 或普通模式
make run
```

## 生产环境部署

### 方式一：直接部署

#### 1. 编译

```bash
# 在本地编译
GOOS=linux GOARCH=amd64 go build -o magickingdom-go main.go

# 或在服务器上编译
make build
```

#### 2. 上传文件

```bash
# 上传到服务器
scp -r bin/ configs/ migrations/ user@server:/path/to/app/
```

#### 3. 配置生产环境

创建 `configs/config.prod.yaml`：

```yaml
server:
  port: 8080
  mode: release

database:
  host: your-db-host
  port: 3306
  user: your-db-user
  password: your-db-password
  dbname: magickingdom

jwt:
  secret: your-production-secret-key
  expire_hours: 168

wechat:
  app_id: your_production_app_id
  app_secret: your_production_app_secret

log:
  level: info
  format: json
  output: file
  file_path: /var/log/magickingdom/app.log
```

#### 4. 运行

```bash
# 使用 systemd 管理
sudo systemctl start magickingdom
sudo systemctl enable magickingdom
```

创建 systemd 服务文件 `/etc/systemd/system/magickingdom.service`：

```ini
[Unit]
Description=MagicKingdom Go Service
After=network.target mysql.service

[Service]
Type=simple
User=www-data
WorkingDirectory=/path/to/app
Environment="CONFIG_PATH=/path/to/app/configs/config.prod.yaml"
ExecStart=/path/to/app/bin/magickingdom-go
Restart=on-failure
RestartSec=5s

[Install]
WantedBy=multi-user.target
```

### 方式二：Docker 部署

#### 1. 构建镜像

```bash
docker build -t magickingdom-go:latest .
```

#### 2. 运行容器

```bash
docker run -d \
  --name magickingdom-app \
  -p 8080:8080 \
  -v /path/to/config.prod.yaml:/root/configs/config.yaml \
  -e CONFIG_PATH=/root/configs/config.yaml \
  magickingdom-go:latest
```

#### 3. 使用 Docker Compose

创建 `docker-compose.prod.yml`：

```yaml
version: '3.8'

services:
  mysql:
    image: mysql:8.0
    restart: always
    environment:
      MYSQL_ROOT_PASSWORD: ${DB_ROOT_PASSWORD}
      MYSQL_DATABASE: magickingdom
    volumes:
      - mysql_data:/var/lib/mysql
    networks:
      - app-network

  app:
    build: .
    restart: always
    ports:
      - "8080:8080"
    environment:
      CONFIG_PATH: /root/configs/config.yaml
    volumes:
      - ./configs/config.prod.yaml:/root/configs/config.yaml
    depends_on:
      - mysql
    networks:
      - app-network

volumes:
  mysql_data:

networks:
  app-network:
    driver: bridge
```

运行：

```bash
docker-compose -f docker-compose.prod.yml up -d
```

### 方式三：Kubernetes 部署

#### 1. 创建 ConfigMap

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: magickingdom-config
data:
  config.yaml: |
    server:
      port: 8080
      mode: release
    # ... 其他配置
```

#### 2. 创建 Deployment

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: magickingdom-app
spec:
  replicas: 3
  selector:
    matchLabels:
      app: magickingdom
  template:
    metadata:
      labels:
        app: magickingdom
    spec:
      containers:
      - name: app
        image: magickingdom-go:latest
        ports:
        - containerPort: 8080
        env:
        - name: CONFIG_PATH
          value: /etc/config/config.yaml
        volumeMounts:
        - name: config
          mountPath: /etc/config
      volumes:
      - name: config
        configMap:
          name: magickingdom-config
```

#### 3. 创建 Service

```yaml
apiVersion: v1
kind: Service
metadata:
  name: magickingdom-service
spec:
  selector:
    app: magickingdom
  ports:
  - protocol: TCP
    port: 80
    targetPort: 8080
  type: LoadBalancer
```

## Nginx 反向代理

```nginx
server {
    listen 80;
    server_name your-domain.com;

    location / {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

## HTTPS 配置

使用 Let's Encrypt：

```bash
# 安装 certbot
sudo apt-get install certbot python3-certbot-nginx

# 获取证书
sudo certbot --nginx -d your-domain.com

# 自动续期
sudo certbot renew --dry-run
```

## 监控和日志

### 1. 日志收集

使用 ELK Stack 或 Loki：

```yaml
# docker-compose.yml 添加
  loki:
    image: grafana/loki:latest
    ports:
      - "3100:3100"
```

### 2. 性能监控

使用 Prometheus + Grafana：

```yaml
  prometheus:
    image: prom/prometheus:latest
    volumes:
      - ./prometheus.yml:/etc/prometheus/prometheus.yml
    ports:
      - "9090:9090"

  grafana:
    image: grafana/grafana:latest
    ports:
      - "3000:3000"
```

## 数据库备份

```bash
# 备份
docker exec magickingdom-mysql mysqldump -u root -p magickingdom > backup.sql

# 恢复
docker exec -i magickingdom-mysql mysql -u root -p magickingdom < backup.sql
```

## 性能优化

1. 启用 Gzip 压缩
2. 使用 Redis 缓存
3. 数据库连接池优化
4. 静态资源 CDN
5. 负载均衡

## 安全建议

1. 使用强密码和密钥
2. 定期更新依赖
3. 启用 HTTPS
4. 配置防火墙
5. 限制 API 访问频率
6. 定期备份数据

