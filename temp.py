import time

# Number of iterations
iterations = 1000

class TimeManager:
    def __enter__(self):
        self.startTime = time.perf_counter_ns()
    def __exit__(self, exc_type, exc_val, exc_tb):
        print(f"Performance counter time: {(( time.perf_counter_ns() - self.startTime)/1000000000):.6f} nanoseconds")

with TimeManager():
    for i in range(iterations):
        # Perform some operation
        x = i * i
