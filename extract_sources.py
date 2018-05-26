from sys import argv
from urllib import request
import json, re, io

if len(argv)<2:
    print("No argument provided. Exiting...")
    exit(1)

try:
    x = re.search('(.+)#(.+)', argv[1])
    data = json.loads(request.urlopen('https://api.github.com/repos/' + x.group(1) + '/issues/' + x.group(2)).read().decode('UTF-8'))["body"]
except:
    try:
        data = request.urlopen(argv[1]).read()
    except:
        try:
            data = open(argv[1],'r').read()
        except:
            print("Data not found. The argument should be a raw issue reference or a reachable filename.")
            exit(1)

# Find files in code blocks
d = []
for x in re.finditer('#>>', data, flags=re.DOTALL):
    d += [x.start()]

for x in range(len(d)-1):
    buf = io.StringIO(data[d[x]:d[x+1]-1])
    l = buf.readline()
    open(l[4:len(l)-2], 'w').write(buf.read())

# Find attached files
for x in re.finditer('\[(.+)\]\((https://github.com/VUnit/vunit/files/.+)\)', data, flags=re.DOTALL):
    l = x.group(1)
    open(l[0:len(l)-4], 'w').write(request.urlopen(x.group(2)).read().decode('UTF-8'))
