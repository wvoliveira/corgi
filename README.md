# Corgi

[![Lint](https://github.com/wvoliveira/corgi/actions/workflows/server.lint.yml/badge.svg)](https://github.com/wvoliveira/corgi/actions/workflows/server.lint.yml)
[![Test](https://github.com/wvoliveira/corgi/actions/workflows/server.test.yml/badge.svg)](https://github.com/wvoliveira/corgi/actions/workflows/server.test.yml)
[![Language](https://img.shields.io/badge/Language-pt--br-blue)](./README_pt-br.md)


Corgi is a link shortener system.

## Features

* **Users** - Registration/Authentication with new users via social network or email/password.
* **Easy** - Corgi is easy and fast. Insert a giant link and get a shortened link.
* **Your own domain** - Reduce links using your own domain and increase click-through rate.
* **Groups** - Manage links as a group, assigning roles to who can change and view information about the links.
* **API** - Use one of the available APIs to manage links effectively.
* **Statistics** - Check the amount of clicks on shortened linkdb.Debug()
* **Shortener** - Use any link, no matter the size. Corgi will always shorten it.
* **Manage** - Optimize and customize each link to take advantage. Use an alias, affiliate programs, create QR code and much more..

Use your own infrastructure to install this link shortener. With several features that will bring you more information about your users.

## Install

Prerequisites:
- Go 1.18+
- Node 16+
- NPM 8+

Build yourself:

```bash
$ make
```

And run:

```bash
$ ./corgi
```
