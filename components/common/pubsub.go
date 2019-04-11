package common

// the publish-subscrib equipment
type IPubSubEquip interface {
	Publish(string, string) error
	Subscribe(string)
	BindSubscriber(ISubscriber)
}

type ISubscriber interface {
	OnSubscribe(string)
	OnSubMessage(string, string)
	OnSubError(string, error)
}
