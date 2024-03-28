package kdtree

import (
	"fmt"
	"math"
	"sort"
)

type Vector struct {
	ID     int
	Values []float32
}

type KDTreeNode interface {
	isKDTreeNode()
}

type Leaf struct {
	Point Vector
}

type Internal struct {
	Left           KDTreeNode
	Right          KDTreeNode
	SplitValue     float32
	SplitDimension int
}

func (_ Leaf) isKDTreeNode()     {}
func (_ Internal) isKDTreeNode() {}

type KDTree struct {
	Root KDTreeNode
}

func BuildTree(points []Vector, depth int) KDTreeNode {
	if len(points) == 1 {
		return Leaf{points[0]}
	}
	dim := depth % len(points[0].Values) // Assuming all points have the same dimension
	sortedPoints := make([]Vector, len(points))
	copy(sortedPoints, points)
	sortByDimension(sortedPoints, dim)
	medianIdx := len(sortedPoints) / 2
	medianValue := sortedPoints[medianIdx].Values[dim]
	return Internal{
		Left:           BuildTree(sortedPoints[:medianIdx], depth+1),
		Right:          BuildTree(sortedPoints[medianIdx:], depth+1),
		SplitValue:     medianValue,
		SplitDimension: dim,
	}
}

func NewKDTree() *KDTree {
	return &KDTree{
		Root: nil,
	}
}

func (kdtree *KDTree) Insert(point Vector) {
	kdtree.Root = insert(kdtree.Root, point, 0)
}

func insert(node KDTreeNode, point Vector, depth int) KDTreeNode {
	if node == nil {
		return Leaf{point}
	}

	switch n := node.(type) {
	case Leaf:
		return Internal{
			Left:           n,
			Right:          Leaf{point},
			SplitValue:     point.Values[depth%len(point.Values)],
			SplitDimension: depth % len(point.Values),
		}
	case Internal:
		dim := depth % len(point.Values)
		if point.Values[dim] < n.SplitValue {
			n.Left = insert(n.Left, point, depth+1)
		} else {
			n.Right = insert(n.Right, point, depth+1)
		}
		return n
	default:
		panic("unexpected node type")
	}
}

func (kdtree KDTree) PrintTree() {
	printTree(kdtree.Root, 0)
}

func printTree(node KDTreeNode, depth int) {
	switch n := node.(type) {
	case Leaf:
		printSpaces(depth)
		printVector(n.Point)
	case Internal:
		printSpaces(depth)
		printInternal(n)
		printTree(n.Left, depth+1)
		printTree(n.Right, depth+1)
	default:
		panic("unexpected node type")
	}
}

func printSpaces(depth int) {
	for i := 0; i < depth; i++ {
		fmt.Print("  ")
	}
}

func printInternal(node Internal) {
	fmt.Printf("Split Dimension: %d, Split Value: %.02f\n", node.SplitDimension, node.SplitValue)
}

func printVector(v Vector) {
	fmt.Printf("ID: %d, Point: %v\n", v.ID, v.Values)
}

func sortByDimension(points []Vector, dim int) {
	sort.Slice(points, func(i, j int) bool {
		return points[i].Values[dim] < points[j].Values[dim]
	})
}

func (kdtree KDTree) GetNeighbours(query Vector, k int) []HeapVector {
	if k <= 0 {
		return nil
	}
	neighbors, _ := kdtree.nearest(query, kdtree.Root, make([]HeapVector, 0, k), math.MaxFloat32, k)
	return neighbors
}

func (kdtree KDTree) nearest(query Vector, node KDTreeNode, neighbors []HeapVector, bestDist float32, k int) ([]HeapVector, float32) {
	switch n := node.(type) {
	case Leaf:
		dist := euclideanDistance(query, n.Point)
		neighbors = pushHeap(neighbors, n.Point, dist, k)
		if len(neighbors) == k {
			bestDist = neighbors[0].Distance
		}
		return neighbors, bestDist
	case Internal:
		var nextNode KDTreeNode
		var otherNode KDTreeNode
		if query.Values[n.SplitDimension] < n.SplitValue {
			nextNode = n.Left
			otherNode = n.Right
		} else {
			nextNode = n.Right
			otherNode = n.Left
		}
		neighbors, bestDist = kdtree.nearest(query, nextNode, neighbors, bestDist, k)
		if math.Abs(float64(query.Values[n.SplitDimension]-n.SplitValue)) < float64(bestDist) {
			neighbors, bestDist = kdtree.nearest(query, otherNode, neighbors, bestDist, k)
		}
		return neighbors, bestDist
	default:
		panic("unexpected node type")
	}
}

func euclideanDistance(p1, p2 Vector) float32 {
	var sum float32
	for i := 0; i < len(p1.Values); i++ {
		sum += (p1.Values[i] - p2.Values[i]) * (p1.Values[i] - p2.Values[i])
	}
	return float32(math.Sqrt(float64(sum)))
}

type HeapVector struct {
	Point    Vector
	Distance float32
}

type HeapVectors []HeapVector

func (h HeapVectors) Len() int           { return len(h) }
func (h HeapVectors) Less(i, j int) bool { return h[i].Distance < h[j].Distance }
func (h HeapVectors) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func pushHeap(neighbors []HeapVector, point Vector, dist float32, k int) []HeapVector {
	heap := make(HeapVectors, len(neighbors), k)
	copy(heap, neighbors)
	heap = heap[:min(len(heap)+1, k)]
	heap[len(heap)-1] = HeapVector{Point: point, Distance: dist}
	sort.Sort(heap)
	return heap
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
