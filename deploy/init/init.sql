-- MySQL 8.0.X 初始化脚本
SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ============================================
-- 创建数据库
-- ============================================
CREATE DATABASE IF NOT EXISTS medical_system 
DEFAULT CHARACTER SET utf8mb4 
DEFAULT COLLATE utf8mb4_unicode_ci;

USE medical_system;

-- ============================================
-- 用户表
-- ============================================
CREATE TABLE IF NOT EXISTS `user` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '用户ID',
    `openid` VARCHAR(255) DEFAULT NULL COMMENT '微信OpenID',
    `username` VARCHAR(50) NOT NULL COMMENT '用户名',
    `password` VARCHAR(255) NOT NULL COMMENT '加密密码',
    `name` VARCHAR(50) DEFAULT NULL COMMENT '姓名',
    `age` INT DEFAULT NULL COMMENT '年龄',
    `phone` VARCHAR(20) DEFAULT NULL COMMENT '手机号',
    `sex` TINYINT DEFAULT NULL COMMENT '性别：1男 2女',
    `email` VARCHAR(100) DEFAULT NULL COMMENT '邮箱',
    `id_card` VARCHAR(18) DEFAULT NULL COMMENT '身份证号',
    `avatar` VARCHAR(500) DEFAULT NULL COMMENT '头像URL',
    `home_address` VARCHAR(255) DEFAULT NULL COMMENT '家庭地址',
    `status` TINYINT DEFAULT 1 COMMENT '状态：0禁用 1启用',
    `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    `deleted_at` DATETIME DEFAULT NULL COMMENT '删除时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_username` (`username`),
    UNIQUE KEY `uk_phone` (`phone`),
    KEY `idx_status` (`status`),
    KEY `idx_deleted` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户表';

-- ============================================
-- 医生表
-- ============================================
CREATE TABLE IF NOT EXISTS `doctor` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '医生ID',
    `username` VARCHAR(50) NOT NULL COMMENT '用户名',
    `password` VARCHAR(255) NOT NULL COMMENT '密码',
    `openid` VARCHAR(255) DEFAULT NULL COMMENT '微信OpenID',
    `name` VARCHAR(50) NOT NULL COMMENT '姓名',
    `age` INT DEFAULT NULL COMMENT '年龄',
    `phone` VARCHAR(20) DEFAULT NULL COMMENT '手机号',
    `sex` TINYINT DEFAULT NULL COMMENT '性别',
    `id_card` VARCHAR(18) DEFAULT NULL COMMENT '身份证号',
    `avatar` VARCHAR(500) DEFAULT NULL COMMENT '头像',
    `employee_id` VARCHAR(20) DEFAULT NULL COMMENT '员工工号',
    `position` TINYINT DEFAULT NULL COMMENT '职务：1主任医师 2副主任医师 3主治医师 4住院医师 5助理医师',
    `designation` TINYINT DEFAULT NULL COMMENT '职称：1院长 2副院长 3科室主任 4医疗组长 5医疗组员',
    `department_id` BIGINT UNSIGNED DEFAULT NULL COMMENT '科室ID',
    `is_expert` TINYINT DEFAULT 0 COMMENT '是否专家：0否 1是',
    `status` TINYINT DEFAULT 2 COMMENT '状态：1工作中 2休息中 3请假中',
    `intro` TEXT COMMENT '个人介绍',
    `work_experience` TEXT COMMENT '工作经历',
    `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP,
    `updated_at` DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `deleted_at` DATETIME DEFAULT NULL,
    PRIMARY KEY (`id`),
    KEY `idx_department` (`department_id`),
    KEY `idx_status` (`status`),
    KEY `idx_is_expert` (`is_expert`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='医生表';

-- ============================================
-- 科室表
-- ============================================
CREATE TABLE IF NOT EXISTS `department` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '科室ID',
    `name` VARCHAR(100) NOT NULL COMMENT '科室名称',
    `introduction` TEXT COMMENT '科室简介',
    `person_in_charge_id` VARCHAR(20) DEFAULT NULL COMMENT '科室负责人工号',
    `phone` VARCHAR(20) DEFAULT NULL COMMENT '科室电话',
    `address` VARCHAR(255) DEFAULT NULL COMMENT '地址',
    `type` TINYINT DEFAULT 1 COMMENT '类型：1门诊科室 2住院科室',
    `staff_count` INT DEFAULT 0 COMMENT '科室人数',
    `status` TINYINT DEFAULT 2 COMMENT '状态：1营业中 2休息中 3已停用',
    `registration_fee` DECIMAL(10,2) DEFAULT 0.00 COMMENT '普通挂号费',
    `expert_registration_fee` DECIMAL(10,2) DEFAULT 0.00 COMMENT '专家挂号费',
    `duty_doctor_id` BIGINT UNSIGNED DEFAULT NULL COMMENT '值班医生ID',
    `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP,
    `updated_at` DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `deleted_at` DATETIME DEFAULT NULL,
    PRIMARY KEY (`id`),
    KEY `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='科室表';

-- ============================================
-- 疾病知识表
-- ============================================
CREATE TABLE IF NOT EXISTS `disease_knowledge` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '疾病ID',
    `name` VARCHAR(100) NOT NULL COMMENT '疾病名称',
    `symptoms` TEXT COMMENT '症状描述',
    `diagnosis` TEXT COMMENT '诊断信息',
    `treatment_plan` TEXT COMMENT '治疗方案',
    `prevention_measures` TEXT COMMENT '预防措施',
    `drug_category` VARCHAR(255) DEFAULT NULL COMMENT '适用药品类型',
    `department` VARCHAR(100) DEFAULT NULL COMMENT '所属科室',
    `severity` TINYINT DEFAULT 1 COMMENT '严重程度：1轻微 2中等 3严重 4危急',
    `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP,
    `updated_at` DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `deleted_at` DATETIME DEFAULT NULL,
    PRIMARY KEY (`id`),
    KEY `idx_name` (`name`),
    KEY `idx_department` (`department`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='疾病知识表';

-- ============================================
-- 药品表
-- ============================================
CREATE TABLE IF NOT EXISTS `drug` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '药品ID',
    `name` VARCHAR(100) NOT NULL COMMENT '药品名称',
    `description` TEXT COMMENT '药品描述',
    `price` DECIMAL(10,2) DEFAULT 0.00 COMMENT '药品价格',
    `quantity_in_stock` INT DEFAULT 0 COMMENT '库存数量',
    `supplier` VARCHAR(100) DEFAULT NULL COMMENT '供应商',
    `category` TINYINT DEFAULT NULL COMMENT '药品分类',
    `use_method` TINYINT DEFAULT NULL COMMENT '使用方式',
    `expiration_date` VARCHAR(50) DEFAULT NULL COMMENT '保质期',
    `contraindication` TEXT COMMENT '禁忌症',
    `side_effects` TEXT COMMENT '副作用',
    `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP,
    `updated_at` DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `deleted_at` DATETIME DEFAULT NULL,
    PRIMARY KEY (`id`),
    KEY `idx_category` (`category`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='药品表';

-- ============================================
-- 病历表
-- ============================================
CREATE TABLE IF NOT EXISTS `medical_record` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '病历ID',
    `medical_number` VARCHAR(50) NOT NULL COMMENT '病历单号',
    `patient_id` BIGINT UNSIGNED NOT NULL COMMENT '患者ID',
    `doctor_id` BIGINT UNSIGNED DEFAULT NULL COMMENT '医生ID',
    `department_id` BIGINT UNSIGNED DEFAULT NULL COMMENT '科室ID',
    `status` TINYINT DEFAULT 1 COMMENT '状态：1挂号中 2排队中 3诊断中 4手术中 5诊断结束 6已取消',
    `symptom` TEXT COMMENT '症状描述',
    `diagnosis` TEXT COMMENT '诊断',
    `treatment_plan` TEXT COMMENT '治疗方案',
    `billing_time` DATETIME DEFAULT NULL COMMENT '开单时间',
    `pay_method` TINYINT DEFAULT NULL COMMENT '支付方式：1微信 2支付宝',
    `pay_status` TINYINT DEFAULT 1 COMMENT '支付状态：1未支付 2已支付 3已退款',
    `amount` DECIMAL(10,2) DEFAULT 0.00 COMMENT '金额',
    `settlement_time` DATETIME DEFAULT NULL COMMENT '结算时间',
    `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP,
    `updated_at` DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `deleted_at` DATETIME DEFAULT NULL,
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_medical_number` (`medical_number`),
    KEY `idx_patient` (`patient_id`),
    KEY `idx_doctor` (`doctor_id`),
    KEY `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='病历表';

-- ============================================
-- AI导诊会话表
-- ============================================
CREATE TABLE IF NOT EXISTS `diagnosis_session` (
    `id` VARCHAR(50) NOT NULL COMMENT '会话ID',
    `patient_id` BIGINT UNSIGNED NOT NULL COMMENT '患者ID',
    `status` TINYINT DEFAULT 1 COMMENT '状态：1收集中 2分析中 3已推荐 4已完成',
    `symptoms` JSON COMMENT '症状列表',
    `conversations` JSON COMMENT '对话历史',
    `recommendation` JSON COMMENT '推荐结果',
    `vector_queries` JSON COMMENT '向量查询记录',
    `confidence` DECIMAL(3,2) DEFAULT 0.00 COMMENT '置信度',
    `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP,
    `updated_at` DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    KEY `idx_patient` (`patient_id`),
    KEY `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='AI导诊会话表';

-- ============================================
-- 手术表
-- ============================================
CREATE TABLE IF NOT EXISTS `surgical` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '手术ID',
    `medical_id` BIGINT UNSIGNED NOT NULL COMMENT '病历ID',
    `patient_id` BIGINT UNSIGNED NOT NULL COMMENT '患者ID',
    `doctor_id` BIGINT UNSIGNED NOT NULL COMMENT '主刀医生ID',
    `department_id` BIGINT UNSIGNED NOT NULL COMMENT '科室ID',
    `type` VARCHAR(50) COMMENT '手术类型',
    `status` TINYINT DEFAULT 1 COMMENT '状态：1待手术 2手术中 3已完成 4已取消',
    `scheduled_time` DATETIME COMMENT '预定时间',
    `duration` INT COMMENT '预计时长(分钟)',
    `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP,
    `updated_at` DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    KEY `idx_medical` (`medical_id`),
    KEY `idx_patient` (`patient_id`),
    KEY `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='手术表';

-- ============================================
-- 挂号表
-- ============================================
CREATE TABLE IF NOT EXISTS `registration` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '挂号ID',
    `patient_id` BIGINT UNSIGNED NOT NULL COMMENT '患者ID',
    `doctor_id` BIGINT UNSIGNED NOT NULL COMMENT '医生ID',
    `department_id` BIGINT UNSIGNED NOT NULL COMMENT '科室ID',
    `type` TINYINT DEFAULT 1 COMMENT '类型：1普通挂号 2专家挂号',
    `status` TINYINT DEFAULT 1 COMMENT '状态：1待确认 2已确认 3已取消 4已完成',
    `appointment_time` DATETIME COMMENT '预约时间',
    `fee` DECIMAL(10,2) DEFAULT 0.00 COMMENT '挂号费',
    `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP,
    `updated_at` DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    KEY `idx_patient` (`patient_id`),
    KEY `idx_doctor` (`doctor_id`),
    KEY `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='挂号表';

-- ============================================
-- 插入初始数据
-- ============================================
INSERT INTO `department` (`name`, `introduction`, `type`, `status`, `registration_fee`, `expert_registration_fee`) VALUES
('内科', '综合内科，诊治各类内科疾病', 1, 1, 20.00, 50.00),
('外科', '综合外科，提供手术治疗', 1, 1, 20.00, 50.00),
('儿科', '儿童疾病诊治', 1, 1, 15.00, 40.00),
('妇产科', '妇科产科疾病诊治', 1, 1, 25.00, 60.00),
('急诊科', '24小时急诊服务', 1, 1, 30.00, 80.00),
('眼科', '眼科疾病诊治', 1, 1, 20.00, 50.00),
('耳鼻喉科', '耳鼻喉疾病诊治', 1, 1, 20.00, 50.00),
('口腔科', '口腔疾病诊治', 1, 1, 25.00, 60.00),
('皮肤科', '皮肤疾病诊治', 1, 1, 20.00, 50.00),
('骨科', '骨科疾病诊治', 1, 1, 25.00, 60.00);

INSERT INTO `disease_knowledge` (`name`, `symptoms`, `diagnosis`, `treatment_plan`, `department`, `severity`) VALUES
('流感', '高热、头痛、全身酸痛、咳嗽、流涕、咽痛', '流感病毒核酸检测阳性', '抗病毒治疗、对症支持治疗、多休息、多喝水', '内科', 2),
('上呼吸道感染', '鼻塞、流涕、咽痛、咳嗽、低热、头痛', '症状诊断，血常规检查', '对症治疗、休息、补充水分、必要时抗生素', '内科', 1),
('急性胃肠炎', '腹痛、腹泻、恶心、呕吐、发热', '症状+粪便检查', '补液、止泻、抗感染治疗、饮食调理', '内科', 2),
('高血压', '头晕、头痛、心悸、视物模糊', '血压测量≥140/90mmHg', '降压药物治疗、生活方式干预、定期监测', '内科', 2),
('糖尿病', '多饮、多尿、多食、体重下降、疲乏', '空腹血糖≥7.0mmol/L', '降糖药物治疗、饮食控制、运动治疗、血糖监测', '内科', 3);

SET FOREIGN_KEY_CHECKS = 1;
