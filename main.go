package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/fatih/color"
	"github.com/gtuk/discordwebhook"
	"github.com/inancgumus/screen"
	"github.com/lukesampson/figlet/figletlib"
	"golang.org/x/term"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"unicode/utf8"
)

type person struct {
	whmessage string
}

var (
	width, _, err = term.GetSize(0)
)

type Data struct {
	Username string
	Messages string
}

func ErrorCheck(err error) {
	if err != nil {
		fmt.Println(Pretty(err.Error()))
	}
}

func NCenter(width int, s string) *bytes.Buffer {
	const half, space = 2, "\u0020"
	var b bytes.Buffer
	n := (width - utf8.RuneCountInString(s)) / half
	fmt.Fprintf(&b, "%s%s", strings.Repeat(space, n), s)
	return &b
}

type erorr struct {
	Global      bool   `json:"global"`
	Message     string `json:"message"`
	Retry_after int    `json:"retry_after"`
}

func Spam(username string, content string, webhook string, amount int) {
	message := discordwebhook.Message{
		Username: &username,
		Content:  &content,
	}
	for i := 1; i <= amount; i++ {

		err := discordwebhook.SendMessage(webhook, message)
		if err != nil {
			if strings.Contains(err.Error(), "rate limit") {
				var errors erorr
				json.Unmarshal([]byte(err.Error()), &errors)
				i = i - 1
			}
		} else {
			fmt.Println(Pretty("- Sent Message: " + content))
		}
	}
}

func Pretty(info string) string {
	pretty := ""
	pretty += color.HiMagentaString("[")
	pretty += color.WhiteString("+")
	pretty += color.HiMagentaString("] ")
	pretty += info
	return pretty
}

func Clear() {
	screen.Clear()
}

func Border() {
	i := 0
	res1 := ""
	for i < width {
		res1 += "─"
		i += 1
	}
	fmt.Println(res1)
}

func Logo() {
	cwd, _ := os.Getwd()
	fontsdir := filepath.Join(cwd, "data")
	ErrorCheck(err)
	f, err := figletlib.GetFontByName(fontsdir, "4max")
	ErrorCheck(err)
	color.Set(color.FgHiMagenta)
	figletlib.PrintMsg("Kyanite", f, width, f.Settings(), "center")
	color.Set(color.FgHiWhite)
	fmt.Println()
	fmt.Println()
	Border()
}

func main() {
	Clear()
	Logo()
	content, err := ioutil.ReadFile("main.json")
	ErrorCheck(err)

	var payload Data
	err = json.Unmarshal(content, &payload)
	ErrorCheck(err)

	var username string = payload.Username
	var message string = payload.Messages
	var webhook string
	var amount int

	fmt.Println()
	fmt.Print(Pretty("Webhook URL: "))
	color.Set(color.FgHiMagenta)
	fmt.Scanln(&webhook)

	fmt.Print(Pretty("Amount: "))
	color.Set(color.FgHiMagenta)
	fmt.Scanln(&amount)
	fmt.Println()
	color.Set(color.FgHiWhite)
	Border()
	Spam(username, message, webhook, amount)
}
