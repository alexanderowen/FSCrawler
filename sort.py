

def main():
    f = open("out", "r")
    l = []
    for line in f:
        l.append(line.strip())

    l.sort()
    f.close()
    f = open("out", "w")
    for line in l:
        f.write(line + "\n")



if __name__ == "__main__":
    main()
