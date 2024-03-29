# MuffinDB

MuffinDB is a simple and efficient vector store database written in Go. It uses a KD-Tree data structure to store and retrieve vectors efficiently. The database supports basic operations like inserting, deleting, updating, and querying vectors based on their proximity to a given query vector.

## Features

- **Fast Vector Insertion**: MuffinDB allows you to insert vectors quickly using the `InsertVector` method.
- **Efficient Nearest Neighbor Search**: The KD-Tree data structure enables fast nearest neighbor searches, allowing you to retrieve the `k` nearest vectors to a given query vector using the `GetKNearestNeighbors` method.
- **Vector Retrieval**: You can retrieve a specific vector by its ID using the `GetVector` method.
- **Vector Deletion**: MuffinDB supports deleting vectors from the database using the `DeleteVector` method.
- **Vector Updating**: You can update an existing vector in the database using the `UpdateVector` method.
- **Persistence**: MuffinDB stores the vector data on disk using Go's built-in `encoding/gob` package, ensuring data persistence across program restarts.

## Getting Started

To get started with MuffinDB, follow these steps:

1. Install Go on your machine if you haven't already.
2. Clone the MuffinDB repository:

```bash
git clone https://github.com/sahildotexe/MuffinDB.git
```

3. Navigate to the project directory:
```bash
cd MuffinDB
```

4. Import the muffin package in your Go code and start using MuffinDB!

## Usage

```go
import "github.com/sahildotexe/MuffinDB/muffin"

func main() {
    // Connect to the vector store
    store := muffin.Connect()

    // Insert a vector
    store.InsertVector("This is a sample text", []float32{0.1, 0.2, 0.3})

    // Get the nearest neighbor to a query vector
    queryVector := []float32{0.15, 0.25, 0.35}
    neighbors := store.GetKNearestNeighbors(queryVector, 1)

    // Print the nearest neighbor's text
    fmt.Println(neighbors[0].Point.Text)
}
```

Contributions to MuffinDB are welcome! If you find any issues or have suggestions for improvements, please open an issue or submit a pull request.

