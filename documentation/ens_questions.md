1. What is the label hash of this domain?
    Q. Does this mean for a given namehash of "a.b.c" find keccak256(a), keccak256(b), and keccak256(c)? Do we know the parent domain and/or owner address?
    1. Watch NewOwner(bytes32 indexed node, bytes32 indexed label, address owner) events of the ENS Registry contract and filter for the root domain and/or owner address to narrow search
    `./vulcanize lightOmniWacther --config=./environments/<config.toml> --starting-block-number=## --ending-block-numer=### --contract-address=0x314159265dD8dbb310642f98f50C066173C1259b --contract-events=NewOwner`
    2. For each node + label pair we find emitted, calculate the keccak256(abi.encodePacked(node, label)) and see if it matches our namehash
    3. If it does, hash(label) is our answer

2. What is the parent domain of this domain?
    1. Watch NewOwner(bytes32 indexed node, bytes32 indexed label, address owner) events of the ENS Registry contract
    `./vulcanize lightOmniWacther --config=./environments/<config.toml> --starting-block-number=## --ending-block-numer=### --contract-address=0x314159265dD8dbb310642f98f50C066173C1259b --contract-events=NewOwner`
    2. Filter for our label (domain) and collect the node (parent domain namehash) that was emitted with it
    3. Call the Registry's resolver(bytes32 node) method for the parent node to find the parent domain's Resolver
    4. Call its Resolver's name(bytes32 node) method for the parent node to find the parent domain's name

3. What are the subdomains of this domain?
    1. Watch NewOwner events of the ENS Registry contract
    `./vulcanize lightOmniWacther --config=./environments/<config.toml> --starting-block-number=## --ending-block-numer=### --contract-address=0x314159265dD8dbb310642f98f50C066173C1259b --contract-events=NewOwner`
    2. Filter for our node (domain) and collect all the labels emitted with it
    3. Calculate subdomain hashes: subnode = keccak256(abi.encodePacked(node, label));
    4. Call the Registry's resolver(bytes32 node) method for a subnode to find the subdomain's Resolver
    5. Call its Resolver's name(bytes32 node) method for a subnode to find the subdomain's name

4. What domains does this address own?
    1. Watch NewOwner(bytes32 indexed node, bytes32 indexed label, address owner) and Transfer(bytes32 indexed node, address owner) events of the ENS Registry contract and filter for the address
    `./vulcanize lightOmniWacther --config=./environments/<config.toml> --starting-block-number=## --ending-block-numer=### --contract-address=0x314159265dD8dbb310642f98f50C066173C1259b --contract-eevents=NewOwner --contract-events=Transfer --event-filter-addresses=<address>`
    2. Generate list of all nodes this address has ever owned
    3. Check which of these they still own at a given blockheight by iterating over the list and calling the owner(bytes32 node) method
    4. Call the Registry's resolver(bytes32 node) method for a node to find the domain's Resolver
    5. Call its Resolver's name(bytes32 node) method for the node to find the domain's name

5. What names point to this address?
    Q. Is this in terms of which ENS nodes point to a given Resolver address? E.g. All nodes where the ENS records[node].resolver == address? Or is this in terms of Resolver records? E.g. All the records[node].names where the Resolver records[node].addr == address
    1. In the former case, watch NewResolver(bytes32 indexed node, address resolver) events of the Registry and filter for the account address
    `./vulcanize lightOmniWacther --config=./environments/<config.toml> --starting-block-number=## --ending-block-numer=### --contract-address=0x1da022710dF5002339274AaDEe8D58218e9D6AB5 --contract-events=NewResolver --event-filter-addresses=<address>`
    2. Generate a list of nodes that have pointed to this resolver address
    3. Check which of these names still point at the address by iterating over the list and calling the resolver(bytes32 node) method
    1. In the latter case, watch AddrChanged(bytes32 indexed node, address a) events of the Resolver and filter for the account address
    `./vulcanize lightOmniWacther --config=./environments/<config.toml> --starting-block-number=## --ending-block-numer=### --contract-address=0x1da022710dF5002339274AaDEe8D58218e9D6AB5 --contract-events=AddrChanged --event-filter-addresses=<address>`
    2. Generate our list of nodes that have pointed towards our address
    3. Check which of these they still own at a given blockheight by iterating over the list and calling the Resolver's addr(bytes32 node) method
    4. We can then fetch the string names of these nodes using the Resolver's name(bytes32 node) method.


Currently the only filtering that can be done during event watching is for addresses and the only methods
that can be polled in an automated fashion are ones that take only address-type arguments (of which there
are less than three) and return a single value. For the sake of answering these questions it would be really helpful if
we could also perform []byte filtering on the events and automate polling of events that take []byte-type arguments. I am
currently working on adding this in, and once it is in you would be able to automate more of the steps in these processes.

E.g. you will be able to run
`./vulcanize lightOmniWacther --config=./environments/<config.toml> --starting-block-number=## --ending-block-numer=### --contract-address=0x314159265dD8dbb310642f98f50C066173C1259b --contract-events=NewOwner --contract-events=Transfer --event-args=<address> --contract-methods=owner`
To automate the process in question 4 through step 3 (it will collect node []byte values emitted from the events it watches and then use those to call the owner method, persisting the results)

Or
`./vulcanize lightOmniWacther --config=./environments/<config.toml> --starting-block-number=## --ending-block-numer=### --contract-address=0x314159265dD8dbb310642f98f50C066173C1259b --contract-events=NewOwner --event-args=<bytes-to-filter-for>`
To provide automated filtering for node []byte values in question 3.