package config

type Config struct {
	Mysql *Mysql
	Gin   *Gin
	Log   *Log
}

type Mysql struct {
	Host     string
	Database string
	User     string
	Password string
}

type Gin struct {
	Mode string
}

type Log struct {
	Level string
}
