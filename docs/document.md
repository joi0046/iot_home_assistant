# Raspberry Pi IoT Gateway

## 要件定義書

**Version:** 1.0
**作成日:** 2026年9月25日
**対象:** Raspberry Pi 4 IoT Gatewayプロジェクト

---

# 1. プロジェクト概要

## 1.1 目的

Raspberry Piを中心としたIoT Gatewayを構築し、接続されたセンサーやデバイスを自動的に認識・管理し、取得したデータをMQTT経由で外部システムへ提供する。

最終的には、

> 接続 → 認識 → 登録 → 取得 → 可視化 → 判断 → 制御

までを可能な限り自動化することを目標とする。

---

# 2. システムの目的

本システムでは、個々のセンサーやデバイスごとに個別のプログラムを書くのではなく、Gateway側でデバイスを抽象化して扱える構成を目指す。

主な目的は以下の通り。

* I²C接続デバイスの自動検出
* センサー種類の自動判定
* センサードライバの統一インターフェース化
* センサー情報の一元管理
* MQTTによるデータ配信
* Home Assistantとの連携
* 将来的なBLE/GPIOデバイスへの対応
* デバイス追加時のGateway本体への変更を最小化

---

# 3. 対象ハードウェア

## 3.1 Gateway

| 項目        | 内容             |
| --------- | -------------- |
| コンピュータ    | Raspberry Pi 4 |
| メモリ       | 8GB            |
| OS        | Linux          |
| Gateway言語 | Go             |
| 通信        | I²C / MQTT     |
| コンテナ      | Docker         |

---

## 3.2 I²Cセンサー

MVPでは以下のセンサーを対象とする。

### BME280

取得対象:

* 温度
* 湿度
* 気圧

I²Cアドレス:

* `0x76`
* `0x77`

### BH1750

取得対象:

* 照度

---

# 4. ソフトウェア構成

システムは以下のコンポーネントで構成する。

```text
Raspberry Pi
│
├── Go IoT Gateway
│   │
│   ├── Discovery
│   │
│   ├── Driver Registry
│   │
│   ├── Device Manager
│   │
│   └── Sensor Drivers
│       ├── BME280
│       └── BH1750
│
├── MQTT Broker
│   └── Mosquitto
│
├── Home Assistant
│
├── InfluxDB
│
└── Grafana
```

---

# 5. Gatewayの責務

Goで実装するGatewayは以下を担当する。

1. デバイスの検出
2. デバイス種類の判定
3. 適切なDriverの選択
4. デバイスの登録
5. センサーデータの取得
6. MQTTへのデータ送信
7. デバイス状態の管理
8. 将来的なGPIO/BLEデバイスの管理

Gateway自身は、データの長期保存や可視化を主目的としない。

---

# 6. デバイス自動検出

## 6.1 I²Cスキャン

Gateway起動時にI²Cバスをスキャンし、接続されているデバイスを検出する。

現在の対象バス:

```text
I²C Bus 1
```

Linux環境では `i2cdetect` 等を利用してI²Cデバイスの存在を確認する。

---

## 6.2 アドレス検出

I²Cスキャンによってデバイスアドレスを取得する。

例:

```text
Scanning I²C bus...

Found device: 0x76
Found device: 0x23
```

検出結果は `driver.DeviceInfo` としてGateway内部へ渡す。

---

# 7. センサー種類の自動判定

検出されたI²Cデバイスについて、対応するDriverが自分自身で対応可能か判定する。

概念:

```text
I²C Device
     │
     ▼
DeviceInfo
     │
     ▼
Driver Registry
     │
     ├── BME280 Driver
     ├── BH1750 Driver
     └── その他のDriver
```

各Driverは `Detect()` を実装する。

例:

```go
func (s *Sensor) Detect(info driver.DeviceInfo) bool
```

---

# 8. Driver設計

センサーごとの処理は個別Driverとして実装する。

基本的なインターフェースは以下とする。

```go
type SensorDriver interface {
    Name() string
    Detect(info DeviceInfo) bool
    Read() (float64, error)
}
```

各Driverは共通インターフェースを実装することで、Device Managerから統一的に扱えるようにする。

---

# 9. BME280 Driver

BME280 DriverはBME280センサーとの通信を担当する。

現在の基本構造:

```go
type Sensor struct{}

func (s *Sensor) Name() string {
    return "BME280"
}

func (s *Sensor) Detect(info driver.DeviceInfo) bool {
    return info.Address == 0x76 ||
           info.Address == 0x77
}
```

将来的にはアドレスだけではなく、センサー内部のChip ID等を確認して、より確実なデバイス判定を行う。

---

# 10. BH1750 Driver

BH1750 DriverはBH1750センサーとの通信を担当する。

主な責務:

* BH1750の検出
* 測定モード設定
* 照度取得
* エラー処理

---

# 11. Driver Registry

Driver Registryでは、Gatewayが利用可能なDriverを管理する。

概念:

```text
Registry
│
├── BME280
└── BH1750
```

Discoveryによってデバイスが発見された場合、Registryから対応するDriverを検索する。

---

# 12. Device Manager

Device ManagerはGateway上で認識されたデバイスを管理する。

責務:

* Deviceの登録
* Deviceの削除
* Device一覧の管理
* Driverとの関連付け
* センサーの取得処理

概念:

```text
Discovery
   ↓
DeviceInfo
   ↓
Registry
   ↓
Driver
   ↓
Device Manager
   ↓
Sensor
```

---

# 13. 自動登録フロー

Gateway起動後は、以下の処理を自動的に行う。

```text
1. Gateway起動
       ↓
2. I²C Bus Scan
       ↓
3. デバイス発見
       ↓
4. DeviceInfo生成
       ↓
5. Driver Registry検索
       ↓
6. Driver Detect()
       ↓
7. 対応Driver決定
       ↓
8. Device Managerへ登録
       ↓
9. センサー値取得
       ↓
10. MQTTへ送信
```

---

# 14. MQTT

Gatewayと外部システム間の通信にはMQTTを使用する。

MQTT BrokerにはMosquittoを使用する。

概念:

```text
Sensor
  ↓
Go Gateway
  ↓
MQTT
  ↓
Mosquitto
  ↓
Home Assistant
```

---

# 15. Home Assistant連携

Home Assistantをセンサーデータの主要な可視化・操作システムとして利用する。

Gatewayはセンサー情報をMQTT経由でHome Assistantへ提供する。

将来的にはMQTT Discovery等を利用し、Gatewayへ新しいセンサーを追加した際にHome Assistant側でも自動的に認識できる構成を目指す。

---

# 16. データ保存・可視化

取得したセンサーデータについて、以下のシステムとの連携を想定する。

### InfluxDB

時系列データの保存を担当する。

### Grafana

蓄積されたデータの分析・可視化を担当する。

### Home Assistant

リアルタイムの状態確認およびスマートホームシステムとの連携を担当する。

---

# 17. Docker

Gateway周辺のサービスはDockerを利用して構築する。

主なコンテナ:

```text
docker
│
├── Go Gateway
├── Mosquitto
├── Home Assistant
├── InfluxDB
└── Grafana
```

サービスの設定は可能な限りコード・設定ファイルとして管理し、環境を再構築できるようにする。

---

# 18. エラー処理

以下の異常を想定する。

* I²Cバスが利用できない
* センサーが応答しない
* センサー通信に失敗する
* 未対応デバイスが検出される
* MQTT Brokerへ接続できない
* センサー読み取りに失敗する

未対応デバイスが検出された場合、Gateway全体を停止させず、Unknown Deviceとして扱える構成を目指す。

---

# 19. ログ

Gatewayは主要な処理についてログを出力する。

例:

```text
Scanning I²C bus...
Found device: 0x76
Detected sensor: BME280
Registered device: BME280
Reading sensor...
Publishing MQTT message...
```

エラー発生時には原因を特定できる情報を出力する。

---

# 20. 拡張性

将来的にはI²C以外のデバイスにも対応する。

対象候補:

```text
I²C
 ├── BME280
 ├── BH1750
 └── その他センサー

BLE
 └── BLEデバイス

GPIO
 └── GPIOデバイス
```

ただし、MVPではI²Cセンサーを優先する。

---

# 21. ESP32について

当初想定していたESP32による中継・UI・GPIO機能については、現在のMVPでは使用しない。

まずRaspberry Pi単体で、

```text
センサー
 ↓
Go Gateway
 ↓
MQTT
 ↓
Home Assistant
```

という基本システムを完成させる。

その後、必要に応じてBLE/GPIO等の機能を拡張する。

---

# 22. MVP

最初の完成目標は以下とする。

### 必須

* [x] Raspberry Pi上でGo Gatewayを起動できる
* [x] I²Cバスをスキャンできる
* [x] I²Cデバイスのアドレスを検出できる
* [x] BME280を認識できる
* [ ] BH1750を認識できる
* [ ] Driver Registryを利用できる
* [ ] Device Managerへ自動登録できる
* [ ] センサー値を取得できる
* [ ] MQTTへ送信できる
* [ ] Home Assistantで確認できる

### 後回し

* BLE
* GPIO
* AIによる判断
* 高度な自動制御
* ESP32連携
* 高度なデバイス自動判定

---

# 23. 完成後の目標

最終的には、新しいセンサーを追加する際にGateway本体の処理を書き換えるのではなく、

```text
新しいDriverを追加
       ↓
Registryへ登録
       ↓
Gatewayが自動検出
       ↓
自動登録
       ↓
自動取得
       ↓
MQTTへ公開
```

という構造を実現する。

これにより、センサーの種類が増えても既存システムへの影響を小さくする。

---

# 24. システム全体像

```text
                  Raspberry Pi
┌─────────────────────────────────────────┐
│                                         │
│  I²C Sensors                            │
│  ┌────────┐    ┌────────┐               │
│  │ BME280 │    │ BH1750 │               │
│  └────┬───┘    └────┬───┘               │
│       │             │                   │
│       └──────┬──────┘                   │
│              ↓                          │
│       ┌──────────────┐                  │
│       │  Discovery   │                  │
│       └──────┬───────┘                  │
│              ↓                          │
│       ┌──────────────┐                  │
│       │DeviceInfo    │                  │
│       └──────┬───────┘                  │
│              ↓                          │
│       ┌──────────────┐                  │
│       │Driver        │                  │
│       │Registry      │                  │
│       └──────┬───────┘                  │
│              ↓                          │
│       ┌──────────────┐                  │
│       │Device        │                  │
│       │Manager       │                  │
│       └──────┬───────┘                  │
│              ↓                          │
│       ┌──────────────┐                  │
│       │ Sensor       │                  │
│       │ Drivers      │                  │
│       └──────┬───────┘                  │
│              ↓                          │
│            MQTT                         │
│              │                          │
└──────────────┼──────────────────────────┘
               ↓
          ┌──────────┐
          │ Mosquitto│
          └────┬─────┘
               ↓
       ┌───────────────┐
       │Home Assistant │
       └───────┬───────┘
               │
       ┌───────┴────────┐
       ↓                ↓
   InfluxDB          Grafana
```

---

# 25. 開発方針

本プロジェクトでは、機能追加よりも先に**責務の分離とインターフェースの統一**を重視する。

特に以下を原則とする。

* Discoveryは「デバイスを発見する」
* Driverは「特定デバイスとの通信を担当する」
* Registryは「Driverを管理する」
* Device Managerは「認識済みデバイスを管理する」
* MQTTは「データを外部へ渡す」
* Home Assistantは「状態の可視化・操作を担当する」

各コンポーネントが別の責務を持つことで、将来的なセンサー追加や通信方式追加に対応しやすい構成とする。

---

# 26. 現在の開発状況

2026年9月25日時点では、Go Gatewayの基本構造を実装中。

現在までに、

* I²Cスキャン
* `DeviceInfo`
* BME280 Driver
* Driver Registry
* Device Manager

などの基礎部分を構築している。

I²Cスキャンについては、単純なアドレス総当たりではなく、`i2cdetect` の結果を正しく解析して実際に検出されたアドレスを取得する方式へ変更した。

今後は、

```text
I²C Discovery
      ↓
Driver Detect
      ↓
自動Register
      ↓
Sensor Read
      ↓
MQTT
```

の一連の処理を完成させる。
