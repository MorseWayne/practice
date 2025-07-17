# Forum 项目部署文档
本项目旨在实现一个论坛系统，前端主要使用js完成，后端使用rust语言完成
<img width="2492" height="1167" alt="image" src="https://github.com/user-attachments/assets/cda2610d-5cbc-4491-bb40-453d3727ed2a" />

## Nginx 反向代理配置
### 1. 配置文件路径
Nginx 站点配置文件位于 `/etc/nginx/sites-available/forum`，内容如下：

```nginx
server {
    listen 80;
    server_name _;  # 无域名时保留，有域名替换为具体域名（如 forum.example.com）

    location / {
        proxy_pass http://localhost:8080;  # 指向本地后端服务（Actix-web 默认端口）
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }
}
```

### 2. 启用配置
- 创建软链接到 `sites-enabled` 目录：
  ```bash
  sudo ln -s /etc/nginx/sites-available/forum /etc/nginx/sites-enabled/
  ```
- 禁用默认配置（可选，避免冲突）：
  ```bash
  sudo rm /etc/nginx/sites-enabled/default  # 删除默认配置软链接
  ```

### 3. 验证配置
- 检查 Nginx 配置语法：
  ```bash
  sudo nginx -t
  ```
- 重启 Nginx 使配置生效：
  ```bash
  sudo nginx -s reload
  ```

## 运行指导
### 1. 启动后端服务
进入后端项目目录，使用 Cargo 启动 Actix-web 服务：
```bash
cd backend/forum_server
cargo run  # 默认监听 0.0.0.0:8080
```

### 2. 测试访问
- 本地测试：访问 `http://localhost` 应返回后端的 "Hello, World!" 响应。
- 远程访问：通过服务器公网 IP 或域名访问（需确保服务器 80 端口已开放防火墙）。

### 3. 常见问题排查
- **显示 Nginx 欢迎页**：检查 `sites-enabled` 目录是否存在默认配置软链接（`default`），删除后重启 Nginx。
- **502 Bad Gateway**：确认后端服务已启动（`curl http://localhost:8080` 测试），或检查 `proxy_pass` 地址是否正确。
- **域名无法解析**：确保 `server_name` 配置与域名解析记录一致（如使用 `forum.example.com`，需在 DNS 服务商添加 A 记录指向服务器 IP）。
