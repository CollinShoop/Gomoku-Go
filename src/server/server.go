package server

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"math/rand"
	"net/http"

	"github.com/collinshoop/gomoku/src/board"
	"github.com/gorilla/websocket"
)

var addr = flag.String("addr", "localhost:8080", "http service address")
var upgrader = websocket.Upgrader{} // use default

var brd = board.NewBoard(15, 15)

func Start() {
	flag.Parse()
	log.SetFlags(0)
	http.HandleFunc("/echo", echo)
	http.HandleFunc("/", home)
	log.Fatal(http.ListenAndServe(*addr, nil))
}

func echo(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()
	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		log.Printf("recv: %s", message)

		brd.Move(rand.Intn(2)+1, rand.Intn(15), rand.Intn(15))
		var msg string
		gameIsOver, player := brd.IsOver()
		if gameIsOver {
			msg = fmt.Sprintf("Game over, Player %d won!", player)
		} else {
			msg = brd.ToStr()
		}

		err = c.WriteMessage(mt, []byte(msg))
		if err != nil {
			log.Println("write:", err)
			break
		}
	}
}

func home(w http.ResponseWriter, r *http.Request) {
	_ = homeTemplate.Execute(w, "ws://"+r.Host+"/echo")
}

// TODO needs to go to a static file
var homeTemplate = template.Must(template.New("").Parse(`
<!DOCTYPE html>
<html>
<head>
<meta charset="utf-8">
<script>  
window.addEventListener("load", function(evt) {
    var output = document.getElementById("output");
    var input = document.getElementById("input");
    var ws = new WebSocket("{{.}}");
	ws.onopen = function(evt) {
		print("OPEN");
	}
	ws.onclose = function(evt) {
		print("CLOSE");
		ws = null;
	}
	ws.onmessage = function(evt) {
		print("RESPONSE: " + evt.data);
	}
	ws.onerror = function(evt) {
		print("ERROR: " + evt.data);
	}
    var print = function(message) {
        var d = document.createElement("div");
        d.textContent = message;
        output.textContent = message;
    };
    document.getElementById("send").onclick = function(evt) {
        if (!ws) {
            return false;
        }
        print("SEND: " + input.value);
        ws.send(input.value);
        return false;
    };

    // close websocket
    // ws.close();
});
</script>
</head>
<body>
<table>
<tr><td valign="top" width="50%">
<p>Click "Open" to create a connection to the server, 
"Send" to send a message to the server and "Close" to close the connection. 
You can change the message and send multiple times.
<p>
<form>
<p><input id="input" type="text" value="Hello world!">
<button id="send">Send</button>
</form>
</td><td valign="top" width="50%">
<pre id="output"></pre>
</td></tr></table>
</body>
</html>
`))
