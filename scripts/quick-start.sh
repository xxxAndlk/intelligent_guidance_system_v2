#!/bin/bash
# 智慧医疗导诊系统 v2 - 快速启动脚本

set -e

echo "=========================================="
echo "  智慧医疗导诊系统 v2 - 快速启动脚本"
echo "=========================================="
echo ""

# 检查目录
PROJECT_ROOT="$(cd "$(dirname "$0")/.." && pwd)"
cd "$PROJECT_ROOT"

echo "📁 项目目录: $PROJECT_ROOT"
echo ""

# 步骤1: 启动基础设施
echo "🚀 步骤1: 启动基础设施..."
cd deploy
docker-compose up -d mysql redis qdrant rabbitmq consul

# 等待MySQL启动
echo "⏳ 等待MySQL启动..."
sleep 30

# 检查服务状态
echo "📊 检查服务状态..."
docker-compose ps

# 步骤2: 初始化数据库
echo ""
echo "🗄️ 步骤2: 初始化数据库..."
docker-compose exec -T mysql mysql -uroot -proot123 medical_system < init/init.sql || echo "数据库可能已初始化"

echo ""
echo "✅ 基础设施启动完成！"
echo ""
echo "访问地址:"
echo "  MySQL:     localhost:3306"
echo "  Redis:     localhost:6379"
echo "  RabbitMQ:  localhost:15672 (admin/admin123)"
echo "  Qdrant:    localhost:6333"
echo "  Consul:    localhost:8500"
echo ""
echo "下一步:"
echo "  1. 生成protobuf代码: make generate"
echo "  2. 运行AI服务: make run-ai"
echo "  3. 运行Patient服务: make run-patient"
echo ""
echo "=========================================="
