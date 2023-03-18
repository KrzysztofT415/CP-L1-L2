package main

import (
	"./app"
	"fmt"
	"math"
	"sort"
)

func main() {
	n, d, k := 6, 7, 5
	//n, d, k := scanParameters()

	printServerChannel := make(chan string, 50)
	donePrinting := make(chan bool)
	go func() {
		for info := range printServerChannel {
			if info == "EOF" {
				break
			}
			fmt.Print(info)
		}
		close(printServerChannel)
		donePrinting <- true
	}()

	graph := app.NewGraph(n, d)

	printGraph(graph, printServerChannel)

	doneTransmitting := make(chan bool)
	packets := make([]*app.Packet, k)
	for i := 0; i < k; i++ {
		packets[i] = app.NewPacket(i)
	}

	printServerChannel <- "START OF TRANSFER\n\n"
	for _, vertex := range graph.VerticesInfo {
		go app.Transmitter(vertex, printServerChannel)
	}
	go app.Receiver(graph.Outfall, doneTransmitting, printServerChannel)
	go app.Sender(graph.Source, packets)

	<-doneTransmitting
	close(doneTransmitting)
	printServerChannel <- "\nEND OF TRANSFER\n---------------\n   LOG\n"

	printGraphVerticesInfo(graph, printServerChannel, packets)
	printPacketsInfo(packets, printServerChannel, n)

	<-donePrinting
	close(donePrinting)
}

func scanParameters() (int, int, int) {
	var n, d, k int
	fmt.Print("N : ")
	_, _ = fmt.Scanf("%d", &n)
	n = int(math.Max(float64(n), 1))
	fmt.Print("D (", (n*(n-3)/2)+1, ") : ")
	_, _ = fmt.Scanf("%d", &d)
	d = int(math.Min(float64(d), float64((n*(n-3))/2)+1))
	fmt.Print("K : ")
	_, _ = fmt.Scanf("%d", &k)
	return n, d, k
}

func printGraph(graph *app.Graph, print chan<- string) {
	print <- "---------------\n  GRAPH\nId:  Links:\n 0 == SOURCE\n"
	for _, v := range graph.VerticesInfo[:len(graph.VerticesInfo)-1] {
		print <- fmt.Sprintf(" %d -> ", v.Id)
		for _, c := range v.NextVerticesInfo {
			print <- fmt.Sprintf("%d ", c.Id)
		}
		print <- "\n"
	}
	print <- fmt.Sprintf(" %d == OUTFALL\n---------------\n", graph.VerticesInfo[len(graph.VerticesInfo)-1].Id)
}

func printGraphVerticesInfo(graph *app.Graph, printServerChannel chan string, packets []*app.Packet) {
	printServerChannel <- "VERTICES :\n"
	for _, v := range graph.VerticesInfo {
		if v.Id/10 == 0 {
			printServerChannel <- " "
		}
		printServerChannel <- fmt.Sprintf("%d |", v.Id)
		sort.Ints(v.PacketsSeenInfo)
		for i := 0; i < len(packets); i++ {
			var seen = false
			for _, packet := range v.PacketsSeenInfo {
				if packet == i {
					seen = true
				} else if packet > i {
					break
				}
			}
			if !seen && i/10 > 0 {
				printServerChannel <- "   "
			} else if !seen {
				printServerChannel <- "  "
			} else {
				printServerChannel <- fmt.Sprintf(" %d", i)
			}

		}
		printServerChannel <- "\n"
	}
}

func printPacketsInfo(packets []*app.Packet, printServerChannel chan string, n int) {
	printServerChannel <- "\nPACKETS :\n"
	for _, p := range packets {
		if p.Id/10 == 0 {
			printServerChannel <- " "
		}
		printServerChannel <- fmt.Sprintf("%d |", p.Id)
		for i := 0; i < n; i++ {
			var seen = false
			for _, vertex := range p.VerticesVisited {
				if vertex == i {
					seen = true
				} else if vertex > i {
					break
				}
			}
			if !seen && i/10 > 0 {
				printServerChannel <- "   "
			} else if !seen {
				printServerChannel <- "  "
			} else {
				printServerChannel <- fmt.Sprintf(" %d", i)
			}

		}
		printServerChannel <- "\n"
	}
	printServerChannel <- "---------------\n"
	printServerChannel <- "EOF"
}
