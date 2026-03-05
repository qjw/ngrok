# 1. 准备

``` bash
# 服务器域名， 后续使用 <subdomain>.king.com 来做内网穿透
export domain="king.com"
# 创建 ca 私钥 
openssl genrsa -out ca.key 4096
# 创建 ca 根证书
openssl req -x509 -new -nodes -key ca.key -subj "/CN=${domain}" -days 3650 -out ca.crt
	
# 创建 客户端私钥
openssl genrsa -out client.key 4096
# 创建客户端csr
openssl req -new -key client.key -subj "/CN=${domain}" -out client.csr
# ca签发客户端证书 (无需sans)
openssl x509 -req -in client.csr -CA ca.crt -CAkey ca.key -CAcreateserial -out client.crt -days 3650

# 创建 服务器私钥
openssl genrsa -out server.key 4096
# 创建服务器csr
openssl req -new -key server.key -subj "/CN=${domain}" -out server.csr
# ca签发服务器证书
openssl x509 -req -in server.csr -CA ca.crt -CAkey ca.key -CAcreateserial -out server.crt -days 3650 \
	-extfile <(printf "subjectAltName=DNS:${domain}")

# 拷贝到项目目录
cp ca.crt assets/client/tls/ngrokroot.crt
cp client.crt assets/client/tls/client.crt
cp client.key assets/client/tls/client.key
cp server.crt assets/server/tls/server.crt
cp server.key assets/server/tls/server.key
rm *.key *.csr *.crt *.srl
```

# 2. 编译
``` bash
# 测试
$ make client
$ make server
# 生产
$ make release-all
$ make release-client
$ make release-server
```

生成的二进制文件

```
build/
├── ngrokd-debug
├── ngrok-debug
├── ngrokd-release
└── ngrok-release
```

# 3. 部署服务端
`ngrokd -domain=dev.domain.com -httpAddr=:8002 -httpsAddr=:9082 -tunnelAddr=:4443`

> 如果是云服务器， 需要放开 80、443、4443。 80、443 通过nginx代理， 4443用于nginx客户端直接连接

通过nginx代理80端口

``` nginx
upstream ngrok {
    server 127.0.0.1:8002;
    keepalive 64;
}
server {
	listen 80;
	server_name *.dev.domain.com;
	location / {

		proxy_set_header X-Real-IP $remote_addr;
		proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
		proxy_set_header Host  $http_host:8002;
		proxy_set_header X-Nginx-Proxy true;
		proxy_set_header Connection "";
		proxy_pass      http://ngrok ;

	}
}
```

通过监听进程
`/etc/supervisor/conf.d/ngrok.conf`

``` conf
[program:ngrok]
command=ngrokd  -domain="dev.domain.com" -httpAddr=":8002" -httpsAddr=":9082" -tunnelAddr=":4443"
process_name=ngrokd
autostart=true
autorestart=true
```

# 4. 客户端
``` conf
server_addr: dev.domain.com:4443
trust_host_root_certs: false
```
`./ngrok -subdomain test -proto=http -config=./ngrok.cfg 54321`

```
ngrok

Tunnel Status                 online
Version                       1.7/1.7
Forwarding                    http://test.dev.domain.com:8002 -> 127.0.0.1:54321
Web Interface                 127.0.0.1:4040
# Conn                        0                                                                                                                    
Avg Conn Time                 0.00ms
```

当公网访问 <http://test.dev.domain.com> 时， 流量自动导入本机 54321 端口

## 4.1. tcp
``` conf
server_addr: dev.domain.com:4443
trust_host_root_certs: false
tunnels:
 lot:
  remote_port: 8888
  proto:
   tcp: 8880
```

`ngrok -config ./ngrok.cfg start lot`

```
ngrok

Tunnel Status                 online
Version                       1.7/1.7
Forwarding                    tcp://dev.domain.com:8888 -> 127.0.0.1:8880
Web Interface                 127.0.0.1:4040
# Conn                        0
Avg Conn Time                 0.00ms
```

## 4.2. 合并
``` conf
server_addr: dev.domain.com:4443
trust_host_root_certs: false
tunnels:
 lot:
  remote_port: 8888
  proto:
   tcp: 8880
 http:
  subdomain: test
  proto:
   http: 127.0.0.1:54321
```

`ngrok -config ./ngrok.cfg start http lot`

```
ngrok

Tunnel Status                 online
Version                       1.7/1.7
Forwarding                    tcp://dev.domain.com:8888 -> 127.0.0.1:8880
Forwarding                    http://test.dev.domain.com:8002 -> 127.0.0.1:54321
Web Interface                 127.0.0.1:4040
# Conn                        0
Avg Conn Time                 0.00ms
```