# 智慧医疗系统 - Go+Kratos微服务迁移架构设计文档

## 一、项目概述

### 1.1 迁移目标

将现有的Java SpringBoot单体应用迁移到Go语言，采用Kratos微服务框架，实现：
- ✅ 微服务架构拆分（DDD领域驱动设计）
- ✅ 强化权限控制系统（RBAC+数据权限）
- ✅ 保留所有现有功能
- ✅ 提升系统性能和可扩展性

### 1.2 技术栈

| 层级 | 技术选型 | 版本 |
|-----|---------|------|
| 框架 | Kratos | v2.x |
| 语言 | Go | 1.21+ |
| 通信 | gRPC + HTTP | - |
| 数据库 | MySQL + Redis | 8.0 / 7.x |
| 消息队列 | RabbitMQ / Kafka | - |
| 注册中心 | Consul / etcd | - |
| 网关 | Kratos Gateway / Nginx | - |
| 链路追踪 | Jaeger / Zipkin | - |
| 监控 | Prometheus + Grafana | - |

### 1.3 迁移范围

| 功能模块 | 原Java | 迁移后Go | 状态 |
|---------|--------|---------|------|
| 患者管理 | UserService | user-service | ⬜ |
| 医生管理 | DoctorService | doctor-service | ⬜ |
| 科室管理 | DepartmentService | department-service | ⬜ |
| 挂号预约 | RegistrationService | registration-service | ⬜ |
| 病历管理 | MedicalService | medical-service | ⬜ |
| 手术管理 | SurgicalService | surgical-service | ⬜ |
| 药品管理 | DrugService | drug-service | ⬜ |
| AI导诊 | DiagnosisService | ai-service | ⬜ |
| 支付结算 | PaymentService | payment-service | ⬜ |
| 权限系统 | JWT拦截器 | auth-service | ⬜ |
| 文件存储 | AliOSS | file-service | ⬜ |
| 消息通知 | - | notification-service | ⬜ |

---

## 二、微服务架构设计

### 2.1 整体架构图

```
┌─────────────────────────────────────────────────────────────────────────┐
│                              接入层                                      │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐                 │
│  │ 微信小程序   │  │   Web管理端  │  │  小程序医生端 │                 │
│  └──────┬───────┘  └──────┬───────┘  └──────┬───────┘                 │
└─────────┼─────────────────┼─────────────────┼─────────────────────────┘
          │                 │                 │
          ▼                 ▼                 ▼
┌─────────────────────────────────────────────────────────────────────────┐
│                          API Gateway (Kratos Gateway)                   │
│  • 统一认证  • 限流熔断  • 路由转发  • 日志记录  • 协议转换              │
└─────────────────────────────────────────────────────────────────────────┘
          │                 │                 │
          ▼                 ▼                 ▼
┌─────────────────────────────────────────────────────────────────────────┐
│                          微服务层                                        │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                          │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐    │
│  │  user-svc   │  │ doctor-svc  │  │   dept-svc  │  │registra-svc │    │
│  │  (患者服务)  │  │  (医生服务)  │  │  (科室服务)  │  │ (挂号服务)  │    │
│  └─────────────┘  └─────────────┘  └─────────────┘  └─────────────┘    │
│                                                                          │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐    │
│  │ medical-svc │  │surgical-svc │  │  drug-svc   │  │  ai-svc     │    │
│  │  (病历服务)  │  │  (手术服务)  │  │  (药品服务)  │  │ (AI导诊)    │    │
│  └─────────────┘  └─────────────┘  └─────────────┘  └─────────────┘    │
│                                                                          │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐    │
│  │payment-svc  │  │  auth-svc   │  │  file-svc   │  │ notify-svc  │    │
│  │  (支付服务)  │  │  (认证服务)  │  │  (文件服务)  │  │ (通知服务)  │    │
│  └─────────────┘  └─────────────┘  └─────────────┘  └─────────────┘    │
│                                                                          │
│  ┌───────────────────────────────────────────────────────────────────┐ │
│  │                    支撑服务 (BFF - Backend For Frontend)         │ │
│  │  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐             │ │
│  │  │  patient-api│  │ doctor-api  │  │  admin-api  │             │ │
│  │  │  (患者端聚合) │  │(医生端聚合) │  │(管理端聚合) │             │ │
│  │  └─────────────┘  └─────────────┘  └─────────────┘             │ │
│  └───────────────────────────────────────────────────────────────────┘ │
│                                                                          │
└─────────────────────────────────────────────────────────────────────────┘
          │                 │                 │
          ▼                 ▼                 ▼
┌─────────────────────────────────────────────────────────────────────────┐
│                          基础设施层                                      │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                          │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐    │
│  │   MySQL     │  │   Redis     │  │ RabbitMQ    │  │   Consul    │    │
│  │ (主从集群)  │  │  (集群)     │  │  (消息队列)  │  │ (服务注册)  │    │
│  └─────────────┘  └─────────────┘  └─────────────┘  └─────────────┘    │
│                                                                          │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐    │
│  │  MinIO      │  │ Elasticsearch│  │   Jaeger    │  │ Prometheus  │    │
│  │ (对象存储)  │  │  (搜索引擎)  │  │  (链路追踪)  │  │  (监控)     │    │
│  └─────────────┘  └─────────────┘  └─────────────┘  └─────────────┘    │
│                                                                          │
└─────────────────────────────────────────────────────────────────────────┘
```

### 2.2 微服务拆分策略

#### 2.2.1 拆分原则

| 原则 | 说明 |
|-----|------|
| **业务边界** | 按业务领域拆分，如患者、医生、挂号等 |
| **数据独立性** | 每个服务拥有独立数据库，避免数据耦合 |
| **通信频率** | 高频通信的服务考虑合并 |
| **事务边界** | 强事务一致性的操作放在同一服务 |

#### 2.2.2 服务依赖图

```
┌─────────────────────────────────────────────────────────────────┐
│                        服务依赖关系                              │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│                         API Gateway                             │
│                              │                                  │
│          ┌───────────────────┼───────────────────┐              │
│          ▼                   ▼                   ▼              │
│  ┌─────────────┐    ┌─────────────┐    ┌─────────────┐        │
│  │ patient-api │    │ doctor-api  │    │  admin-api  │        │
│  │    (BFF)    │    │    (BFF)    │    │    (BFF)    │        │
│  └──────┬──────┘    └──────┬──────┘    └──────┬──────┘        │
│         │                   │                   │               │
│         └───────────────────┼───────────────────┘               │
│                             │                                   │
│    ┌────────────────────────┼────────────────────────┐         │
│    ▼                        ▼                        ▼         │
│ ┌──────┐  ┌──────┐  ┌──────┐  ┌──────┐  ┌──────┐  ┌──────┐   │
│ │user  │  │doctor│  │dept  │  │regis │  │medical│  │surgical│   │
│ │-svc  │  │-svc  │  │-svc  │  │-svc  │  │-svc   │  │-svc   │   │
│ └──┬───┘  └──┬───┘  └──┬───┘  └──┬───┘  └──┬───┘  └──┬───┘   │
│    │         │         │         │         │         │        │
│    └─────────┴─────────┴─────────┴─────────┴─────────┘        │
│                      │                                          │
│                      ▼                                          │
│              ┌─────────────┐                                    │
│              │  auth-svc   │                                    │
│              │  (认证中心)  │                                    │
│              └─────────────┘                                    │
│                                                                  │
│  ┌──────┐  ┌──────┐  ┌──────┐  ┌──────┐                       │
│  │payment│  │ ai   │  │ drug │  │ file │  (独立服务)           │
│  │-svc  │  │-svc  │  │-svc  │  │-svc  │                       │
│  └──────┘  └──────┘  └──────┘  └──────┘                       │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘

依赖规则：
- 上层可以调用下层
- 同级服务不直接调用（通过BFF聚合）
- 所有服务依赖auth-svc进行认证
- 跨服务事务使用Saga模式
```

### 2.3 DDD领域模型设计

#### 2.3.1 领域划分

```
┌─────────────────────────────────────────────────────────────────┐
│                     领域上下文划分                               │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│  ┌─────────────────┐  ┌─────────────────┐  ┌─────────────────┐ │
│  │   患者上下文    │  │   医生上下文    │  │   科室上下文    │ │
│  │  Patient BC     │  │   Doctor BC     │  │ Department BC   │ │
│  ├─────────────────┤  ├─────────────────┤  ├─────────────────┤ │
│  │ 聚合根: Patient │  │ 聚合根: Doctor  │  │ 聚合根: Dept    │ │
│  │ 实体: Address   │  │ 实体: Role      │  │ 实体: DeptDoc   │ │
│  │ 值对象: Phone   │  │ 值对象: Title   │  │ 值对象: Location│ │
│  └─────────────────┘  └─────────────────┘  └─────────────────┘ │
│                                                                  │
│  ┌─────────────────┐  ┌─────────────────┐  ┌─────────────────┐ │
│  │   挂号上下文    │  │   病历上下文    │  │   手术上下文    │ │
│  │Registration BC  │  │  Medical BC     │  │ Surgical BC     │ │
│  ├─────────────────┤  ├─────────────────┤  ├─────────────────┤ │
│  │ 聚合根: Regis   │  │ 聚合根: Record  │  │ 聚合根: Surgical│ │
│  │ 实体: Queue     │  │ 实体: Diagnosis │  │ 实体: Flow      │ │
│  │ 值对象: TimeSlot│  │ 值对象: Prescrip│  │ 值对象: Duration│ │
│  └─────────────────┘  └─────────────────┘  └─────────────────┘ │
│                                                                  │
│  ┌─────────────────┐  ┌─────────────────┐  ┌─────────────────┐ │
│  │   支付上下文    │  │   AI上下文      │  │   药品上下文    │ │
│  │  Payment BC     │  │    AI BC        │  │   Drug BC       │ │
│  ├─────────────────┤  ├─────────────────┤  ├─────────────────┤ │
│  │ 聚合根: Payment │  │ 聚合根: Diagnos │  │ 聚合根: Drug    │ │
│  │ 实体: Refund    │  │ 服务: DeepSeek  │  │ 实体: Stock     │ │
│  │ 值对象: Money   │  │ 值对象: Symptom │  │ 值对象: Category│ │
│  └─────────────────┘  └─────────────────┘  └─────────────────┘ │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
```

#### 2.3.2 聚合根设计

```go
// patient/internal/domain/aggregate/patient.go
package aggregate

import (
    "time"
    
    "github.com/go-kratos/kratos/v2/log"
    "medical-system/user/internal/domain/entity"
    "medical-system/user/internal/domain/event"
    "medical-system/user/internal/domain/vo"
)

// Patient 患者聚合根
type Patient struct {
    // 聚合根ID
    ID int64
    
    // 基础信息
    OpenID      string
    Username    string
    Password    string // 加密存储
    Name        string
    Phone       vo.Phone
    Email       vo.Email
    IDCard      vo.IDCard
    
    // 地址信息（值对象）
    Address vo.Address
    
    // 状态
    Status PatientStatus
    
    // 元数据
    CreatedAt time.Time
    UpdatedAt time.Time
    DeletedAt *time.Time
    
    // 领域事件
    Events []event.DomainEvent
}

// PatientStatus 患者状态
type PatientStatus int32

const (
    PatientStatusNormal   PatientStatus = 1  // 正常
    PatientStatusLocked   PatientStatus = 2  // 锁定
    PatientStatusInactive PatientStatus = 3  // 未激活
)

// 领域方法

// UpdateProfile 更新个人信息
func (p *Patient) UpdateProfile(name string, phone vo.Phone) error {
    if p.Status == PatientStatusLocked {
        return ErrPatientLocked
    }
    
    p.Name = name
    p.Phone = phone
    p.UpdatedAt = time.Now()
    
    // 发布领域事件
    p.Events = append(p.Events, event.NewPatientProfileUpdatedEvent(p.ID))
    
    return nil
}

// ChangePassword 修改密码
func (p *Patient) ChangePassword(oldPwd, newPwd string, encryptFn func(string) string) error {
    if encryptFn(oldPwd) != p.Password {
        return ErrInvalidPassword
    }
    
    p.Password = encryptFn(newPwd)
    p.UpdatedAt = time.Now()
    
    p.Events = append(p.Events, event.NewPatientPasswordChangedEvent(p.ID))
    
    return nil
}

// SoftDelete 软删除
func (p *Patient) SoftDelete() {
    now := time.Now()
    p.DeletedAt = &now
    p.Status = PatientStatusInactive
    p.Events = append(p.Events, event.NewPatientDeletedEvent(p.ID))
}

// IsDeleted 是否已删除
func (p *Patient) IsDeleted() bool {
    return p.DeletedAt != nil
}

var (
    ErrPatientLocked    = errors.New("patient is locked")
    ErrInvalidPassword  = errors.New("invalid password")
)
```

#### 2.3.3 实体设计

```go
// medical/internal/domain/entity/medical_record.go
package entity

import (
    "time"
    
    "medical-system/medical/internal/domain/vo"
)

// MedicalRecord 病历实体（聚合根）
type MedicalRecord struct {
    ID              int64
    MedicalNumber   string          // 病历单号
    PatientID       int64           // 患者ID（跨服务引用）
    DoctorID        int64           // 医生ID（跨服务引用）
    DepartmentID    int64           // 科室ID（跨服务引用）
    
    // 诊断信息
    Symptom       string
    Diagnosis     string
    TreatmentPlan string
    
    // 关联实体
    Diseases []DiseaseItem         // 诊断疾病
    Drugs    []PrescriptionItem    // 处方药品
    
    // 状态
    Status MedicalStatus
    
    // 支付信息（值对象）
    Payment *vo.PaymentInfo
    
    // 时间戳
    BillingTime    *time.Time
    SettlementTime *time.Time
    CreatedAt      time.Time
    UpdatedAt      time.Time
    DeletedAt      *time.Time
}

// MedicalStatus 病历状态
type MedicalStatus int32

const (
    MedicalStatusRegistering MedicalStatus = 1  // 挂号中
    MedicalStatusQueuing     MedicalStatus = 2  // 排队中
    MedicalStatusDiagnosing  MedicalStatus = 3  // 诊断中
    MedicalStatusOperating   MedicalStatus = 4  // 手术中
    MedicalStatusCompleted   MedicalStatus = 5  // 诊断结束
    MedicalStatusCancelled   MedicalStatus = 6  // 已取消
)

// DiseaseItem 疾病项（实体）
type DiseaseItem struct {
    ID          int64
    DiseaseID   int64   // 疾病字典ID
    DiseaseName string  // 冗余存储，避免跨服务查询
    Symptoms    string
    Diagnosis   string
}

// PrescriptionItem 处方项（实体）
type PrescriptionItem struct {
    ID       int64
    DrugID   int64
    DrugName string
    Quantity int32
    Price    vo.Money
    Usage    string
}

// AddDiagnosis 添加诊断
func (m *MedicalRecord) AddDiagnosis(disease DiseaseItem, doctorID int64) error {
    if m.Status != MedicalStatusDiagnosing {
        return ErrInvalidStatus
    }
    
    m.Diseases = append(m.Diseases, disease)
    m.DoctorID = doctorID
    m.UpdatedAt = time.Now()
    
    return nil
}

// AddPrescription 添加处方
func (m *MedicalRecord) AddPrescription(drug PrescriptionItem) error {
    if m.Status != MedicalStatusDiagnosing {
        return ErrInvalidStatus
    }
    
    m.Drugs = append(m.Drugs, drug)
    m.UpdatedAt = time.Now()
    
    return nil
}

// Complete 完成诊断
func (m *MedicalRecord) Complete() error {
    if m.Status != MedicalStatusDiagnosing {
        return ErrInvalidStatus
    }
    
    m.Status = MedicalStatusCompleted
    now := time.Now()
    m.SettlementTime = &now
    m.UpdatedAt = now
    
    return nil
}
```

#### 2.3.4 值对象设计

```go
// common/domain/vo/value_objects.go
package vo

import (
    "errors"
    "regexp"
)

// Phone 手机号值对象
type Phone struct {
    value string
}

func NewPhone(phone string) (Phone, error) {
    if !isValidPhone(phone) {
        return Phone{}, errors.New("invalid phone number")
    }
    return Phone{value: phone}, nil
}

func (p Phone) String() string {
    return p.value
}

func isValidPhone(phone string) bool {
    pattern := `^1[3-9]\d{9}$`
    matched, _ := regexp.MatchString(pattern, phone)
    return matched
}

// IDCard 身份证值对象
type IDCard struct {
    value string
}

func NewIDCard(idCard string) (IDCard, error) {
    if !isValidIDCard(idCard) {
        return IDCard{}, errors.New("invalid ID card")
    }
    return IDCard{value: idCard}, nil
}

func (i IDCard) String() string {
    return i.value
}

func (i IDCard) Masked() string {
    // 脱敏显示：110101********1234
    if len(i.value) != 18 {
        return i.value
    }
    return i.value[:6] + "********" + i.value[14:]
}

// Money 金额值对象
type Money struct {
    amount   int64  // 分为单位，避免浮点数精度问题
    currency string // CNY/USD
}

func NewMoney(yuan float64) Money {
    return Money{
        amount:   int64(yuan * 100),
        currency: "CNY",
    }
}

func (m Money) Yuan() float64 {
    return float64(m.amount) / 100
}

func (m Money) Add(other Money) Money {
    return Money{
        amount:   m.amount + other.amount,
        currency: m.currency,
    }
}

// Address 地址值对象
type Address struct {
    Province string
    City     string
    District string
    Detail   string
}

func (a Address) String() string {
    return a.Province + a.City + a.District + a.Detail
}
```

---

## 三、微服务目录结构

### 3.1 Kratos标准目录结构

```
medical-system/                      # 项目根目录
├── api/                             # API定义（protobuf）
│   ├── patient/                     # 患者服务API
│   │   ├── v1/
│   │   │   ├── patient.proto        # 患者服务proto定义
│   │   │   └── patient.pb.go        # 生成代码
│   │   └── v1.go
│   ├── doctor/                      # 医生服务API
│   ├── department/                  # 科室服务API
│   ├── registration/                # 挂号服务API
│   ├── medical/                     # 病历服务API
│   ├── surgical/                    # 手术服务API
│   ├── drug/                        # 药品服务API
│   ├── payment/                     # 支付服务API
│   ├── ai/                          # AI服务API
│   └── auth/                        # 认证服务API
│
├── app/                             # 应用层（BFF）
│   ├── patient-api/                 # 患者端API
│   │   ├── cmd/
│   │   │   └── main.go
│   │   ├── internal/
│   │   │   ├── service/             # 应用服务（聚合多个领域服务）
│   │   │   ├── handler/             # HTTP处理器
│   │   │   └── router/              # 路由定义
│   │   └── configs/
│   ├── doctor-api/                  # 医生端API
│   └── admin-api/                   # 管理端API
│
├── service/                         # 领域服务（微服务）
│   ├── patient/                     # 患者服务
│   │   ├── cmd/
│   │   │   └── main.go              # 服务入口
│   │   ├── internal/
│   │   │   ├── service/             # 应用服务（实现proto接口）
│   │   │   ├── biz/                 # 业务层（领域逻辑）
│   │   │   │   ├── domain/          # 领域模型
│   │   │   │   │   ├── aggregate/   # 聚合根
│   │   │   │   │   ├── entity/      # 实体
│   │   │   │   │   ├── vo/          # 值对象
│   │   │   │   │   └── event/       # 领域事件
│   │   │   │   └── usecase/         # 用例层
│   │   │   ├── data/                # 数据层（仓储实现）
│   │   │   │   ├── mysql/           # MySQL实现
│   │   │   │   ├── redis/           # Redis实现
│   │   │   │   └── repository.go    # 仓储接口实现
│   │   │   └── server/              # 服务启动配置
│   │   ├── configs/                 # 配置文件
│   │   └── Dockerfile
│   ├── doctor/                      # 医生服务
│   ├── department/                  # 科室服务
│   ├── registration/                # 挂号服务
│   ├── medical/                     # 病历服务
│   ├── surgical/                    # 手术服务
│   ├── drug/                        # 药品服务
│   ├── payment/                     # 支付服务
│   ├── ai/                          # AI服务
│   ├── auth/                        # 认证服务
│   ├── file/                        # 文件服务
│   └── notification/                # 通知服务
│
├── common/                          # 公共库
│   ├── middleware/                  # 中间件
│   │   ├── auth/                    # 认证中间件
│   │   ├── logging/                 # 日志中间件
│   │   ├── tracing/                 # 链路追踪
│   │   └── recovery/                # 异常恢复
│   ├── pkg/                         # 公共包
│   │   ├── crypto/                  # 加密工具
│   │   ├── jwt/                     # JWT工具
│   │   ├── validator/               # 验证器
│   │   └── errors/                  # 错误定义
│   └── proto/                       # 公共proto定义
│
├── deploy/                          # 部署配置
│   ├── docker/                      # Docker配置
│   ├── k8s/                         # Kubernetes配置
│   └── nginx/                       # Nginx配置
│
├── scripts/                         # 脚本
│   ├── generate.sh                  # 代码生成脚本
│   ├── migrate.sh                   # 数据库迁移脚本
│   └── deploy.sh                    # 部署脚本
│
├── Makefile                         # Make配置
├── go.mod                           # Go模块定义
└── README.md                        # 项目说明
```

### 3.2 单服务内部结构

```
patient-service/
├── cmd/
│   └── main.go                      # 服务入口
├── internal/
│   ├── service/                     # 应用服务层
│   │   ├── patient.go               # 患者服务实现
│   │   └── patient_test.go          # 服务测试
│   ├── biz/                         # 业务逻辑层
│   │   ├── domain/                  # 领域模型
│   │   │   ├── aggregate/           # 聚合根
│   │   │   │   └── patient.go
│   │   │   ├── entity/              # 实体
│   │   │   ├── vo/                  # 值对象
│   │   │   │   ├── phone.go
│   │   │   │   ├── address.go
│   │   │   │   └── idcard.go
│   │   │   ├── event/               # 领域事件
│   │   │   │   └── patient_events.go
│   │   │   └── repository/          # 仓储接口
│   │   │       └── patient_repo.go
│   │   ├── usecase/                 # 用例层
│   │   │   ├── register_usecase.go  # 注册用例
│   │   │   └── profile_usecase.go   # 个人信息用例
│   │   └── interceptor/             # 业务拦截器
│   ├── data/                        # 数据层
│   │   ├── mysql/                   # MySQL实现
│   │   │   ├── patient.go           # 仓储实现
│   │   │   └── model.go             # 数据库模型
│   │   ├── redis/                   # Redis实现
│   │   │   └── cache.go
│   │   └── ent/                     # 可选：Ent ORM
│   └── server/                      # 服务配置
│       ├── grpc.go                  # gRPC服务
│       ├── http.go                  # HTTP服务
│       └── wire.go                  # 依赖注入
├── configs/                         # 配置文件
│   ├── config.yaml                  # 主配置
│   └── config_test.yaml             # 测试配置
└── Dockerfile                       # 容器镜像
```

---

## 四、强化权限系统设计

### 4.1 整体架构

```
┌─────────────────────────────────────────────────────────────────┐
│                      权限系统架构                                │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│  ┌─────────────────────────────────────────────────────────┐   │
│  │                      认证层                              │   │
│  │   ┌──────────┐    ┌──────────┐    ┌──────────┐        │   │
│  │   │ OAuth2.0 │    │   JWT    │    │  Refresh │        │   │
│  │   │ 微信登录 │    │ 令牌颁发 │    │ 令牌刷新 │        │   │
│  │   └──────────┘    └──────────┘    └──────────┘        │   │
│  └─────────────────────────────────────────────────────────┘   │
│                              │                                   │
│                              ▼                                   │
│  ┌─────────────────────────────────────────────────────────┐   │
│  │                      授权层                              │   │
│  │   ┌──────────┐    ┌──────────┐    ┌──────────┐        │   │
│  │   │   RBAC   │    │ 数据权限 │    │ API权限  │        │   │
│  │   │ 角色控制 │    │ 科室隔离 │    │ 接口控制 │        │   │
│  │   └──────────┘    └──────────┘    └──────────┘        │   │
│  └─────────────────────────────────────────────────────────┘   │
│                              │                                   │
│                              ▼                                   │
│  ┌─────────────────────────────────────────────────────────┐   │
│  │                      审计层                              │   │
│  │   ┌──────────┐    ┌──────────┐    ┌──────────┐        │   │
│  │   │操作日志  │    │安全审计  │    │异常告警  │        │   │
│  │   │          │    │          │    │          │        │   │
│  │   └──────────┘    └──────────┘    └──────────┘        │   │
│  └─────────────────────────────────────────────────────────┘   │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
```

### 4.2 数据库设计

```sql
-- 用户表（统一认证）
CREATE TABLE sys_user (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    user_type TINYINT NOT NULL COMMENT '用户类型:1-患者 2-医生 3-管理员',
    username VARCHAR(50) NOT NULL,
    password VARCHAR(255) NOT NULL COMMENT '加密密码',
    real_name VARCHAR(50),
    phone VARCHAR(20),
    email VARCHAR(100),
    avatar VARCHAR(255),
    status TINYINT DEFAULT 1 COMMENT '状态:0-禁用 1-启用',
    last_login_at DATETIME,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at DATETIME DEFAULT NULL,
    UNIQUE KEY uk_username (username),
    UNIQUE KEY uk_phone (phone),
    INDEX idx_status (status),
    INDEX idx_deleted (deleted_at)
) ENGINE=InnoDB COMMENT='系统用户表';

-- 角色表
CREATE TABLE sys_role (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    role_code VARCHAR(50) NOT NULL UNIQUE COMMENT '角色编码',
    role_name VARCHAR(100) NOT NULL COMMENT '角色名称',
    role_type TINYINT DEFAULT 1 COMMENT '角色类型:1-系统 2-自定义',
    data_scope TINYINT DEFAULT 4 COMMENT '数据范围:1-全部 2-本部门 3-本部门及下级 4-仅本人',
    description VARCHAR(255),
    status TINYINT DEFAULT 1,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at DATETIME DEFAULT NULL,
    INDEX idx_code (role_code),
    INDEX idx_type (role_type)
) ENGINE=InnoDB COMMENT='角色表';

-- 权限表（资源）
CREATE TABLE sys_permission (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    perm_code VARCHAR(100) NOT NULL UNIQUE COMMENT '权限编码',
    perm_name VARCHAR(100) NOT NULL COMMENT '权限名称',
    resource_type TINYINT NOT NULL COMMENT '资源类型:1-菜单 2-按钮 3-API',
    parent_id BIGINT DEFAULT 0 COMMENT '父权限ID',
    resource_url VARCHAR(255) COMMENT '资源URL/接口',
    http_method VARCHAR(20) COMMENT 'HTTP方法',
    icon VARCHAR(100) COMMENT '图标',
    sort_order INT DEFAULT 0,
    status TINYINT DEFAULT 1,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_parent (parent_id),
    INDEX idx_type (resource_type)
) ENGINE=InnoDB COMMENT='权限表';

-- 用户-角色关联表（多对多）
CREATE TABLE sys_user_role (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    user_id BIGINT NOT NULL,
    role_id BIGINT NOT NULL,
    dept_id BIGINT COMMENT '所属科室（科室角色需要）',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    UNIQUE KEY uk_user_role_dept (user_id, role_id, dept_id),
    FOREIGN KEY (user_id) REFERENCES sys_user(id),
    FOREIGN KEY (role_id) REFERENCES sys_role(id)
) ENGINE=InnoDB COMMENT='用户角色关联表';

-- 角色-权限关联表（多对多）
CREATE TABLE sys_role_permission (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    role_id BIGINT NOT NULL,
    permission_id BIGINT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    UNIQUE KEY uk_role_perm (role_id, permission_id),
    FOREIGN KEY (role_id) REFERENCES sys_role(id),
    FOREIGN KEY (permission_id) REFERENCES sys_permission(id)
) ENGINE=InnoDB COMMENT='角色权限关联表';

-- 数据权限规则表（精细化控制）
CREATE TABLE sys_data_permission (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    user_id BIGINT NOT NULL,
    resource_type VARCHAR(50) NOT NULL COMMENT '资源类型:medical/surgical等',
    data_scope TINYINT NOT NULL COMMENT '数据范围',
    dept_ids JSON COMMENT '可访问科室ID列表',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE KEY uk_user_resource (user_id, resource_type),
    FOREIGN KEY (user_id) REFERENCES sys_user(id)
) ENGINE=InnoDB COMMENT='数据权限规则表';

-- 操作日志表
CREATE TABLE sys_operation_log (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    user_id BIGINT,
    user_type TINYINT,
    username VARCHAR(50),
    operation VARCHAR(100) COMMENT '操作描述',
    method VARCHAR(255) COMMENT '请求方法',
    request_url VARCHAR(500) COMMENT '请求URL',
    request_method VARCHAR(20) COMMENT 'HTTP方法',
    request_params TEXT COMMENT '请求参数',
    response_data TEXT COMMENT '响应数据',
    ip VARCHAR(50) COMMENT 'IP地址',
    user_agent VARCHAR(500) COMMENT 'User-Agent',
    status TINYINT COMMENT '状态:0-失败 1-成功',
    error_msg TEXT COMMENT '错误信息',
    duration_ms INT COMMENT '执行时长(ms)',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_user (user_id),
    INDEX idx_time (created_at),
    INDEX idx_operation (operation)
) ENGINE=InnoDB COMMENT='操作日志表';

-- 登录日志表
CREATE TABLE sys_login_log (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    user_id BIGINT,
    username VARCHAR(50),
    login_type TINYINT COMMENT '登录类型:1-账号 2-微信 3-短信',
    ip VARCHAR(50),
    user_agent VARCHAR(500),
    status TINYINT COMMENT '状态:0-失败 1-成功',
    fail_reason VARCHAR(255),
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_user (user_id),
    INDEX idx_time (created_at)
) ENGINE=InnoDB COMMENT='登录日志表';
```

### 4.3 Go实现

```go
// auth/internal/domain/aggregate/user.go
package aggregate

import (
    "time"
    
    "medical-system/auth/internal/domain/vo"
)

// User 用户聚合根
type User struct {
    ID          int64
    UserType    UserType
    Username    string
    Password    vo.Password
    RealName    string
    Phone       vo.Phone
    Email       vo.Email
    Avatar      string
    Status      UserStatus
    Roles       []Role          // 角色集合
    DataPerms   []DataPermission // 数据权限
    LastLoginAt *time.Time
    CreatedAt   time.Time
    UpdatedAt   time.Time
    DeletedAt   *time.Time
}

type UserType int32

const (
    UserTypePatient UserType = 1
    UserTypeDoctor  UserType = 2
    UserTypeAdmin   UserType = 3
)

type UserStatus int32

const (
    UserStatusDisabled UserStatus = 0
    UserStatusEnabled  UserStatus = 1
    UserStatusLocked   UserStatus = 2
)

// HasPermission 检查是否有权限
func (u *User) HasPermission(permCode string) bool {
    for _, role := range u.Roles {
        for _, perm := range role.Permissions {
            if perm.Code == permCode {
                return true
            }
        }
    }
    return false
}

// HasRole 检查是否有角色
func (u *User) HasRole(roleCode string) bool {
    for _, role := range u.Roles {
        if role.Code == roleCode {
            return true
        }
    }
    return false
}

// CanAccessData 检查数据权限
func (u *User) CanAccessData(resourceType string, deptID int64) bool {
    // 超级管理员
    if u.HasRole("SUPER_ADMIN") {
        return true
    }
    
    // 检查数据权限规则
    for _, dp := range u.DataPerms {
        if dp.ResourceType == resourceType {
            switch dp.Scope {
            case DataScopeAll:
                return true
            case DataScopeDept:
                return dp.DeptID == deptID
            case DataScopeDeptAndSub:
                return dp.ContainsDept(deptID)
            case DataScopeOwn:
                return dp.UserID == u.ID
            }
        }
    }
    
    return false
}

// Role 角色实体
type Role struct {
    ID          int64
    Code        string
    Name        string
    Type        int32
    DataScope   DataScopeType
    Description string
    Permissions []Permission
}

// Permission 权限实体
type Permission struct {
    ID           int64
    Code         string
    Name         string
    ResourceType int32
    ResourceURL  string
    HTTPMethod   string
}

// DataPermission 数据权限
type DataPermission struct {
    ID           int64
    UserID       int64
    ResourceType string
    Scope        DataScopeType
    DeptID       int64
    DeptIDs      []int64
}

type DataScopeType int32

const (
    DataScopeAll        DataScopeType = 1
    DataScopeDept       DataScopeType = 2
    DataScopeDeptAndSub DataScopeType = 3
    DataScopeOwn        DataScopeType = 4
)

func (dp *DataPermission) ContainsDept(deptID int64) bool {
    if dp.DeptID == deptID {
        return true
    }
    for _, id := range dp.DeptIDs {
        if id == deptID {
            return true
        }
    }
    return false
}
```

### 4.4 认证中间件

```go
// common/middleware/auth/auth.go
package auth

import (
    "context"
    "strings"
    
    "github.com/go-kratos/kratos/v2/errors"
    "github.com/go-kratos/kratos/v2/metadata"
    "github.com/go-kratos/kratos/v2/middleware"
    "github.com/go-kratos/kratos/v2/transport"
    
    "medical-system/common/pkg/jwt"
)

// 上下文Key
const (
    ContextKeyUserID   = "user_id"
    ContextKeyUserType = "user_type"
    ContextKeyRoles    = "roles"
    ContextKeyToken    = "token"
)

// UserClaims 用户声明
type UserClaims struct {
    UserID   int64    `json:"user_id"`
    UserType int32    `json:"user_type"`
    Username string   `json:"username"`
    Roles    []string `json:"roles"`
}

// AuthMiddleware 认证中间件
func AuthMiddleware(verifier jwt.Verifier) middleware.Middleware {
    return func(handler middleware.Handler) middleware.Handler {
        return func(ctx context.Context, req interface{}) (interface{}, error) {
            // 从metadata获取token
            md, ok := metadata.FromServerContext(ctx)
            if !ok {
                return nil, errors.Unauthorized("UNAUTHORIZED", "missing metadata")
            }
            
            token := md.Get("Authorization")
            if token == "" {
                return nil, errors.Unauthorized("UNAUTHORIZED", "missing token")
            }
            
            // Bearer token
            if strings.HasPrefix(token, "Bearer ") {
                token = token[7:]
            }
            
            // 验证token
            claims, err := verifier.Verify(token)
            if err != nil {
                return nil, errors.Unauthorized("UNAUTHORIZED", err.Error())
            }
            
            // 注入上下文
            ctx = context.WithValue(ctx, ContextKeyUserID, claims.UserID)
            ctx = context.WithValue(ctx, ContextKeyUserType, claims.UserType)
            ctx = context.WithValue(ctx, ContextKeyRoles, claims.Roles)
            ctx = context.WithValue(ctx, ContextKeyToken, token)
            
            return handler(ctx, req)
        }
    }
}

// PermissionMiddleware 权限中间件
func PermissionMiddleware(requiredPerms ...string) middleware.Middleware {
    return func(handler middleware.Handler) middleware.Handler {
        return func(ctx context.Context, req interface{}) (interface{}, error) {
            // 获取用户角色
            roles, ok := ctx.Value(ContextKeyRoles).([]string)
            if !ok {
                return nil, errors.Forbidden("FORBIDDEN", "no permission")
            }
            
            // 检查权限
            hasPerm := false
            for _, role := range roles {
                if role == "SUPER_ADMIN" {
                    hasPerm = true
                    break
                }
                // 其他权限检查逻辑...
            }
            
            if !hasPerm {
                return nil, errors.Forbidden("FORBIDDEN", "insufficient permission")
            }
            
            return handler(ctx, req)
        }
    }
}

// GetUserID 从上下文获取用户ID
func GetUserID(ctx context.Context) int64 {
    id, _ := ctx.Value(ContextKeyUserID).(int64)
    return id
}

// GetUserType 从上下文获取用户类型
func GetUserType(ctx context.Context) int32 {
    t, _ := ctx.Value(ContextKeyUserType).(int32)
    return t
}
```

---

## 五、Proto定义示例

### 5.1 患者服务定义

```protobuf
// api/patient/v1/patient.proto
syntax = "proto3";

package api.patient.v1;

option go_package = "medical-system/api/patient/v1;v1";
option java_multiple_files = true;
option java_package = "api.patient.v1";

import "google/protobuf/timestamp.proto";
import "google/api/annotations.proto";
import "validate/validate.proto";

// 患者服务
service PatientService {
    // 患者注册
    rpc RegisterPatient(RegisterPatientRequest) returns (RegisterPatientResponse) {
        option (google.api.http) = {
            post: "/api/v1/patient/register"
            body: "*"
        };
    }
    
    // 患者登录
    rpc LoginPatient(LoginPatientRequest) returns (LoginPatientResponse) {
        option (google.api.http) = {
            post: "/api/v1/patient/login"
            body: "*"
        };
    }
    
    // 获取患者信息
    rpc GetPatient(GetPatientRequest) returns (Patient) {
        option (google.api.http) = {
            get: "/api/v1/patient/{id}"
        };
    }
    
    // 更新患者信息
    rpc UpdatePatient(UpdatePatientRequest) returns (Patient) {
        option (google.api.http) = {
            put: "/api/v1/patient/{id}"
            body: "*"
        };
    }
    
    // 患者列表（分页）
    rpc ListPatients(ListPatientsRequest) returns (ListPatientsResponse) {
        option (google.api.http) = {
            post: "/api/v1/patients"
            body: "*"
        };
    }
    
    // 删除患者（软删除）
    rpc DeletePatient(DeletePatientRequest) returns (DeletePatientResponse) {
        option (google.api.http) = {
            delete: "/api/v1/patient/{id}"
        };
    }
}

// 患者信息
message Patient {
    int64 id = 1;
    string openid = 2;
    string username = 3;
    string name = 4 [(validate.rules).string = {min_len: 2, max_len: 50}];
    string phone = 5 [(validate.rules).string = {pattern: "^1[3-9]\\d{9}$"}];
    string email = 6;
    string id_card = 7;
    string avatar = 8;
    Address address = 9;
    int32 status = 10;
    google.protobuf.Timestamp created_at = 11;
    google.protobuf.Timestamp updated_at = 12;
}

// 地址
message Address {
    string province = 1;
    string city = 2;
    string district = 3;
    string detail = 4;
}

// 注册请求
message RegisterPatientRequest {
    string phone = 1 [(validate.rules).string = {pattern: "^1[3-9]\\d{9}$"}];
    string password = 2 [(validate.rules).string = {min_len: 6, max_len: 20}];
    string name = 3 [(validate.rules).string = {min_len: 2, max_len: 50}];
    string code = 4 [(validate.rules).string = {len: 6}]; // 验证码
}

message RegisterPatientResponse {
    Patient patient = 1;
    string token = 2;
}

// 登录请求
message LoginPatientRequest {
    string username = 1;
    string password = 2;
}

message LoginPatientResponse {
    Patient patient = 1;
    string token = 2;
    int64 expires_in = 3; // 过期时间（秒）
}

// 获取患者请求
message GetPatientRequest {
    int64 id = 1 [(validate.rules).int64 = {gt: 0}];
}

// 更新患者请求
message UpdatePatientRequest {
    int64 id = 1 [(validate.rules).int64 = {gt: 0}];
    string name = 2;
    string phone = 3;
    Address address = 4;
}

// 患者列表请求
message ListPatientsRequest {
    int32 page = 1 [(validate.rules).int32 = {gt: 0}];
    int32 page_size = 2 [(validate.rules).int32 = {gt: 0, lte: 100}];
    string keyword = 3; // 搜索关键词
    int32 status = 4;   // 状态筛选
}

message ListPatientsResponse {
    int64 total = 1;
    repeated Patient patients = 2;
}

// 删除患者请求
message DeletePatientRequest {
    int64 id = 1 [(validate.rules).int64 = {gt: 0}];
}

message DeletePatientResponse {
    bool success = 1;
}
```

### 5.2 病历服务定义

```protobuf
// api/medical/v1/medical.proto
syntax = "proto3";

package api.medical.v1;

option go_package = "medical-system/api/medical/v1;v1";

import "google/protobuf/timestamp.proto";
import "google/api/annotations.proto";
import "validate/validate.proto";

// 病历服务
service MedicalService {
    // 创建病历
    rpc CreateMedicalRecord(CreateMedicalRecordRequest) returns (MedicalRecord) {
        option (google.api.http) = {
            post: "/api/v1/medical"
            body: "*"
        };
    }
    
    // 获取病历
    rpc GetMedicalRecord(GetMedicalRecordRequest) returns (MedicalRecord) {
        option (google.api.http) = {
            get: "/api/v1/medical/{id}"
        };
    }
    
    // 更新病历
    rpc UpdateMedicalRecord(UpdateMedicalRecordRequest) returns (MedicalRecord) {
        option (google.api.http) = {
            put: "/api/v1/medical/{id}"
            body: "*"
        };
    }
    
    // 病历列表（支持数据权限）
    rpc ListMedicalRecords(ListMedicalRecordsRequest) returns (ListMedicalRecordsResponse) {
        option (google.api.http) = {
            post: "/api/v1/medicals"
            body: "*"
        };
    }
    
    // 添加诊断
    rpc AddDiagnosis(AddDiagnosisRequest) returns (MedicalRecord) {
        option (google.api.http) = {
            post: "/api/v1/medical/{id}/diagnosis"
            body: "*"
        };
    }
    
    // 添加处方
    rpc AddPrescription(AddPrescriptionRequest) returns (MedicalRecord) {
        option (google.api.http) = {
            post: "/api/v1/medical/{id}/prescription"
            body: "*"
        };
    }
    
    // 完成诊断
    rpc CompleteDiagnosis(CompleteDiagnosisRequest) returns (MedicalRecord) {
        option (google.api.http) = {
            post: "/api/v1/medical/{id}/complete"
        };
    }
}

// 病历记录
message MedicalRecord {
    int64 id = 1;
    string medical_number = 2;
    int64 patient_id = 3;
    int64 doctor_id = 4;
    int64 department_id = 5;
    string symptom = 6;
    string diagnosis = 7;
    string treatment_plan = 8;
    repeated DiagnosisItem diagnoses = 9;
    repeated PrescriptionItem prescriptions = 10;
    PaymentInfo payment = 11;
    int32 status = 12; // 1-挂号中 2-排队中 3-诊断中 4-手术中 5-诊断结束 6-已取消
    google.protobuf.Timestamp billing_time = 13;
    google.protobuf.Timestamp settlement_time = 14;
    google.protobuf.Timestamp created_at = 15;
    google.protobuf.Timestamp updated_at = 16;
}

// 诊断项
message DiagnosisItem {
    int64 id = 1;
    int64 disease_id = 2;
    string disease_name = 3;
    string symptoms = 4;
    string diagnosis = 5;
}

// 处方项
message PrescriptionItem {
    int64 id = 1;
    int64 drug_id = 2;
    string drug_name = 3;
    int32 quantity = 4;
    double price = 5;
    string usage = 6;
}

// 支付信息
message PaymentInfo {
    int32 pay_method = 1; // 1-微信 2-支付宝
    int32 pay_status = 2; // 1-未支付 2-已支付 3-已退款
    double amount = 3;
    google.protobuf.Timestamp pay_time = 4;
}

// 创建病历请求
message CreateMedicalRecordRequest {
    int64 patient_id = 1 [(validate.rules).int64 = {gt: 0}];
    int64 doctor_id = 2 [(validate.rules).int64 = {gt: 0}];
    int64 department_id = 3 [(validate.rules).int64 = {gt: 0}];
    string symptom = 4;
}

// 获取病历请求
message GetMedicalRecordRequest {
    int64 id = 1 [(validate.rules).int64 = {gt: 0}];
}

// 更新病历请求
message UpdateMedicalRecordRequest {
    int64 id = 1 [(validate.rules).int64 = {gt: 0}];
    string symptom = 2;
    string diagnosis = 3;
    string treatment_plan = 4;
}

// 病历列表请求
message ListMedicalRecordsRequest {
    int32 page = 1 [(validate.rules).int32 = {gt: 0}];
    int32 page_size = 2 [(validate.rules).int32 = {gt: 0, lte: 100}];
    int64 patient_id = 3;   // 筛选患者
    int64 doctor_id = 4;    // 筛选医生
    int64 department_id = 5; // 筛选科室
    int32 status = 6;       // 筛选状态
    string start_date = 7;  // 开始日期 YYYY-MM-DD
    string end_date = 8;    // 结束日期 YYYY-MM-DD
}

message ListMedicalRecordsResponse {
    int64 total = 1;
    repeated MedicalRecord records = 2;
}

// 添加诊断请求
message AddDiagnosisRequest {
    int64 id = 1 [(validate.rules).int64 = {gt: 0}];
    int64 disease_id = 2 [(validate.rules).int64 = {gt: 0}];
    string symptoms = 3;
    string diagnosis = 4;
}

// 添加处方请求
message AddPrescriptionRequest {
    int64 id = 1 [(validate.rules).int64 = {gt: 0}];
    repeated PrescriptionItem items = 2 [(validate.rules).repeated = {min_items: 1}];
}

// 完成诊断请求
message CompleteDiagnosisRequest {
    int64 id = 1 [(validate.rules).int64 = {gt: 0}];
}
```

---

## 六、服务间通信设计

### 6.1 同步通信（gRPC）

```go
// 服务调用示例

// registration/internal/biz/usecase/registration_usecase.go
package usecase

import (
    "context"
    
    "github.com/go-kratos/kratos/v2/log"
    "github.com/go-kratos/kratos/v2/transport/grpc"
    
    patientv1 "medical-system/api/patient/v1"
    doctorv1 "medical-system/api/doctor/v1"
    medicalv1 "medical-system/api/medical/v1"
)

// RegistrationUsecase 挂号用例
type RegistrationUsecase struct {
    patientClient patientv1.PatientServiceClient
    doctorClient  doctorv1.DoctorServiceClient
    medicalClient medicalv1.MedicalServiceClient
    repo          RegistrationRepo
    log           *log.Helper
}

// NewRegistrationUsecase 创建用例
func NewRegistrationUsecase(
    patientConn grpc.ClientConn,
    doctorConn grpc.ClientConn,
    medicalConn grpc.ClientConn,
    repo RegistrationRepo,
    logger log.Logger,
) *RegistrationUsecase {
    return &RegistrationUsecase{
        patientClient: patientv1.NewPatientServiceClient(patientConn),
        doctorClient:  doctorv1.NewDoctorServiceClient(doctorConn),
        medicalClient: medicalv1.NewMedicalServiceClient(medicalConn),
        repo:          repo,
        log:           log.NewHelper(logger),
    }
}

// CreateRegistration 创建挂号
func (uc *RegistrationUsecase) CreateRegistration(ctx context.Context, req *CreateRegistrationRequest) (*Registration, error) {
    // 1. 验证患者信息（跨服务调用）
    patient, err := uc.patientClient.GetPatient(ctx, &patientv1.GetPatientRequest{
        Id: req.PatientID,
    })
    if err != nil {
        return nil, err
    }
    
    // 2. 验证医生信息（跨服务调用）
    doctor, err := uc.doctorClient.GetDoctor(ctx, &doctorv1.GetDoctorRequest{
        Id: req.DoctorID,
    })
    if err != nil {
        return nil, err
    }
    
    // 3. 检查医生是否可预约
    if doctor.Status != doctorv1.DoStatusWorking {
        return nil, ErrDoctorNotAvailable
    }
    
    // 4. 创建病历（跨服务调用）
    medical, err := uc.medicalClient.CreateMedicalRecord(ctx, &medicalv1.CreateMedicalRecordRequest{
        PatientId:    req.PatientID,
        DoctorId:     req.DoctorID,
        DepartmentId: doctor.DepartmentId,
        Symptom:      req.Symptom,
    })
    if err != nil {
        return nil, err
    }
    
    // 5. 创建挂号记录
    registration := &Registration{
        PatientID:    req.PatientID,
        DoctorID:     req.DoctorID,
        DepartmentID: doctor.DepartmentId,
        MedicalID:    medical.Id,
        RegisterDate: req.RegisterDate,
        Status:       RegistrationStatusPending,
    }
    
    if err := uc.repo.Create(ctx, registration); err != nil {
        return nil, err
    }
    
    return registration, nil
}
```

### 6.2 异步通信（消息队列）

```go
// 领域事件发布

// medical/internal/domain/event/medical_events.go
package event

import (
    "time"
    
    "github.com/google/uuid"
)

// MedicalCompletedEvent 病历完成事件
type MedicalCompletedEvent struct {
    EventID      string    `json:"event_id"`
    EventType    string    `json:"event_type"`
    EventTime    time.Time `json:"event_time"`
    MedicalID    int64     `json:"medical_id"`
    PatientID    int64     `json:"patient_id"`
    DoctorID     int64     `json:"doctor_id"`
    DepartmentID int64     `json:"department_id"`
    Amount       float64   `json:"amount"`
}

// NewMedicalCompletedEvent 创建事件
func NewMedicalCompletedEvent(medicalID, patientID, doctorID, deptID int64, amount float64) *MedicalCompletedEvent {
    return &MedicalCompletedEvent{
        EventID:      uuid.New().String(),
        EventType:    "medical.completed",
        EventTime:    time.Now(),
        MedicalID:    medicalID,
        PatientID:    patientID,
        DoctorID:     doctorID,
        DepartmentID: deptID,
        Amount:       amount,
    }
}

// medical/internal/service/medical.go
package service

import (
    "encoding/json"
    
    "github.com/go-kratos/kratos/v2/log"
    
    "medical-system/medical/internal/domain/event"
    "medical-system/medical/internal/pkg/mq"
)

// CompleteMedical 完成诊断
func (s *MedicalService) CompleteMedical(ctx context.Context, req *v1.CompleteDiagnosisRequest) (*v1.MedicalRecord, error) {
    // ... 业务逻辑
    
    // 发布领域事件
    evt := event.NewMedicalCompletedEvent(
        medical.ID,
        medical.PatientID,
        medical.DoctorID,
        medical.DepartmentID,
        medical.Payment.Amount,
    )
    
    evtJSON, _ := json.Marshal(evt)
    if err := s.mq.Publish(ctx, "medical.completed", evtJSON); err != nil {
        s.log.Errorf("发布事件失败: %v", err)
        // 失败不阻塞主流程，记录日志即可
    }
    
    return convertToProto(medical), nil
}
```

```go
// payment/internal/service/event_handler.go
package service

import (
    "context"
    "encoding/json"
    
    "medical-system/medical/event"
)

// MedicalEventHandler 病历事件处理器
type MedicalEventHandler struct {
    paymentUsecase *biz.PaymentUsecase
}

// HandleMedicalCompleted 处理病历完成事件
func (h *MedicalEventHandler) HandleMedicalCompleted(ctx context.Context, msg []byte) error {
    var evt event.MedicalCompletedEvent
    if err := json.Unmarshal(msg, &evt); err != nil {
        return err
    }
    
    // 创建支付订单
    return h.paymentUsecase.CreatePayment(ctx, &biz.CreatePaymentRequest{
        MedicalID:    evt.MedicalID,
        PatientID:    evt.PatientID,
        Amount:       evt.Amount,
        Description:  "诊疗费用",
    })
}
```

### 6.3 Saga分布式事务

```go
// registration/internal/biz/saga/registration_saga.go
package saga

import (
    "context"
    "fmt"
    
    "medical-system/common/pkg/saga"
)

// RegistrationSaga 挂号Saga事务
type RegistrationSaga struct {
    patientClient PatientClient
    doctorClient  DoctorClient
    medicalClient MedicalClient
    paymentClient PaymentClient
}

// Execute 执行挂号流程
func (s *RegistrationSaga) Execute(ctx context.Context, req RegistrationRequest) error {
    saga := saga.NewSaga("registration").
        // Step1: 创建病历
        AddStep(
            func() error {
                _, err := s.medicalClient.CreateMedical(ctx, req)
                return err
            },
            func() error {
                // 补偿：删除病历
                return s.medicalClient.DeleteMedical(ctx, req.MedicalID)
            },
        ).
        // Step2: 创建支付订单
        AddStep(
            func() error {
                _, err := s.paymentClient.CreateOrder(ctx, req)
                return err
            },
            func() error {
                // 补偿：取消订单
                return s.paymentClient.CancelOrder(ctx, req.OrderID)
            },
        ).
        // Step3: 更新挂号状态
        AddStep(
            func() error {
                return s.registrationRepo.UpdateStatus(ctx, req.ID, StatusSuccess)
            },
            func() error {
                // 补偿：回滚状态
                return s.registrationRepo.UpdateStatus(ctx, req.ID, StatusFailed)
            },
        )
    
    return saga.Execute(ctx)
}
```

---

## 七、部署架构

### 7.1 Docker Compose本地部署

```yaml
# deploy/docker-compose.yml
version: '3.8'

services:
  # MySQL
  mysql:
    image: mysql:8.0
    container_name: medical-mysql
    environment:
      MYSQL_ROOT_PASSWORD: root123
      MYSQL_DATABASE: medical_system
    ports:
      - "3306:3306"
    volumes:
      - mysql_data:/var/lib/mysql
      - ./init.sql:/docker-entrypoint-initdb.d/init.sql
    networks:
      - medical-network

  # Redis
  redis:
    image: redis:7-alpine
    container_name: medical-redis
    ports:
      - "6379:6379"
    networks:
      - medical-network

  # Consul（服务注册）
  consul:
    image: consul:1.15
    container_name: medical-consul
    ports:
      - "8500:8500"
    command: consul agent -dev -ui -client=0.0.0.0
    networks:
      - medical-network

  # RabbitMQ
  rabbitmq:
    image: rabbitmq:3-management
    container_name: medical-rabbitmq
    ports:
      - "5672:5672"
      - "15672:15672"
    environment:
      RABBITMQ_DEFAULT_USER: admin
      RABBITMQ_DEFAULT_PASS: admin123
    networks:
      - medical-network

  # 患者服务
  patient-service:
    build:
      context: ../service/patient
      dockerfile: Dockerfile
    container_name: medical-patient-service
    environment:
      - SERVER_GRPC_ADDR=0.0.0.0:9001
      - SERVER_HTTP_ADDR=0.0.0.0:8001
      - DATA_DATABASE_DRIVER=mysql
      - DATA_DATABASE_SOURCE=root:root123@tcp(mysql:3306)/medical_system?parseTime=True
      - DATA_REDIS_ADDR=redis:6379
      - REGISTRY_CONSUL_ADDR=consul:8500
    ports:
      - "8001:8001"
      - "9001:9001"
    depends_on:
      - mysql
      - redis
      - consul
    networks:
      - medical-network

  # 医生服务
  doctor-service:
    build:
      context: ../service/doctor
      dockerfile: Dockerfile
    container_name: medical-doctor-service
    environment:
      - SERVER_GRPC_ADDR=0.0.0.0:9002
      - SERVER_HTTP_ADDR=0.0.0.0:8002
    ports:
      - "8002:8002"
      - "9002:9002"
    depends_on:
      - mysql
      - consul
    networks:
      - medical-network

  # 病历服务
  medical-service:
    build:
      context: ../service/medical
      dockerfile: Dockerfile
    container_name: medical-medical-service
    environment:
      - SERVER_GRPC_ADDR=0.0.0.0:9003
      - SERVER_HTTP_ADDR=0.0.0.0:8003
    ports:
      - "8003:8003"
      - "9003:9003"
    depends_on:
      - mysql
      - rabbitmq
      - consul
    networks:
      - medical-network

  # API Gateway
  gateway:
    image: go-kratos/gateway
    container_name: medical-gateway
    volumes:
      - ./gateway.yaml:/etc/gateway.yaml
    ports:
      - "8080:8080"
    command: ["-conf", "/etc/gateway.yaml"]
    depends_on:
      - patient-service
      - doctor-service
      - medical-service
    networks:
      - medical-network

volumes:
  mysql_data:

networks:
  medical-network:
    driver: bridge
```

### 7.2 Kubernetes部署

```yaml
# deploy/k8s/patient-service.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: patient-service
  namespace: medical-system
  labels:
    app: patient-service
spec:
  replicas: 3
  selector:
    matchLabels:
      app: patient-service
  template:
    metadata:
      labels:
        app: patient-service
    spec:
      containers:
      - name: patient-service
        image: registry.cn-beijing.aliyuncs.com/medical/patient-service:v1.0.0
        ports:
        - containerPort: 8001
          name: http
        - containerPort: 9001
          name: grpc
        env:
        - name: SERVER_GRPC_ADDR
          value: "0.0.0.0:9001"
        - name: SERVER_HTTP_ADDR
          value: "0.0.0.0:8001"
        - name: DATA_DATABASE_SOURCE
          valueFrom:
            secretKeyRef:
              name: db-secret
              key: dsn
        resources:
          requests:
            memory: "256Mi"
            cpu: "250m"
          limits:
            memory: "512Mi"
            cpu: "500m"
        livenessProbe:
          httpGet:
            path: /health
            port: 8001
          initialDelaySeconds: 30
          periodSeconds: 10
        readinessProbe:
          httpGet:
            path: /ready
            port: 8001
          initialDelaySeconds: 5
          periodSeconds: 5

---
apiVersion: v1
kind: Service
metadata:
  name: patient-service
  namespace: medical-system
  labels:
    app: patient-service
spec:
  type: ClusterIP
  ports:
  - port: 8001
    name: http
    targetPort: 8001
  - port: 9001
    name: grpc
    targetPort: 9001
  selector:
    app: patient-service

---
apiVersion: autoscaling/v2
kind: HorizontalPodAutoscaler
metadata:
  name: patient-service-hpa
  namespace: medical-system
spec:
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: patient-service
  minReplicas: 3
  maxReplicas: 10
  metrics:
  - type: Resource
    resource:
      name: cpu
      target:
        type: Utilization
        averageUtilization: 70
```

---

## 八、监控与链路追踪

### 8.1 Prometheus指标

```go
// 自定义指标
package metrics

import (
    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promauto"
)

var (
    // HTTP请求计数
    HTTPRequestsTotal = promauto.NewCounterVec(
        prometheus.CounterOpts{
            Name: "http_requests_total",
            Help: "Total number of HTTP requests",
        },
        []string{"service", "method", "endpoint", "status"},
    )
    
    // HTTP请求耗时
    HTTPRequestDuration = promauto.NewHistogramVec(
        prometheus.HistogramOpts{
            Name:    "http_request_duration_seconds",
            Help:    "HTTP request duration in seconds",
            Buckets: prometheus.DefBuckets,
        },
        []string{"service", "endpoint"},
    )
    
    // 业务指标
    RegistrationTotal = promauto.NewCounterVec(
        prometheus.CounterOpts{
            Name: "registration_total",
            Help: "Total number of registrations",
        },
        []string{"status"},
    )
)
```

### 8.2 Jaeger链路追踪

```go
// 链路追踪配置
package tracing

import (
    "github.com/go-kratos/kratos/v2/middleware/tracing"
    "go.opentelemetry.io/otel/trace"
)

// NewTracer 创建链路追踪器
func NewTracer(serviceName string) (trace.Tracer, error) {
    // 配置Jaeger导出器
    // ...
}

// 在服务中使用
httpSrv := http.NewServer(
    http.Middleware(
        tracing.Server(),
    ),
)
```

---

## 九、迁移计划

### 9.1 阶段规划

| 阶段 | 时间 | 任务 | 产出 |
|-----|------|------|------|
| **Phase 1** | Week 1-2 | 基础架构搭建 | 框架、数据库、proto定义 |
| **Phase 2** | Week 3-4 | 核心服务开发 | user-service, doctor-service, auth-service |
| **Phase 3** | Week 5-6 | 业务服务开发 | medical-service, registration-service |
| **Phase 4** | Week 7-8 | 支撑服务开发 | payment-service, ai-service, notification-service |
| **Phase 5** | Week 9 | 测试与集成 | 单元测试、集成测试、性能测试 |
| **Phase 6** | Week 10 | 灰度发布 | 线上验证、数据迁移、回滚方案 |

### 9.2 数据迁移

```go
// 数据迁移工具
package migration

// MigrateFromJava 从Java系统迁移数据
func MigrateFromJava(cfg MigrationConfig) error {
    // 1. 连接Java数据库
    // 2. 读取Java数据
    // 3. 转换数据格式（如JSON字段）
    // 4. 写入Go数据库
    // 5. 验证数据一致性
}
```

---

## 十、总结

本方案将Java SpringBoot单体应用迁移到Go+Kratos微服务架构，主要改进：

1. **微服务拆分**：12个独立服务，DDD领域建模
2. **强化权限**：RBAC+数据权限，AOP实现
3. **性能提升**：Go高性能+gRPC通信
4. **可扩展性**：容器化部署，弹性伸缩
5. **可观测性**：链路追踪+监控告警

**预期收益**：
- 性能提升：3-5倍（延迟降低、吞吐量提升）
- 开发效率：微服务独立部署、独立扩展
- 维护成本：代码量减少40%，模块化更清晰
- 可扩展性：支持水平扩展，应对高并发
