package main

import (
    "fmt"
    "strings"
    "hash/fnv"
)

import "github.com/Workiva/go-datastructures/augmentedtree"

type ChainInterval struct {
    contig string
    start int64
    end int64
}

func (ci *ChainInterval) String() string {
    return fmt.Sprintf("%s:%d-%d", ci.contig, ci.start, ci.end)
}

// Implement Interval interface functions for ChainInterval
func (ci ChainInterval) LowAtDimension(dim uint64) int64 {
    return ci.start
}

func (ci ChainInterval) HighAtDimension(dim uint64) int64 {
    return ci.end
}

func (ci ChainInterval) OverlapsAtDimension(iv augmentedtree.Interval, dim uint64) bool {
    if (iv.LowAtDimension(dim) <= ci.start) && (ci.end <= iv.HighAtDimension(dim)) {
        // self       ================
        // other   =====================
        return true
    } else if (ci.start <= iv.LowAtDimension(dim)) && (iv.LowAtDimension(dim) <= ci.end) {
        // self      ================
        // other         ===============
        return true
    } else if (ci.start <= iv.HighAtDimension(dim)) && (iv.HighAtDimension(dim) <= ci.end) {
        // self      ===============
        // other  =================
        return true
    }
    return false
}

func (ci ChainInterval) ID() uint64 {
    h := fnv.New64a()
    h.Write([]byte(ci.String()))
    return h.Sum64()
}

// done with Interval interface methods for ChainInterval

type ChainBlock struct {
    source *ChainInterval
    target *ChainInterval
}

func (cb *ChainBlock) String() string {
    return fmt.Sprintf("%s -> %s", cb.source, cb.target)
}

func (cb *ChainBlock) GetOverlap(contig string, start, end int64) *ChainInterval {
    ci := new(ChainInterval)
    if contig != cb.source.contig {
        // No overlap due to contig mismatch
        return nil
    }

    ci.contig = cb.target.contig


    var start_adj int64 = 0

    if start > cb.source.start {
        start_adj = start - cb.source.start
    }
    ci.start = cb.target.start + start_adj


    size := end - start
    if end > cb.source.end {
        size -= (end - cb.source.end)
    }
    ci.end = ci.start + size

    return ci
}


type Chain struct {
    score int64
    source_name string
    source_size int64
    source_strand string
    source_start int64
    source_end int64
    target_name string
    target_size int64  
    target_strand string
    target_start int64
    target_end int64
    id string
    blocks []*ChainBlock
}

func (c *Chain) String() string {
    var output []string
    output = append(output, fmt.Sprintf("%s:%d-%d to %s:%d-%d", c.source_name, c.source_start, c.source_end, c.target_name, c.target_start, c.target_end))
    for _, x := range c.blocks {
        output = append(output, fmt.Sprintf("> %s", x))
    }
    return strings.Join(output, "\n")
}

// Implement Interval interface functions for ChainBlock by taking advantage
// of the implemented functions for ChainInterval
func (cb *ChainBlock) LowAtDimension(dim uint64) int64 {
    return cb.source.LowAtDimension(dim)
}

func (cb *ChainBlock) HighAtDimension(dim uint64) int64 {
    return cb.source.HighAtDimension(dim)
}

func (cb *ChainBlock) OverlapsAtDimension(iv augmentedtree.Interval, dim uint64) bool {
    return cb.source.OverlapsAtDimension(iv, dim)
}

func (cb *ChainBlock) ID() uint64 {
    h := fnv.New64a()
    h.Write([]byte(cb.String()))
    return h.Sum64()
}