#!/bin/bash

# 快速构建并部署 panda-wiki-admin 前端到 Docker 容器
# 使用方法: ./light_build.sh [container_name]
# 默认容器名: panda-wiki-nginx

CONTAINER_NAME=${1:-panda-wiki-nginx}

# 设置版本号（应该与后端版本一致）
export VITE_APP_VERSION="3.42.0"

echo "开始构建前端代码..."
echo "版本号: ${VITE_APP_VERSION}"
cd "$(dirname "$0")"

# 检查 pnpm 是否可用，否则尝试使用 corepack 启用
if ! command -v pnpm >/dev/null 2>&1; then
    echo "未检测到 pnpm，尝试使用 corepack 启用 pnpm..."
    if command -v corepack >/dev/null 2>&1; then
        corepack enable && corepack prepare pnpm@10.12.1 --activate
    fi
fi
if ! command -v pnpm >/dev/null 2>&1; then
    echo "错误：未安装 pnpm，请先安装 Node.js 20+ 并启用 pnpm（corepack enable）"
    exit 1
fi

# 如果缺少依赖则先安装
if [ ! -d "node_modules" ]; then
    echo "检测到缺少 node_modules，开始安装依赖..."
    pnpm install || {
        echo "依赖安装失败，请检查网络或重试 pnpm install"
        exit 1
    }
fi

pnpm run build

if [ $? -ne 0 ]; then
    echo "构建失败，请检查错误信息"
    exit 1
fi

echo "构建完成，开始复制文件到容器..."

# 检查容器是否存在
if ! docker ps -a --format '{{.Names}}' | grep -q "^${CONTAINER_NAME}$"; then
    echo "错误: 容器 ${CONTAINER_NAME} 不存在"
    echo "可用容器列表:"
    docker ps -a --format '{{.Names}}'
    exit 1
fi

# 检查容器是否运行中
if ! docker ps --format '{{.Names}}' | grep -q "^${CONTAINER_NAME}$"; then
    echo "警告: 容器 ${CONTAINER_NAME} 未运行，正在启动..."
    docker start ${CONTAINER_NAME}
    sleep 2
fi

# 复制 dist 目录到容器
echo "复制 dist 目录到容器 ${CONTAINER_NAME}:/opt/frontend/dist/..."
if [ -d "dist" ]; then
    docker cp dist/. ${CONTAINER_NAME}:/opt/frontend/dist/
else
    echo "错误: dist 目录不存在"
    exit 1
fi

# 重载 Nginx 配置
echo "重载 Nginx 配置..."
docker exec ${CONTAINER_NAME} nginx -s reload

if [ $? -eq 0 ]; then
    echo "部署完成！"
    echo "Nginx 已重载，新代码已生效"
else
    echo "警告: Nginx 重载失败，可能需要重启容器"
    echo "可以尝试: docker restart ${CONTAINER_NAME}"
fi
