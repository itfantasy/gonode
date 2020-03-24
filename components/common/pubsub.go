package common

// the publish-subscrib equipment
type ISubscribable interface {
	Subscribe(string)
	BindSubscriber(ISubscriber)
}

type ISubscriber interface {
	OnSubscribe(string)
	OnSubMessage(string, string)
	OnSubError(string, error)
}

type IPublisher interface {
	Publish(string, string) error
}
