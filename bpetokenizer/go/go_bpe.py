import ctypes
import json
import os

class GoBPE:
    def __init__(self):
        # Load the Go shared library
        lib_path = os.path.join(os.path.dirname(__file__), "libbpe.so")
        self.lib = ctypes.CDLL(lib_path)
        
        # Define function signatures
        self.lib.TrainBPEWrapper.argtypes = [ctypes.c_char_p, ctypes.c_int]
        self.lib.TrainBPEWrapper.restype = ctypes.c_char_p

    def train(self, text, vocab_size):
        # Call Go function
        result = self.lib.TrainBPEWrapper(text.encode(), vocab_size)
        return json.loads(result.decode())