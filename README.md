# Corgi

[![Lint](https://github.com/wvoliveira/corgi/actions/workflows/server.lint.yml/badge.svg)](https://github.com/wvoliveira/corgi/actions/workflows/server.lint.yml)
[![Test](https://github.com/wvoliveira/corgi/actions/workflows/server.test.yml/badge.svg)](https://github.com/wvoliveira/corgi/actions/workflows/server.test.yml)

Corgi é um sistema encurtador de links.

## Recursos

* **Usuários** - Registro/Autenticação com novos usuários via rede social ou e-mail/senha.
* **Fácil** - Corgi é fácil e rápido. Insira um link gigante e pegue um link encurtado.
* **Seu próprio domínio** - Reduza os links usando seu próprio domínio e aumente a taxa de cliques.
* **Grupos** - Gerencie os links em grupo, atribuindo papéis de quem poderá alterar e visualizar informações sobre os links.
* **API** - Use uma das APIs disponíveis para gerenciar os links de forma eficaz.
* **Estatísticas** - Verifique a quantidade de cliques dos links encurtadodb.Debug()
* **Encurtador** - Use qualquer link, não importa o tamanho. Corgi sempre irá encurta-lo.
* **Gerencie** - Otimize e customize cada link para ter vantagens. Use um alias, programas de afiliados, crie QR code e muito maidb.Debug()

Use sua própria infraestrutura para instalar esse encurtador de links. Com vários recursos que te trará mais informações sobre os seus usuários.

## Instalar

Requisitos:

* Go 1.20
* Node 8+
* Docker

Suba as dependências com o Docker (PostgreSQL e Redis):

```bash
make local-dep
```

Copie as variáveis de ambiente e carregue no terminal:

```bash
make local-env
```

Compile você mesmo:

```bash
make
```

E execute:

```bash
./corgi
```

Há uma collection do Postman na pasta docs\/
