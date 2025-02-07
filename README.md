# ⚡ BitSet ⚡

This is a simple, though very fast, bitset implementation in Go derived from the [yourbasic/bit](github.com/yourbasic/bit) package.

## Installation

```go
go get github.com/KernelPryanic/bitset
```

## Usage

### Creating BitSets

```go
// Create an empty bitset
empty := bitset.New()

// Create a bitset with initial values
set := bitset.New(1, 3, 5, 7)
```

### Basic Operations

```go
// Add elements
set.Add(9)          // Add single element
set.AddRange(2, 5)  // Add range [2,3,4]

// Check membership
exists := set.Contains(3)    // true
absent := set.Contains(6)    // false

// Remove elements
set.Delete(3)               // Remove single element
set.DeleteRange(2, 5)       // Remove range [2,3,4]

// Get set information
size := set.Size()          // Number of elements
isEmpty := set.Empty()      // Check if set is empty
max := set.Max()           // Get maximum element

// Iterate through elements
set.Visit(func(n int) bool {
    fmt.Printf("%d ", n)
    return false  // Return true to stop iteration
})
```

### Set Operations

```go
set1 := bitset.New(1, 2, 3, 4, 5)
set2 := bitset.New(4, 5, 6, 7, 8)

// Create new sets from operations
intersection := bitset.And(set1, set2)  // Elements in both sets
union := bitset.Or(set1, set2)         // Elements in either set
symDiff := bitset.Xor(set1, set2)      // Elements in one set but not both
diff := bitset.AndNot(set1, set2)      // Elements in set1 but not in set2

// Modify existing sets
set1.And(set2)    // Keep only elements present in both sets
set1.Or(set2)     // Add all elements from set2
set1.Xor(set2)    // Keep elements present in one set but not both
set1.AndNot(set2) // Remove all elements present in set2
```

### Navigation

```go
set := bitset.New(1, 3, 5, 7, 9)

// Find next/previous elements
next := set.Next(4)     // Returns 5 (next element after 4)
prev := set.Prev(6)     // Returns 5 (previous element before 6)

// Copy sets
copy := set.Copy()      // Create a new copy
set2 := bitset.New()
set2.Set(set)          // Replace contents of set2 with set
```

### String Representation

```go
set := bitset.New(1, 2, 3, 5, 7, 8, 9, 10)
fmt.Println(set) // Outputs: {1..3 5 7..10}
```
