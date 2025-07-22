# 📊 GitHub Analytics - Exemplos de Dados Brutos

Este arquivo contém exemplos reais dos dados brutos retornados pelos endpoints da API.

## 🚀 Como usar

Os dados brutos permitem máxima flexibilidade no frontend para:
- Criar gráficos customizados
- Mostrar estatísticas detalhadas  
- Integrar com diferentes bibliotecas de visualização
- Processar dados para análises específicas

## 📂 Endpoints Disponíveis

### 1. 🏢 Repositórios - `/api/repositories/{owner}`
**URL:** `GET /api/repositories/microsoft?limit=5`
**Exemplo:** [raw_repositories.json](examples/raw_repositories.json)

```json
{
  "success": true,
  "data": [
    {
      "owner": "microsoft",
      "name": "winget-pkgs",
      "full_name": "microsoft/winget-pkgs",
      "description": "The Microsoft community Windows Package Manager manifest repository",
      "private": false,
      "fork": false,
      "created_at": "2019-07-16T22:21:48Z",
      "updated_at": "2025-07-21T18:52:24Z",
      "language": "PowerShell",
      "stars": 9547,
      "forks": 5242
    }
  ]
}
```

### 2. 💻 Métricas de Commits - `/api/metrics/commits/{owner}/{repo}`
**URL:** `GET /api/metrics/commits/microsoft/vscode?start_date=2025-07-01&end_date=2025-07-21`
**Exemplo:** [raw_commit_metrics.json](examples/raw_commit_metrics.json)

```json
{
  "success": true,
  "data": {
    "total_commits": 814,
    "unique_authors": 86,
    "date_range": {
      "start": "2025-07-01",
      "end": "2025-07-22"
    },
    "leaderboard": [
      {
        "author": "Daniel Imms",
        "email": "2193314+Tyriar@users.noreply.github.com", 
        "commits": 141,
        "additions": 0,
        "deletions": 0
      }
    ],
    "timeline": [
      {
        "date": "2025-07-01",
        "commits": 4
      },
      {
        "date": "2025-07-02", 
        "commits": 3
      }
    ]
  }
}
```

### 3. 🔀 Métricas de Pull Requests - `/api/metrics/prs/{owner}/{repo}`
**URL:** `GET /api/metrics/prs/microsoft/vscode?start_date=2025-07-01&end_date=2025-07-21`
**Exemplo:** [raw_pr_metrics.json](examples/raw_pr_metrics.json)

```json
{
  "success": true,
  "data": {
    "total_prs": 763,
    "open_prs": 81,
    "closed_prs": 103,
    "merged_prs": 579,
    "draft_prs": 78,
    "date_range": {
      "start": "2025-07-01",
      "end": "2025-07-22"
    },
    "leaderboard": [
      {
        "author": "bpasero",
        "total_prs": 75,
        "open_prs": 0,
        "merged_prs": 72,
        "closed_prs": 3,
        "draft_prs": 2,
        "avg_merge_time_hours": 12.83
      }
    ],
    "timeline": [
      {
        "date": "2025-07-01",
        "created": 2,
        "merged": 1,
        "closed": 0
      }
    ]
  }
}
```

### 4. 🏆 Métricas de Contribuições - `/api/metrics/contributions/{owner}/{repo}`
**URL:** `GET /api/metrics/contributions/microsoft/vscode?start_date=2025-07-01&end_date=2025-07-21`

```json
{
  "success": true,
  "data": {
    "total_score": 1250.5,
    "unique_contributors": 86,
    "date_range": {
      "start": "2025-07-01",
      "end": "2025-07-22"
    },
    "leaderboard": [
      {
        "author": "Daniel Imms",
        "email": "2193314+Tyriar@users.noreply.github.com",
        "commits": 141,
        "prs": 15,
        "issues": 5,
        "reviews": 25,
        "score": 245.75,
        "weighted_score": {
          "commits": 141.0,
          "prs": 75.0,
          "issues": 15.0,
          "reviews": 12.5
        }
      }
    ],
    "timeline": [
      {
        "date": "2025-07-01",
        "total_score": 45.5,
        "contributors": 12
      }
    ],
    "weights": {
      "commit": 1.0,
      "pr": 5.0,
      "issue": 3.0,
      "review": 0.5
    }
  }
}
```

---

## 🔧 Vantagens dos Dados Brutos

### ✅ **vs. Dados formatados para Chart.js**

| Aspecto | Dados Brutos | Chart.js Format |
|---------|--------------|-----------------|
| **Flexibilidade** | ⭐⭐⭐⭐⭐ | ⭐⭐ |
| **Reutilização** | ⭐⭐⭐⭐⭐ | ⭐⭐ |
| **Performance** | ⭐⭐⭐ | ⭐⭐⭐⭐⭐ |
| **Facilidade** | ⭐⭐ | ⭐⭐⭐⭐⭐ |

### 🎯 **Casos de Uso Ideais**
- **Dashboard personalizado** - Dados brutos
- **Prototipagem rápida** - Chart.js format  
- **Análise de dados** - Dados brutos
- **Integração múltipla** - Dados brutos
- **Debug/desenvolvimento** - Chart.js format

---
