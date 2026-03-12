🛡️ TTLedger: Blockchain-Based Certificate & AI Audit System
O TTLedger é uma infraestrutura de rede blockchain privada projetada para a emissão, gestão e auditoria de certificados acadêmicos. O sistema utiliza uma arquitetura Non-Custodial, garantindo que as instituições possuam total controle sobre suas chaves criptográficas, aliada a uma camada de Inteligência Artificial para monitoramento de integridade da rede.

🚀 Tecnologias Core
Backend: Go (Golang) - Alta performance e concorrência para o nó da rede.

Blockchain: Estrutura de blocos encadeados com hashing SHA-256 e assinaturas ECDSA.

IA de Auditoria: Google Gemini 1.5 Flash - Análise proativa de padrões e detecção de fraudes.

Database: PostgreSQL - Persistência robusta de estados e logs.

Frontend: HTML5, CSS3 e JavaScript Vanilla (Foco em performance e zero dependências pesadas).

Segurança: Autenticação via JWT com controle de acesso baseado em funções (RBAC).

✨ Funcionalidades Principais
1. Governança & Expansão de Rede
Admin Console: O administrador da rede pode provisionar novas instituições.

Provisionamento On-the-fly: Geração dinâmica de pares de chaves (Pública/Privada) para novos nós, garantindo a descentralização da confiança.

2. Gestão de Certificados
Emissão Criptográfica: Cada certificado é minerado em um bloco com o hash do bloco anterior.

Non-Custodial: As chaves privadas nunca são armazenadas em texto claro, garantindo a soberania da instituição emissora.

Verificação Instantânea: Validação pública de autenticidade através do hash do certificado.

3. Auditoria Inteligente (IA Gemini)
Detecção de Anomalias: IA integrada que analisa o fluxo de blocos em busca de comportamentos suspeitos (ex: emissões em massa, nomes inconsistentes).

Relatórios em Tempo Real: Dashboard de auditoria que traduz dados técnicos da blockchain em insights acionáveis para o administrador.

🛠️ Arquitetura do Projeto
Plaintext
/ttledger
  ├── /api             # Handlers HTTP e definição de rotas
  ├── /blockchain      # Lógica core da corrente e mineração
  ├── /database        # Migrations e conexão com Postgres
  ├── /internal/ia     # Motor de integração com Google Gemini
  ├── /utils           # Ferramentas criptográficas e auxiliares
  ├── /web             # Interface administrativa e de verificação
  ├── main.go          # Ponto de entrada do sistema
  └── .env             # Variáveis sensíveis (API Keys, DB URLs)
🔧 Configuração e Instalação
Clone o repositório:

Bash
git clone https://github.com/seu-usuario/ttledger.git
Configure as Variáveis de Ambiente (.env):
Crie um arquivo .env na raiz com:

Snippet de código
GEMINI_API_KEY=sua_chave_aqui
DB_URL=postgres://user:pass@localhost:5432/cert_chain?sslmode=disable
PORT=8080
Instale as dependências e rode o projeto:

Bash
go mod tidy
go run main.go
🛡️ Segurança e Privacidade
O projeto segue as diretrizes de segurança para redes privadas:

Headers de Segurança: Proteção contra ataques comuns de web.

Segregação de Roles: Distinção clara entre o que um Administrador de Rede e uma Instituição podem realizar.

Auditoria Imutável: Uma vez minerado, o registro é permanente, e qualquer tentativa de manipulação é detectada pela IA.

📈 Roadmap de Evolução
[x] Implementação do Nó Core e Blockchain

[x] Sistema de Expansão de Rede (Admin)

[x] Integração com IA Gemini para Auditoria

[ ] Implementação de Smart Contracts para renovação automática

[ ] App Mobile para carteira de certificados do aluno

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
