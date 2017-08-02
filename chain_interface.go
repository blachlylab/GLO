package glo

import (
    "hash/fnv"
)

import "github.com/Workiva/go-datastructures/augmentedtree"

// Implement Interval interface functions for ChainInterval
func (ci ChainInterval) LowAtDimension(dim uint64) int64 {
    return ci.Start
}

func (ci ChainInterval) HighAtDimension(dim uint64) int64 {
    return ci.End
}

func (ci ChainInterval) OverlapsAtDimension(iv augmentedtree.Interval, dim uint64) bool {
    if (iv.LowAtDimension(dim) <= ci.Start) && (ci.End <= iv.HighAtDimension(dim)) {
        // self       ================
        // other   =====================
        return true
    } else if (ci.Start <= iv.LowAtDimension(dim)) && (iv.LowAtDimension(dim) <= ci.End) {
        // self      ================
        // other         ===============
        return true
    } else if (ci.Start <= iv.HighAtDimension(dim)) && (iv.HighAtDimension(dim) <= ci.End) {
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


// Implement Interval interface functions for ChainBlock by taking advantage
// of the implemented functions for ChainInterval
func (cb *ChainBlock) LowAtDimension(dim uint64) int64 {
    return cb.Source.LowAtDimension(dim)
}

func (cb *ChainBlock) HighAtDimension(dim uint64) int64 {
    return cb.Source.HighAtDimension(dim)
}

func (cb *ChainBlock) OverlapsAtDimension(iv augmentedtree.Interval, dim uint64) bool {
    return cb.Source.OverlapsAtDimension(iv, dim)
}

func (cb *ChainBlock) ID() uint64 {
    h := fnv.New64a()
    h.Write([]byte(cb.String()))
    return h.Sum64()
}