# 智慧医疗系统 - RAG+Agent架构设计文档

## 文档信息

| 属性 | 值 |
|-----|---|
| 版本 | V1.0 |
| 日期 | 2026-03-27 |
| 向量数据库 | Qdrant |
| Agent模式 | 混合模式（规则+Agent） |
| 目标 | RAG检索与Agent智能体架构设计 |

---

## 目录

1. [RAG架构总览](#一rag架构总览)
2. [向量数据库设计](#二向量数据库设计)
3. [RAG应用场景](#三rag应用场景)
4. [Agent架构设计](#四agent架构设计)
5. [混合推理模式](#五混合推理模式)
6. [性能优化](#六性能优化)
7. [安全保障](#七安全保障)

---

## 一、RAG架构总览

### 1.1 什么是RAG

RAG（Retrieval-Augmented Generation，检索增强生成）是将外部知识检索与大型语言模型生成能力相结合的架构。

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                          RAG架构原理                                        │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│   传统LLM                                RAG增强LLM                         │
│   ┌──────────┐                          ┌──────────┐                       │
│   │ 用户输入  │                          │ 用户输入  │                       │
│   └────┬─────┘                          └────┬─────┘                       │
│        │                                     │                              │
│        ▼                                     ▼                              │
│   ┌──────────┐                          ┌──────────┐                       │
│   │   LLM    │                          │ 向量检索  │                       │
│   │  (闭卷)   │                          │  知识库   │                       │
│   └────┬─────┘                          └────┬─────┘                       │
│        │                                     │                              │
│        ▼                                     ▼                              │
│   ┌──────────┐                          ┌──────────┐                       │
│   │  生成回答 │    ◀── 检索结果 ──       │   LLM    │                       │
│   │ (可能幻觉)│        注入上下文        │  (开卷)   │                       │
│   └──────────┘                          └────┬─────┘                       │
│                                              │                              │
│                                              ▼                              │
│                                         ┌──────────┐                       │
│                                         │ 事实性回答│                       │
│                                         │ (有依据)  │                       │
│                                         └──────────┘                       │
│                                                                              │
│   问题：知识截止、幻觉                      优势：实时知识、可溯源             │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 1.2 智慧医疗RAG架构

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                      智慧医疗系统RAG架构全景图                                │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│   接入层                                                                     │
│   ┌─────────────────────────────────────────────────────────┐                │
│   │ 微信小程序  /  Web管理端  /  医生端                      │                │
│   └──────────────────────────┬──────────────────────────────┘                │
│                              │                                               │
│   ┌──────────────────────────┴──────────────────────────────┐                │
│   │                    API Gateway                          │                │
│   │              (统一认证 / 限流 / 路由)                    │                │
│   └──────────────────────────┬──────────────────────────────┘                │
│                              │                                               │
│   应用层                      │                                               │
│   ┌──────────────────────────┴──────────────────────────────┐                │
│   │                   AI Service (ai-svc)                   │                │
│   │  ┌─────────────────────────────────────────────────┐   │                │
│   │  │              RAG Pipeline                        │   │                │
│   │  │  ┌─────────┐  ┌─────────┐  ┌─────────────────┐  │   │                │
│   │  │  │ Query   │──▶│Embedding│──▶│  Vector Search  │  │   │                │
│   │  │  │ Rewrite │  │ (BGE)   │  │   (Qdrant)      │  │   │                │
│   │  │  └─────────┘  └─────────┘  └────────┬────────┘  │   │                │
│   │  │                                      │          │   │                │
│   │  │  ┌───────────────────────────────────┘          │   │                │
│   │  │  │                                              │   │                │
│   │  │  ▼                                              │   │                │
│   │  │  ┌─────────────┐  ┌─────────────┐              │   │                │
│   │  │  │   Rerank    │──▶│  Context    │              │   │                │
│   │  │  │   (Cross)   │  │  Fusion     │              │   │                │
│   │  │  └─────────────┘  └──────┬──────┘              │   │                │
│   │  │                          │                      │   │                │
│   │  │  ┌───────────────────────┘                      │   │                │
│   │  │  │                                              │   │                │
│   │  │  ▼                                              │   │                │
│   │  │  ┌─────────────────────────────────────────┐   │   │                │
│   │  │  │              LLM Generation              │   │   │                │
│   │  │  │  (DeepSeek-V3 / GPT-4 with Context)      │   │   │                │
│   │  │  └─────────────────────────────────────────┘   │   │                │
│   │  └─────────────────────────────────────────────────┘   │                │
│   │                                                          │                │
│   │  ┌─────────────────────────────────────────────────┐   │                │
│   │  │              Agent Layer                         │   │                │
│   │  │  ┌─────────┐  ┌─────────┐  ┌─────────────────┐  │   │                │
│   │  │  │ Planner │──▶│ ReAct   │──▶│    Skills       │  │   │                │
│   │  │  │         │  │ Agent   │  │   (Tools)       │  │   │                │
│   │  │  └─────────┘  └─────────┘  └─────────────────┘  │   │                │
│   │  └─────────────────────────────────────────────────┘   │                │
│   └─────────────────────────────────────────────────────────┘                │
│                              │                                               │
│   数据层                      │                                               │
│   ┌──────────────────────────┴──────────────────────────────┐                │
│   │                                                          │                │
│   │   ┌─────────────┐  ┌─────────────┐  ┌─────────────┐    │                │
│   │   │   Qdrant    │  │    MySQL    │  │    Redis    │    │                │
│   │   │  (Vector)   │  │   (Meta)    │  │   (Cache)   │    │                │
│   │   │             │  │             │  │             │    │                │
│   │   │ • diseases  │  │ • patients  │  │ • sessions  │    │                │
│   │   │ • doctors   │  │ • records   │  │ • vectors   │    │                │
│   │   │ • drugs     │  │ • logs      │  │ • rate_limit│    │                │
│   │   │ • records   │  │             │  │             │    │                │
│   │   └─────────────┘  └─────────────┘  └─────────────┘    │                │
│   │                                                          │                │
│   └─────────────────────────────────────────────────────────┘                │
│                                                                              │
│   外部服务                                                                   │
│   ┌──────────────────────────┬──────────────────────────────┐                │
│   │   SiliconFlow API         │   DeepSeek API             │                │
│   │   (Embedding)             │   (LLM)                    │                │
│   └──────────────────────────┴──────────────────────────────┘                │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 1.3 RAG数据流

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                          RAG数据流图                                        │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│   离线流程（数据准备）                                                        │
│   ═══════════════════                                                         │
│                                                                              │
│   ┌──────────┐    ┌──────────┐    ┌──────────┐    ┌──────────┐             │
│   │  原始数据 │───▶│  文本分块 │───▶│  向量化  │───▶│ 存入向量库│             │
│   │          │    │  (Chunk) │    │ (BGE)   │    │ (Qdrant)│             │
│   └──────────┘    └──────────┘    └──────────┘    └──────────┘             │
│        │                                               │                    │
│   疾病知识库                                           │                    │
│   医生介绍文档                                         │                    │
│   药品说明书      ─────────────────────────────────────┘                    │
│   历史病历                                                                   │
│                                                                              │
│   ┌─────────────────────────────────────────────────────────────────────┐   │
│   │                      Qdrant Collections                              │   │
│   │  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐  ┌────────────┐  │   │
│   │  │  diseases   │  │   doctors   │  │    drugs    │  │  records   │  │   │
│   │  │             │  │             │  │             │  │            │  │   │
│   │  │ id: int64   │  │ id: int64   │  │ id: int64   │  │ id: int64  │  │   │
│   │  │ vector: []  │  │ vector: []  │  │ vector: []  │  │ vector: [] │  │   │
│   │  │ payload: {} │  │ payload: {} │  │ payload: {} │  │ payload: {}│  │   │
│   │  └─────────────┘  └─────────────┘  └─────────────┘  └────────────┘  │   │
│   └─────────────────────────────────────────────────────────────────────┘   │
│                                                                              │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│   在线流程（查询推理）                                                        │
│   ═══════════════════                                                         │
│                                                                              │
│   ┌──────────┐    ┌──────────┐    ┌──────────┐    ┌──────────┐             │
│   │ 用户输入  │───▶│ Query    │───▶│  向量检索 │───▶│ 结果排序 │             │
│   │ "头痛发热"│    │ Rewrite  │    │  (TopK)  │    │ (Rerank) │             │
│   └──────────┘    └──────────┘    └──────────┘    └──────────┘             │
│                                        │                    │               │
│                              ┌─────────┴─────────┐         │               │
│                              ▼                   ▼         │               │
│                        ┌──────────┐        ┌──────────┐    │               │
│                        │ diseases │        │  doctors │    │               │
│                        │   集合   │        │   集合   │    │               │
│                        └──────────┘        └──────────┘    │               │
│                                                            │               │
│   ┌──────────┐    ┌──────────┐    ┌───────────────────────┘               │
│   │  LLM     │◀───│  Prompt  │◀───│  Context Fusion                          │
│   │ 生成回答  │    │ 组装    │    │  (检索结果+用户输入)                      │
│   └──────────┘    └──────────┘    └───────────────────────┘                │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 二、向量数据库设计

### 2.1 Qdrant选型理由

| 特性 | Qdrant | Milvus | Pinecone | Weaviate |
|-----|--------|--------|----------|----------|
| 部署复杂度 | ⭐⭐ 低 | ⭐⭐⭐⭐ 高 | ⭐ 托管 | ⭐⭐⭐ 中 |
| Go客户端 | ⭐⭐⭐⭐ 原生 | ⭐⭐⭐ 第三方 | ⭐⭐ 有限 | ⭐⭐⭐ 一般 |
| 性能 | ⭐⭐⭐⭐ 高 | ⭐⭐⭐⭐ 高 | ⭐⭐⭐ 中 | ⭐⭐⭐ 中 |
| 功能丰富度 | ⭐⭐⭐⭐ 高 | ⭐⭐⭐⭐⭐ 很高 | ⭐⭐⭐ 中 | ⭐⭐⭐⭐ 高 |
| 内存占用 | ⭐⭐⭐⭐ 低 | ⭐⭐⭐ 中 | ⭐⭐⭐⭐ 低 | ⭐⭐⭐ 中 |
| 过滤能力 | ⭐⭐⭐⭐ 强 | ⭐⭐⭐⭐⭐ 很强 | ⭐⭐⭐ 中 | ⭐⭐⭐⭐ 强 |

**选择Qdrant的理由**：
1. **Go原生支持**：官方Go客户端，集成简单
2. **轻量部署**：单容器即可运行，适合开发阶段
3. **高性能**：基于Rust实现，查询延迟低
4. **功能丰富**：支持过滤、分页、批量操作
5. **易于扩展**：后期可迁移到Milvus分布式集群

### 2.2 Collections设计

```go
// Qdrant集合定义

// Collection: diseases
{
    "collection_name": "diseases",
    "vector_size": 1024,  // BGE-large-zh embedding维度
    "distance": "Cosine",
    "payload_schema": {
        "name": {"type": "keyword"},
        "symptoms": {"type": "text"},
        "diagnosis": {"type": "text"},
        "treatment_plan": {"type": "text"},
        "prevention": {"type": "text"},
        "drug_category": {"type": "keyword"},
        "department": {"type": "keyword"},
        "severity": {"type": "integer"}
    }
}

// Collection: doctors
{
    "collection_name": "doctors",
    "vector_size": 1024,
    "distance": "Cosine",
    "payload_schema": {
        "name": {"type": "keyword"},
        "title": {"type": "keyword"},        // 职称
        "position": {"type": "keyword"},     // 职务
        "specialty": {"type": "text"},       // 专业领域
        "intro": {"type": "text"},           // 个人介绍
        "is_expert": {"type": "bool"},
        "department_id": {"type": "integer"},
        "department_name": {"type": "keyword"},
        "status": {"type": "integer"}        // 工作状态
    }
}

// Collection: drugs
{
    "collection_name": "drugs",
    "vector_size": 1024,
    "distance": "Cosine",
    "payload_schema": {
        "name": {"type": "keyword"},
        "description": {"type": "text"},
        "category": {"type": "keyword"},      // 药品分类
        "category_id": {"type": "integer"},
        "usage": {"type": "text"},           // 用法
        "contraindication": {"type": "text"}, // 禁忌
        "side_effects": {"type": "text"},     // 副作用
        "price": {"type": "float"},
        "in_stock": {"type": "bool"}
    }
}

// Collection: medical_records
{
    "collection_name": "medical_records",
    "vector_size": 1024,
    "distance": "Cosine",
    "payload_schema": {
        "patient_id": {"type": "integer"},
        "medical_number": {"type": "keyword"},
        "symptom": {"type": "text"},
        "diagnosis": {"type": "text"},
        "treatment_plan": {"type": "text"},
        "department_id": {"type": "integer"},
        "doctor_id": {"type": "integer"},
        "status": {"type": "integer"},
        "create_time": {"type": "datetime"}
    }
}
```

### 2.3 索引策略

```go
// 索引配置优化

// 1. HNSW索引（近似最近邻搜索）
hnswConfig := &qdrant.HnswConfigDiff{
    M:              16,     // 每个节点的最大连接数
    EfConstruct:    100,    // 构建时的搜索深度
    FullScanThreshold: 10000, // 超过此数量启用全扫描
    OnDisk:         false,   // 是否存储在磁盘
}

// 2. 标量索引（用于过滤）
scalarConfig := &qdrant.ScalarQuantization{
    Type:      qdrant.QuantizationType_Int8,
    Quantile:  0.99,
    AlwaysRam: true,
}

// 3. Payload索引（用于过滤条件）
payloadIndex := map[string]*qdrant.PayloadIndexParams{
    "department_id": {
        IndexParams: &qdrant.PayloadIndexParams_FieldIndexParams{
            FieldIndexParams: &qdrant.FieldIndexParams{
                Type: &qdrant.FieldIndexParams_Integer{
                    Integer: &qdrant.IntegerIndexParams{
                        Lookup: true,
                        Range:  true,
                    },
                },
            },
        },
    },
    "is_expert": {
        IndexParams: &qdrant.PayloadIndexParams_FieldIndexParams{
            FieldIndexParams: &qdrant.FieldIndexParams{
                Type: &qdrant.FieldIndexParams_Bool{
                    Bool: &qdrant.BoolIndexParams{
                        OnDisk: false,
                    },
                },
            },
        },
    },
}
```

### 2.4 数据向量化流程

```go
// 数据向量化管道

// 步骤1：文本构造
type TextBuilder struct {
    Template map[string]string
}

func (b *TextBuilder) BuildDiseaseText(d Disease) string {
    return fmt.Sprintf("疾病名称：%s\n"+
        "症状：%s\n"+
        "诊断：%s\n"+
        "治疗方案：%s\n"+
        "预防措施：%s\n"+
        "适用科室：%s",
        d.Name, d.Symptoms, d.Diagnosis, 
        d.TreatmentPlan, d.Prevention, d.Department)
}

func (b *TextBuilder) BuildDoctorText(d Doctor) string {
    return fmt.Sprintf("医生姓名：%s\n"+
        "职称：%s\n"+
        "专业领域：%s\n"+
        "个人介绍：%s\n"+
        "是否专家：%v",
        d.Name, d.Title, d.Specialty, d.Intro, d.IsExpert)
}

// 步骤2：向量化
type EmbeddingService struct {
    client *siliconflow.Client
    model  string
}

func (s *EmbeddingService) EmbedBatch(ctx context.Context, texts []string) ([][]float32, error) {
    req := &siliconflow.EmbeddingRequest{
        Model: s.model,  // BAAI/bge-large-zh-v1.5
        Input: texts,
    }
    
    resp, err := s.client.CreateEmbeddings(ctx, req)
    if err != nil {
        return nil, err
    }
    
    embeddings := make([][]float32, len(resp.Data))
    for i, data := range resp.Data {
        embeddings[i] = data.Embedding
    }
    
    return embeddings, nil
}

// 步骤3：批量索引
type Indexer struct {
    vectorRepo VectorRepository
    embedder   *EmbeddingService
    batchSize  int
}

func (i *Indexer) IndexDiseases(ctx context.Context, diseases []Disease) error {
    // 分批处理
    for batch := 0; batch < len(diseases); batch += i.batchSize {
        end := batch + i.batchSize
        if end > len(diseases) {
            end = len(diseases)
        }
        
        batchData := diseases[batch:end]
        
        // 构造文本
        texts := make([]string, len(batchData))
        for i, d := range batchData {
            texts[i] = i.BuildDiseaseText(d)
        }
        
        // 向量化
        embeddings, err := i.embedder.EmbedBatch(ctx, texts)
        if err != nil {
            return err
        }
        
        // 批量插入Qdrant
        docs := make([]VectorDocument, len(batchData))
        for i, d := range batchData {
            docs[i] = VectorDocument{
                ID:       int64(d.ID),
                Vector:   embeddings[i],
                Metadata: map[string]string{
                    "name":             d.Name,
                    "symptoms":         d.Symptoms,
                    "diagnosis":        d.Diagnosis,
                    "treatment_plan":   d.TreatmentPlan,
                    "department":       d.Department,
                },
            }
        }
        
        if err := i.vectorRepo.BatchUpsert(ctx, "diseases", docs); err != nil {
            return err
        }
    }
    
    return nil
}
```

---

## 三、RAG应用场景

### 3.1 疾病症状检索

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                        疾病症状检索场景                                       │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│   用户输入："最近头痛、发热、咳嗽，还有点流鼻涕"                              │
│                                                                              │
│   ┌─────────────────────────────────────────────────────────────────────┐   │
│   │                        处理流程                                      │   │
│   ├─────────────────────────────────────────────────────────────────────┤   │
│   │                                                                      │   │
│   │  1. Query预处理                                                      │   │
│   │     - 关键词提取：头痛、发热、咳嗽、流鼻涕                            │   │
│   │     - 意图识别：症状描述                                             │   │
│   │                                                                      │   │
│   │  2. 向量化                                                           │   │
│   │     - 输入："症状：头痛、发热、咳嗽、流鼻涕"                          │   │
│   │     - 输出：1024维向量                                               │   │
│   │                                                                      │   │
│   │  3. 向量检索（Qdrant）                                               │   │
│   │     - 集合：diseases                                                │   │
│   │     - TopK：5                                                       │   │
│   │     - 返回：                                                         │   │
│   │       ┌─────────────────────────────────────────────────────────┐   │   │
│   │       │ 1. 流感 (score: 0.89)                                   │   │   │
│   │       │ 2. 上呼吸道感染 (score: 0.85)                           │   │   │
│   │       │ 3. 普通感冒 (score: 0.82)                               │   │   │
│   │       │ 4. 急性支气管炎 (score: 0.78)                           │   │   │
│   │       │ 5. 肺炎 (score: 0.75)                                   │   │   │
│   │       └─────────────────────────────────────────────────────────┘   │   │
│   │                                                                      │   │
│   │  4. 结果处理                                                         │   │
│   │     - 过滤：score > 0.70                                            │   │
│   │     - 排序：按score降序                                             │   │
│   │     - 去重：相同科室疾病合并                                        │   │
│   │                                                                      │   │
│   │  5. LLM增强（可选）                                                  │   │
│   │     - 将检索结果输入LLM                                            │   │
│   │     - 生成自然语言解释                                             │   │
│   │                                                                      │   │
│   └─────────────────────────────────────────────────────────────────────┘   │
│                                                                              │
│   输出结果：                                                                 │
│   ┌─────────────────────────────────────────────────────────────────────┐   │
│   │ 可能疾病：                                                          │   │
│   │ 1. 流感 (匹配度: 89%)                                               │   │
│   │    - 症状：高热、头痛、全身酸痛、咳嗽、流涕                         │   │
│   │    - 建议：发热门诊就诊，注意休息，多喝水                           │   │
│   │                                                                    │   │
│   │ 2. 上呼吸道感染 (匹配度: 85%)                                       │   │
│   │    - 症状：鼻塞、流涕、咽痛、咳嗽、低热                             │   │
│   │    - 建议：内科或呼吸科就诊                                         │   │
│   │                                                                    │   │
│   │ 3. 普通感冒 (匹配度: 82%)                                           │   │
│   │    - 症状：鼻塞、流涕、打喷嚏、轻微发热                             │   │
│   │    - 建议：多休息，症状加重及时就医                                 │   │
│   └─────────────────────────────────────────────────────────────────────┘   │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 3.2 医生智能推荐

```go
// 医生推荐算法

type DoctorRecommender struct {
    vectorRepo VectorRepository
    doctorRepo DoctorRepository
}

// Recommend 推荐医生
func (r *DoctorRecommender) Recommend(ctx context.Context, symptoms string, preferredDept int64, needExpert bool) ([]DoctorRecommendation, error) {
    // 1. 向量化症状描述
    queryVector, err := r.embedder.Embed(ctx, symptoms)
    if err != nil {
        return nil, err
    }
    
    // 2. 构建过滤条件
    filter := &qdrant.Filter{}
    
    // 2.1 状态过滤：只推荐工作中的医生
    filter.Must = append(filter.Must, &qdrant.Condition{
        Condition: &qdrant.Condition_Field{
            Field: &qdrant.FieldCondition{
                Key: "status",
                Match: &qdrant.Match{
                    MatchValue: &qdrant.Match_Integer{
                        Integer: 1, // 工作中
                    },
                },
            },
        },
    })
    
    // 2.2 专家过滤
    if needExpert {
        filter.Must = append(filter.Must, &qdrant.Condition{
            Condition: &qdrant.Condition_Field{
                Field: &qdrant.FieldCondition{
                    Key: "is_expert",
                    Match: &qdrant.Match{
                        MatchValue: &qdrant.Match_Bool{
                            Bool: true,
                        },
                    },
                },
            },
        })
    }
    
    // 2.3 科室偏好
    if preferredDept > 0 {
        filter.Should = append(filter.Should, &qdrant.Condition{
            Condition: &qdrant.Condition_Field{
                Field: &qdrant.FieldCondition{
                    Key: "department_id",
                    Match: &qdrant.Match{
                        MatchValue: &qdrant.Match_Integer{
                            Integer: preferredDept,
                        },
                    },
                },
            },
        })
    }
    
    // 3. 向量检索
    results, err := r.vectorRepo.SearchWithFilter(ctx, "doctors", queryVector, filter, 10)
    if err != nil {
        return nil, err
    }
    
    // 4. 获取详细信息并计算综合评分
    var recommendations []DoctorRecommendation
    for _, result := range results {
        doctor, err := r.doctorRepo.GetByID(ctx, result.ID)
        if err != nil {
            continue
        }
        
        // 计算综合评分
        score := r.calculateScore(result, doctor, preferredDept)
        
        recommendations = append(recommendations, DoctorRecommendation{
            DoctorID:     doctor.ID,
            DoctorName:   doctor.Name,
            Title:        doctor.Title,
            IsExpert:     doctor.IsExpert,
            DepartmentID: doctor.DepartmentID,
            MatchScore:   score,
            Reason:       r.generateReason(doctor, result.Score),
        })
    }
    
    // 5. 排序并返回Top4
    sort.Slice(recommendations, func(i, j int) bool {
        return recommendations[i].MatchScore > recommendations[j].MatchScore
    })
    
    if len(recommendations) > 4 {
        recommendations = recommendations[:4]
    }
    
    return recommendations, nil
}

// calculateScore 计算综合评分
func (r *DoctorRecommender) calculateScore(vectorResult VectorResult, doctor Doctor, preferredDept int64) float64 {
    score := vectorResult.Score * 0.6  // 向量相似度占60%
    
    // 专家加分
    if doctor.IsExpert {
        score += 0.15
    }
    
    // 职称加分
    switch doctor.Position {
    case PositionChiefPhysician:
        score += 0.1
    case PositionAssociateChief:
        score += 0.08
    case PositionAttending:
        score += 0.05
    }
    
    // 科室匹配加分
    if preferredDept > 0 && doctor.DepartmentID == preferredDept {
        score += 0.1
    }
    
    return min(score, 1.0)
}

// generateReason 生成推荐理由
func (r *DoctorRecommender) generateReason(doctor Doctor, matchScore float64) string {
    reasons := []string{}
    
    if matchScore > 0.8 {
        reasons = append(reasons, "专业领域高度匹配")
    }
    if doctor.IsExpert {
        reasons = append(reasons, "专家医师")
    }
    if doctor.Position <= PositionAssociateChief {
        reasons = append(reasons, "高级职称")
    }
    
    if len(reasons) == 0 {
        return "专业对口"
    }
    
    return strings.Join(reasons, "，")
}
```

### 3.3 相似病历检索

```go
// 相似病历检索
type SimilarMedicalFinder struct {
    vectorRepo VectorRepository
    medicalRepo MedicalRepository
}

// FindSimilar 查找相似病历
func (f *SimilarMedicalFinder) FindSimilar(ctx context.Context, patientID int64, currentSymptoms string, limit int) ([]SimilarMedicalRecord, error) {
    // 1. 当前症状向量化
    queryVector, err := f.embedder.Embed(ctx, currentSymptoms)
    if err != nil {
        return nil, err
    }
    
    // 2. 构建过滤：只搜索该患者的历史病历
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
    
    // 3. 向量检索
    results, err := f.vectorRepo.SearchWithFilter(ctx, "medical_records", queryVector, filter, limit)
    if err != nil {
        return nil, err
    }
    
    // 4. 获取病历详情
    var similarRecords []SimilarMedicalRecord
    for _, result := range results {
        record, err := f.medicalRepo.GetByID(ctx, result.ID)
        if err != nil {
            continue
        }
        
        similarRecords = append(similarRecords, SimilarMedicalRecord{
            MedicalID:     record.ID,
            MedicalNumber: record.MedicalNumber,
            DiagnoseDate:  record.CreatedAt,
            Diagnosis:     record.Diagnosis,
            Treatment:     record.TreatmentPlan,
            Similarity:    result.Score,
            Reference:     r.generateReference(record),
        })
    }
    
    return similarRecords, nil
}

// generateReference 生成参考建议
func (f *SimilarMedicalFinder) generateReference(record MedicalRecord) string {
    return fmt.Sprintf("该患者曾在%s就诊，诊断为%s，治疗方案为%s。",
        record.CreatedAt.Format("YYYY-MM-DD"),
        record.Diagnosis,
        record.TreatmentPlan)
}
```

### 3.4 药品知识检索

```go
// 药品知识检索
type DrugKnowledgeRetriever struct {
    vectorRepo VectorRepository
    drugRepo   DrugRepository
}

// Retrieve 检索药品知识
func (r *DrugKnowledgeRetriever) Retrieve(ctx context.Context, query string, patientAllergies []string) ([]DrugKnowledge, error) {
    // 1. 向量化查询
    queryVector, err := r.embedder.Embed(ctx, query)
    if err != nil {
        return nil, err
    }
    
    // 2. 向量检索
    results, err := r.vectorRepo.Search(ctx, "drugs", queryVector, 10)
    if err != nil {
        return nil, err
    }
    
    // 3. 过滤过敏药品
    var knowledgeList []DrugKnowledge
    for _, result := range results {
        drug, err := r.drugRepo.GetByID(ctx, result.ID)
        if err != nil {
            continue
        }
        
        // 检查过敏
        if r.hasAllergy(drug, patientAllergies) {
            continue
        }
        
        knowledgeList = append(knowledgeList, DrugKnowledge{
            DrugID:      drug.ID,
            DrugName:    drug.Name,
            Description: drug.Description,
            Category:    drug.Category,
            Usage:       drug.Usage,
            Contraindication: drug.Contraindication,
            SideEffects: drug.SideEffects,
            MatchScore:  result.Score,
        })
    }
    
    return knowledgeList, nil
}

// hasAllergy 检查过敏
func (r *DrugKnowledgeRetriever) hasAllergy(drug Drug, allergies []string) bool {
    for _, allergy := range allergies {
        if strings.Contains(drug.Name, allergy) || 
           strings.Contains(drug.Contraindication, allergy) {
            return true
        }
    }
    return false
}
```

---

## 四、Agent架构设计

### 4.1 Agent模式对比

| 模式 | 特点 | 适用场景 | 延迟 | 复杂度 |
|-----|------|---------|------|--------|
| **Rule-based** | 纯规则，无AI | 确定性任务 | 极低 | 低 |
| **Skills** | 工具调用，单步执行 | 简单功能 | 低 | 低 |
| **ReAct** | 思考-行动-观察循环 | 多步推理 | 高 | 中 |
| **Plan-and-Execute** | 先规划后执行 | 复杂任务 | 高 | 高 |
| **混合模式** | 规则+Agent切换 | 综合场景 | 可变 | 中 |

### 4.2 混合Agent架构

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                        混合Agent架构                                         │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│   ┌─────────────────────────────────────────────────────────────────────┐   │
│   │                      Router (路由层)                                 │   │
│   │  根据输入特征决定使用哪种模式                                          │   │
│   ├─────────────────────────────────────────────────────────────────────┤   │
│   │                                                                      │   │
│   │  IF 简单症状 + 明确意图                                               │   │
│   │      THEN 规则模式 (快速响应)                                         │   │
│   │                                                                      │   │
│   │  IF 复杂症状 + 多轮对话 + 需要推理                                     │   │
│   │      THEN Agent模式 (智能推理)                                        │   │
│   │                                                                      │   │
│   │  IF 紧急情况 / 高危症状                                                │   │
│   │      THEN 规则兜底 + Agent增强                                        │   │
│   │                                                                      │   │
│   └────────────────┬────────────────────┬───────────────────────────────┘   │
│                    │                    │                                   │
│                    ▼                    ▼                                   │
│   ┌───────────────────────┐    ┌───────────────────────┐                   │
│   │     规则模式          │    │      Agent模式         │                   │
│   │  ═══════════════════  │    │  ═══════════════════  │                   │
│   │                      │    │                      │                   │
│   │  ┌───────────────┐   │    │  ┌───────────────┐   │                   │
│   │  │ 关键词匹配   │   │    │  │   Planner     │   │                   │
│   │  │ 症状-疾病映射│   │    │  │  (规划器)      │   │                   │
│   │  └───────┬───────┘   │    │  └───────┬───────┘   │                   │
│   │          │           │    │          │           │                   │
│   │  ┌───────┴───────┐   │    │  ┌───────┴───────┐   │                   │
│   │  │ 规则引擎     │   │    │  │   ReAct      │   │                   │
│   │  │ (Drools/Go)│   │    │  │   Agent       │   │                   │
│   │  └───────┬───────┘   │    │  └───────┬───────┘   │                   │
│   │          │           │    │          │           │                   │
│   │  ┌───────┴───────┐   │    │  ┌───────┴───────┐   │                   │
│   │  │ 向量检索     │   │    │  │    Skills     │   │                   │
│   │  │ (TopK固定)  │   │    │  │   (工具调用)  │   │                   │
│   │  └───────────────┘   │    │  └───────────────┘   │                   │
│   │                      │    │                      │                   │
│   │  特点：快速、可控     │    │  特点：智能、灵活    │                   │
│   │  延迟：<100ms         │    │  延迟：1-3s         │                   │
│   │                      │    │                      │                   │
│   └──────────┬───────────┘    └──────────┬───────────┘                   │
│              │                           │                                 │
│              └───────────┬───────────────┘                                 │
│                          ▼                                                 │
│   ┌─────────────────────────────────────────────────────────────────────┐   │
│   │                      Response Fusion (结果融合)                      │   │
│   │  合并规则结果和Agent结果，确保一致性和完整性                          │   │
│   └─────────────────────────────────────────────────────────────────────┘   │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 4.3 ReAct Agent实现

```go
// ReAct Agent实现

// ReActPrompt ReAct提示词模板
const ReActPrompt = `你是一个专业的医疗导诊助手。请根据患者的症状描述，使用以下工具来帮助诊断：

可用工具：
{{tools}}

请按照以下格式进行思考：

思考：我需要分析患者的症状，并搜索相关疾病和医生
行动：search_disease
输入：{"query": "头痛发热咳嗽", "limit": 5}
观察：找到3种可能疾病：1.流感 2.上呼吸道感染 3.普通感冒

思考：我需要根据可能的疾病搜索合适的医生
行动：search_doctor
输入：{"query": "流感 呼吸内科", "limit": 4}
观察：找到4位医生，包括2位专家

思考：我已经收集了足够的信息，可以给出推荐
最终答案：基于您的症状（头痛、发热、咳嗽），最可能的诊断是流感或上呼吸道感染。建议就诊科室：发热门诊或呼吸内科。推荐医生：...

当前患者症状：
{{context}}

请开始思考：`

// ReActAgent ReAct Agent
type ReActAgent struct {
    llmClient LLMClient
    skills    map[string]Skill
    maxSteps  int
}

// Execute 执行ReAct循环
func (a *ReActAgent) Execute(ctx context.Context, task Task) (Result, error) {
    var result Result
    var steps []Step
    
    context := task.Context
    
    for stepNum := 1; stepNum <= a.maxSteps; stepNum++ {
        // 构建提示词
        prompt := a.buildPrompt(context, steps)
        
        // LLM思考
        response, err := a.llmClient.Complete(ctx, prompt)
        if err != nil {
            return result, err
        }
        
        // 解析响应
        thought, action, input, isFinal := a.parseResponse(response)
        
        step := Step{
            Number:  stepNum,
            Thought: thought,
            Action:  action,
            Input:   input,
            IsFinal: isFinal,
        }
        
        if isFinal {
            step.Output = action
            steps = append(steps, step)
            result.Answer = action
            break
        }
        
        // 执行工具
        output, err := a.executeTool(ctx, action, input)
        if err != nil {
            step.Output = fmt.Sprintf("Error: %v", err)
        } else {
            step.Output = output
        }
        
        steps = append(steps, step)
        
        // 更新上下文
        context = fmt.Sprintf("%s\n\n步骤%d:\n思考：%s\n行动：%s\n输入：%s\n观察：%s",
            context, stepNum, thought, action, input, step.Output)
    }
    
    result.Steps = steps
    result.Recommendation = a.generateRecommendation(ctx, context, steps)
    
    return result, nil
}

// parseResponse 解析LLM响应
func (a *ReActAgent) parseResponse(response string) (thought, action, input string, isFinal bool) {
    lines := strings.Split(response, "\n")
    
    for _, line := range lines {
        if strings.HasPrefix(line, "思考：") {
            thought = strings.TrimPrefix(line, "思考：")
        }
        if strings.HasPrefix(line, "行动：") {
            action = strings.TrimPrefix(line, "行动：")
        }
        if strings.HasPrefix(line, "输入：") {
            input = strings.TrimPrefix(line, "输入：")
        }
        if strings.HasPrefix(line, "最终答案：") {
            return "", strings.TrimPrefix(line, "最终答案："), "", true
        }
    }
    
    return thought, action, input, false
}

// executeTool 执行工具
func (a *ReActAgent) executeTool(ctx context.Context, name, input string) (string, error) {
    tool, exists := a.skills[name]
    if !exists {
        return "", fmt.Errorf("tool not found: %s", name)
    }
    
    var params map[string]interface{}
    if err := json.Unmarshal([]byte(input), &params); err != nil {
        return "", err
    }
    
    return tool.Execute(ctx, params)
}
```

### 4.4 Skills工具定义

```go
// Skills定义

// Skill 工具接口
type Skill interface {
    Name() string
    Description() string
    Execute(ctx context.Context, params map[string]interface{}) (string, error)
}

// 1. 症状收集Skill
var SymptomCollectorSkill = &SkillDefinition{
    Name:        "collect_symptom",
    Description: "收集患者症状信息，包括身体部位、症状描述、严重程度",
    Parameters: map[string]Parameter{
        "body_part": {
            Type:        "string",
            Description: "身体部位，如头部、胸部、腹部等",
            Required:    true,
        },
        "description": {
            Type:        "string",
            Description: "症状详细描述",
            Required:    true,
        },
        "severity": {
            Type:        "integer",
            Description: "严重程度：1轻微 2中等 3严重 4危急",
            Required:    true,
        },
    },
}

// 2. 疾病搜索Skill
var DiseaseSearchSkill = &SkillDefinition{
    Name:        "search_disease",
    Description: "根据症状描述搜索相关疾病",
    Parameters: map[string]Parameter{
        "query": {
            Type:        "string",
            Description: "症状描述或疾病相关关键词",
            Required:    true,
        },
        "limit": {
            Type:        "integer",
            Description: "返回结果数量，默认5",
            Required:    false,
        },
    },
}

// 3. 医生搜索Skill
var DoctorSearchSkill = &SkillDefinition{
    Name:        "search_doctor",
    Description: "根据疾病或科室搜索合适的医生",
    Parameters: map[string]Parameter{
        "query": {
            Type:        "string",
            Description: "疾病名称或专业领域",
            Required:    true,
        },
        "need_expert": {
            Type:        "boolean",
            Description: "是否只需要专家医生",
            Required:    false,
        },
        "limit": {
            Type:        "integer",
            Description: "返回结果数量，默认4",
            Required:    false,
        },
    },
}

// 4. 药品搜索Skill
var DrugSearchSkill = &SkillDefinition{
    Name:        "search_drug",
    Description: "根据疾病或症状搜索适用药品",
    Parameters: map[string]Parameter{
        "query": {
            Type:        "string",
            Description: "疾病名称或症状描述",
            Required:    true,
        },
        "limit": {
            Type:        "integer",
            Description: "返回结果数量，默认4",
            Required:    false,
        },
    },
}

// 5. 病历查询Skill
var MedicalHistorySkill = &SkillDefinition{
    Name:        "query_medical_history",
    Description: "查询患者历史病历",
    Parameters: map[string]Parameter{
        "patient_id": {
            Type:        "integer",
            Description: "患者ID",
            Required:    true,
        },
        "query": {
            Type:        "string",
            Description: "当前症状描述，用于匹配历史病历",
            Required:    true,
        },
    },
}
```

---

## 五、混合推理模式

### 5.1 模式切换策略

```go
// 模式路由器
type Router struct {
    ruleEngine    *RuleEngine
    reactAgent    *ReActAgent
    hybridConfig  HybridConfig
}

// Route 路由决策
func (r *Router) Route(ctx context.Context, input UserInput) (Mode, error) {
    // 1. 紧急程度检测
    if r.isEmergency(input) {
        return ModeRuleBased, nil // 紧急情况下使用规则模式，快速响应
    }
    
    // 2. 复杂度评估
    complexity := r.assessComplexity(input)
    
    // 3. 根据复杂度选择模式
    switch complexity {
    case ComplexitySimple:
        // 简单场景：直接规则匹配
        return ModeRuleBased, nil
    case ComplexityMedium:
        // 中等复杂度：规则+向量检索
        return ModeHybrid, nil
    case ComplexityComplex:
        // 复杂场景：启用Agent推理
        return ModeAgent, nil
    default:
        return ModeHybrid, nil
    }
}

// assessComplexity 评估复杂度
func (r *Router) assessComplexity(input UserInput) Complexity {
    score := 0
    
    // 症状数量
    if len(input.Symptoms) > 3 {
        score += 2
    }
    
    // 症状持续时间
    if input.Duration > 7*24*time.Hour { // 超过7天
        score += 2
    }
    
    // 是否多系统症状
    if r.hasMultiSystemSymptoms(input.Symptoms) {
        score += 3
    }
    
    // 是否提及多个疾病
    if r.mentionsDiseases(input.Text) {
        score += 1
    }
    
    // 对话轮数
    if input.ConversationTurns > 3 {
        score += 2
    }
    
    // 根据分数判断复杂度
    if score >= 6 {
        return ComplexityComplex
    } else if score >= 3 {
        return ComplexityMedium
    }
    return ComplexitySimple
}
```

### 5.2 混合推理流程

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                       混合推理执行流程                                       │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│   用户输入："头痛3天了，还有点恶心，昨天量体温37.8度"                         │
│                                                                              │
│   Step 1: 模式识别                                                           │
│   ┌─────────────────────────────────────────────────────────────────────┐   │
│   │ 复杂度评估：                                                        │   │
│   │ - 症状数量：3 (头痛、恶心、发热) → score +2                          │   │
│   │ - 持续时间：3天 → score +2                                           │   │
│   │ - 多系统：是 (神经+消化+全身) → score +3                              │   │
│   │ 总分：7 → 复杂度：Complex                                            │   │
│   │ 选择模式：Agent模式                                                 │   │
│   └─────────────────────────────────────────────────────────────────────┘   │
│                                                                              │
│   Step 2: Agent执行                                                          │
│   ┌─────────────────────────────────────────────────────────────────────┐   │
│   │                                                                      │   │
│   │  思考：患者症状包括头痛、恶心、低热，持续3天                          │   │
│   │        可能是流感、脑膜炎、偏头痛等                                   │   │
│   │        需要搜索相关疾病和医生                                         │   │
│   │                                                                      │   │
│   │  行动：search_disease                                                │   │
│   │  输入：{"query": "头痛 恶心 低热 3天", "limit": 5}                    │   │
│   │                                                                      │   │
│   │  观察：[{"name": "流感", "score": 0.85},                                │   │
│   │        {"name": "病毒性脑膜炎", "score": 0.78},                         │   │
│   │        {"name": "偏头痛", "score": 0.72}]                              │   │
│   │                                                                      │   │
│   │  思考：找到可能的疾病，需要进一步搜索医生                             │   │
│   │                                                                      │   │
│   │  行动：search_doctor                                                 │   │
│   │  输入：{"query": "流感 神经内科", "limit": 4}                          │   │
│   │                                                                      │   │
│   │  观察：[{"name": "张医生", "title": "主任医师", ...}]                   │   │
│   │                                                                      │   │
│   │  思考：已收集足够信息，生成推荐                                       │   │
│   │                                                                      │   │
│   │  最终答案：基于您的症状（头痛、恶心、低热持续3天），                   │   │
│   │          最可能的诊断是流感或病毒性脑膜炎...                           │   │
│   │                                                                      │   │
│   └─────────────────────────────────────────────────────────────────────┘   │
│                                                                              │
│   Step 3: 结果增强                                                           │
│   ┌─────────────────────────────────────────────────────────────────────┐   │
│   │  使用规则验证Agent结果：                                            │   │
│   │  - 检查是否有紧急情况（无）                                         │   │
│   │  - 补充科室信息（神经内科）                                         │   │
│   │  - 添加就诊建议（建议尽快就医，完善检查）                           │   │
│   └─────────────────────────────────────────────────────────────────────┘   │
│                                                                              │
│   Step 4: 返回结果                                                           │
│   ┌─────────────────────────────────────────────────────────────────────┐   │
│   │ 推荐结果：                                                          │   │
│   │ - 可能疾病：流感（85%）、病毒性脑膜炎（78%）                         │   │
│   │ - 推荐科室：神经内科、发热门诊                                       │   │
│   │ - 推荐医生：张医生（主任医师，神经内科）                             │   │
│   │ - 建议：症状持续3天且有发热，建议尽快就医完善检查                     │   │
│   └─────────────────────────────────────────────────────────────────────┘   │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 六、性能优化

### 6.1 缓存策略

```go
// 多级缓存

type CacheManager struct {
    localCache  *ristretto.Cache  // 本地缓存
    redisCache  *redis.Client     // Redis缓存
    vectorCache *VectorCache      // 向量结果缓存
}

// 缓存策略
const (
    CacheStrategyEmbedding = "embedding"  // 向量化结果缓存
    CacheStrategyVector    = "vector"     // 向量检索结果缓存
    CacheStrategyLLM       = "llm"        // LLM响应缓存
)

// GetEmbedding 获取向量（带缓存）
func (c *CacheManager) GetEmbedding(ctx context.Context, text string) ([]float32, error) {
    // 1. 计算缓存key
    key := "emb:" + hashText(text)
    
    // 2. 本地缓存
    if cached, found := c.localCache.Get(key); found {
        return cached.([]float32), nil
    }
    
    // 3. Redis缓存
    if data, err := c.redisCache.Get(ctx, key).Bytes(); err == nil {
        vector := decodeVector(data)
        c.localCache.Set(key, vector, 0)
        return vector, nil
    }
    
    // 4. 调用Embedding API
    vector, err := c.embedder.Embed(ctx, text)
    if err != nil {
        return nil, err
    }
    
    // 5. 写入缓存
    c.localCache.Set(key, vector, 0)
    c.redisCache.Set(ctx, key, encodeVector(vector), 24*time.Hour)
    
    return vector, nil
}

// GetVectorSearch 向量检索（带缓存）
func (c *CacheManager) GetVectorSearch(ctx context.Context, collection string, query string, limit int) ([]VectorResult, error) {
    key := fmt.Sprintf("vec:%s:%s:%d", collection, hashText(query), limit)
    
    // 尝试从缓存获取
    if data, err := c.redisCache.Get(ctx, key).Bytes(); err == nil {
        var results []VectorResult
        if err := json.Unmarshal(data, &results); err == nil {
            return results, nil
        }
    }
    
    // 执行检索
    vector, err := c.GetEmbedding(ctx, query)
    if err != nil {
        return nil, err
    }
    
    results, err := c.vectorRepo.Search(ctx, collection, vector, limit)
    if err != nil {
        return nil, err
    }
    
    // 缓存结果（疾病和医生数据变化较慢，可缓存较长时间）
    data, _ := json.Marshal(results)
    c.redisCache.Set(ctx, key, data, 1*time.Hour)
    
    return results, nil
}
```

### 6.2 并发优化

```go
// 并行检索
func (s *AIService) parallelSearch(ctx context.Context, queryVector []float32) (*SearchResults, error) {
    var results SearchResults
    
    // 使用errgroup并行检索
    var g errgroup.Group
    
    // 并行搜索疾病
    g.Go(func() error {
        diseases, err := s.vectorRepo.SearchDiseases(ctx, queryVector, 5)
        if err != nil {
            return err
        }
        results.Diseases = diseases
        return nil
    })
    
    // 并行搜索医生
    g.Go(func() error {
        doctors, err := s.vectorRepo.SearchDoctors(ctx, queryVector, 4)
        if err != nil {
            return err
        }
        results.Doctors = doctors
        return nil
    })
    
    // 并行搜索药品
    g.Go(func() error {
        drugs, err := s.vectorRepo.SearchDrugs(ctx, queryVector, 4)
        if err != nil {
            return err
        }
        results.Drugs = drugs
        return nil
    })
    
    // 等待所有检索完成
    if err := g.Wait(); err != nil {
        return nil, err
    }
    
    return &results, nil
}
```

---

## 七、安全保障

### 7.1 数据安全

```go
// 数据脱敏
func (s *AIService) maskSensitiveData(rec *Recommendation) {
    // 医生信息脱敏
    for i := range rec.RecommendedDoctors {
        // 保留医生姓名，但脱敏其他信息
        rec.RecommendedDoctors[i].DoctorName = maskName(rec.RecommendedDoctors[i].DoctorName)
    }
    
    // 患者历史病历脱敏
    for i := range rec.SimilarRecords {
        rec.SimilarRecords[i].MedicalNumber = maskMedicalNumber(rec.SimilarRecords[i].MedicalNumber)
    }
}

// HIPAA合规
func (s *AIService) ensureHIPAACompliance(ctx context.Context, patientID int64) error {
    // 1. 验证用户权限
    if !s.hasPermission(ctx, patientID) {
        return ErrUnauthorizedAccess
    }
    
    // 2. 审计日志
    s.auditLog.Log(ctx, AuditEvent{
        Action:    "diagnosis_query",
        PatientID: patientID,
        UserID:    getCurrentUserID(ctx),
        Timestamp: time.Now(),
    })
    
    // 3. 数据加密
    // Qdrant中的向量数据是加密的，只有应用层可以解密
    
    return nil
}
```

### 7.2 医疗安全

```go
// 安全过滤器
type SafetyFilter struct {
    dangerousSymptoms []string
    emergencyKeywords []string
}

// Filter 过滤危险情况
func (f *SafetyFilter) Filter(input UserInput) (*SafetyResult, error) {
    // 1. 检测紧急情况
    for _, keyword := range f.emergencyKeywords {
        if strings.Contains(input.Text, keyword) {
            return &SafetyResult{
                IsSafe:      false,
                Level:       LevelEmergency,
                Message:     "检测到紧急症状，请立即拨打120或前往急诊科！",
                Action:      ActionRedirectEmergency,
            }, nil
        }
    }
    
    // 2. 检测危险症状
    dangerScore := 0
    for _, symptom := range f.dangerousSymptoms {
        if strings.Contains(input.Text, symptom) {
            dangerScore++
        }
    }
    
    if dangerScore >= 3 {
        return &SafetyResult{
            IsSafe:      false,
            Level:       LevelHighRisk,
            Message:     "您的症状可能较为严重，建议尽快就医！",
            Action:      ActionRecommendUrgent,
        }, nil
    }
    
    return &SafetyResult{IsSafe: true}, nil
}
```

---

## 八、总结

### 8.1 架构亮点

1. **RAG增强**：通过向量检索提供事实性支撑，减少LLM幻觉
2. **混合Agent**：平衡响应速度和智能程度，适应不同场景
3. **模块化设计**：Skills工具化，易于扩展和维护
4. **安全保障**：多层安全过滤，符合医疗数据合规要求

### 8.2 性能指标

| 指标 | 目标值 | 优化手段 |
|-----|--------|---------|
| 向量化延迟 | <500ms | 缓存、并发 |
| 向量检索延迟 | <100ms | HNSW索引、Redis缓存 |
| Agent推理延迟 | 1-3s | 并行工具调用 |
| 整体响应时间 | <3s | 异步处理、流式输出 |

### 8.3 演进路线

```
Phase 1 (当前): 基础RAG
├── Qdrant向量库
├── 基础疾病检索
└── 简单Agent

Phase 2: 增强RAG
├── 多路召回（向量+关键词+规则）
├── Rerank重排序
├── 多模态（图文）
└── 记忆机制

Phase 3: 高级Agent
├── 多Agent协作
├── 知识图谱集成
├── 个性化推荐
└── 持续学习
```

---

**文档版本**: V1.0  
**更新日期**: 2026-03-27  
**作者**: 架构组
