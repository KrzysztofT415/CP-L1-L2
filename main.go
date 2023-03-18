package main

import (
	"fmt"
	"math"
	"sort"

	"./app"
)

func main() {
	//n - vertices_num, d - shortcuts, k - packets, b - backtracks, h - health
	n, d, k, b, h := 10, 28, 15, 38, 20
	//n, d, k, b, h := scanParameters()

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

	graph := app.NewGraph(n, d, b)

	printGraph(graph, printServerChannel)

	doneTransmitting := make(chan bool)
	trashChannel := make(chan *app.Packet)
	go app.TrashCounter(trashChannel, k, doneTransmitting)

	packets := make([]*app.Packet, k)
	for i := 0; i < k; i++ {
		packets[i] = app.NewPacket(i, h)
	}

	printServerChannel <- "START OF TRANSFER\n\n"
	for _, vertex := range graph.VerticesInfo {
		go app.Transmitter(vertex, printServerChannel, trashChannel)
		go app.Watcher(vertex)
	}
	go app.Receiver(graph.Outfall, printServerChannel, trashChannel)
	go app.Sender(graph.Source, packets)
	go app.Hunter(graph, printServerChannel)

	<-doneTransmitting
	close(graph.Source)
	close(doneTransmitting)
	printServerChannel <- "\nEND OF TRANSFER\n---------------\n   LOG\n"

	printGraphVerticesInfo(graph, printServerChannel, packets)
	printPacketsInfo(packets, printServerChannel, n)

	<-donePrinting
	close(donePrinting)
}

func scanParameters() (int, int, int, int, int) {
	var n, d, k, b, h int
	fmt.Print("N : ")
	_, _ = fmt.Scanf("%d", &n)
	n = int(math.Max(float64(n), 1))
	fmt.Print("D (", (n*(n-3)/2)+1, ") : ")
	_, _ = fmt.Scanf("%d", &d)
	d = int(math.Min(float64(d), float64((n*(n-3))/2)+1))
	fmt.Print("K : ")
	_, _ = fmt.Scanf("%d", &k)
	fmt.Print("B : ")
	_, _ = fmt.Scanf("%d", &b)
	d = int(math.Min(float64(d), float64(n*(n-1))/2))
	fmt.Print("H : ")
	_, _ = fmt.Scanf("%d", &h)
	return n, d, k, b, h
}

func printGraph(graph *app.Graph, print chan<- string) {
	print <- "---------------\n  GRAPH\nId:  Links:\n 0 == SOURCE\n"
	for _, v := range graph.VerticesInfo {
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
		printServerChannel <- " |"
		for _, p := range v.PacketsSeenInfo {
			printServerChannel <- fmt.Sprintf(" %d", p)
		}
		printServerChannel <- fmt.Sprintf(" | %d packets, %d rounds and %d tries\n", len(v.PacketsSeenInfo), v.Rounds, v.Tries)
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

		printServerChannel <- " |"
		for _, v := range p.VerticesVisited {
			printServerChannel <- fmt.Sprintf(" %d", v)
		}
		if p.Life < 0 {
			printServerChannel <- fmt.Sprintf(" | been in %d vertices, %d lives left, was caught by hunter\n", len(p.VerticesVisited), -p.Life)
		} else if p.Life == 0 {
			printServerChannel <- fmt.Sprintf(" | been in %d vertices, died\n", len(p.VerticesVisited))
		} else {
			printServerChannel <- fmt.Sprintf(" | been in %d vertices, %d lives left, was received\n", len(p.VerticesVisited), p.Life)
		}
	}
	printServerChannel <- "---------------\n"
	printServerChannel <- "EOF"
}
