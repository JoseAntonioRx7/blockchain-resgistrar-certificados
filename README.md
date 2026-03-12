# 🔗 TTLedger – Blockchain Certificate Infrastructure

![go](https://img.shields.io/badge/go-golang-blue)
![status](https://img.shields.io/badge/status-em%20desenvolvimento-green)
![license](https://img.shields.io/badge/license-MIT-blue)

Uma infraestrutura de **blockchain privada para emissão e verificação de certificados acadêmicos**, desenvolvida em **Go (Golang)**.

O sistema registra certificados em uma cadeia de blocos **imutável**, utilizando **criptografia SHA-256 e assinaturas digitais**, com uma camada adicional de **Inteligência Artificial para auditoria e monitoramento da rede**.

Projeto desenvolvido para aprofundar conhecimentos práticos em:

- tecnologias blockchain
- criptografia aplicada
- sistemas distribuídos
- arquitetura backend em Go

---

# 🚀 Funcionalidades Implementadas

- **Estrutura de Blockchain:**  
  Implementação de blocos encadeados contendo `Hash`, `PrevHash`, `Timestamp` e `Data`.

- **Mineração de Blocos:**  
  Algoritmo de **Proof of Work (PoW)** utilizado para validação e adição de novos blocos.

- **Registro de Certificados:**  
  Cada certificado é convertido em um **hash criptográfico** e armazenado na blockchain.

- **Assinaturas Digitais:**  
  Utilização de **ECDSA** para autenticação da instituição emissora.

- **API REST Integrada**

```
GET /blocks
Retorna todos os blocos da blockchain.

POST /certificates
Registra um novo certificado na rede.
```

- **Validação da Blockchain:**  
  Mecanismo que verifica a integridade da cadeia garantindo consistência entre blocos.

- **Auditoria com Inteligência Artificial:**  
  Integração com **Google Gemini** para análise de padrões suspeitos na emissão de certificados.

---

# ⚙️ Tecnologias

- **Linguagem:** Go (Golang)

- **Blockchain**
  - Hashing SHA-256
  - Assinaturas digitais ECDSA
  - Proof of Work

- **Banco de Dados**
  - PostgreSQL

- **Inteligência Artificial**
  - Google Gemini API

- **Arquitetura**
  - REST API

---

# 📦 Estrutura do Projeto

```
ttledger/

api/            # Rotas e handlers HTTP
blockchain/     # Lógica principal da blockchain
database/       # Conexão e migrations PostgreSQL
internal/ia/    # Integração com IA (Gemini)

utils/          # Funções criptográficas
web/            # Interface administrativa

main.go         # Entry point da aplicação
```

---

# ⚙️ Como Executar

### 1️⃣ Clone o repositório

```bash
git clone https://github.com/JoseAntonioRx7/ttledger.git
```

### 2️⃣ Acesse a pasta do projeto

```bash
cd ttledger
```

### 3️⃣ Configure as variáveis de ambiente

Crie um arquivo `.env` na raiz do projeto:

```
GEMINI_API_KEY=sua_chave

DB_URL=postgres://user:pass@localhost:5432/ttledger?sslmode=disable

PORT=8080
```

### 4️⃣ Execute a aplicação

```bash
go mod tidy
go run main.go
```

---

# 🔐 Segurança

O sistema implementa mecanismos de segurança baseados em criptografia:

- **Hashing SHA-256** para integridade dos dados
- **Assinaturas digitais ECDSA** para autenticação das instituições
- **Estrutura imutável de blocos** impedindo alterações históricas
- **Auditoria inteligente** para detecção de anomalias na rede

---

# 📈 Roadmap

- Implementar **rede P2P entre instituições**
- Implementar **identidade descentralizada (DID)**
- Criar **API pública de verificação de certificados**
- Desenvolver **carteira digital de certificados para alunos**

---

# 👨‍💻 Autor

José Antonio Ramos da Silva  

Estudante de Engenharia da Computação  

Interesses:

- Blockchain
- Inteligência Artificial
- Segurança Digital
- Sistemas Distribuídos

GitHub:  
https://github.com/JoseAntonioRx7

---

# 📜 Licença

MIT License
