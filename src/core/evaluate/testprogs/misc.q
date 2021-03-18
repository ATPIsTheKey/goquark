let Mul = fn a, b -> a * b in
    let Double = Mul(2) in Double(6);

let rec Fib = fn n -> if n < 2 then n else Fib(n - 1) + Fib(n - 2) in Fib(8);

let Y = fn f -> (fn x -> x(x))(fn y -> f(y(y))) in
    let NonRecFib fn self, n -> if n < 2 then n else self(n - 1) + self(n - 2) in
        let Fib = fn n -> Y(NonRecFib)(n) in Fib(8);

let rec Ack = fn m, n -> if m == 0 then n + 1 else if n == 0 then Ack(m - 1, 1) else Ack(m - 1, Ack(m, n - 1)) in Ack(2, 2);

let rec Fac = fn n -> if n == 1 then 1 else Fac(n - 1) * n in Fac(4);

let Y = fn f -> (fn x -> x(x))(fn y -> f(y(y))) in
    let NonRecFac = fn self, n -> if n == 1 then 1 else self(n - 1) * n in
        let Fac = fn n -> Y(NonRecFac)(n) in Fac(4);

def Fac4 = (fn f -> (fn x -> f(x(x)))(fn y -> f(y(y))))(fn g -> fn n -> if n == 1 then 1 else g(n - 1) * n)(4);

let rec Gcd = fn a, b -> if b == 0 then a else Gcd(b, a % b) in Gcd(56, 88);

let Y = fn f -> (fn x -> x(x))(fn y -> f(y(y))) in
    let NonRecGcd = fn self, a, b -> if b == 0 then a else self(b, a % b) in
        let Gcd = fn a, b -> Y(Gcd)(a, b) in Gcd(56, 88);
