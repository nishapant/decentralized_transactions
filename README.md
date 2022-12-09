# Decentralized Transactions

This is a distributed system of nodes that maintain transactions and account information about every other node in the network. To ensure that all nodes maintain the same order of transactions, the nodes use *totally ordered multicast*. To implement this, the [ISIS](https://courses.grainger.illinois.edu/ece428/sp2022//assets/slides/lect8-after.pdf#page=28) (Intermediate System to Intermediate System) algorithm can be used.

- The sender multicasts a message to all nodes on the network
- The receiving nodes will send a *proposed priority* back to the sender. This is a sequence number for the number of messages that have been completely delivered. This number is always the highest out of all other self-proposed priorities and agreed priorities of previous messages.
- The sender then chooses the maximum value out of all the proposed priorities as the agreed priority and multicasts this value back out to all other nodes on the network.
- The receiving nodes will update the message with the finalized priority, mark it as delivered, and reorder the message on the queue according to priority. The messages at the front of the queue that can be delivered and have not been yet will be delivered.

This project sends all messages with TCP and is written in Go. To run this, you will need multiple VM's or computers with Go installed. `config.txt` should be updated such that it contains a complete list of all the nodes that are running this program. It should be set up as shown below.

```
3    # number of nodes in the network
node1 sp21-cs425-g01-01.cs.illinois.edu 1234  # {node_name} {ip_address} {port_number}
node2 sp21-cs425-g01-02.cs.illinois.edu 1234
node3 sp21-cs425-g01-03.cs.illinois.edu 1234
```
