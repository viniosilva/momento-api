# Guia de Documentação - Pinnado

Este projeto utiliza um sistema de **personas/agentes** para diferentes contextos de desenvolvimento. Cada documento serve um propósito específico.

## 📚 Como Usar Esta Documentação

### 🎯 **Para Implementação de Código**

Use o **`.cursorrules`** (raiz do projeto):
- Padrões arquiteturais obrigatórios
- Estrutura de camadas (Domain, Application, Infrastructure, Presentation)
- Convenções de código e nomenclatura
- Estratégia de testes
- Sequência de implementação

**Quando usar**: Ao implementar features, seguir padrões, escrever testes.

---

### 📋 **Para Planejar Features (Plan Mode)**

Use **`docs/ARCHITECT_BACKEND_GO.md`**:
- Decompor User Stories em 7 fases
- Checklist de tarefas técnicas
- Definir contratos (structs, interfaces, funções)
- Output: Plano estruturado (não código)

**Quando usar**: Antes de implementar, para criar roteiro de execução.

---


### 💼 **Para Criar Histórias de Usuário**

Use **`docs/PRODUCT_MANAGER.md`**:
- Princípio INVEST
- Lógica de precedência de dados
- Formato de histórias

**Quando usar**: Ao iniciar uma nova feature ou épico.

---

### 📋 **Para Detalhar Requisitos**

Use **`docs/PRODUCT_OWNER.md`**:
- Especificar critérios de aceitação
- Detalhar regras de negócio
- Cenários e edge cases

**Quando usar**: Para expandir histórias em especificações técnicas.

---

### 🧪 **Para Testes e QA**

Use **`docs/QUALITY_ASSURANCE.md`**:
- Estratégias de teste
- Cenários de teste
- Cobertura e qualidade

**Quando usar**: Durante e após implementação para garantir qualidade.

---

### 🔧 **Para Débitos Técnicos**

Use **`docs/TECHNICAL_DEBT.md`**:
- Análise crítica da arquitetura
- Pontos de melhoria priorizados
- Soluções propostas com exemplos
- Roadmap de implementação

**Quando usar**: Para planejamento de refactoring, revisões técnicas, onboarding de novos membros.

---

## 🔄 **Fluxo de Trabalho Recomendado**

```
1. PRODUCT_MANAGER.md
   ↓ (Criar histórias de usuário)
   
2. PRODUCT_OWNER.md
   ↓ (Detalhar requisitos e critérios)
   
3. ARCHITECT_BACKEND_GO.md
   ↓ (Planejar arquitetura técnica)
   
4. .cursorrules
   ↓ (Implementar seguindo padrões)
   
5. QUALITY_ASSURANCE.md
   ↓ (Testar e validar)
   
✅ Feature completa e com qualidade
```

---

## 🎯 **Guia Rápido por Contexto**

| Contexto | Documento Principal | Documento Secundário |
|----------|---------------------|---------------------|
| **Implementando código** | `.cursorrules` | - |
| **Planejando arquitetura** | `ARCHITECT_BACKEND_GO.md` | `.cursorrules` |
| **Criando user stories** | `PRODUCT_MANAGER.md` | - |
| **Detalhando requisitos** | `PRODUCT_OWNER.md` | - |
| **Testando** | `QUALITY_ASSURANCE.md` | `.cursorrules` (seção 7) |
| **Refactoring/Débitos técnicos** | `TECHNICAL_DEBT.md` | `.cursorrules` |
| **Frontend** | `ARCHITECT_FRONTEND_REACT.md` | `GOOD_PRACTICES_FRONTEND_REACT.md` |

---

## 📝 **Dica Pro**

Use o **@mention** no Cursor para referenciar documentos específicos:
- `@.cursorrules` - Para implementação (sempre ativo)
- `@docs/ARCHITECT_BACKEND_GO.md` - Para planejamento
- `@docs/TECHNICAL_DEBT.md` - Para débitos técnicos e refactoring