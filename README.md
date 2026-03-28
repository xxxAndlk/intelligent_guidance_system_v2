# 智慧医疗导诊系统 v2 (Intelligent Guidance System v2)

基于 Go + Kratos + DDD + RAG + Agent 的微服务架构。

## 架构概览

```
┌─────────────────────────────────────────────────────────────┐
│                         API Gateway                          │
│                    (Kratos Gateway)                          │
└────────────────────────┬────────────────────────────────────┘
                         │
    ┌────────────────────┼────────────────────┐
    │                    │                    │
    ▼                    ▼                    ▼
┌─────────┐      ┌─────────────┐      ┌─────────────┐
│patient- │      │  doctor-    │      │   admin-    │
│api      │      │  api        │      │   api       │
│(BFF)    │      │  (BFF)      │      │   (BFF)     │
└────┬────┘      └──────┬──────┘      └──────┬──────┘
     │                  │                    │
     └──────────────────┼────────────────────┘
                        │
    ┌───────────────────┼───────────────────┐
    │                   │                   │
    ▼                   ▼                   ▼
┌──────────┐    ┌──────────┐    ┌──────────┐
│ ai-svc   │    │ patient- │    │ doctor-  │
│(核心)    │    │ svc      │    │ svc      │
└──────────┘    └──────────┘    └──────────┘
     │               │               │
     └───────────────┼───────────────┘
                     │
    ┌────────────────┼────────────────┐
    ▼                ▼                ▼
┌──────────┐  ┌──────────┐  ┌──────────┐
│ medical- │  │ dept-svc │  │ auth-svc │
│ svc      │  │          │  │          │
└──────────┘  └──────────┘  └──────────┘
```

## 技术栈

- **框架**: Kratos v2.7.2
- **语言**: Go 1.21+
- **通信**: gRPC + HTTP/REST
- **数据库**: MySQL 8.0, Redis 7.4
- **消息队列**: RabbitMQ 3.13
- **向量数据库**: Qdrant
- **服务注册**: Consul
- **LLM**: DeepSeek API (via SiliconFlow)
- **向量嵌入**: BAAI/bge-large-zh-v1.5

## 服务列表

| 服务 | 端口 (HTTP/gRPC) | 描述 |
|------|------------------|------|
| ai-service | 8001/9001 | AI导诊核心服务（RAG + Agent） |
| patient-service | 8002/9002 | 患者管理服务 |
| doctor-service | 8003/9003 | 医生管理服务 |
| medical-service | 8004/9004 | 病历管理服务 |
| department-service | 8005/9005 | 科室管理服务 |
| auth-service | 8006/9006 | 认证授权服务 |
| registration-service | 8007/9007 | 挂号预约服务 |
| surgical-service | 8008/9008 | 手术管理服务 |
| drug-service | 8009/9009 | 药品管理服务 |
| payment-service | 8010/9010 | 支付结算服务 |
| notification-service | 8011/9011 | 消息通知服务 |
| file-service | 8012/9012 | 文件存储服务 |

## 快速开始

### 1. 环境要求

- Go 1.21+
- Docker & Docker Compose
- Protocol Buffers compiler

### 2. 安装依赖

```bash
make init
```

### 3. 生成代码

```bash
make generate
```

### 4. 启动基础设施

```bash
cd deploy
docker-compose up -d mysql redis qdrant rabbitmq consul
```

### 5. 运行服务

```bash
# 运行 AI 服务
make run-ai

# 运行患者服务
make run-patient
```

## 目录结构

```
intelligent_guidance_system_v2/
├── api/                    # Protobuf 定义
│   ├── ai/v1/
│   ├── patient/v1/
│   ├── doctor/v1/
│   └── ...
├── app/                    # BFF 层
│   ├── patient-api/
│   ├── doctor-api/
│   └── admin-api/
├── service/                # 领域服务
│   ├── ai/
│   ├── patient/
│   ├── doctor/
│   └── ...
├── common/                 # 公共库
│   ├── middleware/
│   ├── pkg/
│   └── proto/
├── deploy/                 # 部署配置
│   ├── docker-compose.yml
│   ├── k8s/
│   └── config/
└── scripts/                # 脚本
```

## DDD 架构

每个服务内部采用 DDD 分层架构：

```
service/{name}/
├── internal/
│   ├── domain/            # 领域层
│   │   ├── aggregate/     # 聚合根
│   │   ├── entity/        # 实体
│   │   ├── vo/            # 值对象
│   │   ├── event/         # 领域事件
│   │   └── repository/    # 仓储接口
│   ├── biz/               # 应用层
│   │   ├── usecase/       # 用例
│   │   └── dto/           # DTO
│   ├── data/              # 基础设施层
│   │   ├── mysql/
│   │   ├── redis/
│   │   └── vector/        # 向量数据库
│   ├── acl/               # 防腐层
│   ├── service/           # 接入层 (proto实现)
│   └── server/            # 服务配置
```

## 核心功能

### AI 导诊 (核心)

1. **症状收集**: 多轮对话收集患者症状
2. **RAG检索**: 
   - 向量检索相似疾病
   - 向量检索相关医生
   - 向量检索推荐药品
3. **Agent推理**:
   - ReAct 模式多步推理
   - Skills 工具调用
   - 智能推荐科室和医生

### 智能推荐

- 疾病推荐（基于向量相似度）
- 科室推荐
- 医生推荐（考虑职称、专家级别）
- 药品推荐

## API 文档

启动服务后访问：
- AI Service: http://localhost:8001/docs
- Patient Service: http://localhost:8002/docs

## 开发指南

### 添加新服务

1. 在 `api/` 创建 proto 定义
2. 在 `service/` 创建服务目录
3. 实现 domain, biz, data, acl, service 层
4. 配置 wire 依赖注入
5. 添加 Dockerfile

### 生成代码

```bash
# 生成所有 proto
cd api/{service}/v1
protoc --go_out=. --go_opt=paths=source_relative \
       --go-grpc_out=. --go-grpc_opt=paths=source_relative \
       *.proto
```

### 运行测试

```bash
make test
```

## 部署

### Docker Compose

```bash
cd deploy
docker-compose up -d
```

### Kubernetes

```bash
kubectl apply -f deploy/k8s/
```

## 文档

- [完整开发手册](docs/完整开发手册.md)
- [DDD架构分析](docs/DDD架构分析文档.md)
- [Go微服务迁移架构](docs/Go微服务迁移架构设计.md)
- [配置文件参考](docs/配置文件参考.md)

## License

MIT License
