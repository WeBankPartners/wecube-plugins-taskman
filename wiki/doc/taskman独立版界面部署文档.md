# Taskman 独立版本界面部署文档

taskman的独立前端镜像，是为了在办公网就能提请求、做审批什么的，不用进入VDI。需要部署在办公网能开通策略访问、又能访问到MGMT的wecube后台服务的区域。

### docker-compose部署

应用包解压后会有镜像文件：taskman-standalone-ui.tar
```bash
# 导入docker镜像
docker load --input taskman-standalone-ui.tar 
```
#### 环境变量整理

**请按需修正一下环境变量值**

```bash
# 独立版本界面版本
VERSION='v1.1.5'
# 部署的主机ip
HOSTIP='10.0.0.1'
# 网页访问的端口
PORT='8900'
# wecube platform-gateway的ip地址
WECUBE_GATEWAY_HOST='10.0.0.2'
# wecube platform-gateway的端口
WECUBE_GATEWAY_PORT='8005'
```

#### 准备yaml内容

在本地准备以下yaml文件

taskman-standalone-ui.yml

```yaml
version: '3'
services:
  taskman-standalone:
    image: taskman-standalone-ui:{{VERSION}}
    container_name: taskman-standalone-{{VERSION}}
    restart: always
    volumes:
      - /etc/localtime:/etc/localtime
      - /data/app/taskman-ui/log:/var/log/nginx/
    ports:
      - "{{HOSTIP}}:{{PORT}}:8080"
    environment:
      - GATEWAY_HOST={{WECUBE_GATEWAY_HOST}}
      - GATEWAY_PORT={{WECUBE_GATEWAY_PORT}}
      - PUBLIC_DOMAIN={{HOSTIP}}:{{PORT}}
      - TZ=Asia/Shanghai
    command: /bin/bash -c "/etc/nginx/start_taskman.sh"
```

#### 修正yaml内容的值

```bash
# 请修改以下变量为正确值
# 粘贴以上整理的环境变量


sed -i "s/{{VERSION}}/$VERSION/g" taskman-standalone-ui.yml
sed -i "s/{{HOSTIP}}/$HOSTIP/g" taskman-standalone-ui.yml
sed -i "s/{{PORT}}/$PORT/g" taskman-standalone-ui.yml
sed -i "s/{{WECUBE_GATEWAY_HOST}}/$WECUBE_GATEWAY_HOST/g" taskman-standalone-ui.yml
sed -i "s/{{WECUBE_GATEWAY_PORT}}/$WECUBE_GATEWAY_PORT/g" taskman-standalone-ui.yml
```

#### 启动docker容器
- 启动taskman-standalone-ui服务
  ```bash
  # 建映射出来的日志路径
  mkdir -p /data/app/taskman-ui/log
  chmod 777 /data/app/taskman-ui/log
  # 启动容器
  docker-compose -f taskman-standalone-ui.yml up -d
  ```

- 验证服务，访问 http://{{HOSTIP}}:{{PORT}}

