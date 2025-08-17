package notifications

type StatusNotification struct {
	EvseId    int
	Status    string
	ErrorCode string
}

func NewStatusNotification(evseId int, status, errorCode string) StatusNotification {
	return StatusNotification{
		Status:    status,
		EvseId:    evseId,
		ErrorCode: errorCode,
	}
}
