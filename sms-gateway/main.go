package main

import (
	"io"
	"log"
	"net/http"

	//"sms-gateway/pkg/pdutext"

	"github.com/fiorix/go-smpp/smpp"
	"github.com/fiorix/go-smpp/smpp/pdu/pdufield"
	"github.com/fiorix/go-smpp/smpp/pdu/pdutext"
)

func main() {

	tx := &smpp.Transceiver{
		Addr:       "_______",
		User:       "____",
		Passwd:     "_____",
		SystemType: "VMA",
	}
	conn := tx.Bind()
	// check initial connection status
	var status smpp.ConnStatus
	if status = <-conn; status.Error() != nil {
		log.Fatalln("Unable to connect, aborting:", status.Error())
	}
	log.Println("Connection completed, status:", status.Status().String())
	// example of connection checker goroutine
	go func() {
		for c := range conn {
			log.Println("SMPP connection status:", c.Status())
		}
	}()

	var text = "Ш ш	Ъ ъ	Э э	Ю ю	Я я	Ғ ғ	Ӣ ӣ	Қ қ	Ӯ ӯ	Ҳ ҳ	Ҷ ҷ"

	sm, err := tx.Submit(&smpp.ShortMessage{
		Src:                  "SAFECITY",
		Dst:                  "992931441244",
		DstList:              []string{},
		DLs:                  []string{},
		Text:                 pdutext.UCS2([]byte(text)),
		Validity:             0,
		Register:             pdufield.NoDeliveryReceipt,
		SourceAddrTON:        uint8(5),
		SourceAddrNPI:        uint8(0),
		DestAddrTON:          uint8(1),
		DestAddrNPI:          uint8(1),
		ESMClass:             uint8(0),
		ProtocolID:           0,
		PriorityFlag:         0,
		ScheduleDeliveryTime: "",
		ReplaceIfPresentFlag: 0,
		SMDefaultMsgID:       0x03,
		NumberDests:          0,
	})

	if err != nil {
		log.Fatalf("%v \n %v", sm, err)
	}
	// example of sender handler func
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		sm, err := tx.Submit(&smpp.ShortMessage{
			Src:            r.URL.Query().Get("src"),
			Dst:            r.URL.Query().Get("dst"),
			Text:           pdutext.Raw(r.URL.Query().Get("text")),
			Register:       pdufield.NoDeliveryReceipt,
			SourceAddrTON:  uint8(5),
			SourceAddrNPI:  uint8(0),
			DestAddrTON:    uint8(1),
			DestAddrNPI:    uint8(1),
			SMDefaultMsgID: 0x03,
			ESMClass:       uint8(0),
		})
		if err == smpp.ErrNotConnected {
			http.Error(w, "Oops.", http.StatusServiceUnavailable)
			return
		}
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		io.WriteString(w, sm.RespID())
		log.Println(sm.Resp())
		log.Println(sm.Resp().Header().Status)
		log.Println(sm)
	})
	log.Fatal(http.ListenAndServe(":3445", nil))
}
