package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/fatih/color"
	"github.com/gtuk/discordwebhook"
	"github.com/inancgumus/screen"
	"github.com/lukesampson/figlet/figletlib"
	"golang.org/x/term"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
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

func readLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

func ErrorCheck(err error) {
	if err != nil {
		msg := strings.Trim(err.Error(), "\n")
		fmt.Println(Pretty(msg))
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

func Delete(webhook string) {
	client := &http.Client{}

	req, err := http.NewRequest("DELETE", webhook, nil)
	ErrorCheck(err)

	resp, err := client.Do(req)
	ErrorCheck(err)
	defer resp.Body.Close()

	_, err = ioutil.ReadAll(resp.Body)
	ErrorCheck(err)
	var ReplacedWebhook string = strings.Replace(webhook, "https://", "", 0)
	var SplitWebhook []string = strings.Split(ReplacedWebhook, "/")

	if strings.Contains(resp.Status, "204") {
		fmt.Println(Pretty("- Succesfully Deleted Webhook: " + SplitWebhook[5]))
	}
}

func Spam(username string, content string, webhook string, amount int, delete_after bool, wg *sync.WaitGroup) {
	defer wg.Done()
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
	if delete_after {
		go Delete(webhook)
		time.Sleep(2 * time.Second)
	}
}

func main() {
	var wg sync.WaitGroup
	Clear()
	Logo()
	content, err := ioutil.ReadFile("main.json")
	ErrorCheck(err)

	var payload Data
	err = json.Unmarshal(content, &payload)
	ErrorCheck(err)

	var username string = payload.Username
	var message string = payload.Messages
	var amount int
	var delete_string string
	var delete_after bool

	fmt.Println()
	fmt.Print(Pretty("Amount: "))
	color.Set(color.FgHiMagenta)
	fmt.Scanln(&amount)

	fmt.Print(Pretty("Delete After [Y/N]: "))
	color.Set(color.FgHiMagenta)
	fmt.Scanln(&delete_string)

	delete_after = strings.Contains(strings.ToLower(delete_string), "y")
	fmt.Println()
	color.Set(color.FgHiWhite)
	Border()

	lines, err := readLines("webhooks.txt")
	for sex, webhook := range lines {
		wg.Add(sex)
		go Spam(username, message, webhook, amount, delete_after, &wg)
	}
	time.Sleep(500000 * time.Minute)
}
