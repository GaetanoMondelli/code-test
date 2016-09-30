package main

/* author : Gaetano Mondelli
 * id2Map stores the pairs sessionId and Data structure.
* Each JSON sent from client is unmarshalled into an Event structure.
* Depending on the fields 'Event Type' and 'SessionId' of the Event structure,
* the postHandle function updates some fields  of the specific 'Data'
* structure in the map.
*/

import (
	"bytes"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"regexp"
	"strconv"
	"sync"
)

type Event struct {
	EventType  string
	WebsiteUrl string
	SessionId  string
	Time       int
	Pasted     string
	FormId     string
	Width      string
	Height     string
	ResizeFrom Dimension
	ResizeTo   Dimension
}

type Data struct {
	WebsiteUrl         string
	SessionId          string
	ResizeFrom         Dimension
	ResizeTo           Dimension
	CopyAndPaste       map[string]bool // map[fieldId]true
	FormCompletionTime int             // Seconds
}

func (data Data) String() string {
	title := "\nData Structure\n"
	b, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		fmt.Println("error:", err)
	}
	reg := regexp.MustCompile(`^{|}$`)
	prettify := reg.ReplaceAllString(string(b), "")
	if data.FormCompletionTime > 0 {
		title = "\nData Structure COMPLETED\n"
	}
	return title + prettify + "\n"
}

type Dimension struct {
	Width  string
	Height string
}

var id2DataMap = struct {
	sync.RWMutex
	m map[string]Data
}{m: make(map[string]Data)}

func NewData(rows, cols int) *Data {
	d := new(Data)
	d.CopyAndPaste = make(map[string]bool)
	return d
}

func main() {

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		cookie, err := r.Cookie("sessionId")

		if err != nil { //no cookie found
			id := generateId() //generates a SessionId
			cookie := &http.Cookie{
				Name:  "sessionId",
				Value: id,
			}
			http.SetCookie(w, cookie)
		}

		fmt.Println("New sessionId generated", cookie)

		switch r.Method {

		case "GET":
			getHandle(w)
		case "POST":
			postHandle(w, r)
		case "PUT":
			// Update
		case "DELETE":
			// Remove
		default:
			fmt.Println("Error in the method")
		}

	})

	http.HandleFunc("/client/frontend.js", func(w http.ResponseWriter, r *http.Request) {
		frontendJsRoute(w)
	})

	fmt.Println(http.ListenAndServe(":8080", nil))
}

func frontendJsRoute(w http.ResponseWriter) {
	t := template.Must(template.ParseFiles("client/frontend.js"))
	if err := t.ExecuteTemplate(w, "frontend.js", nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func getHandle(w http.ResponseWriter) {
	//fmt.Println("Get")
	t := template.Must(template.ParseFiles("client/index.html"))

	if err := t.ExecuteTemplate(w, "index.html", nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func postHandle(w http.ResponseWriter, r *http.Request) {

	e := json2Event(r)
	id2DataMap.Lock()
	defer id2DataMap.Unlock() // assures lock release (try-catch like)

	if data, ok := id2DataMap.m[e.SessionId]; !ok {
		data.WebsiteUrl = e.WebsiteUrl
		data.SessionId = e.SessionId
		data.CopyAndPaste = make(map[string]bool)
		id2DataMap.m[e.SessionId] = data
	}

	data := id2DataMap.m[e.SessionId]

	if data.FormCompletionTime > 0 {
		fmt.Println("Form already submitted")
		return
	}

	switch e.EventType {
	case "copyAndPaste":
		id2DataMap.m[e.SessionId].CopyAndPaste[e.FormId], _ =
			strconv.ParseBool(e.Pasted)
	case "timeTaken":
		data.FormCompletionTime = e.Time //strconv.Atoi(e.Time)
		id2DataMap.m[e.SessionId] = data

	case "screenResize":
		data.ResizeFrom = e.ResizeFrom
		data.ResizeTo = e.ResizeTo
		id2DataMap.m[e.SessionId] = data

	default:

	}
	fmt.Println(id2DataMap.m[e.SessionId])
	//id2DataMap.Unlock()

}

func json2Event(r *http.Request) Event {
	var e Event
	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)
	s := buf.String()
	if err_m := json.Unmarshal([]byte(s), &e); err_m != nil {
		fmt.Println("Error during the unmarshalling")
	}
	return e
}

func generateId() string {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		fmt.Println("Error: ", err)
		return ""
	}
	uuid := fmt.Sprintf("%X-%X-%X-%X-%X", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
	return uuid
}
