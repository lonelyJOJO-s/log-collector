# logv1
该项目作为log项目（主要为了解决日志的存储和可视化搜索的需求）的一个模块，主要实现了监听日志，抓取配置，发送消息的功能
详细文档参考[日志存储](https://momenta.feishu.cn/docx/SJqNdrsjxozSIXx3h8zc6xDtnGF)
## 功能特性
### 1. 日志监听
实时监听日志文件的最新内容
### 2. 日志配置动态更新
配置信息通过etcd存储，log_collector动态拉取更新配置文件
### 3. 日志回流kafka
监听到的日志信息，配合拉取配置文件中的topic，发送到kafka队列

## 软件架构
1. 整体使用golang开发
2. 消息队列采用kafka(相对于rabbitmq具有处理大流量数据的能力)
3. 配置存储采用etcd
4. 日志搜集采用tail包

## 快速开始
### 1. 安装依赖
#### 1.1 基本依赖
本次实验golang版本为1.20，通过以下命令安装依赖包
```bash
go mod init 
go mod tidy
```
#### 1.2 etcd 和 kafka安装
通过docker-compose.yml文件一键部署
### 2. 运行程序
测试环境
```bash
go run main.go
```
正式环境
```bash
go build -o logv1 
./logv1
```

## extentions:
### docker-compose
部署dockers
```bash
docker-compose up -d
```
修改docker-compose.yml后dockers重新部署
```bash
docker-compose down
docker-compose up -d --build
# -d 选项表示在后台运行。
# --build 选项表示重新构建容器
```


### kafka命令基本工具
1. 进入kafka容器
```bash
sudo docker exec -it kafka /bin/bash
```

2. create topic
```bash
kafka-topics.sh --bootstrap-server localhost:9092 --create --topic --partition 3 --replication-factor 1
```

3. 查询topic 状态
```bash
kafka-topics.sh --bootstrap-server localhost:9092 --describe --topic web_log
```

4. 发送kafka消息
``` bash
kafka-console-producer.sh --broker-list localhost:9092 --topic test
# 注意: 这里--bootstrap-server 变成了--broker-list， 理论上不应该出现kafka内不同脚本的版本不一致，原因待排查
```

5. 查看kafka数据
```bash
docker exec kafka kafka-console-consumer.sh --bootstrap-server localhost:9092 --topic web_log --from-beginning
```
### etcd命令

1. 进入etcd 
```bash
docker exec -it etcd /bin/sh
```

2. 查看所有数据
```bash
etcdctl get --prefix "" --keys-only=false
```


3. 通过命令行客户端将配置信息写入etcd
```
etcdctl --endpoints=http://localhost:2379 put /logv1/124.221.101.169/collect '[{"path":"/var/log/log_test/web.log","topic":"web_log"},{"path":"/var/log/log_test/sys.log","topic":"sys_log"}]'

etcdctl --endpoints=http://localhost:2379 put /logv1/124.221.101.169/collect '[{"path":"/var/log/log_test/web1.log","topic":"web_log"},{"path":"/var/log/log_test/sys.log","topic":"sys_log"}]'

etcdctl --endpoints=http://localhost:2379 put /logv1/collect '[{"path":"/var/log/log_test/web.log","topic":"web_log"},{"path":"/var/log/log_test/sys.log","topic":"sys_log"}]'
```

### 注意项
#### echo >> 的错误
sudo sh -c "echo  >> "
#### 注意需要配值hosts
127.0.0.1 kafka

### todo
- [x] 使用viper构建配置
- [ ] 使用ip配置拉取内容，配置分布式kafka and etcd

### what else
本次实验参考[黄忠德的博客](https://huangzhongde.cn/post/2020-03-03-golang_devops_logAgent_1_write_log_to_kafka/)
相当于一次复现

后续分布式节点dockerfile
```Dockerfile
basic configuration
version: "3"

networks:
  app-kafka:
    driver: bridge

services:
  zookeeper:
    container_name: zookeeper
    hostname: zookeeper
    image: zookeeper:3.4.14
    restart: always
    networks:
      - app-kafka
  kafka-0:
    container_name: kafka-0
    hostname: kafka-0
    image: bitnami/kafka:2.4.0
    restart: always
    
    environment: 
      - KAFKA_CFG_ZOOKEEPER_CONNECT=zookeeper:2181
      - ALLOW_PLAINTEXT_LISTENER=yes
      - KAFKA_CFG_BROKER_ID=0
      - KAFKA_CFG_LISTENER_SECURITY_PROTOCOL_MAP=INTERNAL:PLAINTEXT,EXTERNAL:PLAINTEXT
      - KAFKA_CFG_LISTENERS=INTERNAL://:9092,EXTERNAL://0.0.0.0:19092
      - KAFKA_CFG_ADVERTISED_LISTENERS=INTERNAL://kafka-0:9092,EXTERNAL://localhost:19092
      - KAFKA_CFG_INTER_BROKER_LISTENER_NAME=INTERNAL
    ports:
    - 127.0.0.1:9092:9092
    - 127.0.0.1:19092:19092
    networks:
      - app-kafka

  kafka-1:
    container_name: kafka-1
    hostname: kafka-1
    image: bitnami/kafka:2.4.0
    restart: always
    
    environment: 
      - KAFKA_CFG_ZOOKEEPER_CONNECT=zookeeper:2181
      - ALLOW_PLAINTEXT_LISTENER=yes
      - KAFKA_CFG_BROKER_ID=1
      - KAFKA_CFG_LISTENER_SECURITY_PROTOCOL_MAP=INTERNAL:PLAINTEXT,EXTERNAL:PLAINTEXT
      - KAFKA_CFG_LISTENERS=INTERNAL://:9092,EXTERNAL://0.0.0.0:29092
      - KAFKA_CFG_ADVERTISED_LISTENERS=INTERNAL://kafka-1:9092,EXTERNAL://localhost:29092
      - KAFKA_CFG_INTER_BROKER_LISTENER_NAME=INTERNAL
    ports:
    - 127.0.0.1:9092:9092
    - 127.0.0.1:29092:29092
    networks:
      - app-kafka

  kafka-2:
    container_name: kafka-2
    hostname: kafka-2
    image: bitnami/kafka:2.4.0
    restart: always
    
    environment: 
      - KAFKA_CFG_ZOOKEEPER_CONNECT=zookeeper:2181
      - ALLOW_PLAINTEXT_LISTENER=yes
      - KAFKA_CFG_BROKER_ID=2
      - KAFKA_CFG_LISTENER_SECURITY_PROTOCOL_MAP=INTERNAL:PLAINTEXT,EXTERNAL:PLAINTEXT
      - KAFKA_CFG_LISTENERS=INTERNAL://:9092,EXTERNAL://0.0.0.0:39092
      - KAFKA_CFG_ADVERTISED_LISTENERS=INTERNAL://kafka-2:9092,EXTERNAL://localhost:39092
      - KAFKA_CFG_INTER_BROKER_LISTENER_NAME=INTERNAL
    ports:
    - 127.0.0.1:9092:9092
    - 127.0.0.1:39092:39092
    networks:
      - app-kafka
    
  etcd:
    container_name: etcd
    image: bitnami/etcd:3
    restart: always
    environment:
      - ALLOW_NONE_AUTHENTICATION=yes
    ports: 
    - 127.0.0.1:2379:2379
    networks: 
      - app-kafka

```



