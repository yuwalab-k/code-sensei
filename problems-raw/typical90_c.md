---
contest: 典型90問
atcoder_url: https://atcoder.jp/contests/typical90/tasks/typical90_c
difficulty: 4
tags: [木, 木の直径, BFS, グラフ]
added_at: 2026-05-21
id: typical90_c
file: typical90
---

# 003 - Longest Circular Road（★4）

## 問題文

$N$ 個の都市があり、それぞれの都市に $1$ から $N$ までの番号が付けられています。
また、$N-1$ 本の道路があり、$i$ 本目 $(1 \leq i \leq N-1)$ の道路は都市 $A_i$ と都市 $B_i$ を双方向に結んでいます。
どの都市の間も、いくつかの道路を通って移動可能なものとなっています。

さて、あなたは整数 $u, v$ $(1 \leq u < v \leq N)$ を自由に選び、都市 $u$ と都市 $v$ を双方向に結ぶ道路を 1 本だけ新設することができます。そこで、以下で定められる値を **スコア** とします。

同じ道を 2 度通らずにある都市から同じ都市に戻ってくる経路における、通った道の本数（この値はただ 1 つに定まる）

スコアとして考えられる最大の値を出力してください。

## 制約

- $3 \leq N \leq 100000$
- $1 \leq A_i < B_i \leq N$ $(1 \leq i \leq N-1)$
- どの都市の間も、いくつかの道路を通って移動可能
- 与えられる入力は全て整数

## 入力

```text
N
A_1 B_1
...
A_{N-1} B_{N-1}
```

## 出力

問題文で定義されたスコアとして考えられる最大値を出力してください。

---

## 入力例 1

```text
3
1 2
2 3
```

## 出力例 1

```text
3
```

### 解説

**この入力の意味：** 3つの都市（1・2・3番）が 1−2−3 と一列につながっている。ここに道路を1本追加する。

- 都市 1 と 都市 3 の間に追加 → **1 → 2 → 3 → 1** と一周できる。道の本数 **3本** ← これが最大！
- 都市 1 と 都市 2 の間に追加 → 同じ道が2本になるだけ。一周できるが道は **2本**（微妙）
- 都市 2 と 都市 3 の間に追加 → 同上、**2本**

いちばんスコアが高くなる選び方は 都市1 と 都市3 → 答えは **3**。

---

## 入力例 2

```text
5
1 2
2 3
3 4
3 5
```

## 出力例 2

```text
4
```

### 解説

**この入力の意味：** 5つの都市、道路が 1−2−3−4 と 3−5。木の形。ここに道路を1本追加する。

- 都市 1 と 都市 4 の間 → 1→2→3→4→1、**4本** ← 最大！
- 都市 1 と 都市 5 の間 → 1→2→3→5→1、**4本** ← 同じ最大
- 都市 4 と 都市 5 の間 → 4→3→5→4、**3本**
- 都市 1 と 都市 3 の間 → 1→2→3→1、**3本**

答えは **4**。

---

## 入力例 3

```text
10
1 2
1 3
2 4
4 5
4 6
3 7
7 8
8 9
8 10
```

## 出力例 3

```text
8
```

### 解説

**この入力の意味：** 10都市の木。

木の形を整理すると：
- 1番 → 2番 → 4番 → **5番**（端）
- 1番 → 2番 → 4番 → **6番**（端）
- 1番 → 3番 → 7番 → 8番 → **9番**（端）
- 1番 → 3番 → 7番 → 8番 → **10番**（端）

いちばん遠い2都市は **5番と9番**（または5番と10番、6番と9番、6番と10番）で、その間の道の本数は **7本**。

都市5と都市9の間に道路を新設 → **5→4→2→1→3→7→8→9→5** と一周。道の本数 **8本**。

答えは **8**。

---

## 入力例 4

```text
31
1 2
1 3
2 4
2 5
3 6
3 7
4 8
4 9
5 10
5 11
6 12
6 13
7 14
7 15
8 16
8 17
9 18
9 19
10 20
10 21
11 22
11 23
12 24
12 25
13 26
13 27
14 28
14 29
15 30
15 31
```

## 出力例 4

```text
9
```

### 解説

**この入力の意味：** 31都市の完全二分木（全部で5段）。

```
              1
           /     \
         2         3
        / \       / \
       4   5     6   7
      /\ /\     /\ /\
     8 9 10 11 12 13 14 15
    /\/\/\/\/\/\/\/\/\/\/\/\/\
   16-31（葉）
```

いちばん遠い2都市は **16番と24番**（など葉と葉）で、木の上を通る道の本数は **8本**。

道路を新設すると一周で道は **9本**。

答えは **9**。

---

# 問題の解説

## この問題で何をするの？

3つの都市が 1−2−3 と一列につながっているとします（入力例1）。

ここに道路を1本だけ追加して、「ぐるっと一周できるコース」を作ります。その一周の道の本数をスコアといいます。

1番と3番をつなぐと **1→2→3→1** と3本の道を通って一周でき、スコアは **3**。

つまり「追加する道路をどう選べば、スコア（一周の長さ）が最大になるか」を求める問題です。

## 入力の読み方

```text
3       ← 都市の数が3つ
1 2     ← 1番と2番の都市をつなぐ道路
2 3     ← 2番と3番の都市をつなぐ道路
```

## 出力の読み方

この例では `3` を出力する。

---

# やさしい解説

## イメージ

木に1本だけ道を追加するとき、追加した道と既存の道でぴったり1つの「輪っか」ができます。

輪っかの長さ＝（追加した2都市を既存の道でつないだときの道の本数）＋（新しく追加した道の1本）

なので、**2都市間の既存の道の本数が最大になる2つの都市を選べばいい**ことになります。

木の上でいちばん遠い2点間の距離のことを **「木の直径」** といいます。

## どう解くか

1. 木の直径を求める（2回BFSで求められる）
2. 答え＝直径＋1

**木の直径の求め方（2回BFS）：**

1. 都市1からBFSして、いちばん遠い都市 $u$ を見つける
2. $u$ からBFSして、いちばん遠い都市までの距離を求める → それが直径

---

# 解説

## 方針

**木の直径を求めて、直径＋1を出力する。**

### なぜ「直径＋1」か

木にエッジ $(u, v)$ を1本追加すると、ちょうど1つの単純サイクルができる。そのサイクルの長さは：

$$\text{サイクル長} = \text{dist}(u, v) + 1$$

（`dist(u, v)` は木上での $u$〜$v$ 間の辺数）

スコアを最大化するには $\text{dist}(u, v)$ を最大化すればよい。木上で最大の距離＝木の直径。

### 木の直径を O(N) で求める

BFS を 2 回行う。

1. 任意の頂点 $s$ から BFS → 最遠点 $u$ を求める
2. $u$ から BFS → 最遠点 $v$ までの距離 $d$ を求める
3. $d$ が木の直径

---

# よくあるミス

## 1. 答えを「直径」そのものにしてしまう

新しく1本道路を追加するので、スコアは **直径＋1**。`+1` を忘れずに。

---

## 2. BFSを1回しかしない

1回目のBFS（任意の点から）だけでは直径の端点を正しく求められない場合がある。2回必要。

---

## 3. グラフが大きいときにスタック溢れ（DFSの再帰）

$N \leq 100000$ なのでDFSの再帰が深くなりスタックオーバーフローになることがある。BFS（キュー）か、再帰制限を上げる必要がある。

---

# 学びポイント

- **木の直径** の典型的な求め方（2回BFS）
- 「木に辺を1本追加するとサイクルができる」という性質の活用
- 計算量 O(N) で解ける典型問題

---

# Python 解答

```python
import sys
from collections import deque
input = sys.stdin.readline  # 高速入力に切り替え

def bfs(start, graph, n):
    dist = [-1] * (n + 1)   # 各都市までの距離。-1 は未訪問
    dist[start] = 0          # スタート地点は距離 0
    q = deque([start])       # BFS 用キュー
    farthest = start         # 最遠点（初期値はスタート地点）
    max_dist = 0             # 最遠距離

    while q:
        v = q.popleft()      # キューの先頭から都市を取り出す
        for u in graph[v]:   # 隣接する都市を全部調べる
            if dist[u] == -1:          # まだ訪れていない都市なら
                dist[u] = dist[v] + 1  # 距離を 1 増やして記録
                q.append(u)            # キューに追加
                if dist[u] > max_dist: # より遠い都市を見つけたら更新
                    max_dist = dist[u]
                    farthest = u

    return farthest, max_dist  # 最遠点とその距離を返す

N = int(input())                            # 都市の数
graph = [[] for _ in range(N + 1)]         # 隣接リスト（1-indexed）

for _ in range(N - 1):
    a, b = map(int, input().split())        # 道路の両端を読む
    graph[a].append(b)                      # a→b の道
    graph[b].append(a)                      # b→a の道（双方向）

# 1回目の BFS: 都市 1 から最も遠い都市を見つける
u, _ = bfs(1, graph, N)

# 2回目の BFS: 最遠点 u からさらに最も遠い都市を見つける（＝直径）
_, diameter = bfs(u, graph, N)

print(diameter + 1)  # 直径 + 1 がスコアの最大値
```

---

# C++ 解答

```cpp
#include <bits/stdc++.h>
using namespace std;

// BFS で start から各都市への距離を求め、最遠点と距離を返す
pair<int,int> bfs(int start, vector<vector<int>>& graph, int n) {
    vector<int> dist(n + 1, -1);  // 距離配列（-1 は未訪問）
    dist[start] = 0;              // スタート地点の距離は 0
    queue<int> q;
    q.push(start);                // キューにスタート地点を追加
    int farthest = start;         // 最遠点
    int max_dist = 0;             // 最遠距離

    while (!q.empty()) {
        int v = q.front(); q.pop();           // キューの先頭を取り出す
        for (int u : graph[v]) {              // 隣接都市を調べる
            if (dist[u] == -1) {             // 未訪問なら
                dist[u] = dist[v] + 1;       // 距離を記録
                q.push(u);
                if (dist[u] > max_dist) {    // より遠ければ更新
                    max_dist = dist[u];
                    farthest = u;
                }
            }
        }
    }
    return {farthest, max_dist};  // 最遠点と距離を返す
}

int main() {
    ios::sync_with_stdio(false);  // 高速入力
    cin.tie(nullptr);

    int N;
    cin >> N;                     // 都市の数

    vector<vector<int>> graph(N + 1);  // 隣接リスト（1-indexed）

    for (int i = 0; i < N - 1; i++) {
        int a, b;
        cin >> a >> b;            // 道路の両端
        graph[a].push_back(b);   // 双方向に追加
        graph[b].push_back(a);
    }

    // 1回目の BFS: 都市 1 から最遠点 u を探す
    auto [u, d1] = bfs(1, graph, N);

    // 2回目の BFS: u から最遠点までの距離 = 木の直径
    auto [v, diameter] = bfs(u, graph, N);

    cout << diameter + 1 << "\n";  // 直径 + 1 が答え
}
```

---

# TypeScript 解答

```typescript
import * as readline from 'readline';

const rl = readline.createInterface({ input: process.stdin });
const lines: string[] = [];
rl.on('line', (line) => lines.push(line.trim()));
rl.on('close', () => {
    const N = parseInt(lines[0]);                              // 都市の数
    const graph: number[][] = Array.from({length: N + 1}, () => []);  // 隣接リスト

    for (let i = 1; i < N; i++) {
        const [a, b] = lines[i].split(' ').map(Number);       // 道路の両端
        graph[a].push(b);                                      // 双方向に追加
        graph[b].push(a);
    }

    // BFS で start から最遠点と距離を返す
    function bfs(start: number): [number, number] {
        const dist = new Array(N + 1).fill(-1);  // 距離配列
        dist[start] = 0;
        const q: number[] = [start];             // キュー（配列で代用）
        let qi = 0;                              // キューの読み取り位置
        let farthest = start;
        let maxDist = 0;

        while (qi < q.length) {
            const v = q[qi++];                   // キューの先頭を取り出す
            for (const u of graph[v]) {          // 隣接都市を調べる
                if (dist[u] === -1) {            // 未訪問なら
                    dist[u] = dist[v] + 1;       // 距離を記録
                    q.push(u);
                    if (dist[u] > maxDist) {     // より遠ければ更新
                        maxDist = dist[u];
                        farthest = u;
                    }
                }
            }
        }
        return [farthest, maxDist];
    }

    const [u] = bfs(1);              // 1回目: 都市 1 から最遠点 u を探す
    const [, diameter] = bfs(u);     // 2回目: u から直径を求める
    console.log(diameter + 1);       // 直径 + 1 が答え
});
```

---

# Ruby 解答

```ruby
N = gets.to_i                              # 都市の数
graph = Array.new(N + 1) { [] }           # 隣接リスト（1-indexed）

(N - 1).times do
    a, b = gets.split.map(&:to_i)         # 道路の両端
    graph[a] << b                          # 双方向に追加
    graph[b] << a
end

# BFS で start から最遠点と距離を返す
def bfs(start, graph, n)
    dist = Array.new(n + 1, -1)           # 距離配列（-1 は未訪問）
    dist[start] = 0                        # スタート地点は 0
    q = [start]                            # BFS キュー
    farthest = start
    max_dist = 0

    until q.empty?
        v = q.shift                        # キューの先頭を取り出す
        graph[v].each do |u|
            if dist[u] == -1               # 未訪問なら
                dist[u] = dist[v] + 1      # 距離を記録
                q << u
                if dist[u] > max_dist      # より遠ければ更新
                    max_dist = dist[u]
                    farthest = u
                end
            end
        end
    end
    [farthest, max_dist]
end

u, _ = bfs(1, graph, N)           # 1回目: 都市 1 から最遠点 u を探す
_, diameter = bfs(u, graph, N)    # 2回目: u から直径を求める
puts diameter + 1                  # 直径 + 1 が答え
```

---

# PHP 解答

```php
<?php
$N = (int)trim(fgets(STDIN));             // 都市の数
$graph = array_fill(0, $N + 1, []);       // 隣接リスト（1-indexed）

for ($i = 0; $i < $N - 1; $i++) {
    [$a, $b] = array_map('intval', explode(' ', trim(fgets(STDIN))));  // 道路の両端
    $graph[$a][] = $b;                    // 双方向に追加
    $graph[$b][] = $a;
}

// BFS で start から最遠点と距離を返す
function bfs($start, $graph, $n) {
    $dist = array_fill(0, $n + 1, -1);   // 距離配列（-1 は未訪問）
    $dist[$start] = 0;                    // スタート地点は 0
    $q = [$start];                        // BFS キュー
    $head = 0;                            // キューの読み取り位置
    $farthest = $start;
    $maxDist = 0;

    while ($head < count($q)) {
        $v = $q[$head++];                 // キューの先頭を取り出す
        foreach ($graph[$v] as $u) {
            if ($dist[$u] === -1) {       // 未訪問なら
                $dist[$u] = $dist[$v] + 1; // 距離を記録
                $q[] = $u;
                if ($dist[$u] > $maxDist) { // より遠ければ更新
                    $maxDist = $dist[$u];
                    $farthest = $u;
                }
            }
        }
    }
    return [$farthest, $maxDist];
}

[$u] = bfs(1, $graph, $N);              // 1回目: 都市 1 から最遠点 u を探す
[, $diameter] = bfs($u, $graph, $N);   // 2回目: u から直径を求める
echo $diameter + 1 . "\n";             // 直径 + 1 が答え
```

---

# Rust 解答

```rust
use std::io::{self, BufRead, Write};
use std::collections::VecDeque;

// BFS で start から最遠点と距離を返す
fn bfs(start: usize, graph: &Vec<Vec<usize>>, n: usize) -> (usize, usize) {
    let mut dist = vec![usize::MAX; n + 1];  // 距離配列（MAX は未訪問）
    dist[start] = 0;                          // スタート地点は 0
    let mut q = VecDeque::new();
    q.push_back(start);
    let mut farthest = start;
    let mut max_dist = 0;

    while let Some(v) = q.pop_front() {      // キューの先頭を取り出す
        for &u in &graph[v] {               // 隣接都市を調べる
            if dist[u] == usize::MAX {      // 未訪問なら
                dist[u] = dist[v] + 1;     // 距離を記録
                q.push_back(u);
                if dist[u] > max_dist {    // より遠ければ更新
                    max_dist = dist[u];
                    farthest = u;
                }
            }
        }
    }
    (farthest, max_dist)
}

fn main() {
    let stdin = io::stdin();
    let stdout = io::stdout();
    let mut out = io::BufWriter::new(stdout.lock());  // 出力バッファ

    let mut lines = stdin.lock().lines();
    let n: usize = lines.next().unwrap().unwrap().trim().parse().unwrap();  // 都市の数

    let mut graph = vec![vec![]; n + 1];  // 隣接リスト（1-indexed）

    for _ in 0..n - 1 {
        let line = lines.next().unwrap().unwrap();
        let mut iter = line.split_whitespace();
        let a: usize = iter.next().unwrap().parse().unwrap();  // 道路の一端
        let b: usize = iter.next().unwrap().parse().unwrap();  // 道路の他端
        graph[a].push(b);  // 双方向に追加
        graph[b].push(a);
    }

    let (u, _) = bfs(1, &graph, n);            // 1回目: 都市 1 から最遠点 u を探す
    let (_, diameter) = bfs(u, &graph, n);     // 2回目: u から直径を求める
    writeln!(out, "{}", diameter + 1).unwrap(); // 直径 + 1 が答え
}
```

---

# Perl 解答

```perl
use strict;
use warnings;

my $N = int(<STDIN>);                       # 都市の数
my @graph = map { [] } (0..$N);             # 隣接リスト（1-indexed）

for (1..$N-1) {
    my ($a, $b) = split ' ', <STDIN>;       # 道路の両端
    push @{$graph[$a]}, $b;                 # 双方向に追加
    push @{$graph[$b]}, $a;
}

# BFS で start から最遠点と距離を返す
sub bfs {
    my ($start) = @_;
    my @dist = (-1) x ($N + 1);            # 距離配列（-1 は未訪問）
    $dist[$start] = 0;                      # スタート地点は 0
    my @q = ($start);                       # BFS キュー
    my ($farthest, $max_dist) = ($start, 0);

    while (@q) {
        my $v = shift @q;                   # キューの先頭を取り出す
        for my $u (@{$graph[$v]}) {         # 隣接都市を調べる
            if ($dist[$u] == -1) {          # 未訪問なら
                $dist[$u] = $dist[$v] + 1;  # 距離を記録
                push @q, $u;
                if ($dist[$u] > $max_dist) {# より遠ければ更新
                    $max_dist = $dist[$u];
                    $farthest = $u;
                }
            }
        }
    }
    return ($farthest, $max_dist);
}

my ($u) = bfs(1);                          # 1回目: 都市 1 から最遠点 u を探す
my (undef, $diameter) = bfs($u);           # 2回目: u から直径を求める
print $diameter + 1, "\n";                 # 直径 + 1 が答え
```
