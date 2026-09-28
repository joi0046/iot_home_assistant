import json
import os

from influxdb_client.client.influxdb_client import InfluxDBClient


INFLUX_URL = os.getenv("INFLUX_URL", "http://influxdb:8086")
INFLUX_TOKEN = os.environ["INFLUXDB_TOKEN"]
INFLUX_ORG = os.environ["INFLUXDB_ORG"]
INFLUX_BUCKET = os.environ["INFLUXDB_BUCKET"]


def query_sensor_data():
    query = f'''
from(bucket: "{INFLUX_BUCKET}")
  |> range(start: -5m)
'''

    with InfluxDBClient(
        url=INFLUX_URL,
        token=INFLUX_TOKEN,
        org=INFLUX_ORG,
    ) as client:
        result = client.query_api().query(
            query=query,
            org=INFLUX_ORG,
        )

    readings = []

    for table in result:
        for record in table.records:
            readings.append({
                "time": record.get_time().isoformat(),
                "measurement": record.get_measurement(),
                "field": record.get_field(),
                "value": record.get_value(),
                "device_id": record.values.get("device_id"),
                "device_type": record.values.get("device_type"),
            })

    return readings


def main():
    print("AI Agent starting...")

    readings = query_sensor_data()

    print(f"Retrieved {len(readings)} sensor readings")

    print(json.dumps(
        readings,
        ensure_ascii=False,
        indent=2,
    ))
