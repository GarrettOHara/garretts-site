import sqlite3
import pandas as pd
import json
import os
import logging
import time
import sys
from sklearn.cluster import KMeans
from sklearn.ensemble import IsolationForest

# --- Setup paths ---
# Log cwd and script location for debugging
print("CWD:", os.getcwd())
print("FILE:", os.path.abspath(__file__))

# Robust path to DB
BASE_DIR = os.path.dirname(os.path.abspath(__file__))
DB_PATH = os.path.join(BASE_DIR, "..", "requests.db")
OUTPUT_DIR = os.path.join(BASE_DIR, "..", "static")
os.makedirs(OUTPUT_DIR, exist_ok=True)

# --- Setup logging ---
logging.basicConfig(
    stream=sys.stdout,
    level=logging.DEBUG,
    format="%(asctime)s %(levelname)s [%(filename)s:%(lineno)d] %(message)s"
)

def log_time(func):
    def wrapper(*args, **kwargs):
        start = time.time()
        logging.info(f"Starting {func.__name__}...")
        result = func(*args, **kwargs)
        duration = time.time() - start
        logging.info(f"{func.__name__} completed in {duration:.2f} seconds.")
        return result
    return wrapper

try:
    # --- Load and clean data ---
    logging.info("Loading and cleaning data...")
    conn = sqlite3.connect(DB_PATH)
    df = pd.read_sql_query("SELECT * FROM requests", conn)
    conn.close()
    logging.info(f"Loaded {len(df)} rows.")

    df['visited_at'] = pd.to_datetime(df['visited_at'], errors='coerce', utc=True)
    df = df[~df['ip_address'].str.contains(',')]
    df = df.dropna(subset=['visited_at'])

    df['hour'] = df['visited_at'].dt.hour
    df['day_of_week'] = df['visited_at'].dt.dayofweek
    df['is_desktop'] = (df['device_type'] == 'Desktop').astype(int)
    df['is_mobile'] = (df['device_type'] == 'Mobile').astype(int)

    # Features for clustering and anomalies
    features = df[['hour', 'day_of_week', 'is_desktop', 'is_mobile']]

    @log_time
    def run_clustering():
        logging.info("Starting clustering...")
        kmeans = KMeans(n_clusters=3, random_state=42).fit(features)
        df['cluster'] = kmeans.labels_
        cluster_counts = df['cluster'].value_counts().sort_index().to_dict()
        with open(f"{OUTPUT_DIR}/clusters.json", "w") as f:
            json.dump({"clusters": cluster_counts}, f)

    @log_time
    def run_time_series():
        logging.info("Starting time series analysis...")
        ts = df.set_index('visited_at').resample('h').size()
        ts_rolling = ts.rolling(3).mean().fillna(0)
        with open(f"{OUTPUT_DIR}/time_series.json", "w") as f:
            json.dump({
                "labels": ts.index.strftime("%Y-%m-%d %H:%M").tolist(),
                "values": ts.tolist(),
                "rolling_avg": ts_rolling.tolist()
            }, f)

    @log_time
    def run_anomaly_detection():
        logging.info("Starting anomaly detection...")
        model = IsolationForest(contamination=0.1, random_state=42)
        df['anomaly'] = model.fit_predict(features)
        df['anomaly'] = df['anomaly'].map({1: 0, -1: 1})
        anomalies = df[df['anomaly'] == 1]
        with open(f"{OUTPUT_DIR}/anomalies.json", "w") as f:
            json.dump({
                "timestamps": anomalies['visited_at'].dt.strftime("%Y-%m-%d %H:%M:%S").tolist(),
                "count": int(df['anomaly'].sum())
            }, f)

    # Run each section
    run_clustering()
    run_time_series()
    run_anomaly_detection()

    logging.info("✅ All models executed successfully.")

except Exception as e:
    logging.exception("❌ Error during analysis:")