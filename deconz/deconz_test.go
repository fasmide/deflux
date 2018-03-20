package deconz

import (
	"fmt"
	"log"
	"net/http"
	"testing"

	"github.com/gorilla/websocket"
)

const httpListen string = `localhost:44431`

func TestWebsocket(t *testing.T) {
	// we start of by creating a webserver
	src := &http.Server{Addr: httpListen}

	http.HandleFunc("/ws", testHandler)

	go func() {
		err := src.ListenAndServe()
		if err != nil {
			t.Logf("unable to listen for http requests: %s", err)
			t.Fail()
		}
	}()

	cc := Config{
		wsAddr: fmt.Sprintf("ws://%s/ws", httpListen),
	}
	handler := func(e interface{}) {
		t.Logf("I have received: %+v", e)
		err := src.Shutdown(nil)
		if err != nil {
			t.FailNow()
		}
	}

	d := Deconz{Config: cc}
	d.AddHandler(handler)

	// this blocks
	go func() {
		err := d.Connect()
		if err != nil {
			t.Fail()
		}
	}()

}

var upgrader = websocket.Upgrader{} // use default options

func testHandler(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	c.WriteMessage(websocket.TextMessage, []byte(temperatureEventPayload))
	c.WriteMessage(websocket.TextMessage, []byte(humidityEventPayload))
	c.WriteMessage(websocket.TextMessage, []byte(pressureEventPayload))
	c.WriteMessage(websocket.TextMessage, []byte(smokeDetectorNoFireEventPayload))
	c.WriteMessage(websocket.TextMessage, []byte(floodDetectorFloodDetectedEventPayload))
	c.Close()
}
