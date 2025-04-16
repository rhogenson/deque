# deque: a high-performance slice-backed double-ended queue inspired by Rust's VecDeque

Compared to other popular slice-backed deque implementations, this one

 - is only 32 bytes;
 - uses append to get an optimial growth factor;
 - supports iterating using Go 1.23 iterators;
 - and steals Rust's clever strategy for minimizing the amount of data copied
   on reallocation.
