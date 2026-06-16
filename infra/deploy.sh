#!/bin/bash
# Deploy helper script for SCD infrastructure

set -e

COMPONENT=$1
ACTION=${2:-up}

usage() {
    echo "Usage: ./deploy.sh <component> [action]"
    echo ""
    echo "Components:"
    echo "  kafka       - Kafka broker + Kafka UI"
    echo "  redis       - Redis (local dev only)"
    echo "  ai-server   - AI language detection service"
    echo "  all         - All components"
    echo ""
    echo "Actions:"
    echo "  up          - Start services (default)"
    echo "  down        - Stop services"
    echo "  restart     - Restart services"
    echo "  logs        - Show logs"
    echo "  ps          - Show running containers"
    echo ""
    echo "Examples:"
    echo "  ./deploy.sh kafka up"
    echo "  ./deploy.sh ai-server logs"
    echo "  ./deploy.sh all down"
    exit 1
}

check_env() {
    local dir=$1
    if [ ! -f "$dir/.env" ]; then
        echo "⚠️  Warning: $dir/.env not found"
        echo "   Copy from .env.example: cp $dir/.env.example $dir/.env"
        read -p "   Continue anyway? (y/N) " -n 1 -r
        echo
        if [[ ! $REPLY =~ ^[Yy]$ ]]; then
            exit 1
        fi
    fi
}

deploy_component() {
    local component=$1
    local action=$2
    local dir="infra/$component"

    if [ ! -d "$dir" ]; then
        echo "❌ Component '$component' not found in infra/"
        exit 1
    fi

    echo "🚀 $action $component..."
    check_env "$dir"

    cd "$dir"

    case $action in
        up)
            docker-compose up -d
            echo "✅ $component started"
            ;;
        down)
            docker-compose down
            echo "✅ $component stopped"
            ;;
        restart)
            docker-compose restart
            echo "✅ $component restarted"
            ;;
        logs)
            docker-compose logs -f
            ;;
        ps)
            docker-compose ps
            ;;
        *)
            echo "❌ Unknown action: $action"
            usage
            ;;
    esac

    cd - > /dev/null
}

# Main
if [ -z "$COMPONENT" ]; then
    usage
fi

case $COMPONENT in
    kafka|redis|ai-server)
        deploy_component "$COMPONENT" "$ACTION"
        ;;
    all)
        for comp in kafka redis ai-server; do
            deploy_component "$comp" "$ACTION"
            echo ""
        done
        ;;
    *)
        echo "❌ Unknown component: $COMPONENT"
        usage
        ;;
esac

echo "✅ Done!"
