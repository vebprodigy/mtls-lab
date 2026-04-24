# server.py
import ssl
import socket
import threading

def start_mtls_server(host='localhost', port=4443):
    """
    mTLS server that requires client certificate.
    """
    # Create SSL context
    context = ssl.SSLContext(ssl.PROTOCOL_TLS_SERVER)
    
    # Load server certificate and key
    context.load_cert_chain(
        certfile='server.crt',
        keyfile='server.key'
    )
    
    # Load CA bundle (root + intermediate) for verification
    context.load_verify_locations('ca-bundle.crt')
    
    # CRITICAL: Require client certificate
    context.verify_mode = ssl.CERT_REQUIRED
    
    # Optional: Set up socket
    with socket.socket(socket.AF_INET, socket.SOCK_STREAM) as sock:
        sock.setsockopt(socket.SOL_SOCKET, socket.SO_REUSEADDR, 1)
        sock.bind((host, port))
        sock.listen(1)
        print(f"[SERVER] Listening on {host}:{port}")
        
        # Accept connections with SSL wrapper
        with context.wrap_socket(sock, server_side=True) as ssock:
            while True:
                try:
                    # Accept client connection
                    conn, addr = ssock.accept()
                    print(f"[SERVER] Client connected from {addr}")
                    
                    # Get client certificate info
                    client_cert = conn.getpeercert()
                    print(f"[SERVER] Client certificate subject: {client_cert['subject']}")
                    print(f"[SERVER] Client certificate CN: {client_cert['subject'][0][0][1]}")
                    
                    # Receive and send data
                    data = conn.recv(1024).decode()
                    print(f"[SERVER] Received: {data}")
                    
                    response = f"Echo: {data}"
                    conn.sendall(response.encode())
                    print(f"[SERVER] Sent: {response}")
                    
                    conn.close()
                except Exception as e:
                    print(f"[SERVER] Error: {e}")

if __name__ == '__main__':
    start_mtls_server()