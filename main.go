package main

import (
    tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
    "errors"
    "fmt"
    "net/http"
    "io/ioutil"
    "encoding/json"
    "os"
    "strings"
    "time"
)

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
    return inline
}



var _words []string


func processMessage(update tgbotapi.Update, words *[]string) (error, *tgbotapi.MessageConfig) {
    message := update.Message.Text
    lower := strings.ToLower(message)
    for _, word := range *words {
        if strings.Contains(lower, word) {
            reply := "Не " + word + ", а говно"
            msg := tgbotapi.NewMessage(update.Message.Chat.ID, reply)
            msg.ReplyToMessageID = update.Message.MessageID
            return nil, &msg
        }
    }
    return errors.New("No match"), nil
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

func main() {
    go getWords(&_words)
    key := os.Getenv("KEY")
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
            err, reply := processMessage(update, &_words)
            if err == nil {
                bot.Send(reply)
            }
        }
    }

}
