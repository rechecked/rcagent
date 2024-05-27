package connect

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/rechecked/rcagent/internal/config"
)

type client struct {
	conn  *websocket.Conn
	mutex sync.Mutex
}

type Request struct {
	Id     int64         `json:"id"`
	Path   string        `json:"path"`
	Values config.Values `json:"values"`
}

type Response struct {
	Id       int64  `json:"id"`
	Response string `json:"response"`
}

func WSHandler() {

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	u := url.URL{
		Scheme: "wss",
		Host:   "localhost:8080",
		Path:   "/agents/connect/" + config.GetMachineId(),
	}
	log.Printf("connecting to %s", u.String())

	header := http.Header{}
	header.Add("X-API-Key", config.Settings.Manager.APIKey)

	c, _, err := websocket.DefaultDialer.Dial(u.String(), header)
	if err != nil {
		log.Print("dial:", err)
		return
	}
	defer c.Close()

	done := make(chan struct{})

	client := &client{
		conn: c,
	}

	// Read messages and send the request off to be processed
	go func() {
		defer close(done)
		for {

			r := Request{}
			err := c.ReadJSON(&r)
			if err != nil {
				config.LogDebug("wshandler read:", err)
				return
			}

			// Handle request and return it once it's finished
			go processRequest(client, r)

			config.LogDebugf("wshandler recv: %v", r)
		}
	}()

	// Handle shutdown
	for {
		select {
		case <-done:
			return
		case <-interrupt:
			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				config.LogDebug("wshandler close write:", err)
			}
			return
		}
	}

}

// Process request and once we have the data, we can then return the request
// back to the socket and forwarded
func processRequest(c *client, r Request) {

	data, err := callEndpoint(r)
	if err != nil {
		config.LogDebug(err)
	}

	b, err := json.Marshal(data)
	if err != nil {
		config.LogDebug(err)
	}

	resp := Response{
		Id:       r.Id,
		Response: string(b),
	}

	fmt.Println(resp)

	c.writeJSON(resp)
}

func (c *client) writeJSON(v interface{}) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.conn.WriteJSON(v)
}

func callEndpoint(r Request) (interface{}, error) {

	if r.Path == "" {
		return nil, errors.New("callEndpoint: no path given")
	}

	// Break string at main path (such as status or config api)
	// r.Path = status/memory/virtual
	// path[0] = status
	// path[1] = memory/virtual
	path := strings.SplitN(r.Path, "/", 2)

	endpoint := config.Endpoints[path[1]]
	if endpoint == nil {
		return nil, errors.New("callEndpoint: endpoint not found")
	}

	data := endpoint(r.Values)
	return data, nil
}
