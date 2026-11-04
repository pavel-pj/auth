package client

//go:generate  sh -c "rm -rf mocks && mkdir -p mocks"
//go:generate minimock -i github.com/pavel-pj/platform_lib/pkg/client/db.TxManager -o ./mocks/ -s "_minimock.go"
