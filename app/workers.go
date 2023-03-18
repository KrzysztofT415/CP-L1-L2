package app

import (
	"fmt"
	"math/rand"
	"time"
)

const DELAY = 2

func Sender(source chan<- *Packet, packets []*Packet) {
	for i := 0; i < len(packets); i++ {
		time.Sleep(time.Second * time.Duration(rand.New(rand.NewSource(time.Now().UnixNano())).Intn(DELAY)))
		source <- packets[i]
	}
}

func Receiver(outfall <-chan *Packet, print chan<- string, trash chan<- *Packet) {
	for packet := range outfall {
		print <- fmt.Sprintf("\tPacket received %d\n", packet.Id)
		time.Sleep(time.Second * time.Duration(rand.New(rand.NewSource(time.Now().UnixNano())).Intn(DELAY)))
		trash <- packet
	}
}

func Transmitter(vertex *Vertex, print chan<- string, trash chan<- *Packet) {
	for packet := range vertex.ThisVertexChannel {
		print <- fmt.Sprintf("Packet %d is in %d\n", packet.Id, vertex.Id)

		vertex.PacketsSeenInfo = append(vertex.PacketsSeenInfo, packet.Id)
		packet.VerticesVisited = append(packet.VerticesVisited, vertex.Id)

		time.Sleep(time.Second * time.Duration(rand.New(rand.NewSource(time.Now().UnixNano())).Intn(DELAY)))

		if vertex.IsTrapped {
			print <- fmt.Sprintf("\tPacket %d was caught in %d\n", packet.Id, vertex.Id)
			packet.Life = -packet.Life
			vertex.IsTrapped = false
			vertex.HasPacket = false
			trash <- packet
			continue
		}

		packet.Life--
		if packet.Life == 0 {
			print <- fmt.Sprintf("\tPacket %d died in %d\n", packet.Id, vertex.Id)
			vertex.HasPacket = false
			trash <- packet
			continue
		}

		var nextVertexId = -1
		for {
			vertex.Rounds++

			nextVerticesIds := make([]int, len(vertex.NextVerticesChannels))
			for i := range nextVerticesIds {
				nextVerticesIds[i] = i
			}
			rand.Shuffle(len(nextVerticesIds), func(i, j int) { nextVerticesIds[i], nextVerticesIds[j] = nextVerticesIds[j], nextVerticesIds[i] })

			for _, id := range nextVerticesIds {
				vertex.Tries++
				if len(vertex.NextVerticesAskChannels) < len(vertex.NextVerticesChannels) && id == len(vertex.NextVerticesChannels) - 1 {
					nextVertexId = len(vertex.NextVerticesChannels) - 1
					break
				}

				vertex.NextVerticesAskChannels[id] <- vertex.ThisVertexResponseChannel
				isOccupied := <-vertex.ThisVertexResponseChannel
				if !isOccupied {
					nextVertexId = id
					break
				}
			}
			if nextVertexId != -1 {
				break
			}

			time.Sleep(time.Second * time.Duration(DELAY * 4))
		}

		vertex.NextVerticesChannels[nextVertexId] <- packet
		time.Sleep(time.Millisecond * time.Duration(50))
		vertex.HasPacket = false
	}

	if len(vertex.NextVerticesAskChannels) == len(vertex.NextVerticesChannels) {
		close(vertex.NextVerticesChannels[0])
	}
	close(vertex.ThisVertexAskChannel)
	close(vertex.ThisVertexResponseChannel)
}

func Watcher(vertex *Vertex) {
	for responseChannel := range vertex.ThisVertexAskChannel {
		if vertex.HasPacket {
			responseChannel <- true
		} else {
			vertex.HasPacket = true
			responseChannel <- false
		}
	}
}

func TrashCounter(trash chan *Packet, packetsNum int, done chan bool) {
	var counter = 0
	for range trash {
		counter++
		if counter == packetsNum {
			break
		}
	}
	done <- true
}

func Hunter(graph *Graph, print chan<- string) {
	for {
		time.Sleep(time.Second * time.Duration(rand.New(rand.NewSource(time.Now().UnixNano())).Intn(DELAY * 3) + DELAY * 7))
		vertex := rand.New(rand.NewSource(time.Now().UnixNano())).Intn(len(graph.VerticesInfo))
		print <- fmt.Sprintf("\tTrap set on %d\n", vertex)
		graph.VerticesInfo[vertex].IsTrapped = true
	}
}
