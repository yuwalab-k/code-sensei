---
contest: 典型90問
atcoder_url: https://atcoder.jp/contests/typical90/tasks/typical90_a
difficulty: 4
tags: [二分探索, 貪欲]
added_at: 2026-05-20
---

# 001 - Yokan Party（★4）

## 問題文

左右の長さが `L` cm のようかんがあります。  
`N` 個の切れ目が付けられており、左から `i` 番目の切れ目は左から `Ai` cm の位置にあります。

あなたは `N` 個の切れ目のうち `K` 個を選び、ようかんを `K+1` 個のピースに分割したいです。

このとき、以下の値を **スコア** とします。

> `K+1` 個のピースのうち、最も短いものの長さ

スコアが最大となるように分割する場合に得られるスコアを求めてください。

---

## 制約

- `1 ≤ K ≤ N ≤ 100000`
- `0 < A1 < A2 < ... < AN < L ≤ 10^9`
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

左から 2 番目の切れ目で分割すると：

- 13 cm
- 21 cm

となり、最小値は `13`。

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

---

# 解説

## 方針

この問題は：

> 「最小ピース長」を最大化する

問題。

このタイプは典型的な：

- 最大値の最小化
- 最小値の最大化

であり、**二分探索** を使う。

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
using ll = long long;

int main() {
    int N, K;
    ll L;

    cin >> N >> L;
    cin >> K;

    vector<ll> A(N);

    for (int i = 0; i < N; i++) {
        cin >> A[i];
    }

    auto check = [&](ll x) {
        ll prev = 0;
        int cnt = 0;

        for (int i = 0; i < N; i++) {
            // 今切ったとき
            // 前のピース >= x
            // 最後のピース >= x
            if (A[i] - prev >= x && L - A[i] >= x) {
                cnt++;
                prev = A[i];
            }
        }

        return cnt >= K;
    };

    ll low = 0;
    ll high = L + 1;

    while (high - low > 1) {
        ll mid = (low + high) / 2;

        if (check(mid)) {
            low = mid;
        } else {
            high = mid;
        }
    }

    cout << low << endl;
}
```

---

# Python 解答

```python
N, L = map(int, input().split())
K = int(input())
A = list(map(int, input().split()))

def check(x):
    prev = 0
    cnt = 0

    for a in A:
        if a - prev >= x and L - a >= x:
            cnt += 1
            prev = a

    return cnt >= K

low = 0
high = L + 1

while high - low > 1:
    mid = (low + high) // 2

    if check(mid):
        low = mid
    else:
        high = mid

print(low)
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

競プロ頻出パターンなのでかなり重要。