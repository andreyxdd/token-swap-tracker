package consumer

type Config struct {
	Brokers string
	Topic   string
	GroupId string
	Debug   bool
}
