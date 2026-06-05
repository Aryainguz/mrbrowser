from setuptools import setup, find_packages

with open("README.md") as f:
    long_description = f.read()

setup(
    name="mrbrowser",
    version="0.1.0",
    description="Python SDK for the Mr. Browser automation engine",
    long_description=long_description,
    long_description_content_type="text/markdown",
    author="Mr. Browser Contributors",
    license="Apache-2.0",
    packages=find_packages(),
    python_requires=">=3.9",
    install_requires=[],  # Zero mandatory dependencies — stdlib only
    extras_require={
        "dev": ["pytest>=7.0", "pytest-mock"],
    },
    classifiers=[
        "Development Status :: 3 - Alpha",
        "Intended Audience :: Developers",
        "License :: OSI Approved :: Apache Software License",
        "Programming Language :: Python :: 3",
        "Programming Language :: Python :: 3.9",
        "Programming Language :: Python :: 3.10",
        "Programming Language :: Python :: 3.11",
        "Programming Language :: Python :: 3.12",
        "Topic :: Software Development :: Testing",
        "Topic :: Internet :: WWW/HTTP :: Browsers",
    ],
    keywords=["browser", "automation", "testing", "selenium", "playwright", "self-healing"],
)
