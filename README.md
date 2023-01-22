# PCCISO
内部で自動的に最新のISOファイルを再配布するためのソフトウェアです。

毎日特定の時間に以下のタスクを実行します。
- jsonの保存されたデータと実ファイルの整合性チェック
- 各ページを巡回し、最新ファイルがローカルにダウンロードされているか確認
- 最新の状態に基づいてHTMLページを再生成

## run
`docker compose up -d`

## config
### debug
デバッグログです。  
ダウンロード状況などが表示されます。  
主に時間のかかる処理にログが追加されるので、何も出なくて動いているか不安になるときは入れてみるといいかもしれません。  

### data_dir
データディレクトリを指定します。
このディレクトリをそのまま他のサーバーで静的ホストできます。

### webhook
Discord Webhook URLを指定します。
タスクの完了後に結果が送信されます。
エラーに気づくために入れておくことをおすすめします。
指定されなかったり空文字列であればWebhook機能は無効になります。

### strict
trueの場合、ファイルの整合性チェックでチェックサムを検証します。
falseの場合はファイルの存在とサイズのみ確認します。

## run_time
定期タスクを実行する時刻を指定します。
"12:34"の形式で指定してください。
**Dockerを使用している場合、おそらくUTCです。**

### os
id: URLなどに使われる一意の識別子で、被ってはいけません。  
name: 表示名です。  
color: ページなどの基軸になる色です。  
category: trueだとその1つ下までの子のファイルを含んで表示するようになります。  
#### update
最新ファイルに辿り着く方法を記述します。  
指定しなかった場合はupdate処理が行われません。
url: 指定されたURLをそのまま次のURLとして設定(最初は必ずこれなはず)  
regex: 一つ前に設定されたURLを開き、その中からregexに一致する部分を探し、第一グループを次のURLとして設定  
regex_select: regexとセットで用い、一致したうち何番目を使用するか指定する。negative indexが使用可能。省略した場合0と扱われる。  
最後に次のURLとして指定されたものがダウンロードURLとして扱われます。  

#### update_sum
updateによってたどり着いたファイルのチェックサムの入手方法を記述します。  
省略した場合チェックサムは用いられません。  
updateと同じ内容に加え、以下の内容が追加されます。  
url: updateと同じですが、$FILEURLはupdateで決定したURLに、$FILENAMEはそのURLのファイル名部分に置換されます。  
sum_regex: sum_typeとセットで用います。regexと同じような方法でチェックサムファイルの中のどこがチェックサムか指定します。urlと同じ置換が行われます。  
sum_type: sum_regexとセットで用います。チェックサムの種類を指定します。現在有効なのは`sha256`と`sha512`のみです。  

#### download
ユーザーがダウンロードを要求するときに許容するURLの正規表現  
...だったのですが、ユーザーからの要求でダウンロードする機能は未実装のため現状は無意味です。  

#### child
osのlistを指定できます。  
ここに指定したものは子の扱いになります。  



## Credit
Copy icon is from https://iconscout.com/icon/copy-duplicate-clone-document-zerox-print