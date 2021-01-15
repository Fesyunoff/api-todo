package alert

type Alert interface {
	SentMessage(msg string)
}
