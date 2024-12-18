import struct
import time
import socket

def decode_payload(payload):
    index = 0
    
    # Step 1: Get the count (the first byte)
    count = payload[index]
    index += 1
    print(f"Count (amount): {count}")

    # Step 2: Extract the circular buffer data (the next `count` bytes)
    buffer_data = payload[index:index + count]
    index += count
    print(f"Buffer Data: {list(buffer_data)}")
    
    # Step 3: Extract the sampling frequency (next byte)
    sampling_frequency = payload[index]
    index += 1
    print(f"Sampling Frequency: {sampling_frequency}")
    
    # Step 4: Extract the MAC address (next 8 bytes)
    mac_address = payload[index:index + 8]
    index += 8
    # Convert the MAC address to a human-readable format
    mac_address_str = ':'.join(f"{byte:02x}" for byte in mac_address)
    print(f"MAC Address: {mac_address_str}")
    
    return {
        "count": count,
        "buffer_data": list(buffer_data),
        "sampling_frequency": sampling_frequency,
        "mac_address": mac_address_str,
    }

# Connect to the server
client_sock = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
client_sock.connect(('192.168.135.147', 8080))

# Send the 0x02 byte (this simulates the behavior of the receiver in the C code)
client_sock.sendall(b'\x02')  # Sending byte 0x02

# Receive data from the server (adjust buffer size if needed)
received_data = client_sock.recv(128)

# Print the received data byte by byte
for byte in received_data:
    print(f"Received byte: {byte}")
print(f"Received data: {received_data}")

# Process the received byte array
decoded_data = decode_payload(received_data)

# Print the decoded data
print("Decoded Data:", decoded_data)

# Close the connection
client_sock.close()
