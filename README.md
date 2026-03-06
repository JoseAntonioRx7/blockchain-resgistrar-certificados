# 🔗 Cert-Chain: Sistema de Certificados em Blockchain

> **Plataforma descentralizada para autenticação e verificação de certificados educacionais utilizando blockchain e criptografia digital.**

---

## 📋 Visão Geral do Projeto

**Cert-Chain** é uma solução inovadora para combater fraudes em certificados e diplomas educacionais. A startup implementa uma blockchain completa desenvolvida from scratch em Go, integrando autenticação criptográfica, banco de dados persistente e uma interface web intuitiva.

### Problema que Soluciona

- ❌ Certificados falsificados circulando no mercado educacional
- ❌ Impossibilidade de verificar autenticidade de diplomas rapidamente
- ❌ Falta de transparência entre instituições e empregadores
- ❌ Centralização de informações educacionais em órgãos únicos

### Solução Oferecida

✅ **Prova criptográfica de autenticidade** - Cada certificado gera um hash único e imutável  
✅ **Blockchain descentralizada** - Registro permanente e tamper-proof  
✅ **Verificação pública** - Qualquer pessoa pode validar um certificado sem intermediários  
✅ **Assinatura digital** - Sistema cryptográfico Ed25519 para autenticação institucional  
✅ **Interface intuitiva** - Dashboard web para registrar e verificar certificados  

---

## 🏗️ Arquitetura da Solução

```
┌─────────────────────────────────────────────────────────────┐
│                     Frontend (HTML/JS)                      │
│              (Interface de Usuário + Dashboard)             │
└──────────────────────┬──────────────────────────────────────┘
                       │ HTTP/REST
┌──────────────────────▼──────────────────────────────────────┐
│              Backend Go - API REST (Port :8080)             │
│  ┌─────────────┐    ┌──────────┐    ┌────────────────┐    │
│  │ /register   │    │ /verify  │    │ /list          │    │
│  └─────────────┘    └──────────┘    └────────────────┘    │
└──────────────────────┬──────────────────────────────────────┘
        ┌─────────────┬┴────────────┬─────────────┐
        │             │             │             │
┌───────▼────┐ ┌─────▼──────┐ ┌───▼────────┐ ┌─▼─────────────┐
│ Blockchain │ │ Database   │ │  Utils     │ │  Cryptography│
│ (Chain)    │ │ (PostgreSQL)│ │  (Helper)  │ │ (Ed25519)    │
└────────────┘ └────────────┘ └────────────┘ └──────────────┘
```

### Componentes Principais

#### 1. **Blockchain Engine** (`/blockchain`)
- **`chain.go`** - Núcleo da blockchain com persistência
- **`block.go`** - Estrutura de blocos individuais
- **`pow.go`** - Algoritmo Proof of Work (Difficulty = 4)
- **`transaction.go`** - Estrutura de transações (certificados)

#### 2. **API Layer** (`/api`)
- **`routes.go`** - Definição de endpoints REST
- **`handlers.go`** - Lógica de registrar, verificar e listar certificados

#### 3. **Database** (`/database`)
- **`db.go`** - Integração com PostgreSQL
- Criação automática de tabelas: `certificates` e `blocks`

#### 4. **Utilitários** (`/utils`)
- **`crypto.go`** - Geração de chaves e assinatura digital (Ed25519)
- **`hash.go`** - Hash SHA-256 de arquivos
- **`id.go`** - Geração de IDs únicos

#### 5. **Frontend** (`/web`)
- **`index.html`** - Interface de registro e verificação
- **`scripts.js`** - Lógica de interação
- **`style.css`** - Estilos responsivos

---

## ⚙️ Como Funciona

### Fluxo de Registro de Certificado

```
1. Usuário faz upload de arquivo + dados (aluno, instituição, curso)
   ↓
2. Backend calcula hash SHA-256 do arquivo
   ↓
3. Cria assinatura digital com chave privada da instituição (Ed25519)
   ↓
4. Salva certificado na tabela PostgreSQL
   ↓
5. Cria transação e adiciona ao bloco atual
   ↓
6. Minera novo bloco (Proof of Work - 4 zeros no início)
   ↓
7. Bloco salvo no PostgreSQL com histórico completo
   ↓
8. Retorna hash ao usuário para verificação futura
```

### Fluxo de Verificação

```
1. Usuário fornece hash do certificado
   ↓
2. Sistema busca o hash na blockchain
   ↓
3. Valida presença e integridade do bloco
   ↓
4. Retorna dados: aluno, instituição, curso, assinatura digital
```

---

## 🗂️ Estrutura de Diretórios

```
cert-chain/
├── main.go                 # Ponto de entrada da aplicação
├── go.mod                  # Dependências Go
├── go.sum
├── postgress.sql           # Scripts SQL para banco de dados
├── README.md               # Este arquivo
│
├── api/
│   ├── handlers.go         # Handlers dos endpoints REST
│   └── routes.go           # Definição das rotas
│
├── blockchain/
│   ├── block.go            # Estrutura do bloco
│   ├── chain.go            # Lógica da blockchain
│   ├── pow.go              # Proof of Work
│   └── transaction.go      # Estrutura de transação
│
├── database/
│   └── db.go               # Inicialização e gerenciamento do PostgreSQL
│
├── utils/
│   ├── crypto.go           # Criptografia Ed25519
│   ├── hash.go             # Hashing SHA-256
│   └── id.go               # Geração de IDs únicos
│
└── web/
    ├── index.html          # Interface web
    ├── scripts.js          # Lógica frontend
    └── style.css           # Estilos CSS
```

---

## 🚀 Como Começar

### Pré-requisitos

- **Go 1.25.6+** - [Download](https://golang.org/dl/)
- **PostgreSQL 12+** - [Download](https://www.postgresql.org/download/)
- **Git** - Controle de versão

### 1. Instalação do PostgreSQL

```bash
# Windows (usando instalador)
# Crie um banco de dados chamado: cert_chain
# Usuário: postgres
# Senha: tedcrypto1239
```

### 2. Clonar Repositório

```bash
git clone https://github.com/seu-usuario/cert-chain.git
cd cert-chain
```

### 3. Instalar Dependências Go

```bash
go mod download
go mod tidy
```

### 4. Setup do Banco de Dados

```bash
# Conecte ao PostgreSQL e execute:
psql -U postgres -d cert_chain -f postgress.sql
```

### 5. Configurar Credenciais

Edite em [`database/db.go`](database/db.go):

```go
connStr := "host=localhost port=5432 user=postgres password=SEU_PASSWORD dbname=cert_chain sslmode=disable"
```

### 6. Executar Aplicação

```bash
go run main.go
```

**Output esperado:**
```
Chaves geradas com sucesso!
Public Key: ...
Private Key: ...
Banco de Dados conectado com sucesso!
Tabela 'certificates' pronta para receber dados!
Tabela 'blocks' pronta para receber dados!
Blockchain funcionando!
Servidor rodando em http://localhost:8080
```

Acesse: **http://localhost:8080**

---

## 📡 API REST Endpoints

### 1. Registrar Certificado

**POST** `/register`

```javascript
FormData:
- student_name: "João da Silva"
- institution: "Universidade Federal"
- course: "Engenharia de Software"
- file: <arquivo PDF>

Resposta:
{
  "message": "Certificado registrado com sucesso",
  "hash": "a1b2c3d4e5f6...",
  "id": "550e8400-e29b-41d4-a716-446655440000"
}
```

### 2. Verificar Certificado

**GET** `/verify?hash=a1b2c3d4e5f6`

```javascript
Resposta:
{
  "found": true,
  "student_name": "João da Silva",
  "institution": "Universidade Federal",
  "course": "Engenharia de Software",
  "file_hash": "a1b2c3d4e5f6...",
  "signature": "hexadecimal_signature",
  "timestamp": 1709779200,
  "block_index": 1
}
```

### 3. Listar Certificados

**GET** `/list`

```javascript
Resposta:
{
  "total": 45,
  "certificates": [
    {
      "id": "...",
      "student_name": "Maria Santos",
      "institution": "UFRJ",
      "course": "Ciência da Computação",
      "file_hash": "...",
      "timestamp": 1709779200
    },
    ...
  ]
}
```

---

## 🔐 Segurança & Criptografia

### Mecanismos Implementados

| Componente | Tecnologia | Finalidade |
|-----------|-----------|-----------|
| **Hash de Arquivo** | SHA-256 | Integridade do certificado |
| **Assinatura Digital** | Ed25519 | Autenticação institucional |
| **Proof of Work** | SHA-256 + Difficulty 4 | Imutabilidade blockchain |
| **Chaves Institucionais** | Ed25519 Key Pair | Identificação única da instituição |

### Garantias

✅ Certificado hashado uma única vez  
✅ Assinado com chave privada da instituição  
✅ Registrado em bloco imutável  
✅ Verificável publicamente sem intermediários  

---

## 📊 Estrutura do Banco de Dados

### Tabela: `certificates`

```sql
id VARCHAR(64) PRIMARY KEY           -- ID único da transação
student_name VARCHAR(255)            -- Nome do aluno
institution VARCHAR(255)             -- Instituição de ensino
course VARCHAR(255)                  -- Curso realizado
file_hash TEXT UNIQUE NOT NULL       -- Hash SHA-256 do arquivo
signature TEXT                       -- Assinatura Ed25519
timestamp BIGINT                     -- Timestamp Unix
```

### Tabela: `blocks`

```sql
index INTEGER PRIMARY KEY            -- Número do bloco
timestamp BIGINT                     -- Quando foi minerado
prev_hash TEXT                       -- Hash do bloco anterior
hash TEXT                            -- Hash do bloco atual
nonce BIGINT                         -- Nonce para Proof of Work
transactions JSONB                   -- Array de transações
```

---

## 🧠 Funcionamento Técnico Detalhado

### Processo de Mineração (Proof of Work)

```go
const Difficulty = 4  // Requer 4 zeros no início do hash

// Mining:
// Hash esperado: 0000xxxxxxxxxxxx...
// A cada tentativa, nonce incrementa
// Até encontrar hash com prefixo correto
```

### Carregamento da Blockchain

1. Backend inicia
2. Conecta ao PostgreSQL
3. Busca todos os blocos (ORDER BY index ASC)
4. Se vazio, cria Bloco Gênese
5. Salva estado em memória + banco

### Adição de Novo Bloco

1. Coleta transações pendentes
2. Cria novo bloco com referência ao anterior
3. Executa Proof of Work (Mining)
4. Persiste no PostgreSQL

---

## 🎯 Roadmap & Próximas Features

### Phase 1 (Atual) ✅
- [x] Blockchain core em Go
- [x] API REST básica
- [x] Integração PostgreSQL
- [x] Criptografia Ed25519
- [x] Interface web simples

### Phase 2 (Próximo)
- [ ] Autenticação de instituições (JWT)
- [ ] Dashboard administrativo
- [ ] Exportação de relatórios
- [ ] API de integração para universidades
- [ ] Testes unitários automatizados

### Phase 3 (Futuro)
- [ ] Deploy em Kubernetes
- [ ] Sincronização com blockchain pública (Ethereum)
- [ ] Mobile app (React Native)
- [ ] Sistema de reputação institucional
- [ ] DAO governance

---

## 🧪 Testando Localmente

### 1. Registrar um Certificado

```bash
curl -X POST http://localhost:8080/register \
  -F "student_name=João Silva" \
  -F "institution=UFRJ" \
  -F "course=Engenharia" \
  -F "file=@certificado.pdf"
```

### 2. Verificar Certificado

```bash
curl "http://localhost:8080/verify?hash=a1b2c3d4e5f6..."
```

### 3. Listar Todos

```bash
curl http://localhost:8080/list
```

---

## 🔧 Tecnologias Utilizadas

| Stack | Tecnologia |
|-------|-----------|
| **Backend** | Go 1.25.6 |
| **Database** | PostgreSQL 12+ |
| **Criptografia** | SHA-256, Ed25519 |
| **Frontend** | HTML5, CSS3, Vanilla JavaScript |
| **Formato de Dados** | JSON, JSONB |
| **Protocolo** | HTTP/REST |

### Dependências Go

```
- github.com/lib/pq (PostgreSQL driver)
```

---

## 👥 Visão de Negócio

### Proposta de Valor

1. **Para Instituições** - Melhor reputação e segurança em seus certificados
2. **Para Alunos** - Portabilidade e segurança do histórico acadêmico
3. **Para Empregadores** - Verificação rápida e confiável de credenciais
4. **Para Sociedade** - Redução de fraudes educacionais

### Modelo de Monetização (Futuro)

- 💰 Taxa por certificado registrado (~0.1% - 0.5%)
- 💰 Plano enterprise para universidades
- 💰 API de verificação para recrutadoras
- 💰 Serviços de auditoria e compliance

---

## 📝 Licença

MIT License - Sinta-se livre para usar e contribuir!

---

## 📞 Contato & Suporte

Dúvidas ou sugestões? Entre em contato!

**GitHub Issues**: [cert-chain/issues](https://github.com/seu-usuario/cert-chain/issues)

---

## 🙏 Agradecimentos

Desenvolvido como projeto de aprendizado em blockchain com foco em aplicações reais.

