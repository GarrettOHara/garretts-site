import sqlite3
import pandas as pd
import json
import os
import logging
import time
import requests
from collections import Counter, defaultdict

API_TOKEN = os.getenv("IPINFO_TOKEN")

# --- Setup paths ---
DB_PATH = "../requests.db"
OUTPUT_DIR = "../static/charts"
LOG_FILE = "../analysis/analysis.log"
os.makedirs(OUTPUT_DIR, exist_ok=True)

# --- Setup logging ---
logging.basicConfig(
    filename=LOG_FILE,
    level=logging.INFO,
    format="%(asctime)s %(levelname)s: %(message)s"
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

def get_geolocation(ip_address):
    """Get geolocation data for an IP address using ipinfo.io API"""
    try:
        logging.info(f"Getting geolocation for IP {ip_address}...")
        url = f"https://ipinfo.io/{ip_address}/json?token={API_TOKEN}"
        response = requests.get(url)
        if response.status_code == 200:
            data = response.json()
            return {
                'ip': ip_address,
                'country': data.get('country', 'Unknown'),
                'region': data.get('region', 'Unknown'),
                'city': data.get('city', 'Unknown'),
                'loc': data.get('loc', '0,0')
            }
        else:
            logging.warning(f"Failed to get geolocation for IP {ip_address}: {response.status_code}")
            return {
                'ip': ip_address,
                'country': 'Unknown',
                'region': 'Unknown',
                'city': 'Unknown',
                'loc': '0,0'
            }
    except Exception as e:
        logging.error(f"Error getting geolocation for IP {ip_address}: {str(e)}")
        return {
            'ip': ip_address,
            'country': 'Unknown',
            'region': 'Unknown',
            'city': 'Unknown',
            'loc': '0,0'
        }

@log_time
def process_ip_geolocations():
    """Process all IPs from the database and get their geolocation data"""
    try:
        # Connect to the database
        conn = sqlite3.connect(DB_PATH)
        
        # Get the data from the requests table
        query = "SELECT id, ip_address, user_agent, device_type, visited_at FROM requests ORDER BY id DESC"
        df = pd.read_sql_query(query, conn)
        conn.close()
        
        # Clean the data
        df = df[~df['ip_address'].str.contains(',', na=False)]
        df = df.dropna(subset=['ip_address'])
        
        # Get unique IP addresses to avoid redundant API calls
        unique_ips = df['ip_address'].unique()
        logging.info(f"Found {len(unique_ips)} unique IP addresses to process")
        
        # Get geolocation for each unique IP
        geo_data = {}
        for ip in unique_ips:
            if ip and ip.strip():  # Ensure IP is not empty
                geo_data[ip] = get_geolocation(ip)
                # Add a small delay to avoid rate limiting
                time.sleep(0.1)
        
        # Add geolocation data to the dataframe
        df['country'] = df['ip_address'].map(lambda ip: geo_data.get(ip, {}).get('country', 'Unknown'))
        df['region'] = df['ip_address'].map(lambda ip: geo_data.get(ip, {}).get('region', 'Unknown'))
        df['city'] = df['ip_address'].map(lambda ip: geo_data.get(ip, {}).get('city', 'Unknown'))
        df['loc'] = df['ip_address'].map(lambda ip: geo_data.get(ip, {}).get('loc', '0,0'))
        
        # Parse the location into latitude and longitude
        df[['latitude', 'longitude']] = df['loc'].str.split(',', expand=True).astype(float)
        
        # Convert timestamp to datetime
        df['visited_at'] = pd.to_datetime(df['visited_at'], errors='coerce', utc=True)
        
        # Generate data for visualizations
        generate_country_chart(df)
        generate_device_by_country_chart(df)
        generate_visitor_map_data(df)
        
        logging.info("✅ Geolocation analysis completed successfully")
        return df
        
    except Exception as e:
        logging.exception(f"❌ Error during geolocation analysis: {str(e)}")
        return None

@log_time
def generate_country_chart(df):
    """Generate data for country distribution chart"""
    country_counts = df['country'].value_counts().head(10)
    
    chart_data = {
        "labels": country_counts.index.tolist(),
        "values": country_counts.values.tolist()
    }
    
    with open(f"{OUTPUT_DIR}/country_distribution.json", "w") as f:
        json.dump(chart_data, f)
    
    logging.info(f"Country distribution data saved to {OUTPUT_DIR}/country_distribution.json")

@log_time
def generate_device_by_country_chart(df):
    """Generate data for device type by country chart"""
    # Group by country and device type
    device_by_country = df.groupby(['country', 'device_type']).size().unstack(fill_value=0)
    
    # Get top 10 countries by total visits
    top_countries = df['country'].value_counts().head(10).index.tolist()
    device_by_country = device_by_country.loc[top_countries]
    
    # Prepare data for Chart.js
    chart_data = {
        "labels": device_by_country.index.tolist(),
        "desktop": device_by_country.get('Desktop', pd.Series(0, index=device_by_country.index)).tolist(),
        "mobile": device_by_country.get('Mobile', pd.Series(0, index=device_by_country.index)).tolist()
    }
    
    with open(f"{OUTPUT_DIR}/device_by_country.json", "w") as f:
        json.dump(chart_data, f)
    
    logging.info(f"Device by country data saved to {OUTPUT_DIR}/device_by_country.json")

@log_time
def generate_visitor_map_data(df):
    """Generate data for visitor map visualization"""
    # Group by country and count
    country_data = df.groupby('country').agg({
        'id': 'count',
        'latitude': 'first',
        'longitude': 'first'
    }).reset_index()
    
    # Prepare data for the map
    map_data = []
    for _, row in country_data.iterrows():
        if row['latitude'] != 0 and row['longitude'] != 0:
            map_data.append({
                'country': row['country'],
                'count': int(row['id']),
                'lat': float(row['latitude']),
                'lng': float(row['longitude'])
            })
    
    with open(f"{OUTPUT_DIR}/visitor_map.json", "w") as f:
        json.dump(map_data, f)
    
    logging.info(f"Visitor map data saved to {OUTPUT_DIR}/visitor_map.json")

if __name__ == "__main__":
    logging.info("🌎 Starting IP geolocation analysis...")
    process_ip_geolocations()
