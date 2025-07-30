# 🚀 GitHub Analytics - Scripts de Execução

Este arquivo documenta todos os scripts e comandos disponíveis para executar o projeto.

## 📋 Resumo dos Scripts

| Script | Descrição | Uso |
|--------|-----------|-----|
| `start_project.sh` | 🎯 **Script principal** - Setup completo e execução | `./start_project.sh` |
| `dev_server.sh` | 🔧 Servidor de desenvolvimento com restart manual | `./dev_server.sh` |
| `test_runner.sh` | 🧪 Suite de testes da API | `./test_runner.sh` |
| `test_raw_endpoints.sh` | 📊 Testa endpoints e salva exemplos | `./test_raw_endpoints.sh` |
| `setup.sh` | ⚙️ Apenas configuração inicial | `./setup.sh` |

## 🎯 Cenários de Uso

### 1️⃣ **Primeira Execução (Recomendado)**
```bash
./start_project.sh
```
- ✅ Verifica dependências
- ✅ Configura ambiente (.env)
- ✅ Escolhe modo de execução
- ✅ Interface amigável

### 2️⃣ **Desenvolvimento com 2 Terminais**

**Terminal 1 - Servidor:**
```bash
./dev_server.sh
```

**Terminal 2 - Testes:**
```bash
./test_runner.sh
```

### 3️⃣ **Execução Rápida (Make)**
```bash
# Setup + execução
make start          # Docker Compose
make full-start     # Local
make dev            # Desenvolvimento

# Gerenciamento
make stop           # Parar
make logs           # Ver logs
make help           # Ajuda
```

### 4️⃣ **Docker Compose**
```bash
# Usando make
make start
make logs
make stop

# Comandos diretos
docker-compose up -d
docker-compose logs -f
docker-compose down
```

## 🔧 Comandos Make Disponíveis

### 🚀 **Principais**
```bash
make start          # Docker Compose (recomendado)
make stop           # Parar Docker Compose  
make logs           # Ver logs
make full-start     # Execução local
```

### 🛠️ **Desenvolvimento**
```bash
make setup          # Configurar projeto
make build          # Compilar binário
make dev            # Modo desenvolvimento
make clean          # Limpar arquivos
```

### 🐳 **Docker**
```bash
make docker-build   # Build da imagem
make docker-run     # Executar container
make docker-stop    # Parar containers
```

### 🎬 **Demo & Testes**
```bash
make demo           # Demo com docker/compose
make test           # Executar testes
```

## 🎮 Fluxo de Trabalho Recomendado

### **Para Desenvolvimento:**
1. **Setup inicial:**
   ```bash
   ./start_project.sh
   # Escolher opção 4 (Setup only)
   ```

2. **Desenvolvimento ativo:**
   ```bash
   # Terminal 1
   ./dev_server.sh

   # Terminal 2  
   ./test_runner.sh
   ```

3. **Testes automáticos:**
   ```bash
   ./test_raw_endpoints.sh
   ```

### **Para Produção/Demo:**
```bash
make start    # Docker Compose
make logs     # Monitorar
make stop     # Quando terminar
```

## 📊 URLs Importantes

Após iniciar o servidor:

| Serviço | URL | Descrição |
|---------|-----|-----------|
| **Frontend** | http://localhost:8080 | Interface de debug |
| **API** | http://localhost:8080/api | Endpoints REST |
| **Health** | http://localhost:8080/health | Status do servidor |

### **Endpoints de Exemplo:**
```bash
# Repositórios
curl "http://localhost:8080/api/repositories/microsoft?limit=5"

# Métricas de commits
curl "http://localhost:8080/api/metrics/commits/microsoft/vscode?start_date=2025-07-01&end_date=2025-07-21"

# Gráficos Chart.js
curl "http://localhost:8080/api/charts/commits-timeline/microsoft/vscode?start_date=2025-07-01&end_date=2025-07-21"
```

## ⚙️ Configuração

### **Arquivo .env:**
```bash
GITHUB_TOKEN=your_github_token_here  # OBRIGATÓRIO
DATE_RANGE_START=2025-07-01         # Opcional
DATE_RANGE_END=2025-07-31           # Opcional  
PORT=8080                           # Opcional
```

### **Como obter GitHub Token:**
1. Acesse: https://github.com/settings/tokens
2. "Generate new token (classic)"
3. Selecione scopes: `repo`, `user`
4. Copie o token para o .env

## 🆘 Troubleshooting

### **"Port 8080 already in use"**
```bash
# Verificar o que está usando a porta
lsof -ti:8080

# Usar porta diferente
PORT=8081 ./dev_server.sh
# ou
PORT=8081 make dev
```

### **"GITHUB_TOKEN not set"**
```bash
# Editar .env
nano .env
# ou
vim .env

# Verificar configuração
cat .env
```

### **"Server not responding"**
```bash
# Verificar se está rodando
curl http://localhost:8080/health

# Ver logs (Docker)
make logs

# Verificar dependências
go mod tidy
```

### **"Permission denied"**
```bash
# Dar permissão aos scripts
chmod +x *.sh
```

## 🎯 Quick Start

**Para usuários novos:**
```bash
git clone <repo>
cd magic_ops
./start_project.sh
```

**Para desenvolvedores:**
```bash
# Terminal 1
./dev_server.sh

# Terminal 2
./test_runner.sh
```

**Para produção:**
```bash
make start
make logs
```
