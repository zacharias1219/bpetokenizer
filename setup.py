import os
from setuptools import setup, find_packages
from setuptools.command.build_ext import build_ext
import subprocess

here = os.path.abspath(os.path.dirname(__file__))

# Get the code version
version = {}
with open(os.path.join(here, "bpetokenizer/version.py")) as f:
    exec(f.read(), version)
__version__ = version["__version__"]

# Custom build command to compile Go code
class BuildGoExtension(build_ext):
    def run(self):
        # Compile Go code into a shared library
        go_dir = os.path.join(here, "bpetokenizer/go")
        subprocess.check_call(["go", "build", "-buildmode=c-shared", "-o", os.path.join(go_dir, "libbpe.so"), os.path.join(go_dir, "bpe.go")])

with open("README.md", "r", encoding="utf-8") as f:
    long_description = f.read()


setup(
    name="bpetokenizer",
    version=__version__,
    description="A Byte Pair Encoding (BPE) tokenizer with Go optimizations.",
    long_description=long_description,
    long_description_content_type="text/markdown",
    url="https://github.com/Hk669/bpetokenizer",
    author="Hrushikesh Dokala",
    author_email="hrushi669@gmail.com",
    license="MIT",
    packages=find_packages(include=["bpetokenizer"]),
    package_data={
        'bpetokenizer': [
            'pretrained/wi17k_base/wi17k_base.json',
            'go/libbpe.so',  # Include the Go shared library
        ],
    },
    classifiers=[
        "License :: OSI Approved :: MIT License",
        "Programming Language :: Python :: 3",
        "Operating System :: OS Independent",
    ],
    install_requires=["regex"],
    extras_require={
        "dev": ["pytest", "twine"],
    },
    python_requires=">=3.9,<3.13",
    cmdclass={'build_ext': BuildGoExtension},
    ext_modules=[],  # Placeholder for Go compilation
)