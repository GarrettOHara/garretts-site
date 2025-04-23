import json
import os
import matplotlib.pyplot as plt
from datetime import datetime

OUTPUT_DIR = "./output_data"

def plot_time_series():
    with open(os.path.join(OUTPUT_DIR, "time_series.json")) as f:
        data = json.load(f)

    labels = [datetime.strptime(ts, "%Y-%m-%d %H:%M") for ts in data["labels"]]
    values = data["values"]
    rolling = data["rolling_avg"]

    plt.figure(figsize=(12, 6))
    plt.plot(labels, values, label="Visits per Hour")
    plt.plot(labels, rolling, label="3-hr Rolling Avg", linestyle='--')
    plt.title("Website Visits Over Time")
    plt.xlabel("Timestamp")
    plt.ylabel("Hits")
    plt.legend()
    plt.xticks(rotation=45)
    plt.tight_layout()
    plt.show()

def plot_clusters():
    with open(os.path.join(OUTPUT_DIR, "clusters.json")) as f:
        data = json.load(f)

    clusters = data["clusters"]
    labels = list(map(str, clusters.keys()))
    sizes = list(clusters.values())

    plt.figure(figsize=(6, 6))
    plt.pie(sizes, labels=labels, autopct='%1.1f%%')
    plt.title("Clustering of Visitor Behavior")
    plt.show()

def plot_anomalies():
    with open(os.path.join(OUTPUT_DIR, "anomalies.json")) as f:
        data = json.load(f)

    timestamps = [datetime.strptime(ts, "%Y-%m-%d %H:%M:%S") for ts in data["timestamps"]]

    plt.figure(figsize=(10, 4))
    plt.hist(timestamps, bins=20, edgecolor='black')
    plt.title(f"Anomalies Detected ({data['count']} total)")
    plt.xlabel("Timestamp")
    plt.ylabel("Count")
    plt.xticks(rotation=45)
    plt.tight_layout()
    plt.show()

if __name__ == "__main__":
    print("📊 Loading and displaying charts...")
    plot_time_series()
    plot_clusters()
    plot_anomalies()
