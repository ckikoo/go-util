package seq

import (
    "sync"
    "sync/atomic"
)

//======转换,添加或修改内部元素========

// MapParallel 每个元素转换为any,使用 Sync() 保证消费不竞争
// order 是否保持顺序,1尽量保持顺序(可能消费竞争),大于1强制保持顺序(约等于加锁)
// order 第二个参数,并发数
func (t Seq[T]) MapParallel(f func(T) any, order ...int) Seq[any] {
    o := 0
    sl := 0
    if len(order) > 0 {
        o = order[0]
    }
    if len(order) > 1 {
        sl = order[1]
    }
    if o > 0 {
        p := NewParallel(sl)
        l := sync.NewCond(&sync.Mutex{})
        return func(c func(any)) {
            var currentIndex int32 = 1
            var id int32
            t(func(t T) {
                var id = atomic.AddInt32(&id, 1)
                p.Add(func() {
                    a := f(t)
                    l.L.Lock()
                    for atomic.LoadInt32(&currentIndex) != id {
                        l.Wait()
                    }
                    atomic.AddInt32(&currentIndex, 1)
                    defer l.Broadcast()
                    if o > 1 {
                        defer l.L.Unlock()
                    } else {
                        l.L.Unlock()
                    }
                    c(a)
                })
            })
            p.Wait()
        }
    } else {
        return t.Parallel(sl).Map(f)
    }
}

// Map 每个元素转换为any
func (t Seq[T]) Map(f func(T) any) Seq[any] {
    return func(c func(any)) { t(func(t T) { c(f(t)) }) }
}

// MapFlat 每个元素转换为Seq,并扁平化
func (t Seq[T]) MapFlat(f func(T) Seq[any]) Seq[any] {
    return func(c func(any)) { t(func(t T) { f(t).ForEach(c) }) }
}

// MapString 每个元素转换为字符串
func (t Seq[T]) MapString(f func(T) string) Seq[string] {
    return func(c func(string)) { t(func(t T) { c(f(t)) }) }
}

func (t Seq[T]) MergeBiInt(iterator Iterator[int]) BiSeq[int, T] {
    return BiFrom(func(f1 func(int, T)) {
        t.ConsumeTillStop(func(t T) {
            x, exist := iterator()
            if !exist {
                panic(&Stop)
            }
            f1(x, t)
        })
    })
}

func (t Seq[T]) MergeBiIntRight(iterator Iterator[int]) BiSeq[T, int] {
    return BiFrom(func(f1 func(T, int)) {
        t.ConsumeTillStop(func(t T) {
            x, exist := iterator()
            if !exist {
                panic(&Stop)
            }
            f1(t, x)
        })
    })
}

func (t Seq[T]) MergeBiAny(iterator Iterator[any]) BiSeq[any, T] {
    return BiFrom(func(f1 func(any, T)) {
        t.ConsumeTillStop(func(t T) {
            x, exist := iterator()
            if !exist {
                panic(&Stop)
            }
            f1(x, t)
        })
    })
}

func (t Seq[T]) MergeBiAnyRight(iterator Iterator[any]) BiSeq[any, T] {
    return BiFrom(func(f1 func(any, T)) {
        t.ConsumeTillStop(func(t T) {
            x, exist := iterator()
            if !exist {
                panic(&Stop)
            }
            f1(x, t)
        })
    })
}

// MapBiSerialNumber 为每个元素生成一个序列号,并转换为BiSeq,参数rang为规则,参考 IteratorInt
func (t Seq[T]) MapBiSerialNumber(rang ...int) BiSeq[int, T] {
    it := IteratorInt(rang...)
    return BiFrom(func(f1 func(int, T)) {
        t.ConsumeTillStop(func(t T) {
            i, exist := it()
            if !exist {
                panic(&Stop)
            }
            f1(i, t)
        })
    })
}

// MapBiInt 每个元素获取一个int,并转换为BiSeq
func (t Seq[T]) MapBiInt(f func(T) int) BiSeq[int, T] {
    return BiFrom(func(f1 func(int, T)) { t(func(t T) { f1(f(t), t) }) })
}

// MapBiString 每个元素获取一个String,并转换为BiSeq
func (t Seq[T]) MapBiString(f func(T) string) BiSeq[string, T] {
    return BiFrom(func(f1 func(string, T)) { t(func(t T) { f1(f(t), t) }) })
}

// MapBiAny 每个元素获取一个值,并转换为BiSeq
func (t Seq[T]) MapBiAny(f func(T) any) BiSeq[any, T] {
    return BiFrom(func(f1 func(any, T)) { t(func(t T) { f1(f(t), t) }) })
}

// MapBiAnyRight 每个元素获取一个值,并转换为BiSeq
func (t Seq[T]) MapBiAnyRight(f func(T) any) BiSeq[T, any] {
    return BiFrom(func(f1 func(T, any)) { t(func(t T) { f1(t, f(t)) }) })
}

// MapInt 每个元素转换为int
func (t Seq[T]) MapInt(f func(T) int) Seq[int] {
    return func(c func(int)) { t(func(t T) { c(f(t)) }) }
}

// MapFloat32 每个元素转换为float32
func (t Seq[T]) MapFloat32(f func(T) float32) Seq[float32] {
    return func(c func(float32)) { t(func(t T) { c(f(t)) }) }
}

// MapFloat64 每个元素转换为float64
func (t Seq[T]) MapFloat64(f func(T) float64) Seq[float64] {
    return func(c func(float64)) { t(func(t T) { c(f(t)) }) }
}

// MapSliceN 每n个元素合并为[]T,由于golang泛型问题,不能使用[]Seq[T]
func (t Seq[T]) MapSliceN(n int) Seq[[]any] {
    return t.MapSliceF(func(t T, ts []T) bool { return len(ts) == n })
}

// MapSliceF 自定义元素合并为[]T,由于golang泛型问题,不能使用[]Seq[T]
func (t Seq[T]) MapSliceF(f func(T, []T) bool) Seq[[]any] {
    return func(c func([]any)) {
        var ts []T
        t(func(t T) {
            ts = append(ts, t)
            if f(t, ts) {
                c(FromSlice(ts).Map(AnyT[T]).ToSlice())
                ts = nil
            }
        })
        if len(ts) > 0 {
            c(FromSlice(ts).Map(AnyT[T]).ToSlice())
        }
    }
}

// Join 合并多个Seq
func (t Seq[T]) Join(seqs ...Seq[T]) Seq[T] {
    return func(c func(T)) {
        t(func(t T) { c(t) })
        for _, seq := range seqs {
            seq(func(t T) { c(t) })
        }
    }
}

// Add 添加单个元素
func (t Seq[T]) Add(ts ...T) Seq[T] {
    return func(c func(T)) {
        t(func(t T) { c(t) })
        for _, e := range ts {
            c(e)
        }
    }
}

// AddF 添加单个元素
func (t Seq[T]) AddF(cast func(any) T, ts ...any) Seq[T] {
    return func(c func(T)) {
        t(func(t T) { c(t) })
        for _, e := range ts {
            c(cast(e))
        }
    }
}

// JoinF 合并Seq
func (t Seq[T]) JoinF(seq Seq[any], cast func(any) T) Seq[T] {
    return func(c func(T)) {
        t(func(t T) { c(t) })
        seq(func(t any) { c(cast(t)) })
    }
}
