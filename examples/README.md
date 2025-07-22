# 📊 Dados Brutos - GitHub Analytics

## 📁 Arquivos de Exemplo

Este diretório contém exemplos reais de respostas da API em formato JSON bruto.

### 🏢 Repositórios
- **Arquivo:** `raw_repositories.json`
- **Endpoint:** `GET /api/repositories/{owner}?limit=5`
- **Descrição:** Lista repositórios de um owner com metadados

### 💻 Commits  
- **Arquivo:** `raw_commit_metrics.json`
- **Endpoint:** `GET /api/metrics/commits/{owner}/{repo}`
- **Descrição:** Métricas detalhadas de commits, timeline e leaderboard

### 🔀 Pull Requests
- **Arquivo:** `raw_pr_metrics.json` 
- **Endpoint:** `GET /api/metrics/prs/{owner}/{repo}`
- **Descrição:** Estatísticas de PRs, autores e timeline

### 🏆 Contribuições
- **Arquivo:** `raw_contribution_metrics.json`
- **Endpoint:** `GET /api/metrics/contributions/{owner}/{repo}`
- **Descrição:** Score ponderado de contribuições por desenvolvedor

## 🔧 Como Usar no Frontend

### JavaScript Vanilla
```javascript
async function fetchRawData(endpoint) {
    const response = await fetch(`/api/${endpoint}`);
    const result = await response.json();
    
    if (result.success) {
        // Processar dados brutos
        console.log(result.data);
        return result.data;
    }
}

// Exemplo de uso
const commits = await fetchRawData('metrics/commits/microsoft/vscode?start_date=2025-07-01&end_date=2025-07-21');
```

### React
```jsx
const [data, setData] = useState(null);

useEffect(() => {
    fetch('/api/metrics/commits/microsoft/vscode?start_date=2025-07-01&end_date=2025-07-21')
        .then(r => r.json())
        .then(result => {
            if (result.success) {
                setData(result.data);
            }
        });
}, []);
```

### Chart.js Custom
```javascript
function createCustomChart(rawData) {
    const ctx = document.getElementById('myChart').getContext('2d');
    
    new Chart(ctx, {
        type: 'line',
        data: {
            labels: rawData.timeline.map(t => t.date),
            datasets: [{
                label: 'Commits',
                data: rawData.timeline.map(t => t.commits),
                borderColor: 'rgb(75, 192, 192)',
                tension: 0.1
            }]
        },
        options: {
            responsive: true,
            plugins: {
                title: {
                    display: true,
                    text: `Total: ${rawData.total_commits} commits`
                }
            }
        }
    });
}
```

## 🎯 Vantagens dos Dados Brutos

### ✅ Flexibilidade Total
- Criar qualquer tipo de visualização
- Combinar múltiplas métricas
- Filtrar dados customizados
- Análises estatísticas

### ✅ Reutilização
- Um endpoint, múltiplos usos
- Independente de biblioteca UI
- Cache mais eficiente
- Processamento no cliente

### ✅ Debugging
- Ver dados originais
- Validar cálculos
- Investigar problemas
- Auditoria de dados

## 📈 Comparação com Chart.js Format

| Aspecto | Dados Brutos | Chart.js Format |
|---------|--------------|-----------------|
| Flexibilidade | ⭐⭐⭐⭐⭐ | ⭐⭐ |
| Performance | ⭐⭐⭐ | ⭐⭐⭐⭐⭐ |
| Reutilização | ⭐⭐⭐⭐⭐ | ⭐⭐ |
| Facilidade | ⭐⭐ | ⭐⭐⭐⭐⭐ |

## 🚀 Script de Teste

Execute `../test_raw_endpoints.sh` para testar todos os endpoints automaticamente:

```bash
./test_raw_endpoints.sh
```

Este script irá:
1. Fazer requisições para todos os endpoints
2. Salvar as respostas em arquivos de teste
3. Mostrar o status de cada requisição
4. Permitir validar os dados rapidamente
