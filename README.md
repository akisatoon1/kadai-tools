# kadai-tools - 課題をやるときの作業を楽にする

## ツールの概要
某大学のCプログラミング演習をやるときの定型作業を、このツールのコマンドで代用できます。具体的には、
  
- 提出用ファイルを作る
- プログラムをコンパイルして実行
- 課題を提出
  
のような作業を、ツールのコマンドを叩くだけで完了します。

## 準備
実行ファイルと同じディレクトリに、`tmpl.c`と`login.json`があります。無い場合は追加してください。それぞれのファイルに役割があります。

### tmpl.c
```c:tmpl.c
// kadai{{.Num}}{{.Level}}

#include <stdio.h>

int main (void) {
    return 0;
}
```
C言語のテンプレートをこのファイルに書いてください。
  
`kadai01a`の`01`の部分を埋め込みたいときは`{{.Num}}`、`a`の部分を埋め込みたいときは`{{.Level}}`を使用してください。

### settings.json
```json:settings.json
{
    "username": "your student ID",
    "password": "your password",
    "lang": "c or c++"
}
```
`username`と`password`はmanabaに提出するときに使います。総合認証に使う学籍番号とパスワードをそれぞれの`""`の中に書いてください。
  
`lang`はファイル生成時の拡張子や`kadai debug`時に使うコンパイラを決めます。現在`c`と`c++`が使用でき、`c`のときはgcc、`c++`のときはg++を使用します。
## 使い方

### make
```
~$ kadai make 01 a x
./01/ 作成
./tmpl.c 読み込み
./01/kadai01a.c 作成
./01/kadai01a.c テンプレート書き込み
./01/inputFiles/inputa1.txt 作成
./01/kadai01x.c 作成
./01/kadai01x.c テンプレート書き込み
./01/inputFiles/inputx1.txt 作成
完了
```
```
~$ ls -R 01
01:
inputFiles  kadai01a.c  kadai01x.c

01/inputFiles:
inputa1.txt  inputx1.txt
```
最初の引数はディレクトリ名、残りの引数はそれぞれの課題を表します。
  
それぞれ作成した`*.c`ファイルに`tmpl.c`で設定したテンプレートを書き込みます。
### debug
```
~/01$ kadai debug a
compile: gcc -Wall -o kadai01a kadai01a.c -lm
Output of inputFiles/inputa1.txt:
Hello World!
```
引数はデバッグする課題を指定します。標準入力には`./inputFiles/inputa1.txt`の中身を使います。
  
`inputa2.txt` `inputa3.txt`のようなファイルを自身で作成すれば、それぞれのファイルに対してプログラムを実行します。
### submit
```
~/01$ kadai submit a
```
引数はmanabaに提出する課題を指定します。
  
```
~/01$ kadai submit -r a
```
rオプションで再提出します。