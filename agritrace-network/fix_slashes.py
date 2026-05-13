import os
import glob

def fix_config_files():
    """Fix Windows line endings (\\r\\n) to Unix (\\n) in all config.yaml files"""
    count = 0
    for filepath in glob.iglob('crypto-config/**/config.yaml', recursive=True):
        with open(filepath, 'rb') as f:
            content = f.read()
        
        # Fix backslashes and Windows line endings
        new_content = content.replace(b'\\', b'/').replace(b'\r\n', b'\n')
        
        if new_content != content:
            with open(filepath, 'wb') as f:
                f.write(new_content)
            count += 1
            print(f"Fixed: {filepath}")
    
    # Also fix org-level config.yaml
    for filepath in glob.iglob('crypto-config/**/msp/config.yaml', recursive=True):
        with open(filepath, 'rb') as f:
            content = f.read()
        
        new_content = content.replace(b'\\', b'/').replace(b'\r\n', b'\n')
        
        if new_content != content:
            with open(filepath, 'wb') as f:
                f.write(new_content)
            count += 1
            print(f"Fixed: {filepath}")
    
    print(f"\nTotal files fixed: {count}")

if __name__ == "__main__":
    fix_config_files()
