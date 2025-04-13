# container: the missing piece of the Go standard library

container implements efficient slice-backed data structures that would probably
have been included in Go's standard library if generics had been available from
the start. Package deque implements a double-ended queue inspired by Rust's
wonderful VecDeque type. Package heap is a reimagining of the standard library
container/heap with a generics-first implementation.
