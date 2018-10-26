### Carrot

### Example

```go
import (
  "fmt"
  "github.com/PMoneda/carrot"
)


func main(){
  config := carrot.ConnectionConfig{
    Host:     os.Getenv("RABBITMQ_HOST"),
    Username: os.Getenv("RABBITMQ_USERNAME"),
    Password: os.Getenv("RABBITMQ_PASSWORD"),
    VHost:    "/",
  }
  exchangeName := "my_exchange"
  var conn *carrot.BrokerClient
  builder = carrot.NewBuilder(conn)
  builder.UseVHost("plataforma_v1.0")
  builder.DeclareTopicExchange(exchangeName)
  builder.DeclareTopicExchange(exchangeName + "_error")
  builder.UpdateTopicPermission(config.Username, exchangeName)
  builder.UpdateTopicPermission(config.Username, exchangeName+"_error")
  subConn, _ := carrot.NewBrokerClient(&config)
  subscriber = carrot.NewSubscriber(subConn)
  pubConn, _ := carrot.NewBrokerClient(&config)
  publisher = carrot.NewPublisher(pubConn)
  pickerConn, _ := carrot.NewBrokerClient(&config)
  picker = carrot.NewPicker(pickerConn)
}



```
