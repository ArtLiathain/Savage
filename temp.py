import psutil
import time

def simulate_memory_usage():
    process = psutil.Process()  # Get the current process
    memory_usage = []  # List to hold large data for simulation

    print("Simulating memory usage. Press Ctrl+C to stop.")

    try:
        while True:
            # Simulate memory usage by appending large lists to memory_usage
            memory_usage.append([0] * 10**7)  # Allocate 10 million integers (about 40MB)

            # Get current memory usage of the process
            memory_info = process.memory_info()

            # Print memory usage in MB
            print(f"Memory Usage: {memory_info.rss / 1024 / 1024:.2f} MB")
            
            # Sleep for a short period before the next iteration
            time.sleep(1)
    
    except KeyboardInterrupt:
        print("\nSimulation stopped by user.")
        print("Final Memory Usage:", memory_info.rss / 1024 / 1024, "MB")

if __name__ == "__main__":
    simulate_memory_usage()
