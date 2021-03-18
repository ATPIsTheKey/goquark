@include math

def Double = let Mul = fn a, b -> a * b in Mul(2);

def rec Fac = fn n -> if n == 1 then 1 else Fac(n - 1) * n in Fac(4);

def Fac' = let Y = fn f -> (fn x -> x(x))(fn y -> f(y(y))) in
    let NonRecFac = fn self, n -> if n == 1 then 1 else self(n - 1) * n in
        fn n -> Y(NonRecFac)(n);

def rec Fib = fn n -> if n < 2 then n else Fib(n - 1) + Fib(n - 2);

def Fib' = let Y = fn f -> (fn x -> x(x))(fn y -> f(y(y))) in
    let NonRecFib fn self, n -> if n < 2 then n else self(n - 1) + self(n - 2) in
        fn n -> Y(NonRecFib)(n),

def rec Gcd = fn a, b -> if b == 0 then a else Gcd(b, a % b);

def Gcd' = let Y = fn f -> (fn x -> x(x))(fn y -> f(y(y))) in
    let NonRecGcd = fn self, a, b -> if b == 0 then a else self(b, a % b) in
        fn a, b -> Y(NonRecGcd)(a, b);

def rec Ack = fn m, n -> if m == 0 then n + 1 else if n == 0 then Ack(m - 1, 1) else Ack(m - 1, Ack(m, n - 1))

def Ack' = let Y = fn f -> (fn x -> x(x))(fn y -> f(y(y))) in
    let NonRecAck = fn self, m, n -> if m == 0 then n + 1 else if n == 0 then self(m - 1, 1) else self(m - 1, self(m, n - 1)) in
        fn m, n -> Y(NonRecAck)(m, n);
