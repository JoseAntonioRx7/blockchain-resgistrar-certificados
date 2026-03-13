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

- **Arquitetura Non-Custodial:**
  A infraestrutura não retém as chaves privadas. As instituições emissoras possuem total soberania sobre suas credenciais de assinatura.

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
