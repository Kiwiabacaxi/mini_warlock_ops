# GitHub Analytics Backend

Um sistema completo de coleta e análise de métricas do GitHub usando Go, GraphQL e frontend de debug.

## 🚀 Funcionalidades

### Métricas de Commits
- **Leaderboard de commits**: Ranking de colaboradores por número de commits
- **Timeline de commits**: Evolução temporal dos commits
- Dados agregados: total de commits, autores únicos, adições/deleções

### Métricas de Pull Requests
- **Timeline de PRs**: Criação e merge de PRs ao longo do tempo
- **Leaderboard de PRs**: Ranking de colaboradores por PRs criados
- Filtros por status: Open, Closed, Merged, Draft
- Métricas: tempo médio para merge, taxa de aprovação

### Métricas de Contribuições (NOVO)
- **Leaderboard de contribuições**: Ranking ponderado de contribuidores
- **Sistema de pontuação**: Diferentes pesos para diferentes tipos de contribuição
  - Commits: 2.0 pontos
  - Reviews de PR: 2.0 pontos
  - Aprovações: 2.0 pontos
  - Issues: 1.0 ponto
  - Comentários: 0.5 pontos
  - Solicitações de mudança: 0.5 pontos

## 🛠️ Instalação

### Pré-requisitos
- Go 1.21+
- Git
- Token de acesso pessoal do GitHub

### Configuração

1. **Clone o repositório**:
```bash
git clone <repository-url>
cd github-analytics-backend
```

2. **Configure as variáveis de ambiente**:
```bash
cp .env.example .env
```

3. **Edite o arquivo `.env`**:
```env
GITHUB_TOKEN=your_personal_access_token_here
DATE_RANGE_START=2025-07-01
DATE_RANGE_END=2025-07-31
PORT=8080
```

4. **Instale as dependências**:
```bash
go mod tidy
```

5. **Execute o servidor**:
```bash
go run main.go
```

## 🔑 Token do GitHub

Para obter um token de acesso pessoal:

1. Acesse: https://github.com/settings/tokens
2. Clique em "Generate new token (classic)"
3. Selecione os escopos necessários:
   - `repo` - Acesso completo a repositórios
   - `user` - Acesso aos dados do usuário
4. Copie o token e adicione ao arquivo `.env`

## 📊 API Endpoints

### Repositórios
```
GET /api/repositories/{owner}?limit=50
```

### Métricas
```
GET /api/metrics/commits/{owner}/{repo}?start_date=2025-07-01&end_date=2025-07-31
GET /api/metrics/prs/{owner}/{repo}?start_date=2025-07-01&end_date=2025-07-31
GET /api/metrics/contributions/{owner}/{repo}?start_date=2025-07-01&end_date=2025-07-31
```

### Gráficos (Chart.js format)
```
GET /api/charts/commits-leaderboard/{owner}/{repo}
GET /api/charts/commits-timeline/{owner}/{repo}
GET /api/charts/prs-timeline/{owner}/{repo}
GET /api/charts/prs-leaderboard/{owner}/{repo}
GET /api/charts/contributions-leaderboard/{owner}/{repo}?repos=repo1,repo2,repo3
```

### Parâmetros de Query
- `start_date` e `end_date`: Formato YYYY-MM-DD
- `status`: Para PRs (open/closed/merged/draft)
- `repos`: Lista de repositórios separados por vírgula (para contribuições)

## 🖥️ Frontend de Debug

Acesse: http://localhost:8080

O frontend oferece:
- ✅ Interface simples para testar todos os endpoints
- ✅ Seleção de datas personalizáveis
- ✅ Visualização de JSON raw
- ✅ Gráficos interativos (Chart.js)
- ✅ Estatísticas resumidas
- ✅ Tratamento de erros

## 📈 Exemplos de Uso

### Analisar repositório Docker Compose
```bash
curl "http://localhost:8080/api/metrics/commits/docker/compose?start_date=2025-07-01&end_date=2025-07-31"
```

### Gerar gráfico de contribuições
```bash
curl "http://localhost:8080/api/charts/contributions-leaderboard/docker/compose"
```

### Múltiplos repositórios
```bash
curl "http://localhost:8080/api/charts/contributions-leaderboard/docker/compose?repos=compose,cli,buildx"
```

## 🏗️ Estrutura do Projeto

```
github-analytics-backend/
├── main.go                 # Servidor principal
├── models/
│   └── models.go          # Estruturas de dados
├── github/
│   └── client.go          # Cliente GitHub GraphQL/REST
├── services/
│   └── metrics.go         # Lógica de cálculo de métricas
├── handlers/
│   └── handlers.go        # Handlers HTTP
├── frontend/
│   └── index.html         # Interface de debug
├── .env                   # Configuração (não commitado)
├── .env.example          # Template de configuração
├── go.mod                # Dependências Go
└── README.md             # Este arquivo
```

## 🔧 Desenvolvimento

### Executar em modo de desenvolvimento
```bash
go run main.go
```

### Build para produção
```bash
go build -o github-analytics-backend main.go
./github-analytics-backend
```

### Testar endpoints específicos
```bash
# Health check
curl http://localhost:8080/health

# Listar repositórios
curl http://localhost:8080/api/repositories/docker

# Métricas de commits
curl "http://localhost:8080/api/metrics/commits/docker/compose?start_date=2025-07-01&end_date=2025-07-31"
```

## ⚠️ Limitações e Considerações

- **Rate Limiting**: A API do GitHub tem limites de taxa (5000 req/hora para token autenticado)
- **Repositórios Grandes**: Repositórios com muitos commits podem demorar para processar
- **Datas**: O sistema funciona melhor com períodos de até 3 meses
- **Tokens**: Certifique-se de que o token tem as permissões necessárias

## 🐛 Troubleshooting

### Token inválido
```
Error: failed to fetch repositories: 401 Unauthorized
```
**Solução**: Verifique se o token está correto e tem as permissões necessárias.

### Repositório não encontrado
```
Error: failed to fetch commits: 404 Not Found
```
**Solução**: Verifique se o repositório existe e é acessível com seu token.

### Rate limit excedido
```
Error: API rate limit exceeded
```
**Solução**: Aguarde ou use um token com limite maior.

## 📝 TODO

- [ ] Cache Redis para melhorar performance
- [ ] Suporte a webhooks para atualizações em tempo real
- [ ] Exportação de dados para CSV/Excel
- [ ] Dashboard mais avançado com filtros
- [ ] Suporte a organizações GitHub
- [ ] Métricas de issues e discussões

## 🤝 Contribuindo

1. Fork o projeto
2. Crie uma branch para sua feature (`git checkout -b feature/AmazingFeature`)
3. Commit suas mudanças (`git commit -m 'Add some AmazingFeature'`)
4. Push para a branch (`git push origin feature/AmazingFeature`)
5. Abra um Pull Request

## 📄 Licença

Este projeto está licenciado sob a MIT License - veja o arquivo LICENSE para detalhes.
