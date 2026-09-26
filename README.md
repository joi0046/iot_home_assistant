# iot_home_assistant

ラズパイを使用したIoTゲートウェイシステム。

センサーやデバイスの接続・認識・データ収集・可視化・制御をできるだけ自動化して、最終的にローカルAIによる判断・制御までを行うことを目指す。

技術スタックとして <br>
    1. Golang <br>
    2. home assistant <br>
    3. InfluxDB <br>
    4. Grafana <br>
    5. Hermes Agent <br>
    6. Bonsai <br>

ゴールとして <br>
    1. センサーの自動認識 <br>
    2. センサーDriverの共通化 <br>
    3. MQTTによる疎結合なシステム <br>
    4. Home Assistantへの自動登録 <br>
    5. センサーデータの長期保存・可視化 <br>
    6. ESP32との連携 <br>
    7. ローカルAIによるIoT制御 <br>

テストを実行

```bash
go test ./...
```

```bash
go test ./internal/protocol
```

docker composeを使用してコンテナを起動
 
```bash
cd iot_home_assistant
docker compose -f docker/compose.yml up -d
```

ゲートウェイを起動

```bash
cd app/gateway/
go run cmd/gateway/main.go
```
