---
contest: 典型90問
atcoder_url: https://atcoder.jp/contests/typical90/tasks/typical90_a
difficulty: 4
tags: [二分探索, 貪欲]
added_at: 2026-05-20
---

# 001 - Yokan Party（★4）

## 問題文

左右の長さが $L$ cm のようかんがあります。  
$N$ 個の切れ目が付けられており、左から $i$ 番目の切れ目は左から $A_i$ cm の位置にあります。

あなたは $N$ 個の切れ目のうち $K$ 個を選び、ようかんを $K+1$ 個のピースに分割したいです。

このとき、以下の値を **スコア** とします。

> $K+1$ 個のピースのうち、最も短いものの長さ

スコアが最大となるように分割する場合に得られるスコアを求めてください。

---

## 制約

- $1 \leq K \leq N \leq 100000$
- $0 < A_1 < A_2 < \cdots < A_N < L \leq 10^9$
- 入力はすべて整数

---

## 入力

```text
N L
K
A1 A2 ... AN
```

---

## 出力

求めるスコアを出力してください。

---

## 入力例 1

```text
3 34
1
8 13 26
```

## 出力例 1

```text
13
```

### 解説

**この入力の意味：** 34cm のようかんに切れ目が **3か所**（8cm・13cm・26cm）あって、その中から **1か所** を選んで切る。

全部の切り方を比べると：

- 8cm で切る → ピースは **8cm** と 26cm → いちばん短いのは **8cm**
- **13cm で切る** → ピースは **13cm** と 21cm → いちばん短いのは **13cm** ← これが最大！
- 26cm で切る → ピースは **8cm** と 26cm → いちばん短いのは **8cm**

いちばん短いピースが最も長くなる切り方は 13cm → 答えは **13**。

---

## 入力例 2

```text
7 45
2
7 11 16 20 28 34 38
```

## 出力例 2

```text
12
```

### 解説

**この入力の意味：** 45cm のようかんに切れ目が **7か所**（7・11・16・20・28・34・38cm）あって、その中から **2か所** 選んで3つに切る。

16cm と 28cm で切ると：

- 0〜16cm → **16cm**
- 16〜28cm → **12cm**
- 28〜45cm → **17cm**

いちばん短いのは **12cm** → 答えは **12**。

---

## 入力例 3

```text
3 100
1
28 54 81
```

## 出力例 3

```text
46
```

### 解説

**この入力の意味：** 100cm のようかんに切れ目が **3か所**（28cm・54cm・81cm）あって、**1か所** 選んで2つに切る。

全部の切り方：

- 28cm で切る → **28cm** と 72cm → 最小 **28cm**
- **54cm で切る** → **46cm** と 46cm → 最小 **46cm** ← これが最大！
- 81cm で切る → 19cm と **81cm** → 最小 **19cm**

真ん中に近い 54cm で切ると両方が均等になる → 答えは **46**。

---

## 入力例 4

```text
3 100
2
28 54 81
```

## 出力例 4

```text
26
```

### 解説

**この入力の意味：** 100cm のようかんに切れ目が **3か所**（28cm・54cm・81cm）あって、**2か所** 選んで3つに切る。

28cm と 54cm で切ると：

- 0〜28cm → **28cm**
- 28〜54cm → **26cm**
- 54〜100cm → **46cm**

いちばん短いのは **26cm** → 答えは **26**。

---

## 入力例 5

```text
20 1000
4
51 69 102 127 233 295 350 388 417 466 469 523 553 587 720 739 801 855 926 954
```

## 出力例 5

```text
170
```

### 解説

**この入力の意味：** 1000cm のようかんに切れ目が **20か所** あって、その中から **4か所** 選んで5つに切る。

組み合わせがとても多いので全部試すのは大変。プログラムで最良の切り方を探す。

最もうまく切ったとき、いちばん短いピースは **170cm** → 答えは **170**。

---

# 問題の解説

## この問題で何をするの？

34cm のようかんがあって、切れ目が 3 か所（8cm・13cm・26cm）に入っています。

この中から **1 か所だけ** 選んで切ります。すると 2 つのピースになります。

- 8cm で切る → 8cm と 26cm → 短い方は **8cm**
- 13cm で切る → 13cm と 21cm → 短い方は **13cm**
- 26cm で切る → 26cm と 8cm → 短い方は **8cm**

「短い方のピース」ができるだけ長くなるように切りたい。  
いちばん長くなるのは **13cm で切ったとき（答え：13）**。

## 入力の読み方

```text
3 34     ← 「切れ目が 3 か所」「ようかんは 34cm」
1        ← 「1 か所だけ選んで切る」
8 13 26  ← 切れ目の場所（左の端から何 cm か）
```

## 出力の読み方

「いちばん短いピース」の長さを 1 つ出力する。この例では `13`。

---

# やさしい解説

## イメージ

ようかんを何個かのピースに切るとき、「いちばん短いピース」をできるだけ長くしたい。

たとえば `x = 13 cm` のとき：「全部のピースを 13 cm 以上にできる？」という **YES/NO の質問** に変える。

この YES/NO の境目を **二分探索** で素早く見つける。

## どうやって切るか

左から歩いていって「もう `x` cm 歩いた！」と思ったら切る。でも最後の残りも `x` cm 以上なければダメ。

できるだけ **左で切るほど後ろに余裕ができる** から、このやり方（貪欲法）が正しい。

# 解説

## 方針

「最小ピース長」を最大化したい → **答えを決めて判定する** 方針に変える。

「全ピースを `x` cm 以上にできるか？」という判定が **できる/できない** の境目を二分探索で求める。

---

# 判定問題に変換

「すべてのピース長を `x` 以上にできるか？」

を考える。

例えば：

- `x = 10` → 可能
- `x = 20` → 不可能

なら、

```text
可能 可能 可能 不可能 不可能
```

のように単調性があるので二分探索できる。

---

# Greedy 判定

左から順番に見て：

- 前回切った位置から `x` 以上になったら切る

を繰り返す。

ただし最後のピースも `x` 以上必要。

---

# なぜ Greedy でよい？

できるだけ左で切るほど：

- 後ろに余裕ができる
- 次のピースを作りやすい

ため。

---

# 計算量

二分探索は：

```text
O(log L)
```

各判定は：

```text
O(N)
```

なので全体：

```text
O(N log L)
```

---

# C++ 解答

```cpp
#include <bits/stdc++.h>
using namespace std;
using ll = long long;  // long long を ll と短く書けるようにする（大きな数に対応）

int main() {
    int N, K;
    ll L;

    cin >> N >> L;  // 切れ目の数 N と ようかんの長さ L を読み込む
    cin >> K;       // 選ぶ切れ目の数 K を読み込む

    vector<ll> A(N);  // 切れ目の位置を N 個分用意する

    for (int i = 0; i < N; i++) {
        cin >> A[i];  // 各切れ目の位置を読み込む
    }

    // 「全ピースを x cm 以上にできるか？」を判定する関数
    // できれば true、できなければ false を返す
    auto check = [&](ll x) {
        ll prev = 0;  // 直前に切った位置（最初は左端 = 0）
        int cnt = 0;  // 実際に切った回数

        for (int i = 0; i < N; i++) {
            // A[i] - prev >= x : 今切るとこのピースが x 以上になる
            // L - A[i] >= x   : 最後のピース（右端まで）も x 以上になる
            if (A[i] - prev >= x && L - A[i] >= x) {
                cnt++;        // 切る
                prev = A[i];  // 切った位置を更新
            }
        }

        return cnt >= K;  // K 回以上切れれば OK
    };

    ll low = 0;      // low は「この値なら確実にできる」（答えの候補）
    ll high = L + 1; // high は「この値ではできない」

    // 二分探索：low と high の差が 1 になるまで繰り返す
    while (high - low > 1) {
        ll mid = (low + high) / 2;  // 真ん中の値を試す

        if (check(mid)) {
            low = mid;   // mid でできた → もっと大きくできるかも
        } else {
            high = mid;  // mid ではできない → 小さくする
        }
    }

    cout << low << endl;  // low が答え（最大スコア）
}
```

---

# Python 解答

```python
N, L = map(int, input().split())  # 切れ目の数 N、ようかんの長さ L を読み込む
K = int(input())                  # 選ぶ切れ目の数 K を読み込む
A = list(map(int, input().split()))  # 各切れ目の位置をリストで読み込む

# 「全ピースを x cm 以上にできるか？」を判定する関数
def check(x):
    prev = 0  # 直前に切った位置（最初は左端 = 0）
    cnt = 0   # 実際に切った回数

    for a in A:
        # a - prev >= x : 今切るとこのピースが x 以上
        # L - a >= x   : 残り（右端まで）も x 以上
        if a - prev >= x and L - a >= x:
            cnt += 1  # 切る
            prev = a  # 切った位置を更新

    return cnt >= K  # K 回以上切れれば OK

low = 0      # low は「この値なら確実にできる」
high = L + 1 # high は「この値ではできない」

# 二分探索：low と high の差が 1 になるまで繰り返す
while high - low > 1:
    mid = (low + high) // 2  # 真ん中の値を試す（// は整数除算）

    if check(mid):
        low = mid   # mid でできた → もっと大きくできるかも
    else:
        high = mid  # mid ではできない → 小さくする

print(low)  # low が答え（最大スコア）
```

---

# よくあるミス

## 1. 最後のピース確認忘れ

これを忘れやすい。

```cpp
L - A[i] >= x
```

---

## 2. `cnt == K` にしてしまう

正しくは：

```cpp
cnt >= K
```

---

## 3. 二分探索境界ミス

安全な持ち方：

```cpp
low = OK
high = NG
```

---

# 学びポイント

この問題で重要なのは：

- 「最小値を最大化」
- 判定問題化
- 二分探索
- Greedy

の組み合わせ。

---

# TypeScript 解答

```typescript
import * as readline from 'readline';

// 標準入力を1行ずつ受け取る（AtCoder の Node.js 典型パターン）
const rl = readline.createInterface({ input: process.stdin });
const lines: string[] = [];
rl.on('line', (line) => lines.push(line.trim()));  // 1行読むたびに配列に追加
rl.on('close', () => {
  // 全行読み終わったら処理開始
  const [N, L] = lines[0].split(' ').map(Number);  // 1行目：N と L
  const K = Number(lines[1]);                       // 2行目：K
  const A = lines[2].split(' ').map(Number);        // 3行目：切れ目の位置リスト

  // 「全ピースを x cm 以上にできるか？」を判定する関数
  // number 型の x を受け取り、boolean（true/false）を返す
  const check = (x: number): boolean => {
    let prev = 0, cnt = 0;  // prev：直前に切った位置、cnt：切った回数
    for (const a of A) {
      // このピース(a - prev)が x 以上 かつ 残り(L - a)が x 以上なら切る
      if (a - prev >= x && L - a >= x) {
        cnt++;
        prev = a;
      }
    }
    return cnt >= K;  // K 回以上切れれば OK
  };

  let low = 0, high = L + 1;  // 二分探索の範囲を初期化
  while (high - low > 1) {
    const mid = Math.floor((low + high) / 2);  // 整数で真ん中を計算
    if (check(mid)) low = mid;   // できた → low を上げる
    else high = mid;             // できない → high を下げる
  }
  console.log(low);  // 答えを出力
});
```

---

# Ruby 解答

```ruby
N, L = gets.split.map(&:to_i)  # 1行目：N と L（文字列 → 整数に変換）
K = gets.to_i                  # 2行目：K
A = gets.split.map(&:to_i)     # 3行目：切れ目の位置リスト

# 「全ピースを x cm 以上にできるか？」を判定するラムダ（無名関数）
check = ->(x) {
  prev = 0  # 直前に切った位置
  cnt = 0   # 切った回数
  A.each do |a|
    # このピース(a - prev)が x 以上 かつ 残り(L - a)が x 以上なら切る
    if a - prev >= x && L - a >= x
      cnt += 1
      prev = a
    end
  end
  cnt >= K  # Ruby では最後の式が return 値になる
}

low = 0      # low：できる側の境界
high = L + 1 # high：できない側の境界
while high - low > 1
  mid = (low + high) / 2  # 真ん中を試す
  if check.call(mid)      # ラムダは .call() で呼び出す
    low = mid   # できた → low を上げる
  else
    high = mid  # できない → high を下げる
  end
end
puts low  # 答えを出力
```

---

# PHP 解答

```php
<?php
// 1行目を読んで空白で分割し、整数に変換して N と L に代入
[$N, $L] = array_map('intval', explode(' ', trim(fgets(STDIN))));
$K = (int)trim(fgets(STDIN));                                    // 2行目：K
$A = array_map('intval', explode(' ', trim(fgets(STDIN))));      // 3行目：切れ目リスト

// 「全ピースを x cm 以上にできるか？」を判定する関数
// PHP は関数内から外の変数を直接使えないので $A, $L, $K を引数で渡す
function check($A, $L, $K, $x) {
    $prev = 0;  // 直前に切った位置
    $cnt = 0;   // 切った回数
    foreach ($A as $a) {
        // このピース($a - $prev)が x 以上 かつ 残り($L - $a)が x 以上なら切る
        if ($a - $prev >= $x && $L - $a >= $x) {
            $cnt++;
            $prev = $a;
        }
    }
    return $cnt >= $K;  // K 回以上切れれば true
}

$low = 0;       // できる側の境界
$high = $L + 1; // できない側の境界
while ($high - $low > 1) {
    $mid = intdiv($low + $high, 2);  // intdiv で整数除算（切り捨て）
    if (check($A, $L, $K, $mid)) $low = $mid;   // できた → low を上げる
    else $high = $mid;                           // できない → high を下げる
}
echo $low . "\n";  // 答えを出力
```

---

# Rust 解答

```rust
use std::io::{self, BufRead};  // 標準入力を読むためにインポート

// 「全ピースを x cm 以上にできるか？」を判定する関数
// &[i64] はスライス（配列の参照）、-> bool は返り値の型
fn check(a: &[i64], l: i64, k: usize, x: i64) -> bool {
    let mut prev = 0i64;    // 直前に切った位置（i64 型で初期化）
    let mut cnt = 0usize;   // 切った回数（usize は配列添字などに使う符号なし整数）
    for &ai in a {          // a の各要素を ai として取り出す（& でコピー）
        // このピース(ai - prev)が x 以上 かつ 残り(l - ai)が x 以上なら切る
        if ai - prev >= x && l - ai >= x {
            cnt += 1;
            prev = ai;
        }
    }
    cnt >= k  // Rust では最後の式がセミコロンなしだと return 値になる
}

fn main() {
    let stdin = io::stdin();
    let mut lines = stdin.lock().lines();  // 標準入力を行単位で読む

    // 1行目：空白区切りで N と L を読み込む
    // unwrap() はエラーがないと仮定して値を取り出す（競プロでは通常 OK）
    let line1: Vec<i64> = lines.next().unwrap().unwrap()
        .split_whitespace().map(|s| s.parse().unwrap()).collect();
    let (_n, l) = (line1[0] as usize, line1[1]);  // N は使わないので _ を付ける

    // 2行目：K を読み込む
    let k: usize = lines.next().unwrap().unwrap().trim().parse().unwrap();

    // 3行目：切れ目の位置リストを読み込む
    let a: Vec<i64> = lines.next().unwrap().unwrap()
        .split_whitespace().map(|s| s.parse().unwrap()).collect();

    let mut low = 0i64;   // できる側の境界
    let mut high = l + 1; // できない側の境界
    while high - low > 1 {
        let mid = (low + high) / 2;  // 真ん中を試す
        // Rust の if は式なので、結果を直接代入できる
        if check(&a, l, k, mid) { low = mid; } else { high = mid; }
    }
    println!("{}", low);  // 答えを出力
}
```

---

# Perl 解答

```perl
use strict;    # 変数宣言を強制（バグを防ぐ）
use warnings;  # 警告を有効化（バグを見つけやすくする）

# 1行目：空白で分割して N と L に代入
my ($N, $L) = split ' ', <STDIN>;
my $K = int(<STDIN>);        # 2行目：K（int で整数に変換）
my @A = split ' ', <STDIN>;  # 3行目：切れ目リスト（@ は配列の印）

# 「全ピースを x cm 以上にできるか？」を判定するサブルーチン
sub check {
    my ($x) = @_;           # 引数 x を受け取る（Perl の引数は @_ 経由）
    my ($prev, $cnt) = (0, 0);  # 直前に切った位置と切った回数を 0 で初期化
    for my $a (@A) {
        # このピース($a - $prev)が x 以上 かつ 残り($L - $a)が x 以上なら切る
        if ($a - $prev >= $x && $L - $a >= $x) {
            $cnt++;
            $prev = $a;
        }
    }
    return $cnt >= $K;  # K 回以上切れれば 1（true）を返す
}

my ($low, $high) = (0, $L + 1);  # 二分探索の範囲を初期化
while ($high - $low > 1) {
    my $mid = int(($low + $high) / 2);  # int で整数除算
    if (check($mid)) { $low = $mid; }   # できた → low を上げる
    else             { $high = $mid; }  # できない → high を下げる
}
print "$low\n";  # 答えを出力
```

競プロ頻出パターンなのでかなり重要。