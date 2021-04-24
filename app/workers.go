package app

import (
	"fmt"
	"math/rand"
	"time"
)

const DELAY = 5

func Sender(source chan<- *Packet, packets []*Packet) {
	for i := 0; i < len(packets); i++ {
		time.Sleep(time.Second * time.Duration(rand.New(rand.NewSource(time.Now().UnixNano())).Intn(DELAY)))
		source <- packets[i]
	}
	close(source)
}

func Receiver(outfall <-chan *Packet, done chan<- bool, print chan<- string) {
	for packet := range outfall {
		print <- fmt.Sprintf("\tPacket received %d\n", packet.Id)
		time.Sleep(time.Second * time.Duration(rand.New(rand.NewSource(time.Now().UnixNano())).Intn(DELAY)))
	}
	done <- true
}

func Transmitter(vertex *Vertex, print chan<- string) {
	for packet := range vertex.ThisVertexChannel {
		print <- fmt.Sprintf("Packet %d is in %d\n", packet.Id, vertex.Id)

		vertex.PacketsSeenInfo = append(vertex.PacketsSeenInfo, packet.Id)
		packet.VerticesVisited = append(packet.VerticesVisited, vertex.Id)

		time.Sleep(time.Second * time.Duration(rand.New(rand.NewSource(time.Now().UnixNano())).Intn(DELAY)))

		r := rand.New(rand.NewSource(time.Now().UnixNano())).Intn(len(vertex.NextVerticesChannels))
		vertex.NextVerticesChannels[r] <- packet

		time.Sleep(time.Millisecond * time.Duration(50))
	}
	close(vertex.NextVerticesChannels[0])
}
