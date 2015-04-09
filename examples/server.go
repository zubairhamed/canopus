package main
import "github.com/zubairhamed/goap"

func main() {
    server := goap.NewLocalServer()

    server.OnDiscover(request, response) {

    }

    server.OnError(request, error, errorCode) {

    }
}
