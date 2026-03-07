# 🔗 TTLedger

### Trusted Transcript Ledger

> **Infraestrutura descentralizada para autenticação e verificação de certificados educacionais utilizando blockchain e criptografia digital.**

---

# 📋 Visão Geral

**TTLedger** é uma plataforma baseada em **blockchain** projetada para registrar, autenticar e verificar **certificados educacionais e credenciais acadêmicas** de forma segura, transparente e imutável.

O projeto busca resolver um problema crescente no mercado educacional global: **fraudes em diplomas e certificados**.

A solução combina:

* blockchain própria desenvolvida em **Golang**
* assinaturas criptográficas **Ed25519**
* banco de dados **PostgreSQL**
* interface web de verificação pública

O objetivo é construir uma **infraestrutura global confiável para verificação de credenciais educacionais**.

---

# 🚨 Problema

Fraudes em certificados e diplomas são um problema mundial.

Principais desafios atuais:

* falsificação de diplomas
* certificados digitais manipulados
* verificação manual por universidades
* processos lentos de autenticação
* baixa interoperabilidade entre instituições

Empresas e universidades frequentemente precisam **verificar manualmente credenciais educacionais**, gerando:

* atrasos
* custos administrativos
* risco de fraude

---

# 💡 Solução

O **TTLedger** cria um **registro descentralizado e imutável de certificados educacionais**.

Cada certificado gera um **hash criptográfico único** que é registrado na blockchain.

Esse registro permite que qualquer pessoa verifique a autenticidade de um certificado sem depender da instituição emissora.

Principais características:

* 🔐 **Prova criptográfica de autenticidade**
* ⛓️ **Registro imutável em blockchain**
* 🌍 **Verificação pública e descentralizada**
* 🏫 **Assinatura digital institucional**
* 📄 **Validação instantânea de certificados**

---

# ⚙️ Como Funciona

### Fluxo de Registro de Certificado

```
1. Instituição emite certificado
   ↓
2. Sistema gera hash SHA-256 do documento
   ↓
3. Instituição assina digitalmente com chave privada (Ed25519)
   ↓
4. Certificado registrado na blockchain
   ↓
5. Bloco minerado utilizando Proof of Work
   ↓
6. Registro torna-se permanente
```

---

### Fluxo de Verificação

```
1. Usuário fornece o hash do certificado
   ↓
2. Sistema consulta a blockchain
   ↓
3. Verifica integridade do bloco
   ↓
4. Retorna dados do certificado registrado
```

Resultado:

```
✔ Certificado válido
ou
✖ Certificado inválido
```

---

# 🏗️ Arquitetura do Sistema

```
┌──────────────────────────────────────┐
│        Frontend (HTML / JS)          │
│     Dashboard + Verificação Web      │
└───────────────┬──────────────────────┘
                │
            HTTP / REST
                │
┌───────────────▼──────────────────────┐
│         Backend API (Go)             │
│     Registro | Verificação | List    │
└───────────────┬──────────────────────┘
                │
     ┌──────────┼───────────┬──────────┐
     │          │           │          │
┌────▼────┐ ┌───▼────┐ ┌────▼────┐ ┌───▼────────┐
│Blockchain│ │Database│ │ Utils  │ │Cryptography│
│ Engine   │ │Postgres│ │Helper  │ │Ed25519     │
└──────────┘ └────────┘ └────────┘ └────────────┘
```

---

# 🧩 Componentes do Sistema

## Blockchain Engine

Responsável por:

* criação de blocos
* validação da cadeia
* mineração (Proof of Work)
* registro imutável de transações

Arquivos principais:

```
/blockchain
  block.go
  chain.go
  pow.go
  transaction.go
```

---

## API Backend

Camada responsável por expor funcionalidades via REST.

Endpoints principais:

```
POST /register
GET /verify
GET /list
```

Arquivos:

```
/api
  handlers.go
  routes.go
```

---

## Banco de Dados

Persistência de certificados e blocos.

Tecnologia:

```
PostgreSQL
```

Tabelas principais:

* certificates
* blocks

---

## Utilitários

Funções auxiliares do sistema:

```
/utils
  crypto.go
  hash.go
  id.go
```

Responsáveis por:

* assinatura digital
* geração de hash
* geração de IDs únicos

---

# 📂 Estrutura do Projeto

```
ttledger/
│
├── main.go
├── go.mod
├── go.sum
├── README.md
│
├── api/
│   ├── handlers.go
│   └── routes.go
│
├── blockchain/
│   ├── block.go
│   ├── chain.go
│   ├── pow.go
│   └── transaction.go
│
├── database/
│   └── db.go
│
├── utils/
│   ├── crypto.go
│   ├── hash.go
│   └── id.go
│
└── web/
    ├── index.html
    ├── scripts.js
    └── style.css
```

---

# 🔐 Segurança

O TTLedger utiliza múltiplas camadas de segurança.

| Tecnologia    | Função                      |
| ------------- | --------------------------- |
| SHA-256       | Integridade de certificados |
| Ed25519       | Assinatura digital          |
| Proof of Work | Imutabilidade da blockchain |
| Hash chaining | Integridade dos blocos      |

Garantias:

* certificados não podem ser alterados
* registros são permanentes
* verificação pública e transparente

---

# 📡 API REST

### Registrar Certificado

```
POST /register
```

Dados enviados:

```
student_name
institution
course
file
```

Resposta:

```
{
 "message": "Certificado registrado",
 "hash": "...",
 "id": "..."
}
```

---

### Verificar Certificado

```
GET /verify?hash=HASH
```

Resposta:

```
{
 "found": true,
 "student_name": "...",
 "institution": "...",
 "course": "...",
 "block_index": 1
}
```

---

### Listar Certificados

```
GET /list
```

---

# 🚀 Como Executar

### Pré-requisitos

* Go 1.25+
* PostgreSQL 12+
* Git

---

### Instalar dependências

```
go mod download
go mod tidy
```

---

### Executar o projeto

```
go run main.go
```

Servidor disponível em:

```
http://localhost:8080
```

---

# 🎯 Roadmap

### Phase 1 (MVP)

* blockchain core
* API REST
* registro de certificados
* verificação pública

### Phase 2

* autenticação institucional
* dashboard administrativo
* geração automática de QR Code

### Phase 3

* rede de nós distribuídos
* integração com universidades
* API pública

### Phase 4

* integração com blockchain pública
* mobile app
* identidade digital descentralizada

---

# 💼 Visão de Negócio

O TTLedger pode operar como **SaaS para instituições educacionais**.

Possíveis modelos de receita:

* taxa por certificado registrado
* planos institucionais
* API para recrutadores
* auditoria de credenciais

---

# 🌍 Impacto

Benefícios esperados:

* redução de fraude educacional
* confiança digital em diplomas
* verificação instantânea de credenciais
* interoperabilidade entre instituições

---

# 👤 Autor

**José Antonio Ramos da Silva**
Estudante de Engenharia da Computação

Interesses:

* blockchain
* inteligência artificial
* segurança digital
* infraestrutura distribuída

---

# 📜 Licença

MIT License

---

# 🤝 Contribuição

Contribuições são bem-vindas.

Você pode colaborar com:

* melhorias no código
* auditoria de segurança
* novas funcionalidades
* integrações educacionais

---

# 📬 Contato

GitHub:

https://github.com/JoseAntonioRx7
