# client.py
import ssl
import socket

def mtls_client(host='localhost', port=4443, message="Hello mTLS!"):
    """
    mTLS client that sends its certificate to server.
    """
    # Create SSL context
    context = ssl.SSLContext(ssl.PROTOCOL_TLS_CLIENT)
    
    # Load client certificate and key
    context.load_cert_chain(
        certfile='client.crt',
        keyfile='client.key'
    )
    
    # Load CA bundle (root + intermediate) for verification
    context.load_verify_locations('ca-bundle.crt')
    
    # Verify server certificate
    context.check_hostname = True  # Verify hostname matches cert
    context.verify_mode = ssl.CERT_REQUIRED
    
    # Connect to server
    with socket.socket(socket.AF_INET, socket.SOCK_STREAM) as sock:
        with context.wrap_socket(sock, server_hostname=host) as ssock:
            try:
                ssock.connect((host, port))
                print(f"[CLIENT] Connected to {host}:{port}")
                
                # Get server certificate info
                server_cert = ssock.getpeercert()
                print(f"[CLIENT] Server certificate subject: {server_cert['subject']}")
                
                # Send message
                ssock.sendall(message.encode())
                print(f"[CLIENT] Sent: {message}")
                
                # Receive response
                response = ssock.recv(1024).decode()
                print(f"[CLIENT] Received: {response}")
                
            except ssl.SSLError as e:
                print(f"[CLIENT] SSL Error: {e}")
            except Exception as e:
                print(f"[CLIENT] Error: {e}")

if __name__ == '__main__':
    mtls_client()