# Immortal's Blog

## 部署
nginx可以复制`config/nginx`内的`nginx.conf`
修改其中的
```
server_name blog.local;  # 修改为你自己的域名或 IP
```
然后使用
```
sudo certbot --nginx
```
提供HTTPS支持
