package main

import (
	"context"
	"crypto/rand"
	"fmt"
	"github.com/adshao/go-binance/v2"
	"github.com/bwmarrin/discordgo"
	tb "gopkg.in/tucnak/telebot.v2"
	"io/ioutil"
	"log"
	"math/big"
	"os"
	"strconv"
	"strings"
	"time"
)

var (
	// Universal markup builders.
	menu = &tb.ReplyMarkup{ResizeReplyKeyboard: true}

	// Reply buttons.
	btnadd    = menu.Text("Θ 추가")
	btnremove = menu.Text("Θ 제거")
	btnlist   = menu.Text("Θ 리스트")
	btnshutdown = menu.Text("Θ 종료")
	btncancle = menu.Text("Θ 취소")
	btn2fa = menu.Text("Θ 2FA 인증")
)
var (
	apiKey = "BINANCE_API_KEY"
	secretKey = "BINANCE_SECRET_KEY"
)
var (
	// Universal markup builders.
	way = &tb.ReplyMarkup{ResizeReplyKeyboard: true}

	// Reply buttons.
	btnbuy    = way.Text("매수")
	btnstoplos = way.Text("손절")
	btnsell   = way.Text("익절")
	btnreserv   = way.Text("익절예약")
)
var (
	// Universal markup builders.
	removeop = &tb.ReplyMarkup{ResizeReplyKeyboard: true}

	// Reply buttons.
	btn1    = removeop.Text("제거 1")
	btn2 = removeop.Text("제거 2")
	btn3   = removeop.Text("제거 3")
	btn4   = removeop.Text("제거 4")
	btn5    = removeop.Text("제거 5")
	btn6 = removeop.Text("제거 6")
	btn7   = removeop.Text("제거 7")
	btn8   = removeop.Text("제거 8")
	btn9   = removeop.Text("제거 9")
	btn10   = removeop.Text("제거 10")
	btnall   = removeop.Text("제거 모두")
)

var user = &tb.User{ID: TELEGRAM_USER_ID}
var poller = &tb.LongPoller{Timeout: 15 * time.Second}
var bp *tb.Bot
var msg *tb.Update
var spamProtected = tb.NewMiddlewarePoller(poller, func(upd *tb.Update) bool {

	if upd.Message == nil {
		return true
	}
	msg = upd
	data, _ := ioutil.ReadFile("./auth.txt")

	if string(data) == "0" {
		if strings.Contains(upd.Message.Text,"2FA 인증") == true{
			return true
		}
		twfa, _ := ioutil.ReadFile("./authy.txt")

		if msg.Message.Text != string(twfa){
			_, _ = bp.Send(user, "번호가 올바르지 않습니다.")
			return false
		} else{
			menu.Reply(
				menu.Row(btnadd),
				menu.Row(btnremove),
				menu.Row(btnlist),
				menu.Row(btnshutdown),
			)
			go remove()
			_ = ioutil.WriteFile("./auth.txt", []byte("1"), os.FileMode(644))
			_ = ioutil.WriteFile("./authy.txt", []byte("0"), os.FileMode(644))
			_, _ = bp.Send(user, "5분간 권한이 부여됩니다", menu)
		}
	}

	if strings.Contains(string(data), "1") == true {

		resp := "1"
		if strings.Contains(string(data), resp) == true {
			if strings.Contains(msg.Message.Text, "추가") == true {
				_ = ioutil.WriteFile("./cache.txt", []byte(""), os.FileMode(644))
				way.Reply(
					way.Row(btnbuy),
					way.Row(btnstoplos),
					way.Row(btnsell),
					way.Row(btnreserv),
					menu.Row(btncancle),
				)
				_, _ = bp.Send(user, "매수 방법을 선택해주세요", way)
			}
			if msg.Message.Text == "매수" || msg.Message.Text == "손절" || msg.Message.Text == "익절" || msg.Message.Text == "익절예약" {
				way.Reply(
					way.Row(btnbuy),
					way.Row(btnstoplos),
					way.Row(btnsell),
					way.Row(btnreserv),
					menu.Row(btncancle),
				)
				_, _ = bp.Send(user, "코인의 티커를 입력해주세요", way)
				_ = ioutil.WriteFile("./cache.txt", []byte(msg.Message.Text+","), os.FileMode(644))
				return true
			} else if strings.Contains(msg.Message.Text, " 제거") == true {
				removeop.Reply(
					removeop.Row(btn1),
					removeop.Row(btn2),
					removeop.Row(btn3),
					removeop.Row(btn4),
					removeop.Row(btn5),
					removeop.Row(btn6),
					removeop.Row(btn7),
					removeop.Row(btn8),
					removeop.Row(btn9),
					removeop.Row(btn10),
					removeop.Row(btnall),
					menu.Row(btncancle),
				)
				_, _ = bp.Send(user, "리스트를 선택해주세요", removeop)




			} else if strings.Contains(msg.Message.Text, "제거 ") == true {

				// 모두
				if strings.Contains(msg.Message.Text, "모두") {
					err := ioutil.WriteFile("./data.txt", []byte(""), os.FileMode(644))
					if err != nil {
						_, _ = bp.Send(user, "예약된 주문들을 모두 삭제하는데 오류가 발생했습니다", menu)
						return true
					}
					_, _ = bp.Send(user, "주문을 모두 삭제하였습니다.", menu)
					return true
				}

				// 번호
				message_ := strings.TrimLeft(msg.Message.Text, "제거 ")
				data, _ := ioutil.ReadFile("./data.txt")
				data_ := strings.Split(string(data), "\n")
				num, _ := strconv.Atoi(message_)

				// num보다 예약된 주문 적으면 초기화
				if len(data_) < num {
					_, _ = bp.Send(user, "해당 번호의 주문이 없습니다", menu)
					return true
				}
				// 데이타에서 주문 삭제
				fmt.Println(data)
				data_ = append(data_[:num-1], data_[num:]...)
				fmt.Println(len(data_))
				fmt.Println(data_)
				text := strings.Join(data_[:], "\n")
				_ = ioutil.WriteFile("./data.txt", []byte(text), os.FileMode(644))
				//제거 끝 -- 문제 없음

				//리스트 시작
				data, _ = ioutil.ReadFile("./data.txt")
				if string(data) == ""{
					_, _ = bp.Send(user, "예약된 주문이 없습니다", menu)
					return true
				}
				data_ = strings.Split(string(data), "\n")
				for i, v := range data_ {
					// 매수 타임스탬프 계산 X
					if strings.Contains(v, "매수") == true {
						data__ := strings.Split(v, ",")
						data__ = data__[:4]
						v = strings.Join(data__[:], ",")
					}
					data_[i] = strconv.Itoa(i+1) + "." + v
				}
				text = strings.Join(data_[:], "\n")
				_, _ = bp.Send(user, "주문 리스트\n\n\n"+text)
				removeop.Reply(
					menu.Row(btncancle),
					removeop.Row(btn1),
					removeop.Row(btn2),
					removeop.Row(btn3),
					removeop.Row(btn4),
					removeop.Row(btn5),
					removeop.Row(btn6),
					removeop.Row(btn7),
					removeop.Row(btn8),
					removeop.Row(btn9),
					removeop.Row(btn10),
					removeop.Row(btnall),
				)
				_, _ = bp.Send(user, "삭제완료!", removeop)
				return true





			} else if strings.Contains(msg.Message.Text, "리스트") {

				data, _ = ioutil.ReadFile("./data.txt")
				if string(data) == ""{
					_, _ = bp.Send(user, "예약된 주문이 없습니다", menu)
					return true
				}
				data_ := strings.Split(string(data), "\n")
				for i, v := range data_ {
					// 매수 타임스탬프 계산 X
					if strings.Contains(v, "매수") == true {
						data__ := strings.Split(v, ",")
						data__ = data__[:4]
						v = strings.Join(data__[:], ",")
					}
					data_[i] = strconv.Itoa(i+1) + "." + v
				}
				text := strings.Join(data_[:], "\n")
				_, _ = bp.Send(user, "주문 리스트\n\n\n"+text, menu)



			} else if strings.Contains(msg.Message.Text, "종료") {
				removenow()

			} else if strings.Contains(msg.Message.Text, "취소") {
				_, _ = bp.Send(user, "작업을 취소합니다", menu)
			} else {

				cache, _ := ioutil.ReadFile("./cache.txt")
				if string(cache) == "" {
					return true
				}
				ord := strings.Split(string(cache), ",")
				if strings.Contains(string(cache), "익절예약") == true {
					if len(ord) == 2 {
						client := binance.NewClient(apiKey, secretKey)
						_, err := client.NewAveragePriceService().Symbol(strings.ToUpper(msg.Message.Text) + "USDT").Do(context.Background())
						if err != nil {
							fmt.Println(err)
							_, _ = bp.Send(user, "없는 코인입니다.")
							return true
						}
						text := string(cache) + msg.Message.Text + ","
						_ = ioutil.WriteFile("./cache.txt", []byte(text), os.FileMode(644))
						_, _ = bp.Send(user, "목표가격을 입력해주세요")
						return true

					} else if len(ord) == 3 {
						message_ := strings.TrimRight(string(cache), ",")
						text := message_ + "," + msg.Message.Text + ","
						_ = ioutil.WriteFile("./cache.txt", []byte(text), os.FileMode(644))
						_, _ = bp.Send(user, "익절가를 입력해주세요")
						return true
					} else if len(ord) == 4 {
						message_ := strings.TrimRight(string(cache), ",")
						text := message_ + "," + msg.Message.Text + ","
						_ = ioutil.WriteFile("./cache.txt", []byte(text), os.FileMode(644))
						_, _ = bp.Send(user, "개수를 입력해주세요")

						return true
					} else if len(ord) == 5 {
						data, _ := ioutil.ReadFile("./data.txt")
						message_ := strings.TrimRight(string(cache), ",")
						text := message_ + "," + msg.Message.Text
						if string(data) != "" {
							data_ := strings.Split(string(data), "\n")
							data_ = append(data_, text)
							text = strings.Join(data_[:], "\n")
						}
						_ = ioutil.WriteFile("./data.txt", []byte(text), os.FileMode(644))
						_ = ioutil.WriteFile("./cache.txt", []byte(""), os.FileMode(644))
						_, _ = bp.Send(user, "추가 완료", menu)
						return true
					}
					return true
				}
				if len(ord) == 2 {
					client := binance.NewClient(apiKey, secretKey)
					_, err := client.NewAveragePriceService().Symbol(strings.ToUpper(msg.Message.Text) + "USDT").Do(context.Background())
					if err != nil {
						fmt.Println(err)
						_, _ = bp.Send(user, "없는 코인입니다.")
						return true
					}
					text := string(cache) + msg.Message.Text + ","
					_ = ioutil.WriteFile("./cache.txt", []byte(text), os.FileMode(644))
					_, _ = bp.Send(user, "가격을 입력해주세요")
					return true

				}  else if len(ord) == 3 {
					message_ := strings.TrimRight(string(cache), ",")
					text := message_ + "," + msg.Message.Text + ","
					_ = ioutil.WriteFile("./cache.txt", []byte(text), os.FileMode(644))
					_, _ = bp.Send(user, "개수를 입력해주세요")
					return true
				} else if len(ord) == 4 {
					data, _ := ioutil.ReadFile("./data.txt")
					message_ := strings.TrimRight(string(cache), ",")
					text := message_ + "," + msg.Message.Text
					if string(data) != ""{
					data_ := strings.Split(string(data), "\n")
					data_ = append(data_, text)
					text = strings.Join(data_[:], "\n")
					}
					_ = ioutil.WriteFile("./data.txt", []byte(text), os.FileMode(644))
					_ = ioutil.WriteFile("./cache.txt", []byte(""), os.FileMode(644))
					_, _ = bp.Send(user, "추가 완료", menu)
					return true
				}
			}
		}
	}

		 return true
	})

func main() {
	_, err := os.Open("/auth.txt")
	if err != nil {
		_, _ = os.Create("./auth.txt")
	}
	_, err = os.Open("/authy.txt")
	if err != nil {
		_, _ = os.Create("./authy.txt")
	}
	_, err = os.Open("/cache.txt")
	if err != nil {
		_, _ = os.Create("./cache.txt")
	}
	_, err = os.Open("/data.txt")
	if err != nil {
		_, _ = os.Create("./data.txt")
	}
	bot, err := tb.NewBot(tb.Settings{
		Token:  "TELEGRAM_BOT_TOKEN",
		Poller: spamProtected,
	})
	bp = bot
	if err != nil {
		log.Fatal(err)
		return
	}
	_ = ioutil.WriteFile("./auth.txt", []byte("0"), os.FileMode(644))
	menu.Reply(
		menu.Row(btn2fa),
		)
	_, _ = bot.Send(user, "2FA 인증번호 발행", menu)
	discord, err := discordgo.New("Bot " + "DISCORD_BOT_TOKEN")
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}
	err = discord.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}
	bot.Handle(&btn2fa, func(m *tb.Message) {
		a := genRandNum(111111, 999999)
		_ = ioutil.WriteFile("./authy.txt", []byte(strconv.Itoa(int(a))), os.FileMode(644))
		_, err = discord.ChannelMessageSend("CHANNEL_NUM", strconv.Itoa(int(a)))
		if err != nil {
			fmt.Println("Error while sending chat", err)
			return
		}
		_, _ = ioutil.ReadFile("./authy.txt")
		go remove2fa()
	})
	bot.Start()
}
func remove() {
	time.Sleep(time.Second * 300)
	_ = ioutil.WriteFile("./auth.txt", []byte("0"), os.FileMode(644))
	_ = ioutil.WriteFile("./cache.txt", []byte(""), os.FileMode(644))
	menu.Reply(
		menu.Row(btn2fa),
	)
	_, _ = bp.Send(user, "인증이 만료되었습니다.", menu)
}
func removenow() {
	_ = ioutil.WriteFile("./auth.txt", []byte("0"), os.FileMode(644))
	_ = ioutil.WriteFile("./cache.txt", []byte(""), os.FileMode(644))
	menu.Reply(
		menu.Row(btn2fa),
	)
	_, _ = bp.Send(user, "종료되었습니다.", menu)
}
func genRandNum(min, max int64) int64 {
	bg := big.NewInt(max - min)

	n, err := rand.Int(rand.Reader, bg)
	if err != nil {
		panic(err)
	}

	return n.Int64() + min
}
func remove2fa() {
	time.Sleep(time.Second * 33)
	_ = ioutil.WriteFile("./authy.txt", []byte(""), os.FileMode(644))
}
