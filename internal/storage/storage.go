package storage

type BotStorage interface {
	Insert(model interface{}) (interface{}, error)
	Select(model interface{}, tableName string, condition string) error
}
