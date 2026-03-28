# 智慧医疗系统 - DDD架构分析文档

## 文档信息

| 属性 | 值 |
|-----|---|
| 版本 | V1.0 |
| 日期 | 2026-03-27 |
| 技术栈 | Go + Kratos + DDD |
| 目标 | 微服务领域驱动设计建模 |

---

## 目录

1. [DDD概述](#一ddd概述)
2. [领域上下文划分](#二领域上下文划分)
3. [领域模型详细设计](#三领域模型详细设计)
4. [聚合设计](#四聚合设计)
5. [领域事件](#五领域事件)
6. [防腐层设计](#六防腐层设计)
7. [服务间协作](#七服务间协作)
8. [DDD代码规范](#八ddd代码规范)

---

## 一、DDD概述

### 1.1 什么是DDD

领域驱动设计（Domain-Driven Design）是一种以业务领域为核心的软件设计方法。在智慧医疗系统中，DDD帮助我们：

- **统一语言**：业务、产品、开发使用同一套术语
- **边界清晰**：每个微服务有明确的业务边界
- **高内聚低耦合**：核心业务逻辑集中在领域层
- **可演进**：领域模型随业务变化而演进

### 1.2 为什么医疗系统需要DDD

医疗系统的业务复杂度极高：

```
┌─────────────────────────────────────────────────────────────┐
│                    医疗系统业务复杂度                          │
├─────────────────────────────────────────────────────────────┤
│                                                              │
│  业务流程复杂                    规则繁多                     │
│  ├── 挂号 → 排队 → 诊断 → 处方 → 结算                        │
│  ├── 手术安排 → 术前检查 → 手术 → 恢复                        │
│  └── 病历生命周期管理                                        │
│                                                              │
│  多角色协作                      数据一致性要求高              │
│  ├── 患者、医生、护士、药师、管理员                          │
│  ├── 科室主任、专家、普通医生                                │
│  └── 权限分级（本人/科室/全院）                              │
│                                                              │
│  合规要求                      业务知识密集                    │
│  ├── 医疗数据隐私保护                                        │
│  ├── 诊断记录不可篡改                                        │
│  └── 医学术语标准化                                          │
│                                                              │
└─────────────────────────────────────────────────────────────┘
```

### 1.3 DDD在微服务中的应用

```
┌─────────────────────────────────────────────────────────────┐
│              DDD与微服务的映射关系                            │
├─────────────────────────────────────────────────────────────┤
│                                                              │
│  DDD概念              微服务实现                              │
│  ─────────────────────────────────────────                   │
│  限界上下文(Bounded    微服务边界                             │
│  Context)                                                    │
│                                                              │
│  领域模型              服务的核心业务逻辑                      │
│  (Domain Model)                                              │
│                                                              │
│  领域服务              跨聚合的业务操作                        │
│  (Domain Service)                                            │
│                                                              │
│  应用服务              用例编排、事务管理                      │
│  (Application Service)                                       │
│                                                              │
│  领域事件              服务间异步通信                        │
│  (Domain Event)                                              │
│                                                              │
│  防腐层                外部服务适配器                        │
│  (Anti-Corruption Layer)                                     │
│                                                              │
└─────────────────────────────────────────────────────────────┘
```

---

## 二、领域上下文划分

### 2.1 限界上下文总览

基于智慧医疗系统的业务分析，识别出12个限界上下文：

```
┌──────────────────────────────────────────────────────────────────────┐
│                       限界上下文全景图                                │
├──────────────────────────────────────────────────────────────────────┤
│                                                                       │
│   ┌──────────────┐  ┌──────────────┐  ┌──────────────┐               │
│   │  患者上下文   │  │  医生上下文   │  │  科室上下文   │               │
│   │ Patient BC   │  │  Doctor BC   │  │Department BC │               │
│   └──────┬───────┘  └──────┬───────┘  └──────┬───────┘               │
│          │                 │                 │                        │
│          └─────────────────┼─────────────────┘                        │
│                            │                                          │
│   ┌──────────────┐  ┌──────┴──────┐  ┌──────────────┐               │
│   │  挂号上下文   │  │   病历上下文 │  │  手术上下文   │               │
│   │Registration │  │  Medical BC │  │ Surgical BC  │               │
│   │     BC       │  │             │  │              │               │
│   └──────┬───────┘  └──────┬──────┘  └──────┬───────┘               │
│          │                 │                 │                        │
│          └─────────────────┼─────────────────┘                        │
│                            │                                          │
│   ┌──────────────┐  ┌──────┴──────┐  ┌──────────────┐               │
│   │  药品上下文   │  │  支付上下文  │  │   AI上下文   │               │
│   │   Drug BC    │  │ Payment BC │  │    AI BC     │               │
│   └──────┬───────┘  └──────┬──────┘  └──────┬───────┘               │
│          │                 │                 │                        │
│          │                 │                 │                        │
│   ┌──────────────┐  ┌──────────────┐  ┌──────────────┐               │
│   │  权限上下文   │  │  通知上下文  │  │  文件上下文  │               │
│   │   Auth BC    │  │ Notification│  │    File BC   │               │
│   │              │  │     BC      │  │              │               │
│   └──────────────┘  └──────────────┘  └──────────────┘               │
│                                                                       │
└──────────────────────────────────────────────────────────────────────┘
```

### 2.2 上下文详细说明

| 序号 | 限界上下文 | 英文名称 | 核心业务 | 对应微服务 |
|-----|-----------|---------|---------|-----------|
| 1 | 患者上下文 | Patient BC | 患者信息管理、健康档案 | user-service |
| 2 | 医生上下文 | Doctor BC | 医生档案、排班、职称管理 | doctor-service |
| 3 | 科室上下文 | Department BC | 科室管理、医生分配 | department-service |
| 4 | 挂号上下文 | Registration BC | 普通挂号、专家挂号、排队 | registration-service |
| 5 | 病历上下文 | Medical BC | 病历创建、诊断、处方 | medical-service |
| 6 | 手术上下文 | Surgical BC | 手术安排、进度跟踪 | surgical-service |
| 7 | 药品上下文 | Drug BC | 药品信息、库存、分类 | drug-service |
| 8 | 支付上下文 | Payment BC | 挂号费、诊疗费结算 | payment-service |
| 9 | AI上下文 | AI BC | 智能导诊、RAG检索、Agent推理 | ai-service |
| 10 | 权限上下文 | Auth BC | 用户认证、角色权限、数据权限 | auth-service |
| 11 | 通知上下文 | Notification BC | 消息推送、通知管理 | notification-service |
| 12 | 文件上下文 | File BC | 文件上传、存储、管理 | file-service |

### 2.3 上下文关系图

```
┌──────────────────────────────────────────────────────────────────────┐
│                       上下文映射关系                                  │
├──────────────────────────────────────────────────────────────────────┤
│                                                                       │
│   合作关系 (Partnership)         共享内核 (Shared Kernel)             │
│   ═══════════════════════       ════════════════════════              │
│   病历 ←──────→ 挂号              患者/医生/科室 共享用户信息           │
│   病历 ←──────→ 手术              认证中心统一提供                      │
│                                                                       │
│   客户-供应商                    防腐层 (ACL)                         │
│   ═══════════════════            ═══════════════                      │
│   挂号 ───────→ 患者              AI服务 ──ACL──→ 向量数据库            │
│   挂号 ───────→ 医生              AI服务 ──ACL──→ 外部LLM API           │
│   挂号 ───────→ 科室              病历 ──ACL──→ 支付服务                │
│   病历 ───────→ 药品                                                   │
│                                                                       │
│   发布-订阅 (Pub-Sub)                                                   │
│   ═══════════════════                                                   │
│   病历 ──发布事件──→ 支付                                               │
│   病历 ──发布事件──→ 通知                                               │
│   挂号 ──发布事件──→ 通知                                               │
│                                                                       │
└──────────────────────────────────────────────────────────────────────┘
```

---

## 三、领域模型详细设计

### 3.1 患者上下文 (Patient BC)

#### 3.1.1 领域模型图

```
┌───────────────────────────────────────────────────────────────┐
│                    患者上下文领域模型                           │
├───────────────────────────────────────────────────────────────┤
│                                                                │
│   ┌───────────────────────────────────────┐                   │
│   │        Patient (聚合根)                │                   │
│   ├───────────────────────────────────────┤                   │
│   │ - ID: PatientID                       │                   │
│   │ - OpenID: string                      │                   │
│   │ - BasicInfo: PatientBasicInfo         │                   │
│   │ - Contact: ContactInfo                │                   │
│   │ - Status: PatientStatus               │                   │
│   │ - MedicalHistory: []MedicalHistory    │                   │
│   │ - Allergies: []Allergy                │                   │
│   │ - CreateTime: time.Time               │
│   │ - UpdateTime: time.Time               │
│   ├───────────────────────────────────────┤                   │
│   │ + UpdateProfile()                     │                   │
│   │ + ChangePassword()                    │                   │
│   │ + AddAllergy()                        │                   │
│   │ + AddMedicalHistory()                 │                   │
│   │ + SoftDelete()                        │                   │
│   └───────────────────────────────────────┘                   │
│                           │                                    │
│           ┌───────────────┴───────────────┐                    │
│           ▼                               ▼                    │
│   ┌───────────────┐               ┌───────────────┐           │
│   │  PatientBasic │               │  ContactInfo  │           │
│   │    Info       │               │    (VO)       │           │
│   │   (VO)        │               ├───────────────┤           │
│   ├───────────────┤               │ - Phone       │           │
│   │ - Name        │               │ - Email       │           │
│   │ - IDCard      │               │ - Address     │           │
│   │ - Age         │               │ - Emergency   │           │
│   │ - Gender      │               │   Contact     │           │
│   │ - BirthDate   │               └───────────────┘           │
│   └───────────────┘                                            │
│                                                                │
│   ┌───────────────┐               ┌───────────────┐           │
│   │MedicalHistory │               │   Allergy     │           │
│   │   (Entity)    │               │   (Entity)    │           │
│   ├───────────────┤               ├───────────────┤           │
│   │ - ID          │               │ - ID          │           │
│   │ - DiseaseName │               │ - Allergen    │           │
│   │ - DiagnoseDate│               │ - Severity    │           │
│   │ - Treatment   │               │ - Reaction    │           │
│   └───────────────┘               └───────────────┘           │
│                                                                │
└───────────────────────────────────────────────────────────────┘
```

#### 3.1.2 聚合根：Patient

```go
// service/patient/internal/domain/aggregate/patient.go
package aggregate

import (
    "time"
    
    "github.com/go-kratos/kratos/v2/log"
    "medical-system/patient/internal/domain/entity"
    "medical-system/patient/internal/domain/event"
    "medical-system/patient/internal/domain/vo"
)

// Patient 患者聚合根
type Patient struct {
    ID        PatientID
    OpenID    string
    BasicInfo vo.PatientBasicInfo
    Contact   vo.ContactInfo
    Status    PatientStatus
    
    // 关联实体
    MedicalHistory []entity.MedicalHistory
    Allergies      []entity.Allergy
    
    // 元数据
    CreatedAt time.Time
    UpdatedAt time.Time
    DeletedAt *time.Time
    
    // 领域事件
    Events []event.DomainEvent
}

// PatientID 患者ID（值对象）
type PatientID int64

// PatientStatus 患者状态
type PatientStatus int32

const (
    PatientStatusNormal   PatientStatus = 1  // 正常
    PatientStatusLocked   PatientStatus = 2  // 锁定
    PatientStatusInactive PatientStatus = 3  // 未激活
)

// 领域方法

// UpdateProfile 更新个人信息
func (p *Patient) UpdateProfile(info vo.PatientBasicInfo, contact vo.ContactInfo) error {
    if p.Status == PatientStatusLocked {
        return ErrPatientLocked
    }
    
    p.BasicInfo = info
    p.Contact = contact
    p.UpdatedAt = time.Now()
    
    // 发布领域事件
    p.Events = append(p.Events, event.NewPatientProfileUpdatedEvent(p.ID))
    
    return nil
}

// AddAllergy 添加过敏史
func (p *Patient) AddAllergy(allergen string, severity entity.AllergySeverity, reaction string) error {
    allergy := entity.Allergy{
        ID:        entity.AllergyID(time.Now().UnixNano()),
        Allergen:  allergen,
        Severity:  severity,
        Reaction:  reaction,
        CreatedAt: time.Now(),
    }
    
    p.Allergies = append(p.Allergies, allergy)
    p.UpdatedAt = time.Now()
    
    p.Events = append(p.Events, event.NewPatientAllergyAddedEvent(p.ID, allergy.ID))
    
    return nil
}

// AddMedicalHistory 添加病史
func (p *Patient) AddMedicalHistory(diseaseName string, diagnoseDate time.Time, treatment string) error {
    history := entity.MedicalHistory{
        ID:           entity.HistoryID(time.Now().UnixNano()),
        DiseaseName:  diseaseName,
        DiagnoseDate: diagnoseDate,
        Treatment:    treatment,
        CreatedAt:    time.Now(),
    }
    
    p.MedicalHistory = append(p.MedicalHistory, history)
    p.UpdatedAt = time.Now()
    
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
    ErrPatientLocked   = errors.New("patient is locked")
    ErrInvalidInfo     = errors.New("invalid patient info")
)
```

#### 3.1.3 值对象：PatientBasicInfo

```go
// service/patient/internal/domain/vo/patient_basic_info.go
package vo

import (
    "errors"
    "time"
)

// PatientBasicInfo 患者基本信息值对象
type PatientBasicInfo struct {
    Name      string
    IDCard    IDCard
    Age       int32
    Gender    Gender
    BirthDate time.Time
}

// Validate 验证基本信息
func (p PatientBasicInfo) Validate() error {
    if p.Name == "" {
        return errors.New("name is required")
    }
    if p.Age < 0 || p.Age > 150 {
        return errors.New("invalid age")
    }
    if err := p.IDCard.Validate(); err != nil {
        return err
    }
    return nil
}

// Gender 性别枚举
type Gender int32

const (
    GenderUnknown Gender = 0
    GenderMale    Gender = 1
    GenderFemale  Gender = 2
)

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

func (i IDCard) Validate() error {
    if !isValidIDCard(i.value) {
        return errors.New("invalid ID card format")
    }
    return nil
}

func (i IDCard) Masked() string {
    // 脱敏显示：110101********1234
    if len(i.value) != 18 {
        return i.value
    }
    return i.value[:6] + "********" + i.value[14:]
}

func isValidIDCard(idCard string) bool {
    // 简化验证：长度18位
    return len(idCard) == 18
}
```

#### 3.1.4 实体：MedicalHistory

```go
// service/patient/internal/domain/entity/medical_history.go
package entity

import "time"

// HistoryID 病史ID
type HistoryID int64

// MedicalHistory 病史实体
type MedicalHistory struct {
    ID           HistoryID
    DiseaseName  string
    DiagnoseDate time.Time
    Treatment    string
    CreatedAt    time.Time
}
```

---

### 3.2 医生上下文 (Doctor BC)

#### 3.2.1 领域模型图

```
┌───────────────────────────────────────────────────────────────┐
│                    医生上下文领域模型                           │
├───────────────────────────────────────────────────────────────┤
│                                                                │
│   ┌───────────────────────────────────────┐                   │
│   │         Doctor (聚合根)                │                   │
│   ├───────────────────────────────────────┤                   │
│   │ - ID: DoctorID                        │                   │
│   │ - EmployeeID: string                  │                   │
│   │ - BasicInfo: DoctorBasicInfo          │                   │
│   │ - Professional: ProfessionalInfo      │                   │
│   │ - DepartmentID: DepartmentID          │                   │
│   │ - Status: DoctorStatus                │                   │
│   │ - Schedule: []Schedule                │                   │
│   │ - Roles: []Role                       │                   │
│   ├───────────────────────────────────────┤                   │
│   │ + ChangeDepartment()                  │                   │
│   │ + UpdateProfessional()                │                   │
│   │ + SetSchedule()                       │                   │
│   │ + AssignRole()                        │                   │
│   │ + IsExpert() bool                     │                   │
│   └───────────────────────────────────────┘                   │
│                           │                                    │
│           ┌───────────────┼───────────────┐                    │
│           ▼               ▼               ▼                    │
│   ┌───────────┐    ┌───────────┐    ┌───────────┐             │
│   │ DoctorBasic│    │Professional│    │  Schedule │             │
│   │   Info    │    │   Info     │    │  (Entity) │             │
│   │   (VO)    │    │    (VO)    │    ├───────────┤             │
│   ├───────────┤    ├───────────┤    │ - Date    │             │
│   │ - Name    │    │ - Position │    │ - TimeSlots│            │
│   │ - Phone   │    │ - Title    │    │ - IsAvailable│          │
│   │ - Avatar  │    │ - Specialty│    └───────────┘             │
│   └───────────┘    │ - IsExpert │                              │
│                    └───────────┘                              │
│                                                                │
│   ┌───────────┐               ┌───────────┐                   │
│   │   Role    │               │Department │                   │
│   │ (Entity)  │               │Reference  │                   │
│   ├───────────┤               │   (ID)    │                   │
│   │ - Code    │               └───────────┘                   │
│   │ - Name    │                                                │
│   │ - DataScope│                                               │
│   └───────────┘                                                │
│                                                                │
└───────────────────────────────────────────────────────────────┘
```

#### 3.2.2 聚合根：Doctor

```go
// service/doctor/internal/domain/aggregate/doctor.go
package aggregate

import (
    "time"
    
    "medical-system/doctor/internal/domain/entity"
    "medical-system/doctor/internal/domain/event"
    "medical-system/doctor/internal/domain/vo"
)

// DoctorID 医生ID
type DoctorID int64

// Doctor 医生聚合根
type Doctor struct {
    ID         DoctorID
    EmployeeID string
    BasicInfo  vo.DoctorBasicInfo
    Professional vo.ProfessionalInfo
    
    // 外键引用（其他上下文的ID）
    DepartmentID DepartmentID
    
    Status DoctorStatus
    
    // 关联实体
    Schedules []entity.Schedule
    Roles     []entity.Role
    
    // 元数据
    CreatedAt time.Time
    UpdatedAt time.Time
    DeletedAt *time.Time
    
    Events []event.DomainEvent
}

// DoctorStatus 医生工作状态
type DoctorStatus int32

const (
    DoctorStatusWorking DoctorStatus = 1  // 工作中
    DoctorStatusResting DoctorStatus = 2  // 休息中
    DoctorStatusLeave   DoctorStatus = 3  // 请假中
)

// ChangeDepartment 更换科室
func (d *Doctor) ChangeDepartment(newDeptID DepartmentID) error {
    if d.Status == DoctorStatusLeave {
        return ErrDoctorOnLeave
    }
    
    oldDeptID := d.DepartmentID
    d.DepartmentID = newDeptID
    d.UpdatedAt = time.Now()
    
    d.Events = append(d.Events, event.NewDoctorDepartmentChangedEvent(
        d.ID, oldDeptID, newDeptID,
    ))
    
    return nil
}

// UpdateProfessional 更新职称信息
func (d *Doctor) UpdateProfessional(info vo.ProfessionalInfo) error {
    d.Professional = info
    d.UpdatedAt = time.Now()
    
    d.Events = append(d.Events, event.NewDoctorProfessionalUpdatedEvent(d.ID))
    
    return nil
}

// SetSchedule 设置排班
func (d *Doctor) SetSchedule(date time.Time, timeSlots []entity.TimeSlot) error {
    // 查找是否已有排班
    for i, schedule := range d.Schedules {
        if isSameDay(schedule.Date, date) {
            d.Schedules[i].TimeSlots = timeSlots
            d.Schedules[i].UpdatedAt = time.Now()
            return nil
        }
    }
    
    // 新建排班
    newSchedule := entity.Schedule{
        ID:        entity.ScheduleID(time.Now().UnixNano()),
        Date:      date,
        TimeSlots: timeSlots,
        CreatedAt: time.Now(),
        UpdatedAt: time.Now(),
    }
    d.Schedules = append(d.Schedules, newSchedule)
    
    return nil
}

// AssignRole 分配角色
func (d *Doctor) AssignRole(role entity.Role) error {
    // 检查是否已有该角色
    for _, r := range d.Roles {
        if r.Code == role.Code {
            return ErrRoleAlreadyAssigned
        }
    }
    
    d.Roles = append(d.Roles, role)
    d.UpdatedAt = time.Now()
    
    d.Events = append(d.Events, event.NewDoctorRoleAssignedEvent(d.ID, role.Code))
    
    return nil
}

// IsExpert 是否专家
func (d *Doctor) IsExpert() bool {
    return d.Professional.IsExpert
}

// CanTreat 是否可以接诊
func (d *Doctor) CanTreat() bool {
    return d.Status == DoctorStatusWorking
}

var (
    ErrDoctorOnLeave       = errors.New("doctor is on leave")
    ErrRoleAlreadyAssigned = errors.New("role already assigned")
)

func isSameDay(t1, t2 time.Time) bool {
    return t1.Year() == t2.Year() && t1.Month() == t2.Month() && t1.Day() == t2.Day()
}
```

#### 3.2.3 值对象：ProfessionalInfo

```go
// service/doctor/internal/domain/vo/professional_info.go
package vo

// ProfessionalInfo 专业信息值对象
type ProfessionalInfo struct {
    Position  Position  // 职务
    Title     Title     // 职称
    Specialty string    // 专业领域
    IsExpert  bool      // 是否专家
    Intro     string    // 个人简介
}

// Position 职务枚举
type Position int32

const (
    PositionChiefPhysician     Position = 1  // 主任医师
    PositionAssociateChief     Position = 2  // 副主任医师
    PositionAttending          Position = 3  // 主治医师
    PositionResident           Position = 4  // 住院医师
    PositionAssistant          Position = 5  // 助理医师
)

func (p Position) String() string {
    switch p {
    case PositionChiefPhysician:
        return "主任医师"
    case PositionAssociateChief:
        return "副主任医师"
    case PositionAttending:
        return "主治医师"
    case PositionResident:
        return "住院医师"
    case PositionAssistant:
        return "助理医师"
    default:
        return "未知"
    }
}

// Title 职称枚举
type Title int32

const (
    TitleDean           Title = 1  // 院长
    TitleViceDean       Title = 2  // 副院长
    TitleDeptHead       Title = 3  // 科室主任
    TitleMedicalLeader  Title = 4  // 医疗组长
    TitleMedicalMember  Title = 5  // 医疗组员
)

func (t Title) String() string {
    switch t {
    case TitleDean:
        return "院长"
    case TitleViceDean:
        return "副院长"
    case TitleDeptHead:
        return "科室主任"
    case TitleMedicalLeader:
        return "医疗组长"
    case TitleMedicalMember:
        return "医疗组员"
    default:
        return ""
    }
}

// GetDisplayTitle 获取显示职称
func (p ProfessionalInfo) GetDisplayTitle() string {
    return p.Position.String() + "/" + p.Title.String()
}
```

---

### 3.3 病历上下文 (Medical BC)

#### 3.3.1 领域模型图

```
┌─────────────────────────────────────────────────────────────────────┐
│                      病历上下文领域模型                              │
├─────────────────────────────────────────────────────────────────────┤
│                                                                      │
│   ┌─────────────────────────────────────────┐                       │
│   │         MedicalRecord (聚合根)           │                       │
│   ├─────────────────────────────────────────┤                       │
│   │ - ID: MedicalRecordID                   │                       │
│   │ - MedicalNumber: string                 │                       │
│   │ - PatientID: PatientID                  │                       │
│   │ - DoctorID: DoctorID                    │                       │
│   │ - DepartmentID: DepartmentID            │                       │
│   │ - Status: MedicalStatus                 │                       │
│   │ - Symptom: string                       │                       │
│   │ - Diagnosis: string                     │                       │
│   │ - TreatmentPlan: string                 │                       │
│   │ - Diseases: []DiseaseItem               │                       │
│   │ - Prescriptions: []PrescriptionItem     │                       │
│   │ - Payment: PaymentInfo                  │                       │
│   │ - Timeline: []StatusChange              │                       │
│   ├─────────────────────────────────────────┤                       │
│   │ + AddDisease()                          │                       │
│   │ + AddPrescription()                     │                       │
│   │ + UpdateStatus()                        │                       │
│   │ + Complete()                            │                       │
│   │ + Cancel()                              │                       │
│   │ + CalculateTotal()                      │                       │
│   └─────────────────────────────────────────┘                       │
│                           │                                          │
│           ┌───────────────┼───────────────┐                         │
│           ▼               ▼               ▼                         │
│   ┌───────────┐    ┌───────────┐    ┌───────────┐                  │
│   │DiseaseItem│    │Prescription│    │ StatusChange              │
│   │ (Entity)  │    │   Item     │    │  (Entity)                  │
│   ├───────────┤    │  (Entity)  │    ├───────────┤                  │
│   │ - ID      │    ├───────────┤    │ - ID      │                  │
│   │ - DiseaseID│   │ - ID      │    │ - From    │                  │
│   │ - Name    │    │ - DrugID  │    │ - To      │                  │
│   │ - Symptoms│    │ - Name    │    │ - Time    │                  │
│   │ - Diagnosis│   │ - Quantity│    │ - Operator│                  │
│   │ - Treatment│   │ - Price   │    └───────────┘                  │
│   └───────────┘    │ - Usage   │                                   │
│                    └───────────┘                                   │
│   ┌───────────┐                                                    │
│   │ PaymentInfo│                                                   │
│   │   (VO)    │                                                    │
│   ├───────────┤                                                    │
│   │ - Amount  │                                                    │
│   │ - Method  │                                                    │
│   │ - Status  │                                                    │
│   │ - PayTime │                                                    │
│   └───────────┘                                                    │
│                                                                      │
└─────────────────────────────────────────────────────────────────────┘
```

#### 3.3.2 聚合根：MedicalRecord

```go
// service/medical/internal/domain/aggregate/medical_record.go
package aggregate

import (
    "time"
    
    "medical-system/medical/internal/domain/entity"
    "medical-system/medical/internal/domain/event"
    "medical-system/medical/internal/domain/vo"
)

// MedicalRecordID 病历ID
type MedicalRecordID int64

// MedicalRecord 病历聚合根
type MedicalRecord struct {
    ID            MedicalRecordID
    MedicalNumber string
    
    // 外键引用
    PatientID    PatientID
    DoctorID     DoctorID
    DepartmentID DepartmentID
    
    // 诊断信息
    Symptom       string
    Diagnosis     string
    TreatmentPlan string
    
    // 关联实体
    Diseases      []entity.DiseaseItem
    Prescriptions []entity.PrescriptionItem
    
    // 支付信息（值对象）
    Payment *vo.PaymentInfo
    
    // 状态跟踪
    Status   MedicalStatus
    Timeline []entity.StatusChange
    
    // 时间戳
    BillingTime    *time.Time
    SettlementTime *time.Time
    CreatedAt      time.Time
    UpdatedAt      time.Time
    DeletedAt      *time.Time
    
    Events []event.DomainEvent
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

func (s MedicalStatus) String() string {
    switch s {
    case MedicalStatusRegistering:
        return "挂号中"
    case MedicalStatusQueuing:
        return "排队中"
    case MedicalStatusDiagnosing:
        return "诊断中"
    case MedicalStatusOperating:
        return "手术中"
    case MedicalStatusCompleted:
        return "诊断结束"
    case MedicalStatusCancelled:
        return "已取消"
    default:
        return "未知"
    }
}

// AddDisease 添加诊断疾病
func (m *MedicalRecord) AddDisease(disease entity.DiseaseItem) error {
    if !m.CanAddDisease() {
        return ErrInvalidStatusForAddingDisease
    }
    
    m.Diseases = append(m.Diseases, disease)
    m.UpdatedAt = time.Now()
    
    return nil
}

// AddPrescription 添加处方
func (m *MedicalRecord) AddPrescription(prescription entity.PrescriptionItem) error {
    if !m.CanAddPrescription() {
        return ErrInvalidStatusForPrescription
    }
    
    m.Prescriptions = append(m.Prescriptions, prescription)
    m.UpdatedAt = time.Now()
    
    return nil
}

// UpdateStatus 更新状态
func (m *MedicalRecord) UpdateStatus(newStatus MedicalStatus, operatorID int64, reason string) error {
    if !m.IsValidStatusTransition(m.Status, newStatus) {
        return ErrInvalidStatusTransition
    }
    
    oldStatus := m.Status
    m.Status = newStatus
    m.UpdatedAt = time.Now()
    
    // 记录状态变更
    change := entity.StatusChange{
        ID:        entity.StatusChangeID(time.Now().UnixNano()),
        From:      oldStatus,
        To:        newStatus,
        Time:      time.Now(),
        OperatorID: operatorID,
        Reason:    reason,
    }
    m.Timeline = append(m.Timeline, change)
    
    // 发布领域事件
    m.Events = append(m.Events, event.NewMedicalStatusChangedEvent(
        m.ID, oldStatus, newStatus, operatorID,
    ))
    
    return nil
}

// Complete 完成诊断
func (m *MedicalRecord) Complete(operatorID int64) error {
    if m.Status != MedicalStatusDiagnosing && m.Status != MedicalStatusOperating {
        return ErrCannotComplete
    }
    
    now := time.Now()
    m.SettlementTime = &now
    
    return m.UpdateStatus(MedicalStatusCompleted, operatorID, "诊断完成")
}

// Cancel 取消病历
func (m *MedicalRecord) Cancel(operatorID int64, reason string) error {
    if m.Status == MedicalStatusCompleted {
        return ErrCannotCancelCompleted
    }
    
    return m.UpdateStatus(MedicalStatusCancelled, operatorID, reason)
}

// CalculateTotal 计算总费用
func (m *MedicalRecord) CalculateTotal() float64 {
    var total float64
    for _, p := range m.Prescriptions {
        total += p.Price * float64(p.Quantity)
    }
    return total
}

// CanAddDisease 是否可以添加疾病诊断
func (m *MedicalRecord) CanAddDisease() bool {
    return m.Status == MedicalStatusDiagnosing
}

// CanAddPrescription 是否可以添加处方
func (m *MedicalRecord) CanAddPrescription() bool {
    return m.Status == MedicalStatusDiagnosing
}

// IsValidStatusTransition 状态转换是否有效
func (m *MedicalRecord) IsValidStatusTransition(from, to MedicalStatus) bool {
    // 定义合法的状态转换
    validTransitions := map[MedicalStatus][]MedicalStatus{
        MedicalStatusRegistering: {MedicalStatusQueuing, MedicalStatusCancelled},
        MedicalStatusQueuing:     {MedicalStatusDiagnosing, MedicalStatusCancelled},
        MedicalStatusDiagnosing:  {MedicalStatusOperating, MedicalStatusCompleted, MedicalStatusCancelled},
        MedicalStatusOperating:   {MedicalStatusCompleted, MedicalStatusCancelled},
    }
    
    validNextStates, exists := validTransitions[from]
    if !exists {
        return false
    }
    
    for _, validState := range validNextStates {
        if validState == to {
            return true
        }
    }
    return false
}

var (
    ErrInvalidStatusForAddingDisease = errors.New("current status cannot add disease")
    ErrInvalidStatusForPrescription  = errors.New("current status cannot add prescription")
    ErrInvalidStatusTransition       = errors.New("invalid status transition")
    ErrCannotComplete                = errors.New("cannot complete medical record")
    ErrCannotCancelCompleted         = errors.New("cannot cancel completed medical record")
)
```

#### 3.3.3 实体：PrescriptionItem

```go
// service/medical/internal/domain/entity/prescription_item.go
package entity

// PrescriptionItemID 处方项ID
type PrescriptionItemID int64

// PrescriptionItem 处方项实体
type PrescriptionItem struct {
    ID       PrescriptionItemID
    DrugID   DrugID
    DrugName string
    Quantity int32
    Price    float64
    Usage    string
}

// CalculateSubtotal 计算小计
func (p PrescriptionItem) CalculateSubtotal() float64 {
    return p.Price * float64(p.Quantity)
}
```

---

### 3.4 AI上下文 (AI BC) - 核心创新

#### 3.4.1 领域模型图

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                           AI上下文领域模型                                    │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│   ┌─────────────────────────────────────────────────────────┐               │
│   │               DiagnosisSession (聚合根)                  │               │
│   ├─────────────────────────────────────────────────────────┤               │
│   │ - ID: SessionID                                         │               │
│   │ - PatientID: PatientID                                  │               │
│   │ - Symptoms: []Symptom                                   │               │
│   │ - Conversations: []Conversation                         │               │
│   │ - Status: SessionStatus                                 │               │
│   │ - Recommendation: *Recommendation                       │               │
│   │ - VectorQueries: []VectorQuery                          │               │
│   ├─────────────────────────────────────────────────────────┤               │
│   │ + AddSymptom()                                          │               │
│   │ + AddConversation()                                     │               │
│   │ + GenerateRecommendation()                              │               │
│   │ + Complete()                                            │               │
│   │ + GetContextForRAG() string                             │               │
│   └─────────────────────────────────────────────────────────┘               │
│                           │                                                  │
│           ┌───────────────┼───────────────┐                                 │
│           ▼               ▼               ▼                                 │
│   ┌───────────┐    ┌───────────┐    ┌───────────┐                          │
│   │  Symptom  │    │Conversation│    │ VectorQuery                         │
│   │ (Entity)  │    │  (Entity)  │    │  (Entity)                          │
│   ├───────────┤    ├───────────┤    ├───────────┤                          │
│   │ - ID      │    │ - ID      │    │ - ID      │                          │
│   │ - BodyPart│    │ - Role    │    │ - Query   │                          │
│   │ - Description│ │ - Content │    │ - Results │                          │
│   │ - Severity│    │ - Time    │    │ - Source  │                          │
│   └───────────┘    └───────────┘    └───────────┘                          │
│                                                                              │
│   ┌─────────────────────────────────────────────────────────┐               │
│   │              Recommendation (值对象)                     │               │
│   ├─────────────────────────────────────────────────────────┤               │
│   │ - SuspectedDiseases: []DiseaseRecommendation            │               │
│   │ - RecommendedDepartments: []DepartmentRecommendation    │               │
│   │ - RecommendedDoctors: []DoctorRecommendation            │               │
│   │ - RecommendedDrugs: []DrugRecommendation                │               │
│   │ - Analysis: string                                      │               │
│   │ - Confidence: float64                                   │               │
│   └─────────────────────────────────────────────────────────┘               │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

#### 3.4.2 聚合根：DiagnosisSession

```go
// service/ai/internal/domain/aggregate/diagnosis_session.go
package aggregate

import (
    "strings"
    "time"
    
    "medical-system/ai/internal/domain/entity"
    "medical-system/ai/internal/domain/event"
    "medical-system/ai/internal/domain/vo"
)

// SessionID 会话ID
type SessionID string

// DiagnosisSession AI导诊会话聚合根
type DiagnosisSession struct {
    ID       SessionID
    PatientID PatientID
    
    // 症状收集
    Symptoms []entity.Symptom
    
    // 对话历史
    Conversations []entity.Conversation
    
    // 向量检索记录
    VectorQueries []entity.VectorQuery
    
    // 推荐结果
    Recommendation *vo.Recommendation
    
    Status SessionStatus
    
    CreatedAt time.Time
    UpdatedAt time.Time
    
    Events []event.DomainEvent
}

// SessionStatus 会话状态
type SessionStatus int32

const (
    SessionStatusCollecting   SessionStatus = 1  // 收集中
    SessionStatusAnalyzing    SessionStatus = 2  // 分析中
    SessionStatusRecommended  SessionStatus = 3  // 已推荐
    SessionStatusCompleted    SessionStatus = 4  // 已完成
)

// AddSymptom 添加症状
func (s *DiagnosisSession) AddSymptom(bodyPart, description string, severity entity.SeverityLevel) error {
    if s.Status != SessionStatusCollecting {
        return ErrCannotAddSymptom
    }
    
    symptom := entity.Symptom{
        ID:          entity.SymptomID(time.Now().UnixNano()),
        BodyPart:    bodyPart,
        Description: description,
        Severity:    severity,
        CreatedAt:   time.Now(),
    }
    
    s.Symptoms = append(s.Symptoms, symptom)
    s.UpdatedAt = time.Now()
    
    return nil
}

// AddConversation 添加对话
func (s *DiagnosisSession) AddConversation(role entity.ConversationRole, content string) {
    conv := entity.Conversation{
        ID:        entity.ConversationID(time.Now().UnixNano()),
        Role:      role,
        Content:   content,
        Timestamp: time.Now(),
    }
    
    s.Conversations = append(s.Conversations, conv)
    s.UpdatedAt = time.Now()
}

// RecordVectorQuery 记录向量查询
func (s *DiagnosisSession) RecordVectorQuery(query string, source entity.VectorSource, results []entity.VectorResult) {
    vq := entity.VectorQuery{
        ID:        entity.VectorQueryID(time.Now().UnixNano()),
        Query:     query,
        Source:    source,
        Results:   results,
        Timestamp: time.Now(),
    }
    
    s.VectorQueries = append(s.VectorQueries, vq)
}

// GenerateRecommendation 生成推荐
func (s *DiagnosisSession) GenerateRecommendation(rec vo.Recommendation) error {
    if len(s.Symptoms) == 0 {
        return ErrNoSymptoms
    }
    
    s.Recommendation = &rec
    s.Status = SessionStatusRecommended
    s.UpdatedAt = time.Now()
    
    s.Events = append(s.Events, event.NewDiagnosisRecommendationGeneratedEvent(
        s.ID, s.PatientID, rec,
    ))
    
    return nil
}

// GetContextForRAG 获取RAG上下文
func (s *DiagnosisSession) GetContextForRAG() string {
    var context strings.Builder
    
    // 添加症状描述
    context.WriteString("症状描述：\n")
    for _, symptom := range s.Symptoms {
        context.WriteString("- ")
        if symptom.BodyPart != "" {
            context.WriteString(symptom.BodyPart + ": ")
        }
        context.WriteString(symptom.Description)
        context.WriteString(" (严重程度: " + symptom.Severity.String() + ")\n")
    }
    
    // 添加对话历史摘要
    if len(s.Conversations) > 0 {
        context.WriteString("\n对话历史：\n")
        for _, conv := range s.Conversations {
            context.WriteString(conv.Role.String() + ": " + conv.Content + "\n")
        }
    }
    
    return context.String()
}

// Complete 完成会话
func (s *DiagnosisSession) Complete() error {
    if s.Recommendation == nil {
        return ErrNoRecommendation
    }
    
    s.Status = SessionStatusCompleted
    s.UpdatedAt = time.Now()
    
    s.Events = append(s.Events, event.NewDiagnosisSessionCompletedEvent(s.ID, s.PatientID))
    
    return nil
}

// IsReadyForAnalysis 是否准备好分析
func (s *DiagnosisSession) IsReadyForAnalysis() bool {
    return len(s.Symptoms) > 0
}

var (
    ErrCannotAddSymptom = errors.New("cannot add symptom in current status")
    ErrNoSymptoms       = errors.New("no symptoms collected")
    ErrNoRecommendation = errors.New("no recommendation generated")
)
```

#### 3.4.3 值对象：Recommendation

```go
// service/ai/internal/domain/vo/recommendation.go
package vo

// Recommendation 推荐结果值对象
type Recommendation struct {
    SuspectedDiseases      []DiseaseRecommendation
    RecommendedDepartments []DepartmentRecommendation
    RecommendedDoctors     []DoctorRecommendation
    RecommendedDrugs       []DrugRecommendation
    Analysis               string
    Confidence             float64
    GeneratedAt            int64
}

// DiseaseRecommendation 疾病推荐
type DiseaseRecommendation struct {
    DiseaseID   int64
    DiseaseName string
    Symptoms    string
    Diagnosis   string
    Probability float64
}

// DepartmentRecommendation 科室推荐
type DepartmentRecommendation struct {
    DepartmentID   int64
    DepartmentName string
    Introduction   string
    Reason         string
    MatchScore     float64
}

// DoctorRecommendation 医生推荐
type DoctorRecommendation struct {
    DoctorID     int64
    DoctorName   string
    Title        string
    IsExpert     bool
    DepartmentID int64
    Intro        string
    MatchScore   float64
}

// DrugRecommendation 药品推荐
type DrugRecommendation struct {
    DrugID      int64
    DrugName    string
    Description string
    Category    string
    Usage       string
    Reason      string
}

// GetTopDisease 获取最可能的疾病
func (r Recommendation) GetTopDisease() *DiseaseRecommendation {
    if len(r.SuspectedDiseases) == 0 {
        return nil
    }
    
    top := &r.SuspectedDiseases[0]
    for i := range r.SuspectedDiseases {
        if r.SuspectedDiseases[i].Probability > top.Probability {
            top = &r.SuspectedDiseases[i]
        }
    }
    return top
}

// GetExpertDoctors 获取专家医生
func (r Recommendation) GetExpertDoctors() []DoctorRecommendation {
    var experts []DoctorRecommendation
    for _, doc := range r.RecommendedDoctors {
        if doc.IsExpert {
            experts = append(experts, doc)
        }
    }
    return experts
}
```

#### 3.4.4 实体：Symptom

```go
// service/ai/internal/domain/entity/symptom.go
package entity

import "time"

// SymptomID 症状ID
type SymptomID int64

// Symptom 症状实体
type Symptom struct {
    ID          SymptomID
    BodyPart    string    // 身体部位
    Description string    // 症状描述
    Severity    SeverityLevel
    Duration    string    // 持续时间
    CreatedAt   time.Time
}

// SeverityLevel 严重程度
type SeverityLevel int32

const (
    SeverityMild    SeverityLevel = 1  // 轻微
    SeverityModerate SeverityLevel = 2  // 中等
    SeveritySevere   SeverityLevel = 3  // 严重
    SeverityCritical SeverityLevel = 4  // 危急
)

func (s SeverityLevel) String() string {
    switch s {
    case SeverityMild:
        return "轻微"
    case SeverityModerate:
        return "中等"
    case SeveritySevere:
        return "严重"
    case SeverityCritical:
        return "危急"
    default:
        return "未知"
    }
}
```

---

## 四、聚合设计

### 4.1 聚合设计原则

在智慧医疗系统中，聚合设计遵循以下原则：

```
┌─────────────────────────────────────────────────────────────┐
│                    聚合设计原则                              │
├─────────────────────────────────────────────────────────────┤
│                                                              │
│  1. 业务一致性边界                                           │
│     聚合内的实体必须保持一致性，一起保存                       │
│     示例：MedicalRecord + Diseases + Prescriptions            │
│                                                              │
│  2. 聚合根唯一入口                                           │
│     外部只能通过聚合根访问聚合内的实体                         │
│     示例：只能通过 Patient 访问 MedicalHistory                │
│                                                              │
│  3. 小聚合优先                                               │
│     一个聚合应该只有一个聚合根，3-5个实体为宜                  │
│     示例：DiagnosisSession 包含 Symptoms、Conversations       │
│                                                              │
│  4. 跨聚合通过ID引用                                         │
│     不直接引用其他聚合的实体，只引用ID                         │
│     示例：MedicalRecord 中 DoctorID 而不是 Doctor 对象        │
│                                                              │
│  5. 事务边界                                                 │
│     一个事务只修改一个聚合                                    │
│     跨聚合使用领域事件                                        │
│                                                              │
└─────────────────────────────────────────────────────────────┘
```

### 4.2 聚合一览表

| 限界上下文 | 聚合根 | 包含实体 | 业务规则 |
|-----------|--------|---------|---------|
| 患者上下文 | Patient | MedicalHistory, Allergy | 患者信息完整、过敏史追踪 |
| 医生上下文 | Doctor | Schedule, Role | 排班不冲突、角色权限 |
| 科室上下文 | Department | DepartmentDoctor | 科室医生分配 |
| 挂号上下文 | Registration | - | 号源管理、排队逻辑 |
| 病历上下文 | MedicalRecord | DiseaseItem, PrescriptionItem, StatusChange | 状态流转、处方配伍 |
| 手术上下文 | Surgical | SurgicalDoctor, SurgicalFlow | 手术流程、参与医生 |
| 药品上下文 | Drug | DrugStock | 库存管理、效期管理 |
| AI上下文 | DiagnosisSession | Symptom, Conversation, VectorQuery | 症状收集、对话上下文 |

---

## 五、领域事件

### 5.1 事件总览

```
┌─────────────────────────────────────────────────────────────────────┐
│                        领域事件全景图                                │
├─────────────────────────────────────────────────────────────────────┤
│                                                                      │
│   患者上下文事件                                                      │
│   ┌──────────────────────────────────────┐                          │
│   │ PatientProfileUpdated                 │                          │
│   │ PatientAllergyAdded                   │                          │
│   │ PatientDeleted                        │                          │
│   └────────────────┬─────────────────────┘                          │
│                    │                                                  │
│                    ▼                                                  │
│   病历上下文事件          通知上下文                                  │
│   ┌────────────────┐      ┌────────────────┐                         │
│   │MedicalCreated  │──────▶│ Notification   │                         │
│   │StatusChanged   │      │ Sent           │                         │
│   │Completed       │      └────────────────┘                         │
│   └────────────────┘                                                  │
│                                                                      │
│   挂号上下文事件          病历上下文                                  │
│   ┌────────────────┐      ┌────────────────┐                         │
│   │RegistrationCreated│───▶│ MedicalCreate  │                         │
│   │RegistrationCancelled│  │ (自动创建病历)  │                         │
│   └────────────────┘      └────────────────┘                         │
│                                                                      │
│   AI上下文事件                                                        │
│   ┌──────────────────────────────────────┐                          │
│   │ DiagnosisRecommendationGenerated      │                          │
│   │ DiagnosisSessionCompleted             │                          │
│   │ VectorQueryPerformed                  │                          │
│   └──────────────────────────────────────┘                          │
│                                                                      │
└─────────────────────────────────────────────────────────────────────┘
```

### 5.2 核心领域事件定义

```go
// common/domain/event/medical_events.go
package event

import (
    "time"
    
    "medical-system/ai/internal/domain/vo"
)

// MedicalStatusChangedEvent 病历状态变更事件
type MedicalStatusChangedEvent struct {
    EventID    string    `json:"event_id"`
    EventType  string    `json:"event_type"`
    EventTime  time.Time `json:"event_time"`
    MedicalID  int64     `json:"medical_id"`
    OldStatus  int32     `json:"old_status"`
    NewStatus  int32     `json:"new_status"`
    OperatorID int64     `json:"operator_id"`
}

func NewMedicalStatusChangedEvent(medicalID int64, oldStatus, newStatus int32, operatorID int64) *MedicalStatusChangedEvent {
    return &MedicalStatusChangedEvent{
        EventID:    generateEventID(),
        EventType:  "medical.status_changed",
        EventTime:  time.Now(),
        MedicalID:  medicalID,
        OldStatus:  oldStatus,
        NewStatus:  newStatus,
        OperatorID: operatorID,
    }
}

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

func NewMedicalCompletedEvent(medicalID, patientID, doctorID, deptID int64, amount float64) *MedicalCompletedEvent {
    return &MedicalCompletedEvent{
        EventID:      generateEventID(),
        EventType:    "medical.completed",
        EventTime:    time.Now(),
        MedicalID:    medicalID,
        PatientID:    patientID,
        DoctorID:     doctorID,
        DepartmentID: deptID,
        Amount:       amount,
    }
}
```

```go
// service/ai/internal/domain/event/diagnosis_events.go
package event

import (
    "time"
    
    "medical-system/ai/internal/domain/vo"
)

// DiagnosisRecommendationGeneratedEvent 导诊推荐生成事件
type DiagnosisRecommendationGeneratedEvent struct {
    EventID       string            `json:"event_id"`
    EventType     string            `json:"event_type"`
    EventTime     time.Time         `json:"event_time"`
    SessionID     string            `json:"session_id"`
    PatientID     int64             `json:"patient_id"`
    Recommendation vo.Recommendation `json:"recommendation"`
}

func NewDiagnosisRecommendationGeneratedEvent(sessionID string, patientID int64, rec vo.Recommendation) *DiagnosisRecommendationGeneratedEvent {
    return &DiagnosisRecommendationGeneratedEvent{
        EventID:       generateEventID(),
        EventType:     "diagnosis.recommendation_generated",
        EventTime:     time.Now(),
        SessionID:     sessionID,
        PatientID:     patientID,
        Recommendation: rec,
    }
}

// DiagnosisSessionCompletedEvent 导诊会话完成事件
type DiagnosisSessionCompletedEvent struct {
    EventID   string    `json:"event_id"`
    EventType string    `json:"event_type"`
    EventTime time.Time `json:"event_time"`
    SessionID string    `json:"session_id"`
    PatientID int64     `json:"patient_id"`
}

func NewDiagnosisSessionCompletedEvent(sessionID string, patientID int64) *DiagnosisSessionCompletedEvent {
    return &DiagnosisSessionCompletedEvent{
        EventID:   generateEventID(),
        EventType: "diagnosis.session_completed",
        EventTime: time.Now(),
        SessionID: sessionID,
        PatientID: patientID,
    }
}

// VectorQueryPerformedEvent 向量查询执行事件
type VectorQueryPerformedEvent struct {
    EventID   string    `json:"event_id"`
    EventType string    `json:"event_type"`
    EventTime time.Time `json:"event_time"`
    SessionID string    `json:"session_id"`
    Query     string    `json:"query"`
    Source    string    `json:"source"`
    ResultCount int     `json:"result_count"`
}

func NewVectorQueryPerformedEvent(sessionID, query, source string, resultCount int) *VectorQueryPerformedEvent {
    return &VectorQueryPerformedEvent{
        EventID:     generateEventID(),
        EventType:   "diagnosis.vector_query_performed",
        EventTime:   time.Now(),
        SessionID:   sessionID,
        Query:       query,
        Source:      source,
        ResultCount: resultCount,
    }
}
```

---

## 六、防腐层设计

### 6.1 防腐层概念

```
┌─────────────────────────────────────────────────────────────┐
│                     防腐层 (ACL) 设计                        │
├─────────────────────────────────────────────────────────────┤
│                                                              │
│   目的：保护领域模型不受外部变化影响                           │
│                                                              │
│   ┌──────────┐     ┌──────────┐     ┌──────────┐           │
│   │ 领域层   │────▶│  防腐层   │────▶│ 外部服务  │           │
│   │          │     │ (Adapter)│     │          │           │
│   │ 领域模型 │◀────│          │◀────│ RPC/HTTP │           │
│   └──────────┘     └──────────┘     └──────────┘           │
│                                                              │
│   防腐层职责：                                                │
│   1. 协议转换：gRPC/HTTP ↔ 领域对象                           │
│   2. 数据映射：外部DTO ↔ 领域实体                             │
│   3. 错误处理：外部错误 → 领域错误                            │
│   4. 重试/熔断：网络层保护                                    │
│                                                              │
└─────────────────────────────────────────────────────────────┘
```

### 6.2 AI服务防腐层实现

```go
// service/ai/internal/acl/patient_acl.go
package acl

import (
    "context"
    
    "github.com/go-kratos/kratos/v2/log"
    patientv1 "medical-system/api/patient/v1"
    "medical-system/ai/internal/domain/aggregate"
)

// PatientACL 患者服务防腐层接口
type PatientACL interface {
    GetPatient(ctx context.Context, patientID int64) (*aggregate.Patient, error)
    GetPatientMedicalHistory(ctx context.Context, patientID int64) ([]aggregate.MedicalHistory, error)
}

// patientACL 患者服务防腐层实现
type patientACL struct {
    client patientv1.PatientServiceClient
    log    *log.Helper
}

// NewPatientACL 创建患者服务防腐层
func NewPatientACL(client patientv1.PatientServiceClient, logger log.Logger) PatientACL {
    return &patientACL{
        client: client,
        log:    log.NewHelper(logger),
    }
}

// GetPatient 获取患者信息
func (a *patientACL) GetPatient(ctx context.Context, patientID int64) (*aggregate.Patient, error) {
    resp, err := a.client.GetPatient(ctx, &patientv1.GetPatientRequest{
        Id: patientID,
    })
    if err != nil {
        a.log.Errorf("获取患者信息失败: %v", err)
        return nil, ErrPatientServiceUnavailable
    }
    
    // 转换为领域对象
    return convertProtoToPatient(resp), nil
}

// GetPatientMedicalHistory 获取患者病史
func (a *patientACL) GetPatientMedicalHistory(ctx context.Context, patientID int64) ([]aggregate.MedicalHistory, error) {
    resp, err := a.client.GetPatientMedicalHistory(ctx, &patientv1.GetPatientMedicalHistoryRequest{
        PatientId: patientID,
    })
    if err != nil {
        a.log.Errorf("获取患者病史失败: %v", err)
        return nil, err
    }
    
    // 转换为领域对象列表
    var histories []aggregate.MedicalHistory
    for _, h := range resp.Histories {
        histories = append(histories, convertProtoToMedicalHistory(h))
    }
    
    return histories, nil
}

// convertProtoToPatient 将proto转换为领域对象
func convertProtoToPatient(p *patientv1.Patient) *aggregate.Patient {
    return &aggregate.Patient{
        ID:     aggregate.PatientID(p.Id),
        OpenID: p.Openid,
        BasicInfo: vo.PatientBasicInfo{
            Name:   p.Name,
            IDCard: vo.IDCard{value: p.IdCard},
            Age:    p.Age,
            Gender: vo.Gender(p.Gender),
        },
        // ... 其他字段映射
    }
}

var (
    ErrPatientServiceUnavailable = errors.New("patient service unavailable")
)
```

### 6.3 向量数据库防腐层

```go
// service/ai/internal/acl/vector_db_acl.go
package acl

import (
    "context"
    
    "github.com/go-kratos/kratos/v2/log"
    "github.com/qdrant/go-client/qdrant"
    "medical-system/ai/internal/domain/entity"
)

// VectorDBACL 向量数据库防腐层接口
type VectorDBACL interface {
    SearchDiseases(ctx context.Context, queryVector []float32, limit int) ([]entity.VectorResult, error)
    SearchDoctors(ctx context.Context, queryVector []float32, limit int) ([]entity.VectorResult, error)
    SearchDrugs(ctx context.Context, queryVector []float32, limit int) ([]entity.VectorResult, error)
    SearchMedicalRecords(ctx context.Context, patientID int64, queryVector []float32, limit int) ([]entity.VectorResult, error)
    UpsertDocument(ctx context.Context, collection string, doc entity.VectorDocument) error
}

// vectorDBACL 向量数据库防腐层实现
type vectorDBACL struct {
    client *qdrant.Client
    log    *log.Helper
}

// NewVectorDBACL 创建向量数据库防腐层
func NewVectorDBACL(client *qdrant.Client, logger log.Logger) VectorDBACL {
    return &vectorDBACL{
        client: client,
        log:    log.NewHelper(logger),
    }
}

// SearchDiseases 搜索疾病
func (a *vectorDBACL) SearchDiseases(ctx context.Context, queryVector []float32, limit int) ([]entity.VectorResult, error) {
    return a.searchCollection(ctx, "diseases", queryVector, limit)
}

// SearchDoctors 搜索医生
func (a *vectorDBACL) SearchDoctors(ctx context.Context, queryVector []float32, limit int) ([]entity.VectorResult, error) {
    return a.searchCollection(ctx, "doctors", queryVector, limit)
}

// SearchDrugs 搜索药品
func (a *vectorDBACL) SearchDrugs(ctx context.Context, queryVector []float32, limit int) ([]entity.VectorResult, error) {
    return a.searchCollection(ctx, "drugs", queryVector, limit)
}

// SearchMedicalRecords 搜索历史病历
func (a *vectorDBACL) SearchMedicalRecords(ctx context.Context, patientID int64, queryVector []float32, limit int) ([]entity.VectorResult, error) {
    // 添加过滤条件：只搜索该患者的历史病历
    filter := &qdrant.Filter{
        Must: []*qdrant.Condition{
            {
                Condition: &qdrant.Condition_Field{
                    Field: &qdrant.FieldCondition{
                        Key: "patient_id",
                        Match: &qdrant.Match{
                            MatchValue: &qdrant.Match_Integer{
                                Integer: patientID,
                            },
                        },
                    },
                },
            },
        },
    }
    
    return a.searchWithFilter(ctx, "medical_records", queryVector, filter, limit)
}

// searchCollection 搜索集合
func (a *vectorDBACL) searchCollection(ctx context.Context, collection string, queryVector []float32, limit int) ([]entity.VectorResult, error) {
    return a.searchWithFilter(ctx, collection, queryVector, nil, limit)
}

// searchWithFilter 带过滤条件的搜索
func (a *vectorDBACL) searchWithFilter(ctx context.Context, collection string, queryVector []float32, filter *qdrant.Filter, limit int) ([]entity.VectorResult, error) {
    resp, err := a.client.Query(ctx, &qdrant.QueryPoints{
        CollectionName: collection,
        Query: &qdrant.Query{
            Variant: &qdrant.Query_Nearest{
                Nearest: &qdrant.VectorInput{
                    Variant: &qdrant.VectorInput_Vector{
                        Vector: &qdrant.Vector{
                            Data: queryVector,
                        },
                    },
                },
            },
        },
        Filter: filter,
        Limit:  uint64(limit),
        WithPayload: &qdrant.WithPayloadSelector{
            SelectorOptions: &qdrant.WithPayloadSelector_Enable{
                Enable: true,
            },
        },
    })
    if err != nil {
        a.log.Errorf("向量搜索失败 [collection=%s]: %v", collection, err)
        return nil, err
    }
    
    // 转换为领域对象
    var results []entity.VectorResult
    for _, point := range resp {
        results = append(results, entity.VectorResult{
            ID:       point.Id.GetNum(),
            Score:    point.Score,
            Payload:  convertPayload(point.Payload),
            Source:   entity.VectorSource(collection),
        })
    }
    
    return results, nil
}

// UpsertDocument 插入或更新文档
func (a *vectorDBACL) UpsertDocument(ctx context.Context, collection string, doc entity.VectorDocument) error {
    _, err := a.client.Upsert(ctx, &qdrant.UpsertPoints{
        CollectionName: collection,
        Points: []*qdrant.PointStruct{
            {
                Id: &qdrant.PointId{
                    PointIdOptions: &qdrant.PointId_Num{
                        Num: doc.ID,
                    },
                },
                Vectors: &qdrant.Vectors{
                    VectorsOptions: &qdrant.Vectors_Vector{
                        Vector: &qdrant.Vector{
                            Data: doc.Vector,
                        },
                    },
                },
                Payload: doc.Metadata,
            },
        },
    })
    
    return err
}
```

---

## 七、服务间协作

### 7.1 协作模式

```
┌─────────────────────────────────────────────────────────────────────┐
│                        服务间协作模式                                │
├─────────────────────────────────────────────────────────────────────┤
│                                                                      │
│   模式1: 同步调用 (Synchronous)                                      │
│   ═══════════════════════════════                                   │
│   使用场景：需要立即获取结果的操作                                    │
│   技术：gRPC / HTTP                                                  │
│   示例：挂号服务 → 患者服务（验证患者存在）                           │
│                                                                      │
│   ┌──────────┐      gRPC      ┌──────────┐                          │
│   │Registration│──────────────▶│ Patient  │                          │
│   │  Service   │◀──────────────│ Service  │                          │
│   └──────────┘   GetPatient   └──────────┘                          │
│                                                                      │
├─────────────────────────────────────────────────────────────────────┤
│                                                                      │
│   模式2: 异步事件 (Asynchronous)                                     │
│   ═══════════════════════════════                                   │
│   使用场景：松耦合、最终一致性                                        │
│   技术：RabbitMQ / Kafka                                             │
│   示例：病历完成 → 支付服务 + 通知服务                                │
│                                                                      │
│   ┌──────────┐      Event      ┌──────────┐                         │
│   │  Medical  │────────────────▶│ Payment  │                         │
│   │  Service  │  MedicalCompleted│ Service │                         │
│   └──────────┘                 └──────────┘                         │
│        │                                                      │      │
│        │      Event      ┌──────────┐                         │      │
│        └────────────────▶│Notification│                        │      │
│          MedicalCompleted│  Service   │                        │      │
│                          └──────────┘                        │      │
│                                                                      │
├─────────────────────────────────────────────────────────────────────┤
│                                                                      │
│   模式3: Saga分布式事务                                              │
│   ═══════════════════════════════                                   │
│   使用场景：跨服务长事务，需要补偿机制                                │
│   技术：Saga编排 / 事件溯源                                          │
│   示例：挂号 → 创建病历 → 创建支付订单                                │
│                                                                      │
│   Step1: 创建病历                                                    │
│   Step2: 创建支付订单                                                │
│   Step3: 更新挂号状态                                                │
│                                                                      │
│   失败补偿：                                                         │
│   - Step2失败 → 补偿Step1（删除病历）                                │
│   - Step3失败 → 补偿Step2（取消订单）+ 补偿Step1                     │
│                                                                      │
└─────────────────────────────────────────────────────────────────────┘
```

### 7.2 AI服务协作示例

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    AI导诊服务协作流程                                       │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│   用户输入症状                                                                │
│       │                                                                       │
│       ▼                                                                       │
│   ┌───────────────┐                                                          │
│   │  AI Service   │                                                          │
│   │  (Orchestrator)│                                                         │
│   └───────┬───────┘                                                          │
│           │                                                                   │
│     ┌─────┴─────┬─────────────┬─────────────┐                                │
│     ▼           ▼             ▼             ▼                                │
│ ┌──────┐   ┌──────┐     ┌──────────┐   ┌──────────┐                         │
│ │Patient│   │Vector│     │  LLM     │   │Disease   │                         │
│ │Service│   │ DB   │     │  API     │   │Service   │                         │
│ └──┬───┘   └──┬───┘     └────┬─────┘   └────┬─────┘                         │
│    │          │              │              │                                │
│    │  1.获取患者信息          │              │                                │
│    │─────────▶│              │              │                                │
│    │◀─────────│              │              │                                │
│    │          │              │              │                                │
│    │     2.向量检索疾病       │              │                                │
│    │          │─────────────▶│              │                                │
│    │          │◀─────────────│              │                                │
│    │          │              │              │                                │
│    │     3.向量检索医生       │              │                                │
│    │          │─────────────▶│              │                                │
│    │          │◀─────────────│              │                                │
│    │          │              │              │                                │
│    │     4.LLM推理（ReAct）   │              │                                │
│    │          │              │─────────────▶│                                │
│    │          │              │◀─────────────│                                │
│    │          │              │              │                                │
│    │     5.查询疾病详情       │              │                                │
│    │────────────────────────────────────────▶│                                │
│    │◀────────────────────────────────────────│                                │
│    │          │              │              │                                │
│    │          │              │              │                                │
│    └──────────┴──────────────┴──────────────┘                                │
│           │                                                                   │
│           ▼                                                                   │
│   组装Recommendation返回给用户                                                │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 八、DDD代码规范

### 8.1 目录结构规范

```
service/{service-name}/
├── cmd/
│   └── main.go                    # 服务入口
├── internal/
│   ├── domain/                    # 领域层 ⭐ 核心业务逻辑
│   │   ├── aggregate/             # 聚合根
│   │   ├── entity/                # 实体
│   │   ├── vo/                    # 值对象
│   │   ├── event/                 # 领域事件
│   │   ├── repository/            # 仓储接口
│   │   └── service/               # 领域服务
│   ├── biz/                       # 业务层（应用层）
│   │   ├── usecase/               # 用例
│   │   └── dto/                   # DTO
│   ├── data/                      # 基础设施层
│   │   ├── mysql/                 # MySQL实现
│   │   ├── redis/                 # Redis实现
│   │   ├── vector/                # 向量数据库实现
│   │   └── repository/            # 仓储实现
│   ├── acl/                       # 防腐层
│   │   ├── patient_acl.go
│   │   ├── doctor_acl.go
│   │   └── vector_db_acl.go
│   ├── service/                   # 应用服务（实现proto接口）
│   │   └── {service}.go
│   └── server/                    # 服务启动配置
│       ├── grpc.go
│       ├── http.go
│       └── wire.go
├── configs/                       # 配置文件
└── Dockerfile
```

### 8.2 命名规范

| 类型 | 命名规则 | 示例 |
|-----|---------|------|
| 聚合根 | {Domain} + Aggregate | PatientAggregate |
| 实体 | 名词，单数 | MedicalRecord, DiseaseItem |
| 值对象 | 描述性名词 | PatientBasicInfo, Address |
| 领域事件 | {Domain} + {Action} + Event | MedicalStatusChangedEvent |
| 领域服务 | {Domain} + DomainService | MedicalDomainService |
| 仓储接口 | {Aggregate} + Repository | PatientRepository |
| 防腐层 | {External} + ACL | PatientACL |

### 8.3 编码规范

```go
// 1. 聚合根必须包含ID和业务方法
type Patient struct {
    ID PatientID  // 必须有ID
    // ... 其他字段
}

// 业务方法，不是setter
func (p *Patient) UpdateProfile(info PatientBasicInfo) error {
    // 业务验证
    if err := info.Validate(); err != nil {
        return err
    }
    // 状态检查
    if p.Status == PatientStatusLocked {
        return ErrPatientLocked
    }
    // 执行业务逻辑
    p.BasicInfo = info
    p.UpdatedAt = time.Now()
    // 发布领域事件
    p.Events = append(p.Events, event.NewPatientProfileUpdatedEvent(p.ID))
    return nil
}

// 2. 值对象不可变，创建时验证
type IDCard struct {
    value string
}

func NewIDCard(idCard string) (IDCard, error) {
    if !isValidIDCard(idCard) {
        return IDCard{}, errors.New("invalid ID card")
    }
    return IDCard{value: idCard}, nil
}

// 3. 领域错误统一定义
var (
    ErrPatientLocked   = errors.New("patient is locked")
    ErrInvalidIDCard   = errors.New("invalid ID card")
    ErrInvalidStatusTransition = errors.New("invalid status transition")
)

// 4. 仓储接口定义在领域层
type PatientRepository interface {
    Get(ctx context.Context, id PatientID) (*Patient, error)
    Save(ctx context.Context, patient *Patient) error
    Update(ctx context.Context, patient *Patient) error
    Delete(ctx context.Context, id PatientID) error
    FindByPhone(ctx context.Context, phone string) (*Patient, error)
}

// 5. 领域事件在聚合方法中发布
func (p *Patient) SoftDelete() {
    now := time.Now()
    p.DeletedAt = &now
    p.Status = PatientStatusInactive
    // 发布领域事件
    p.Events = append(p.Events, event.NewPatientDeletedEvent(p.ID))
}
```

---

## 九、总结

### 9.1 DDD架构核心收益

```
┌─────────────────────────────────────────────────────────────┐
│                    DDD架构带来的收益                         │
├─────────────────────────────────────────────────────────────┤
│                                                              │
│  1. 业务与技术对齐                                            │
│     - 领域模型直接反映业务概念                                │
│     - 业务人员和技术人员使用统一语言                          │
│     - 需求变更可以快速定位到代码位置                          │
│                                                              │
│  2. 高内聚低耦合                                              │
│     - 每个限界上下文职责单一                                  │
│     - 核心业务逻辑集中在领域层                                │
│     - 外部变化通过防腐层隔离                                  │
│                                                              │
│  3. 可测试性                                                  │
│     - 领域对象纯内存操作，易于单元测试                        │
│     - 不依赖框架，业务逻辑独立                                │
│     - 可以快速验证业务规则                                    │
│                                                              │
│  4. 可演进性                                                  │
│     - 新功能添加不破坏现有结构                                │
│     - 领域模型随业务发展而演进                                │
│     - 微服务拆分有明确的边界指导                              │
│                                                              │
│  5. 团队并行                                                  │
│     - 不同上下文可以独立开发                                  │
│     - 上下文间通过契约（proto）协作                           │
│     - 降低团队间沟通成本                                      │
│                                                              │
└─────────────────────────────────────────────────────────────┘
```

### 9.2 实施建议

1. **从核心域开始**：优先实现Medical BC和AI BC
2. **逐步演进**：不要一次性完美设计，允许模型迭代
3. **保持简单**：不要为了DDD而DDD，避免过度设计
4. **持续重构**：定期审视领域模型，消除技术债务
5. **文档同步**：保持代码与文档的一致性

---

**文档版本**: V1.0  
**更新日期**: 2026-03-27  
**作者**: 架构组
