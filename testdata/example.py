import requests

def fetch_data(url):
    try:
        response = requests.get(url)
        if response.status_code == 200:
            return response.json()
        else:
            return None
    except Exception as e:
        print(f"Error fetching data: {e}")
        return None

def process_data(data):
    if not data:
        return "No data to process"
    
    results = []
    for item in data:
        if "value" in item:
            results.append(item["value"] * 2)
    
    return results

def main():
    url = "https://jsonplaceholder.typicode.com/posts"
    data = fetch_data(url)
    processed = process_data(data)
    
    if processed:
        print("Processed Data:", processed)
    else:
        print("No valid data found.")

if __name__ == "__main__":
    main()
