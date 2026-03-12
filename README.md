🔗 TTLedger – Blockchain Certificate Infrastructure






Infraestrutura de blockchain privada para emissão e verificação de certificados acadêmicos, desenvolvida em Go (Golang).

O sistema registra certificados em uma cadeia de blocos imutável, utilizando criptografia SHA-256 e assinaturas digitais, com uma camada adicional de auditoria baseada em Inteligência Artificial para monitoramento da rede.

Projeto desenvolvido para aprofundar conhecimentos em:

blockchain

criptografia

sistemas distribuídos

arquitetura backend em Go

🚀 Funcionalidades Implementadas

• Estrutura de Blockchain:
Implementação de blocos encadeados contendo Hash, PrevHash, Timestamp e Data.

• Mineração de Blocos:
Algoritmo de Proof of Work (PoW) para validação e adição de novos blocos.

• Registro de Certificados:
Cada certificado é convertido em hash criptográfico e armazenado na blockchain.

• Assinaturas Digitais:
Uso de ECDSA para autenticação da instituição emissora.

• API REST Integrada:

GET /blocks
Retorna todos os blocos da blockchain.

POST /certificates
Registra um novo certificado na rede.


• Verificação de Integridade:
Mecanismo de validação da blockchain garantindo consistência entre blocos.

• Auditoria com Inteligência Artificial:
Integração com Google Gemini para análise de padrões suspeitos na emissão de certificados.

⚙️ Tecnologias

• Linguagem: Go (Golang)

• Blockchain:

Hashing SHA-256

Assinaturas digitais ECDSA

Proof of Work

• Banco de Dados:
PostgreSQL

• Inteligência Artificial:
Google Gemini API

• Arquitetura:
REST API

🏗 Estrutura do Projeto
ttledger/

api/            # Rotas e handlers HTTP
blockchain/     # Lógica da blockchain
database/       # Conexão com PostgreSQL
internal/ia/    # Integração com IA (Gemini)

utils/          # Funções criptográficas
web/            # Interface administrativa

main.go         # Entry point da aplicação

⚙️ Como Executar
Clone o repositório
git clone https://github.com/JoseAntonioRx7/ttledger.git

Acesse a pasta do projeto
cd ttledger

Configure as variáveis de ambiente

Crie um arquivo .env:

GEMINI_API_KEY=sua_chave

DB_URL=postgres://user:pass@localhost:5432/ttledger?sslmode=disable

PORT=8080

Execute a aplicação
go mod tidy

go run main.go

🔐 Segurança

O sistema implementa mecanismos de segurança baseados em criptografia:

• Hashing SHA-256 para integridade dos dados
• Assinaturas ECDSA para autenticação das instituições
• Estrutura imutável de blocos para impedir alterações históricas
• Auditoria inteligente para detectar anomalias na rede

📈 Roadmap

• Implementar rede P2P entre instituições
• Implementar identidade descentralizada (DID)
• Criar API pública de verificação de certificados
• Desenvolver carteira digital de certificados para alunos

👨‍💻 Autor

José Antonio Ramos da Silva

Estudante de Engenharia da Computação
Focado em:

Blockchain

Inteligência Artificial

Segurança Digital

Sistemas Distribuídos

GitHub:
https://github.com/JoseAntonioRx7
