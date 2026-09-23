# iot_home_assistant

ラズパイを使用したIoTゲートウェイシステム。

センサーやデバイスの接続・認識・データ収集・可視化・制御をできるだけ自動化して、最終的にローカルAIによる判断・制御までを行うことを目指す。

技術スタックとして
Golang
home assistant
InfluxDB
Grafana
Hermes Agent
Bonsai

ゴールとして
    1. センサーの自動認識
    2. センサーDriverの共通化
    3. MQTTによる疎結合なシステム
    4. Home Assistantへの自動登録
    5. センサーデータの長期保存・可視化
    6. ESP32との連携
    7. ローカルAIによるIoT制御
