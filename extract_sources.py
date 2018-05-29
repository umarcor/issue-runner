from sys import argv
from urllib.request import urlopen
import json, re, io

if len(argv)<2:
    print("No argument provided. Exiting...")
    exit(1)

#---

try:
    x = re.search('(.+)#(.+)', argv[1])
    data = json.loads(urlopen('https://api.github.com/repos/' + x.group(1) + '/issues/' + x.group(2)).read().decode('UTF-8'))["body"]
except:
    try:
        data = urlopen(argv[1]).read()
    except:
        try:
            data = open('../'+argv[1],'r').read()
        except:
            print("Data not found. The argument should be a raw issue reference or a reachable filename. Check your connection for network issues.")
            exit(1)

#---

# Find files in code blocks
d = []
for x in re.finditer('#>>', data, flags=re.DOTALL):
    d += [x.start()]

for x in range(len(d)-1):
    buf = io.StringIO(re.sub('\r', r'', data[d[x]:d[x+1]-1]))
    filename = buf.readline()[3:-1]
    if filename[0]==' ':
        filename=filename[1:]
    print("Get", filename)
    open(filename, 'w').write(buf.read())

#---

def x_invalid(filename, url):
    print("Invalid extension: [", filename, "](", url, ")")

def x_get(filename, url):
    print("Get [" + filename + "](" + url + ")")
    open(filename, 'wb').write(urlopen(url).read())

def x_txt(filename, url):
    x_get(filename[0:-4], url)

def x_tar(filename, url, mode):
    x_get(filename, url)
    print("Extract", filename)
    import tarfile
    t = tarfile.TarFile(filename, mode)
    t.extractall('.')
    t.close()

def x_tgz(filename, url):
    x_tar(filename, url, 'r:gz')

def x_tbz(filename, url):
    x_tar(filename, url, 'r:bz2')

def x_txz(filename, url):
    x_tar(filename, url, 'r:xz')

def x_zip(filename, url):
    x_get(filename, url)
    print("Extract", filename)
    import zipfile
    z = zipfile.ZipFile(filename, 'r')
    z.extractall('.')
    z.close()

switcher = {
    "txt": x_txt,
    "zip": x_zip,
    ".gz": x_tgz,
    "tgz": x_tgz,
    ".bz": x_tbz,
    "tbz": x_tbz,
    ".xz": x_txz,
    "txz": x_txz,
}

# Find attached files
for x in re.finditer('\[([^\s].+)\]\((https://github.com/.+/files/.+)\)', data, flags=re.DOTALL):
    switcher.get(x.group(1)[-3:], x_invalid)(x.group(1), x.group(2))
