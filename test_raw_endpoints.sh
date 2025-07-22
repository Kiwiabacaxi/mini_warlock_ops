#!/bin/bash

# 🚀 GitHub Analytics - Script de Testes dos Dados Brutos
# Este script faz requisições para todos os endpoints e salva os exemplos

API_BASE="http://localhost:8081/api"
OWNER="microsoft"
REPO="vscode"
START_DATE="2025-07-01"
END_DATE="2025-07-21"

echo "🔍 Testando todos os endpoints de dados brutos..."
echo "📅 Período: $START_DATE a $END_DATE"
echo "🏢 Owner: $OWNER"
echo "📦 Repo: $REPO"
echo ""

# Test 1: Repositories
echo "1️⃣ Testando Repositórios..."
curl -s "$API_BASE/repositories/$OWNER?limit=5" | jq . | tee test_repositories.json
echo ""

# Test 2: Commit Metrics  
echo "2️⃣ Testando Métricas de Commits..."
curl -s "$API_BASE/metrics/commits/$OWNER/$REPO?start_date=$START_DATE&end_date=$END_DATE" | jq . | head -50 | tee test_commit_metrics.json
echo ""

# Test 3: PR Metrics
echo "3️⃣ Testando Métricas de Pull Requests..."
curl -s "$API_BASE/metrics/prs/$OWNER/$REPO?start_date=$START_DATE&end_date=$END_DATE" | jq . | head -50 | tee test_pr_metrics.json  
echo ""

# Test 4: Contribution Metrics
echo "4️⃣ Testando Métricas de Contribuições..."
curl -s "$API_BASE/metrics/contributions/$OWNER/$REPO?start_date=$START_DATE&end_date=$END_DATE" | jq . | tee test_contribution_metrics.json
echo ""

echo "✅ Todos os testes concluídos!"
echo "📁 Arquivos gerados:"
echo "   - test_repositories.json"
echo "   - test_commit_metrics.json" 
echo "   - test_pr_metrics.json"
echo "   - test_contribution_metrics.json"
