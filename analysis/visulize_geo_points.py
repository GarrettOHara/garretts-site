import json
import os
import matplotlib.pyplot as plt
import folium
# from mpl_toolkits.basemap import Basemap

DATA_PATH = "../static"
MAP_FILE = "visitor_map.html"

def plot_country_distribution():
    with open(os.path.join(DATA_PATH, "country_distribution.json")) as f:
        data = json.load(f)

    labels = data["labels"]
    values = data["values"]

    plt.figure(figsize=(10, 6))
    plt.bar(labels, values, color="skyblue")
    plt.title("Top 10 Visitor Countries")
    plt.ylabel("Visits")
    plt.xlabel("Country")
    plt.xticks(rotation=45)
    plt.tight_layout()
    plt.show()

def plot_device_by_country():
    with open(os.path.join(DATA_PATH, "device_by_country.json")) as f:
        data = json.load(f)

    labels = data["labels"]
    desktop = data["desktop"]
    mobile = data["mobile"]

    x = range(len(labels))

    plt.figure(figsize=(10, 6))
    plt.bar(x, desktop, width=0.4, label="Desktop", align="center")
    plt.bar([i + 0.4 for i in x], mobile, width=0.4, label="Mobile", align="center")
    plt.xticks([i + 0.2 for i in x], labels, rotation=45)
    plt.title("Device Type by Country")
    plt.xlabel("Country")
    plt.ylabel("Visits")
    plt.legend()
    plt.tight_layout()
    plt.show()

# def plot_visitor_map():
#     with open(os.path.join(OUTPUT_DIR, "visitor_map.json")) as f:
#         data = json.load(f)
# 
#     plt.figure(figsize=(12, 7))
#     m = Basemap(projection="robin", lon_0=0, resolution="c")
#     m.drawcoastlines()
#     m.drawcountries()
#     m.drawmapboundary()
# 
#     for entry in data:
#         lat = entry["lat"]
#         lng = entry["lng"]
#         size = entry["count"] * 2  # Scale marker size
#         x, y = m(lng, lat)
#         m.plot(x, y, 'ro', markersize=size)
# 
#     plt.title("Visitor Map")
#     plt.show()

def plot_visitor_map():
    with open(os.path.join(DATA_PATH, "visitor_map.json")) as f:
        data = json.load(f)

    m = folium.Map(location=[20, 0], zoom_start=2)

    for entry in data:
        lat = entry['lat']
        lng = entry['lng']
        country = entry['country']
        count = entry['count']
        folium.CircleMarker(
            location=[lat, lng],
            radius=min(count * 0.8, 10),
            popup=f"{country}: {count} visitors",
            color='crimson',
            fill=True,
            fill_opacity=0.7
        ).add_to(m)

    m.save(os.path.join(DATA_PATH, MAP_FILE))
    print(f"🌍 Visitor map saved to {DATA_PATH}/{MAP_FILE}")

if __name__ == "__main__":
    print("📊 Displaying geolocation visualizations...")
    plot_country_distribution()
    plot_device_by_country()
    plot_visitor_map()