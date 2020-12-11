package config

var flagsOpts = []flagOpt{
	{
		optName:         FLAG_KEY_SERVER_HOST,
		optDefaultValue: "0.0.0.0",
		optUsage:        "server listen host",
	},
	{
		optName:         FLAG_KEY_SERVER_PORT,
		optDefaultValue: 8081,
		optUsage:        "server listen port",
	},
	{
		optName:         FLAG_KEY_MYSQL_HOST,
		optDefaultValue: "127.0.0.1",
		optUsage:        "mysql host",
	},
	{
		optName:         FLAG_KEY_MYSQL_PASSWORD,
		optDefaultValue: "",
		optUsage:        "mysql password",
	},
	{
		optName:         FLAG_KEY_MYSQL_USER,
		optDefaultValue: "root",
		optUsage:        "mysql user",
	},
	{
		optName:         FLAG_KEY_GIN_MODE,
		optDefaultValue: "debug",
		optUsage:        "gin mode",
	},
	{
		optName:         FLAG_KEY_LOG_LEVEL,
		optDefaultValue: "info",
		optUsage:        "log level",
	},
}
