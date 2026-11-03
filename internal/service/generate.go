package service

//go:generate  sh -c "rm -rf mocks && mkdir -p mocks"
//go:generate  minimock  -i UserService -o ./mocks/ -s "_minimock.go"
//go:generate  minimock  -i LogService -o ./mocks/ -s "_minimock.go"
