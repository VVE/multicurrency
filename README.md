# multicurrency
This is an example to study Go bots for Telegram messanger. 

The bot can take into account the operations of clients using their cryptocurrency wallets.
Customers can:
* add and subtract money in the wallet,
* remove currency from the wallet,
* check the balance of the cryptocurrency wallet.

The bot does not conduct real money transactions and does not serve currency conversions.

To do this, they enter the following commands:
* `ADD <currency name> <amount>` - add the amount,
* `SUB <currency name> <amount>` - withdraw the amount,
* `SHOW` – show balance,
* `DEL <currency name>` - delete the currency.

Here:
<currency name> is cryptocurrency code,
<amount> is the amount in cryptocurrency.

The exchange rates are obtained from binance.com using the API:
* Request is like: https://api.binance.com/api/v3/ticker/price?symbol=BTCRUB

Here the first 3 characters in the symbol parameter are the cryptocurrency code (BTC is bitcoin, ETH is ethereum),
the second 3 characters are the code of the currency in which we are evaluating.
* Response is like: {"symbol":"BTCRUB","price":"3640630.00000000"}

The bot is not case sensitive.

The balance for each cryptocurrency is displayed in units of this cryptocurrency and in rubles, indicating the current exchange rate for reference.

The Russian currency can be replaced with another by replacing constants.

If the balance is insufficient for withdrawal, the bot reports this.

The bot does not allow you to delete cryptocurrency with a non-zero balance.

# Running the code
* set bot token into return operator in getToken() function
* compile and run the bot program: for Windows OS type
`C:\Users\Vova\go\src\multicurrency>go run main.go`
in terminal, similarly for another OS.
The cursor will move to the beginning of the next line, but nothing will be output.
* Launch Telegram messenger, find the bot there by any of its names: multicurrency or multi_currency_test_bot.
* Enter commands one at a time and get answers.
* To stop the bot program, enter Ctrl+C in the terminal.
* Exit the Telegram chat-as usual.
The order of the bot program and the chat termination is not important. But at the end of the program, the balance is not saved, since the database is stored in memory.

Implementation details are:
for simplicity, the database an in-memory structure:
```
{
<user code>:{
<currency code>: <amount>,
<currency code>: <amount>},
},
<user code>:{
<currency code>: <amount>,
<currency code>: <amount>},
},
}
```
Here:
`<user code>` is a natural number - userID in Telegram,
`<currency code>` is a string,
`<sum>` is a non-negative real number.

The bot was created under the leadership of Valeriy Kabisov, Senior Golang Developer, during a 3-day intensive training, organized by Skillbox on August 19-21, 2021.

