# kubectl updatecontext

[![GitHub Release](https://img.shields.io/github/release/borisputerka/updatecontext.svg?style=flat)](https://github.com/borisputerka/updatecontext/releases)
[![Go Report Card](https://goreportcard.com/badge/github.com/borisputerka/updatecontext)](https://goreportcard.com/report/github.com/borisputerka/updatecontext)

Plugin for kubectl that will create and update kubernetes contexts based on namespaces. For every namespace it will create context. After namespace is deleted, updatecontext will delete specific context.
