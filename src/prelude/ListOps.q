def rec Head = fn list -> if list == [] then nil else list !! 0;

def rec _Tail = fn i, len, list -> if i == len then [] else [list !! i] ++ _Tail(i + 1, len, list);

def rec Tail = fn list -> _Tail(1, Len(list), list);

def rec Foldl = fn f, acc, list -> if list == [] then acc else Foldl(f, f(acc, Head(list)), Tail(list));

def rec Foldr = fn f, acc, list -> if list == [] then acc else f(Head(list), Foldr(f, acc, Tail(list)));

def rec Map = fn f, list -> if list == [] then [] else [f(Head(list))] ++ Map(f, Tail(list));

def rec Filter = fn f, list -> Foldl(fn e, acc -> ([e] if f(e) else []) ++ acc, [], list);

def rec Sort = fn comp, list -> let tail = Tail(list), head = Head(list) in
    Filter(tail, comp(head)) ++ [head] ++ Filter(tail, fn -> not comp(head));
