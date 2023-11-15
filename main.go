package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	gt "github.com/bas24/googletranslatefree"
)

func main() {
	var check int8 = 0
	botToken := "TOKEN"
	botApi := "https://api.telegram.org/bot"
	botUrl := botApi + botToken
	offset := 0
	for {
		updates, err := getUpdates(botUrl, offset)
		fmt.Println(updates)
		if err != nil {
			log.Println(err)
		}
		for _, update := range updates {
			err := respond(update, botUrl, &check)
			fmt.Println(err)
			if err != nil {
				log.Println(err)
			}
			offset = update.UpdateId + 1
		}
	}
}

func getUpdates(boturl string, offset int) ([]Update, error) {
	resp, err := http.Get(boturl + "/getUpdates" + "?offset=" + strconv.Itoa(offset))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var restResponse RestResponse
	err = json.Unmarshal(body, &restResponse)
	if err != nil {
		return nil, err
	}
	return restResponse.Result, nil
}

func respond(update Update, botUrl string, check *int8) error {
	switch update.Message.Text {
	case "/start":
		err := CommandStart(update, botUrl)
		return err
	case "Английский":
		err := CommandEnglish(update, botUrl, check)
		return err
	case "Русский":
		err := CommandRussian(update, botUrl, check)
		return err
	default:
		err := CommandError(update, botUrl, check)
		return err
	}

}

func CommandStart(update Update, botUrl string) error {
	var sendMessage SendMessage
	var replykeyboardmarkup ReplyKeyboardMarkup
	var keyboardbutton []Keyboardbutton
	var button1 Keyboardbutton
	var button2 Keyboardbutton
	button1.Text = "Английский"
	button2.Text = "Русский"
	keyboardbutton = append(keyboardbutton, button1)
	keyboardbutton = append(keyboardbutton, button2)
	replykeyboardmarkup.ResizeKeyboard = true
	replykeyboardmarkup.OneTimeKeyboard = true
	replykeyboardmarkup.Keyboard = append(replykeyboardmarkup.Keyboard, keyboardbutton)
	sendMessage.ReplyMarkup = replykeyboardmarkup
	sendMessage.ChatId = update.Message.Chat.Id
	sendMessage.Text = "Привет " + update.Message.From.FirstName + ", вам необходимо выбрать с какого языка вы будете осуществлять перевод!"
	buf, err := json.Marshal(sendMessage)
	if err != nil {
		return err
	}
	_, err = http.Post(botUrl+"/sendMessage", "application/json", bytes.NewBuffer(buf))
	if err != nil {
		return err
	}
	return nil
}

func CommandEnglish(update Update, botUrl string, check *int8) error {
	var sendMessage SendMessage2
	sendMessage.ChatId = update.Message.Chat.Id
	sendMessage.Text = "Отлично, следующие английские фразы будут переведенны на русский язык:"
	buf, err := json.Marshal(sendMessage)
	if err != nil {
		return err
	}
	_, err = http.Post(botUrl+"/sendMessage", "application/json", bytes.NewBuffer(buf))
	if err != nil {
		return err
	}
	*check = 1
	return nil
}
func CommandRussian(update Update, botUrl string, check *int8) error {
	var sendMessage SendMessage2
	sendMessage.ChatId = update.Message.Chat.Id
	sendMessage.Text = "Отлично, следующие русские фразы будут переведенны на английский язык:"
	buf, err := json.Marshal(sendMessage)
	if err != nil {
		return err
	}
	_, err = http.Post(botUrl+"/sendMessage", "application/json", bytes.NewBuffer(buf))
	if err != nil {
		return err
	}
	*check = 2
	return nil
}
func CommandError(update Update, botUrl string, check *int8) error {
	switch *check {
	case 1:
		// перевод с английского на русский
		var sendMessage SendMessage2
		sendMessage.ChatId = update.Message.Chat.Id
		result := englishTranslate(update.Message.Text)
		sendMessage.Text = result
		buf, err := json.Marshal(sendMessage)
		if err != nil {
			return err
		}
		_, err = http.Post(botUrl+"/sendMessage", "application/json", bytes.NewBuffer(buf))
		if err != nil {
			return err
		}

	case 2:
		// перевод с русского на английский
		var sendMessage SendMessage2
		sendMessage.ChatId = update.Message.Chat.Id
		result := russianTranslate(update.Message.Text)
		sendMessage.Text = result
		buf, err := json.Marshal(sendMessage)
		if err != nil {
			return err
		}
		_, err = http.Post(botUrl+"/sendMessage", "application/json", bytes.NewBuffer(buf))
		if err != nil {
			return err
		}
	default:
		var sendMessage SendMessage2
		sendMessage.ChatId = update.Message.Chat.Id
		sendMessage.Text = "Пожалуйста выберите валидный язык!"
		buf, err := json.Marshal(sendMessage)
		if err != nil {
			return err
		}
		_, err = http.Post(botUrl+"/sendMessage", "application/json", bytes.NewBuffer(buf))
		if err != nil {
			return err
		}
	}
	return nil
}

func englishTranslate(text string) string {
	result, _ := gt.Translate(text, "en", "ru")
	return result
}
func russianTranslate(text string) string {
	result, _ := gt.Translate(text, "ru", "en")
	return result
}
