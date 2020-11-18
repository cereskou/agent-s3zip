# agent-s3zip
## AWS Lambda function
this function will zip s3 file and upload to another bucket.<br>

### Create AWS Lambda Function
![create](https://github.com/cereskou/agent-s3zip/blob/main/images/lambda-create.png)
1. Create a lambda function named S3ServiceTar.<br>
1. Upload the main.zip(built by gobuild.cmd).
1. change the function setup. 
![setup](https://github.com/cereskou/agent-s3zip/blob/main/images/lambda-set.png)

### Setup S3 bucket event (sorry for japanese)
1. AWSコンソールページに、「S3バケット」を選択します。
1. 「イベント」ボタンをクリックします。
1. 「通知の追加」をクリックして、入力画面を表示します。
1. 名前が「S3ServiceTar」を入力します。（名前が任意）
1. イベントが「PUT」をチェックします。
1. プレフィックスが「作業用パス」を入力します。（ファイル一覧.jsonxファイルのアップロード先）
1. サフィックスが .jsonx を入力します。
1. 送信先が「Lambda関数」を選択します。
1. Lambda関数 S3ServiceTarを選択し、保存します。
![s3bucket](https://github.com/cereskou/agent-s3zip/blob/main/images/s3bucket.png)

### Function flow
![flow](https://github.com/cereskou/agent-s3zip/blob/main/images/zip-flow.png)
