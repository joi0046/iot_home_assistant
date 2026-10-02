# iot_home_assistant

[日本語](#概要) | [English](#english-overview)

## English Overview

An IoT gateway system running on Raspberry Pi. It auto-detects I²C sensors,
collects data, and publishes it over MQTT so that Home Assistant, InfluxDB /
Grafana, and other systems can consume it in a loosely coupled way.

- Go gateway: I²C discovery → driver registry → MQTT publish (`gateway/v1/readings`)
- Supported sensors: HDC1000 (temperature / humidity), MPU6050 (acceleration / gyroscope / temperature)
- Stack: Go, Home Assistant, InfluxDB, Grafana, Mosquitto (MQTT)
- See [docs/architecture.md](docs/architecture.md) and [docs/protocol/mqtt-v1.md](docs/protocol/mqtt-v1.md) for details.

## 概要

ラズパイを使用したIoTゲートウェイシステム。

センサーやデバイスの接続・認識・データ収集・可視化・制御をできるだけ自動化して、最終的にローカルAIによる判断・制御までを行うことを目指す。

ゴールとして <br>
    1. センサーの自動認識 <br>
    2. センサーDriverの共通化 <br>
    3. MQTTによる疎結合なシステム <br>
    4. Home Assistantへの自動登録 <br>
    5. センサーデータの長期保存・可視化 <br>
    6. ESP32との連携 <br>

技術スタックとして <br>
    1. Golang <br>
    2. Home Assistant <br>
    3. InfluxDB <br>
    4. Grafana <br>
    5. Mosquitto (MQTT Broker) <br>

詳しくは [docs/architecture.md](docs/architecture.md) を参照。

## 対応センサー

| センサー | 取得データ | 状態 |
| ------- | --------- | ---- |
| HDC1000 | 温度・湿度 | 対応済み |
| MPU6050 | 加速度・ジャイロ・温度 | 対応済み |
| BME280 | 温度・湿度・気圧 | 予定 (Driverはスタブのみ) |
| BH1750 | 照度 | 予定 (未実装) |

MQTTペイロード形式は [docs/protocol/mqtt-v1.md](docs/protocol/mqtt-v1.md) を参照。

## 前提条件

- Raspberry Pi (I²Cが有効化されていること)
- Docker / Docker Compose
- Go 1.27以上 (ゲートウェイを直接実行する場合のみ)

## 起動手順

### 1. 環境変数の設定

```bash
cp .env.example .env
```

必要に応じて `.env` の値を編集する (InfluxDBの認証情報など)。

### 2. 周辺サービスの起動

```bash
docker compose --env-file .env -f docker/compose.yml up -d --build
```



## サービス一覧

docker composeで起動する主なサービス:

| サービス | URL / ポート | 用途 |
| ------- | ------------ | ---- |
| Home Assistant | http://\<raspberry-pi\>:8123 | 状態表示・操作 |
| Grafana | http://\<raspberry-pi\>:3000 | データ可視化 |
| InfluxDB | http://\<raspberry-pi\>:8086 | 時系列データ保存 |
| Mosquitto | \<raspberry-pi\>:1883 | MQTT Broker |

注意: `services/mqtt/config/mosquitto.conf` は `allow_anonymous true`
(ローカル実験用設定) になっている。公開ネットワークでは使用しないこと。

## テストを実行

```bash
cd app/gateway/
go test ./...
```

```bash
go test ./internal/protocol
```

## ライセンス

MIT License ([LICENSE](LICENSE))
