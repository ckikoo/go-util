package seq

import (
    "fmt"
    "math/rand"
    "testing"
    "time"
)

func init() {
    rand.Seed(time.Now().UnixNano())
}
func preTest(t *testing.T) {
    t.Parallel()
}

func Test1(t *testing.T) {
    preTest(t)
    seq := FromSlice([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10})
    ok1 := 0
    ok2 := 0
    ok3 := 0
    ok4 := 1
    CastAnyT(
        seq.OnEach(func(i int) {
            ok1++
        }).Take(50).Filter(func(i int) bool {
            return i%2 == 0
        }).OnEach(func(i int) {
            ok2++
        }).MapFlat(func(i int) Seq[any] {
            return FromSlice([]any{i, i + 1})
        }), 0,
    ).ForEach(func(i int) {
        ok3++
        ok4++
        if ok4 != i {
            t.Fail()
        }
    })
    if ok1 != 10 || ok2 != 5 || ok3 != 10 {
        t.Fail()
    }
    ok4 = 0
    seq.ForEach(func(i int) {
        ok4++
        if ok4 != i {
            t.Fail()
        }
    })
}

func TestFromIntSeq(t *testing.T) {
    preTest(t)
    seq := FromIntSeq(1, 10)
    ok := 0
    seq.ForEach(func(i int) {
        ok++
    })
    if ok != 10 {
        t.Fail()
    }
}

func TestTake(t *testing.T) {
    preTest(t)
    seq := FromIntSeq(0, 9)
    var r []int
    seq.Take(5).ForEach(func(i int) { r = append(r, i) })
    if len(r) != 5 {
        t.Fail()
    }
    for i := 0; i < 5; i++ {
        if r[i] != i {
            t.Fail()
        }
    }
}

func TestDrop(t *testing.T) {
    preTest(t)
    seq := FromIntSeq(0, 9)
    var r []int
    seq.Drop(5).ForEach(func(i int) { r = append(r, i) })
    if len(r) != 5 {
        t.Fail()
    }
    for i := 0; i < 5; i++ {
        if r[i] != i+5 {
            t.Fail()
        }
    }
}

func TestDropTake(t *testing.T) {
    preTest(t)
    seq := FromIntSeq()
    var r []int
    seq.Drop(5).Take(5).ForEach(func(i int) { r = append(r, i) })
    if len(r) != 5 {
        t.Fail()
    }
    for i := 0; i < 5; i++ {
        if r[i] != i+5 {
            t.Fail()
        }
    }
}

func TestRand(t *testing.T) {
    preTest(t)
    slice := FromRandIntSeq().OnEach(func(i int) {
        fmt.Println("", i)
    }).Filter(func(i int) bool {
        return i%2 == 0
    }).Drop(10).Take(5).ToSlice()
    if len(slice) != 5 {
        t.Fail()
    }
    fmt.Println(slice)
}

func TestSeq_Complete(t *testing.T) {
    preTest(t)
    s := FromIntSeq().Take(1000).MapBiSerialNumber(100).Cache()
    {
        it := IteratorInt()
        s.SeqV().ForEach(func(i int) {
            i2, _ := it()
            if i != i2 {
                t.Fail()
            }
        })
    }
    {
        it := IteratorInt(100)
        s.SeqK().ForEach(func(i int) {
            i2, _ := it()
            if i != i2 {
                t.Fail()
            }
        })
    }
}
