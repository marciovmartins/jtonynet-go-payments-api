# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

---

## [Unreleased]
### Added

### Fixed

## [0.0.4] - 2024-20-10
### Added

- [Issue 4](https://github.com/jtonynet/go-payments-api/issues/8)
  - L2 finalizada
  - Adiciona a lógica para debitar o saldo de CASH se o saldo da categoria principal for insuficiente.
  - melhorias nos testes
  - melhorias na arquitetura

## [0.0.3] - 2024-20-10
### Added

- [Issue 3](https://github.com/jtonynet/go-payments-api/issues/6)
  - L1 finalizada
  - Testes de integração na routes, unitários de repository e routes adicionados
  -  Acertos no Github Actions
  - Adicionado Swagger
  - Adicionado hot reload
  


## [0.0.2] - 2024-09-10
### Added

- [Issue 2](https://github.com/jtonynet/go-payments-api/issues/4)
  - Criar estrutura base para a implementação das tarefas/cards L1 a L3


## [0.0.1] - 2024-05-10
### Added

- [Issue 2](https://github.com/jtonynet/go-payments-api/issues/2)
   - Banco Postgres
   - API
   - docker compose
   - Dockerfile da API


## [0.0.0] - 2024-05-10
### Added

- [Issue 1](https://github.com/users/jtonynet/projects/7/views/1?pane=issue&itemId=82288146)
  - [Kanban Project View Iniciado](https://github.com/users/jtonynet/projects/7) com o commit inicial. 
  - Documentação base: Readme Rico, [Diagramas Mermaid](https://github.com/jtonynet/go-products-api/tree/main#diagrams), ADRs: [0001: Registro de Decisões de Arquitetura (ADR)](./docs/architecture/decisions/registro-de-decisoes-de-arquitetura.md) e [0002: Go, Gin, Gorm e PostgreSQL com Arquitetura Hexagonal e TDD](./docs/architecture/decisions/0002-go-gin-gorm-e-postgres-com-arquitetura-hexagonal-tdd.md).
  - Sabemos o que fazer, graças às definições do arquivo __README.md__. Sabemos como fazer graças aos __ADRs__ e documentações vinculadas. Devemos nos organizar em estrutura __Kanban__, guiados pelo modelo Agile, em nosso __Github Project__, e dar o devido prosseguimento às tarefas.

[0.0.4]: https://github.com/jtonynet/go-payments-api/compare/v0.0.3...v0.0.4
[0.0.3]: https://github.com/jtonynet/go-payments-api/compare/v0.0.2...v0.0.3
[0.0.2]: https://github.com/jtonynet/go-payments-api/compare/v0.0.1...v0.0.2
[0.0.1]: https://github.com/jtonynet/go-payments-api/compare/v0.0.0...v0.0.1
[0.0.0]: https://github.com/jtonynet/go-payments-api/releases/tag/v0.0.0
