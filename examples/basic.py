import sys
import time

try:
    from mrbrowser import MrBrowser, MrBrowserError
except ImportError:
    print("Error: mrbrowser SDK not installed.")
    print("Run `pip install -e ./sdk/python` to install it locally.")
    sys.exit(1)

def main():
    print("Connecting to Mr. Browser server at localhost:7331...")
    
    # Using context manager ensures the session is closed when done
    try:
        with MrBrowser(host="localhost", port=7331) as browser:
            if not browser.ping():
                print("Error: Mr. Browser server is not running on localhost:7331.")
                print("Please start it first.")
                sys.exit(1)
                
            print(f"Connected. Server version: {browser.version()}")
            
            print("\n1. Navigating to Example.com...")
            page = browser.open("https://example.com")
            
            print(f"   URL: {page.url}")
            print(f"   Title: {page.title}")
            
            print("\n2. Inspecting page elements...")
            elements = page.inspect(visible_only=True)
            print(f"   Found {len(elements)} visible elements")
            
            for el in elements[:3]:
                print(f"   - <{el.get('type')}>: {el.get('text', '')!r}")
                
            print("\n3. Capturing screenshot...")
            page.screenshot(save_to="example_python_sdk.png")
            print("   Saved to example_python_sdk.png")
            
            print("\n4. Closing session...")
            
    except MrBrowserError as e:
        print(f"\nError: {e}")
        
    print("\nDone.")

if __name__ == "__main__":
    main()
