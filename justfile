
URL := "github.com/masnyjimmy/gofig"

[script]
init-mod path:
    cd {{path}}
    go mod init {{URL + path}}

update-version ver:
    git tag {{ver}}
    git push origin {{ver}}