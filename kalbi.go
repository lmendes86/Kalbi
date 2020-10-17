package kalbi

import (
	"github.com/KalbiProject/Kalbi/log"
	"github.com/KalbiProject/Kalbi/sip/dialog"
	"github.com/KalbiProject/Kalbi/sip/message"
	"github.com/KalbiProject/Kalbi/sip/transaction"
	"github.com/KalbiProject/Kalbi/transport"
	"github.com/KalbiProject/Kalbi/interfaces"
	"github.com/KalbiProject/Kalbi/sip/event"
)

//NewSipStack  creates new sip stack
func NewSipStack(Name string) *SipStack {
	stack := new(SipStack)
	stack.Name = Name
	stack.TransManager = transaction.NewTransactionManager()
	stack.TransportChannel = make(chan event.S)
    
	return stack
}

//SipStack has multiple protocol listning points
type SipStack struct {
	Name             string
	ListeningPoints  []interfaces.ListeningPoint
	OutputPoint      chan message.SipMsg
	InputPoint       chan message.SipMsg
	Alive            bool
	TransManager     *transaction.TransactionManager
	Dialogs          []dialog.Dialog
	RequestChannel   chan interfaces.Transaction
	ResponseChannel  chan interfaces.Transaction
	TransportChannel chan interfaces.SipEventObject
	sipListener      interfaces.SipListener
}

func (ed *SipStack) GetTransactionManager() *transaction.TransactionManager {
	return ed.TransManager
}

//CreateListenPoint creates listening point to the event dispatcher
func (ed *SipStack) CreateListenPoint(protocol string, host string, port int) transport.ListeningPoint {
	listenpoint := transport.NewTransportListenPoint(protocol, host, port)
	listenpoint.SetTransportChannel(ed.TransportChannel)
	ed.ListeningPoints = append(ed.ListeningPoints, listenpoint)
	return listenpoint
}

func (ed *SipStack) CreateRequestsChannel() chan interfaces.Transaction {
	Channel := make(chan interfaces.Transaction)
	ed.RequestChannel = Channel
	ed.TransManager.RequestChannel = Channel
	return Channel
}

func (ed *SipStack) CreateResponseChannel() chan interfaces.Transaction {
	Channel := make(chan interfaces.Transaction)
	ed.ResponseChannel = Channel
	ed.TransManager.ResponseChannel = Channel
	return Channel
}

func (ed *SipStack) SetSipListener(listener interfaces.SipListener){
	ed.sipListener = listener

}

func (ed *SipStack) IsAlive() bool {
	return ed.Alive
}

func (ed *SipStack) Stop() {
	log.Log.Info("Stopping SIPStack...")
	ed.Alive = false
}

//Start starts the sip stack
func (ed *SipStack) Start() {
	log.Log.Info("Starting SIPStack...")
	ed.TransManager.ListeningPoint = ed.ListeningPoints[0]
	ed.Alive = true
	for _, listeningPoint := range ed.ListeningPoints {
         go listeningPoint.Start()
	}

	for ed.Alive == true {
			msg := <-ed.TransportChannel
			event := ed.TransManager.Handle(msg)
			message := event.GetSipMessage
			if message.Req.StatusCode != nil {
                ed.sipListener.HandleResponses(event)
			}else if message.Req.Method != nil {
                ed.sipListener.HandleRequests(event)
			}

	}
}
