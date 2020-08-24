module github.com/aibotsoft/trade

go 1.14

require (
	github.com/aibotsoft/gen v0.0.0-20200531091936-c4d5d714bf82
	github.com/aibotsoft/micro v0.0.0-20200606052507-83958c4d3f36
	github.com/google/uuid v1.1.1
	github.com/stretchr/testify v1.6.0
	go.uber.org/zap v1.15.0
)

replace github.com/aibotsoft/micro => ../micro
replace github.com/aibotsoft/gen => ../gen