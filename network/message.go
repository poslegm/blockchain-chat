package network

// TODO сообщение, идущее по сети, должно быть сериализованно в массив байтов;
// TODO сделать интерфейс для преобразования сообщения от клиента в вид для сети;
// TODO так же по сети могут передаваться блоки, надо держать это в голове;

const MESSAGE = "MESSAGE"
const REQUEST = "REQUEST"
type NetworkMessage struct {
	MessageType string
	IP string

	// TODO убрать и использовть дату
	Receiver string
	Sender string
	Text string

	data []byte
}

func CreateTextMessage(receiver, sender, text string) NetworkMessage {
	return NetworkMessage{
		Receiver:receiver,
		Sender:sender,
		Text:text,
	}
}