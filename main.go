package main

import (
    tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
    "errors"
    "fmt"
    "strings"
    "os"
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

var words []string = []string {
    "пхп",
    "php",
    "вуе",
    "vue",
    "яндекс",
    " го ",
    " go ",
    "golang",
    "голанг",
    "питон",
    "python",
    "стартап",
    "карпрайс",
    "ангуляр",
    "angular",
    "реакт",
    "react",
    "джаваскрипт",
    "дзюба",
}


func processMessage(update tgbotapi.Update) (error, *tgbotapi.MessageConfig) {
    message := update.Message.Text
    lower := strings.ToLower(message)
    for _, word := range words {
        if strings.Contains(lower, word) {
            reply := "Не " + word + ", а говно"
            msg := tgbotapi.NewMessage(update.Message.Chat.ID, reply)
            msg.ReplyToMessageID = update.Message.MessageID
            return nil, &msg
        }
    }
    return errors.New("No match"), nil
}

func main() {
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
            err, reply := processMessage(update)
            if err == nil {
                bot.Send(reply)
            }
        }
    }

}
