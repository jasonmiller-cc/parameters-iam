module github.com/jasonmiller-cc/parameters-iam

go 1.27.1

require github.com/jasonmiller-cc/parameters-core v0.0.0

require (
	github.com/golang-jwt/jwt/v5 v5.3.1 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/jasonmiller-cc/parameters-core => ../parameters-core
