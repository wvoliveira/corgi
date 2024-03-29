# Corgi

[![Lint](https://github.com/wvoliveira/corgi/actions/workflows/lint.yaml/badge.svg)](https://github.com/wvoliveira/corgi/actions/workflows/lint.yaml)
[![Test](https://github.com/wvoliveira/corgi/actions/workflows/tests.yaml/badge.svg)](https://github.com/wvoliveira/corgi/actions/workflows/tests.yaml)

Corgi é um sistema encurtador de links.

## Recursos

* **Usuários** - Registro/Autenticação com novos usuários via rede social ou e-mail/senha.
* **Fácil** - Corgi é fácil e rápido. Insira um link gigante e pegue um link encurtado.
* **Seu próprio domínio** - Reduza os links usando seu próprio domínio e aumente a taxa de cliques.
* **Grupos** - Gerencie os links em grupo, atribuindo papéis de quem poderá alterar e visualizar informações sobre os links.
* **API** - Use uma das APIs disponíveis para gerenciar os links de forma eficaz.
* **Estatísticas** - Verifique a quantidade de cliques dos links encurtado.
* **Encurtador** - Use qualquer link, não importa o tamanho. Corgi sempre irá encurta-lo.
* **Gerencie** - Otimize e customize cada link para ter vantagens. Use um alias, programas de afiliados, crie QR code e muito mais.

Use sua própria infraestrutura para instalar esse encurtador de links. Com vários recursos que te trará mais informações sobre os seus usuários.

## Desenvolvimento local

### API

Requisitos para a API:

* Go 1.20+
* Docker 23+

Suba as dependências com o Docker (PostgreSQL e Redis):

```bash
make dev-dep
```

Copie as variáveis de ambiente e carregue no terminal:

```bash
make dev-env
export $(cat .env)
```

Inicie a API:

```bash
make dev-run
```

Há uma collection do Postman na pasta docs\/

### Frontend

Requisitos para o Frontend:

* Corgi API (passo anterior)
* Node 14+

Vá para a pasta "web" e instale as dependências:

```bash
cd web
npm install
```

Inicie no modo desenvolvedor:

```bash
npm run dev
```
