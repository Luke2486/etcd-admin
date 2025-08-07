#!/bin/bash

# 数据库迁移脚本

set -e

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 项目根目录
PROJECT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
BACKEND_DIR="$PROJECT_ROOT/backend"

echo -e "${GREEN}=== etcd-admin Database Migration Tool ===${NC}"

# 检查参数
case "$1" in
    "up"|"migrate")
        echo -e "${YELLOW}Running database migration...${NC}"
        cd "$BACKEND_DIR"
        go run cmd/server/main.go -migrate
        echo -e "${GREEN}Migration completed successfully!${NC}"
        ;;
    "down"|"rollback")
        echo -e "${YELLOW}Rolling back database migration...${NC}"
        cd "$BACKEND_DIR"
        go run cmd/server/main.go -rollback
        echo -e "${GREEN}Migration rollback completed successfully!${NC}"
        ;;
    "reset")
        echo -e "${YELLOW}Resetting database (rollback + migrate)...${NC}"
        cd "$BACKEND_DIR"
        echo -e "${YELLOW}Step 1: Rolling back...${NC}"
        go run cmd/server/main.go -rollback
        echo -e "${YELLOW}Step 2: Migrating...${NC}"
        go run cmd/server/main.go -migrate
        echo -e "${GREEN}Database reset completed successfully!${NC}"
        ;;
    *)
        echo "Usage: $0 {up|migrate|down|rollback|reset}"
        echo ""
        echo "Commands:"
        echo "  up|migrate    - Run database migration"
        echo "  down|rollback - Rollback database migration"
        echo "  reset         - Rollback and then migrate (full reset)"
        echo ""
        echo "Examples:"
        echo "  $0 migrate    # Run migration"
        echo "  $0 rollback   # Rollback migration"
        echo "  $0 reset      # Reset database"
        exit 1
        ;;
esac
