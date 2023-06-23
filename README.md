# MAXIBAN

## TL:DR
Are you a skeptic concerned about Bitcoin nodes using your local network, ISP, or proxy service? MAXIBAN   
Are you a crypto widow or widower sick of paying for a static IP every month, or annoyed about lost savings? MAXIBAN   
Is your flatmate always telling you that Bitcoin isn't Crypto? MAXIBAN   
Are you concerned about the environmental impacts of Proof of Work Consensus protocols? MAXIBAN   
Do you know what ponzi means? MAXIBAN   
Consume MaxiBan once daily until symptoms subside. Side effects may include smug.

## What is MaxiBan?
MaxiBan is a security tool intended to discourage Bitcoin nodes operating from the users public IP address.   
This tool will connect and handshake to every listening IPv4 Bitcoin node and send a message that will get your public IP "discouraged" by the node for 24 hours.   
As per the Bitcoin Core documentation, discouragement is an anti-DOS feature that does several things:
*   As a non-listening node: Prefered for eviction when a node receives other connections.
*   As a listening node: No incoming connections. Your address is not shared by other nodes.

MaxiBan takes several minutes to run, and must be run again in 24 hours.

## Requirements

[Go](http://golang.org) 1.17 or newer.

## Installation

https://github.com/btcsuite/btcd/releases

#### Linux/BSD/MacOSX/POSIX - Build from Source

- Install Go according to the installation instructions here:
  http://golang.org/doc/install

- Ensure Go was installed properly and is a supported version:

```bash
$ go version
$ go env GOROOT GOPATH
```

- Clone this repository and navigate into directory
```bash
$ git clone https://github.com/doublethink/maxiban.git
$ cd maxiban
```

- Install dependencies, build, and run
```bash
$ go install
$ go build -o maxiban main
$ ./maxiban
```

## FAQ

### What is a Bitcoin node?
A Bitcoin node is a computer that stores a copy of the blockchain, and share addresses, transactions, and blocks with other nodes.  
The network of nodes is what makes Bitcoin work.

### What is a listening node?
A listening node is a Bitcoin node that accepts incoming connections in addition to outgoing connections to other nodes.   
Without them, the Bitcoin network would not be possible as a static IP and port are needed.   

### What are the implications for non-listening nodes?
You can still connect to listening nodes, but will get priority kicked from busy nodes.

### What are the implications for listening nodes?   
In addition to the above. Listening nodes will not connect to you, and will not share your address with other nodes.   

### What if I am using a proxy or TOR?   
The final public IP of MaxiBans route is the IP that is discouraged.   
If you are running from home with no static IP, the CG-NAT egress IP will be discouraged.   
If the OS running MaxiBan is using a SOCK5 proxy, the remote IP will be discouraged.   
If the OS running MaxiBan is using TOR, the IP of the exit relay will be discouraged.    
MaxiBan does not currently support onion addresses but its worth mentioning only other TOR nodes can connect to onion addresses...  
MaxiBan has not been tested through a proxy or TOR.

### Is this an exploit?
Not really. Its a clever use of game mechanics.   
Discourage is an anti-DOS feature implemented by the godking Satoshi himself in 2011.   
As it was implemented at the time, this would have been devastating for non-listening nodes as it prohibited discouraged IPs from connecting to listening nodes.

### Did you responsibly disclose your findings?
Yes, I emailed Satoshi but he didn't reply   

### Should I run this tool while connected to the network of a major Bitcoin mining operation or BaaS?
Yes, MaxiBan helps stabalise the grid and encourages investment in renewable energy...   

### Would it take down the network if an entity capable of spoofing TCP connections recursively ran this against all listening node IPs?
Yes.   
but its just code...

### Why doesn't MaxiBan support IPv6?
What is IPv6?   
Just kidding. It can, I just ran out of time.   

### If Bitcoin network traffic was encrypted would this still be possible?
Yes it would.   
But encryption would make Bitcoin censorship resistant, unlike today.

### Is it ethical to potentially limit others access to the network using the same public IP as you?
Is is ethical to consume as much energy as a medium sized country on a speculative asset?