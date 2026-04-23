# Feature-Spec Workflow

## O que é a "feature-spec mencionada"
Arquivo `.md` na pasta `features/` que descreve uma funcionalidade a ser implementada. Contém: título, descrição, User Story, Acceptance Criteria e/o Link JIRA.

---

## Passo 1: Gerar Cenários Gherkin
Para cada Acceptance Criteria, criar um cenário de teste no formato:
```
Cenário: [Nome do cenário]
Dado que [pré-condição]
Quando [ação do usuário]
Então [resultado esperado]
```

**Exemplo:**
```
Cenário: Criar usuário com email válido
Dado que o formulário de registro está aberto
Quando o usuário preenche "email" com "joao@exemplo.com"
Então o sistema aceita o email sem erros
```

---

## Passo 2: Agrupar Cenários
Organizar cenários em dois grupos:
- **Positivo (Happy Path):** Cenários que representam o fluxo desejável/sucesso
- **Negativo (Erro):** Cenários que validam tratamento de erros e edge cases

Manter os cenários juntos no mesmo arquivo (não separar em arquivos distintos).

---

## Passo 3: Gerar Arquivo BDD Feature Scenario
Criar arquivo `{xx}-{feature_name}.md` na pasta da feature com estrutura:

```markdown
# {Feature Name}

## User Story
[Nome] Como [tipo] Eu quero [ação] Para [benefício]

##cenário Positivo
### Cenário 1: [Nome]
Given [contexto]
When [ação]
Then [resultado]

## Cenário de Erro
### Cenário 2: [Nome]
Given [contexto]
When [ação]
Then [resultado]

## Especificação Técnica
- Endpoint: `POST /api/users`
- Request: `{ "name": string, "email": string }`
- Response: `201 Created` ou `400 Bad Request`
- Validações: email formato válido, name obrigatório
```

**Nomear com prefixo numérico** (`01-create-user.md`, `02-update-user.md`) para ordenação.

---

## Passo 4: Analisar Bootstrap (00-initial-setup.md)
Antes de implementar as issues, identificar o que é necessário para rodar os testes **sem dependências entre si**:

- **Fixtures/Seeds:** Dados base (usuários, roles, configurações) que Toda Issue precisa
- **Mocks:** Stub de serviços externos (email, pagamento, etc.)
- **Setup Comum:** Criar tabelas, limpar banco, configurar variáveis de ambiente

Gerar arquivo `00-initial-setup.md` listando:
```
## Dependências Globais
- Usuário admin padrão
- Roles: ADMIN, USER
- Configuração de email fake

## Fixtures Necessárias
- `users.json`: 3 usuários para login
- `roles.json`: 2 roles base
```