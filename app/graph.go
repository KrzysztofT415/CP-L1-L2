package app

import (
	"math/rand"
	"time"
)

type Graph struct {
	Source  chan<- *Packet
	Outfall <-chan *Packet

	VerticesInfo []*Vertex
}

func NewGraph(n int, d int) *Graph {
	newGraph := new(Graph)

	//CREATING VERTICES WITH BASIC EDGES (i, i + 1)

	newGraph.VerticesInfo = make([]*Vertex, n)
	newGraph.VerticesInfo[n-1] = NewVertex(n - 1)
	for i := n - 2; i >= 0; i-- {
		newGraph.VerticesInfo[i] = NewVertex(i)
		newGraph.VerticesInfo[i].NextVerticesInfo = append(newGraph.VerticesInfo[i].NextVerticesInfo, newGraph.VerticesInfo[i+1])
	}

	//ADDING RANDOM EDGES - SHORTCUTS (j, k)

	var added [][]int
	for i := 0; i < d; i++ {
		w := rand.New(rand.NewSource(time.Now().UnixNano())).Intn(n - 2)
		v := rand.New(rand.NewSource(time.Now().UnixNano())).Intn(n-w-2) + 2

		var wasUsed = false
		for i := range added {
			if added[i][0] == w && added[i][1] == v {
				wasUsed = true
			}
		}
		if !wasUsed {
			added = append(added, []int{w, v})
			newGraph.VerticesInfo[w].NextVerticesInfo = append(newGraph.VerticesInfo[w].NextVerticesInfo, newGraph.VerticesInfo[w+v])
			i++
		}
		i--
	}

	//CONNECTING CHANNELS

	for _, v := range newGraph.VerticesInfo {

		var channels = make([]chan<- *Packet, len(v.NextVerticesInfo))
		for i, c := range v.NextVerticesInfo {
			channels[i] = c.ThisVertexChannel
		}
		v.NextVerticesChannels = channels
	}

	newGraph.Source = newGraph.VerticesInfo[0].ThisVertexChannel
	output := make(chan *Packet, 0)
	newGraph.VerticesInfo[n-1].NextVerticesChannels = []chan<- *Packet{output}
	newGraph.Outfall = output

	return newGraph
}
