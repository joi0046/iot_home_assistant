# IoT Gateway Architecture

## 1. Purpose

このプロジェクトは、IoT GatewayとAIエージェントやデータベースなどの外部システムを簡単に接続・連携するための通信規格・基盤を提供する。

IoT Gatewayは、Raspberry Piなどのデバイス上で動作し、センサーやIoTデバイスからデータを取得する。

取得したデータは、MQTTなどを利用して外部システムへ提供する。

主な利用者として、以下のシステムを想定する。

* AIエージェント
* データベース
* Home Assistant
* その他のIoTサービス

## 2. System Overview

```text
                    ┌──────────────────┐
                    │   IoT Devices    │
                    │                  │
                    │ HDC1000          │
                    │ BME280           │
                    │ BH1750           │
                    └────────┬─────────┘
                             │
                            I²C
                             │
                             ▼
                    ┌──────────────────┐
                    │   IoT Gateway    │
                    │                  │
                    │ Discovery        │
                    │ Driver Registry  │
                    │ Device Manager   │
                    │ Sensor Drivers   │
                    └────────┬─────────┘
                             │
                          MQTT / API
                             │
             ┌───────────────┼────────────────┐
             │               │                │
             ▼               ▼                ▼
      ┌────────────┐  ┌─────────────┐  ┌──────────────┐
      │AI Agent    │  │ Time-series │  │Home Assistant│
      │            │  │ Database    │  │              │
      │分析・判断  │  │データ保存   │  │操作・表示    │
      └─────┬──────┘  └─────────────┘  └──────────────┘
            │
            │ 状況判断
            ▼
      Home Assistant
         を操作
```

## 3. Software Architecture

IoT Gateway内部は、センサーの検出・認識とデータ取得を分離する。

```text
cmd/gateway
      │
      ▼
   Gateway
      │
      ├── Discovery
      │      │
      │      ▼
      │   DeviceInfo
      │
      ├── Registry
      │      │
      │      ▼
      │   SensorDriver
      │      ├── HDC1000
      │      ├── BME280
      │      └── BH1750
      │
      └── Device Manager
             │
             ▼
          Reading
             │
             ▼
            MQTT
```

## 4. Device Discovery

GatewayはI²Cバスをスキャンし、接続されているデバイスを自動的に検出する。

検出したデバイスについて、I²Cアドレスなどの情報を`DeviceInfo`として扱う。

その後、Driver Registryから対応するDriverを検索し、認識可能なデバイスであれば登録する。

これにより、Gatewayにセンサーを接続した際に、個別の設定をできるだけ減らす。

## 5. Driver

各センサーは共通の`SensorDriver`インターフェースを実装する。

```text
SensorDriver
 ├── Name()
 ├── Detect()
 └── Read()
```

`Detect()`によってデバイスがそのDriverに対応しているか判断し、`Read()`によってセンサーのデータを取得する。

取得したデータは共通の`Reading`形式に変換する。

```text
Sensor
   │
   ▼
Reading
 ├── Temperature
 ├── Humidity
 └── Lux
```

## 6. Data Flow

センサーから外部システムまでのデータは、以下の流れで処理する。

```text
I²C Sensor
    │
    ▼
Discovery
    │
    ▼
DeviceInfo
    │
    ▼
Driver Detection
    │
    ▼
Device Manager
    │
    ▼
Sensor.Read()
    │
    ▼
Reading
    │
    ▼
MQTT
    │
    ├──────────────► AI Agent
    │
    ├──────────────► Database
    │
    └──────────────► Home Assistant
```

## 7. AI Agent

AI AgentはGatewayから提供されたセンサーデータを取得し、その情報を理解・分析する。

AI Agent自身がセンサーとの通信を直接行うのではなく、GatewayがIoTデバイスとの通信を担当する。

AI Agentは取得したデータを利用して状況を判断し、必要に応じてHome Assistantなどの外部システムを操作する。

```text
Sensor
  ↓
IoT Gateway
  ↓
AI Agent
  ↓
状況分析・判断
  ↓
Home Assistant
  ↓
デバイス操作
```

## 8. Database

センサーから取得したデータは、時系列データとしてデータベースに保存する。

これにより、現在の状態だけではなく、過去のデータを利用できるようにする。

```text
Sensor
  ↓
Gateway
  ↓
MQTT
  ↓
Database
  ↓
Time-series Data
```

保存されたデータは、後から分析や可視化、AI Agentによる判断などに利用する。

## 9. Home Assistant

Home AssistantはIoTデバイスの状態表示や操作を担当する。

Gatewayはセンサーとの低レベルな通信を担当し、Home Assistantとの連携にはMQTTなどの標準的な通信方式を利用する。

これにより、GatewayとHome Assistantを疎結合にする。

## 10. Communication

Gatewayと外部システム間の通信にはMQTTを利用する。

MQTTを利用することで、AI Agent、Database、Home AssistantなどをGatewayから独立して接続できるようにする。

```text
             ┌── AI Agent
             │
Gateway ─ MQTT ─ Database
             │
             └── Home Assistant
```

## 11. Design Goals

このプロジェクトでは以下を目標とする。

* IoTデバイスの自動検出
* センサーDriverの共通化
* 外部システムからGatewayを簡単に利用できること
* AI Agentとの連携
* 時系列データベースとの連携
* Home Assistantとの連携
* 手動設定をできるだけ減らすこと
* 新しいセンサーをDriver追加だけで対応できること
* Gatewayと利用側システムを疎結合にすること

## 12. Future

今後は以下の機能を検討する。

* MQTT Discovery
* センサーDriverの追加
* GatewayのHTTP/API
* Gatewayの状態・デバイス情報取得
* AI Agent向けのデータ取得インターフェース
* 時系列データベースとの標準的な連携方法
* 複数Gatewayへの対応
* Gateway間通信
* センサーのホットプラグ対応

````

これ、**かなり方向性が変わった**と思う。

最初は「I²Cセンサーを自動検出するGateway」だったけど、今の定義だと、

> **IoTの物理世界と、AI・DB・Home Assistantなどのソフトウェア世界の間をつなぐ標準化されたGateway**

になってる。

そして重要なのが、**AI AgentにI²Cの知識を持たせない**こと。

```text
        AI Agent
           │
     「温度どうなってる？」
           │
           ▼
      Gateway / MQTT
           │
           ▼
       HDC1000
           │
           ▼
        24.3°C
````

AI側は「HDC1000」「I²Cアドレス0x40」「レジスタ0xFE」なんて知らなくていい。

逆にGateway側も「AIがどう判断するか」は知らなくていい。

この**責務分離**が、このプロジェクトの設計上かなり大事なところになりそう。
