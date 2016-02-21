
package main

import (
  "fmt"
  "net/http"
  // "io"
  "io/ioutil"
  "encoding/json"
  "time"
  "log"
  "golang.org/x/net/websocket"
)

var comments_tracker = make(map[string]string)
// type WebSockMW chan

var socket_ch = make(chan Comment, 100)

type Comment struct {
  Data struct {
    ApprovedBy          interface{}   `json:"approved_by"`
    Archived            bool          `json:"archived"`
    Author              string        `json:"author"`
    AuthorFlairCSSClass interface{}   `json:"author_flair_css_class"`
    AuthorFlairText     interface{}   `json:"author_flair_text"`
    BannedBy            interface{}   `json:"banned_by"`
    Body                string        `json:"body"`
    BodyHTML            string        `json:"body_html"`
    Controversiality    int           `json:"controversiality"`
    Created             int           `json:"created"`
    CreatedUtc          int           `json:"created_utc"`
    Distinguished       interface{}   `json:"distinguished"`
    Downs               int           `json:"downs"`
    Edited              bool          `json:"edited"`
    Gilded              int           `json:"gilded"`
    ID                  string        `json:"id"`
    Likes               interface{}   `json:"likes"`
    LinkAuthor          string        `json:"link_author"`
    LinkID              string        `json:"link_id"`
    LinkTitle           string        `json:"link_title"`
    LinkURL             string        `json:"link_url"`
    ModReports          []interface{} `json:"mod_reports"`
    Name                string        `json:"name"`
    NumReports          interface{}   `json:"num_reports"`
    Over18              bool          `json:"over_18"`
    ParentID            string        `json:"parent_id"`
    Quarantine          bool          `json:"quarantine"`
    RemovalReason       interface{}   `json:"removal_reason"`
    Replies             string        `json:"replies"`
    ReportReasons       interface{}   `json:"report_reasons"`
    Saved               bool          `json:"saved"`
    Score               int           `json:"score"`
    ScoreHidden         bool          `json:"score_hidden"`
    Stickied            bool          `json:"stickied"`
    Subreddit           string        `json:"subreddit"`
    SubredditID         string        `json:"subreddit_id"`
    Ups                 int           `json:"ups"`
    UserReports         []interface{} `json:"user_reports"`
  } `json:"data"`
  Kind string `json:"kind"`
}


type Comments struct {
    Data struct {
        After    interface{} `json:"after"`
        Before   interface{} `json:"before"`
        Children []Comment `json:"children"`
        Modhash string `json:"modhash"`
    } `json:"data"`
    Kind string `json:"kind"`
}

func getComments(user string, comments_ch chan Comment) {
  ticker := time.NewTicker(5 * time.Second)
  for {
    select {
      case <- ticker.C:
        url := fmt.Sprintf("https://api.reddit.com/user/%s/comments?limit=1&before=%s", user, comments_tracker[user])
        resp, err := http.Get(url)
        if err == nil {
          body, err := ioutil.ReadAll(resp.Body)
          if err == nil {
            resp.Body.Close()
            var comments Comments
            json.Unmarshal(body, &comments)
            for _, comment := range comments.Data.Children {
              comments_ch <- comment
              comments_tracker[user] = comment.Data.Name
            }
          }
        }
    }
  }
}

func processComments(socket_ch, comments_ch chan Comment) {
  for {
    select {
      case comment := <- comments_ch:
        socket_ch <- comment
    }
  }
}

func WebsocketHandler(ws *websocket.Conn) {
  for {
    select {
    case comment := <- socket_ch:
      websocket.JSON.Send(ws, comment)
    }
  }
}

func main() {
  comments_ch := make(chan Comment, 100)
  users := make([]string, 0)
  users_list = "Poem_for_your_sprog Shitty_Watercolour "
  users = append(users, "Poem_for_your_sprog")
  users = append(users, "Shitty_Watercolour")
  users = append(users, "golangbottest")

  for _, user := range users {
    go getComments(user, comments_ch)
    go processComments(socket_ch, comments_ch)
  }

  http.Handle("/", websocket.Handler(WebsocketHandler))
  if err := http.ListenAndServe(":1234", nil); err != nil {
    log.Fatal("ListenAndServe:", err)
  }
}
