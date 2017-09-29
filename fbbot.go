package main

import (
    "encoding/json"
    "fmt"
    _ "io/ioutil"
    "log"
    "net/http"
)

const (
    FbMsgWebhookToken = "verify_which_cinema_is_ooo"
)

// PAGE ACCESS TOKEN
// EAABqwvSKZCKcBALlTpZBERZC9Ww3z2YUxHAco8jz3ZAQ8lYMBchvGaQZC7hbGMSOMTmi8pxBhuGt1mHmlPqvDVj4dBQYE0Fa2xvLymNbJhJ5dXWqvXtjpV2SN5PSEoi3SMA2IlqOog2DHVekfk0jnmhf5xdcSYMIPhGKAhDYPwcEZCAfM5lAY9

type MessagingEvent struct {
    Entries []Entry   `json:"entry"` 
}

type Entry struct {
    Id         string                  `json:"id"`
    Time       int64                   `json:"time"`
    Messagings []Messaging             `json:"messaging"`
}

type Messaging struct {
    Sender    Sender     `json:"sender"`
    Recipient Recipient  `json:"recipient"`
    Timestamp int64      `json:"timestamp"`
    Message   Message    `json:"message,omitempty"`
    Postback  Postback   `json:"postback,omitempty"`
}

type Sender struct {
    Id string `json:"id"`
}

type Recipient struct {
    Id string `json:"id"`
}

type Message struct {
    Mid            string          `json:"mid"`
    Text           string          `json:"text,omitempty"`
    Attachments    []Attachment    `json:"attachments,omitempty"`
    QuickReply     QuickReply      `json:"quick_reply,omitempty"`
}

type Attachment struct {
    Type    string `json:"type,omitempty"`
    Payload string `json:"payload,omitempty"`
}

type QuickReply struct {
    Payload string `json:"payload,omitempty"`
}

type Postback struct {
    Title       string               `json:"title,omitempty"`
    Payload     string               `json:"payload"`
    Referral    PostbackReferral     `json:"referral,omitempty"`
}

type PostbackReferral struct {
    Ref       string    `json:"ref,omitempty"`
    Source    string    `json:"source"`
    Type      string    `json:"type"`
}

type ResponseMeta struct {
    Status int      `json:"status"`
    Msg    string   `json:"msg"`
}

func main() {
    app := http.NewServeMux()
    app.HandleFunc("/wciooo/bot/webhook", fbwebhookHandler)
    log.Print("Listening on port 8080")
    log.Fatal(http.ListenAndServe(":8080", app))
}

func fbwebhookHandler(w http.ResponseWriter, req *http.Request) {
    // [Start check the method]
    if req.Method == http.MethodGet {
        verifyWebhook(w, req)
    } else if req.Method == http.MethodPost {
        messagingEvents := req.Body

        me := MessagingEvent{}
        jerr := json.NewDecoder(messagingEvents).Decode(&me)
        if jerr != nil {
            fmt.Println("error:", jerr)
        }
        // fmt.Println(me)
        fbWebhookCb(&me)
    } else {
        w.Header().Set("Content-Type", "application/json; charset=utf-8")
        w.WriteHeader(http.StatusMethodNotAllowed)
        // Only support GET method
        resMeta := ResponseMeta{
            Status: http.StatusMethodNotAllowed,
            Msg: "Request method must be 'GET'",
        }
        resMetaJson, err := json.Marshal(resMeta)
        if err != nil {
            fmt.Println("error:", err)
        }
        w.Write(resMetaJson)
    }
    // [End check the method]
        
}

func verifyWebhook(w http.ResponseWriter, req *http.Request) {
    // [Start facebook webhook challenge]
    log.Print("Get /webhook querystring")
    qs := req.URL.Query()
    if qs.Get("hub.verify_token") == FbMsgWebhookToken {
        w.WriteHeader(http.StatusOK)
        w.Write([]byte(qs.Get("hub.challenge")))
    } else {
        w.WriteHeader(http.StatusOK)
        w.Write([]byte("Error, wrong token"))
    }
}

func fbWebhookCb(me *MessagingEvent) {
    // fmt.Println(me)
    for k, v := range me.Entries  {
        fmt.Printf("key: %v\n", k)
        fmt.Printf("value: %v\n", v.Messagings)
        for _, msg := range v.Messagings {
            fmt.Printf("%u\n", msg)
            fmt.Printf("%u\n", msg.Sender.Id)
        }
    }
}
