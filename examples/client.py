#!/usr/bin/env python3
import socket
import struct

sock = socket.socket(socket.AF_UNIX, socket.SOCK_STREAM)
sock.connect('/tmp/lock00.sock')

while True:
    data = b'python'
    wirebytes = struct.pack('b6s', len(data), data)
    sock.send(wirebytes)
    res = sock.recv(2)
    print("Sent %s and got %s" % (wirebytes, res))
