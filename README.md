# 🔗 TTLedger – Blockchain Certificate Infrastructure

![go](https://img.shields.io/badge/go-golang-blue)
![status](https://img.shields.io/badge/status-em%20desenvolvimento-green)
![license](https://img.shields.io/badge/license-MIT-blue)

Uma infraestrutura de **blockchain privada para emissão e verificação de certificados acadêmicos**, desenvolvida em **Go (Golang)**.

O sistema registra certificados em uma cadeia de blocos **imutável**, utilizando **criptografia SHA-256 e assinaturas digitais**, com uma camada adicional de **Inteligência Artificial para auditoria e monitoramento da rede**.

Projeto desenvolvido para aprofundar conhecimentos práticos em:

- Tecnologias Blockchain e Web3
- Criptografia Aplicada (Non-Custodial)
- Sistemas Distribuídos
- Arquitetura Backend em Go
- Integração de IA Generativa em Infraestrutura

---

# 🚀 Funcionalidades Implementadas

- **Arquitetura Non-Custodial:** A infraestrutura não retém as chaves privadas. As instituições emissoras possuem total soberania sobre suas credenciais de assinatura.

- **Estrutura de Blockchain:** Implementação de blocos encadeados contendo `Hash`, `PrevHash`, `Timestamp` e `Data`.

- **Mineração de Blocos:** Algoritmo de **Proof of Work (PoW)** utilizado para validação e adição de novos blocos à rede.

- **Registro de Certificados:** Cada certificado é convertido em um **hash criptográfico** e armazenado na blockchain.

- **Assinaturas Digitais:** Utilização de **ECDSA** para autenticação criptográfica da instituição emissora.

- **API REST Integrada:**

```text
GET /blocks
Retorna todos os blocos validados da blockchain.

POST /certificates
Registra e minera um novo certificado na rede.

GET /api/admin/audit-network
Aciona o Oráculo de IA para analisar os últimos registros e gerar um relatório de segurança em tempo real.

Validação da Blockchain: Mecanismo que verifica a integridade da cadeia garantindo consistência matemática entre os blocos.

Auditoria com Inteligência Artificial: Integração com a API do Google Gemini (1.5 Flash) para análise de padrões suspeitos e detecção de anomalias na emissão de certificados.

⚙️ Tecnologias
Linguagem: Go (Golang)

Blockchain

Hashing SHA-256

Assinaturas digitais ECDSA

Proof of Work (PoW)

Banco de Dados

PostgreSQL

Inteligência Artificial

Google Gemini 1.5 Flash API (Oráculo de Segurança)

Arquitetura

REST API

RBAC (Role-Based Access Control) com JWT

📦 Estrutura do Projeto

ttledger/

api/            # Rotas, middlewares e handlers HTTP
blockchain/     # Lógica core da blockchain, blocos e mineração
database/       # Conexão e migrations PostgreSQL
internal/ia/    # Integração e prompts para a IA (Gemini)
utils/          # Funções criptográficas e auxiliares
web/            # Interface administrativa e dashboard
main.go         # Entry point da aplicação

🔐 Segurança
O sistema implementa mecanismos de segurança focados em descentralização de confiança:

Abordagem Non-Custodial: Chaves privadas nunca trafegam em texto claro no banco de dados.

Hashing SHA-256: Garantia de integridade dos dados registrados.

Assinaturas digitais ECDSA: Autenticação inquestionável das instituições emissoras.

Estrutura imutável de blocos: Impede a manipulação de históricos passados.

Auditoria Inteligente (IA): Detecção proativa de comportamentos de emissão anômalos.

📈 Roadmap
Implementar rede P2P entre instituições

Implementar identidade descentralizada (DID)

Criar API pública de verificação de certificados

Desenvolver carteira digital de certificados para alunos

👨‍💻 Autor
José Antonio Ramos da Silva

Estudante de Engenharia da Computação

Interesses:

Blockchain & Web3

Inteligência Artificial

Segurança Digital

Sistemas Distribuídos

GitHub:

https://github.com/JoseAntonioRx7

📜 Licença
MIT License
