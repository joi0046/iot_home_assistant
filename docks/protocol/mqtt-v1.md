## MQTT Protocol v1

Go IoT Gateway と外部サービス間でセンサーデータを交換するための通信規格。

本プロトコルは、Home Assistant、データベース、AI Agent などの外部システムが、センサーの具体的な実装やドライバを意識せずにGatewayからデータを取得できることを目的とする。

### Topic

センサーデータは以下のMQTT TopicにPublishする。

```text
gateway/v1/readings
```

すべてのセンサーの計測値を同一TopicへPublishする。

### Message Format

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
      "value": 28.99505615234375,
      "unit": "°C"
    },
    "humidity": {
      "value": 71.563720703125,
      "unit": "%"
    }
  },
  "timestamp": "2026-09-26T18:47:07.20007484+09:00"
}
```

### Fields

#### `version`

プロトコルのバージョン。

現在のバージョンは `"1"`。

```json
"version": "1"
```

#### `device`

データを取得したデバイスの情報。

| Field  | Type   | Description       |
| ------ | ------ | ----------------- |
| `id`   | string | Gateway内でのデバイス識別子 |
| `type` | string | センサーの種類           |

Example:

```json
"device": {
  "id": "i2c-1/0x40",
  "type": "HDC1000"
}
```

##### Device ID

Device IDは以下の形式で自動生成する。

```text
{bus}/{address}
```

Example:

```text
i2c-1/0x40
```

これにより、同じI²Cアドレスを使用するデバイスが異なるバスに存在する場合でも識別できる。

デバイスIDは手動登録ではなく、Gatewayのデバイス検出機構から生成する。

### `data`

センサーが取得した物理量を格納する。

データ項目名にはセンサー固有の名前ではなく、物理量を表す名前を使用する。

Example:

```json
"data": {
  "temperature": {
    "value": 28.99505615234375,
    "unit": "°C"
  },
  "humidity": {
    "value": 71.563720703125,
    "unit": "%"
  }
}
```

各データ項目は以下の形式とする。

| Field   | Type   | Description |
| ------- | ------ | ----------- |
| `value` | number | 計測値         |
| `unit`  | string | 計測値の単位      |

### Standard Data Keys

現在定義している物理量のキーは以下。

| Key           | Physical Quantity | Example Unit |
| ------------- | ----------------- | ------------ |
| `temperature` | 温度                | `°C`         |
| `humidity`    | 湿度                | `%`          |
| `pressure`    | 気圧                | `hPa`        |
| `illuminance` | 照度                | `lx`         |
| `co2`         | CO₂濃度             | `ppm`        |
| `voltage`     | 電圧                | `V`          |
| `current`     | 電流                | `A`          |
| `power`       | 電力                | `W`          |

新しいセンサーを追加する場合も、可能な限り既存の物理量キーを使用する。

例えば、照度センサーの値は `lux` ではなく `illuminance` とする。

```json
"illuminance": {
  "value": 320.5,
  "unit": "lx"
}
```

`lux` は単位であり、物理量の名称ではないためである。

### `timestamp`

センサー値を取得した時刻。

ISO 8601形式の日時文字列として記録する。

```json
"timestamp": "2026-09-26T18:47:07.20007484+09:00"
```

`timestamp` は `data` の外側に配置する。

### Design Goals

本プロトコルは以下を目的とする。

* センサーの種類に依存しない共通データ形式
* センサーをGatewayへ自動登録できる構成
* Home Assistantから容易に利用できること
* データベースで時系列データとして保存しやすいこと
* AI Agentがセンサーデータを機械的に解釈しやすいこと
* 将来的なプロトコル拡張に対応できること

### Data Flow

```text
Sensor
  ↓
Go IoT Gateway
  ↓
gateway/v1/readings
  ↓
MQTT
  ├──→ Home Assistant
  ├──→ Database
  └──→ AI Agent
```

現在のv1では、Gatewayから外部サービスへの**計測データ配信**を対象とする。

デバイス制御などの逆方向通信については、将来のプロトコル拡張で扱う。
