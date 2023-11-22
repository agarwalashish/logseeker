import time
import random

def generate_log_message():
    messages = [
        "System update completed successfully",
        "User admin logged in from IP 192.168.1.100",
        "Failed to connect to database: timeout error",
        "Disk usage exceeded 80% on /dev/sda1",
        "Starting scheduled backup process",
        "Network interface eth0 went down",
        "Unauthorized access attempt detected",
        "Successfully installed application package",
        "Error parsing configuration file /etc/app/config.yml",
        "Memory usage within normal parameters",
    ]
    return random.choice(messages)

def generate_log_line(timestamp):
    log_level = random.choice(["INFO", "ERROR", "WARN", "DEBUG"])
    message = generate_log_message()
    return f"{timestamp} {log_level} {message}"

def generate_log_file(filename, size_in_mb):
    target_size = size_in_mb * 1024 * 1024  # Convert MB to Bytes
    current_size = 0
    start_time = time.time()

    with open(filename, "w") as file:
        while current_size < target_size:
            timestamp = time.strftime("%Y-%m-%d %H:%M:%S", time.localtime(start_time))
            line = generate_log_line(timestamp) + "\n"
            file.write(line)
            current_size += len(line)
            start_time += 0.001  # Increment the timestamp by 1 millisecond

    print(f"Generated log file of size: {current_size} bytes")

filename = input("Enter the filename for the log file: ")
size_in_mb = float(input("Enter the size of the file in MB: "))

# Usage
generate_log_file(f"../logs/{filename}", 1)  # Generates a 1GB log file
