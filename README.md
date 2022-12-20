# 1. 准备
`make deps` 安装必要的工具， 如 `go-bindata`

验证 `go-bindata`

`make assets` 打包静态资源和证书文件， 见

+ ngrok/client/assets/assets_debug.go
+ ngrok/server/assets/assets_debug.go

> 对应的release版本是 */assets_release.go

## 1.1. 重新生成证书

``` bash
export DOMAIN=dev.domain.com

openssl genrsa -out rootCA.key 4096

openssl req -x509 -new -nodes -key rootCA.key -subj "/CN=${DOMAIN}" -days 5000 -out rootCA.pem

cp rootCA.pem assets/client/tls/ngrokroot.crt

openssl genrsa -out device.key 4096

openssl req -new -key device.key -subj "/CN=${DOMAIN}" -out device.csr

openssl x509 -req -in device.csr -CA rootCA.pem -CAkey rootCA.key -CAcreateserial -out device.crt -days 5000

cp device.crt assets/server/tls/snakeoil.crt

cp device.key assets/server/tls/snakeoil.key
```

> 生产环境如果要求是真正的可信任证书， 将`ngrok/client/release.go` 的 `useInsecureSkipVerify` 返回 false

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