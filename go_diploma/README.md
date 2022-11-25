### Для корректной работы приложения требуется выполнить следующие шаги:

1. Запустить эмулятор **go run skillbox-diploma/main.go**;
2. В проектном файле **go_diploma/settings.toml** свойстве [Emu].Path прописать полный путь до директории в которой эмулятор сгенерировал датафайлы. Директория самого эмулятора;
3. Перейти в папку проекта, собрать и запустить командой: **go run go_diploma/cmd/main.go**;
4. Чтобы проверить проектный сервис на доступность требуется выполнить get-запрос  на endpoint **http://127.0.0.1:8585/**. Если сервис рабочий и настроен, вернётся "ОК";
5. Выполнить get-запрос по умолчанию на endpoint **http://127.0.0.1:8585/test**. Достаточно в браузере перейти на указанный адрес или использовать другой клиент, например, утилиту **curl** или **Postman**.


#### Ожидаемый результат:

1. Произойдёт сбор данных из настроенных источников;
2. Вернётся структура типа service.ResultT сконвертированная в формат JSON заполненная данными и статусом равным 'true' в нормальном случае. В случае сбоя вернётся эта структура без данных, со статусом равным 'false' и описанием ошибки.


#### Пример запроса:

curl: **curl http://127.0.0.1:8585/test**


#### Примеры ответов:

##### Сбой:
```json
{
	"status": false,
	"data": {
		"sms": null,
		"mms": null,
		"voice_call": null,
		"email": null,
		"billing": {
			"CreateCustomer": false,
			"Purchase": false,
			"Payout": false,
			"Recurring": false,
			"FraudControl": false,
			"CheckoutPage": false
		},
		"support": null,
		"incident": null
	},
	"error": "open C:/Users/yunusov/go/skillbox-diploma/sms.data1: The system cannot find the file specified."
}
```

##### Успешно:
```json
{
	"status": true,
	"data": {
		"sms": [
			[
				{
					"Country": "AT",
					"Bandwidth": "9",
					"ResponseTime": "1263",
					"Provider": "Topolo"
				},
				{
					"Country": "BG",
					"Bandwidth": "59",
					"ResponseTime": "1791",
					"Provider": "Rond"
				}
			],
			[
				{
					"Country": "Saint Barthélemy",
					"Bandwidth": "65",
					"ResponseTime": "1096",
					"Provider": "Kildy"
				},
				{
					"Country": "New Zealand",
					"Bandwidth": "40",
					"ResponseTime": "95",
					"Provider": "Kildy"
				}
			]
		],
		"mms": [
			[
				{
					"country": "AT",
					"provider": "Topolo",
					"bandwidth": "69",
					"response_time": "924"
				},
				{
					"country": "BG",
					"provider": "Rond",
					"bandwidth": "12",
					"response_time": "1583"
				}
			],
			[
				{
					"country": "Saint Barthélemy",
					"provider": "Kildy",
					"bandwidth": "98",
					"response_time": "508"
				},
				{
					"country": "New Zealand",
					"provider": "Kildy",
					"bandwidth": "87",
					"response_time": "945"
				}
			]
		],
		"voice_call": [
			{
				"Country": "RU",
				"Load": 13,
				"ResponseTime": 1420,
				"Provider": "TransparentCalls",
				"Stability": 0.78,
				"TtfbClearence": 193,
				"CallDuration": 86,
				"UnknowValue": 28
			},
			{
				"Country": "US",
				"Load": 69,
				"ResponseTime": 504,
				"Provider": "E-Voice",
				"Stability": 0.96,
				"TtfbClearence": 136,
				"CallDuration": 3,
				"UnknowValue": 50
			}
		],
		"email": {
			"AT": [
				[
					{
						"Country": "AT",
						"Provider": "GMX",
						"DeliveryTime": 63
					},
					{
						"Country": "AT",
						"Provider": "Orange",
						"DeliveryTime": 64
					},
					{
						"Country": "AT",
						"Provider": "Yahoo",
						"DeliveryTime": 91
					}
				],
				[
					{
						"Country": "AT",
						"Provider": "Mail.ru",
						"DeliveryTime": 588
					},
					{
						"Country": "AT",
						"Provider": "Live",
						"DeliveryTime": 549
					},
					{
						"Country": "AT",
						"Provider": "Gmail",
						"DeliveryTime": 350
					}
				]
			],
			"BG": [
				[
					{
						"Country": "BG",
						"Provider": "GMX",
						"DeliveryTime": 30
					},
					{
						"Country": "BG",
						"Provider": "Yandex",
						"DeliveryTime": 35
					},
					{
						"Country": "BG",
						"Provider": "Yahoo",
						"DeliveryTime": 76
					}
				],
				[
					{
						"Country": "BG",
						"Provider": "Comcast",
						"DeliveryTime": 579
					},
					{
						"Country": "BG",
						"Provider": "MSN",
						"DeliveryTime": 505
					},
					{
						"Country": "BG",
						"Provider": "Mail.ru",
						"DeliveryTime": 457
					}
				]
			]
		},
		"billing": {
			"CreateCustomer": true,
			"Purchase": true,
			"Payout": false,
			"Recurring": true,
			"FraudControl": true,
			"CheckoutPage": true
		},
		"support": [
			2,
			15
		],
		"incident": [
			{
				"topic": "SMS delivery in EU",
				"status": "closed"
			},
			{
				"topic": "MMS connection stability",
				"status": "closed"
			},
			{
				"topic": "Voice call connection purity",
				"status": "closed"
			},
			{
				"topic": "Checkout page is down",
				"status": "closed"
			},
			{
				"topic": "Support overload",
				"status": "closed"
			},
			{
				"topic": "Buy phone number not working in US",
				"status": "closed"
			},
			{
				"topic": "API Slow latency",
				"status": "closed"
			}
		]
	},
	"error": ""
}
```
