package main

import (
    tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
    "github.com/jpillora/go-ogle-analytics"
    "errors"
    "fmt"
    "net/http"
    "io/ioutil"
    "encoding/json"
    "os"
    "strings"
    "time"
)

var _words []string
var gaClient *ga.Client

func processQuery(update tgbotapi.Update) (tgbotapi.InlineConfig) {
    var results []interface{}
    query := update.InlineQuery.Query
    result := "Не " + query + ", а говно"
    if query != "" {
        article := tgbotapi.NewInlineQueryResultArticle(update.InlineQuery.ID, result, result)
        results = append(results, article)
    }
    inline := tgbotapi.InlineConfig{
        InlineQueryID: update.InlineQuery.ID,
        IsPersonal: true,
        CacheTime: 0,
        Results: results,
    }
    gaClient.Send(ga.NewEvent("Inline", "Govno").Label(query))
    return inline
}

func getWords(words *[]string) {
    for {
        response, err := http.Get("https://raw.githubusercontent.com/alexeimoisseev/NeGovnoABot/master/words.json")
        if err != nil {
            fmt.Println("Error getting words")
            fmt.Println(err)
            continue
        }
        defer response.Body.Close()
        contents, err := ioutil.ReadAll(response.Body)
        if err != nil {
            fmt.Println("Error reading body stream")
            fmt.Println(err)
        }
        err = json.Unmarshal(contents, words)
        if err != nil {
            fmt.Println("Error parsing json")
            fmt.Println(err)
        }
        time.Sleep(60 * time.Second)
    }
}

func createReply(update tgbotapi.Update) (error, *tgbotapi.MessageConfig) {
    message := update.Message.Text
    lower := strings.ToLower(message)
    if strings.Contains(lower, "кадыров") {
        reply := "Извинись!"
        msg := tgbotapi.NewMessage(update.Message.Chat.ID, reply)
        msg.ReplyToMessageID = update.Message.MessageID
        gaClient.Send(ga.NewEvent("Message", "Sorry"))
        return nil, &msg
    }

    for _, word := range _words {
        if strings.Contains(lower, word) {
            reply := "Не " + word + ", а говно"
            msg := tgbotapi.NewMessage(update.Message.Chat.ID, reply)
            msg.ReplyToMessageID = update.Message.MessageID
            gaClient.Send(ga.NewEvent("Message", "Govno").Label(word))
            return nil, &msg
        }
    }
    return errors.New("No match"), nil
}

func main() {
    go getWords(&_words)
    key := os.Getenv("KEY")
    gaId := os.Getenv("GA")
    _client, err := ga.NewClient(gaId)
    if err != nil {
        panic(err)
    }
    gaClient = _client

    bot, err := tgbotapi.NewBotAPI(key)
    if err != nil {
        panic(err)
    }
    u := tgbotapi.NewUpdate(0)
    updates, err := bot.GetUpdatesChan(u)
    for update := range updates {
        fmt.Println(update)
        if update.InlineQuery != nil {
            inline := processQuery(update)
            bot.AnswerInlineQuery(inline)
        }

        if update.Message != nil {
            err, reply := createReply(update)
            if err == nil {
                bot.Send(reply)
            }
        }
    }

}
