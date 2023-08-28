package seq

//========转换为 BiSeq ================

func (t Seq[T]) MapFlatBi2AnyAny(f func(T) BiSeq[any, any]) BiSeq[any, any] {
    return BiFrom(func(f1 func(any, any)) { t(func(t T) { f(t)(func(k, v any) { f1(k, v) }) }) })
}
func (t Seq[T]) MapFlatBi2StringString(f func(T) BiSeq[string, string]) BiSeq[string, string] {
    return BiFrom(func(f1 func(string, string)) { t(func(t T) { f(t)(func(k, v string) { f1(k, v) }) }) })
}

func (t Seq[T]) MapFlatBi2StringAny(f func(T) BiSeq[string, any]) BiSeq[string, any] {
    return BiFrom(func(f1 func(string, any)) { t(func(t T) { f(t)(func(k string, v any) { f1(k, v) }) }) })
}

func (t Seq[T]) MapFlatBi2IntAny(f func(T) BiSeq[int, any]) BiSeq[int, any] {
    return BiFrom(func(f1 func(int, any)) { t(func(t T) { f(t)(func(k int, v any) { f1(k, v) }) }) })
}
