# kubectl updatecontext

[![GitHub Release](https://img.shields.io/github/release/borisputerka/updatecontext.svg?style=flat)](https://github.com/borisputerka/updatecontext/releases)
[![Go Report Card](https://goreportcard.com/badge/github.com/borisputerka/updatecontext)](https://goreportcard.com/report/github.com/borisputerka/updatecontext)

Plugin for kubectl that will create and update kubernetes contexts based on namespaces. For every namespace it will create context. After namespace is deleted, updatecontext will delete specific context. 
You can then switch between contexts using [kubectx](https://github.com/ahmetb/kubectx) or using [fzf](https://github.com/junegunn/fzf) (see instruction below). Once you have contexts created, you will no longer need to use `kubectx` with `kubens` when accessing namespace in different cluster.

# Installation and usage

> **Info**: This plugin is not yet in [krew-index](https://github.com/kubernetes-sigs/krew-index) repository. Please use installation as described below

1. Install plugin
    * using krew manifest in this repository
        ```
        $ kubectl krew install --manifest=deploy/krew/plugin.yaml
        ```
    or
    * using `Makefile`
        ```
        $ make bin
        ```

2. Use plugin
    ```
    $ kubectl updatecontext
    ```

# Use fzf to switch contexts

Add these lines into your `.bashrc` or `.zshrc`
```
#fzf inline alias
alias _inline_fzf="fzf --multi --ansi -i -1 --height=50% --reverse -0 --header-lines=1 --inline-info --border"

#kubernetes contexts switcher
kcs() {
    local context="$(kubectl config get-contexts | _inline_fzf | awk '{print $1}')"
    eval kubectl config set current-context "${context}"
}
```

Now use `kcs` within your terminal.
```
$ kcs
```

