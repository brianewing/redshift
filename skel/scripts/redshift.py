import sys, os

MAX_BUFFER_SIZE = 65536

def log(msg):
    sys.stderr.write(msg + "\n")

def run(effect):
    while True:
        frame = read()
        result = effect(frame)
        write(result or frame)

def read():
    buf = os.read(sys.stdin.fileno(), MAX_BUFFER_SIZE)
    return unpack(buf)

def write(frame):
    os.write(sys.stdout.fileno(), pack(frame))

def pack(frame):
    return bytearray(color for led in frame for color in led)

def unpack(buf):
    return [buf[i:i+3] for i in xrange(0, len(buf), 3)]
