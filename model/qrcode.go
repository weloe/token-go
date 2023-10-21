package model

type QRCodeState int

// QRCode State
const (
	WaitScan    QRCodeState = 1
	WaitAuth    QRCodeState = 2
	ConfirmAuth QRCodeState = 3
	CancelAuth  QRCodeState = 4
	Expired     QRCodeState = 5
)

type QRCode struct {
	Id      string
	State   QRCodeState
	LoginId string
	Ticket  string
}

func NewQRCode(id string) *QRCode {
	return &QRCode{Id: id, State: WaitScan}
}
