import datetime

from binance.client import Client
from binance.exceptions import *
import time, random
from telegram.ext import Updater

while True:
    try:
        client = Client(api_key='BINANCE_API_KEY',
                        api_secret='BINANCE_SECRET_KEY')
        break
    except Exception:
        print(str(Exception))
        time.sleep(10)
        pass

user_id = "TELEGRAM_USER_ID"
updater = Updater(token='TELEGRAM_TOKEN', use_context=True)


def createex():
    try:
        f = open("/data.txt", 'r')
        f.close()
    except Exception:
        f = open("/data.txt", 'w')
        f.close()


def removefromorder(row):
    with open("/data.txt", 'r') as f:
        lines = f.readlines()
        lines = list(map(lambda s: s.strip(), lines))
    if type(row) == list:
        for a in row:
            del lines[int(a)]
            num = 0
            for c in row:
                row[num] = c - 1
                num += 1
    lines = "\n".join(map(str, lines))
    with open("/data.txt", 'w') as f:
        f.write(lines)
    clean()


def discountamount(num, amount):
    def read():
        with open("/data.txt", 'r') as f:
            lines = f.readlines()
        lines = list(map(lambda s: s.strip(), lines))
        data = lines[num]
        data = data.split(',')
        return data, lines

    def write(data, lines):
        data = ",".join(map(str, data))
        lines[num] = data
        lines = "\n".join(map(str, lines))
        with open("/data.txt", 'w') as f:
            f.write(lines)
        clean()

    data, lines = read()
    try:
        print(data[4])
    except:
        data.append(int(datetime.datetime.now().timestamp()) + random.randrange(20,45) * 60)
        write(data, lines)
        data, lines = read()

    if datetime.datetime.fromtimestamp(float(data[4])).minute != int(datetime.datetime.now().minute):
        return '-1'
    del data[4]
    data.append(int(datetime.datetime.now().timestamp()) + random.randrange(20,45) * 60)
    cache = data[3]
    data[3] = str(float(data[3]) - float(amount))
    lst = []
    lst.append(num)
    if float(data[3]) <= 0:
        removefromorder(lst)
        return float(cache)
    write(data, lines)
    return amount


def removetargetprice(row):
    with open("/data.txt", 'r') as f:
        lines = f.readlines()
        lines = list(map(lambda s: s.strip(), lines))
    b = lines[row].split(',')
    b[2] = 'o'
    lines[row] = ",".join(map(str, b))
    lines = "\n".join(map(str, lines))
    with open("/data.txt", 'w') as f:
        f.write(lines)
    clean()


def clean():
    f = open("/data.txt", 'r')
    line = f.read()
    f.close()
    while line.find("\n\n") != -1:
        line = line.replace("\n\n", "\n")
    if line != '':
        if line[len(line) - 1:len(line)] == '\n':
            line = line[:len(line) - 2]
    f = open("/data.txt", 'w')
    f.write(line)


def buy_limit_test(coin, amount):
    client.create_test_order(
        symbol=coin + 'USDT',
        side=Client.SIDE_BUY,
        type=Client.ORDER_TYPE_MARKET,
        quantity=float(amount))


def buy_limit(coin, amount):
    client.create_order(
        symbol=coin + 'USDT',
        side=Client.SIDE_BUY,
        type=Client.ORDER_TYPE_MARKET,
        quantity=float(amount))


def sell_market_test(coin, amount):
    client.create_test_order(
        symbol=coin + 'USDT',
        side=Client.SIDE_SELL,
        type=Client.ORDER_TYPE_MARKET,
        quantity=float(amount))


def sell_market(coin, amount):
    client.create_order(
        symbol=coin + 'USDT',
        side=Client.SIDE_SELL,
        type=Client.ORDER_TYPE_MARKET,
        quantity=float(amount))


def loop():
    def orderfunc(son, coin, amount, pri, num2, number, marketprice):
        maa = ""
        b = 0
        if son == 1:
            maa = '손절'
        if son == 2:
            maa = '익절'
        if son == 3:
            maa = '익절예약'
        if son == 4:
            maa = '매수'
            b = round(random.randrange(20, 45) / float(marketprice), 4)
        na = 3
        chance = 4
        while True:
            if chance == 0:
                break
            try:
                if son == 4:
                    buy_limit_test(coin, b)
                    a = discountamount(num2, b)
                    if a == '-1':
                        break
                    # 실주문 들어가야함
                    buy_limit_test(coin, a)
                    updater.bot.sendMessage(user_id, coin + " 코인의" + pri + " 가격에" + str(a) + "개 매수 주문을 완료하였습니다.")
                    break
                else:
                    sell_market_test(coin, amount)
                    number.append(num2)
                    # 실주문 들어가야함
                    updater.bot.sendMessage(user_id,
                                            coin + " 코인의" + pri + " 가격에" + str(amount) + "개 " + maa + " 주문을 완료하였습니다.")
                    break
            except BinanceAPIException as a:
                if a.message == 'Filter failure: LOT_SIZE':
                    na -= 1
                    if son == 4:
                        b = round(float(b), na)
                    else:
                        amount = round(float(amount), na)
                elif a.message.find("'quantity'; legal range is") == True:
                    removefromorder(row=num2)
                    updater.bot.sendMessage(user_id,
                                            coin + " 코인의" + pri + " 가격에" + str(amount) + "개 " + maa + " 주문을 삭제하였습니다. - 수량 오류")
                    break
                else:
                    updater.bot.sendMessage(user_id,
                                            "Binance API Exception " + a.message + " " + maa + "주문을 수행하는 도중 오류가 발생하였습니다.")
                    break
            except Exception as e:
                updater.bot.sendMessage(user_id,
                                        "Exception " + str(e) + " " + maa + "주문을 수행하는 도중 오류가 발생하였습니다.")
                break
            chance -= 1
        return num2, number

    def price():
        with open("/data.txt", 'r') as f:
            lines = f.readlines()
            lines = list(map(lambda s: s.strip(), lines))
        if lines == []:
            time.sleep(10)
            return
        price = client.get_all_tickers()
        number = []
        num2 = 0
        try:
            for order in lines:
                order = order.split(',')
                coin = order[1].upper()
                pri = order[2]
                amount = order[3]
                # 퍼센트 빼고 개수로 변환하는 작업
                if amount.find('%') != -1 and order != '매수':
                    amount = amount.replace('%', '')
                    amount = float(client.get_asset_balance(asset=coin)['free']) * (float(amount) / 100)
                    amount = round(float(amount) / float(pri), 3)
                elif amount.find('%') != -1:
                    amount = amount.replace('%', '')
                    amount = float(client.get_asset_balance(asset='USDT')['free']) * (float(amount) / 100)
                    amount = round(float(amount) / float(pri), 3)

                if order[0] == '손절':
                    marketprice = 0
                    for cc in price:
                        if cc['symbol'] == coin + 'USDT':
                            marketprice = float(cc['price'])
                            break
                    if marketprice < float(pri):
                        num2, number = orderfunc(1, coin, amount, pri, num2, number, None)
                    num2 += 1

                elif order[0] == '익절':
                    marketprice = 0
                    for cc in price:
                        if cc['symbol'] == coin + 'USDT':
                            marketprice = float(cc['price'])
                            break
                    if marketprice >= float(pri):
                        num2, number = orderfunc(2, coin, amount, pri, num2, number, None)
                    num2 += 1

                elif order[0] == '매수':
                    marketprice = 0
                    for cc in price:
                        if cc['symbol'] == coin + 'USDT':
                            marketprice = float(cc['price'])
                            break
                    if marketprice < float(pri):
                        num2, number = orderfunc(4, coin, amount, pri, num2, number, marketprice)
                    num2 += 1

                elif order[0] == '익절예약':
                    targetpri = order[2]
                    sellpri = order[3]
                    amount = order[4]

                    if amount.find('%') != -1:
                        amount = amount.replace('%', '')
                        amount = float(client.get_asset_balance(asset=coin)['free']) * (float(amount) / 100)
                        amount = round(amount / sellpri, 3)
                    marketprice = 0
                    for cc in price:
                        if cc['symbol'] == coin + 'USDT':
                            marketprice = float(cc['price'])
                            break
                    if targetpri == 'o':
                        if float(sellpri) > marketprice:
                            num2, number = orderfunc(3, coin, amount, pri, num2, number, None)
                    elif marketprice >= float(targetpri):
                        removetargetprice(num2)
                        return
                    num2 += 1
                time.sleep(0.3)
            if number != []:
                removefromorder(number)



        except Exception:
            if number != []:
                removefromorder(number)
            return

    while True:
        try:
            price()
        except KeyboardInterrupt:
            return
        except:
            pass


createex()
loop()
