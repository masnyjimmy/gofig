
URL := "github.com/masnyjimmy/gofig"

[script]
init-mod path:
    cd {{path}}
    go mod init {{URL + path}}