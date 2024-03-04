# BUPTreasure
> 本项目是一个基于Go语言Gin框架的校园抽奖系统服务端

## 项目简介
本系统主要功能包括
- 服务于报名小程序端：处理抽奖用户的报名请求
- 服务于抽奖网页端：创建抽奖活动、抽奖主功能、获取参与者和中奖者信息

## 项目结构
```
BUPTreasure
├── .github
│   └── workflows
│       └── Publish_Deploy.yml
├── buffer
│   └── lottery-config-users.json
├── internal
│   ├── ApiDTO
│   │   └── api.go
│   └── myDB
│       └── db-operation.go
├── main.go
├── go.mod
├── go.sum
├── Dockerfile
├── .gitignore
└── README.md
```

## 项目部署
### 1. 服务器环境
- 操作系统：Ubuntu 20.04.2 LTS

### 2. 项目依赖
- Go 1.21.5
- github.com/gin-gonic/gin
- github.com/go-sql-driver/mysql

### 3. 项目部署
- Docker部署
    - 项目Dockerfile：[Dockerfile](./Dockerfile)
    - 项目Docker云镜像：[buptreasure-server](https://hub.docker.com/r/acs991314/buptreasure-server)
    - 项目Docker Compose配置信息
      ```yaml
        version: '3'
        services:
          mysql-container:
            image: mysql:latest
            container_name: mysql-container
            networks:
              - bupt-net
            ports:
              - "3306:3306"
            environment:
              MYSQL_ROOT_PASSWORD: your_password
              MYSQL_DATABASE: BUPTreasure
            restart: unless-stopped
          buptreasure-server:
            image: acs991314/buptreasure-server:latest
            container_name: buptreasure-server
            ports:
              - "9000:8080"
            networks:
              - bupt-net
            volumes:
              - /home/nginx/html/temp:/app/buffer
            depends_on:
              - mysql-container
          buptreasure-client:
            image: acs991314/buptreasure-client:latest
            container_name: buptreasure-client
            ports:
            - "8080:8080"
            networks:
            - bupt-net
            volumes:
            - /home/nginx/html/temp:/app/src/buffer
            depends_on:
            - mysql-container
          nginx:
            image: nginx:latest
            container_name: nginx
            ports:
              - "443:443"
              - "80:80"
            volumes:
              - /home/nginx/log:/var/log/nginx
              - /home/nginx/html:/usr/share/nginx/html
              - /home/nginx/ssl:/etc/nginx/ssl
              - /home/nginx/conf/nginx.conf:/etc/nginx/nginx.conf
              - /home/nginx/conf/certs:/etc/nginx/conf/certs
              - /home/nginx/conf/conf.d:/etc/nginx/conf.d
            networks:
              - bupt-net
            depends_on:
              - buptreasure-client
              - buptreasure-server
            links:
              - buptreasure-client
              - buptreasure-server
        
        networks:
          bupt-net:
            external: true
      ```
- Github Actions自动化部署
    - 项目Github Actions配置文件：[Publish_Deploy.yml](./.github/workflows/Master_Push.yml)