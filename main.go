package main

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"github.com/pion/webrtc/v3"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/signaling", signalingHandler)
	log.Fatal(http.ListenAndServe(":4444", nil))
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func signalingHandler(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()
	for {

		// client connected create offer
		pc1, _ := webrtc.NewPeerConnection(webrtc.Configuration{
			ICEServers: []webrtc.ICEServer{
				{URLs: []string{"stun:stun.l.google.com:19302"}},
			},
		})

		pc1.OnConnectionStateChange(func(state webrtc.PeerConnectionState) {
			log.Println("connection state change to: " + state.String())
		})

		pc1.OnICEConnectionStateChange(func(state webrtc.ICEConnectionState) {
			log.Println("ice connection state change to: " + state.String())
		})

		pc1.CreateDataChannel("dc", nil)

		initialOffer, err := pc1.CreateOffer(nil)
		if err != nil {
			panic(err)
		}

		gathering := webrtc.GatheringCompletePromise(pc1)

		if err = pc1.SetLocalDescription(initialOffer); err != nil {
			return
		}

		<-gathering

		initialOffer = *pc1.LocalDescription()

		offer, _ := json.Marshal(initialOffer)

		c.WriteMessage(websocket.TextMessage, offer)

		// wait for answer
		_, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}

		answer := webrtc.SessionDescription{}
		json.Unmarshal(message, &answer)

		if err = pc1.SetRemoteDescription(answer); err != nil {
			log.Println("failed to set remote description: " + err.Error())
			break
		}

		// wait for new offer
		_, message, err = c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}

		newOffer := webrtc.SessionDescription{}
		json.Unmarshal(message, &newOffer)

		pc1.SetRemoteDescription(newOffer)

		answer, err = pc1.CreateAnswer(nil)
		if err != nil {
			log.Println("failed to create new answer: " + err.Error())
			break
		}

		answerBytes, _ := json.Marshal(answer)

		c.WriteMessage(websocket.TextMessage, answerBytes)

		select {}
	}
}
