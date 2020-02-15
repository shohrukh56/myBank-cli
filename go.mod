module github.com/shohrukh56/myBank-cli

go 1.13

require (
	github.com/mattn/go-sqlite3 v2.0.3+incompatible
	github.com/shohrukh56/myBank-core v0.0.0-20200210093103-3532bae9c046
)

replace github.com/shohrukh56/myBank-core => ../myBank-core
