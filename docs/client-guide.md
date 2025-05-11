# Руководство по созданию клиентов для сервиса отправки push-уведомлений
## Общие сведения о сервисе push-уведомлений
Сервис push-уведомлений может быть использован для отправки сообщений клиентам, которые подключены к нему с помощью Websockets.
## Технологии необходимые для подключения к сервису
Для подключения к сервису необходимо использовать любую библиотеку по работе с Websockets со стандартными настройками, а также библиотека для работы с сообщениями protobuf - разбор сообщений и их формирование.
## Формат сообщений
Ниже приведены возможные сообщения в рамках канала Websockets. При обмене сообщений сервер всегда отправляет ServerMessage и всегда ожидает получения ClientMessage. Все возможные вариации осуществляются с помощью данных сущностей.
```protobuf
message ServerMessage {
    oneof server_message {
        PushNotification push_notification = 1;    
    }
}

message PushNotification {
    string id = 1; // uuid format
    string created_at = 2; // rfc 3339 format
    string title = 3;
    string body = 4;
}

message ClientMessage {
    oneof client_message {
        DeliveryAck delivery_ack = 1;
    }
}

message DeliveryAck {
    enum Status {
        STATUS_ACK = 0;
        STATUS_NACK = 1;
    }

    string push_id = 1; // uuid format
    Status status = 2;
}
```
## Протокол обмена сообщениями
После подключения к Websockets клиент обязан находиться на данном соединении бесконечное количество времени, пока не случится разрыв соединения. За счёт этого обеспечивается мгновенная доставка сообщения. В самом начале может не приходить никаких сообщений, но как только появится запрос на отправку push-уведомления, сервер отправит ServerMessage с соответствующим наполнением.

Если клиент не будет отвечать на первое отправленное ServerMessage, то сервер произведёт переотправку сообщения. После некоторого количества попыток и отсутствия ответа сервер насильно разорвёт соединение.

Сервер также может разорвать соединение если будут отсутствовать служебные сообщения ping-pong.

После разрыва соединения и восстановления сети клиент обязан восстановить соединение с сервером.
## Пример реализации
Ниже приведён пример реализации на языке программирования Go. Он демонстрирует разбор сообщения от сервера, его возможную обработку и ответ на данное сообщение.
```go
{
    // msg является массивом байт, полученными из Websockets
    var serverMessage protocol.ServerMessage
    err := proto.Unmarshal(msg, &serverMessage)
    if err != nil {
        slog.Error("error with unmarshalling msg: " + err.Error())
        continue
    }

    switch val := serverMessage.ServerMessage.(type) {
    case *protocol.ServerMessage_PushNotification:
        slog.Info(fmt.Sprintf("received push: %+v", val.PushNotification))

        err := sendAck(writeChan, val.PushNotification.Id)
        if err != nil {
            slog.Error("error with sending ack: " + err.Error())
        }
    }
}

func sendAck(writeChan chan []byte, pushID string) error {
	clientMessage := protocol.ClientMessage{
		ClientMessage: &protocol.ClientMessage_DeliveryAck{
			DeliveryAck: &protocol.DeliveryAck{
				PushId: pushID,
				Status: protocol.DeliveryAck_STATUS_ACK,
			},
		},
	}

	msg, err := proto.Marshal(&clientMessage)
	if err != nil {
		return errors.New("error with marshalling client message: " + err.Error())
	}

	sendMessage(msg) // реализация отправки в Websockets

	return nil
}
```