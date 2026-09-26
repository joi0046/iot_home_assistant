# MQTT Protocol v1

## 1. 概要

Go IoT Gatewayでは、センサーから取得したデータを外部サービスへ提供するためにMQTTを使用する。

本プロトコルは、AI Agent、データベース、Home Assistantなどの外部システムが、Gatewayから提供されるセンサーデータを共通の形式で利用できることを目的とする。

Gatewayはセンサー固有のI²C通信を担当し、外部システムにはセンサーの種類に依存しない統一されたReading形式でデータを提供する。

## 2. 通信構成

```text
Sensor / I²C
      ↓
Go IoT Gateway
      ↓
    MQTT
      ↓
 ┌────┼──────────────┐
 ↓    ↓              ↓
Home  Database     AI Agent
Assistant             ↓
                  Home Assistant API
```

## 3. MQTT Topic

センサーの測定値は以下のTopicへPublishする。

```text
gateway/v1/readings
```

### Topic構成

```text
gateway/
└── v1/
    └── readings
```

* `gateway`: Go IoT Gatewayが提供するデータ
* `v1`: Protocol version 1
* `readings`: センサーの測定値

Protocolのバージョンが変更された場合は、Topicのバージョンも変更する。

例:

```text
gateway/v2/readings
```

## 4. Payload

PayloadはJSON形式とする。

```json
{
  "version": "1",
  "device": {
    "id": "i2c-1/0x40",
    "type": "HDC1000"
  },
  "data": {
    "temperature": {
      "value": 31.6,
      "unit": "°C"
    },
    "humidity": {
      "value": 66.9,
      "unit": "%"
    }
  },
  "timestamp": "2026-09-26T20:18:12.611629168+09:00"
}
```

## 5. version

Protocolのバージョンを表す。

```json
"version": "1"
```

Protocol v1では文字列として扱う。

## 6. device

センサーを識別するための情報。

```json
"device": {
  "id": "i2c-1/0x40",
  "type": "HDC1000"
}
```

### device.id

デバイスの物理的な接続先を基準として識別する。

形式:

```text
{bus}/{address}
```

例:

```text
i2c-1/0x40
i2c-1/0x76
```

### device.type

センサーの種類を表す。

例:

```text
HDC1000
BME280
BH1750
```

## 7. data

センサーから取得した測定値を格納する。

```json
"data": {
  "temperature": {
    "value": 31.6,
    "unit": "°C"
  },
  "humidity": {
    "value": 66.9,
    "unit": "%"
  }
}
```

`data`は固定されたフィールドではなく、センサーが提供する測定項目に応じて動的に構成される。

これにより、センサーごとに異なる測定項目を同一のProtocolで扱える。

## 8. Measurement Key

測定項目は物理量を表す名前を使用する。

| Key           | 意味      | Unit  |
| ------------- | ------- | ----- |
| `temperature` | 温度      | `°C`  |
| `humidity`    | 湿度      | `%`   |
| `pressure`    | 気圧      | `hPa` |
| `illuminance` | 照度      | `lx`  |
| `co2`         | 二酸化炭素濃度 | `ppm` |
| `voltage`     | 電圧      | `V`   |
| `current`     | 電流      | `A`   |
| `power`       | 電力      | `W`   |

センサー固有の名称ではなく、可能な限り一般的な物理量の名称を使用する。

例えばBH1750の測定値は、

```text
lux
```

ではなく、

```text
illuminance
```

として扱う。

## 9. Value

各測定項目は以下の形式で表現する。

```json
{
  "value": 31.6,
  "unit": "°C"
}
```

### value

測定値。

JSON上では数値として扱う。

### unit

測定値の単位。

単位は各測定値に含めることで、受信側がセンサーの種類を知らなくても値の意味を判断できるようにする。

## 10. timestamp

測定を行った時刻を表す。

```json
"timestamp": "2026-09-26T20:18:12.611629168+09:00"
```

ISO 8601形式を使用する。

`timestamp`は`data`の外側に配置する。

## 11. データ処理の責務

Gatewayは以下を担当する。

```text
I²C通信
 ↓
センサー検出
 ↓
センサー読み取り
 ↓
Reading生成
 ↓
Protocol形式へ変換
 ↓
JSON化
 ↓
MQTT Publish
```

一方、受信側はセンサー固有のI²C通信を意識せず、Protocolに従ってデータを利用する。

## 12. 想定される利用者

Protocol v1では、主に以下のシステムからの利用を想定する。

### Home Assistant

MQTTを利用してセンサー状態を取得する。

### Database

センサーの時系列データを保存する。

```text
timestamp
device_id
measurement
value
unit
```

### AI Agent

センサー情報を取得し、環境の状態を分析する。

例:

```text
temperature = 31.6 °C
humidity = 66.9 %
```

などのデータを利用して状況を判断する。

## 13. Protocol v1の範囲

Protocol v1では、Gatewayから外部システムへの**測定値の提供**を対象とする。

```text
Gateway → MQTT → External System
```

デバイスを制御するための逆方向通信についてはProtocol v1の対象外とする。

将来的には、

```text
AI Agent
    ↓
  MQTT
    ↓
Gateway
    ↓
Device
```

のような制御機能を追加する可能性がある。

## 14. 設計方針

Protocol v1では以下を重視する。

* センサーの種類に依存しない
* 外部システムから扱いやすい
* 人間がJSONを見ても意味を理解できる
* AI Agentが解釈しやすい
* Databaseへ保存しやすい
* Protocolのバージョンアップに対応できる
* Gateway内部の実装と外部Protocolを分離する

これにより、Go IoT Gatewayをセンサーと外部システムの間にある共通データ基盤として利用できるようにする。
