// fibonacci.js
const fs = require('fs');

// The python script to be executed
const pythonScript = `
import sys
import json
import csv

def fib(n):
    if n <= 0: return 0
    elif n == 1: return 1
    a, b = 0, 1
    for _ in range(2, n + 1):
        a, b = b, a + b
    return b

def process_file(filepath):
    try:
        # Try JSON first
        with open(filepath, 'r') as f:
             data = json.load(f)
             start = data.get('start', 0)
             end = data.get('end', 0)
             return start, end
    except:
        pass
        
    try:
        # Try CSV next
        with open(filepath, 'r') as f:
             reader = csv.reader(f)
             headers = next(reader) # skip headers
             row = next(reader)
             return int(row[0]), int(row[1])
    except:
        pass
        
    return 0, 0

def main():
    if len(sys.argv) < 2:
        print("Usage: python script.py <input_file_path>")
        sys.exit(1)
        
    filepath = sys.argv[1]
    start, end = process_file(filepath)
    
    print(f"Calculating Fibonacci from {start} to {end}")
    
    results = []
    for i in range(start, end + 1):
        results.append(fib(i))
        
    output = {
         "start": start,
         "end": end,
         "fibonacci_sequence": results
    }
    
    print(json.dumps(output))

if __name__ == "__main__":
    main()
`;

// Write the python script to a file
fs.writeFileSync('script.py', pythonScript);

console.log("Python script generated successfully: script.py");
