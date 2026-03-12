🔗 TTLedger
Blockchain-Based Certificate & AI Audit System

TTLedger é uma infraestrutura de blockchain privada projetada para emissão, gestão e auditoria de certificados acadêmicos.

A plataforma utiliza uma arquitetura Non-Custodial, garantindo que as instituições mantenham controle total sobre suas chaves criptográficas, combinada com uma camada de Inteligência Artificial para auditoria e monitoramento de integridade da rede.

O objetivo é combater fraudes em diplomas e certificados educacionais, oferecendo verificação pública, segura e instantânea.

🚀 Tecnologias Utilizadas
Backend

Golang (Go) — Alta performance e concorrência para o nó da rede

PostgreSQL — Persistência de dados, estados e logs

Blockchain

Estrutura de blocos encadeados

Hashing SHA-256

Assinaturas digitais ECDSA

Sistema de mineração de blocos

Inteligência Artificial

Google Gemini 1.5 Flash

Análise proativa da blockchain

Detecção de padrões suspeitos e fraudes

Frontend

HTML5

CSS3

JavaScript Vanilla

Foco em:

performance

simplicidade

zero dependências pesadas

Segurança

Autenticação via JWT

RBAC (Role-Based Access Control)

Headers de segurança contra ataques web

✨ Funcionalidades Principais
🏛️ Governança e Expansão de Rede
Admin Console

O administrador da rede pode provisionar novas instituições educacionais dentro da rede blockchain.

Provisionamento On-the-Fly

O sistema gera pares de chaves criptográficas (Pública/Privada) dinamicamente para novos nós da rede.

Isso garante:

descentralização

segurança

autonomia institucional

🎓 Gestão de Certificados
Emissão Criptográfica

Cada certificado emitido é registrado em um bloco da blockchain, contendo:

hash do certificado

hash do bloco anterior

assinatura digital da instituição emissora

Arquitetura Non-Custodial

As chaves privadas nunca são armazenadas em texto claro, garantindo soberania criptográfica das instituições.

Verificação Pública

Qualquer pessoa pode verificar a autenticidade de um certificado através do hash registrado na blockchain.

🤖 Auditoria Inteligente com IA

O TTLedger possui uma camada de auditoria baseada em IA para monitoramento contínuo da integridade da rede.

Detecção de Anomalias

A IA analisa padrões como:

emissões massivas de certificados

inconsistência em nomes ou instituições

padrões suspeitos de atividade

Relatórios em Tempo Real

Dashboard administrativo que traduz dados técnicos da blockchain em insights acionáveis para o administrador da rede.

🏗️ Arquitetura do Projeto
/ttledger
│
├── api/                # Handlers HTTP e definição de rotas
├── blockchain/         # Lógica core da blockchain e mineração
├── database/           # Conexão e migrations do PostgreSQL
├── internal/
│   └── ia/             # Integração com Google Gemini
│
├── utils/              # Funções criptográficas e utilitários
├── web/                # Interface administrativa e verificação
│
├── main.go             # Ponto de entrada da aplicação
└── .env                # Variáveis de ambiente


🛡️ Segurança e Privacidade

O projeto segue boas práticas de segurança para redes privadas:

Headers de Segurança

Proteção contra ataques comuns da web como:

XSS

Clickjacking

Content Injection

Controle de Acesso (RBAC)

Separação clara entre:

Administrador da Rede

Instituições emissoras

Auditoria Imutável

Após minerado, um bloco torna-se imutável.

Qualquer tentativa de manipulação:

altera o hash

quebra a cadeia

é detectada pela camada de auditoria.

📈 Roadmap de Evolução
Implementado

 Nó core da blockchain

 Estrutura de blocos encadeados

 Sistema de expansão de rede

 Integração com IA Gemini para auditoria

Em Desenvolvimento

 Smart Contracts para renovação automática de certificados

 Identidade descentralizada para alunos

 API pública de verificação

 App mobile para carteira de certificados

🌍 Impacto Esperado

O TTLedger busca resolver problemas críticos no sistema educacional:

Redução de fraudes em diplomas

Confiança digital em credenciais acadêmicas

Verificação instantânea de certificados

Interoperabilidade entre instituições

👤 Autor

José Antonio Ramos da Silva
Estudante de Engenharia da Computação

Interesses:

Blockchain

Inteligência Artificial

Segurança Digital

Infraestrutura Distribuída

GitHub
https://github.com/JoseAntonioRx7

📜 Licença

Este projeto está sob a licença MIT License.

🤝 Contribuição

Contribuições são bem-vindas.

Você pode colaborar com:

melhorias no código

auditoria de segurança

novas funcionalidades

integrações educacionais

📬 Contato

GitHub
https://github.com/JoseAntonioRx7
