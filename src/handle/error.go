package handle

//RedisError 错误信息.
type RedisError string

const (
	Redis_NoAuth   RedisError = "NOAUTH Authentication required"
	Redis_AuthFail            = "ERR invalid password"
)
